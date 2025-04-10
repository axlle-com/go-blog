package config

import (
	"fmt"
	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type config struct {
	*sync.Mutex
	rootDir  string
	env      string
	port     string
	logLevel int

	redisHost string
	redisPort string

	dialector      string
	dbHost         string
	dbPort         string
	dbName         string
	dbUser         string
	dbPassword     string
	dbNameTest     string
	dbUserTest     string
	dbPasswordTest string

	keyJWT    string
	keyCookie string

	uploadsPath string
	srcFolder   string
}

var (
	once     sync.Once
	instance *config
)

func LoadConfig() (err error) {
	once.Do(func() {
		instance = &config{Mutex: &sync.Mutex{}}

		rootDir, err := instance.root()
		if err != nil {
			err = fmt.Errorf("ошибка определения корневой директории: %w", err)
			return
		}

		err = godotenv.Load(filepath.Join(rootDir, ".env"))
		if err != nil {
			err = fmt.Errorf("ошибка при загрузке .env файла: %v", err)
			return
		}

		instance.env = getEnv("ENV", "dev")
		instance.port = getEnv("PORT", "3000")

		logLevel := getEnv("LOG_LEVEL", "6")
		instance.logLevel, err = strconv.Atoi(logLevel)
		if err != nil {
			instance.logLevel = 6
		}

		instance.keyJWT = getEnv("KEY_JWT", "")
		instance.keyCookie = getEnv("KEY_COOKIE", "")

		instance.redisHost = getEnv("REDIS_HOST", "127.0.0.1")
		instance.redisPort = getEnv("REDIS_PORT", "6380")

		instance.dialector = getEnv("DIALECTOR", "postgres")

		instance.dbHost = getEnv("DB_HOST", "127.0.0.1")
		instance.dbPort = getEnv("DB_PORT", "5432")
		instance.dbName = getEnv("DB_NAME", "cms_main")
		instance.dbUser = getEnv("DB_USER", "postgres")
		instance.dbPassword = getEnv("DB_PASSWORD", "secret")

		instance.dbNameTest = getEnv("DB_NAME_TEST", "cms_test")
		instance.dbUserTest = getEnv("DB_USER_TEST", "postgres")
		instance.dbPasswordTest = getEnv("DB_PASSWORD_TEST", "secret")

		instance.uploadsPath = getEnv("FILE_UPLOADS_PATH", "/public/uploads/")
		instance.srcFolder = getEnv("FILE_SRC_FOLDER", "src")
	})
	return
}

func Config() contracts.Config {
	if instance == nil {
		if err := LoadConfig(); err != nil {
			panic("Ошибка загрузки конфигурации: " + err.Error())
		}
	}
	return instance
}

func (c *config) SetTestENV() {
	c.Lock()
	defer c.Unlock()

	c.env = "test"
}

func (c *config) IsTest() bool {
	return c.env == "test"
}

func (c *config) DBUrl() string {
	var dsn string
	if c.dialector == "postgres" {
		dsn = "host=" + c.dbHost +
			" user=" + c.dbUser +
			" password=" + c.dbPassword +
			" dbname=" + c.dbName +
			" port=" + c.dbPort +
			" sslmode=disable TimeZone=Europe/Moscow"
	}
	return dsn
}

func (c *config) DBUrlTest() string {
	var dsn string
	if c.dialector == "postgres" {
		dsn = "host=" + c.dbHost +
			" user=" + c.dbUserTest +
			" password=" + c.dbPasswordTest +
			" dbname=" + c.dbNameTest +
			" port=" + c.dbPort +
			" sslmode=disable TimeZone=Europe/Moscow"
	}
	return dsn
}

func (c *config) RedisHost() string {
	if c.IsTest() {
		s := c.redisHost + ":" + c.redisPort
		return s
	}
	s := c.redisHost + ":" + c.redisPort
	return s
}

func (c *config) KeyCookie() []byte {
	if c.IsTest() {
		s := []byte(c.keyCookie)
		return s
	}
	s := []byte(c.keyCookie)
	return s
}

func (c *config) KeyJWT() []byte {
	if c.IsTest() {
		s := []byte(c.keyJWT)
		return s
	}
	s := []byte(c.keyJWT)
	return s
}

func (c *config) SessionsName() string {
	if c.IsTest() {
		s := "web_session_test"
		return s
	}
	s := "web_session"
	return s
}

func (c *config) Port() string {
	if c.IsTest() {
		return c.port
	}
	return c.port
}

func (c *config) UploadPath() string {
	if c.IsTest() {
		return c.uploadsPath + "test/"
	}
	return c.uploadsPath
}

func (c *config) SrcFolder() string {
	if c.IsTest() {
		root, err := c.root()
		if err != nil {
			return ""
		}
		return root + "/" + c.srcFolder
	}
	return c.srcFolder
}

func (c *config) SrcFolderBuilder(s string) string {
	return c.SrcFolder() + "/" + strings.TrimLeft(s, " /")
}

func (c *config) UserSessionKey(s string) string {
	key := "user_session_"
	if c.IsTest() {
		return key + "test_" + s
	}
	return key + s
}

func (c *config) SessionKey(s string) string {
	return "session_" + s
}

func (c *config) LogLevel() int {
	return c.logLevel
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *config) root() (string, error) {
	c.Lock()
	defer c.Unlock()

	if c.rootDir != "" {
		return c.rootDir, nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			c.rootDir = dir
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("не удалось найти корневую директорию модуля")
		}
		dir = parent
	}
}
