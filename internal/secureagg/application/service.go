package application

import (
	"context"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/secureagg/domain"
	"sync"
	"time"
)

type Service struct {
	mu    sync.RWMutex
	masks map[string]domain.Mask
}

func New() *Service { return &Service{masks: map[string]domain.Mask{}} }
func (s *Service) Register(ctx context.Context, m domain.Mask) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if err := m.Validate(); err != nil {
		return fmt.Errorf("mask: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.masks[m.ID]; ok {
		return fmt.Errorf("mask already registered")
	}
	s.masks[m.ID] = m
	return nil
}
func (s *Service) Reveal(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.masks[id]
	if !ok {
		return fmt.Errorf("mask not found")
	}
	m.Reveal()
	s.masks[id] = m
	return nil
}
func (s *Service) List(_ context.Context, round string) []domain.Mask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []domain.Mask{}
	for _, m := range s.masks {
		if round == "" || m.RoundID == round {
			out = append(out, m)
		}
	}
	return out
}
func (s *Service) Expired(m domain.Mask, now time.Time) bool {
	return now.Sub(m.CreatedAt) > time.Hour && !m.Revealed
}
