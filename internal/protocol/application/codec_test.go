package application

import (
	"context"
	"testing"
	"time"
	"sync"
	"github.com/example/federated-learning-coordinator/internal/protocol/domain"
	infra "github.com/example/federated-learning-coordinator/internal/protocol/infrastructure"
)

func TestBatchPublishErrorDoesNotHang(t *testing.T) {
	items := []domain.Message{{ID: "ok", Kind: domain.Hello, RoundID: "r"}, {ID: "bad", Kind: domain.RoundState, RoundID: "r"}}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond); defer cancel()
	codec := New()
	var wg sync.WaitGroup
	start := make(chan struct{})
	errs := make(chan error, 2)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; _, err := codec.Batch(ctx, items); errs <- err }()
	}
	close(start)
	wg.Wait()
	first, second := <-errs, <-errs
	if first == nil || second == nil { t.Fatal("expected validation error") }
	bus := infra.New()
	busStart := make(chan struct{})
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) { defer wg.Done(); <-busStart; _ = bus.Publish(context.Background(), domain.Message{ID: string(rune('a'+i)), Kind: domain.Hello, RoundID: "r"}) }(i)
	}
	close(busStart)
	for i := 0; i < 100; i++ { _ = bus.List(context.Background(), "r") }
	wg.Wait()
	canceled, cancelPublish := context.WithCancel(context.Background())
	cancelPublish()
	if err := bus.Publish(canceled, items[0]); err == nil { t.Fatal("canceled publish succeeded") }
}
