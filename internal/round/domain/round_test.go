package domain_test

import (
	"github.com/example/federated-learning-coordinator/internal/round/adapter"
	"github.com/example/federated-learning-coordinator/internal/round/domain"
	"testing"
)

func TestRoundRetryReachesCompleted(t *testing.T) {
	if !domain.Transition(domain.Aggregating, domain.Completed) {
		t.Fatal("aggregating should complete")
	}
	r := domain.Round{Status: domain.Aggregating}
	r.MarkCompleted("d")
	if r.Status != domain.Completed {
		t.Fatal("completion mark lost")
	}
	r.Status = domain.Collecting
	r.MarkAggregating()
	if r.Status != domain.Collecting {
		t.Fatal("invalid promotion accepted")
	}
	if len(adapter.Statuses()) != 5 {
		t.Fatalf("status list incomplete: %v", adapter.Statuses())
	}
	b, _ := adapter.Encode(domain.Round{Status: domain.Completed})
	decoded, _ := adapter.Decode(b)
	if decoded.Status != domain.Completed {
		t.Fatalf("completed status lost in adapter: %v", decoded.Status)
	}
	r.Status = domain.Aborted
	r.AggregateDigest = ""
	r.MarkCompleted("ignored")
	if r.Status != domain.Aborted || r.AggregateDigest != "" {
		t.Fatal("aborted round changed")
	}
}
