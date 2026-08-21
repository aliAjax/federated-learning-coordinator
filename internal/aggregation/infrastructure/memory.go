package infrastructure

import (
	"context"
	"github.com/example/federated-learning-coordinator/internal/aggregation/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items map[string]domain.Result
}

func New() *Memory { return &Memory{items: map[string]domain.Result{}} }
func (m *Memory) Save(_ context.Context, r domain.Result) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[r.RoundID] = r
}
func (m *Memory) Get(_ context.Context, id string) (domain.Result, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.items[id]
	return r, ok
}
func (m *Memory) List(_ context.Context) []domain.Result {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.Result{}
	for _, r := range m.items {
		out = append(out, r)
	}
	return out
}
