package grpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

// GetProduct retrieves a single product by ID
func (s *Server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	s.logger.Debug().Int64("product_id", req.ProductId).Msg("GetProduct called")

	// Validation
	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product ID must be greater than 0")
	}

	// Prepare storefront filter
	var storefrontID *int64
	if req.StorefrontId != nil && *req.StorefrontId > 0 {
		storefrontID = req.StorefrontId
	}

	// Get from service
	product, err := s.service.GetProduct(ctx, req.ProductId, storefrontID)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", req.ProductId).Msg("failed to get product")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get product: %v", err))
	}

	// Convert to proto
	protoProduct := ProductToProto(product)

	return &pb.ProductResponse{Product: protoProduct}, nil
}

// GetProductsBySKUs retrieves products by list of SKUs
func (s *Server) GetProductsBySKUs(ctx context.Context, req *pb.GetProductsBySKUsRequest) (*pb.ProductsResponse, error) {
	s.logger.Debug().Int("sku_count", len(req.Skus)).Msg("GetProductsBySKUs called")

	// Validation
	if len(req.Skus) == 0 {
		return nil, status.Error(codes.InvalidArgument, "SKU list cannot be empty")
	}

	if len(req.Skus) > 100 {
		return nil, status.Error(codes.InvalidArgument, "cannot request more than 100 SKUs at once")
	}

	// Prepare storefront filter
	var storefrontID *int64
	if req.StorefrontId != nil && *req.StorefrontId > 0 {
		storefrontID = req.StorefrontId
	}

	// Get from service
	products, err := s.service.GetProductsBySKUs(ctx, req.Skus, storefrontID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get products by SKUs")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get products: %v", err))
	}

	// Convert to proto
	protoProducts := make([]*pb.Product, 0, len(products))
	for _, p := range products {
		protoProducts = append(protoProducts, ProductToProto(p))
	}

	s.logger.Debug().Int("found_count", len(products)).Msg("products retrieved by SKUs")
	return &pb.ProductsResponse{
		Products:   protoProducts,
		TotalCount: int32(len(products)),
	}, nil
}

// GetProductsByIDs retrieves products by list of IDs
func (s *Server) GetProductsByIDs(ctx context.Context, req *pb.GetProductsByIDsRequest) (*pb.ProductsResponse, error) {
	s.logger.Debug().Int("id_count", len(req.ProductIds)).Msg("GetProductsByIDs called")

	// Validation
	if len(req.ProductIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "product ID list cannot be empty")
	}

	if len(req.ProductIds) > 100 {
		return nil, status.Error(codes.InvalidArgument, "cannot request more than 100 product IDs at once")
	}

	for _, id := range req.ProductIds {
		if id <= 0 {
			return nil, status.Error(codes.InvalidArgument, "all product IDs must be greater than 0")
		}
	}

	// Prepare storefront filter
	var storefrontID *int64
	if req.StorefrontId != nil && *req.StorefrontId > 0 {
		storefrontID = req.StorefrontId
	}

	// Get from service
	products, err := s.service.GetProductsByIDs(ctx, req.ProductIds, storefrontID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get products by IDs")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get products: %v", err))
	}

	// Convert to proto
	protoProducts := make([]*pb.Product, 0, len(products))
	for _, p := range products {
		protoProducts = append(protoProducts, ProductToProto(p))
	}

	s.logger.Debug().Int("found_count", len(products)).Msg("products retrieved by IDs")
	return &pb.ProductsResponse{
		Products:   protoProducts,
		TotalCount: int32(len(products)),
	}, nil
}

// ListProducts retrieves paginated list of products
func (s *Server) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ProductsResponse, error) {
	s.logger.Debug().Int64("storefront_id", req.StorefrontId).Int32("page", req.Page).Int32("page_size", req.PageSize).Msg("ListProducts called")

	// Validation
	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront ID must be greater than 0")
	}

	if req.Page <= 0 {
		return nil, status.Error(codes.InvalidArgument, "page must be greater than 0")
	}

	if req.PageSize <= 0 || req.PageSize > 100 {
		return nil, status.Error(codes.InvalidArgument, "page size must be between 1 and 100")
	}

	// Get is_active_only (defaults to false)
	isActiveOnly := false
	if req.IsActiveOnly != nil {
		isActiveOnly = *req.IsActiveOnly
	}

	// Get from service
	products, totalCount, err := s.service.ListProducts(ctx, req.StorefrontId, int(req.Page), int(req.PageSize), isActiveOnly)
	if err != nil {
		s.logger.Error().Err(err).Int64("storefront_id", req.StorefrontId).Msg("failed to list products")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list products: %v", err))
	}

	// Convert to proto
	protoProducts := make([]*pb.Product, 0, len(products))
	for _, p := range products {
		protoProducts = append(protoProducts, ProductToProto(p))
	}

	s.logger.Debug().Int("count", len(products)).Int("total", totalCount).Msg("products listed")
	return &pb.ProductsResponse{
		Products:   protoProducts,
		TotalCount: int32(totalCount),
	}, nil
}

// GetVariant retrieves a single variant by ID
func (s *Server) GetVariant(ctx context.Context, req *pb.GetVariantRequest) (*pb.VariantResponse, error) {
	s.logger.Debug().Int64("variant_id", req.VariantId).Msg("GetVariant called")

	// Validation
	if req.VariantId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "variant ID must be greater than 0")
	}

	// Prepare product filter
	var productID *int64
	if req.ProductId != nil && *req.ProductId > 0 {
		productID = req.ProductId
	}

	// Get from service
	variant, err := s.service.GetVariant(ctx, req.VariantId, productID)
	if err != nil {
		s.logger.Error().Err(err).Int64("variant_id", req.VariantId).Msg("failed to get variant")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "variant not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get variant: %v", err))
	}

	// Convert to proto
	protoVariant := ProductVariantToProto(variant)

	return &pb.VariantResponse{Variant: protoVariant}, nil
}

// GetVariantsByProductID retrieves all variants for a product
func (s *Server) GetVariantsByProductID(ctx context.Context, req *pb.GetVariantsByProductIDRequest) (*pb.ProductVariantsResponse, error) {
	s.logger.Debug().Int64("product_id", req.ProductId).Msg("GetVariantsByProductID called")

	// Validation
	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product ID must be greater than 0")
	}

	// Get is_active_only (defaults to false)
	isActiveOnly := false
	if req.IsActiveOnly != nil {
		isActiveOnly = *req.IsActiveOnly
	}

	// Get from service
	variants, err := s.service.GetVariantsByProductID(ctx, req.ProductId, isActiveOnly)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", req.ProductId).Msg("failed to get variants")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get variants: %v", err))
	}

	// Convert to proto
	protoVariants := make([]*pb.ProductVariant, 0, len(variants))
	for _, v := range variants {
		protoVariants = append(protoVariants, ProductVariantToProto(v))
	}

	s.logger.Debug().Int("count", len(variants)).Msg("variants retrieved")
	return &pb.ProductVariantsResponse{
		Variants: protoVariants,
	}, nil
}

// CreateProduct creates a new product
func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	s.logger.Debug().
		Int64("storefront_id", req.StorefrontId).
		Str("name", req.Name).
		Msg("CreateProduct called")

	// Validation
	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront ID must be greater than 0")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "product name cannot be empty")
	}

	if req.Price < 0 {
		return nil, status.Error(codes.InvalidArgument, "price must be non-negative")
	}

	if req.Currency == "" {
		return nil, status.Error(codes.InvalidArgument, "currency cannot be empty")
	}

	if req.CategoryId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "category ID must be greater than 0")
	}

	if req.StockQuantity < 0 {
		return nil, status.Error(codes.InvalidArgument, "stock quantity must be non-negative")
	}

	// Convert proto to domain input
	input := ProtoToCreateProductInput(req)

	// Create product via service
	product, err := s.service.CreateProduct(ctx, input)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create product")

		// Check for specific errors with placeholders
		errMsg := err.Error()
		if errMsg == "products.sku_duplicate" {
			return nil, status.Error(codes.AlreadyExists, "products.sku_duplicate")
		}

		// Generic error
		return nil, status.Error(codes.Internal, "products.create_failed")
	}

	// Convert domain product to proto
	protoProduct := ProductToProto(product)

	s.logger.Info().Int64("product_id", product.ID).Msg("product created successfully")
	return &pb.ProductResponse{Product: protoProduct}, nil
}

// UpdateProduct updates an existing product
func (s *Server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	s.logger.Debug().
		Int64("product_id", req.ProductId).
		Int64("storefront_id", req.StorefrontId).
		Msg("UpdateProduct called")

	// Validation
	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product ID must be greater than 0")
	}

	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront ID must be greater than 0")
	}

	// Validate at least one field is being updated
	hasUpdate := req.Name != nil || req.Description != nil || req.Price != nil ||
		req.Sku != nil || req.Barcode != nil || req.StockQuantity != nil ||
		req.StockStatus != nil || req.IsActive != nil || req.Attributes != nil ||
		req.HasIndividualLocation != nil || req.IndividualAddress != nil ||
		req.IndividualLatitude != nil || req.IndividualLongitude != nil ||
		req.LocationPrivacy != nil || req.ShowOnMap != nil

	if !hasUpdate {
		return nil, status.Error(codes.InvalidArgument, "at least one field must be specified for update")
	}

	// Validate price if provided
	if req.Price != nil && *req.Price < 0 {
		return nil, status.Error(codes.InvalidArgument, "price must be non-negative")
	}

	// Validate stock_quantity if provided
	if req.StockQuantity != nil && *req.StockQuantity < 0 {
		return nil, status.Error(codes.InvalidArgument, "stock quantity must be non-negative")
	}

	// Validate stock_status if provided
	if req.StockStatus != nil {
		validStatuses := map[string]bool{
			"in_stock": true, "low_stock": true, "out_of_stock": true, "pre_order": true,
		}
		if !validStatuses[*req.StockStatus] {
			return nil, status.Error(codes.InvalidArgument, "invalid stock status")
		}
	}

	// Convert proto to domain input
	input := ProtoToUpdateProductInput(req)

	// Update product via service
	product, err := s.service.UpdateProduct(ctx, req.ProductId, req.StorefrontId, input)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", req.ProductId).Msg("failed to update product")

		// Check for specific errors with placeholders
		errMsg := err.Error()
		if errMsg == "products.not_found" {
			return nil, status.Error(codes.NotFound, "products.not_found")
		}
		if errMsg == "products.sku_duplicate" {
			return nil, status.Error(codes.AlreadyExists, "products.sku_duplicate")
		}

		// Generic error
		return nil, status.Error(codes.Internal, "products.update_failed")
	}

	// Convert domain product to proto
	protoProduct := ProductToProto(product)

	s.logger.Info().Int64("product_id", product.ID).Msg("product updated successfully")
	return &pb.ProductResponse{Product: protoProduct}, nil
}

// DeleteProduct deletes a product (soft or hard delete)
func (s *Server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	s.logger.Debug().
		Int64("product_id", req.ProductId).
		Int64("storefront_id", req.StorefrontId).
		Bool("hard_delete", req.HardDelete).
		Msg("DeleteProduct called")

	// Validation
	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product ID must be greater than 0")
	}

	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront ID must be greater than 0")
	}

	// Delete product via service
	variantsDeleted, err := s.service.DeleteProduct(ctx, req.ProductId, req.StorefrontId, req.HardDelete)
	if err != nil {
		s.logger.Error().Err(err).Int64("product_id", req.ProductId).Msg("failed to delete product")

		// Check for specific errors with placeholders
		errMsg := err.Error()
		if errMsg == "products.not_found" {
			return nil, status.Error(codes.NotFound, "products.not_found")
		}
		if errMsg == "products.has_active_orders" {
			return nil, status.Error(codes.FailedPrecondition, "products.has_active_orders")
		}

		// Generic error
		return nil, status.Error(codes.Internal, "products.delete_failed")
	}

	deleteMsg := "Product soft deleted successfully"
	if req.HardDelete {
		deleteMsg = "Product hard deleted successfully"
	}

	s.logger.Info().
		Int64("product_id", req.ProductId).
		Int32("variants_deleted", variantsDeleted).
		Msg("product deleted successfully")

	return &pb.DeleteProductResponse{
		Success:         true,
		Message:         &deleteMsg,
		VariantsDeleted: variantsDeleted,
	}, nil
}

// CreateProductVariant creates a new product variant
func (s *Server) CreateProductVariant(ctx context.Context, req *pb.CreateProductVariantRequest) (*pb.CreateProductVariantResponse, error) {
	s.logger.Debug().
		Int64("product_id", req.ProductId).
		Msg("CreateProductVariant called")

	// Validation
	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product ID must be greater than 0")
	}

	if req.StockQuantity < 0 {
		return nil, status.Error(codes.InvalidArgument, "stock quantity must be non-negative")
	}

	// Convert proto to domain input
	input := ProtoToCreateVariantInput(req)

	// Create variant via service
	variant, err := s.service.CreateProductVariant(ctx, input)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create variant")

		// Check for specific errors with placeholders
		errMsg := err.Error()
		switch errMsg {
		case "variants.sku_duplicate":
			return nil, status.Error(codes.AlreadyExists, "variants.sku_duplicate")
		case "variants.product_not_found":
			return nil, status.Error(codes.NotFound, "variants.product_not_found")
		case "variants.product_no_variants":
			return nil, status.Error(codes.FailedPrecondition, "variants.product_no_variants")
		case "variants.invalid_product_id":
			return nil, status.Error(codes.InvalidArgument, "variants.invalid_product_id")
		case "variants.invalid_stock_quantity":
			return nil, status.Error(codes.InvalidArgument, "variants.invalid_stock_quantity")
		default:
			// Generic error
			return nil, status.Error(codes.Internal, "variants.create_failed")
		}
	}

	// Convert domain variant to proto
	protoVariant := ProductVariantToProto(variant)

	s.logger.Info().Int64("variant_id", variant.ID).Msg("variant created successfully")
	return &pb.CreateProductVariantResponse{Variant: protoVariant}, nil
}

// UpdateProductVariant updates an existing product variant
func (s *Server) UpdateProductVariant(ctx context.Context, req *pb.UpdateProductVariantRequest) (*pb.UpdateProductVariantResponse, error) {
	s.logger.Debug().
		Int64("variant_id", req.VariantId).
		Int64("product_id", req.ProductId).
		Msg("UpdateProductVariant called")

	// Validation
	if req.VariantId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "variant ID must be greater than 0")
	}

	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product ID must be greater than 0")
	}

	// Convert proto to domain input
	input := ProtoToUpdateVariantInput(req)

	// Update variant via service
	variant, err := s.service.UpdateProductVariant(ctx, req.VariantId, req.ProductId, input)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to update variant")

		// Check for specific errors with placeholders
		errMsg := err.Error()
		switch errMsg {
		case "variants.not_found":
			return nil, status.Error(codes.NotFound, "variants.not_found")
		case "variants.sku_duplicate":
			return nil, status.Error(codes.AlreadyExists, "variants.sku_duplicate")
		case "variants.last_variant":
			return nil, status.Error(codes.FailedPrecondition, "variants.last_variant")
		case "variants.invalid_variant_id":
			return nil, status.Error(codes.InvalidArgument, "variants.invalid_variant_id")
		case "variants.invalid_product_id":
			return nil, status.Error(codes.InvalidArgument, "variants.invalid_product_id")
		default:
			// Generic error
			return nil, status.Error(codes.Internal, "variants.update_failed")
		}
	}

	// Convert domain variant to proto
	protoVariant := ProductVariantToProto(variant)

	s.logger.Info().Int64("variant_id", variant.ID).Msg("variant updated successfully")
	return &pb.UpdateProductVariantResponse{Variant: protoVariant}, nil
}

// DeleteProductVariant deletes a product variant
func (s *Server) DeleteProductVariant(ctx context.Context, req *pb.DeleteProductVariantRequest) (*pb.DeleteProductVariantResponse, error) {
	s.logger.Debug().
		Int64("variant_id", req.VariantId).
		Int64("product_id", req.ProductId).
		Msg("DeleteProductVariant called")

	// Validation
	if req.VariantId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "variant ID must be greater than 0")
	}

	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product ID must be greater than 0")
	}

	// Delete variant via service
	err := s.service.DeleteProductVariant(ctx, req.VariantId, req.ProductId)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to delete variant")

		// Check for specific errors with placeholders
		errMsg := err.Error()
		switch errMsg {
		case "variants.not_found":
			return nil, status.Error(codes.NotFound, "variants.not_found")
		case "variants.invalid_variant_id":
			return nil, status.Error(codes.InvalidArgument, "variants.invalid_variant_id")
		case "variants.invalid_product_id":
			return nil, status.Error(codes.InvalidArgument, "variants.invalid_product_id")
		default:
			// Generic error
			return nil, status.Error(codes.Internal, "variants.delete_failed")
		}
	}

	s.logger.Info().
		Int64("variant_id", req.VariantId).
		Int64("product_id", req.ProductId).
		Msg("variant deleted successfully")

	return &pb.DeleteProductVariantResponse{Success: true}, nil
}
