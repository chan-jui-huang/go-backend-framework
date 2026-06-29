package middleware

import (
	"net/http"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/config"
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/random"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type CsrfMiddleware struct {
	logger *zap.Logger
	config *config.CsrfConfig
}

func NewCsrfMiddleware(
	logger *zap.Logger,
	config *config.CsrfConfig,
) *CsrfMiddleware {
	return &CsrfMiddleware{
		logger: logger,
		config: config,
	}
}

func (m *CsrfMiddleware) Handle() gin.HandlerFunc {
	skipPaths := map[string]bool{
		"/skip-path": true,
	}

	return func(c *gin.Context) {
		setCsrfToken(c, m.config)
		if isReadingHttpMethod(c) ||
			skipPaths[c.Request.URL.Path] ||
			verifyCsrfToken(c, m.config.Cookie.Name, m.config.Header) {
			c.Next()
			return
		}

		errResp := response.NewErrorResponse(response.Forbidden, errors.New("csrf token mismatch"), nil, response.DebugMode(c))
		m.logger.Warn(response.Forbidden, errResp.MakeLogFields(c)...)
		c.AbortWithStatusJSON(errResp.StatusCode(), errResp)
	}
}

func setCsrfToken(c *gin.Context, config *config.CsrfConfig) {
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
