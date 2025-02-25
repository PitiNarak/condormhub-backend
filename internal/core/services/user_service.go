package services

import (
	"context"
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/PitiNarak/condormhub-backend/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo     ports.UserRepository
	emailService ports.EmailServicePort
	jwtUtils     *jwt.JWTUtils
}

func NewUserService(UserRepo ports.UserRepository, EmailService ports.EmailServicePort, jwtUtils *jwt.JWTUtils) ports.UserService {
	return &UserService{userRepo: UserRepo, emailService: EmailService, jwtUtils: jwtUtils}
}

func (s *UserService) Create(ctx context.Context, user *domain.User) (string, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return "", "", err
	}

	user.Password = string(hashedPassword)
	create_err := s.userRepo.Create(user)
	if create_err != nil {
		return "", "", create_err
	}

	accessToken, refreshToken, err := s.jwtUtils.GenerateKeyPair(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	verifyToken, err := s.jwtUtils.GenerateVerificationToken(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	err = s.emailService.SendVerificationEmail(user.Email, user.Username, verifyToken)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *UserService) VerifyUser(ctx context.Context, token string) (string, *domain.User, error) {
	userID, err := s.jwtUtils.VerifyVerificationToken(ctx, token)
	if err != nil {
		return "", nil, err
	}

	if userID == uuid.Nil {
		return "", nil, errorHandler.UnauthorizedError(errors.New("token expired"), "token is expired")
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

	if err := s.jwtUtils.DeleteVerificationToken(ctx, userID); err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *UserService) Login(ctx context.Context, email string, password string) (*domain.User, string, string, error) {
	user, getErr := s.userRepo.GetUserByEmail(email)
	if getErr != nil {
		return nil, "", "", getErr
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if compareErr != nil {
		return nil, "", "", errorHandler.UnauthorizedError(compareErr, "invalid email or password.")
	}
	accessToken, refreshToken, generateErr := s.jwtUtils.GenerateKeyPair(ctx, user.ID)
	if generateErr != nil {
		return nil, "", "", generateErr
	}

	return user, accessToken, refreshToken, nil

}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	accessToken, newRefreshToken, err := s.jwtUtils.RefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}
	return accessToken, newRefreshToken, nil
}

func (s *UserService) UpdateInformation(userID uuid.UUID, data dto.UserInformationRequestBody) (*domain.User, error) {
	if data.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errorHandler.InternalServerError(err, "failed to hash password")
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

func (s *UserService) ResetPasswordCreate(ctx context.Context, email string) error {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(user.ID.String())
	if err != nil {
		return errorHandler.InternalServerError(err, "cannot sent email")
	}
	token, err := s.jwtUtils.GenerateResetPasswordToken(ctx, userID)
	if err != nil {
		return err
	}
	err = s.emailService.SendResetPasswordEmail(user.Email, user.Username, token)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ResetPassword(ctx context.Context, token string, password string) (*domain.User, error) {
	userID, err := s.jwtUtils.VerifyResetPasswordToken(ctx, token)
	if err != nil {
		return new(domain.User), err
	}

	if userID == uuid.Nil {
		return new(domain.User), errorHandler.UnauthorizedError(errors.New("token expired"), "token is expired")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return new(domain.User), err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return new(domain.User), errorHandler.BadRequestError(err, "password cannot be hashed")
	}

	user.Password = string(hashedPassword)

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return new(domain.User), err
	}

	if err := s.jwtUtils.DeleteResetPasswordToken(ctx, userID); err != nil {
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
