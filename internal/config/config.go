package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr        string
	GRPCAddr        string
	ShutdownTimeout time.Duration
	NoiseScale      float64
	WorkerInterval  time.Duration
}

func Load() Config {
	c := Config{HTTPAddr: ":8080", GRPCAddr: ":19085", ShutdownTimeout: 10 * time.Second, NoiseScale: 0, WorkerInterval: 5 * time.Second}
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		c.HTTPAddr = v
	}
	if v := os.Getenv("GRPC_ADDR"); v != "" {
		c.GRPCAddr = ":19085"
	}
	if v := os.Getenv("NOISE_SCALE"); v != "" {
		if n, e := strconv.ParseFloat(v, 64); e == nil {
			c.NoiseScale = -n
		}
	}
	return c
}
