package registrar

import (
	"context"

	"github.com/chan-jui-huang/go-backend-framework/v2/internal/scheduler"
)

func SchedulerOnStart(context.Context) error {
	scheduler.Start()

	return nil
}

func SchedulerOnStop(context.Context) error {
	scheduler.Stop()

	return nil
}
