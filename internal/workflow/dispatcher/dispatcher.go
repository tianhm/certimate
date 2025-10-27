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

	"github.com/samber/lo"

	"github.com/certimate-go/certimate/internal/app"
	"github.com/certimate-go/certimate/internal/domain"
	"github.com/certimate-go/certimate/internal/repository"
	"github.com/certimate-go/certimate/internal/workflow/engine"
	"github.com/certimate-go/certimate/pkg/logging"
)

var maxWorkers = 1

func init() {
	envMaxWorkers := os.Getenv("CERTIMATE_WORKFLOW_MAX_WORKERS")
	if n, _ := strconv.Atoi(envMaxWorkers); n > 0 {
		maxWorkers = n
	} else {
		maxWorkers = runtime.GOMAXPROCS(0)
		if maxWorkers == 0 {
			maxWorkers = max(1, runtime.NumCPU())
		}
	}
}

type WorkflowDispatcher interface {
	GetStatistics() Statistics

	Bootup(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Start(ctx context.Context, runId string) error
	Cancel(ctx context.Context, runId string) error
}

type Statistics struct {
	Concurrency      int
	PendingRunIds    []string
	ProcessingRunIds []string
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

	syslog *slog.Logger
}

var _ WorkflowDispatcher = (*workflowDispatcher)(nil)

func (wd *workflowDispatcher) GetStatistics() Statistics {
	wd.taskMtx.RLock()
	defer wd.taskMtx.RUnlock()

	stats := Statistics{
		Concurrency:      wd.concurrency,
		PendingRunIds:    make([]string, 0),
		ProcessingRunIds: make([]string, 0),
	}
	for _, pendingRunId := range wd.pendingRunQueue {
		stats.PendingRunIds = append(stats.PendingRunIds, pendingRunId)
	}
	for _, processingRunId := range wd.processingTasks {
		stats.ProcessingRunIds = append(stats.ProcessingRunIds, processingRunId.RunId)
	}

	return stats
}

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
	wd.taskMtx.RLock()
	if _, exists := wd.processingTasks[runId]; exists {
		wd.taskMtx.RUnlock()
		return fmt.Errorf("workflow run %s is already processing", runId)
	}
	for _, pendingRunId := range wd.pendingRunQueue {
		if pendingRunId == runId {
			wd.taskMtx.RUnlock()
			return fmt.Errorf("workflow run %s is already in the queue", runId)
		}
	}
	wd.taskMtx.RUnlock()

	wd.taskMtx.Lock()
	wd.pendingRunQueue = append(wd.pendingRunQueue, runId)
	wd.taskMtx.Unlock()

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
		return fmt.Errorf("workflow run #%s is already completed", workflowRun.Id)
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

		wd.syslog.Info(fmt.Sprintf("workflow run #%s was canceled", task.RunId))
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
	var workflow *domain.Workflow
	var workflowRun *domain.WorkflowRun
	var err error

	// 捕获 panic
	defer func() {
		if r := recover(); r != nil {
			wd.syslog.Error(fmt.Sprintf("workflow dispatcher panic: %v", r), slog.String("workflowId", task.WorkflowId), slog.String("runId", task.RunId))
			slog.Default().Error(fmt.Sprintf("workflow dispatcher panic: %v, stack trace: %s", r, string(debug.Stack())), slog.String("workflowId", task.WorkflowId), slog.String("runId", task.RunId))

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
		wd.taskMtx.Lock()
		delete(wd.processingTasks, task.RunId)
		wd.taskMtx.Unlock()

		go func() { wd.tryNextAsync() }()
	}()

	// 查询运行实体，并级联更新状态
	if workflowRun, err = wd.workflowRunRepo.GetById(task.ctx, task.RunId); err != nil {
		if !(errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			wd.syslog.Error(fmt.Sprintf("failed to get workflow run #%s record", task.RunId), slog.Any("error", err))
		}
		return
	} else {
		if workflowRun.Status == domain.WorkflowRunStatusTypePending {
			workflowRun.Status = domain.WorkflowRunStatusTypeProcessing
			wd.workflowRunRepo.SaveWithCascading(task.ctx, workflowRun)
		} else {
			// WTF? That should be impossible!
			return
		}
	}

	// 查询工作流实体
	workflow, err = wd.workflowRepo.GetById(task.ctx, workflowRun.WorkflowId)
	if err != nil {
		wd.syslog.Error(fmt.Sprintf("failed to get workflow #%s record", workflowRun.WorkflowId), slog.Any("error", err))
		return
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
		if errors.Is(err, engine.ErrTerminated) || errors.Is(err, engine.ErrBlocksException) {
			return nil
		}

		log := domain.WorkflowLog{}
		log.WorkflowId = task.WorkflowId
		log.RunId = task.RunId
		log.NodeId = node.Id
		log.NodeName = node.Data.Name
		log.Timestamp = time.Now().UnixMilli()
		log.Level = int32(slog.LevelError)
		log.Message = err.Error()
		log.CreatedAt = time.Now()
		logsBuf = append(logsBuf, log)

		if _, err := wd.workflowLogRepo.Save(ctx, &log); err != nil {
			wd.syslog.Error(err.Error())
		}

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
		log.Data = record.Data()
		log.CreatedAt = time.Now()
		logsBuf = append(logsBuf, log)

		if _, err := wd.workflowLogRepo.Save(ctx, &log); err != nil {
			wd.syslog.Error(err.Error())
		}

		return nil
	})

	// 执行工作流
	wd.syslog.Info(fmt.Sprintf("workflow run #%s (work#%s) started", task.RunId, task.WorkflowId))
	we.Invoke(task.ctx, engine.WorkflowExecution{
		WorkflowId:   workflowRun.WorkflowId,
		WorkflowName: workflow.Name,
		RunId:        workflowRun.Id,
		RunTrigger:   workflowRun.Trigger,
		Graph:        workflowRun.Graph,
	})
	wd.syslog.Info(fmt.Sprintf("workflow run #%s (work#%s) stopped", task.RunId, task.WorkflowId))
}

func (wd *workflowDispatcher) tryNextAsync() {
	wd.taskMtx.RLock()

	for _, pendingRunId := range wd.pendingRunQueue {
		workflowRun, err := wd.workflowRunRepo.GetById(context.Background(), pendingRunId)
		if err != nil {
			wd.syslog.Error(fmt.Sprintf("failed to get workflow run #%s record", pendingRunId), slog.Any("error", err))
			continue
		}

		var hasSameWorkflowTask bool // 相同 Workflow 的任务同一时间只能有一个 Run 在执行
		for _, processingTask := range wd.processingTasks {
			if processingTask.WorkflowId == workflowRun.WorkflowId {
				hasSameWorkflowTask = true
				break
			}
		}

		if hasSameWorkflowTask {
			wd.syslog.Warn(fmt.Sprintf("workflow run #%s is pending, because tasks that belonging to the same workflow #%s already exists", workflowRun.Id, workflowRun.WorkflowId))
		} else if len(wd.processingTasks) >= wd.concurrency && wd.concurrency > 0 {
			wd.syslog.Warn(fmt.Sprintf("workflow run #%s is pending, because the maximum concurrency (limit: %d) has been reached", pendingRunId, wd.concurrency))
		} else {
			wd.taskMtx.RUnlock()

			wd.taskMtx.Lock()
			ctxRun, ctxCancel := context.WithCancel(context.Background())
			task := &taskInfo{WorkflowId: workflowRun.WorkflowId, RunId: workflowRun.Id, ctx: ctxRun, cancel: ctxCancel}
			wd.pendingRunQueue = lo.Filter(wd.pendingRunQueue, func(s string, _ int) bool { return s != pendingRunId })
			wd.processingTasks[pendingRunId] = task
			wd.syslog.Info(fmt.Sprintf("workflow run #%s (work#%s) is being dispatched ...", task.RunId, task.WorkflowId))
			wd.taskMtx.Unlock()

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

		syslog: app.GetLogger(),
	}
}
