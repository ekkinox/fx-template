package fxtracertest

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func AssertHasTraceSpan(t testing.TB, expectedAttributes ...attribute.KeyValue) bool {
	if !GetTestTraceExporterInstance().HasSpan(expectedAttributes) {
		t.Errorf("cannot find trace span with attributes %v", expectedAttributes)

		return false
	}

	return true
}
