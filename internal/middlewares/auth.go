package middlewares

import (
	"errors"
	"strings"

	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/jwt"
	"github.com/gofiber/fiber/v2"
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
	userID, err := a.jwtUtils.VerifyAccessToken(ctx.Context(), token)
	if err != nil {
		return errorHandler.UnauthorizedError(err, "Invalid token")
	}

	user, err := a.userRepo.GetUserByID(userID)
	if err != nil {
		return errorHandler.UnauthorizedError(err, "User not found")
	}

	ctx.Locals("userID", userID)
	ctx.Locals("user", user)

	return ctx.Next()
}
