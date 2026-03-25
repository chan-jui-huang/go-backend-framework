package middleware

import (
	"fmt"
	"time"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func AccessLogger() gin.HandlerFunc {
	skipPaths := map[string]bool{
		"/skip-path": true,
	}
	accessLogger := deps.AccessLogger()

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
			accessLogger.Info(message, fields...)
		case status >= 400 && status < 500:
			accessLogger.Warn(message, fields...)
		case status >= 500:
			accessLogger.Error(message, fields...)
		}
	}
}
