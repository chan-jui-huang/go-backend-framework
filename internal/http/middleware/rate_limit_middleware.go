package middleware

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

type RateLimitConfig struct {
	PutTokenRate rate.Limit
	BurstNumber  int
}

func RateLimit(config RateLimitConfig) gin.HandlerFunc {
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
		logger.Error(response.ServiceUnavailable, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
	}
}
