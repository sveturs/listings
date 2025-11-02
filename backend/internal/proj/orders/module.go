package orders

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
	"backend/internal/proj/delivery/grpcclient"
	"backend/internal/proj/orders/adapters"
	"backend/internal/proj/orders/handler"
	"backend/internal/proj/orders/service"
	"backend/internal/storage"
	"backend/internal/storage/postgres"
	"backend/pkg/logger"
)

// Module представляет модуль заказов
type Module struct {
	Handler *handler.OrdersHandler
}

// NewModule создает новый модуль заказов со всеми зависимостями
func NewModule(db storage.Storage, deliveryClient *grpcclient.Client) (*Module, error) {
	// Получаем репозитории из storage
	orderRepo := db.Order().(postgres.OrderRepositoryInterface)
	cartRepo := db.Cart().(postgres.CartRepositoryInterface)
	inventoryRepo := db.Inventory().(postgres.InventoryRepositoryInterface)
	storefrontRepo := db.Storefront().(service.StorefrontRepositoryInterface)

	// Получаем postgresDB для адаптера и транзакций
	postgresDB, ok := db.(*postgres.Database)
	if !ok {
		return nil, fmt.Errorf("expected postgres.Database, got %T", db)
	}

	// Создаем адаптер для работы с продуктами
	productRepo := adapters.NewProductRepositoryAdapter(postgresDB)

	// Создаем сервисы
	log := logger.New()
	inventoryManager := service.NewInventoryManager(inventoryRepo, nil, *log)

	// Создаем OrderService с deliveryClient для создания shipments
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo, storefrontRepo, inventoryManager, deliveryClient, *log)
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
