package job

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"scheduler.job",
	)
}
