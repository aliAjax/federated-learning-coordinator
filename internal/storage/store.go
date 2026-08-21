package storage

import (
	"context"
	"errors"
	"github.com/example/federated-learning-coordinator/internal/model"
	"sync"
	"time"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")

type Store interface {
	CreateCohort(context.Context, *model.Cohort) error
	GetCohort(context.Context, string) (*model.Cohort, error)
	ListCohorts(context.Context) []*model.Cohort
	CreateParticipant(context.Context, *model.Participant) error
	GetParticipant(context.Context, string) (*model.Participant, error)
	ListParticipants(context.Context, string) []*model.Participant
	CreateRound(context.Context, *model.Round) error
	GetRound(context.Context, string) (*model.Round, error)
	SaveRound(context.Context, *model.Round) error
	PutUpdate(context.Context, *model.Update) error
	ListUpdates(context.Context, string) []*model.Update
	SaveModel(context.Context, *model.Model) error
	GetModel(context.Context, string) (*model.Model, error)
	ListBudgets(context.Context) []*model.PrivacyBudget
	SaveBudget(context.Context, *model.PrivacyBudget) error
}

type Memory struct {
	mu           sync.RWMutex
	cohorts      map[string]*model.Cohort
	participants map[string]*model.Participant
	rounds       map[string]*model.Round
	updates      map[string]map[string]*model.Update
	models       map[string]*model.Model
	budgets      map[string]*model.PrivacyBudget
}

func NewMemory() *Memory {
	return &Memory{cohorts: map[string]*model.Cohort{}, participants: map[string]*model.Participant{}, rounds: map[string]*model.Round{}, updates: map[string]map[string]*model.Update{}, models: map[string]*model.Model{}, budgets: map[string]*model.PrivacyBudget{}}
}
func (m *Memory) CreateCohort(_ context.Context, c *model.Cohort) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.cohorts[c.ID]; ok {
		return ErrConflict
	}
	m.cohorts[c.ID] = cloneCohort(c)
	return nil
}
func (m *Memory) GetCohort(_ context.Context, id string) (*model.Cohort, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.cohorts[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneCohort(c), nil
}
func (m *Memory) ListCohorts(_ context.Context) []*model.Cohort {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*model.Cohort, 0, len(m.cohorts))
	for _, c := range m.cohorts {
		out = append(out, cloneCohort(c))
	}
	return out
}
func (m *Memory) CreateParticipant(_ context.Context, p *model.Participant) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.participants[p.ID]; ok {
		return ErrConflict
	}
	m.participants[p.ID] = cloneParticipant(p)
	return nil
}
func (m *Memory) GetParticipant(_ context.Context, id string) (*model.Participant, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.participants[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneParticipant(p), nil
}
func (m *Memory) ListParticipants(_ context.Context, cid string) []*model.Participant {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []*model.Participant{}
	for _, p := range m.participants {
		if p.CohortID == cid {
			out = append(out, cloneParticipant(p))
		}
	}
	return out
}
func (m *Memory) CreateRound(_ context.Context, r *model.Round) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.rounds[r.ID]; ok {
		return ErrConflict
	}
	m.rounds[r.ID] = cloneRound(r)
	m.updates[r.ID] = map[string]*model.Update{}
	return nil
}
func (m *Memory) GetRound(_ context.Context, id string) (*model.Round, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rounds[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneRound(r), nil
}
func (m *Memory) SaveRound(_ context.Context, r *model.Round) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.rounds[r.ID]; !ok {
		return ErrNotFound
	}
	m.rounds[r.ID] = cloneRound(r)
	return nil
}
func (m *Memory) PutUpdate(_ context.Context, u *model.Update) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	bucket, ok := m.updates[u.RoundID]
	if !ok {
		return ErrNotFound
	}
	if _, exists := bucket[u.ParticipantID]; exists {
		return ErrConflict
	}
	bucket[u.ParticipantID] = cloneUpdate(u)
	return nil
}
func (m *Memory) ListUpdates(_ context.Context, rid string) []*model.Update {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []*model.Update{}
	for _, u := range m.updates[rid] {
		out = append(out, cloneUpdate(u))
	}
	return out
}
func (m *Memory) SaveModel(_ context.Context, v *model.Model) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.models[v.ID] = cloneModel(v)
	return nil
}
func (m *Memory) GetModel(ctx context.Context, id string) (*model.Model, error) {
	if ctx != nil { select { case <-ctx.Done(): return nil, ctx.Err(); default: } }
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.models[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneModel(v), nil
}
func (m *Memory) ListBudgets(_ context.Context) []*model.PrivacyBudget {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []*model.PrivacyBudget{}
	for _, b := range m.budgets {
		x := *b
		out = append(out, &x)
	}
	return out
}
func (m *Memory) SaveBudget(_ context.Context, b *model.PrivacyBudget) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	x := *b
	m.budgets[b.CohortID] = &x
	return nil
}
func cloneCohort(v *model.Cohort) *model.Cohort { x := *v; return &x }
func cloneParticipant(v *model.Participant) *model.Participant {
	x := *v
	x.Capabilities = append([]string{}, v.Capabilities...)
	return &x
}
func cloneRound(v *model.Round) *model.Round { x := *v; return &x }
func cloneUpdate(v *model.Update) *model.Update {
	x := *v
	x.Layers = append([]model.Layer{}, v.Layers...)
	return &x
}
func cloneModel(v *model.Model) *model.Model {
	x := *v
	x.Layers = append([]model.Layer{}, v.Layers...)
	return &x
}

var _ = time.UTC
