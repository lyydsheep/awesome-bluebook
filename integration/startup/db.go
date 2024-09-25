package startup

import (
	config2 "awesome-bluebook/config"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	var config config2.DBConfig
	err := viper.UnmarshalKey("db", &config)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(mysql.Open(config.DSN))
	if err != nil {
		panic(err)
	}
	return db
}
