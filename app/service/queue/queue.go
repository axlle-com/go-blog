package queue

import (
	"container/heap"
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
)

type Queue struct {
	mu            *sync.Mutex
	cond          *sync.Cond    // ожидание при пустой очереди
	wake          chan struct{} // «появилась более ранняя задача»
	priorityQueue priorityQueue
	closing       bool
	handlers      map[string][]contract.QueueHandler

	wg sync.WaitGroup // ждём воркеры при Close()
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

func (q *Queue) SetHandlers(handlers map[string][]contract.QueueHandler) {
	q.mu.Lock()
	q.handlers = handlers
	q.mu.Unlock()
}

func (q *Queue) Enqueue(job contract.Job, delay time.Duration) {
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
	if ctx == nil || n <= 0 {
		return
	}

	for i := 0; i < n; i++ {
		q.wg.Add(1)
		go func() {
			defer q.wg.Done()
			q.worker(ctx)
		}()
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

	q.wg.Wait()
	logger.Info("[queue] Close")
}

func (q *Queue) next() (*queueItem, bool) {
	if len(q.priorityQueue) == 0 {
		return nil, false
	}
	return q.priorityQueue[0], true
}

func (q *Queue) worker(ctx context.Context) {
	logger.Info("[queue][worker] Start")
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("[queue][worker] panic: %v\n%s", r, debug.Stack())
		}
	}()

	for {
		q.mu.Lock()

		item, ok := q.next()
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
		if now.Before(item.runAt) {
			wait := time.Until(item.runAt)
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

		handlers := q.safeHandlersFor(item.job.GetQueue())
		if len(handlers) == 0 {
			logger.Errorf("[queue][%s] handlers not found, Data: %s", item.job.GetName(), string(item.job.GetData()))
			continue
		}

		// Общий таймаут на job (чтобы один job не висел вечно)
		jobTimeout := 30 * time.Second
		jobContext, cancel := context.WithTimeout(ctx, jobTimeout)

		for _, handler := range handlers {
			// Таймаут на каждый handler короче общего
			perHandlerTimeout := 15 * time.Second
			if err := runHandlerSafe(jobContext, handler, item.job.GetData(), perHandlerTimeout); err != nil {
				logger.Errorf("[queue][%s] handler error: %v", item.job.GetName(), err)
				// Здесь нужно включить retry/Backoff/DLQ (см. комментарии ниже)
			}
		}
		cancel()

		logger.Debugf("[queue][%s] Duration: %.2fms", item.job.GetName(), item.job.Duration())
	}
}

func (q *Queue) safeHandlersFor(queueName string) []contract.QueueHandler {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.handlers[queueName]
}

// runHandlerSafe запускает handler с перехватом panics и таймаутом.
func runHandlerSafe(ctx context.Context, handler contract.QueueHandler, payload []byte, timeout time.Duration) error {
	handlerContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errCh <- fmt.Errorf("panic: %v\n%s", r, debug.Stack())
			}
		}()
		handler.Run(payload) // старый контракт: Run([]byte)
		errCh <- nil
	}()

	select {
	case <-handlerContext.Done():
		return handlerContext.Err() // DeadlineExceeded / Canceled
	case err := <-errCh:
		return err
	}
}
