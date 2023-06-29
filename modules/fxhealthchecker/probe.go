package fxhealthchecker

import (
	"context"
)

type HealthCheckerProbe interface {
	Name() string
	Check(ctx context.Context) *HealthCheckerProbeResult
}

type HealthCheckerProbeResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewHealthCheckerProbeResult(success bool, message string) *HealthCheckerProbeResult {
	return &HealthCheckerProbeResult{
		Success: success,
		Message: message,
	}
}
