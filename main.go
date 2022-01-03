package main

import (
	auth "funding/auth"
	"funding/campaign"
	"funding/env"
	"funding/handler"
	"funding/middleware"
	"funding/payment"
	"funding/transaction"
	"funding/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open(env.Dsn()), &gorm.Config{})
	// db, err := gorm.Open(mysql.Open(env.Dsn()), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewJwtService()
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService, paymentService)

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email-checker", userHandler.IsEmailExist)
	api.POST("/avatars", middleware.AuthMiddleware(userService, authService), userHandler.UploadAvatar)
	api.GET("/users/fetch", middleware.AuthMiddleware(userService, authService), userHandler.FetchUser)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", middleware.AuthMiddleware(userService, authService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", middleware.AuthMiddleware(userService, authService), campaignHandler.UpdateCampaign)
	api.POST("/campaigns-images", middleware.AuthMiddleware(userService, authService), campaignHandler.UploadImage)

	api.GET("/campaign/:id/transactions", middleware.AuthMiddleware(userService, authService), transactionHandler.GetCampaignsTransaction)
	api.GET("/transactions", middleware.AuthMiddleware(userService, authService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", middleware.AuthMiddleware(userService, authService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	router.Run()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
