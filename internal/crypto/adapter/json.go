package adapter

import (
	"encoding/json"
	"github.com/example/federated-learning-coordinator/internal/crypto/domain"
)

func Encode(s domain.Signature) ([]byte, error) { return json.Marshal(s) }
func Decode(b []byte) (domain.Signature, error) {
	var s domain.Signature
	e := json.Unmarshal(b, &s)
	return s, e
}
func Algorithm(s domain.Signature) string { return s.Algorithm }
