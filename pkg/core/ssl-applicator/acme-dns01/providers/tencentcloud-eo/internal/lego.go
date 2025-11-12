package internal

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcteo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"golang.org/x/net/idna"
)

const (
	envNamespace = "TENCENTCLOUDEO_"

	EnvSecretID  = envNamespace + "SECRET_ID"
	EnvSecretKey = envNamespace + "SECRET_KEY"
	EnvZoneID    = envNamespace + "ZONE_ID"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvHTTPTimeout        = envNamespace + "HTTP_TIMEOUT"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

type Config struct {
	SecretID  string
	SecretKey string
	ZoneID    string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
	HTTPTimeout        time.Duration
}

type DNSProvider struct {
	client *TeoClient
	config *Config

	recordIDs   map[string]*string
	recordIDsMu sync.Mutex
}

func NewDefaultConfig() *Config {
	return &Config{
		TTL:                env.GetOrDefaultInt(EnvTTL, 300),
		PropagationTimeout: env.GetOrDefaultSecond(EnvPropagationTimeout, 2*time.Minute),
		PollingInterval:    env.GetOrDefaultSecond(EnvPollingInterval, dns01.DefaultPollingInterval),
		HTTPTimeout:        env.GetOrDefaultSecond(EnvHTTPTimeout, 30*time.Second),
	}
}

func NewDNSProvider() (*DNSProvider, error) {
	values, err := env.Get(EnvSecretID, EnvSecretKey, EnvZoneID)
	if err != nil {
		return nil, fmt.Errorf("tencentcloud-eo: %w", err)
	}

	config := NewDefaultConfig()
	config.SecretID = values[EnvSecretID]
	config.SecretKey = values[EnvSecretKey]
	config.ZoneID = values[EnvSecretKey]

	return NewDNSProviderConfig(config)
}

func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("tencentcloud-eo: the configuration of the DNS provider is nil")
	}

	credential := common.NewCredential(config.SecretID, config.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqTimeout = int(math.Round(config.HTTPTimeout.Seconds()))
	client, err := NewTeoClient(credential, "", cpf)
	if err != nil {
		return nil, err
	}

	return &DNSProvider{
		client:      client,
		config:      config,
		recordIDs:   make(map[string]*string),
		recordIDsMu: sync.Mutex{},
	}, nil
}

func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	punnyCoded, err := idna.ToASCII(dns01.UnFqdn(info.EffectiveFQDN))
	if err != nil {
		return fmt.Errorf("tencentcloud-eo: fail to convert punycode: %w", err)
	}

	// REF: https://cloud.tencent.com/document/product/1552/80720
	teoCreateDnsRecordReq := tcteo.NewCreateDnsRecordRequest()
	teoCreateDnsRecordReq.ZoneId = common.StringPtr(d.config.ZoneID)
	teoCreateDnsRecordReq.Name = common.StringPtr(punnyCoded)
	teoCreateDnsRecordReq.Type = common.StringPtr("TXT")
	teoCreateDnsRecordReq.Content = common.StringPtr(info.Value)
	teoCreateDnsRecordReq.TTL = common.Int64Ptr(int64(d.config.TTL))
	teoCreateDnsRecordResp, err := d.client.CreateDnsRecord(teoCreateDnsRecordReq)
	if err != nil {
		return fmt.Errorf("tencentcloud-eo: error when create record: %w", err)
	}

	d.recordIDsMu.Lock()
	d.recordIDs[token] = teoCreateDnsRecordResp.Response.RecordId
	d.recordIDsMu.Unlock()

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	d.recordIDsMu.Lock()
	recordID, ok := d.recordIDs[token]
	d.recordIDsMu.Unlock()
	if !ok {
		return fmt.Errorf("tencentcloud-eo: unknown record ID for '%s'", info.EffectiveFQDN)
	}

	// REF: https://cloud.tencent.com/document/product/1552/80718
	teoDeleteDnsRecordReq := tcteo.NewDeleteDnsRecordsRequest()
	teoDeleteDnsRecordReq.ZoneId = common.StringPtr(d.config.ZoneID)
	teoDeleteDnsRecordReq.RecordIds = []*string{recordID}
	_, err := d.client.DeleteDnsRecords(teoDeleteDnsRecordReq)
	if err != nil {
		return fmt.Errorf("tencentcloud-eo: error when delete record: %w", err)
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}
