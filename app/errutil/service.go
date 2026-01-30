package errutil

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type ErrUtil struct {
	mu     sync.RWMutex
	counts map[string]int
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

func (e *ErrUtil) AddFormat(format string, args ...any) {
	e.Add(fmt.Errorf(format, args...))
}

func (e *ErrUtil) Any() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.counts) > 0
}

func (e *ErrUtil) Count() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.counts)
}

// Не сбрасывает, просто возвращает ошибку (или nil)
func (e *ErrUtil) Error() error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if len(e.counts) == 0 {
		return nil
	}
	keys := make([]string, 0, len(e.counts))
	for k := range e.counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, msg := range keys {
		if c := e.counts[msg]; c > 1 {
			parts = append(parts, fmt.Sprintf("%s [x%d]", msg, c))
		} else {
			parts = append(parts, msg)
		}
	}

	return errors.New(strings.Join(parts, "; "))
}

// ErrorAndReset Собирает ошибку и СБРАСЫВАЕТ состояние
func (e *ErrUtil) ErrorAndReset() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.counts) == 0 {
		return nil
	}
	keys := make([]string, 0, len(e.counts))
	for k := range e.counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, msg := range keys {
		if c := e.counts[msg]; c > 1 {
			parts = append(parts, fmt.Sprintf("%s [x%d]", msg, c))
		} else {
			parts = append(parts, msg)
		}
	}

	e.counts = make(map[string]int)

	return errors.New(strings.Join(parts, "; "))
}

// Merge Схлопывает набор ошибок в одну (игнорируя nil)
func (e *ErrUtil) Merge(errs ...error) error {
	m := New()

	for _, err := range errs {
		m.Add(err)
	}

	return m.Error()
}
