package infrastructure

import (
	"context"
	"errors"
	"github.com/example/federated-learning-coordinator/internal/participant/domain"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items map[string]domain.Participant
}

func New() *Memory { return &Memory{items: map[string]domain.Participant{}} }
func (m *Memory) Save(_ context.Context, p domain.Participant) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.ID == "" {
		return errors.New("empty participant id")
	}
	m.items[p.ID] = p
	return nil
}
func (m *Memory) Find(_ context.Context, id string) (domain.Participant, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.items[id]
	if !ok {
		return p, errors.New("participant not found")
	}
	return p, nil
}
func (m *Memory) List(_ context.Context, cid string) []domain.Participant {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.Participant{}
	for _, p := range m.items {
		if cid == "" || p.CohortID == cid {
			out = append(out, p)
		}
	}
	return out
}
