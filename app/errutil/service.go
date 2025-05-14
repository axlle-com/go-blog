package errutil

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type ErrUtil struct {
	counts map[string]int
	mu     sync.RWMutex
}

func New() *ErrUtil {
	return &ErrUtil{counts: make(map[string]int)}
}

func (e *ErrUtil) Add(err error) {
	if err == nil {
		return
	}
	e.mu.Lock()
	e.counts[err.Error()]++
	e.mu.Unlock()
}

func (e *ErrUtil) Merge(errs ...error) error {
	merged := New()
	for _, err := range errs {
		merged.Add(err)
	}
	return merged.Error()
}

func (e *ErrUtil) Error() error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if len(e.counts) == 0 {
		return nil
	}

	parts := make([]string, 0, len(e.counts))
	for msg, count := range e.counts {
		if count > 1 {
			parts = append(parts, fmt.Sprintf("%s [x%d]", msg, count))
		} else {
			parts = append(parts, msg)
		}
	}
	return errors.New(strings.Join(parts, "; "))
}
