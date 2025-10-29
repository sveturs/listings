package repository

import (
	"context"

	"backend/internal/domain/models"
)

// OrderRepositoryInterface интерфейс для работы с заказами маркетплейса
type OrderRepositoryInterface interface {
	Create(ctx context.Context, order *models.MarketplaceOrder) error
	GetByID(ctx context.Context, id int64) (*models.MarketplaceOrder, error)
	GetByPaymentTransactionID(ctx context.Context, transactionID int64) (*models.MarketplaceOrder, error)
	UpdateStatus(ctx context.Context, orderID int64, newStatus models.MarketplaceOrderStatus, reason string, userID *int64) error
	GetOrdersForAutoCapture(ctx context.Context) ([]*models.MarketplaceOrder, error)
	GetBuyerOrders(ctx context.Context, buyerID int64, limit, offset int) ([]*models.MarketplaceOrder, int, error)
	GetSellerOrders(ctx context.Context, sellerID int64, limit, offset int) ([]*models.MarketplaceOrder, int, error)
	UpdateShippingInfo(ctx context.Context, orderID int64, shippingMethod string) error
	AddMessage(ctx context.Context, message *models.OrderMessage) error
	GetOrderMessages(ctx context.Context, orderID int64) ([]*models.OrderMessage, error)
}

// MarketplaceRepositoryInterface интерфейс для работы с листингами
type MarketplaceRepositoryInterface interface {
	GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error)
}
