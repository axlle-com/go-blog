package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/axlle-com/blog/app/logger"
	app "github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/axlle-com/blog/pkg/blog/models"
)

func NewUpdatePostJob(
	post *models.Post,
) contracts.Job {
	return &UpdatePostJob{
		post:  post,
		start: time.Now(),
	}
}

type UpdatePostJob struct {
	data  []byte
	post  *models.Post
	start time.Time
}

func (j *UpdatePostJob) Run(ctx context.Context) error {
	return nil
}

func (j *UpdatePostJob) GetData() []byte {
	if j.data != nil {
		return j.data
	}

	raw, err := json.Marshal(j.post)
	if err != nil {
		logger.Errorf("[AnalyticsJob][GetData] Error: %v", err)
		return nil
	}
	j.data = app.NewEnvelopeQueue().ConvertData(j.GetAction(), string(raw))
	return j.data
}

func (j *UpdatePostJob) GetName() string {
	return "posts"
}

func (j *UpdatePostJob) GetQueue() string {
	return "posts"
}

func (j *UpdatePostJob) GetAction() string {
	return "update"
}

func (j *UpdatePostJob) Duration() float64 {
	elapsed := time.Since(j.start)
	return float64(elapsed.Nanoseconds()) / 1e6
}
