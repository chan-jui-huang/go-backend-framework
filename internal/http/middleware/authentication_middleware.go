package middleware

import (
	"net/url"
	"strings"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/pkg/user"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/authentication"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

func Authenticate() gin.HandlerFunc {
	logger := deps.Logger()
	authenticator := deps.Authenticator()
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorizationHeader, "Bearer") {
			errResp := response.NewErrorResponse(response.Unauthorized, errors.New("jwt authentication failed"), nil)
			logger.Warn(response.Unauthorized, errResp.MakeLogFields(c.Request)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}

		accessTokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
		subject, err := verifyAccessToken(authenticator, accessTokenString)
		if err != nil {
			errResp := response.NewErrorResponse(response.Unauthorized, err, nil)
			logger.Warn(response.Unauthorized, errResp.MakeLogFields(c.Request)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}

		values, err := url.ParseQuery(subject)
		if err != nil {
			errResp := response.NewErrorResponse(response.Unauthorized, errors.WithStack(err), nil)
			logger.Warn(response.Unauthorized, errResp.MakeLogFields(c.Request)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}

		decoder := schema.NewDecoder()
		userQuery := &user.Query{}
		if err := decoder.Decode(userQuery, values); err != nil {
			errResp := response.NewErrorResponse(response.Unauthorized, errors.WithStack(err), nil)
			logger.Warn(response.Unauthorized, errResp.MakeLogFields(c.Request)...)
			c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
			return
		}
		c.Set("user_id", userQuery.UserId)

		c.Next()
	}
}

func verifyAccessToken(authenticator *authentication.Authenticator, accessTokenString string) (string, error) {
	accessToken, err := authenticator.VerifyJwt(accessTokenString)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if !accessToken.Valid {
		return "", errors.New("access token is invalid")
	}

	claims := accessToken.Claims.(jwt.MapClaims)
	if claims["aud"].([]any)[0].(string) != "access" {
		return "", errors.New("jwt aud claim is not [access] string")
	}

	return claims["sub"].(string), nil
}
