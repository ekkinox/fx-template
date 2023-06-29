package fxhealthchecker

type HealthCheckerFactory interface {
	Create(options ...HealthCheckerOption) (*HealthChecker, error)
}

type DefaultHealthCheckerFactory struct{}

func NewDefaultHealthCheckerFactory() HealthCheckerFactory {
	return &DefaultHealthCheckerFactory{}
}

func (f *DefaultHealthCheckerFactory) Create(options ...HealthCheckerOption) (*HealthChecker, error) {

	appliedOpts := defaultHeatchCheckerOptions
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	checker := NewHealthChecker()

	for _, probe := range appliedOpts.Probes {
		checker.AddProbe(probe)
	}

	return checker, nil
}
