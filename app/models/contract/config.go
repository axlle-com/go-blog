package contract

import "gorm.io/gorm"

type Config interface {
	AppHost() string
	Port() string
	LogLevel() int
	SetTestENV()
	IsTest() bool
	IsLocal() bool

	DBUrl() string
	DBUrlTest() string
	SetGORM(*gorm.DB)
	GetGORM() *gorm.DB

	RedisHost() string
	RedisPassword() string
	StoreIsRedis() bool

	KeyCookie() []byte
	KeyJWT() []byte
	SessionsName() string
	SessionKey(string) string
	UserSessionKey(string) string

	Layout() string
	Root() string
	UploadPath() string
	DataFolder(...string) string
	SrcFolder(...string) string
	StaticFolder(...string) string

	SMTPActive() bool
	SMTPPort() int
	SMTPHost() string
	SMTPUsername() string
	SMTPPassword() string

	NotifyEmail() string
}
