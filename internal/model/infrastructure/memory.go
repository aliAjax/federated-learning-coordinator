package infrastructure

import (
	"context"
	"errors"
	"github.com/example/federated-learning-coordinator/internal/model/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items map[string]domain.Model
}

func New() *Memory { return &Memory{items: map[string]domain.Model{}} }
func (m *Memory) Save(_ context.Context, v domain.Model) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[v.ID] = v
	return nil
}
func (m *Memory) Find(_ context.Context, id string) (domain.Model, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.items[id]
	if !ok {
		return v, errors.New("model not found")
	}
	return v, nil
}
func (m *Memory) List(_ context.Context, cid string) []domain.Model {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.Model{}
	for _, v := range m.items {
		if cid == "" || v.CohortID == cid {
			out = append(out, v)
		}
	}
	return out
}
