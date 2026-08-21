package tensor_test

import (
	"github.com/example/federated-learning-coordinator/internal/aggregation"
	"github.com/example/federated-learning-coordinator/internal/model"
	"github.com/example/federated-learning-coordinator/internal/tensor"
	"testing"
)

func TestClipDoesNotPolluteOriginalLayers(t *testing.T) {
	in := []model.Layer{{Name: "w", DType: "float32", Shape: []int{2}, Values: []float64{3, 4}}}
	out, _ := tensor.Clip(in, 1)
	if out[0].Values[0] == in[0].Values[0] || in[0].Values[0] != 3 {
		t.Fatalf("input polluted: %#v", in)
	}
	if tensor.Compatible(in, []model.Layer{{Name: "w", DType: "float64", Shape: []int{2}, Values: []float64{3, 4}}}) {
		t.Fatal("dtype mismatch accepted")
	}
	if got := tensor.CloneLayers(in); &got[0].Values[0] == &in[0].Values[0] {
		t.Fatal("clone shares values")
	}
	first := &model.Update{Layers: []model.Layer{{Name: "w", DType: "float32", Shape: []int{2}, Values: []float64{1, 2}}}}
	second := &model.Update{Layers: []model.Layer{{Name: "w", DType: "float32", Shape: []int{2}, Values: []float64{3, 4}}}}
	got, err := (&aggregation.Mean{}).Aggregate([]*model.Update{first, second})
	if err != nil || len(got) != 1 || got[0].DType != "float32" || got[0].Values[0] != 2 || got[0].Values[1] != 3 {
		t.Fatalf("bad aggregate: %#v %v", got, err)
	}
	if first.Layers[0].Values[0] != 1 || first.Layers[0].Shape[0] != 2 {
		t.Fatalf("aggregate polluted source: %#v", first.Layers[0])
	}
	shape := []int{2}
	cloned := tensor.CloneLayers([]model.Layer{{Name: "shape", Shape: shape, Values: []float64{1, 2}}})
	shape[0] = 9
	if cloned[0].Shape[0] != 2 {
		t.Fatal("clone shares shape metadata")
	}
}
