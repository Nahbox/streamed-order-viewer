package order

import (
	"github.com/Nahbox/streamed-order-viewer/service/internal/models"
)

// TODO: add mock tests
type DB interface {
	GetAllOrders() ([]models.OrderDetails, error)
	AddOrder(orderDetails *models.OrderDetails) error
	GetOrderById(userId int) (*models.OrderDetails, error)
	GetOrderDataById(orderId int) (*models.Delivery, *models.Payment, []models.Item, error)
}
