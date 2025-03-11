package repository

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/database"
	"github.com/PitiNarak/condormhub-backend/pkg/apperror"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *database.Database
}

func NewOrderRepository(db *database.Database) ports.OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(order *domain.Order) error {
	if err := r.db.Create(order).Error; err != nil {
		return apperror.InternalServerError(err, "failed to create order")
	}

	if err := r.db.Preload("LeasingHistory").First(order, order.ID).Error; err != nil {
		return apperror.InternalServerError(err, "failed to preload leasing history")
	}

	return nil
}

func (r *OrderRepository) GetByID(orderID uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	if err := r.db.
		Where("id = ?", orderID).
		Preload("LeasingHistory").
		Preload("LeasingHistory.Dorm").
		Preload("LeasingHistory.Lessee").
		Preload("PaidTransaction").
		First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFoundError(err, "order not found")
		}
		return nil, apperror.InternalServerError(err, "database error retrieving order")
	}
	return &order, nil
}

func (r *OrderRepository) GetUnpaidByUserID(userID uuid.UUID, limit int, page int) ([]domain.Order, int, int, error) {
	var orders []domain.Order
	query := r.db.
		Joins("JOIN leasing_histories ON leasing_histories.id = orders.leasing_history_id").
		Where("leasing_histories.lessee_id = ?", userID).
		Where("orders.paid_transaction_id IS NULL")

	totalPage, totalRows, err := r.db.Paginate(orders, query, limit, page, "create_at desc")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, apperror.NotFoundError(err, "order not found")
		}
		return nil, 0, 0, apperror.InternalServerError(err, "failed to orders")
	}

	return orders, totalPage, totalRows, nil
}

func (r *OrderRepository) Update(order *domain.Order) error {
	if err := r.db.Model(order).Where("id = ?", order.ID).Updates(order).Error; err != nil {
		return apperror.InternalServerError(err, "failed to update order")
	}
	return nil
}

func (r *OrderRepository) Delete(orderID uuid.UUID) error {
	if err := r.db.Where("id = ?", orderID).Delete(&domain.Order{}).Error; err != nil {
		return apperror.InternalServerError(err, "failed to delete order")
	}
	return nil
}
