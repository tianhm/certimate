package engine

import (
	"fmt"
	"log/slog"
	"maps"
	"time"

	"github.com/certimate-go/certimate/internal/applicant"
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
)

type bizApplyNodeExecutor struct {
	nodeExecutor

	certificateRepo certificateRepository
	wfoutputRepo    workflowOutputRepository
}

func (ne *bizApplyNodeExecutor) Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
	execRes := &NodeExecutionResult{}

	nodeCfg := execCtx.Node.Data.Config.AsBizApply()
	ne.logger.Info("ready to request certificate ...", slog.Any("config", nodeCfg))

	// 查询上次执行结果
	lastOutput, lastCertificate, err := ne.getLastOutputArtifacts(execCtx)
	if err != nil {
		return execRes, err
	}

	// 检测是否可以跳过本次执行
	if skippable, reason := ne.checkCanSkip(execCtx, lastOutput, lastCertificate); skippable {
		ne.logger.Info(fmt.Sprintf("skip this application, because %s", reason))

		execRes.AddVariable(execCtx.Node.Id, stateVariableKeyNodeSkipped, true, "boolean")
		return execRes, nil
	} else if reason != "" {
		ne.logger.Info(fmt.Sprintf("re-apply, because %s", reason))
	} else {
		ne.logger.Info("no found last issued certificate, begin to apply")
	}

	// 初始化申请器
	// TODO: 解耦
	applicant, err := applicant.NewWithWorkflowNode(applicant.ApplicantWithWorkflowNodeConfig{
		WorkflowId: execCtx.WorkflowId,
		Node:       execCtx.Node,
		Logger:     ne.logger,
	})
	if err != nil {
		ne.logger.Warn("failed to create applicant provider")
		return execRes, err
	}

	// 申请证书
	applyResult, err := applicant.Apply(execCtx.ctx)
	if err != nil {
		ne.logger.Warn("failed to obtain certificate")
		return execRes, err
	}

	// 解析证书并生成实体
	certX509, err := xcert.ParseCertificateFromPEM(applyResult.FullChainCertificate)
	if err != nil {
		ne.logger.Warn("failed to parse certificate, may be the CA responded error")
		return execRes, err
	}

	certificate := &domain.Certificate{
		Source:            domain.CertificateSourceTypeRequest,
		Certificate:       applyResult.FullChainCertificate,
		PrivateKey:        applyResult.PrivateKey,
		IssuerCertificate: applyResult.IssuerCertificate,
		ACMEAccountUrl:    applyResult.ACMEAccountUrl,
		ACMECertUrl:       applyResult.ACMECertUrl,
		ACMECertStableUrl: applyResult.ACMECertStableUrl,
	}
	certificate.PopulateFromX509(certX509)

	// 保存执行结果
	// TODO: 解耦
	output := &domain.WorkflowOutput{
		WorkflowId: execCtx.WorkflowId,
		RunId:      execCtx.RunId,
		NodeId:     execCtx.Node.Id,
		NodeConfig: execCtx.Node.Data.Config,
		Succeeded:  true,
		Outputs: []*domain.WorkflowOutputEntry{
			{
				Name:      "certificate",
				Type:      "certificate",
				Value:     certificate.Id,
				ValueType: "string",
			},
		},
	}
	if _, err := ne.wfoutputRepo.SaveWithCertificate(execCtx.ctx, output, certificate); err != nil {
		ne.logger.Warn("failed to save node output")
		return execRes, err
	}

	// 保存 ARI 记录
	if lastCertificate != nil && applyResult.ARIReplaced {
		lastCertificate.ACMERenewed = true
		ne.certificateRepo.Save(execCtx.ctx, lastCertificate)
	}

	// 记录中间结果
	execRes.AddVariable(execCtx.Node.Id, stateVariableKeyNodeSkipped, false, "boolean")
	execRes.AddVariable(execCtx.Node.Id, stateVariableKeyCertificateValidity, true, "boolean")
	execRes.AddVariable(execCtx.Node.Id, stateVariableKeyCertificateDaysLeft, int32(time.Until(certificate.ValidityNotAfter).Hours()/24), "number")

	ne.logger.Info("application completed")
	return execRes, nil
}

func (ne *bizApplyNodeExecutor) getLastOutputArtifacts(execCtx *NodeExecutionContext) (*domain.WorkflowOutput, *domain.Certificate, error) {
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

func (ne *bizApplyNodeExecutor) checkCanSkip(execCtx *NodeExecutionContext, lastOutput *domain.WorkflowOutput, lastCertificate *domain.Certificate) (_skip bool, _reason string) {
	thisNodeCfg := execCtx.Node.Data.Config.AsBizApply()

	if lastOutput != nil && lastOutput.Succeeded {
		// 比较和上次申请时的关键配置（即影响证书签发的）参数是否一致
		lastNodeCfg := lastOutput.NodeConfig.AsBizApply()

		if thisNodeCfg.Domains != lastNodeCfg.Domains {
			return false, "the configuration item 'Domains' changed"
		}
		if thisNodeCfg.ContactEmail != lastNodeCfg.ContactEmail {
			return false, "the configuration item 'ContactEmail' changed"
		}
		if thisNodeCfg.Provider != lastNodeCfg.Provider {
			return false, "the configuration item 'Provider' changed"
		}
		if thisNodeCfg.ProviderAccessId != lastNodeCfg.ProviderAccessId {
			return false, "the configuration item 'ProviderAccessId' changed"
		}
		if !maps.Equal(thisNodeCfg.ProviderConfig, lastNodeCfg.ProviderConfig) {
			return false, "the configuration item 'ProviderConfig' changed"
		}
		if thisNodeCfg.CAProvider != lastNodeCfg.CAProvider {
			return false, "the configuration item 'CAProvider' changed"
		}
		if thisNodeCfg.CAProviderAccessId != lastNodeCfg.CAProviderAccessId {
			return false, "the configuration item 'CAProviderAccessId' changed"
		}
		if !maps.Equal(thisNodeCfg.CAProviderConfig, lastNodeCfg.CAProviderConfig) {
			return false, "the configuration item 'CAProviderConfig' changed"
		}
		if thisNodeCfg.KeyAlgorithm != lastNodeCfg.KeyAlgorithm {
			return false, "the configuration item 'KeyAlgorithm' changed"
		}
	}

	if lastCertificate != nil {
		renewalInterval := time.Duration(thisNodeCfg.SkipBeforeExpiryDays) * time.Hour * 24
		expirationTime := time.Until(lastCertificate.ValidityNotAfter)
		daysLeft := int(expirationTime.Hours() / 24)
		if expirationTime > renewalInterval {
			// TODO: 优化此处逻辑，[checkCanSkip] 方法不应该修改中间结果，违背单一职责
			// ne.outputs[variableKeyCertificateValidity] = strconv.FormatBool(true)
			// ne.outputs[variableKeyCertificateDaysLeft] = strconv.FormatInt(int64(daysLeft), 10)
			return true, fmt.Sprintf("the last issued certificate #%s expires in %d day(s), next renewal will be in %d day(s)", lastCertificate.Id, daysLeft, thisNodeCfg.SkipBeforeExpiryDays)
		}

		return false, fmt.Sprintf("the last issued certificate #%s expires in %d day(s)", lastCertificate.Id, daysLeft)
	}

	return false, ""
}

func newBizApplyNodeExecutor() NodeExecutor {
	return &bizApplyNodeExecutor{
		nodeExecutor:    nodeExecutor{logger: slog.Default()},
		certificateRepo: repository.NewCertificateRepository(),
		wfoutputRepo:    repository.NewWorkflowOutputRepository(),
	}
}
