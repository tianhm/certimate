package dispatcher

import (
	"context"
	"log/slog"
	"time"

	"github.com/certimate-go/certimate/internal/app"
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/workflow/engine"
	"github.com/certimate-go/certimate/pkg/logging"
)

type workflowInvoker struct {
	workflowId    string
	workflowGraph *domain.WorkflowGraph
	runId         string
	logs          []domain.WorkflowLog

	workflowLogRepo workflowLogRepository
}

func newWorkflowInvokerWithData(workflowLogRepo workflowLogRepository, data *WorkflowWorkerData) *workflowInvoker {
	if data == nil {
		panic("workflow dispatcher: the worker data is nil")
	}

	return &workflowInvoker{
		workflowId:    data.WorkflowId,
		workflowGraph: data.WorkflowGraph,
		runId:         data.RunId,
		logs:          make([]domain.WorkflowLog, 0),

		workflowLogRepo: workflowLogRepo,
	}
}

func (w *workflowInvoker) Invoke(ctx context.Context) error {
	we := engine.NewWorkflowEngine()
	// TODO: 优化日志记录
	we.OnNodeError(func(ctx context.Context, node *engine.Node, err error) error {
		log := domain.WorkflowLog{}
		log.WorkflowId = w.workflowId
		log.RunId = w.runId
		log.NodeId = node.Id
		log.NodeName = node.Data.Name
		log.Timestamp = time.Now().UnixMilli()
		log.Level = int32(slog.LevelError)
		log.Message = err.Error()
		log.CreatedAt = time.Now()
		if _, err := w.workflowLogRepo.Save(ctx, &log); err != nil {
			app.GetLogger().Error(err.Error())
		}

		w.logs = append(w.logs, log)
		return nil
	})
	we.OnNodeLogging(func(ctx context.Context, node *engine.Node, record logging.Record) error {
		log := domain.WorkflowLog{}
		log.WorkflowId = w.workflowId
		log.RunId = w.runId
		log.NodeId = node.Id
		log.NodeName = node.Data.Name
		log.Timestamp = record.Time.UnixMilli()
		log.Level = int32(record.Level)
		log.Message = record.Message
		log.Data = record.Data
		log.CreatedAt = record.Time
		if _, err := w.workflowLogRepo.Save(ctx, &log); err != nil {
			app.GetLogger().Error(err.Error())
		}

		w.logs = append(w.logs, log)
		return nil
	})

	err := we.Invoke(ctx, w.workflowId, w.runId, w.workflowGraph)
	if err != nil {
		return err
	}

	return nil
}

func (w *workflowInvoker) GetLogs() domain.WorkflowLogs {
	return w.logs
}
