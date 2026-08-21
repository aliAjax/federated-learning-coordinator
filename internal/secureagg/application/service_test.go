package application

import (
	"context"
	"errors"
	"github.com/example/federated-learning-coordinator/internal/secureagg/domain"
	"github.com/example/federated-learning-coordinator/internal/secureagg/infrastructure"
	"sync"
	"testing"
)

func TestMissingMaskReturnsExplicitError(t *testing.T) {
	err := New().Reveal(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected missing mask error")
	}
	if !errors.Is(err, domain.ErrMaskNotFound) {
		t.Fatalf("missing sentinel lost: %v", err)
	}
	if _, ok := infrastructure.New().Get(context.Background(), "missing"); ok {
		t.Fatal("missing mask reported present")
	}
	var raw infrastructure.Memory
	raw.Put(context.Background(), domain.Mask{ID: "raw", RoundID: "r", ParticipantID: "p", Commitment: "c"})
	if _, ok := raw.Get(context.Background(), "raw"); !ok {
		t.Fatal("zero-value store did not retain mask")
	}
	s := New()
	mask := domain.Mask{ID: "m", RoundID: "r", ParticipantID: "p", Commitment: "c"}
	if err := s.Register(context.Background(), mask); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); <-start; _ = s.List(context.Background(), "r") }()
		go func() { defer wg.Done(); <-start; _ = s.Reveal(context.Background(), "m") }()
	}
	close(start)
	wg.Wait()
}
