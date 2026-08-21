package grpcapi

import (
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/storage"
)

func mapRoundLookup(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("round lookup: %w", err)
}
func mapRoundSubmit(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("round submit: %w", err)
}
func mapRoundStorage(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("storage: %w", err)
}

var _ = storage.ErrNotFound
