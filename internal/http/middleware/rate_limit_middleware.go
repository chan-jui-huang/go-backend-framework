package middleware

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/config"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type RateLimitMiddleware struct {
	logger  *zap.Logger
	limiter *rate.Limiter
}

func NewRateLimitMiddleware(
	logger *zap.Logger,
	config *config.RateLimitConfig,
) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		logger: logger,
		limiter: rate.NewLimiter(
			config.PutTokenRate,
			config.BurstNumber,
		),
	}
}

func (m *RateLimitMiddleware) Handle() gin.HandlerFunc {
	skipPaths := map[string]bool{
		"/skip-path": true,
	}

	return func(c *gin.Context) {
		if skipPaths[c.Request.URL.Path] || m.limiter.Allow() {
			c.Next()
			return
		}
		errResp := response.NewErrorResponse(response.ServiceUnavailable, errors.New("token bucket is empty"), nil, response.DebugMode(c))
		m.logger.Error(response.ServiceUnavailable, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
	}
}
