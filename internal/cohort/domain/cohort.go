package domain

import (
	"errors"
	"math"
	"strings"
	"time"
)

type Status string

const (
	Enabled  Status = "enabled"
	Disabled Status = "disabled"
)

type Cohort struct {
	ID, Name, ModelVersion           string
	Status                           Status
	MinParticipants, MaxParticipants int
	Epsilon, Delta, ClipNorm         float64
	CreatedAt                        time.Time
}

func (c Cohort) Validate() error {
	if strings.TrimSpace(c.ID) == "" || strings.TrimSpace(c.Name) == "" || c.ModelVersion == "" {
		return errors.New("cohort identity required")
	}
	if c.MinParticipants < 1 || c.MaxParticipants < c.MinParticipants {
		return errors.New("invalid participant bounds")
	}
	if c.Epsilon <= 0 || c.Delta <= 0 || c.ClipNorm <= 0 {
		return errors.New("privacy parameters must be positive")
	}
	if math.IsInf(c.Epsilon, 0) || math.IsNaN(c.Epsilon) {
		return errors.New("invalid epsilon")
	}
	return nil
}
func (c Cohort) Accepts(n int) bool {
	return c.Status == Enabled && n >= c.MinParticipants && n <= c.MaxParticipants
}
func (c *Cohort) Disable() { c.Status = Disabled }
func (c *Cohort) Enable()  { c.Status = Enabled }
