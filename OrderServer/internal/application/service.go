package application

import (
	"ITK_Code/m/v2/internal/core/order/service"

	"go.uber.org/zap"
)

type OrderService struct {
	log zap.Logger

	order service.Order
}

func NewOrderService(log zap.Logger, order service.Order) *OrderService {
	return &OrderService{log: log, order: order}
}
