package adapter

import (
	"encoding/json"
	"github.com/example/federated-learning-coordinator/internal/worker/domain"
)

func Encode(j domain.Job) ([]byte, error) { return json.Marshal(j) }
func Decode(b []byte) (domain.Job, error) { var j domain.Job; e := json.Unmarshal(b, &j); return j, e }
