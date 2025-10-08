package orders

import (
	"fmt"

	"backend/internal/middleware"
	storefrontsOpensearch "backend/internal/proj/b2c/storage/opensearch"
	"backend/internal/proj/orders/adapters"
	"backend/internal/proj/orders/handler"
	"backend/internal/proj/orders/service"
	"backend/internal/storage"
	opensearchClient "backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"
	"backend/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// Module представляет модуль заказов
type Module struct {
	Handler *handler.OrdersHandler
}

// NewModule создает новый модуль заказов со всеми зависимостями
func NewModule(db storage.Storage, opensearchCfg *opensearchClient.Config) (*Module, error) {
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

	// Создаем OpenSearch репозиторий для обновления остатков
	var productSearchRepo storefrontsOpensearch.ProductSearchRepository
	if opensearchCfg != nil && opensearchCfg.URL != "" {
		osClient, err := opensearchClient.NewOpenSearchClient(*opensearchCfg)
		if err != nil {
			// Логируем ошибку, но продолжаем без OpenSearch
			log := logger.New()
			log.Error("Failed to create OpenSearch client for orders module: %v", err)
		} else {
			// Используем правильный индекс для товаров витрин
			// ВАЖНО: Используем тот же индекс что и для marketplace (унифицированный поиск)
			productSearchRepo = storefrontsOpensearch.NewProductRepository(osClient, "marketplace_listings")
		}
	}

	// Создаем сервисы
	log := logger.New()
	inventoryManager := service.NewInventoryManager(inventoryRepo, nil, *log)

	// Передаем productSearchRepo в OrderService
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo, storefrontRepo, inventoryManager, productSearchRepo, *log)
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
