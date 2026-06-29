package scheduler

import (
	"github.com/chan-jui-huang/go-backend-package/v2/pkg/scheduler"
	"go.uber.org/zap"
)

type Scheduler struct {
	logger *zap.Logger
	jobs   []scheduler.Job
}

func NewScheduler(logger *zap.Logger, jobs []scheduler.Job) *Scheduler {
	return &Scheduler{
		logger: logger,
		jobs:   jobs,
	}
}

func (s *Scheduler) Start() {
	scheduler.Scheduler.QueueJobs(s.jobs)
	scheduler.Scheduler.Start()
	s.logger.Info("scheduler is started")
}

func (s *Scheduler) Stop() {
	<-scheduler.Scheduler.Stop().Done()
	s.logger.Info("scheduler is stopped")
}
