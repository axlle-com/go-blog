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

	// UploadPath return /uploads/
	UploadPath() string
	DataFolder(...string) string
	Root() string
	SrcFolder() string
	SrcFolderBuilder(...string) string
	PublicFolderBuilder(...string) string
	Layout() string

	SMTPActive() bool
	SMTPPort() int
	SMTPHost() string
	SMTPUsername() string
	SMTPPassword() string

	NotifyEmail() string
}
