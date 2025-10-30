package larkbot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/go-resty/resty/v2"

	"github.com/certimate-go/certimate/pkg/core"
)

type NotifierProviderConfig struct {
	// 飞书机器人 Webhook 地址。
	WebhookUrl string `json:"webhookUrl"`
}

type NotifierProvider struct {
	config     *NotifierProviderConfig
	logger     *slog.Logger
	httpClient *resty.Client
}

var _ core.Notifier = (*NotifierProvider)(nil)

func NewNotifierProvider(config *NotifierProviderConfig) (*NotifierProvider, error) {
	if config == nil {
		return nil, errors.New("the configuration of the notifier provider is nil")
	}

	client := resty.New().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "certimate")

	return &NotifierProvider{
		config:     config,
		logger:     slog.Default(),
		httpClient: client,
	}, nil
}

func (n *NotifierProvider) SetLogger(logger *slog.Logger) {
	if logger == nil {
		n.logger = slog.New(slog.DiscardHandler)
	} else {
		n.logger = logger
	}
}

func (n *NotifierProvider) Notify(ctx context.Context, subject string, message string) (*core.NotifyResult, error) {
	webhookUrl, err := url.Parse(n.config.WebhookUrl)
	if err != nil {
		return nil, fmt.Errorf("lark api error: invalid webhook url: %w", err)
	} else {
		const hostname = "open.larksuite.com"
		const hostname_cn = "open.feishu.cn"
		if webhookUrl.Hostname() != hostname && webhookUrl.Hostname() != hostname_cn {
			n.logger.Warn(fmt.Sprintf("the webhook url hostname is not '%s' or '%s', please make sure it is correct", hostname, hostname_cn))
		}
	}

	// REF: https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot
	// REF: https://open.larksuite.com/document/client-docs/bot-v3/add-custom-bot
	var result struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}
	req := n.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]any{
			"msg_type": "text",
			"content": map[string]string{
				"text": subject + "\n\n" + message,
			},
		})
	resp, err := req.Post(webhookUrl.String())
	if err != nil {
		return nil, fmt.Errorf("lark api error: failed to send request: %w", err)
	} else if resp.IsError() {
		return nil, fmt.Errorf("lark api error: unexpected status code: %d, resp: %s", resp.StatusCode(), resp.String())
	} else if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("lark api error: failed to unmarshal response: %w", err)
	} else if result.Code != 0 {
		return nil, fmt.Errorf("lark api error: code='%d', msg='%s'", result.Code, result.Message)
	}

	return &core.NotifyResult{}, nil
}
