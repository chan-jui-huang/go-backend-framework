package permission

import (
	"database/sql"
	"net/http"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PermissionUpdateRequest struct {
	Name     string `json:"name" structs:"name" binding:"required"`
	HttpApis []struct {
		Path   string `json:"path" binding:"required"`
		Method string `json:"method" binding:"required"`
	} `json:"http_apis" structs:"-" binding:"required,min=1,dive"`
}

type PermissionUpdateData struct {
	PermissionData `mapstructure:",squash"`
	HttpApis       []HttpApiData `json:"http_apis" mapstructure:"http_apis" validate:"required"`
}

// @tags admin-permission
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true " "
// @param Authorization header string true " "
// @param id path string true " "
// @param request body permission.PermissionUpdateRequest true " "
// @success 200 {object} response.Response{data=permission.PermissionUpdateData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/admin/permission/{id} [put]
func Update(c *gin.Context) {
	reqBody := new(PermissionUpdateRequest)
	logger := service.Registry.Get("logger").(*zap.Logger)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	casbinRules := make([]gormadapter.CasbinRule, len(reqBody.HttpApis))
	for i := 0; i < len(casbinRules); i++ {
		casbinRules[i].Ptype = "p"
		casbinRules[i].V0 = reqBody.Name
		casbinRules[i].V1 = reqBody.HttpApis[i].Path
		casbinRules[i].V2 = reqBody.HttpApis[i].Method
	}

	err := database.NewTx().Transaction(func(tx *gorm.DB) error {
		p, err := permission.Get(tx.Clauses(clause.Locking{Strength: "UPDATE"}), "id = ?", c.Param("id"))
		if err != nil {
			return err
		}

		if err := permission.DeleteCasbinRule(tx, "ptype = ? AND v0 = ?", "p", p.Name); err != nil {
			return err
		}

		if rows, err := permission.Update(tx, c.Param("id"), structs.Map(reqBody)); err != nil || rows != 1 {
			tx.Rollback()
			return err
		}

		if err := permission.CreateCasbinRule(tx, casbinRules); err != nil {
			return err
		}

		_, err = permission.UpdateCasbinRule(
			tx,
			map[string]any{"v1": reqBody.Name},
			"ptype = ? AND v1 = ?",
			"g",
			p.Name,
		)
		if err != nil {
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

	p, err := permission.Get(database.NewTx(), "id = ?", c.Param("id"))
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := PermissionGetData{
		HttpApis: make([]HttpApiData, len(casbinRules)),
	}
	data.PermissionData.Fill(p)
	for i := 0; i < len(data.HttpApis); i++ {
		data.HttpApis[i].Fill(casbinRules[i])
	}

	c.JSON(http.StatusOK, response.NewResponse(data))
}
