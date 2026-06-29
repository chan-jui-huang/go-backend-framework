package permission

import (
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/model"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/permission"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PermissionDeleteRequest struct {
	Ids []uint `json:"ids" binding:"required"`
}
type DeleteHandler struct {
	database       *gorm.DB
	casbinEnforcer *casbin.SyncedCachedEnforcer
	logger         *zap.Logger
}

func NewDeleteHandler(database *gorm.DB, casbinEnforcer *casbin.SyncedCachedEnforcer, logger *zap.Logger) *DeleteHandler {
	return &DeleteHandler{
		database: database, casbinEnforcer: casbinEnforcer, logger: logger}
}

// @tags admin-permission
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true " "
// @param Authorization header string true " "
// @param id path string true " "
// @param request body permission.PermissionDeleteRequest true " "
// @success 204
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/admin/permission [delete]
func (h *DeleteHandler) Handle(c *gin.Context) {
	reqBody := new(PermissionDeleteRequest)
	logger := h.logger
	if err := c.ShouldBindBodyWithJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), response.MakeValidationErrorContext(err), response.DebugMode(c))
		logger.Warn(errResp.Message, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	permissions, err := permission.GetMany(database.NewTx(h.database), "id IN ?", reqBody.Ids)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil, response.DebugMode(c))
		logger.Warn(errResp.Message, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	err = database.NewTx(h.database).Transaction(func(tx *gorm.DB) error {
		if err := permission.Delete(tx, "id IN ?", reqBody.Ids); err != nil {
			return err
		}

		names := lo.Map(permissions, func(p model.Permission, _ int) string {
			return p.Name
		})
		if err := permission.DeleteCasbinRule(tx, "ptype = ? AND v0 IN ?", "p", names); err != nil {
			return err
		}

		if err := permission.DeleteCasbinRule(tx, "ptype = ? AND v1 IN ?", "g", names); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil, response.DebugMode(c))
		logger.Warn(errResp.Message, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	enforcer := h.casbinEnforcer
	if err := enforcer.LoadPolicy(); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil, response.DebugMode(c))
		logger.Warn(errResp.Message, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	c.Status(http.StatusNoContent)
}
