package middleware

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend/internal/logger"
)

// Список путей, которые не нужно логировать (слишком шумные)
var noiseLogPaths = []string{
	"/api/v1/admin/marketplace-translations/status",
	"/api/v1/admin/c2c-translations/status",
	"/api/v1/admin/tests/runs/", // Frequent polling during test execution
	"/ws/chat",                  // WebSocket chat connections (frequent reconnects)
}

// shouldSkipLogging проверяет нужно ли пропустить логирование для данного пути
func shouldSkipLogging(path string) bool {
	for _, noisePath := range noiseLogPaths {
		if strings.HasPrefix(path, noisePath) {
			return true
		}
	}
	return false
}

func (m *Middleware) Logger() fiber.Handler {
	masker := NewSensitiveDataMasker()

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Логируем запрос с маскированием sensitive данных
		path := c.Path()
		method := c.Method()

		// Пропускаем логирование для шумных endpoints
		skipLogging := shouldSkipLogging(path)

		if !skipLogging {
			// Маскируем query параметры если есть sensitive данные
			queryParams := c.Request().URI().QueryString()
			if len(queryParams) > 0 {
				maskedQuery := masker.Mask(string(queryParams))
				logger.Info().Str("method", method).Str("path", path).Str("query", maskedQuery).Msg("REQUEST")
			} else {
				logger.Info().Str("method", method).Str("path", path).Msg("REQUEST")
			}
		}

		// Логируем тело запроса только для POST/PUT/PATCH с маскированием
		if !skipLogging && (method == httpMethodPost || method == httpMethodPut || method == httpMethodPatch) {
			body := c.Body()
			if len(body) > 0 && len(body) < 10000 { // Не логируем большие тела запросов
				// Пытаемся распарсить как JSON для красивого вывода
				var jsonBody interface{}
				if err := json.Unmarshal(body, &jsonBody); err == nil {
					maskedBody := masker.Mask(string(body))
					logger.Info().Str("body", maskedBody).Msg("REQUEST BODY")
				}
			}
		}

		err := c.Next()

		// Логируем результат только если не skipLogging
		if !skipLogging {
			logger.Info().
				Str("method", method).
				Str("path", path).
				Int("status", c.Response().StatusCode()).
				Dur("duration", time.Since(start)).
				Msg("RESPONSE")
		}

		return err
	}
}
