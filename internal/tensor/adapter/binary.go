package adapter

import (
	"encoding/binary"
	"fmt"
)

func EncodeFloat64s(values []float64) []byte {
	b := make([]byte, 8*len(values))
	for i, v := range values {
		binary.BigEndian.PutUint64(b[i*8:], uint64(v))
	}
	return b
}
func DecodeFloat64s(b []byte) ([]float64, error) {
	if len(b)%8 != 0 {
		return nil, fmt.Errorf("float payload alignment")
	}
	out := make([]float64, len(b)/8)
	for i := range out {
		out[i] = float64(binary.BigEndian.Uint64(b[i*8:]))
	}
	return out, nil
}
func Size(values []float64) int { return len(values) * 8 }
