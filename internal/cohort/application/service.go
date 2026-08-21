package application

import (
	"context"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/cohort/domain"
	"time"
)

type Repository interface {
	Save(context.Context, domain.Cohort) error
	Find(context.Context, string) (domain.Cohort, error)
	List(context.Context) []domain.Cohort
}
type Service struct {
	repo Repository
	now  func() time.Time
}

func New(r Repository) *Service { return &Service{repo: r, now: time.Now} }
func (s *Service) Create(ctx context.Context, c domain.Cohort) error {
	if c.Status == "" {
		c.Status = domain.Enabled
	}
	if err := c.Validate(); err != nil {
		return fmt.Errorf("cohort validation: %w", err)
	}
	c.CreatedAt = s.now().UTC()
	return s.repo.Save(ctx, c)
}
func (s *Service) Disable(ctx context.Context, id string) error {
	c, e := s.repo.Find(ctx, id)
	if e != nil {
		return e
	}
	c.Disable()
	return s.repo.Save(ctx, c)
}
func (s *Service) Enable(ctx context.Context, id string) error {
	c, e := s.repo.Find(ctx, id)
	if e != nil {
		return e
	}
	c.Enable()
	return s.repo.Save(ctx, c)
}
