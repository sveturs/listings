package service

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/proj/storefronts/parsers"
)

const (
	// Status values
	statusFailed = "failed"

	// File types
	fileTypeXML = "xml"
	fileTypeCSV = "csv"
	fileTypeZIP = "zip"
)

// ImportService handles product import operations
type ImportService struct {
	productService *ProductService
	// Add repository for import jobs when needed
}

// NewImportService creates a new import service
func NewImportService(productService *ProductService) *ImportService {
	return &ImportService{
		productService: productService,
	}
}

// ImportFromURL downloads and imports products from a URL
func (s *ImportService) ImportFromURL(ctx context.Context, req models.ImportRequest) (*models.ImportJob, error) {
	if req.FileURL == nil {
		return nil, fmt.Errorf("file URL is required")
	}

	// Create import job
	job := &models.ImportJob{
		StorefrontID: req.StorefrontID,
		FileType:     req.FileType,
		FileURL:      req.FileURL,
		Status:       "pending",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Download file
	data, _, err := s.downloadFile(ctx, *req.FileURL)
	if err != nil {
		job.Status = statusFailed
		job.ErrorMessage = &[]string{fmt.Sprintf("Failed to download file: %v", err)}[0]
		return job, err
	}

	// Update job status
	job.Status = "processing"
	startTime := time.Now()
	job.StartedAt = &startTime

	// Process file based on type
	var products []models.ImportProductRequest
	var validationErrors []models.ImportValidationError

	switch req.FileType {
	case fileTypeXML:
		products, validationErrors, err = s.processXMLData(data, req.StorefrontID)
	case fileTypeCSV:
		products, validationErrors, err = s.processCSVData(data, req.StorefrontID)
	case fileTypeZIP:
		products, validationErrors, err = s.processZIPData(data, req.StorefrontID)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", req.FileType)
	}

	if err != nil {
		job.Status = statusFailed
		job.ErrorMessage = &[]string{err.Error()}[0]
		return job, err
	}

	// Update job with totals
	job.TotalRecords = len(products) + len(validationErrors)
	job.FailedRecords = len(validationErrors)

	// Import products
	successCount, importErrors := s.importProducts(ctx, products, req.StorefrontID, req.UpdateMode)

	job.ProcessedRecords = len(products)
	job.SuccessfulRecords = successCount
	job.FailedRecords += len(importErrors)

	// Complete job
	completedTime := time.Now()
	job.CompletedAt = &completedTime
	job.Status = "completed"

	if len(importErrors) > 0 {
		// Store errors in job (simplified - in real implementation, store in separate table)
		errorMsg := fmt.Sprintf("Import completed with %d errors", len(importErrors))
		job.ErrorMessage = &errorMsg
	}

	return job, nil
}

// ImportFromFile imports products from uploaded file data
func (s *ImportService) ImportFromFile(ctx context.Context, fileData []byte, req models.ImportRequest) (*models.ImportJob, error) {
	// Create import job
	job := &models.ImportJob{
		StorefrontID: req.StorefrontID,
		FileType:     req.FileType,
		Status:       "processing",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if req.FileName != nil {
		job.FileName = *req.FileName
	}

	startTime := time.Now()
	job.StartedAt = &startTime

	// Process file based on type
	var products []models.ImportProductRequest
	var validationErrors []models.ImportValidationError
	var err error

	switch req.FileType {
	case fileTypeXML:
		products, validationErrors, err = s.processXMLData(fileData, req.StorefrontID)
	case fileTypeCSV:
		products, validationErrors, err = s.processCSVData(fileData, req.StorefrontID)
	case fileTypeZIP:
		products, validationErrors, err = s.processZIPData(fileData, req.StorefrontID)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", req.FileType)
	}

	if err != nil {
		job.Status = statusFailed
		job.ErrorMessage = &[]string{err.Error()}[0]
		return job, err
	}

	// Update job with totals
	job.TotalRecords = len(products) + len(validationErrors)
	job.FailedRecords = len(validationErrors)

	// Import products
	successCount, importErrors := s.importProducts(ctx, products, req.StorefrontID, req.UpdateMode)

	job.ProcessedRecords = len(products)
	job.SuccessfulRecords = successCount
	job.FailedRecords += len(importErrors)

	// Complete job
	completedTime := time.Now()
	job.CompletedAt = &completedTime
	job.Status = "completed"

	if len(importErrors) > 0 {
		errorMsg := fmt.Sprintf("Import completed with %d errors", len(importErrors))
		job.ErrorMessage = &errorMsg
	}

	return job, nil
}

// downloadFile downloads a file from URL
func (s *ImportService) downloadFile(ctx context.Context, url string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Логирование ошибки закрытия Body
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return data, resp.Header.Get("Content-Type"), nil
}

// processXMLData processes XML data
func (s *ImportService) processXMLData(data []byte, storefrontID int) ([]models.ImportProductRequest, []models.ImportValidationError, error) {
	parser := parsers.NewXMLParser(storefrontID)

	// Try Digital Vision format first
	products, errors, err := parser.ParseDigitalVisionXML(data)
	if err != nil {
		// Log the Digital Vision parsing error for debugging
		fmt.Printf("Digital Vision XML parsing failed: %v\n", err)
		// If Digital Vision format fails, try generic XML (but it's not implemented)
		return nil, nil, fmt.Errorf("digital Vision XML parsing failed: %v (generic XML parsing not implemented)", err)
	}

	return products, errors, nil
}

// processCSVData processes CSV data
func (s *ImportService) processCSVData(data []byte, storefrontID int) ([]models.ImportProductRequest, []models.ImportValidationError, error) {
	parser := parsers.NewCSVParser(storefrontID)
	reader := bytes.NewReader(data)
	return parser.ParseCSV(reader)
}

// processZIPData processes ZIP archive data
func (s *ImportService) processZIPData(data []byte, storefrontID int) ([]models.ImportProductRequest, []models.ImportValidationError, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read ZIP archive: %w", err)
	}

	var allProducts []models.ImportProductRequest
	var allErrors []models.ImportValidationError

	// Process each file in the ZIP
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		// Open file
		rc, err := file.Open()
		if err != nil {
			continue
		}

		// Read file content
		fileData, err := io.ReadAll(rc)
		func() {
			if err := rc.Close(); err != nil {
				fmt.Printf("Failed to close file: %v", err)
			}
		}()
		if err != nil {
			continue
		}

		// Determine file type by extension
		ext := strings.ToLower(filepath.Ext(file.Name))
		var products []models.ImportProductRequest
		var errors []models.ImportValidationError

		switch ext {
		case ".xml":
			products, errors, err = s.processXMLData(fileData, storefrontID)
		case ".csv":
			products, errors, err = s.processCSVData(fileData, storefrontID)
		default:
			continue // Skip unsupported files
		}

		if err == nil {
			allProducts = append(allProducts, products...)
			allErrors = append(allErrors, errors...)
		}
	}

	return allProducts, allErrors, nil
}

// importProducts imports validated products
func (s *ImportService) importProducts(ctx context.Context, products []models.ImportProductRequest, storefrontID int, updateMode string) (int, []error) {
	var successCount int
	var errors []error

	for _, product := range products {
		err := s.importSingleProduct(ctx, product, storefrontID, updateMode)
		if err != nil {
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	return successCount, errors
}

// importSingleProduct imports a single product
func (s *ImportService) importSingleProduct(ctx context.Context, importProduct models.ImportProductRequest, storefrontID int, updateMode string) error {
	// Check if product already exists by SKU or external ID
	var existingProduct *models.StorefrontProduct
	var err error

	if importProduct.SKU != "" {
		// existingProduct, err = s.productService.GetProductBySKU(ctx, storefrontID, importProduct.SKU)
		// For now, treat all as new products
		_ = err
	}

	switch updateMode {
	case "create_only":
		if existingProduct != nil {
			return fmt.Errorf("product with SKU %s already exists", importProduct.SKU)
		}
		return s.createProduct(ctx, importProduct, storefrontID)

	case "update_only":
		if existingProduct == nil {
			return fmt.Errorf("product with SKU %s not found", importProduct.SKU)
		}
		return s.updateProduct(ctx, existingProduct.ID, importProduct, storefrontID)

	case "upsert":
		if existingProduct != nil {
			return s.updateProduct(ctx, existingProduct.ID, importProduct, storefrontID)
		}
		return s.createProduct(ctx, importProduct, storefrontID)

	default:
		return s.createProduct(ctx, importProduct, storefrontID)
	}
}

// createProduct creates a new product
func (s *ImportService) createProduct(ctx context.Context, importProduct models.ImportProductRequest, storefrontID int) error {
	createReq := models.CreateProductRequest{
		Name:          importProduct.Name,
		Description:   importProduct.Description,
		Price:         importProduct.Price,
		Currency:      importProduct.Currency,
		CategoryID:    importProduct.CategoryID,
		SKU:           &importProduct.SKU,
		Barcode:       &importProduct.Barcode,
		StockQuantity: importProduct.StockQuantity,
		IsActive:      importProduct.IsActive,
		Attributes:    importProduct.Attributes,
	}

	// Create product using import-specific method (bypasses ownership check)
	product, err := s.productService.CreateProductForImport(ctx, storefrontID, &createReq)
	if err != nil {
		fmt.Printf("Failed to create product %s: %v\n", importProduct.Name, err)
		return fmt.Errorf("failed to create product: %w", err)
	}
	fmt.Printf("Successfully created product: %s (ID: %d)\n", product.Name, product.ID)

	// Add images if provided
	if len(importProduct.ImageURLs) > 0 {
		for i, imageURL := range importProduct.ImageURLs {
			imageReq := models.StorefrontProductImage{
				StorefrontProductID: product.ID,
				ImageURL:            imageURL,
				DisplayOrder:        i + 1,
				IsDefault:           i == 0,
			}
			// Add image (would need to implement this in product service)
			_ = imageReq
		}
	}

	return nil
}

// updateProduct updates an existing product
func (s *ImportService) updateProduct(ctx context.Context, productID int, importProduct models.ImportProductRequest, storefrontID int) error {
	updateReq := models.UpdateProductRequest{
		Name:          &importProduct.Name,
		Description:   &importProduct.Description,
		Price:         &importProduct.Price,
		CategoryID:    &importProduct.CategoryID,
		SKU:           &importProduct.SKU,
		Barcode:       &importProduct.Barcode,
		StockQuantity: &importProduct.StockQuantity,
		IsActive:      &importProduct.IsActive,
		Attributes:    importProduct.Attributes,
	}

	// Update product using import-specific method (bypasses ownership check)
	err := s.productService.UpdateProductForImport(ctx, storefrontID, productID, &updateReq)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

// GetCSVTemplate returns a CSV template for product import
func (s *ImportService) GetCSVTemplate() [][]string {
	parser := parsers.NewCSVParser(0)
	return parser.GenerateCSVTemplate()
}

// ValidateImportFile validates an import file without importing
func (s *ImportService) ValidateImportFile(ctx context.Context, fileData []byte, fileType string, storefrontID int) (*models.ImportJobStatus, error) {
	var products []models.ImportProductRequest
	var validationErrors []models.ImportValidationError
	var err error

	switch fileType {
	case fileTypeXML:
		products, validationErrors, err = s.processXMLData(fileData, storefrontID)
	case fileTypeCSV:
		products, validationErrors, err = s.processCSVData(fileData, storefrontID)
	case fileTypeZIP:
		products, validationErrors, err = s.processZIPData(fileData, storefrontID)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}

	if err != nil {
		return nil, err
	}

	status := &models.ImportJobStatus{
		Status:            "validated",
		TotalRecords:      len(products) + len(validationErrors),
		ProcessedRecords:  len(products),
		SuccessfulRecords: len(products),
		FailedRecords:     len(validationErrors),
		Progress:          100.0,
	}

	return status, nil
}

// GetJobs returns list of import jobs for a storefront
func (s *ImportService) GetJobs(ctx context.Context, storefrontID int, status string, limit, offset int) (*models.ImportJobsResponse, error) {
	// TODO: Implement database queries for import jobs
	// For now, return mock data
	jobs := []models.ImportJob{
		{
			ID:                1,
			StorefrontID:      storefrontID,
			FileName:          "test.xml",
			FileType:          "xml",
			Status:            "completed",
			TotalRecords:      2,
			ProcessedRecords:  2,
			SuccessfulRecords: 0,
			FailedRecords:     2,
			ErrorMessage:      stringPtr("Import completed with errors"),
			CreatedAt:         time.Now().Add(-time.Hour),
			UpdatedAt:         time.Now().Add(-time.Minute * 30),
		},
	}

	// Filter by status if provided
	if status != "" {
		var filteredJobs []models.ImportJob
		for _, job := range jobs {
			if job.Status == status {
				filteredJobs = append(filteredJobs, job)
			}
		}
		jobs = filteredJobs
	}

	return &models.ImportJobsResponse{
		Jobs:  jobs,
		Total: len(jobs),
	}, nil
}

// GetJobDetails returns detailed information about an import job
func (s *ImportService) GetJobDetails(ctx context.Context, jobID int) (*models.ImportJob, error) {
	// TODO: Implement database query for job details
	// For now, return mock data
	job := &models.ImportJob{
		ID:                jobID,
		StorefrontID:      4, // Mock storefront ID
		FileName:          "test.xml",
		FileType:          "xml",
		Status:            "completed",
		TotalRecords:      2,
		ProcessedRecords:  2,
		SuccessfulRecords: 0,
		FailedRecords:     2,
		ErrorMessage:      stringPtr("Import completed with errors"),
		CreatedAt:         time.Now().Add(-time.Hour),
		UpdatedAt:         time.Now().Add(-time.Minute * 30),
	}

	return job, nil
}

// GetJobStatus returns status of an import job
func (s *ImportService) GetJobStatus(ctx context.Context, jobID int) (*models.ImportJobStatus, error) {
	// TODO: Implement database query for job status
	// For now, return mock data
	status := &models.ImportJobStatus{
		Status:            "completed",
		TotalRecords:      2,
		ProcessedRecords:  2,
		SuccessfulRecords: 0,
		FailedRecords:     2,
		Progress:          100.0,
	}

	return status, nil
}

// CancelJob cancels a running import job
func (s *ImportService) CancelJob(ctx context.Context, jobID int) error {
	// TODO: Implement job cancellation logic
	return nil
}

// RetryJob retries a failed import job
func (s *ImportService) RetryJob(ctx context.Context, jobID int) (*models.ImportJob, error) {
	// TODO: Implement job retry logic
	// For now, return a new mock job
	job := &models.ImportJob{
		ID:                jobID + 100, // New job ID
		StorefrontID:      4,           // Mock storefront ID
		FileName:          "test.xml",
		FileType:          "xml",
		Status:            "pending",
		TotalRecords:      0,
		ProcessedRecords:  0,
		SuccessfulRecords: 0,
		FailedRecords:     0,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	return job, nil
}

// Helper function to convert string to string pointer
func stringPtr(s string) *string {
	return &s
}
