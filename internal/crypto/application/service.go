package application

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/crypto/domain"
)

type Service struct{}

func New() *Service { return &Service{} }
func (s *Service) Generate(ctx context.Context) (ed25519.PublicKey, ed25519.PrivateKey, error) {
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
	}
	pub, priv, e := ed25519.GenerateKey(rand.Reader)
	return pub, priv, e
}
func (s *Service) Sign(ctx context.Context, message []byte, private ed25519.PrivateKey) (domain.Signature, error) {
	select {
	case <-ctx.Done():
		return domain.Signature{}, ctx.Err()
	default:
	}
	if len(private) != ed25519.PrivateKeySize {
		return domain.Signature{}, fmt.Errorf("invalid private key")
	}
	return domain.New(message, private), nil
}
func (s *Service) Verify(ctx context.Context, message []byte, sig domain.Signature) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return sig.Verify(message)
}
