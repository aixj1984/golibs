package zlog

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	ctx, _ = otel.Tracer("foo").Start(context.Background(), "bar")
)

func TestMain(m *testing.M) {

	InitLogger(&Config{
		LogPath:    "./log/test.log",
		AppName:    "log-sample",
		Level:      -1,
		MaxSize:    1024,
		MaxAge:     3,
		MaxBackups: 4,
		Compress:   false,
	})
	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	m.Run()
}

func TestInit(t *testing.T) {
	Info("info log", Fields{"abc": 11})
	Warnf("warn log : %d", 12)
	DebugO("debug object", Fields{"abc": 13})
}

func TestLogStruct(t *testing.T) {
	Logger().WithContext(ctx).Info("log start", nil)
	Logger().WithContext(ctx).Debug("log conf", Fields{
		"conf1": GetConfig(),
		"conf2": GetConfig(),
	})
	Debugf("config [%+v]", GetConfig())
}
