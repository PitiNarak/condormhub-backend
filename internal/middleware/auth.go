package middleware

import (
	"errors"
	"strings"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
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
		return apperror.UnauthorizedError(errors.New("request without authorization header"), "Authorization header is required")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return apperror.UnauthorizedError(errors.New("invalid authorization header"), "Authorization header is invalid")
	}

	token := authHeader[7:]
	userID, err := a.jwtUtils.VerifyAccessToken(ctx.Context(), token)
	if err != nil {
		return apperror.UnauthorizedError(err, "Invalid token")
	}

	user, err := a.userRepo.GetUserByID(userID)
	if err != nil {
		return apperror.UnauthorizedError(err, "User not found")
	}

	if user.Banned {
		return apperror.ForbiddenError(errors.New("account banned"), "account banned")
	}

	ctx.Locals("userID", userID)
	ctx.Locals("user", user)

	return ctx.Next()
}

func (a *AuthMiddleware) RequireAdmin(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	if user.Role != domain.AdminRole {
		return apperror.ForbiddenError(errors.New("admin access required"), "Admin access required")
	}
	return c.Next()
}
