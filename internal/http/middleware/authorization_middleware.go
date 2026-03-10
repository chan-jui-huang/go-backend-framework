package middleware

import (
	"strconv"

	"github.com/casbin/casbin/v3"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Authorize() gin.HandlerFunc {
	logger := service.Registry.Get("logger").(*zap.Logger)
	return func(c *gin.Context) {
		userId := c.GetUint("user_id")
		enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
		ok, err := enforcer.Enforce(strconv.FormatUint(uint64(userId), 10), c.Request.URL.Path, c.Request.Method)
		if err != nil {
			errResp := response.NewErrorResponse(response.Forbidden, errors.WithStack(err), nil)
			logger.Warn(response.Forbidden, errResp.MakeLogFields(c.Request)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}

		if !ok {
			errResp := response.NewErrorResponse(response.Forbidden, errors.New("casbin authorization failed"), nil)
			logger.Warn(response.Forbidden, errResp.MakeLogFields(c.Request)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}

		c.Next()
	}
}
