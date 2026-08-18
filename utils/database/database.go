package database

import (
	"fmt"
	"restoran/config"
	"time"

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

	var DB *gorm.DB
	var err error

	for i := 0; i < 15; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			logrus.Info("Database connected successfully")
			return DB
		}
		logrus.Warnf("Cannot connect to database (attempt %d/15): %v. Retrying in 2s...", i+1, err)
		time.Sleep(2 * time.Second)
	}

	logrus.Fatalf("Model : cannot connect to database after retries, %v", err)
	return nil
}
