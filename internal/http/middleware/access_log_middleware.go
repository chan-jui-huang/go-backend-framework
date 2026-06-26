package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type AccessLogMiddleware struct {
	logger *zap.Logger
}

func NewAccessLogMiddleware(logger *zap.Logger) *AccessLogMiddleware {
	return &AccessLogMiddleware{
		logger: logger,
	}
}

func (m *AccessLogMiddleware) Handle() gin.HandlerFunc {
	skipPaths := map[string]bool{
		"/skip-path": true,
	}

	return func(c *gin.Context) {
		now := time.Now()
		path := c.Request.URL.Path
		c.Next()

		if skipPaths[path] {
			return
		}
		latency := time.Since(now)
		status := c.Writer.Status()
		fields := []zapcore.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("referer", c.Request.Referer()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
		}

		message := fmt.Sprintf("%s %s", c.Request.Method, path)
		switch {
		case status < 400:
			m.logger.Info(message, fields...)
		case status >= 400 && status < 500:
			m.logger.Warn(message, fields...)
		case status >= 500:
			m.logger.Error(message, fields...)
		}
	}
}
