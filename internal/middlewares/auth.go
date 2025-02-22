package middlewares

import (
	"errors"
	"strings"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	jwtUtils *jwt.JWTUtils
	userRepo ports.UserRepository
}

func NewAuthMiddleware(jwtUtils *jwt.JWTUtils, userRepo ports.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtUtils: jwtUtils,
		userRepo: userRepo,
	}
}

func (a *AuthMiddleware) Auth(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		return errorHandler.UnauthorizedError(errors.New("request without authorization header"), "Authorization header is required")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return errorHandler.UnauthorizedError(errors.New("invalid authorization header"), "Authorization header is invalid")
	}

	token := authHeader[7:]
	claims, err := a.jwtUtils.DecodeJWT(token)
	if err != nil {
		return errorHandler.UnauthorizedError(err, "Invalid token")
	}

	if claims.GetExp() < time.Now().Unix() {
		return errorHandler.UnauthorizedError(errors.New("token expired"), "Token is expired")
	}

	uuid, err := uuid.Parse(claims.GetUserID())
	if err != nil {
		return errorHandler.UnauthorizedError(err, "Invalid user ID")
	}

	user, err := a.userRepo.GetUserByID(uuid)
	if err != nil {
		return errorHandler.UnauthorizedError(err, "User not found")
	}

	ctx.Locals("userID", claims.GetUserID())
	ctx.Locals("user", user)

	return ctx.Next()
}
