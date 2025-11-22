// Package service implements business logic for the listings microservice.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/repository"
)

// AttributeServiceImpl implements AttributeService interface
type AttributeServiceImpl struct {
	repo         repository.AttributeRepository
	cache        *AttributeCache
	validator    *AttributeValidator
	stdValidator *validator.Validate
	logger       zerolog.Logger
}

// NewAttributeService creates a new attribute service
func NewAttributeService(
	repo repository.AttributeRepository,
	redisClient redis.UniversalClient,
	logger zerolog.Logger,
) AttributeService {
	stdValidator := validator.New()

	return &AttributeServiceImpl{
		repo:         repo,
		cache:        NewAttributeCache(redisClient, logger),
		validator:    NewAttributeValidator(),
		stdValidator: stdValidator,
		logger:       logger.With().Str("component", "attribute_service").Logger(),
	}
}

// CreateAttribute creates a new attribute with comprehensive validation
func (s *AttributeServiceImpl) CreateAttribute(ctx context.Context, input *domain.CreateAttributeInput) (*domain.Attribute, error) {
	// 1. Validate input using struct validator
	if err := s.stdValidator.Struct(input); err != nil {
		s.logger.Warn().Err(err).Msg("attribute input validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 2. Business logic validation
	if err := s.validateAttributeInput(input); err != nil {
		s.logger.Warn().Err(err).Msg("attribute business validation failed")
		return nil, err
	}

	// 3. Check code uniqueness
	existing, err := s.repo.GetByCode(ctx, input.Code)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("attribute with code '%s' already exists", input.Code)
	}

	// 4. Create attribute in repository
	attr, err := s.repo.Create(ctx, input)
	if err != nil {
		s.logger.Error().Err(err).Str("code", input.Code).Msg("failed to create attribute")
		return nil, fmt.Errorf("failed to create attribute: %w", err)
	}

	// 5. Cache the new attribute
	if err := s.cache.SetAttribute(ctx, attr); err != nil {
		s.logger.Warn().Err(err).Int32("id", attr.ID).Msg("failed to cache new attribute")
		// Don't fail on cache error
	}

	s.logger.Info().Int32("id", attr.ID).Str("code", attr.Code).Msg("attribute created successfully")
	return attr, nil
}

// UpdateAttribute updates an existing attribute with validation
func (s *AttributeServiceImpl) UpdateAttribute(ctx context.Context, id int32, input *domain.UpdateAttributeInput) (*domain.Attribute, error) {
	// 1. Validate input
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}

	// 2. Get existing attribute
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("attribute with id %d not found", id)
	}

	// 3. Business validation for update
	if err := s.validateAttributeUpdate(existing, input); err != nil {
		s.logger.Warn().Err(err).Int32("id", id).Msg("attribute update validation failed")
		return nil, err
	}

	// 4. Update attribute in repository
	updated, err := s.repo.Update(ctx, id, input)
	if err != nil {
		s.logger.Error().Err(err).Int32("id", id).Msg("failed to update attribute")
		return nil, fmt.Errorf("failed to update attribute: %w", err)
	}

	// 5. Invalidate cache
	if err := s.cache.InvalidateAttribute(ctx, id, existing.Code); err != nil {
		s.logger.Warn().Err(err).Int32("id", id).Msg("failed to invalidate attribute cache")
	}

	s.logger.Info().Int32("id", id).Msg("attribute updated successfully")
	return updated, nil
}

// DeleteAttribute soft-deletes an attribute
func (s *AttributeServiceImpl) DeleteAttribute(ctx context.Context, id int32) error {
	// 1. Get existing attribute to get code for cache invalidation
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get attribute: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("attribute with id %d not found", id)
	}

	// 2. Soft delete (set is_active=false)
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error().Err(err).Int32("id", id).Msg("failed to delete attribute")
		return fmt.Errorf("failed to delete attribute: %w", err)
	}

	// 3. Invalidate cache
	if err := s.cache.InvalidateAttribute(ctx, id, existing.Code); err != nil {
		s.logger.Warn().Err(err).Int32("id", id).Msg("failed to invalidate attribute cache")
	}

	s.logger.Info().Int32("id", id).Str("code", existing.Code).Msg("attribute deleted successfully")
	return nil
}

// GetAttribute retrieves an attribute by ID or code
func (s *AttributeServiceImpl) GetAttribute(ctx context.Context, identifier string) (*domain.Attribute, error) {
	// Try to parse as ID first
	if id, err := strconv.ParseInt(identifier, 10, 32); err == nil {
		return s.GetAttributeByID(ctx, int32(id))
	}

	// Otherwise treat as code
	return s.GetAttributeByCode(ctx, identifier)
}

// GetAttributeByID retrieves an attribute by ID (cached)
func (s *AttributeServiceImpl) GetAttributeByID(ctx context.Context, id int32) (*domain.Attribute, error) {
	// 1. Try cache first
	cached, err := s.cache.GetAttribute(ctx, id)
	if err == nil && cached != nil {
		return cached, nil
	}

	// 2. Cache miss - fetch from repository
	attr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute by id: %w", err)
	}
	if attr == nil {
		return nil, fmt.Errorf("attribute with id %d not found", id)
	}

	// 3. Cache the result
	if err := s.cache.SetAttribute(ctx, attr); err != nil {
		s.logger.Warn().Err(err).Int32("id", id).Msg("failed to cache attribute")
	}

	return attr, nil
}

// GetAttributeByCode retrieves an attribute by code (cached)
func (s *AttributeServiceImpl) GetAttributeByCode(ctx context.Context, code string) (*domain.Attribute, error) {
	// 1. Try cache first
	cached, err := s.cache.GetAttributeByCode(ctx, code)
	if err == nil && cached != nil {
		return cached, nil
	}

	// 2. Cache miss - fetch from repository
	attr, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute by code: %w", err)
	}
	if attr == nil {
		return nil, fmt.Errorf("attribute with code '%s' not found", code)
	}

	// 3. Cache the result
	if err := s.cache.SetAttribute(ctx, attr); err != nil {
		s.logger.Warn().Err(err).Str("code", code).Msg("failed to cache attribute")
	}

	return attr, nil
}

// ListAttributes lists attributes with filters and pagination
func (s *AttributeServiceImpl) ListAttributes(ctx context.Context, filter *domain.ListAttributesFilter) ([]*domain.Attribute, int64, error) {
	// Validate filter
	if filter == nil {
		filter = &domain.ListAttributesFilter{
			Limit:  20,
			Offset: 0,
		}
	}

	// Ensure pagination limits
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	// Fetch from repository
	attrs, total, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list attributes")
		return nil, 0, fmt.Errorf("failed to list attributes: %w", err)
	}

	return attrs, total, nil
}

// LinkAttributeToCategory links an attribute to a category with optional overrides
func (s *AttributeServiceImpl) LinkAttributeToCategory(
	ctx context.Context,
	categoryID int32,
	attributeID int32,
	settings *domain.CategoryAttributeSettings,
) error {
	// 1. Validate attribute exists
	attr, err := s.GetAttributeByID(ctx, attributeID)
	if err != nil {
		return fmt.Errorf("attribute not found: %w", err)
	}

	// 2. Validate settings
	if settings == nil {
		settings = &domain.CategoryAttributeSettings{
			IsEnabled: true,
			SortOrder: 0,
		}
	}

	// 3. Link in repository
	_, err = s.repo.LinkToCategory(ctx, categoryID, attributeID, settings)
	if err != nil {
		s.logger.Error().Err(err).
			Int32("category_id", categoryID).
			Int32("attribute_id", attributeID).
			Msg("failed to link attribute to category")
		return fmt.Errorf("failed to link attribute to category: %w", err)
	}

	// 4. Invalidate category cache
	if err := s.cache.InvalidateCategory(ctx, categoryID); err != nil {
		s.logger.Warn().Err(err).Int32("category_id", categoryID).Msg("failed to invalidate category cache")
	}

	s.logger.Info().
		Int32("category_id", categoryID).
		Int32("attribute_id", attributeID).
		Str("attribute_code", attr.Code).
		Msg("attribute linked to category")

	return nil
}

// UpdateCategoryAttribute updates category-specific attribute settings
func (s *AttributeServiceImpl) UpdateCategoryAttribute(
	ctx context.Context,
	catAttrID int32,
	settings *domain.CategoryAttributeSettings,
) error {
	// 1. Validate settings
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}

	// 2. Update in repository
	catAttr, err := s.repo.UpdateCategoryAttribute(ctx, catAttrID, settings)
	if err != nil {
		s.logger.Error().Err(err).Int32("cat_attr_id", catAttrID).Msg("failed to update category attribute")
		return fmt.Errorf("failed to update category attribute: %w", err)
	}

	// 3. Invalidate category cache
	if err := s.cache.InvalidateCategory(ctx, catAttr.CategoryID); err != nil {
		s.logger.Warn().Err(err).Int32("category_id", catAttr.CategoryID).Msg("failed to invalidate category cache")
	}

	s.logger.Info().Int32("cat_attr_id", catAttrID).Int32("category_id", catAttr.CategoryID).Msg("category attribute updated")
	return nil
}

// UnlinkAttributeFromCategory removes attribute-category association
func (s *AttributeServiceImpl) UnlinkAttributeFromCategory(ctx context.Context, categoryID int32, attributeID int32) error {
	// 1. Unlink in repository
	if err := s.repo.UnlinkFromCategory(ctx, categoryID, attributeID); err != nil {
		s.logger.Error().Err(err).
			Int32("category_id", categoryID).
			Int32("attribute_id", attributeID).
			Msg("failed to unlink attribute from category")
		return fmt.Errorf("failed to unlink attribute from category: %w", err)
	}

	// 2. Invalidate category cache
	if err := s.cache.InvalidateCategory(ctx, categoryID); err != nil {
		s.logger.Warn().Err(err).Int32("category_id", categoryID).Msg("failed to invalidate category cache")
	}

	s.logger.Info().
		Int32("category_id", categoryID).
		Int32("attribute_id", attributeID).
		Msg("attribute unlinked from category")

	return nil
}

// GetCategoryAttributes retrieves all attributes for a category with effective values (cached)
func (s *AttributeServiceImpl) GetCategoryAttributes(
	ctx context.Context,
	categoryID int32,
	filter *domain.GetCategoryAttributesFilter,
) ([]*domain.CategoryAttribute, error) {
	// 1. Try cache first (if no filter specified - full cache)
	if filter == nil || s.isEmptyFilter(filter) {
		cached, err := s.cache.GetCategoryAttributes(ctx, categoryID)
		if err == nil && cached != nil {
			return cached, nil
		}
	}

	// 2. Cache miss or filtered query - fetch from repository
	attrs, err := s.repo.GetCategoryAttributes(ctx, categoryID, filter)
	if err != nil {
		s.logger.Error().Err(err).Int32("category_id", categoryID).Msg("failed to get category attributes")
		return nil, fmt.Errorf("failed to get category attributes: %w", err)
	}

	// 3. Cache the result (only if no filter)
	if filter == nil || s.isEmptyFilter(filter) {
		if err := s.cache.SetCategoryAttributes(ctx, categoryID, attrs); err != nil {
			s.logger.Warn().Err(err).Int32("category_id", categoryID).Msg("failed to cache category attributes")
		}
	}

	return attrs, nil
}

// GetListingAttributes retrieves all attribute values for a listing (cached)
func (s *AttributeServiceImpl) GetListingAttributes(ctx context.Context, listingID int32) ([]*domain.ListingAttributeValue, error) {
	// 1. Try cache first
	cached, err := s.cache.GetListingAttributes(ctx, listingID)
	if err == nil && cached != nil {
		return cached, nil
	}

	// 2. Cache miss - fetch from repository
	values, err := s.repo.GetListingValues(ctx, listingID)
	if err != nil {
		s.logger.Error().Err(err).Int32("listing_id", listingID).Msg("failed to get listing attributes")
		return nil, fmt.Errorf("failed to get listing attributes: %w", err)
	}

	// 3. Cache the result
	if err := s.cache.SetListingAttributes(ctx, listingID, values); err != nil {
		s.logger.Warn().Err(err).Int32("listing_id", listingID).Msg("failed to cache listing attributes")
	}

	return values, nil
}

// SetListingAttributes validates and sets attribute values for a listing
func (s *AttributeServiceImpl) SetListingAttributes(
	ctx context.Context,
	listingID int32,
	values []domain.SetListingAttributeValue,
) error {
	// Note: Category validation should be done by caller (they need to know listing's category)
	// This method focuses on type validation

	// 1. Validate all values
	for i, val := range values {
		// Get attribute metadata
		attr, err := s.GetAttributeByID(ctx, val.AttributeID)
		if err != nil {
			return fmt.Errorf("invalid attribute_id at index %d: %w", i, err)
		}

		// Validate value based on type
		if err := s.validateAttributeValue(&val, attr); err != nil {
			return fmt.Errorf("validation failed for attribute '%s' at index %d: %w", attr.Code, i, err)
		}
	}

	// 2. Set values in repository (replaces existing)
	if err := s.repo.SetListingValues(ctx, listingID, values); err != nil {
		s.logger.Error().Err(err).Int32("listing_id", listingID).Msg("failed to set listing attributes")
		return fmt.Errorf("failed to set listing attributes: %w", err)
	}

	// 3. Invalidate cache
	if err := s.cache.InvalidateListing(ctx, listingID); err != nil {
		s.logger.Warn().Err(err).Int32("listing_id", listingID).Msg("failed to invalidate listing cache")
	}

	s.logger.Info().Int32("listing_id", listingID).Int("count", len(values)).Msg("listing attributes set successfully")
	return nil
}

// ValidateAttributeValues validates attribute values according to type and rules
func (s *AttributeServiceImpl) ValidateAttributeValues(
	ctx context.Context,
	categoryID int32,
	values []domain.SetListingAttributeValue,
) error {
	// 1. Get category attributes to validate against
	catAttrs, err := s.GetCategoryAttributes(ctx, categoryID, nil)
	if err != nil {
		return fmt.Errorf("failed to get category attributes: %w", err)
	}

	// Build map of attribute_id -> category_attribute for quick lookup
	catAttrMap := make(map[int32]*domain.CategoryAttribute)
	for _, ca := range catAttrs {
		catAttrMap[ca.AttributeID] = ca
	}

	// 2. Validate each value
	for i, val := range values {
		// Check if attribute is enabled for this category
		catAttr, exists := catAttrMap[val.AttributeID]
		if !exists {
			return fmt.Errorf("attribute_id %d at index %d is not linked to category %d", val.AttributeID, i, categoryID)
		}

		if !catAttr.IsEnabled {
			return fmt.Errorf("attribute '%s' at index %d is disabled for this category", catAttr.Attribute.Code, i)
		}

		// Validate using effective values (category overrides or attribute defaults)
		if err := s.validateCategoryAttributeValue(&val, catAttr); err != nil {
			return fmt.Errorf("validation failed for attribute '%s' at index %d: %w", catAttr.Attribute.Code, i, err)
		}
	}

	// 3. Check for required attributes
	for _, catAttr := range catAttrs {
		if catAttr.GetEffectiveIsRequired() {
			found := false
			for _, val := range values {
				if val.AttributeID == catAttr.AttributeID {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("required attribute '%s' is missing", catAttr.Attribute.Code)
			}
		}
	}

	return nil
}

// GetCategoryVariantAttributes retrieves variant attributes for a category (cached)
func (s *AttributeServiceImpl) GetCategoryVariantAttributes(ctx context.Context, categoryID int32) ([]*domain.VariantAttribute, error) {
	// 1. Try cache first
	cached, err := s.cache.GetCategoryVariantAttributes(ctx, categoryID)
	if err == nil && cached != nil {
		return cached, nil
	}

	// 2. Cache miss - fetch from repository
	attrs, err := s.repo.GetCategoryVariantAttributes(ctx, categoryID)
	if err != nil {
		s.logger.Error().Err(err).Int32("category_id", categoryID).Msg("failed to get variant attributes")
		return nil, fmt.Errorf("failed to get variant attributes: %w", err)
	}

	// 3. Cache the result
	if err := s.cache.SetCategoryVariantAttributes(ctx, categoryID, attrs); err != nil {
		s.logger.Warn().Err(err).Int32("category_id", categoryID).Msg("failed to cache variant attributes")
	}

	return attrs, nil
}

// InvalidateAttributeCache invalidates cache for an attribute
func (s *AttributeServiceImpl) InvalidateAttributeCache(ctx context.Context, attributeID int32) error {
	// Get attribute to get code
	attr, err := s.repo.GetByID(ctx, attributeID)
	if err != nil {
		return fmt.Errorf("failed to get attribute: %w", err)
	}
	if attr == nil {
		return fmt.Errorf("attribute not found")
	}

	return s.cache.InvalidateAttribute(ctx, attributeID, attr.Code)
}

// InvalidateCategoryCache invalidates all attribute caches for a category
func (s *AttributeServiceImpl) InvalidateCategoryCache(ctx context.Context, categoryID int32) error {
	return s.cache.InvalidateCategory(ctx, categoryID)
}

// InvalidateListingCache invalidates attribute cache for a listing
func (s *AttributeServiceImpl) InvalidateListingCache(ctx context.Context, listingID int32) error {
	return s.cache.InvalidateListing(ctx, listingID)
}

// Private validation helper methods

// validateAttributeInput validates business rules for attribute creation
func (s *AttributeServiceImpl) validateAttributeInput(input *domain.CreateAttributeInput) error {
	// Validate attribute type
	validTypes := map[domain.AttributeType]bool{
		domain.AttributeTypeText:        true,
		domain.AttributeTypeTextarea:    true,
		domain.AttributeTypeNumber:      true,
		domain.AttributeTypeBoolean:     true,
		domain.AttributeTypeSelect:      true,
		domain.AttributeTypeMultiselect: true,
		domain.AttributeTypeDate:        true,
		domain.AttributeTypeColor:       true,
		domain.AttributeTypeSize:        true,
	}

	if !validTypes[input.AttributeType] {
		return fmt.Errorf("invalid attribute type: %s", input.AttributeType)
	}

	// Select/Multiselect must have options
	if (input.AttributeType == domain.AttributeTypeSelect || input.AttributeType == domain.AttributeTypeMultiselect) &&
		len(input.Options) == 0 {
		return fmt.Errorf("select/multiselect attributes must have at least one option")
	}

	// Validate i18n fields have at least one language
	if len(input.Name) == 0 {
		return fmt.Errorf("name must have at least one language")
	}
	if len(input.DisplayName) == 0 {
		return fmt.Errorf("display_name must have at least one language")
	}

	return nil
}

// validateAttributeUpdate validates business rules for attribute updates
func (s *AttributeServiceImpl) validateAttributeUpdate(existing *domain.Attribute, input *domain.UpdateAttributeInput) error {
	// If type is being changed, validate it
	if input.AttributeType != nil {
		// Cannot change type if attribute is already used in categories or listings
		// This is a simplified check - in production you'd query usage
		if *input.AttributeType != existing.AttributeType {
			return fmt.Errorf("cannot change attribute type after creation")
		}
	}

	// If changing to select/multiselect, ensure options are provided
	targetType := existing.AttributeType
	if input.AttributeType != nil {
		targetType = *input.AttributeType
	}

	if targetType == domain.AttributeTypeSelect || targetType == domain.AttributeTypeMultiselect {
		// Check if options will be empty after update
		hasOptions := len(existing.Options) > 0
		if input.Options != nil {
			hasOptions = len(*input.Options) > 0
		}
		if !hasOptions {
			return fmt.Errorf("select/multiselect attributes must have at least one option")
		}
	}

	return nil
}

// validateAttributeValue validates a single attribute value
func (s *AttributeServiceImpl) validateAttributeValue(val *domain.SetListingAttributeValue, attr *domain.Attribute) error {
	// Extract the actual value based on attribute type
	var value interface{}

	switch attr.AttributeType {
	case domain.AttributeTypeText, domain.AttributeTypeTextarea, domain.AttributeTypeSize:
		if val.ValueText == nil {
			return fmt.Errorf("value_text is required for type %s", attr.AttributeType)
		}
		value = *val.ValueText

	case domain.AttributeTypeNumber:
		if val.ValueNumber == nil {
			return fmt.Errorf("value_number is required for type number")
		}
		value = *val.ValueNumber

	case domain.AttributeTypeBoolean:
		if val.ValueBoolean == nil {
			return fmt.Errorf("value_boolean is required for type boolean")
		}
		value = *val.ValueBoolean

	case domain.AttributeTypeSelect:
		if val.ValueText == nil {
			return fmt.Errorf("value_text is required for type select")
		}
		value = *val.ValueText

	case domain.AttributeTypeMultiselect:
		if val.ValueJSON == nil {
			return fmt.Errorf("value_json is required for type multiselect")
		}
		// Convert ValueJSON to []string for validation
		if arr, ok := val.ValueJSON["values"].([]interface{}); ok {
			value = arr
		} else {
			return fmt.Errorf("multiselect value_json must contain 'values' array")
		}

	case domain.AttributeTypeDate:
		if val.ValueDate == nil {
			return fmt.Errorf("value_date is required for type date")
		}
		value = *val.ValueDate

	case domain.AttributeTypeColor:
		if val.ValueText == nil {
			return fmt.Errorf("value_text is required for type color")
		}
		value = *val.ValueText

	default:
		return fmt.Errorf("unsupported attribute type: %s", attr.AttributeType)
	}

	// Validate using validator
	return s.validator.ValidateValue(attr.AttributeType, value, attr.ValidationRules, attr.Options)
}

// validateCategoryAttributeValue validates a value using category-specific rules
func (s *AttributeServiceImpl) validateCategoryAttributeValue(val *domain.SetListingAttributeValue, catAttr *domain.CategoryAttribute) error {
	if catAttr.Attribute == nil {
		return fmt.Errorf("category attribute is missing attribute metadata")
	}

	// Get effective validation rules and options
	rules := catAttr.GetEffectiveValidationRules()
	options := catAttr.GetEffectiveOptions()

	// Extract value based on type
	var value interface{}
	attr := catAttr.Attribute

	switch attr.AttributeType {
	case domain.AttributeTypeText, domain.AttributeTypeTextarea, domain.AttributeTypeSize:
		if val.ValueText == nil {
			return fmt.Errorf("value_text is required")
		}
		value = *val.ValueText

	case domain.AttributeTypeNumber:
		if val.ValueNumber == nil {
			return fmt.Errorf("value_number is required")
		}
		value = *val.ValueNumber

	case domain.AttributeTypeBoolean:
		if val.ValueBoolean == nil {
			return fmt.Errorf("value_boolean is required")
		}
		value = *val.ValueBoolean

	case domain.AttributeTypeSelect:
		if val.ValueText == nil {
			return fmt.Errorf("value_text is required")
		}
		value = *val.ValueText

	case domain.AttributeTypeMultiselect:
		if val.ValueJSON == nil {
			return fmt.Errorf("value_json is required")
		}
		// Convert to appropriate format for validation
		if arr, ok := val.ValueJSON["values"].([]interface{}); ok {
			value = arr
		} else if str, ok := val.ValueJSON["values"].(string); ok {
			// Handle case where it's stored as JSON string
			var arr []string
			if err := json.Unmarshal([]byte(str), &arr); err != nil {
				return fmt.Errorf("invalid multiselect format: %w", err)
			}
			value = arr
		} else {
			return fmt.Errorf("multiselect value_json must contain 'values' array")
		}

	case domain.AttributeTypeDate:
		if val.ValueDate == nil {
			return fmt.Errorf("value_date is required")
		}
		value = *val.ValueDate

	case domain.AttributeTypeColor:
		if val.ValueText == nil {
			return fmt.Errorf("value_text is required")
		}
		value = *val.ValueText
	}

	// Validate using effective rules
	return s.validator.ValidateValue(attr.AttributeType, value, rules, options)
}

// isEmptyFilter checks if a filter has any criteria set
func (s *AttributeServiceImpl) isEmptyFilter(filter *domain.GetCategoryAttributesFilter) bool {
	return filter.IsEnabled == nil &&
		filter.IsRequired == nil &&
		filter.IsSearchable == nil &&
		filter.IsFilterable == nil
}
