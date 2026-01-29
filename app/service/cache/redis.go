package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	client "github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

func NewRedisCache(cfg contract.Config) contract.Cache {
	logger.Info("[Cache] Redis is up, using Redis")

	c := &redisClient{config: cfg}
	c.client = client.NewClient(&client.Options{
		Addr:         cfg.RedisHost(),
		Password:     cfg.RedisPassword(),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	return c
}

func PingRedisCache(cfg contract.Config) error {
	rdb := client.NewClient(&client.Options{
		Addr:         cfg.RedisHost(),
		Password:     cfg.RedisPassword(),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return rdb.Ping(ctx).Err()
}

type redisClient struct {
	client *client.Client
	config contract.Config
}

func (r *redisClient) AddCache(key, value string) {
	if err := r.client.Set(context.Background(), key, value, 0).Err(); err != nil {
		logger.Errorf("[RedisClient][AddCache] Error: %v", err)
	}
}

func (r *redisClient) AddCacheTTL(key, value string, ttl time.Duration) {
	if err := r.client.Set(context.Background(), key, value, ttl).Err(); err != nil {
		logger.Errorf("[RedisClient][AddCacheTTL] Error: %v", err)
	}
}

func (r *redisClient) GetCache(key string) (string, bool) {
	value, err := r.client.Get(context.Background(), key).Result()
	if errors.Is(err, client.Nil) {
		return "", false
	}

	if err != nil {
		logger.Errorf("[RedisClient][GetCache] Error: %value", err)
		return "", false
	}

	return value, true
}

func (r *redisClient) DeleteCache(key string) {
	if err := r.client.Del(context.Background(), key).Err(); err != nil {
		logger.Errorf("[RedisClient][DeleteCache] Error: %v", err)
	}
}

func (r *redisClient) DeleteByPrefix(prefix string) {
	if prefix == "" {
		return
	}
	var cursor uint64
	pattern := prefix + "*"
	for {
		keys, next, err := r.client.Scan(context.Background(), cursor, pattern, 1000).Result()
		if err != nil {
			logger.Errorf("[RedisClient][DeleteByPrefix] Error: %v", err)
			return
		}
		for _, key := range keys {
			r.DeleteCache(key)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
}

func (r *redisClient) GetUserKey(id uint) string {
	return fmt.Sprintf(r.config.UserSessionKey("%d"), id)
}

func (r *redisClient) AddUserSession(id uint, sessionID string) {
	r.AddCache(r.GetUserKey(id), sessionID)
}

func (r *redisClient) ResetUserSession(userID uint) error {
	sessionID, err := r.client.Get(context.Background(), r.GetUserKey(userID)).Result()
	if errors.Is(err, client.Nil) {
		return nil
	}
	if err != nil {
		return err
	}

	if err := r.client.Del(context.Background(), sessionID).Err(); err != nil {
		return err
	}

	r.DeleteCache(r.GetUserKey(userID))

	return nil
}

// ResetUsersSession — очищает user->session и сами session ключи
func (r *redisClient) ResetUsersSession() {
	var cursor uint64
	for {
		keys, next, err := r.client.Scan(context.Background(), cursor, r.config.UserSessionKey("*"), 1000).Result()
		if err != nil {
			logger.Errorf("[RedisClient][ResetUsersSession] Error: %v", err)
			return
		}

		for _, key := range keys {
			r.DeleteCache(key)
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}

	cursor = 0
	for {
		keys, next, err := r.client.Scan(context.Background(), cursor, r.config.SessionKey("*"), 1000).Result()
		if err != nil {
			logger.Errorf("[RedisClient][ResetUsersSession] Error: %v", err)
			return
		}

		for _, key := range keys {
			r.DeleteCache(key)
		}

		cursor = next
		if cursor == 0 {
			break
		}
	}
}
