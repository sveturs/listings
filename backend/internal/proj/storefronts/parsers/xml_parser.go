package parsers

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"backend/internal/domain/models"
)

// XMLParser handles parsing of XML product files
type XMLParser struct {
	storefrontID int
}

// NewXMLParser creates a new XML parser
func NewXMLParser(storefrontID int) *XMLParser {
	return &XMLParser{
		storefrontID: storefrontID,
	}
}

// ParseDigitalVisionXML parses Digital Vision XML format
func (p *XMLParser) ParseDigitalVisionXML(xmlData []byte) ([]models.ImportProductRequest, []models.ImportValidationError, error) {
	var catalog models.DigitalVisionCatalog
	if err := xml.Unmarshal(xmlData, &catalog); err != nil {
		return nil, nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	products := make([]models.ImportProductRequest, 0, len(catalog.Products))
	var validationErrors []models.ImportValidationError

	for i, dvProduct := range catalog.Products {
		product, errs := p.convertDigitalVisionProduct(dvProduct, i+1)
		if len(errs) > 0 {
			validationErrors = append(validationErrors, errs...)
			continue
		}
		products = append(products, product)
	}

	return products, validationErrors, nil
}

// convertDigitalVisionProduct converts Digital Vision product to ImportProductRequest
func (p *XMLParser) convertDigitalVisionProduct(dv models.DigitalVisionProduct, lineNumber int) (models.ImportProductRequest, []models.ImportValidationError) {
	var errors []models.ImportValidationError
	product := models.ImportProductRequest{}

	// Required fields validation
	if strings.TrimSpace(dv.Naziv) == "" {
		errors = append(errors, models.ImportValidationError{
			Field:   "naziv",
			Message: "Product name is required",
			Value:   dv.Naziv,
		})
	} else {
		product.Name = cleanCDATA(dv.Naziv)
	}

	// Parse prices
	retailPrice, err := parsePrice(dv.MpCena)
	if err != nil {
		errors = append(errors, models.ImportValidationError{
			Field:   "mpCena",
			Message: fmt.Sprintf("Invalid retail price: %v", err),
			Value:   dv.MpCena,
		})
	} else {
		product.Price = retailPrice
	}

	// Parse wholesale price (optional)
	if strings.TrimSpace(dv.VpCena) != "" {
		wholesalePrice, err := parsePrice(dv.VpCena)
		if err != nil {
			errors = append(errors, models.ImportValidationError{
				Field:   "vpCena",
				Message: fmt.Sprintf("Invalid wholesale price: %v", err),
				Value:   dv.VpCena,
			})
		} else {
			product.WholesalePrice = &wholesalePrice
		}
	}

	// Set basic fields
	product.ExternalID = dv.ID
	product.SKU = dv.Sifra
	product.Description = cleanCDATA(dv.Opis)
	product.Currency = "RSD" // Serbian Dinar
	product.Barcode = dv.BarKod

	// Parse availability
	product.IsActive = dv.Dostupan == "1"
	product.OnSale = dv.NaAkciji == "1"

	// Set default stock quantity (since it's not in the XML format)
	if product.IsActive {
		product.StockQuantity = 1 // Default available quantity
	} else {
		product.StockQuantity = 0
	}

	// Parse images
	if len(dv.Slike.Slika) > 0 {
		product.ImageURLs = make([]string, 0, len(dv.Slike.Slika))
		for _, imageURL := range dv.Slike.Slika {
			cleanURL := cleanCDATA(imageURL)
			if cleanURL != "" {
				product.ImageURLs = append(product.ImageURLs, cleanURL)
			}
		}
	}

	// Set attributes from additional fields
	product.Attributes = map[string]interface{}{
		"kategorija1":     cleanCDATA(dv.Kategorija1),
		"kategorija2":     cleanCDATA(dv.Kategorija2),
		"kategorija3":     cleanCDATA(dv.Kategorija3),
		"uvoznik":         dv.Uvoznik,
		"godina_uvoza":    dv.GodinaUvoza,
		"zemlja_porekla":  dv.ZemljaPorekla,
	}

	// Category mapping will be handled by the service layer
	// For now, set a default category ID (will be mapped later)
	product.CategoryID = 1 // Default category

	return product, errors
}

// ParseGenericXML parses generic XML format
func (p *XMLParser) ParseGenericXML(xmlData []byte) ([]models.ImportProductRequest, []models.ImportValidationError, error) {
	// Implementation for generic XML format
	// This can be extended based on other XML formats encountered
	return nil, nil, fmt.Errorf("generic XML parsing not implemented yet")
}

// Helper functions

// cleanCDATA removes CDATA wrapper from XML content
func cleanCDATA(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "<![CDATA[") && strings.HasSuffix(content, "]]>") {
		content = content[9 : len(content)-3]
	}
	return strings.TrimSpace(content)
}

// parsePrice parses price string to float64
func parsePrice(priceStr string) (float64, error) {
	priceStr = strings.TrimSpace(priceStr)
	if priceStr == "" {
		return 0, fmt.Errorf("empty price")
	}

	// Remove any non-numeric characters except dots and commas
	cleanPrice := strings.ReplaceAll(priceStr, ",", ".")
	
	// Parse as float
	price, err := strconv.ParseFloat(cleanPrice, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid price format: %s", priceStr)
	}

	if price < 0 {
		return 0, fmt.Errorf("price cannot be negative: %f", price)
	}

	return price, nil
}

// ValidateProduct validates a product before import
func (p *XMLParser) ValidateProduct(product models.ImportProductRequest) []models.ImportValidationError {
	var errors []models.ImportValidationError

	if strings.TrimSpace(product.Name) == "" {
		errors = append(errors, models.ImportValidationError{
			Field:   "name",
			Message: "Product name is required",
		})
	}

	if product.Price <= 0 {
		errors = append(errors, models.ImportValidationError{
			Field:   "price",
			Message: "Product price must be greater than 0",
			Value:   product.Price,
		})
	}

	if product.Currency == "" {
		errors = append(errors, models.ImportValidationError{
			Field:   "currency",
			Message: "Currency is required",
		})
	}

	if len(product.Currency) != 3 {
		errors = append(errors, models.ImportValidationError{
			Field:   "currency",
			Message: "Currency must be 3 characters long",
			Value:   product.Currency,
		})
	}

	if product.StockQuantity < 0 {
		errors = append(errors, models.ImportValidationError{
			Field:   "stock_quantity",
			Message: "Stock quantity cannot be negative",
			Value:   product.StockQuantity,
		})
	}

	if product.WholesalePrice != nil && *product.WholesalePrice < 0 {
		errors = append(errors, models.ImportValidationError{
			Field:   "wholesale_price",
			Message: "Wholesale price cannot be negative",
			Value:   *product.WholesalePrice,
		})
	}

	if product.SalePrice != nil && *product.SalePrice < 0 {
		errors = append(errors, models.ImportValidationError{
			Field:   "sale_price",
			Message: "Sale price cannot be negative",
			Value:   *product.SalePrice,
		})
	}

	return errors
}