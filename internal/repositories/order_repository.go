package repositories

import (
	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/error_handler"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

type OrderFindOption func(r *OrderRepository)

func NewOrderRepository(db *gorm.DB) ports.OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *domain.Order) *error_handler.ErrorHandler {
	err := r.db.Create(order).Error
	if err != nil {
		return error_handler.InternalServerError(err, "Failed to create order")
	}
	return nil
}

func (r *OrderRepository) Update(order *domain.Order) *error_handler.ErrorHandler {
	err := r.db.Model(order).Where("session_id = ?", order.SessionID).Updates(order).Error
	if err != nil {
		return error_handler.InternalServerError(err, "Failed to create order")
	}
	return nil
}

// func (r *OrderRepository) Delete(orderId uuid.UUID) *error_handler.ErrorHandler {
// 	err := r.db.Where("id = ?", orderId).Delete(&domain.Order{}).Error
// 	if err != nil {
// 		return error_handler.InternalServerError(err, "Failed to delete order")
// 	}
// 	return nil
// }

// func (r *OrderRepository) Find(options ...OrderFindOption) ([]domain.Order, *error_handler.ErrorHandler) {
// 	var orders []domain.Order
// 	query := r.db
// 	for _, option := range options {
// 		option(r)
// 	}
// 	err := query.Find(&orders).Error
// 	if err != nil {
// 		return nil, error_handler.BadRequestError(err, "Failed to get orders")
// 	}
// 	return orders, nil
// }

// func (r *OrderRepository) GetById(orderId uuid.UUID) (*domain.Order, *error_handler.ErrorHandler) {
// 	var order domain.Order
// 	err := r.db.Where("id = ?", orderId).First(&order).Error
// 	if err != nil {
// 		return nil, error_handler.NotFoundError(err, "Order not found")
// 	}
// 	return &order, nil
// }

// func (r *OrderRepository) GetByLessorId(lessorId uuid.UUID) ([]domain.Order, *error_handler.ErrorHandler) {
// 	var orders []domain.Order
// 	err := r.db.Where("lessor_id = ?", lessorId).Find(&orders).Error
// 	if err != nil {
// 		return nil, error_handler.BadRequestError(err, "Failed to get orders")
// 	}
// 	return orders, nil
// }

// func (r *OrderRepository) GetByLesseeId(lesseeId uuid.UUID) ([]domain.Order, *error_handler.ErrorHandler) {
// 	var orders []domain.Order
// 	err := r.db.Where("lessee_id = ?", lesseeId).Find(&orders).Error
// 	if err != nil {
// 		return nil, error_handler.BadRequestError(err, "Failed to get orders")
// 	}
// 	return orders, nil
// }

// func (r *OrderRepository) GetByDormitoryId(dormitoryId uuid.UUID) ([]domain.Order, *error_handler.ErrorHandler) {
// 	var orders []domain.Order
// 	err := r.db.Where("dormitory_id = ?", dormitoryId).Find(&orders).Error
// 	if err != nil {
// 		return nil, error_handler.BadRequestError(err, "Failed to get orders")
// 	}
// 	return orders, nil
// }

// func (r *OrderRepository) GetAll() ([]domain.Order, *error_handler.ErrorHandler) {
// 	var orders []domain.Order
// 	err := r.db.Find(&orders).Error
// 	if err != nil {
// 		return nil, error_handler.BadRequestError(err, "Failed to get orders")
// 	}
// 	return orders, nil
// }
