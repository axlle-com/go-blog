package models

import (
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
)

func Store(cfg contract.Config) redis.Store {
	var store sessions.Store
	var err error
	if cfg.IsTest() || !cfg.StoreIsRedis() {
		store = memstore.NewStore(cfg.KeyCookie())
		logger.Info("[Store] Using memstore")
		return store
	}

	store, err = redis.NewStore(10, "tcp", cfg.RedisHost(), "", cfg.RedisPassword(), cfg.KeyCookie())
	if err != nil {
		logger.Errorf("[Store] Error: %v, Started memstore", err)
		store = memstore.NewStore(cfg.KeyCookie())
		return store
	}

	logger.Info("[Store] Using redis")
	store.Options(sessions.Options{
		MaxAge: 86400 * 7,
		Path:   "/",
	})
	return store
}
