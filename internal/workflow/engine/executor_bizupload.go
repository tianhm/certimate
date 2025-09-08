package engine

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-resty/resty/v2"

	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
	xcert "github.com/certimate-go/certimate/pkg/utils/cert"
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

const (
	BizUploadSourceForm  = "form"
	BizUploadSourceLocal = "local"
	BizUploadSourceURL   = "url"
)

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
	} else if lastCertificate != nil {
		ne.logger.Info("no found last uploaded certificate, begin to upload")
	} else {
		ne.logger.Info("try to upload")
	}

	// 获取证书及私钥
	var certPEM, privkeyPEM string
	switch nodeCfg.Source {
	case BizUploadSourceForm:
		{
			certPEM = nodeCfg.Certificate
			privkeyPEM = nodeCfg.PrivateKey
		}

	case BizUploadSourceLocal:
		{
			certData, err := os.ReadFile(nodeCfg.Certificate)
			if err != nil {
				return execRes, fmt.Errorf("failed to read certificate file from local path: %w", err)
			} else {
				certPEM = string(certData)
			}

			privkeyData, err := os.ReadFile(nodeCfg.PrivateKey)
			if err != nil {
				return execRes, fmt.Errorf("failed to read private key file from local path: %w", err)
			} else {
				privkeyPEM = string(privkeyData)
			}
		}

	case BizUploadSourceURL:
		{
			client := resty.New()
			client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

			certResp, err := client.NewRequest().Get(nodeCfg.Certificate)
			if err != nil || certResp.IsError() {
				return execRes, fmt.Errorf("failed to download certificate from URL: %w", err)
			} else {
				certPEM = string(certResp.Body())
			}

			privkeyResp, err := client.NewRequest().Get(nodeCfg.PrivateKey)
			if err != nil || privkeyResp.IsError() {
				return execRes, fmt.Errorf("failed to download private key from URL: %w", err)
			} else {
				privkeyPEM = string(privkeyResp.Body())
			}
		}

	default:
		return execRes, fmt.Errorf("unsupported upload source: '%s'", nodeCfg.Source)
	}

	// 验证证书
	certX509, err := certcrypto.ParsePEMCertificate([]byte(certPEM))
	if err != nil {
		return execRes, err
	} else if time.Now().After(certX509.NotAfter) {
		ne.logger.Warn(fmt.Sprintf("the uploaded certificate has expired at %s", certX509.NotAfter.UTC().Format(time.RFC3339)))
	}

	// 验证私钥
	privkey, err := certcrypto.ParsePEMPrivateKey([]byte(privkeyPEM))
	if err != nil {
		return nil, err
	} else {
		matched := false
		switch pub := certX509.PublicKey.(type) {
		case *rsa.PublicKey:
			p, ok := privkey.(*rsa.PrivateKey)
			matched = ok && pub.Equal(p.Public())
		case *ecdsa.PublicKey:
			p, ok := privkey.(*ecdsa.PrivateKey)
			matched = ok && pub.Equal(p.Public())
		case ed25519.PublicKey:
			p, ok := privkey.(ed25519.PrivateKey)
			matched = ok && pub.Equal(p.Public())
		default:
			matched = false
		}

		if !matched {
			return nil, fmt.Errorf("the uploaded private key does not match the uploaded certificate")
		}
	}

	// 二次检测是否可以跳过执行
	if lastCertificate != nil {
		lastCertX509, err := xcert.ParseCertificateFromPEM(lastCertificate.Certificate)
		if err == nil && xcert.EqualCertificates(certX509, lastCertX509) {
			ne.logger.Info("skip this uploading, because the last uploaded certificate already exists")
			return execRes, nil
		}
	}

	// 保存证书实体
	certificate := &domain.Certificate{
		Source:         domain.CertificateSourceTypeUpload,
		WorkflowId:     execCtx.WorkflowId,
		WorkflowRunId:  execCtx.RunId,
		WorkflowNodeId: execCtx.Node.Id,
	}
	certificate.PopulateFromPEM(certPEM, privkeyPEM)
	if certificate, err := ne.certificateRepo.Save(execCtx.ctx, certificate); err != nil {
		ne.logger.Warn("could not save certificate")
		return execRes, err
	} else {
		ne.logger.Info("certificate saved", slog.String("recordId", certificate.Id))
	}

	// 节点输出
	execRes.AddOutputWithPersistent(stateIOTypeRef, "certificate", fmt.Sprintf("%s#%s", domain.CollectionNameCertificate, certificate.Id), "string")
	execRes.AddVariableWithScope(execCtx.Node.Id, stateVarKeyNodeSkipped, false, "boolean")
	execRes.AddVariableWithScope(execCtx.Node.Id, stateVarKeyCertificateValidity, time.Now().After(certificate.ValidityNotAfter), "boolean")
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

		if thisNodeCfg.Source != lastNodeCfg.Source {
			return false, "the configuration item 'Source' changed"
		}

		switch thisNodeCfg.Source {
		case BizUploadSourceForm:
			if strings.TrimSpace(thisNodeCfg.Certificate) != strings.TrimSpace(lastNodeCfg.Certificate) {
				return false, "the configuration item 'Certificate' changed"
			}
			if strings.TrimSpace(thisNodeCfg.PrivateKey) != strings.TrimSpace(lastNodeCfg.PrivateKey) {
				return false, "the configuration item 'PrivateKey' changed"
			}

		default:
			// 本地或远程文件来源，需实际下载后才能比较
			return false, ""
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
