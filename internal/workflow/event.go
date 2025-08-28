package workflow

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pocketbase/pocketbase/core"

	"github.com/certimate-go/certimate/internal/app"
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/domain/dtos"
	"github.com/certimate-go/certimate/internal/repository"
)

func Register() {
	pb := app.GetApp()
	pb.OnRecordCreateRequest(domain.CollectionNameWorkflow).BindFunc(func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		if err := onWorkflowRecordCreateOrUpdate(e.Request.Context(), e.Record); err != nil {
			app.GetLogger().Error(err.Error())
			return err
		}

		return nil
	})
	pb.OnRecordUpdateRequest(domain.CollectionNameWorkflow).BindFunc(func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		if err := onWorkflowRecordCreateOrUpdate(e.Request.Context(), e.Record); err != nil {
			app.GetLogger().Error(err.Error())
			return err
		}

		return nil
	})
	pb.OnRecordDeleteRequest(domain.CollectionNameWorkflow).BindFunc(func(e *core.RecordRequestEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		if err := onWorkflowRecordDelete(e.Request.Context(), e.Record); err != nil {
			app.GetLogger().Error(err.Error())
			return err
		}

		return nil
	})
}

func onWorkflowRecordCreateOrUpdate(ctx context.Context, record *core.Record) error {
	scheduler := app.GetScheduler()

	// 向数据库插入/更新时，同时更新定时任务
	workflowId := record.Id
	enabled := record.GetBool("enabled")
	trigger := record.GetString("trigger")

	// 如果非定时触发或未启用，移除定时任务
	if !enabled || trigger != string(domain.WorkflowTriggerTypeScheduled) {
		scheduler.Remove(fmt.Sprintf("workflow#%s", workflowId))
		return nil
	}

	// 反之，重新添加定时任务
	err := scheduler.Add(fmt.Sprintf("workflow#%s", workflowId), record.GetString("triggerCron"), func() {
		workflowSrv := NewWorkflowService(repository.NewWorkflowRepository(), repository.NewWorkflowRunRepository(), repository.NewSettingsRepository())
		workflowSrv.StartRun(ctx, &dtos.WorkflowStartRunReq{
			WorkflowId: workflowId,
			RunTrigger: domain.WorkflowTriggerTypeScheduled,
		})
	})
	if err != nil {
		app.GetLogger().Error(fmt.Sprintf("failed to register cron job for workflow #%s", record.Id), slog.Any("error", err))
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	return nil
}

func onWorkflowRecordDelete(_ context.Context, record *core.Record) error {
	scheduler := app.GetScheduler()

	// 从数据库删除时，同时移除定时任务
	workflowId := record.Id
	scheduler.Remove(fmt.Sprintf("workflow#%s", workflowId))

	return nil
}
