package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

const SessionsName = "web_session"

type Config struct {
	RedisHost  string
	RedisPort  string
	Port       string
	DBUrl      string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	KeyJWT     string
	KeyCookie  string
}

var (
	once     sync.Once
	instance *Config
)

func LoadConfig() (err error) {
	once.Do(func() {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Ошибка при загрузке .env файла: %v", err)
		}

		instance = &Config{}
		instance.Port = os.Getenv("PORT")
		instance.KeyJWT = os.Getenv("KEY_JWT")
		instance.KeyCookie = os.Getenv("KEY_COOKIE")
		instance.RedisHost = os.Getenv("REDIS_HOST")
		instance.RedisPort = os.Getenv("REDIS_PORT")
		instance.DBHost = os.Getenv("POSTGRES_HOST")
		instance.DBPort = os.Getenv("POSTGRES_PORT")
		instance.DBUser = os.Getenv("POSTGRES_USER")
		instance.DBPassword = os.Getenv("POSTGRES_PASSWORD")
		instance.DBUrl = "postgres://" + instance.DBUser + ":" + instance.DBPassword + "@" + instance.DBHost + ":" + instance.DBPort + "/postgres"
	})
	return
}

func GetConfig() *Config {
	if instance == nil {
		if err := LoadConfig(); err != nil {
			panic("Ошибка загрузки конфигурации: " + err.Error())
		}
	}
	return instance
}
