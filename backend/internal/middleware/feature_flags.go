package middleware

import (
	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/pkg/utils"
)

// FeatureFlagsMiddleware проверяет доступность функций на основе feature flags
type FeatureFlagsMiddleware struct {
	featureFlags *config.FeatureFlags
}

// NewFeatureFlagsMiddleware создает новый middleware для feature flags
func NewFeatureFlagsMiddleware(featureFlags *config.FeatureFlags) *FeatureFlagsMiddleware {
	return &FeatureFlagsMiddleware{
		featureFlags: featureFlags,
	}
}

// CheckUnifiedAttributes проверяет доступность унифицированных атрибутов для пользователя
func (m *FeatureFlagsMiddleware) CheckUnifiedAttributes() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, _ := authMiddleware.GetUserID(c)

		// Проверяем, включена ли функция для пользователя
		if !m.featureFlags.ShouldUseUnifiedAttributes(userID) {
			// Если функция отключена, возвращаем ошибку или перенаправляем на старую версию
			if c.Path() == "/api/v2/c2c/categories/:category_id/attributes" {
				// Можно перенаправить на v1
				c.Path("/api/v1/c2c/categories/" + c.Params("category_id") + "/attributes")
			} else {
				return utils.SendError(c, fiber.StatusNotImplemented, "errors.featureNotAvailable")
			}
		}

		// Логируем использование новой системы если включено
		if m.featureFlags.LogAttributeSystemCalls {
			logger.Info().
				Int("user_id", userID).
				Str("path", c.Path()).
				Str("method", c.Method()).
				Msg("Unified attributes feature accessed")
		}

		return c.Next()
	}
}

// CheckFeaturePercentage проверяет процент включения функции
func (m *FeatureFlagsMiddleware) CheckFeaturePercentage(feature string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, _ := authMiddleware.GetUserID(c)

		// Проверяем процент включения
		percentage := m.featureFlags.GetFeaturePercentage(feature)
		if percentage < 100 {
			userGroup := userID % 100
			if userGroup >= percentage {
				return utils.SendError(c, fiber.StatusNotImplemented, "errors.featureNotAvailable")
			}
		}

		return c.Next()
	}
}

// LogFeatureUsage логирует использование функций
func (m *FeatureFlagsMiddleware) LogFeatureUsage(feature string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if m.featureFlags.LogAttributeSystemCalls {
			userID, _ := authMiddleware.GetUserID(c)
			logger.Info().
				Str("feature", feature).
				Int("user_id", userID).
				Str("path", c.Path()).
				Str("method", c.Method()).
				Msg("Feature accessed")
		}
		return c.Next()
	}
}

// DynamicVersionRouting выбирает версию API на основе feature flags
func (m *FeatureFlagsMiddleware) DynamicVersionRouting() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, _ := authMiddleware.GetUserID(c)

		// Если путь содержит /api/unified/ - проверяем доступность
		if c.Path()[:12] == "/api/unified" {
			if m.featureFlags.ShouldUseUnifiedAttributes(userID) {
				// Заменяем на v2
				newPath := "/api/v2" + c.Path()[12:]
				c.Path(newPath)
			} else {
				// Заменяем на v1
				newPath := "/api/v1" + c.Path()[12:]
				c.Path(newPath)
			}
		}

		return c.Next()
	}
}
