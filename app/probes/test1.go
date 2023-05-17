package probes

import (
	"fmt"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
)

type Test1Probe struct {
	config *fxconfig.Config
}

func NewTest1Probe(config *fxconfig.Config) *Test1Probe {
	return &Test1Probe{
		config: config,
	}
}

func (p *Test1Probe) Name() string {
	return "probe1"
}

func (p *Test1Probe) Check() *fxhealthchecker.ProbeResult {
	return fxhealthchecker.NewProbeResult(
		true,
		fmt.Sprintf("probe 1 - %s", p.config.AppConfig.Name),
	)
}
