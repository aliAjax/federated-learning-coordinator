package infrastructure

import (
	"context"
	"github.com/example/federated-learning-coordinator/internal/crypto/domain"
	"sync"
)

type KeyStore struct {
	mu    sync.RWMutex
	items map[string][]byte
}

func New() *KeyStore { return &KeyStore{items: map[string][]byte{}} }
func (k *KeyStore) Put(_ context.Context, id string, key []byte) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.items[id] = append([]byte(nil), key...)
}
func (k *KeyStore) Get(_ context.Context, id string) ([]byte, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	v, ok := k.items[id]
	return append([]byte(nil), v...), ok
}
func (_ *KeyStore) Describe(s domain.Signature) string { return s.Algorithm }
