package service

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/proj/storefronts/parsers"
	"backend/internal/services"
	"backend/internal/storage/postgres"

	_ "golang.org/x/image/webp"
)

const (
	// Status values
	statusPending    = "pending"
	statusProcessing = "processing"
	statusCompleted  = "completed"
	statusFailed     = "failed"
	statusCanceled   = "canceled"

	// File types
	fileTypeXML = "xml"
	fileTypeCSV = "csv"
	fileTypeZIP = "zip"

	// Update modes
	updateModeCreateOnly = "create_only"
	updateModeUpdateOnly = "update_only"
	updateModeUpsert     = "upsert"
)

// ImportService handles product import operations
type ImportService struct {
	productService         *ProductService
	jobsRepo               postgres.ImportJobsRepositoryInterface
	queueManager           *ImportQueueManager
	imageService           *services.ImageService
	categoryMappingService *CategoryMappingService
}

// NewImportService creates a new import service
func NewImportService(
	productService *ProductService,
	jobsRepo postgres.ImportJobsRepositoryInterface,
	imageService *services.ImageService,
	categoryMappingService *CategoryMappingService,
) *ImportService {
	return &ImportService{
		productService:         productService,
		jobsRepo:               jobsRepo,
		queueManager:           nil, // Will be set via SetQueueManager
		imageService:           imageService,
		categoryMappingService: categoryMappingService,
	}
}

// SetQueueManager sets the queue manager for async processing
func (s *ImportService) SetQueueManager(queueManager *ImportQueueManager) {
	s.queueManager = queueManager
}

// GetQueueManager returns the queue manager
func (s *ImportService) GetQueueManager() *ImportQueueManager {
	return s.queueManager
}

// SetCategoryMappingService sets the category mapping service
func (s *ImportService) SetCategoryMappingService(categoryMappingSvc *CategoryMappingService) {
	s.categoryMappingService = categoryMappingSvc
}

// ImportFromURL downloads and imports products from a URL
func (s *ImportService) ImportFromURL(ctx context.Context, userID int, req models.ImportRequest) (*models.ImportJob, error) {
	if req.FileURL == nil {
		return nil, fmt.Errorf("file URL is required")
	}

	// Create import job in database
	job := &models.ImportJob{
		StorefrontID: req.StorefrontID,
		UserID:       userID,
		FileType:     req.FileType,
		FileURL:      req.FileURL,
		Status:       "pending",
	}

	// Save job to database
	if err := s.jobsRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create import job: %w", err)
	}

	// Download file
	data, _, err := s.downloadFile(ctx, *req.FileURL)
	if err != nil {
		job.Status = statusFailed
		errorMsg := fmt.Sprintf("Failed to download file: %v", err)
		job.ErrorMessage = &errorMsg
		_ = s.jobsRepo.Update(ctx, job)
		return job, err
	}

	// Update job status to processing
	job.Status = statusProcessing
	startTime := time.Now()
	job.StartedAt = &startTime
	if err := s.jobsRepo.Update(ctx, job); err != nil {
		return job, fmt.Errorf("failed to update job status: %w", err)
	}

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
		job.Status = statusFailed
		errorMsg := fmt.Sprintf("unsupported file type: %s", req.FileType)
		job.ErrorMessage = &errorMsg
		_ = s.jobsRepo.Update(ctx, job)
		return job, errors.New(errorMsg)
	}

	if err != nil {
		job.Status = statusFailed
		errorMsg := err.Error()
		job.ErrorMessage = &errorMsg
		_ = s.jobsRepo.Update(ctx, job)
		return job, err
	}

	// Update job with totals
	job.TotalRecords = len(products) + len(validationErrors)
	job.FailedRecords = len(validationErrors)

	// Save validation errors to import_errors table
	for i, valErr := range validationErrors {
		importError := &models.ImportError{
			JobID:        job.ID,
			LineNumber:   i + 1,
			FieldName:    valErr.Field,
			ErrorMessage: valErr.Message,
			RawData:      fmt.Sprintf("%v", valErr.Value),
		}
		_ = s.jobsRepo.AddError(ctx, importError)
	}

	// Import products
	successCount, importErrors := s.importProducts(ctx, job.ID, products, req.StorefrontID, req.UpdateMode)

	job.ProcessedRecords = len(products)
	job.SuccessfulRecords = successCount
	job.FailedRecords += len(importErrors)

	// Complete job
	completedTime := time.Now()
	job.CompletedAt = &completedTime
	job.Status = statusCompleted

	if len(importErrors) > 0 {
		errorMsg := fmt.Sprintf("Import completed with %d errors", len(importErrors))
		job.ErrorMessage = &errorMsg
	}

	// Update job in database
	if err := s.jobsRepo.Update(ctx, job); err != nil {
		return job, fmt.Errorf("failed to update job: %w", err)
	}

	return job, nil
}

// ImportFromFile imports products from uploaded file data
func (s *ImportService) ImportFromFile(ctx context.Context, userID int, fileData []byte, req models.ImportRequest) (*models.ImportJob, error) {
	// Create import job in database
	job := &models.ImportJob{
		StorefrontID: req.StorefrontID,
		UserID:       userID,
		FileType:     req.FileType,
		Status:       "pending",
	}

	if req.FileName != nil {
		job.FileName = *req.FileName
	}

	// Save job to database
	if err := s.jobsRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create import job: %w", err)
	}

	// Update job status to processing
	job.Status = statusProcessing
	startTime := time.Now()
	job.StartedAt = &startTime
	if err := s.jobsRepo.Update(ctx, job); err != nil {
		return job, fmt.Errorf("failed to update job status: %w", err)
	}

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
		job.Status = statusFailed
		errorMsg := fmt.Sprintf("unsupported file type: %s", req.FileType)
		job.ErrorMessage = &errorMsg
		_ = s.jobsRepo.Update(ctx, job)
		return job, errors.New(errorMsg)
	}

	if err != nil {
		job.Status = statusFailed
		errorMsg := err.Error()
		job.ErrorMessage = &errorMsg
		_ = s.jobsRepo.Update(ctx, job)
		return job, err
	}

	// Update job with totals
	job.TotalRecords = len(products) + len(validationErrors)
	job.FailedRecords = len(validationErrors)

	// Save validation errors to import_errors table
	for i, valErr := range validationErrors {
		importError := &models.ImportError{
			JobID:        job.ID,
			LineNumber:   i + 1,
			FieldName:    valErr.Field,
			ErrorMessage: valErr.Message,
			RawData:      fmt.Sprintf("%v", valErr.Value),
		}
		_ = s.jobsRepo.AddError(ctx, importError)
	}

	// Import products
	successCount, importErrors := s.importProducts(ctx, job.ID, products, req.StorefrontID, req.UpdateMode)

	job.ProcessedRecords = len(products)
	job.SuccessfulRecords = successCount
	job.FailedRecords += len(importErrors)

	// Complete job
	completedTime := time.Now()
	job.CompletedAt = &completedTime
	job.Status = statusCompleted

	if len(importErrors) > 0 {
		errorMsg := fmt.Sprintf("Import completed with %d errors", len(importErrors))
		job.ErrorMessage = &errorMsg
	}

	// Update job in database
	if err := s.jobsRepo.Update(ctx, job); err != nil {
		return job, fmt.Errorf("failed to update job: %w", err)
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
			_ = err // Explicitly ignore error
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
		return nil, nil, fmt.Errorf("digital Vision XML parsing failed: %w (generic XML parsing not implemented)", err)
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

// importProducts imports validated products using batch processing
func (s *ImportService) importProducts(ctx context.Context, jobID int, products []models.ImportProductRequest, storefrontID int, updateMode string) (int, []error) {
	var successCount int
	var allErrors []error

	const batchSize = 100

	// Process products in batches
	for batchStart := 0; batchStart < len(products); batchStart += batchSize {
		batchEnd := batchStart + batchSize
		if batchEnd > len(products) {
			batchEnd = len(products)
		}

		batch := products[batchStart:batchEnd]
		count, errors := s.importProductsBatch(ctx, jobID, batch, storefrontID, updateMode, batchStart)
		successCount += count
		allErrors = append(allErrors, errors...)
	}

	return successCount, allErrors
}

// importProductsBatch imports a batch of products
func (s *ImportService) importProductsBatch(ctx context.Context, jobID int, products []models.ImportProductRequest, storefrontID int, updateMode string, offset int) (int, []error) {
	var successCount int
	var errors []error

	// Collect all SKUs
	skus := make([]string, 0, len(products))
	productsBySKU := make(map[string]*models.ImportProductRequest)
	productsWithoutSKU := make([]*models.ImportProductRequest, 0)

	for i := range products {
		if products[i].SKU != "" {
			skus = append(skus, products[i].SKU)
			productsBySKU[products[i].SKU] = &products[i]
		} else {
			productsWithoutSKU = append(productsWithoutSKU, &products[i])
		}
	}

	// Get all existing products in one query
	existingProducts := make(map[string]*models.StorefrontProduct)
	if len(skus) > 0 {
		existing, err := s.productService.GetProductsBySKUs(ctx, storefrontID, skus)
		if err != nil {
			// Log error but continue
			fmt.Printf("Warning: failed to get existing products: %v\n", err)
		} else {
			existingProducts = existing
		}
	}

	// Separate products into new and existing
	var newProducts []*models.CreateProductRequest
	var updateProducts []struct {
		product *models.StorefrontProduct
		request *models.ImportProductRequest
	}

	// Process products with SKU
	for sku, importProduct := range productsBySKU {
		if existing, found := existingProducts[sku]; found {
			// Product exists
			switch updateMode {
			case "create_only":
				errors = append(errors, fmt.Errorf("product with SKU %s already exists", sku))
				lineNumber := offset + 1 // Will be corrected below
				_ = s.jobsRepo.AddError(ctx, &models.ImportError{
					JobID:        jobID,
					LineNumber:   lineNumber,
					FieldName:    sku,
					ErrorMessage: fmt.Sprintf("product with SKU %s already exists", sku),
					RawData:      importProduct.Name,
				})
			case updateModeUpdateOnly, updateModeUpsert:
				updateProducts = append(updateProducts, struct {
					product *models.StorefrontProduct
					request *models.ImportProductRequest
				}{existing, importProduct})
			}
		} else {
			// Product doesn't exist
			if updateMode == updateModeUpdateOnly {
				errors = append(errors, fmt.Errorf("product with SKU %s not found", sku))
				lineNumber := offset + 1
				_ = s.jobsRepo.AddError(ctx, &models.ImportError{
					JobID:        jobID,
					LineNumber:   lineNumber,
					FieldName:    sku,
					ErrorMessage: fmt.Sprintf("product with SKU %s not found", sku),
					RawData:      importProduct.Name,
				})
			} else {
				// Resolve category
				if err := s.resolveCategoryID(ctx, importProduct, storefrontID); err != nil {
					fmt.Printf("Warning: failed to resolve category for product %s: %v\n", importProduct.Name, err)
				}

				newProducts = append(newProducts, &models.CreateProductRequest{
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
				})
			}
		}
	}

	// Add products without SKU to new products list
	for _, importProduct := range productsWithoutSKU {
		if updateMode != "update_only" {
			// Resolve category
			if err := s.resolveCategoryID(ctx, importProduct, storefrontID); err != nil {
				fmt.Printf("Warning: failed to resolve category for product %s: %v\n", importProduct.Name, err)
			}

			newProducts = append(newProducts, &models.CreateProductRequest{
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
			})
		}
	}

	// Batch create new products
	if len(newProducts) > 0 {
		createdProducts, err := s.productService.BatchCreateProductsForImport(ctx, storefrontID, newProducts)
		if err != nil {
			fmt.Printf("Failed to batch create products: %v\n", err)
			errors = append(errors, err)
		} else {
			successCount += len(createdProducts)
			fmt.Printf("Successfully batch created %d products\n", len(createdProducts))

			// Import images for created products
			// TODO: This is sequential, could be optimized further
			for i, product := range createdProducts {
				if i < len(newProducts) {
					// Find corresponding import request with ImageURLs
					var imageURLs []string
					// Match by name or SKU
					if product.SKU != nil && *product.SKU != "" {
						for _, importProduct := range productsBySKU {
							if importProduct.SKU == *product.SKU {
								imageURLs = importProduct.ImageURLs
								break
							}
						}
					}
					if len(imageURLs) == 0 {
						for _, importProduct := range productsWithoutSKU {
							if importProduct.Name == product.Name {
								imageURLs = importProduct.ImageURLs
								break
							}
						}
					}

					if len(imageURLs) > 0 {
						_ = s.importProductImages(ctx, product.ID, imageURLs)
					}
				}
			}
		}
	}

	// Update existing products (still sequential, but usually fewer updates than creates)
	for _, item := range updateProducts {
		err := s.updateProduct(ctx, item.product.ID, *item.request, storefrontID)
		if err != nil {
			errors = append(errors, err)
			_ = s.jobsRepo.AddError(ctx, &models.ImportError{
				JobID:        jobID,
				LineNumber:   offset + 1,
				FieldName:    item.request.SKU,
				ErrorMessage: err.Error(),
				RawData:      item.request.Name,
			})
		} else {
			successCount++
		}
	}

	return successCount, errors
}

// importSingleProduct imports a single product
//
//nolint:unused // Legacy function, kept for potential future use
func (s *ImportService) importSingleProduct(ctx context.Context, importProduct models.ImportProductRequest, storefrontID int, updateMode string) error {
	// Check if product already exists by SKU or external ID
	var existingProduct *models.StorefrontProduct
	var err error

	if importProduct.SKU != "" {
		existingProduct, err = s.productService.GetProductBySKU(ctx, storefrontID, importProduct.SKU)
		// Если товар не найден - это нормально, просто создадим новый
		if err != nil && !errors.Is(err, postgres.ErrStorefrontProductNotFound) {
			return fmt.Errorf("failed to check existing product: %w", err)
		}
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
//
//nolint:unused // Legacy function, kept for potential future use
func (s *ImportService) createProduct(ctx context.Context, importProduct models.ImportProductRequest, storefrontID int) error {
	// Resolve category if original_category is present in attributes
	if err := s.resolveCategoryID(ctx, &importProduct, storefrontID); err != nil {
		// Log warning but continue with default category
		fmt.Printf("Warning: failed to resolve category for product %s: %v\n", importProduct.Name, err)
	}

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

	// Import images if provided
	if len(importProduct.ImageURLs) > 0 {
		err := s.importProductImages(ctx, product.ID, importProduct.ImageURLs)
		if err != nil {
			// Логируем ошибку, но не прерываем импорт товара
			// Товар уже создан, и частичная загрузка изображений допустима
			fmt.Printf("Failed to import all images for product %d: %v\n", product.ID, err)
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

// PreviewImport previews first N rows of import file with validation
func (s *ImportService) PreviewImport(ctx context.Context, fileData []byte, fileType string, storefrontID int, previewLimit int) (*models.ImportPreviewResponse, error) {
	if previewLimit <= 0 {
		previewLimit = 10 // Default preview limit
	}

	switch fileType {
	case fileTypeCSV:
		return s.previewCSVData(fileData, storefrontID, previewLimit)
	case fileTypeXML:
		return s.previewXMLData(fileData, storefrontID, previewLimit)
	case fileTypeZIP:
		return nil, fmt.Errorf("preview not supported for ZIP files")
	default:
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}
}

// previewCSVData previews CSV file data
func (s *ImportService) previewCSVData(csvData []byte, storefrontID int, previewLimit int) (*models.ImportPreviewResponse, error) {
	// Parse CSV to get products and validation errors
	reader := bytes.NewReader(csvData)
	parser := parsers.NewCSVParser(storefrontID)
	products, validationErrors, err := parser.ParseCSV(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	response := &models.ImportPreviewResponse{
		FileType:     fileTypeCSV,
		PreviewRows:  make([]models.ImportPreviewRow, 0),
		TotalRows:    len(products) + len(validationErrors),
		ValidationOK: len(validationErrors) == 0,
	}

	// Convert first N products to preview rows
	for i, product := range products {
		if i >= previewLimit {
			break
		}

		rowData := map[string]interface{}{
			"sku":            product.SKU,
			"name":           product.Name,
			"price":          product.Price,
			"currency":       product.Currency,
			"category_id":    product.CategoryID,
			"description":    product.Description,
			"stock_quantity": product.StockQuantity,
			"barcode":        product.Barcode,
			"is_active":      product.IsActive,
		}

		previewRow := models.ImportPreviewRow{
			LineNumber: i + 2, // +1 for 0-index, +1 for header
			Data:       rowData,
			Errors:     []models.ImportValidationError{},
			IsValid:    true,
		}

		response.PreviewRows = append(response.PreviewRows, previewRow)
	}

	// Add validation errors to response (if any)
	if len(validationErrors) > 0 && len(response.PreviewRows) < previewLimit {
		// Create error rows for first few errors
		for i, validErr := range validationErrors {
			if len(response.PreviewRows) >= previewLimit {
				break
			}

			errorRow := models.ImportPreviewRow{
				LineNumber: len(products) + i + 2,
				Data:       map[string]interface{}{"error": "Validation failed"},
				Errors:     []models.ImportValidationError{validErr},
				IsValid:    false,
			}

			response.PreviewRows = append(response.PreviewRows, errorRow)
			response.ValidationOK = false
		}
	}

	if !response.ValidationOK {
		response.ErrorSummary = fmt.Sprintf("Found %d validation errors in file", len(validationErrors))
	}

	return response, nil
}

// previewXMLData previews XML file data
func (s *ImportService) previewXMLData(xmlData []byte, storefrontID int, previewLimit int) (*models.ImportPreviewResponse, error) {
	// Parse XML using existing parser
	products, validationErrors, err := s.processXMLData(xmlData, storefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	response := &models.ImportPreviewResponse{
		FileType:     fileTypeXML,
		PreviewRows:  make([]models.ImportPreviewRow, 0),
		TotalRows:    len(products) + len(validationErrors),
		ValidationOK: true,
	}

	// Convert products to preview rows (limit to previewLimit)
	for i, product := range products {
		if i >= previewLimit {
			break
		}

		rowData := map[string]interface{}{
			"sku":         product.SKU,
			"name":        product.Name,
			"price":       product.Price,
			"currency":    product.Currency,
			"category_id": product.CategoryID,
			"description": product.Description,
		}

		previewRow := models.ImportPreviewRow{
			LineNumber: i + 1,
			Data:       rowData,
			Errors:     []models.ImportValidationError{},
			IsValid:    true,
		}

		response.PreviewRows = append(response.PreviewRows, previewRow)
	}

	// Add validation errors to preview rows
	errorMap := make(map[int][]models.ImportValidationError)
	if len(validationErrors) > 0 {
		errorMap[0] = append(errorMap[0], validationErrors...)
	}

	for lineNum, errors := range errorMap {
		if lineNum >= previewLimit {
			continue
		}

		// Find or create preview row for this line
		found := false
		for i := range response.PreviewRows {
			if response.PreviewRows[i].LineNumber == lineNum+1 {
				response.PreviewRows[i].Errors = errors
				response.PreviewRows[i].IsValid = false
				response.ValidationOK = false
				found = true
				break
			}
		}

		if !found && len(response.PreviewRows) < previewLimit {
			previewRow := models.ImportPreviewRow{
				LineNumber: lineNum + 1,
				Data:       map[string]interface{}{},
				Errors:     errors,
				IsValid:    false,
			}
			response.PreviewRows = append(response.PreviewRows, previewRow)
			response.ValidationOK = false
		}
	}

	if !response.ValidationOK {
		response.ErrorSummary = "Found validation errors in XML data"
	}

	return response, nil
}

// countInvalidRows counts rows with validation errors
// nolint:unused // Used in preview logic
func countInvalidRows(rows []models.ImportPreviewRow) int {
	count := 0
	for _, row := range rows {
		if !row.IsValid {
			count++
		}
	}
	return count
}

// GetJobs returns list of import jobs for a storefront
func (s *ImportService) GetJobs(ctx context.Context, storefrontID int, status string, limit, offset int) (*models.ImportJobsResponse, error) {
	filter := &postgres.ImportJobFilter{
		Limit:     limit,
		Offset:    offset,
		SortBy:    "created_at",
		SortOrder: "DESC",
	}

	if status != "" {
		filter.Status = &status
	}

	jobs, total, err := s.jobsRepo.GetByStorefront(ctx, storefrontID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get import jobs: %w", err)
	}

	// Convert []*models.ImportJob to []models.ImportJob
	jobsList := make([]models.ImportJob, len(jobs))
	for i, job := range jobs {
		jobsList[i] = *job
	}

	return &models.ImportJobsResponse{
		Jobs:  jobsList,
		Total: total,
	}, nil
}

// GetJobDetails returns detailed information about an import job
func (s *ImportService) GetJobDetails(ctx context.Context, jobID int) (*models.ImportJob, error) {
	job, err := s.jobsRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job details: %w", err)
	}

	return job, nil
}

// GetJobStatus returns status of an import job
func (s *ImportService) GetJobStatus(ctx context.Context, jobID int) (*models.ImportJobStatus, error) {
	job, err := s.jobsRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job status: %w", err)
	}

	// Get errors for this job
	errors, err := s.jobsRepo.GetErrors(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get job errors: %w", err)
	}

	// Calculate progress
	var progress float64
	if job.TotalRecords > 0 {
		progress = float64(job.ProcessedRecords) / float64(job.TotalRecords) * 100
	}

	status := &models.ImportJobStatus{
		ID:                job.ID,
		Status:            job.Status,
		Progress:          progress,
		TotalRecords:      job.TotalRecords,
		ProcessedRecords:  job.ProcessedRecords,
		SuccessfulRecords: job.SuccessfulRecords,
		FailedRecords:     job.FailedRecords,
		Errors:            convertToImportErrors(errors),
		StartedAt:         job.StartedAt,
		CompletedAt:       job.CompletedAt,
	}

	return status, nil
}

// convertToImportErrors converts []*models.ImportError to []models.ImportError
func convertToImportErrors(errors []*models.ImportError) []models.ImportError {
	result := make([]models.ImportError, len(errors))
	for i, err := range errors {
		result[i] = *err
	}
	return result
}

// CancelJob cancels a running import job
func (s *ImportService) CancelJob(ctx context.Context, jobID int) error {
	// Check if job exists and is in cancellable state
	job, err := s.jobsRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	// Only pending or processing jobs can be canceled
	if job.Status != statusPending && job.Status != statusProcessing {
		return fmt.Errorf("job cannot be canceled in %s state", job.Status)
	}

	// Update job status to canceled
	return s.jobsRepo.UpdateStatus(ctx, jobID, statusCanceled)
}

// RetryJob retries a failed import job
func (s *ImportService) RetryJob(ctx context.Context, jobID int) (*models.ImportJob, error) {
	// Get the original job
	originalJob, err := s.jobsRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("failed to get original job: %w", err)
	}

	// Only failed jobs can be retried
	if originalJob.Status != statusFailed && originalJob.Status != statusCanceled {
		return nil, fmt.Errorf("only failed or canceled jobs can be retried")
	}

	// Create a new job based on the original
	newJob := &models.ImportJob{
		StorefrontID: originalJob.StorefrontID,
		UserID:       originalJob.UserID,
		FileName:     originalJob.FileName,
		FileType:     originalJob.FileType,
		FileURL:      originalJob.FileURL,
		Status:       "pending",
	}

	// Save new job to database
	if err := s.jobsRepo.Create(ctx, newJob); err != nil {
		return nil, fmt.Errorf("failed to create retry job: %w", err)
	}

	return newJob, nil
}

// ImportFromURLAsync asynchronously downloads and imports products from a URL
func (s *ImportService) ImportFromURLAsync(ctx context.Context, userID int, req models.ImportRequest) (*models.ImportJob, error) {
	if req.FileURL == nil {
		return nil, fmt.Errorf("file URL is required")
	}

	// Check if queue manager is available
	if s.queueManager == nil || !s.queueManager.IsRunning() {
		// Fallback to synchronous import if queue manager is not available
		return s.ImportFromURL(ctx, userID, req)
	}

	// Create import job in database with pending status
	job := &models.ImportJob{
		StorefrontID: req.StorefrontID,
		UserID:       userID,
		FileType:     req.FileType,
		FileURL:      req.FileURL,
		Status:       "pending",
	}

	// Save job to database
	if err := s.jobsRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create import job: %w", err)
	}

	// Download file in background
	go func(parentCtx context.Context) {
		// Create a new context with timeout for background processing (detached from parent)
		bgCtx, cancel := context.WithTimeout(parentCtx, 30*time.Minute)
		defer cancel()

		data, _, err := s.downloadFile(bgCtx, *req.FileURL)
		if err != nil {
			// Update job status to failed
			job.Status = statusFailed
			errorMsg := fmt.Sprintf("Failed to download file: %v", err)
			job.ErrorMessage = &errorMsg
			_ = s.jobsRepo.Update(bgCtx, job)
			return
		}

		// Create task for worker pool
		task := &ImportJobTask{
			JobID:        job.ID,
			UserID:       userID,
			StorefrontID: req.StorefrontID,
			FileData:     data,
			FileType:     req.FileType,
			UpdateMode:   req.UpdateMode,
			FileURL:      req.FileURL,
		}

		// Enqueue job
		if err := s.queueManager.EnqueueJob(task); err != nil {
			// Update job status to failed if enqueue fails
			job.Status = statusFailed
			errorMsg := fmt.Sprintf("Failed to enqueue job: %v", err)
			job.ErrorMessage = &errorMsg
			_ = s.jobsRepo.Update(bgCtx, job)
		}
	}(ctx)

	// Return job immediately
	return job, nil
}

// ImportFromFileAsync asynchronously imports products from uploaded file data
func (s *ImportService) ImportFromFileAsync(ctx context.Context, userID int, fileData []byte, req models.ImportRequest) (*models.ImportJob, error) {
	// Check if queue manager is available
	if s.queueManager == nil || !s.queueManager.IsRunning() {
		// Fallback to synchronous import if queue manager is not available
		return s.ImportFromFile(ctx, userID, fileData, req)
	}

	// Create import job in database with pending status
	job := &models.ImportJob{
		StorefrontID: req.StorefrontID,
		UserID:       userID,
		FileType:     req.FileType,
		Status:       "pending",
	}

	if req.FileName != nil {
		job.FileName = *req.FileName
	}

	// Save job to database
	if err := s.jobsRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create import job: %w", err)
	}

	// Create task for worker pool
	task := &ImportJobTask{
		JobID:        job.ID,
		UserID:       userID,
		StorefrontID: req.StorefrontID,
		FileData:     fileData,
		FileType:     req.FileType,
		UpdateMode:   req.UpdateMode,
		FileName:     req.FileName,
	}

	// Enqueue job
	if err := s.queueManager.EnqueueJob(task); err != nil {
		// Update job status to failed if enqueue fails
		job.Status = statusFailed
		errorMsg := fmt.Sprintf("Failed to enqueue job: %v", err)
		job.ErrorMessage = &errorMsg
		_ = s.jobsRepo.Update(ctx, job)
		return job, fmt.Errorf("failed to enqueue job: %w", err)
	}

	// Return job immediately
	return job, nil
}

// downloadImage downloads an image from a URL
func (s *ImportService) downloadImage(ctx context.Context, url string) ([]byte, string, error) {
	// Валидация URL
	if url == "" {
		return nil, "", fmt.Errorf("empty image URL")
	}

	// Создание HTTP запроса
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	// HTTP клиент с таймаутом
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Выполнение запроса
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download image: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Failed to close response body: %v\n", err)
		}
	}()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Проверка Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, "", fmt.Errorf("invalid content type: %s (expected image/*)", contentType)
	}

	// Проверка размера файла (максимум 10MB)
	const maxImageSize = 10 * 1024 * 1024
	if resp.ContentLength > maxImageSize {
		return nil, "", fmt.Errorf("image size exceeds maximum limit of 10MB")
	}

	// Чтение данных изображения
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Дополнительная проверка размера после загрузки
	if len(imageData) > maxImageSize {
		return nil, "", fmt.Errorf("image size exceeds maximum limit of 10MB")
	}

	// Валидация что это действительно изображение (попытка декодировать)
	_, _, err = image.DecodeConfig(bytes.NewReader(imageData))
	if err != nil {
		return nil, "", fmt.Errorf("invalid image format: %w", err)
	}

	// Определение расширения файла
	ext := filepath.Ext(url)
	if ext == "" || len(ext) > 5 {
		// Если расширение не найдено в URL, определяем по Content-Type
		switch contentType {
		case "image/jpeg", "image/jpg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg" // Fallback
		}
	}

	return imageData, ext, nil
}

// bytesFileAdapter wraps bytes.Reader to implement multipart.File interface
type bytesFileAdapter struct {
	*bytes.Reader
}

// Close implements multipart.File
func (b *bytesFileAdapter) Close() error {
	return nil // bytes.Reader doesn't need closing
}

// importProductImages imports images for a product from URLs
func (s *ImportService) importProductImages(ctx context.Context, productID int, imageURLs []string) error {
	if s.imageService == nil {
		return fmt.Errorf("image service not initialized")
	}

	if len(imageURLs) == 0 {
		return nil // Нет изображений для загрузки
	}

	// Загружаем каждое изображение
	for i, imageURL := range imageURLs {
		// Скачиваем изображение
		imageData, ext, err := s.downloadImage(ctx, imageURL)
		if err != nil {
			// Логируем ошибку, но продолжаем загрузку остальных изображений
			fmt.Printf("Failed to download image from %s: %v\n", imageURL, err)
			continue
		}

		// Создаем multipart.File адаптер
		fileReader := &bytesFileAdapter{Reader: bytes.NewReader(imageData)}

		// Создаем multipart.FileHeader
		filename := fmt.Sprintf("import_image_%d%s", i+1, ext)
		fileHeader := &multipart.FileHeader{
			Filename: filename,
			Size:     int64(len(imageData)),
			Header:   make(map[string][]string),
		}
		// Определяем Content-Type на основе расширения
		contentType := "image/jpeg"
		switch ext {
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		case ".webp":
			contentType = "image/webp"
		}
		fileHeader.Header["Content-Type"] = []string{contentType}

		// Создаем запрос для загрузки изображения
		uploadRequest := &services.UploadImageRequest{
			EntityType:   services.ImageTypeStorefrontProduct,
			EntityID:     productID,
			File:         fileReader,
			FileHeader:   fileHeader,
			IsMain:       i == 0, // Первое изображение делаем главным
			DisplayOrder: i + 1,
		}

		// Загружаем изображение через ImageService
		_, err = s.imageService.UploadImage(ctx, uploadRequest)
		if err != nil {
			// Логируем ошибку, но продолжаем загрузку остальных изображений
			fmt.Printf("Failed to upload image %s for product %d: %v\n", imageURL, productID, err)
			continue
		}

		fmt.Printf("Successfully imported image %d/%d for product %d from %s\n", i+1, len(imageURLs), productID, imageURL)
	}

	return nil
}

// resolveCategoryID resolves category ID from original_category attribute using CategoryMappingService
func (s *ImportService) resolveCategoryID(ctx context.Context, importProduct *models.ImportProductRequest, storefrontID int) error {
	// Skip if category mapping service is not available
	if s.categoryMappingService == nil {
		return nil
	}

	// Check if original_category is present in attributes
	if importProduct.Attributes == nil {
		return nil
	}

	originalCategory, exists := importProduct.Attributes["original_category"]
	if !exists {
		return nil
	}

	// Convert to string
	categoryPath, ok := originalCategory.(string)
	if !ok || categoryPath == "" {
		return nil
	}

	// Use CategoryMappingService to resolve category
	categoryID, err := s.categoryMappingService.GetOrCreateMapping(
		ctx,
		storefrontID,
		categoryPath,
		importProduct.Name,
		importProduct.Description,
	)
	if err != nil {
		// Log warning but use fallback category (1001 - Elektronika)
		// This ensures import doesn't fail due to category resolution issues
		fmt.Printf("Warning: failed to resolve category '%s' for product %s: %v. Using fallback category 1001.\n",
			categoryPath, importProduct.Name, err)
		importProduct.CategoryID = 1001
		return nil
	}

	// Update category ID
	importProduct.CategoryID = categoryID

	return nil
}
