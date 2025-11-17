package grpc

import (
	"time"

	categoriespb "github.com/sveturs/listings/api/proto/categories/v1"
	listingspb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/domain"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DomainToProtoListing converts domain.Listing to protobuf Listing
func DomainToProtoListing(listing *domain.Listing) *listingspb.Listing {
	if listing == nil {
		return nil
	}

	pbListing := &listingspb.Listing{
		Id:               listing.ID,
		Uuid:             listing.UUID,
		UserId:           listing.UserID,
		Title:            listing.Title,
		Price:            listing.Price,
		Currency:         listing.Currency,
		CategoryId:       listing.CategoryID,
		Status:           listing.Status,
		Visibility:       listing.Visibility,
		Quantity:         listing.Quantity,
		ViewsCount:       listing.ViewsCount,
		FavoritesCount:   listing.FavoritesCount,
		CreatedAt:        listing.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        listing.UpdatedAt.Format(time.RFC3339),
		IsDeleted:        listing.IsDeleted,
		OriginalLanguage: listing.OriginalLanguage,
	}

	// Optional fields
	if listing.StorefrontID != nil {
		pbListing.StorefrontId = listing.StorefrontID
	}

	if listing.Description != nil {
		pbListing.Description = listing.Description
	}

	if listing.SKU != nil {
		pbListing.Sku = listing.SKU
	}

	if listing.PublishedAt != nil {
		publishedStr := listing.PublishedAt.Format(time.RFC3339)
		pbListing.PublishedAt = &publishedStr
	}

	if listing.DeletedAt != nil {
		deletedStr := listing.DeletedAt.Format(time.RFC3339)
		pbListing.DeletedAt = &deletedStr
	}

	// Translations
	pbListing.Translations = DomainTranslationsToProto(listing)

	// Relations
	if len(listing.Images) > 0 {
		pbListing.Images = make([]*listingspb.ListingImage, len(listing.Images))
		for i, img := range listing.Images {
			pbListing.Images[i] = DomainToProtoImage(img)
		}
	}

	if len(listing.Attributes) > 0 {
		pbListing.Attributes = make([]*listingspb.ListingAttribute, len(listing.Attributes))
		for i, attr := range listing.Attributes {
			pbListing.Attributes[i] = DomainToProtoAttribute(attr)
		}
	}

	if len(listing.Tags) > 0 {
		pbListing.Tags = listing.Tags
	}

	if listing.Location != nil {
		pbListing.Location = DomainToProtoLocation(listing.Location)
	}

	return pbListing
}

// DomainToProtoImage converts domain.ListingImage to protobuf ListingImage
func DomainToProtoImage(img *domain.ListingImage) *listingspb.ListingImage {
	if img == nil {
		return nil
	}

	pbImage := &listingspb.ListingImage{
		Id:           img.ID,
		ListingId:    img.ListingID,
		Url:          img.URL,
		DisplayOrder: img.DisplayOrder,
		IsPrimary:    img.IsPrimary,
		CreatedAt:    img.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    img.UpdatedAt.Format(time.RFC3339),
	}

	if img.StoragePath != nil {
		pbImage.StoragePath = img.StoragePath
	}

	if img.ThumbnailURL != nil {
		pbImage.ThumbnailUrl = img.ThumbnailURL
	}

	if img.Width != nil {
		pbImage.Width = img.Width
	}

	if img.Height != nil {
		pbImage.Height = img.Height
	}

	if img.FileSize != nil {
		pbImage.FileSize = img.FileSize
	}

	if img.MimeType != nil {
		pbImage.MimeType = img.MimeType
	}

	return pbImage
}

// DomainToProtoAttribute converts domain.ListingAttribute to protobuf ListingAttribute
func DomainToProtoAttribute(attr *domain.ListingAttribute) *listingspb.ListingAttribute {
	if attr == nil {
		return nil
	}

	return &listingspb.ListingAttribute{
		Id:             attr.ID,
		ListingId:      attr.ListingID,
		AttributeKey:   attr.AttributeKey,
		AttributeValue: attr.AttributeValue,
		CreatedAt:      attr.CreatedAt.Format(time.RFC3339),
	}
}

// DomainToProtoLocation converts domain.ListingLocation to protobuf ListingLocation
func DomainToProtoLocation(loc *domain.ListingLocation) *listingspb.ListingLocation {
	if loc == nil {
		return nil
	}

	pbLocation := &listingspb.ListingLocation{
		Id:        loc.ID,
		ListingId: loc.ListingID,
		CreatedAt: loc.CreatedAt.Format(time.RFC3339),
		UpdatedAt: loc.UpdatedAt.Format(time.RFC3339),
	}

	if loc.Country != nil {
		pbLocation.Country = loc.Country
	}

	if loc.City != nil {
		pbLocation.City = loc.City
	}

	if loc.PostalCode != nil {
		pbLocation.PostalCode = loc.PostalCode
	}

	if loc.AddressLine1 != nil {
		pbLocation.AddressLine1 = loc.AddressLine1
	}

	if loc.AddressLine2 != nil {
		pbLocation.AddressLine2 = loc.AddressLine2
	}

	if loc.Latitude != nil {
		pbLocation.Latitude = loc.Latitude
	}

	if loc.Longitude != nil {
		pbLocation.Longitude = loc.Longitude
	}

	return pbLocation
}

// ProtoToCreateListingInput converts CreateListingRequest to domain.CreateListingInput
func ProtoToCreateListingInput(req *listingspb.CreateListingRequest) *domain.CreateListingInput {
	if req == nil {
		return nil
	}

	input := &domain.CreateListingInput{
		UserID:     req.UserId,
		Title:      req.Title,
		Price:      req.Price,
		Currency:   req.Currency,
		CategoryID: req.CategoryId,
		Quantity:   req.Quantity,
	}

	// Automatically determine SourceType based on StorefrontID presence
	if req.StorefrontId != nil {
		input.StorefrontID = req.StorefrontId
		input.SourceType = "b2c"
	} else {
		input.SourceType = "c2c"
	}

	if req.Description != nil {
		input.Description = req.Description
	}

	if req.Sku != nil {
		input.SKU = req.Sku
	}

	// Translations
	if len(req.Translations) > 0 {
		input.Translations = ProtoTranslationsToMap(req.Translations)
	}

	if req.OriginalLanguage != nil && *req.OriginalLanguage != "" {
		input.OriginalLanguage = *req.OriginalLanguage
	}

	return input
}

// ProtoToUpdateListingInput converts UpdateListingRequest to domain.UpdateListingInput
func ProtoToUpdateListingInput(req *listingspb.UpdateListingRequest) *domain.UpdateListingInput {
	if req == nil {
		return nil
	}

	input := &domain.UpdateListingInput{}

	if req.Title != nil {
		input.Title = req.Title
	}

	if req.Description != nil {
		input.Description = req.Description
	}

	if req.Price != nil {
		input.Price = req.Price
	}

	if req.Quantity != nil {
		input.Quantity = req.Quantity
	}

	if req.Status != nil {
		input.Status = req.Status
	}

	return input
}

// ProtoToListListingsFilter converts ListListingsRequest to domain.ListListingsFilter
func ProtoToListListingsFilter(req *listingspb.ListListingsRequest) *domain.ListListingsFilter {
	if req == nil {
		return nil
	}

	filter := &domain.ListListingsFilter{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	if req.UserId != nil {
		filter.UserID = req.UserId
	}

	if req.StorefrontId != nil {
		filter.StorefrontID = req.StorefrontId
	}

	if req.CategoryId != nil {
		filter.CategoryID = req.CategoryId
	}

	if req.Status != nil {
		filter.Status = req.Status
	}

	if req.MinPrice != nil {
		filter.MinPrice = req.MinPrice
	}

	if req.MaxPrice != nil {
		filter.MaxPrice = req.MaxPrice
	}

	return filter
}

// ProtoToSearchListingsQuery converts SearchListingsRequest to domain.SearchListingsQuery
func ProtoToSearchListingsQuery(req *listingspb.SearchListingsRequest) *domain.SearchListingsQuery {
	if req == nil {
		return nil
	}

	query := &domain.SearchListingsQuery{
		Query:  req.Query,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	if req.CategoryId != nil {
		query.CategoryID = req.CategoryId
	}

	if req.MinPrice != nil {
		query.MinPrice = req.MinPrice
	}

	if req.MaxPrice != nil {
		query.MaxPrice = req.MaxPrice
	}

	return query
}

// DomainToProtoCategory converts domain.Category to protobuf Category
func DomainToProtoCategory(cat *domain.Category) *listingspb.Category {
	if cat == nil {
		return nil
	}

	pbCategory := &listingspb.Category{
		Id:           cat.ID,
		Name:         cat.Name,
		Slug:         cat.Slug,
		IsActive:     cat.IsActive,
		ListingCount: cat.ListingCount,
		SortOrder:    cat.SortOrder,
		Level:        cat.Level,
		Translations: make(map[string]string), // Placeholder for translations
		HasCustomUi:  cat.HasCustomUI,
		CreatedAt:    cat.CreatedAt.Format(time.RFC3339),
	}

	// Optional fields
	if cat.ParentID != nil {
		pbCategory.ParentId = cat.ParentID
	}

	if cat.Icon != nil {
		pbCategory.Icon = cat.Icon
	}

	if cat.Description != nil {
		pbCategory.Description = cat.Description
	}

	if cat.CustomUIComponent != nil {
		pbCategory.CustomUiComponent = cat.CustomUIComponent
	}

	return pbCategory
}

// DomainToProtoCategoryTree converts domain.CategoryTreeNode to protobuf CategoryTreeNode
func DomainToProtoCategoryTree(node *domain.CategoryTreeNode) *listingspb.CategoryTreeNode {
	if node == nil {
		return nil
	}

	pbNode := &listingspb.CategoryTreeNode{
		Id:            node.ID,
		Name:          node.Name,
		Slug:          node.Slug,
		Level:         node.Level,
		Path:          node.Path,
		ListingCount:  node.ListingCount,
		ChildrenCount: node.ChildrenCount,
		Translations:  make(map[string]string), // Placeholder for translations
		HasCustomUi:   node.HasCustomUI,
		CreatedAt:     node.CreatedAt,
	}

	// Optional fields
	if node.Icon != nil {
		pbNode.Icon = node.Icon
	}

	if node.ParentID != nil {
		pbNode.ParentId = node.ParentID
	}

	if node.CustomUIComponent != nil {
		pbNode.CustomUiComponent = node.CustomUIComponent
	}

	// Convert children recursively
	if len(node.Children) > 0 {
		pbNode.Children = make([]*listingspb.CategoryTreeNode, len(node.Children))
		for i, child := range node.Children {
			childCopy := child // Create copy to avoid pointer issues
			pbNode.Children[i] = DomainToProtoCategoryTree(&childCopy)
		}
	}

	return pbNode
}

// DomainToProtoVariant converts domain.ListingVariant to protobuf ListingVariant
func DomainToProtoVariant(variant *domain.ListingVariant) *listingspb.ListingVariant {
	if variant == nil {
		return nil
	}

	pbVariant := &listingspb.ListingVariant{
		Id:         variant.ID,
		ListingId:  variant.ListingID,
		Sku:        variant.SKU,
		IsActive:   variant.IsActive,
		Attributes: variant.Attributes,
	}

	// Optional fields
	if variant.Price != nil {
		pbVariant.Price = variant.Price
	}

	if variant.Stock != nil {
		pbVariant.Stock = variant.Stock
	}

	if variant.ImageURL != nil {
		pbVariant.ImageUrl = variant.ImageURL
	}

	if variant.CreatedAt != nil {
		createdStr := variant.CreatedAt.Format(time.RFC3339)
		pbVariant.CreatedAt = &createdStr
	}

	if variant.UpdatedAt != nil {
		updatedStr := variant.UpdatedAt.Format(time.RFC3339)
		pbVariant.UpdatedAt = &updatedStr
	}

	return pbVariant
}

// ProtoToVariantInput converts listingspb.VariantInput to domain.ListingVariant
func ProtoToVariantInput(input *listingspb.VariantInput, listingID int64) *domain.ListingVariant {
	if input == nil {
		return nil
	}

	variant := &domain.ListingVariant{
		ListingID:  listingID,
		SKU:        input.Sku,
		IsActive:   input.IsActive,
		Attributes: input.Attributes,
	}

	if input.Price != nil {
		variant.Price = input.Price
	}

	if input.Stock != nil {
		variant.Stock = input.Stock
	}

	if input.ImageUrl != nil {
		variant.ImageURL = input.ImageUrl
	}

	return variant
}

// ProtoToAddImageInput converts listingspb.AddImageRequest to domain.ListingImage
func ProtoToAddImageInput(req *listingspb.AddImageRequest) *domain.ListingImage {
	if req == nil {
		return nil
	}

	image := &domain.ListingImage{
		ListingID:    req.ListingId,
		URL:          req.Url,
		DisplayOrder: req.DisplayOrder,
		IsPrimary:    req.IsPrimary,
	}

	if req.StoragePath != nil {
		image.StoragePath = req.StoragePath
	}

	if req.ThumbnailUrl != nil {
		image.ThumbnailURL = req.ThumbnailUrl
	}

	if req.Width != nil {
		image.Width = req.Width
	}

	if req.Height != nil {
		image.Height = req.Height
	}

	if req.FileSize != nil {
		image.FileSize = req.FileSize
	}

	if req.MimeType != nil {
		image.MimeType = req.MimeType
	}

	return image
}

// ProductToProto converts domain.Product to listingspb.Product
func ProductToProto(p *domain.Product) *listingspb.Product {
	if p == nil {
		return nil
	}

	pbProduct := &listingspb.Product{
		Id:                    p.ID,
		StorefrontId:          p.StorefrontID,
		Name:                  p.Name,
		Description:           p.Description,
		Price:                 p.Price,
		Currency:              p.Currency,
		CategoryId:            p.CategoryID,
		StockQuantity:         p.StockQuantity,
		StockStatus:           p.StockStatus,
		IsActive:              p.IsActive,
		ViewCount:             p.ViewCount,
		SoldCount:             p.SoldCount,
		HasIndividualLocation: p.HasIndividualLocation,
		ShowOnMap:             p.ShowOnMap,
		HasVariants:           p.HasVariants,
	}

	// Convert optional fields
	if p.SKU != nil {
		pbProduct.Sku = p.SKU
	}

	if p.Barcode != nil {
		pbProduct.Barcode = p.Barcode
	}

	if p.IndividualAddress != nil {
		pbProduct.IndividualAddress = p.IndividualAddress
	}

	if p.IndividualLatitude != nil {
		pbProduct.IndividualLatitude = p.IndividualLatitude
	}

	if p.IndividualLongitude != nil {
		pbProduct.IndividualLongitude = p.IndividualLongitude
	}

	if p.LocationPrivacy != nil {
		pbProduct.LocationPrivacy = p.LocationPrivacy
	}

	// Convert timestamps
	if !p.CreatedAt.IsZero() {
		pbProduct.CreatedAt = timestamppb.New(p.CreatedAt)
	}

	if !p.UpdatedAt.IsZero() {
		pbProduct.UpdatedAt = timestamppb.New(p.UpdatedAt)
	}

	// Convert JSONB attributes to google.protobuf.Struct
	if p.Attributes != nil {
		if attrStruct, err := structpb.NewStruct(p.Attributes); err == nil {
			pbProduct.Attributes = attrStruct
		}
	}

	// Convert variants
	if len(p.Variants) > 0 {
		pbProduct.Variants = make([]*listingspb.ProductVariant, len(p.Variants))
		for i, v := range p.Variants {
			pbProduct.Variants[i] = ProductVariantToProto(&v)
		}
	}

	return pbProduct
}

// ProductVariantToProto converts domain.ProductVariant to listingspb.ProductVariant
func ProductVariantToProto(v *domain.ProductVariant) *listingspb.ProductVariant {
	if v == nil {
		return nil
	}

	pbVariant := &listingspb.ProductVariant{
		Id:            v.ID,
		ProductId:     v.ProductID,
		StockQuantity: v.StockQuantity,
		StockStatus:   v.StockStatus,
		IsActive:      v.IsActive,
		IsDefault:     v.IsDefault,
		ViewCount:     v.ViewCount,
		SoldCount:     v.SoldCount,
	}

	// Convert optional fields
	if v.SKU != nil {
		pbVariant.Sku = v.SKU
	}

	if v.Barcode != nil {
		pbVariant.Barcode = v.Barcode
	}

	if v.Price != nil {
		pbVariant.Price = v.Price
	}

	if v.CompareAtPrice != nil {
		pbVariant.CompareAtPrice = v.CompareAtPrice
	}

	if v.CostPrice != nil {
		pbVariant.CostPrice = v.CostPrice
	}

	if v.LowStockThreshold != nil {
		pbVariant.LowStockThreshold = v.LowStockThreshold
	}

	if v.Weight != nil {
		pbVariant.Weight = v.Weight
	}

	// Convert timestamps
	if !v.CreatedAt.IsZero() {
		pbVariant.CreatedAt = timestamppb.New(v.CreatedAt)
	}

	if !v.UpdatedAt.IsZero() {
		pbVariant.UpdatedAt = timestamppb.New(v.UpdatedAt)
	}

	// Convert JSONB fields to google.protobuf.Struct
	if v.VariantAttributes != nil {
		if attrStruct, err := structpb.NewStruct(v.VariantAttributes); err == nil {
			pbVariant.VariantAttributes = attrStruct
		}
	}

	if v.Dimensions != nil {
		if dimStruct, err := structpb.NewStruct(v.Dimensions); err == nil {
			pbVariant.Dimensions = dimStruct
		}
	}

	return pbVariant
}

// ProtoToCreateProductInput converts listingspb.CreateProductRequest to domain.CreateProductInput
func ProtoToCreateProductInput(req *listingspb.CreateProductRequest) *domain.CreateProductInput {
	if req == nil {
		return nil
	}

	input := &domain.CreateProductInput{
		StorefrontID:  req.StorefrontId,
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		Currency:      req.Currency,
		CategoryID:    req.CategoryId,
		StockQuantity: req.StockQuantity,
	}

	// Optional fields
	if req.Sku != nil {
		input.SKU = req.Sku
	}

	if req.Barcode != nil {
		input.Barcode = req.Barcode
	}

	// Convert google.protobuf.Struct to map[string]interface{}
	if req.Attributes != nil {
		input.Attributes = req.Attributes.AsMap()
	}

	// Location fields - bool fields are not optional in proto
	input.HasIndividualLocation = req.HasIndividualLocation
	input.ShowOnMap = req.ShowOnMap
	input.HasVariants = req.HasVariants

	if req.IndividualAddress != nil {
		input.IndividualAddress = req.IndividualAddress
	}

	if req.IndividualLatitude != nil {
		input.IndividualLatitude = req.IndividualLatitude
	}

	if req.IndividualLongitude != nil {
		input.IndividualLongitude = req.IndividualLongitude
	}

	if req.LocationPrivacy != nil {
		input.LocationPrivacy = req.LocationPrivacy
	}

	return input
}

// ProtoToUpdateProductInput converts listingspb.UpdateProductRequest to domain.UpdateProductInput
func ProtoToUpdateProductInput(req *listingspb.UpdateProductRequest) *domain.UpdateProductInput {
	if req == nil {
		return nil
	}

	input := &domain.UpdateProductInput{}

	// Optional fields - only set if provided
	if req.Name != nil {
		input.Name = req.Name
	}

	if req.Description != nil {
		input.Description = req.Description
	}

	if req.Price != nil {
		input.Price = req.Price
	}

	if req.IsActive != nil {
		input.IsActive = req.IsActive
	}

	// Convert google.protobuf.Struct to map[string]interface{}
	if req.Attributes != nil {
		input.Attributes = req.Attributes.AsMap()
	}

	// Location fields
	if req.HasIndividualLocation != nil {
		boolValue := *req.HasIndividualLocation
		input.HasIndividualLocation = &boolValue
	}

	if req.IndividualAddress != nil {
		input.IndividualAddress = req.IndividualAddress
	}

	if req.IndividualLatitude != nil {
		input.IndividualLatitude = req.IndividualLatitude
	}

	if req.IndividualLongitude != nil {
		input.IndividualLongitude = req.IndividualLongitude
	}

	if req.LocationPrivacy != nil {
		input.LocationPrivacy = req.LocationPrivacy
	}

	if req.ShowOnMap != nil {
		boolValue := *req.ShowOnMap
		input.ShowOnMap = &boolValue
	}

	return input
}

// ProtoToCreateVariantInput converts listingspb.CreateProductVariantRequest to domain.CreateVariantInput
func ProtoToCreateVariantInput(req *listingspb.CreateProductVariantRequest) *domain.CreateVariantInput {
	if req == nil {
		return nil
	}

	input := &domain.CreateVariantInput{
		ProductID:     req.ProductId,
		StockQuantity: req.StockQuantity,
		IsDefault:     req.IsDefault,
	}

	// Optional string fields
	if req.Sku != nil {
		input.SKU = req.Sku
	}
	if req.Barcode != nil {
		input.Barcode = req.Barcode
	}

	// Optional numeric fields
	if req.Price != nil {
		input.Price = req.Price
	}
	if req.CompareAtPrice != nil {
		input.CompareAtPrice = req.CompareAtPrice
	}
	if req.CostPrice != nil {
		input.CostPrice = req.CostPrice
	}
	if req.Weight != nil {
		input.Weight = req.Weight
	}
	if req.LowStockThreshold != nil {
		input.LowStockThreshold = req.LowStockThreshold
	}

	// Convert google.protobuf.Struct to map[string]interface{}
	if req.VariantAttributes != nil {
		input.VariantAttributes = req.VariantAttributes.AsMap()
	}

	if req.Dimensions != nil {
		input.Dimensions = req.Dimensions.AsMap()
	}

	return input
}

// ProtoToUpdateVariantInput converts listingspb.UpdateProductVariantRequest to domain.UpdateVariantInput
func ProtoToUpdateVariantInput(req *listingspb.UpdateProductVariantRequest) *domain.UpdateVariantInput {
	if req == nil {
		return nil
	}

	input := &domain.UpdateVariantInput{}

	// Optional string fields
	if req.Sku != nil {
		input.SKU = req.Sku
	}
	if req.Barcode != nil {
		input.Barcode = req.Barcode
	}

	// Optional numeric fields
	if req.Price != nil {
		input.Price = req.Price
	}
	if req.CompareAtPrice != nil {
		input.CompareAtPrice = req.CompareAtPrice
	}
	if req.CostPrice != nil {
		input.CostPrice = req.CostPrice
	}
	if req.Weight != nil {
		input.Weight = req.Weight
	}
	if req.StockQuantity != nil {
		input.StockQuantity = req.StockQuantity
	}
	if req.LowStockThreshold != nil {
		input.LowStockThreshold = req.LowStockThreshold
	}

	// Optional boolean fields
	if req.IsActive != nil {
		input.IsActive = req.IsActive
	}
	if req.IsDefault != nil {
		input.IsDefault = req.IsDefault
	}

	// Convert google.protobuf.Struct to map[string]interface{}
	if req.VariantAttributes != nil {
		input.VariantAttributes = req.VariantAttributes.AsMap()
	}

	if req.Dimensions != nil {
		input.Dimensions = req.Dimensions.AsMap()
	}

	return input
}

// ProtoToBulkVariantInputs converts slice of listingspb.ProductVariantInput to slice of domain.CreateVariantInput
func ProtoToBulkVariantInputs(productID int64, protoInputs []*listingspb.ProductVariantInput) []*domain.CreateVariantInput {
	if len(protoInputs) == 0 {
		return nil
	}

	inputs := make([]*domain.CreateVariantInput, 0, len(protoInputs))

	for _, protoInput := range protoInputs {
		if protoInput == nil {
			continue
		}

		input := &domain.CreateVariantInput{
			ProductID:     productID,
			StockQuantity: protoInput.StockQuantity,
			IsDefault:     protoInput.IsDefault,
		}

		// Optional string fields
		if protoInput.Sku != nil {
			input.SKU = protoInput.Sku
		}
		if protoInput.Barcode != nil {
			input.Barcode = protoInput.Barcode
		}

		// Optional numeric fields
		if protoInput.Price != nil {
			input.Price = protoInput.Price
		}
		if protoInput.CompareAtPrice != nil {
			input.CompareAtPrice = protoInput.CompareAtPrice
		}
		if protoInput.CostPrice != nil {
			input.CostPrice = protoInput.CostPrice
		}
		if protoInput.Weight != nil {
			input.Weight = protoInput.Weight
		}
		if protoInput.LowStockThreshold != nil {
			input.LowStockThreshold = protoInput.LowStockThreshold
		}

		// Convert google.protobuf.Struct to map[string]interface{}
		if protoInput.VariantAttributes != nil {
			input.VariantAttributes = protoInput.VariantAttributes.AsMap()
		}

		if protoInput.Dimensions != nil {
			input.Dimensions = protoInput.Dimensions.AsMap()
		}

		inputs = append(inputs, input)
	}

	return inputs
}

// ProtoToProductInput converts listingspb.ProductInput to domain.CreateProductInput
func ProtoToProductInput(p *listingspb.ProductInput, storefrontID int64) *domain.CreateProductInput {
	if p == nil {
		return nil
	}

	input := &domain.CreateProductInput{
		StorefrontID:          storefrontID,
		Name:                  p.Name,
		Description:           p.Description,
		Price:                 p.Price,
		Currency:              p.Currency,
		CategoryID:            p.CategoryId,
		StockQuantity:         p.StockQuantity,
		HasIndividualLocation: p.GetHasIndividualLocation(),
		ShowOnMap:             p.GetShowOnMap(),
	}

	// Handle optional fields
	if p.Sku != nil {
		input.SKU = p.Sku
	}

	if p.Barcode != nil {
		input.Barcode = p.Barcode
	}

	if p.Attributes != nil {
		input.Attributes = p.Attributes.AsMap()
	}

	if p.IndividualAddress != nil {
		input.IndividualAddress = p.IndividualAddress
	}

	if p.IndividualLatitude != nil {
		input.IndividualLatitude = p.IndividualLatitude
	}

	if p.IndividualLongitude != nil {
		input.IndividualLongitude = p.IndividualLongitude
	}

	if p.LocationPrivacy != nil {
		input.LocationPrivacy = p.LocationPrivacy
	}

	return input
}

// ProtoToBulkProductInputs converts repeated listingspb.ProductInput to domain.CreateProductInput slice
func ProtoToBulkProductInputs(protoInputs []*listingspb.ProductInput, storefrontID int64) []*domain.CreateProductInput {
	if len(protoInputs) == 0 {
		return nil
	}

	inputs := make([]*domain.CreateProductInput, 0, len(protoInputs))
	for _, protoInput := range protoInputs {
		if protoInput == nil {
			continue
		}
		inputs = append(inputs, ProtoToProductInput(protoInput, storefrontID))
	}

	return inputs
}

// ProtoToBulkUpdateInput converts listingspb.ProductUpdateInput to domain.BulkUpdateProductInput
func ProtoToBulkUpdateInput(req *listingspb.ProductUpdateInput) *domain.BulkUpdateProductInput {
	if req == nil {
		return nil
	}

	input := &domain.BulkUpdateProductInput{
		ProductID: req.ProductId,
	}

	// Optional fields - only set if provided
	if req.Name != nil {
		input.Name = req.Name
	}
	if req.Description != nil {
		input.Description = req.Description
	}
	if req.Price != nil {
		input.Price = req.Price
	}
	if req.Sku != nil {
		input.SKU = req.Sku
	}
	if req.Barcode != nil {
		input.Barcode = req.Barcode
	}
	if req.IsActive != nil {
		input.IsActive = req.IsActive
	}

	// Convert google.protobuf.Struct to map[string]interface{}
	if req.Attributes != nil {
		input.Attributes = req.Attributes.AsMap()
	}

	// Convert FieldMask to []string
	if req.UpdateMask != nil && len(req.UpdateMask.Paths) > 0 {
		input.UpdateMask = req.UpdateMask.Paths
	}

	return input
}

// ProtoToBulkUpdateInputs converts repeated listingspb.ProductUpdateInput to domain.BulkUpdateProductInput slice
func ProtoToBulkUpdateInputs(protoUpdates []*listingspb.ProductUpdateInput) []*domain.BulkUpdateProductInput {
	if len(protoUpdates) == 0 {
		return nil
	}

	inputs := make([]*domain.BulkUpdateProductInput, 0, len(protoUpdates))
	for _, protoUpdate := range protoUpdates {
		if protoUpdate == nil {
			continue
		}
		input := ProtoToBulkUpdateInput(protoUpdate)
		if input != nil {
			inputs = append(inputs, input)
		}
	}

	return inputs
}

// StorefrontToProto converts domain.Storefront to listingspb.Storefront
func StorefrontToProto(sf *domain.Storefront) *listingspb.Storefront {
	if sf == nil {
		return nil
	}

	// Convert *string to string for Country (required in proto)
	country := ""
	if sf.Country != nil {
		country = *sf.Country
	}

	protoSF := &listingspb.Storefront{
		Id:             sf.ID,
		UserId:         sf.UserID,
		Slug:           sf.Slug,
		Name:           sf.Name,
		Country:        country,
		IsActive:       sf.IsActive,
		IsVerified:     sf.IsVerified,
		Rating:         sf.Rating,
		ReviewsCount:   sf.ReviewsCount,
		ProductsCount:  sf.ProductsCount,
		SalesCount:     sf.SalesCount,
		ViewsCount:     sf.ViewsCount,
		FollowersCount: sf.FollowersCount,
	}

	// Optional fields
	if sf.Description != nil {
		protoSF.Description = sf.Description
	}

	if sf.LogoURL != nil {
		protoSF.LogoUrl = sf.LogoURL
	}

	if sf.BannerURL != nil {
		protoSF.BannerUrl = sf.BannerURL
	}

	if sf.Phone != nil {
		protoSF.Phone = sf.Phone
	}

	if sf.Email != nil {
		protoSF.Email = sf.Email
	}

	if sf.Website != nil {
		protoSF.Website = sf.Website
	}

	if sf.Address != nil {
		protoSF.Address = sf.Address
	}

	if sf.City != nil {
		protoSF.City = sf.City
	}

	if sf.PostalCode != nil {
		protoSF.PostalCode = sf.PostalCode
	}

	if sf.Latitude != nil {
		protoSF.Latitude = sf.Latitude
	}

	if sf.Longitude != nil {
		protoSF.Longitude = sf.Longitude
	}

	// Timestamps
	if !sf.CreatedAt.IsZero() {
		protoSF.CreatedAt = sf.CreatedAt.Format(time.RFC3339)
	}

	if !sf.UpdatedAt.IsZero() {
		protoSF.UpdatedAt = sf.UpdatedAt.Format(time.RFC3339)
	}

	return protoSF
}

// DomainTranslationsToProto converts domain listing translation fields to proto map
func DomainTranslationsToProto(listing *domain.Listing) map[string]*listingspb.ListingFieldTranslations {
	if listing == nil {
		return nil
	}

	result := make(map[string]*listingspb.ListingFieldTranslations)

	// Collect all unique language codes
	languages := make(map[string]bool)
	for lang := range listing.TitleTranslations {
		languages[lang] = true
	}
	for lang := range listing.DescriptionTranslations {
		languages[lang] = true
	}
	for lang := range listing.LocationTranslations {
		languages[lang] = true
	}
	for lang := range listing.CityTranslations {
		languages[lang] = true
	}
	for lang := range listing.CountryTranslations {
		languages[lang] = true
	}

	// Build translations for each language
	for lang := range languages {
		fields := &listingspb.ListingFieldTranslations{}

		if title, ok := listing.TitleTranslations[lang]; ok && title != "" {
			fields.Title = &title
		}
		if desc, ok := listing.DescriptionTranslations[lang]; ok && desc != "" {
			fields.Description = &desc
		}
		if loc, ok := listing.LocationTranslations[lang]; ok && loc != "" {
			fields.Location = &loc
		}
		if city, ok := listing.CityTranslations[lang]; ok && city != "" {
			fields.City = &city
		}
		if country, ok := listing.CountryTranslations[lang]; ok && country != "" {
			fields.Country = &country
		}

		result[lang] = fields
	}

	return result
}

// ProtoTranslationsToMap converts proto translations to domain format map[string]map[string]string
func ProtoTranslationsToMap(protoTranslations map[string]*listingspb.ListingFieldTranslations) map[string]map[string]string {
	if protoTranslations == nil {
		return nil
	}

	result := make(map[string]map[string]string)
	for lang, fields := range protoTranslations {
		if fields == nil {
			continue
		}

		langFields := make(map[string]string)
		if fields.Title != nil && *fields.Title != "" {
			langFields["title"] = *fields.Title
		}
		if fields.Description != nil && *fields.Description != "" {
			langFields["description"] = *fields.Description
		}
		if fields.Location != nil && *fields.Location != "" {
			langFields["location"] = *fields.Location
		}
		if fields.City != nil && *fields.City != "" {
			langFields["city"] = *fields.City
		}
		if fields.Country != nil && *fields.Country != "" {
			langFields["country"] = *fields.Country
		}

		if len(langFields) > 0 {
			result[lang] = langFields
		}
	}

	return result
}

// ApplyTranslation applies translations to a listing for a specific language
func ApplyTranslation(listing *domain.Listing, lang string) {
	if listing == nil || lang == "" || lang == listing.OriginalLanguage {
		return
	}

	// Apply title translation
	if title, ok := listing.TitleTranslations[lang]; ok && title != "" {
		listing.Title = title
	}

	// Apply description translation
	if desc, ok := listing.DescriptionTranslations[lang]; ok && desc != "" {
		listing.Description = &desc
	}

	// Apply location translation if Location is loaded
	if listing.Location != nil {
		if loc, ok := listing.LocationTranslations[lang]; ok && loc != "" {
			listing.Location.Country = &loc
		}
		if city, ok := listing.CityTranslations[lang]; ok && city != "" {
			listing.Location.City = &city
		}
		if country, ok := listing.CountryTranslations[lang]; ok && country != "" {
			listing.Location.Country = &country
		}
	}
}

// ========================================================================================
// CategoryService proto converters (separate from ListingsService proto converters above)
// ========================================================================================

// DomainToCategoryServiceProtoCategory converts domain.Category to categoriespb.Category
func DomainToCategoryServiceProtoCategory(cat *domain.Category) *categoriespb.Category {
	if cat == nil {
		return nil
	}

	pbCat := &categoriespb.Category{
		Id:          int32(cat.ID),
		Name:        cat.Name,
		Slug:        cat.Slug,
		CreatedAt:   timestamppb.New(cat.CreatedAt),
		HasCustomUi: cat.HasCustomUI,
		SortOrder:   cat.SortOrder,
		Level:       cat.Level,
		Count:       cat.ListingCount,
		IsActive:    cat.IsActive,
	}

	// Optional fields
	if cat.ParentID != nil {
		parentID := int32(*cat.ParentID)
		pbCat.ParentId = &parentID
	}
	if cat.Icon != nil {
		pbCat.Icon = cat.Icon
	}
	if cat.CustomUIComponent != nil {
		pbCat.CustomUiComponent = cat.CustomUIComponent
	}
	if cat.ExternalID != nil {
		pbCat.ExternalId = cat.ExternalID
	}
	if cat.Description != nil {
		pbCat.Description = cat.Description
	}
	if cat.SEOTitle != nil {
		pbCat.SeoTitle = cat.SEOTitle
	}
	if cat.SEODescription != nil {
		pbCat.SeoDescription = cat.SEODescription
	}
	if cat.SEOKeywords != nil {
		pbCat.SeoKeywords = cat.SEOKeywords
	}
	if cat.TitleEn != nil {
		pbCat.TitleEn = cat.TitleEn
	}
	if cat.TitleRu != nil {
		pbCat.TitleRu = cat.TitleRu
	}
	if cat.TitleSr != nil {
		pbCat.TitleSr = cat.TitleSr
	}

	return pbCat
}

// DomainToCategoryServiceProtoCategoryTree converts domain.CategoryTreeNode to categoriespb.CategoryTree
func DomainToCategoryServiceProtoCategoryTree(node *domain.CategoryTreeNode) *categoriespb.CategoryTree {
	if node == nil {
		return nil
	}

	// Convert CategoryTreeNode to Category struct first
	cat := &domain.Category{
		ID:                node.ID,
		Name:              node.Name,
		Slug:              node.Slug,
		ParentID:          node.ParentID,
		Icon:              node.Icon,
		IsActive:          true, // TreeNode doesn't have IsActive, assume true
		ListingCount:      node.ListingCount,
		SortOrder:         0, // TreeNode doesn't have SortOrder
		Level:             node.Level,
		HasCustomUI:       node.HasCustomUI,
		CustomUIComponent: node.CustomUIComponent,
	}

	pbTree := &categoriespb.CategoryTree{
		Category: DomainToCategoryServiceProtoCategory(cat),
	}

	// Convert children
	if len(node.Children) > 0 {
		pbTree.Subcategories = make([]*categoriespb.CategoryTree, len(node.Children))
		for i := range node.Children {
			pbTree.Subcategories[i] = DomainToCategoryServiceProtoCategoryTree(&node.Children[i])
		}
	}

	return pbTree
}

// ProtoToCategoryServiceCreateDomain converts CreateCategoryRequest to domain.Category
func ProtoToCategoryServiceCreateDomain(req *categoriespb.CreateCategoryRequest) *domain.Category {
	cat := &domain.Category{
		Name:      req.Name,
		Slug:      req.Slug,
		SortOrder: req.SortOrder,
		IsActive:  true, // New categories are active by default
	}

	// Optional fields
	if req.ParentId != nil {
		parentID := int64(*req.ParentId)
		cat.ParentID = &parentID
	}
	if req.Icon != nil {
		cat.Icon = req.Icon
	}
	if req.Description != nil {
		cat.Description = req.Description
	}
	if req.CustomUiComponent != nil {
		cat.CustomUIComponent = req.CustomUiComponent
	}
	if req.SeoTitle != nil {
		cat.SEOTitle = req.SeoTitle
	}
	if req.SeoDescription != nil {
		cat.SEODescription = req.SeoDescription
	}
	if req.SeoKeywords != nil {
		cat.SEOKeywords = req.SeoKeywords
	}
	if req.TitleEn != nil {
		cat.TitleEn = req.TitleEn
	}
	if req.TitleRu != nil {
		cat.TitleRu = req.TitleRu
	}
	if req.TitleSr != nil {
		cat.TitleSr = req.TitleSr
	}

	return cat
}

// ProtoToCategoryServiceUpdateDomain converts UpdateCategoryRequest to domain.Category
func ProtoToCategoryServiceUpdateDomain(req *categoriespb.UpdateCategoryRequest) *domain.Category {
	cat := &domain.Category{
		ID: int64(req.Id),
	}

	// Optional fields - only set if provided
	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.Slug != nil {
		cat.Slug = *req.Slug
	}
	if req.ParentId != nil {
		parentID := int64(*req.ParentId)
		cat.ParentID = &parentID
	}
	if req.Icon != nil {
		cat.Icon = req.Icon
	}
	if req.Description != nil {
		cat.Description = req.Description
	}
	if req.CustomUiComponent != nil {
		cat.CustomUIComponent = req.CustomUiComponent
	}
	if req.SortOrder != nil {
		cat.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		cat.IsActive = *req.IsActive
	}
	if req.SeoTitle != nil {
		cat.SEOTitle = req.SeoTitle
	}
	if req.SeoDescription != nil {
		cat.SEODescription = req.SeoDescription
	}
	if req.SeoKeywords != nil {
		cat.SEOKeywords = req.SeoKeywords
	}
	if req.TitleEn != nil {
		cat.TitleEn = req.TitleEn
	}
	if req.TitleRu != nil {
		cat.TitleRu = req.TitleRu
	}
	if req.TitleSr != nil {
		cat.TitleSr = req.TitleSr
	}

	return cat
}
