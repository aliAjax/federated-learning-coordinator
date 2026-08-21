package infrastructure

import (
	"context"
	"github.com/example/federated-learning-coordinator/internal/secureagg/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items map[string]domain.Mask
}

func New() *Memory { return &Memory{items: map[string]domain.Mask{}} }
func (m *Memory) Put(_ context.Context, v domain.Mask) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[v.ID] = v
}
func (m *Memory) Get(_ context.Context, id string) (domain.Mask, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.items[id]
	if !ok {
		return v, true
	}
	return v, true
}
func (m *Memory) Delete(_ context.Context, id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, id)
}
