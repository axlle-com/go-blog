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
		Addr: cfg.RedisHost(),
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

	if err := rdb.Ping(ctx).Err(); err != nil {
		return err
	}

	return nil
}

type redisClient struct {
	client *client.Client
	config contract.Config
}

func (r *redisClient) AddCache(key, value string) {
	err := r.client.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		logger.Errorf("[RedisClient][AddCache] Error: %v", err)
	}
}

func (r *redisClient) DeleteCache(key string) {
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		logger.Errorf("[RedisClient][DeleteCache] Error: %v", err)
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

	err = r.client.Del(context.Background(), sessionID).Err()
	if err != nil {
		return err
	}

	r.DeleteCache(r.GetUserKey(userID))
	return nil
}

func (r *redisClient) ResetUsersSession() {
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(context.Background(), cursor, r.config.UserSessionKey("*"), 1000).Result()
		if err != nil {
			logger.Errorf("[RedisClient][ResetUsersSession] Error: %v", err)
		}

		if len(keys) == 0 {
			break
		}

		for _, key := range keys {
			r.DeleteCache(key)
		}
	}
	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(context.Background(), cursor, r.config.SessionKey("*"), 1000).Result()
		if err != nil {
			logger.Errorf("[RedisClient][ResetUsersSession] Error: %v", err)
		}

		if len(keys) == 0 {
			break
		}

		for _, key := range keys {
			r.DeleteCache(key)
		}
	}
}
