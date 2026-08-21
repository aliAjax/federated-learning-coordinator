package application

import (
	"context"
	"errors"
	"fmt"
	base "github.com/example/federated-learning-coordinator/internal/aggregation"
	"github.com/example/federated-learning-coordinator/internal/aggregation/domain"
	"github.com/example/federated-learning-coordinator/internal/model"
)

type Service struct{ aggregator base.Aggregator }

func NewService(a base.Aggregator) *Service { return &Service{aggregator: a} }
func (s *Service) Run(ctx context.Context, round string, updates []*model.Update) (result domain.Result, layers []model.Layer, err error) {
	select {
	case <-ctx.Done():
		return domain.Result{}, nil, ctx.Err()
	default:
	}
	accepted := 0
	rejected := 0
	filtered := []*model.Update{}
	for _, u := range updates {
		if u.Anomaly {
			rejected++
			continue
		}
		accepted++
		filtered = append(filtered, u)
	}
	var e error
	layers, e = s.aggregator.Aggregate(filtered)
	if e != nil {
		return domain.Result{}, nil, fmt.Errorf("aggregate updates: %w", e)
	}
	result = domain.Result{RoundID: round, Digest: "pending", Algorithm: s.aggregator.Name(), Accepted: accepted, Rejected: rejected}
	if accepted == 0 {
		return domain.Result{}, nil, errors.New("no accepted updates")
	}
	return result, layers, nil
}
