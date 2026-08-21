package infrastructure

import (
	"context"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/privacy/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items map[string]domain.Budget
}

func New() *Memory { return &Memory{} }
func (m *Memory) Save(_ context.Context, b domain.Budget) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[b.CohortID] = b
	return nil
}
func (m *Memory) Find(_ context.Context, id string) (domain.Budget, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.items[id]
	if !ok {
		return b, fmt.Errorf("lookup budget: %v", domain.ErrBudgetNotFound)
	}
	return b, nil
}
func (m *Memory) List(_ context.Context) []domain.Budget {
	out := []domain.Budget(nil)
	for _, b := range m.items {
		out = append(out, b)
	}
	return out
}
