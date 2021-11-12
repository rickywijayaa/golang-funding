package auth

import (
	"encoding/base64"
	"funding/env"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userID int) (string, error)
}

type jwtService struct {
}

func NewJwtService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {
	payLoad := jwt.MapClaims{}
	payLoad["user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payLoad)
	hashJwt := base64.StdEncoding.EncodeToString([]byte(env.GetEnv("JWTSecret")))
	signedToken, err := token.SignedString([]byte(hashJwt))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}
