package round

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/aggregation"
	"github.com/example/federated-learning-coordinator/internal/anomaly"
	"github.com/example/federated-learning-coordinator/internal/model"
	"github.com/example/federated-learning-coordinator/internal/storage"
	"github.com/example/federated-learning-coordinator/internal/tensor"
	"time"
)

var ErrValidation = errors.New("validation failed")
var ErrInvalidState = errors.New("invalid state transition")

type Service struct {
	store    storage.Store
	agg      aggregation.Aggregator
	detector *anomaly.Detector
	now      func() time.Time
}

func NewService(s storage.Store, a aggregation.Aggregator, d *anomaly.Detector) *Service {
	return &Service{store: s, agg: a, detector: d, now: time.Now}
}
func ID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return prefix + "_" + hex.EncodeToString(b)
}
func (s *Service) CreateCohort(ctx context.Context, q model.CohortRequest) (*model.Cohort, error) {
	if q.Name == "" || q.ModelVersion == "" {
		return nil, fmt.Errorf("%w: name and model_version required", ErrValidation)
	}
	if q.MinParticipants <= 0 {
		q.MinParticipants = 2
	}
	if q.MaxParticipants < q.MinParticipants {
		q.MaxParticipants = q.MinParticipants
	}
	if q.Epsilon <= 0 {
		q.Epsilon = 8
	}
	if q.Delta <= 0 {
		q.Delta = .000001
	}
	if q.ClipNorm <= 0 {
		q.ClipNorm = 5
	}
	c := &model.Cohort{ID: ID("cohort"), Name: q.Name, ModelVersion: q.ModelVersion, Status: "active", MinParticipants: q.MinParticipants, MaxParticipants: q.MaxParticipants, Epsilon: q.Epsilon, Delta: q.Delta, ClipNorm: q.ClipNorm, CreatedAt: s.now().UTC()}
	if err := s.store.CreateCohort(ctx, c); err != nil {
		return nil, fmt.Errorf("create cohort: %w", err)
	}
	b := &model.PrivacyBudget{CohortID: c.ID, EpsilonLimit: c.Epsilon, Delta: c.Delta, NoiseMultiplier: .5, UpdatedAt: s.now().UTC()}
	if err := s.store.SaveBudget(ctx, b); err != nil {
		return nil, fmt.Errorf("create privacy budget: %w", err)
	}
	return c, nil
}
func (s *Service) Register(ctx context.Context, q model.ParticipantRequest) (*model.Participant, error) {
	c, err := s.store.GetCohort(ctx, q.CohortID)
	if err != nil {
		return nil, fmt.Errorf("cohort: %w", err)
	}
	if c.Status != "active" || q.Name == "" {
		return nil, fmt.Errorf("%w: inactive cohort or empty name", ErrValidation)
	}
	if len(q.Capabilities) == 0 {
		q.Capabilities = []string{"float32"}
	}
	p := &model.Participant{ID: ID("participant"), CohortID: c.ID, Name: q.Name, Status: "active", Token: ID("token"), Capabilities: q.Capabilities, CreatedAt: s.now().UTC()}
	if err := s.store.CreateParticipant(ctx, p); err != nil {
		return nil, fmt.Errorf("register participant: %w", err)
	}
	return p, nil
}
func (s *Service) CreateRound(ctx context.Context, q model.RoundRequest) (*model.Round, error) {
	c, err := s.store.GetCohort(ctx, q.CohortID)
	if err != nil {
		return nil, fmt.Errorf("cohort: %w", err)
	}
	if q.ModelVersion == "" {
		q.ModelVersion = c.ModelVersion
	}
	if q.ModelVersion != c.ModelVersion {
		return nil, fmt.Errorf("%w: model version mismatch", ErrValidation)
	}
	if q.MinParticipants <= 0 {
		q.MinParticipants = c.MinParticipants
	}
	if q.MaxParticipants <= 0 {
		q.MaxParticipants = c.MaxParticipants
	}
	if q.MinParticipants > q.MaxParticipants {
		return nil, fmt.Errorf("%w: participant limits", ErrValidation)
	}
	if q.TimeoutSeconds <= 0 {
		q.TimeoutSeconds = 600
	}
	now := s.now().UTC()
	r := &model.Round{ID: ID("round"), CohortID: c.ID, ModelVersion: q.ModelVersion, Status: "collecting", MinParticipants: q.MinParticipants, MaxParticipants: q.MaxParticipants, CreatedAt: now, Deadline: now.Add(time.Duration(q.TimeoutSeconds) * time.Second)}
	if err := s.store.CreateRound(ctx, r); err != nil {
		return nil, fmt.Errorf("create round: %w", err)
	}
	return r, nil
}
func (s *Service) GetRound(ctx context.Context, id string) (*model.Round, error) {
	return s.store.GetRound(ctx, id)
}
func (s *Service) Abort(ctx context.Context, id string) (*model.Round, error) {
	r, err := s.store.GetRound(ctx, id)
	if err != nil {
		return nil, err
	}
	if r.Status == "completed" {
		return nil, ErrInvalidState
	}
	r.Status = "aborted"
	if err := s.store.SaveRound(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}
func (s *Service) Submit(ctx context.Context, rid string, q model.UpdateRequest) (*model.Update, error) {
	r, err := s.store.GetRound(ctx, rid)
	if err != nil {
		return nil, err
	}
	if r.Status != "collecting" && r.Status != "masked" {
		return nil, ErrInvalidState
	}
	p, err := s.store.GetParticipant(ctx, q.ParticipantID)
	if err != nil {
		return nil, err
	}
	if p.CohortID != r.CohortID || p.Status != "active" {
		return nil, ErrValidation
	}
	if err := tensor.Validate(q.Layers, 1_000_000); err != nil {
		return nil, err
	}
	c, err := s.store.GetCohort(ctx, r.CohortID)
	if err != nil {
		return nil, err
	}
	clipped, originalNorm := tensor.Clip(q.Layers, c.ClipNorm)
	check := s.detector.Check(q.Layers, nil)
	u := &model.Update{ID: ID("update"), RoundID: r.ID, ParticipantID: p.ID, MaskID: q.MaskID, Layers: clipped, PayloadHash: tensor.Digest(clipped), CreatedAt: s.now().UTC(), Anomaly: check.Anomalous, ClippedNorm: originalNorm}
	if err := s.store.PutUpdate(ctx, u); err != nil {
		return nil, fmt.Errorf("save update: %w", err)
	}
	updates := s.store.ListUpdates(ctx, r.ID)
	r.UpdateCount = len(updates)
	if r.UpdateCount >= r.MinParticipants {
		r.Status = "masked"
	}
	if err := s.store.SaveRound(ctx, r); err != nil {
		return nil, err
	}
	return u, nil
}
func (s *Service) Aggregate(ctx context.Context, rid string) (*model.Model, error) {
	r, err := s.store.GetRound(ctx, rid)
	if err != nil {
		return nil, err
	}
	if r.Status != "masked" {
		return nil, ErrInvalidState
	}
	updates := s.store.ListUpdates(ctx, r.ID)
	accepted := make([]*model.Update, 0, len(updates))
	for _, u := range updates {
		if !u.Anomaly {
			accepted = append(accepted, u)
		}
	}
	if len(accepted) < r.MinParticipants {
		return nil, fmt.Errorf("%w: insufficient accepted updates", ErrValidation)
	}
	r.Status = "aggregating"
	if err := s.store.SaveRound(ctx, r); err != nil {
		return nil, err
	}
	layers, err := s.agg.Aggregate(accepted)
	if err != nil {
		return nil, fmt.Errorf("aggregate: %w", err)
	}
	budgets := s.store.ListBudgets(ctx)
	var budget *model.PrivacyBudget
	for _, b := range budgets {
		if b.CohortID == r.CohortID {
			budget = b
			break
		}
	}
	if budget == nil || budget.EpsilonUsed+1 > budget.EpsilonLimit {
		r.Status = "masked"
		_ = s.store.SaveRound(ctx, r)
		return nil, errors.New("privacy budget exhausted")
	}
	budget.EpsilonUsed++
	budget.UpdatedAt = s.now().UTC()
	_ = s.store.SaveBudget(ctx, budget)
	digest := tensor.Digest(layers)
	now := s.now().UTC()
	v := &model.Model{ID: r.ID, CohortID: r.CohortID, Version: r.ModelVersion + "-" + r.ID, Digest: digest, Status: "candidate", Layers: layers, CreatedAt: &now}
	if err := s.store.SaveModel(ctx, v); err != nil {
		return nil, err
	}
	r.Status = "completed"
	r.AggregateDigest = digest
	if err := s.store.SaveRound(ctx, r); err != nil {
		return nil, err
	}
	return v, nil
}
func (s *Service) Publish(ctx context.Context, id string) (*model.Model, error) {
	v, err := s.store.GetModel(ctx, id)
	if err != nil {
		return nil, err
	}
	if v.Status != "candidate" && v.Status != "published" {
		return nil, ErrInvalidState
	}
	now := s.now().UTC()
	v.Status = "published"
	v.PublishedAt = &now
	if err := s.store.SaveModel(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}
func (s *Service) GetModel(ctx context.Context, id string) (*model.Model, error) {
	return s.store.GetModel(ctx, id)
}
func (s *Service) Budgets(ctx context.Context) []*model.PrivacyBudget {
	return s.store.ListBudgets(ctx)
}
