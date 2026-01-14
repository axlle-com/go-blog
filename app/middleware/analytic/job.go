package analytic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models"
	"github.com/axlle-com/blog/app/service/queue"
)

func NewAnalyticsJob(
	event AnalyticsEvent,
) *AnalyticsJob {
	return &AnalyticsJob{
		event: &event,
		start: time.Now(),
	}
}

type AnalyticsJob struct {
	start time.Time
	event *AnalyticsEvent
	data  []byte
}

func (j *AnalyticsJob) Run(ctx context.Context) error {
	return nil
}

func (j *AnalyticsJob) GetData() []byte {
	if j.data != nil {
		return j.data
	}

	raw, err := json.Marshal(j.event)
	if err != nil {
		logger.Errorf("[AnalyticsJob][GetData] Error: %v", err)
		return nil
	}
	j.data = models.NewEnvelopeQueue().ConvertData(queue.Create, string(raw))
	return j.data
}

func (j *AnalyticsJob) GetName() string {
	return "analytics"
}

func (j *AnalyticsJob) GetQueue() string {
	return "analytics"
}

func (j *AnalyticsJob) GetAction() string {
	return queue.Create
}

func (j *AnalyticsJob) Duration() float64 {
	elapsed := time.Since(j.start)
	return float64(elapsed.Nanoseconds()) / 1e6
}
