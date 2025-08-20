package engine

import (
	"context"
	"slices"
	"sync"
)

type WorkflowContext struct {
	WorkflowId string
	RunId      string
	RunGraph   *Graph

	engine    WorkflowEngine
	variables WorkflowIOsManager
	inputs    WorkflowIOsManager

	ctx context.Context
}

func (c *WorkflowContext) SetExecutingWorkflow(workflowId string, runId string, runGraph *Graph) *WorkflowContext {
	c.WorkflowId = workflowId
	c.RunId = runId
	c.RunGraph = runGraph
	return c
}

func (c *WorkflowContext) SetEngine(engine WorkflowEngine) *WorkflowContext {
	c.engine = engine
	return c
}

func (c *WorkflowContext) SetVariablesManager(variables WorkflowIOsManager) *WorkflowContext {
	c.variables = variables
	return c
}

func (c *WorkflowContext) SetInputsManager(inputs WorkflowIOsManager) *WorkflowContext {
	c.inputs = inputs
	return c
}

func (c *WorkflowContext) SetContext(ctx context.Context) *WorkflowContext {
	c.ctx = ctx
	return c
}

func (c *WorkflowContext) Clone() *WorkflowContext {
	return &WorkflowContext{
		WorkflowId: c.WorkflowId,
		RunId:      c.RunId,
		RunGraph:   c.RunGraph,

		engine:    c.engine,
		variables: c.variables,
		inputs:    c.inputs,

		ctx: c.ctx,
	}
}

type WorkflowIOsManager interface {
	All() []NodeIOEntry
	Erase()

	Add(entry NodeIOEntry)
	Get(name string) (*NodeIOEntry, bool)
	GetScoped(scope string, name string) (*NodeIOEntry, bool)
	Take(name string) (*NodeIOEntry, bool)
	TakeScoped(scope string, name string) (*NodeIOEntry, bool)
	Remove(name string) bool
	RemoveScoped(scope string, name string) bool
}

type workflowIOsManager struct {
	entriesMtx sync.RWMutex
	entries    []NodeIOEntry
}

var _ WorkflowIOsManager = (*workflowIOsManager)(nil)

func (m *workflowIOsManager) All() []NodeIOEntry {
	m.entriesMtx.RLock()
	defer m.entriesMtx.RUnlock()

	if m.entries == nil {
		return make([]NodeIOEntry, 0)
	}

	return slices.Clone(m.entries)
}

func (m *workflowIOsManager) Erase() {
	m.entriesMtx.Lock()
	defer m.entriesMtx.Unlock()

	m.entries = make([]NodeIOEntry, 0)
}

func (m *workflowIOsManager) Add(entry NodeIOEntry) {
	m.entriesMtx.Lock()
	defer m.entriesMtx.Unlock()

	if m.entries == nil {
		m.entries = make([]NodeIOEntry, 0)
	}

	for i, item := range m.entries {
		if item.Scope == entry.Scope && item.Key == entry.Key {
			m.entries[i] = entry
			return
		}
	}
	m.entries = append(m.entries, entry)
}

func (m *workflowIOsManager) Get(name string) (*NodeIOEntry, bool) {
	return m.GetScoped("", name)
}

func (m *workflowIOsManager) GetScoped(scope string, name string) (*NodeIOEntry, bool) {
	m.entriesMtx.RLock()
	defer m.entriesMtx.RUnlock()

	if m.entries == nil {
		return nil, false
	}

	for _, item := range m.entries {
		if item.Scope == scope && item.Key == name {
			return &item, true
		}
	}
	return nil, false
}

func (m *workflowIOsManager) Take(name string) (*NodeIOEntry, bool) {
	return m.TakeScoped("", name)
}

func (m *workflowIOsManager) TakeScoped(scope, name string) (*NodeIOEntry, bool) {
	m.entriesMtx.Lock()
	defer m.entriesMtx.Unlock()

	if m.entries == nil {
		return nil, false
	}

	for i, item := range m.entries {
		if item.Scope == scope && item.Key == name {
			m.entries = slices.Delete(m.entries, i, i+1)
			return &item, true
		}
	}
	return nil, false
}

func (m *workflowIOsManager) Remove(name string) bool {
	return m.RemoveScoped("", name)
}

func (m *workflowIOsManager) RemoveScoped(scope string, name string) bool {
	_, ok := m.TakeScoped(scope, name)
	return ok
}

func newWorkflowIOsManager() WorkflowIOsManager {
	return &workflowIOsManager{
		entries: make([]NodeIOEntry, 0),
	}
}

const (
	wfIOTypeCertificate = "certificate"
)

const (
	wfIOKeyCertificate = "certificate"
)

const (
	wfVariableKeyCertificateValidity = "certificate.validity" // ValueType: "bool"
	wfVariableKeyCertificateDaysLeft = "certificate.daysLeft" // ValueType: "int32"
	wfVariableKeyNodeSkipped         = "node.skipped"         // ValueType: "bool"
)
