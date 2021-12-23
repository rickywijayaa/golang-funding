package middleware

import (
	"funding/auth"
	"funding/helper"
	"funding/user"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userService user.Service, jwtService auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.APIFailedResponse(
				"Unauthorized",
				http.StatusUnauthorized,
				nil,
			))
			return
		}

		tokenString := ""
		headerToken := strings.Split(authHeader, " ")
		if len(headerToken) == 2 {
			tokenString = headerToken[1]
		}

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.APIFailedResponse(
				"Unauthorized",
				http.StatusUnauthorized,
				nil,
			))
			return
		}

		payload, isValid := token.Claims.(jwt.MapClaims)
		if !isValid || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.APIFailedResponse(
				"Unauthorized",
				http.StatusUnauthorized,
				nil,
			))
			return
		}

		userID := int(payload["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.APIFailedResponse(
				"Unauthorized",
				http.StatusUnauthorized,
				nil,
			))
			return
		}

		c.Set("current_user", user)
	}
}
