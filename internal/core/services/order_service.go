package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/PitiNarak/condormhub-backend/pkg/errorHandler"
	"github.com/google/uuid"
)

type OrderService struct {
	orderRepository          ports.OrderRepository
	leasingHistoryRepository ports.LeasingHistoryRepository
}

func NewOrderService(orderRepository ports.OrderRepository, leasingHistoryRepository ports.LeasingHistoryRepository) ports.OrderService {
	return &OrderService{orderRepository: orderRepository, leasingHistoryRepository: leasingHistoryRepository}
}

func (s *OrderService) CreateOrder(leasingHistoryID uuid.UUID) *errorHandler.ErrorHandler {
	leasingHistory, err := s.leasingHistoryRepository.GetByID(leasingHistoryID)
	if err != nil {
		return errorHandler.NotFoundError(err, err.Error())
	}

	if !leasingHistory.End.IsZero() {
		return errorHandler.BadRequestError(errors.New("leasing history has ended"), "leasing history has ended")
	}

	order := &domain.Order{
		LeasingHistoryID: leasingHistoryID,
		Price:            int64(leasingHistory.Dorm.Price),
		Type:             domain.MonthlyBillOrderType,
	}

	return s.orderRepository.Create(order)
}

func (s *OrderService) GetOrderByID(orderID uuid.UUID) (*domain.Order, *errorHandler.ErrorHandler) {
	order, err := s.orderRepository.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetUnpaidOrderByUserID(userID uuid.UUID, limit int, page int) ([]domain.Order, int, int, *errorHandler.ErrorHandler) {
	orders, totalPage, totalRows, err := s.orderRepository.GetUnpaidByUserID(userID, limit, page)
	if err != nil {
		return nil, 0, 0, err
	}
	return orders, totalPage, totalRows, nil
}

func (s *OrderService) UpdateOrder(order *domain.Order) *errorHandler.ErrorHandler {
	return s.orderRepository.Update(order)
}

func (s *OrderService) DeleteOrder(orderID uuid.UUID) *errorHandler.ErrorHandler {
	return s.orderRepository.Delete(orderID)
}
