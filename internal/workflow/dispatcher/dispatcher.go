package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/certimate-go/certimate/internal/app"
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
	"github.com/certimate-go/certimate/internal/workflow/engine"
	"github.com/certimate-go/certimate/pkg/logging"
)

var maxWorkers = 1

func init() {
	envMaxWorkers := os.Getenv("CERTIMATE_WORKFLOW_MAX_WORKERS")
	if n, err := strconv.Atoi(envMaxWorkers); err != nil && n > 0 {
		maxWorkers = n
	} else {
		maxWorkers = runtime.GOMAXPROCS(0)
		if maxWorkers == 0 {
			maxWorkers = max(1, runtime.NumCPU())
		}
	}
}

type WorkflowDispatcher interface {
	Bootup(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Start(ctx context.Context, runId string) error
	Cancel(ctx context.Context, runId string) error
}

type workflowDispatcher struct {
	booted      bool
	concurrency int

	taskMtx         sync.RWMutex
	pendingRunQueue []string
	processingTasks map[string]*taskInfo // Key: RunId

	workflowRepo    workflowRepository
	workflowRunRepo workflowRunRepository
	workflowLogRepo workflowLogRepository
}

var _ WorkflowDispatcher = (*workflowDispatcher)(nil)

func (wd *workflowDispatcher) Bootup(ctx context.Context) error {
	if wd.booted {
		return errors.New("could not re-bootup")
	}

	wd.taskMtx.Lock()
	defer wd.taskMtx.Unlock()

	if _, err := app.GetDB().NewQuery("UPDATE workflow SET lastRunStatus = 'canceled' WHERE lastRunStatus = 'pending' OR lastRunStatus = 'processing'").Execute(); err != nil {
		return err
	}
	if _, err := app.GetDB().NewQuery("UPDATE workflow_run SET status = 'canceled' WHERE status = 'pending' OR status = 'processing'").Execute(); err != nil {
		return err
	}

	wd.booted = true
	return nil
}

func (wd *workflowDispatcher) Shutdown(ctx context.Context) error {
	if !wd.booted {
		return errors.New("could not re-shutdown")
	}

	wd.taskMtx.Lock()
	defer wd.taskMtx.Unlock()

	for runId, task := range wd.processingTasks {
		task.cancel()
		delete(wd.processingTasks, runId)
	}

	wd.booted = false
	wd.pendingRunQueue = make([]string, 0)
	wd.processingTasks = make(map[string]*taskInfo)
	return nil
}

func (wd *workflowDispatcher) Start(ctx context.Context, runId string) error {
	wd.taskMtx.Lock()
	defer wd.taskMtx.Unlock()

	if _, exists := wd.processingTasks[runId]; exists {
		return errors.New("workflow run is already processing")
	}

	for _, pendingRunId := range wd.pendingRunQueue {
		if pendingRunId == runId {
			return errors.New("workflow run is already in the queue")
		}
	}

	wd.pendingRunQueue = append(wd.pendingRunQueue, runId)
	go func() { wd.tryNextAsync() }()

	return nil
}

func (wd *workflowDispatcher) Cancel(ctx context.Context, runId string) error {
	wd.taskMtx.Lock()
	defer wd.taskMtx.Unlock()

	workflowRun, err := wd.workflowRunRepo.GetById(ctx, runId)
	if err != nil {
		return err
	} else if workflowRun.Status != domain.WorkflowRunStatusTypePending && workflowRun.Status != domain.WorkflowRunStatusTypeProcessing {
		return errors.New("workflow run is already completed")
	}

	workflow, err := wd.workflowRepo.GetById(ctx, workflowRun.WorkflowId)
	if err != nil {
		return err
	}

	workflowRun.Status = domain.WorkflowRunStatusTypeCanceled
	if workflow.LastRunId == workflowRun.Id {
		_, err := wd.workflowRunRepo.SaveWithCascading(ctx, workflowRun)
		if err != nil {
			return err
		}
	} else {
		_, err := wd.workflowRunRepo.Save(ctx, workflowRun)
		if err != nil {
			return err
		}
	}

	if task, exists := wd.processingTasks[runId]; exists {
		task.cancel()
		delete(wd.processingTasks, runId)
	}

	for i, pendingRunId := range wd.pendingRunQueue {
		if pendingRunId == runId {
			wd.pendingRunQueue = append(wd.pendingRunQueue[:i], wd.pendingRunQueue[i+1:]...)
			break
		}
	}

	go func() { wd.tryNextAsync() }()

	return nil
}

func (wd *workflowDispatcher) tryExecuteAsync(task *taskInfo) {
	var workflowRun *domain.WorkflowRun

	// 捕获 panic
	defer func() {
		if r := recover(); r != nil {
			slog.Default().Warn(fmt.Sprintf("workflow dispatcher panic: %v, stack trace: %s", r, string(debug.Stack())), slog.Any("workflowId", task.WorkflowId), slog.Any("runId", task.RunId))
			app.GetLogger().Error(fmt.Sprintf("workflow dispatcher panic: %v", r), slog.Any("workflowId", task.WorkflowId), slog.Any("runId", task.RunId))

			if workflowRun != nil {
				workflowRun.Status = domain.WorkflowRunStatusTypeFailed
				workflowRun.EndedAt = time.Now()
				workflowRun.Error = fmt.Sprintf("workflow dispatcher panic: %v", r)
				if _, err := wd.workflowRunRepo.SaveWithCascading(context.Background(), workflowRun); err != nil {
					log.Default().Println("failed to save workflow run after panic", slog.Any("error", err))
				}
			}
		}
	}()

	// 尝试继续执行等待队列中的任务
	defer func() {
		delete(wd.processingTasks, task.RunId)
		wd.tryNextAsync()
	}()

	// 查询运行实体，并级联更新状态
	if run, err := wd.workflowRunRepo.GetById(task.ctx, task.RunId); err != nil {
		if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			app.GetLogger().Error(fmt.Sprintf("failed to get workflow run #%s record", task.RunId), slog.Any("error", err))
		}
		return
	} else {
		workflowRun = run

		if run.Status == domain.WorkflowRunStatusTypePending {
			run.Status = domain.WorkflowRunStatusTypeProcessing
			wd.workflowRunRepo.SaveWithCascading(task.ctx, run)
		} else {
			// WTF? That should be impossible!
			return
		}
	}

	// 初始化工作流引擎
	logsBuf := make(domain.WorkflowLogs, 0)
	we := engine.NewWorkflowEngine()
	we.OnEnd(func(ctx context.Context) error {
		if errmsg := logsBuf.ErrorString(); errmsg == "" {
			workflowRun.Status = domain.WorkflowRunStatusTypeSucceeded
			workflowRun.EndedAt = time.Now()
		} else {
			workflowRun.Status = domain.WorkflowRunStatusTypeFailed
			workflowRun.EndedAt = time.Now()
			workflowRun.Error = errmsg
		}
		wd.workflowRunRepo.SaveWithCascading(task.ctx, workflowRun)
		return nil
	})
	we.OnError(func(ctx context.Context, err error) error {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			workflowRun.Status = domain.WorkflowRunStatusTypeCanceled
			wd.workflowRunRepo.SaveWithCascading(context.Background(), workflowRun)
		} else {
			workflowRun.Status = domain.WorkflowRunStatusTypeFailed
			workflowRun.EndedAt = time.Now()
			workflowRun.Error = err.Error()
			wd.workflowRunRepo.SaveWithCascading(task.ctx, workflowRun)
		}
		return nil
	})
	we.OnNodeError(func(ctx context.Context, node *engine.Node, err error) error {
		log := domain.WorkflowLog{}
		log.WorkflowId = task.WorkflowId
		log.RunId = task.RunId
		log.NodeId = node.Id
		log.NodeName = node.Data.Name
		log.Timestamp = time.Now().UnixMilli()
		log.Level = int32(slog.LevelError)
		log.Message = err.Error()
		log.CreatedAt = time.Now()
		if _, err := wd.workflowLogRepo.Save(ctx, &log); err != nil {
			app.GetLogger().Error(err.Error())
		}

		logsBuf = append(logsBuf, log)
		return nil
	})
	we.OnNodeLogging(func(ctx context.Context, node *engine.Node, record logging.Record) error {
		log := domain.WorkflowLog{}
		log.WorkflowId = task.WorkflowId
		log.RunId = task.RunId
		log.NodeId = node.Id
		log.NodeName = node.Data.Name
		log.Timestamp = record.Time.UnixMilli()
		log.Level = int32(record.Level)
		log.Message = record.Message
		log.Data = record.Data
		log.CreatedAt = record.Time
		if _, err := wd.workflowLogRepo.Save(ctx, &log); err != nil {
			app.GetLogger().Error(err.Error())
		}

		logsBuf = append(logsBuf, log)
		return nil
	})

	// 执行工作流
	app.GetLogger().Info(fmt.Sprintf("start to invoke workflow run #%s", task.RunId))
	we.Invoke(task.ctx, workflowRun.WorkflowId, workflowRun.Id, &workflowRun.Graph.WorkflowGraph)
}

func (wd *workflowDispatcher) tryNextAsync() {
	wd.taskMtx.RLock()

	for i, pendingRunId := range wd.pendingRunQueue {
		workflowRun, err := wd.workflowRunRepo.GetById(context.Background(), pendingRunId)
		if err != nil {
			app.GetLogger().Error(fmt.Sprintf("failed to get workflow run #%s record", pendingRunId), slog.Any("error", err))
			continue
		}

		var hasSameWorkflowTask bool // 相同 Workflow 的任务同一时间只能有一个 Run 在执行
		for _, task := range wd.processingTasks {
			if task.WorkflowId == workflowRun.WorkflowId {
				hasSameWorkflowTask = true
				break
			}
		}

		if !hasSameWorkflowTask && len(wd.processingTasks) < wd.concurrency {
			wd.taskMtx.RUnlock()
			wd.taskMtx.Lock()
			defer wd.taskMtx.Unlock()

			ctxRun, ctxCancel := context.WithCancel(context.Background())
			task := &taskInfo{WorkflowId: workflowRun.WorkflowId, RunId: workflowRun.Id, ctx: ctxRun, cancel: ctxCancel}
			wd.pendingRunQueue = append(wd.pendingRunQueue[:i], wd.pendingRunQueue[i+1:]...)
			wd.processingTasks[pendingRunId] = task
			go func() { wd.tryExecuteAsync(task) }()
			return
		}
	}

	wd.taskMtx.RUnlock()
}

func newWorkflowDispatcher() WorkflowDispatcher {
	return &workflowDispatcher{
		concurrency: maxWorkers,

		pendingRunQueue: make([]string, 0),
		processingTasks: make(map[string]*taskInfo),

		workflowRepo:    repository.NewWorkflowRepository(),
		workflowRunRepo: repository.NewWorkflowRunRepository(),
		workflowLogRepo: repository.NewWorkflowLogRepository(),
	}
}
