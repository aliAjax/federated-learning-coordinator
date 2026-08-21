package domain

import (
	"errors"
	"strings"
	"time"
)

type Status string

const (
	Candidate Status = "candidate"
	Published Status = "published"
	Revoked   Status = "revoked"
)

type Model struct {
	ID, CohortID, Version, Digest string
	Status                        Status
	Layers                        []Layer
	CreatedAt, PublishedAt        *time.Time
}
type Layer struct {
	Name, DType string
	Shape       []int
	Values      []float64
}

func (m Model) Validate() error {
	if m.ID == "" || m.CohortID == "" || m.Version == "" || m.Digest == "" {
		return errors.New("model identity required")
	}
	if len(m.Layers) == 0 {
		return errors.New("model layers required")
	}
	return nil
}
func (l Layer) Validate() error {
	if strings.TrimSpace(l.Name) == "" || len(l.Shape) == 0 || len(l.Values) == 0 {
		return errors.New("invalid layer")
	}
	if l.DType != "float32" && l.DType != "float64" && l.DType != "fixed32" {
		return errors.New("unsupported dtype")
	}
	return nil
}
func (m *Model) Publish(now time.Time) error {
	if m.Status != Candidate && m.Status != Published {
		return errors.New("model cannot publish")
	}
	m.Status = Published
	m.PublishedAt = &now
	return nil
}
func (m *Model) Revoke() { m.Status = Revoked }
