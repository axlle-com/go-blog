package contracts

type Config interface {
	Port() string
	LogLevel() int
	SetTestENV()
	IsTest() bool
	DBUrl() string
	DBUrlTest() string
	RedisHost() string
	KeyCookie() []byte
	KeyJWT() []byte
	SessionsName() string
	SessionKey(string) string
	UserSessionKey(string) string
	UploadPath() string
	SrcFolder() string
	SrcFolderBuilder(string) string
}
