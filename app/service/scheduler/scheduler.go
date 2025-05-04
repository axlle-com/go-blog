package scheduler

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/app/service/scheduler/job"
	fileProvider "github.com/axlle-com/blog/pkg/file/provider"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	queue  contracts.Queue
	config contracts.Config
	cron   *cron.Cron

	file fileProvider.FileProvider
}

func NewScheduler(
	config contracts.Config,
	queue contracts.Queue,
	file fileProvider.FileProvider,
) contracts.Scheduler {
	return &Scheduler{
		config: config,
		queue:  queue,
		cron:   cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger))),
		file:   file,
	}
}

func (s *Scheduler) Start() {
	logger.Info("[Scheduler] Start")

	addFunc, err := s.cron.AddFunc("@every 1m", s.enqueueDeleteFiles)
	if err != nil {
		logger.Errorf("[Scheduler][Start] Func : %v,Error : %v", addFunc, err)
	}

	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	logger.Info("[Scheduler] Stop")
}

func (s *Scheduler) enqueueDeleteFiles() {
	s.queue.Enqueue(job.NewDeleteFiles(s.file), 0)
}
