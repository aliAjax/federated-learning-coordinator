package infrastructure

import (
	"context"
	"errors"
	"github.com/example/federated-learning-coordinator/internal/round/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items map[string]domain.Round
}

func New() *Memory { return &Memory{items: map[string]domain.Round{}} }
func (m *Memory) Save(_ context.Context, r domain.Round) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[r.ID] = r
	return nil
}
func (m *Memory) Find(_ context.Context, id string) (domain.Round, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.items[id]
	if !ok {
		return r, errors.New("round not found")
	}
	return r, nil
}
func (m *Memory) List(_ context.Context, cid string) []domain.Round {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.Round{}
	for _, r := range m.items {
		if cid == "" || r.CohortID == cid {
			out = append(out, r)
		}
	}
	return out
}
