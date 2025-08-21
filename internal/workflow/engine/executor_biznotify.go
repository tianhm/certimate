package engine

import (
	"fmt"
	"log/slog"

	"github.com/certimate-go/certimate/internal/notify"
	"github.com/certimate-go/certimate/internal/repository"
)

type bizNotifyNodeExecutor struct {
	nodeExecutor

	settingsRepo settingsRepository
}

func (ne *bizNotifyNodeExecutor) Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
	execRes := newNodeExecutionResult(execCtx.Node)

	nodeCfg := execCtx.Node.Data.Config.AsBizNotify()
	ne.logger.Info("ready to send notification ...", slog.Any("config", nodeCfg))

	// 检测是否可以跳过本次执行
	if skippable := ne.checkCanSkip(execCtx); skippable {
		ne.logger.Info(fmt.Sprintf("skip this notification, because all the previous nodes have been skipped"))
		return execRes, nil
	}

	// 初始化通知器
	deployer, err := notify.NewWithWorkflowNode(notify.NotifierWithWorkflowNodeConfig{
		Node:    execCtx.Node,
		Logger:  ne.logger,
		Subject: nodeCfg.Subject,
		Message: nodeCfg.Message,
	})
	if err != nil {
		ne.logger.Warn("failed to create notifier provider")
		return execRes, err
	}

	// 推送通知
	if err := deployer.Notify(execCtx.ctx); err != nil {
		ne.logger.Warn("failed to send notification")
		return execRes, err
	}

	ne.logger.Info("notification completed")
	return execRes, nil
}

func (ne *bizNotifyNodeExecutor) checkCanSkip(execCtx *NodeExecutionContext) (_skip bool) {
	thisNodeCfg := execCtx.Node.Data.Config.AsBizNotify()
	if !thisNodeCfg.SkipOnAllPrevSkipped {
		return false
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
	return total > 0 && skipped == total
}

func newBizNotifyNodeExecutor() NodeExecutor {
	return &bizNotifyNodeExecutor{
		nodeExecutor: nodeExecutor{logger: slog.Default()},
		settingsRepo: repository.NewSettingsRepository(),
	}
}
