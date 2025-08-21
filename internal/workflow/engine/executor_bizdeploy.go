package engine

import (
	"fmt"
	"log/slog"
	"maps"
	"strings"

	"github.com/certimate-go/certimate/internal/deployer"
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
)

/**
 * Result Variables:
 *   - node.skipped: boolean
 */
type bizDeployNodeExecutor struct {
	nodeExecutor

	certificateRepo certificateRepository
	wfoutputRepo    workflowOutputRepository
}

func (ne *bizDeployNodeExecutor) Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
	execRes := newNodeExecutionResult(execCtx.Node)

	nodeCfg := execCtx.Node.Data.Config.AsBizDeploy()
	ne.logger.Info("ready to deploy certificate ...", slog.Any("config", nodeCfg))

	// 查询上次执行结果
	lastOutput, err := ne.getLastOutputArtifacts(execCtx)
	if err != nil {
		return execRes, err
	}

	// 获取前序节点输出证书
	var inputCertificate *domain.Certificate
	if inputState, ok := execCtx.inputs.Get(nodeCfg.CertificateOutputNodeId, "certificate"); ok {
		if inputStateValue, ok := inputState.Value.(string); ok {
			s := strings.Split(inputStateValue, "#")
			if len(s) == 2 {
				certificate, err := ne.certificateRepo.GetById(execCtx.ctx, s[1])
				if err != nil {
					ne.logger.Warn("failed to get input certificate")
					return execRes, err
				}

				inputCertificate = certificate
			}
		}
	}
	if inputCertificate == nil {
		return execRes, fmt.Errorf("invalid input certificate")
	}

	// 检测是否可以跳过本次执行
	if lastOutput != nil && inputCertificate.CreatedAt.Before(lastOutput.UpdatedAt) {
		if skippable, reason := ne.checkCanSkip(execCtx, lastOutput); skippable {
			ne.logger.Info(fmt.Sprintf("skip this deployment, because %s", reason))

			execRes.AddVariableWithScope(execCtx.Node.Id, stateVarKeyNodeSkipped, true, "boolean")
			return execRes, nil
		} else if reason != "" {
			ne.logger.Info(fmt.Sprintf("re-deploy, because %s", reason))
		}
	}

	// 初始化部署器
	// TODO: 解耦
	deployer, err := deployer.NewWithWorkflowNode(deployer.DeployerWithWorkflowNodeConfig{
		Node:           execCtx.Node,
		Logger:         ne.logger,
		CertificatePEM: inputCertificate.Certificate,
		PrivateKeyPEM:  inputCertificate.PrivateKey,
	})
	if err != nil {
		ne.logger.Warn("failed to create deployer provider")
		return execRes, err
	}

	// 部署证书
	if err := deployer.Deploy(execCtx.ctx); err != nil {
		ne.logger.Warn("failed to deploy certificate")
		return execRes, err
	}

	// 节点输出
	execRes.outputForced = true
	execRes.AddVariableWithScope(execCtx.Node.Id, stateVarKeyNodeSkipped, false, "boolean")

	ne.logger.Info("deployment completed")
	return execRes, nil
}

func (ne *bizDeployNodeExecutor) getLastOutputArtifacts(execCtx *NodeExecutionContext) (*domain.WorkflowOutput, error) {
	lastOutput, err := ne.wfoutputRepo.GetByNodeId(execCtx.ctx, execCtx.Node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		return nil, fmt.Errorf("failed to get last output record of node #%s: %w", execCtx.Node.Id, err)
	}

	return lastOutput, nil
}

func (ne *bizDeployNodeExecutor) checkCanSkip(execCtx *NodeExecutionContext, lastOutput *domain.WorkflowOutput) (_skip bool, _reason string) {
	thisNodeCfg := execCtx.Node.Data.Config.AsBizDeploy()

	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次部署时的关键配置（即影响证书部署的）参数是否一致
		lastNodeCfg := lastOutput.NodeConfig.AsBizDeploy()

		if thisNodeCfg.ProviderAccessId != lastNodeCfg.ProviderAccessId {
			return false, "the configuration item 'ProviderAccessId' changed"
		}
		if !maps.Equal(thisNodeCfg.ProviderConfig, lastNodeCfg.ProviderConfig) {
			return false, "the configuration item 'ProviderConfig' changed"
		}

		if thisNodeCfg.SkipOnLastSucceeded {
			return true, "the last deployment already completed"
		}
	}

	return false, ""
}

func newBizDeployNodeExecutor() NodeExecutor {
	return &bizDeployNodeExecutor{
		nodeExecutor:    nodeExecutor{logger: slog.Default()},
		certificateRepo: repository.NewCertificateRepository(),
		wfoutputRepo:    repository.NewWorkflowOutputRepository(),
	}
}
