package handler

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"backend/internal/domain/models"
	"backend/internal/proj/storefronts/parsers"
	"backend/internal/proj/storefronts/service"

	"github.com/gofiber/fiber/v2"
)

const (
	fileTypeXML = "xml"
	fileTypeCSV = "csv"
)

// AnalyzeCategoriesRequest represents the request for category analysis
type AnalyzeCategoriesRequest struct {
	FileType      string                        `json:"file_type"`                // xml, csv
	Products      []models.ImportProductRequest `json:"products"`                 // Sample products with categories
	CategoryPaths []string                      `json:"category_paths,omitempty"` // Unique category paths from file
}

// CategoryMappingDTO represents a single category mapping for frontend
type CategoryMappingDTO struct {
	ExternalCategory              string `json:"external_category"`
	SuggestedInternalCategoryID   *int   `json:"suggested_internal_category_id"`
	SuggestedInternalCategoryName string `json:"suggested_internal_category_name,omitempty"`
	Confidence                    string `json:"confidence"` // "high", "medium", "low"
	Reasoning                     string `json:"reasoning,omitempty"`
	IsApproved                    bool   `json:"is_approved"`
	RequiresNewCategory           bool   `json:"requires_new_category"`
}

// QualitySummaryDTO represents mapping quality statistics
type QualitySummaryDTO struct {
	HighConfidence      int `json:"high_confidence"`
	MediumConfidence    int `json:"medium_confidence"`
	LowConfidence       int `json:"low_confidence"`
	RequiresNewCategory int `json:"requires_new_category"`
}

// AnalyzeCategoriesResponse represents the response with AI category mapping suggestions
type AnalyzeCategoriesResponse struct {
	TotalCategories    int                  `json:"total_categories"`
	Mappings           []CategoryMappingDTO `json:"mappings"`
	QualitySummary     QualitySummaryDTO    `json:"quality_summary"`
	UnmappedCategories []string             `json:"unmapped_categories"`
	ProcessingTimeMs   int64                `json:"processing_time_ms"`
}

// AnalyzeCategories analyzes categories in import file and provides AI mapping suggestions
// @Summary Analyze categories in import file
// @Description Analyze categories and get AI mapping suggestions before import
// @Tags storefronts,import
// @Accept multipart/form-data
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param file formData file true "Import file"
// @Param file_type formData string true "File type" Enums(xml,csv,zip)
// @Success 200 {object} AnalyzeCategoriesResponse "Category analysis result"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/import/analyze-categories [post]
func (h *ImportHandler) AnalyzeCategories(c *fiber.Ctx) error {
	// Get storefront ID from path
	storefrontID, err := c.ParamsInt("storefront_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid storefront_id",
			Message: err.Error(),
		})
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "File upload failed",
			Message: err.Error(),
		})
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Failed to open uploaded file",
			Message: err.Error(),
		})
	}
	defer func() {
		if err := src.Close(); err != nil {
			fmt.Printf("Failed to close file: %v", err)
		}
	}()

	// Read file data
	fileData, err := io.ReadAll(src)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Failed to read uploaded file",
			Message: err.Error(),
		})
	}

	// Get file type from form
	fileType := c.FormValue("file_type")
	if fileType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "file_type is required",
		})
	}

	// Parse file to extract products
	var products []models.ImportProductRequest
	var parseErrors []models.ImportValidationError

	switch fileType {
	case fileTypeXML:
		xmlParser := parsers.NewXMLParser(storefrontID)
		// Try Digital Vision format first
		products, parseErrors, err = xmlParser.ParseDigitalVisionXML(fileData)
		if err != nil {
			// Try generic XML format
			products, parseErrors, err = xmlParser.ParseGenericXML(fileData)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
					Error:   "Failed to parse XML file",
					Message: err.Error(),
				})
			}
		}
	case fileTypeCSV:
		csvParser := parsers.NewCSVParser(storefrontID)
		reader := bytes.NewReader(fileData)
		products, parseErrors, err = csvParser.ParseCSV(reader)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "Failed to parse CSV file",
				Message: err.Error(),
			})
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Unsupported file type. Supported types: xml, csv",
		})
	}

	// Log parse errors (but continue with analysis)
	if len(parseErrors) > 0 {
		fmt.Printf("Parse warnings: %d\n", len(parseErrors))
	}

	// Extract unique categories and sample products
	uniqueCategories := make(map[string]bool)
	categoryToProduct := make(map[string]models.ImportProductRequest)

	// Extract categories from products
	for _, product := range products {
		// Try to extract category from Attributes map
		if product.Attributes != nil {
			var categoryPath string

			// Try category_path first (generic format)
			if catPath, ok := product.Attributes["category_path"].(string); ok && catPath != "" {
				categoryPath = catPath
			} else {
				// Try Digital Vision format (kategorija1, kategorija2, kategorija3)
				cat1, _ := product.Attributes["kategorija1"].(string)
				cat2, _ := product.Attributes["kategorija2"].(string)
				cat3, _ := product.Attributes["kategorija3"].(string)

				// Build category path from available levels
				categories := []string{}
				if cat1 != "" {
					categories = append(categories, cat1)
				}
				if cat2 != "" {
					categories = append(categories, cat2)
				}
				if cat3 != "" {
					categories = append(categories, cat3)
				}

				if len(categories) > 0 {
					categoryPath = categories[len(categories)-1] // Use deepest category
				}
			}

			if categoryPath != "" {
				uniqueCategories[categoryPath] = true
				// Keep first product as sample for this category
				if _, exists := categoryToProduct[categoryPath]; !exists {
					categoryToProduct[categoryPath] = product
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

	// Convert map to array of DTOs
	mappings := make([]CategoryMappingDTO, 0, len(suggestions))
	unmappedCategories := make([]string, 0)

	for externalCat, suggestion := range suggestions {
		// Determine confidence level
		confidence := "low"
		if suggestion.ConfidenceScore >= 0.8 {
			confidence = "high"
		} else if suggestion.ConfidenceScore >= 0.5 {
			confidence = "medium"
		}

		// Build reasoning string from steps
		reasoning := ""
		if len(suggestion.ReasoningSteps) > 0 {
			reasoning = suggestion.ReasoningSteps[0] // Use first reasoning step
		}

		// Check if category was mapped
		var categoryID *int
		categoryName := ""
		requiresNew := false

		if suggestion.SuggestedCategoryID > 0 {
			categoryID = &suggestion.SuggestedCategoryID
			categoryName = suggestion.SuggestedCategoryPath
		} else {
			requiresNew = true
			unmappedCategories = append(unmappedCategories, externalCat)
		}

		mappings = append(mappings, CategoryMappingDTO{
			ExternalCategory:              externalCat,
			SuggestedInternalCategoryID:   categoryID,
			SuggestedInternalCategoryName: categoryName,
			Confidence:                    confidence,
			Reasoning:                     reasoning,
			IsApproved:                    confidence == "high", // Auto-approve high confidence
			RequiresNewCategory:           requiresNew,
		})
	}

	// Build quality summary
	qualitySummary := QualitySummaryDTO{
		HighConfidence:      len(quality.HighConfidence),
		MediumConfidence:    len(quality.MediumConfidence),
		LowConfidence:       len(quality.LowConfidence),
		RequiresNewCategory: len(unmappedCategories),
	}

	response := AnalyzeCategoriesResponse{
		TotalCategories:    len(categories),
		Mappings:           mappings,
		QualitySummary:     qualitySummary,
		UnmappedCategories: unmappedCategories,
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
	Name              string   `json:"name"`
	ValueType         string   `json:"value_type"`          // Type of value: string, number, boolean, enum
	SampleValues      []string `json:"sample_values"`       // Example values
	Frequency         int      `json:"frequency"`           // How many products have this attribute
	IsStandard        bool     `json:"is_standard"`         // Is this a standard marketplace attribute
	SuggestedMapping  string   `json:"suggested_mapping"`   // Suggested mapping to standard attribute
	IsVariantDefining bool     `json:"is_variant_defining"` // Could be used for variants (color, size, etc.)
}

// detectValueType determines the type of attribute based on sample values
func detectValueType(values []string) string {
	if len(values) == 0 {
		return "string"
	}

	// Check if all values are numbers
	allNumbers := true
	for _, v := range values {
		if _, err := strconv.ParseFloat(v, 64); err != nil {
			allNumbers = false
			break
		}
	}
	if allNumbers {
		return "number"
	}

	// Check if all values are booleans
	allBooleans := true
	for _, v := range values {
		lower := strings.ToLower(strings.TrimSpace(v))
		if lower != "true" && lower != "false" && lower != "да" && lower != "нет" && lower != "yes" && lower != "no" {
			allBooleans = false
			break
		}
	}
	if allBooleans {
		return "boolean"
	}

	// Check if limited unique values (enum)
	uniqueValues := make(map[string]bool)
	for _, v := range values {
		uniqueValues[v] = true
	}
	if len(uniqueValues) <= 10 && len(values) > len(uniqueValues)*2 {
		// If there are <=10 unique values and they repeat (frequency > 2x unique count)
		return "enum"
	}

	return "string"
}

// isVariantDefiningAttribute checks if attribute name suggests it could define variants
func isVariantDefiningAttribute(name string) bool {
	nameLower := strings.ToLower(name)

	variantKeywords := []string{
		"color", "colour", "boja", "цвет", //nolint:misspell // "colour" - британский вариант написания
		"size", "velicina", "veličina", "размер",
		"material", "materijal", "материал",
		"style", "stil", "стиль",
		"model", "модель",
		"capacity", "kapacitet", "capacité", "емкость",
		"version", "verzija", "версия",
		"type", "tip", "тип",
	}

	for _, keyword := range variantKeywords {
		if strings.Contains(nameLower, keyword) {
			return true
		}
	}

	return false
}

// AnalyzeAttributesResponse represents the response with detected attributes
type AnalyzeAttributesResponse struct {
	Attributes                []DetectedAttribute `json:"attributes"`
	TotalAttributes           int                 `json:"total_attributes"`
	VariantDefiningAttributes []string            `json:"variant_defining_attributes"`
	TotalProducts             int                 `json:"total_products"`
	ProcessingTimeMs          int64               `json:"processing_time_ms"`
}

// AnalyzeAttributes analyzes attributes in import file
// @Summary Analyze attributes in import file
// @Description Detect and analyze product attributes before import
// @Tags storefronts,import
// @Accept multipart/form-data
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param file formData file true "Import file"
// @Param file_type formData string true "File type" Enums(xml,csv,zip)
// @Success 200 {object} AnalyzeAttributesResponse "Attribute analysis result"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/import/analyze-attributes [post]
func (h *ImportHandler) AnalyzeAttributes(c *fiber.Ctx) error {
	// Get storefront ID from path
	storefrontID, err := c.ParamsInt("storefront_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid storefront_id",
			Message: err.Error(),
		})
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "File upload failed",
			Message: err.Error(),
		})
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Failed to open uploaded file",
			Message: err.Error(),
		})
	}
	defer func() {
		if err := src.Close(); err != nil {
			fmt.Printf("Failed to close file: %v", err)
		}
	}()

	// Read file data
	fileData, err := io.ReadAll(src)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Failed to read uploaded file",
			Message: err.Error(),
		})
	}

	// Get file type from form
	fileType := c.FormValue("file_type")
	if fileType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "file_type is required",
		})
	}

	// Parse file to extract products
	var products []models.ImportProductRequest
	var parseErrors []models.ImportValidationError

	switch fileType {
	case "xml":
		xmlParser := parsers.NewXMLParser(storefrontID)
		products, parseErrors, err = xmlParser.ParseDigitalVisionXML(fileData)
		if err != nil {
			products, parseErrors, err = xmlParser.ParseGenericXML(fileData)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
					Error:   "Failed to parse XML file",
					Message: err.Error(),
				})
			}
		}
	case "csv":
		csvParser := parsers.NewCSVParser(storefrontID)
		reader := bytes.NewReader(fileData)
		products, parseErrors, err = csvParser.ParseCSV(reader)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "Failed to parse CSV file",
				Message: err.Error(),
			})
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Unsupported file type. Supported types: xml, csv",
		})
	}

	// Log parse errors
	if len(parseErrors) > 0 {
		fmt.Printf("Parse warnings: %d\n", len(parseErrors))
	}

	detectedAttrs := make([]DetectedAttribute, 0)

	// Analyze products to detect common attributes from Attributes map
	attributeFrequency := make(map[string]int)
	attributeExamples := make(map[string][]string)
	attributeFirstValue := make(map[string]interface{}) // Для маппинга

	fmt.Printf("[DEBUG] AnalyzeAttributes: Total products parsed: %d\n", len(products))

	for _, product := range products {
		// Analyze custom attributes from product.Attributes
		if product.Attributes != nil {
			fmt.Printf("[DEBUG] Product %s has %d attributes\n", product.SKU, len(product.Attributes))
			for attrName, attrValue := range product.Attributes {
				fmt.Printf("[DEBUG]   - %s: %v (type: %T)\n", attrName, attrValue, attrValue)

				// Skip standard fields and category fields
				if attrName == "category_path" ||
					attrName == "kategorija1" ||
					attrName == "kategorija2" ||
					attrName == "kategorija3" {
					fmt.Printf("[DEBUG]     SKIPPED (category field)\n")
					continue
				}

				attributeFrequency[attrName]++

				// Сохраняем первое значение для маппинга
				if _, exists := attributeFirstValue[attrName]; !exists {
					attributeFirstValue[attrName] = attrValue
				}

				// Add example value (up to 3 examples)
				if len(attributeExamples[attrName]) < 3 {
					if strVal, ok := attrValue.(string); ok {
						attributeExamples[attrName] = append(attributeExamples[attrName], strVal)
					}
				}
			}
		}
	}

	fmt.Printf("[DEBUG] Detected %d unique attributes after filtering\n", len(attributeFrequency))

	// Используем AttributeMapper для маппинга атрибутов
	fmt.Printf("[DEBUG] Starting attribute mapping for %d attributes\n", len(attributeFrequency))
	for attrName, frequency := range attributeFrequency {
		fmt.Printf("[DEBUG] Mapping attribute: %s (frequency: %d)\n", attrName, frequency)
		// Мапим через AttributeMapper
		mappedAttr, err := h.attributeMapper.MapExternalAttribute(
			c.Context(),
			attrName,
			attributeFirstValue[attrName],
			nil, // categoryID optional
		)

		suggestedMap := "custom_attribute"
		isStandard := false

		if err == nil && mappedAttr != nil {
			// Если нашли маппинг в unified_attributes
			if !mappedAttr.IsNewAttribute {
				suggestedMap = mappedAttr.Code
				isStandard = true
			} else {
				// Если это новый атрибут - используем предложенный код
				suggestedMap = mappedAttr.SuggestedCode
			}
		}

		sampleValues := attributeExamples[attrName]
		valueType := detectValueType(sampleValues)
		isVariantDefining := isVariantDefiningAttribute(attrName)

		detectedAttrs = append(detectedAttrs, DetectedAttribute{
			Name:              attrName,
			ValueType:         valueType,
			SampleValues:      sampleValues,
			Frequency:         frequency,
			IsStandard:        isStandard,
			SuggestedMapping:  suggestedMap,
			IsVariantDefining: isVariantDefining,
		})
	}

	// Collect variant defining attribute names
	variantDefiningAttrs := make([]string, 0)
	for _, attr := range detectedAttrs {
		if attr.IsVariantDefining {
			variantDefiningAttrs = append(variantDefiningAttrs, attr.Name)
		}
	}

	response := AnalyzeAttributesResponse{
		Attributes:                detectedAttrs,
		TotalAttributes:           len(detectedAttrs),
		VariantDefiningAttributes: variantDefiningAttrs,
		TotalProducts:             len(products),
		ProcessingTimeMs:          0,
	}

	fmt.Printf("[DEBUG] Final response: %d detected attributes, %d total products\n", len(detectedAttrs), len(products))
	for i, attr := range detectedAttrs {
		fmt.Printf("[DEBUG]   Attribute %d: %s (type: %s, samples: %v, frequency: %d, standard: %v, variant: %v)\n",
			i+1, attr.Name, attr.ValueType, attr.SampleValues, attr.Frequency, attr.IsStandard, attr.IsVariantDefining)
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
	BaseName        string   `json:"base_name"`
	VariantCount    int      `json:"variant_count"`
	VariantAttr     string   `json:"variant_attr"`   // e.g., "color", "size", "model"
	VariantValues   []string `json:"variant_values"` // e.g., ["Black", "White", "Red"]
	ProductSKUs     []string `json:"product_skus"`
	ConfidenceScore float64  `json:"confidence_score"`  // 0.0-1.0
	ShouldGroupAuto bool     `json:"should_group_auto"` // true if confidence > 0.85
}

// DetectVariantsResponse represents the response with detected variant groups
type DetectVariantsResponse struct {
	VariantGroups      []VariantGroup `json:"variant_groups"`
	TotalProducts      int            `json:"total_products"`
	GroupedProducts    int            `json:"grouped_products"`
	StandaloneProducts int            `json:"standalone_products"`
	ProcessingTimeMs   int64          `json:"processing_time_ms"`
}

// DetectVariants detects potential product variants in import file
// @Summary Detect product variants
// @Description Detect and group product variants before import
// @Tags storefronts,import
// @Accept multipart/form-data
// @Produce json
// @Param storefront_id path int true "Storefront ID"
// @Param file formData file true "Import file"
// @Param file_type formData string true "File type" Enums(xml,csv,zip)
// @Success 200 {object} DetectVariantsResponse "Variant detection result"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{storefront_id}/import/detect-variants [post]
func (h *ImportHandler) DetectVariants(c *fiber.Ctx) error {
	// Get storefront ID from path
	storefrontID, err := c.ParamsInt("storefront_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid storefront_id",
			Message: err.Error(),
		})
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "File upload failed",
			Message: err.Error(),
		})
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Failed to open uploaded file",
			Message: err.Error(),
		})
	}
	defer func() {
		if err := src.Close(); err != nil {
			fmt.Printf("Failed to close file: %v", err)
		}
	}()

	// Read file data
	fileData, err := io.ReadAll(src)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Failed to read uploaded file",
			Message: err.Error(),
		})
	}

	// Get file type from form
	fileType := c.FormValue("file_type")
	if fileType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "file_type is required",
		})
	}

	// Parse file to extract products
	var products []models.ImportProductRequest
	var parseErrors []models.ImportValidationError

	switch fileType {
	case "xml":
		xmlParser := parsers.NewXMLParser(storefrontID)
		products, parseErrors, err = xmlParser.ParseDigitalVisionXML(fileData)
		if err != nil {
			products, parseErrors, err = xmlParser.ParseGenericXML(fileData)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
					Error:   "Failed to parse XML file",
					Message: err.Error(),
				})
			}
		}
	case "csv":
		csvParser := parsers.NewCSVParser(storefrontID)
		reader := bytes.NewReader(fileData)
		products, parseErrors, err = csvParser.ParseCSV(reader)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "Failed to parse CSV file",
				Message: err.Error(),
			})
		}
	default:
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Unsupported file type. Supported types: xml, csv",
		})
	}

	// Log parse errors
	if len(parseErrors) > 0 {
		fmt.Printf("Parse warnings: %d\n", len(parseErrors))
	}

	// TODO: Implement variant detection logic
	// This will be implemented in Phase 2 (Variant Detector)

	// For now, return empty result
	response := DetectVariantsResponse{
		VariantGroups:      []VariantGroup{},
		TotalProducts:      len(products),
		GroupedProducts:    0,
		StandaloneProducts: len(products),
		ProcessingTimeMs:   0,
	}

	return c.JSON(response)
}

// AnalyzeClientCategoriesRequest represents request for analyzing unique client categories
type AnalyzeClientCategoriesRequest struct {
	FileType   string                       `json:"file_type"`
	Categories []service.ClientCategoryInfo `json:"categories"` // Categories with product counts
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
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
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
