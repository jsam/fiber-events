package fiberEvents

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

const (
	EventLevelLog = iota
	EventLevelDebug
	EventLevelTrace
	EventLevelInfo
	EventLevelNotice
	EventLevelWarning
	EventLevelError
	EventLevelCritical
	EventLevelAlert
	EventLevelEmergency
)

type PlainEvent struct {
	EventType string      `json:"event_type"`
	Level     int         `json:"level"`
	Message   string      `json:"message"`
	Payload   interface{} `json:"payload"`
}

type PlainStream struct {
	level int
	out   io.Writer
}

func NewStream(out io.Writer, level int) *PlainStream {
	return &PlainStream{
		out:   out,
		level: level,
	}
}

func (stream *PlainStream) eventPlainAppend(eventType string, level int, message string, payload interface{}) {
	if stream.level > level {
		return
	}

	plain := &PlainEvent{
		EventType: eventType,
		Level:     level,
		Message:   message,
		Payload:   payload,
	}

	logStr, err := json.Marshal(plain)
	if err != nil {
		log.Println(err)
		return
	}

	record := fmt.Sprintf("%s\n", logStr)
	stream.out.Write([]byte(record))
}

func (stream *PlainStream) Emergency(message string, payload interface{}) {
	stream.eventPlainAppend("emergency", EventLevelEmergency, message, payload)
}

func (stream *PlainStream) Alert(message string, payload interface{}) {
	stream.eventPlainAppend("alert", EventLevelAlert, message, payload)
}

func (stream *PlainStream) Critical(message string, payload interface{}) {
	stream.eventPlainAppend("critical", EventLevelCritical, message, payload)
}

func (stream *PlainStream) Error(message string, payload interface{}) {
	stream.eventPlainAppend("error", EventLevelError, message, payload)
}

func (stream *PlainStream) Warning(message string, payload interface{}) {
	stream.eventPlainAppend("warning", EventLevelWarning, message, payload)
}

func (stream *PlainStream) Notice(message string, payload interface{}) {
	stream.eventPlainAppend("notice", EventLevelNotice, message, payload)
}

func (stream *PlainStream) Info(message string, payload interface{}) {
	stream.eventPlainAppend("info", EventLevelInfo, message, payload)
}

func (stream *PlainStream) Trace(message string, payload interface{}) {
	stream.eventPlainAppend("trace", EventLevelTrace, message, payload)
}

func (stream *PlainStream) Debug(message string, payload interface{}) {
	stream.eventPlainAppend("debug", EventLevelDebug, message, payload)
}

func (stream *PlainStream) Log(level string, message string, payload interface{}) {
	stream.eventPlainAppend("log", EventLevelLog, message, payload)
}
