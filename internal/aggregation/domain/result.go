package domain

import (
	"errors"
	"time"
)

type Result struct {
	RoundID, Digest, Algorithm string
	Accepted, Rejected         int
	CreatedAt                  time.Time
}

func (r Result) Validate() error {
	if r.RoundID == "" || r.Digest == "" || r.Algorithm == "" {
		return errors.New("aggregation result incomplete")
	}
	if r.Accepted < 1 {
		return errors.New("no accepted updates")
	}
	return nil
}
func (r Result) AcceptanceRate() float64 {
	total := r.Accepted + r.Rejected
	if total == 0 {
		return 0
	}
	return float64(r.Accepted) / float64(total)
}
