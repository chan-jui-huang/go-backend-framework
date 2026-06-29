package middleware

import "github.com/gin-gonic/gin"

type GlobalMiddlewares struct {
	accessLog       *AccessLogMiddleware
	recover         *RecoverMiddleware
	csrf            *CsrfMiddleware
	responseContext *ResponseContextMiddleware
}

func NewGlobalMiddlewares(
	accessLog *AccessLogMiddleware,
	recover *RecoverMiddleware,
	csrf *CsrfMiddleware,
	responseContext *ResponseContextMiddleware,
) *GlobalMiddlewares {
	return &GlobalMiddlewares{
		accessLog:       accessLog,
		recover:         recover,
		csrf:            csrf,
		responseContext: responseContext,
	}
}

func (m *GlobalMiddlewares) Attach(router *gin.Engine) {
	router.Use(
		m.accessLog.Handle(),
		m.recover.Handle(),
		m.responseContext.Handle(),
		m.csrf.Handle(),
	)
}
