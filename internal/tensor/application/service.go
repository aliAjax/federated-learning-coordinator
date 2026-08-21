package application

import (
	"context"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/tensor/domain"
)

type Service struct{}

func New() *Service { return &Service{} }
func (s *Service) Validate(ctx context.Context, schema domain.Schema) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if err := schema.Validate(); err != nil {
		return fmt.Errorf("tensor schema: %w", err)
	}
	if schema.Count() > 10_000_000 {
		return fmt.Errorf("tensor schema too large")
	}
	return nil
}
func (s *Service) Compatible(a, b domain.Schema) bool {
	if a.Version != b.Version || len(a.Layers) != len(b.Layers) {
		return false
	}
	for i := range a.Layers {
		if a.Layers[i].Name != b.Layers[i].Name || a.Layers[i].DType != b.Layers[i].DType {
			return false
		}
	}
	return true
}
