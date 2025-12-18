// Package service contains business logic services for the listings microservice.
// This file implements SKU generation for product variants.
package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/vondi-global/listings/internal/domain"
)

// SKUGenerator generates unique SKU codes for product variants
type SKUGenerator struct {
	categoryPrefixes map[string]string
}

// NewSKUGenerator creates a new SKU generator instance
func NewSKUGenerator() *SKUGenerator {
	return &SKUGenerator{
		categoryPrefixes: map[string]string{
			// Clothing & Fashion
			"clothing":      "CLO",
			"mens-clothing": "MCL",
			"womens-clothing": "WCL",
			"kids-clothing": "KCL",
			"underwear":     "UND",
			"sportswear":    "SPO",
			"accessories":   "ACC",

			// Footwear
			"shoes":         "SHO",
			"mens-shoes":    "MSH",
			"womens-shoes":  "WSH",
			"kids-shoes":    "KSH",
			"sneakers":      "SNK",
			"boots":         "BOT",

			// Electronics
			"electronics":   "ELE",
			"smartphones":   "PHN",
			"laptops":       "LAP",
			"tablets":       "TAB",
			"headphones":    "HDP",
			"cameras":       "CAM",
			"smartwatches":  "SWT",

			// Home & Garden
			"home":          "HOM",
			"furniture":     "FUR",
			"kitchen":       "KIT",
			"bedding":       "BED",
			"garden":        "GAR",

			// Beauty & Health
			"beauty":        "BEA",
			"cosmetics":     "COS",
			"skincare":      "SKN",
			"haircare":      "HAR",
			"perfume":       "PRF",

			// Sports & Outdoors
			"sports":        "SPT",
			"fitness":       "FIT",
			"camping":       "CMP",
			"cycling":       "CYC",

			// Default
			"general":       "GEN",
		},
	}
}

// VariantAttributeForSKU represents attribute data needed for SKU generation
type VariantAttributeForSKU struct {
	Code       string  // e.g., "clothing_size", "color", "storage_capacity"
	ValueLabel string  // e.g., "M", "Black", "256GB"
	ValueID    *string // Optional: ID of the value
}

// GenerateSKU generates a unique SKU for a product variant
// Format: PREFIX-PRODID-ATTR1-ATTR2-...
// Example: CLO-A1B2C3-M-BLK (Clothing, product A1B2C3, size M, color Black)
func (g *SKUGenerator) GenerateSKU(productID string, categoryCode string, attributes []VariantAttributeForSKU) string {
	// Get category prefix (default: GEN)
	prefix := g.getCategoryPrefix(categoryCode)

	// Extract short product ID (first 6 chars of UUID)
	shortProductID := g.extractShortID(productID)

	// Build attribute suffixes
	suffixes := g.buildAttributeSuffixes(attributes)

	// Join all parts with hyphen
	parts := []string{prefix, shortProductID}
	parts = append(parts, suffixes...)

	return strings.Join(parts, "-")
}

// GenerateSKUFromDomain is a convenience method that works with domain.VariantAttributeValueV2
func (g *SKUGenerator) GenerateSKUFromDomain(productID string, categoryCode string, attrs []*domain.VariantAttributeValueV2) string {
	skuAttrs := make([]VariantAttributeForSKU, len(attrs))

	for i, attr := range attrs {
		skuAttrs[i] = VariantAttributeForSKU{
			Code:       g.getAttributeCode(attr.AttributeID), // Would need attribute lookup
			ValueLabel: g.getAttributeValueLabel(attr),
		}
	}

	return g.GenerateSKU(productID, categoryCode, skuAttrs)
}

// ValidateSKU checks if a SKU follows the correct format
func (g *SKUGenerator) ValidateSKU(sku string) error {
	if sku == "" {
		return fmt.Errorf("SKU cannot be empty")
	}

	// SKU format: PREFIX-XXXXX-... (minimum 3 parts)
	parts := strings.Split(sku, "-")
	if len(parts) < 2 {
		return fmt.Errorf("SKU must have at least PREFIX-PRODID format")
	}

	// Prefix must be 3 uppercase letters
	if len(parts[0]) != 3 {
		return fmt.Errorf("SKU prefix must be exactly 3 characters")
	}

	matched, _ := regexp.MatchString("^[A-Z]{3}$", parts[0])
	if !matched {
		return fmt.Errorf("SKU prefix must be 3 uppercase letters")
	}

	// Product ID must be 6 alphanumeric chars
	if len(parts[1]) != 6 {
		return fmt.Errorf("SKU product ID must be exactly 6 characters")
	}

	matched, _ = regexp.MatchString("^[A-Z0-9]{6}$", strings.ToUpper(parts[1]))
	if !matched {
		return fmt.Errorf("SKU product ID must be 6 alphanumeric characters")
	}

	return nil
}

// getCategoryPrefix returns the prefix for a category code
func (g *SKUGenerator) getCategoryPrefix(categoryCode string) string {
	// Normalize category code
	normalized := strings.ToLower(strings.TrimSpace(categoryCode))

	if prefix, ok := g.categoryPrefixes[normalized]; ok {
		return prefix
	}

	// Try to extract from path (e.g., "clothing/mens/tshirts" -> "clothing")
	if strings.Contains(normalized, "/") {
		parts := strings.Split(normalized, "/")
		if prefix, ok := g.categoryPrefixes[parts[0]]; ok {
			return prefix
		}
	}

	// Default prefix
	return "GEN"
}

// extractShortID extracts first 6 characters from UUID
func (g *SKUGenerator) extractShortID(productID string) string {
	// Remove hyphens from UUID
	cleaned := strings.ReplaceAll(productID, "-", "")

	// Take first 6 characters and uppercase
	if len(cleaned) >= 6 {
		return strings.ToUpper(cleaned[:6])
	}

	// Fallback: pad with zeros
	return strings.ToUpper(fmt.Sprintf("%-6s", cleaned))
}

// buildAttributeSuffixes builds SKU suffixes from variant attributes
func (g *SKUGenerator) buildAttributeSuffixes(attributes []VariantAttributeForSKU) []string {
	suffixes := make([]string, 0)

	for _, attr := range attributes {
		suffix := g.formatAttributeSuffix(attr)
		if suffix != "" {
			suffixes = append(suffixes, suffix)
		}
	}

	return suffixes
}

// formatAttributeSuffix formats a single attribute into SKU suffix
func (g *SKUGenerator) formatAttributeSuffix(attr VariantAttributeForSKU) string {
	code := strings.ToLower(attr.Code)
	value := strings.TrimSpace(attr.ValueLabel)

	switch code {
	case "clothing_size", "shoe_size_eu", "size":
		// Size: use as-is (e.g., "M", "42", "XL")
		return strings.ToUpper(value)

	case "color", "colour":
		// Color: first 3 letters, uppercase (e.g., "Black" -> "BLK")
		if len(value) >= 3 {
			return strings.ToUpper(value[:3])
		}
		return strings.ToUpper(value)

	case "storage_capacity", "memory", "capacity":
		// Storage: extract number + unit (e.g., "256 GB" -> "256GB", "1 TB" -> "1TB")
		return g.extractCapacity(value)

	case "material":
		// Material: first 3 letters (e.g., "Leather" -> "LEA")
		if len(value) >= 3 {
			return strings.ToUpper(value[:3])
		}
		return strings.ToUpper(value)

	case "weight", "volume":
		// Weight/Volume: extract number + unit (e.g., "500g" -> "500G")
		return g.extractMeasurement(value)

	default:
		// Generic: first 3 chars or full value if shorter
		cleaned := regexp.MustCompile(`[^A-Za-z0-9]`).ReplaceAllString(value, "")
		if len(cleaned) >= 3 {
			return strings.ToUpper(cleaned[:3])
		}
		return strings.ToUpper(cleaned)
	}
}

// extractCapacity extracts storage capacity (e.g., "256 GB" -> "256GB")
func (g *SKUGenerator) extractCapacity(value string) string {
	// Remove spaces and normalize
	normalized := strings.ToUpper(strings.ReplaceAll(value, " ", ""))

	// Extract number and unit
	re := regexp.MustCompile(`(\d+)\s*(GB|TB|MB)`)
	matches := re.FindStringSubmatch(normalized)

	if len(matches) >= 3 {
		return matches[1] + matches[2]
	}

	// Fallback: just extract numbers
	re = regexp.MustCompile(`\d+`)
	number := re.FindString(normalized)
	if number != "" {
		return number
	}

	return normalized
}

// extractMeasurement extracts weight/volume measurement (e.g., "500 g" -> "500G")
func (g *SKUGenerator) extractMeasurement(value string) string {
	normalized := strings.ToUpper(strings.ReplaceAll(value, " ", ""))

	// Extract number and unit
	re := regexp.MustCompile(`(\d+)\s*(G|KG|ML|L)`)
	matches := re.FindStringSubmatch(normalized)

	if len(matches) >= 3 {
		return matches[1] + matches[2]
	}

	// Fallback: just the value
	return normalized
}

// getAttributeCode would lookup attribute code by ID (requires attribute repository)
// For now, returns placeholder
func (g *SKUGenerator) getAttributeCode(attributeID int32) string {
	// TODO: Implement attribute lookup by ID
	// This would require injecting AttributeRepository
	return "attr"
}

// getAttributeValueLabel extracts the label from VariantAttributeValueV2
func (g *SKUGenerator) getAttributeValueLabel(attr *domain.VariantAttributeValueV2) string {
	if attr.ValueText != nil {
		return *attr.ValueText
	}
	if attr.ValueNumber != nil {
		return fmt.Sprintf("%.0f", *attr.ValueNumber)
	}
	if attr.ValueBoolean != nil {
		if *attr.ValueBoolean {
			return "YES"
		}
		return "NO"
	}
	return "UNK"
}
