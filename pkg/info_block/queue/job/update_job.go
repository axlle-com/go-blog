package job

import (
	"context"
	"encoding/json"
	"time"

	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/app/models/dto"
)

func NewInfoBlockJob(
	infoBlock *dto.Collection,
	action string,
) contracts.Job {
	return &UpdateInfoBlockJob{
		collection: infoBlock,
		action:     action,
		start:      time.Now(),
	}
}

type UpdateInfoBlockJob struct {
	data       []byte
	collection *dto.Collection
	action     string
	start      time.Time
}

func (j *UpdateInfoBlockJob) Run(ctx context.Context) error {
	return nil
}

func (j *UpdateInfoBlockJob) GetData() []byte {
	if j.data != nil {
		return j.data
	}

	raw, err := json.Marshal(j.collection)
	if err != nil {
		logger.Errorf("[info_block][UpdateInfoBlockJob][GetData] Error: %v", err)
		return nil
	}
	j.data = app.NewEnvelopeQueue().ConvertData(j.GetAction(), string(raw))
	return j.data
}

func (j *UpdateInfoBlockJob) GetName() string {
	return "info_blocks"
}

func (j *UpdateInfoBlockJob) GetQueue() string {
	return "info_blocks"
}

func (j *UpdateInfoBlockJob) GetAction() string {
	return j.action
}

func (j *UpdateInfoBlockJob) Duration() float64 {
	elapsed := time.Since(j.start)
	return float64(elapsed.Nanoseconds()) / 1e6
}
