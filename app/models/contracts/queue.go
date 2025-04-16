package contracts

import (
	"context"
	"time"
)

type Queue interface {
	Enqueue(job Job, delay time.Duration)
	StartWorkers(ctx context.Context, n int)
	Close()
}
