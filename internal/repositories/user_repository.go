package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) ports.UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *domain.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		return error_handler.InternalServerError(result.Error, "Failed to save user to database")
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, error_handler.NotFoundError(result.Error, "user not found")
	}

	return &user, nil
}

func (r *UserRepo) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, error_handler.InternalServerError(result.Error, "user not found")
	}
	return &user, result.Error
}

func (r *UserRepo) UpdateUser(user domain.User) error {
	result := r.db.Model(&user).Updates(user)

	return result.Error
}

func (r *UserRepo) UpdateInformation(userID uuid.UUID, data dto.UserInformationRequestBody) error {
	err := r.db.Model(&domain.User{}).Where("id = ?", userID).Updates(data).Error
	if err != nil {
		return error_handler.InternalServerError(err, "Failed to update user information")
	}

	return nil
}
