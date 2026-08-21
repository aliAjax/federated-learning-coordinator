package infrastructure

import (
	"context"
	"github.com/example/federated-learning-coordinator/internal/worker/domain"
	"sync"
)

type Queue struct {
	mu    sync.Mutex
	items []domain.Job
}

func New() *Queue { return &Queue{items: []domain.Job{}} }
func (q *Queue) Push(_ context.Context, j domain.Job) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, j)
}
func (q *Queue) Pop(_ context.Context) (domain.Job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return domain.Job{}, false
	}
	j := q.items[0]
	q.items = q.items[1:]
	return j, true
}
func (q *Queue) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }
