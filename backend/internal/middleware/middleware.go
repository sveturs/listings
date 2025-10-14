// backend/internal/middleware/middleware.go
package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/monitoring"
	globalService "backend/internal/proj/global/service"
	pkglogger "backend/pkg/logger"
	"backend/pkg/utils"

	authService "github.com/sveturs/auth/pkg/http/service"
)

type Middleware struct {
	config      *config.Config
	services    globalService.ServicesInterface
	metrics     *monitoring.MetricsCollector
	authService *authService.AuthService
	jwtParserMW fiber.Handler
}

// JWTParser возвращает middleware для парсинга JWT токенов из auth service
func (m *Middleware) JWTParser() fiber.Handler {
	return m.jwtParserMW
}

func NewMiddleware(cfg *config.Config, services globalService.ServicesInterface, authSvc *authService.AuthService, jwtParser fiber.Handler) *Middleware {
	return &Middleware{
		config:      cfg,
		services:    services,
		metrics:     monitoring.NewMetricsCollector(pkglogger.GetLogger()),
		authService: authSvc,
		jwtParserMW: jwtParser,
	}
}

// ErrorHandler обрабатывает все ошибки приложения
func (m *Middleware) ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	logger.Error().Err(err).Msg("Error handling request")

	return utils.ErrorResponse(c, code, message)
}

// GetJWTParser возвращает JWT parser middleware из sveturs/auth
func (m *Middleware) GetJWTParser() fiber.Handler {
	return m.jwtParserMW
}

// GetAuthService возвращает auth service из sveturs/auth
func (m *Middleware) GetAuthService() *authService.AuthService {
	return m.authService
}
