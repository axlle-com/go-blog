package analytic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func NewAnalyticsJob(event AnalyticsEvent) *AnalyticsJob {
	return &AnalyticsJob{event}
}

type AnalyticsJob struct {
	ev AnalyticsEvent
}

func (p *AnalyticsJob) Run(ctx context.Context) error {
	str, _ := toMap(p.ev)
	fmt.Println(time.Now().Format("15:04:05"), ":", str)
	return nil
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
