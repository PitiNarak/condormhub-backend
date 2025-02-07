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
