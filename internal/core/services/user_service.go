package services

import (
	"context"
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/PitiNarak/condormhub-backend/pkg/email"
	"github.com/PitiNarak/condormhub-backend/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo     ports.UserRepository
	emailService email.Email
	jwtUtils     *jwt.JWTUtils
}

func NewUserService(UserRepo ports.UserRepository, EmailService email.Email, jwtUtils *jwt.JWTUtils) *UserService {
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

func (s *UserService) VerifyUser(ctx context.Context, token string) (*domain.User, string, string, error) {
	userID, err := s.jwtUtils.VerifyVerificationToken(ctx, token)
	if err != nil {
		return nil, "", "", err
	}

	if userID == uuid.Nil {
		return nil, "", "", apperror.UnauthorizedError(errors.New("token expired"), "token is expired")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil || user.ID == uuid.Nil {
		return nil, "", "", err
	}

	user.IsVerified = true

	updateErr := s.userRepo.UpdateUser(user)
	if updateErr != nil {
		return nil, "", "", updateErr
	}

	if err := s.jwtUtils.DeleteVerificationToken(ctx, userID); err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := s.jwtUtils.GenerateKeyPair(ctx, userID)
	if err != nil {
		return nil, "", "", apperror.InternalServerError(err, "generate key failed")
	}

	return user, accessToken, refreshToken, nil
}

func (s *UserService) Login(ctx context.Context, email string, password string) (*domain.User, string, string, error) {
	user, getErr := s.userRepo.GetUserByEmail(email)
	if getErr != nil {
		return nil, "", "", getErr
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if compareErr != nil {
		return nil, "", "", apperror.UnauthorizedError(compareErr, "invalid email or password.")
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
			return nil, apperror.InternalServerError(err, "failed to hash password")
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

func (s *UserService) GetUserByID(id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetUserByID(id)
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
		return apperror.InternalServerError(err, "cannot sent email")
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

func (s *UserService) ResetPassword(ctx context.Context, token string, password string) (*domain.User, string, string, error) {
	userID, err := s.jwtUtils.VerifyResetPasswordToken(ctx, token)
	if err != nil {
		return nil, "", "", err
	}

	if userID == uuid.Nil {
		return nil, "", "", apperror.UnauthorizedError(errors.New("token expired"), "token is expired")
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, "", "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", apperror.BadRequestError(err, "password cannot be hashed")
	}

	user.Password = string(hashedPassword)

	if err = s.userRepo.UpdateUser(user); err != nil {
		return nil, "", "", err
	}

	if err := s.jwtUtils.DeleteResetPasswordToken(ctx, userID); err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := s.jwtUtils.GenerateKeyPair(ctx, userID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *UserService) DeleteAccount(userID uuid.UUID) error {
	err := s.userRepo.DeleteAccount(userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ResendVerificationEmailService(ctx context.Context, email string) error {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	verifyToken, err := s.jwtUtils.GenerateVerificationToken(ctx, user.ID)
	if err != nil {
		return err
	}

	if err = s.emailService.SendVerificationEmail(user.Email, user.Username, verifyToken); err != nil {
		return err
	}
	return nil
}
