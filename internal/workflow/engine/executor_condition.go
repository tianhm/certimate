package engine

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/samber/lo"
)

type conditionNodeExecutor struct {
	nodeExecutor
}

func (ne *conditionNodeExecutor) Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
	var engine *workflowEngine
	if we, ok := execCtx.engine.(*workflowEngine); !ok {
		panic("impossible!")
	} else {
		engine = we
	}

	execRes := &NodeExecutionResult{}

	errs := make([]error, 0)
	blocks := lo.Filter(execCtx.Node.Blocks, func(n *Node, _ int) bool { return n.Type == NodeTypeBranchBlock })
	for _, node := range blocks {
		select {
		case <-execCtx.ctx.Done():
			return execRes, execCtx.ctx.Err()
		default:
		}

		err := engine.executeNode(execCtx.Clone(), node)
		if err != nil {
			if errors.Is(err, errInterrupted) {
				return execRes, err
			}
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return execRes, fmt.Errorf("error occurred when executing child nodes: %w", errors.Join(errs...))
	}

	return execRes, nil
}

func newConditionNodeExecutor() NodeExecutor {
	return &conditionNodeExecutor{
		nodeExecutor: nodeExecutor{logger: slog.Default()},
	}
}

type branchBlockNodeExecutor struct {
	nodeExecutor
}

func (ne *branchBlockNodeExecutor) Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error) {
	execRes := &NodeExecutionResult{}

	nodeCfg := execCtx.Node.Data.Config.AsBranchBlock()
	if nodeCfg.Expression == nil {
		ne.logger.Info("enter this branch without any conditions")
	} else {
		variables := lo.Reduce(execCtx.variables.All(), func(acc map[string]map[string]any, entry NodeIOEntry, _ int) map[string]map[string]any {
			if _, ok := acc[entry.Scope]; !ok {
				acc[entry.Scope] = make(map[string]any)
			}

			// 这里需要把所有值都转换为字符串形式，因为 Expression.Eval 仅支持字符串类型的值
			var value any
			if entry.ValueType == "number" {
				value = fmt.Sprintf("%d", entry.Value)
			} else if entry.ValueType == "boolean" {
				value = strconv.FormatBool(entry.Value.(bool))
			} else {
				value = entry.Value
			}
			acc[entry.Scope][entry.Key] = value

			return acc
		}, make(map[string]map[string]any))

		rs, err := nodeCfg.Expression.Eval(variables)
		if err != nil {
			ne.logger.Warn(fmt.Sprintf("failed to eval expr: %+v", err))
			return execRes, err
		}

		if rs.Value == false {
			ne.logger.Info("skip this branch, because condition not met")
			return execRes, nil
		} else {
			ne.logger.Info("enter this branch, because condition met")
		}
	}

	if engine, ok := execCtx.engine.(*workflowEngine); !ok {
		panic("impossible!")
	} else {
		if err := engine.executeBlocks(execCtx.Clone(), execCtx.Node.Blocks); err != nil {
			return execRes, err
		}
	}

	return execRes, nil
}

func newBranchBlockNodeExecutor() NodeExecutor {
	return &branchBlockNodeExecutor{
		nodeExecutor: nodeExecutor{logger: slog.Default()},
	}
}
