package permission

import (
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/permission"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PermissionGetData struct {
	PermissionData `mapstructure:",squash"`
	HttpApis       []HttpApiData `json:"http_apis" mapstructure:"http_apis" validate:"required"`
}

// @tags admin-permission
// @accept json
// @produce json
// @param Authorization header string true " "
// @param id path string true " "
// @success 200 {object} response.Response{data=permission.PermissionGetData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/admin/permission/{id} [get]
func Get(c *gin.Context) {
	p, err := permission.Get(database.NewTx(), "id = ?", c.Param("id"))
	logger := service.Registry.Get("logger").(*zap.Logger)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	casbinRules, err := permission.GetCasbinRules(database.NewTx(), "ptype = ? AND v0 = ?", "p", p.Name)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger.Warn(errResp.Message, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := PermissionCreateData{
		HttpApis: make([]HttpApiData, len(casbinRules)),
	}
	data.PermissionData.Fill(p)
	for i := 0; i < len(data.HttpApis); i++ {
		data.HttpApis[i].Fill(casbinRules[i])
	}

	c.JSON(http.StatusOK, response.NewResponse(data))
}
