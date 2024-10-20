package db

import (
	"github.com/axlle-com/blog/pkg/common/models"
	"github.com/axlle-com/blog/pkg/common/models/contracts"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

func RedisStore(address string, password string, keyPairs []byte) redis.Store {
	store, err := redis.NewStore(10, "tcp", address, password, keyPairs)
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
