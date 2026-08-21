package domain

import (
	"errors"
	"time"
)

type Status string

const (
	Queued    Status = "queued"
	Running   Status = "running"
	Succeeded Status = "succeeded"
	Failed    Status = "failed"
)

type Job struct {
	ID, Kind    string
	Status      Status
	Attempts    int
	MaxAttempts int
	NextRun     time.Time
	LastError   string
	CreatedAt   time.Time
}

func (j Job) Validate() error {
	if j.ID == "" || j.Kind == "" {
		return errors.New("job identity required")
	}
	if j.MaxAttempts < 1 {
		return errors.New("max attempts required")
	}
	return nil
}
func (j *Job) Retry(now time.Time, err error) {
	j.Attempts++
	j.LastError = err.Error()
	if j.Attempts >= j.MaxAttempts {
		j.Status = Failed
		return
	}
	j.Status = Queued
	j.NextRun = now.Add(time.Duration(j.Attempts*j.Attempts) * time.Second)
}
func (j *Job) Start()    { j.Status = Running }
func (j *Job) Complete() { j.Status = Succeeded }
