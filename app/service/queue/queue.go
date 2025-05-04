package queue

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contracts"
)

type Queue struct {
	mu            *sync.Mutex
	cond          *sync.Cond    // ожидание, когда очередь пуста
	wake          chan struct{} // сигнал «появилась более ранняя задача»
	priorityQueue priorityQueue
	closing       bool
}

func NewQueue() *Queue {
	mu := &sync.Mutex{}
	q := &Queue{
		mu:            mu,
		wake:          make(chan struct{}, 1),
		priorityQueue: make(priorityQueue, 0),
	}
	q.cond = sync.NewCond(mu)
	return q
}

// Enqueue ставит задачу в очередь; delay==0 — немедленно.
// Если воркер спит до будущего runAt, посылаем wake‑сигнал,
// чтобы он пересмотрел приоритеты.
func (q *Queue) Enqueue(job contracts.Job, delay time.Duration) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closing {
		return
	}

	it := &queueItem{job: job, runAt: time.Now().Add(delay)}
	heap.Push(&q.priorityQueue, it)

	// Разбудим одного спящего воркера, если очередь была пуста.
	q.cond.Signal()

	// Неблокирующий wake — если канал полон, значит воркеры уже проснутся.
	select {
	case q.wake <- struct{}{}:
	default:
	}
}

func (q *Queue) Start(ctx context.Context, n int) {
	if ctx == nil {
		return
	}

	for i := 0; i < n; i++ {
		go q.worker(ctx)
	}
}

func (q *Queue) Close() {
	q.mu.Lock()
	q.closing = true
	q.cond.Broadcast()
	// wake может быть пуст, но broadcast гарантирует, что воркеры проверят closing.
	select {
	case q.wake <- struct{}{}:
	default:
	}
	q.mu.Unlock()
	logger.Info("[Queue] Close")
}

func (q *Queue) next() (*queueItem, bool) {
	if len(q.priorityQueue) == 0 {
		return nil, false
	}
	return q.priorityQueue[0], true
}

func (q *Queue) worker(ctx context.Context) {
	logger.Info("[Queue][worker] Start")
	for {
		q.mu.Lock()

		it, ok := q.next()
		if !ok {
			if q.closing {
				q.mu.Unlock()
				return
			}
			q.cond.Wait()
			q.mu.Unlock()
			continue
		}

		now := time.Now()
		if now.Before(it.runAt) {
			wait := time.Until(it.runAt)
			q.mu.Unlock()

			timer := time.NewTimer(wait)
			select {
			case <-timer.C: // наступило время runAt
			case <-q.wake: // появилась более ранняя задача — пересмотреть очередь
			case <-ctx.Done():
				timer.Stop()
				return
			}
			timer.Stop()
			continue // вернёмся, захватим лок и пересчитаем состояние
		}

		// Пора исполнять
		heap.Pop(&q.priorityQueue)
		q.mu.Unlock()

		if err := it.job.Run(ctx); err != nil {
			logger.Errorf("[Queue][%s] Error: %v, Data: %s", it.job.GetName(), err, string(it.job.GetData()))
		}
		logger.Debugf("[Queue][%s] Duration: %.2fms", it.job.GetName(), it.job.Duration())
	}
}
