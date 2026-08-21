package domain

import (
	"encoding/json"
	"errors"
	"time"
)

type Kind string

const (
	Hello        Kind = "hello"
	RoundState   Kind = "round_state"
	MaskedUpdate Kind = "masked_update"
	Recovery     Kind = "recovery"
	Ack          Kind = "ack"
	Error        Kind = "error"
)

type Message struct {
	ID            string
	Kind          Kind
	RoundID       string
	ParticipantID string
	Payload       json.RawMessage
	CreatedAt     time.Time
}

var ErrIdentityRequired = errors.New("message identity required")
var ErrParticipantRequired = errors.New("participant required")

func (m Message) Validate() error {
	if m.ID == "" || m.RoundID == "" || m.Kind == "" {
		return ErrIdentityRequired
	}
	if m.Kind != Hello && m.ParticipantID == "" {
		return ErrParticipantRequired
	}
	return nil
}
func (m Message) IsControl() bool {
	return m.Kind == Hello || m.Kind == RoundState || m.Kind == Ack || m.Kind == Error
}
