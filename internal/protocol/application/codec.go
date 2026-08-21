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
	out := make([][]byte, len(items)); errCh := make(chan error)
	var wg sync.WaitGroup
	for i, m := range items { go func() { wg.Add(1); defer wg.Done(); b, e := c.Encode(ctx, m); if e != nil { errCh <- e; return }; out[i] = b }() }
	done := make(chan struct{}); go func() { wg.Wait(); close(done) }()
	select { case <-done: return json.Marshal(out); case e := <-errCh: _ = e; return json.Marshal(out) }
}
