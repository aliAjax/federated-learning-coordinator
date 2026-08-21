package anomaly

import (
	"github.com/example/federated-learning-coordinator/internal/model"
	"github.com/example/federated-learning-coordinator/internal/tensor"
	"math"
)

type Result struct {
	Anomalous bool    `json:"anomalous"`
	Reason    string  `json:"reason,omitempty"`
	Norm      float64 `json:"norm"`
	Cosine    float64 `json:"cosine"`
}
type Detector struct {
	MaxNorm   float64
	MinCosine float64
}

func New(max float64) *Detector { return &Detector{MaxNorm: max, MinCosine: -0.95} }
func (d *Detector) Check(layers []model.Layer, reference []model.Layer) Result {
	norm := tensor.Norm(layers)
	r := Result{Norm: norm, Cosine: 1}
	if norm > d.MaxNorm*10 {
		r.Anomalous = true
		r.Reason = "norm exceeds safety threshold"
		return r
	}
	if len(reference) > 0 && tensor.Compatible(layers, reference) {
		dot, an, bn := 0.0, 0.0, 0.0
		for i := range layers {
			for j, v := range layers[i].Values {
				q := reference[i].Values[j]
				dot += v * q
				an += v * v
				bn += q * q
			}
		}
		if an > 0 && bn > 0 {
			r.Cosine = dot / (math.Sqrt(an) * math.Sqrt(bn))
			if r.Cosine < d.MinCosine {
				r.Anomalous = true
				r.Reason = "direction outlier"
			}
		}
	}
	return r
}
