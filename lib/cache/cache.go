package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"pmsGo/lib/config"
	"time"
)

var Prefix string
var Expire int
var Redis *redis.Client
var Cache *cache.Cache

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Auth,
		DB:       config.Config.Redis.Database,
	})
	Cache = cache.New(&cache.Options{
		Redis:      Redis,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
		Marshal: func(i interface{}) ([]byte, error) {
			return Marshal(i)
		},
		Unmarshal: func(bytes []byte, i interface{}) error {
			return Unmarshal(bytes, i)
		},
	})
	Prefix = config.Config.Cache.Prefix
	Expire = config.Config.Cache.Expire
}

func Key(key string) string {
	return Prefix + key
}

func Ttl(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

func Marshal(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func Unmarshal(bytes []byte, i interface{}) error {
	return json.Unmarshal(bytes, i)
}

func Get(key string) (interface{}, error) {
	key = Key(key)
	var val interface{}
	err := Cache.Get(context.Background(), key, &val)
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

func Set(key string, value interface{}, ttl int) error {
	key = Key(key)
	item := &cache.Item{
		Ctx:   context.Background(),
		Key:   key,
		Value: value,
	}
	if ttl == 0 {
		ttl = Expire
	}
	if ttl > 0 {
		item.TTL = Ttl(ttl)
	}
	err := Cache.Set(item)
	if err != nil {
		return err
	}
	return nil
}

func SetNX(key string, value interface{}, ttl int) error {
	key = Key(key)
	err := Cache.Set(&cache.Item{
		Ctx:   context.Background(),
		Key:   key,
		Value: value,
		SetNX: true,
		TTL:   Ttl(ttl),
	})
	if err != nil {
		return err
	}
	return nil
}

func SetEX(key string, value interface{}, ttl int) error {
	key = Key(key)
	if ttl == 0 {
		ttl = Expire
	}
	err := Cache.Set(&cache.Item{
		Ctx:   context.Background(),
		Key:   key,
		Value: value,
		SetNX: true,
		TTL:   Ttl(ttl),
	})
	if err != nil {
		return err
	}
	return nil
}

func Push(key string, value interface{}) error {
	key = Key(key)
	marshal, err := Marshal(value)
	if err != nil {
		return err
	}
	result := Redis.LPush(context.Background(), key, marshal)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func Pop(key string) (value interface{}, err error) {
	key = Key(key)
	result := Redis.RPop(context.TODO(), key)
	if result.Err() != nil {
		return nil, result.Err()
	}
	data, err := result.Result()
	if err != nil {
		return nil, err
	}
	err = Unmarshal([]byte(data), &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
