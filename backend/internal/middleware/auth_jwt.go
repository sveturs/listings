// backend/internal/middleware/auth_jwt.go
package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend/internal/logger"
	"backend/pkg/jwt"
	"backend/pkg/utils"
)

// AuthRequiredJWT - основной метод аутентификации через JWT
// Поддерживает как Bearer токены в заголовке, так и fallback на session cookies
func (m *Middleware) AuthRequiredJWT(c *fiber.Ctx) error {
	logger.Info().Str("path", c.Path()).Msg("AuthRequiredJWT middleware called")

	// Временное решение: пропускаем публичные маршруты storefronts и аналитики
	path := c.Path()

	// Пропускаем все публичные API маршруты
	if strings.HasPrefix(path, "/api/v1/public/") {
		logger.Info().Str("path", path).Msg("Skipping auth for public API routes")
		return c.Next()
	}

	// Пропускаем маршруты аналитики для записи событий (но НЕ метрики - они для админов)
	if strings.HasPrefix(path, "/api/v1/analytics/track") ||
		strings.HasPrefix(path, "/api/v1/analytics/event") ||
		strings.HasPrefix(path, "/api/v1/analytics/sessions/") {
		logger.Info().Str("path", path).Msg("Skipping auth for public analytics routes")
		return c.Next()
	}

	// Пропускаем публичные маршруты marketplace для автодополнения и поиска
	if strings.HasPrefix(path, "/api/v1/marketplace/") && c.Method() == httpMethodGet {
		logger.Info().Str("path", path).Str("method", c.Method()).Msg("Checking marketplace route")
		if strings.HasSuffix(path, "/suggestions") ||
			strings.Contains(path, "/search/autocomplete") ||
			strings.HasSuffix(path, "/category-suggestions") ||
			path == "/api/v1/marketplace/search" {
			logger.Info().Str("path", path).Msg("Skipping auth for public marketplace search routes")
			return c.Next()
		}
		logger.Info().Str("path", path).Msg("Marketplace route requires auth")
	}

	// Пропускаем публичные GIS маршруты
	if strings.HasPrefix(path, "/api/v1/gis") {
		method := c.Method()
		// Публичные GET routes
		if method == httpMethodGet && (strings.HasSuffix(path, "/search") ||
			strings.HasSuffix(path, "/search/radius") ||
			strings.HasSuffix(path, "/clusters") ||
			strings.HasSuffix(path, "/nearby") ||
			strings.Contains(path, "/listings/") && strings.HasSuffix(path, "/location") ||
			strings.HasSuffix(path, "/districts") ||
			strings.Contains(path, "/districts/") ||
			strings.HasSuffix(path, "/municipalities") ||
			strings.Contains(path, "/municipalities/") ||
			strings.Contains(path, "/search/by-district/") ||
			strings.Contains(path, "/search/by-municipality/") ||
			strings.HasSuffix(path, "/cities")) {
			logger.Info().Str("path", path).Msg("Skipping auth for public GIS routes")
			return c.Next()
		}
		// Публичные POST routes для cities
		if method == httpMethodPost && strings.HasSuffix(path, "/cities/visible") {
			logger.Info().Str("path", path).Msg("Skipping auth for public GIS cities routes")
			return c.Next()
		}
		// Публичные Geocoding API routes (Phase 2)
		if strings.Contains(path, "/geocode/") && ((method == httpMethodGet && (strings.HasSuffix(path, "/suggestions") ||
			strings.HasSuffix(path, "/search") ||
			strings.HasSuffix(path, "/reverse") ||
			strings.HasSuffix(path, "/stats"))) ||
			(method == httpMethodPost && strings.HasSuffix(path, "/validate"))) {
			logger.Info().Str("path", path).Msg("Skipping auth for public GIS geocoding routes")
			return c.Next()
		}
	}

	// Определяем, является ли это публичным маршрутом
	isPublicRoute := false

	if strings.HasPrefix(path, "/api/v1/storefronts") && !strings.Contains(path, "/my") {
		// Проверяем что это не защищенные маршруты
		method := c.Method()

		// Пропускаем роуты корзины - они используют OptionalAuth в orders module
		if strings.Contains(path, "/cart") {
			logger.Info().Str("path", path).Msg("Skipping auth for cart route - handled by orders module")
			return c.Next()
		}

		if method == "GET" && (strings.Contains(path, "/slug/") || strings.HasSuffix(path, "/storefronts") ||
			strings.Contains(path, "/search") || strings.Contains(path, "/nearby") ||
			strings.Contains(path, "/map") || strings.Contains(path, "/building") ||
			strings.Contains(path, "/staff") || strings.Contains(path, "/products/")) {
			logger.Info().Str("path", path).Msg("Public storefront route detected")
			isPublicRoute = true
		}
		if method == "POST" && strings.Contains(path, "/view") {
			logger.Info().Str("path", path).Msg("Public storefront view tracking route")
			isPublicRoute = true
		}
		// Проверка на ID маршрут
		parts := strings.Split(path, "/")
		if len(parts) == 5 && parts[3] == "storefronts" && method == "GET" {
			logger.Info().Str("path", path).Msg("Public storefront by ID route")
			isPublicRoute = true
		}
	}

	// Приоритет 1: Проверяем JWT токен в заголовке Authorization
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		// Извлекаем токен из заголовка "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == bearerScheme {
			tokenString := parts[1]

			// Сначала пробуем валидировать как токен от auth service (RS256)
			if m.authServicePubKey != nil {
				authClaims, err := jwt.ValidateAuthServiceToken(tokenString, m.authServicePubKey)
				if err == nil {
					// Токен от auth service валиден
					logger.Info().
						Int("user_id", authClaims.UserID).
						Str("email", authClaims.Email).
						Str("provider", authClaims.Provider).
						Msg("Auth service token validated")
					
					// Проверяем что пользователь существует
					user, err := m.services.User().GetUserByID(c.Context(), authClaims.UserID)
					if err != nil || user == nil {
						logger.Info().
							Int("user_id", authClaims.UserID).
							Msg("User not found for auth service token")
						return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.user_not_found")
					}
					
					// Устанавливаем данные пользователя в контекст
					c.Locals("user_id", authClaims.UserID)
					c.Locals("user_email", authClaims.Email)
					c.Locals("auth_method", "jwt_auth_service")
					c.Locals("jwt_token", tokenString)
					
					// Проверяем роли
					isAdmin := false
					for _, role := range authClaims.Roles {
						if role == "admin" {
							isAdmin = true
							break
						}
					}
					// Дополнительная проверка по email
					if !isAdmin {
						isAdmin, _ = m.services.User().IsUserAdmin(c.Context(), authClaims.Email)
					}
					c.Locals("is_admin", isAdmin)
					
					// Обновляем последнее посещение асинхронно
					go func() {
						ctx := context.Background()
						_ = m.services.User().UpdateLastSeen(ctx, authClaims.UserID)
					}()
					
					return c.Next()
				}
			}

			// Если RS256 не сработал, пробуем HS256 (старый формат)
			claims, err := jwt.ValidateToken(tokenString, m.config.JWTSecret)
			if err != nil {
				logger.Error().
					Err(err).
					Str("auth_type", "bearer").
					Msg("JWT validation failed")
				return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.invalid_token")
			}

			// Проверяем что пользователь существует
			user, err := m.services.User().GetUserByID(c.Context(), claims.UserID)
			if err != nil || user == nil {
				logger.Info().
					Int("user_id", claims.UserID).
					Msg("User not found for JWT")
				return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.user_not_found")
			}

			// Устанавливаем данные пользователя в контекст
			c.Locals("user_id", claims.UserID)
			c.Locals("user_email", claims.Email)
			c.Locals("auth_method", "jwt")
			c.Locals("jwt_token", tokenString)

			// Проверяем, является ли пользователь администратором
			isAdmin, err := m.services.User().IsUserAdmin(c.Context(), claims.Email)
			if err != nil {
				logger.Warn().
					Err(err).
					Str("email", claims.Email).
					Msg("Failed to check admin status")
				// В случае ошибки считаем, что пользователь не админ
				isAdmin = false
			}
			c.Locals("is_admin", isAdmin)

			if isAdmin {
				logger.Info().
					Int("user_id", claims.UserID).
					Str("email", claims.Email).
					Msg("Admin user authenticated")
			}

			// Обновляем последнее посещение асинхронно
			go func() {
				ctx := context.Background()
				_ = m.services.User().UpdateLastSeen(ctx, claims.UserID)
			}()

			return c.Next()
		}
	}

	// Приоритет 2: Проверяем JWT в query параметре (для WebSocket)
	tokenFromQuery := c.Query("token")
	if tokenFromQuery != "" {
		// Сначала пробуем валидировать как токен от auth service (RS256)
		if m.authServicePubKey != nil {
			authClaims, err := jwt.ValidateAuthServiceToken(tokenFromQuery, m.authServicePubKey)
			if err == nil {
				// Токен от auth service валиден
				c.Locals("user_id", authClaims.UserID)
				c.Locals("user_email", authClaims.Email)
				c.Locals("auth_method", "jwt_query_auth_service")
				c.Locals("jwt_token", tokenFromQuery)
				
				// Проверяем роли
				isAdmin := false
				for _, role := range authClaims.Roles {
					if role == "admin" {
						isAdmin = true
						break
					}
				}
				if !isAdmin {
					isAdmin, _ = m.services.User().IsUserAdmin(c.Context(), authClaims.Email)
				}
				c.Locals("is_admin", isAdmin)
				
				go func() {
					ctx := context.Background()
					_ = m.services.User().UpdateLastSeen(ctx, authClaims.UserID)
				}()
				
				return c.Next()
			}
		}

		// Если RS256 не сработал, пробуем HS256
		claims, err := jwt.ValidateToken(tokenFromQuery, m.config.JWTSecret)
		if err != nil {
			logger.Error().
				Err(err).
				Str("auth_type", "query").
				Msg("JWT validation failed")
		} else {
			// JWT из query валиден
			c.Locals("user_id", claims.UserID)
			c.Locals("user_email", claims.Email)
			c.Locals("auth_method", "jwt_query")
			c.Locals("jwt_token", tokenFromQuery)

			// Проверяем, является ли пользователь администратором
			isAdmin, err := m.services.User().IsUserAdmin(c.Context(), claims.Email)
			if err != nil {
				logger.Warn().
					Err(err).
					Str("email", claims.Email).
					Msg("Failed to check admin status")
				isAdmin = false
			}
			c.Locals("is_admin", isAdmin)

			go func() {
				ctx := context.Background()
				_ = m.services.User().UpdateLastSeen(ctx, claims.UserID)
			}()

			return c.Next()
		}
	}

	// Приоритет 3: Проверяем JWT в cookie (для web клиентов)
	jwtCookie := c.Cookies("jwt_token")
	if jwtCookie != "" {

		claims, err := jwt.ValidateToken(jwtCookie, m.config.JWTSecret)
		if err != nil {
			logger.Error().
				Err(err).
				Str("auth_type", "cookie").
				Msg("JWT validation failed")
			// Очищаем невалидную cookie
			c.Cookie(&fiber.Cookie{
				Name:     "jwt_token",
				Value:    "",
				Path:     "/",
				Domain:   m.config.GetCookieDomain(),
				MaxAge:   -1,
				Secure:   m.config.GetCookieSecure(),
				HTTPOnly: true,
				SameSite: m.config.GetCookieSameSite(),
			})
		} else {
			// JWT из cookie валиден
			c.Locals("user_id", claims.UserID)
			c.Locals("user_email", claims.Email)
			c.Locals("auth_method", "jwt_cookie")
			c.Locals("jwt_token", jwtCookie)

			// Проверяем, является ли пользователь администратором
			isAdmin, err := m.services.User().IsUserAdmin(c.Context(), claims.Email)
			if err != nil {
				logger.Warn().
					Err(err).
					Str("email", claims.Email).
					Msg("Failed to check admin status")
				isAdmin = false
			}
			c.Locals("is_admin", isAdmin)

			go func() {
				ctx := context.Background()
				_ = m.services.User().UpdateLastSeen(ctx, claims.UserID)
			}()

			return c.Next()
		}
	}

	// Приоритет 4: Fallback на session cookie для обратной совместимости
	sessionToken := c.Cookies("session_token")
	if sessionToken != "" {
		logger.Info().Msg("Falling back to session authentication")

		session, err := m.services.Auth().GetSession(c.Context(), sessionToken)
		if err != nil {
			logger.Error().
				Err(err).
				Str("auth_type", "session").
				Msg("Session validation failed")
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.invalid_session")
		}

		if session == nil || session.UserID == 0 {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.invalid_session")
		}

		// Генерируем JWT токен для этой сессии
		newJWT, err := jwt.GenerateTokenWithDuration(
			session.UserID,
			session.Email,
			m.config.JWTSecret,
			time.Duration(m.config.JWTExpirationHours)*time.Hour,
		)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to generate JWT from session")
		} else {
			// Устанавливаем JWT cookie для будущих запросов
			c.Cookie(&fiber.Cookie{
				Name:     "jwt_token",
				Value:    newJWT,
				Path:     "/",
				Domain:   m.config.GetCookieDomain(),
				MaxAge:   m.config.JWTExpirationHours * 3600,
				Secure:   m.config.GetCookieSecure(),
				HTTPOnly: true,
				SameSite: m.config.GetCookieSameSite(),
			})

			// Добавляем JWT в заголовок ответа для API клиентов
			c.Set("X-Auth-Token", newJWT)
		}

		c.Locals("user_id", session.UserID)
		c.Locals("user_email", session.Email)
		c.Locals("session_token", sessionToken)
		c.Locals("auth_method", "session_fallback")
		c.Locals("auth_provider", session.Provider)

		// Проверяем, является ли пользователь администратором
		isAdmin, err := m.services.User().IsUserAdmin(c.Context(), session.Email)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("email", session.Email).
				Msg("Failed to check admin status")
			isAdmin = false
		}
		c.Locals("is_admin", isAdmin)

		go func() {
			ctx := context.Background()
			_ = m.services.User().UpdateLastSeen(ctx, session.UserID)
		}()

		return c.Next()
	}

	// Если ни один метод аутентификации не сработал
	// Проверяем, является ли это публичным маршрутом
	if isPublicRoute {
		// Для публичных маршрутов разрешаем доступ без аутентификации
		logger.Info().Str("path", path).Msg("Allowing unauthenticated access to public route")
		return c.Next()
	}

	if c.Get("Upgrade") == "websocket" {
		logger.Info().
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Str("connection_type", "websocket").
			Msg("SECURITY: Unauthorized connection attempt")
	}

	return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.authentication_required")
}

// OptionalAuthJWT - опциональная аутентификация, не требует токена но извлекает данные если есть
func (m *Middleware) OptionalAuthJWT(c *fiber.Ctx) error {
	// Пробуем извлечь JWT из заголовка
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == bearerScheme {
			claims, err := jwt.ValidateToken(parts[1], m.config.JWTSecret)
			if err == nil {
				c.Locals("user_id", claims.UserID)
				c.Locals("user_email", claims.Email)
				c.Locals("auth_method", "jwt")

				// Проверяем админ статус
				isAdmin, err := m.services.User().IsUserAdmin(c.Context(), claims.Email)
				if err != nil {
					isAdmin = false
				}
				c.Locals("is_admin", isAdmin)
			}
		}
	}

	// Если нет JWT в заголовке, проверяем cookie
	if c.Locals("user_id") == nil {
		jwtCookie := c.Cookies("jwt_token")
		if jwtCookie != "" {
			claims, err := jwt.ValidateToken(jwtCookie, m.config.JWTSecret)
			if err == nil {
				c.Locals("user_id", claims.UserID)
				c.Locals("user_email", claims.Email)
				c.Locals("auth_method", "jwt_cookie")

				// Проверяем админ статус
				isAdmin, err := m.services.User().IsUserAdmin(c.Context(), claims.Email)
				if err != nil {
					isAdmin = false
				}
				c.Locals("is_admin", isAdmin)
			}
		}
	}

	return c.Next()
}

// RefreshJWT - обновление JWT токена
func (m *Middleware) RefreshJWT(c *fiber.Ctx) error {
	// Получаем текущий токен
	var currentToken string

	// Проверяем заголовок
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == bearerScheme {
			currentToken = parts[1]
		}
	}

	// Или из cookie
	if currentToken == "" {
		currentToken = c.Cookies("jwt_token")
	}

	if currentToken == "" {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.token_required")
	}

	// Валидируем текущий токен
	claims, err := jwt.ValidateToken(currentToken, m.config.JWTSecret)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.invalid_token")
	}

	// Генерируем новый токен
	newToken, err := jwt.GenerateTokenWithDuration(
		claims.UserID,
		claims.Email,
		m.config.JWTSecret,
		time.Duration(m.config.JWTExpirationHours)*time.Hour,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "users.auth.error.token_generation_failed")
	}

	// Устанавливаем новый токен в cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    newToken,
		Path:     "/",
		Domain:   m.config.GetCookieDomain(),
		MaxAge:   m.config.JWTExpirationHours * 3600,
		Secure:   m.config.GetCookieSecure(),
		HTTPOnly: true,
		SameSite: m.config.GetCookieSameSite(),
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"token":      newToken,
			"expires_in": m.config.JWTExpirationHours * 3600,
		},
	})
}
