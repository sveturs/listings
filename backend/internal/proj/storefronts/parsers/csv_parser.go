package parsers

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"backend/internal/domain/models"
)

// CSVParser handles parsing of CSV product files
type CSVParser struct {
	storefrontID int
}

// NewCSVParser creates a new CSV parser
func NewCSVParser(storefrontID int) *CSVParser {
	return &CSVParser{
		storefrontID: storefrontID,
	}
}

// ParseCSV parses CSV data and returns products and validation errors
func (p *CSVParser) ParseCSV(csvData io.Reader) ([]models.ImportProductRequest, []models.ImportValidationError, error) {
	reader := csv.NewReader(csvData)
	reader.Comma = ','
	reader.TrimLeadingSpace = true

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	// Create header map for column index lookup
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.ToLower(strings.TrimSpace(header))] = i
	}

	// Validate required headers
	if err := p.validateHeaders(headerMap); err != nil {
		return nil, nil, fmt.Errorf("invalid CSV headers: %w", err)
	}

	products := make([]models.ImportProductRequest, 0)
	var validationErrors []models.ImportValidationError
	lineNumber := 1 // Start from 1 (header is line 0)

	// Read data rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read CSV line %d: %w", lineNumber+1, err)
		}

		lineNumber++
		product, errs := p.parseCSVRecord(record, headerMap, lineNumber)
		
		if len(errs) > 0 {
			validationErrors = append(validationErrors, errs...)
			continue
		}

		products = append(products, product)
	}

	return products, validationErrors, nil
}

// validateHeaders checks if required headers are present
func (p *CSVParser) validateHeaders(headerMap map[string]int) error {
	requiredHeaders := []string{"name", "price", "currency"}
	
	var missingHeaders []string
	for _, header := range requiredHeaders {
		if _, exists := headerMap[header]; !exists {
			missingHeaders = append(missingHeaders, header)
		}
	}

	if len(missingHeaders) > 0 {
		return fmt.Errorf("missing required headers: %s", strings.Join(missingHeaders, ", "))
	}

	return nil
}

// parseCSVRecord parses a single CSV record into ImportProductRequest
func (p *CSVParser) parseCSVRecord(record []string, headerMap map[string]int, lineNumber int) (models.ImportProductRequest, []models.ImportValidationError) {
	var errors []models.ImportValidationError
	product := models.ImportProductRequest{}

	// Helper function to get field value safely
	getField := func(fieldName string) string {
		if idx, exists := headerMap[fieldName]; exists && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	// Helper function to add validation error
	addError := func(field, message string, value interface{}) {
		errors = append(errors, models.ImportValidationError{
			Field:   field,
			Message: fmt.Sprintf("Line %d: %s", lineNumber, message),
			Value:   value,
		})
	}

	// Parse required fields
	product.Name = getField("name")
	if product.Name == "" {
		addError("name", "Product name is required", "")
	}

	// Parse price
	priceStr := getField("price")
	if priceStr == "" {
		addError("price", "Price is required", "")
	} else {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			addError("price", fmt.Sprintf("Invalid price format: %v", err), priceStr)
		} else if price <= 0 {
			addError("price", "Price must be greater than 0", price)
		} else {
			product.Price = price
		}
	}

	// Parse currency
	product.Currency = getField("currency")
	if product.Currency == "" {
		product.Currency = "RSD" // Default to Serbian Dinar
	} else if len(product.Currency) != 3 {
		addError("currency", "Currency must be 3 characters long", product.Currency)
	}

	// Parse optional fields
	product.SKU = getField("sku")
	product.Description = getField("description")
	product.Barcode = getField("barcode")

	// Parse wholesale price
	wholesalePriceStr := getField("wholesale_price")
	if wholesalePriceStr != "" {
		wholesalePrice, err := strconv.ParseFloat(wholesalePriceStr, 64)
		if err != nil {
			addError("wholesale_price", fmt.Sprintf("Invalid wholesale price format: %v", err), wholesalePriceStr)
		} else if wholesalePrice < 0 {
			addError("wholesale_price", "Wholesale price cannot be negative", wholesalePrice)
		} else {
			product.WholesalePrice = &wholesalePrice
		}
	}

	// Parse stock quantity
	stockStr := getField("stock_quantity")
	if stockStr == "" {
		product.StockQuantity = 0
	} else {
		stock, err := strconv.Atoi(stockStr)
		if err != nil {
			addError("stock_quantity", fmt.Sprintf("Invalid stock quantity format: %v", err), stockStr)
		} else if stock < 0 {
			addError("stock_quantity", "Stock quantity cannot be negative", stock)
		} else {
			product.StockQuantity = stock
		}
	}

	// Parse boolean fields
	product.IsActive = parseBool(getField("is_active"), true) // Default to active
	product.OnSale = parseBool(getField("on_sale"), false)

	// Parse sale price
	salePriceStr := getField("sale_price")
	if salePriceStr != "" {
		salePrice, err := strconv.ParseFloat(salePriceStr, 64)
		if err != nil {
			addError("sale_price", fmt.Sprintf("Invalid sale price format: %v", err), salePriceStr)
		} else if salePrice < 0 {
			addError("sale_price", "Sale price cannot be negative", salePrice)
		} else {
			product.SalePrice = &salePrice
		}
	}

	// Parse image URLs (comma-separated)
	imageURLStr := getField("image_url")
	if imageURLStr != "" {
		imageURLs := strings.Split(imageURLStr, ",")
		product.ImageURLs = make([]string, 0, len(imageURLs))
		for _, url := range imageURLs {
			url = strings.TrimSpace(url)
			if url != "" {
				product.ImageURLs = append(product.ImageURLs, url)
			}
		}
	}

	// Parse category
	categoryStr := getField("category")
	if categoryStr != "" {
		// For now, set default category ID (will be mapped by service layer)
		product.CategoryID = 1
		
		// Store original category in attributes for mapping
		if product.Attributes == nil {
			product.Attributes = make(map[string]interface{})
		}
		product.Attributes["original_category"] = categoryStr
	} else {
		product.CategoryID = 1 // Default category
	}

	// Parse additional attributes
	brand := getField("brand")
	if brand != "" {
		if product.Attributes == nil {
			product.Attributes = make(map[string]interface{})
		}
		product.Attributes["brand"] = brand
	}

	model := getField("model")
	if model != "" {
		if product.Attributes == nil {
			product.Attributes = make(map[string]interface{})
		}
		product.Attributes["model"] = model
	}

	countryOfOrigin := getField("country_of_origin")
	if countryOfOrigin != "" {
		if product.Attributes == nil {
			product.Attributes = make(map[string]interface{})
		}
		product.Attributes["country_of_origin"] = countryOfOrigin
	}

	return product, errors
}

// parseBool converts string to boolean with default value
func parseBool(value string, defaultValue bool) bool {
	if value == "" {
		return defaultValue
	}
	
	value = strings.ToLower(value)
	switch value {
	case "true", "1", "yes", "y", "da", "активан":
		return true
	case "false", "0", "no", "n", "ne", "неактиван":
		return false
	default:
		return defaultValue
	}
}

// GetSupportedHeaders returns the list of supported CSV headers
func (p *CSVParser) GetSupportedHeaders() []string {
	return []string{
		"sku",                // Product SKU/Code
		"name",              // Product name (required)
		"description",       // Product description
		"price",             // Retail price (required)
		"wholesale_price",   // Wholesale price (optional)
		"currency",          // Currency code (required, defaults to RSD)
		"category",          // Category name
		"stock_quantity",    // Stock quantity
		"barcode",           // Product barcode
		"image_url",         // Image URLs (comma-separated)
		"is_active",         // Active status (true/false)
		"on_sale",           // On sale status (true/false)
		"sale_price",        // Sale price (optional)
		"brand",             // Product brand
		"model",             // Product model
		"country_of_origin", // Country of origin
	}
}

// GenerateCSVTemplate generates a CSV template with headers and example data
func (p *CSVParser) GenerateCSVTemplate() [][]string {
	headers := p.GetSupportedHeaders()
	exampleData := []string{
		"SKU001",                           // sku
		"Example Product Name",             // name
		"This is an example product",       // description
		"1000.00",                         // price
		"800.00",                          // wholesale_price
		"RSD",                             // currency
		"Electronics",                     // category
		"50",                              // stock_quantity
		"1234567890123",                   // barcode
		"https://example.com/image1.jpg,https://example.com/image2.jpg", // image_url
		"true",                            // is_active
		"false",                           // on_sale
		"900.00",                          // sale_price
		"ExampleBrand",                    // brand
		"Model123",                        // model
		"Serbia",                          // country_of_origin
	}

	return [][]string{headers, exampleData}
}