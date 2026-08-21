package adapter

import (
	"fmt"
	"github.com/example/federated-learning-coordinator/internal/cohort/domain"
	"strings"
)

func ParseLines(lines []string) (domain.Cohort, error) {
	c := domain.Cohort{Status: domain.Enabled}
	for _, line := range lines {
		p := strings.SplitN(line, ":", 2)
		if len(p) != 2 {
			continue
		}
		k, v := strings.TrimSpace(p[0]), strings.TrimSpace(p[1])
		switch k {
		case "id":
			c.ID = v
		case "name":
			c.Name = v
		case "model_version":
			c.ModelVersion = v
		}
	}
	if err := c.Validate(); err != nil {
		return c, fmt.Errorf("cohort config: %w", err)
	}
	return c, nil
}
func Format(c domain.Cohort) []string {
	return []string{"id: " + c.ID, "name: " + c.Name, "model_version: " + c.ModelVersion}
}
