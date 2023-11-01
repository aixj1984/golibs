package zlog

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	viper "github.com/aixj1984/golibs/conf"
)

type Config struct {
	LogPath    string `yaml:"logPath" json:"logPath"`
	AppName    string `yaml:"appName" json:"appName"`
	Debug      bool   `yaml:"debug" json:"debug"`
	Level      int8   `yaml:"level" json:"level"`
	MaxSize    int    `yaml:"maxSize" json:"maxSize"`
	MaxAge     int    `yaml:"maxAge" json:"maxAge"`
	MaxBackups int    `yaml:"maxBackups" json:"maxBackups"`
	Compress   bool   `yaml:"compress" json:"compress"`
}

var (
	h    *Entry
	conf *Config
)

func init() {

	logCfg, err := viper.GetSubCfg[Config]("log")
	if err != nil {
		fmt.Printf("unable to get config, %s\n", err.Error())
		return
	} else {
		InitLogger(logCfg)
	}
}

func InitLogger(config *Config) {
	if h != nil {
		return
	}

	conf = config
	if config.MaxSize < 1 {
		config.MaxSize = 1
	}
	if config.MaxAge < 1 {
		config.MaxAge = 1
	}
	if config.MaxBackups < 1 {
		config.MaxBackups = 1
	}
	hook := lumberjack.Logger{
		Filename:   config.LogPath,    // 日志文件路径
		MaxSize:    config.MaxSize,    // 每个日志文件保存的大小 单位:M
		MaxAge:     config.MaxAge,     // 文件最多保存多少天
		MaxBackups: config.MaxBackups, // 日志文件最多保存多少个备份
		Compress:   false,             // 是否压缩
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "file",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()

	atomicLevel.SetLevel(zapcore.Level(config.Level))
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	// 如果是开发环境，同时在控制台上也输出
	if config.Debug {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writes...), atomicLevel),
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(1)
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	field := zap.Fields(zap.String("appName", config.AppName))
	// 构造日志
	h = NewEntry(zap.New(core, caller, callerSkip, development, field).WithOptions(zap.AddCallerSkip(1)))
}

func Empty() bool {
	return h == nil
}

// 不建议再获取对象，除非要调用zap的特殊功能
func Logger() *Entry {
	return h
}

func GetConfig() *Config {
	return conf
}

func Debug(msg string, fields Fields) {
	h.Debug(msg, fields)
}

func DebugO(msg string, object interface{}) {
	h.Debug(msg, object)
}

func Debugf(format string, args ...interface{}) {
	h.Debug("", Fields{"content": fmt.Sprintf(format, args...)})
}

func Info(msg string, fields Fields) {
	h.Info(msg, fields)
}

func InfoO(msg string, object interface{}) {
	h.Info(msg, object)
}

func Infof(format string, args ...interface{}) {
	h.Info("", Fields{"content": fmt.Sprintf(format, args...)})
}

func Warn(msg string, fields Fields) {
	h.Warn(msg, fields)
}

func WarnO(msg string, object interface{}) {
	h.Warn(msg, object)
}

func Warnf(format string, args ...interface{}) {
	h.Warn("", Fields{"content": fmt.Sprintf(format, args...)})
}

func Error(msg string, fields Fields) {
	h.Error(msg, fields)
}

func ErrorO(msg string, object interface{}) {
	h.Error(msg, object)
}

func Errorf(format string, args ...interface{}) {
	h.Error("", Fields{"content": fmt.Sprintf(format, args...)})
}

func DPanic(msg string, fields Fields) {
	h.DPanic(msg, fields)
}

func DPanicO(msg string, object interface{}) {
	h.DPanic(msg, object)
}

func DPanicf(format string, args ...interface{}) {
	h.DPanic("", Fields{"content": fmt.Sprintf(format, args...)})
}

func Panic(msg string, fields Fields) {
	h.Panic(msg, fields)
}

func PanicO(msg string, object interface{}) {
	h.Panic(msg, object)
}

func Panicf(format string, args ...interface{}) {
	h.Panic("", Fields{"content": fmt.Sprintf(format, args...)})
}

func Fatal(msg string, fields Fields) {
	h.Fatal(msg, fields)
}

func FatalO(msg string, object interface{}) {
	h.Fatal(msg, object)
}

func Fatalf(format string, args ...interface{}) {
	h.Fatal("", Fields{"content": fmt.Sprintf(format, args...)})
}
