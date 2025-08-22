package engine

import (
	"log/slog"
)

type endNodeExecutor struct {
	nodeExecutor
}

func (e *endNodeExecutor) Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
	execRes := newNodeExecutionResult(execCtx.Node)
	execRes.Interrupted = true

	return execRes, nil
}

func newEndNodeExecutor() NodeExecutor {
	return &endNodeExecutor{
		nodeExecutor: nodeExecutor{logger: slog.Default()},
	}
}
