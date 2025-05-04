package contracts

import (
	"context"
	"time"
)

type Queue interface {
	Enqueue(job Job, delay time.Duration)
	Start(ctx context.Context, n int)
	Close()
}
