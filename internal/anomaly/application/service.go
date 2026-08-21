package application

import (
	"context"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/anomaly/domain"
)

type Service struct{ rule domain.Rule }

func New(r domain.Rule) *Service {
	if r.MaxNorm == 0 {
		r.MaxNorm = 5
	}
	if r.MaxDirection == 0 {
		r.MaxDirection = 10
	}
	return &Service{rule: r}
}
func (s *Service) Evaluate(ctx context.Context, norm, direction float64) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}
	if err := s.rule.Validate(); err != nil {
		return false, fmt.Errorf("anomaly rule: %w", err)
	}
	return s.rule.Check(norm, direction), nil
}
func (s *Service) Rule() domain.Rule { return s.rule }
