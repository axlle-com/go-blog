package contracts

type Config interface {
	Port() string
	LogLevel() int
	SetTestENV()
	IsTest() bool

	DBUrl() string
	DBUrlTest() string

	RedisHost() string
	RedisPassword() string
	StoreIsRedis() bool

	KeyCookie() []byte
	KeyJWT() []byte
	SessionsName() string
	SessionKey(string) string
	UserSessionKey(string) string

	UploadPath() string
	RuntimeFolder(s string) string
	SrcFolder() string
	SrcFolderBuilder(string) string

	SMTPActive() bool
	SMTPPort() int
	SMTPHost() string
	SMTPUsername() string
	SMTPPassword() string
}
