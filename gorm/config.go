package gorm

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Alias        string        `yaml:"alias" json:"alias"`
	Driver       string        `yaml:"driver" json:"driver"`
	Server       string        `yaml:"server" json:"server"`
	Port         int           `yaml:"port" json:"port"`
	Database     string        `yaml:"database" json:"database"`
	User         string        `yaml:"user" json:"user"`
	Password     string        `yaml:"password" json:"password"`
	MaxIdleConns int           `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns int           `yaml:"maxOpenConns" json:"maxOpenConns"`
	Charset      string        `yaml:"charset" json:"charset"`
	TimeZone     string        `yaml:"timezone" json:"timezone"`
	MaxLeftTime  time.Duration `yaml:"maxLeftTime" json:"maxLeftTime"`
}

func authConfig(conf *Config) (err error) {

	if len(conf.Driver) == 0 {
		conf.Driver = defaultDatabase
	}
	if conf.Port == 0 {
		conf.Port = defaultPort
	}
	if len(conf.User) == 0 || len(conf.Password) == 0 {
		err = fmt.Errorf("User or  Password is empty")
		return
	}
	if len(conf.Server) == 0 {
		err = fmt.Errorf("server addr is empty")
		return
	}
	if len(conf.Database) == 0 {
		err = fmt.Errorf("database is empty")
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
