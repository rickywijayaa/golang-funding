package main

import (
	auth "funding/auth"
	"funding/campaign"
	env "funding/env"
	"funding/handler"
	"funding/middleware"
	"funding/transaction"
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
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewJwtService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := gin.Default()
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email-checker", userHandler.IsEmailExist)
	api.POST("/avatars", middleware.AuthMiddleware(userService, authService), userHandler.UploadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", middleware.AuthMiddleware(userService, authService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", middleware.AuthMiddleware(userService, authService), campaignHandler.UpdateCampaign)
	api.POST("/campaigns-images", middleware.AuthMiddleware(userService, authService), campaignHandler.UploadImage)

	api.GET("/campaign/:id/transactions", middleware.AuthMiddleware(userService, authService), transactionHandler.GetCampaignsTransaction)
	api.GET("/transactions", middleware.AuthMiddleware(userService, authService), transactionHandler.GetUserTransactions)

	router.Run()
}
