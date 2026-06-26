package scheduler

import (
	"context"

	"github.com/chan-jui-huang/go-backend-framework/v3/internal/scheduler/job"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"scheduler",
		job.NewModule(),
		fx.Provide(
			fx.Annotate(
				NewScheduler,
				fx.ParamTags(``, `group:"scheduler_jobs"`),
			),
		),
		fx.Invoke(
			fx.Annotate(
				func() {},
				fx.OnStart(SchedulerOnStart),
				fx.OnStop(SchedulerOnStop),
			),
		),
	)
}

func SchedulerOnStart(_ context.Context, scheduler *Scheduler) error {
	scheduler.Start()

	return nil
}

func SchedulerOnStop(_ context.Context, scheduler *Scheduler) error {
	scheduler.Stop()

	return nil
}
