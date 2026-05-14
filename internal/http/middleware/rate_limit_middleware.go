package middleware

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/config"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

func RateLimit(config config.RateLimitConfig) gin.HandlerFunc {
	limiter := rate.NewLimiter(
		config.PutTokenRate,
		config.BurstNumber,
	)
	skipPaths := map[string]bool{
		"/skip-path": true,
	}
	logger := deps.Logger()

	return func(c *gin.Context) {
		if skipPaths[c.Request.URL.Path] || limiter.Allow() {
			c.Next()
			return
		}
		errResp := response.NewErrorResponse(response.ServiceUnavailable, errors.New("token bucket is empty"), nil)
		logger.Error(response.ServiceUnavailable, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
	}
}
