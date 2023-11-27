package gorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	DB          = gorm.DB          //nolint
	Model       = gorm.Model       //nolint
	Association = gorm.Association //nolint
	DeletedAt   = gorm.DeletedAt   //nolint
)

// LogLevel 日志等级
type LogLevel = logger.LogLevel

const (
	// LogLevelSilent is logger silent level
	LogLevelSilent = logger.Silent
	// LogLevelError is logger error level
	LogLevelError = logger.Error
	// LogLevelWarn is logger warn level
	LogLevelWarn = logger.Warn
	// LogLevelInfo is logger info level
	LogLevelInfo = logger.Info
)
