package cache

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
)

type inMemoryCache struct {
	mu     sync.RWMutex
	data   map[string]string
	exp    map[string]time.Time
	config contract.Config
}

func NewInMemoryCache(cfg contract.Config) contract.Cache {
	logger.Info("[Cache] Using inMemoryCache")
	return &inMemoryCache{
		data:   make(map[string]string),
		exp:    make(map[string]time.Time),
		config: cfg,
	}
}

func (m *inMemoryCache) AddCache(key, value string) {
	m.mu.Lock()
	m.data[key] = value
	delete(m.exp, key)
	m.mu.Unlock()
}

func (m *inMemoryCache) AddCacheTTL(key, value string, ttl time.Duration) {
	m.mu.Lock()
	m.data[key] = value
	if ttl > 0 {
		m.exp[key] = time.Now().Add(ttl)
	} else {
		delete(m.exp, key)
	}
	m.mu.Unlock()
}

func (m *inMemoryCache) GetCache(key string) (string, bool) {
	m.mu.RLock()
	v, ok := m.data[key]
	exp, hasExp := m.exp[key]
	m.mu.RUnlock()

	if !ok {
		return "", false
	}

	if hasExp && time.Now().After(exp) {
		m.mu.Lock()
		exp2, hasExp2 := m.exp[key]
		if hasExp2 && time.Now().After(exp2) {
			delete(m.data, key)
			delete(m.exp, key)
			m.mu.Unlock()
			return "", false
		}
		m.mu.Unlock()
	}

	return v, true
}

func (m *inMemoryCache) DeleteCache(key string) {
	m.mu.Lock()
	delete(m.data, key)
	delete(m.exp, key)
	m.mu.Unlock()
}

func (m *inMemoryCache) GetUserKey(id uint) string {
	return fmt.Sprintf(m.config.UserSessionKey("%d"), id)
}

func (m *inMemoryCache) AddUserSession(id uint, sessionID string) {
	m.AddCache(m.GetUserKey(id), sessionID)
}

func (m *inMemoryCache) ResetUserSession(userID uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	userKey := m.GetUserKey(userID)

	sessionID, ok := m.data[userKey]
	if !ok {
		return nil
	}

	delete(m.data, sessionID)
	delete(m.exp, sessionID)

	delete(m.data, userKey)
	delete(m.exp, userKey)

	return nil
}

func (m *inMemoryCache) ResetUsersSession() {
	m.mu.Lock()
	defer m.mu.Unlock()

	prefix1 := m.config.UserSessionKey("")
	prefix2 := m.config.SessionKey("")
	for k := range m.data {
		if strings.HasPrefix(k, prefix1) || strings.HasPrefix(k, prefix2) {
			delete(m.data, k)
			delete(m.exp, k)
		}
	}
}
