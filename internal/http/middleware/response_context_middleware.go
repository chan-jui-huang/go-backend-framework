package middleware

import (
	"github.com/chan-jui-huang/go-backend-framework/v3/internal/http/response"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/gin-gonic/gin"
)

type ResponseContextMiddleware struct {
	debugMode bool
}

func NewResponseContextMiddleware(booterConfig *booter.Config) *ResponseContextMiddleware {
	return &ResponseContextMiddleware{
		debugMode: booterConfig.Debug,
	}
}

func (m *ResponseContextMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.SetDebugMode(c, m.debugMode)
		c.Next()
	}
}
