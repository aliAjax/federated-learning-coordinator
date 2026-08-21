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

func (m Message) Validate() error {
	if m.ID == "" || m.RoundID == "" || m.Kind == "" {
		return errors.New("message identity required")
	}
	if m.Kind != Hello && m.ParticipantID == "" {
		return errors.New("participant required")
	}
	return nil
}
func (m Message) IsControl() bool {
	return m.Kind == Hello || m.Kind == RoundState || m.Kind == Ack || m.Kind == Error
}
