package fxlogger

import (
	"errors"
	"fmt"
)

type TestLogRecord struct {
	fields map[string]interface{}
}

func NewTestLogRecord(fields map[string]interface{}) *TestLogRecord {
	return &TestLogRecord{
		fields: fields,
	}
}

func (r *TestLogRecord) GetLevel() (string, error) {
	return r.GetFieldValueAsString("level")
}

func (r *TestLogRecord) GetMessage() (string, error) {
	return r.GetFieldValueAsString("message")
}

func (r *TestLogRecord) GetFieldValueAsString(name string) (string, error) {
	value, err := r.GetFieldValue(name)
	if err != nil {
		return "", err
	}

	return value.(string), nil
}

func (r *TestLogRecord) GetFieldValueAsInt(name string) (int, error) {
	value, err := r.GetFieldValue(name)
	if err != nil {
		return 0, err
	}

	return value.(int), nil
}

func (r *TestLogRecord) GetFieldValueAsFloat32(name string) (float32, error) {
	value, err := r.GetFieldValue(name)
	if err != nil {
		return 0, err
	}

	return value.(float32), nil
}

func (r *TestLogRecord) GetFieldValue(name string) (interface{}, error) {
	value, ok := r.fields[name]
	if ok {
		return value, nil
	} else {
		return "", errors.New(fmt.Sprintf("field %s not found in record", name))
	}
}
