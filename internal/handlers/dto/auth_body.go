package dto

import "github.com/PitiNarak/condormhub-backend/internal/core/domain"

type RegisterRequestBody struct {
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	UserName string `json:"username" validate:"required"`
}

type LoginRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenWithUserInformationResponseBody struct {
	AccessToken     string      `json:"accessToken"`
	RefreshToken    string      `json:"refreshToken"`
	UserInformation domain.User `json:"userInformation"`
}
