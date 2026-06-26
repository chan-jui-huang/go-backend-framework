package user

import (
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/user"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/argon2"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserUpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" structs:"-"`
	Password        string `json:"password" binding:"required,gte=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=1234567890" structs:"password"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password" structs:"-"`
}
type UpdatePasswordHandler struct {
	database *gorm.DB
	logger   *zap.Logger
}

func NewUpdatePasswordHandler(database *gorm.DB, logger *zap.Logger) *UpdatePasswordHandler {
	return &UpdatePasswordHandler{
		database: database, logger: logger}
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true " "
// @param Authorization header string true " "
// @param request body user.UserUpdatePasswordRequest true " "
// @success 204
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed)"
// @failure 401 {object} response.ErrorResponse "code: 401-001(Unauthorized)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/user/password [put]
func (h *UpdatePasswordHandler) Handle(c *gin.Context) {
	logger := h.logger
	reqBody := new(UserUpdatePasswordRequest)
	if err := c.ShouldBindBodyWithJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), response.MakeValidationErrorContext(err), response.DebugMode(c))
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c)...)
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

	if !argon2.VerifyArgon2IdHash(reqBody.CurrentPassword, u.Password) {
		errResp := response.NewErrorResponse(response.PasswordIsWrong, errors.New(response.PasswordIsWrong), nil, response.DebugMode(c))
		logger.Warn(response.PasswordIsWrong, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	reqBody.Password = argon2.MakeArgon2IdHash(reqBody.Password)
	values := structs.Map(reqBody)
	if _, err := user.Update(database.NewTx(h.database), c.GetUint("user_id"), values); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil, response.DebugMode(c))
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	c.Status(http.StatusNoContent)
}
