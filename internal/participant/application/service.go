package application

import (
	"context"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/participant/domain"
	"time"
)

type Repository interface {
	Save(context.Context, domain.Participant) error
	Find(context.Context, string) (domain.Participant, error)
	List(context.Context, string) []domain.Participant
}
type Service struct {
	repo Repository
	now  func() time.Time
}

func New(repo Repository) *Service { return &Service{repo: repo, now: time.Now} }
func (s *Service) Register(ctx context.Context, p domain.Participant) error {
	p.Capabilities = domain.NormalizeCapabilities(p.Capabilities)
	if err := p.Validate(); err != nil {
		return fmt.Errorf("validate participant: %w", err)
	}
	p.CreatedAt = s.now().UTC()
	return s.repo.Save(ctx, p)
}
func (s *Service) Pause(ctx context.Context, id string) error {
	p, err := s.repo.Find(ctx, id)
	if err != nil {
		return err
	}
	p.Status = domain.Paused
	return s.repo.Save(ctx, p)
}
func (s *Service) Resume(ctx context.Context, id string) error {
	p, err := s.repo.Find(ctx, id)
	if err != nil {
		return err
	}
	p.Status = domain.Active
	return s.repo.Save(ctx, p)
}
func (s *Service) Remove(ctx context.Context, id string) error {
	p, err := s.repo.Find(ctx, id)
	if err != nil {
		return err
	}
	p.Status = domain.Removed
	return s.repo.Save(ctx, p)
}
func (s *Service) Heartbeat(ctx context.Context, id string) error {
	p, err := s.repo.Find(ctx, id)
	if err != nil {
		return err
	}
	now := s.now().UTC()
	p.Touch(now)
	return s.repo.Save(ctx, p)
}
