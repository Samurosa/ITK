package application

import (
	"ITK_Code/m/v2/internal/core/order/repository"

	"go.uber.org/zap"
)

type OrderService struct {
	log *zap.Logger

	repository repository.OrderRepository
}

func NewOrderService(log *zap.Logger, repository repository.OrderRepository) *OrderService {
	return &OrderService{log: log, repository: repository}
}
