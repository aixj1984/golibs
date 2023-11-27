// Package zlog is a wrapper for zap.
package zlog

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Entry 是对zap的一个封装
type Entry struct {
	*zap.Logger
	fields map[string]interface{}
}

// NewEntry 通过传入zap的logger对象，构造一个entry的对象
func NewEntry(l *zap.Logger) *Entry {
	locFields := make(map[string]interface{})

	return &Entry{
		Logger: l,
		fields: locFields,
	}
}

// WithContext 通过上下文获取跟踪ID的信息，构造一个实例
func (e *Entry) WithContext(ctx context.Context) *Entry {
	newEntry := NewEntry(e.Logger)
	if traceID := traceIDFromContext(ctx); len(traceID) != 0 {
		newEntry.fields["trace_id"] = traceID
	}

	return newEntry
}

// WithEvent 通过事件信息，构造一个实例
func (e *Entry) WithEvent(event string) *Entry {
	newEntry := NewEntry(e.Logger)
	newEntry.fields["event"] = event

	return newEntry
}

// WithField 向实例中添加日志元素
func (e *Entry) WithField(key string, value interface{}) *Entry {
	e.fields[key] = value

	return e
}

// WithFields 向实例中添加多个日志元素
func (e *Entry) WithFields(fields ...zapcore.Field) *Entry {
	for _, field := range fields {
		e.fields[field.Key] = field.String
	}

	return e
}

// WithError 向实例中添加err
func (e *Entry) WithError(err error) *Entry {
	if err == nil {
		return e
	}

	e.fields["err"] = err.Error()

	return e
}

// Fields 获取实例中的所有元素
func (e *Entry) Fields() []zapcore.Field {
	fields := make([]zapcore.Field, 0)
	for key, val := range e.fields {
		fields = append(fields, zap.Any(key, val))
	}

	return fields
}

// Debug 输出debug级别的日志
func (e *Entry) Debug(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).Debug(msg, zap.Reflect("content", fields))
}

// Info 输出info级别的日志
func (e *Entry) Info(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).Info(msg, zap.Reflect("content", fields))
}

// Warn 输出warn级别的日志
func (e *Entry) Warn(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).Warn(msg, zap.Reflect("content", fields))
}

// Error 输出error级别的日志
func (e *Entry) Error(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).Error(msg, zap.Reflect("content", fields))
}

// DPanic 输出DPanic级别的日志,同时进程退出
func (e *Entry) DPanic(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).DPanic(msg, zap.Reflect("content", fields))
}

// Panic 输出panic级别的日志,同时进程退出
func (e *Entry) Panic(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).Panic(msg, zap.Reflect("content", fields))
}

// Fatal 输出fatal级别的日志,同时进程退出
func (e *Entry) Fatal(msg string, fields interface{}) {
	e.Logger.With(e.Fields()...).Fatal(msg, zap.Reflect("content", fields))
}

/*
func getLogCaller(skipLevel ...int) string {
	level := 4
	if len(skipLevel) > 0 && skipLevel[0] > 0 {
		level = skipLevel[0]
	}
	_, file, line, _ := runtime.Caller(level)

	location := fmt.Sprintf("%s:%d", file, line)
	return location
}
*/

// Debug 输出debug级别的日志
func traceIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}
