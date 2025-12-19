// Package grpc implements gRPC handlers for the listings microservice.
// This file contains the VariantHandler implementation.
package grpc

import (
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/service"
	pb "github.com/vondi-global/listings/api/proto/variants/v1"
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
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
}

// CreateVariant creates a new variant
func (h *VariantHandler) CreateVariant(ctx context.Context, req *pb.CreateVariantRequest) (*pb.VariantResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
}

// UpdateVariant updates an existing variant
func (h *VariantHandler) UpdateVariant(ctx context.Context, req *pb.UpdateVariantRequest) (*pb.VariantResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
}

// DeleteVariant deletes a variant
func (h *VariantHandler) DeleteVariant(ctx context.Context, req *pb.DeleteVariantRequest) (*pb.DeleteVariantResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
}

// GetVariantBySku retrieves a variant by SKU
func (h *VariantHandler) GetVariantBySku(ctx context.Context, req *pb.GetVariantBySkuRequest) (*pb.VariantResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
}

// FindVariantByAttributes finds a variant by attribute values
func (h *VariantHandler) FindVariantByAttributes(ctx context.Context, req *pb.FindVariantByAttributesRequest) (*pb.VariantResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
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
