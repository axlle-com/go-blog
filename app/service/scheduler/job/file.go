package job

import (
	"context"
	"github.com/axlle-com/blog/app/models/contracts"
	fileProvider "github.com/axlle-com/blog/pkg/file/provider"
	"time"
)

type DeleteFiles struct {
	start time.Time
	file  fileProvider.FileProvider
}

func NewDeleteFiles(file fileProvider.FileProvider) contracts.Job {
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
	return nil
}

func (j *DeleteFiles) GetName() string {
	return "DeleteFiles"
}

func (j *DeleteFiles) Duration() float64 {
	elapsed := time.Since(j.start)
	return float64(elapsed.Nanoseconds()) / 1e6
}
