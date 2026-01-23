package contract

import "time"

type Cache interface {
	AddCache(key, value string)
	AddCacheTTL(key, value string, ttl time.Duration)
	GetCache(key string) (string, bool)
	GetUserKey(id uint) string
	DeleteCache(key string)
	AddUserSession(id uint, sessionID string)
	ResetUserSession(userID uint) error
	ResetUsersSession()
}
