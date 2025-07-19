package handler

import (
	"net/http"
	"strconv"

	"backend/internal/services"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// ImageHandler - обработчик для работы с изображениями товаров витрин
type ImageHandler struct {
	imageService *services.ImageService
}

// NewImageHandler создает новый ImageHandler
func NewImageHandler(imageService *services.ImageService) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
	}
}

// UploadProductImage загружает изображение для товара витрины
// @Summary Upload image for storefront product
// @Description Uploads a new image for a storefront product
// @Tags storefront-images
// @Accept multipart/form-data
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param product_id path int true "Product ID"
// @Param image formData file true "Image file"
// @Param is_main formData bool false "Set as main image"
// @Param display_order formData int false "Display order"
// @Security Bearer
// @Success 200 {object} utils.SuccessResponseSwag{data=services.UploadImageResponse} "Image uploaded successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 413 {object} utils.ErrorResponseSwag "File too large"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/slug/{slug}/products/{product_id}/images [post]
func (h *ImageHandler) UploadProductImage(c *fiber.Ctx) error {
	// Получение ID товара
	productIDStr := c.Params("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "storefronts.invalid_product_id")
	}

	// Получение файла из формы
	file, err := c.FormFile("image")
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "storefronts.no_image_file")
	}

	// Открытие файла
	src, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "storefronts.file_open_error")
	}
	defer src.Close()

	// Получение дополнительных параметров
	isMain := c.FormValue("is_main") == "true"
	displayOrder := 0
	if displayOrderStr := c.FormValue("display_order"); displayOrderStr != "" {
		if order, err := strconv.Atoi(displayOrderStr); err == nil {
			displayOrder = order
		}
	}

	// Создание запроса для загрузки изображения
	uploadRequest := &services.UploadImageRequest{
		EntityType:   services.ImageTypeStorefrontProduct,
		EntityID:     productID,
		File:         src,
		FileHeader:   file,
		IsMain:       isMain,
		DisplayOrder: displayOrder,
	}

	// Загрузка изображения
	response, err := h.imageService.UploadImage(c.Context(), uploadRequest)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "storefronts.upload_failed")
	}

	return utils.SuccessResponse(c, response)
}

// GetProductImages получает все изображения товара
// @Summary Get product images
// @Description Returns all images for a specific product
// @Tags storefront-images
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param product_id path int true "Product ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]services.UploadImageResponse} "Product images"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 404 {object} utils.ErrorResponseSwag "Product not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/slug/{slug}/products/{product_id}/images [get]
func (h *ImageHandler) GetProductImages(c *fiber.Ctx) error {
	// Получение ID товара
	productIDStr := c.Params("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "storefronts.invalid_product_id")
	}

	// Получение изображений
	images, err := h.imageService.GetImagesByEntity(c.Context(), services.ImageTypeStorefrontProduct, productID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "storefronts.get_images_failed")
	}

	// Конвертация в формат ответа
	responses := make([]services.UploadImageResponse, 0, len(images))
	for _, img := range images {
		responses = append(responses, services.UploadImageResponse{
			ID:           img.GetID(),
			ImageURL:     img.GetImageURL(),
			ThumbnailURL: img.GetThumbnailURL(),
			PublicURL:    img.GetImageURL(),
			IsMain:       img.IsMainImage(),
			DisplayOrder: img.GetDisplayOrder(),
		})
	}

	return utils.SuccessResponse(c, responses)
}

// DeleteProductImage удаляет изображение товара
// @Summary Delete product image
// @Description Deletes a specific image of a product
// @Tags storefront-images
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param product_id path int true "Product ID"
// @Param image_id path int true "Image ID"
// @Security Bearer
// @Success 200 {object} utils.SuccessResponseSwag "Image deleted successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} utils.ErrorResponseSwag "Image not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/slug/{slug}/products/{product_id}/images/{image_id} [delete]
func (h *ImageHandler) DeleteProductImage(c *fiber.Ctx) error {
	// Получение ID изображения
	imageIDStr := c.Params("image_id")
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "storefronts.invalid_image_id")
	}

	// Удаление изображения
	err = h.imageService.DeleteImage(c.Context(), imageID, services.ImageTypeStorefrontProduct)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "storefronts.delete_failed")
	}

	return utils.SuccessResponse(c, nil)
}

// SetMainProductImage устанавливает изображение как главное
// @Summary Set main product image
// @Description Sets a specific image as the main image for a product
// @Tags storefront-images
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param product_id path int true "Product ID"
// @Param image_id path int true "Image ID"
// @Security Bearer
// @Success 200 {object} utils.SuccessResponseSwag "Main image set successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 404 {object} utils.ErrorResponseSwag "Image not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/slug/{slug}/products/{product_id}/images/{image_id}/main [post]
func (h *ImageHandler) SetMainProductImage(c *fiber.Ctx) error {
	// Получение ID товара и изображения
	productIDStr := c.Params("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "storefronts.invalid_product_id")
	}

	imageIDStr := c.Params("image_id")
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "storefronts.invalid_image_id")
	}

	// Установка главного изображения
	err = h.imageService.SetMainImage(c.Context(), imageID, services.ImageTypeStorefrontProduct, productID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "storefronts.set_main_failed")
	}

	return utils.SuccessResponse(c, nil)
}

// ImageOrderUpdate структура для обновления порядка изображений
type ImageOrderUpdate struct {
	ID           int `json:"id"`
	DisplayOrder int `json:"display_order"`
}

// UpdateImageOrder обновляет порядок отображения изображений
// @Summary Update image display order
// @Description Updates the display order of product images
// @Tags storefront-images
// @Accept json
// @Produce json
// @Param slug path string true "Storefront slug"
// @Param product_id path int true "Product ID"
// @Param request body []handler.ImageOrderUpdate true "Image order updates"
// @Security Bearer
// @Success 200 {object} utils.SuccessResponseSwag "Image order updated successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/slug/{slug}/products/{product_id}/images/order [put]
func (h *ImageHandler) UpdateImageOrder(c *fiber.Ctx) error {
	var updates []ImageOrderUpdate
	if err := c.BodyParser(&updates); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "storefronts.invalid_request")
	}

	// Обновление порядка для каждого изображения
	for _, update := range updates {
		// Здесь можно добавить валидацию, что изображение принадлежит товару
		// Но для простоты пока пропускаем

		// Обновление порядка отображения
		err := h.imageService.GetRepo().UpdateDisplayOrder(c.Context(), update.ID, update.DisplayOrder)
		if err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "storefronts.update_order_failed")
		}
	}

	return utils.SuccessResponse(c, nil)
}
