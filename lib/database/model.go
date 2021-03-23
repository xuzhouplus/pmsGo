package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"pmsGo/lib/config"
	"strconv"
)

type Model struct {
	db *gorm.DB
}

func Query(model interface{}) *gorm.DB {
	dsn := config.Config.Database.Username + ":" + config.Config.Database.Password + "@tcp(" + config.Config.Database.Host + ":" + strconv.Itoa(config.Config.Database.Port) + ")/" + config.Config.Database.Database + "?charset=" + config.Config.Database.Charset + "&parseTime=True&loc=Local"
	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.Config.Database.Prefix,
		},
	})
	if err != nil {
		fmt.Errorf("unable to connect to database:%err", err)
	}
	if config.Config.Site.Debug {
		return db.Debug().Model(model)
	}
	return db.Model(model)
}
