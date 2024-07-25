package db

import (
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

func InitRedis(cfg *config.Config) redis.Store {
	store, err := redis.NewStore(10, "tcp", "127.0.0.1:6379", "", []byte(cfg.KeyCookie))
	if err != nil {
		panic(err)
	}
	store.Options(sessions.Options{
		MaxAge: 86400 * 7,
		Path:   "/",
	})
	return store
}
