package infrastructure

import (
	"context"
	"github.com/example/federated-learning-coordinator/internal/protocol/domain"
	"sync"
)

type Bus struct {
	mu       sync.Mutex
	messages []domain.Message
}

func New() *Bus { return &Bus{messages: []domain.Message{}} }
func (b *Bus) Publish(ctx context.Context, m domain.Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.messages = append(b.messages, m)
	return nil
}
func (b *Bus) List(_ context.Context, round string) []domain.Message {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]domain.Message, 0, len(b.messages))
	for _, m := range b.messages {
		if round == "" || m.RoundID == round {
			out = append(out, m)
		}
	}
	return out
}
func (b *Bus) Count() int { b.mu.Lock(); defer b.mu.Unlock(); return len(b.messages) }
