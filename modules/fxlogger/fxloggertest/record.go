package fxloggertest

import (
	"errors"
	"fmt"
	"time"
)

type TestLogRecord struct {
	attributes map[string]interface{}
}

func NewTestLogRecord(attributes map[string]interface{}) *TestLogRecord {
	return &TestLogRecord{
		attributes: attributes,
	}
}

func (r *TestLogRecord) GetLevel() (string, error) {
	value, err := r.GetAttribute("level")
	if err != nil {
		return "", err
	}

	return value.(string), nil
}

func (r *TestLogRecord) GetMessage() (string, error) {
	value, err := r.GetAttribute("message")
	if err != nil {
		return "", err
	}

	return value.(string), nil
}

func (r *TestLogRecord) GetService() (string, error) {
	value, err := r.GetAttribute("service")
	if err != nil {
		return "", err
	}

	return value.(string), nil
}

func (r *TestLogRecord) GetTime() (time.Time, error) {
	value, err := r.GetAttribute("time")
	if err != nil {
		return time.Now(), err
	}

	return time.Unix(value.(int64), 0), nil
}

func (r *TestLogRecord) GetAttribute(name string) (interface{}, error) {
	value, ok := r.attributes[name]
	if ok {
		return value, nil
	} else {
		return "", errors.New(fmt.Sprintf("attribute %s not found in record", name))
	}
}

func (r *TestLogRecord) MatchAttributes(expectedAttributes map[string]interface{}) (bool, error) {

	match := true

	if len(expectedAttributes) == 0 {
		return false, nil
	}

	for expectedAttributeName, expectedAttributeValue := range expectedAttributes {
		value, ok := r.attributes[expectedAttributeName]
		if ok {
			match = match && value == expectedAttributeValue
		} else {
			return false, errors.New(fmt.Sprintf("attribute %s not found in record", expectedAttributeName))
		}
	}

	return match, nil
}
