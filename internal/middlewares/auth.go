package middlewares

import (
	"errors"
	"strings"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	jwtUtils *utils.JWTUtils
	userRepo ports.UserRepository
}

func NewAuthMiddleware(jwtUtils *utils.JWTUtils, userRepo ports.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwtUtils: jwtUtils,
		userRepo: userRepo,
	}
}

func (a *AuthMiddleware) Auth(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		return error_handler.UnauthorizedError(errors.New("request without authorization header"), "Authorization header is required")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return error_handler.UnauthorizedError(errors.New("invalid authorization header"), "Authorization header is invalid")
	}

	token := authHeader[7:]
	claims, err := a.jwtUtils.DecodeJWT(token)
	if err != nil {
		return error_handler.UnauthorizedError(err, "Invalid token")
	}

	if claims.GetExp() < time.Now().Unix() {
		return error_handler.UnauthorizedError(errors.New("token expired"), "Token is expired")
	}

	uuid, err := uuid.Parse(claims.GetUserID())
	if err != nil {
		return error_handler.UnauthorizedError(err, "Invalid user ID")
	}

	user, err := a.userRepo.GetUser(uuid)
	if err != nil {
		return error_handler.UnauthorizedError(err, "User not found")
	}

	ctx.Locals("userID", claims.GetUserID())
	ctx.Locals("user", user)

	return ctx.Next()
}
