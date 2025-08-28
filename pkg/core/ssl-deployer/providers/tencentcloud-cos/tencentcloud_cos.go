package tencentcloudcos

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/certimate-go/certimate/pkg/core"
	sslmgrsp "github.com/certimate-go/certimate/pkg/core/ssl-manager/providers/tencentcloud-ssl"
)

type SSLDeployerProviderConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 腾讯云地域。
	Region string `json:"region"`
	// 存储桶名。
	Bucket string `json:"bucket"`
	// 自定义域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type SSLDeployerProvider struct {
	config     *SSLDeployerProviderConfig
	logger     *slog.Logger
	sdkClient  *wSDKClients
	sslManager core.SSLManager
}

var _ core.SSLDeployer = (*SSLDeployerProvider)(nil)

type wSDKClients struct {
	SSL *tcssl.Client
}

func NewSSLDeployerProvider(config *SSLDeployerProviderConfig) (*SSLDeployerProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the ssl deployer provider is nil")
	}

	clients, err := createSDKClients(config.SecretId, config.SecretKey, config.Region)
	if err != nil {
		return nil, fmt.Errorf("could not create sdk client: %w", err)
	}

	sslmgr, err := sslmgrsp.NewSSLManagerProvider(&sslmgrsp.SSLManagerProviderConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create ssl manager: %w", err)
	}

	return &SSLDeployerProvider{
		config:     config,
		logger:     slog.Default(),
		sdkClient:  clients,
		sslManager: sslmgr,
	}, nil
}

func (d *SSLDeployerProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		d.logger = slog.New(slog.DiscardHandler)
	} else {
		d.logger = logger
	}

	d.sslManager.SetLogger(logger)
}

func (d *SSLDeployerProvider) Deploy(ctx context.Context, certPEM string, privkeyPEM string) (*core.SSLDeployResult, error) {
	if d.config.Bucket == "" {
		return nil, errors.New("config `bucket` is required")
	}
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 上传证书
	upres, err := d.sslManager.Upload(ctx, certPEM, privkeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to upload certificate file: %w", err)
	} else {
		d.logger.Info("ssl certificate uploaded", slog.Any("result", upres))
	}

	// 避免多次部署，否则会报错 https://github.com/certimate-go/certimate/issues/897#issuecomment-3182904098
	if bind, _ := d.checkIsBind(ctx, upres.CertId); bind {
		d.logger.Info("ssl certificate already deployed")
		return &core.SSLDeployResult{}, nil
	}

	// 证书部署到 COS 实例
	// REF: https://cloud.tencent.com/document/api/400/91667
	deployCertificateInstanceReq := tcssl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(upres.CertId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("cos")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s|%s|%s", d.config.Region, d.config.Bucket, d.config.Domain)})
	deployCertificateInstanceResp, err := d.sdkClient.SSL.DeployCertificateInstance(deployCertificateInstanceReq)
	d.logger.Debug("sdk request 'ssl.DeployCertificateInstance'", slog.Any("request", deployCertificateInstanceReq), slog.Any("response", deployCertificateInstanceResp))
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'ssl.DeployCertificateInstance': %w", err)
	}

	// 循环获取部署任务详情，等待任务状态变更
	// REF: https://cloud.tencent.com/document/api/400/91658
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		describeHostDeployRecordDetailReq := tcssl.NewDescribeHostDeployRecordDetailRequest()
		describeHostDeployRecordDetailReq.DeployRecordId = common.StringPtr(fmt.Sprintf("%d", *deployCertificateInstanceResp.Response.DeployRecordId))
		describeHostDeployRecordDetailResp, err := d.sdkClient.SSL.DescribeHostDeployRecordDetail(describeHostDeployRecordDetailReq)
		d.logger.Debug("sdk request 'ssl.DescribeHostDeployRecordDetail'", slog.Any("request", describeHostDeployRecordDetailReq), slog.Any("response", describeHostDeployRecordDetailResp))
		if err != nil {
			return nil, fmt.Errorf("failed to execute sdk request 'ssl.DescribeHostDeployRecordDetail': %w", err)
		}

		var pendingCount, runningCount, succeededCount, failedCount, totalCount int64
		if describeHostDeployRecordDetailResp.Response.TotalCount == nil {
			return nil, errors.New("unexpected tencentcloud deployment job status")
		} else {
			if describeHostDeployRecordDetailResp.Response.PendingTotalCount != nil {
				pendingCount = *describeHostDeployRecordDetailResp.Response.PendingTotalCount
			}
			if describeHostDeployRecordDetailResp.Response.RunningTotalCount != nil {
				runningCount = *describeHostDeployRecordDetailResp.Response.RunningTotalCount
			}
			if describeHostDeployRecordDetailResp.Response.SuccessTotalCount != nil {
				succeededCount = *describeHostDeployRecordDetailResp.Response.SuccessTotalCount
			}
			if describeHostDeployRecordDetailResp.Response.FailedTotalCount != nil {
				failedCount = *describeHostDeployRecordDetailResp.Response.FailedTotalCount
			}
			if describeHostDeployRecordDetailResp.Response.TotalCount != nil {
				totalCount = *describeHostDeployRecordDetailResp.Response.TotalCount
			}

			if succeededCount+failedCount == totalCount {
				if failedCount > 0 {
					return nil, fmt.Errorf("tencentcloud deployment job failed (succeeded: %d, failed: %d, total: %d)", succeededCount, failedCount, totalCount)
				}
				break
			}
		}

		d.logger.Info(fmt.Sprintf("waiting for tencentcloud deployment job completion (pending: %d, running: %d, succeeded: %d, failed: %d, total: %d) ...", pendingCount, runningCount, succeededCount, failedCount, totalCount))
		time.Sleep(time.Second * 5)
	}

	return &core.SSLDeployResult{}, nil
}

func (d *SSLDeployerProvider) checkIsBind(ctx context.Context, cloudCertId string) (bool, error) {
	// 创建证书关联云资源异步任务
	// REF: https://cloud.tencent.com/document/api/400/97759
	createCertificateBindResourceSyncTaskReq := tcssl.NewCreateCertificateBindResourceSyncTaskRequest()
	createCertificateBindResourceSyncTaskReq.CertificateIds = []*string{common.StringPtr(cloudCertId)}
	createCertificateBindResourceSyncTaskReq.IsCache = common.Uint64Ptr(0)
	createCertificateBindResourceSyncTaskResp, err := d.sdkClient.SSL.CreateCertificateBindResourceSyncTask(createCertificateBindResourceSyncTaskReq)
	d.logger.Debug("sdk request 'ssl.CreateCertificateBindResourceSyncTask'", slog.Any("request", createCertificateBindResourceSyncTaskReq), slog.Any("response", createCertificateBindResourceSyncTaskResp))
	if err != nil {
		return false, fmt.Errorf("failed to execute sdk request 'ssl.CreateCertificateBindResourceSyncTask': %w", err)
	}

	// 查询证书关联云资源任务结果
	// REF: https://cloud.tencent.com/document/api/400/97758
	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}

		describeCertificateBindResourceTaskDetailReq := tcssl.NewDescribeCertificateBindResourceTaskDetailRequest()
		describeCertificateBindResourceTaskDetailReq.TaskId = createCertificateBindResourceSyncTaskResp.Response.CertTaskIds[0].TaskId
		describeCertificateBindResourceTaskDetailReq.ResourceTypes = []*string{common.StringPtr("cos")}
		describeCertificateBindResourceTaskDetailReq.Regions = []*string{common.StringPtr(d.config.Region)}
		describeCertificateBindResourceTaskDetailReq.Offset = common.StringPtr("0")
		describeCertificateBindResourceTaskDetailReq.Limit = common.StringPtr("100")
		describeCertificateBindResourceTaskDetailResp, err := d.sdkClient.SSL.DescribeCertificateBindResourceTaskDetail(describeCertificateBindResourceTaskDetailReq)
		d.logger.Debug("sdk request 'ssl.DescribeCertificateBindResourceTaskDetail'", slog.Any("request", describeCertificateBindResourceTaskDetailReq), slog.Any("response", describeCertificateBindResourceTaskDetailResp))
		if err != nil {
			return false, fmt.Errorf("failed to execute sdk request 'ssl.DescribeCertificateBindResourceTaskDetail': %w", err)
		}

		if describeCertificateBindResourceTaskDetailResp.Response.Status == nil || *describeCertificateBindResourceTaskDetailResp.Response.Status == 2 {
			return false, errors.New("unexpected tencentcloud query task status")
		} else if *describeCertificateBindResourceTaskDetailResp.Response.Status == 1 {
			for _, record := range describeCertificateBindResourceTaskDetailResp.Response.COS {
				for _, instance := range record.InstanceList {
					if instance.Bucket == nil || *instance.Bucket != d.config.Bucket {
						continue
					}
					if instance.Domain == nil || *instance.Domain != d.config.Domain {
						continue
					}
					if instance.Status == nil || *instance.Status != "ENABLED" {
						continue
					}
					return true, nil
				}
			}
			return false, nil
		}

		d.logger.Info("waiting for tencentcloud query task completion")
		time.Sleep(time.Second * 5)
	}
}

func createSDKClients(secretId, secretKey, region string) (*wSDKClients, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tcssl.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &wSDKClients{
		SSL: client,
	}, nil
}
