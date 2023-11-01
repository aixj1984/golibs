package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	viper "github.com/aixj1984/golibs/conf"

	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	parentSpanGormKey    = "opentracingParentSpan"
	parentSpanGormCtxKey = "opentracingParentSpanCtx"
	spanGormKey          = "opentracingSpan"
	spanDuration         = "opentracingSpanDuration"
)

var (
	defaultDatabase     = "mysql"
	connStrTmpl         = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s"
	defaultMaxOpenConns = 200
	defaultMaxIdleConns = 60
	defaultMaxLeftTime  = 300 * time.Second
	defaultCharset      = "utf8mb4"
	defaultPort         = 3306
	defaultTimeZone     = "Local"
	gormEngine          *Engine
)

type Engine struct {
	gorm *gorm.DB
}

func init() {
	dbCfg, err := viper.GetSubCfg[Config]("mysqlDB")
	if err != nil {
		fmt.Printf("unable to get config, %s", err.Error())
		return
	} else {
		New(dbCfg)
	}
}

func GetEngine() *Engine {
	return gormEngine
}

// New 实例化新的Gorm实例
func New(conf *Config) *Engine {
	if gormEngine != nil {
		return gormEngine
	}
	err := authConfig(conf)
	if err != nil {
		panic(err)
	}
	if strings.TrimSpace(conf.Charset) == "" {
		conf.Charset = defaultCharset
	}
	if strings.TrimSpace(conf.TimeZone) == "" {
		conf.TimeZone = defaultTimeZone
	}
	gormConf := &gorm.Config{}
	dsn := fmt.Sprintf(connStrTmpl,
		conf.User,
		conf.Password,
		conf.Server,
		conf.Port,
		conf.Database,
		conf.Charset,
		conf.TimeZone)
	db, err := gorm.Open(mysql.Open(dsn), gormConf)
	if err != nil {
		panic(err)
	}

	gormEngine = &Engine{db}
	gormEngine.wrapLog()
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetConnMaxLifetime(conf.MaxLeftTime)
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

	addGormCallbacks(db)
	return gormEngine
}

func (db *Engine) Context(ctx context.Context) *gorm.DB {
	if ctx == nil {
		fmt.Println("Engine no ctx")
		return db.gorm
	}
	parentSpan := trace.SpanFromContext(ctx)
	return db.gorm.WithContext(ctx).Set(parentSpanGormKey, parentSpan).Set(parentSpanGormCtxKey, ctx)
}

func (db *Engine) GetDB() *gorm.DB {
	return db.gorm
}

func (db *Engine) SetLogMode(mode bool) {
	if !mode {
		db.gorm.Logger.LogMode(LogLevelSilent)
	}
}

func (db *Engine) SetLogLevel(level LogLevel) {
	db.gorm.Logger.LogMode(level)
}

func Context(ctx context.Context) *gorm.DB {
	if gormEngine == nil {
		panic(fmt.Errorf("must init gorm.New"))
	}
	parentSpan := trace.SpanFromContext(ctx)
	return gormEngine.gorm.Set(parentSpanGormKey, parentSpan).Set(parentSpanGormCtxKey, ctx)
}

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
		//xlog.Infoln("no parentSpanCtx")
		return
	}
	db.Set(spanDuration, time.Now())
	_, span := otel.Tracer("GORM-V2-SQL").Start(parentSpanCtx.(context.Context), db.Statement.Name())
	db.Set(spanGormKey, span)
}

func (c *callbacks) after(scope *gorm.DB) {
	t, ok := scope.Get(spanDuration)
	if !ok {
		t = time.Now()
	}
	if span, ok := scope.Get(spanGormKey); ok {
		vars, _ := json.Marshal(scope.Statement.Vars)
		span.(trace.Span).SetAttributes(attribute.Key("db.statement").String(string(vars)))
		span.(trace.Span).SetAttributes(attribute.Key("db.sql").String(scope.Statement.SQL.String()))
		if scope.Statement.Error != nil {
			span.(trace.Span).SetAttributes(attribute.Key("db.err").String(scope.Statement.Error.Error()))
		}
		span.(trace.Span).SetAttributes(attribute.Key("db.took μs").Int64(time.Since(t.(time.Time)).Microseconds()))
		defer span.(trace.Span).End()
	} else {
		//xlog.Infoln("no span")
	}
}

func registerCallbacks(db *gorm.DB, name string, c *callbacks) {
	beforeName := fmt.Sprintf("tracing:%v_before", name)
	afterName := fmt.Sprintf("tracing:%v_after", name)
	gormCallbackName := fmt.Sprintf("gorm:%v", name)
	switch name {
	case "create":
		_ = db.Callback().Create().Before(gormCallbackName).Register(beforeName, c.beforeCreate)
		_ = db.Callback().Create().After(gormCallbackName).Register(afterName, c.afterCreate)
	case "query":
		_ = db.Callback().Query().Before(gormCallbackName).Register(beforeName, c.beforeQuery)
		_ = db.Callback().Query().After(gormCallbackName).Register(afterName, c.afterQuery)
	case "update":
		_ = db.Callback().Update().Before(gormCallbackName).Register(beforeName, c.beforeUpdate)
		_ = db.Callback().Update().After(gormCallbackName).Register(afterName, c.afterUpdate)
	case "delete":
		_ = db.Callback().Delete().Before(gormCallbackName).Register(beforeName, c.beforeDelete)
		_ = db.Callback().Delete().After(gormCallbackName).Register(afterName, c.afterDelete)
	case "row_query":
		_ = db.Callback().Row().Before(gormCallbackName).Register(beforeName, c.beforeRowQuery)
		_ = db.Callback().Row().After(gormCallbackName).Register(afterName, c.afterRowQuery)
	}
}
