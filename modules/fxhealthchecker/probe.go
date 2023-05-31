package fxhealthchecker

import (
	"context"
	"go.uber.org/fx"
)

type Probe interface {
	Name() string
	Check(ctx context.Context) *ProbeResult
}

type ProbeResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewProbeResult(success bool, message string) *ProbeResult {
	return &ProbeResult{
		Success: success,
		Message: message,
	}
}

func AsProbe(p any) any {
	return fx.Annotate(
		p,
		fx.As(new(Probe)),
		fx.ResultTags(`group:"health-checker-probes"`),
	)
}
