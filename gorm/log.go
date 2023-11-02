package gorm

import (
	"time"

	"gorm.io/gorm/logger"

	"github.com/aixj1984/golibs/zlog"
)

type LoggerFunc func(...interface{})

func (f LoggerFunc) Print(args ...interface{}) { f(args...) }

type Writer struct {
}

func (w Writer) Printf(format string, args ...interface{}) {
	zlog.Infof(format, args...)
}

func (db *Engine) wrapLog() {
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
