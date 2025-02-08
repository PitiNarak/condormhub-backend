package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTConfig struct {
	JWTSecretKey string `env:"SECRET,required"`
	Expiration   int    `env:"EXPIRATION_HOURS,required"`
}

func GenerateJWT(userID uuid.UUID, config *JWTConfig) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID.String()
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(config.Expiration)).Unix()

	return token.SignedString([]byte(config.JWTSecretKey))
}

func DecodeJWT(inputToken string, config *JWTConfig) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(inputToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecretKey), nil
	})

	if err != nil || !token.Valid {
		return new(jwt.MapClaims), err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return new(jwt.MapClaims), err
	}
	return &claims, nil
}
