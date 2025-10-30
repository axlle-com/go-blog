package cache

import (
	"github.com/axlle-com/blog/app/config"
	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
)

func NewCache() contract.Cache {
	cfg := config.Config()

	if cfg.IsTest() || !cfg.StoreIsRedis() {
		return NewInMemoryCache(cfg)
	}

	if err := PingRedisCache(cfg); err != nil {
		logger.Errorf("[Cache] Redis ping failed, falling back to in-memory: %v", err)
		return NewInMemoryCache(cfg)
	}

	return NewRedisCache(cfg)
}
