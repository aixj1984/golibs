package gorm

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Config 是gorm的配置文件字段定义
type Config struct {
	Alias        string        `mapstructure:"alias" json:"alias"`
	Driver       string        `mapstructure:"driver" json:"driver"`
	Server       string        `mapstructure:"server" json:"server"`
	Port         int           `mapstructure:"port" json:"port"`
	Database     string        `mapstructure:"database" json:"database"`
	User         string        `mapstructure:"user" json:"user"`
	Password     string        `mapstructure:"password" json:"password"`
	MaxIdleConns int           `mapstructure:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns int           `mapstructure:"maxOpenConns" json:"maxOpenConns"`
	Charset      string        `mapstructure:"charset" json:"charset"`
	TimeZone     string        `mapstructure:"timezone" json:"timezone"`
	MaxLeftTime  time.Duration `mapstructure:"maxLeftTime" json:"maxLeftTime"`
}

func authConfig(conf *Config) (err error) {
	if len(conf.Driver) == 0 {
		conf.Driver = defaultDatabase
	}
	if conf.Port == 0 {
		conf.Port = defaultPort
	}
	if len(conf.User) == 0 || len(conf.Password) == 0 {
		err = errors.New("User or  Password is empty")
		return
	}

	if len(conf.Server) == 0 {
		err = errors.New("server addr is empty")
		return
	}
	if len(conf.Database) == 0 {
		err = errors.New("database is empty")
		return
	}
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = defaultMaxIdleConns
	}
	if conf.MaxLeftTime == 0 {
		conf.MaxLeftTime = defaultMaxLeftTime
	}
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = defaultMaxOpenConns
	}

	if strings.TrimSpace(conf.Charset) == "" {
		conf.Charset = defaultCharset
	}
	if strings.TrimSpace(conf.TimeZone) == "" {
		conf.TimeZone = defaultTimeZone
	}

	return
}
