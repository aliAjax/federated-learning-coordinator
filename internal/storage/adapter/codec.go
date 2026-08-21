package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/model"
)

func EncodeUpdate(u *model.Update) ([]byte, error) { return json.Marshal(u) }
func DecodeUpdate(b []byte) (*model.Update, error) {
	var u model.Update
	if err := json.Unmarshal(b, &u); err != nil {
		return nil, fmt.Errorf("update payload: %w", err)
	}
	return &u, nil
}
func EncodeLayers(l []model.Layer) ([]byte, error) { return json.Marshal(l) }
func DecodeLayers(b []byte) ([]model.Layer, error) {
	var l []model.Layer
	e := json.Unmarshal(b, &l)
	return l, e
}
