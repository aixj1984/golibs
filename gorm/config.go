package gorm

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Config 是gorm的配置文件字段定义
type Config struct {
	Alias        string        `mapstructure:"alias" json:"alias" yaml:"alias" comment:"数据库别名"`
	Driver       string        `mapstructure:"driver" json:"driver" yaml:"driver" comment:"数据库驱动"`
	Server       string        `mapstructure:"server" json:"server" yaml:"server" comment:"数据库服务器地址"`
	Port         int           `mapstructure:"port" json:"port" yaml:"port" comment:"数据库端口"`
	Database     string        `mapstructure:"database" json:"database" yaml:"database" comment:"数据库名称"`
	User         string        `mapstructure:"user" json:"user" yaml:"user" comment:"数据库用户名"`
	Password     string        `mapstructure:"password" json:"password" yaml:"password" comment:"数据库密码"`
	MaxIdleConns int           `mapstructure:"maxIdleConns" json:"maxIdleConns" yaml:"maxIdleConns" comment:"最大空闲连接数"`
	MaxOpenConns int           `mapstructure:"maxOpenConns" json:"maxOpenConns" yaml:"maxOpenConns" comment:"最大打开连接数"`
	Charset      string        `mapstructure:"charset" json:"charset" yaml:"charset" comment:"字符集"`
	TimeZone     string        `mapstructure:"timezone" json:"timezone" yaml:"timezone" comment:"时区"`
	MaxLeftTime  time.Duration `mapstructure:"maxLeftTime" json:"maxLeftTime" yaml:"maxLeftTime" comment:"最大连接时间 0h20m30s"`
}

func authConfig(conf *Config) (err error) {
	if len(conf.Driver) == 0 {
		conf.Driver = defaultDatabase
	}
	if conf.Port == 0 {
		conf.Port = defaultPort
	}
	if conf.Driver != "sqlite" && (len(conf.User) == 0 || len(conf.Password) == 0) {
		err = errors.New("User or  Password is empty")
		return
	}

	if conf.Driver != "sqlite" && len(conf.Server) == 0 {
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
