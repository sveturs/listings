package handler

import (
	"io"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"backend/internal/domain/models"
	"backend/internal/proj/storefronts/service"
)

// ImportHandler handles product import endpoints
type ImportHandler struct {
	importService *service.ImportService
}

// NewImportHandler creates a new import handler
func NewImportHandler(importService *service.ImportService) *ImportHandler {
	return &ImportHandler{
		importService: importService,
	}
}

// ImportFromURL imports products from a URL
// @Summary Import products from URL
// @Description Import products from a URL (supports XML, CSV, ZIP formats)
// @Tags storefronts,import
// @Accept json
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param request body models.ImportRequest true "Import request"
// @Success 200 {object} models.ImportJob "Import job created"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/import/url [post]
func (h *ImportHandler) ImportFromURL(c *fiber.Ctx) error {
	// Try to get storefront ID from locals first (for slug-based routes)
	var storefrontID int
	var err error
	
	if id, ok := c.Locals("storefrontID").(int); ok {
		storefrontID = id
	} else {
		// Fall back to path parameter
		storefrontID, err = strconv.Atoi(c.Params("storefront_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "Invalid storefront ID",
				Message: err.Error(),
			})
		}
	}

	var req models.ImportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Set storefront ID from path parameter
	req.StorefrontID = storefrontID

	// Validate request
	if req.FileURL == nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "file_url is required",
		})
	}

	if req.FileType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "file_type is required",
		})
	}

	// Set default values
	if req.UpdateMode == "" {
		req.UpdateMode = "upsert"
	}
	if req.CategoryMappingMode == "" {
		req.CategoryMappingMode = "auto"
	}

	// Start import job
	job, err := h.importService.ImportFromURL(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to start import",
			Message: err.Error(),
		})
	}

	return c.JSON(job)
}

// ImportFromFile imports products from uploaded file
// @Summary Import products from file
// @Description Import products from uploaded file (supports XML, CSV, ZIP formats)
// @Tags storefronts,import
// @Accept multipart/form-data
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param file formData file true "Import file"
// @Param file_type formData string true "File type" Enums(xml,csv,zip)
// @Param update_mode formData string false "Update mode" Enums(create_only,update_only,upsert) default(upsert)
// @Param category_mapping_mode formData string false "Category mapping mode" Enums(auto,manual,skip) default(auto)
// @Success 200 {object} models.ImportJob "Import job created"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/import/file [post]
func (h *ImportHandler) ImportFromFile(c *fiber.Ctx) error {
	// Try to get storefront ID from locals first (for slug-based routes)
	var storefrontID int
	var err error
	
	if id, ok := c.Locals("storefrontID").(int); ok {
		storefrontID = id
	} else {
		// Fall back to path parameter
		storefrontID, err = strconv.Atoi(c.Params("storefront_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "Invalid storefront ID",
				Message: err.Error(),
			})
		}
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "File upload failed",
			Message: err.Error(),
		})
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Failed to open uploaded file",
			Message: err.Error(),
		})
	}
	defer src.Close()

	// Read file data
	fileData, err := io.ReadAll(src)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Failed to read uploaded file",
			Message: err.Error(),
		})
	}

	// Parse form data
	req := models.ImportRequest{
		StorefrontID: storefrontID,
		FileType:     c.FormValue("file_type"),
		UpdateMode:   c.FormValue("update_mode"),
		CategoryMappingMode: c.FormValue("category_mapping_mode"),
	}

	// Validate file type
	if req.FileType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "file_type is required",
		})
	}

	// Set default values
	if req.UpdateMode == "" {
		req.UpdateMode = "upsert"
	}
	if req.CategoryMappingMode == "" {
		req.CategoryMappingMode = "auto"
	}

	// Set filename
	fileName := file.Filename
	req.FileName = &fileName

	// Start import job
	job, err := h.importService.ImportFromFile(c.Context(), fileData, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to start import",
			Message: err.Error(),
		})
	}

	return c.JSON(job)
}

// ValidateImportFile validates import file without importing
// @Summary Validate import file
// @Description Validate import file structure and data without actually importing
// @Tags storefronts,import
// @Accept multipart/form-data
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param file formData file true "Import file"
// @Param file_type formData string true "File type" Enums(xml,csv,zip)
// @Success 200 {object} models.ImportJobStatus "Validation result"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/import/validate [post]
func (h *ImportHandler) ValidateImportFile(c *fiber.Ctx) error {
	// Try to get storefront ID from locals first (for slug-based routes)
	var storefrontID int
	var err error
	
	if id, ok := c.Locals("storefrontID").(int); ok {
		storefrontID = id
	} else {
		// Fall back to path parameter
		storefrontID, err = strconv.Atoi(c.Params("storefront_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "Invalid storefront ID",
				Message: err.Error(),
			})
		}
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "File upload failed",
			Message: err.Error(),
		})
	}

	// Get file type
	fileType := c.FormValue("file_type")
	if fileType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "file_type is required",
		})
	}

	// Open and read file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Failed to open uploaded file",
			Message: err.Error(),
		})
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Failed to read uploaded file",
			Message: err.Error(),
		})
	}

	// Validate file
	status, err := h.importService.ValidateImportFile(c.Context(), fileData, fileType, storefrontID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "File validation failed",
			Message: err.Error(),
		})
	}

	return c.JSON(status)
}

// GetCSVTemplate returns CSV template for product import
// @Summary Get CSV import template
// @Description Get CSV template with headers and example data for product import
// @Tags storefronts,import
// @Produce text/csv
// @Success 200 {string} string "CSV template"
// @Router /api/v1/storefronts/import/csv-template [get]
func (h *ImportHandler) GetCSVTemplate(c *fiber.Ctx) error {
	template := h.importService.GetCSVTemplate()
	
	// Convert to CSV format
	var csvContent string
	for _, row := range template {
		for j, cell := range row {
			if j > 0 {
				csvContent += ","
			}
			// Escape quotes and wrap in quotes if contains comma or quotes
			if containsCommaOrQuote(cell) {
				cell = `"` + strings.ReplaceAll(cell, `"`, `""`) + `"`
			}
			csvContent += cell
		}
		csvContent += "\n"
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="product_import_template.csv"`)
	
	return c.SendString(csvContent)
}

// GetImportFormats returns supported import formats and their descriptions
// @Summary Get supported import formats
// @Description Get information about supported import formats and their requirements
// @Tags storefronts,import
// @Produce json
// @Success 200 {object} map[string]interface{} "Import formats information"
// @Router /api/v1/storefronts/import/formats [get]
func (h *ImportHandler) GetImportFormats(c *fiber.Ctx) error {
	formats := map[string]interface{}{
		"xml": map[string]interface{}{
			"description": "XML format supporting Digital Vision Serbian standard",
			"file_extensions": []string{".xml"},
			"sample_structure": map[string]interface{}{
				"root": "artikli",
				"item": "artikal",
				"fields": []string{"id", "sifra", "naziv", "kategorija1", "kategorija2", "kategorija3", 
					"uvoznik", "godinaUvoza", "zemljaPorekla", "vpCena", "mpCena", "dostupan", 
					"naAkciji", "opis", "barKod", "slike"},
			},
		},
		"csv": map[string]interface{}{
			"description": "CSV format with customizable headers",
			"file_extensions": []string{".csv"},
			"required_headers": []string{"name", "price", "currency"},
			"optional_headers": []string{"sku", "description", "wholesale_price", "category", 
				"stock_quantity", "barcode", "image_url", "is_active", "on_sale", "sale_price", 
				"brand", "model", "country_of_origin"},
			"encoding": "UTF-8",
			"delimiter": ",",
		},
		"zip": map[string]interface{}{
			"description": "ZIP archive containing XML or CSV files",
			"file_extensions": []string{".zip"},
			"supported_contents": []string{"xml", "csv"},
			"note": "All supported files in the archive will be processed",
		},
	}

	return c.JSON(map[string]interface{}{
		"supported_formats": formats,
		"update_modes": []string{"create_only", "update_only", "upsert"},
		"category_mapping_modes": []string{"auto", "manual", "skip"},
		"max_file_size": "100MB",
		"max_products_per_import": 10000,
	})
}

// Helper function to check if string contains comma or quote
func containsCommaOrQuote(s string) bool {
	return strings.Contains(s, ",") || strings.Contains(s, `"`)
}