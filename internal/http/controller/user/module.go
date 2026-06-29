package user

import "go.uber.org/fx"

// NewModule provides user handlers.
func NewModule() fx.Option {
	return fx.Module(
		"http.controller.user",
		fx.Provide(
			NewGetMeHandler,
			NewLoginHandler,
			NewRegisterHandler,
			NewUpdateHandler,
			NewUpdatePasswordHandler,
		),
	)
}
