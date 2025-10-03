package user

import (
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/user"
	"github.com/chan-jui-huang/go-backend-package/pkg/argon2"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type UserUpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" structs:"-"`
	Password        string `json:"password" binding:"required,gte=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=1234567890" structs:"password"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password" structs:"-"`
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
func UpdatePassword(c *gin.Context) {
	logger := service.Registry.Get("logger").(*zap.Logger)
	reqBody := new(UserUpdatePasswordRequest)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	u, err := user.Get(database.NewTx(), "id = ?", c.GetUint("user_id"))
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	if !argon2.VerifyArgon2IdHash(reqBody.CurrentPassword, u.Password) {
		errResp := response.NewErrorResponse(response.PasswordIsWrong, errors.New(response.PasswordIsWrong), nil)
		logger.Warn(response.PasswordIsWrong, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	reqBody.Password = argon2.MakeArgon2IdHash(reqBody.Password)
	values := structs.Map(reqBody)
	if _, err := user.Update(database.NewTx(), c.GetUint("user_id"), values); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	c.Status(http.StatusNoContent)
}
