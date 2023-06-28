package fxtracer

import (
	"strings"
)

type Exporter int

const (
	Noop Exporter = iota
	Memory
	Stdout
	OtlpGrpc
)

func (e Exporter) String() string {
	switch e {
	case Memory:
		return "memory"
	case Stdout:
		return "stdout"
	case OtlpGrpc:
		return "otlp-grpc"
	default:
		return "noop"
	}
}

func GetExporter(exporter string) Exporter {
	switch strings.ToLower(exporter) {
	case "noop":
		return Noop
	case "memory":
		return Memory
	case "stdout":
		return Stdout
	case "otlp-grpc":
		return OtlpGrpc
	default:
		return Noop
	}
}
