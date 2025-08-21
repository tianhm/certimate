package engine

import (
	"slices"
	"sync"
)

type StateEntry struct {
	Scope     string // 零值时表示全局的，否则表示指定节点的
	Type      string // 仅表示输入输出有值，表示变量无值
	Key       string
	Value     any
	ValueType string `options:"string | number | boolean"`
}

type StateManager interface {
	All() []StateEntry
	Erase()

	Add(entry StateEntry)
	Get(name string) (*StateEntry, bool)
	GetScoped(scope string, name string) (*StateEntry, bool)
	Take(name string) (*StateEntry, bool)
	TakeScoped(scope string, name string) (*StateEntry, bool)
	Remove(name string) bool
	RemoveScoped(scope string, name string) bool
}

type stateManager struct {
	entriesMtx sync.RWMutex
	entries    []StateEntry
}

var _ StateManager = (*stateManager)(nil)

func (m *stateManager) All() []StateEntry {
	m.entriesMtx.RLock()
	defer m.entriesMtx.RUnlock()

	if m.entries == nil {
		return make([]StateEntry, 0)
	}

	return slices.Clone(m.entries)
}

func (m *stateManager) Erase() {
	m.entriesMtx.Lock()
	defer m.entriesMtx.Unlock()

	m.entries = make([]StateEntry, 0)
}

func (m *stateManager) Add(entry StateEntry) {
	m.entriesMtx.Lock()
	defer m.entriesMtx.Unlock()

	if m.entries == nil {
		m.entries = make([]StateEntry, 0)
	}

	for i, item := range m.entries {
		if item.Scope == entry.Scope && item.Key == entry.Key {
			m.entries[i] = entry
			return
		}
	}
	m.entries = append(m.entries, entry)
}

func (m *stateManager) Get(name string) (*StateEntry, bool) {
	return m.GetScoped("", name)
}

func (m *stateManager) GetScoped(scope string, name string) (*StateEntry, bool) {
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

func (m *stateManager) Take(name string) (*StateEntry, bool) {
	return m.TakeScoped("", name)
}

func (m *stateManager) TakeScoped(scope, name string) (*StateEntry, bool) {
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

func (m *stateManager) Remove(name string) bool {
	return m.RemoveScoped("", name)
}

func (m *stateManager) RemoveScoped(scope string, name string) bool {
	_, ok := m.TakeScoped(scope, name)
	return ok
}

func newStateManager() StateManager {
	return &stateManager{
		entries: make([]StateEntry, 0),
	}
}

const (
	stateIOTypeCertificate = "certificate"
)

const (
	stateIOKeyCertificate = "certificate"
)

const (
	stateVarKeyNodeSkipped         = "node.skipped"         // ValueType: "boolean"
	stateVarKeyCertificateValidity = "certificate.validity" // ValueType: "boolean"
	stateVarKeyCertificateDaysLeft = "certificate.daysLeft" // ValueType: "number"
)
