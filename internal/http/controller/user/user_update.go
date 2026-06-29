package user

import (
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/user"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserUpdateRequest struct {
	Name  string `json:"name" binding:"required" structs:"name"`
	Email string `json:"email" binding:"required" structs:"email"`
}
type UpdateHandler struct {
	database *gorm.DB
	logger   *zap.Logger
}

func NewUpdateHandler(database *gorm.DB, logger *zap.Logger) *UpdateHandler {
	return &UpdateHandler{
		database: database, logger: logger}
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true " "
// @param Authorization header string true " "
// @param request body user.UserUpdateRequest true " "
// @success 200 {object} response.Response{data=user.UserData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/user [put]
func (h *UpdateHandler) Handle(c *gin.Context) {
	logger := h.logger
	reqBody := new(UserUpdateRequest)
	if err := c.ShouldBindBodyWithJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), response.MakeValidationErrorContext(err), response.DebugMode(c))
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	values := structs.Map(reqBody)
	_, err := user.Update(database.NewTx(h.database), c.GetUint("user_id"), values)
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil, response.DebugMode(c))
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	u, err := user.Get(database.NewTx(h.database), "id = ?", c.GetUint("user_id"))
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil, response.DebugMode(c))
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := UserData{}
	data.Fill(u)
	respBody := response.NewResponse(data)
	c.JSON(http.StatusOK, respBody)
}
