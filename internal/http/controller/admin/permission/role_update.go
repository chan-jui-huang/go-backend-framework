package permission

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleUpdateRequest struct {
	Name          string `json:"name" structs:"name" binding:"required"`
	IsPublic      bool   `json:"is_public" structs:"is_public"`
	PermissionIds []uint `json:"permission_ids" structs:"-" binding:"required,min=1"`
}

// @tags admin-permission
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true " "
// @param Authorization header string true " "
// @param id path string true " "
// @param request body permission.RoleUpdateRequest true " "
// @success 200 {object} response.Response{data=permission.RoleData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/admin/role/{id} [put]
func UpdateRole(c *gin.Context) {
	reqBody := new(RoleUpdateRequest)
	logger := service.Registry.Get("logger").(*zap.Logger)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	rolePermissions := make([]model.RolePermission, len(reqBody.PermissionIds))
	err := database.NewTx().Transaction(func(tx *gorm.DB) error {
		role, err := permission.GetRole(
			tx.Table("roles").Preload("Permissions").Preload("Users").Clauses(clause.Locking{Strength: "UPDATE"}),
			"id = ?",
			c.Param("id"),
		)
		if err != nil {
			return err
		}

		if err := permission.DeleteRolePermission(tx, "role_id = ?", role.Id); err != nil {
			return err
		}

		userIds := lo.Map(role.Users, func(u model.User, _ int) uint {
			return u.Id
		})
		permissionNames := lo.Map(role.Permissions, func(p model.Permission, _ int) string {
			return p.Name
		})
		if len(userIds) > 0 {
			if err := permission.DeleteCasbinRule(tx, "ptype = ? AND v0 IN ? AND v1 IN ?", "g", userIds, permissionNames); err != nil {
				return err
			}
		}

		if rows, err := permission.UpdateRole(tx, role.Id, structs.Map(reqBody)); err != nil || rows != 1 {
			tx.Rollback()
			return err
		}

		for i := 0; i < len(rolePermissions); i++ {
			rolePermissions[i].RoleId = role.Id
			rolePermissions[i].PermissionId = reqBody.PermissionIds[i]
		}
		if err := permission.CreateRolePermission(tx, rolePermissions); err != nil {
			return err
		}

		if len(userIds) == 0 {
			return nil
		}

		permissions, err := permission.GetMany(tx, "id IN ?", reqBody.PermissionIds)
		if err != nil {
			return err
		}
		casbinRules := []gormadapter.CasbinRule{}
		for i := 0; i < len(userIds); i++ {
			userId := fmt.Sprintf("%d", userIds[i])
			for j := 0; j < len(permissions); j++ {
				casbinRule := gormadapter.CasbinRule{
					Ptype: "g",
					V0:    userId,
					V1:    permissions[j].Name,
				}
				casbinRules = append(casbinRules, casbinRule)
			}
		}
		if err := permission.CreateCasbinRule(tx, casbinRules); err != nil {
			return err
		}

		return nil
	})
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
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

	role, err := permission.GetRole(database.NewTx("Permissions"), "id = ?", c.Param("id"))
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := RoleData{}
	data.Fill(role)
	c.JSON(http.StatusOK, response.NewResponse(data))
}
