package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/core/ports"

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
		return result.Error
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
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
