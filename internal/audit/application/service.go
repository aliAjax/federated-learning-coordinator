package application

import (
	"context"
	"github.com/example/federated-learning-coordinator/internal/audit/domain"
	"time"
)

type Repository interface {
	Append(context.Context, domain.Event) error
	List(context.Context, string) []domain.Event
	LastHash(context.Context) string
}
type Service struct {
	repo Repository
	now  func() time.Time
}

func New(r Repository) *Service { return &Service{repo: r, now: time.Now} }
func (s *Service) Record(ctx context.Context, typ, actor, subject string, payload map[string]any) error {
	e := domain.Event{Type: typ, Actor: actor, Subject: subject, Payload: payload, CreatedAt: s.now().UTC()}
	e.Seal(s.repo.LastHash(ctx))
	return s.repo.Append(ctx, e)
}
func (s *Service) Verify(ctx context.Context) bool {
	events := s.repo.List(ctx, "")
	prev := ""
	for _, e := range events {
		if e.PreviousHash != prev || !e.Verify() {
			return false
		}
		prev = e.Hash
	}
	return true
}
