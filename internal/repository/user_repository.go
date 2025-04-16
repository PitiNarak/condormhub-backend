package repository

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
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
		return nil, apperror.NotFoundError(result.Error, "user not found")
	}
	return &user, result.Error
}

func (r *UserRepo) UpdateUser(user *domain.User) error {
	result := r.db.Model(&user).Save(user)
	if result.Error != nil {
		return apperror.InternalServerError(result.Error, "failed to update database")
	}
	return nil
}

func (r *UserRepo) UpdateInformation(userID uuid.UUID, data domain.User) error {

	if err := r.db.Model(&domain.User{}).Where("id = ?", userID).Updates(data).Error; err != nil {
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

func (r *UserRepo) GetLessorIncome(lessorID uuid.UUID) (float64, error) {
	var income float64
	err := r.db.Model(&domain.LeasingHistory{}).Joins("JOIN dorms ON dorms.id = leasing_histories.dorm_id").
		Where("dorms.owner_id = ?", lessorID).
		Where("leasing_histories.end IS NULL").
		Select("COALESCE(SUM(leasing_histories.price), 0)").
		Scan(&income).Error
	if err != nil {
		return 0, apperror.InternalServerError(err, "failed to calculate lessor's income")
	}
	return income, nil
}

func (r *UserRepo) GetPending(limit int, page int) ([]domain.User, int, int, error) {
	var pending []domain.User
	query := r.db.Where("is_student_verified = ?", domain.StatusPending)

	totalPages, totalRows, err := r.db.Paginate(&pending, query, limit, page, "update_at DESC")
	if err != nil {
		return nil, 0, 0, apperror.InternalServerError(err, "Failed to load lessee with pending verification")
	}
	return pending, totalPages, totalRows, nil
}
