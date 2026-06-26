package middleware

import (
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"http.middleware",
		fx.Provide(
			fx.Annotate(
				NewAccessLogMiddleware,
				fx.ParamTags(`name:"logger.access"`),
			),
			NewAuthenticationMiddleware,
			NewAuthorizationMiddleware,
			NewRateLimitMiddleware,
			NewRecoverMiddleware,
			NewCsrfMiddleware,
			NewResponseContextMiddleware,
			NewGlobalMiddlewares,
		),
	)
}
