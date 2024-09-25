package db

import (
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

func InitRedis(cfg *config.Config) redis.Store {
	store, err := redis.NewStore(10, "tcp", cfg.RedisHost+":"+cfg.RedisPort, "", []byte(cfg.KeyCookie))
	if err != nil {
		panic(err)
	}
	store.Options(sessions.Options{
		MaxAge: 86400 * 7,
		Path:   "/",
	})
	return store
}

func Cache() contracts.Cache {
	return models.Redis()
}
