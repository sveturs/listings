package grpc

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/sveturs/listings/api/proto/attributes/v1"
	"github.com/sveturs/listings/internal/domain"
)

// ============================================================================
// Attribute Type Conversions
// ============================================================================

// ProtoToAttributeType converts proto AttributeType to domain AttributeType
func ProtoToAttributeType(pt pb.AttributeType) domain.AttributeType {
	switch pt {
	case pb.AttributeType_ATTRIBUTE_TYPE_TEXT:
		return domain.AttributeTypeText
	case pb.AttributeType_ATTRIBUTE_TYPE_TEXTAREA:
		return domain.AttributeTypeTextarea
	case pb.AttributeType_ATTRIBUTE_TYPE_NUMBER:
		return domain.AttributeTypeNumber
	case pb.AttributeType_ATTRIBUTE_TYPE_BOOLEAN:
		return domain.AttributeTypeBoolean
	case pb.AttributeType_ATTRIBUTE_TYPE_SELECT:
		return domain.AttributeTypeSelect
	case pb.AttributeType_ATTRIBUTE_TYPE_MULTISELECT:
		return domain.AttributeTypeMultiselect
	case pb.AttributeType_ATTRIBUTE_TYPE_DATE:
		return domain.AttributeTypeDate
	case pb.AttributeType_ATTRIBUTE_TYPE_COLOR:
		return domain.AttributeTypeColor
	case pb.AttributeType_ATTRIBUTE_TYPE_SIZE:
		return domain.AttributeTypeSize
	default:
		return domain.AttributeTypeText
	}
}

// AttributeTypeToProto converts domain AttributeType to proto AttributeType
func AttributeTypeToProto(dt domain.AttributeType) pb.AttributeType {
	switch dt {
	case domain.AttributeTypeText:
		return pb.AttributeType_ATTRIBUTE_TYPE_TEXT
	case domain.AttributeTypeTextarea:
		return pb.AttributeType_ATTRIBUTE_TYPE_TEXTAREA
	case domain.AttributeTypeNumber:
		return pb.AttributeType_ATTRIBUTE_TYPE_NUMBER
	case domain.AttributeTypeBoolean:
		return pb.AttributeType_ATTRIBUTE_TYPE_BOOLEAN
	case domain.AttributeTypeSelect:
		return pb.AttributeType_ATTRIBUTE_TYPE_SELECT
	case domain.AttributeTypeMultiselect:
		return pb.AttributeType_ATTRIBUTE_TYPE_MULTISELECT
	case domain.AttributeTypeDate:
		return pb.AttributeType_ATTRIBUTE_TYPE_DATE
	case domain.AttributeTypeColor:
		return pb.AttributeType_ATTRIBUTE_TYPE_COLOR
	case domain.AttributeTypeSize:
		return pb.AttributeType_ATTRIBUTE_TYPE_SIZE
	default:
		return pb.AttributeType_ATTRIBUTE_TYPE_UNSPECIFIED
	}
}

// ProtoToAttributePurpose converts proto AttributePurpose to domain AttributePurpose
func ProtoToAttributePurpose(pp pb.AttributePurpose) domain.AttributePurpose {
	switch pp {
	case pb.AttributePurpose_ATTRIBUTE_PURPOSE_REGULAR:
		return domain.AttributePurposeRegular
	case pb.AttributePurpose_ATTRIBUTE_PURPOSE_VARIANT:
		return domain.AttributePurposeVariant
	case pb.AttributePurpose_ATTRIBUTE_PURPOSE_BOTH:
		return domain.AttributePurposeBoth
	default:
		return domain.AttributePurposeRegular
	}
}

// AttributePurposeToProto converts domain AttributePurpose to proto AttributePurpose
func AttributePurposeToProto(dp domain.AttributePurpose) pb.AttributePurpose {
	switch dp {
	case domain.AttributePurposeRegular:
		return pb.AttributePurpose_ATTRIBUTE_PURPOSE_REGULAR
	case domain.AttributePurposeVariant:
		return pb.AttributePurpose_ATTRIBUTE_PURPOSE_VARIANT
	case domain.AttributePurposeBoth:
		return pb.AttributePurpose_ATTRIBUTE_PURPOSE_BOTH
	default:
		return pb.AttributePurpose_ATTRIBUTE_PURPOSE_UNSPECIFIED
	}
}

// ============================================================================
// Struct/Map Conversions (i18n, JSONB)
// ============================================================================

// StructToMap converts proto Struct to map[string]interface{}
func StructToMap(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}
	return s.AsMap()
}

// StructToStringMap converts proto Struct to map[string]string (for i18n)
func StructToStringMap(s *structpb.Struct) map[string]string {
	if s == nil || s.Fields == nil {
		return make(map[string]string)
	}

	result := make(map[string]string)
	for key, value := range s.Fields {
		if value.GetStringValue() != "" {
			result[key] = value.GetStringValue()
		}
	}
	return result
}

// MapToStruct converts map[string]interface{} to proto Struct
func MapToStruct(m map[string]interface{}) (*structpb.Struct, error) {
	if m == nil {
		return nil, nil
	}
	return structpb.NewStruct(m)
}

// StringMapToStruct converts map[string]string to proto Struct
func StringMapToStruct(m map[string]string) (*structpb.Struct, error) {
	if m == nil || len(m) == 0 {
		return structpb.NewStruct(map[string]interface{}{})
	}

	fields := make(map[string]interface{})
	for k, v := range m {
		fields[k] = v
	}
	return structpb.NewStruct(fields)
}

// StructToAttributeOptions converts proto Struct to []AttributeOption
func StructToAttributeOptions(s *structpb.Struct) ([]domain.AttributeOption, error) {
	if s == nil {
		return nil, nil
	}

	// Proto stores options as {"options": [{"value": "...", "label": {...}}, ...]}
	// Or as an array directly
	asMap := s.AsMap()

	// Try to extract as array
	var optionsArray []interface{}
	if arr, ok := asMap["options"].([]interface{}); ok {
		optionsArray = arr
	} else if len(asMap) > 0 {
		// If not wrapped, assume the entire struct is an array
		// This is a simplified conversion - in practice, you may need more complex logic
		return nil, nil
	}

	if len(optionsArray) == 0 {
		return nil, nil
	}

	result := make([]domain.AttributeOption, 0, len(optionsArray))
	for _, item := range optionsArray {
		optMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		opt := domain.AttributeOption{}

		if value, ok := optMap["value"].(string); ok {
			opt.Value = value
		}

		if labelMap, ok := optMap["label"].(map[string]interface{}); ok {
			opt.Label = make(map[string]string)
			for k, v := range labelMap {
				if strVal, ok := v.(string); ok {
					opt.Label[k] = strVal
				}
			}
		}

		result = append(result, opt)
	}

	return result, nil
}

// AttributeOptionsToStruct converts []AttributeOption to proto Struct
func AttributeOptionsToStruct(opts []domain.AttributeOption) (*structpb.Struct, error) {
	if opts == nil || len(opts) == 0 {
		return nil, nil
	}

	optionsArray := make([]interface{}, len(opts))
	for i, opt := range opts {
		optMap := map[string]interface{}{
			"value": opt.Value,
			"label": opt.Label,
		}
		optionsArray[i] = optMap
	}

	return structpb.NewStruct(map[string]interface{}{
		"options": optionsArray,
	})
}

// ============================================================================
// Domain to Proto Conversions
// ============================================================================

// DomainAttributeToProto converts domain.Attribute to proto Attribute
func DomainAttributeToProto(attr *domain.Attribute) *pb.Attribute {
	if attr == nil {
		return nil
	}

	pbAttr := &pb.Attribute{
		Id:                  attr.ID,
		Code:                attr.Code,
		AttributeType:       AttributeTypeToProto(attr.AttributeType),
		Purpose:             AttributePurposeToProto(attr.Purpose),
		IsSearchable:        attr.IsSearchable,
		IsFilterable:        attr.IsFilterable,
		IsRequired:          attr.IsRequired,
		IsVariantCompatible: attr.IsVariantCompatible,
		AffectsStock:        attr.AffectsStock,
		AffectsPrice:        attr.AffectsPrice,
		ShowInCard:          attr.ShowInCard,
		IsActive:            attr.IsActive,
		SortOrder:           attr.SortOrder,
		Icon:                attr.Icon,
		CreatedAt:           timestamppb.New(attr.CreatedAt),
		UpdatedAt:           timestamppb.New(attr.UpdatedAt),
	}

	// Convert i18n fields
	if name, err := StringMapToStruct(attr.Name); err == nil {
		pbAttr.Name = name
	}
	if displayName, err := StringMapToStruct(attr.DisplayName); err == nil {
		pbAttr.DisplayName = displayName
	}

	// Convert JSONB fields
	if options, err := AttributeOptionsToStruct(attr.Options); err == nil {
		pbAttr.Options = options
	}
	if validationRules, err := MapToStruct(attr.ValidationRules); err == nil {
		pbAttr.ValidationRules = validationRules
	}
	if uiSettings, err := MapToStruct(attr.UISettings); err == nil {
		pbAttr.UiSettings = uiSettings
	}

	return pbAttr
}

// DomainToProtoCategoryAttribute converts domain.CategoryAttribute to proto CategoryAttribute
func DomainToProtoCategoryAttribute(ca *domain.CategoryAttribute) *pb.CategoryAttribute {
	if ca == nil {
		return nil
	}

	pbCatAttr := &pb.CategoryAttribute{
		Id:          ca.ID,
		CategoryId:  ca.CategoryID,
		AttributeId: ca.AttributeID,
		SortOrder:   ca.SortOrder,
		IsActive:    ca.IsActive,
		CreatedAt:   timestamppb.New(ca.CreatedAt),
		UpdatedAt:   timestamppb.New(ca.UpdatedAt),
	}

	// Embed attribute if loaded
	if ca.Attribute != nil {
		pbCatAttr.Attribute = DomainAttributeToProto(ca.Attribute)
	}

	// Nullable overrides
	if ca.IsRequired != nil {
		pbCatAttr.IsEnabled = &ca.IsActive
		pbCatAttr.IsRequired = ca.IsRequired
	}
	if ca.IsSearchable != nil {
		pbCatAttr.IsSearchable = ca.IsSearchable
	}
	if ca.IsFilterable != nil {
		pbCatAttr.IsFilterable = ca.IsFilterable
	}

	// JSONB fields
	if options, err := AttributeOptionsToStruct(ca.CategorySpecificOptions); err == nil {
		pbCatAttr.CategorySpecificOptions = options
	}
	if rules, err := MapToStruct(ca.CustomValidationRules); err == nil {
		pbCatAttr.CustomValidationRules = rules
	}
	if settings, err := MapToStruct(ca.CustomUISettings); err == nil {
		pbCatAttr.CustomUiSettings = settings
	}

	return pbCatAttr
}

// DomainToProtoListingAttributeValue converts domain.ListingAttributeValue to proto ListingAttributeValue
func DomainToProtoListingAttributeValue(lav *domain.ListingAttributeValue) *pb.ListingAttributeValue {
	if lav == nil {
		return nil
	}

	pbValue := &pb.ListingAttributeValue{
		Id:          lav.ID,
		ListingId:   lav.ListingID,
		AttributeId: lav.AttributeID,
		CreatedAt:   timestamppb.New(lav.CreatedAt),
		UpdatedAt:   timestamppb.New(lav.UpdatedAt),
	}

	// Embed attribute if loaded
	if lav.Attribute != nil {
		pbValue.Attribute = DomainAttributeToProto(lav.Attribute)
	}

	// Set value based on type (oneof)
	if lav.ValueText != nil {
		pbValue.Value = &pb.ListingAttributeValue_ValueText{ValueText: *lav.ValueText}
	} else if lav.ValueNumber != nil {
		pbValue.Value = &pb.ListingAttributeValue_ValueNumber{ValueNumber: *lav.ValueNumber}
	} else if lav.ValueBoolean != nil {
		pbValue.Value = &pb.ListingAttributeValue_ValueBoolean{ValueBoolean: *lav.ValueBoolean}
	} else if lav.ValueDate != nil {
		pbValue.Value = &pb.ListingAttributeValue_ValueDate{ValueDate: lav.ValueDate.Format(time.RFC3339)}
	} else if lav.ValueJSON != nil {
		if jsonStruct, err := MapToStruct(lav.ValueJSON); err == nil {
			pbValue.Value = &pb.ListingAttributeValue_ValueJson{ValueJson: jsonStruct}
		}
	}

	return pbValue
}

// DomainToProtoVariantAttribute converts domain.VariantAttribute to proto VariantAttribute
func DomainToProtoVariantAttribute(va *domain.VariantAttribute) *pb.VariantAttribute {
	if va == nil {
		return nil
	}

	pbVariant := &pb.VariantAttribute{
		Id:           va.ID,
		CategoryId:   va.CategoryID,
		AttributeId:  va.AttributeID,
		IsRequired:   va.IsRequired,
		AffectsPrice: va.AffectsPrice,
		AffectsStock: va.AffectsStock,
		SortOrder:    va.SortOrder,
		DisplayAs:    va.DisplayAs,
		IsActive:     va.IsActive,
		CreatedAt:    timestamppb.New(va.CreatedAt),
		UpdatedAt:    timestamppb.New(va.UpdatedAt),
	}

	// Embed attribute if loaded
	if va.Attribute != nil {
		pbVariant.Attribute = DomainAttributeToProto(va.Attribute)
	}

	return pbVariant
}

// ============================================================================
// Proto to Domain Conversions (Input DTOs)
// ============================================================================

// ProtoToCreateAttributeInput converts proto CreateAttributeRequest to domain CreateAttributeInput
func ProtoToCreateAttributeInput(req *pb.CreateAttributeRequest) (*domain.CreateAttributeInput, error) {
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	input := &domain.CreateAttributeInput{
		Code:                req.Code,
		Name:                StructToStringMap(req.Name),
		DisplayName:         StructToStringMap(req.DisplayName),
		AttributeType:       ProtoToAttributeType(req.AttributeType),
		Purpose:             ProtoToAttributePurpose(req.Purpose),
		ValidationRules:     StructToMap(req.ValidationRules),
		UISettings:          StructToMap(req.UiSettings),
		IsSearchable:        req.IsSearchable,
		IsFilterable:        req.IsFilterable,
		IsRequired:          req.IsRequired,
		IsVariantCompatible: req.IsVariantCompatible,
		AffectsStock:        req.AffectsStock,
		AffectsPrice:        req.AffectsPrice,
		ShowInCard:          req.ShowInCard,
		SortOrder:           req.SortOrder,
		Icon:                req.Icon,
	}

	// Convert options
	if req.Options != nil {
		options, err := StructToAttributeOptions(req.Options)
		if err != nil {
			return nil, fmt.Errorf("failed to convert options: %w", err)
		}
		input.Options = options
	}

	return input, nil
}

// ProtoToUpdateAttributeInput converts proto UpdateAttributeRequest to domain UpdateAttributeInput
func ProtoToUpdateAttributeInput(req *pb.UpdateAttributeRequest) *domain.UpdateAttributeInput {
	if req == nil {
		return nil
	}

	input := &domain.UpdateAttributeInput{}

	if req.Name != nil {
		nameMap := StructToStringMap(req.Name)
		input.Name = &nameMap
	}
	if req.DisplayName != nil {
		displayNameMap := StructToStringMap(req.DisplayName)
		input.DisplayName = &displayNameMap
	}
	if req.Options != nil {
		options, _ := StructToAttributeOptions(req.Options)
		input.Options = &options
	}
	if req.ValidationRules != nil {
		rules := StructToMap(req.ValidationRules)
		input.ValidationRules = &rules
	}
	if req.UiSettings != nil {
		settings := StructToMap(req.UiSettings)
		input.UISettings = &settings
	}
	if req.IsSearchable != nil {
		input.IsSearchable = req.IsSearchable
	}
	if req.IsFilterable != nil {
		input.IsFilterable = req.IsFilterable
	}
	if req.IsRequired != nil {
		input.IsRequired = req.IsRequired
	}
	if req.ShowInCard != nil {
		input.ShowInCard = req.ShowInCard
	}
	if req.SortOrder != nil {
		input.SortOrder = req.SortOrder
	}
	if req.Icon != nil {
		input.Icon = req.Icon
	}
	if req.IsActive != nil {
		input.IsActive = req.IsActive
	}

	return input
}

// ProtoToCategoryAttributeSettings converts proto LinkAttributeToCategoryRequest to domain CategoryAttributeSettings
func ProtoToCategoryAttributeSettings(req *pb.LinkAttributeToCategoryRequest) *domain.CategoryAttributeSettings {
	if req == nil {
		return nil
	}

	settings := &domain.CategoryAttributeSettings{
		IsEnabled: true, // Default to enabled when linking
		SortOrder: req.SortOrder,
	}

	if req.IsEnabled != nil {
		settings.IsEnabled = *req.IsEnabled
	}
	if req.IsRequired != nil {
		settings.IsRequired = req.IsRequired
	}
	if req.IsSearchable != nil {
		settings.IsSearchable = req.IsSearchable
	}
	if req.IsFilterable != nil {
		settings.IsFilterable = req.IsFilterable
	}
	if req.CategorySpecificOptions != nil {
		options, _ := StructToAttributeOptions(req.CategorySpecificOptions)
		settings.CategorySpecificOptions = &options
	}
	if req.CustomValidationRules != nil {
		rules := StructToMap(req.CustomValidationRules)
		settings.CustomValidationRules = &rules
	}
	if req.CustomUiSettings != nil {
		uiSettings := StructToMap(req.CustomUiSettings)
		settings.CustomUISettings = &uiSettings
	}

	return settings
}

// ProtoToUpdateCategoryAttributeSettings converts proto UpdateCategoryAttributeRequest to domain CategoryAttributeSettings
func ProtoToUpdateCategoryAttributeSettings(req *pb.UpdateCategoryAttributeRequest) *domain.CategoryAttributeSettings {
	if req == nil {
		return nil
	}

	settings := &domain.CategoryAttributeSettings{}

	if req.IsEnabled != nil {
		settings.IsEnabled = *req.IsEnabled
	}
	if req.IsRequired != nil {
		settings.IsRequired = req.IsRequired
	}
	if req.IsSearchable != nil {
		settings.IsSearchable = req.IsSearchable
	}
	if req.IsFilterable != nil {
		settings.IsFilterable = req.IsFilterable
	}
	if req.SortOrder != nil {
		settings.SortOrder = *req.SortOrder
	}
	if req.CategorySpecificOptions != nil {
		options, _ := StructToAttributeOptions(req.CategorySpecificOptions)
		settings.CategorySpecificOptions = &options
	}
	if req.CustomValidationRules != nil {
		rules := StructToMap(req.CustomValidationRules)
		settings.CustomValidationRules = &rules
	}
	if req.CustomUiSettings != nil {
		uiSettings := StructToMap(req.CustomUiSettings)
		settings.CustomUISettings = &uiSettings
	}

	return settings
}

// ProtoToSetListingAttributeValue converts proto AttributeValueInput to domain SetListingAttributeValue
func ProtoToSetListingAttributeValue(input *pb.AttributeValueInput) (*domain.SetListingAttributeValue, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}

	value := &domain.SetListingAttributeValue{
		AttributeID: input.AttributeId,
	}

	// Extract value from oneof
	switch v := input.Value.(type) {
	case *pb.AttributeValueInput_ValueText:
		value.ValueText = &v.ValueText
	case *pb.AttributeValueInput_ValueNumber:
		value.ValueNumber = &v.ValueNumber
	case *pb.AttributeValueInput_ValueBoolean:
		value.ValueBoolean = &v.ValueBoolean
	case *pb.AttributeValueInput_ValueDate:
		// Parse ISO 8601 date
		dateVal, err := time.Parse(time.RFC3339, v.ValueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		value.ValueDate = &dateVal
	case *pb.AttributeValueInput_ValueJson:
		value.ValueJSON = StructToMap(v.ValueJson)
	default:
		return nil, fmt.Errorf("no value provided")
	}

	return value, nil
}

// ProtoAttributeToCreateInput converts proto Attribute (for bulk import) to domain CreateAttributeInput
func ProtoAttributeToCreateInput(pbAttr *pb.Attribute) (*domain.CreateAttributeInput, error) {
	if pbAttr == nil {
		return nil, fmt.Errorf("attribute cannot be nil")
	}

	input := &domain.CreateAttributeInput{
		Code:                pbAttr.Code,
		Name:                StructToStringMap(pbAttr.Name),
		DisplayName:         StructToStringMap(pbAttr.DisplayName),
		AttributeType:       ProtoToAttributeType(pbAttr.AttributeType),
		Purpose:             ProtoToAttributePurpose(pbAttr.Purpose),
		ValidationRules:     StructToMap(pbAttr.ValidationRules),
		UISettings:          StructToMap(pbAttr.UiSettings),
		IsSearchable:        pbAttr.IsSearchable,
		IsFilterable:        pbAttr.IsFilterable,
		IsRequired:          pbAttr.IsRequired,
		IsVariantCompatible: pbAttr.IsVariantCompatible,
		AffectsStock:        pbAttr.AffectsStock,
		AffectsPrice:        pbAttr.AffectsPrice,
		ShowInCard:          pbAttr.ShowInCard,
		SortOrder:           pbAttr.SortOrder,
		Icon:                pbAttr.Icon,
	}

	// Convert options
	if pbAttr.Options != nil {
		options, err := StructToAttributeOptions(pbAttr.Options)
		if err != nil {
			return nil, fmt.Errorf("failed to convert options: %w", err)
		}
		input.Options = options
	}

	return input, nil
}
