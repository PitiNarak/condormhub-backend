package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo     ports.UserRepository
	EmailService ports.EmailServicePort
	Config       *utils.JWTConfig
}

func NewUserService(UserRepo ports.UserRepository, EmailService ports.EmailServicePort, config *utils.JWTConfig) ports.UserService {
	return &UserService{UserRepo: UserRepo, EmailService: EmailService, Config: config}
}

func (s *UserService) Create(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	create_err := s.UserRepo.Create(user)
	if create_err != nil {
		return create_err
	}
	err = s.EmailService.SendVerificationEmail(user.Email, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) VerifyUser(token string) error {
	claims, err := utils.DecodeJWT(token, s.Config)
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
	user, err := s.UserRepo.GetUser(userID)
	if err != nil || user.ID == uuid.Nil {
		return err
	}

	user.IsVerified = true
	return s.UserRepo.UpdateUser(user)
}
