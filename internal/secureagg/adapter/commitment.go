package adapter

import (
	"crypto/rand"
	"github.com/example/federated-learning-coordinator/internal/secureagg/domain"
)

func NewSeed(n int) ([]byte, error) { b := make([]byte, n); _, e := rand.Read(b); return b, e }
func NewCommitment(n int) (string, []byte, error) {
	b, e := NewSeed(n)
	if e != nil {
		return "", nil, e
	}
	return domain.Commitment(b), b, nil
}
func Verify(c string, seed []byte) bool { return c == domain.Commitment(seed) }
