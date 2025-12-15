package job

import (
	"context"
	"time"

	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contract"
	appProvider "github.com/axlle-com/blog/app/models/provider"
)

type DeleteFiles struct {
	start time.Time
	file  appProvider.FileProvider
}

func NewDeleteFiles(file appProvider.FileProvider) contract.Job {
	return &DeleteFiles{
		file:  file,
		start: time.Now(),
	}
}

func (j *DeleteFiles) Run(ctx context.Context) error {
	j.file.RevisionReceived()
	return nil
}

func (j *DeleteFiles) GetData() []byte {
	return app.NewEnvelopeQueue().ConvertData(j.GetAction(), "")
}

func (j *DeleteFiles) GetName() string {
	return "files"
}

func (j *DeleteFiles) GetQueue() string {
	return "files"
}

func (j *DeleteFiles) GetAction() string {
	return "revision_received"
}

func (j *DeleteFiles) Duration() float64 {
	elapsed := time.Since(j.start)
	return float64(elapsed.Nanoseconds()) / 1e6
}
