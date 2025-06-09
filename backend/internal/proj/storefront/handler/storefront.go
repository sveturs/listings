// Package handler
// backend/internal/proj/storefront/handler/storefront.go
package handler

import (
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/pkg/utils"
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
// @Summary Create new storefront
// @Description Creates a new storefront for the authenticated user
// @Tags storefronts
// @Accept json
// @Produce json
// @Param body body models.StorefrontCreate true "Storefront data"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.Storefront} "Storefront created successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "storefront.invalidRequest"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 402 {object} utils.ErrorResponseSwag "storefront.insufficientFunds"
// @Failure 500 {object} utils.ErrorResponseSwag "storefront.createError"
// @Security BearerAuth
// @Router /api/v1/storefronts [post]
func (h *StorefrontHandler) CreateStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var create models.StorefrontCreate
	if err := c.BodyParser(&create); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefront.invalidRequest")
	}

	// Валидация
	if create.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefront.nameRequired")
	}

	storefront, err := h.services.Storefront().CreateStorefront(c.Context(), userID, &create)
	if err != nil {
		if err.Error() == "insufficient funds" {
			return utils.ErrorResponse(c, fiber.StatusPaymentRequired, "storefront.insufficientFunds")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefront.createError")
	}

	return utils.SuccessResponse(c, storefront)
}

// GetUserStorefronts возвращает все витрины пользователя
// @Summary Get user storefronts
// @Description Returns all storefronts owned by the authenticated user
// @Tags storefronts
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.Storefront} "List of user storefronts"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "storefront.getError"
// @Security BearerAuth
// @Router /api/v1/storefronts [get]
func (h *StorefrontHandler) GetUserStorefronts(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	storefronts, err := h.services.Storefront().GetUserStorefronts(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefront.getError")
	}

	return utils.SuccessResponse(c, storefronts)
}

// GetStorefront возвращает витрину по ID
// @Summary Get storefront by ID
// @Description Returns a specific storefront by ID
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.Storefront} "Storefront details"
// @Failure 400 {object} utils.ErrorResponseSwag "storefront.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "storefront.getError"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id} [get]
func (h *StorefrontHandler) GetStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefront.invalidId")
	}

	storefront, err := h.services.Storefront().GetStorefrontByID(c.Context(), id, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefront.getError")
	}

	return utils.SuccessResponse(c, storefront)
}

// GetCategoryMappings возвращает сопоставления категорий для источника импорта
// @Summary Get category mappings for import source
// @Description Returns category mappings for a specific import source
// @Tags import-sources
// @Accept json
// @Produce json
// @Param id path int true "Import source ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]int} "Category mappings"
// @Failure 400 {object} utils.ErrorResponseSwag "importSource.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "importSource.getCategoryMappingsError"
// @Security BearerAuth
// @Router /api/v1/import-sources/{id}/category-mappings [get]
func (h *StorefrontHandler) GetCategoryMappings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	sourceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "importSource.invalidId")
	}

	mappings, err := h.services.Storefront().GetCategoryMappings(c.Context(), sourceID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.getCategoryMappingsError")
	}

	return utils.SuccessResponse(c, mappings)
}

// UpdateCategoryMappings обновляет сопоставления категорий для источника импорта
// @Summary Update category mappings for import source
// @Description Updates category mappings for a specific import source
// @Tags import-sources
// @Accept json
// @Produce json
// @Param id path int true "Import source ID"
// @Param body body map[string]int true "Category mappings"
// @Success 200 {object} utils.SuccessResponseSwag{data=CategoryMappingsUpdateResponse} "Update successful"
// @Failure 400 {object} utils.ErrorResponseSwag "importSource.invalidId,importSource.invalidMappingsFormat"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "importSource.updateCategoryMappingsError"
// @Security BearerAuth
// @Router /api/v1/import-sources/{id}/category-mappings [put]
func (h *StorefrontHandler) UpdateCategoryMappings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	sourceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "importSource.invalidId")
	}

	var mappings map[string]int
	if err := c.BodyParser(&mappings); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "importSource.invalidMappingsFormat")
	}

	err = h.services.Storefront().UpdateCategoryMappings(c.Context(), sourceID, userID, mappings)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.updateCategoryMappingsError")
	}

	return utils.SuccessResponse(c, &CategoryMappingsUpdateResponse{
		Message: "importSource.categoryMappingsUpdated",
	})
}

// GetImportedCategories возвращает список категорий, которые были импортированы этим источником
// @Summary Get imported categories
// @Description Returns list of categories that were imported by this source
// @Tags import-sources
// @Accept json
// @Produce json
// @Param id path int true "Import source ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]string} "List of imported categories"
// @Failure 400 {object} utils.ErrorResponseSwag "importSource.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "importSource.getImportedCategoriesError"
// @Security BearerAuth
// @Router /api/v1/import-sources/{id}/imported-categories [get]
func (h *StorefrontHandler) GetImportedCategories(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	sourceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "importSource.invalidId")
	}

	categories, err := h.services.Storefront().GetImportedCategories(c.Context(), sourceID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.getImportedCategoriesError")
	}

	return utils.SuccessResponse(c, categories)
}

// ApplyCategoryMappings применяет сопоставления категорий к импортированным товарам
// @Summary Apply category mappings
// @Description Applies category mappings to imported listings
// @Tags import-sources
// @Accept json
// @Produce json
// @Param id path int true "Import source ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=ApplyCategoryMappingsResponse} "Mappings applied successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "importSource.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "importSource.applyCategoryMappingsError"
// @Security BearerAuth
// @Router /api/v1/import-sources/{id}/apply-mappings [post]
func (h *StorefrontHandler) ApplyCategoryMappings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	sourceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "importSource.invalidId")
	}

	updatedCount, err := h.services.Storefront().ApplyCategoryMappings(c.Context(), sourceID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.applyCategoryMappingsError")
	}

	return utils.SuccessResponse(c, &ApplyCategoryMappingsResponse{
		Message:      "importSource.categoriesUpdated",
		UpdatedCount: updatedCount,
	})
}

// GetPublicStorefront возвращает публичные данные витрины по ID
// @Summary Get public storefront
// @Description Returns public data of a storefront (no authentication required)
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.Storefront} "Public storefront data"
// @Failure 400 {object} utils.ErrorResponseSwag "storefront.invalidId"
// @Failure 403 {object} utils.ErrorResponseSwag "storefront.notActive"
// @Failure 404 {object} utils.ErrorResponseSwag "storefront.notFound"
// @Failure 500 {object} utils.ErrorResponseSwag "storefront.getError"
// @Router /api/v1/public/storefronts/{id} [get]
func (h *StorefrontHandler) GetPublicStorefront(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefront.invalidId")
	}

	storefront, err := h.services.Storefront().GetPublicStorefrontByID(c.Context(), id)
	if err != nil {
		// Проверяем различные типы ошибок
		if err.Error() == "storefront not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "storefront.notFound")
		}
		if err.Error() == "storefront is not active" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefront.notActive")
		}

		// Логируем детали ошибки для отладки
		logger.Error().Err(err).Int("storefront_id", id).Msg("Error getting public storefront")

		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefront.getError")
	}

	return utils.SuccessResponse(c, storefront)
}

// UpdateStorefront обновляет информацию о витрине
// @Summary Update storefront
// @Description Updates storefront information
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param body body models.Storefront true "Updated storefront data"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Update successful"
// @Failure 400 {object} utils.ErrorResponseSwag "storefront.invalidId,storefront.invalidRequest"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "storefront.updateError"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id} [put]
func (h *StorefrontHandler) UpdateStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefront.invalidId")
	}

	var storefront models.Storefront
	if err := c.BodyParser(&storefront); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefront.invalidRequest")
	}
	storefront.ID = id

	err = h.services.Storefront().UpdateStorefront(c.Context(), &storefront, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefront.updateError")
	}

	return utils.SuccessResponse(c, &MessageResponse{Message: "storefront.updateSuccess"})
}

// DeleteStorefront удаляет витрину
// @Summary Delete storefront
// @Description Deletes a storefront
// @Tags storefronts
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Delete successful"
// @Failure 400 {object} utils.ErrorResponseSwag "storefront.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "storefront.deleteError"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id} [delete]
func (h *StorefrontHandler) DeleteStorefront(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefront.invalidId")
	}

	err = h.services.Storefront().DeleteStorefront(c.Context(), id, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefront.deleteError")
	}

	return utils.SuccessResponse(c, &MessageResponse{Message: "storefront.deleteSuccess"})
}

// CreateImportSource создаёт новый источник импорта
// @Summary Create import source
// @Description Creates a new import source for storefront
// @Tags import-sources
// @Accept json
// @Produce json
// @Param body body models.ImportSourceCreate true "Import source data"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.ImportSource} "Import source created"
// @Failure 400 {object} utils.ErrorResponseSwag "storefront.invalidRequest"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "importSource.createError"
// @Security BearerAuth
// @Router /api/v1/import-sources [post]
func (h *StorefrontHandler) CreateImportSource(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var create models.ImportSourceCreate
	if err := c.BodyParser(&create); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefront.invalidRequest")
	}

	source, err := h.services.Storefront().CreateImportSource(c.Context(), &create, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.createError")
	}

	return utils.SuccessResponse(c, source)
}

// GetImportSources возвращает источники импорта для витрины
// @Summary Get import sources for storefront
// @Description Returns all import sources for a specific storefront
// @Tags import-sources
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.ImportSource} "List of import sources"
// @Failure 400 {object} utils.ErrorResponseSwag "storefront.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "importSource.getError"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/import-sources [get]
func (h *StorefrontHandler) GetImportSources(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefront.invalidId")
	}

	sources, err := h.services.Storefront().GetImportSources(c.Context(), storefrontID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.getError")
	}

	return utils.SuccessResponse(c, sources)
}

// UpdateImportSource обновляет источник импорта
// @Summary Update import source
// @Description Updates import source information
// @Tags import-sources
// @Accept json
// @Produce json
// @Param id path int true "Import source ID"
// @Param body body models.ImportSource true "Updated import source data"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Update successful"
// @Failure 400 {object} utils.ErrorResponseSwag "importSource.invalidId,storefront.invalidRequest"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "importSource.updateError"
// @Security BearerAuth
// @Router /api/v1/import-sources/{id} [put]
func (h *StorefrontHandler) UpdateImportSource(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "importSource.invalidId")
	}

	var source models.ImportSource
	if err := c.BodyParser(&source); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefront.invalidRequest")
	}
	source.ID = id

	err = h.services.Storefront().UpdateImportSource(c.Context(), &source, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.updateError")
	}

	return utils.SuccessResponse(c, &MessageResponse{Message: "importSource.updateSuccess"})
}

// DeleteImportSource удаляет источник импорта
// @Summary Delete import source
// @Description Deletes an import source
// @Tags import-sources
// @Accept json
// @Produce json
// @Param id path int true "Import source ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Delete successful"
// @Failure 400 {object} utils.ErrorResponseSwag "importSource.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "importSource.deleteError"
// @Security BearerAuth
// @Router /api/v1/import-sources/{id} [delete]
func (h *StorefrontHandler) DeleteImportSource(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "importSource.invalidId")
	}

	err = h.services.Storefront().DeleteImportSource(c.Context(), id, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.deleteError")
	}

	return utils.SuccessResponse(c, &MessageResponse{Message: "importSource.deleteSuccess"})
}

// RunImport запускает импорт данных
// @Summary Run import
// @Description Starts import process for an import source
// @Tags import-sources
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Import source ID"
// @Param file formData file false "CSV file for import"
// @Param images_zip formData file false "ZIP file with images"
// @Param xml_zip formData file false "ZIP file with XML data"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.ImportHistory} "Import completed"
// @Failure 400 {object} utils.ErrorResponseSwag "importSource.invalidId,importSource.notFound,importSource.noDataSource"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "importSource.runImportError"
// @Security BearerAuth
// @Router /api/v1/import-sources/{id}/import [post]
func (h *StorefrontHandler) RunImport(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	sourceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "importSource.invalidId")
	}

	// Отладочный лог
	logger.Debug().Int("source_id", sourceID).Int("user_id", userID).Msg("Running import")

	// Сначала проверим, существует ли источник импорта
	source, err := h.services.Storefront().GetImportSourceByID(c.Context(), sourceID, userID)
	if err != nil {
		logger.Error().Err(err).Int("source_id", sourceID).Msg("Error fetching source")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "importSource.notFound")
	}

	// Отладочный лог
	logger.Debug().Interface("source", source).Msg("Found import source")

	var history *models.ImportHistory
	var csvFile, zipFile multipart.File
	var xmlZipFile multipart.File

	// Проверяем тип контента запроса
	contentType := c.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		// Это запрос с файлами
		form, err := c.MultipartForm()
		if err != nil && err != fiber.ErrUnprocessableEntity {
			logger.Error().Err(err).Msg("Error processing form file")
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "importSource.fileProcessingError")
		}

		// Обработка файлов из формы если они есть
		if form != nil {
			// Проверяем файл CSV
			csvFiles := form.File["file"]
			if len(csvFiles) > 0 {
				csvFileHeader := csvFiles[0]
				logger.Debug().Str("filename", csvFileHeader.Filename).Int64("size", csvFileHeader.Size).Msg("Processing uploaded CSV file")

				var err error
				csvFile, err = csvFileHeader.Open()
				if err != nil {
					logger.Error().Err(err).Msg("Error opening CSV file")
					return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.csvOpenError")
				}
				defer csvFile.Close()
			}

			// Проверяем ZIP файл с изображениями
			zipFiles := form.File["images_zip"]
			if len(zipFiles) > 0 {
				zipFileHeader := zipFiles[0]
				logger.Debug().Str("filename", zipFileHeader.Filename).Int64("size", zipFileHeader.Size).Msg("Processing uploaded ZIP file")

				var err error
				zipFile, err = zipFileHeader.Open()
				if err != nil {
					logger.Error().Err(err).Msg("Error opening ZIP file")
					return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.zipOpenError")
				}
				defer zipFile.Close()
			}

			// Проверяем ZIP файл для XML импорта
			zipXmlFiles := form.File["xml_zip"]
			if len(zipXmlFiles) > 0 {
				zipXmlHeader := zipXmlFiles[0]
				logger.Debug().Str("filename", zipXmlHeader.Filename).Int64("size", zipXmlHeader.Size).Msg("Processing uploaded XML ZIP file")

				var err error
				xmlZipFile, err = zipXmlHeader.Open()
				if err != nil {
					logger.Error().Err(err).Msg("Error opening XML ZIP file")
					return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.xmlZipOpenError")
				}
				defer xmlZipFile.Close()

				// Определяем тип содержимого файла
				extension := strings.ToLower(filepath.Ext(zipXmlHeader.Filename))
				if extension == ".zip" {
					history, err = h.services.Storefront().ImportXMLFromZip(c.Context(), sourceID, xmlZipFile, userID)
					if err != nil {
						logger.Error().Err(err).Msg("Error importing XML from ZIP")
						return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.xmlImportError")
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
			logger.Error().Err(err).Msg("Error importing CSV")
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.csvImportError")
		}
	} else if source.URL != "" {
		// Если CSV файл не загружен, но есть URL в источнике
		// Проверяем, если URL оканчивается на .zip, предполагаем, что это XML в ZIP
		if strings.HasSuffix(strings.ToLower(source.URL), ".zip") {
			logger.Debug().Int("source_id", sourceID).Str("url", source.URL).Msg("Detected ZIP URL")
			// Загружаем ZIP-архив
			resp, err := http.Get(source.URL)
			if err != nil {
				logger.Error().Err(err).Msg("Error downloading ZIP from URL")
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.downloadError")
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				logger.Error().Str("status", resp.Status).Msg("Bad status when downloading ZIP")
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.downloadError")
			}

			history, err = h.services.Storefront().ImportXMLFromZip(c.Context(), sourceID, resp.Body, userID)
			if err != nil {
				logger.Error().Err(err).Msg("Error importing XML from ZIP URL")
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.xmlImportError")
			}
		} else {
			history, err = h.services.Storefront().RunImport(c.Context(), sourceID, userID)
			if err != nil {
				logger.Error().Err(err).Msg("Error running import from URL")
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.runImportError")
			}
		}
	} else {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "importSource.noDataSource")
	}

	// Отладочный лог
	logger.Debug().Interface("history", history).Msg("Import completed successfully")

	return utils.SuccessResponse(c, history)
}

// GetImportHistory возвращает историю импорта
// @Summary Get import history
// @Description Returns import history for an import source
// @Tags import-sources
// @Accept json
// @Produce json
// @Param id path int true "Import source ID"
// @Param limit query int false "Limit number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.ImportHistory} "Import history"
// @Failure 400 {object} utils.ErrorResponseSwag "importSource.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "importSource.getHistoryError"
// @Security BearerAuth
// @Router /api/v1/import-sources/{id}/history [get]
func (h *StorefrontHandler) GetImportHistory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	sourceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "importSource.invalidId")
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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "importSource.getHistoryError")
	}

	return utils.SuccessResponse(c, history)
}
