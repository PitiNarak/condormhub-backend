package services

import (
	"errors"
	"time"

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
	jwtUtils     *utils.JWTUtils
}

func NewUserService(UserRepo ports.UserRepository, EmailService ports.EmailServicePort, jwtUtils *utils.JWTUtils) ports.UserService {
	return &UserService{userRepo: UserRepo, emailService: EmailService, jwtUtils: jwtUtils}
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

func (s *UserService) VerifyUser(token string) (string, *domain.User, error) {
	claims, err := s.jwtUtils.DecodeJWT(token)
	if err != nil {
		return "", nil, error_handler.UnauthorizedError(err, "invalid token")
	}

	if claims.GetExp() < time.Now().Unix() {
		return "", nil, error_handler.UnauthorizedError(errors.New("token expired"), "token is expired")
	}

	userID, err := uuid.Parse(claims.GetUserID())
	if err != nil {
		return "", nil, error_handler.UnauthorizedError(err, "invalid user ID")
	}
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil || user.ID == uuid.Nil {
		return "", nil, err
	}

	user.IsVerified = true

	updateErr := s.userRepo.UpdateUser(user)
	if updateErr != nil {
		return "", nil, updateErr
	}
	return token, user, nil
}

func (s *UserService) Login(email string, password string) (*domain.User, string, error) {
	user, getErr := s.userRepo.GetUserByEmail(email)
	if getErr != nil {
		return nil, "", getErr
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if compareErr != nil {
		return nil, "", error_handler.UnauthorizedError(compareErr, "invalid email or password.")
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
			return nil, error_handler.InternalServerError(err, "failed to hash password")
		}
		data.Password = string(hashedPassword)
	}

	err := s.userRepo.UpdateInformation(userID, data)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.userRepo.GetUserByID(userID)
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

func (s *UserService) ResetPassword(token string, password string) (*domain.User, error) {
	claims, err := s.jwtUtils.DecodeJWT(token)
	if err != nil {
		return new(domain.User), err
	}
	userIDstr := claims.UserID
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return new(domain.User), error_handler.InternalServerError(err, "cannot parse uuid")
	}
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return new(domain.User), err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return new(domain.User), error_handler.BadRequestError(err, "password cannot be hashed")
	}
	user.Password = string(hashedPassword)
	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return new(domain.User), err
	}
	return user, nil
}

func (s *UserService) DeleteAccount(userID uuid.UUID) error {
	err := s.userRepo.DeleteAccount(userID)
	if err != nil {
		return err
	}
	return nil
}
