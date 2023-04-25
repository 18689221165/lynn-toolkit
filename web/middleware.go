package web

import (
	"bytes"
	"github.com/18689221165/lynn-toolkit/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

type config struct {
	TimeFormat string
	UTC        bool
	SkipPaths  []string
}

type accessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w accessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

// GinZapLog 用zap包记录日志,可打印post请求form及响应内容
func GinZapLog(logger *zap.Logger, runMode string) gin.HandlerFunc {
	conf := config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  nil,
	}
	skipPaths := make(map[string]bool, len(conf.SkipPaths))
	for _, path := range conf.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		bw := accessLogWriter{ResponseWriter: c.Writer, body: bytes.NewBufferString("")}
		c.Writer = bw

		c.Next()

		if _, ok := skipPaths[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)
			if conf.UTC {
				end = end.UTC()
			}

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				for _, e := range c.Errors.Errors() {
					logger.Error(e)
				}
			} else {
				fields := []zapcore.Field{
					zap.Int("status", c.Writer.Status()),
					zap.String("method", c.Request.Method),
					zap.String("path", path),
					zap.String("query", query),
				}
				if runMode == "debug" {
					// 只打印非上传文件的Post表单内容
					if (c.Request.Method == "POST" || c.Request.Method == "PUT") && !strings.Contains(c.ContentType(), "multipart/form-data") {
						fields = append(fields, zap.String("form", c.Request.PostForm.Encode()))
					}
					fields = append(fields, zap.String("resp", bw.body.String()))
				}

				if conf.TimeFormat != "" {
					fields = append(fields, zap.String("time", end.Format(conf.TimeFormat)))
				}
				fields = append(fields, zap.Duration("latency", latency))

				logger.Info(path, fields...)
			}
		}
	}
}

// RecoveryWithZap 拦截处理请求的goroute中发生的panic
// copy了gin_zap包的代码，修改了painc后返回自定义的消息包
func RecoveryWithZap(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, service.ErrInternalServerError)
			}
		}()
		c.Next()
	}
}

// IgnoreIndexAndFavicon 忽略index和favicon.ico请求
func IgnoreIndexAndFavicon() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/favicon.ico" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		if path == "/" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
