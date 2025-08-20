package engine

import (
	"context"
	"log/slog"
	"sync"
)

type NodeExecutor interface {
	withLogger

	Execute(execCtx *NodeExecutionContext) (*NodeExecutionResult, error)
}

type nodeExecutor struct {
	logger *slog.Logger
}

func (e *nodeExecutor) SetLogger(logger *slog.Logger) {
	e.logger = logger
}

type NodeExecutionContext struct {
	WorkflowContext

	Node *Node
}

func (c *NodeExecutionContext) SetExecutingWorkflow(workflowId string, runId string, runGraph *Graph) *NodeExecutionContext {
	c.WorkflowContext.SetExecutingWorkflow(workflowId, runId, runGraph)
	return c
}

func (c *NodeExecutionContext) SetExecutingNode(node *Node) *NodeExecutionContext {
	c.Node = node
	return c
}

func (c *NodeExecutionContext) SetEngine(engine WorkflowEngine) *NodeExecutionContext {
	c.WorkflowContext.SetEngine(engine)
	return c
}

func (c *NodeExecutionContext) SetVariablesManager(variablesMgr WorkflowIOsManager) *NodeExecutionContext {
	c.variables = variablesMgr
	return c
}

func (c *NodeExecutionContext) SetInputsManager(inputsMgr WorkflowIOsManager) *NodeExecutionContext {
	c.inputs = inputsMgr
	return c
}

func (c *NodeExecutionContext) SetContext(ctx context.Context) *NodeExecutionContext {
	c.WorkflowContext.SetContext(ctx)
	return c
}

type NodeExecutionResult struct {
	Interrupted bool // 是否中断执行（通常由 End 节点主动触发）

	variablesMtx sync.RWMutex
	Variables    []NodeIOEntry
	outputsMtx   sync.RWMutex
	Outputs      []NodeIOEntry
}

func (r *NodeExecutionResult) AddVariable(scope string, key string, value any, valueType string) {
	r.AddVariableEntry(NodeIOEntry{
		Scope:     scope,
		Type:      "",
		Key:       key,
		Value:     value,
		ValueType: valueType,
	})
}

func (r *NodeExecutionResult) AddVariableEntry(entry NodeIOEntry) {
	r.variablesMtx.Lock()
	defer r.variablesMtx.Unlock()

	if r.Variables == nil {
		r.Variables = make([]NodeIOEntry, 0)
	}

	for i, item := range r.Variables {
		if item.Scope == entry.Scope && item.Key == entry.Key {
			r.Variables[i] = entry
			return
		}
	}
	r.Variables = append(r.Variables, entry)
}

func (r *NodeExecutionResult) AddOutput(scope string, type_ string, key string, value any, valueType string) {
	r.AddOutputEntry(NodeIOEntry{
		Scope:     scope,
		Type:      type_,
		Key:       key,
		Value:     value,
		ValueType: valueType,
	})
}

func (r *NodeExecutionResult) AddOutputEntry(entry NodeIOEntry) {
	r.outputsMtx.Lock()
	defer r.outputsMtx.Unlock()

	if r.Outputs == nil {
		r.Outputs = make([]NodeIOEntry, 0)
	}

	for i, t := range r.Outputs {
		if t.Scope == entry.Scope && t.Key == entry.Key {
			r.Outputs[i] = entry
			return
		}
	}
	r.Outputs = append(r.Outputs, entry)
}
