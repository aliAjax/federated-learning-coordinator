package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/model"
)

var ErrUnavailable = errors.New("postgres adapter unavailable in local mode")

type Repository struct{ DSN string }

func New(dsn string) *Repository { return &Repository{DSN: dsn} }
func (r *Repository) Ping(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if r.DSN == "" {
		return fmt.Errorf("postgres ping: %v", ErrUnavailable)
	}
	return nil
}
func (r *Repository) SaveModel(context.Context, *model.Model) error { return fmt.Errorf("postgres save: %v", ErrUnavailable) }
func (r *Repository) LoadModel(context.Context, string) (*model.Model, error) {
	return nil, fmt.Errorf("postgres load: %v", ErrUnavailable)
}
func (r *Repository) Migrate(context.Context) error { return fmt.Errorf("postgres migrate: %v", ErrUnavailable) }
