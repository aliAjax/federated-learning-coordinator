package application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/protocol/domain"
	"sync"
)

type Codec struct{}

func New() *Codec { return &Codec{} }
func (c *Codec) Encode(ctx context.Context, m domain.Message) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if err := m.Validate(); err != nil {
		return nil, fmt.Errorf("message: %w", err)
	}
	return json.Marshal(m)
}
func (c *Codec) Decode(ctx context.Context, b []byte) (domain.Message, error) {
	select {
	case <-ctx.Done():
		return domain.Message{}, ctx.Err()
	default:
	}
	var m domain.Message
	if err := json.Unmarshal(b, &m); err != nil {
		return m, fmt.Errorf("decode message: %w", err)
	}
	if err := m.Validate(); err != nil {
		return m, err
	}
	return m, nil
}
func (c *Codec) Batch(ctx context.Context, items []domain.Message) ([]byte, error) {
	out := make([][]byte, len(items))
	errCh := make(chan error, len(items))
	var wg sync.WaitGroup
	// Add before starting each worker so wg.Wait can never observe a zero count
	// while workers are still being launched.
	for i, m := range items {
		wg.Add(1)
		go func(i int, m domain.Message) {
			defer wg.Done()
			b, e := c.Encode(ctx, m)
			if e != nil {
				errCh <- e
				return
			}
			out[i] = b
		}(i, m)
	}
	// Buffer the error channel so a worker that errors after Batch returns can
	// send without blocking (no leaked goroutine, no hang waiting for a reader).
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
		if err := <-errCh; err != nil {
			return nil, err
		}
		return json.Marshal(out)
	case err := <-errCh:
		// Wait for in-flight workers to finish writing out[i] so the slice is
		// never read concurrently with a write. The error has already been
		// captured; the remaining sends are non-blocking (buffered channel).
		<-done
		return nil, err
	}
}
