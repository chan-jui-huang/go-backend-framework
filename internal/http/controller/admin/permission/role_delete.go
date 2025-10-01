package permission

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoleDeleteRequest struct {
	Ids []uint `json:"ids" binding:"required"`
}

// @tags admin-permission
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true "csrf token"
// @param Authorization header string true "bearer token"
// @param id path string true "id"
// @param request body permission.RoleDeleteRequest true " "
// @success 204
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/admin/role [delete]
func DeleteRoles(c *gin.Context) {
	reqBody := new(RoleDeleteRequest)
	logger := service.Registry.Get("logger").(*zap.Logger)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	roles, err := permission.GetRoles(database.NewTx("Permissions", "Users"), "id IN ?", reqBody.Ids)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	err = database.NewTx().Transaction(func(tx *gorm.DB) error {
		if err := permission.DeleteRolePermission(tx, "role_id IN ?", reqBody.Ids); err != nil {
			return err
		}

		if err := permission.DeleteUserRole(tx, "role_id IN ?", reqBody.Ids); err != nil {
			return err
		}

		if err := permission.DeleteRole(tx, "id IN ?", reqBody.Ids); err != nil {
			return err
		}

		names := []string{}
		userIds := []uint{}
		for i := 0; i < len(roles); i++ {
			for j := 0; j < len(roles[i].Permissions); j++ {
				names = append(names, roles[i].Permissions[j].Name)
			}
			for j := 0; j < len(roles[i].Users); j++ {
				userIds = append(userIds, roles[i].Users[j].Id)
			}
		}
		if len(userIds) == 0 {
			return nil
		}
		if err := permission.DeleteCasbinRule(tx, "ptype = ? AND v0 IN ? AND v1 IN ?", "g", userIds, names); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	enforcer := service.Registry.Get("casbinEnforcer").(*casbin.SyncedCachedEnforcer)
	if err := enforcer.LoadPolicy(); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	c.Status(http.StatusNoContent)
}
