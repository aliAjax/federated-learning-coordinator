package adapter

import (
	"encoding/json"
	"github.com/example/federated-learning-coordinator/internal/round/domain"
)

func Encode(r domain.Round) ([]byte, error) { return json.Marshal(r) }
func Decode(b []byte) (domain.Round, error) {
	var r domain.Round
	e := json.Unmarshal(b, &r)
	return r, e
}
func Statuses() []domain.Status {
	return []domain.Status{domain.Collecting, domain.Masked, domain.Aggregating, domain.Completed, domain.Aborted}
}
