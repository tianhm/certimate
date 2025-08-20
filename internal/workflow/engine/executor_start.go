package engine

import (
	"log/slog"
)

type startNodeExecutor struct {
	nodeExecutor
}

func (e *startNodeExecutor) Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
	execRes := &NodeExecutionResult{}

	return execRes, nil
}

func newStartNodeExecutor() NodeExecutor {
	return &startNodeExecutor{
		nodeExecutor: nodeExecutor{logger: slog.Default()},
	}
}
