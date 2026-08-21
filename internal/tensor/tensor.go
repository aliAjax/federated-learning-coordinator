package tensor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/example/federated-learning-coordinator/internal/model"
	"math"
)

var ErrInvalidTensor = errors.New("invalid tensor")

func Validate(layers []model.Layer, maxValues int) error {
	if len(layers) == 0 {
		return ErrInvalidTensor
	}
	seen := map[string]bool{}
	total := 0
	for _, l := range layers {
		if l.Name == "" || seen[l.Name] || len(l.Shape) == 0 {
			return ErrInvalidTensor
		}
		seen[l.Name] = true
		if l.DType != "float32" && l.DType != "float64" && l.DType != "fixed32" && l.DType != "sparse" {
			return ErrInvalidTensor
		}
		size := 1
		for _, d := range l.Shape {
			if d <= 0 || d > 1_000_000 {
				return ErrInvalidTensor
			}
			size *= d
		}
		if size != len(l.Values) {
			return ErrInvalidTensor
		}
		for _, v := range l.Values {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return ErrInvalidTensor
			}
		}
		total += size
		if total > maxValues {
			return ErrInvalidTensor
		}
	}
	return nil
}
func Norm(layers []model.Layer) float64 {
	sum := 0.0
	for _, l := range layers {
		for _, v := range l.Values {
			sum += v * v
		}
	}
	return math.Sqrt(sum)
}
func Clip(layers []model.Layer, limit float64) ([]model.Layer, float64) {
	norm := Norm(layers)
	factor := 1.0
	if norm > limit && limit > 0 {
		factor = limit / norm
	}
	out := make([]model.Layer, len(layers))
	for i, l := range layers {
		out[i] = model.Layer{Name: l.Name, DType: l.DType, Shape: append([]int{}, l.Shape...), Values: make([]float64, len(l.Values))}
		for j, v := range l.Values {
			out[i].Values[j] = v * factor
		}
	}
	return out, norm
}
func Digest(layers []model.Layer) string {
	b, _ := json.Marshal(layers)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func Compatible(a, b []model.Layer) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Name != b[i].Name || a[i].DType != b[i].DType || len(a[i].Shape) != len(b[i].Shape) {
			return false
		}
		for j := range a[i].Shape {
			if a[i].Shape[j] != b[i].Shape[j] {
				return false
			}
		}
	}
	return true
}
func Count(l []model.Layer) int {
	n := 0
	for _, x := range l {
		n += len(x.Values)
	}
	return n
}
func Difference(a, b []model.Layer) float64 {
	if !Compatible(a, b) {
		return math.Inf(1)
	}
	sum := 0.0
	for i := range a {
		for j := range a[i].Values {
			d := a[i].Values[j] - b[i].Values[j]
			sum += d * d
		}
	}
	return math.Sqrt(sum)
}
