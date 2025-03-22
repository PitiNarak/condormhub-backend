package repository

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"

	"github.com/google/uuid"
)

type UserRepo struct {
	db *database.Database
}

func NewUserRepo(db *database.Database) ports.UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *domain.User) error {

	exitsUser := r.db.Model(&domain.User{}).Where("username = ?", user.Username).Or("email = ?", user.Email).First(&domain.User{})
	if exitsUser.RowsAffected > 0 {
		return apperror.BadRequestError(nil, "username already exists or email already exists")
	}

	err := r.db.Create(&user).Error

	if err != nil {
		return apperror.InternalServerError(err, "failed to save user to database")
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, apperror.NotFoundError(result.Error, "user not found")
	}

	return &user, nil
}

func (r *UserRepo) GetUserByID(userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, apperror.InternalServerError(result.Error, "user not found")
	}
	return &user, result.Error
}

func (r *UserRepo) UpdateUser(user *domain.User) error {
	result := r.db.Model(&user).Updates(user)
	if result.Error != nil {
		return apperror.InternalServerError(result.Error, "failed to update database")
	}
	return nil
}

func (r *UserRepo) UpdateInformation(userID uuid.UUID, data dto.UserInformationRequestBody) error {

	lifestyles := make([]domain.Lifestyle, len(data.Lifestyles))
	for i, v := range data.Lifestyles {
		lifestyles[i] = domain.Lifestyle(v)
	}

	user := domain.User{
		Username:        data.Username,
		Password:        data.Password,
		Firstname:       data.Firstname,
		Lastname:        data.Lastname,
		NationalID:      data.NationalID,
		Gender:          data.Gender,
		BirthDate:       data.BirthDate,
		Lifestyles:      lifestyles,
		PhoneNumber:     data.PhoneNumber,
		StudentEvidence: data.StudentEvidence,
	}

	if err := r.db.Model(&domain.User{}).Where("id = ?", userID).Updates(user).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update user information")
	}

	return nil
}

func (r *UserRepo) DeleteAccount(userID uuid.UUID) error {
	var user domain.User
	// TODO: Cascade delete?
	result := r.db.Delete(&user, userID)
	if result.Error != nil {
		return apperror.InternalServerError(result.Error, "cannot delete user")
	}
	return nil
}
