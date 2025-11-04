package grpc

import (
	"time"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/domain"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DomainToProtoListing converts domain.Listing to protobuf Listing
func DomainToProtoListing(listing *domain.Listing) *pb.Listing {
	if listing == nil {
		return nil
	}

	pbListing := &pb.Listing{
		Id:             listing.ID,
		Uuid:           listing.UUID,
		UserId:         listing.UserID,
		Title:          listing.Title,
		Price:          listing.Price,
		Currency:       listing.Currency,
		CategoryId:     listing.CategoryID,
		Status:         listing.Status,
		Visibility:     listing.Visibility,
		Quantity:       listing.Quantity,
		ViewsCount:     listing.ViewsCount,
		FavoritesCount: listing.FavoritesCount,
		CreatedAt:      listing.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      listing.UpdatedAt.Format(time.RFC3339),
		IsDeleted:      listing.IsDeleted,
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

	// Relations
	if len(listing.Images) > 0 {
		pbListing.Images = make([]*pb.ListingImage, len(listing.Images))
		for i, img := range listing.Images {
			pbListing.Images[i] = DomainToProtoImage(img)
		}
	}

	if len(listing.Attributes) > 0 {
		pbListing.Attributes = make([]*pb.ListingAttribute, len(listing.Attributes))
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
func DomainToProtoImage(img *domain.ListingImage) *pb.ListingImage {
	if img == nil {
		return nil
	}

	pbImage := &pb.ListingImage{
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
func DomainToProtoAttribute(attr *domain.ListingAttribute) *pb.ListingAttribute {
	if attr == nil {
		return nil
	}

	return &pb.ListingAttribute{
		Id:             attr.ID,
		ListingId:      attr.ListingID,
		AttributeKey:   attr.AttributeKey,
		AttributeValue: attr.AttributeValue,
		CreatedAt:      attr.CreatedAt.Format(time.RFC3339),
	}
}

// DomainToProtoLocation converts domain.ListingLocation to protobuf ListingLocation
func DomainToProtoLocation(loc *domain.ListingLocation) *pb.ListingLocation {
	if loc == nil {
		return nil
	}

	pbLocation := &pb.ListingLocation{
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
func ProtoToCreateListingInput(req *pb.CreateListingRequest) *domain.CreateListingInput {
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

	if req.StorefrontId != nil {
		input.StorefrontID = req.StorefrontId
	}

	if req.Description != nil {
		input.Description = req.Description
	}

	if req.Sku != nil {
		input.SKU = req.Sku
	}

	return input
}

// ProtoToUpdateListingInput converts UpdateListingRequest to domain.UpdateListingInput
func ProtoToUpdateListingInput(req *pb.UpdateListingRequest) *domain.UpdateListingInput {
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
func ProtoToListListingsFilter(req *pb.ListListingsRequest) *domain.ListListingsFilter {
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
func ProtoToSearchListingsQuery(req *pb.SearchListingsRequest) *domain.SearchListingsQuery {
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
func DomainToProtoCategory(cat *domain.Category) *pb.Category {
	if cat == nil {
		return nil
	}

	pbCategory := &pb.Category{
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
func DomainToProtoCategoryTree(node *domain.CategoryTreeNode) *pb.CategoryTreeNode {
	if node == nil {
		return nil
	}

	pbNode := &pb.CategoryTreeNode{
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
		pbNode.Children = make([]*pb.CategoryTreeNode, len(node.Children))
		for i, child := range node.Children {
			childCopy := child // Create copy to avoid pointer issues
			pbNode.Children[i] = DomainToProtoCategoryTree(&childCopy)
		}
	}

	return pbNode
}

// DomainToProtoVariant converts domain.ListingVariant to protobuf ListingVariant
func DomainToProtoVariant(variant *domain.ListingVariant) *pb.ListingVariant {
	if variant == nil {
		return nil
	}

	pbVariant := &pb.ListingVariant{
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

// ProtoToVariantInput converts pb.VariantInput to domain.ListingVariant
func ProtoToVariantInput(input *pb.VariantInput, listingID int64) *domain.ListingVariant {
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

// ProtoToAddImageInput converts pb.AddImageRequest to domain.ListingImage
func ProtoToAddImageInput(req *pb.AddImageRequest) *domain.ListingImage {
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

// ProductToProto converts domain.Product to pb.Product
func ProductToProto(p *domain.Product) *pb.Product {
	if p == nil {
		return nil
	}

	pbProduct := &pb.Product{
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
		pbProduct.Variants = make([]*pb.ProductVariant, len(p.Variants))
		for i, v := range p.Variants {
			pbProduct.Variants[i] = ProductVariantToProto(&v)
		}
	}

	return pbProduct
}

// ProductVariantToProto converts domain.ProductVariant to pb.ProductVariant
func ProductVariantToProto(v *domain.ProductVariant) *pb.ProductVariant {
	if v == nil {
		return nil
	}

	pbVariant := &pb.ProductVariant{
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

// ProtoToCreateProductInput converts pb.CreateProductRequest to domain.CreateProductInput
func ProtoToCreateProductInput(req *pb.CreateProductRequest) *domain.CreateProductInput {
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

	// Location fields
	if req.HasIndividualLocation != nil {
		input.HasIndividualLocation = *req.HasIndividualLocation
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
		input.ShowOnMap = *req.ShowOnMap
	}

	return input
}

// ProtoToUpdateProductInput converts pb.UpdateProductRequest to domain.UpdateProductInput
func ProtoToUpdateProductInput(req *pb.UpdateProductRequest) *domain.UpdateProductInput {
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

	if req.StockQuantity != nil {
		input.StockQuantity = req.StockQuantity
	}

	if req.StockStatus != nil {
		input.StockStatus = req.StockStatus
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
		input.HasIndividualLocation = req.HasIndividualLocation
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
		input.ShowOnMap = req.ShowOnMap
	}

	return input
}

// ProtoToCreateVariantInput converts pb.CreateProductVariantRequest to domain.CreateVariantInput
func ProtoToCreateVariantInput(req *pb.CreateProductVariantRequest) *domain.CreateVariantInput {
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

// ProtoToUpdateVariantInput converts pb.UpdateProductVariantRequest to domain.UpdateVariantInput
func ProtoToUpdateVariantInput(req *pb.UpdateProductVariantRequest) *domain.UpdateVariantInput {
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
	if req.StockStatus != nil {
		input.StockStatus = req.StockStatus
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
