package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) ResetPasswordCreate(email string) error {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(user.ID.String())
	if err != nil {
		return error_handler.InternalServerError(err, "cannot sent email")
	}
	token, err := s.jwtUtils.GenerateJWT(userID)
	if err != nil {
		return err
	}
	err = s.emailService.SendResetPasswordEmail(user.Email, user.Username, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ResetPasswordResponse(token string, password string) error {
	claims, err := utils.DecodeJWT(token, s.config)
	if err != nil {
		return err
	}
	userIDstr, ok := (*claims)["user_id"].(string)
	if !ok {
		return errors.New("cannot get user_id")
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return err
	}
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	err = s.userRepo.UpdateUser(*user)
	if err != nil {
		return err
	}
	return nil
}
