package config

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type config struct {
	*sync.Mutex
	env            string
	redisHost      string
	redisPort      string
	port           string
	dialector      string
	dbHost         string
	dbPort         string
	dbName         string
	dbUser         string
	dbPassword     string
	dbNameTest     string
	dbUserTest     string
	dbPasswordTest string
	keyJWT         string
	keyCookie      string
	uploadsPath    string
	uploadsFolder  string
}

var (
	once     sync.Once
	instance *config
)

func LoadConfig() (err error) {
	once.Do(func() {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Ошибка при загрузке .env файла: %v", err)
		}

		instance = &config{Mutex: &sync.Mutex{}}
		instance.env = os.Getenv("ENV")

		instance.port = os.Getenv("PORT")

		instance.keyJWT = os.Getenv("KEY_JWT")
		instance.keyCookie = os.Getenv("KEY_COOKIE")

		instance.redisHost = os.Getenv("REDIS_HOST")
		instance.redisPort = os.Getenv("REDIS_PORT")

		instance.dialector = os.Getenv("DIALECTOR")
		if instance.dialector == "" {
			instance.dialector = "postgres"
		}

		instance.dbHost = os.Getenv("POSTGRES_HOST")
		instance.dbPort = os.Getenv("POSTGRES_PORT")
		instance.dbUser = os.Getenv("POSTGRES_USER")
		instance.dbName = os.Getenv("POSTGRES_DB")
		instance.dbPassword = os.Getenv("POSTGRES_PASSWORD")

		instance.dbUserTest = os.Getenv("POSTGRES_USER_TEST")
		instance.dbNameTest = os.Getenv("POSTGRES_DB_TEST")
		instance.dbPasswordTest = os.Getenv("POSTGRES_PASSWORD_TEST")

		instance.uploadsPath = os.Getenv("FILE_UPLOADS_PATH")
		if instance.uploadsPath == "" {
			instance.uploadsPath = "/public/uploads/"
		}
		instance.uploadsFolder = os.Getenv("FILE_UPLOADS_FOLDER")
		if instance.uploadsFolder == "" {
			instance.uploadsFolder = "src"
		}
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
	c.env = "test"
	c.Unlock()
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
		return root + "/" + c.uploadsFolder
	}
	return c.uploadsFolder
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

func (c *config) root() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("не удалось найти корневую директорию модуля")
		}

		dir = parent
	}
}
