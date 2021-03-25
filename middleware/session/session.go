package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"pmsGo/lib/config"
	"strconv"
)

func Register() gin.HandlerFunc {
	store, _ := redis.NewStoreWithDB(config.Config.Session.Idle, "tcp", config.Config.Redis.Host+":"+strconv.Itoa(config.Config.Redis.Port), config.Config.Redis.Auth, strconv.Itoa(config.Config.Redis.Database), []byte(config.Config.Session.Secret))
	return sessions.Sessions(config.Config.Session.Name, store)
}
