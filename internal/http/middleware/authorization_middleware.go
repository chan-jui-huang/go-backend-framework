package middleware

import (
	"strconv"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func Authorize() gin.HandlerFunc {
	logger := deps.Logger()
	return func(c *gin.Context) {
		userId := c.GetUint("user_id")
		enforcer := deps.CasbinEnforcer()
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
