package analytic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/pkg/analytic/provider"
)

func NewAnalyticsJob(
	event AnalyticsEvent,
	analyticProvider provider.AnalyticProvider,
) *AnalyticsJob {
	return &AnalyticsJob{
		event:            &event,
		analyticProvider: analyticProvider,
		start:            time.Now(),
	}
}

type AnalyticsJob struct {
	start            time.Time
	event            *AnalyticsEvent
	analyticProvider provider.AnalyticProvider
}

func (j *AnalyticsJob) Run(ctx context.Context) error {
	_, err := j.analyticProvider.SaveForm(j.event)
	if err != nil {
		return err
	}
	return nil
}

func (j *AnalyticsJob) GetData() []byte {
	raw, err := json.Marshal(j.event)
	if err != nil {
		logger.Errorf("[AnalyticsJob][GetData] Error: %v", err)
		return nil
	}

	return raw
}

func (j *AnalyticsJob) GetName() string {
	return "Analytics"
}

func (j *AnalyticsJob) Duration() float64 {
	elapsed := time.Since(j.start)
	return float64(elapsed.Nanoseconds()) / 1e6
}

func toMap(v any) (map[string]any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}
