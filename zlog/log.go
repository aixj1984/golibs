package zlog

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	viper "github.com/aixj1984/golibs/conf"
)

// Config 是log文件的参数配置
type Config struct {
	LogPath    string `mapstructure:"logPath" yaml:"logPath" json:"logPath" comment:"日志文件路径"`
	AppName    string `mapstructure:"appName" yaml:"appName" json:"appName" comment:"应用名称"`
	Debug      bool   `mapstructure:"debug" yaml:"debug" json:"debug" comment:"是否开启调试模式"`
	Level      int8   `mapstructure:"level" yaml:"level" json:"level" comment:"日志级别"`
	MaxSize    int    `mapstructure:"maxSize" yaml:"maxSize" json:"maxSize" comment:"每个日志文件保存的大小 单位:M"`
	MaxAge     int    `mapstructure:"maxAge" yaml:"maxAge" json:"maxAge" comment:"文件最多保存多少天"`
	MaxBackups int    `mapstructure:"maxBackups" yaml:"maxBackups" json:"maxBackups" comment:"日志文件最多保存多少个备份"`
	Compress   bool   `mapstructure:"compress" yaml:"compress" json:"compress" comment:"是否压缩"`
}

var (
	mLog  *Entry
	mConf *Config
)

func init() {
	logCfg, err := viper.GetSubCfg[Config]("log")
	if err != nil {
		fmt.Printf("unable to get config, %s\n", err.Error())
		return
	}
	InitLogger(logCfg)
}

// InitLogger 通过传入的config，来初始化日志对象
func InitLogger(config *Config) {
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
	writes := []zapcore.WriteSyncer{zapcore.AddSync(&hook)}

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
	mLog = NewEntry(zap.New(core, caller, callerSkip, development, field).WithOptions(zap.AddCallerSkip(1)))

	mConf = config
}

// Empty 是将当前的日志对象设置为null
func Empty() bool {
	return mLog == nil
}

// Logger 获取当前的日志对象
func Logger() *Entry {
	// 不建议再获取对象，除非要调用zap的特殊功能
	return mLog
}

// GetConfig 获取当前的配置文件信息
func GetConfig() *Config {
	return mConf
}

// Debug 输出debug级别的日志
func Debug(msg string, fields Fields) {
	mLog.Debug(msg, fields)
}

// DebugO 输出debug级别的任意对象到日志
func DebugO(msg string, object interface{}) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, object)
		return
	}
	mLog.Debug(msg, object)
}

// Debugf 输出debug级别的format日志
func Debugf(format string, args ...interface{}) {
	if mLog == nil {
		fmt.Printf(format+"\n", args...)
		return
	}
	mLog.Debug("", Fields{"content": fmt.Sprintf(format, args...)})
}

// Info 输出info级别的日志
func Info(msg string, fields Fields) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, fields)
		return
	}
	mLog.Info(msg, fields)
}

// InfoO 输出info级别的任意对象到日志
func InfoO(msg string, object interface{}) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, object)
		return
	}
	mLog.Info(msg, object)
}

// Infof 输出info级别的format日志
func Infof(format string, args ...interface{}) {
	if mLog == nil {
		fmt.Printf(format+"\n", args...)
		return
	}
	mLog.Info("", Fields{"content": fmt.Sprintf(format, args...)})
}

// Warn 输出warn级别的日志
func Warn(msg string, fields Fields) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, fields)
		return
	}
	mLog.Warn(msg, fields)
}

// WarnO 输出warn级别的任意对象到日志
func WarnO(msg string, object interface{}) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, object)
		return
	}
	mLog.Warn(msg, object)
}

// Warnf 输出warn级别的format日志
func Warnf(format string, args ...interface{}) {
	mLog.Warn("", Fields{"content": fmt.Sprintf(format, args...)})
}

// Error 输出error级别的日志
func Error(msg string, fields Fields) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, fields)
		return
	}
	mLog.Error(msg, fields)
}

// ErrorO 输出error级别的任意对象到日志
func ErrorO(msg string, object interface{}) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, object)
		return
	}
	mLog.Error(msg, object)
}

// Errorf 输出error级别的format日志
func Errorf(format string, args ...interface{}) {
	if mLog == nil {
		fmt.Printf(format+"\n", args...)
		return
	}
	mLog.Error("", Fields{"content": fmt.Sprintf(format, args...)})
}

// DPanic 输出dpanic级别的日志,同时进程退出
func DPanic(msg string, fields Fields) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, fields)
		return
	}
	mLog.DPanic(msg, fields)
}

// DPanicO 输出dpanic级别的任意对象到日志,同时进程退出
func DPanicO(msg string, object interface{}) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, object)
		return
	}
	mLog.DPanic(msg, object)
}

// DPanicf 输出dpanic级别的format日志,同时进程退出
func DPanicf(format string, args ...interface{}) {
	if mLog == nil {
		fmt.Printf(format+"\n", args...)
		return
	}
	mLog.DPanic("", Fields{"content": fmt.Sprintf(format, args...)})
}

// Panic 输出fatal级别的日志,同时进程退出
func Panic(msg string, fields Fields) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, fields)
		return
	}
	mLog.Panic(msg, fields)
}

// PanicO 输出panic级别的任意对象到日志,同时进程退出
func PanicO(msg string, object interface{}) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, object)
		return
	}
	mLog.Panic(msg, object)
}

// Panicf 输出panic级别的format日志,同时进程退出
func Panicf(format string, args ...interface{}) {
	if mLog == nil {
		fmt.Printf(format+"\n", args...)
		return
	}
	mLog.Panic("", Fields{"content": fmt.Sprintf(format, args...)})
}

// Fatal 输出fatal级别的日志,同时进程退出
func Fatal(msg string, fields Fields) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, fields)
		return
	}
	mLog.Fatal(msg, fields)
}

// FatalO 输出fatal级别的任意对象到日志,同时进程退出
func FatalO(msg string, object interface{}) {
	if mLog == nil {
		fmt.Printf("%s : %+v\n", msg, object)
		return
	}
	mLog.Fatal(msg, object)
}

// Fatalf 输出fatal级别的format日志,同时进程退出
func Fatalf(format string, args ...interface{}) {
	if mLog == nil {
		fmt.Printf(format+"\n", args...)
		return
	}
	mLog.Fatal("", Fields{"content": fmt.Sprintf(format, args...)})
}
