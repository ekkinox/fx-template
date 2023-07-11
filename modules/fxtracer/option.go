package fxtracer

type options struct {
	Name      string
	Collector string
	Exporter  Exporter
}

var defaultTracerProviderOptions = options{
	Name:      "default",
	Collector: "",
	Exporter:  Noop,
}

type TracerProviderOption func(o *options)

func WithName(n string) TracerProviderOption {
	return func(o *options) {
		o.Name = n
	}
}

func WithCollector(c string) TracerProviderOption {
	return func(o *options) {
		o.Collector = c
	}
}

func WithExporter(e Exporter) TracerProviderOption {
	return func(o *options) {
		o.Exporter = e
	}
}
