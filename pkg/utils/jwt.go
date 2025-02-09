package utils

import (
	"errors"
	"fmt"
	"time"

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
	return jwtToken.SignedString([]byte(j.Config.JWTSecretKey))
}

func (j *JWTUtils) DecodeJWT(inputToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(inputToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Config.JWTSecretKey), nil
	})

	if err != nil || !token.Valid {
		return new(JWTClaims), errors.New("parse token failed")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return new(JWTClaims), errors.New("invalid token")
	}

	return claims, nil
}

// Deprecated: Use JWTUtils.DecodeJWT instead.
func GenerateJWT(userID uuid.UUID, config *JWTConfig) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID.String()
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(config.Expiration)).Unix()

	return token.SignedString([]byte(config.JWTSecretKey))
}

// Deprecated: Use JWTUtils.DecodeJWT instead.
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
