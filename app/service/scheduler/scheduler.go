package scheduler

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	appProvider "github.com/axlle-com/blog/app/models/provider"
	"github.com/axlle-com/blog/app/service/scheduler/job"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	queue  contract.Queue
	config contract.Config
	cron   *cron.Cron

	file appProvider.FileProvider
}

func NewScheduler(
	config contract.Config,
	queue contract.Queue,
	file appProvider.FileProvider,
) contract.Scheduler {
	return &Scheduler{
		config: config,
		queue:  queue,
		cron:   cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger))),
		file:   file,
	}
}

func (s *Scheduler) Start() {
	logger.Info("[scheduler] Start")

	addFunc, err := s.cron.AddFunc("@every 1h", s.enqueueDeleteFiles)
	if err != nil {
		logger.Errorf("[scheduler][Start] Func : %v,Error : %v", addFunc, err)
	}

	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	logger.Info("[scheduler][Stop] Stop")
}

func (s *Scheduler) enqueueDeleteFiles() {
	s.queue.Enqueue(job.NewDeleteFiles(s.file), 0)
}
