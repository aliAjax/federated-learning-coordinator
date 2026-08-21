package application

import (
	"context"
	"errors"
	"github.com/example/federated-learning-coordinator/internal/privacy/domain"
	"github.com/example/federated-learning-coordinator/internal/privacy/infrastructure"
	"testing"
	"time"
)

func TestDefaultBudgetDoesNotPanic(t *testing.T) {
	r := infrastructure.New()
	if err := r.Save(context.Background(), domain.Budget{CohortID: "c1", EpsilonLimit: 2, UpdatedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	s := New(r)
	defer func() {
		if v := recover(); v != nil {
			t.Fatalf("panic: %v", v)
		}
	}()
	if _, err := s.Consume(context.Background(), "c1", 1); err != nil {
		t.Fatal(err)
	}
	if got := r.List(context.Background()); got == nil {
		t.Fatal("nil budget list")
	}
	if _, err := s.Remaining(context.Background(), "missing"); !errors.Is(err, domain.ErrBudgetNotFound) {
		t.Fatalf("remaining lookup lost: %v", err)
	}
	if _, err := s.Consume(context.Background(), "missing", 1); !errors.Is(err, domain.ErrBudgetNotFound) {
		t.Fatalf("missing budget chain lost: %v", err)
	}
	if _, err := s.Consume(nil, "c1", 1); err != nil {
		t.Fatal(err)
	}
	if got := (domain.Budget{EpsilonLimit: 1, EpsilonUsed: 3}).Remaining(); got != 0 {
		t.Fatalf("remaining budget was not clamped: %v", got)
	}
	start := make(chan struct{})
	done := make(chan struct{})
	go func() {
		<-start
		for i := 0; i < 100; i++ {
			_ = r.List(context.Background())
		}
		close(done)
	}()
	close(start)
	for i := 0; i < 100; i++ {
		_ = r.Save(context.Background(), domain.Budget{CohortID: "x", EpsilonLimit: 2})
	}
	<-done
}
