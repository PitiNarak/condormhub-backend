package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/handlers/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"

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

	exitsUser := r.db.Model(&domain.User{}).Where("username = ?", user.Username).Or("email = ?", user.Email).First(&domain.User{})
	if exitsUser.RowsAffected > 0 {
		return errorHandler.BadRequestError(nil, "username already exists or email already exists")
	}

	err := r.db.Create(&user).Error

	if err != nil {
		return errorHandler.InternalServerError(err, "failed to save user to database")
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, errorHandler.NotFoundError(result.Error, "user not found")
	}

	return &user, nil
}

func (r *UserRepo) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, errorHandler.InternalServerError(result.Error, "user not found")
	}
	return &user, result.Error
}

func (r *UserRepo) UpdateUser(user *domain.User) error {
	result := r.db.Model(&user).Updates(user)
	if result.Error != nil {
		return errorHandler.InternalServerError(result.Error, "failed to update database")
	}
	return nil
}

func (r *UserRepo) UpdateInformation(userID uuid.UUID, data dto.UserInformationRequestBody) error {
	err := r.db.Model(&domain.User{}).Where("id = ?", userID).Updates(data).Error
	if err != nil {
		return errorHandler.InternalServerError(err, "failed to update user information")
	}

	return nil
}

func (r *UserRepo) DeleteAccount(userID uuid.UUID) error {
	var user domain.User
	result := r.db.Delete(&user, userID)
	if result.Error != nil {
		return errorHandler.InternalServerError(result.Error, "cannot delete user")
	}
	return nil
}
