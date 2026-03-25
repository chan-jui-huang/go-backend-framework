package scheduler

import (
	"github.com/chan-jui-huang/go-backend-framework/v2/internal/deps"
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/scheduler"
)

func backlogJobs() {
	scheduler.Scheduler.BacklogJobs(map[string]scheduler.Job{
		// "example": job.NewExampleJob(),
	})
}

func Start() {
	backlogJobs()
	scheduler.Scheduler.Start()

	logger := deps.Logger()
	logger.Info("scheduler is started")
}

func Stop() {
	<-scheduler.Scheduler.Stop().Done()

	logger := deps.Logger()
	logger.Info("scheduler is stopped")
}
