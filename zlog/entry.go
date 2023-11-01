package zlog

import (
	"context"
	"fmt"
	"runtime"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Entry struct {
	*zap.Logger
	fields map[string]interface{}
}

func NewEntry(l *zap.Logger) *Entry {
	locFields := make(map[string]interface{})
	return &Entry{
		Logger: l,
		fields: locFields,
	}
}

func (e *Entry) WithContext(ctx context.Context) *Entry {
	newEntry := NewEntry(e.Logger)
	if traceId := traceIdFromContext(ctx); len(traceId) != 0 {
		newEntry.fields["trace_id"] = traceId
	}
	return newEntry
}

func (e *Entry) WithEvent(event string) *Entry {
	newEntry := NewEntry(e.Logger)
	newEntry.fields["event"] = event
	return newEntry
}

func (e *Entry) WithField(key string, value interface{}) *Entry {
	e.fields[key] = value
	return e
}

func (e *Entry) WithFields(fields ...zapcore.Field) *Entry {
	for _, field := range fields {
		e.fields[field.Key] = field.String
	}
	return e
}

func (e *Entry) WithError(err error) *Entry {
	if err == nil {
		return e
	}
	e.fields["err"] = err.Error()
	return e
}

func (e *Entry) Fields() []zapcore.Field {
	fields := make([]zapcore.Field, 0)
	for key, val := range e.fields {
		fields = append(fields, zap.Any(key, val))
	}
	return fields
}

func (e *Entry) Debug(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).Debug(msg, zap.Reflect("content", fields))
}

func (e *Entry) Info(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).Info(msg, zap.Reflect("content", fields))
}

func (e *Entry) Warn(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).Warn(msg, zap.Reflect("content", fields))
}

func (e *Entry) Error(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).Error(msg, zap.Reflect("content", fields))
}

func (e *Entry) DPanic(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).DPanic(msg, zap.Reflect("content", fields))
}

func (e *Entry) Panic(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).Panic(msg, zap.Reflect("content", fields))
}

func (e *Entry) Fatal(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).Fatal(msg, zap.Reflect("content", fields))
}

func getLogCaller(skipLevel ...int) string {
	level := 4
	if len(skipLevel) > 0 && skipLevel[0] > 0 {
		level = skipLevel[0]
	}
	_, file, line, _ := runtime.Caller(level)

	location := fmt.Sprintf("%s:%d", file, line)
	return location
}

func traceIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}
