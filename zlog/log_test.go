package zlog

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gin-gonic/gin"
	require "github.com/stretchr/testify/require"
)

var ctx, _ = otel.Tracer("foo").Start(context.Background(), "bar")

func TestMain(m *testing.M) {
	m.Run()
}

func TestInit(t *testing.T) {
	InitLogger(&Config{
		LogPath:    "./log/test.log",
		AppName:    "log-sample",
		Level:      -1,
		MaxSize:    1024,
		MaxAge:     3,
		MaxBackups: 4,
		Compress:   false,
	})

	Empty()

	InitLogger(&Config{
		LogPath:    "./log/test.log",
		AppName:    "log-sample",
		Level:      -1,
		MaxSize:    0,
		MaxAge:     0,
		MaxBackups: 0,
		Compress:   true,
	})
}

func TestLog(t *testing.T) {
	Debug("info log", Fields{"abc": 11})
	Info("info log", Fields{"abc": 11})
	Warn("info log", Fields{"abc": 11})
	Error("info log", Fields{"abc": 11})

	Debugf("format log : %d", 12)
	Infof("format log : %d", 12)
	Warnf("format log : %d", 12)
	Errorf("format log : %d", 12)

	DebugO("object log ", time.Now())
	InfoO("object log ", time.Now())
	WarnO("object log ", time.Now())
	ErrorO("object log ", time.Now())
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("A panic was recovered", r)
		}
	}()
	Panic("info log", Fields{"abc": 11})
}

func TestPanicf(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("A panic was recovered", r)
		}
	}()
	Panicf("info log : %s ", "tttt")
}

func TestPanicO(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("A panic was recovered", r)
		}
	}()
	PanicO("info log", Fields{"abc": 11})
}

func TestDPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("A panic was recovered", r)
		}
	}()
	DPanic("info log", Fields{"abc": 11})
}

func TestDPanicf(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("A panic was recovered", r)
		}
	}()
	DPanicf("info log : %s ", "tttt")
}

func TestDPanicO(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("A panic was recovered", r)
		}
	}()
	DPanicO("info log", Fields{"abc": 11})
}

func TestLogStruct(t *testing.T) {
	Logger().WithContext(ctx).Info("log start", nil)
	Logger().WithContext(ctx).Debug("log conf", Fields{
		"conf1": GetConfig(),
		"conf2": GetConfig(),
	})
	Debugf("config [%+v]", GetConfig())
}

func TestLogWith(t *testing.T) {
	Logger().WithContext(ctx).WithError(errors.New("test"))
	Logger().WithEvent("test event").WithFields(zapcore.Field{}).WithField("abc", 123)
	Logger().WithContext(ctx).Debug("log conf", Fields{
		"conf1": GetConfig(),
		"conf2": GetConfig(),
	})
	Debugf("config [%+v]", GetConfig())
}

func TestGinLogger(t *testing.T) {
	// 初始化gin引擎
	router := gin.Default()

	// 注册GinLogger中间件
	router.Use((&Entry{Logger: zap.NewExample()}).GinLogger())

	// 创建一个路由用于测试
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 创建一个请求
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	require.NoError(t, err)

	// 使用httptest创建一个ResponseRecorder
	w := httptest.NewRecorder()

	// 使用router去处理这个请求
	router.ServeHTTP(w, req)

	// 检查状态码是否正确
	require.Equal(t, http.StatusOK, w.Code)

	// 检查响应体是否正确
	require.Equal(t, "ok", w.Body.String())
}

func TestGinRecovery(t *testing.T) {
	buffer := new(bytes.Buffer)
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(buffer),
		zap.InfoLevel,
	))

	// 初始化gin引擎
	router := gin.Default()

	// 注册GinRecovery中间件
	router.Use((&Entry{Logger: logger}).GinRecovery(true))

	// 创建一个路由用于测试
	router.GET("/test", func(c *gin.Context) {
		panic("test panic")
	})

	// 创建一个请求
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	require.NoError(t, err)

	// 使用httptest创建一个ResponseRecorder
	w := httptest.NewRecorder()

	// 使用router去处理这个请求
	router.ServeHTTP(w, req)

	// 检查状态码是否正确
	require.Equal(t, http.StatusInternalServerError, w.Code)

	// 检查日志中是否包含了panic的信息
	require.Contains(t, buffer.String(), "test panic")
}
