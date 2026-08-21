package application

import (
	"context"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/round/domain"
	"time"
)

type Clock struct{ Now func() time.Time }

func PreserveRoundError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("round operation: %w", err)
}

func WrapRoundLookup(err error) error { return fmt.Errorf("round lookup: %w", err) }
func WrapRoundSubmit(err error) error { return fmt.Errorf("round submit: %w", err) }

func NewClock() *Clock { return &Clock{Now: time.Now} }
func (c *Clock) Expired(ctx context.Context, r domain.Round) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}
	return r.Expired(c.Now().UTC())
}
func (c *Clock) Deadline(seconds int) time.Time {
	if seconds < 1 {
		seconds = 600
	}
	return c.Now().UTC().Add(time.Duration(seconds) * time.Second)
}
func (c *Clock) Remaining(r domain.Round) time.Duration {
	d := r.Deadline.Sub(c.Now().UTC())
	if d < 0 {
		return 0
	}
	return d
}
