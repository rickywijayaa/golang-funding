package main

import (
	auth "funding/auth"
	"funding/campaign"
	env "funding/env"
	"funding/handler"
	"funding/middleware"
	"funding/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open(env.Dsn()), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewJwtService()

	userHandler := handler.NewUserHandler(userService, authService)

	// campaign, _ := campaignRepository.FindByUserID(21)
	// for _, value := range campaign {
	// 	fmt.Println(value)
	// 	if len(value.CampaignImages) > 0 {
	// 		for _, img := range value.CampaignImages {
	// 			fmt.Println(img)
	// 		}
	// 	}
	// }

	router := gin.Default()
	api := router.Group("/api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email-checker", userHandler.IsEmailExist)
	api.POST("/avatars", middleware.AuthMiddleware(userService, authService), userHandler.UploadAvatar)

	router.Run()
}
