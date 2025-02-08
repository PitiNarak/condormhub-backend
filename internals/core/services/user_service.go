package services

import (
	"os"
	"time"

	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/core/ports"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo ports.UserRepository
}

func NewUserService(UserRepo ports.UserRepository) ports.UserService {
	return &UserService{UserRepo: UserRepo}
}

func (s *UserService) Create(user domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	create_err := s.UserRepo.Create(user)
	if create_err != nil {
		return create_err
	}

	return nil
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

	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil

}

func (s *UserService) Update(user domain.User) (domain.User, error) {
	updatedUser, err := s.UserRepo.Update(user)
	return updatedUser, err
}
