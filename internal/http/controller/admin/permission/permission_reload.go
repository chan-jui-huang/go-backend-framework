package permission

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// @tags admin-permission
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param Authorization header string true "bearer token"
// @success 204
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/admin/permission/reload [post]
func Reload(c *gin.Context) {
	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	if err := enforcer.LoadPolicy(); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger := service.Registry.Get("logger").(*zap.Logger)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	c.Status(http.StatusNoContent)
}
