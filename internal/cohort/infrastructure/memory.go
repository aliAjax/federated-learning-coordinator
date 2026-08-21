package infrastructure

import (
	"context"
	"errors"
	"github.com/example/federated-learning-coordinator/internal/cohort/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items map[string]domain.Cohort
}

func New() *Memory { return &Memory{items: map[string]domain.Cohort{}} }
func (m *Memory) Save(_ context.Context, c domain.Cohort) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.items[c.ID]; ok {
		return errors.New("cohort exists")
	}
	m.items[c.ID] = c
	return nil
}
func (m *Memory) Find(_ context.Context, id string) (domain.Cohort, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.items[id]
	if !ok {
		return c, errors.New("cohort not found")
	}
	return c, nil
}
func (m *Memory) List(_ context.Context) []domain.Cohort {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.Cohort{}
	for _, c := range m.items {
		out = append(out, c)
	}
	return out
}
