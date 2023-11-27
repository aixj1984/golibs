package zlog

import (
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/* 日志处理 */

// GinLogger 是给gin框架提供访问日志输出的中间件
func (e *Entry) GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)

		e.Logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery 是给gin框架提供访异常回复时的日志中间件；错误处理，也可以不写自己处理错误使用gin写好的错误处理
func (e *Entry) GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool

				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne.Err, &se) {
						if strings.Contains(strings.ToLower(ne.Err.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(ne.Err.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, httpErr := httputil.DumpRequest(c.Request, false)

				if httpErr != nil {
					e.Logger.Error("httputil.DumpRequest", zap.Any("error", httpErr))
				}

				if brokenPipe {
					e.Logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)

					err := c.Error(err.(error))
					if err != nil {
						e.Logger.Error("c.Error", zap.Any("error", httpErr))
					}

					c.Abort()

					return
				}

				if stack {
					e.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					e.Logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
