package domain

import (
	"errors"
	"strings"
)

type Schema struct {
	Version string
	Layers  []Layer
}
type Layer struct {
	Name, DType string
	Shape       []int
}

func (s Schema) Validate() error {
	if strings.TrimSpace(s.Version) == "" || len(s.Layers) == 0 {
		return errors.New("schema requires version and layers")
	}
	seen := map[string]bool{}
	for _, l := range s.Layers {
		if l.Name == "" || seen[l.Name] || len(l.Shape) == 0 {
			return errors.New("invalid schema layer")
		}
		seen[l.Name] = true
	}
	return nil
}
func (s Schema) Has(name string) bool {
	for _, l := range s.Layers {
		if l.Name == name {
			return true
		}
	}
	return false
}
func (s Schema) Count() int {
	n := 0
	for _, l := range s.Layers {
		size := 1
		for _, d := range l.Shape {
			size *= d
		}
		n += size
	}
	return n
}
