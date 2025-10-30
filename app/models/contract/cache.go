package contract

type Cache interface {
	AddCache(key, value string)
	DeleteCache(key string)
	GetUserKey(id uint) string
	AddUserSession(id uint, sessionID string)
	ResetUserSession(userID uint) error
	ResetUsersSession()
}
