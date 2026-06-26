package user

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/user"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/argon2"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/authentication"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type LoginHandler struct {
	database      *gorm.DB
	authenticator *authentication.
			Authenticator
	logger *zap.Logger
}

func NewLoginHandler(database *gorm.DB, authenticator *authentication.
	Authenticator, logger *zap.Logger) *LoginHandler {
	return &LoginHandler{
		database: database, authenticator: authenticator, logger: logger}
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
func (h *LoginHandler) Handle(c *gin.Context) {
	logger := h.logger
	reqBody := new(UserLoginRequest)
	if err := c.ShouldBindBodyWithJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), response.MakeValidationErrorContext(err), response.DebugMode(c))
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	u, err := user.Get(database.NewTx(h.database), "email = ?", reqBody.Email)
	if err != nil {
		errResp := response.NewErrorResponse(response.EmailIsWrong, err, nil, response.DebugMode(c))
		logger.Warn(response.EmailIsWrong, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}
	if !argon2.VerifyArgon2IdHash(reqBody.Password, u.Password) {
		errResp := response.NewErrorResponse(response.PasswordIsWrong, errors.New(response.PasswordIsWrong), nil, response.DebugMode(c))
		logger.Warn(response.PasswordIsWrong, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	encoder := schema.NewEncoder()
	values := url.Values{}
	userQuery := user.Query{
		UserId: u.Id,
	}
	if err := encoder.Encode(userQuery, values); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil, response.DebugMode(c))
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	authenticator := h.authenticator
	accessToken, err := authenticator.IssueAccessToken(values.Encode())
	if err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, errors.WithStack(err), nil, response.DebugMode(c))
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	data := TokenData{}
	data.Fill(fmt.Sprintf("Bearer %s", accessToken))
	respBody := response.NewResponse(data)
	c.JSON(http.StatusOK, respBody)
}
