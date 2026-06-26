package user

import (
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GetMeHandler struct {
	database *gorm.DB
	logger   *zap.Logger
}

// NewGetMeHandler creates the handler for the current user API.
func NewGetMeHandler(database *gorm.DB, logger *zap.Logger) *GetMeHandler {
	return &GetMeHandler{
		database: database,
		logger:   logger,
	}
}

// @tags user
// @accept json
// @produce json
// @param Authorization header string true " "
// @success 200 {object} response.Response{data=UserData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/user/me [get]
func (h *GetMeHandler) Handle(c *gin.Context) {
	u, err := user.Get(
		h.database.Preload("Roles.Permissions"),
		"id = ?",
		c.GetUint("user_id"),
	)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil, response.DebugMode(c))
		h.logger.Warn(response.BadRequest, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := UserData{}
	data.Fill(u)
	respBody := response.NewResponse(data)
	c.JSON(http.StatusOK, respBody)
}
