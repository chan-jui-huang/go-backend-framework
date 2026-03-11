package middleware

import (
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/http/response"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/random"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type CsrfConfig struct {
	Cookie struct {
		Name     string
		Path     string
		Domain   string
		MaxAge   int
		Secure   bool
		HttpOnly bool
		SameSite http.SameSite
	}
	Header string
}

func VerifyCsrfToken(config CsrfConfig) gin.HandlerFunc {
	skipPaths := map[string]bool{
		"/skip-path": true,
	}
	logger := deps.Logger()

	return func(c *gin.Context) {
		setCsrfToken(c, config)
		if isReadingHttpMethod(c) ||
			skipPaths[c.Request.URL.Path] ||
			verifyCsrfToken(c, config.Cookie.Name, config.Header) {
			c.Next()
			return
		}

		errResp := response.NewErrorResponse(response.Forbidden, errors.New("csrf token mismatch"), nil)
		logger.Warn(response.Forbidden, errResp.MakeLogFields(c.Request)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
	}
}

func setCsrfToken(c *gin.Context, config CsrfConfig) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.Cookie.Name,
		Value:    random.RandomString(20),
		Path:     config.Cookie.Path,
		Domain:   config.Cookie.Domain,
		MaxAge:   config.Cookie.MaxAge,
		Secure:   config.Cookie.Secure,
		HttpOnly: config.Cookie.HttpOnly,
		SameSite: config.Cookie.SameSite,
	})
}

func isReadingHttpMethod(c *gin.Context) bool {
	methods := map[string]bool{
		http.MethodGet:     true,
		http.MethodHead:    true,
		http.MethodOptions: true,
	}
	return methods[c.Request.Method]
}

func verifyCsrfToken(c *gin.Context, cookieName string, header string) bool {
	csrfCookie, _ := c.Cookie(cookieName)
	csrfHeader := c.GetHeader(header)
	if csrfCookie == csrfHeader &&
		csrfCookie != "" &&
		csrfHeader != "" {

		return true
	}

	return false
}
