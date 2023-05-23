package probes

import (
	"context"
	"fmt"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
)

type Test2Probe struct {
	config *fxconfig.Config
}

func NewTest2Probe(config *fxconfig.Config) *Test2Probe {
	return &Test2Probe{
		config: config,
	}
}

func (p *Test2Probe) Name() string {
	return "probe2"
}

func (p *Test2Probe) Check(ctx context.Context) *fxhealthchecker.ProbeResult {
	return fxhealthchecker.NewProbeResult(
		true,
		fmt.Sprintf("success test probe 2 - %s", p.config.AppName()),
	)
}
