package job

import (
	"context"
	"encoding/json"
	"time"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/app/models/dto"
)

func NewGalleryJob(
	collection *dto.Collection,
	action string,
) contracts.Job {
	return &UpdateGalleryJob{
		collection: collection,
		action:     action,
		start:      time.Now(),
	}
}

type UpdateGalleryJob struct {
	data       []byte
	collection *dto.Collection
	action     string
	start      time.Time
}

func (j *UpdateGalleryJob) Run(ctx context.Context) error {
	return nil
}

func (j *UpdateGalleryJob) GetData() []byte {
	if j.data != nil {
		return j.data
	}

	raw, err := json.Marshal(j.collection)
	if err != nil {
		logger.Errorf("[gallery][UpdateGalleryJob][GetData] Error: %v", err)
		return nil
	}
	j.data = models.NewEnvelopeQueue().ConvertData(j.GetAction(), string(raw))
	return j.data
}

func (j *UpdateGalleryJob) GetName() string {
	return "galleries"
}

func (j *UpdateGalleryJob) GetQueue() string {
	return "galleries"
}

func (j *UpdateGalleryJob) GetAction() string {
	return j.action
}

func (j *UpdateGalleryJob) Duration() float64 {
	elapsed := time.Since(j.start)
	return float64(elapsed.Nanoseconds()) / 1e6
}
