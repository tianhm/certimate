package huaweicloudobs

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/certimate-go/certimate/pkg/core"
)

type SSLDeployerProviderConfig struct {
	// 华为云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 华为云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 华为云 Bucket 对应的 Endpoint。
	Endpoint string `json:"endpoint"`
	// 华为云 OBS 桶名。
	Bucket string `json:"bucket"`
	// 自定义域名。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config *SSLDeployerProviderConfig
	logger *slog.Logger
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	config.Endpoint = strings.TrimPrefix(strings.TrimPrefix(config.Endpoint, "http://"), "https://")

	return &SSLDeployerProvider{
		config: config,
		logger: slog.Default(),
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}
}

// REF: https://support.huaweicloud.com/usermanual-obs/obs_06_3200.html
// REF: https://support.huaweicloud.com/api-obs/obs_04_0059.html
func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	if d.config.Domain == "" {
		return nil, fmt.Errorf("config `domain` is required")
	}

	url := fmt.Sprintf("https://%s.%s/?customdomain=%s", d.config.Bucket, d.config.Endpoint, d.config.Domain)
	bodyXML := fmt.Sprintf(`
<CustomDomainConfiguration>
	<Name>%s</Name>
	<Certificate>%s</Certificate>
	<CertificateChain>%s</CertificateChain>
	<PrivateKey>%s</PrivateKey>
</CustomDomainConfiguration>`,
		d.config.Bucket+"_"+d.config.Domain, certPEM, certPEM, privkeyPEM,
	)

	// 计算 Content-MD5（Base64 编码）
	md5sum := md5.Sum([]byte(bodyXML))
	contentMD5 := base64.StdEncoding.EncodeToString(md5sum[:])

	// 日期
	date := time.Now().UTC().Format(http.TimeFormat)

	// 构造签名字符串
	method := "PUT"
	contentType := "application/xml"
	canonicalizedResource := fmt.Sprintf("/%s/?customdomain=%s", d.config.Bucket, d.config.Domain)
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s\n%s", method, contentMD5, contentType, date, canonicalizedResource)

	// HMAC-SHA1 签名
	h := hmac.New(sha1.New, []byte(d.config.SecretAccessKey))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Authorization
	authHeader := fmt.Sprintf("OBS %s:%s", d.config.AccessKeyId, signature)

	// 创建请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(bodyXML)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Date", date)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-MD5", contentMD5)
	req.Header.Set("Content-Type", contentType)

	// 请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body := new(bytes.Buffer)
		body.ReadFrom(resp.Body)
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, body.String())
	}
	return &core.SSLDeployResult{}, nil
}
