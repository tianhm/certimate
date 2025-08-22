package engine

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
)

/**
 * Result Variables:
 *   - node.skipped: boolean
 *   - certificate.validity: boolean
 *   - certificate.daysLeft: number
 *
 * Result Outputs:
 *   - ref: certificate: string
 */
type bizUploadNodeExecutor struct {
	nodeExecutor

	certificateRepo certificateRepository
	wfoutputRepo    workflowOutputRepository
}

func (ne *bizUploadNodeExecutor) Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
	execRes := newNodeExecutionResult(execCtx.Node)

	nodeCfg := execCtx.Node.Data.Config.AsBizUpload()
	ne.logger.Info("ready to upload certiticate ...", slog.Any("config", nodeCfg))

	// 查询上次执行结果
	lastOutput, lastCertificate, err := ne.getLastOutputArtifacts(execCtx)
	if err != nil {
		return execRes, err
	} else if lastCertificate != nil {
		execRes.AddOutput(stateIOTypeRef, "certificate", fmt.Sprintf("%s#%s", domain.CollectionNameCertificate, lastCertificate.Id), "string")
		execRes.AddVariableWithScope(execCtx.Node.Id, stateVarKeyCertificateValidity, time.Now().After(lastCertificate.ValidityNotAfter), "boolean")
		execRes.AddVariableWithScope(execCtx.Node.Id, stateVarKeyCertificateDaysLeft, int32(time.Until(lastCertificate.ValidityNotAfter).Hours()/24), "number")
	}

	// 检测是否可以跳过本次执行
	if skippable, reason := ne.checkCanSkip(execCtx, lastOutput, lastCertificate); skippable {
		ne.logger.Info(fmt.Sprintf("skip this uploading, because %s", reason))

		execRes.AddVariableWithScope(execCtx.Node.Id, stateVarKeyNodeSkipped, true, "boolean")
		return execRes, nil
	} else if reason != "" {
		ne.logger.Info(fmt.Sprintf("re-upload, because %s", reason))
	} else {
		ne.logger.Info("no found last uploaded certificate, begin to upload")
	}

	// 保存证书实体
	certificate := &domain.Certificate{
		Source:         domain.CertificateSourceTypeUpload,
		WorkflowId:     execCtx.WorkflowId,
		WorkflowRunId:  execCtx.RunId,
		WorkflowNodeId: execCtx.Node.Id,
	}
	certificate.PopulateFromPEM(nodeCfg.Certificate, nodeCfg.PrivateKey)
	if certificate, err := ne.certificateRepo.Save(execCtx.ctx, certificate); err != nil {
		ne.logger.Warn("failed to save certificate")
		return execRes, err
	} else {
		ne.logger.Info("certificate saved", slog.String("recordId", certificate.Id))
	}

	// 节点输出
	execRes.AddOutputWithPersistent(stateIOTypeRef, "certificate", fmt.Sprintf("%s#%s", domain.CollectionNameCertificate, certificate.Id), "string")
	execRes.AddVariableWithScope(execCtx.Node.Id, stateVarKeyNodeSkipped, false, "boolean")
	execRes.AddVariableWithScope(execCtx.Node.Id, stateVarKeyCertificateValidity, true, "boolean")
	execRes.AddVariableWithScope(execCtx.Node.Id, stateVarKeyCertificateDaysLeft, int32(time.Until(certificate.ValidityNotAfter).Hours()/24), "number")

	ne.logger.Info("uploading completed")
	return execRes, nil
}

func (ne *bizUploadNodeExecutor) getLastOutputArtifacts(execCtx *NodeExecutionContext) (*domain.WorkflowOutput, *domain.Certificate, error) {
	lastOutput, err := ne.wfoutputRepo.GetByNodeId(execCtx.ctx, execCtx.Node.Id)
	if err != nil && !domain.IsRecordNotFoundError(err) {
		return nil, nil, fmt.Errorf("failed to get last output record of node #%s: %w", execCtx.Node.Id, err)
	}

	if lastOutput != nil {
		lastCertificate, err := ne.certificateRepo.GetByWorkflowRunIdAndNodeId(execCtx.ctx, lastOutput.RunId, lastOutput.NodeId)
		if err != nil && !domain.IsRecordNotFoundError(err) {
			return lastOutput, nil, fmt.Errorf("failed to get last certificate record of node #%s: %w", execCtx.Node.Id, err)
		}

		return lastOutput, lastCertificate, nil
	}

	return lastOutput, nil, nil
}

func (ne *bizUploadNodeExecutor) checkCanSkip(execCtx *NodeExecutionContext, lastOutput *domain.WorkflowOutput, lastCertificate *domain.Certificate) (_skip bool, _reason string) {
	thisNodeCfg := execCtx.Node.Data.Config.AsBizUpload()

	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次上传时的关键配置（即影响证书上传的）参数是否一致
		lastNodeCfg := lastOutput.NodeConfig.AsBizUpload()

		if strings.TrimSpace(thisNodeCfg.Certificate) != strings.TrimSpace(lastNodeCfg.Certificate) {
			return false, "the configuration item 'Certificate' changed"
		}
		if strings.TrimSpace(thisNodeCfg.PrivateKey) != strings.TrimSpace(lastNodeCfg.PrivateKey) {
			return false, "the configuration item 'PrivateKey' changed"
		}
	}

	if lastCertificate != nil {
		return true, "the last uploaded certificate already exists"
	}

	return false, ""
}

func newBizUploadNodeExecutor() NodeExecutor {
	return &bizUploadNodeExecutor{
		nodeExecutor:    nodeExecutor{logger: slog.Default()},
		certificateRepo: repository.NewCertificateRepository(),
		wfoutputRepo:    repository.NewWorkflowOutputRepository(),
	}
}
