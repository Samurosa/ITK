package application

import (
	"ITK_Code/m/v2/internal/core/order/repository"
	"ITK_Code/m/v2/internal/core/spot/client"

	"go.uber.org/zap"
)

type OrderService struct {
	log *zap.Logger

	repository repository.OrderRepository

	spotProvider client.SpotProvider
}

func NewOrderService(log *zap.Logger, repository repository.OrderRepository, spotProvider client.SpotProvider) *OrderService {
	return &OrderService{log: log, repository: repository, spotProvider: spotProvider}
}
