package orders

import (
	"fmt"
	
	"backend/internal/middleware"
	"backend/internal/proj/orders/handler"
	"backend/internal/proj/orders/service"
	"backend/internal/storage"
	"backend/internal/storage/postgres"
	"backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// Module представляет модуль заказов
type Module struct {
	Handler *handler.OrdersHandler
}

// NewModule создает новый модуль заказов со всеми зависимостями
func NewModule(db storage.Storage) (*Module, error) {
	// Получаем репозитории из storage
	orderRepo := db.Order().(postgres.OrderRepositoryInterface)
	cartRepo := db.Cart().(postgres.CartRepositoryInterface)
	inventoryRepo := db.Inventory().(postgres.InventoryRepositoryInterface)
	storefrontRepo := db.Storefront().(service.StorefrontRepositoryInterface)

	// Создаем сервисы
	log := logger.New()
	inventoryManager := service.NewInventoryManager(inventoryRepo, nil, *log)

	// Пока используем nil для productRepo - TODO: реализовать позже
	orderService := service.NewOrderService(orderRepo, cartRepo, nil, storefrontRepo, inventoryManager, *log)

	// Получаем sqlx.DB для транзакций
	postgresDB, ok := db.(*postgres.Database)
	if !ok {
		return nil, fmt.Errorf("expected postgres.Database, got %T", db)
	}
	sqlxDB := postgresDB.GetSQLXDB()

	// Создаем handler с поддержкой транзакций
	ordersHandler := handler.NewOrdersHandler(orderService, inventoryManager, sqlxDB)

	return &Module{
		Handler: ordersHandler,
	}, nil
}

// GetPrefix возвращает префикс для маршрутов
func (m *Module) GetPrefix() string {
	return "/api/v1/orders"
}

// RegisterRoutes регистрирует маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	m.Handler.RegisterRoutes(app, mw)
	return nil
}
