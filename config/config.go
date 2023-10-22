package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server_Port     string
	DB_Username     string
	DB_Password     string
	DB_Port         string
	DB_Host         string
	DB_Name         string
	Secret          string
	CDN_Cloud_Name  string
	CDN_API_Key     string
	CDN_API_Secret  string
	CDN_Folder_Name string
	MT_Server_Key   string
}

func InitConfig() *Config {
	var result = new(Config)
	result = loadConfig()

	if result == nil {
		logrus.Fatal("Config : Cannot start program, failed to load configuration")
		return nil
	}

	return result
}

func loadConfig() *Config {
	var result = new(Config)

	err := godotenv.Load(".env")
	if err != nil {
		logrus.Error("Config: Cannot load config file,", err.Error())
		return nil
	}

	if value, found := os.LookupEnv("SERVER"); found {
		result.Server_Port = value
	}
	if value, found := os.LookupEnv("DBUSER"); found {
		result.DB_Username = value
	}
	if value, found := os.LookupEnv("DBPASS"); found {
		result.DB_Password = value
	}
	if value, found := os.LookupEnv("DBPORT"); found {
		result.DB_Port = value
	}
	if value, found := os.LookupEnv("DBHOST"); found {
		result.DB_Host = value
	}
	if value, found := os.LookupEnv("DBNAME"); found {
		result.DB_Name = value
	}
	if value, found := os.LookupEnv("SECRET"); found {
		result.Secret = value
	}
	if value, found := os.LookupEnv("CLOUDINARY_CLOUD_NAME"); found {
		result.CDN_Cloud_Name = value
	}
	if value, found := os.LookupEnv("CLOUDINARY_API_KEY"); found {
		result.CDN_API_Key = value
	}
	if value, found := os.LookupEnv("CLOUDINARY_API_SECRET"); found {
		result.CDN_API_Secret = value
	}
	if value, found := os.LookupEnv("CLOUDINARY_UPLOAD_FOLDER_NAME"); found {
		result.CDN_Folder_Name = value
	}
	if value, found := os.LookupEnv("MIDTRANS_SERVER_KEY"); found {
		result.MT_Server_Key = value
	}

	return result
}
