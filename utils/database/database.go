package database

import (
	"fmt"
	"restoran/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(config config.Config) *gorm.DB {
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DB_Username,
		config.DB_Password,
		config.DB_Host,
		config.DB_Port,
		config.DB_Name)

	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Error("Model : cannot connect to database, ", err.Error())
		return nil
	}

	return DB
}
