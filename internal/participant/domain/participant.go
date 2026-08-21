package domain

import (
	"errors"
	"strings"
	"time"
)

type Status string

const (
	Active  Status = "active"
	Paused  Status = "paused"
	Removed Status = "removed"
)

type Participant struct {
	ID, CohortID, Name string
	Status             Status
	Capabilities       []string
	CreatedAt          time.Time
	LastSeen           *time.Time
}

func (p Participant) Validate() error {
	if strings.TrimSpace(p.ID) == "" || strings.TrimSpace(p.CohortID) == "" {
		return errors.New("participant identity required")
	}
	if p.Name == "" {
		return errors.New("participant name required")
	}
	if len(p.Capabilities) == 0 {
		return errors.New("participant capabilities required")
	}
	return nil
}
func (p Participant) CanUpload(modelVersion string, cap string) bool {
	if p.Status != Active {
		return false
	}
	for _, v := range p.Capabilities {
		if v == cap || v == modelVersion {
			return true
		}
	}
	return false
}
func (p *Participant) Touch(now time.Time) { p.LastSeen = &now }
func NormalizeCapabilities(items []string) []string {
	out := []string{}
	seen := map[string]bool{}
	for _, v := range items {
		v = strings.TrimSpace(strings.ToLower(v))
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}
