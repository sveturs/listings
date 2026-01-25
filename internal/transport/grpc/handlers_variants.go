// Package grpc implements gRPC handlers for the listings microservice.
// This file contains the VariantHandler implementation.
package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/vondi-global/listings/api/proto/variants/v1"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/service"
)

// VariantHandler handles gRPC requests for variant operations
type VariantHandler struct {
	pb.UnimplementedVariantServiceServer
	variantService *service.VariantService
	logger         zerolog.Logger
}

// NewVariantHandler creates a new variant gRPC handler
func NewVariantHandler(variantService *service.VariantService, logger zerolog.Logger) *VariantHandler {
	return &VariantHandler{
		variantService: variantService,
		logger:         logger.With().Str("component", "variant_grpc_handler").Logger(),
	}
}

// ListVariants retrieves all variants for a product
func (h *VariantHandler) ListVariants(ctx context.Context, req *pb.ListVariantsRequest) (*pb.ListVariantsResponse, error) {
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	variants, err := h.variantService.ListByProduct(ctx, req.ProductId)
	if err != nil {
		h.logger.Error().Err(err).Str("product_id", req.ProductId).Msg("failed to list variants")
		return nil, status.Error(codes.Internal, "failed to list variants")
	}

	// Convert domain models to proto
	protoVariants := make([]*pb.Variant, len(variants))
	for i, v := range variants {
		protoVariants[i] = h.variantToProto(v)
	}

	return &pb.ListVariantsResponse{
		Variants: protoVariants,
		Total:    int32(len(variants)),
	}, nil
}

// ReserveStock reserves stock for an order
func (h *VariantHandler) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest) (*pb.ReserveStockResponse, error) {
	if req.VariantId == "" || req.OrderId == "" {
		return nil, status.Error(codes.InvalidArgument, "variant_id and order_id are required")
	}

	if req.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
	}

	serviceReq := &service.ReserveStockRequest{
		VariantID:  req.VariantId,
		OrderID:    req.OrderId,
		Quantity:   req.Quantity,
		TTLMinutes: req.TtlMinutes,
	}

	resp, err := h.variantService.ReserveStock(ctx, serviceReq)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to reserve stock")
		return nil, status.Error(codes.Internal, "failed to reserve stock")
	}

	return &pb.ReserveStockResponse{
		Success:        resp.Success,
		ReservationId:  resp.ReservationID,
		AvailableAfter: resp.AvailableAfter,
		ErrorMessage:   resp.ErrorMessage,
	}, nil
}

// ReleaseStock releases a stock reservation
func (h *VariantHandler) ReleaseStock(ctx context.Context, req *pb.ReleaseStockRequest) (*pb.ReleaseStockResponse, error) {
	if req.ReservationId == "" {
		return nil, status.Error(codes.InvalidArgument, "reservation_id is required")
	}

	err := h.variantService.ReleaseStock(ctx, req.ReservationId)
	if err != nil {
		h.logger.Error().Err(err).Str("reservation_id", req.ReservationId).Msg("failed to release stock")
		return nil, status.Error(codes.Internal, "failed to release stock")
	}

	return &pb.ReleaseStockResponse{Success: true}, nil
}

// ConfirmStockDeduction confirms a reservation and deducts stock
func (h *VariantHandler) ConfirmStockDeduction(ctx context.Context, req *pb.ConfirmStockDeductionRequest) (*pb.ConfirmStockDeductionResponse, error) {
	if req.ReservationId == "" {
		return nil, status.Error(codes.InvalidArgument, "reservation_id is required")
	}

	err := h.variantService.ConfirmStockDeduction(ctx, req.ReservationId)
	if err != nil {
		h.logger.Error().Err(err).Str("reservation_id", req.ReservationId).Msg("failed to confirm stock deduction")
		return nil, status.Error(codes.Internal, "failed to confirm stock deduction")
	}

	return &pb.ConfirmStockDeductionResponse{Success: true}, nil
}

// GetVariant retrieves a single variant by ID
func (h *VariantHandler) GetVariant(ctx context.Context, req *pb.GetVariantRequest) (*pb.VariantResponse, error) {
	// 1. Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// 2. Call service
	variant, err := h.variantService.GetByID(ctx, req.Id)
	if err != nil {
		if err == domain.ErrVariantNotFound {
			return nil, status.Error(codes.NotFound, "variant not found")
		}
		h.logger.Error().Err(err).Str("id", req.Id).Msg("failed to get variant")
		return nil, status.Error(codes.Internal, "failed to get variant")
	}

	// 3. Convert domain → proto
	protoVariant := h.variantToProto(variant)

	// 4. Return
	return &pb.VariantResponse{
		Variant: protoVariant,
	}, nil
}

// CreateVariant creates a new variant
func (h *VariantHandler) CreateVariant(ctx context.Context, req *pb.CreateVariantRequest) (*pb.VariantResponse, error) {
	// 1. Validate request
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	productUUID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product_id format")
	}

	if req.CategoryCode == "" {
		return nil, status.Error(codes.InvalidArgument, "category_code is required for SKU generation")
	}

	// 2. Convert proto → domain
	input := &domain.CreateVariantInputV2{
		ProductID:      productUUID,
		SKU:            req.Sku, // Empty is OK, auto-generated
		Price:          protoDoubleToPointer(req.Price),
		CompareAtPrice: protoDoubleToPointer(req.CompareAtPrice),
		StockQuantity:  req.InitialStock,
		LowStockAlert:  req.LowStockAlert,
		WeightGrams:    protoDoubleToPointer(req.WeightGrams),
		Barcode:        protoStringToPointer(req.Barcode),
		IsDefault:      req.IsDefault,
		Position:       req.Position,
		Attributes:     make([]domain.CreateVariantAttributeValue, len(req.Attributes)),
	}

	// Convert attributes
	for i, attr := range req.Attributes {
		input.Attributes[i] = domain.CreateVariantAttributeValue{
			AttributeID:  attr.AttributeId,
			ValueText:    protoStringToPointer(attr.ValueText),
			ValueNumber:  protoDoubleToPointer(attr.ValueNumber),
			ValueBoolean: protoBoolToPointer(attr.ValueBoolean),
		}
	}

	// 3. Call service
	variant, err := h.variantService.CreateVariant(ctx, input, req.CategoryCode)
	if err != nil {
		h.logger.Error().Err(err).Str("product_id", req.ProductId).Msg("failed to create variant")
		return nil, status.Error(codes.Internal, "failed to create variant")
	}

	// 4. Convert domain → proto
	protoVariant := h.variantToProto(variant)

	// 5. Return
	return &pb.VariantResponse{
		Variant: protoVariant,
	}, nil
}

// UpdateVariant updates an existing variant
func (h *VariantHandler) UpdateVariant(ctx context.Context, req *pb.UpdateVariantRequest) (*pb.VariantResponse, error) {
	// 1. Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// 2. Convert proto → domain
	input := &domain.UpdateVariantInputV2{
		SKU:            protoStringToPointer(req.Sku),
		Price:          protoDoubleToPointer(req.Price),
		CompareAtPrice: protoDoubleToPointer(req.CompareAtPrice),
		StockQuantity:  protoInt32ToPointer(req.StockQuantity),
		LowStockAlert:  protoInt32ToPointer(req.LowStockAlert),
		WeightGrams:    protoDoubleToPointer(req.WeightGrams),
		Barcode:        protoStringToPointer(req.Barcode),
		IsDefault:      protoBoolToPointer(req.IsDefault),
		Position:       protoInt32ToPointer(req.Position),
		Status:         protoStringToPointer(req.Status),
	}

	// 3. Call service
	variant, err := h.variantService.Update(ctx, req.Id, input)
	if err != nil {
		if err == domain.ErrVariantNotFound {
			return nil, status.Error(codes.NotFound, "variant not found")
		}
		h.logger.Error().Err(err).Str("id", req.Id).Msg("failed to update variant")
		return nil, status.Error(codes.Internal, "failed to update variant")
	}

	// 4. Convert domain → proto
	protoVariant := h.variantToProto(variant)

	// 5. Return
	return &pb.VariantResponse{
		Variant: protoVariant,
	}, nil
}

// DeleteVariant deletes a variant
func (h *VariantHandler) DeleteVariant(ctx context.Context, req *pb.DeleteVariantRequest) (*pb.DeleteVariantResponse, error) {
	// 1. Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// 2. Call service
	err := h.variantService.Delete(ctx, req.Id)
	if err != nil {
		if err == domain.ErrVariantNotFound {
			return nil, status.Error(codes.NotFound, "variant not found")
		}
		h.logger.Error().Err(err).Str("id", req.Id).Msg("failed to delete variant")
		return nil, status.Error(codes.Internal, "failed to delete variant")
	}

	// 3. Return success
	return &pb.DeleteVariantResponse{
		Success: true,
	}, nil
}

// GetVariantBySku retrieves a variant by SKU
func (h *VariantHandler) GetVariantBySku(ctx context.Context, req *pb.GetVariantBySkuRequest) (*pb.VariantResponse, error) {
	// 1. Validate request
	if req.Sku == "" {
		return nil, status.Error(codes.InvalidArgument, "sku is required")
	}

	// 2. Call service
	variant, err := h.variantService.GetBySKU(ctx, req.Sku)
	if err != nil {
		if err == domain.ErrVariantNotFound {
			return nil, status.Error(codes.NotFound, "variant not found")
		}
		h.logger.Error().Err(err).Str("sku", req.Sku).Msg("failed to get variant by SKU")
		return nil, status.Error(codes.Internal, "failed to get variant")
	}

	// 3. Convert domain → proto
	protoVariant := h.variantToProto(variant)

	// 4. Return
	return &pb.VariantResponse{
		Variant: protoVariant,
	}, nil
}

// FindVariantByAttributes finds a variant by attribute values
func (h *VariantHandler) FindVariantByAttributes(ctx context.Context, req *pb.FindVariantByAttributesRequest) (*pb.VariantResponse, error) {
	// 1. Validate request
	if req.ProductId == "" {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}

	if len(req.Attributes) == 0 {
		return nil, status.Error(codes.InvalidArgument, "attributes are required")
	}

	productUUID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product_id format")
	}

	// 2. Convert proto map to domain map (int32 -> string to int32 -> interface{})
	domainAttrs := make(map[int32]interface{}, len(req.Attributes))
	for attrID, value := range req.Attributes {
		domainAttrs[attrID] = value
	}

	filter := &domain.FindVariantByAttributesFilter{
		ProductID:  productUUID,
		Attributes: domainAttrs,
	}

	// 3. Call service
	variant, err := h.variantService.FindByAttributes(ctx, filter)
	if err != nil {
		if err == domain.ErrVariantNotFound {
			return nil, status.Error(codes.NotFound, "variant not found")
		}
		h.logger.Error().Err(err).Str("product_id", req.ProductId).Msg("failed to find variant by attributes")
		return nil, status.Error(codes.Internal, "failed to find variant")
	}

	// 4. Convert domain → proto
	protoVariant := h.variantToProto(variant)

	// 5. Return
	return &pb.VariantResponse{
		Variant: protoVariant,
	}, nil
}

// variantToProto converts domain.ProductVariantV2 to proto Variant
func (h *VariantHandler) variantToProto(v *domain.ProductVariantV2) *pb.Variant {
	variant := &pb.Variant{
		Id:                v.ID.String(),
		ProductId:         v.ProductID.String(),
		Sku:               v.SKU,
		StockQuantity:     v.StockQuantity,
		ReservedQuantity:  v.ReservedQuantity,
		AvailableQuantity: v.GetAvailableQuantity(),
		LowStockAlert:     v.LowStockAlert,
		IsDefault:         v.IsDefault,
		Position:          v.Position,
		Status:            string(v.Status),
		CreatedAt:         timestamppb.New(v.CreatedAt),
		UpdatedAt:         timestamppb.New(v.UpdatedAt),
	}

	if v.Price != nil {
		variant.Price = v.Price
	}

	if v.CompareAtPrice != nil {
		variant.CompareAtPrice = v.CompareAtPrice
	}

	if v.WeightGrams != nil {
		variant.WeightGrams = v.WeightGrams
	}

	if v.Barcode != nil {
		variant.Barcode = v.Barcode
	}

	// Convert attributes
	if v.Attributes != nil {
		attrs := make([]*pb.VariantAttribute, len(v.Attributes))
		for i, attr := range v.Attributes {
			attrs[i] = &pb.VariantAttribute{
				Id:          attr.ID.String(),
				AttributeId: attr.AttributeID,
			}

			if attr.ValueText != nil {
				attrs[i].ValueText = attr.ValueText
			}
			if attr.ValueNumber != nil {
				attrs[i].ValueNumber = attr.ValueNumber
			}
			if attr.ValueBoolean != nil {
				attrs[i].ValueBoolean = attr.ValueBoolean
			}
		}
		variant.Attributes = attrs
	}

	return variant
}

// Helper functions for proto optional fields

// protoDoubleToPointer converts optional proto double to *float64
func protoDoubleToPointer(val *float64) *float64 {
	if val == nil {
		return nil
	}
	return val
}

// protoStringToPointer converts optional proto string to *string
func protoStringToPointer(val *string) *string {
	if val == nil {
		return nil
	}
	return val
}

// protoBoolToPointer converts optional proto bool to *bool
func protoBoolToPointer(val *bool) *bool {
	if val == nil {
		return nil
	}
	return val
}

// protoInt32ToPointer converts optional proto int32 to *int32
func protoInt32ToPointer(val *int32) *int32 {
	if val == nil {
		return nil
	}
	return val
}
