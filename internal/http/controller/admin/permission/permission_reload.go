package permission

import (
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ReloadHandler struct {
	casbinEnforcer *casbin.SyncedCachedEnforcer
	logger         *zap.Logger
}

func NewReloadHandler(casbinEnforcer *casbin.SyncedCachedEnforcer, logger *zap.Logger) *ReloadHandler {
	return &ReloadHandler{casbinEnforcer: casbinEnforcer, logger: logger}
}

// @tags admin-permission
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true " "
// @param Authorization header string true " "
// @success 204
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/admin/permission/reload [post]
func (h *ReloadHandler) Handle(c *gin.Context) {
	enforcer := h.casbinEnforcer
	if err := enforcer.LoadPolicy(); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil, response.DebugMode(c))
		logger := h.logger
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	c.Status(http.StatusNoContent)
}
