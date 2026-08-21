package application

import (
	"context"
	"github.com/example/federated-learning-coordinator/internal/worker/domain"
	"sync"
	"time"
)

type Handler func(context.Context, domain.Job) error
type Runner struct {
	mu       sync.Mutex
	jobs     map[string]domain.Job
	handlers map[string]Handler
}

func New() *Runner { return &Runner{jobs: map[string]domain.Job{}, handlers: map[string]Handler{}} }
func (r *Runner) Register(kind string, h Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[kind] = h
}
func (r *Runner) Enqueue(j domain.Job) error {
	if err := j.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[j.ID] = j
	return nil
}
func (r *Runner) Run(ctx context.Context, id string) error {
	r.mu.Lock()
	j, ok := r.jobs[id]
	h := r.handlers[j.Kind]
	if ok {
		j.Start()
		r.jobs[id] = j
	}
	r.mu.Unlock()
	if !ok || h == nil {
		return context.Canceled
	}
	err := h(ctx, j)
	r.mu.Lock()
	defer r.mu.Unlock()
	if err != nil {
		j.Retry(time.Now().UTC(), err)
	} else {
		j.Complete()
	}
	r.jobs[id] = j
	return err
}
func (r *Runner) Get(id string) (domain.Job, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	j, ok := r.jobs[id]
	return j, ok
}
