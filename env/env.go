package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", GetDBUser(),
		GetDBPassword(), GetDBUrl(), GetDBName())
}

func GetEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error get env")
	}

	return os.Getenv(key)
}

func GetDBUser() string {
	return GetEnv("DB_USER")
}

func GetDBPassword() string {
	return GetEnv("DB_PASSWORD")
}

func GetDBUrl() string {
	return GetEnv("DB_URL")
}

func GetDBName() string {
	return GetEnv("DB_NAME")
}
