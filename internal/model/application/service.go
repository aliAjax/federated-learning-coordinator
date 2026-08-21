package application

import (
	"context"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/model/domain"
	"time"
)

type Repository interface {
	Save(context.Context, domain.Model) error
	Find(context.Context, string) (domain.Model, error)
}
type Service struct {
	repo Repository
	now  func() time.Time
}

func New(r Repository) *Service { return &Service{repo: r, now: time.Now} }
func (s *Service) Publish(ctx context.Context, id string) (domain.Model, error) {
	m, e := s.repo.Find(ctx, id)
	if e != nil {
		return m, e
	}
	if e = m.Publish(s.now().UTC()); e != nil {
		return m, fmt.Errorf("publish model: %w", e)
	}
	return m, s.repo.Save(ctx, m)
}
func (s *Service) Revoke(ctx context.Context, id string) (domain.Model, error) {
	m, e := s.repo.Find(ctx, id)
	if e != nil {
		return m, e
	}
	m.Revoke()
	return m, s.repo.Save(ctx, m)
}
func (s *Service) Validate(ctx context.Context, m domain.Model) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return m.Validate()
}
