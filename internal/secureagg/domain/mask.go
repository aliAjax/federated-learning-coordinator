package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

var ErrMaskNotFound = errors.New("mask not found")

type Mask struct {
	ID, RoundID, ParticipantID, Commitment string
	CreatedAt                              time.Time
	Revealed                               bool
}

func (m Mask) Validate() error {
	if m.ID == "" || m.RoundID == "" || m.ParticipantID == "" || m.Commitment == "" {
		return errors.New("mask identity required")
	}
	return nil
}
func Commitment(seed []byte) string { h := sha256.Sum256(seed); return hex.EncodeToString(h[:]) }
func (m *Mask) Reveal()             { m.Revealed = true }
func (m Mask) CanRecover() bool     { return !m.Revealed }
