package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
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

	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, error_handler.NotFoundError(err, "User not found with the provided email.")
		}
		return nil, error_handler.InternalServerError(err, "Error retrieving user data.")
	}

	return &user, result.Error
}

func (r *UserRepo) GetUser(userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, result.Error
}

func (r *UserRepo) UpdateUser(user domain.User) error {
	result := r.db.Model(&user).Updates(user)

	return result.Error
}

func (r *UserRepo) Update(email string, updateInfo domain.UpdateInfo) error {
	var user domain.User

	// Find the user by email
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return err // Return error if user not found
	}
	result := r.db.Model(&user).Updates(updateInfo)

	return result.Error
}
