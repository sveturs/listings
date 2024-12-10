//backend/internal/handlers/marketplace.go
package handlers

import (
    "backend/internal/domain/models"
    "backend/internal/services"
    "backend/pkg/utils"
    "github.com/gofiber/fiber/v2"
    "strconv"
    "log"
)

type MarketplaceHandler struct {
    services services.ServicesInterface
}

func NewMarketplaceHandler(services services.ServicesInterface) *MarketplaceHandler {
    return &MarketplaceHandler{
        services: services,
    }
}

func (h *MarketplaceHandler) CreateListing(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    
    var listing models.MarketplaceListing
    if err := c.BodyParser(&listing); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input format")
    }
    
    listing.UserID = userID
    
    listingID, err := h.services.Marketplace().CreateListing(c.Context(), &listing)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error creating listing")
    }
    
    return utils.SuccessResponse(c, fiber.Map{
        "id": listingID,
        "message": "Listing created successfully",
    })
}

func (h *MarketplaceHandler) GetListings(c *fiber.Ctx) error {
    filters := map[string]string{
        
        "category_id": c.Query("category_id"),
        "city":       c.Query("city"),
        "min_price":  c.Query("min_price"),
        "max_price":  c.Query("max_price"),
        // другие фильтры...
    }
    
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "20"))
    offset := (page - 1) * limit
    
    listings, total, err := h.services.Marketplace().GetListings(c.Context(), filters, limit, offset)
    if err != nil {
        log.Printf("Error getting listings: %v", err) // Добавьте это
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching listings")
    }

    log.Printf("Found %d listings", len(listings)) // И это
    
    return utils.SuccessResponse(c, fiber.Map{
        "data": listings,
        "meta": fiber.Map{
            "total": total,
            "page": page,
            "limit": limit,
        },
    })
}
func (h *MarketplaceHandler) AddToFavorites(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    listingID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
    }

    err = h.services.Marketplace().AddToFavorites(c.Context(), userID, listingID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error adding to favorites")
    }

    return utils.SuccessResponse(c, fiber.Map{
        "message": "Added to favorites successfully",
    })
}
func (h *MarketplaceHandler) RemoveFromFavorites(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    listingID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
    }

    err = h.services.Marketplace().RemoveFromFavorites(c.Context(), userID, listingID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error removing from favorites")
    }

    return utils.SuccessResponse(c, fiber.Map{
        "message": "Removed from favorites successfully",
    })
}
// GetListing - получение объявления по ID
func (h *MarketplaceHandler) GetListing(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
    }

    listing, err := h.services.Marketplace().GetListingByID(c.Context(), id)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusNotFound, "Listing not found")
    }

    return utils.SuccessResponse(c, listing)
}

// UpdateListing - обновление объявления
func (h *MarketplaceHandler) UpdateListing(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    listingID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
    }

    var listing models.MarketplaceListing
    if err := c.BodyParser(&listing); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid input format")
    }

    listing.ID = listingID
    listing.UserID = userID

    err = h.services.Marketplace().UpdateListing(c.Context(), &listing)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error updating listing")
    }

    return utils.SuccessResponse(c, fiber.Map{
        "message": "Listing updated successfully",
    })
}

// DeleteListing - удаление объявления
func (h *MarketplaceHandler) DeleteListing(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    listingID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
    }

    err = h.services.Marketplace().DeleteListing(c.Context(), listingID, userID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error deleting listing")
    }

    return utils.SuccessResponse(c, fiber.Map{
        "message": "Listing deleted successfully",
    })
}

// UploadImages - загрузка изображений для объявления
func (h *MarketplaceHandler) UploadImages(c *fiber.Ctx) error {
    listingID, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid listing ID")
    }

    form, err := c.MultipartForm()
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Error getting files")
    }

    files := form.File["images"]
    isMain := len(files) > 0

    var uploadedImages []models.MarketplaceImage
    for _, file := range files {
        fileName, err := h.services.Marketplace().ProcessImage(file)
        if err != nil {
            return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error processing image")
        }

        image := models.MarketplaceImage{
            ListingID:   listingID,
            FilePath:    fileName,
            FileName:    file.Filename,
            FileSize:    int(file.Size),
            ContentType: file.Header.Get("Content-Type"),
            IsMain:      isMain,
        }

        imageID, err := h.services.Marketplace().AddListingImage(c.Context(), &image)
        if err != nil {
            return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error saving image information")
        }

        image.ID = imageID
        uploadedImages = append(uploadedImages, image)
        isMain = false
    }

    return utils.SuccessResponse(c, uploadedImages)
}

// GetCategories - получение списка категорий
func (h *MarketplaceHandler) GetCategories(c *fiber.Ctx) error {
    categories, err := h.services.Marketplace().GetCategories(c.Context())
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Error fetching categories")
    }

    return utils.SuccessResponse(c, categories)
}