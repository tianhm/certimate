package engine

import (
	"fmt"
	"log/slog"
	"maps"

	"github.com/certimate-go/certimate/internal/deployer"
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
)

type bizDeployNodeExecutor struct {
	nodeExecutor

	certificateRepo certificateRepository
	wfrunRepo       workflowRunRepository
	wfoutputRepo    workflowOutputRepository
}

func (ne *bizDeployNodeExecutor) Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
	execRes := &NodeExecutionResult{}

	nodeCfg := execCtx.Node.GetConfigForBizDeploy()
	ne.logger.Info("ready to deploy certificate ...", slog.Any("config", nodeCfg))

	// 查询上次执行结果
	lastOutput, err := ne.getLastOutputArtifacts(execCtx)
	if err != nil {
		return execRes, err
	}

	// 获取前序节点输出证书
	// TODO: 利用输入而非查库
	certificate, err := ne.certificateRepo.GetByWorkflowIdAndNodeId(execCtx.ctx, execCtx.WorkflowId, nodeCfg.CertificateOutputNodeId)
	if err != nil {
		ne.logger.Warn("invalid certificate output")
		return execRes, err
	}

	// 检测是否可以跳过本次执行
	if lastOutput != nil && certificate.CreatedAt.Before(lastOutput.UpdatedAt) {
		if skippable, reason := ne.checkCanSkip(execCtx, lastOutput); skippable {
			ne.logger.Info(fmt.Sprintf("skip this deployment, because %s", reason))

			execRes.AddVariable(execCtx.Node.Id, wfVariableKeyNodeSkipped, true, "boolean")
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
		CertificatePEM: certificate.Certificate,
		PrivateKeyPEM:  certificate.PrivateKey,
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

	// 保存执行结果
	output := &domain.WorkflowOutput{
		WorkflowId: execCtx.WorkflowId,
		RunId:      execCtx.RunId,
		NodeId:     execCtx.Node.Id,
		Succeeded:  true,
	}
	if _, err := ne.wfoutputRepo.Save(execCtx.ctx, output); err != nil {
		ne.logger.Warn("failed to save node output")
		return execRes, err
	}

	// 记录中间结果
	execRes.AddVariable(execCtx.Node.Id, wfVariableKeyNodeSkipped, false, "boolean")

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
	thisNodeCfg := execCtx.Node.GetConfigForBizDeploy()

	if lastOutput != nil && lastOutput.Succeeded {
		lastRun, err := ne.wfrunRepo.GetById(execCtx.ctx, lastOutput.RunId)
		if err != nil {
			return true, "failed to get last run"
		}
		lastNode, lastNodeExists := lastRun.Graph.GetNodeById(lastOutput.NodeId)
		if !lastNodeExists {
			return true, "failed to get last run node"
		}

		// 比较和上次部署时的关键配置（即影响证书部署的）参数是否一致
		lastNodeCfg := lastNode.GetConfigForBizDeploy()

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
		wfrunRepo:       repository.NewWorkflowRunRepository(),
		wfoutputRepo:    repository.NewWorkflowOutputRepository(),
	}
}
