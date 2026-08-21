package application

import (
	"context"
	"errors"
	"testing"
	"github.com/example/federated-learning-coordinator/internal/model"
	"github.com/example/federated-learning-coordinator/internal/privacy"
)

type failingAggregator struct{}
var primaryErr = errors.New("primary aggregate failure")
func (failingAggregator) Name() string { return "failing" }
func (failingAggregator) Aggregate([]*model.Update) ([]model.Layer, error) { return nil, primaryErr }
func TestBatchAggregationKeepsPrimaryError(t *testing.T) {
	_, _, err := NewService(failingAggregator{}).Run(context.Background(), "r", []*model.Update{{ID: "u", Layers: []model.Layer{{Name: "w", DType: "float32", Shape: []int{1}, Values: []float64{1}}}}})
	if err == nil || !errors.Is(err, primaryErr) { t.Fatalf("primary error lost: %v", err) }
	a := privacy.Accountant{EpsilonLimit: 1}; if !errors.Is(a.Charge(0), privacy.ErrInvalidPrivacyCharge) { t.Fatal("invalid charge sentinel lost") }; if !errors.Is(a.Charge(2), privacy.ErrBudgetExhausted) { t.Fatal("budget sentinel lost") }
	if a.Remaining() < 0 { t.Fatal("remaining budget went negative") }
}
