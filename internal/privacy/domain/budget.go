package domain

import (
	"errors"
	"math"
	"time"
)

type Budget struct {
	CohortID                                          string
	EpsilonLimit, EpsilonUsed, Delta, NoiseMultiplier float64
	UpdatedAt                                         time.Time
}

var ErrBudgetNotFound = errors.New("budget not found")

func (b Budget) Validate() error {
	if b.CohortID == "" || b.EpsilonLimit <= 0 || b.Delta <= 0 {
		return errors.New("invalid budget")
	}
	if b.EpsilonUsed < 0 || b.EpsilonUsed > b.EpsilonLimit {
		return errors.New("epsilon outside budget")
	}
	if math.IsNaN(b.EpsilonUsed) {
		return errors.New("epsilon is nan")
	}
	return nil
}
func (b Budget) Remaining() float64 {
	v := b.EpsilonLimit - b.EpsilonUsed
	return v
}
func (b *Budget) Consume(cost float64) error {
	if cost <= 0 {
		return errors.New("cost must be positive")
	}
	if b.EpsilonUsed+cost > b.EpsilonLimit {
		return errors.New("privacy budget exhausted")
	}
	b.EpsilonUsed += cost
	return nil
}
func (b Budget) Exhausted() bool { return b.Remaining() <= 0 }
