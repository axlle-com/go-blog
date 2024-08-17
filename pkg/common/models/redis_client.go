package models

import (
	"fmt"
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/axlle-com/blog/pkg/common/logger"
	client "github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

const UserSessionKey = "user_session_"

type Cache interface {
	AddCache(key, value string)
	DeleteCache(key string)
	GetUserKey(id uint) string
	AddUserSession(id uint, sessionID string)
	ResetUserSession(userID uint) error
	ResetUsersSession()
}

func NewRedisClient() Cache {
	cfg := config.GetConfig()
	c := &redisClient{}
	c.client = client.NewClient(&client.Options{
		Addr: cfg.RedisHost + ":" + cfg.RedisPort,
	})

	return c
}

type redisClient struct {
	client *client.Client
}

func (r *redisClient) AddCache(key, value string) {
	err := r.client.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		logger.Fatal(err)
	}
}

func (r *redisClient) DeleteCache(key string) {
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		logger.Fatal(err)
	}
}

func (r *redisClient) GetUserKey(id uint) string {
	return fmt.Sprintf(UserSessionKey+"%d", id)
}

func (r *redisClient) AddUserSession(id uint, sessionID string) {
	r.AddCache(r.GetUserKey(id), sessionID)
}

func (r *redisClient) ResetUserSession(userID uint) error {
	sessionID, err := r.client.Get(context.Background(), r.GetUserKey(userID)).Result()
	if err == client.Nil {
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
		keys, cursor, err = r.client.Scan(context.Background(), cursor, UserSessionKey+"*", 1000).Result()
		if err != nil {
			logger.Fatal(err)
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
		keys, cursor, err = r.client.Scan(context.Background(), cursor, "session_*", 1000).Result()
		if err != nil {
			logger.Fatal(err)
		}

		if len(keys) == 0 {
			break
		}

		for _, key := range keys {
			r.DeleteCache(key)
		}
	}
}
