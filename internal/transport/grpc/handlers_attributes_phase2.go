package grpc

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	attributessvcv1 "github.com/vondi-global/listings/api/proto/attributessvc/v1"
	"github.com/vondi-global/listings/internal/repository/postgres"
)

// AttributeServicePhase2Server реализует Phase 2 gRPC AttributeService methods
type AttributeServicePhase2Server struct {
	attributessvcv1.UnimplementedAttributeServiceServer
	attrRepo *postgres.AttributeRepository
}

// NewAttributeServicePhase2Server создаёт новый Phase 2 AttributeService server
func NewAttributeServicePhase2Server(attrRepo *postgres.AttributeRepository) *AttributeServicePhase2Server {
	return &AttributeServicePhase2Server{
		attrRepo: attrRepo,
	}
}

// GetAttributesByCategory получает атрибуты для категории с наследованием
func (s *AttributeServicePhase2Server) GetAttributesByCategory(
	ctx context.Context,
	req *attributessvcv1.GetAttributesByCategoryRequest,
) (*attributessvcv1.GetAttributesByCategoryResponse, error) {
	// Валидация
	if req.CategoryId == "" {
		return nil, status.Error(codes.InvalidArgument, "category_id is required")
	}

	locale := req.Locale
	if locale == "" {
		locale = "sr" // default
	}

	// Получить атрибуты через Repository (с наследованием)
	attrs, err := s.attrRepo.GetByCategoryID(ctx, req.CategoryId, locale)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get attributes: %v", err)
	}

	// Маппинг domain → proto
	protoAttrs := make([]*attributessvcv1.Attribute, len(attrs))
	for i, attr := range attrs {
		protoAttr := &attributessvcv1.Attribute{
			Id:                  attr.ID,
			Code:                attr.Code,
			Name:                attr.Name[locale],
			DisplayName:         attr.DisplayName[locale],
			AttributeType:       string(attr.AttributeType),
			Purpose:             string(attr.Purpose),
			IsSearchable:        attr.IsSearchable,
			IsFilterable:        attr.IsFilterable,
			IsRequired:          attr.IsRequired,
			IsVariantCompatible: attr.IsVariantCompatible,
			AffectsStock:        attr.AffectsStock,
			AffectsPrice:        attr.AffectsPrice,
			SortOrder:           attr.SortOrder,
		}

		// Опционально загрузить значения
		if req.IncludeValues {
			values, err := s.attrRepo.GetValues(ctx, attr.ID, locale)
			if err == nil && len(values) > 0 {
				protoAttr.Values = make([]*attributessvcv1.AttributeValue, len(values))
				for j, val := range values {
					metadataJSON, _ := json.Marshal(val.Metadata)
					protoAttr.Values[j] = &attributessvcv1.AttributeValue{
						Id:           val.ID,
						AttributeId:  val.AttributeID,
						Value:        val.Value,
						Label:        val.Label[locale],
						MetadataJson: string(metadataJSON),
						SortOrder:    val.SortOrder,
					}
				}
			}
		}

		protoAttrs[i] = protoAttr
	}

	return &attributessvcv1.GetAttributesByCategoryResponse{
		Attributes: protoAttrs,
	}, nil
}

// GetAttributeValues получает предопределённые значения атрибута
func (s *AttributeServicePhase2Server) GetAttributeValues(
	ctx context.Context,
	req *attributessvcv1.GetAttributeValuesRequest,
) (*attributessvcv1.GetAttributeValuesResponse, error) {
	if req.AttributeId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "attribute_id is required")
	}

	locale := req.Locale
	if locale == "" {
		locale = "sr"
	}

	values, err := s.attrRepo.GetValues(ctx, req.AttributeId, locale)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get attribute values: %v", err)
	}

	protoValues := make([]*attributessvcv1.AttributeValue, len(values))
	for i, val := range values {
		metadataJSON, _ := json.Marshal(val.Metadata)
		protoValues[i] = &attributessvcv1.AttributeValue{
			Id:           val.ID,
			AttributeId:  val.AttributeID,
			Value:        val.Value,
			Label:        val.Label[locale],
			MetadataJson: string(metadataJSON),
			SortOrder:    val.SortOrder,
		}
	}

	return &attributessvcv1.GetAttributeValuesResponse{
		Values: protoValues,
	}, nil
}
