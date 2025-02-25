package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/internal/databases"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) ports.OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(order *domain.Order) *errorHandler.ErrorHandler {
	if err := r.db.Create(&order).Preload("LeasingHistory").Error; err != nil {
		return errorHandler.InternalServerError(err, "failed to create order")
	}
	return nil
}

func (r *OrderRepository) GetByID(orderID uuid.UUID) (*domain.Order, *errorHandler.ErrorHandler) {
	var order domain.Order
	if err := r.db.Where("id = ?", orderID).Preload("PaidTransaction").First(&order).Error; err != nil {
		return nil, errorHandler.NotFoundError(err, "order not found")
	}
	return &order, nil
}

func (r *OrderRepository) GetUnpaidByUserID(userID uuid.UUID, limit int, page int) ([]domain.Order, int, int, *errorHandler.ErrorHandler) {
	var orders []domain.Order
	query := r.db.Preload("LeasingHistory").Preload("LeasingHistory.Lessee").Where("leasing_history.lessee_id = ?", userID).Where("paid_transaction_id IS NULL")
	scope, totalPage, totalRows, err := databases.Paginate(orders, query, limit, page, "create_at desc")
	if err != nil {
		return nil, 0, 0, errorHandler.BadRequestError(err, err.Error())
	}

	if err := query.Scopes(scope).Find(&orders).Error; err != nil {
		return nil, 0, 0, errorHandler.InternalServerError(err, "failed to get unpaid orders")
	}
	return orders, totalPage, totalRows, nil
}

func (r *OrderRepository) Update(order *domain.Order) *errorHandler.ErrorHandler {
	if err := r.db.Save(&order).Error; err != nil {
		return errorHandler.InternalServerError(err, "failed to update order")
	}
	return nil
}

func (r *OrderRepository) Delete(orderID uuid.UUID) *errorHandler.ErrorHandler {
	if err := r.db.Where("id = ?", orderID).Delete(&domain.Order{}).Error; err != nil {
		return errorHandler.InternalServerError(err, "failed to delete order")
	}
	return nil
}
