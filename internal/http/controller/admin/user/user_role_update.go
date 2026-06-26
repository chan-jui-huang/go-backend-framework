package user

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/model"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/user"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRoleUpdateRequest struct {
	UserId  uint   `json:"user_id" binding:"required"`
	RoleIds []uint `json:"role_ids" binding:"required"`
}
type UpdateUserRoleHandler struct {
	database       *gorm.DB
	casbinEnforcer *casbin.SyncedCachedEnforcer
	logger         *zap.Logger
}

func NewUpdateUserRoleHandler(database *gorm.DB, casbinEnforcer *casbin.SyncedCachedEnforcer, logger *zap.Logger) *UpdateUserRoleHandler {
	return &UpdateUserRoleHandler{
		database: database, casbinEnforcer: casbinEnforcer, logger: logger}
}

// @tags admin-user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true " "
// @param Authorization header string true " "
// @param request body user.UserRoleUpdateRequest true " "
// @success 200 {object} response.Response{data=user.UserData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed), 400-005(permission is repeat)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/admin/user-role [put]
func (h *UpdateUserRoleHandler) Handle(c *gin.Context) {
	reqBody := new(UserRoleUpdateRequest)
	logger := h.logger
	if err := c.ShouldBindBodyWithJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), response.MakeValidationErrorContext(err), response.DebugMode(c))
		logger.Warn(errResp.Message, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	roles, err := permission.GetRoles(database.NewTx(h.database, "Permissions"), "id IN ?", reqBody.RoleIds)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil, response.DebugMode(c))
		logger.Warn(errResp.Message, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	permissions := []model.Permission{}
	for i := 0; i < len(roles); i++ {
		permissions = append(permissions, roles[i].Permissions...)
	}
	if len(permissions) != len(lo.Union(permissions)) {
		errResp := response.NewErrorResponse(response.PermissionIsRepeat, errors.WithStack(err), nil, response.DebugMode(c))
		logger.Warn(errResp.Message, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	userId := fmt.Sprintf("%d", reqBody.UserId)
	casbinRules := make([]gormadapter.CasbinRule, len(permissions))
	for i := 0; i < len(permissions); i++ {
		casbinRules[i].Ptype = "g"
		casbinRules[i].V0 = userId
		casbinRules[i].V1 = permissions[i].Name
	}

	userRoles := make([]model.UserRole, len(reqBody.RoleIds))
	for i := 0; i < len(reqBody.RoleIds); i++ {
		userRoles[i].UserId = reqBody.UserId
		userRoles[i].RoleId = reqBody.RoleIds[i]
	}
	err = database.NewTx(h.database).Transaction(func(tx *gorm.DB) error {
		if err := permission.DeleteUserRole(tx, "user_id = ?", reqBody.UserId); err != nil {
			return err
		}

		if err := permission.DeleteCasbinRule(tx, "ptype = ? AND v0 = ?", "g", reqBody.UserId); err != nil {
			return err
		}

		if len(reqBody.RoleIds) == 0 {
			return nil
		}

		if err := permission.CreateUserRole(tx, userRoles); err != nil {
			return err
		}

		if err := permission.CreateCasbinRule(tx, casbinRules); err != nil {
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

	u, err := user.Get(database.NewTx(h.database, "Roles"), "id = ?", reqBody.UserId)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil, response.DebugMode(c))
		logger.Warn(errResp.Message, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := UserData{}
	data.Fill(u)
	c.JSON(http.StatusOK, response.NewResponse(data))
}
