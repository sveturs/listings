package handler

import (
	"context"
	"strconv"
	"strings"

	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/proj/b2c/common"
	"backend/internal/proj/b2c/service"

	"github.com/gofiber/fiber/v2"
)

// ProductHandler handles HTTP requests for storefront products
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetProducts retrieves products for a storefront
// @Summary Get storefront products
// @Description Returns paginated list of products for a storefront
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param category_id query int false "Filter by category ID"
// @Param search query string false "Search in name and description"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param stock_status query string false "Stock status filter (in_stock, low_stock, out_of_stock)"
// @Param is_active query bool false "Filter by active status"
// @Param sku query string false "Filter by SKU"
// @Param barcode query string false "Filter by barcode"
// @Param sort_by query string false "Sort by field (name, price, created_at, stock_quantity)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param offset query int false "Number of items to skip"
// @Success 200 {object} []models.StorefrontProduct "List of products"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/b2c_stores/{slug}/products [get]
func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	storefrontID, err := getStorefrontIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid storefront ID",
		})
	}

	filter := models.ProductFilter{
		StorefrontID: storefrontID,
		Limit:        20, // default limit
	}

	// Parse query parameters
	if categoryID := c.QueryInt("category_id"); categoryID > 0 {
		filter.CategoryID = &categoryID
	}

	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}

	if minPrice, err := strconv.ParseFloat(c.Query("min_price"), 64); err == nil {
		filter.MinPrice = &minPrice
	}

	if maxPrice, err := strconv.ParseFloat(c.Query("max_price"), 64); err == nil {
		filter.MaxPrice = &maxPrice
	}

	if stockStatus := c.Query("stock_status"); stockStatus != "" {
		filter.StockStatus = &stockStatus
	}

	if isActive := c.Query("is_active"); isActive != "" {
		// Если передан параметр "all" - не фильтруем по is_active (показываем все)
		if isActive != "all" {
			active := isActive == boolValueTrue
			filter.IsActive = &active
		}
		// Если "all" - filter.IsActive остается nil (показываем все товары)
	} else {
		// Для публичного API (не-админов) показываем только активные товары
		isAdmin := authMiddleware.IsAdmin(c)
		if !isAdmin {
			activeTrue := true
			filter.IsActive = &activeTrue
		}
		// Для админов: если is_active не указан, показываем все товары (nil)
	}

	if sku := c.Query("sku"); sku != "" {
		filter.SKU = &sku
	}

	if barcode := c.Query("barcode"); barcode != "" {
		filter.Barcode = &barcode
	}

	filter.SortBy = c.Query("sort_by", "created_at")
	filter.SortOrder = c.Query("sort_order", "desc")

	if limit := c.QueryInt("limit"); limit > 0 {
		filter.Limit = limit
	}

	filter.Offset = c.QueryInt("offset")

	products, err := h.productService.GetProducts(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get products",
		})
	}

	return c.JSON(products)
}

// GetProduct retrieves a single product
// @Summary Get a storefront product
// @Description Returns details of a specific product
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param id path int true "Product ID"
// @Success 200 {object} models.StorefrontProduct "Product details"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/b2c_stores/{slug}/products/{id} [get]
func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	storefrontID, err := getStorefrontIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid storefront ID",
		})
	}

	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	product, err := h.productService.GetProduct(c.Context(), storefrontID, productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get product",
		})
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(product)
}

// GetProductByID retrieves a single product by its ID without requiring storefront slug
// @Summary Get a storefront product by ID
// @Description Returns details of a specific product using only product ID
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.StorefrontProduct "Product details"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/b2c_stores/products/{id} [get]
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	product, err := h.productService.GetProductByID(c.Context(), productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get product",
		})
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(product)
}

// CreateProduct creates a new product
// @Summary Create a storefront product
// @Description Creates a new product for the storefront with optional variants support
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param product body models.CreateProductRequest true "Product data with optional variants"
// @Success 201 {object} models.StorefrontProduct "Created product with variants"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/b2c_stores/{slug}/products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	storefrontID, err := getStorefrontIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid storefront ID",
		})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req models.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Обработка пустых строк для SKU и Barcode
	if req.SKU != nil && *req.SKU == "" {
		req.SKU = nil
	}
	if req.Barcode != nil && *req.Barcode == "" {
		req.Barcode = nil
	}

	product, err := h.productService.CreateProduct(c.Context(), storefrontID, userID, &req)
	if err != nil {
		// Возвращаем 400 для ошибок валидации и владения
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "ownership validation failed") ||
			strings.Contains(errorMsg, "invalid request:") ||
			strings.Contains(errorMsg, "unauthorized:") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "validation or ownership error",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// UpdateProduct updates an existing product
// @Summary Update a storefront product
// @Description Updates product details
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param id path int true "Product ID"
// @Param product body models.UpdateProductRequest true "Product update data"
// @Success 200 {object} models.SuccessResponse "Product updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/b2c_stores/{slug}/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	storefrontID, err := getStorefrontIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid storefront ID",
		})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var req models.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.productService.UpdateProduct(c.Context(), storefrontID, productID, userID, &req); err != nil {
		logger.Error().
			Err(err).
			Int("storefrontID", storefrontID).
			Int("productID", productID).
			Int("userID", userID).
			Interface("request", req).
			Msg("Failed to update product")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product updated successfully",
	})
}

// DeleteProduct deletes a product
// @Summary Delete a storefront product
// @Description Permanently deletes a product
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param id path int true "Product ID"
// @Success 200 {object} models.SuccessResponse "Product deleted successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/b2c_stores/{slug}/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	storefrontID, err := getStorefrontIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid storefront ID",
		})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	// Передаем контекст с информацией об администраторе
	isAdmin := authMiddleware.IsAdmin(c)
	ctx := context.WithValue(context.Background(), common.ContextKeyIsAdmin, isAdmin)

	if err := h.productService.DeleteProduct(ctx, storefrontID, productID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}

// HardDeleteProduct permanently deletes a product
// @Summary Hard delete a storefront product
// @Description Permanently deletes a product and all related data (reviews, favorites, images, etc.)
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param id path int true "Product ID"
// @Param hard query bool true "Must be true for hard delete"
// @Success 200 {object} models.SuccessResponse "Product permanently deleted"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/b2c_stores/{slug}/products/{id} [delete]
func (h *ProductHandler) HardDeleteProduct(c *fiber.Ctx) error {
	storefrontID, err := getStorefrontIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid storefront ID",
		})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	// Передаем контекст с информацией об администраторе
	isAdmin := authMiddleware.IsAdmin(c)
	ctx := context.WithValue(context.Background(), common.ContextKeyIsAdmin, isAdmin)

	if err := h.productService.HardDeleteProduct(ctx, storefrontID, productID, userID); err != nil {
		logger.Error().
			Err(err).
			Int("storefront_id", storefrontID).
			Int("product_id", productID).
			Int("user_id", userID).
			Msg("Failed to hard delete product")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hard delete product",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product permanently deleted",
	})
}

// UpdateInventory updates product stock
// @Summary Update product inventory
// @Description Records inventory movement and updates stock
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param id path int true "Product ID"
// @Param inventory body models.UpdateInventoryRequest true "Inventory update data"
// @Success 200 {object} models.SuccessResponse "Inventory updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Product not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/b2c_stores/{slug}/products/{id}/inventory [post]
func (h *ProductHandler) UpdateInventory(c *fiber.Ctx) error {
	storefrontID, err := getStorefrontIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid storefront ID",
		})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var req models.UpdateInventoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.productService.UpdateInventory(c.Context(), storefrontID, productID, userID, &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update inventory",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Inventory updated successfully",
	})
}

// GetProductStats returns product statistics
// @Summary Get product statistics
// @Description Returns statistics about storefront products
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Success 200 {object} models.ProductStats "Product statistics"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/b2c_stores/{slug}/products/stats [get]
func (h *ProductHandler) GetProductStats(c *fiber.Ctx) error {
	storefrontID, err := getStorefrontIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid storefront ID",
		})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	stats, err := h.productService.GetProductStats(c.Context(), storefrontID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get product stats",
		})
	}

	return c.JSON(stats)
}

// Bulk operation handlers

// BulkCreateProducts creates multiple products
// @Summary Bulk create products
// @Description Create multiple products in a single request
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param slug path string true "Storefront slug"
// @Param body body models.BulkCreateProductsRequest true "Products to create"
// @Success 200 {object} models.BulkCreateProductsResponse "Bulk creation result"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/b2c_stores/{slug}/products/bulk/create [post]
func (h *ProductHandler) BulkCreateProducts(c *fiber.Ctx) error {
	storefrontID, err := getStorefrontIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid storefront ID",
		})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req models.BulkCreateProductsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.productService.BulkCreateProducts(c.Context(), storefrontID, userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to bulk create products",
		})
	}

	return c.JSON(response)
}

// BulkUpdateProducts updates multiple products
// @Summary Bulk update products
// @Description Update multiple products in a single request
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param slug path string true "Storefront slug"
// @Param body body models.BulkUpdateProductsRequest true "Products to update"
// @Success 200 {object} models.BulkUpdateProductsResponse "Bulk update result"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/b2c_stores/{slug}/products/bulk/update [put]
func (h *ProductHandler) BulkUpdateProducts(c *fiber.Ctx) error {
	storefrontID, err := getStorefrontIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid storefront ID",
		})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req models.BulkUpdateProductsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.productService.BulkUpdateProducts(c.Context(), storefrontID, userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to bulk update products",
		})
	}

	return c.JSON(response)
}

// BulkDeleteProducts deletes multiple products
// @Summary Bulk delete products
// @Description Delete multiple products in a single request
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param slug path string true "Storefront slug"
// @Param body body models.BulkDeleteProductsRequest true "Product IDs to delete"
// @Success 200 {object} models.BulkDeleteProductsResponse "Bulk deletion result"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/b2c_stores/{slug}/products/bulk/delete [delete]
func (h *ProductHandler) BulkDeleteProducts(c *fiber.Ctx) error {
	storefrontID, err := getStorefrontIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid storefront ID",
		})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req models.BulkDeleteProductsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Передаем контекст с информацией об администраторе
	isAdmin := authMiddleware.IsAdmin(c)
	ctx := context.WithValue(context.Background(), common.ContextKeyIsAdmin, isAdmin)

	response, err := h.productService.BulkDeleteProducts(ctx, storefrontID, userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to bulk delete products",
		})
	}

	return c.JSON(response)
}

// BulkUpdateStatus updates status of multiple products
// @Summary Bulk update product status
// @Description Update active/inactive status of multiple products
// @Tags storefront-products
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param slug path string true "Storefront slug"
// @Param body body models.BulkUpdateStatusRequest true "Product IDs and status"
// @Success 200 {object} models.BulkUpdateStatusResponse "Bulk status update result"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/b2c_stores/{slug}/products/bulk/status [put]
func (h *ProductHandler) BulkUpdateStatus(c *fiber.Ctx) error {
	storefrontID, err := getStorefrontIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid storefront ID",
		})
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req models.BulkUpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.productService.BulkUpdateStatus(c.Context(), storefrontID, userID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to bulk update status",
		})
	}

	return c.JSON(response)
}

// Helper functions to get context values
func getStorefrontIDFromContext(c *fiber.Ctx) (int, error) {
	storefrontID, ok := c.Locals("storefrontID").(int)
	if !ok {
		return 0, fiber.NewError(fiber.StatusBadRequest, "storefront ID not found in context")
	}
	return storefrontID, nil
}

func getUserIDFromContext(c *fiber.Ctx) (int, error) {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "user ID not found in context")
	}
	return userID, nil
}
