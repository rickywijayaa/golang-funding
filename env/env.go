package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Dsn() string {
	// return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", GetDBUser(),
	// 	GetDBPassword(), GetDBUrl(), GetDBName())
	return fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		GetDBHost(), GetDBUser(), GetDBPassword(), GetDBPort(), GetDBName())
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

func GetDBHost() string {
	return GetEnv("DB_HOST")
}

func GetDBName() string {
	return GetEnv("DB_NAME")
}

func GetDBPort() string {
	return GetEnv("DB_PORT")
}

func GetClientKey() string {
	return GetEnv("ClientKey")
}

func GetServerKey() string {
	return GetEnv("ServerKey")
}
