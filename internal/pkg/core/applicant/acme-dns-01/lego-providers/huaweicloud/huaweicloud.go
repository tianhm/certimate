package huaweicloud

import (
	"time"

	"github.com/go-acme/lego/v4/challenge"
	hwc "github.com/go-acme/lego/v4/providers/dns/huaweicloud"
)

type HuaweiCloudApplicantConfig struct {
	AccessKeyId           string `json:"accessKeyId"`
	SecretAccessKey       string `json:"secretAccessKey"`
	Region                string `json:"region"`
	DnsPropagationTimeout int32  `json:"dnsPropagationTimeout,omitempty"`
	DnsTTL                int32  `json:"dnsTTL,omitempty"`
}

func NewChallengeProvider(config *HuaweiCloudApplicantConfig) (challenge.Provider, error) {
	if config == nil {
		panic("config is nil")
	}

	region := config.Region
	if region == "" {
		// 华为云的 SDK 要求必须传一个区域，实际上 DNS-01 流程里用不到，但不传会报错
		region = "cn-north-1"
	}

	providerConfig := hwc.NewDefaultConfig()
	providerConfig.AccessKeyID = config.AccessKeyId
	providerConfig.SecretAccessKey = config.SecretAccessKey
	providerConfig.Region = region
	if config.DnsPropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.DnsPropagationTimeout) * time.Second
	}
	if config.DnsTTL != 0 {
		providerConfig.TTL = config.DnsTTL
	}

	provider, err := hwc.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
