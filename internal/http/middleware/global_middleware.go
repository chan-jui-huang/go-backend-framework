package middleware

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/gin-gonic/gin"
)

func AttachGlobalMiddleware(router *gin.Engine) {
	handlerFunctions := []gin.HandlerFunc{
		AccessLogger(),
		Recover(),
		VerifyCsrfToken(deps.CsrfConfig()),
	}

	router.Use(
		handlerFunctions...,
	)
}
