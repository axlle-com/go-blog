package config

import (
	"github.com/spf13/viper"
	"sync"
)

const SessionsName = "web_session"

type Config struct {
	Port      string `mapstructure:"PORT"`
	DBUrl     string `mapstructure:"DB_URL"`
	KeyJWT    string `mapstructure:"KEY_JWT"`
	KeyCookie string `mapstructure:"KEY_COOKIE"`
}

var (
	once     sync.Once
	instance *Config
)

func LoadConfig() (err error) {
	once.Do(func() {
		viper.AddConfigPath("./pkg/common/envs")
		viper.SetConfigName("dev")
		viper.SetConfigType("env")
		viper.AutomaticEnv()

		if err = viper.ReadInConfig(); err != nil {
			return
		}

		instance = &Config{}
		err = viper.Unmarshal(instance)
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
