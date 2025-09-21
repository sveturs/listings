package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"backend/internal/logger"
	"backend/internal/service/authclient"

	"github.com/gofiber/fiber/v2"
)

// ServiceProvider предоставляет доступ к сервисам
// Используем интерфейс для избежания циклических зависимостей
type ServiceProvider interface {
	User() interface {
		IsUserAdmin(ctx context.Context, email string) (bool, error)
	}
}

type AuthProxyMiddleware struct {
	authClient *authclient.Client
	httpClient *http.Client
	enabled    bool
	baseURL    string
	services   ServiceProvider
}

func NewAuthProxyMiddleware(services ServiceProvider) *AuthProxyMiddleware {
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:28080"
	}

	// Проверяем переменную USE_AUTH_SERVICE
	useAuthService := os.Getenv("USE_AUTH_SERVICE")
	enabled := useAuthService == "true"

	return &AuthProxyMiddleware{
		authClient: authclient.NewClient(authServiceURL),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Не следовать за редиректами автоматически
				return http.ErrUseLastResponse
			},
		},
		enabled:  enabled,
		baseURL:  authServiceURL,
		services: services,
	}
}

func (m *AuthProxyMiddleware) ProxyToAuthService() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		// Проверяем, относится ли запрос к Auth Service
		// Проксируем ВСЕ auth запросы включая OAuth
		if !strings.HasPrefix(path, "/api/v1/auth/") && !strings.HasPrefix(path, "/auth/") {
			return c.Next()
		}

		// Если Auth Service отключен, пропускаем
		if !m.enabled {
			return c.Next()
		}

		// Логируем все проксируемые запросы для отладки
		queryString := string(c.Request().URI().QueryString())
		logger.Info().
			Str("path", path).
			Str("method", c.Method()).
			Str("queryString", queryString).
			Str("baseURL", m.baseURL).
			Msg("[AUTH_PROXY] === START PROXYING REQUEST TO AUTH SERVICE ===")

		// Создаем новый HTTP запрос к Auth Service
		targetURL := m.baseURL + path

		// Добавляем query parameters
		if c.Request().URI().QueryString() != nil {
			targetURL += "?" + string(c.Request().URI().QueryString())
		}

		logger.Info().
			Str("targetURL", targetURL).
			Int("targetURLLength", len(targetURL)).
			Msg("[AUTH_PROXY] Target URL constructed")

		// Получаем тело запроса
		body := c.Body()

		// Создаем HTTP запрос
		req, err := http.NewRequestWithContext(c.Context(), c.Method(), targetURL, bytes.NewReader(body))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create proxy request",
			})
		}

		// Копируем заголовки, включая cookies, но исключаем несовместимые с HTTP/2
		// HTTP/2 не поддерживает следующие заголовки: Connection, Upgrade (кроме WebSocket),
		// Keep-Alive, Proxy-Connection, Transfer-Encoding, TE
		c.Request().Header.VisitAll(func(key, value []byte) {
			headerName := string(key)
			headerNameLower := strings.ToLower(headerName)

			// Пропускаем заголовки, несовместимые с HTTP/2
			if headerNameLower == "connection" ||
				headerNameLower == "keep-alive" ||
				headerNameLower == "proxy-connection" ||
				headerNameLower == "transfer-encoding" ||
				headerNameLower == "te" ||
				// Upgrade пропускаем только если это не WebSocket запрос
				(headerNameLower == "upgrade" && !strings.Contains(strings.ToLower(string(value)), "websocket")) {
				return
			}

			req.Header.Set(headerName, string(value))
		})

		// Копируем cookies из заголовка Cookie если он есть
		if cookieHeader := c.Get("Cookie"); cookieHeader != "" {
			req.Header.Set("Cookie", cookieHeader)
		}

		// Добавляем заголовки проксирования
		req.Header.Set("X-Forwarded-For", c.IP())
		req.Header.Set("X-Forwarded-Host", c.Hostname())
		req.Header.Set("X-Forwarded-Proto", c.Protocol())

		// Добавляем информацию о frontend URL для правильного редиректа после OAuth
		// Сначала проверяем query параметр returnTo
		if returnTo := c.Query("returnTo"); returnTo != "" {
			req.Header.Set("X-Frontend-URL", returnTo)
		} else {
			// Если returnTo не указан, используем FRONTEND_URL из конфигурации
			frontendURL := os.Getenv("FRONTEND_URL")
			if frontendURL != "" {
				req.Header.Set("X-Frontend-URL", frontendURL)
			}
		}

		// Выполняем запрос
		resp, err := m.httpClient.Do(req)
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error": fmt.Sprintf("auth service request failed: %v", err),
			})
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		// Читаем ответ
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to read auth service response",
			})
		}

		// Копируем заголовки ответа, но НЕ перезаписываем CORS заголовки
		for key, values := range resp.Header {
			keyLower := strings.ToLower(key)

			// Пропускаем CORS заголовки - они должны устанавливаться основным приложением
			if strings.HasPrefix(keyLower, "access-control-") {
				continue
			}

			// Специальная обработка для Set-Cookie
			if keyLower == "set-cookie" {
				// Set-Cookie заголовки должны обрабатываться отдельно для каждого cookie
				for _, cookie := range values {
					c.Response().Header.Add("Set-Cookie", cookie)
				}
			} else {
				// Для остальных заголовков
				for _, value := range values {
					c.Set(key, value)
				}
			}
		}

		// Для OAuth редиректов - возвращаем Location напрямую
		if resp.StatusCode == 302 || resp.StatusCode == 301 || resp.StatusCode == 303 || resp.StatusCode == 307 || resp.StatusCode == 308 {
			if location := resp.Header.Get("Location"); location != "" {
				logger.Info().
					Int("statusCode", resp.StatusCode).
					Str("location", location).
					Int("locationLength", len(location)).
					Msg("[AUTH_PROXY] Redirect response from Auth Service")

				// Возвращаем редирект напрямую для браузера
				// Fiber.Redirect должен работать для любых URL
				return c.Status(resp.StatusCode).Redirect(location)
			}
		}

		// Если это запрос /auth/session, добавляем поле is_admin
		if path == "/api/v1/auth/session" && resp.StatusCode == 200 && m.services != nil {
			var sessionData map[string]interface{}
			if err := json.Unmarshal(respBody, &sessionData); err == nil {
				// Получаем данные пользователя из ответа
				if userData, ok := sessionData["user"].(map[string]interface{}); ok {
					// Проверяем email пользователя
					if email, ok := userData["email"].(string); ok && email != "" {
						// Проверяем, является ли пользователь администратором в локальной БД
						ctx := c.Context()
						isAdmin, err := m.services.User().IsUserAdmin(ctx, email)
						if err == nil {
							userData["is_admin"] = isAdmin

							logger.Info().
								Str("email", email).
								Bool("is_admin", isAdmin).
								Msg("[AUTH_PROXY] Added is_admin field from local DB")
						}
					}
				}

				// Перекодируем ответ обратно в JSON
				modifiedBody, _ := json.Marshal(sessionData)
				respBody = modifiedBody
			}
		}

		logger.Info().
			Int("statusCode", resp.StatusCode).
			Int("responseBodyLength", len(respBody)).
			Msg("[AUTH_PROXY] === END PROXYING REQUEST TO AUTH SERVICE ===")

		// Возвращаем ответ
		c.Status(resp.StatusCode)
		return c.Send(respBody)
	}
}

func (m *AuthProxyMiddleware) ValidateTokenWithAuthService() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Always validate tokens through Auth Service
		// Извлекаем токен из заголовка Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Убираем префикс "Bearer "
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		// Временное логирование токена для отладки - ПОЛНЫЙ ТОКЕН, НЕ БОИМСЯ!
		logger.Info().
			Str("path", c.Path()).
			Int("token_length", len(token)).
			Str("full_token", token).
			Msg("[AuthProxy] Validating token")

		// Валидируем токен через Auth Service
		validateResp, err := m.authClient.ValidateToken(c.Context(), token)
		if err != nil {
			logger.Error().
				Err(err).
				Str("path", c.Path()).
				Str("token", token).
				Msg("[AuthProxy] Token validation failed")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		logger.Info().
			Int("userID", int(validateResp.UserID)).
			Str("email", validateResp.Email).
			Bool("valid", validateResp.Valid).
			Msg("[AuthProxy] Token validated successfully")

		if !validateResp.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token is not valid",
			})
		}

		// Сохраняем информацию о пользователе в контексте
		c.Locals("userID", validateResp.UserID)
		c.Locals("userEmail", validateResp.Email)
		c.Locals("tokenExpiresAt", validateResp.ExpiresAt)

		return c.Next()
	}
}
