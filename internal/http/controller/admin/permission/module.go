package permission

import "go.uber.org/fx"

// NewModule provides permission admin handlers.
func NewModule() fx.Option {
	return fx.Module("http.controller.admin.permission")
}
