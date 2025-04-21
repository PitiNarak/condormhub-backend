package services

import (
	"errors"

	"github.com/PitiNarak/condormhub-backend/internal/core/domain"
	"github.com/PitiNarak/condormhub-backend/internal/core/ports"
	"github.com/google/uuid"
	"github.com/yokeTH/go-pkg/apperror"
)

type OrderService struct {
	orderRepository          ports.OrderRepository
	leasingHistoryRepository ports.LeasingHistoryRepository
}

func NewOrderService(orderRepository ports.OrderRepository, leasingHistoryRepository ports.LeasingHistoryRepository) ports.OrderService {
	return &OrderService{orderRepository: orderRepository, leasingHistoryRepository: leasingHistoryRepository}
}

func (s *OrderService) CreateOrder(leasingHistoryID uuid.UUID) (*domain.Order, error) {
	leasingHistory, err := s.leasingHistoryRepository.GetByID(leasingHistoryID)
	if err != nil {
		return nil, apperror.NotFoundError(err, err.Error())
	}

	if !leasingHistory.End.IsZero() {
		return nil, apperror.BadRequestError(errors.New("leasing history has ended"), "leasing history has ended")
	}

	order := &domain.Order{
		LeasingHistoryID: leasingHistoryID,
		Price:            int64(leasingHistory.Dorm.Price),
		Type:             domain.MonthlyBillOrderType,
	}

	if err := s.orderRepository.Create(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(orderID uuid.UUID) (*domain.Order, error) {
	order, err := s.orderRepository.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetUnpaidOrderByUserID(userID uuid.UUID, limit int, page int) ([]domain.Order, int, int, error) {
	orders, totalPage, totalRows, err := s.orderRepository.GetUnpaidByUserID(userID, limit, page)
	if err != nil {
		return nil, 0, 0, err
	}
	return orders, totalPage, totalRows, nil
}

func (s *OrderService) UpdateOrder(order *domain.Order) error {
	return s.orderRepository.Update(order)
}

func (s *OrderService) DeleteOrder(orderID uuid.UUID) error {
	return s.orderRepository.Delete(orderID)
}
