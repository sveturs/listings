package listings

import (
	"context"
	"fmt"
	"time"

	pb "github.com/sveturs/listings/api/proto/listings/v1"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the Listings gRPC service client
type Client struct {
	client pb.ListingsServiceClient
	conn   *grpc.ClientConn
	logger zerolog.Logger
}

// NewClient creates a new Listings gRPC client
func NewClient(address string, logger zerolog.Logger) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to listings service at %s: %w", address, err)
	}

	return &Client{
		client: pb.NewListingsServiceClient(conn),
		conn:   conn,
		logger: logger.With().Str("component", "listings_client").Logger(),
	}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	return c.conn.Close()
}

// Stock management methods

// DecrementStock decrements stock for multiple listing items
func (c *Client) DecrementStock(ctx context.Context, items []*pb.StockItem, orderID string) (*pb.DecrementStockResponse, error) {
	c.logger.Info().
		Str("order_id", orderID).
		Int("items_count", len(items)).
		Msg("Decrementing stock")

	resp, err := c.client.DecrementStock(ctx, &pb.DecrementStockRequest{
		Items:   items,
		OrderId: &orderID,
	})
	if err != nil {
		c.logger.Error().Err(err).
			Str("order_id", orderID).
			Msg("Failed to decrement stock")
		return nil, err
	}

	if !resp.Success {
		c.logger.Warn().
			Str("error", *resp.Error).
			Str("order_id", orderID).
			Msg("Stock decrement failed")
	}

	return resp, nil
}

// RollbackStock rolls back stock decrement (e.g., when order is canceled)
func (c *Client) RollbackStock(ctx context.Context, items []*pb.StockItem, orderID string) error {
	c.logger.Warn().
		Str("order_id", orderID).
		Int("items_count", len(items)).
		Msg("Rolling back stock")

	resp, err := c.client.RollbackStock(ctx, &pb.RollbackStockRequest{
		Items:   items,
		OrderId: &orderID,
	})
	if err != nil {
		c.logger.Error().Err(err).
			Str("order_id", orderID).
			Msg("Failed to rollback stock")
		return err
	}

	if !resp.Success {
		c.logger.Error().
			Str("order_id", orderID).
			Msg("Stock rollback returned success=false")
	}

	return nil
}

// CheckStockAvailability checks if sufficient stock is available for all items
func (c *Client) CheckStockAvailability(ctx context.Context, items []*pb.StockItem) (*pb.CheckStockAvailabilityResponse, error) {
	c.logger.Debug().
		Int("items_count", len(items)).
		Msg("Checking stock availability")

	resp, err := c.client.CheckStockAvailability(ctx, &pb.CheckStockAvailabilityRequest{
		Items: items,
	})
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to check stock availability")
		return nil, err
	}

	if !resp.AllAvailable {
		c.logger.Warn().
			Interface("items", resp.Items).
			Msg("Some items are not available")
	}

	return resp, nil
}

// ============================================================================
// Product CRUD Operations
// ============================================================================

// CreateProduct creates a new B2C product
// Transaction: Creates product + initializes stock record
func (c *Client) CreateProduct(ctx context.Context, storefrontID int64, product *pb.ProductInput) (*pb.Product, error) {
	c.logger.Info().
		Int64("storefront_id", storefrontID).
		Str("product_name", product.Name).
		Msg("Creating product")

	resp, err := c.client.CreateProduct(ctx, &pb.CreateProductRequest{
		StorefrontId:          storefrontID,
		Name:                  product.Name,
		Description:           product.Description,
		Price:                 product.Price,
		Currency:              product.Currency,
		CategoryId:            product.CategoryId,
		Sku:                   product.Sku,
		Barcode:               product.Barcode,
		StockQuantity:         product.StockQuantity,
		IsActive:              product.IsActive,
		Attributes:            product.Attributes,
		HasIndividualLocation: product.HasIndividualLocation,
		IndividualAddress:     product.IndividualAddress,
		IndividualLatitude:    product.IndividualLatitude,
		IndividualLongitude:   product.IndividualLongitude,
		LocationPrivacy:       product.LocationPrivacy,
		ShowOnMap:             product.ShowOnMap,
		HasVariants:           product.HasVariants,
	})
	if err != nil {
		c.logger.Error().Err(err).
			Int64("storefront_id", storefrontID).
			Str("product_name", product.Name).
			Msg("Failed to create product")
		return nil, err
	}

	c.logger.Info().
		Int64("product_id", resp.Product.Id).
		Str("product_name", resp.Product.Name).
		Msg("Product created successfully")

	return resp.Product, nil
}

// UpdateProduct updates an existing product (partial update via field mask)
// Supports field masking for efficient updates
func (c *Client) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	c.logger.Info().
		Int64("product_id", req.ProductId).
		Int64("storefront_id", req.StorefrontId).
		Msg("Updating product")

	resp, err := c.client.UpdateProduct(ctx, req)
	if err != nil {
		c.logger.Error().Err(err).
			Int64("product_id", req.ProductId).
			Msg("Failed to update product")
		return nil, err
	}

	c.logger.Info().
		Int64("product_id", resp.Product.Id).
		Msg("Product updated successfully")

	return resp.Product, nil
}

// DeleteProduct removes a product (soft delete by default)
// Cascade behavior: Soft-deletes all variants, images, reservations
func (c *Client) DeleteProduct(ctx context.Context, productID, storefrontID int64, hardDelete bool) (*pb.DeleteProductResponse, error) {
	c.logger.Warn().
		Int64("product_id", productID).
		Int64("storefront_id", storefrontID).
		Bool("hard_delete", hardDelete).
		Msg("Deleting product")

	resp, err := c.client.DeleteProduct(ctx, &pb.DeleteProductRequest{
		ProductId:    productID,
		StorefrontId: storefrontID,
		HardDelete:   hardDelete,
	})
	if err != nil {
		c.logger.Error().Err(err).
			Int64("product_id", productID).
			Msg("Failed to delete product")
		return nil, err
	}

	if !resp.Success {
		c.logger.Error().
			Int64("product_id", productID).
			Str("message", *resp.Message).
			Msg("Product deletion returned success=false")
	} else {
		c.logger.Info().
			Int64("product_id", productID).
			Int32("variants_deleted", resp.VariantsDeleted).
			Msg("Product deleted successfully")
	}

	return resp, nil
}

// ============================================================================
// Product Variant CRUD Operations
// ============================================================================

// CreateProductVariant creates a new variant for a product
// Transaction: Creates variant + initializes stock record
func (c *Client) CreateProductVariant(ctx context.Context, productID int64, variant *pb.ProductVariantInput) (*pb.ProductVariant, error) {
	c.logger.Info().
		Int64("product_id", productID).
		Msg("Creating product variant")

	resp, err := c.client.CreateProductVariant(ctx, &pb.CreateProductVariantRequest{
		ProductId:         productID,
		Sku:               variant.Sku,
		Barcode:           variant.Barcode,
		Price:             variant.Price,
		CompareAtPrice:    variant.CompareAtPrice,
		CostPrice:         variant.CostPrice,
		StockQuantity:     variant.StockQuantity,
		LowStockThreshold: variant.LowStockThreshold,
		VariantAttributes: variant.VariantAttributes,
		Weight:            variant.Weight,
		Dimensions:        variant.Dimensions,
		IsActive:          variant.IsActive,
		IsDefault:         variant.IsDefault,
	})
	if err != nil {
		c.logger.Error().Err(err).
			Int64("product_id", productID).
			Msg("Failed to create product variant")
		return nil, err
	}

	c.logger.Info().
		Int64("variant_id", resp.Variant.Id).
		Int64("product_id", productID).
		Msg("Product variant created successfully")

	return resp.Variant, nil
}

// UpdateProductVariant updates an existing product variant
// Supports partial updates via field mask
func (c *Client) UpdateProductVariant(ctx context.Context, req *pb.UpdateProductVariantRequest) (*pb.ProductVariant, error) {
	c.logger.Info().
		Int64("variant_id", req.VariantId).
		Int64("product_id", req.ProductId).
		Msg("Updating product variant")

	resp, err := c.client.UpdateProductVariant(ctx, req)
	if err != nil {
		c.logger.Error().Err(err).
			Int64("variant_id", req.VariantId).
			Msg("Failed to update product variant")
		return nil, err
	}

	c.logger.Info().
		Int64("variant_id", resp.Variant.Id).
		Msg("Product variant updated successfully")

	return resp.Variant, nil
}

// DeleteProductVariant removes a variant from a product
func (c *Client) DeleteProductVariant(ctx context.Context, variantID, productID int64) error {
	c.logger.Warn().
		Int64("variant_id", variantID).
		Int64("product_id", productID).
		Msg("Deleting product variant")

	resp, err := c.client.DeleteProductVariant(ctx, &pb.DeleteProductVariantRequest{
		VariantId: variantID,
		ProductId: productID,
	})
	if err != nil {
		c.logger.Error().Err(err).
			Int64("variant_id", variantID).
			Msg("Failed to delete product variant")
		return err
	}

	if !resp.Success {
		c.logger.Error().
			Int64("variant_id", variantID).
			Str("message", *resp.Message).
			Msg("Product variant deletion returned success=false")
		return fmt.Errorf("variant deletion failed: %s", *resp.Message)
	}

	c.logger.Info().
		Int64("variant_id", variantID).
		Msg("Product variant deleted successfully")

	return nil
}

// ============================================================================
// Bulk Product Operations
// ============================================================================

// BulkCreateProducts creates multiple products in a single transaction
// Recommended batch size: 1-1000 items
func (c *Client) BulkCreateProducts(ctx context.Context, storefrontID int64, products []*pb.ProductInput) (*pb.BulkCreateProductsResponse, error) {
	if len(products) == 0 {
		return nil, fmt.Errorf("products list cannot be empty")
	}
	if len(products) > 1000 {
		return nil, fmt.Errorf("batch size exceeds maximum limit of 1000 items (got %d)", len(products))
	}

	c.logger.Info().
		Int64("storefront_id", storefrontID).
		Int("products_count", len(products)).
		Msg("Bulk creating products")

	resp, err := c.client.BulkCreateProducts(ctx, &pb.BulkCreateProductsRequest{
		StorefrontId: storefrontID,
		Products:     products,
	})
	if err != nil {
		c.logger.Error().Err(err).
			Int64("storefront_id", storefrontID).
			Int("products_count", len(products)).
			Msg("Failed to bulk create products")
		return nil, err
	}

	c.logger.Info().
		Int32("successful", resp.SuccessfulCount).
		Int32("failed", resp.FailedCount).
		Msg("Bulk product creation completed")

	if resp.FailedCount > 0 {
		c.logger.Warn().
			Int32("failed_count", resp.FailedCount).
			Interface("errors", resp.Errors).
			Msg("Some products failed to create")
	}

	return resp, nil
}

// BulkUpdateProducts updates multiple products in a single transaction
// Supports partial updates via field mask per product
func (c *Client) BulkUpdateProducts(ctx context.Context, storefrontID int64, updates []*pb.ProductUpdateInput) (*pb.BulkUpdateProductsResponse, error) {
	if len(updates) == 0 {
		return nil, fmt.Errorf("updates list cannot be empty")
	}
	if len(updates) > 1000 {
		return nil, fmt.Errorf("batch size exceeds maximum limit of 1000 items (got %d)", len(updates))
	}

	c.logger.Info().
		Int64("storefront_id", storefrontID).
		Int("updates_count", len(updates)).
		Msg("Bulk updating products")

	resp, err := c.client.BulkUpdateProducts(ctx, &pb.BulkUpdateProductsRequest{
		StorefrontId: storefrontID,
		Updates:      updates,
	})
	if err != nil {
		c.logger.Error().Err(err).
			Int64("storefront_id", storefrontID).
			Int("updates_count", len(updates)).
			Msg("Failed to bulk update products")
		return nil, err
	}

	c.logger.Info().
		Int32("successful", resp.SuccessfulCount).
		Int32("failed", resp.FailedCount).
		Msg("Bulk product update completed")

	if resp.FailedCount > 0 {
		c.logger.Warn().
			Int32("failed_count", resp.FailedCount).
			Interface("errors", resp.Errors).
			Msg("Some products failed to update")
	}

	return resp, nil
}

// BulkDeleteProducts deletes multiple products in a single transaction
// Cascade behavior: Deletes all associated variants, images, reservations
func (c *Client) BulkDeleteProducts(ctx context.Context, storefrontID int64, productIDs []int64, hardDelete bool) (*pb.BulkDeleteProductsResponse, error) {
	if len(productIDs) == 0 {
		return nil, fmt.Errorf("product IDs list cannot be empty")
	}
	if len(productIDs) > 1000 {
		return nil, fmt.Errorf("batch size exceeds maximum limit of 1000 items (got %d)", len(productIDs))
	}

	c.logger.Warn().
		Int64("storefront_id", storefrontID).
		Int("products_count", len(productIDs)).
		Bool("hard_delete", hardDelete).
		Msg("Bulk deleting products")

	resp, err := c.client.BulkDeleteProducts(ctx, &pb.BulkDeleteProductsRequest{
		StorefrontId: storefrontID,
		ProductIds:   productIDs,
		HardDelete:   hardDelete,
	})
	if err != nil {
		c.logger.Error().Err(err).
			Int64("storefront_id", storefrontID).
			Int("products_count", len(productIDs)).
			Msg("Failed to bulk delete products")
		return nil, err
	}

	c.logger.Info().
		Int32("successful", resp.SuccessfulCount).
		Int32("failed", resp.FailedCount).
		Int32("variants_deleted", resp.VariantsDeleted).
		Msg("Bulk product deletion completed")

	if resp.FailedCount > 0 {
		c.logger.Warn().
			Int32("failed_count", resp.FailedCount).
			Interface("errors", resp.Errors).
			Msg("Some products failed to delete")
	}

	return resp, nil
}

// BulkCreateProductVariants creates multiple variants for a product
// Recommended for product imports with size/color matrices
func (c *Client) BulkCreateProductVariants(ctx context.Context, productID int64, variants []*pb.ProductVariantInput) (*pb.BulkCreateProductVariantsResponse, error) {
	if len(variants) == 0 {
		return nil, fmt.Errorf("variants list cannot be empty")
	}
	if len(variants) > 1000 {
		return nil, fmt.Errorf("batch size exceeds maximum limit of 1000 items (got %d)", len(variants))
	}

	c.logger.Info().
		Int64("product_id", productID).
		Int("variants_count", len(variants)).
		Msg("Bulk creating product variants")

	resp, err := c.client.BulkCreateProductVariants(ctx, &pb.BulkCreateProductVariantsRequest{
		ProductId: productID,
		Variants:  variants,
	})
	if err != nil {
		c.logger.Error().Err(err).
			Int64("product_id", productID).
			Int("variants_count", len(variants)).
			Msg("Failed to bulk create product variants")
		return nil, err
	}

	c.logger.Info().
		Int32("successful", resp.SuccessfulCount).
		Int32("failed", resp.FailedCount).
		Msg("Bulk variant creation completed")

	if resp.FailedCount > 0 {
		c.logger.Warn().
			Int32("failed_count", resp.FailedCount).
			Interface("errors", resp.Errors).
			Msg("Some variants failed to create")
	}

	return resp, nil
}
