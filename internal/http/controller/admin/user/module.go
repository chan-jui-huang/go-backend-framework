package user

import "go.uber.org/fx"

// NewModule provides user admin handlers.
func NewModule() fx.Option {
	return fx.Module(
		"http.controller.admin.user",
		fx.Provide(
			NewUpdateUserRoleHandler,
		),
	)
}
