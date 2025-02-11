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
	userRepo     ports.UserRepository
	emailService ports.EmailServicePort
	config       *utils.JWTConfig
	jwtUtils     *utils.JWTUtils
}

func NewUserService(UserRepo ports.UserRepository, EmailService ports.EmailServicePort, jwtUtils *utils.JWTUtils, config *utils.JWTConfig) ports.UserService {
	return &UserService{userRepo: UserRepo, emailService: EmailService, config: config, jwtUtils: jwtUtils}
}

func (s *UserService) Create(user *domain.User) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	user.Password = string(hashedPassword)
	create_err := s.userRepo.Create(user)
	if create_err != nil {
		return "", create_err
	}

	token, generateErr := s.jwtUtils.GenerateJWT(user.ID)
	if generateErr != nil {
		return "", generateErr
	}

	err = s.emailService.SendVerificationEmail(user.Email, user.Username, token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) VerifyUser(token string) error {
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
	if err != nil || user.ID == uuid.Nil {
		return err
	}

	user.IsVerified = true
	return s.userRepo.UpdateUser(*user)
}

func (s *UserService) Login(email string, password string) (*domain.User, string, error) {
	user, getErr := s.userRepo.GetUserByEmail(email)
	if getErr != nil {
		return nil, "", getErr
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if compareErr != nil {
		return nil, "", error_handler.UnauthorizedError(compareErr, "Invalid email or password.")
	}
	token, generateErr := s.jwtUtils.GenerateJWT(user.ID)
	if generateErr != nil {
		return nil, "", generateErr
	}

	return user, token, nil

}

func (s *UserService) UpdateInformation(userID uuid.UUID, data dto.UserInformationRequestBody) (*domain.User, error) {
	if data.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, error_handler.InternalServerError(err, "Failed to hash password")
		}
		data.Password = string(hashedPassword)
	}

	err := s.userRepo.UpdateInformation(userID, data)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

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
