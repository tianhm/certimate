package engine

import (
	"crypto/x509"
	"fmt"
	"log/slog"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/certimate-go/certimate/internal/repository"
	xhttp "github.com/certimate-go/certimate/pkg/utils/http"
	xtls "github.com/certimate-go/certimate/pkg/utils/tls"
)

type bizMonitorNodeExecutor struct {
	nodeExecutor

	certificateRepo certificateRepository
}

func (ne *bizMonitorNodeExecutor) Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
	execRes := &NodeExecutionResult{}

	nodeCfg := execCtx.Node.Data.Config.AsBizMonitor()
	ne.logger.Info("ready to monitor certificate ...", slog.Any("config", nodeCfg))

	targetAddr := net.JoinHostPort(nodeCfg.Host, strconv.Itoa(int(nodeCfg.Port)))
	if nodeCfg.Port == 0 {
		targetAddr = net.JoinHostPort(nodeCfg.Host, "443")
	}

	targetDomain := nodeCfg.Domain
	if targetDomain == "" {
		targetDomain = nodeCfg.Host
	}

	ne.logger.Info(fmt.Sprintf("retrieving certificate at %s (domain: %s)", targetAddr, targetDomain))

	const MAX_ATTEMPTS = 3
	const RETRY_INTERVAL = 2 * time.Second
	var err error
	var certs []*x509.Certificate
	for attempt := 0; attempt < MAX_ATTEMPTS; attempt++ {
		if attempt > 0 {
			ne.logger.Info(fmt.Sprintf("retry %d time(s) ...", attempt, targetAddr))

			select {
			case <-execCtx.ctx.Done():
				return execRes, execCtx.ctx.Err()
			case <-time.After(RETRY_INTERVAL):
			}
		}

		certs, err = ne.tryRetrievePeerCertificates(execCtx, targetAddr, targetDomain, nodeCfg.RequestPath)
		if err == nil {
			break
		}
	}

	if err != nil {
		ne.logger.Warn("failed to monitor certificate")
		return execRes, err
	} else {
		if len(certs) == 0 {
			ne.logger.Warn("no ssl certificates retrieved in http response")

			execRes.AddVariable(execCtx.Node.Id, wfVariableKeyCertificateValidity, false, "boolean")
			execRes.AddVariable(execCtx.Node.Id, wfVariableKeyCertificateDaysLeft, 0, "number")
		} else {
			cert := certs[0] // 只取证书链中的第一个证书，即服务器证书
			ne.logger.Info(fmt.Sprintf("ssl certificate retrieved (serial='%s', subject='%s', issuer='%s', not_before='%s', not_after='%s', sans='%s')",
				cert.SerialNumber, cert.Subject.String(), cert.Issuer.String(),
				cert.NotBefore.Format(time.RFC3339), cert.NotAfter.Format(time.RFC3339),
				strings.Join(cert.DNSNames, ";")),
			)

			now := time.Now()
			isCertPeriodValid := now.Before(cert.NotAfter) && now.After(cert.NotBefore)
			isCertHostMatched := true
			if err := cert.VerifyHostname(targetDomain); err != nil {
				isCertHostMatched = false
			}

			validated := isCertPeriodValid && isCertHostMatched
			daysLeft := int(math.Floor(cert.NotAfter.Sub(now).Hours() / 24))
			execRes.AddVariable(execCtx.Node.Id, wfVariableKeyCertificateValidity, validated, "boolean")
			execRes.AddVariable(execCtx.Node.Id, wfVariableKeyCertificateDaysLeft, daysLeft, "number")

			if validated {
				ne.logger.Info(fmt.Sprintf("the certificate is valid, and will expire in %d day(s)", daysLeft))
			} else {
				if !isCertHostMatched {
					ne.logger.Warn("the certificate is invalid, because it is not matched the host")
				} else if !isCertPeriodValid {
					ne.logger.Warn("the certificate is invalid, because it is either expired or not yet valid")
				} else {
					ne.logger.Warn("the certificate is invalid")
				}
			}
		}
	}

	ne.logger.Info("monitoring completed")
	return execRes, nil
}

func (ne *bizMonitorNodeExecutor) tryRetrievePeerCertificates(execCtx *NodeExecutionContext, addr, domain, requestPath string) ([]*x509.Certificate, error) {
	transport := xhttp.NewDefaultTransport()
	transport.TLSClientConfig = xtls.NewInsecureConfig()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	url := fmt.Sprintf("https://%s/%s", addr, strings.TrimLeft(requestPath, "/"))
	req, err := http.NewRequestWithContext(execCtx.ctx, http.MethodHead, url, nil)
	if err != nil {
		err = fmt.Errorf("failed to create http request: %w", err)
		ne.logger.Warn(err.Error())
		return nil, err
	}

	req.Header.Set("Host", domain)
	req.Header.Set("User-Agent", "certimate")
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to send http request: %w", err)
		ne.logger.Warn(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.TLS == nil || len(resp.TLS.PeerCertificates) == 0 {
		return make([]*x509.Certificate, 0), nil
	}
	return resp.TLS.PeerCertificates, nil
}

func newBizMonitorNodeExecutor() NodeExecutor {
	return &bizMonitorNodeExecutor{
		nodeExecutor:    nodeExecutor{logger: slog.Default()},
		certificateRepo: repository.NewCertificateRepository(),
	}
}
