package infrastructure

import (
	"context"
	"github.com/example/federated-learning-coordinator/internal/anomaly/domain"
	"sync"
)

type Memory struct {
	mu   sync.RWMutex
	rule domain.Rule
}

func New(r domain.Rule) *Memory                         { return &Memory{rule: r} }
func (m *Memory) Save(_ context.Context, r domain.Rule) { m.mu.Lock(); defer m.mu.Unlock(); m.rule = r }
func (m *Memory) Get(_ context.Context) domain.Rule {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.rule
}
