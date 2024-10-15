// Package gorm is a wrapper for gorm.
package gorm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	viper "github.com/aixj1984/golibs/conf"

	"github.com/aixj1984/golibs/zlog"

	// Import MySQL database driver
	// _ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/clickhouse"

	// Import PostgreSQL database driver
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"

	// Import SQLite3 database driver
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	//"gorm.io/driver/sqlite"
	"github.com/glebarez/sqlite"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"gorm.io/gorm"
)

const (
	parentSpanGormKey    = "opentracingParentSpan"
	parentSpanGormCtxKey = "opentracingParentSpanCtx"
	spanGormKey          = "opentracingSpan"
	spanDuration         = "opentracingSpanDuration"
	defaultEngine        = "default"
)

var (
	defaultDatabase     = "mysql"
	mysqlConnStrTmpl    = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s"
	pgConnStrTmpl       = "host=%s port=%s user=%s dbname=%s password=%s TimeZone=%s"
	ckConnStrTmpl       = "clickhouse://%s:%s@%s:%d/%s?dial_timeout=30s&max_execution_time=300"
	defaultMaxOpenConns = 200
	defaultMaxIdleConns = 60
	defaultMaxLeftTime  = 300 * time.Second
	defaultCharset      = "utf8mb4"
	defaultPort         = 3306
	defaultTimeZone     = "Local"
	engineMap           map[string]*Engine
)

// Engine 是gorm的一个封装类
type Engine struct {
	gorm *gorm.DB
}

func init() {
	engineMap = make(map[string]*Engine)
	dbCfg, err := viper.GetSubCfg[Config]("gorm")
	if err != nil {
		fmt.Printf("unable to get config, %s", err.Error())

		return
	}
	RegisterDataBase(defaultEngine, dbCfg)
}

// GetEngine 通过别名，获取DB的实例
func GetEngine(aliasNames ...string) *Engine {
	if len(aliasNames) == 0 {
		return engineMap[defaultEngine]
	}
	if _, ok := engineMap[aliasNames[0]]; ok {
		return engineMap[aliasNames[0]]
	}
	return nil
}

// RegisterDataBase 注册一个别名的DB
func RegisterDataBase(aliasName string, conf *Config) {
	if len(engineMap) == 0 && aliasName != defaultEngine {
		panic("please set defalut db first")
	}

	db := NewEngine(conf)
	if db != nil {
		engineMap[aliasName] = db
	}
}

// NewEngine 实例化新的Gorm实例
func NewEngine(conf *Config) *Engine {
	err := authConfig(conf)
	if err != nil {
		panic(err)
	}

	gormConf := &gorm.Config{}
	var tempDB *gorm.DB
	switch conf.Driver {
	case "mysql":
		dsn := fmt.Sprintf(mysqlConnStrTmpl,
			conf.User,
			conf.Password,
			conf.Server,
			conf.Port,
			conf.Database,
			conf.Charset,
			conf.TimeZone)

		mysqlConfig := mysql.Config{
			DriverName:                conf.Driver,
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         255,   // string 类型字段的默认长度
			SkipInitializeWithVersion: false, // 根据版本自动配置
		}

		tempDB, err = gorm.Open(mysql.New(mysqlConfig), gormConf)
		if err != nil {
			panic(err)
		}

	case "postgres":
		dsn := fmt.Sprintf(pgConnStrTmpl,
			conf.Server,
			conf.Port,
			conf.User,
			conf.Password,
			conf.Database,
			conf.TimeZone)
		pgConfig := postgres.Config{
			DriverName: conf.Driver,
			DSN:        dsn, // DSN data source name
		}
		tempDB, err = gorm.Open(postgres.New(pgConfig), gormConf)
		if err != nil {
			panic(err)
		}

	case "sqlite":
		tempDB, err = gorm.Open(sqlite.Open(conf.Database), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			panic("panic code: 155")
		}

	case "clickhouse":
		dsn := fmt.Sprintf(ckConnStrTmpl,
			conf.User,
			conf.Password,
			conf.Server,
			conf.Port,
			conf.Database)
		tempDB, err = gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}

	default:
		panic("error db driver")
	}

	fmt.Println("DB connection successful!")

	newEngine := &Engine{tempDB}
	newEngine.WrapLog()

	sqlDB, err := tempDB.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetConnMaxLifetime(conf.MaxLeftTime)
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

	addGormCallbacks(tempDB)
	return newEngine
}

// Context 设置db查询是的上下文
func (db *Engine) Context(ctx context.Context) *gorm.DB {
	parentSpan := trace.SpanFromContext(ctx)
	return db.gorm.WithContext(ctx).Set(parentSpanGormKey, parentSpan).Set(parentSpanGormCtxKey, ctx)
}

// GetDB 获取当前实例中的db对象
func (db *Engine) GetDB() *gorm.DB {
	return db.gorm
}

// SetLogMode 设置日志模式开关
func (db *Engine) SetLogMode(mode bool) {
	if !mode {
		db.gorm.Logger.LogMode(LogLevelSilent)
	}
}

// SetLogLevel 设置输出日志的级别
func (db *Engine) SetLogLevel(level LogLevel) {
	db.gorm.Logger.LogMode(level)
}

/*
func Context(ctx context.Context, aliasName string) *gorm.DB {
	if GetEngine(aliasName) == nil {
		panic(fmt.Errorf("must init gorm.New"))
	}
	parentSpan := trace.SpanFromContext(ctx)
	return GetEngine(aliasName).gorm.Set(parentSpanGormKey, parentSpan).Set(parentSpanGormCtxKey, ctx)
}*/

func addGormCallbacks(db *gorm.DB) {
	callbacks := newCallbacks()
	registerCallbacks(db, "create", callbacks)
	registerCallbacks(db, "query", callbacks)
	registerCallbacks(db, "update", callbacks)
	registerCallbacks(db, "delete", callbacks)
	registerCallbacks(db, "row_query", callbacks)
}

type callbacks struct{}

func newCallbacks() *callbacks {
	return &callbacks{}
}

func (c *callbacks) beforeCreate(scope *gorm.DB)   { c.before(scope) }
func (c *callbacks) afterCreate(scope *gorm.DB)    { c.after(scope) }
func (c *callbacks) beforeQuery(scope *gorm.DB)    { c.before(scope) }
func (c *callbacks) afterQuery(scope *gorm.DB)     { c.after(scope) }
func (c *callbacks) beforeUpdate(scope *gorm.DB)   { c.before(scope) }
func (c *callbacks) afterUpdate(scope *gorm.DB)    { c.after(scope) }
func (c *callbacks) beforeDelete(scope *gorm.DB)   { c.before(scope) }
func (c *callbacks) afterDelete(scope *gorm.DB)    { c.after(scope) }
func (c *callbacks) beforeRowQuery(scope *gorm.DB) { c.before(scope) }
func (c *callbacks) afterRowQuery(scope *gorm.DB)  { c.after(scope) }

func (c *callbacks) before(db *gorm.DB) {
	parentSpanCtx, ok := db.Get(parentSpanGormCtxKey)
	if !ok {
		// xlog.Infoln("no parentSpanCtx")
		return
	}
	db.Set(spanDuration, time.Now())
	_, span := otel.Tracer("GORM-V2-SQL").Start(parentSpanCtx.(context.Context), db.Statement.Name())
	db.Set(spanGormKey, span)
}

func (c *callbacks) after(scope *gorm.DB) {
	tempTime, ok := scope.Get(spanDuration)
	if !ok {
		tempTime = time.Now()
	}
	if span, ok := scope.Get(spanGormKey); ok {
		vars, err := json.Marshal(scope.Statement.Vars)
		if err != nil {
			zlog.Errorf("json.Marshal Error : %s", err.Error())
		} else {
			span.(trace.Span).SetAttributes(attribute.Key("db.statement").String(string(vars)))
			span.(trace.Span).SetAttributes(attribute.Key("db.sql").String(scope.Statement.SQL.String()))
			if scope.Statement.Error != nil {
				span.(trace.Span).SetAttributes(attribute.Key("db.err").String(scope.Statement.Error.Error()))
			}
			span.(trace.Span).SetAttributes(attribute.Key("db.took μs").Int64(time.Since(tempTime.(time.Time)).Microseconds()))
		}

		defer span.(trace.Span).End()
	}
}

func registerCallbacks(db *gorm.DB, name string, c *callbacks) {
	beforeName := fmt.Sprintf("tracing:%v_before", name)
	afterName := fmt.Sprintf("tracing:%v_after", name)
	gormCallbackName := fmt.Sprintf("gorm:%v", name)
	switch name {
	case "create":
		_ = db.Callback().Create().Before(gormCallbackName).Register(beforeName, c.beforeCreate) //nolint:errcheck,staticcheck
		_ = db.Callback().Create().After(gormCallbackName).Register(afterName, c.afterCreate)    //nolint:errcheck,staticcheck
	case "query":
		_ = db.Callback().Query().Before(gormCallbackName).Register(beforeName, c.beforeQuery) //nolint:errcheck,staticcheck
		_ = db.Callback().Query().After(gormCallbackName).Register(afterName, c.afterQuery)    //nolint:errcheck,staticcheck
	case "update":
		_ = db.Callback().Update().Before(gormCallbackName).Register(beforeName, c.beforeUpdate) //nolint:errcheck,staticcheck
		_ = db.Callback().Update().After(gormCallbackName).Register(afterName, c.afterUpdate)    //nolint:errcheck,staticcheck
	case "delete":
		_ = db.Callback().Delete().Before(gormCallbackName).Register(beforeName, c.beforeDelete) //nolint:errcheck,staticcheck
		_ = db.Callback().Delete().After(gormCallbackName).Register(afterName, c.afterDelete)    //nolint:errcheck,staticcheck
	case "row_query":
		_ = db.Callback().Row().Before(gormCallbackName).Register(beforeName, c.beforeRowQuery) //nolint:errcheck,staticcheck
		_ = db.Callback().Row().After(gormCallbackName).Register(afterName, c.afterRowQuery)    //nolint:errcheck,staticcheck
	}
}
