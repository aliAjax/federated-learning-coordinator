package storage

import (
	"context"
	"sync"
	"testing"
	"github.com/example/federated-learning-coordinator/internal/model"
	aggregationinfra "github.com/example/federated-learning-coordinator/internal/aggregation/infrastructure"
	"github.com/example/federated-learning-coordinator/internal/aggregation/domain"
)

func TestConcurrentUpdateSnapshot(t *testing.T) {
	m := NewMemory()
	r := &model.Round{ID: "r1", CohortID: "c1", ModelVersion: "v1", MinParticipants: 1, MaxParticipants: 4}
	if err := m.CreateRound(context.Background(), r); err != nil { t.Fatal(err) }
	var wg sync.WaitGroup
	start := make(chan struct{})
	readerDone := make(chan struct{})
	go func() {
		defer close(readerDone)
		for i := 0; i < 1000; i++ {
			_ = m.ListUpdates(context.Background(), r.ID)
		}
	}()
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_ = m.PutUpdate(context.Background(), &model.Update{ID: string(rune('u'+i)), RoundID: r.ID, ParticipantID: string(rune('p'+i)), Layers: []model.Layer{{Name: "w", DType: "float32", Shape: []int{1}, Values: []float64{float64(i)}}}})
		}(i)
	}
	close(start)
	wg.Wait()
	<-readerDone
	if got := len(m.ListUpdates(context.Background(), r.ID)); got != 8 { t.Fatalf("got %d updates", got) }
	results := aggregationinfra.New()
	for i := 0; i < 20; i++ { results.Save(context.Background(), domain.Result{RoundID: string(rune(i)), Digest: "d", Algorithm: "a", Accepted: 1}) }
	go func() { for i := 0; i < 1000; i++ { results.List(context.Background()) } }()
	for i := 0; i < 1000; i++ { results.Save(context.Background(), domain.Result{RoundID: string(rune(i)), Digest: "d", Algorithm: "a", Accepted: 1}) }
}
