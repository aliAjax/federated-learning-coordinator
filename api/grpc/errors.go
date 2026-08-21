package grpcapi

import (
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/storage"
)

func mapRoundLookup(err error) error { return fmt.Errorf("round lookup: %v", err) }
func mapRoundSubmit(err error) error { return fmt.Errorf("round submit: %v", err) }
func mapRoundStorage(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("storage: %v", err)
}

var _ = storage.ErrNotFound
