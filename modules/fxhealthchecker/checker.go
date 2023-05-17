package fxhealthchecker

import "context"

type CheckerResult struct {
	Success       bool                    `json:"success"`
	ProbesResults map[string]*ProbeResult `json:"results"`
}

type Checker struct {
	probes []Probe
}

func NewChecker() *Checker {
	return &Checker{
		probes: []Probe{},
	}
}

func (c *Checker) AddProbe(p Probe) *Checker {
	c.probes = append(c.probes, p)

	return c
}

func (c *Checker) Run(ctx context.Context) *CheckerResult {

	success := true
	probeResults := map[string]*ProbeResult{}

	for _, p := range c.probes {

		pr := p.Check(ctx)

		success = success && pr.Success
		probeResults[p.Name()] = pr
	}

	return &CheckerResult{
		Success:       success,
		ProbesResults: probeResults,
	}
}
