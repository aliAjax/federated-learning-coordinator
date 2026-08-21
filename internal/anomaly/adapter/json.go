package adapter

import (
	"encoding/json"
	"github.com/example/federated-learning-coordinator/internal/anomaly/domain"
)

func Encode(r domain.Rule) ([]byte, error) { return json.Marshal(r) }
func Decode(b []byte) (domain.Rule, error) {
	var r domain.Rule
	e := json.Unmarshal(b, &r)
	return r, e
}
