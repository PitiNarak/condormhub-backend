package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/PitiNarak/condormhub-backend/pkg/redis"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/yokeTH/go-pkg/apperror"
)

type JWTConfig struct {
	JWTSecretKey           string `env:"SECRET,required"`
	AccessTokenExpiration  int    `env:"EXPIRATION_HOURS,required"`
	RefreshTokenExpiration int    `env:"REFRESH_EXPIRATION_HOURS,required"`
}

type JWTUtils struct {
	Config *JWTConfig
	Redis  *redis.Redis
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

func NewJWTUtils(config *JWTConfig, redis *redis.Redis) *JWTUtils {
	return &JWTUtils{Config: config, Redis: redis}
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

func (j *JWTUtils) GenerateJWT(userID uuid.UUID, exp int) (string, error) {
	claims := &JWTClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(exp))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString([]byte(j.Config.JWTSecretKey))
	if err != nil {
		return "", apperror.InternalServerError(err, "cannot generate token")
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
		return new(JWTClaims), apperror.UnauthorizedError(err, "parse token failed")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return new(JWTClaims), apperror.UnauthorizedError(err, "invalid token")
	}

	return claims, nil
}

func (j *JWTUtils) GenerateKeyPair(ctx context.Context, userID uuid.UUID) (string, string, error) {
	accessToken, err := j.GenerateJWT(userID, j.Config.AccessTokenExpiration)
	if err != nil {
		return "", "", err
	}
	err = j.Redis.SetAccessToken(ctx, userID, accessToken, time.Hour*time.Duration(j.Config.AccessTokenExpiration))
	if err != nil {
		return "", "", apperror.InternalServerError(err, "cannot set access token")
	}

	refreshToken, err := j.GenerateJWT(userID, j.Config.RefreshTokenExpiration)
	if err != nil {
		return "", "", err
	}
	err = j.Redis.SetRefreshToken(ctx, userID, refreshToken, time.Hour*time.Duration(j.Config.RefreshTokenExpiration))
	if err != nil {
		return "", "", apperror.InternalServerError(err, "cannot set refresh token")
	}

	return accessToken, refreshToken, nil
}

func (j *JWTUtils) VerifyAccessToken(ctx context.Context, accessToken string) (uuid.UUID, error) {
	claims, err := j.DecodeJWT(accessToken)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, apperror.InternalServerError(err, "cannot parse user id")
	}

	token, err := j.Redis.GetAccessToken(ctx, userID)
	if err != nil {
		return uuid.Nil, apperror.InternalServerError(err, "cannot get access token")
	}

	if token != accessToken {
		return uuid.Nil, apperror.UnauthorizedError(nil, "invalid access token")
	}

	return userID, nil
}

func (j *JWTUtils) VerifyRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error) {
	claims, err := j.DecodeJWT(refreshToken)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, apperror.InternalServerError(err, "cannot parse user id")
	}

	token, err := j.Redis.GetRefreshToken(ctx, userID)
	if err != nil {
		return uuid.Nil, apperror.InternalServerError(err, "cannot get refresh token")
	}

	if token != refreshToken {
		return uuid.Nil, apperror.UnauthorizedError(nil, "invalid refresh token")
	}

	return userID, nil
}

func (j *JWTUtils) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	userID, err := j.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := j.GenerateKeyPair(ctx, userID)
	if err != nil {
		return "", "", err
	}
	err = j.Redis.SetAccessToken(ctx, userID, accessToken, time.Hour*time.Duration(j.Config.AccessTokenExpiration))
	if err != nil {
		return "", "", apperror.InternalServerError(err, "cannot set access token")
	}

	return accessToken, refreshToken, nil
}

func (j *JWTUtils) GenerateResetPasswordToken(ctx context.Context, userID uuid.UUID) (string, error) {
	resetToken, err := j.GenerateJWT(userID, 24)
	if err != nil {
		return "", err
	}
	err = j.Redis.SetResetToken(ctx, userID, resetToken, time.Hour*24)
	if err != nil {
		return "", apperror.InternalServerError(err, "cannot set reset token")
	}
	return resetToken, nil
}

func (j *JWTUtils) VerifyResetPasswordToken(ctx context.Context, resetToken string) (uuid.UUID, error) {
	claims, err := j.DecodeJWT(resetToken)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, apperror.InternalServerError(err, "cannot parse user id")
	}

	token, err := j.Redis.GetResetToken(ctx, userID)
	if err != nil {
		return uuid.Nil, apperror.UnauthorizedError(err, "token is expired or token is used")
	}

	if token != resetToken {
		return uuid.Nil, apperror.UnauthorizedError(nil, "invalid reset token")
	}

	return userID, nil
}

func (j *JWTUtils) DeleteResetPasswordToken(ctx context.Context, userID uuid.UUID) error {
	err := j.Redis.DeleteResetToken(ctx, userID)
	if err != nil {
		return apperror.InternalServerError(err, "cannot delete reset token")
	}
	return nil
}

func (j *JWTUtils) GenerateVerificationToken(ctx context.Context, userID uuid.UUID) (string, error) {
	verificationToken, err := j.GenerateJWT(userID, 24)
	if err != nil {
		return "", err
	}
	err = j.Redis.SetVerificationToken(ctx, userID, verificationToken, time.Hour*24)
	if err != nil {
		return "", apperror.InternalServerError(err, "cannot set verification token")
	}
	return verificationToken, nil
}

func (j *JWTUtils) VerifyVerificationToken(ctx context.Context, verificationToken string) (uuid.UUID, error) {
	claims, err := j.DecodeJWT(verificationToken)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, apperror.InternalServerError(err, "cannot parse user id")
	}

	token, err := j.Redis.GetVerificationToken(ctx, userID)
	if err != nil {
		return uuid.Nil, apperror.UnauthorizedError(err, "token is expired or token is used")
	}

	if token != verificationToken {
		return uuid.Nil, apperror.UnauthorizedError(nil, "invalid verification token")
	}

	return userID, nil
}

func (j *JWTUtils) DeleteVerificationToken(ctx context.Context, userID uuid.UUID) error {
	err := j.Redis.DeleteVerificationToken(ctx, userID)
	if err != nil {
		return apperror.InternalServerError(err, "cannot delete verification token")
	}
	return nil
}
