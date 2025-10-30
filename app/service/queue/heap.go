package queue

import (
	"time"

	"github.com/axlle-com/blog/app/models/contract"
)

type queueItem struct {
	job   contract.Job
	runAt time.Time
	index int // нужен для container/heap
}

type priorityQueue []*queueItem

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].runAt.Before(pq[j].runAt) }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i]; pq[i].index = i; pq[j].index = j }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*queueItem)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return it
}
