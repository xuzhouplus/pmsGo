package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"pmsGo/lib/config"
	"strconv"
)

// Register 注册session中间件
func Register() gin.HandlerFunc {
	//创建redis连接
	store, err := redis.NewStoreWithDB(config.Config.Session.Idle, "tcp", config.Config.Redis.Host+":"+strconv.Itoa(config.Config.Redis.Port), config.Config.Redis.Auth, strconv.Itoa(config.Config.Redis.Database), []byte(config.Config.Session.Secret))
	if err != nil {
		panic(err.Error())
	}
	err = redis.SetKeyPrefix(store, config.Config.Session.Prefix)
	if err != nil {
		panic(err.Error())
	}
	//配置session存储为redis
	return sessions.Sessions(config.Config.Session.Name, store)
}
