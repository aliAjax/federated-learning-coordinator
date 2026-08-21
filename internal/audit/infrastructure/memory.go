package infrastructure

import (
	"context"
	"github.com/example/federated-learning-coordinator/internal/audit/domain"
	"sync"
)

type Memory struct {
	mu     sync.RWMutex
	events []domain.Event
}

func New() *Memory { return &Memory{events: []domain.Event{}} }
func (m *Memory) Append(_ context.Context, e domain.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, e)
	return nil
}
func (m *Memory) List(_ context.Context, _ string) []domain.Event {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.Event(nil), m.events...)
}
func (m *Memory) LastHash(_ context.Context) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.events) == 0 {
		return ""
	}
	return m.events[len(m.events)-1].Hash
}
