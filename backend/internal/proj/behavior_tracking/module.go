package behavior_tracking

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/interfaces"
	"backend/internal/middleware"
	"backend/internal/proj/behavior_tracking/handler"
	"backend/internal/proj/behavior_tracking/service"
	"backend/internal/proj/behavior_tracking/storage/postgres"
)

// Module представляет модуль поведенческой аналитики
type Module struct {
	handler     *handler.BehaviorTrackingHandler
	jwtParserMW fiber.Handler
}

// NewModule создает новый модуль поведенческой аналитики
func NewModule(ctx context.Context, pool *pgxpool.Pool, jwtParserMW fiber.Handler) *Module {
	// Создаем репозиторий
	repo := postgres.NewBehaviorTrackingRepository(pool)

	// Создаем сервис
	svc := service.NewBehaviorTrackingService(ctx, repo)

	// Создаем обработчик
	h := handler.NewBehaviorTrackingHandler(svc)

	return &Module{
		handler:     h,
		jwtParserMW: jwtParserMW,
	}
}

// RegisterRoutes регистрирует маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error {
	// Публичный API для трекинга событий
	api := app.Group("/api/v1/analytics")

	// Endpoint для отслеживания событий (может использоваться без авторизации)
	api.Post("/track", m.handler.TrackEvent)

	// Endpoint для получения событий сессии (публичный)
	api.Get("/sessions/:session_id/events", m.handler.GetSessionEvents)

	// Защищенные endpoints (требуют авторизации)
	protected := api.Group("/",
		m.jwtParserMW,
		authMiddleware.RequireAuth())

	// События пользователя (доступны самому пользователю и админам)
	protected.Get("/users/:user_id/events", m.handler.GetUserEvents)

	// Админские endpoints
	admin := api.Group("/",
		m.jwtParserMW,
		authMiddleware.RequireAuth("admin"))

	// Обновление агрегированных метрик
	admin.Post("/metrics/update", m.handler.UpdateMetrics)

	return nil
}

// GetPrefix возвращает префикс модуля для логирования
func (m *Module) GetPrefix() string {
	return "behavior_tracking"
}

// Ensure Module implements RouteRegistrar interface
var _ interfaces.RouteRegistrar = (*Module)(nil)
