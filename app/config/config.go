package config

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/axlle-com/blog/app/models/contracts"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type config struct {
	*sync.Mutex
	rootDir  string
	env      string
	appHost  string
	appPort  string
	logLevel int

	store         string
	redisHost     string
	redisPort     string
	redisPassword string

	dialector      string
	dbHost         string
	dbPort         string
	dbName         string
	dbUser         string
	dbPassword     string
	dbNameTest     string
	dbUserTest     string
	dbPasswordTest string
	dbGORM         *gorm.DB

	keyJWT    string
	keyCookie string

	uploadsPath string
	srcFolder   string
	layout      string

	runtimeFolder string

	smtpActive   bool
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string

	notifyEmail string
}

var (
	once     sync.Once
	instance *config
)

func LoadConfig() (err error) {
	once.Do(func() {
		instance = &config{Mutex: &sync.Mutex{}}

		var rootDir string
		rootDir, err = instance.root()
		if err != nil {
			err = fmt.Errorf("ошибка определения корневой директории: %w", err)
			return
		}

		err = godotenv.Load(filepath.Join(rootDir, ".env"))
		if err != nil {
			err = fmt.Errorf("ошибка при загрузке .env файла: %v. Скопируйте файл .env.example в .env", err)
			return
		}

		instance.env = getEnv("ENV", "local")
		instance.appHost = getEnv("APP_HOST", "local")
		instance.appPort = getEnv("APP_PORT", "3000")

		logLevel := getEnv("LOG_LEVEL", "6")
		instance.logLevel, err = strconv.Atoi(logLevel)
		if err != nil {
			instance.logLevel = 6
		}

		instance.keyJWT = getEnv("KEY_JWT", "")
		instance.keyCookie = getEnv("KEY_COOKIE", "")

		instance.store = getEnv("CASH_STORE", "redis")
		instance.redisHost = getEnv("REDIS_HOST", "127.0.0.1")
		instance.redisPort = getEnv("REDIS_PORT", "6380")
		instance.redisPassword = getEnv("REDIS_PASSWORD", "")

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
		instance.layout = getEnv("LAYOUT", "")

		instance.runtimeFolder = getEnv("RUNTIME_FOLDER", "runtime")

		instance.smtpHost = getEnv("SMTP_HOST", "")
		smtpPort := getEnv("SMTP_PORT", "2525")
		instance.smtpPort, err = strconv.Atoi(smtpPort)
		if err != nil {
			instance.smtpPort = 6
		}
		instance.smtpUsername = getEnv("SMTP_USERNAME", "")
		instance.smtpPassword = getEnv("SMTP_PASSWORD", "")

		smtpActiveTemp := getEnv("SMTP_ACTIVE", "0")
		smtpActive, err := strconv.Atoi(smtpActiveTemp)
		if err != nil {
			instance.smtpActive = false
		}
		instance.smtpActive = smtpActive == 1

		instance.notifyEmail = getEnv("NOTIFY_EMAIL", "")
		if instance.notifyEmail == "" {
			instance.smtpActive = false
		}
	})
	return
}

func Config() contracts.Config {
	if instance == nil {
		if err := LoadConfig(); err != nil {
			log.Fatalf("\x1b[1;91m%s\x1b[0m", "Ошибка загрузки конфигурации: "+err.Error())
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

func (c *config) IsLocal() bool {
	return c.env != "prod" && c.env != "dev"
}

func (c *config) DBUrl() string {
	if c.IsTest() {
		return c.DBUrlTest()
	}

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
			" appPort=" + c.dbPort +
			" sslmode=disable TimeZone=Europe/Moscow"
	}
	return dsn
}

func (c *config) SetGORM(db *gorm.DB) {
	c.Lock()
	defer c.Unlock()

	c.dbGORM = db
}

func (c *config) GetGORM() *gorm.DB {
	return c.dbGORM
}

func (c *config) RedisHost() string {
	if c.IsTest() {
		s := c.redisHost + ":" + c.redisPort
		return s
	}
	s := c.redisHost + ":" + c.redisPort
	return s
}

func (c *config) RedisPassword() string {
	return c.redisPassword
}

func (c *config) StoreIsRedis() bool {
	return c.store == "redis"
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

func (c *config) AppHost() string {
	host := strings.TrimSpace(c.appHost)
	port := strings.TrimPrefix(strings.TrimSpace(c.appPort), ":")

	// Не добавляем "дефолтные" порты
	if port == "" || port == "80" || port == "8080" {
		return host
	}

	// Если это полный URL — аккуратно добавим порт к u.Host
	if u, err := url.Parse(host); err == nil && u.Scheme != "" {
		// если порт уже указан — ничего не меняем
		if _, _, err := net.SplitHostPort(u.Host); err == nil {
			return host
		}
		u.Host = net.JoinHostPort(u.Hostname(), port)
		return u.String()
	}

	// Сырой хост (в т.ч. IPv6 в квадратных скобках). Если порт уже есть — вернуть как есть.
	if _, _, err := net.SplitHostPort(host); err == nil {
		return host
	}
	return net.JoinHostPort(host, port)
}

func (c *config) Port() string {
	if c.IsTest() {
		return ":" + strings.TrimPrefix(strings.TrimSpace(c.appPort), ":")
	}
	return ":" + strings.TrimPrefix(strings.TrimSpace(c.appPort), ":")
}

func (c *config) UploadPath() string {
	path := "/" + strings.Trim(c.uploadsPath, "/") + "/"
	if c.IsTest() {
		return path + "test/"
	}
	return path
}

func (c *config) RuntimeFolder(s string) string {
	if s != "" {
		return c.runtimeFolder + "/" + strings.Trim(s, " /")
	}
	return c.runtimeFolder
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

func (c *config) Layout() string {
	if c.layout == "" {
		return "default"
	}
	return c.layout
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

func (c *config) SMTPActive() bool {
	return c.smtpActive
}

func (c *config) SMTPPort() int {
	return c.smtpPort
}

func (c *config) SMTPHost() string {
	return c.smtpHost
}

func (c *config) SMTPUsername() string {
	return c.smtpUsername
}

func (c *config) SMTPPassword() string {
	return c.smtpPassword
}

func (c *config) NotifyEmail() string {
	return c.notifyEmail
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
