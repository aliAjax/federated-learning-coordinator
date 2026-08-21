package application

import (
	"context"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/privacy/domain"
	"time"
)

type Repository interface {
	Save(context.Context, domain.Budget) error
	Find(context.Context, string) (domain.Budget, error)
}
type Service struct {
	repo Repository
	now  func() time.Time
}

func New(r Repository) *Service { return &Service{repo: r, now: time.Now} }
func (s *Service) Consume(ctx context.Context, id string, cost float64) (domain.Budget, error) {
	b, e := s.repo.Find(ctx, id)
	if e != nil {
		return b, e
	}
	if e = b.Consume(cost); e != nil {
		return b, fmt.Errorf("consume privacy budget: %w", e)
	}
	b.UpdatedAt = s.now().UTC()
	return b, s.repo.Save(ctx, b)
}
func (s *Service) Remaining(ctx context.Context, id string) (float64, error) {
	b, e := s.repo.Find(ctx, id)
	if e != nil {
		return 0, e
	}
	return b.Remaining(), nil
}
func (s *Service) Validate(ctx context.Context, b domain.Budget) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return b.Validate()
}
