package postgres

import (
	"context"
	"errors"
	"github.com/example/federated-learning-coordinator/internal/protocol/domain"
	"testing"
)

func TestMissingModelPreservesNotFound(t *testing.T) {
	_, err := New("dsn").LoadModel(context.Background(), "missing")
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("error chain lost: %v", err)
	}
	checks := []error{New("").Ping(context.Background()), New("").SaveModel(context.Background(), nil), New("").Migrate(context.Background())}
	for _, check := range checks {
		if !errors.Is(check, ErrUnavailable) {
			t.Fatalf("chain lost: %v", check)
		}
	}
	if !errors.Is((domain.Message{ID: "m", RoundID: "r", Kind: domain.RoundState}).Validate(), domain.ErrParticipantRequired) {
		t.Fatal("participant sentinel lost")
	}
	if !errors.Is((domain.Message{}).Validate(), domain.ErrIdentityRequired) {
		t.Fatal("identity sentinel lost")
	}
}
