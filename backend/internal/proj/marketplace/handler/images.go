// backend/internal/proj/marketplace/handler/images.go
package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ImagesHandler обрабатывает запросы, связанные с изображениями объявлений
type ImagesHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

// NewImagesHandler создает новый обработчик изображений
func NewImagesHandler(services globalService.ServicesInterface) *ImagesHandler {
	return &ImagesHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

// UploadImages загружает изображения для объявления
func (h *ImagesHandler) UploadImages(c *fiber.Ctx) error {
	// Добавим явные логи для отладки
	log.Printf("НАЧАЛО ЗАГРУЗКИ ИЗОБРАЖЕНИЙ. Получен запрос: %s %s", c.Method(), c.Path())

	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Ошибка получения user_id из контекста: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}
	log.Printf("Аутентифицированный пользователь ID: %d", userID)

	// Проверяем, пришли ли какие-то файлы
	form, err := c.MultipartForm()
	if err != nil {
		log.Printf("Ошибка получения MultipartForm: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Ошибка получения файлов")
	}

	files := form.File["file"]
	if len(files) == 0 {
		// Попробуем альтернативное имя поля
		files = form.File["files"]
		if len(files) == 0 {
			// Проверим все поля
			log.Printf("Поиск файлов в form.File с ключами: %v", getMapKeys(form.File))
			for key, values := range form.File {
				log.Printf("Поле %s содержит %d файлов", key, len(values))
				if len(values) > 0 {
					files = values
					break
				}
			}

			if len(files) == 0 {
				log.Printf("Не найдено файлов в запросе")
				return utils.ErrorResponse(c, fiber.StatusBadRequest, "Не найдено файлов для загрузки")
			}
		}
	}

	log.Printf("Найдено %d файлов для загрузки", len(files))

	// Получаем ID объявления из параметров
	listingIDStr := c.FormValue("listing_id")
	if listingIDStr == "" {
		// Попробуем из параметров URL
		listingIDStr = c.Params("id")
		if listingIDStr == "" {
			log.Printf("Ошибка: ID объявления не указан")
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Не указан ID объявления")
		}
	}

	listingID, err := strconv.Atoi(listingIDStr)
	if err != nil {
		log.Printf("Ошибка преобразования ID объявления '%s': %v", listingIDStr, err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
	}
	log.Printf("ID объявления для загрузки изображений: %d", listingID)

	// Получаем информацию об объявлении для проверки владельца
	listing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		log.Printf("Ошибка получения объявления с ID %d: %v", listingID, err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Объявление не найдено")
	}

	// Проверяем, владеет ли пользователь объявлением
	if listing.UserID != userID {
		log.Printf("Доступ запрещен: UserID %d не владеет объявлением %d (владелец: %d)",
			userID, listingID, listing.UserID)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Вы не являетесь владельцем этого объявления")
	}

	// Определяем главное изображение
	mainImageIndex := 0
	if mainImageIndexStr := c.FormValue("main_image_index"); mainImageIndexStr != "" {
		if idx, err := strconv.Atoi(mainImageIndexStr); err == nil {
			mainImageIndex = idx
			log.Printf("Установлен индекс главного изображения: %d", mainImageIndex)
		}
	}

	var uploadedImages []models.MarketplaceImage
	for i, file := range files {
		log.Printf("Обработка файла %d: %s (размер: %d, тип: %s)",
			i, file.Filename, file.Size, file.Header.Get("Content-Type"))

		// Проверка типа файла и размера
		if file.Size > 10*1024*1024 {
			log.Printf("Файл слишком большой: %d байт", file.Size)
			continue // Пропускаем слишком большие файлы
		}

		contentType := file.Header.Get("Content-Type")
		if !strings.HasPrefix(contentType, "image/") {
			log.Printf("Неподдерживаемый тип файла: %s", contentType)
			continue // Пропускаем файлы не-изображения
		}

		// Загружаем изображение
		isMain := i == mainImageIndex
		log.Printf("Загрузка изображения: listingID=%d, isMain=%v", listingID, isMain)
		image, err := h.marketplaceService.UploadImage(c.Context(), file, listingID, isMain)
		if err != nil {
			log.Printf("Ошибка загрузки изображения: %v", err)
			continue
		}

		log.Printf("Image successfully uploaded: ID=%d, Path=%s", image.ID, image.FilePath)

		uploadedImages = append(uploadedImages, *image)
	}

	// Переиндексируем объявление с загруженными изображениями
	log.Printf("Переиндексация объявления %d с новыми изображениями", listingID)
	fullListing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		log.Printf("Ошибка получения полной информации об объявлении для переиндексации: %v", err)
	} else {
		if err := h.marketplaceService.Storage().IndexListing(c.Context(), fullListing); err != nil {
			log.Printf("Ошибка переиндексации объявления после загрузки изображений: %v", err)
		} else {
			log.Printf("Успешно переиндексировано объявление %d с %d изображениями",
				listingID, len(fullListing.Images))
		}
	}

	log.Printf("ЗАВЕРШЕНА ЗАГРУЗКА ИЗОБРАЖЕНИЙ: загружено %d из %d файлов",
		len(uploadedImages), len(files))

	// Возвращаем успешный результат с информацией о загруженных изображениях
	return utils.SuccessResponse(c, fiber.Map{
		"message": "Изображения успешно загружены",
		"images":  uploadedImages,
		"count":   len(uploadedImages),
	})
}

// Вспомогательная функция для получения ключей из map
func getMapKeys(m map[string][]*multipart.FileHeader) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// DeleteImage удаляет изображение
func (h *ImagesHandler) DeleteImage(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем ID изображения из параметров URL
	imageID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID изображения")
	}

	// Получаем информацию об изображении
	image, err := h.services.Storage().GetListingImageByID(c.Context(), imageID)
	if err != nil {
		log.Printf("Failed to get image with ID %d: %v", imageID, err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Изображение не найдено")
	}

	// Получаем информацию об объявлении
	listing, err := h.marketplaceService.GetListingByID(c.Context(), image.ListingID)
	if err != nil {
		log.Printf("Failed to get listing with ID %d: %v", image.ListingID, err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Объявление не найдено")
	}

	// Проверяем владельца объявления
	if listing.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Вы не являетесь владельцем этого объявления")
	}

	// Удаляем изображение
	err = h.marketplaceService.DeleteImage(c.Context(), imageID)
	if err != nil {
		log.Printf("Failed to delete image with ID %d: %v", imageID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось удалить изображение")
	}

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, fiber.Map{
		"message": "Изображение успешно удалено",
	})
}

// ModerateImage проверяет изображение на запрещенный контент
func (h *ImagesHandler) ModerateImage(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем пользователя по ID для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user with ID %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось проверить права администратора")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		log.Printf("User %d is not admin: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Требуются права администратора")
	}

	// Получаем файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Failed to get file from request: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Ошибка загрузки файла")
	}

	// Проверяем размер файла (ограничение 10MB)
	if file.Size > 10*1024*1024 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Размер файла превышает 10MB")
	}

	// Проверяем тип файла
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неподдерживаемый тип файла. Разрешены только изображения.")
	}

	// Создаем временный файл для сохранения загруженного изображения
	tempDir := os.TempDir()
	tempFilePath := filepath.Join(tempDir, fmt.Sprintf("moderate_image_%d", time.Now().UnixNano()))
	if err := c.SaveFile(file, tempFilePath); err != nil {
		log.Printf("Failed to save file: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось сохранить файл")
	}
	defer os.Remove(tempFilePath) // Удаляем временный файл после завершения

	// Возвращаем заглушку для модерации изображения (полная реализация требует Vision API)
	// TODO: Реализовать проверку с использованием Vision API

	// Возвращаем результаты проверки
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"labels":            []string{"image", "photo"},
			"prohibited_labels": []string{},
			"has_prohibited":    false,
		},
	})
}

// EnhancePreview создает предпросмотр улучшенного изображения
func (h *ImagesHandler) EnhancePreview(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем ID изображения из параметров
	imageID, err := strconv.Atoi(c.FormValue("image_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID изображения")
	}

	// Получаем информацию об изображении
	image, err := h.services.Storage().GetListingImageByID(c.Context(), imageID)
	if err != nil {
		log.Printf("Failed to get image with ID %d: %v", imageID, err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Изображение не найдено")
	}

	// Получаем информацию об объявлении
	listing, err := h.marketplaceService.GetListingByID(c.Context(), image.ListingID)
	if err != nil {
		log.Printf("Failed to get listing with ID %d: %v", image.ListingID, err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Объявление не найдено")
	}

	// Проверяем владельца объявления
	if listing.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Вы не являетесь владельцем этого объявления")
	}

	// Получаем тип улучшения
	enhancementType := c.FormValue("enhancement_type")
	if enhancementType == "" {
		enhancementType = "quality" // По умолчанию улучшаем качество
	}

	// Создаем предпросмотр улучшенного изображения
	// Здесь должен быть код для обработки изображения и создания предпросмотра
	// Пока возвращаем заглушку
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"preview_url": fmt.Sprintf("https://example.com/preview/%d/%s", imageID, enhancementType),
		},
	})
}

// EnhanceImages улучшает изображения объявления
func (h *ImagesHandler) EnhanceImages(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Printf("Failed to get user_id from context: %v", c.Locals("user_id"))
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем ID объявления из параметров
	listingID, err := strconv.Atoi(c.FormValue("listing_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
	}

	// Получаем информацию об объявлении
	listing, err := h.marketplaceService.GetListingByID(c.Context(), listingID)
	if err != nil {
		log.Printf("Failed to get listing with ID %d: %v", listingID, err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Объявление не найдено")
	}

	// Проверяем владельца объявления
	if listing.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Вы не являетесь владельцем этого объявления")
	}

	// Получаем список ID изображений для улучшения
	var imageIDs []int
	if err := c.BodyParser(&imageIDs); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат запроса")
	}

	// Проверяем, что указанные изображения действительно принадлежат данному объявлению
	for _, imageID := range imageIDs {
		image, err := h.services.Storage().GetListingImageByID(c.Context(), imageID)
		if err != nil {
			log.Printf("Failed to get image with ID %d: %v", imageID, err)
			return utils.ErrorResponse(c, fiber.StatusNotFound, fmt.Sprintf("Изображение %d не найдено", imageID))
		}

		if image.ListingID != listingID {
			return utils.ErrorResponse(c, fiber.StatusForbidden, fmt.Sprintf("Изображение %d не принадлежит указанному объявлению", imageID))
		}
	}

	// Получаем тип улучшения
	enhancementType := c.FormValue("enhancement_type")
	if enhancementType == "" {
		enhancementType = "quality" // По умолчанию улучшаем качество
	}

	// Запускаем процесс улучшения изображений
	// Здесь должен быть код для обработки изображений
	// Пока возвращаем заглушку
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": fmt.Sprintf("Улучшение изображений для объявления %d запущено", listingID),
			"job_id":  fmt.Sprintf("enhance_%d_%d", listingID, time.Now().Unix()),
		},
	})
}
