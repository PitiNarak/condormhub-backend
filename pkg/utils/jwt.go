package utils

import (
	"fmt"
	"time"

	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTConfig struct {
	JWTSecretKey string `env:"SECRET,required"`
	Expiration   int    `env:"EXPIRATION_HOURS,required"`
}

type JWTUtils struct {
	Config *JWTConfig
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"userID"`
}

type JWTClaimsInterface interface {
	GetUserID() string
	GetExp() int64
	GetIat() int64
}

func NewJWTUtils(config *JWTConfig) *JWTUtils {
	return &JWTUtils{Config: config}
}

func (j *JWTClaims) GetUserID() string {
	return j.UserID
}

func (j *JWTClaims) GetExp() int64 {
	return j.ExpiresAt.Unix()
}

func (j *JWTClaims) GetIat() int64 {
	return j.IssuedAt.Unix()
}

func (j *JWTUtils) GenerateJWT(userID uuid.UUID) (string, error) {
	claims := &JWTClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.Config.Expiration))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString([]byte(j.Config.JWTSecretKey))
	if err != nil {
		return "", errorHandler.InternalServerError(err, "cannot generate token")
	}
	return tokenString, nil
}

func (j *JWTUtils) DecodeJWT(inputToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(inputToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Config.JWTSecretKey), nil
	})

	if err != nil || !token.Valid {
		return new(JWTClaims), errorHandler.UnauthorizedError(err, "parse token failed")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return new(JWTClaims), errorHandler.UnauthorizedError(err, "invalid token")
	}

	return claims, nil
}
