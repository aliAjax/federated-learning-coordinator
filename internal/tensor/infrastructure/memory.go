package infrastructure

import (
	"context"
	"sync"
)

type Memory struct {
	mu    sync.RWMutex
	items map[string][]byte
}

func New() *Memory { return &Memory{items: map[string][]byte{}} }
func (m *Memory) Put(_ context.Context, k string, b []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[k] = append([]byte(nil), b...)
}
func (m *Memory) Get(_ context.Context, k string) ([]byte, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.items[k]
	return append([]byte(nil), b...), ok
}
func (m *Memory) Delete(_ context.Context, k string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, k)
}
