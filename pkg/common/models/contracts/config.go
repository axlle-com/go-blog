package contracts

type Config interface {
	SetTestENV()
	IsTest() bool
	DBUrl() string
	DBUrlTest() string
	RedisHost() string
	KeyCookie() []byte
	KeyJWT() []byte
	SessionsName() string
	Port() string
	UploadPath() string
	SrcFolder() string
	SrcFolderBuilder(string) string
	UserSessionKey(string) string
	SessionKey(string) string
}
