package listings

import (
	"backend/internal/domain/models"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

// ProtoToUnifiedListing конвертирует proto Listing в UnifiedListing
func ProtoToUnifiedListing(protoListing *pb.Listing) (*models.UnifiedListing, error) {
	if protoListing == nil {
		return nil, nil
	}

	// Базовая информация
	description := ""
	if protoListing.Description != nil {
		description = *protoListing.Description
	}

	unified := &models.UnifiedListing{
		ID:          int(protoListing.Id),
		UserID:      int(protoListing.UserId),
		Title:       protoListing.Title,
		Description: description,
		Price:       protoListing.Price,
		CategoryID:  int(protoListing.CategoryId),
		Status:      protoListing.Status,
		ViewsCount:  int(protoListing.ViewsCount),
		// SourceType всегда "c2c" для listings из микросервиса
		SourceType: "c2c",
	}

	// Storefront ID (опционально)
	if protoListing.StorefrontId != nil {
		storefrontID := int(*protoListing.StorefrontId)
		unified.StorefrontID = &storefrontID
	}

	// Images - конвертируем в UnifiedImage
	if len(protoListing.Images) > 0 {
		unified.Images = make([]models.UnifiedImage, 0, len(protoListing.Images))
		for _, img := range protoListing.Images {
			thumbnailURL := ""
			if img.ThumbnailUrl != nil {
				thumbnailURL = *img.ThumbnailUrl
			}

			unified.Images = append(unified.Images, models.UnifiedImage{
				ID:           int(img.Id),
				URL:          img.Url,
				ThumbnailURL: thumbnailURL,
				IsMain:       img.IsPrimary,
				DisplayOrder: int(img.DisplayOrder),
			})
		}
	}

	// Location - маппинг в поля верхнего уровня
	if protoListing.Location != nil {
		loc := protoListing.Location
		if loc.Country != nil {
			unified.Country = *loc.Country
		}
		if loc.City != nil {
			unified.City = *loc.City
		}
		if loc.AddressLine1 != nil {
			unified.Location = *loc.AddressLine1
		}
		if loc.Latitude != nil {
			lat := *loc.Latitude
			unified.Latitude = &lat
		}
		if loc.Longitude != nil {
			lon := *loc.Longitude
			unified.Longitude = &lon
		}
	}

	return unified, nil
}

// UnifiedToProtoCreateRequest конвертирует UnifiedListing в CreateListingRequest
func UnifiedToProtoCreateRequest(unified *models.UnifiedListing) (*pb.CreateListingRequest, error) {
	if unified == nil {
		return nil, nil
	}

	req := &pb.CreateListingRequest{
		UserId:     int64(unified.UserID),
		Title:      unified.Title,
		Price:      unified.Price,
		Currency:   "RSD", // Default currency for monolith
		CategoryId: int64(unified.CategoryID),
		Quantity:   1, // Default quantity
	}

	// Optional fields
	if unified.Description != "" {
		req.Description = &unified.Description
	}

	if unified.StorefrontID != nil {
		storefrontID := int64(*unified.StorefrontID)
		req.StorefrontId = &storefrontID
	}

	// SKU не поддерживается в монолите через UnifiedListing

	return req, nil
}

// UnifiedToProtoUpdateRequest конвертирует UnifiedListing в UpdateListingRequest
func UnifiedToProtoUpdateRequest(unified *models.UnifiedListing) (*pb.UpdateListingRequest, error) {
	if unified == nil {
		return nil, nil
	}

	req := &pb.UpdateListingRequest{
		Id:     int64(unified.ID),
		UserId: int64(unified.UserID),
	}

	// Optional fields - обновляем только если заданы
	if unified.Title != "" {
		req.Title = &unified.Title
	}

	if unified.Description != "" {
		req.Description = &unified.Description
	}

	if unified.Price > 0 {
		req.Price = &unified.Price
	}

	// Quantity не поддерживается в монолите через UnifiedListing

	if unified.Status != "" {
		req.Status = &unified.Status
	}

	return req, nil
}

// Helper functions для работы с optional полями

func stringPtrFromOptional(opt *string) *string {
	if opt == nil {
		return nil
	}
	return opt
}

func float64PtrFromOptional(opt *float64) *float64 {
	if opt == nil {
		return nil
	}
	return opt
}
