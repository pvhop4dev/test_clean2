package http

import (
	"time"

	"clean-arch-go/internal/pkg/config"
)

type HTTPConfig struct {
	Port        string
	GRPCPort    string
	Env         string
	RateLimit   int
	RateBurst   int
	ShutdownTimeout time.Duration
}

func NewHTTPConfig(cfg *config.Config) *HTTPConfig {
	return &HTTPConfig{
		Port:        cfg.App.Port,
		GRPCPort:    cfg.App.GRPCPort,
		Env:         cfg.App.Env,
		RateLimit:   cfg.RateLimit.Limit,
		RateBurst:   cfg.RateLimit.Burst,
		ShutdownTimeout: time.Second * 5,
	}
}
