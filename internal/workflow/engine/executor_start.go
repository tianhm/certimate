package engine

import (
	"log/slog"
)

type startNodeExecutor struct {
	nodeExecutor
}

func (ne *startNodeExecutor) Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
	execRes := newNodeExecutionResult(execCtx.Node)

	ne.logger.Info("")

	return execRes, nil
}

func newStartNodeExecutor() NodeExecutor {
	return &startNodeExecutor{
		nodeExecutor: nodeExecutor{logger: slog.Default()},
	}
}
