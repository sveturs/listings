package models

import (
	"encoding/xml"
	"time"
)

// ImportJob represents a bulk import job
type ImportJob struct {
	ID           int       `json:"id" db:"id"`
	StorefrontID int       `json:"storefront_id" db:"storefront_id"`
	UserID       int       `json:"user_id" db:"user_id"`
	FileName     string    `json:"file_name" db:"file_name"`
	FileType     string    `json:"file_type" db:"file_type"` // xml, csv, zip
	FileURL      *string   `json:"file_url,omitempty" db:"file_url"` // for URL imports
	Status       string    `json:"status" db:"status"` // pending, processing, completed, failed
	TotalRecords int       `json:"total_records" db:"total_records"`
	ProcessedRecords int   `json:"processed_records" db:"processed_records"`
	SuccessfulRecords int  `json:"successful_records" db:"successful_records"`
	FailedRecords    int   `json:"failed_records" db:"failed_records"`
	ErrorMessage     *string `json:"error_message,omitempty" db:"error_message"`
	StartedAt        *time.Time `json:"started_at,omitempty" db:"started_at"`
	CompletedAt      *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// ImportError represents an error during import
type ImportError struct {
	ID        int    `json:"id" db:"id"`
	JobID     int    `json:"job_id" db:"job_id"`
	LineNumber int   `json:"line_number" db:"line_number"`
	FieldName  string `json:"field_name" db:"field_name"`
	ErrorMessage string `json:"error_message" db:"error_message"`
	RawData    string `json:"raw_data" db:"raw_data"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// DigitalVisionProduct represents the Serbian XML format from Digital Vision
type DigitalVisionProduct struct {
	ID          string `xml:"id" json:"id"`
	Sifra       string `xml:"sifra" json:"sifra"` // SKU/Code
	Naziv       string `xml:"naziv" json:"naziv"` // Name
	Kategorija1 string `xml:"kategorija1" json:"kategorija1"` // Category 1
	Kategorija2 string `xml:"kategorija2" json:"kategorija2"` // Category 2  
	Kategorija3 string `xml:"kategorija3" json:"kategorija3"` // Category 3
	Uvoznik     string `xml:"uvoznik" json:"uvoznik"` // Importer
	GodinaUvoza string `xml:"godinaUvoza" json:"godinaUvoza"` // Import year
	ZemljaPorekla string `xml:"zemljaPorekla" json:"zemljaPorekla"` // Country of origin
	VpCena      string `xml:"vpCena" json:"vpCena"` // Wholesale price
	MpCena      string `xml:"mpCena" json:"mpCena"` // Retail price
	Dostupan    string `xml:"dostupan" json:"dostupan"` // Available (1/0)
	NaAkciji    string `xml:"naAkciji" json:"naAkciji"` // On sale (1/0)
	Opis        string `xml:"opis" json:"opis"` // Description
	BarKod      string `xml:"barKod" json:"barKod"` // Barcode
	Slike       struct {
		Slika []string `xml:"slika" json:"slika"` // Images
	} `xml:"slike" json:"slike"`
}

// DigitalVisionCatalog represents the root XML structure
type DigitalVisionCatalog struct {
	XMLName  xml.Name               `xml:"artikli"`
	Products []DigitalVisionProduct `xml:"artikal"`
}

// CategoryMapping represents mapping between import categories and local categories
type CategoryMapping struct {
	ID              int    `json:"id" db:"id"`
	StorefrontID    int    `json:"storefront_id" db:"storefront_id"`
	ImportCategory1 string `json:"import_category1" db:"import_category1"`
	ImportCategory2 string `json:"import_category2" db:"import_category2"`
	ImportCategory3 string `json:"import_category3" db:"import_category3"`
	LocalCategoryID int    `json:"local_category_id" db:"local_category_id"`
	IsActive        bool   `json:"is_active" db:"is_active"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// ImportRequest represents a request to start an import
type ImportRequest struct {
	StorefrontID int    `json:"storefront_id" validate:"required"`
	FileType     string `json:"file_type" validate:"required,oneof=xml csv zip"`
	FileURL      *string `json:"file_url,omitempty" validate:"omitempty,url"`
	FileName     *string `json:"file_name,omitempty"`
	UpdateMode   string `json:"update_mode" validate:"oneof=create_only update_only upsert"` // create_only, update_only, upsert
	CategoryMappingMode string `json:"category_mapping_mode" validate:"oneof=auto manual skip"` // auto, manual, skip
}

// CSVProduct represents a product in CSV format
type CSVProduct struct {
	SKU           string  `csv:"sku" json:"sku"`
	Name          string  `csv:"name" json:"name"`
	Description   string  `csv:"description" json:"description"`
	Price         float64 `csv:"price" json:"price"`
	WholesalePrice *float64 `csv:"wholesale_price" json:"wholesale_price,omitempty"`
	Currency      string  `csv:"currency" json:"currency"`
	Category      string  `csv:"category" json:"category"`
	StockQuantity int     `csv:"stock_quantity" json:"stock_quantity"`
	Barcode       string  `csv:"barcode" json:"barcode"`
	ImageURL      string  `csv:"image_url" json:"image_url"`
	IsActive      bool    `csv:"is_active" json:"is_active"`
	OnSale        bool    `csv:"on_sale" json:"on_sale"`
	SalePrice     *float64 `csv:"sale_price" json:"sale_price,omitempty"`
	Brand         string  `csv:"brand" json:"brand"`
	Model         string  `csv:"model" json:"model"`
	CountryOfOrigin string `csv:"country_of_origin" json:"country_of_origin"`
}

// ImportValidationError represents validation errors during import
type ImportValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ImportProductRequest represents a processed product ready for import
type ImportProductRequest struct {
	ExternalID    string                 `json:"external_id,omitempty"`
	SKU           string                 `json:"sku,omitempty"`
	Name          string                 `json:"name" validate:"required"`
	Description   string                 `json:"description"`
	Price         float64                `json:"price" validate:"required,min=0"`
	WholesalePrice *float64              `json:"wholesale_price,omitempty" validate:"omitempty,min=0"`
	Currency      string                 `json:"currency" validate:"required,len=3"`
	CategoryID    int                    `json:"category_id" validate:"required"`
	StockQuantity int                    `json:"stock_quantity" validate:"min=0"`
	Barcode       string                 `json:"barcode,omitempty"`
	ImageURLs     []string               `json:"image_urls,omitempty"`
	IsActive      bool                   `json:"is_active"`
	OnSale        bool                   `json:"on_sale"`
	SalePrice     *float64               `json:"sale_price,omitempty" validate:"omitempty,min=0"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
}

// ImportJobStatus represents the status of an import job
type ImportJobStatus struct {
	ID                int                `json:"id"`
	Status            string             `json:"status"`
	Progress          float64            `json:"progress"` // 0-100
	TotalRecords      int                `json:"total_records"`
	ProcessedRecords  int                `json:"processed_records"`
	SuccessfulRecords int                `json:"successful_records"`
	FailedRecords     int                `json:"failed_records"`
	Errors            []ImportError      `json:"errors,omitempty"`
	StartedAt         *time.Time         `json:"started_at,omitempty"`
	CompletedAt       *time.Time         `json:"completed_at,omitempty"`
	EstimatedTimeLeft *int               `json:"estimated_time_left,omitempty"` // seconds
}

// ImportSummary represents a summary of completed import
type ImportSummary struct {
	JobID             int       `json:"job_id"`
	TotalRecords      int       `json:"total_records"`
	SuccessfulRecords int       `json:"successful_records"`
	FailedRecords     int       `json:"failed_records"`
	NewProducts       int       `json:"new_products"`
	UpdatedProducts   int       `json:"updated_products"`
	SkippedProducts   int       `json:"skipped_products"`
	ProcessingTime    string    `json:"processing_time"`
	CompletedAt       time.Time `json:"completed_at"`
}