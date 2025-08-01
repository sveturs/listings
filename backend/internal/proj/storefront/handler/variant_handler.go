package handler

import (
	"net/http"
	"strconv"

	"backend/internal/proj/storefront/repository"
	"backend/internal/proj/storefront/types"

	"github.com/gofiber/fiber"
)

type VariantHandler struct {
	variantRepo *repository.VariantRepository
}

func NewVariantHandler(variantRepo *repository.VariantRepository) *VariantHandler {
	return &VariantHandler{
		variantRepo: variantRepo,
	}
}

// GetVariantAttributes godoc
// @Summary Get all variant attributes
// @Description Returns all available variant attributes (color, size, etc.)
// @Tags variants
// @Accept json
// @Produce json
// @Success 200 {array} interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/variants/attributes [get]
func (h *VariantHandler) GetVariantAttributes(c *fiber.Ctx) error {
	attributes, err := h.variantRepo.GetVariantAttributes(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get variant attributes",
		})
	}

	return c.JSON(attributes)
}

// GetVariantAttributeValues godoc
// @Summary Get values for a variant attribute
// @Description Returns all possible values for a specific variant attribute
// @Tags variants
// @Accept json
// @Produce json
// @Param attribute_id path int true "Attribute ID"
// @Success 200 {array} interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/variants/attributes/{attribute_id}/values [get]
func (h *VariantHandler) GetVariantAttributeValues(c *fiber.Ctx) error {
	attributeID, err := strconv.Atoi(c.Params("attribute_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid attribute ID",
		})
	}

	values, err := h.variantRepo.GetVariantAttributeValues(c.Context(), attributeID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attribute values",
		})
	}

	return c.JSON(values)
}

// CreateVariant godoc
// @Summary Create a new product variant
// @Description Creates a new variant for an existing product
// @Tags variants
// @Accept json
// @Produce json
// @Param variant body interface{} true "Variant data"
// @Success 201 {object} interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/variants [post]
func (h *VariantHandler) CreateVariant(c *fiber.Ctx) error {
	var req types.CreateVariantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	variant, err := h.variantRepo.CreateVariant(c.Context(), &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create variant",
		})
	}

	return c.Status(http.StatusCreated).JSON(variant)
}

// GetVariantsByProductID godoc
// @Summary Get variants for a product
// @Description Returns all variants for a specific product
// @Tags variants
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {array} interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/products/{product_id}/variants [get]
func (h *VariantHandler) GetVariantsByProductID(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("product_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	variants, err := h.variantRepo.GetVariantsByProductID(c.Context(), productID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get variants",
		})
	}

	return c.JSON(variants)
}

// GetVariantByID godoc
// @Summary Get variant by ID
// @Description Returns a specific variant with all details including images
// @Tags variants
// @Accept json
// @Produce json
// @Param variant_id path int true "Variant ID"
// @Success 200 {object} interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/variants/{variant_id} [get]
func (h *VariantHandler) GetVariantByID(c *fiber.Ctx) error {
	variantID, err := strconv.Atoi(c.Params("variant_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid variant ID",
		})
	}

	variant, err := h.variantRepo.GetVariantByID(c.Context(), variantID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Variant not found",
		})
	}

	return c.JSON(variant)
}

// GenerateVariants godoc
// @Summary Auto-generate variants for a product
// @Description Automatically generates all possible variants based on attribute matrix
// @Tags variants
// @Accept json
// @Produce json
// @Param request body interface{} true "Generation parameters"
// @Success 201 {array} interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/variants/generate [post]
func (h *VariantHandler) GenerateVariants(c *fiber.Ctx) error {
	var req types.GenerateVariantsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate that attribute matrix is not empty
	if len(req.AttributeMatrix) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Attribute matrix cannot be empty",
		})
	}

	variants, err := h.variantRepo.GenerateVariants(c.Context(), &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate variants: " + err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":        "Variants generated successfully",
		"variants_count": len(variants),
		"variants":       variants,
	})
}

// BulkCreateVariants godoc
// @Summary Create multiple variants at once
// @Description Creates multiple variants for a product in a single request
// @Tags variants
// @Accept json
// @Produce json
// @Param request body interface{} true "Bulk creation data"
// @Success 201 {array} interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/variants/bulk [post]
func (h *VariantHandler) BulkCreateVariants(c *fiber.Ctx) error {
	var req types.BulkCreateVariantsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(req.Variants) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "No variants provided",
		})
	}

	var createdVariants []types.ProductVariant
	for _, variantReq := range req.Variants {
		variantReq.ProductID = req.ProductID // Ensure product ID matches
		variant, err := h.variantRepo.CreateVariant(c.Context(), &variantReq)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create variant: " + err.Error(),
			})
		}
		createdVariants = append(createdVariants, *variant)
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message":        "Variants created successfully",
		"variants_count": len(createdVariants),
		"variants":       createdVariants,
	})
}

// SetupProductAttributes godoc
// @Summary Configure attributes for a product
// @Description Allows seller to configure which attributes to use and add custom values
// @Tags variants
// @Accept json
// @Produce json
// @Param request body interface{} true "Attribute configuration"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/products/attributes/setup [post]
func (h *VariantHandler) SetupProductAttributes(c *fiber.Ctx) error {
	var req types.SetupProductAttributesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := h.variantRepo.SetupProductAttributes(c.Context(), &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to setup product attributes: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product attributes configured successfully",
	})
}

// GetProductAttributes godoc
// @Summary Get configured attributes for a product
// @Description Returns all configured attributes for a seller's product
// @Tags variants
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {array} interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/products/{product_id}/attributes [get]
func (h *VariantHandler) GetProductAttributes(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("product_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	attributes, err := h.variantRepo.GetProductAttributes(c.Context(), productID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get product attributes",
		})
	}

	return c.JSON(attributes)
}

// GetAvailableAttributesForCategory godoc
// @Summary Get available attributes for a category
// @Description Returns all attributes that can be used for products in a specific category
// @Tags variants
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {array} interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/categories/{category_id}/attributes [get]
func (h *VariantHandler) GetAvailableAttributesForCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.Atoi(c.Params("category_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid category ID",
		})
	}

	attributes, err := h.variantRepo.GetAvailableAttributesForCategory(c.Context(), categoryID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get available attributes",
		})
	}

	return c.JSON(attributes)
}

// DeleteVariant godoc
// @Summary Delete a variant
// @Description Soft deletes a variant by setting is_active to false
// @Tags variants
// @Accept json
// @Produce json
// @Param variant_id path int true "Variant ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/variants/{variant_id} [delete]
func (h *VariantHandler) DeleteVariant(c *fiber.Ctx) error {
	variantID, err := strconv.Atoi(c.Params("variant_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid variant ID",
		})
	}

	err = h.variantRepo.DeleteVariant(c.Context(), variantID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete variant",
		})
	}

	return c.JSON(fiber.Map{
		"message":    "Variant deleted successfully",
		"variant_id": variantID,
	})
}

// UpdateVariant godoc
// @Summary Update a variant
// @Description Updates an existing variant's information
// @Tags variants
// @Accept json
// @Produce json
// @Param variant_id path int true "Variant ID"
// @Param variant body types.UpdateVariantRequest true "Updated variant data"
// @Success 200 {object} types.ProductVariant
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/variants/{variant_id} [put]
func (h *VariantHandler) UpdateVariant(c *fiber.Ctx) error {
	variantID, err := strconv.Atoi(c.Params("variant_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid variant ID",
		})
	}

	var req types.UpdateVariantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	variant, err := h.variantRepo.UpdateVariant(c.Context(), variantID, &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update variant",
		})
	}

	return c.JSON(variant)
}

// GetVariantMatrix godoc
// @Summary Get variant matrix for a product
// @Description Returns all possible variant combinations and existing variants
// @Tags variants
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} types.VariantMatrixResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/products/{product_id}/variant-matrix [get]
func (h *VariantHandler) GetVariantMatrix(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("product_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	matrix, err := h.variantRepo.GetVariantMatrix(c.Context(), productID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get variant matrix",
		})
	}

	return c.JSON(matrix)
}

// BulkUpdateStock godoc
// @Summary Bulk update stock quantities
// @Description Updates stock quantities for multiple variants at once
// @Tags variants
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Param request body types.BulkUpdateStockRequest true "Stock updates"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/products/{product_id}/variants/bulk-update-stock [post]
func (h *VariantHandler) BulkUpdateStock(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("product_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var req types.BulkUpdateStockRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	updatedCount, err := h.variantRepo.BulkUpdateStock(c.Context(), productID, &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update stock",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Stock updated successfully",
		"updated_count": updatedCount,
	})
}

// GetVariantAnalytics godoc
// @Summary Get variant analytics
// @Description Returns analytics data for product variants
// @Tags variants
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} types.VariantAnalyticsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/products/{product_id}/variants/analytics [get]
func (h *VariantHandler) GetVariantAnalytics(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("product_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	analytics, err := h.variantRepo.GetVariantAnalytics(c.Context(), productID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get analytics",
		})
	}

	return c.JSON(analytics)
}

// ImportVariants godoc
// @Summary Import variants from CSV
// @Description Imports multiple variants from a CSV file
// @Tags variants
// @Accept multipart/form-data
// @Produce json
// @Param product_id path int true "Product ID"
// @Param file formData file true "CSV file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/products/{product_id}/variants/import [post]
func (h *VariantHandler) ImportVariants(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("product_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open file",
		})
	}
	defer src.Close()

	// Read file content
	fileData := make([]byte, file.Size)
	_, err = src.Read(fileData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read file",
		})
	}

	importedCount, err := h.variantRepo.ImportVariants(c.Context(), productID, fileData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to import variants: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Variants imported successfully",
		"imported_count": importedCount,
	})
}

// ExportVariants godoc
// @Summary Export variants to CSV
// @Description Exports all variants of a product to CSV format
// @Tags variants
// @Accept json
// @Produce text/csv
// @Param product_id path int true "Product ID"
// @Success 200 {file} binary
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/storefront/products/{product_id}/variants/export [get]
func (h *VariantHandler) ExportVariants(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("product_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	csvData, fileName, err := h.variantRepo.ExportVariants(c.Context(), productID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to export variants",
		})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename="+fileName)
	
	c.Send(csvData)
	return nil
}
