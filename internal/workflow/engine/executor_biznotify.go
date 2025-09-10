package engine

import (
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/internal/notify"
	"github.com/certimate-go/certimate/internal/repository"
)

type bizNotifyNodeExecutor struct {
	nodeExecutor

	accessRepo   accessRepository
	settingsRepo settingsRepository
}

func (ne *bizNotifyNodeExecutor) Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
	execRes := newNodeExecutionResult(execCtx.Node)

	nodeCfg := execCtx.Node.Data.Config.AsBizNotify()
	ne.logger.Info("ready to send notification ...", slog.Any("config", nodeCfg))

	// 检测是否可以跳过本次执行
	if skippable, reason := ne.checkCanSkip(execCtx); skippable {
		ne.logger.Info(fmt.Sprintf("skip this application, because %s", reason))
		return execRes, nil
	}

	// 读取部署提供商授权
	providerAccessConfig := make(map[string]any)
	if nodeCfg.ProviderAccessId != "" {
		if access, err := ne.accessRepo.GetById(execCtx.ctx, nodeCfg.ProviderAccessId); err != nil {
			return nil, fmt.Errorf("failed to get access #%s record: %w", nodeCfg.ProviderAccessId, err)
		} else {
			providerAccessConfig = access.Config
		}
	}

	// 推送通知
	notifier := notify.NewClient(notify.WithLogger(ne.logger))
	notifyReq := &notify.SendNotificationRequest{
		Provider:               nodeCfg.Provider,
		ProviderAccessConfig:   providerAccessConfig,
		ProviderExtendedConfig: nodeCfg.ProviderConfig,
		Subject:                nodeCfg.Subject,
		Message:                nodeCfg.Message,
	}
	if _, err := notifier.SendNotification(execCtx.ctx, notifyReq); err != nil {
		ne.logger.Warn("could not send notification")
		return execRes, err
	}

	ne.logger.Info("notification completed")
	return execRes, nil
}

func (ne *bizNotifyNodeExecutor) checkCanSkip(execCtx *NodeExecutionContext) (_skip bool, _reason string) {
	thisNodeCfg := execCtx.Node.Data.Config.AsBizNotify()
	if !thisNodeCfg.SkipOnAllPrevSkipped {
		return false, ""
	}

	var total, skipped int32
	for _, variable := range execCtx.variables.All() {
		if variable.Scope != "" && variable.Key == stateVarKeyNodeSkipped {
			total++
			if variable.Value == true {
				skipped++
			}
		}
	}
	if total == 0 || skipped != total {
		return false, ""
	}

	return true, "all the previous nodes have been skipped"
}

func newBizNotifyNodeExecutor() NodeExecutor {
	return &bizNotifyNodeExecutor{
		nodeExecutor: nodeExecutor{logger: slog.Default()},
		accessRepo:   repository.NewAccessRepository(),
		settingsRepo: repository.NewSettingsRepository(),
	}
}
