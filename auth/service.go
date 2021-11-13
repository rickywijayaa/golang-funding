package auth

import (
	"encoding/base64"
	"errors"
	"funding/env"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(encodedToken string) (*jwt.Token, error)
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

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, isValidAlogrithm := token.Method.(*jwt.SigningMethodHMAC)

		if !isValidAlogrithm {
			return nil, errors.New("Invalid Token")
		}

		hashJwt := base64.StdEncoding.EncodeToString([]byte(env.GetEnv("JWTSecret")))
		return []byte(hashJwt), nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}
