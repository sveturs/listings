package service

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
	marketplaceOpenSearch "backend/internal/proj/marketplace/storage/opensearch"
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
	productService          *ProductService
	jobsRepo                postgres.ImportJobsRepositoryInterface
	queueManager            *ImportQueueManager
	imageService            *services.ImageService
	categoryMappingService  *CategoryMappingService
	variantDetector         *VariantDetector
	attributeMapper         *AttributeMapper
	dbStorage               ImportStorage                                     // Database storage interface для доступа к БД
	marketplaceSearchRepo   marketplaceOpenSearch.MarketplaceSearchRepository // OpenSearch repository для marketplace listings
}

// ImportStorage interface for import service (минимальный интерфейс для избежания конфликтов)
type ImportStorage interface {
	GetMarketplaceListingsForReindex(ctx context.Context, limit int) ([]*models.MarketplaceListing, error)
	ResetMarketplaceListingsReindexFlag(ctx context.Context, listingIDs []int) error
}

// NewImportService creates a new import service
func NewImportService(
	productService *ProductService,
	jobsRepo postgres.ImportJobsRepositoryInterface,
	imageService *services.ImageService,
	categoryMappingService *CategoryMappingService,
	attributeMapper *AttributeMapper,
	dbStorage ImportStorage,
	marketplaceSearchRepo marketplaceOpenSearch.MarketplaceSearchRepository,
) *ImportService {
	return &ImportService{
		productService:         productService,
		jobsRepo:               jobsRepo,
		queueManager:           nil, // Will be set via SetQueueManager
		imageService:           imageService,
		categoryMappingService: categoryMappingService,
		variantDetector:        NewVariantDetector(),
		attributeMapper:        attributeMapper,
		dbStorage:              dbStorage,
		marketplaceSearchRepo:  marketplaceSearchRepo,
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
			case updateModeCreateOnly:
				errors = append(errors, fmt.Errorf("product with SKU %s already exists", sku))
				lineNumber := offset + 1 // Will be corrected below
				_ = s.jobsRepo.AddError(ctx, &models.ImportError{
					JobID:        jobID,
					LineNumber:   lineNumber,
					FieldName:    "sku", // Field that caused the error
					ErrorMessage: fmt.Sprintf("Product with SKU '%s' already exists (mode: create_only). Use 'update_only' or 'upsert' mode to update existing products.", sku),
					RawData:      fmt.Sprintf("name=%s, sku=%s", importProduct.Name, sku),
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
					FieldName:    "sku", // Field that caused the error
					ErrorMessage: fmt.Sprintf("Product with SKU '%s' not found (mode: update_only). Use 'create_only' or 'upsert' mode to create new products.", sku),
					RawData:      fmt.Sprintf("name=%s, sku=%s", importProduct.Name, sku),
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
				FieldName:    "product", // General product update error
				ErrorMessage: fmt.Sprintf("Failed to update product with SKU '%s': %v", item.request.SKU, err),
				RawData:      fmt.Sprintf("name=%s, sku=%s", item.request.Name, item.request.SKU),
			})
		} else {
			successCount++

			// Import images for updated products
			if len(item.request.ImageURLs) > 0 {
				_ = s.importProductImages(ctx, item.product.ID, item.request.ImageURLs)
			}
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

// downloadImage downloads an image from a URL with retry logic for TLS issues
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

	// Попытка 1: Стандартный HTTP клиент с поддержкой устаревших TLS серверов
	// ВАЖНО: InsecureSkipVerify используется только для загрузки внешних изображений
	// из доверенных источников (прайсов поставщиков). Для production рекомендуется
	// настроить прокси или обновить TLS конфигурацию поставщиков.
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,     // Игнорируем проблемы сертификатов
			MinVersion:         tls.VersionTLS10, // Разрешаем старые версии TLS
			MaxVersion:         tls.VersionTLS13, // До новейших
		},
	}
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	// Выполнение запроса
	resp, err := client.Do(req)
	if err != nil {
		// Попытка 2: Fallback на системный curl для обхода Go TLS ограничений
		// Это необходимо для серверов с устаревшими DH ключами (например, digitalvision.rs)
		if strings.Contains(err.Error(), "tls") || strings.Contains(err.Error(), "handshake") {
			logger.Warn().
				Str("url", url).
				Err(err).
				Msg("TLS handshake failed, trying fallback with system curl")

			// Используем системный curl с --insecure для обхода TLS проблем
			imageData, ext, curlErr := s.downloadImageWithCurl(ctx, url)
			if curlErr != nil {
				logger.Error().
					Str("url", url).
					Err(curlErr).
					Msg("Curl fallback also failed")
				return nil, "", fmt.Errorf("failed to download image (both Go and curl failed): %w", err)
			}

			logger.Info().
				Str("url", url).
				Msg("Successfully downloaded image using curl fallback")
			return imageData, ext, nil
		}
		return nil, "", fmt.Errorf("failed to download image: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close response body")
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

// downloadImageWithCurl downloads image using system curl as fallback for TLS issues
// This is necessary for servers with weak DH keys that Go's TLS client cannot handle
func (s *ImportService) downloadImageWithCurl(ctx context.Context, url string) ([]byte, string, error) {
	// Создаем временный файл для сохранения изображения
	tmpFile, err := os.CreateTemp("", "import_image_*")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath) // Удаляем временный файл после использования

	// Выполняем curl с пониженным security level для обхода weak DH key проблемы
	// ВАЖНО: Это workaround для серверов с устаревшей TLS конфигурацией
	// DEFAULT:@SECLEVEL=0 - понижает OpenSSL security level до 0 (разрешает weak DH keys)
	// --insecure: игнорирует проблемы SSL сертификатов
	// --max-time 30: таймаут 30 секунд
	// --location: следует редиректам
	// --fail: возвращает ошибку при HTTP ошибках (4xx, 5xx)
	cmd := exec.CommandContext(ctx, "curl",
		"--ciphers", "DEFAULT:@SECLEVEL=0",
		"--insecure",
		"--silent",
		"--show-error",
		"--max-time", "30",
		"--location",
		"--fail",
		"--output", tmpPath,
		url,
	)
	// Устанавливаем пустой OPENSSL_CONF для избежания конфликтов с системной конфигурацией
	cmd.Env = append(os.Environ(), "OPENSSL_CONF=/dev/null")

	// Выполняем команду
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, "", fmt.Errorf("curl command failed: %w (output: %s)", err, string(output))
	}

	// Читаем загруженный файл
	imageData, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read downloaded file: %w", err)
	}

	// Проверка размера файла (максимум 10MB)
	const maxImageSize = 10 * 1024 * 1024
	if len(imageData) > maxImageSize {
		return nil, "", fmt.Errorf("image size exceeds maximum limit of 10MB")
	}

	// Валидация что это действительно изображение
	config, format, err := image.DecodeConfig(bytes.NewReader(imageData))
	if err != nil {
		return nil, "", fmt.Errorf("invalid image format: %w", err)
	}

	// Определение расширения на основе формата
	ext := "." + format
	if ext == ".jpeg" {
		ext = ".jpg"
	}

	logger.Info().
		Str("url", url).
		Str("format", format).
		Int("width", config.Width).
		Int("height", config.Height).
		Int("size_bytes", len(imageData)).
		Msg("Successfully downloaded image using curl")

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

// convertImportProductsToVariants конвертирует ImportProductRequest в ProductVariant
func (s *ImportService) convertImportProductsToVariants(products []models.ImportProductRequest) []*ProductVariant {
	variants := make([]*ProductVariant, 0, len(products))

	for _, p := range products {
		variant := &ProductVariant{
			Name:          p.Name,
			SKU:           p.SKU,
			Price:         p.Price,
			StockQuantity: p.StockQuantity,
		}

		// Добавляем изображение если есть
		if len(p.ImageURLs) > 0 {
			variant.ImageURL = p.ImageURLs[0]
		}

		// Сохраняем оригинальные атрибуты
		variant.OriginalAttributes = make(map[string]interface{})
		if p.Description != "" {
			variant.OriginalAttributes["description"] = p.Description
		}
		if p.Barcode != "" {
			variant.OriginalAttributes["barcode"] = p.Barcode
		}
		variant.OriginalAttributes["category_id"] = p.CategoryID

		variants = append(variants, variant)
	}

	return variants
}

// groupAndDetectVariants группирует товары в варианты используя VariantDetector
func (s *ImportService) groupAndDetectVariants(products []models.ImportProductRequest) []*VariantGroup {
	// Конвертируем в ProductVariant
	variants := s.convertImportProductsToVariants(products)

	// Группируем через detector
	groups := s.variantDetector.GroupProducts(variants)

	return groups
}

// importVariantGroup импортирует группу вариантов как один товар с вариантами
func (s *ImportService) importVariantGroup(
	ctx context.Context,
	group *VariantGroup,
	storefrontID int,
) error {
	if len(group.Variants) == 0 {
		return fmt.Errorf("variant group has no variants")
	}

	// 1. Создаем parent product из первого варианта группы
	firstVariant := group.Variants[0]

	// Извлекаем данные из OriginalAttributes
	description := ""
	if desc, ok := firstVariant.OriginalAttributes["description"].(string); ok {
		description = desc
	}

	categoryID := 0
	if catID, ok := firstVariant.OriginalAttributes["category_id"].(int); ok {
		categoryID = catID
	}

	// Используем базовое название группы для parent product
	createReq := models.CreateProductRequest{
		Name:          group.BaseName,
		Description:   description,
		Price:         firstVariant.Price, // цена по умолчанию из первого варианта
		Currency:      "RSD",              // TODO: извлекать из OriginalAttributes
		CategoryID:    categoryID,
		SKU:           &group.BaseName, // используем base name как SKU родителя
		StockQuantity: 0,               // суммарное количество будет в вариантах
		IsActive:      true,
		Attributes:    make(map[string]interface{}),
	}

	// Создаем parent product
	product, err := s.productService.CreateProductForImport(ctx, storefrontID, &createReq)
	if err != nil {
		return fmt.Errorf("failed to create parent product for variant group %s: %w", group.BaseName, err)
	}

	fmt.Printf("Created parent product for variant group: %s (ID: %d, %d variants)\n",
		group.BaseName, product.ID, len(group.Variants))

	// 2. Создаем варианты товара
	variantRequests := make([]*models.CreateProductVariantRequest, 0, len(group.Variants))

	for i, variant := range group.Variants {
		// Конвертируем VariantAttributes в JSONB
		variantAttrsJSON := models.JSONB{}
		for key, val := range variant.VariantAttributes {
			variantAttrsJSON[key] = val
		}

		// Извлекаем barcode если есть
		var barcode *string
		if bc, ok := variant.OriginalAttributes["barcode"].(string); ok && bc != "" {
			barcode = &bc
		}

		// Статус склада
		stockStatus := "in_stock"
		if variant.StockQuantity == 0 {
			stockStatus = "out_of_stock"
		}

		variantReq := &models.CreateProductVariantRequest{
			ProductID:         product.ID,
			SKU:               &variant.SKU,
			Barcode:           barcode,
			Price:             &variant.Price,
			CompareAtPrice:    nil, // TODO: если есть sale price
			CostPrice:         nil,
			StockQuantity:     variant.StockQuantity,
			StockStatus:       stockStatus,
			LowStockThreshold: nil,
			VariantAttributes: variantAttrsJSON,
			Weight:            nil,
			Dimensions:        models.JSONB{},
			IsActive:          true,
			IsDefault:         i == 0, // первый вариант делаем default
		}

		variantRequests = append(variantRequests, variantReq)
	}

	// Batch создание вариантов
	createdVariants, err := s.productService.storage.BatchCreateProductVariants(ctx, variantRequests)
	if err != nil {
		return fmt.Errorf("failed to create variants for product %d: %w", product.ID, err)
	}

	fmt.Printf("Created %d variants for product %s (ID: %d)\n",
		len(createdVariants), product.Name, product.ID)

	// 3. Добавляем изображения к вариантам (если есть)
	imageRequests := make([]*models.CreateProductVariantImageRequest, 0)

	for i, variant := range group.Variants {
		if variant.ImageURL != "" && i < len(createdVariants) {
			createdVariantID := createdVariants[i].ID

			imageReq := &models.CreateProductVariantImageRequest{
				VariantID:    createdVariantID,
				ImageURL:     variant.ImageURL,
				ThumbnailURL: nil,
				AltText:      nil,
				DisplayOrder: 0,    // первое изображение
				IsMain:       true, // первое изображение делаем главным
			}
			imageRequests = append(imageRequests, imageReq)
		}
	}

	if len(imageRequests) > 0 {
		_, err = s.productService.storage.BatchCreateProductVariantImages(ctx, imageRequests)
		if err != nil {
			// Логируем ошибку, но не прерываем импорт
			fmt.Printf("Warning: failed to import some variant images for product %d: %v\n", product.ID, err)
		} else {
			fmt.Printf("Imported %d variant images for product %s\n", len(imageRequests), product.Name)
		}
	}

	return nil
}

// IndexPendingProducts индексирует товары витрины в OpenSearch после импорта
// Вызывается после успешного импорта для инкрементальной индексации
func (s *ImportService) IndexPendingProducts(ctx context.Context, storefrontID int, batchSize int) error {
	if batchSize <= 0 {
		batchSize = 100
	}

	// Получаем все активные товары витрины (они уже помечены needs_reindex=true триггером)
	filter := models.ProductFilter{
		StorefrontID: storefrontID,
		IsActive:     boolPtr(true),
		Limit:        1000, // большой лимит для всех товаров
		Offset:       0,
	}

	products, err := s.productService.storage.GetStorefrontProducts(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to fetch storefront products: %w", err)
	}

	if len(products) == 0 {
		return nil // нет товаров для индексации
	}

	// Индексируем товары батчами
	for i := 0; i < len(products); i += batchSize {
		end := i + batchSize
		if end > len(products) {
			end = len(products)
		}

		batch := products[i:end]
		if err := s.indexProductsBatch(ctx, batch); err != nil {
			return fmt.Errorf("failed to index batch %d-%d: %w", i, end, err)
		}
	}

	return nil
}

// indexProductsBatch индексирует batch товаров в OpenSearch
func (s *ImportService) indexProductsBatch(ctx context.Context, products []*models.StorefrontProduct) error {
	if len(products) == 0 {
		return nil
	}

	// Индексируем в OpenSearch через searchRepo
	searchRepo := s.productService.searchRepo
	if searchRepo != nil {
		if err := searchRepo.BulkIndexProducts(ctx, products); err != nil {
			return fmt.Errorf("failed to bulk index products: %w", err)
		}
		fmt.Printf("Successfully indexed %d products in OpenSearch\n", len(products))
	}

	return nil
}

// IndexMarketplaceListings индексирует marketplace listings с needs_reindex=true в OpenSearch
// Вызывается после успешного импорта для инкрементальной индексации marketplace_listings
func (s *ImportService) IndexMarketplaceListings(ctx context.Context, batchSize int) error {
	if batchSize <= 0 {
		batchSize = 100
	}

	// Получаем listings с needs_reindex=true из БД
	listings, err := s.dbStorage.GetMarketplaceListingsForReindex(ctx, 1000)
	if err != nil {
		return fmt.Errorf("failed to fetch marketplace listings for reindex: %w", err)
	}

	if len(listings) == 0 {
		logger.Debug().Msg("No marketplace listings need reindexing")
		return nil
	}

	logger.Info().
		Int("count", len(listings)).
		Msg("Indexing marketplace listings in OpenSearch")

	// Индексируем listings батчами
	for i := 0; i < len(listings); i += batchSize {
		end := i + batchSize
		if end > len(listings) {
			end = len(listings)
		}

		batch := listings[i:end]
		if err := s.indexMarketplaceListingsBatch(ctx, batch); err != nil {
			return fmt.Errorf("failed to index marketplace listings batch %d-%d: %w", i, end, err)
		}
	}

	// Сбрасываем флаг needs_reindex после успешной индексации
	listingIDs := make([]int, len(listings))
	for i, listing := range listings {
		listingIDs[i] = listing.ID
	}

	if err := s.dbStorage.ResetMarketplaceListingsReindexFlag(ctx, listingIDs); err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to reset needs_reindex flag for marketplace listings")
		// Не возвращаем ошибку, так как индексация прошла успешно
	} else {
		logger.Info().
			Int("count", len(listingIDs)).
			Msg("Reset needs_reindex flag for marketplace listings")
	}

	return nil
}

// indexMarketplaceListingsBatch индексирует batch marketplace listings в OpenSearch
func (s *ImportService) indexMarketplaceListingsBatch(ctx context.Context, listings []*models.MarketplaceListing) error {
	if len(listings) == 0 {
		return nil
	}

	// Индексируем в OpenSearch через marketplace search repo
	if s.marketplaceSearchRepo != nil {
		if err := s.marketplaceSearchRepo.BulkIndexListings(ctx, listings); err != nil {
			return fmt.Errorf("failed to bulk index marketplace listings: %w", err)
		}

		logger.Info().
			Int("count", len(listings)).
			Msg("Successfully indexed marketplace listings in OpenSearch")
	} else {
		logger.Warn().Msg("Marketplace search repository is nil, skipping indexing")
	}

	return nil
}

// boolPtr helper для создания *bool
func boolPtr(b bool) *bool {
	return &b
}
