package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"pmsGo/lib/config"
	"pmsGo/lib/log"
	"time"
)

var DB *gorm.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", config.Config.Database.Username, config.Config.Database.Password, config.Config.Database.Host, config.Config.Database.Database, config.Config.Database.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.Config.Database.Prefix,
		},
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Panicf("unable to connect to database:%err \n", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Panicf("init db connect failed:%err", err)
	}
	if config.Config.Database.MaxIdleConnect > 0 {
		sqlDB.SetMaxIdleConns(config.Config.Database.MaxIdleConnect)
	}
	if config.Config.Database.MaxOpenConnect > 0 {
		sqlDB.SetMaxOpenConns(config.Config.Database.MaxOpenConnect)
	}
	if config.Config.Database.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(config.Config.Database.ConnMaxIdleTime))
	}
	if config.Config.Database.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(config.Config.Database.ConnMaxLifetime))
	}
	if config.Config.Site.Debug {
		DB = db.Debug()
	} else {
		DB = db
	}
}
