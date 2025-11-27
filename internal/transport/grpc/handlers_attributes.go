package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/vondi-global/listings/api/proto/attributes/v1"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/service"
)

// AttributeServer implements gRPC AttributeServiceServer
// This is embedded in the main Server struct which implements both ListingsService and AttributeService
type AttributeServer struct {
	pb.UnimplementedAttributeServiceServer
	attrService service.AttributeService
}

// CreateAttribute creates a new attribute definition
func (s *Server) CreateAttribute(ctx context.Context, req *pb.CreateAttributeRequest) (*pb.CreateAttributeResponse, error) {
	s.logger.Debug().
		Str("code", req.Code).
		Str("attribute_type", req.AttributeType.String()).
		Msg("CreateAttribute called")

	// Validate request
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}
	if req.AttributeType == pb.AttributeType_ATTRIBUTE_TYPE_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "attribute_type is required")
	}

	// Convert proto to domain input
	input, err := ProtoToCreateAttributeInput(req)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to convert proto to domain input")
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid input: %v", err))
	}

	// Call service layer
	attr, err := s.attrService.CreateAttribute(ctx, input)
	if err != nil {
		s.logger.Error().Err(err).Str("code", req.Code).Msg("failed to create attribute")
		return nil, convertServiceError(err, "create attribute")
	}

	// Convert domain to proto response
	pbAttr := DomainAttributeToProto(attr)

	s.logger.Info().Int32("id", attr.ID).Str("code", attr.Code).Msg("attribute created successfully")
	return &pb.CreateAttributeResponse{
		Attribute: pbAttr,
	}, nil
}

// UpdateAttribute updates an existing attribute
func (s *Server) UpdateAttribute(ctx context.Context, req *pb.UpdateAttributeRequest) (*pb.UpdateAttributeResponse, error) {
	s.logger.Debug().
		Int32("id", req.Id).
		Msg("UpdateAttribute called")

	// Validate request
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be greater than 0")
	}

	// Convert proto to domain input
	input := ProtoToUpdateAttributeInput(req)

	// Call service layer
	attr, err := s.attrService.UpdateAttribute(ctx, req.Id, input)
	if err != nil {
		s.logger.Error().Err(err).Int32("id", req.Id).Msg("failed to update attribute")
		return nil, convertServiceError(err, "update attribute")
	}

	// Convert domain to proto response
	pbAttr := DomainAttributeToProto(attr)

	s.logger.Info().Int32("id", attr.ID).Msg("attribute updated successfully")
	return &pb.UpdateAttributeResponse{
		Attribute: pbAttr,
	}, nil
}

// DeleteAttribute soft-deletes an attribute (sets is_active=false)
func (s *Server) DeleteAttribute(ctx context.Context, req *pb.DeleteAttributeRequest) (*pb.DeleteAttributeResponse, error) {
	s.logger.Debug().
		Int32("id", req.Id).
		Msg("DeleteAttribute called")

	// Validate request
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be greater than 0")
	}

	// Call service layer
	err := s.attrService.DeleteAttribute(ctx, req.Id)
	if err != nil {
		s.logger.Error().Err(err).Int32("id", req.Id).Msg("failed to delete attribute")
		return nil, convertServiceError(err, "delete attribute")
	}

	s.logger.Info().Int32("id", req.Id).Msg("attribute deleted successfully")
	return &pb.DeleteAttributeResponse{
		Success: true,
		Message: "Attribute deleted successfully",
	}, nil
}

// GetAttribute retrieves a single attribute by ID or code
func (s *Server) GetAttribute(ctx context.Context, req *pb.GetAttributeRequest) (*pb.GetAttributeResponse, error) {
	s.logger.Debug().Msg("GetAttribute called")

	var identifier string
	var attr *domain.Attribute
	var err error

	// Handle oneof identifier (id or code)
	switch req.Identifier.(type) {
	case *pb.GetAttributeRequest_Id:
		id := req.GetId()
		if id <= 0 {
			return nil, status.Error(codes.InvalidArgument, "id must be greater than 0")
		}
		s.logger.Debug().Int32("id", id).Msg("getting attribute by ID")
		attr, err = s.attrService.GetAttributeByID(ctx, id)
		identifier = fmt.Sprintf("id=%d", id)

	case *pb.GetAttributeRequest_Code:
		code := req.GetCode()
		if code == "" {
			return nil, status.Error(codes.InvalidArgument, "code cannot be empty")
		}
		s.logger.Debug().Str("code", code).Msg("getting attribute by code")
		attr, err = s.attrService.GetAttributeByCode(ctx, code)
		identifier = fmt.Sprintf("code=%s", code)

	default:
		return nil, status.Error(codes.InvalidArgument, "either id or code must be provided")
	}

	if err != nil {
		s.logger.Error().Err(err).Str("identifier", identifier).Msg("failed to get attribute")
		return nil, convertServiceError(err, "get attribute")
	}

	if attr == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("attribute not found: %s", identifier))
	}

	// Convert domain to proto response
	pbAttr := DomainAttributeToProto(attr)

	s.logger.Debug().Int32("id", attr.ID).Str("identifier", identifier).Msg("attribute retrieved successfully")
	return &pb.GetAttributeResponse{
		Attribute: pbAttr,
	}, nil
}

// ListAttributes lists all attributes with optional filtering
func (s *Server) ListAttributes(ctx context.Context, req *pb.ListAttributesRequest) (*pb.ListAttributesResponse, error) {
	s.logger.Debug().
		Int32("page", req.Page).
		Int32("page_size", req.PageSize).
		Msg("ListAttributes called")

	// Set defaults
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Convert proto filter to domain filter
	filter := &domain.ListAttributesFilter{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}

	// Apply optional filters
	if req.AttributeType != nil && *req.AttributeType != pb.AttributeType_ATTRIBUTE_TYPE_UNSPECIFIED {
		attrType := ProtoToAttributeType(*req.AttributeType)
		filter.AttributeType = &attrType
	}
	if req.Purpose != nil && *req.Purpose != pb.AttributePurpose_ATTRIBUTE_PURPOSE_UNSPECIFIED {
		purpose := ProtoToAttributePurpose(*req.Purpose)
		filter.Purpose = &purpose
	}
	if req.IsActive != nil {
		filter.IsActive = req.IsActive
	}
	if req.IsSearchable != nil {
		filter.IsSearchable = req.IsSearchable
	}
	if req.IsFilterable != nil {
		filter.IsFilterable = req.IsFilterable
	}

	// Call service layer
	attributes, total, err := s.attrService.ListAttributes(ctx, filter)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list attributes")
		return nil, convertServiceError(err, "list attributes")
	}

	// Convert domain to proto
	pbAttrs := make([]*pb.Attribute, len(attributes))
	for i, attr := range attributes {
		pbAttrs[i] = DomainAttributeToProto(attr)
	}

	// Calculate pagination
	totalPages := int32((total + int64(pageSize) - 1) / int64(pageSize))

	s.logger.Info().
		Int("count", len(attributes)).
		Int64("total", total).
		Msg("attributes listed successfully")

	return &pb.ListAttributesResponse{
		Attributes: pbAttrs,
		TotalCount: int32(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// LinkAttributeToCategory links an attribute to a category with optional overrides
func (s *Server) LinkAttributeToCategory(ctx context.Context, req *pb.LinkAttributeToCategoryRequest) (*pb.LinkAttributeToCategoryResponse, error) {
	s.logger.Debug().
		Int32("category_id", req.CategoryId).
		Int32("attribute_id", req.AttributeId).
		Msg("LinkAttributeToCategory called")

	// Validate request
	if req.CategoryId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "category_id must be greater than 0")
	}
	if req.AttributeId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "attribute_id must be greater than 0")
	}

	// Convert proto to domain settings
	settings := ProtoToCategoryAttributeSettings(req)

	// Call service layer
	err := s.attrService.LinkAttributeToCategory(ctx, req.CategoryId, req.AttributeId, settings)
	if err != nil {
		s.logger.Error().Err(err).
			Int32("category_id", req.CategoryId).
			Int32("attribute_id", req.AttributeId).
			Msg("failed to link attribute to category")
		return nil, convertServiceError(err, "link attribute to category")
	}

	// Get the created relationship
	catAttrs, err := s.attrService.GetCategoryAttributes(ctx, req.CategoryId, &domain.GetCategoryAttributesFilter{})
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get category attributes after linking")
		return nil, convertServiceError(err, "get category attributes")
	}

	// Find the newly linked attribute
	var linkedAttr *domain.CategoryAttribute
	for _, ca := range catAttrs {
		if ca.AttributeID == req.AttributeId {
			linkedAttr = ca
			break
		}
	}

	if linkedAttr == nil {
		return nil, status.Error(codes.Internal, "failed to retrieve linked attribute")
	}

	// Convert to proto
	pbCatAttr := DomainToProtoCategoryAttribute(linkedAttr)

	s.logger.Info().
		Int32("category_id", req.CategoryId).
		Int32("attribute_id", req.AttributeId).
		Msg("attribute linked to category successfully")

	return &pb.LinkAttributeToCategoryResponse{
		CategoryAttribute: pbCatAttr,
	}, nil
}

// UpdateCategoryAttribute updates category-specific attribute settings
func (s *Server) UpdateCategoryAttribute(ctx context.Context, req *pb.UpdateCategoryAttributeRequest) (*pb.UpdateCategoryAttributeResponse, error) {
	s.logger.Debug().
		Int32("category_id", req.CategoryId).
		Int32("attribute_id", req.AttributeId).
		Msg("UpdateCategoryAttribute called")

	// Validate request
	if req.CategoryId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "category_id must be greater than 0")
	}
	if req.AttributeId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "attribute_id must be greater than 0")
	}

	// Get existing category attribute to find the ID
	catAttrs, err := s.attrService.GetCategoryAttributes(ctx, req.CategoryId, &domain.GetCategoryAttributesFilter{})
	if err != nil {
		return nil, convertServiceError(err, "get category attributes")
	}

	var catAttrID int32
	for _, ca := range catAttrs {
		if ca.AttributeID == req.AttributeId {
			catAttrID = ca.ID
			break
		}
	}

	if catAttrID == 0 {
		return nil, status.Error(codes.NotFound, "category attribute relationship not found")
	}

	// Convert proto to domain settings
	settings := ProtoToUpdateCategoryAttributeSettings(req)

	// Call service layer
	err = s.attrService.UpdateCategoryAttribute(ctx, catAttrID, settings)
	if err != nil {
		s.logger.Error().Err(err).
			Int32("category_id", req.CategoryId).
			Int32("attribute_id", req.AttributeId).
			Msg("failed to update category attribute")
		return nil, convertServiceError(err, "update category attribute")
	}

	// Get updated relationship
	catAttrs, err = s.attrService.GetCategoryAttributes(ctx, req.CategoryId, &domain.GetCategoryAttributesFilter{})
	if err != nil {
		return nil, convertServiceError(err, "get updated category attributes")
	}

	var updatedAttr *domain.CategoryAttribute
	for _, ca := range catAttrs {
		if ca.AttributeID == req.AttributeId {
			updatedAttr = ca
			break
		}
	}

	if updatedAttr == nil {
		return nil, status.Error(codes.Internal, "failed to retrieve updated attribute")
	}

	// Convert to proto
	pbCatAttr := DomainToProtoCategoryAttribute(updatedAttr)

	s.logger.Info().
		Int32("category_id", req.CategoryId).
		Int32("attribute_id", req.AttributeId).
		Msg("category attribute updated successfully")

	return &pb.UpdateCategoryAttributeResponse{
		CategoryAttribute: pbCatAttr,
	}, nil
}

// UnlinkAttributeFromCategory removes the attribute-category relationship
func (s *Server) UnlinkAttributeFromCategory(ctx context.Context, req *pb.UnlinkAttributeFromCategoryRequest) (*pb.UnlinkAttributeFromCategoryResponse, error) {
	s.logger.Debug().
		Int32("category_id", req.CategoryId).
		Int32("attribute_id", req.AttributeId).
		Msg("UnlinkAttributeFromCategory called")

	// Validate request
	if req.CategoryId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "category_id must be greater than 0")
	}
	if req.AttributeId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "attribute_id must be greater than 0")
	}

	// Call service layer
	err := s.attrService.UnlinkAttributeFromCategory(ctx, req.CategoryId, req.AttributeId)
	if err != nil {
		s.logger.Error().Err(err).
			Int32("category_id", req.CategoryId).
			Int32("attribute_id", req.AttributeId).
			Msg("failed to unlink attribute from category")
		return nil, convertServiceError(err, "unlink attribute from category")
	}

	s.logger.Info().
		Int32("category_id", req.CategoryId).
		Int32("attribute_id", req.AttributeId).
		Msg("attribute unlinked from category successfully")

	return &pb.UnlinkAttributeFromCategoryResponse{
		Success: true,
		Message: "Attribute unlinked from category successfully",
	}, nil
}

// GetCategoryAttributes retrieves all attributes for a specific category
func (s *Server) GetCategoryAttributes(ctx context.Context, req *pb.GetCategoryAttributesRequest) (*pb.GetCategoryAttributesResponse, error) {
	s.logger.Debug().
		Int32("category_id", req.CategoryId).
		Bool("include_inactive", req.IncludeInactive).
		Msg("GetCategoryAttributes called")

	// Validate request
	if req.CategoryId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "category_id must be greater than 0")
	}

	// Build filter
	filter := &domain.GetCategoryAttributesFilter{}
	if !req.IncludeInactive {
		enabled := true
		filter.IsEnabled = &enabled
	}

	// Call service layer
	catAttrs, err := s.attrService.GetCategoryAttributes(ctx, req.CategoryId, filter)
	if err != nil {
		s.logger.Error().Err(err).Int32("category_id", req.CategoryId).Msg("failed to get category attributes")
		return nil, convertServiceError(err, "get category attributes")
	}

	// Convert domain to proto
	pbCatAttrs := make([]*pb.CategoryAttribute, len(catAttrs))
	for i, ca := range catAttrs {
		pbCatAttrs[i] = DomainToProtoCategoryAttribute(ca)
	}

	s.logger.Info().
		Int32("category_id", req.CategoryId).
		Int("count", len(catAttrs)).
		Msg("category attributes retrieved successfully")

	return &pb.GetCategoryAttributesResponse{
		CategoryAttributes: pbCatAttrs,
	}, nil
}

// GetCategoryVariantAttributes retrieves variant attributes for a category
func (s *Server) GetCategoryVariantAttributes(ctx context.Context, req *pb.GetCategoryVariantAttributesRequest) (*pb.GetCategoryVariantAttributesResponse, error) {
	s.logger.Debug().
		Int32("category_id", req.CategoryId).
		Bool("include_inactive", req.IncludeInactive).
		Msg("GetCategoryVariantAttributes called")

	// Validate request
	if req.CategoryId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "category_id must be greater than 0")
	}

	// Call service layer
	variantAttrs, err := s.attrService.GetCategoryVariantAttributes(ctx, req.CategoryId)
	if err != nil {
		s.logger.Error().Err(err).Int32("category_id", req.CategoryId).Msg("failed to get variant attributes")
		return nil, convertServiceError(err, "get variant attributes")
	}

	// Filter inactive if requested
	if !req.IncludeInactive {
		activeVariantAttrs := make([]*domain.VariantAttribute, 0)
		for _, va := range variantAttrs {
			if va.IsActive {
				activeVariantAttrs = append(activeVariantAttrs, va)
			}
		}
		variantAttrs = activeVariantAttrs
	}

	// Convert domain to proto
	pbVariantAttrs := make([]*pb.VariantAttribute, len(variantAttrs))
	for i, va := range variantAttrs {
		pbVariantAttrs[i] = DomainToProtoVariantAttribute(va)
	}

	s.logger.Info().
		Int32("category_id", req.CategoryId).
		Int("count", len(variantAttrs)).
		Msg("variant attributes retrieved successfully")

	return &pb.GetCategoryVariantAttributesResponse{
		VariantAttributes: pbVariantAttrs,
	}, nil
}

// GetListingAttributes retrieves all attribute values for a listing
func (s *Server) GetListingAttributes(ctx context.Context, req *pb.GetListingAttributesRequest) (*pb.GetListingAttributesResponse, error) {
	s.logger.Debug().
		Int32("listing_id", req.ListingId).
		Msg("GetListingAttributes called")

	// Validate request
	if req.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing_id must be greater than 0")
	}

	// Call service layer
	attrValues, err := s.attrService.GetListingAttributes(ctx, req.ListingId)
	if err != nil {
		s.logger.Error().Err(err).Int32("listing_id", req.ListingId).Msg("failed to get listing attributes")
		return nil, convertServiceError(err, "get listing attributes")
	}

	// Filter by specific attribute IDs if requested
	if len(req.AttributeIds) > 0 {
		filteredValues := make([]*domain.ListingAttributeValue, 0)
		attrIDMap := make(map[int32]bool)
		for _, id := range req.AttributeIds {
			attrIDMap[id] = true
		}
		for _, av := range attrValues {
			if attrIDMap[av.AttributeID] {
				filteredValues = append(filteredValues, av)
			}
		}
		attrValues = filteredValues
	}

	// Convert domain to proto
	pbAttrValues := make([]*pb.ListingAttributeValue, len(attrValues))
	for i, av := range attrValues {
		pbAttrValues[i] = DomainToProtoListingAttributeValue(av)
	}

	s.logger.Info().
		Int32("listing_id", req.ListingId).
		Int("count", len(attrValues)).
		Msg("listing attributes retrieved successfully")

	return &pb.GetListingAttributesResponse{
		AttributeValues: pbAttrValues,
	}, nil
}

// SetListingAttributes sets/updates attribute values for a listing
func (s *Server) SetListingAttributes(ctx context.Context, req *pb.SetListingAttributesRequest) (*pb.SetListingAttributesResponse, error) {
	s.logger.Debug().
		Int32("listing_id", req.ListingId).
		Int("values_count", len(req.Values)).
		Msg("SetListingAttributes called")

	// Validate request
	if req.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing_id must be greater than 0")
	}
	if len(req.Values) == 0 {
		return nil, status.Error(codes.InvalidArgument, "values cannot be empty")
	}

	// Convert proto to domain values
	domainValues := make([]domain.SetListingAttributeValue, len(req.Values))
	for i, v := range req.Values {
		domainValue, err := ProtoToSetListingAttributeValue(v)
		if err != nil {
			s.logger.Warn().Err(err).Int("index", i).Msg("failed to convert attribute value")
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid value at index %d: %v", i, err))
		}
		domainValues[i] = *domainValue
	}

	// Call service layer
	err := s.attrService.SetListingAttributes(ctx, req.ListingId, domainValues)
	if err != nil {
		s.logger.Error().Err(err).Int32("listing_id", req.ListingId).Msg("failed to set listing attributes")
		return nil, convertServiceError(err, "set listing attributes")
	}

	// Get updated attributes
	attrValues, err := s.attrService.GetListingAttributes(ctx, req.ListingId)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get updated listing attributes")
		return nil, convertServiceError(err, "get updated listing attributes")
	}

	// Convert domain to proto
	pbAttrValues := make([]*pb.ListingAttributeValue, len(attrValues))
	for i, av := range attrValues {
		pbAttrValues[i] = DomainToProtoListingAttributeValue(av)
	}

	s.logger.Info().
		Int32("listing_id", req.ListingId).
		Int("count", len(attrValues)).
		Msg("listing attributes set successfully")

	return &pb.SetListingAttributesResponse{
		AttributeValues:  pbAttrValues,
		ValidationErrors: []*pb.ValidationError{}, // No errors if we got here
	}, nil
}

// ValidateAttributeValues validates attribute values before saving
func (s *Server) ValidateAttributeValues(ctx context.Context, req *pb.ValidateAttributeValuesRequest) (*pb.ValidateAttributeValuesResponse, error) {
	s.logger.Debug().
		Int32("category_id", req.CategoryId).
		Int("values_count", len(req.Values)).
		Msg("ValidateAttributeValues called")

	// Validate request
	if req.CategoryId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "category_id must be greater than 0")
	}

	// Convert proto to domain values
	domainValues := make([]domain.SetListingAttributeValue, len(req.Values))
	for i, v := range req.Values {
		domainValue, err := ProtoToSetListingAttributeValue(v)
		if err != nil {
			s.logger.Warn().Err(err).Int("index", i).Msg("failed to convert attribute value")
			return &pb.ValidateAttributeValuesResponse{
				IsValid: false,
				Errors: []*pb.ValidationError{
					{
						AttributeId:  v.AttributeId,
						ErrorCode:    "invalid_value",
						ErrorMessage: fmt.Sprintf("Invalid value format: %v", err),
					},
				},
			}, nil
		}
		domainValues[i] = *domainValue
	}

	// Call service layer
	err := s.attrService.ValidateAttributeValues(ctx, req.CategoryId, domainValues)
	if err != nil {
		// Validation failed - return validation errors
		s.logger.Debug().Err(err).Msg("validation failed")

		return &pb.ValidateAttributeValuesResponse{
			IsValid: false,
			Errors: []*pb.ValidationError{
				{
					ErrorCode:    "validation_failed",
					ErrorMessage: err.Error(),
				},
			},
		}, nil
	}

	// Validation successful
	s.logger.Info().
		Int32("category_id", req.CategoryId).
		Int("values_count", len(req.Values)).
		Msg("attribute values validated successfully")

	return &pb.ValidateAttributeValuesResponse{
		IsValid: true,
		Errors:  []*pb.ValidationError{},
	}, nil
}

// BulkImportAttributes imports multiple attributes at once (for migration)
func (s *Server) BulkImportAttributes(ctx context.Context, req *pb.BulkImportAttributesRequest) (*pb.BulkImportAttributesResponse, error) {
	s.logger.Info().
		Int("total_attributes", len(req.Attributes)).
		Bool("skip_duplicates", req.SkipDuplicates).
		Bool("update_existing", req.UpdateExisting).
		Msg("BulkImportAttributes called")

	var totalImported int32
	var totalSkipped int32
	var totalUpdated int32
	var totalFailed int32
	var errorMessages []string

	for i, pbAttr := range req.Attributes {
		// Convert proto to create input
		input, err := ProtoAttributeToCreateInput(pbAttr)
		if err != nil {
			s.logger.Warn().Err(err).Int("index", i).Msg("failed to convert attribute")
			totalFailed++
			errorMessages = append(errorMessages, fmt.Sprintf("Index %d (%s): %v", i, pbAttr.Code, err))
			continue
		}

		// Check if attribute exists
		existing, err := s.attrService.GetAttributeByCode(ctx, input.Code)
		if err == nil && existing != nil {
			// Attribute exists
			if req.SkipDuplicates && !req.UpdateExisting {
				s.logger.Debug().Str("code", input.Code).Msg("skipping duplicate attribute")
				totalSkipped++
				continue
			}

			if req.UpdateExisting {
				// Update existing
				updateInput := CreateToUpdateInput(input)
				_, err = s.attrService.UpdateAttribute(ctx, existing.ID, updateInput)
				if err != nil {
					s.logger.Error().Err(err).Str("code", input.Code).Msg("failed to update attribute")
					totalFailed++
					errorMessages = append(errorMessages, fmt.Sprintf("Update %s: %v", input.Code, err))
					continue
				}
				totalUpdated++
				s.logger.Debug().Str("code", input.Code).Msg("attribute updated")
				continue
			}
		}

		// Create new attribute
		_, err = s.attrService.CreateAttribute(ctx, input)
		if err != nil {
			s.logger.Error().Err(err).Str("code", input.Code).Msg("failed to create attribute")
			totalFailed++
			errorMessages = append(errorMessages, fmt.Sprintf("Create %s: %v", input.Code, err))
			continue
		}

		totalImported++
		s.logger.Debug().Str("code", input.Code).Msg("attribute imported")
	}

	s.logger.Info().
		Int32("imported", totalImported).
		Int32("updated", totalUpdated).
		Int32("skipped", totalSkipped).
		Int32("failed", totalFailed).
		Msg("bulk import completed")

	return &pb.BulkImportAttributesResponse{
		TotalImported: totalImported,
		TotalSkipped:  totalSkipped,
		TotalUpdated:  totalUpdated,
		TotalFailed:   totalFailed,
		ErrorMessages: errorMessages,
	}, nil
}

// convertServiceError converts service layer errors to gRPC status codes
func convertServiceError(err error, operation string) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// NotFound errors
	if contains(errMsg, "not found") {
		return status.Error(codes.NotFound, fmt.Sprintf("%s: %v", operation, err))
	}

	// Validation errors
	if contains(errMsg, "validation") || contains(errMsg, "invalid") {
		return status.Error(codes.InvalidArgument, fmt.Sprintf("%s: %v", operation, err))
	}

	// Already exists errors
	if contains(errMsg, "already exists") || contains(errMsg, "duplicate") {
		return status.Error(codes.AlreadyExists, fmt.Sprintf("%s: %v", operation, err))
	}

	// Permission errors
	if contains(errMsg, "permission") || contains(errMsg, "unauthorized") {
		return status.Error(codes.PermissionDenied, fmt.Sprintf("%s: %v", operation, err))
	}

	// Default to Internal error
	return status.Error(codes.Internal, fmt.Sprintf("%s: %v", operation, err))
}

// CreateToUpdateInput converts CreateAttributeInput to UpdateAttributeInput for bulk import
func CreateToUpdateInput(input *domain.CreateAttributeInput) *domain.UpdateAttributeInput {
	return &domain.UpdateAttributeInput{
		Name:            &input.Name,
		DisplayName:     &input.DisplayName,
		Options:         &input.Options,
		ValidationRules: &input.ValidationRules,
		UISettings:      &input.UISettings,
		IsSearchable:    &input.IsSearchable,
		IsFilterable:    &input.IsFilterable,
		IsRequired:      &input.IsRequired,
		ShowInCard:      &input.ShowInCard,
		SortOrder:       &input.SortOrder,
		Icon:            &input.Icon,
	}
}
