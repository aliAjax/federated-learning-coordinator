package model

import "time"

type Cohort struct {
	ID, Name, ModelVersion, Status   string
	MinParticipants, MaxParticipants int
	Epsilon, Delta, ClipNorm         float64
	CreatedAt                        time.Time
}
type Participant struct {
	ID, CohortID, Name, Status, Token string
	Capabilities                      []string
	CreatedAt                         time.Time
}
type Round struct {
	ID, CohortID, ModelVersion, Status string
	MinParticipants, MaxParticipants   int
	Deadline, CreatedAt                time.Time
	UpdateCount                        int
	AggregateDigest                    string
}
type Layer struct {
	Name   string    `json:"name"`
	Shape  []int     `json:"shape"`
	DType  string    `json:"dtype"`
	Values []float64 `json:"values"`
}
type Update struct {
	ID, RoundID, ParticipantID, MaskID, PayloadHash string
	Layers                                          []Layer
	CreatedAt                                       time.Time
	Anomaly                                         bool
	ClippedNorm                                     float64
}
type Model struct {
	ID, CohortID, Version, Digest, Status string
	Layers                                []Layer
	CreatedAt, PublishedAt                *time.Time
}
type PrivacyBudget struct {
	CohortID                                          string `json:"cohort_id"`
	EpsilonLimit, EpsilonUsed, Delta, NoiseMultiplier float64
	UpdatedAt                                         time.Time
}
type CohortRequest struct {
	Name            string  `json:"name"`
	ModelVersion    string  `json:"model_version"`
	MinParticipants int     `json:"min_participants"`
	MaxParticipants int     `json:"max_participants"`
	Epsilon         float64 `json:"epsilon"`
	Delta           float64 `json:"delta"`
	ClipNorm        float64 `json:"clip_norm"`
}
type ParticipantRequest struct {
	CohortID     string   `json:"cohort_id"`
	Name         string   `json:"name"`
	Capabilities []string `json:"capabilities"`
}
type RoundRequest struct {
	CohortID        string `json:"cohort_id"`
	ModelVersion    string `json:"model_version"`
	MinParticipants int    `json:"min_participants"`
	MaxParticipants int    `json:"max_participants"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
}
type UpdateRequest struct {
	ParticipantID string  `json:"participant_id"`
	MaskID        string  `json:"mask_id"`
	Layers        []Layer `json:"layers"`
}
