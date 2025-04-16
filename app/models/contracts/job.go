package contracts

import "context"

type Job interface {
	Run(ctx context.Context) error
}
