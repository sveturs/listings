// Package adapters содержит адаптеры для конвертации между различными представлениями listings
package adapters

import (
	"encoding/json"
	"time"

	"backend/internal/domain/models"
)

// C2CToUnified конвертирует C2C listing в унифицированную модель
func C2CToUnified(c2c *models.MarketplaceListing) (*models.UnifiedListing, error) {
	if c2c == nil {
		return nil, nil
	}

	unified := &models.UnifiedListing{
		ID:           c2c.ID,
		SourceType:   "c2c",
		UserID:       c2c.UserID,
		CategoryID:   c2c.CategoryID,
		Title:        c2c.Title,
		Description:  c2c.Description,
		Price:        c2c.Price,
		Condition:    c2c.Condition,
		Status:       c2c.Status,
		Location:     c2c.Location,
		Latitude:     c2c.Latitude,
		Longitude:    c2c.Longitude,
		City:         c2c.City,
		Country:      c2c.Country,
		ViewsCount:   c2c.ViewsCount,
		ShowOnMap:    c2c.ShowOnMap,
		OriginalLang: c2c.OriginalLanguage,
		CreatedAt:    c2c.CreatedAt,
		UpdatedAt:    c2c.UpdatedAt,
		StorefrontID: c2c.StorefrontID, // может быть nil для обычных C2C
	}

	// Metadata
	if c2c.ExternalID != "" || c2c.Metadata != nil {
		unified.Metadata = make(map[string]interface{})
		if c2c.ExternalID != "" {
			unified.Metadata["external_id"] = c2c.ExternalID
		}
		if c2c.Metadata != nil {
			// Merge existing metadata
			for k, v := range c2c.Metadata {
				unified.Metadata[k] = v
			}
		}
	}

	// Convert images (MarketplaceImage type)
	if len(c2c.Images) > 0 {
		unified.Images = make([]models.UnifiedImage, len(c2c.Images))
		for i, img := range c2c.Images {
			unified.Images[i] = models.UnifiedImage{
				ID:           img.ID,
				URL:          img.PublicURL,
				ThumbnailURL: img.ThumbnailURL,
				IsMain:       img.IsMain,
				DisplayOrder: img.DisplayOrder,
			}
		}
	}

	// Convert images to JSON for DB storage
	if len(unified.Images) > 0 {
		imagesJSON, err := json.Marshal(unified.Images)
		if err != nil {
			return nil, err
		}
		unified.ImagesJSON = imagesJSON
	}

	return unified, nil
}

// B2CToUnified конвертирует B2C product (StorefrontProduct) в унифицированную модель
func B2CToUnified(b2c *models.StorefrontProduct, storefront *models.Storefront) (*models.UnifiedListing, error) {
	if b2c == nil {
		return nil, nil
	}

	// Определяем location из product или storefront
	var location string
	var latitude, longitude *float64
	var city, country string

	if b2c.HasIndividualLocation && b2c.IndividualAddress != nil {
		location = *b2c.IndividualAddress
		latitude = b2c.IndividualLatitude
		longitude = b2c.IndividualLongitude
	} else if storefront != nil {
		if storefront.Address != nil {
			location = *storefront.Address
		}
		if storefront.Latitude != nil {
			latitude = storefront.Latitude
		}
		if storefront.Longitude != nil {
			longitude = storefront.Longitude
		}
		if storefront.City != nil {
			city = *storefront.City
		}
		if storefront.Country != nil {
			country = *storefront.Country
		}
	}

	// Определяем status
	status := "inactive"
	if b2c.IsActive {
		status = "active"
	}

	// User ID из storefront
	var userID int
	if storefront != nil {
		userID = storefront.UserID
	}

	unified := &models.UnifiedListing{
		ID:           b2c.ID,
		SourceType:   "b2c",
		UserID:       userID,
		CategoryID:   b2c.CategoryID,
		Title:        b2c.Name,
		Description:  b2c.Description,
		Price:        b2c.Price,
		Condition:    "new", // B2C всегда новое
		Status:       status,
		Location:     location,
		Latitude:     latitude,
		Longitude:    longitude,
		City:         city,
		Country:      country,
		ViewsCount:   b2c.ViewCount,
		ShowOnMap:    b2c.ShowOnMap,
		OriginalLang: "sr", // Default для B2C
		CreatedAt:    b2c.CreatedAt,
		UpdatedAt:    b2c.UpdatedAt,
		StorefrontID: &b2c.StorefrontID,
	}

	// Metadata для B2C включает больше информации
	unified.Metadata = map[string]interface{}{
		"source":         "storefront",
		"storefront_id":  b2c.StorefrontID,
		"stock_quantity": b2c.StockQuantity,
		"stock_status":   b2c.StockStatus,
		"currency":       b2c.Currency,
		"has_variants":   b2c.HasVariants,
		"sold_count":     b2c.SoldCount,
	}

	if b2c.SKU != nil {
		unified.Metadata["sku"] = *b2c.SKU
	}
	if b2c.Barcode != nil {
		unified.Metadata["barcode"] = *b2c.Barcode
	}

	// Include attributes если есть (JSONB type is already map[string]interface{})
	if len(b2c.Attributes) > 0 {
		unified.Metadata["attributes"] = b2c.Attributes
	}

	// Convert images
	if len(b2c.Images) > 0 {
		unified.Images = make([]models.UnifiedImage, len(b2c.Images))
		for i, img := range b2c.Images {
			unified.Images[i] = models.UnifiedImage{
				ID:           img.ID,
				URL:          img.ImageURL,
				ThumbnailURL: img.ThumbnailURL,
				IsMain:       img.IsDefault,
				DisplayOrder: img.DisplayOrder,
			}
		}
	}

	// Convert images to JSON for DB storage
	if len(unified.Images) > 0 {
		imagesJSON, err := json.Marshal(unified.Images)
		if err != nil {
			return nil, err
		}
		unified.ImagesJSON = imagesJSON
	}

	return unified, nil
}

// UnifiedToC2C конвертирует унифицированную модель обратно в C2C listing
// Используется для обратной совместимости и миграции данных
func UnifiedToC2C(unified *models.UnifiedListing) (*models.MarketplaceListing, error) {
	if unified == nil || unified.SourceType != "c2c" {
		return nil, nil
	}

	c2c := &models.MarketplaceListing{
		ID:               unified.ID,
		UserID:           unified.UserID,
		CategoryID:       unified.CategoryID,
		Title:            unified.Title,
		Description:      unified.Description,
		Price:            unified.Price,
		Condition:        unified.Condition,
		Status:           unified.Status,
		Location:         unified.Location,
		Latitude:         unified.Latitude,
		Longitude:        unified.Longitude,
		City:             unified.City,
		Country:          unified.Country,
		ViewsCount:       unified.ViewsCount,
		ShowOnMap:        unified.ShowOnMap,
		OriginalLanguage: unified.OriginalLang,
		CreatedAt:        unified.CreatedAt,
		UpdatedAt:        unified.UpdatedAt,
		StorefrontID:     unified.StorefrontID,
	}

	// Extract external_id from metadata
	if unified.Metadata != nil {
		if externalID, ok := unified.Metadata["external_id"].(string); ok {
			c2c.ExternalID = externalID
		}

		// Convert metadata back to map (already correct type)
		c2c.Metadata = unified.Metadata
	}

	// Convert images (to MarketplaceImage type)
	if len(unified.Images) > 0 {
		c2c.Images = make([]models.MarketplaceImage, len(unified.Images))
		for i, img := range unified.Images {
			c2c.Images[i] = models.MarketplaceImage{
				ID:           img.ID,
				ListingID:    unified.ID,
				PublicURL:    img.URL,
				ThumbnailURL: img.ThumbnailURL,
				IsMain:       img.IsMain,
				DisplayOrder: img.DisplayOrder,
			}
		}
	}

	return c2c, nil
}

// UnifiedToB2C конвертирует унифицированную модель обратно в B2C product
// Используется для обратной совместимости и миграции данных
func UnifiedToB2C(unified *models.UnifiedListing) (*models.StorefrontProduct, error) {
	if unified == nil || unified.SourceType != "b2c" {
		return nil, nil
	}

	if unified.StorefrontID == nil {
		// B2C должен иметь storefront_id
		return nil, nil
	}

	// Determine if location is individual or from storefront
	hasIndividualLocation := false
	var individualAddress *string
	var individualLatitude, individualLongitude *float64

	if unified.Metadata != nil {
		if _, ok := unified.Metadata["individual_location"]; ok {
			hasIndividualLocation = true
			if unified.Location != "" {
				individualAddress = &unified.Location
			}
			individualLatitude = unified.Latitude
			individualLongitude = unified.Longitude
		}
	}

	b2c := &models.StorefrontProduct{
		ID:                    unified.ID,
		StorefrontID:          *unified.StorefrontID,
		Name:                  unified.Title,
		Description:           unified.Description,
		Price:                 unified.Price,
		CategoryID:            unified.CategoryID,
		StockQuantity:         0, // Default
		StockStatus:           "in_stock",
		IsActive:              unified.Status == "active",
		ViewCount:             unified.ViewsCount,
		SoldCount:             0,
		CreatedAt:             unified.CreatedAt,
		UpdatedAt:             unified.UpdatedAt,
		HasIndividualLocation: hasIndividualLocation,
		IndividualAddress:     individualAddress,
		IndividualLatitude:    individualLatitude,
		IndividualLongitude:   individualLongitude,
		ShowOnMap:             unified.ShowOnMap,
		HasVariants:           false,
		Currency:              "USD", // Default
	}

	// Extract metadata
	if unified.Metadata != nil {
		if stockQty, ok := unified.Metadata["stock_quantity"].(float64); ok {
			b2c.StockQuantity = int(stockQty)
		}
		if stockStatus, ok := unified.Metadata["stock_status"].(string); ok {
			b2c.StockStatus = stockStatus
		}
		if currency, ok := unified.Metadata["currency"].(string); ok {
			b2c.Currency = currency
		}
		if hasVariants, ok := unified.Metadata["has_variants"].(bool); ok {
			b2c.HasVariants = hasVariants
		}
		if soldCount, ok := unified.Metadata["sold_count"].(float64); ok {
			b2c.SoldCount = int(soldCount)
		}
		if sku, ok := unified.Metadata["sku"].(string); ok {
			b2c.SKU = &sku
		}
		if barcode, ok := unified.Metadata["barcode"].(string); ok {
			b2c.Barcode = &barcode
		}

		// Extract attributes (JSONB is already map[string]interface{})
		if attrs, ok := unified.Metadata["attributes"].(map[string]interface{}); ok {
			b2c.Attributes = models.JSONB(attrs)
		}
	}

	// Convert images
	if len(unified.Images) > 0 {
		b2c.Images = make([]models.StorefrontProductImage, len(unified.Images))
		for i, img := range unified.Images {
			b2c.Images[i] = models.StorefrontProductImage{
				ID:                  img.ID,
				StorefrontProductID: unified.ID,
				ImageURL:            img.URL,
				ThumbnailURL:        img.ThumbnailURL,
				IsDefault:           img.IsMain,
				DisplayOrder:        img.DisplayOrder,
				CreatedAt:           time.Now(),
			}
		}
	}

	return b2c, nil
}

// BatchC2CToUnified конвертирует batch C2C listings в unified
func BatchC2CToUnified(c2cListings []*models.MarketplaceListing) ([]*models.UnifiedListing, error) {
	if len(c2cListings) == 0 {
		return nil, nil
	}

	unified := make([]*models.UnifiedListing, 0, len(c2cListings))
	for _, c2c := range c2cListings {
		u, err := C2CToUnified(c2c)
		if err != nil {
			return nil, err
		}
		if u != nil {
			unified = append(unified, u)
		}
	}

	return unified, nil
}

// BatchB2CToUnified конвертирует batch B2C products в unified
func BatchB2CToUnified(b2cProducts []*models.StorefrontProduct, storefronts map[int]*models.Storefront) ([]*models.UnifiedListing, error) {
	if len(b2cProducts) == 0 {
		return nil, nil
	}

	unified := make([]*models.UnifiedListing, 0, len(b2cProducts))
	for _, b2c := range b2cProducts {
		var storefront *models.Storefront
		if storefronts != nil {
			storefront = storefronts[b2c.StorefrontID]
		}

		u, err := B2CToUnified(b2c, storefront)
		if err != nil {
			return nil, err
		}
		if u != nil {
			unified = append(unified, u)
		}
	}

	return unified, nil
}
