package httpapi

import "go.uber.org/fx"

// NewModule provides HTTP API admin handlers.
func NewModule() fx.Option {
	return fx.Module("http.controller.admin.http_api", fx.Provide(NewSearchHandler))
}
