package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/model/domain"
)

func Decode(b []byte) (domain.Model, error) {
	var m domain.Model
	if err := json.Unmarshal(b, &m); err != nil {
		return m, fmt.Errorf("model json: %w", err)
	}
	return m, nil
}
func Encode(m domain.Model) ([]byte, error) { return json.Marshal(m) }
func Clone(m domain.Model) domain.Model {
	m.Layers = append([]domain.Layer(nil), m.Layers...)
	return m
}
