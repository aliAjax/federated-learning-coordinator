package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Event struct {
	ID, Type, Actor, Subject string
	Payload                  map[string]any
	PreviousHash, Hash       string
	CreatedAt                time.Time
}

func (e Event) Compute(prev string) string {
	b, _ := json.Marshal(struct {
		Type, Actor, Subject string
		Payload              map[string]any
		Prev                 string
	}{e.Type, e.Actor, e.Subject, e.Payload, prev})
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func (e *Event) Seal(prev string) { e.PreviousHash = prev; e.Hash = e.Compute(prev) }
func (e Event) Verify() bool      { return e.Hash != "" && e.Hash == e.Compute(e.PreviousHash) }
