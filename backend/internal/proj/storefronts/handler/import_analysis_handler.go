package handler

import (
	"backend/internal/domain/models"
	"backend/internal/proj/storefronts/service"

	"github.com/gofiber/fiber/v2"
)

// AnalyzeCategoriesRequest represents the request for category analysis
type AnalyzeCategoriesRequest struct {
	FileType      string                          `json:"file_type"` // xml, csv
	Products      []models.ImportProductRequest   `json:"products"`  // Sample products with categories
	CategoryPaths []string                        `json:"category_paths,omitempty"` // Unique category paths from file
}

// AnalyzeCategoriesResponse represents the response with AI category mapping suggestions
type AnalyzeCategoriesResponse struct {
	TotalCategories     int                                          `json:"total_categories"`
	MappingSuggestions  map[string]*service.CategoryMappingSuggestion `json:"mapping_suggestions"` // external_cat -> suggestion
	MappingQuality      *service.MappingQuality                      `json:"mapping_quality"`
	ProcessingTimeMs    int64                                        `json:"processing_time_ms"`
}

// AnalyzeCategories analyzes categories in import file and provides AI mapping suggestions
// @Summary Analyze categories in import file
// @Description Analyze categories and get AI mapping suggestions before import
// @Tags storefronts,import
// @Accept json
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param request body AnalyzeCategoriesRequest true "Analysis request"
// @Success 200 {object} AnalyzeCategoriesResponse "Category analysis result"
// @Failure 400 {object} backend_internal_domain_models.ErrorResponse "Bad request"
// @Failure 401 {object} backend_internal_domain_models.ErrorResponse "Unauthorized"
// @Failure 403 {object} backend_internal_domain_models.ErrorResponse "Forbidden"
// @Failure 500 {object} backend_internal_domain_models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/import/analyze-categories [post]
func (h *ImportHandler) AnalyzeCategories(c *fiber.Ctx) error {
	var err error

	var req AnalyzeCategoriesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// Extract unique categories and sample products
	uniqueCategories := make(map[string]bool)
	categoryToProduct := make(map[string]models.ImportProductRequest)

	// Use category paths provided directly (for now, we just map these)
	if len(req.CategoryPaths) > 0 {
		for _, catPath := range req.CategoryPaths {
			uniqueCategories[catPath] = true
		}
	}

	// If products provided with category info in Attributes, extract them
	for _, product := range req.Products {
		// Try to extract category from Attributes map
		if product.Attributes != nil {
			if catPath, ok := product.Attributes["category_path"].(string); ok && catPath != "" {
				uniqueCategories[catPath] = true
				// Keep first product as sample for this category
				if _, exists := categoryToProduct[catPath]; !exists {
					categoryToProduct[catPath] = product
				}
			}
		}
	}

	// Convert to slice
	categories := make([]string, 0, len(uniqueCategories))
	for cat := range uniqueCategories {
		categories = append(categories, cat)
	}

	// Batch map categories
	suggestions, err := h.aiCategoryMapper.BatchMapCategories(c.Context(), categories, categoryToProduct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to analyze categories",
			Message: err.Error(),
		})
	}

	// Analyze mapping quality
	quality := h.aiCategoryMapper.AnalyzeMappingQuality(suggestions)

	response := AnalyzeCategoriesResponse{
		TotalCategories:    len(categories),
		MappingSuggestions: suggestions,
		MappingQuality:     quality,
		ProcessingTimeMs:   0, // TODO: track time
	}

	return c.JSON(response)
}

// AnalyzeAttributesRequest represents the request for attribute analysis
type AnalyzeAttributesRequest struct {
	FileType string                        `json:"file_type"` // xml, csv
	Products []models.ImportProductRequest `json:"products"`  // Sample products
}

// DetectedAttribute represents a detected attribute from import file
type DetectedAttribute struct {
	Name         string   `json:"name"`
	Examples     []string `json:"examples"`      // Example values
	Frequency    int      `json:"frequency"`     // How many products have this attribute
	IsStandard   bool     `json:"is_standard"`   // Is this a standard marketplace attribute
	SuggestedMap string   `json:"suggested_map"` // Suggested mapping to standard attribute
}

// AnalyzeAttributesResponse represents the response with detected attributes
type AnalyzeAttributesResponse struct {
	DetectedAttributes []DetectedAttribute `json:"detected_attributes"`
	TotalProducts      int                 `json:"total_products"`
	ProcessingTimeMs   int64               `json:"processing_time_ms"`
}

// AnalyzeAttributes analyzes attributes in import file
// @Summary Analyze attributes in import file
// @Description Detect and analyze product attributes before import
// @Tags storefronts,import
// @Accept json
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param request body AnalyzeAttributesRequest true "Analysis request"
// @Success 200 {object} AnalyzeAttributesResponse "Attribute analysis result"
// @Failure 400 {object} backend_internal_domain_models.ErrorResponse "Bad request"
// @Failure 401 {object} backend_internal_domain_models.ErrorResponse "Unauthorized"
// @Failure 403 {object} backend_internal_domain_models.ErrorResponse "Forbidden"
// @Failure 500 {object} backend_internal_domain_models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/import/analyze-attributes [post]
func (h *ImportHandler) AnalyzeAttributes(c *fiber.Ctx) error {
	var req AnalyzeAttributesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// TODO: Implement attribute detection logic
	// For now, return mock data showing what attributes we support

	detectedAttrs := make([]DetectedAttribute, 0)

	// Analyze products to detect common attributes from Attributes map
	attributeFrequency := make(map[string]int)
	attributeExamples := make(map[string][]string)

	for _, product := range req.Products {
		// Analyze custom attributes from product.Attributes
		if product.Attributes != nil {
			for attrName, attrValue := range product.Attributes {
				// Skip standard fields
				if attrName == "category_path" {
					continue
				}

				attributeFrequency[attrName]++

				// Add example value (up to 3 examples)
				if len(attributeExamples[attrName]) < 3 {
					if strVal, ok := attrValue.(string); ok {
						attributeExamples[attrName] = append(attributeExamples[attrName], strVal)
					}
				}
			}
		}
	}

	// Convert to DetectedAttribute format
	for attrName, frequency := range attributeFrequency {
		detectedAttrs = append(detectedAttrs, DetectedAttribute{
			Name:         attrName,
			Examples:     attributeExamples[attrName],
			Frequency:    frequency,
			IsStandard:   false,
			SuggestedMap: "custom_attribute",
		})
	}

	response := AnalyzeAttributesResponse{
		DetectedAttributes: detectedAttrs,
		TotalProducts:      len(req.Products),
		ProcessingTimeMs:   0,
	}

	return c.JSON(response)
}

// DetectVariantsRequest represents the request for variant detection
type DetectVariantsRequest struct {
	FileType string                        `json:"file_type"`
	Products []models.ImportProductRequest `json:"products"`
}

// VariantGroup represents a group of product variants
type VariantGroup struct {
	BaseName         string   `json:"base_name"`
	VariantCount     int      `json:"variant_count"`
	VariantAttr      string   `json:"variant_attr"`       // e.g., "color", "size", "model"
	VariantValues    []string `json:"variant_values"`     // e.g., ["Black", "White", "Red"]
	ProductSKUs      []string `json:"product_skus"`
	ConfidenceScore  float64  `json:"confidence_score"`   // 0.0-1.0
	ShouldGroupAuto  bool     `json:"should_group_auto"`  // true if confidence > 0.85
}

// DetectVariantsResponse represents the response with detected variant groups
type DetectVariantsResponse struct {
	VariantGroups       []VariantGroup `json:"variant_groups"`
	TotalProducts       int            `json:"total_products"`
	GroupedProducts     int            `json:"grouped_products"`
	StandaloneProducts  int            `json:"standalone_products"`
	ProcessingTimeMs    int64          `json:"processing_time_ms"`
}

// DetectVariants detects potential product variants in import file
// @Summary Detect product variants
// @Description Detect and group product variants before import
// @Tags storefronts,import
// @Accept json
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param request body DetectVariantsRequest true "Detection request"
// @Success 200 {object} DetectVariantsResponse "Variant detection result"
// @Failure 400 {object} backend_internal_domain_models.ErrorResponse "Bad request"
// @Failure 401 {object} backend_internal_domain_models.ErrorResponse "Unauthorized"
// @Failure 403 {object} backend_internal_domain_models.ErrorResponse "Forbidden"
// @Failure 500 {object} backend_internal_domain_models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/import/detect-variants [post]
func (h *ImportHandler) DetectVariants(c *fiber.Ctx) error {
	var req DetectVariantsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	// TODO: Implement variant detection logic
	// This will be implemented in Phase 2 (Variant Detector)

	// For now, return empty result
	response := DetectVariantsResponse{
		VariantGroups:      []VariantGroup{},
		TotalProducts:      len(req.Products),
		GroupedProducts:    0,
		StandaloneProducts: len(req.Products),
		ProcessingTimeMs:   0,
	}

	return c.JSON(response)
}

// AnalyzeClientCategoriesRequest represents request for analyzing unique client categories
type AnalyzeClientCategoriesRequest struct {
	FileType   string                               `json:"file_type"`
	Categories []service.ClientCategoryInfo         `json:"categories"` // Categories with product counts
}

// AnalyzeClientCategoriesResponse represents response with category insights
type AnalyzeClientCategoriesResponse struct {
	Analysis *service.CategoryAnalysisResult `json:"analysis"`
}

// AnalyzeClientCategories analyzes client's categories and suggests new categories
// @Summary Analyze client categories
// @Description Analyze unique client categories and suggest new marketplace categories
// @Tags storefronts,import
// @Accept json
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param request body AnalyzeClientCategoriesRequest true "Analysis request"
// @Success 200 {object} AnalyzeClientCategoriesResponse "Category analysis"
// @Failure 400 {object} backend_internal_domain_models.ErrorResponse "Bad request"
// @Failure 401 {object} backend_internal_domain_models.ErrorResponse "Unauthorized"
// @Failure 403 {object} backend_internal_domain_models.ErrorResponse "Forbidden"
// @Failure 500 {object} backend_internal_domain_models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/import/analyze-client-categories [post]
func (h *ImportHandler) AnalyzeClientCategories(c *fiber.Ctx) error {
	var req AnalyzeClientCategoriesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	analysis, err := h.aiCategoryAnalyzer.AnalyzeClientCategories(c.Context(), req.Categories)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to analyze categories",
			Message: err.Error(),
		})
	}

	return c.JSON(AnalyzeClientCategoriesResponse{
		Analysis: analysis,
	})
}
