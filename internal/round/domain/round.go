package domain

import (
	"errors"
	"time"
)

type Status string

const (
	Collecting  Status = "collecting"
	Masked      Status = "masked"
	Aggregating Status = "aggregating"
	Completed   Status = "completed"
	Aborted     Status = "aborted"
)

type Round struct {
	ID, CohortID, ModelVersion                    string
	Status                                        Status
	MinParticipants, MaxParticipants, UpdateCount int
	CreatedAt, Deadline                           time.Time
	AggregateDigest                               string
}

func (r Round) Validate() error {
	if r.ID == "" || r.CohortID == "" || r.ModelVersion == "" {
		return errors.New("round identity required")
	}
	if r.MinParticipants < 1 || r.MaxParticipants < r.MinParticipants {
		return errors.New("invalid round bounds")
	}
	if r.Deadline.Before(r.CreatedAt) {
		return errors.New("deadline before creation")
	}
	return nil
}
func (r Round) Expired(now time.Time) bool {
	return now.After(r.Deadline) && r.Status != Completed && r.Status != Aborted
}
func Transition(from, to Status) bool {
	switch from {
	case Collecting:
		return to == Masked || to == Aborted
	case Masked:
		return to == Aggregating || to == Aborted
	case Aggregating:
		return to == Masked || to == Aborted
	case Completed, Aborted:
		return false
	}
	return false
}
func (r *Round) MarkMasked() {
	if r.Status == Collecting {
		r.Status = Masked
	}
}
func (r *Round) MarkAggregating() {
	r.Status = Aggregating
}
func (r *Round) MarkCompleted(digest string) {
	r.AggregateDigest = digest
}
