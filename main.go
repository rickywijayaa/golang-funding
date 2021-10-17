package main

import (
	"fmt"
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

	var users []user.User
	db.Find(&users)

	for _, user := range users {
		fmt.Println(user.Name)
	}

}
