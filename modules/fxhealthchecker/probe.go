package fxhealthchecker

import "go.uber.org/fx"

type Probe interface {
	Name() string
	Check() *ProbeResult
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

func RegisterProbe(p any) any {
	return fx.Annotate(
		p,
		fx.As(new(Probe)),
		fx.ResultTags(`group:"health-checker-probes"`),
	)
}
