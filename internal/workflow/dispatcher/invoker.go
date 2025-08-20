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

func (w *workflowInvoker) processNode(ctx context.Context, graph *domain.WorkflowGraph) error {
	// current := graph
	// for current != nil {
	// 	select {
	// 	case <-ctx.Done():
	// 		return ctx.Err()
	// 	default:
	// 	}

	// 	if current.Type == domain.WorkflowNodeTypeCondition || current.Type == domain.WorkflowNodeTypeTryCatch {
	// 		for _, branch := range current.Branches {
	// 			if err := w.processNode(ctx, &branch); err != nil {
	// 				// 并行分支的某一分支发生错误时，忽略此错误，继续执行其他分支
	// 				if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
	// 					continue
	// 				}
	// 				return err
	// 			}
	// 		}
	// 	}

	// 	var processor nodes.NodeProcessor
	// 	var procErr error
	// 	for {
	// 		if current.Type != domain.WorkflowNodeTypeCondition && current.Type != domain.WorkflowNodeTypeTryCatch {
	// 			processor, procErr = nodes.GetProcessor(current)
	// 			if procErr != nil {
	// 				panic(procErr)
	// 			}

	// 			processor.SetLogger(slog.New(logging.NewHookHandler(&logging.HookHandlerOptions{
	// 				Level: slog.LevelDebug,
	// 				WriteFunc: func(ctx context.Context, record *logging.Record) error {
	// 					log := domain.WorkflowLog{}
	// 					log.WorkflowId = w.workflowId
	// 					log.RunId = w.runId
	// 					log.NodeId = current.Id
	// 					log.NodeName = current.Name
	// 					log.Timestamp = record.Time.UnixMilli()
	// 					log.Level = int32(record.Level)
	// 					log.Message = record.Message
	// 					log.Data = record.Data
	// 					log.CreatedAt = record.Time
	// 					if _, err := w.workflowLogRepo.Save(ctx, &log); err != nil {
	// 						return err
	// 					}

	// 					w.logs = append(w.logs, log)
	// 					return nil
	// 				},
	// 			})))

	// 			procErr = processor.Process(ctx)
	// 			if procErr != nil {
	// 				if current.Type != domain.WorkflowNodeTypeBranchBlock {
	// 					processor.GetLogger().Error(procErr.Error())
	// 				}
	// 				break
	// 			}

	// 			nodeOutputs := processor.GetOutputs()
	// 			if len(nodeOutputs) > 0 {
	// 				ctx = nodes.AddNodeOutput(ctx, current.Id, nodeOutputs)
	// 			}
	// 		}

	// 		break
	// 	}

	// 	// TODO: 优化可读性
	// 	if procErr != nil && current.Type == domain.WorkflowNodeTypeBranchBlock {
	// 		current = nil
	// 		procErr = nil
	// 		return nil
	// 	} else if procErr != nil && current.Next != nil && current.Next.Type != domain.WorkflowNodeTypeTryCatch {
	// 		return procErr
	// 	} else if procErr != nil && current.Next != nil && current.Next.Type == domain.WorkflowNodeTypeTryCatch {
	// 		current = w.getBranchByType(current.Next.Branches, domain.WorkflowNodeTypeCatchBlock)
	// 	} else if procErr == nil && current.Next != nil && current.Next.Type == domain.WorkflowNodeTypeTryCatch {
	// 		current = w.getBranchByType(current.Next.Branches, domain.WorkflowNodeTypeTryBlock)
	// 	} else {
	// 		current = current.Next
	// 	}
	// }

	return nil
}
