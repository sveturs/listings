// backend/internal/proj/storefront/handler/storefront.go
package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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
// GetCategoryMappings возвращает сопоставления категорий для источника импорта
func (h *StorefrontHandler) GetCategoryMappings(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    sourceID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid import source ID")
    }
    
    mappings, err := h.services.Storefront().GetCategoryMappings(c.Context(), sourceID, userID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get category mappings")
    }
    
    return utils.SuccessResponse(c, mappings)
}

// UpdateCategoryMappings обновляет сопоставления категорий для источника импорта
func (h *StorefrontHandler) UpdateCategoryMappings(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    sourceID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid import source ID")
    }
    
    var mappings map[string]int
    if err := c.BodyParser(&mappings); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid mappings format")
    }
    
    err = h.services.Storefront().UpdateCategoryMappings(c.Context(), sourceID, userID, mappings)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update category mappings")
    }
    
    return utils.SuccessResponse(c, fiber.Map{
        "message": "Category mappings updated successfully",
    })
}
// GetImportedCategories возвращает список категорий, которые были импортированы этим источником
func (h *StorefrontHandler) GetImportedCategories(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    sourceID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid import source ID")
    }
    
    categories, err := h.services.Storefront().GetImportedCategories(c.Context(), sourceID, userID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get imported categories")
    }
    
    return utils.SuccessResponse(c, categories)
}

// ApplyCategoryMappings применяет сопоставления категорий к импортированным товарам
func (h *StorefrontHandler) ApplyCategoryMappings(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    sourceID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid import source ID")
    }
    
    updatedCount, err := h.services.Storefront().ApplyCategoryMappings(c.Context(), sourceID, userID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to apply category mappings: %s", err.Error()))
    }
    
    return utils.SuccessResponse(c, fiber.Map{
        "message": fmt.Sprintf("Successfully updated categories for %d listings", updatedCount),
        "updated_count": updatedCount,
    })
}

// GetPublicStorefront возвращает публичные данные витрины по ID
func (h *StorefrontHandler) GetPublicStorefront(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid storefront ID")
	}

	storefront, err := h.services.Storefront().GetPublicStorefrontByID(c.Context(), id)
	if err != nil {
		// Проверяем различные типы ошибок
		if err.Error() == "storefront not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Storefront not found")
		}
		if err.Error() == "storefront is not active" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Storefront is not active")
		}

		// Логируем детали ошибки для отладки
		fmt.Printf("Error getting public storefront %d: %v\n", id, err)

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

	var history *models.ImportHistory
	var csvFile, zipFile multipart.File
	var xmlZipFile multipart.File

	// Проверяем тип контента запроса
	contentType := c.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		// Это запрос с файлами
		form, err := c.MultipartForm()
		if err != nil && err != fiber.ErrUnprocessableEntity {
			log.Printf("Error processing form file: %v", err)
			return utils.ErrorResponse(c, fiber.StatusBadRequest, fmt.Sprintf("Error processing file: %v", err))
		}

		// Обработка файлов из формы если они есть
		if form != nil {
			// Проверяем файл CSV
			csvFiles := form.File["file"]
			if len(csvFiles) > 0 {
				csvFileHeader := csvFiles[0]
				log.Printf("Processing uploaded CSV file: %s, size: %d", csvFileHeader.Filename, csvFileHeader.Size)

				var err error
				csvFile, err = csvFileHeader.Open()
				if err != nil {
					log.Printf("Error opening CSV file: %v", err)
					return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Error opening CSV file: %v", err))
				}
				defer csvFile.Close()
			}

			// Проверяем ZIP файл с изображениями
			zipFiles := form.File["images_zip"]
			if len(zipFiles) > 0 {
				zipFileHeader := zipFiles[0]
				log.Printf("Processing uploaded ZIP file: %s, size: %d", zipFileHeader.Filename, zipFileHeader.Size)

				var err error
				zipFile, err = zipFileHeader.Open()
				if err != nil {
					log.Printf("Error opening ZIP file: %v", err)
					return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Error opening ZIP file: %v", err))
				}
				defer zipFile.Close()
			}

			// Проверяем ZIP файл для XML импорта
			zipXmlFiles := form.File["xml_zip"]
			if len(zipXmlFiles) > 0 {
				zipXmlHeader := zipXmlFiles[0]
				log.Printf("Processing uploaded XML ZIP file: %s, size: %d", zipXmlHeader.Filename, zipXmlHeader.Size)

				var err error
				xmlZipFile, err = zipXmlHeader.Open()
				if err != nil {
					log.Printf("Error opening XML ZIP file: %v", err)
					return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Error opening XML ZIP file: %v", err))
				}
				defer xmlZipFile.Close()

				// Определяем тип содержимого файла
				extension := strings.ToLower(filepath.Ext(zipXmlHeader.Filename))
				if extension == ".zip" {
					history, err = h.services.Storefront().ImportXMLFromZip(c.Context(), sourceID, xmlZipFile, userID)
					if err != nil {
						log.Printf("Error importing XML from ZIP: %v", err)
						return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to import XML from ZIP: %v", err))
					}
					return utils.SuccessResponse(c, history)
				}
			}
		}
	}

	if csvFile != nil {
		// Если у нас есть CSV, импортируем с ним
		history, err = h.services.Storefront().ImportCSV(c.Context(), sourceID, csvFile, zipFile, userID)
		if err != nil {
			log.Printf("Error importing CSV: %v", err)
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to import CSV: %v", err))
		}
	} else if source.URL != "" {
		// Если CSV файл не загружен, но есть URL в источнике
		// Проверяем, если URL оканчивается на .zip, предполагаем, что это XML в ZIP
		if strings.HasSuffix(strings.ToLower(source.URL), ".zip") {
			log.Printf("Detected ZIP URL for source ID %d: %s", sourceID, source.URL)
			// Загружаем ZIP-архив
			resp, err := http.Get(source.URL)
			if err != nil {
				log.Printf("Error downloading ZIP from URL: %v", err)
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to download ZIP from URL: %v", err))
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Printf("Bad status when downloading ZIP: %s", resp.Status)
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to download ZIP: %s", resp.Status))
			}

			history, err = h.services.Storefront().ImportXMLFromZip(c.Context(), sourceID, resp.Body, userID)
			if err != nil {
				log.Printf("Error importing XML from ZIP URL: %v", err)
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to import XML from ZIP URL: %v", err))
			}
		} else {
			history, err = h.services.Storefront().RunImport(c.Context(), sourceID, userID)
			if err != nil {
				log.Printf("Error running import from URL: %v", err)
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, fmt.Sprintf("Failed to run import: %v", err))
			}
		}
	} else {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No CSV file uploaded, no XML ZIP file uploaded, and no URL configured for import")
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
