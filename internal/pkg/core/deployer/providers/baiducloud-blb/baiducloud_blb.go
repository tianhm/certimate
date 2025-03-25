﻿package baiducloudblb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	bceblb "github.com/baidubce/bce-sdk-go/services/blb"
	xerrors "github.com/pkg/errors"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/baiducloud-cas"
	"github.com/usual2970/certimate/internal/pkg/utils/sliceutil"
)

type DeployerConfig struct {
	// 百度智能云 AccessKeyId。
	AccessKeyId string `json:"accessKeyId"`
	// 百度智能云 SecretAccessKey。
	SecretAccessKey string `json:"secretAccessKey"`
	// 百度智能云区域。
	Region string `json:"region"`
	// 部署资源类型。
	ResourceType ResourceType `json:"resourceType"`
	// 负载均衡实例 ID。
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_LISTENER] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty"`
	// 负载均衡监听端口。
	// 部署资源类型为 [RESOURCE_TYPE_LISTENER] 时必填。
	ListenerPort int32 `json:"listenerPort,omitempty"`
	// SNI 域名（支持泛域名）。
	// 部署资源类型为 [RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_LISTENER] 时选填。
	Domain string `json:"domain,omitempty"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      *slog.Logger
	sdkClient   *bceblb.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.AccessKeyId, config.SecretAccessKey, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk client")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		AccessKeyId:     config.AccessKeyId,
		SecretAccessKey: config.SecretAccessKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &DeployerProvider{
		config:      config,
		logger:      slog.Default(),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger *slog.Logger) deployer.Deployer {
	if logger == nil {
		d.logger = slog.Default()
	} else {
		d.logger = logger
	}
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	// 上传证书到 CAS
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 根据部署资源类型决定部署方式
	switch d.config.ResourceType {
	case RESOURCE_TYPE_LOADBALANCER:
		if err := d.deployToLoadbalancer(ctx, upres.CertId); err != nil {
			return nil, err
		}

	case RESOURCE_TYPE_LISTENER:
		if err := d.deployToListener(ctx, upres.CertId); err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported resource type: %s", d.config.ResourceType)
	}

	return &deployer.DeployResult{}, nil
}

func (d *DeployerProvider) deployToLoadbalancer(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}

	// 查询 BLB 实例详情
	// REF: https://cloud.baidu.com/doc/BLB/s/njwvxnv79#describeloadbalancerdetail%E6%9F%A5%E8%AF%A2blb%E5%AE%9E%E4%BE%8B%E8%AF%A6%E6%83%85
	describeLoadBalancerDetailResp, err := d.sdkClient.DescribeLoadBalancerDetail(d.config.LoadbalancerId)
	d.logger.Debug("sdk request 'blb.DescribeLoadBalancerAttribute'", slog.String("blbId", d.config.LoadbalancerId), slog.Any("response", describeLoadBalancerDetailResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'blb.DescribeLoadBalancerDetail'")
	}

	// 获取全部 HTTPS/SSL 监听端口
	listeners := make([]struct {
		Type string
		Port int32
	}, 0)
	for _, listener := range describeLoadBalancerDetailResp.Listener {
		if listener.Type == "HTTPS" || listener.Type == "SSL" {
			listenerPort, err := strconv.Atoi(listener.Port)
			if err != nil {
				continue
			}

			listeners = append(listeners, struct {
				Type string
				Port int32
			}{
				Type: listener.Type,
				Port: int32(listenerPort),
			})
		}
	}

	// 遍历更新监听证书
	if len(listeners) == 0 {
		d.logger.Info("no blb listeners to deploy")
	} else {
		d.logger.Info("found https/ssl listeners to deploy", slog.Any("listeners", listeners))
		var errs []error

		for _, listener := range listeners {
			if err := d.updateListenerCertificate(ctx, d.config.LoadbalancerId, listener.Type, listener.Port, cloudCertId); err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	}

	return nil
}

func (d *DeployerProvider) deployToListener(ctx context.Context, cloudCertId string) error {
	if d.config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}
	if d.config.ListenerPort == 0 {
		return errors.New("config `listenerPort` is required")
	}

	// 查询监听
	// REF: https://cloud.baidu.com/doc/BLB/s/yjwvxnvl6#describealllisteners%E6%9F%A5%E8%AF%A2%E6%89%80%E6%9C%89%E7%9B%91%E5%90%AC
	describeAllListenersRequest := &bceblb.DescribeListenerArgs{
		ListenerPort: uint16(d.config.ListenerPort),
	}
	describeAllListenersResp, err := d.sdkClient.DescribeAllListeners(d.config.LoadbalancerId, describeAllListenersRequest)
	d.logger.Debug("sdk request 'blb.DescribeAllListeners'", slog.String("blbId", d.config.LoadbalancerId), slog.Any("request", describeAllListenersRequest), slog.Any("response", describeAllListenersResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'blb.DescribeAllListeners'")
	}

	// 获取全部 HTTPS/SSL 监听端口
	listeners := make([]struct {
		Type string
		Port int32
	}, 0)
	for _, listener := range describeAllListenersResp.AllListenerList {
		if listener.ListenerType == "HTTPS" || listener.ListenerType == "SSL" {
			listeners = append(listeners, struct {
				Type string
				Port int32
			}{
				Type: listener.ListenerType,
				Port: int32(listener.ListenerPort),
			})
		}
	}

	// 遍历更新监听证书
	if len(listeners) == 0 {
		d.logger.Info("no blb listeners to deploy")
	} else {
		d.logger.Info("found https/ssl listeners to deploy", slog.Any("listeners", listeners))
		var errs []error

		for _, listener := range listeners {
			if err := d.updateListenerCertificate(ctx, d.config.LoadbalancerId, listener.Type, listener.Port, cloudCertId); err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	}

	return nil
}

func (d *DeployerProvider) updateListenerCertificate(ctx context.Context, cloudLoadbalancerId string, cloudListenerType string, cloudListenerPort int32, cloudCertId string) error {
	switch strings.ToUpper(cloudListenerType) {
	case "HTTPS":
		return d.updateHttpsListenerCertificate(ctx, cloudLoadbalancerId, cloudListenerPort, cloudCertId)
	case "SSL":
		return d.updateSslListenerCertificate(ctx, cloudLoadbalancerId, cloudListenerPort, cloudCertId)
	default:
		return fmt.Errorf("unsupported listener type: %s", cloudListenerType)
	}
}

func (d *DeployerProvider) updateHttpsListenerCertificate(ctx context.Context, cloudLoadbalancerId string, cloudHttpsListenerPort int32, cloudCertId string) error {
	// 查询 HTTPS 监听器
	// REF: https://cloud.baidu.com/doc/BLB/s/yjwvxnvl6#describehttpslisteners%E6%9F%A5%E8%AF%A2https%E7%9B%91%E5%90%AC%E5%99%A8
	describeHTTPSListenersReq := &bceblb.DescribeListenerArgs{
		ListenerPort: uint16(cloudHttpsListenerPort),
		MaxKeys:      1,
	}
	describeHTTPSListenersResp, err := d.sdkClient.DescribeHTTPSListeners(cloudLoadbalancerId, describeHTTPSListenersReq)
	d.logger.Debug("sdk request 'blb.DescribeHTTPSListeners'", slog.String("blbId", cloudLoadbalancerId), slog.Any("request", describeHTTPSListenersReq), slog.Any("response", describeHTTPSListenersResp))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'blb.DescribeHTTPSListeners'")
	} else if len(describeHTTPSListenersResp.ListenerList) == 0 {
		return fmt.Errorf("listener %s:%d not found", cloudLoadbalancerId, cloudHttpsListenerPort)
	}

	if d.config.Domain == "" {
		// 未指定 SNI，只需部署到监听器

		// 更新 HTTPS 监听器
		// REF: https://cloud.baidu.com/doc/BLB/s/yjwvxnvl6#updatehttpslistener%E6%9B%B4%E6%96%B0https%E7%9B%91%E5%90%AC%E5%99%A8
		updateHTTPSListenerReq := &bceblb.UpdateHTTPSListenerArgs{
			ListenerPort: uint16(cloudHttpsListenerPort),
			CertIds:      []string{cloudCertId},
		}
		err := d.sdkClient.UpdateHTTPSListener(cloudLoadbalancerId, updateHTTPSListenerReq)
		d.logger.Debug("sdk request 'blb.UpdateHTTPSListener'", slog.Any("request", updateHTTPSListenerReq))
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'blb.UpdateHTTPSListener'")
		}
	} else {
		// 指定 SNI，需部署到扩展域名

		// 更新 HTTPS 监听器
		// REF: https://cloud.baidu.com/doc/BLB/s/yjwvxnvl6#updatehttpslistener%E6%9B%B4%E6%96%B0https%E7%9B%91%E5%90%AC%E5%99%A8
		updateHTTPSListenerReq := &bceblb.UpdateHTTPSListenerArgs{
			ListenerPort: uint16(cloudHttpsListenerPort),
			AdditionalCertDomains: sliceutil.Map(describeHTTPSListenersResp.ListenerList[0].AdditionalCertDomains, func(domain bceblb.AdditionalCertDomainsModel) bceblb.AdditionalCertDomainsModel {
				if domain.Host == d.config.Domain {
					return bceblb.AdditionalCertDomainsModel{
						Host:   domain.Host,
						CertId: cloudCertId,
					}
				}

				return bceblb.AdditionalCertDomainsModel{
					Host:   domain.Host,
					CertId: domain.CertId,
				}
			}),
		}
		err := d.sdkClient.UpdateHTTPSListener(cloudLoadbalancerId, updateHTTPSListenerReq)
		d.logger.Debug("sdk request 'blb.UpdateHTTPSListener'", slog.Any("request", updateHTTPSListenerReq))
		if err != nil {
			return xerrors.Wrap(err, "failed to execute sdk request 'blb.UpdateHTTPSListener'")
		}
	}

	return nil
}

func (d *DeployerProvider) updateSslListenerCertificate(ctx context.Context, cloudLoadbalancerId string, cloudHttpsListenerPort int32, cloudCertId string) error {
	// 更新 SSL 监听器
	// REF: https://cloud.baidu.com/doc/BLB/s/yjwvxnvl6#updatessllistener%E6%9B%B4%E6%96%B0ssl%E7%9B%91%E5%90%AC%E5%99%A8
	updateSSLListenerReq := &bceblb.UpdateSSLListenerArgs{
		ListenerPort: uint16(cloudHttpsListenerPort),
		CertIds:      []string{cloudCertId},
	}
	err := d.sdkClient.UpdateSSLListener(cloudLoadbalancerId, updateSSLListenerReq)
	d.logger.Debug("sdk request 'blb.UpdateSSLListener'", slog.Any("request", updateSSLListenerReq))
	if err != nil {
		return xerrors.Wrap(err, "failed to execute sdk request 'blb.UpdateSSLListener'")
	}

	return nil
}

func createSdkClient(accessKeyId, secretAccessKey, region string) (*bceblb.Client, error) {
	endpoint := ""
	if region != "" {
		endpoint = fmt.Sprintf("blb.%s.baidubce.com", region)
	}

	client, err := bceblb.NewClient(accessKeyId, secretAccessKey, endpoint)
	if err != nil {
		return nil, err
	}

	return client, nil
}
