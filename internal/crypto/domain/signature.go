package domain

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type Signature struct {
	Algorithm        string
	PublicKey, Value []byte
	MessageDigest    string
}

func New(message []byte, private ed25519.PrivateKey) Signature {
	h := sha256.Sum256(message)
	return Signature{Algorithm: "ed25519", PublicKey: private.Public().(ed25519.PublicKey), Value: ed25519.Sign(private, message), MessageDigest: hex.EncodeToString(h[:])}
}
func (s Signature) Verify(message []byte) error {
	if s.Algorithm != "ed25519" || len(s.PublicKey) != ed25519.PublicKeySize {
		return errors.New("unsupported signature")
	}
	if !ed25519.Verify(s.PublicKey, message, s.Value) {
		return errors.New("signature verification failed")
	}
	return nil
}
func Digest(message []byte) string { h := sha256.Sum256(message); return hex.EncodeToString(h[:]) }
