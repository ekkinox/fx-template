package fxtracertest

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func AssertHasTraceSpan(t testing.TB, expectedName string, expectedAttributes ...attribute.KeyValue) bool {
	if !GetTestTraceExporterInstance().HasSpan(expectedName, expectedAttributes) {
		t.Errorf("cannot find trace span with name %s and attributes %+v", expectedName, expectedAttributes)

		return false
	}

	return true
}
