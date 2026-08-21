package adapter

import (
	"encoding/json"
	"github.com/example/federated-learning-coordinator/internal/privacy/domain"
)

func Encode(b domain.Budget) ([]byte, error) { return json.Marshal(b) }
func Decode(data []byte) (domain.Budget, error) {
	var b domain.Budget
	e := json.Unmarshal(data, &b)
	return b, e
}
func Summary(b domain.Budget) map[string]float64 {
	return map[string]float64{"limit": b.EpsilonLimit, "used": b.EpsilonUsed, "remaining": b.Remaining(), "delta": b.Delta}
}
