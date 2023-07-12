package fxhttpservertest

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func AssertRecordedResponseCode(t testing.TB, recorder *httptest.ResponseRecorder, expectedCode int) bool {
	if recorder.Code != expectedCode {
		t.Errorf("response code %d does not match expected response code %d", recorder.Code, expectedCode)

		return false
	}

	return true
}

func AssertRecordedResponseBody(t testing.TB, recorder *httptest.ResponseRecorder, expectedBody string) bool {
	if !strings.Contains(recorder.Body.String(), expectedBody) {
		t.Errorf("response body does not contain %s", expectedBody)

		return false
	}

	return true
}
