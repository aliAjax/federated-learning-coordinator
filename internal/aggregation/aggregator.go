package aggregation

import (
	"errors"
	"github.com/example/federated-learning-coordinator/internal/model"
	"github.com/example/federated-learning-coordinator/internal/privacy"
	"github.com/example/federated-learning-coordinator/internal/tensor"
)

var ErrIncompatible = errors.New("incompatible updates")

type Aggregator interface {
	Aggregate([]*model.Update) ([]model.Layer, error)
	Name() string
}
type Mean struct {
	Noise      privacy.Noise
	NoiseScale float64
}

func (m *Mean) Name() string { return "secure_mean" }
func (m *Mean) Aggregate(updates []*model.Update) ([]model.Layer, error) {
	if len(updates) == 0 {
		return nil, ErrIncompatible
	}
	base := updates[0].Layers
	for _, u := range updates[1:] {
		if !tensor.Compatible(base, u.Layers) {
			return nil, ErrIncompatible
		}
	}
	out := make([]model.Layer, len(base))
	for i, l := range base {
		out[i] = model.Layer{Name: l.Name, Shape: append([]int{}, l.Shape...), DType: l.DType, Values: make([]float64, len(l.Values))}
		for _, u := range updates {
			for j, v := range u.Layers[i].Values {
				out[i].Values[j] += v
			}
		}
		for j := range out[i].Values {
			out[i].Values[j] /= float64(len(updates))
		}
		if m.Noise != nil {
			out[i].Values = m.Noise.Apply(out[i].Values, m.NoiseScale)
		}
	}
	return out, nil
}
