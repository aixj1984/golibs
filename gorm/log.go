package gorm

import (
	"time"

	"gorm.io/gorm/logger"

	"github.com/aixj1984/golibs/zlog"
)

// Writer 重新定义gorm的writer类
type Writer struct{}

// Printf 是gorm日志输出的实现
func (w Writer) Printf(format string, args ...interface{}) {
	zlog.Infof(format, args...)
}

// WrapLog 更新gorm的log实现
func (db *Engine) WrapLog() {
	if zlog.Logger() == nil {
		return
	}

	newLogger := logger.New(
		Writer{},
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Info,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                  // Disable color
		},
	)
	db.gorm.Logger = newLogger
}
