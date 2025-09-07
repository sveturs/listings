package server

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ProxyMinIO проксирует запросы к MinIO для статических файлов
func (s *Server) ProxyMinIO(c *fiber.Ctx) error {
	// Получаем путь после /listings/
	path := c.Params("*")
	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid path",
		})
	}

	// Формируем URL для MinIO
	minioURL := fmt.Sprintf("http://%s/listings/%s", s.cfg.FileStorage.MinioEndpoint, path)

	log.Debug().
		Str("path", path).
		Str("minio_url", minioURL).
		Msg("Proxying request to MinIO")

	// Создаем HTTP запрос к MinIO
	req, err := http.NewRequest("GET", minioURL, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create MinIO request")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create request",
		})
	}

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch from MinIO")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch file",
		})
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Warn().
			Int("status", resp.StatusCode).
			Str("path", path).
			Msg("MinIO returned non-OK status")
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Копируем заголовки
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		c.Set("Content-Type", contentType)
	}

	contentLength := resp.Header.Get("Content-Length")
	if contentLength != "" {
		c.Set("Content-Length", contentLength)
	}

	// Устанавливаем заголовки кеширования
	c.Set("Cache-Control", "public, max-age=604800") // 7 дней

	// Определяем тип контента по расширению файла, если Content-Type не установлен
	if contentType == "" {
		ext := strings.ToLower(path[strings.LastIndex(path, ".")+1:])
		switch ext {
		case "jpg", "jpeg":
			c.Set("Content-Type", "image/jpeg")
		case "png":
			c.Set("Content-Type", "image/png")
		case "gif":
			c.Set("Content-Type", "image/gif")
		case "webp":
			c.Set("Content-Type", "image/webp")
		case "pdf":
			c.Set("Content-Type", "application/pdf")
		default:
			c.Set("Content-Type", "application/octet-stream")
		}
	}

	// Копируем тело ответа
	_, err = io.Copy(c.Response().BodyWriter(), resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to copy response body")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send file",
		})
	}

	return nil
}

// ProxyChatFiles проксирует запросы к MinIO для файлов чата
func (s *Server) ProxyChatFiles(c *fiber.Ctx) error {
	// Получаем путь после /chat-files/
	path := c.Params("*")
	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid path",
		})
	}

	// Формируем URL для MinIO используя конфигурируемое имя bucket'а
	bucketName := s.cfg.FileStorage.MinioChatBucket
	if bucketName == "" {
		bucketName = "chat-files" // Fallback на дефолтное значение
	}
	minioURL := fmt.Sprintf("http://%s/%s/%s", s.cfg.FileStorage.MinioEndpoint, bucketName, path)

	log.Debug().
		Str("path", path).
		Str("minio_url", minioURL).
		Msg("Proxying chat file request to MinIO")

	// Используем тот же механизм проксирования
	req, err := http.NewRequest("GET", minioURL, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create request",
		})
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch file",
		})
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Копируем заголовки и тело
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		c.Set("Content-Type", ct)
	}
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		c.Set("Content-Length", cl)
	}
	c.Set("Cache-Control", "public, max-age=604800")

	_, err = io.Copy(c.Response().BodyWriter(), resp.Body)
	return err
}

// ProxyStorefrontProducts проксирует запросы к MinIO для изображений товаров витрин
func (s *Server) ProxyStorefrontProducts(c *fiber.Ctx) error {
	// Получаем путь после /storefront-products/
	path := c.Params("*")
	if path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid path",
		})
	}

	// Формируем URL для MinIO используя конфигурируемое имя bucket'а
	bucketName := s.cfg.FileStorage.MinioStorefrontBucket
	if bucketName == "" {
		bucketName = "storefront-products" // Fallback на дефолтное значение
	}
	minioURL := fmt.Sprintf("http://%s/%s/%s", s.cfg.FileStorage.MinioEndpoint, bucketName, path)

	log.Debug().
		Str("path", path).
		Str("minio_url", minioURL).
		Msg("Proxying storefront product image request to MinIO")

	// Используем тот же механизм проксирования
	req, err := http.NewRequest("GET", minioURL, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create request",
		})
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch file",
		})
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error": "File not found",
		})
	}

	// Копируем заголовки и тело
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		c.Set("Content-Type", ct)
	}
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		c.Set("Content-Length", cl)
	}
	c.Set("Cache-Control", "public, max-age=604800") // 7 дней

	// Определяем тип контента по расширению файла, если Content-Type не установлен
	if resp.Header.Get("Content-Type") == "" {
		ext := strings.ToLower(path[strings.LastIndex(path, ".")+1:])
		switch ext {
		case "jpg", "jpeg":
			c.Set("Content-Type", "image/jpeg")
		case "png":
			c.Set("Content-Type", "image/png")
		case "gif":
			c.Set("Content-Type", "image/gif")
		case "webp":
			c.Set("Content-Type", "image/webp")
		default:
			c.Set("Content-Type", "application/octet-stream")
		}
	}

	_, err = io.Copy(c.Response().BodyWriter(), resp.Body)
	return err
}
