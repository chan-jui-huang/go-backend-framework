package permission

import (
	"fmt"
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/model"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/chan-jui-huang/go-backend-package/pkg/pagination"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PermissionSearchRequest struct {
	Name string `form:"name" structs:"name,omitempty" schema:"name,omitempty"`
	pagination.PaginationRequest
}

type PermissionSearchData struct {
	Permissions                   []PermissionData `json:"permissions" mapstructure:"permissions" validate:"required"`
	pagination.PaginationResponse `mapstructure:",squash"`
}

// @tags admin-permission
// @accept json
// @produce json
// @param Authorization header string true " "
// @param request body permission.PermissionSearchRequest true " "
// @success 200 {object} response.Response{data=permission.PermissionSearchData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/admin/permission [get]
func Search(c *gin.Context) {
	queryString := new(PermissionSearchRequest)
	logger := service.Registry.Get("logger").(*zap.Logger)
	if err := c.ShouldBindQuery(queryString); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), response.MakeValidationErrorContext(err))
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	wm := pagination.WhereConditionMap{
		"name": func(db *gorm.DB, value any) {
			db.Where("name LIKE ?", fmt.Sprintf("%%%v%%", value))
		},
	}
	paginator := pagination.NewPaginator(
		database.NewTxByTable("permissions"),
		wm,
		nil,
		queryString.Page,
		queryString.PerPage,
	)

	permissions := []model.Permission{}
	db := paginator.AddWhereConditions(structs.Map(queryString)).
		Execute(&permissions)
	if err := db.Error; err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := PermissionSearchData{Permissions: make([]PermissionData, len(permissions))}
	for i := 0; i < len(data.Permissions); i++ {
		data.Permissions[i].Fill(&permissions[i])
	}
	data.Total, data.LastPage = paginator.GetTotalAndLastPage()
	c.JSON(http.StatusOK, response.NewResponse(data))
}
