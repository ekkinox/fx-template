package fxtracertest

import (
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var once sync.Once
var testTraceExporter *TestTraceExporter

type TestTraceExporter struct {
	*tracetest.InMemoryExporter
}

func GetTestTraceExporterInstance() *TestTraceExporter {
	once.Do(func() {
		testTraceExporter = &TestTraceExporter{
			tracetest.NewInMemoryExporter(),
		}
	})

	return testTraceExporter
}

func (e *TestTraceExporter) GetExporter() *tracetest.InMemoryExporter {
	return e.InMemoryExporter
}

func (e *TestTraceExporter) ClearSpans() *TestTraceExporter {
	e.InMemoryExporter.Reset()

	return e
}

func (e *TestTraceExporter) GetSpans() tracetest.SpanStubs {
	return e.InMemoryExporter.GetSpans()
}

func (e *TestTraceExporter) HasSpan(expectedName string, expectedAttributes []attribute.KeyValue) bool {
	for _, span := range e.InMemoryExporter.GetSpans() {
		if span.Name == expectedName {
			for _, expectedAttribute := range expectedAttributes {
				for _, spanAttribute := range span.Attributes {
					if spanAttribute.Key == expectedAttribute.Key && spanAttribute.Value == expectedAttribute.Value {
						return true
					}
				}
			}
		}
	}

	return false
}
