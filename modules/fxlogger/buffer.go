package fxlogger

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
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

func (b *TestLogBuffer) AsString() string {
	return b.Buffer.String()
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

func (b *TestLogBuffer) HasRecord(level string, containedMessage string) bool {
	records, err := b.GetRecords()
	if err != nil {
		return false
	}

	for _, record := range records {
		recordLevel, _ := record.GetLevel()
		recordMessage, _ := record.GetMessage()

		if recordLevel == level && strings.Contains(recordMessage, containedMessage) {
			return true
		}
	}

	return false
}
