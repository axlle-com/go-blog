package contract

import "context"

type Job interface {
	Run(ctx context.Context) error
	GetData() []byte
	GetName() string
	GetQueue() string
	GetAction() string
	Duration() float64
}
