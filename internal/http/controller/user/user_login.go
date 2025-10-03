package user

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/user"
	"github.com/chan-jui-huang/go-backend-package/pkg/argon2"
	"github.com/chan-jui-huang/go-backend-package/pkg/authentication"
	"github.com/chan-jui-huang/go-backend-package/pkg/booter/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true " "
// @param request body user.UserLoginRequest true " "
// @success 200 {object} response.Response{data=user.TokenData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed), 400-003(email is wrong), 400-004(password is wrong)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/user/login [post]
func Login(c *gin.Context) {
	logger := service.Registry.Get("logger").(*zap.Logger)
	reqBody := new(UserLoginRequest)
	if err := c.ShouldBindJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), nil)
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	u, err := user.Get(database.NewTx(), "email = ?", reqBody.Email)
	if err != nil {
		errResp := response.NewErrorResponse(response.EmailIsWrong, err, nil)
		logger.Warn(response.EmailIsWrong, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}
	if !argon2.VerifyArgon2IdHash(reqBody.Password, u.Password) {
		errResp := response.NewErrorResponse(response.PasswordIsWrong, errors.New(response.PasswordIsWrong), nil)
		logger.Warn(response.PasswordIsWrong, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	encoder := schema.NewEncoder()
	values := url.Values{}
	userQuery := user.Query{
		UserId: u.Id,
	}
	if err := encoder.Encode(userQuery, values); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	authenticator := service.Registry.Get("authentication.authenticator").(*authentication.Authenticator)
	accessToken, err := authenticator.IssueAccessToken(values.Encode())
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil)
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := TokenData{}
	data.Fill(fmt.Sprintf("Bearer %s", accessToken))
	respBody := response.NewResponse(data)
	c.JSON(http.StatusOK, respBody)
}
