package service

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/axlle-com/blog/pkg/settings/models"
	"github.com/axlle-com/blog/pkg/settings/repository"
)

type Repository = repository.Repository

type Service struct {
	repo  Repository
	mu    sync.RWMutex
	cache map[string]cachedItem
	ttl   time.Duration
}

type cachedItem struct {
	s   models.Setting
	exp time.Time
}

func NewService(r Repository) *Service {
	return &Service{repo: r, cache: make(map[string]cachedItem), ttl: 5 * time.Minute}
}

func (s *Service) cacheKey(ns, key, scope string) string { return ns + "|" + key + "|" + scope }

func (s *Service) getCached(ns, key, scope string) (models.Setting, bool) {
	k := s.cacheKey(ns, key, scope)
	now := time.Now()
	s.mu.RLock()
	if it, ok := s.cache[k]; ok && it.exp.After(now) {
		s.mu.RUnlock()
		return it.s, true
	}
	s.mu.RUnlock()

	st, err := s.repo.Get(ns, key, scope)
	if err != nil || st == nil {
		return models.Setting{}, false
	}

	s.mu.Lock()
	s.cache[k] = cachedItem{s: *st, exp: now.Add(s.ttl)}
	s.mu.Unlock()
	return *st, true
}

func (s *Service) Invalidate(ns, key, scope string) {
	k := s.cacheKey(ns, key, scope)
	s.mu.Lock()
	delete(s.cache, k)
	s.mu.Unlock()
}

func (s *Service) GetString(ns, key, scope string) (string, bool) {
	st, ok := s.getCached(ns, key, scope)
	if !ok {
		return "", false
	}
	var v string
	_ = json.Unmarshal(st.Value, &v)
	if v == "" {
		return "", false
	}
	return v, true
}

func (s *Service) GetBool(ns, key, scope string) (bool, bool) {
	st, ok := s.getCached(ns, key, scope)
	if !ok {
		return false, false
	}
	var v bool
	_ = json.Unmarshal(st.Value, &v)
	return v, true
}

func (s *Service) GetJSON(ns, key, scope string, out any) bool {
	st, ok := s.getCached(ns, key, scope)
	if !ok {
		return false
	}
	if err := json.Unmarshal(st.Value, out); err != nil {
		return false
	}
	return true
}
