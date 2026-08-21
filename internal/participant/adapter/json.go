package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/participant/domain"
)

type Request struct {
	ID, CohortID, Name string
	Capabilities       []string `json:"capabilities"`
}

func Decode(b []byte) (domain.Participant, error) {
	var r Request
	if err := json.Unmarshal(b, &r); err != nil {
		return domain.Participant{}, fmt.Errorf("participant json: %w", err)
	}
	return domain.Participant{ID: r.ID, CohortID: r.CohortID, Name: r.Name, Capabilities: r.Capabilities, Status: domain.Active}, nil
}
func Encode(p domain.Participant) ([]byte, error) { return json.Marshal(p) }
func Clone(p domain.Participant) domain.Participant {
	p.Capabilities = append([]string(nil), p.Capabilities...)
	return p
}
