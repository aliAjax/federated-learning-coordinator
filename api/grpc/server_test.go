package grpcapi

import (
	"context"
	"errors"
	"github.com/example/federated-learning-coordinator/internal/round"
	appclock "github.com/example/federated-learning-coordinator/internal/round/application"
	"github.com/example/federated-learning-coordinator/internal/storage"
	"testing"
)

func TestGRPCMissingRoundKeepsNotFound(t *testing.T) {
	a := &API{App: round.NewService(storage.NewMemory(), nil, nil)}
	_, err := a.GetRound(context.Background(), &RoundRequest{ID: "missing"})
	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("error chain lost: %v", err)
	}
	_, err = a.SubmitUpdate(context.Background(), &UpdateRequest{RoundID: "missing"})
	if !errors.Is(err, storage.ErrNotFound) {
		t.Fatalf("submit chain lost: %v", err)
	}
	if !errors.Is(wrapGRPCError(storage.ErrNotFound), storage.ErrNotFound) {
		t.Fatal("wrapper lost chain")
	}
	if !errors.Is(appclock.PreserveRoundError(storage.ErrNotFound), storage.ErrNotFound) {
		t.Fatal("clock wrapper lost chain")
	}
	if !errors.Is(appclock.WrapRoundLookup(storage.ErrNotFound), storage.ErrNotFound) || !errors.Is(appclock.WrapRoundSubmit(storage.ErrNotFound), storage.ErrNotFound) {
		t.Fatal("round helper lost chain")
	}
	if !errors.Is(mapRoundStorage(storage.ErrNotFound), storage.ErrNotFound) {
		t.Fatal("storage mapper lost chain")
	}
	for i := 0; i < 10; i++ {
		if _, err := a.GetRound(context.Background(), &RoundRequest{ID: "missing"}); !errors.Is(err, storage.ErrNotFound) {
			t.Fatalf("round %d lost chain: %v", i, err)
		}
	}
}
