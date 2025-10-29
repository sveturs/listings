package repository

import (
	"context"

	"backend/internal/domain/models"
)

// PostgresOrderAdapter адаптер для postgres репозитория заказов
type PostgresOrderAdapter struct {
	repo interface {
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
}

// NewPostgresOrderAdapter создает адаптер для postgres репозитория
func NewPostgresOrderAdapter(repo interface{}) OrderRepositoryInterface {
	// Проверяем что репозиторий реализует нужные методы
	if r, ok := repo.(interface {
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
	}); ok {
		return &PostgresOrderAdapter{repo: r}
	}

	panic("repo does not implement required methods")
}

// Реализация интерфейса OrderRepositoryInterface
func (a *PostgresOrderAdapter) Create(ctx context.Context, order *models.MarketplaceOrder) error {
	return a.repo.Create(ctx, order)
}

func (a *PostgresOrderAdapter) GetByID(ctx context.Context, id int64) (*models.MarketplaceOrder, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *PostgresOrderAdapter) GetByPaymentTransactionID(ctx context.Context, transactionID int64) (*models.MarketplaceOrder, error) {
	return a.repo.GetByPaymentTransactionID(ctx, transactionID)
}

func (a *PostgresOrderAdapter) UpdateStatus(ctx context.Context, orderID int64, newStatus models.MarketplaceOrderStatus, reason string, userID *int64) error {
	return a.repo.UpdateStatus(ctx, orderID, newStatus, reason, userID)
}

func (a *PostgresOrderAdapter) GetOrdersForAutoCapture(ctx context.Context) ([]*models.MarketplaceOrder, error) {
	return a.repo.GetOrdersForAutoCapture(ctx)
}

func (a *PostgresOrderAdapter) GetBuyerOrders(ctx context.Context, buyerID int64, limit, offset int) ([]*models.MarketplaceOrder, int, error) {
	return a.repo.GetBuyerOrders(ctx, buyerID, limit, offset)
}

func (a *PostgresOrderAdapter) GetSellerOrders(ctx context.Context, sellerID int64, limit, offset int) ([]*models.MarketplaceOrder, int, error) {
	return a.repo.GetSellerOrders(ctx, sellerID, limit, offset)
}

func (a *PostgresOrderAdapter) UpdateShippingInfo(ctx context.Context, orderID int64, shippingMethod string) error {
	return a.repo.UpdateShippingInfo(ctx, orderID, shippingMethod)
}

func (a *PostgresOrderAdapter) AddMessage(ctx context.Context, message *models.OrderMessage) error {
	return a.repo.AddMessage(ctx, message)
}

func (a *PostgresOrderAdapter) GetOrderMessages(ctx context.Context, orderID int64) ([]*models.OrderMessage, error) {
	return a.repo.GetOrderMessages(ctx, orderID)
}

// PostgresMarketplaceAdapter адаптер для postgres marketplace репозитория
type PostgresMarketplaceAdapter struct {
	repo interface {
		GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error)
	}
}

// NewPostgresMarketplaceAdapter создает адаптер для postgres marketplace репозитория
func NewPostgresMarketplaceAdapter(repo interface{}) MarketplaceRepositoryInterface {
	// Проверяем что репозиторий реализует нужные методы
	if r, ok := repo.(interface {
		GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error)
	}); ok {
		return &PostgresMarketplaceAdapter{repo: r}
	}

	panic("repo does not implement GetListingByID method")
}

func (a *PostgresMarketplaceAdapter) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	return a.repo.GetListingByID(ctx, id)
}
