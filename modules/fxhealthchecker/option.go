package fxhealthchecker

type options struct {
	Probes []HealthCheckerProbe
}

var defaultHeatchCheckerOptions = options{
	Probes: []HealthCheckerProbe{},
}

type HealthCheckerOption func(o *options)

func WithProbes(p ...HealthCheckerProbe) HealthCheckerOption {
	return func(o *options) {
		o.Probes = p
	}
}
