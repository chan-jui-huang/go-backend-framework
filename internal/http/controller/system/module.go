package system

import "go.uber.org/fx"

// NewModule provides system handlers.
func NewModule() fx.Option {
	return fx.Module(
		"http.controller.system",
		fx.Provide(
			NewPingHandler,
			NewSwaggerHandler,
		),
	)
}
