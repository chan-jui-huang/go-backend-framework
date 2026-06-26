package scheduler

import (
	schedulerpkg "github.com/chan-jui-huang/go-backend-package/v2/pkg/scheduler"
	"go.uber.org/zap"
)

type Scheduler struct {
	logger *zap.Logger
	jobs   []schedulerpkg.Job
}

func NewScheduler(logger *zap.Logger, jobs []schedulerpkg.Job) *Scheduler {
	return &Scheduler{
		logger: logger,
		jobs:   jobs,
	}
}

func (s *Scheduler) Start() {
	schedulerpkg.Scheduler.QueueJobs(s.jobs)
	schedulerpkg.Scheduler.Start()
	s.logger.Info("scheduler is started")
}

func (s *Scheduler) Stop() {
	<-schedulerpkg.Scheduler.Stop().Done()
	s.logger.Info("scheduler is stopped")
}
