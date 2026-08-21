package adapter

import (
	"encoding/json"
	"github.com/example/federated-learning-coordinator/internal/aggregation/domain"
)

func Encode(r domain.Result) ([]byte, error) { return json.Marshal(r) }
func Decode(b []byte) (domain.Result, error) {
	var r domain.Result
	e := json.Unmarshal(b, &r)
	return r, e
}
