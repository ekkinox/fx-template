package fxtracer

type options struct {
	Name      string
	Exporter  Exporter
	Collector string
}

var defaultTracerProviderOptions = options{
	Name:      "default",
	Exporter:  Noop,
	Collector: "",
}

type TracerProviderOption func(o *options)

func WithName(n string) TracerProviderOption {
	return func(o *options) {
		o.Name = n
	}
}

func WithExporter(e Exporter) TracerProviderOption {
	return func(o *options) {
		o.Exporter = e
	}
}

func WithCollector(c string) TracerProviderOption {
	return func(o *options) {
		o.Collector = c
	}
}
