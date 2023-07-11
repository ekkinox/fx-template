package fxloggertest

import "testing"

func AssertHasLogRecord(t testing.TB, expectedAttributes map[string]interface{}) bool {

	hasRecord, err := GetTestLogBufferInstance().HasRecord(expectedAttributes)
	if err != nil {
		t.Errorf("error while asserting log record attributes: %v", err)

		return false
	}

	if !hasRecord {
		t.Errorf("cannot find log record with attributes %+v", expectedAttributes)

		return false
	}

	return true
}
