package postgres

import (
	"context"
	"errors"
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
		return ErrUnavailable
	}
	return nil
}
func (r *Repository) SaveModel(context.Context, *model.Model) error { return ErrUnavailable }
func (r *Repository) LoadModel(context.Context, string) (*model.Model, error) {
	return nil, ErrUnavailable
}
func (r *Repository) Migrate(context.Context) error { return ErrUnavailable }
