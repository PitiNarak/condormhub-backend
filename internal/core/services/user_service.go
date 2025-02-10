package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"github.com/PitiNarak/condormhub-backend/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo     ports.UserRepository
	EmailService ports.EmailServicePort
	Config       *utils.JWTConfig
	jwtUtils     *utils.JWTUtils
}

func NewUserService(UserRepo ports.UserRepository, EmailService ports.EmailServicePort, jwtUtils *utils.JWTUtils, config *utils.JWTConfig) ports.UserService {
	return &UserService{UserRepo: UserRepo, EmailService: EmailService, Config: config, jwtUtils: jwtUtils}
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
	err = s.EmailService.SendVerificationEmail(user.Email, user.Username, user.ID)
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
	return s.UserRepo.UpdateUser(*user)
}

func (s *UserService) Login(email string, password string) (string, error) {
	user, getErr := s.UserRepo.GetUserByEmail(email)
	if getErr != nil {
		return "", getErr
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if compareErr != nil {
		return "", compareErr
	}
	token, generateErr := s.jwtUtils.GenerateJWT(user.ID)
	if generateErr != nil {
		return "", generateErr
	}

	return token, nil

}

func (s *UserService) UpdateInformation(userID uuid.UUID, data dto.UserRequestBody) (*domain.User, error) {
	if data.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, error_handler.InternalServerError(err, "Failed to hash password")
		}
		data.Password = string(hashedPassword)
	}

	err := s.UserRepo.UpdateInformation(userID, data)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.UserRepo.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
