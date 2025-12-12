package httpapi

import (
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

type HttpApiSearchRequest struct {
	Method string `form:"method" schema:"method,omitempty" structs:"method,omitempty"`
	Path   string `form:"path" schema:"path,omitempty" structs:"path,omitempty"`
	pagination.PaginationRequest
}

type HttpApiSearchData struct {
	HttpApis                      []HttpApiData `json:"http_apis" mapstructure:"http_apis" validate:"required"`
	pagination.PaginationResponse `mapstructure:",squash"`
}

// @tags admin-http-api
// @accept json
// @produce json
// @param Authorization header string true " "
// @param "query string" query httpapi.HttpApiSearchRequest true " "
// @success 200 {object} response.Response{data=httpapi.HttpApiSearchData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/admin/http-api [get]
func Search(c *gin.Context) {
	queryString := new(HttpApiSearchRequest)
	logger := service.Registry.Get("logger").(*zap.Logger)
	if err := c.ShouldBindQuery(queryString); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), response.MakeValidationErrorContext(err))
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	whereConditions := pagination.WhereConditionMap{
		"method": func(db *gorm.DB, value any) {
			db.Where("method = ?", value)
		},
		"path": func(db *gorm.DB, value any) {
			db.Where("path = ?", value)
		},
	}
	paginator := pagination.NewPaginator(
		database.NewTxByTable("http_apis"),
		whereConditions,
		nil,
		queryString.Page,
		queryString.PerPage,
	)
	paginator.AddWhereConditions(structs.Map(queryString))

	httpApis := []model.HttpApi{}
	if err := paginator.Execute(&httpApis).Error; err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := HttpApiSearchData{HttpApis: make([]HttpApiData, len(httpApis))}
	for i := 0; i < len(data.HttpApis); i++ {
		data.HttpApis[i].Fill(&httpApis[i])
	}
	data.Total, data.LastPage = paginator.GetTotalAndLastPage()
	c.JSON(http.StatusOK, response.NewResponse(data))
}
