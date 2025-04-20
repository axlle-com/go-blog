package models

import (
	"fmt"
	"github.com/axlle-com/blog/app/models/contracts"
	"strings"
	"sync"
)

type inMemoryCache struct {
	mu     sync.RWMutex
	data   map[string]string
	config contracts.Config
}

func NewInMemoryCache(cfg contracts.Config) contracts.Cache {
	return &inMemoryCache{
		data:   make(map[string]string),
		config: cfg,
	}
}

func (m *inMemoryCache) AddCache(key, value string) {
	m.mu.Lock()
	m.data[key] = value
	m.mu.Unlock()
}

func (m *inMemoryCache) DeleteCache(key string) {
	m.mu.Lock()
	delete(m.data, key)
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
	delete(m.data, userKey)
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
		}
	}
}
