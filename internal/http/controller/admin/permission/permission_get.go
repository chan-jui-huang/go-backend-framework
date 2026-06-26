package permission

import (
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/permission"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PermissionGetData struct {
	PermissionData `mapstructure:",squash"`
	HttpApis       []HttpApiData `json:"http_apis" mapstructure:"http_apis" validate:"required"`
}
type GetHandler struct {
	database *gorm.DB
	logger   *zap.Logger
}

func NewGetHandler(database *gorm.DB, logger *zap.Logger) *GetHandler {
	return &GetHandler{
		database: database, logger: logger}
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
func (h *GetHandler) Handle(c *gin.Context) {
	p, err := permission.Get(database.NewTx(h.database), "id = ?", c.Param("id"))
	logger := h.logger
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil, response.DebugMode(c))
		logger.Warn(errResp.Message, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	casbinRules, err := permission.GetCasbinRules(database.NewTx(h.database), "ptype = ? AND v0 = ?", "p", p.Name)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil, response.DebugMode(c))
		logger.Warn(errResp.Message, errResp.MakeLogFields(c)...)
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
