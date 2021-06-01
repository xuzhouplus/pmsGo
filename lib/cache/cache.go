package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"pmsGo/lib/config"
	"time"
)

type redisCache struct {
	Prefix string
	Expire int
	Redis  *redis.Client
	Cache  *cache.Cache
}

var Cache *redisCache

func init() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Auth,
		DB:       config.Config.Redis.Database,
	})
	cacheClient := cache.New(&cache.Options{
		Redis:      redisClient,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	Cache = &redisCache{
		Prefix: config.Config.Redis.Prefix,
		Expire: config.Config.Redis.Expire,
		Redis:  redisClient,
		Cache:  cacheClient,
	}
}

func (rc redisCache) Key(key string) string {
	return rc.Prefix + key
}

func (rc redisCache) Get(key string) (interface{}, error) {
	key = rc.Key(key)
	var val interface{}
	error := rc.Cache.Get(context.Background(), key, &val)
	if error == redis.Nil {
		return nil, nil
	} else if error != nil {
		return nil, error
	}
	return val, nil
}

func (rc redisCache) Set(key string, value interface{}) error {
	key = rc.Key(key)
	err := rc.Cache.Set(&cache.Item{
		Ctx:   context.Background(),
		Key:   key,
		Value: value,
	})
	if err != nil {
		return err
	}
	return nil
}

func (rc redisCache) SetNX(key string, value interface{}, ttl int) error {
	key = rc.Key(key)
	err := rc.Cache.Set(&cache.Item{
		Ctx:   context.Background(),
		Key:   key,
		Value: value,
		SetNX: true,
		TTL:   time.Duration(ttl),
	})
	if err != nil {
		return err
	}
	return nil
}

func (rc redisCache) SetEX(key string, value interface{}, ttl int) error {
	key = rc.Key(key)
	if ttl == 0 {
		ttl = rc.Expire
	}
	err := rc.Cache.Set(&cache.Item{
		Ctx:   context.Background(),
		Key:   key,
		Value: value,
		SetNX: true,
		TTL:   time.Duration(ttl),
	})
	if err != nil {
		return err
	}
	return nil
}
