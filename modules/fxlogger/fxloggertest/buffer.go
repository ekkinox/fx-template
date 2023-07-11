package fxloggertest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"sync"
)

var once sync.Once
var testLogBuffer *TestLogBuffer

type TestLogBuffer struct {
	bytes.Buffer
}

func GetTestLogBufferInstance() *TestLogBuffer {
	once.Do(func() {
		testLogBuffer = &TestLogBuffer{}
	})

	return testLogBuffer
}

func (b *TestLogBuffer) GetBuffer() *bytes.Buffer {
	return &b.Buffer
}

func (b *TestLogBuffer) ClearRecords() *TestLogBuffer {
	b.Buffer.Reset()

	return b
}

func (b *TestLogBuffer) GetRecords() ([]*TestLogRecord, error) {
	var records []*TestLogRecord

	scanner := bufio.NewScanner(b)

	for scanner.Scan() {
		var fields map[string]interface{}

		err := json.Unmarshal(scanner.Bytes(), &fields)
		if err != nil {
			return nil, err
		}

		records = append(records, NewTestLogRecord(fields))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (b *TestLogBuffer) HasRecord(expectedAttributes map[string]interface{}) (bool, error) {
	records, err := b.GetRecords()
	if err != nil {
		return false, err
	}

	for _, record := range records {
		match, err := record.MatchAttributes(expectedAttributes)
		if err != nil {
			return false, err
		}
		if match {
			return match, err
		}
	}

	return false, nil
}
