package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"log"
)

type StorefrontHandler struct {
	services globalService.ServicesInterface
}

func NewStorefrontHandler(services globalService.ServicesInterface) *StorefrontHandler {
	return &StorefrontHandler{
		services: services,
	}
}

// CreateStorefront создаёт новую витрину
func (h *StorefrontHandler) CreateStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var create models.StorefrontCreate
	if err := c.BodyParser(&create); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format")
	}

	// Валидация
	if create.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Name is required")
	}

	storefront, err := h.services.Storefront().CreateStorefront(c.Context(), userID, &create)
	if err != nil {
		if err.Error() == "insufficient funds" {
			return utils.ErrorResponse(c, fiber.StatusPaymentRequired, "Недостаточно средств для создания витрины")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to create storefront: %s", err.Error()))
	}

	return utils.SuccessResponse(c, storefront)
}

// GetUserStorefronts возвращает все витрины пользователя
func (h *StorefrontHandler) GetUserStorefronts(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	storefronts, err := h.services.Storefront().GetUserStorefronts(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get storefronts")
	}

	return utils.SuccessResponse(c, storefronts)
}

// GetStorefront возвращает витрину по ID
func (h *StorefrontHandler) GetStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid storefront ID")
	}

	storefront, err := h.services.Storefront().GetStorefrontByID(c.Context(), id, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get storefront")
	}

	return utils.SuccessResponse(c, storefront)
}

// UpdateStorefront обновляет информацию о витрине
func (h *StorefrontHandler) UpdateStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid storefront ID")
	}

	var storefront models.Storefront
	if err := c.BodyParser(&storefront); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format")
	}
	storefront.ID = id

	err = h.services.Storefront().UpdateStorefront(c.Context(), &storefront, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update storefront")
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Storefront updated successfully"})
}

// DeleteStorefront удаляет витрину
func (h *StorefrontHandler) DeleteStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid storefront ID")
	}

	err = h.services.Storefront().DeleteStorefront(c.Context(), id, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete storefront")
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Storefront deleted successfully"})
}

// CreateImportSource создаёт новый источник импорта
func (h *StorefrontHandler) CreateImportSource(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var create models.ImportSourceCreate
	if err := c.BodyParser(&create); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format")
	}

	source, err := h.services.Storefront().CreateImportSource(c.Context(), &create, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create import source")
	}

	return utils.SuccessResponse(c, source)
}

// GetImportSources возвращает источники импорта для витрины
func (h *StorefrontHandler) GetImportSources(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid storefront ID")
	}

	sources, err := h.services.Storefront().GetImportSources(c.Context(), storefrontID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get import sources")
	}

	return utils.SuccessResponse(c, sources)
}

// UpdateImportSource обновляет источник импорта
func (h *StorefrontHandler) UpdateImportSource(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid import source ID")
	}

	var source models.ImportSource
	if err := c.BodyParser(&source); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request format")
	}
	source.ID = id

	err = h.services.Storefront().UpdateImportSource(c.Context(), &source, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update import source")
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Import source updated successfully"})
}

// DeleteImportSource удаляет источник импорта
func (h *StorefrontHandler) DeleteImportSource(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid import source ID")
	}

	err = h.services.Storefront().DeleteImportSource(c.Context(), id, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete import source")
	}

	return utils.SuccessResponse(c, fiber.Map{"message": "Import source deleted successfully"})
}

// RunImport запускает импорт данных
func (h *StorefrontHandler) RunImport(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    sourceID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid import source ID")
    }

    // Отладочный лог
    log.Printf("Running import for source ID: %d, user ID: %d", sourceID, userID)

    // Сначала проверим, существует ли источник импорта
    source, err := h.services.Storefront().GetImportSourceByID(c.Context(), sourceID, userID)
    if err != nil {
        log.Printf("Error fetching source %d: %v", sourceID, err)
        return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Import source not found or access denied: %v", err))
    }

    // Отладочный лог
    log.Printf("Found import source: %+v", source)

    // Получение файла из формы, если он есть
    file, err := c.FormFile("file")
    if err != nil && err != fiber.ErrUnprocessableEntity {
        log.Printf("Error processing form file: %v", err)
        return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Error processing file: %v", err))
    }

    var history *models.ImportHistory

    if file != nil {
        // Отладочный лог
        log.Printf("Processing uploaded file: %s, size: %d", file.Filename, file.Size)
        
        // Обработка загруженного файла
        fileHandle, err := file.Open()
        if err != nil {
            log.Printf("Error opening file: %v", err)
            return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Error opening file: %v", err))
        }
        defer fileHandle.Close()

        history, err = h.services.Storefront().ImportCSV(c.Context(), sourceID, fileHandle, userID)
        if err != nil {
            log.Printf("Error importing CSV: %v", err)
            return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to import CSV: %v", err))
        }
    } else {
        // Проверяем, есть ли URL для импорта
        if source.URL == "" {
            return utils.ErrorResponse(c, fiber.StatusBadRequest, "No file uploaded and no URL configured for import")
        }
        
        // Запуск импорта по URL
        history, err = h.services.Storefront().RunImport(c.Context(), sourceID, userID)
        if err != nil {
            log.Printf("Error running import: %v", err)
            return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to run import: %v", err))
        }
    }

    // Отладочный лог
    log.Printf("Import completed successfully, history: %+v", history)

    return utils.SuccessResponse(c, history)
}

// GetImportHistory возвращает историю импорта
func (h *StorefrontHandler) GetImportHistory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	sourceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid import source ID")
	}

	limit := 10
	offset := 0
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	history, err := h.services.Storefront().GetImportHistory(c.Context(), sourceID, userID, limit, offset)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get import history")
	}

	return utils.SuccessResponse(c, history)
}