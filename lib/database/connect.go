package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"pmsGo/lib/config"
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", config.Config.Database.Username, config.Config.Database.Password, config.Config.Database.Host, config.Config.Database.Database, config.Config.Database.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.Config.Database.Prefix,
		},
	})
	if err != nil {
		fmt.Errorf("unable to connect to database:%err", err)
	}
	if config.Config.Site.Debug {
		DB = db.Debug()
	} else {
		DB = db
	}
}
