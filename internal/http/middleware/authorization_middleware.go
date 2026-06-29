package middleware

import (
	"strconv"

	"github.com/casbin/casbin/v3"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type AuthorizationMiddleware struct {
	logger   *zap.Logger
	enforcer *casbin.SyncedCachedEnforcer
}

func NewAuthorizationMiddleware(
	logger *zap.Logger,
	enforcer *casbin.SyncedCachedEnforcer,
) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		logger:   logger,
		enforcer: enforcer,
	}
}

func (m *AuthorizationMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetUint("user_id")
		ok, err := m.enforcer.Enforce(strconv.FormatUint(uint64(userId), 10), c.Request.URL.Path, c.Request.Method)
		if err != nil {
			errResp := response.NewErrorResponse(response.Forbidden, errors.WithStack(err), nil, response.DebugMode(c))
			m.logger.Warn(response.Forbidden, errResp.MakeLogFields(c)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}

		if !ok {
			errResp := response.NewErrorResponse(response.Forbidden, errors.New("casbin authorization failed"), nil, response.DebugMode(c))
			m.logger.Warn(response.Forbidden, errResp.MakeLogFields(c)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}

		c.Next()
	}
}
