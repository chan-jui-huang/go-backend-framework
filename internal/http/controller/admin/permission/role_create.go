package permission

import (
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoleCreateRequest struct {
	Name          string `json:"name" binding:"required"`
	IsPublic      bool   `json:"is_public"`
	PermissionIds []uint `json:"permission_ids" binding:"required,min=1"`
}

// @tags admin-permission
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true " "
// @param Authorization header string true " "
// @param request body permission.RoleCreateRequest true " "
// @success 200 {object} response.Response{data=permission.RoleData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/admin/role [post]
func CreateRole(c *gin.Context) {
	reqBody := new(RoleCreateRequest)
	logger := service.Registry.Get("logger").(*zap.Logger)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	role := &model.Role{
		Name:     reqBody.Name,
		IsPublic: reqBody.IsPublic,
	}
	rolePermissions := make([]model.RolePermission, len(reqBody.PermissionIds))
	err := database.NewTx().Transaction(func(tx *gorm.DB) error {
		if err := permission.CreateRole(tx, role); err != nil {
			return err
		}

		for i := 0; i < len(rolePermissions); i++ {
			rolePermissions[i].RoleId = role.Id
			rolePermissions[i].PermissionId = reqBody.PermissionIds[i]
		}
		if err := permission.CreateRolePermission(tx, rolePermissions); err != nil {
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

	role, err = permission.GetRole(database.NewTx("Permissions"), "id = ?", role.Id)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := RoleData{
		Permissions: make([]PermissionData, len(role.Permissions)),
	}
	data.Fill(role)
	for i := 0; i < len(data.Permissions); i++ {
		data.Permissions[i].Fill(&role.Permissions[i])
	}

	c.JSON(http.StatusOK, response.NewResponse(data))
}
