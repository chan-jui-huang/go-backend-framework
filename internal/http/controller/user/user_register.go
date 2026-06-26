package user

import (
	"fmt"
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/database"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/model"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/pkg/user"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/argon2"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/authentication"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=8,containsany=abcdefghijklmnopqrstuvwxyz,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=1234567890"`
}
type RegisterHandler struct {
	database      *gorm.DB
	authenticator *authentication.
			Authenticator
	logger *zap.Logger
}

func NewRegisterHandler(database *gorm.DB, authenticator *authentication.
	Authenticator, logger *zap.Logger) *RegisterHandler {
	return &RegisterHandler{
		database: database, authenticator: authenticator, logger: logger}
}

// @tags user
// @accept json
// @produce json
// @param X-XSRF-TOKEN header string true " "
// @param request body user.UserRegisterRequest true " "
// @success 200 {object} response.Response{data=user.TokenData}
// @failure 400 {object} response.ErrorResponse "code: 400-001(Bad Request), 400-002(request validation failed)"
// @failure 403 {object} response.ErrorResponse "code: 403-001(Forbidden)"
// @failure 500 {object} response.ErrorResponse "code: 500-001(Internal Server Error)"
// @router /api/user/register [post]
func (h *RegisterHandler) Handle(c *gin.Context) {
	logger := h.logger
	reqBody := new(UserRegisterRequest)
	if err := c.ShouldBindBodyWithJSON(reqBody); err != nil {
		errResp := response.NewErrorResponse(response.RequestValidationFailed, errors.WithStack(err), response.MakeValidationErrorContext(err), response.DebugMode(c))
		logger.Warn(response.RequestValidationFailed, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	u := &model.User{
		Name:     reqBody.Name,
		Email:    reqBody.Email,
		Password: argon2.MakeArgon2IdHash(reqBody.Password),
	}
	if err := user.Create(database.NewTx(h.database), u); err != nil {
		errResp := response.NewErrorResponse(response.BadRequest, err, nil, response.DebugMode(c))
		logger.Warn(response.BadRequest, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
		return
	}

	authenticator := h.authenticator
	accessToken, err := authenticator.IssueAccessToken(fmt.Sprintf("%v", u.Id))
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
