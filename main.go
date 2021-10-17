package main

import (
	env "funding/env"
	"funding/user"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open(env.Dsn()), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	user := user.User{
		Name: "Test",
	}

	userRepository.Save(user)
}
