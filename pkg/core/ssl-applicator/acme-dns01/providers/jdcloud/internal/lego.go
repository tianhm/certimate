package internal

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/platform/config/env"
	jdcore "github.com/jdcloud-api/jdcloud-sdk-go/core"
	jddnsapi "github.com/jdcloud-api/jdcloud-sdk-go/services/domainservice/apis"
	jddnsclient "github.com/jdcloud-api/jdcloud-sdk-go/services/domainservice/client"
	jddnsmodel "github.com/jdcloud-api/jdcloud-sdk-go/services/domainservice/models"
)

const (
	envNamespace = "JDCLOUD_"

	EnvAccessKeyID     = envNamespace + "ACCESS_KEY_ID"
	EnvAccessKeySecret = envNamespace + "ACCESS_KEY_SECRET"
	EnvRegionId        = envNamespace + "REGION_ID"

	EnvTTL                = envNamespace + "TTL"
	EnvPropagationTimeout = envNamespace + "PROPAGATION_TIMEOUT"
	EnvPollingInterval    = envNamespace + "POLLING_INTERVAL"
	EnvHTTPTimeout        = envNamespace + "HTTP_TIMEOUT"
)

var _ challenge.ProviderTimeout = (*DNSProvider)(nil)

type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	RegionId        string

	PropagationTimeout time.Duration
	PollingInterval    time.Duration
	TTL                int
	HTTPTimeout        time.Duration
}

type DNSProvider struct {
	client *jddnsclient.DomainserviceClient
	config *Config

	recordIDs   map[string]int
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
	values, err := env.Get(EnvAccessKeyID, EnvAccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("jdcloud: %w", err)
	}

	config := NewDefaultConfig()
	config.AccessKeyID = values[EnvAccessKeyID]
	config.AccessKeySecret = values[EnvAccessKeySecret]
	config.RegionId = values[EnvRegionId]

	return NewDNSProviderConfig(config)
}

func NewDNSProviderConfig(config *Config) (*DNSProvider, error) {
	if config == nil {
		return nil, errors.New("jdcloud: the configuration of the DNS provider is nil")
	}

	clientCredentials := jdcore.NewCredentials(config.AccessKeyID, config.AccessKeySecret)
	client := jddnsclient.NewDomainserviceClient(clientCredentials)
	clientConfig := &client.Config
	clientConfig.SetTimeout(config.HTTPTimeout)
	client.SetConfig(clientConfig)
	client.DisableLogger()

	return &DNSProvider{
		client:      client,
		config:      config,
		recordIDs:   make(map[string]int),
		recordIDsMu: sync.Mutex{},
	}, nil
}

func (d *DNSProvider) Present(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("jdcloud: could not find zone for domain %q: %w", domain, err)
	}

	subDomain, err := dns01.ExtractSubDomain(info.EffectiveFQDN, authZone)
	if err != nil {
		return fmt.Errorf("jdcloud: %w", err)
	}

	zone, err := d.getDNSZone(dns01.UnFqdn(authZone))
	if err != nil {
		return fmt.Errorf("jdcloud: error when list zones: %w", err)
	}

	// REF: https://docs.jdcloud.com/cn/jd-cloud-dns/api/createresourcerecord
	jddnsCreateResourceRecordReq := jddnsapi.NewCreateResourceRecordRequest(d.config.RegionId, fmt.Sprintf("%d", zone.Id), &jddnsmodel.AddRR{
		Type:       "TXT",
		HostRecord: subDomain,
		HostValue:  info.Value,
		Ttl:        int(d.config.TTL),
		ViewValue:  -1,
	})
	jddnsCreateResourceRecordResp, err := d.client.CreateResourceRecord(jddnsCreateResourceRecordReq)
	if err != nil {
		return fmt.Errorf("jdcloud: error when create record: %w", err)
	}

	d.recordIDsMu.Lock()
	d.recordIDs[token] = jddnsCreateResourceRecordResp.Result.DataList.Id
	d.recordIDsMu.Unlock()

	return nil
}

func (d *DNSProvider) CleanUp(domain, token, keyAuth string) error {
	info := dns01.GetChallengeInfo(domain, keyAuth)

	authZone, err := dns01.FindZoneByFqdn(info.EffectiveFQDN)
	if err != nil {
		return fmt.Errorf("jdcloud: could not find zone for domain %q: %w", domain, err)
	}

	d.recordIDsMu.Lock()
	recordID, ok := d.recordIDs[token]
	d.recordIDsMu.Unlock()
	if !ok {
		return fmt.Errorf("jdcloud: unknown record ID for '%s'", info.EffectiveFQDN)
	}

	zone, err := d.getDNSZone(dns01.UnFqdn(authZone))
	if err != nil {
		return fmt.Errorf("jdcloud: error when list zones: %w", err)
	}

	// REF: https://docs.jdcloud.com/cn/jd-cloud-dns/api/deleteresourcerecord
	jddnsDeleteResourceRecordReq := jddnsapi.NewDeleteResourceRecordRequest(d.config.RegionId, fmt.Sprintf("%d", zone.Id), fmt.Sprintf("%d", recordID))
	_, err = d.client.DeleteResourceRecord(jddnsDeleteResourceRecordReq)
	if err != nil {
		return fmt.Errorf("jdcloud: error when delete record: %w", err)
	}

	return nil
}

func (d *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return d.config.PropagationTimeout, d.config.PollingInterval
}

func (d *DNSProvider) getDNSZone(zoneName string) (*jddnsmodel.DomainInfo, error) {
	pageNumber := 1
	pageSize := 10
	for {
		// REF: https://docs.jdcloud.com/cn/jd-cloud-dns/api/describedomains
		jddnsDescribeDomainsReq := jddnsapi.NewDescribeDomainsRequest(d.config.RegionId, pageNumber, pageSize)
		jddnsDescribeDomainsReq.SetDomainName(zoneName)

		jddnsDescribeDomainsResp, err := d.client.DescribeDomains(jddnsDescribeDomainsReq)
		if err != nil {
			return nil, err
		}

		for _, item := range jddnsDescribeDomainsResp.Result.DataList {
			if item.DomainName == zoneName {
				return &item, nil
			}
		}

		if len(jddnsDescribeDomainsResp.Result.DataList) < pageSize {
			break
		}

		pageNumber++
	}

	return nil, fmt.Errorf("jdcloud: zone %s not found", zoneName)
}
