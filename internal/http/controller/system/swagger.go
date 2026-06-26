package system

import (
	_ "github.com/chan-jui-huang/go-backend-framework/v3/docs"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/booter"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerHandler struct {
	enabled bool
	handler gin.HandlerFunc
}

func NewSwaggerHandler(booterConfig *booter.Config) *SwaggerHandler {
	return &SwaggerHandler{
		enabled: booterConfig.Debug,
		handler: ginSwagger.WrapHandler(swaggerFiles.Handler),
	}
}

func (h *SwaggerHandler) Enabled() bool {
	return h.enabled
}

func (h *SwaggerHandler) Handle(c *gin.Context) {
	h.handler(c)
}
