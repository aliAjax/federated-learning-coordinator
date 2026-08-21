package domain

import (
	"errors"
	"math"
)

type Rule struct{ MaxNorm, MaxDirection float64 }

func (r Rule) Validate() error {
	if r.MaxNorm <= 0 || r.MaxDirection <= 0 || math.IsNaN(r.MaxNorm) {
		return errors.New("invalid anomaly rule")
	}
	return nil
}
func (r Rule) Check(norm, direction float64) bool {
	return norm > r.MaxNorm || direction > r.MaxDirection
}
