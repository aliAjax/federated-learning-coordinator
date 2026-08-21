package application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/protocol/domain"
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
	for i, m := range items {
		b, e := c.Encode(ctx, m)
		if e != nil {
			return nil, e
		}
		out[i] = b
	}
	return json.Marshal(out)
}
