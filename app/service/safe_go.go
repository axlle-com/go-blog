package service

import (
	"runtime/debug"
	"sync"

	"github.com/axlle-com/blog/app/logger"
)

func SafeGo(wg *sync.WaitGroup, fn func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("goroutine panic: %v\n%s", r, debug.Stack())
			}
		}()
		fn()
	}()
}
