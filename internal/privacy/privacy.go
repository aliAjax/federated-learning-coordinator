package privacy

import (
	"errors"
	"math"
	"math/rand"
	"sync"
)

var ErrBudgetExhausted = errors.New("privacy budget exhausted")

type Noise interface {
	Apply([]float64, float64) []float64
	Name() string
}
type Gaussian struct {
	mu sync.Mutex
	r  *rand.Rand
}

func NewGaussian(seed int64) *Gaussian { return &Gaussian{r: rand.New(rand.NewSource(seed))} }
func (g *Gaussian) Name() string       { return "gaussian" }
func (g *Gaussian) Apply(v []float64, s float64) []float64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]float64, len(v))
	for i, x := range v {
		out[i] = x + g.r.NormFloat64()*s
	}
	return out
}

type Laplace struct {
	mu sync.Mutex
	r  *rand.Rand
}

func NewLaplace(seed int64) *Laplace { return &Laplace{r: rand.New(rand.NewSource(seed))} }
func (l *Laplace) Name() string      { return "laplace" }
func (l *Laplace) Apply(v []float64, s float64) []float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]float64, len(v))
	for i, x := range v {
		u := l.r.Float64() - .5
		out[i] = x - s*math.Copysign(1, u)*math.Log(1-2*math.Abs(u))
	}
	return out
}

type Accountant struct{ EpsilonLimit, EpsilonUsed, Delta float64 }

func (a *Accountant) Charge(cost float64) error {
	if cost <= 0 {
		return errors.New("invalid privacy charge")
	}
	if a.EpsilonUsed+cost > a.EpsilonLimit {
		return ErrBudgetExhausted
	}
	a.EpsilonUsed += cost
	return nil
}
func (a Accountant) Remaining() float64 { return math.Max(0, a.EpsilonLimit-a.EpsilonUsed) }
