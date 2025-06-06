// backend/internal/middleware/auth_jwt.go
package middleware

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend/pkg/jwt"
	"backend/pkg/utils"
)

// AuthRequiredJWT - основной метод аутентификации через JWT
// Поддерживает как Bearer токены в заголовке, так и fallback на session cookies
func (m *Middleware) AuthRequiredJWT(c *fiber.Ctx) error {
	log.Printf("AuthRequiredJWT middleware called for path: %s", c.Path())

	// Приоритет 1: Проверяем JWT токен в заголовке Authorization
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		// Извлекаем токен из заголовка "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString := parts[1]

			// Валидируем JWT токен
			claims, err := jwt.ValidateToken(tokenString, m.config.JWTSecret)
			if err != nil {
				log.Printf("AuthRequiredJWT: JWT validation failed: %v", err)
				return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.invalid_token")
			}

			// Проверяем что пользователь существует
			user, err := m.services.User().GetUserByID(c.Context(), claims.UserID)
			if err != nil || user == nil {
				log.Printf("AuthRequiredJWT: User not found for ID %d", claims.UserID)
				return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.user_not_found")
			}

			// Устанавливаем данные пользователя в контекст
			c.Locals("user_id", claims.UserID)
			c.Locals("user_email", claims.Email)
			c.Locals("auth_method", "jwt")
			c.Locals("jwt_token", tokenString)

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
		
		claims, err := jwt.ValidateToken(tokenFromQuery, m.config.JWTSecret)
		if err != nil {
			log.Printf("AuthRequiredJWT: JWT query validation failed: %v", err)
		} else {
			// JWT из query валиден
			c.Locals("user_id", claims.UserID)
			c.Locals("user_email", claims.Email)
			c.Locals("auth_method", "jwt_query")
			c.Locals("jwt_token", tokenFromQuery)
			
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
			log.Printf("AuthRequiredJWT: JWT cookie validation failed: %v", err)
			// Очищаем невалидную cookie
			c.Cookie(&fiber.Cookie{
				Name:     "jwt_token",
				Value:    "",
				Path:     "/",
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
		log.Printf("AuthRequiredJWT: Falling back to session authentication")
		
		session, err := m.services.Auth().GetSession(c.Context(), sessionToken)
		if err != nil {
			log.Printf("AuthRequiredJWT: Session validation failed: %v", err)
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
			log.Printf("AuthRequiredJWT: Failed to generate JWT from session: %v", err)
		} else {
			// Устанавливаем JWT cookie для будущих запросов
			c.Cookie(&fiber.Cookie{
				Name:     "jwt_token",
				Value:    newJWT,
				Path:     "/",
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

		go func() {
			ctx := context.Background()
			_ = m.services.User().UpdateLastSeen(ctx, session.UserID)
		}()

		return c.Next()
	}

	// Если ни один метод аутентификации не сработал
	if c.Get("Upgrade") == "websocket" {
		log.Printf("SECURITY: Unauthorized WebSocket connection attempt from IP: %s, User-Agent: %s", 
			c.IP(), c.Get("User-Agent"))
	}
	
	return utils.ErrorResponse(c, fiber.StatusUnauthorized, "users.auth.error.authentication_required")
}

// OptionalAuthJWT - опциональная аутентификация, не требует токена но извлекает данные если есть
func (m *Middleware) OptionalAuthJWT(c *fiber.Ctx) error {
	// Пробуем извлечь JWT из заголовка
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := jwt.ValidateToken(parts[1], m.config.JWTSecret)
			if err == nil {
				c.Locals("user_id", claims.UserID)
				c.Locals("user_email", claims.Email)
				c.Locals("auth_method", "jwt")
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
		if len(parts) == 2 && parts[0] == "Bearer" {
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
		MaxAge:   m.config.JWTExpirationHours * 3600,
		Secure:   m.config.GetCookieSecure(),
		HTTPOnly: true,
		SameSite: m.config.GetCookieSameSite(),
	})
	
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"token": newToken,
			"expires_in": m.config.JWTExpirationHours * 3600,
		},
	})
}