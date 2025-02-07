package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internals/core/domain"
	"github.com/PitiNarak/condormhub-backend/internals/core/ports"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) ports.UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user domain.User) error {
	result := r.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *UserRepo) GetUserViaEmail(email string) (domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)
	return user, result.Error
}

func (r *UserRepo) Update(user domain.User) (domain.User, error) {
	myUser, err := r.GetUserViaEmail(user.Email)
	if err != nil {
		return user, err
	}
	user.ID = myUser.ID
	result := r.db.Save(&user)
	return user, result.Error
}
