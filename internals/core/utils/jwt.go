package utils

import (
	"time"

	"github.com/PitiNarak/condormhub-backend/internals/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateJWT(userID uuid.UUID, config *config.JWTConfig) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID.String()
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(config.Expiration)).Unix()

	return token.SignedString([]byte(config.JWTSecretKey))
}
