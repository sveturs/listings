// backend/internal/proj/storefront/service/storefront.go
package service

import (
	"archive/zip"
	"backend/internal/domain/models"
	"backend/internal/storage"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"encoding/xml"
	"github.com/microcosm-cc/bluemonday"


)

const (
	StorefrontCreationCost = 15000.0 // стоимость создания витрины
)

type StorefrontService struct {
	storage storage.Storage
}

func NewStorefrontService(storage storage.Storage) StorefrontServiceInterface {
	return &StorefrontService{
		storage: storage,
	}
}

// CreateStorefront создает новую витрину с проверкой баланса
func (s *StorefrontService) CreateStorefront(ctx context.Context, userID int, create *models.StorefrontCreate) (*models.Storefront, error) {
	// Получаем баланс пользователя
	balance, err := s.storage.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user balance: %w", err)
	}

	// Проверяем, хватает ли средств
	if balance.Balance < StorefrontCreationCost {
		return nil, fmt.Errorf("insufficient funds: required %.2f, available %.2f", StorefrontCreationCost, balance.Balance)
	}

	// Начинаем транзакцию
	tx, err := s.storage.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Создаем транзакцию списания средств
	now := time.Now()
	transaction := &models.BalanceTransaction{
		UserID:        userID,
		Type:          "service_payment",
		Amount:        StorefrontCreationCost,
		Currency:      "RSD",
		Status:        "completed",
		PaymentMethod: "balance",
		Description:   "Создание витрины магазина",
		CreatedAt:     now,
		CompletedAt:   &now,
	}

	transactionID, err := s.storage.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Обновляем баланс пользователя
	err = s.storage.UpdateBalance(ctx, userID, -StorefrontCreationCost)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Генерируем уникальный slug
	slug := generateSlug(create.Name)

	// Создаем витрину
	storefront := &models.Storefront{
		UserID:                userID,
		Name:                  create.Name,
		Description:           create.Description,
		Slug:                  slug,
		Status:                "active",
		CreationTransactionID: &transactionID,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	// Сохраняем витрину в БД
	storefrontID, err := s.storage.CreateStorefront(ctx, storefront)
	if err != nil {
		return nil, fmt.Errorf("failed to create storefront: %w", err)
	}

	storefront.ID = storefrontID

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return storefront, nil
}

// GetUserStorefronts возвращает все витрины пользователя
func (s *StorefrontService) GetUserStorefronts(ctx context.Context, userID int) ([]models.Storefront, error) {
	return s.storage.GetUserStorefronts(ctx, userID)
}

// GetStorefrontByID возвращает витрину по ID
func (s *StorefrontService) GetStorefrontByID(ctx context.Context, id int, userID int) (*models.Storefront, error) {
	storefront, err := s.storage.GetStorefrontByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверяем права доступа
	if storefront.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return storefront, nil
}

func (s *StorefrontService) GetPublicStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	storefront, err := s.storage.GetStorefrontByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if storefront.Status != "active" {
		return nil, fmt.Errorf("storefront is not active")
	}

	return storefront, nil
}

// UpdateStorefront обновляет информацию о витрине
func (s *StorefrontService) UpdateStorefront(ctx context.Context, storefront *models.Storefront, userID int) error {
	// Проверяем права доступа
	existing, err := s.storage.GetStorefrontByID(ctx, storefront.ID)
	if err != nil {
		return err
	}

	if existing.UserID != userID {
		return fmt.Errorf("access denied")
	}

	return s.storage.UpdateStorefront(ctx, storefront)
}

// DeleteStorefront удаляет витрину
func (s *StorefrontService) DeleteStorefront(ctx context.Context, id int, userID int) error {
	// Проверяем права доступа
	existing, err := s.storage.GetStorefrontByID(ctx, id)
	if err != nil {
		return err
	}

	if existing.UserID != userID {
		return fmt.Errorf("access denied")
	}

	return s.storage.DeleteStorefront(ctx, id)
}

// CreateImportSource создает новый источник импорта
func (s *StorefrontService) CreateImportSource(ctx context.Context, source *models.ImportSourceCreate, userID int) (*models.ImportSource, error) {
	// Проверяем права доступа к витрине
	storefront, err := s.storage.GetStorefrontByID(ctx, source.StorefrontID)
	if err != nil {
		return nil, err
	}

	if storefront.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	importSource := &models.ImportSource{
		StorefrontID: source.StorefrontID,
		Type:         source.Type,
		URL:          source.URL,
		AuthData:     source.AuthData,
		Schedule:     source.Schedule,
		Mapping:      source.Mapping,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	sourceID, err := s.storage.CreateImportSource(ctx, importSource)
	if err != nil {
		return nil, err
	}

	importSource.ID = sourceID
	return importSource, nil
}

// UpdateImportSource обновляет источник импорта
func (s *StorefrontService) UpdateImportSource(ctx context.Context, source *models.ImportSource, userID int) error {
	// Проверяем права доступа
	existing, err := s.storage.GetImportSourceByID(ctx, source.ID)
	if err != nil {
		return err
	}

	// Получаем информацию о витрине
	storefront, err := s.storage.GetStorefrontByID(ctx, existing.StorefrontID)
	if err != nil {
		return err
	}

	if storefront.UserID != userID {
		return fmt.Errorf("access denied")
	}

	return s.storage.UpdateImportSource(ctx, source)
}

// DeleteImportSource удаляет источник импорта
func (s *StorefrontService) DeleteImportSource(ctx context.Context, id int, userID int) error {
	// Проверяем права доступа
	existing, err := s.storage.GetImportSourceByID(ctx, id)
	if err != nil {
		return err
	}

	// Получаем информацию о витрине
	storefront, err := s.storage.GetStorefrontByID(ctx, existing.StorefrontID)
	if err != nil {
		return err
	}

	if storefront.UserID != userID {
		return fmt.Errorf("access denied")
	}

	return s.storage.DeleteImportSource(ctx, id)
}

// GetImportSources возвращает источники импорта для витрины
func (s *StorefrontService) GetImportSources(ctx context.Context, storefrontID int, userID int) ([]models.ImportSource, error) {
	// Проверяем права доступа к витрине
	storefront, err := s.storage.GetStorefrontByID(ctx, storefrontID)
	if err != nil {
		return nil, err
	}

	if storefront.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return s.storage.GetImportSources(ctx, storefrontID)
}

// Проверка доступности URL перед импортом
func (s *StorefrontService) checkURLAccessibility(url string) error {
    client := &http.Client{
        Timeout: 10 * time.Second,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if len(via) >= 10 {
                return fmt.Errorf("too many redirects")
            }
            return nil
        },
    }
    
    req, err := http.NewRequest("HEAD", url, nil)
    if err != nil {
        return fmt.Errorf("error creating HTTP request: %w", err)
    }
    
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("error checking URL accessibility: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("URL returned unexpected status code: %d", resp.StatusCode)
    }
    
    return nil
}

// Обновленная функция RunImport с проверкой доступности URL
func (s *StorefrontService) RunImport(ctx context.Context, sourceID int, userID int) (*models.ImportHistory, error) {
    // Получаем информацию об источнике
    source, err := s.storage.GetImportSourceByID(ctx, sourceID)
    if err != nil {
        return nil, fmt.Errorf("error getting import source: %w", err)
    }
    
    // Проверяем права доступа
    storefront, err := s.storage.GetStorefrontByID(ctx, source.StorefrontID)
    if err != nil {
        return nil, fmt.Errorf("error getting storefront: %w", err)
    }
    
    if storefront.UserID != userID {
        return nil, fmt.Errorf("access denied")
    }
    
    // Проверяем наличие URL
    if source.URL == "" {
        return nil, fmt.Errorf("no URL configured for import source")
    }
    
    // Проверяем доступность URL
    if err := s.checkURLAccessibility(source.URL); err != nil {
        // Если URL локальный для localhost, предлагаем альтернативу
        if strings.Contains(source.URL, "localhost") || strings.Contains(source.URL, "127.0.0.1") {
            log.Printf("Обнаружен локальный URL %s, который может быть недоступен из контейнера", source.URL)
            return nil, fmt.Errorf("localhost URL detected which may not be accessible from container. Try using host.docker.internal instead of localhost or IP address of your host machine: %w", err)
        }
        return nil, fmt.Errorf("URL is not accessible: %w", err)
    }
    
    // Создаем запись в истории импорта
    history := &models.ImportHistory{
        SourceID:  sourceID,
        Status:    "pending",
        StartedAt: time.Now(),
    }
    
    historyID, err := s.storage.CreateImportHistory(ctx, history)
    if err != nil {
        return nil, fmt.Errorf("error creating import history: %w", err)
    }
    history.ID = historyID


	// Загружаем данные с удаленного URL
	client := &http.Client{
		Timeout: 60 * time.Second, // Увеличенный таймаут для больших файлов
	}

	// Запрашиваем файл с сервера
	resp, err := client.Get(source.URL)
	if err != nil {
		history.Status = "failed"
		history.Log = fmt.Sprintf("Error downloading file from URL %s: %v", source.URL, err)
		finishTime := time.Now()
		history.FinishedAt = &finishTime
		s.storage.UpdateImportHistory(ctx, history)
		return history, fmt.Errorf("error downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		history.Status = "failed"
		history.Log = fmt.Sprintf("Error response from URL %s: %s", source.URL, resp.Status)
		finishTime := time.Now()
		history.FinishedAt = &finishTime
		s.storage.UpdateImportHistory(ctx, history)
		return history, fmt.Errorf("error response from URL: %s", resp.Status)
	}

	// Проверяем тип контента
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/csv") &&
		!strings.Contains(contentType, "application/csv") &&
		!strings.Contains(contentType, "text/plain") {
		history.Status = "failed"
		history.Log = fmt.Sprintf("Invalid content type: %s. Expected CSV file.", contentType)
		finishTime := time.Now()
		history.FinishedAt = &finishTime
		s.storage.UpdateImportHistory(ctx, history)
		return history, fmt.Errorf("invalid content type: %s", contentType)
	}

	// Обновляем статус
	history.Status = "in_progress"
	s.storage.UpdateImportHistory(ctx, history)

	// Запускаем импорт из CSV
	updatedHistory, err := s.ImportCSV(ctx, sourceID, resp.Body, nil, userID)
	if err != nil {
		if updatedHistory == nil {
			history.Status = "failed"
			history.Log = fmt.Sprintf("Error importing CSV: %v", err)
			finishTime := time.Now()
			history.FinishedAt = &finishTime
			s.storage.UpdateImportHistory(ctx, history)
			return history, fmt.Errorf("error importing CSV: %w", err)
		}
		return updatedHistory, err
	}

	return updatedHistory, nil
}

// ImportCSV импортирует данные из CSV с опциональной поддержкой ZIP-архива для изображений
func (s *StorefrontService) ImportCSV(ctx context.Context, sourceID int, reader io.Reader, zipFile io.Reader, userID int) (*models.ImportHistory, error) {
	// Получаем информацию об источнике
	source, err := s.storage.GetImportSourceByID(ctx, sourceID)
	if err != nil {
		return nil, fmt.Errorf("error getting import source: %w", err)
	}

	// Проверяем права доступа
	storefront, err := s.storage.GetStorefrontByID(ctx, source.StorefrontID)
	if err != nil {
		return nil, fmt.Errorf("error getting storefront: %w", err)
	}

	if storefront.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// Создаем историю импорта
	history := &models.ImportHistory{
		SourceID:  sourceID,
		Status:    "in_progress",
		StartedAt: time.Now(),
	}

	historyID, err := s.storage.CreateImportHistory(ctx, history)
	if err != nil {
		return nil, fmt.Errorf("error creating import history: %w", err)
	}
	history.ID = historyID

	// Инициализируем ZIP-архив, если он был предоставлен
	var zipReader *zip.Reader
	if zipFile != nil {
		// Читаем все содержимое в буфер, так как zip.NewReader требует io.ReaderAt
		zipData, err := ioutil.ReadAll(zipFile)
		if err != nil {
			history.Status = "failed"
			history.Log = fmt.Sprintf("Failed to read ZIP archive: %v", err)
			finishTime := time.Now()
			history.FinishedAt = &finishTime
			s.storage.UpdateImportHistory(ctx, history)
			return history, fmt.Errorf("failed to read ZIP archive: %w", err)
		}

		// Создаем zip.Reader из буфера
		zipReader, err = zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
		if err != nil {
			history.Status = "failed"
			history.Log = fmt.Sprintf("Failed to parse ZIP archive: %v", err)
			finishTime := time.Now()
			history.FinishedAt = &finishTime
			s.storage.UpdateImportHistory(ctx, history)
			return history, fmt.Errorf("failed to parse ZIP archive: %w", err)
		}

		log.Printf("ZIP archive loaded successfully with %d files", len(zipReader.File))
	}

	// Читаем CSV файл
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'             // Используем точку с запятой как разделитель
	csvReader.LazyQuotes = true       // Разрешаем нестрогие кавычки
	csvReader.TrimLeadingSpace = true // Убираем начальные пробелы

	// Читаем заголовок
	headers, err := csvReader.Read()
	if err != nil {
		history.Status = "failed"
		history.Log = fmt.Sprintf("Failed to read CSV header: %v", err)
		finishTime := time.Now()
		history.FinishedAt = &finishTime
		s.storage.UpdateImportHistory(ctx, history)
		return history, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Логируем заголовки
	log.Printf("CSV Import: Headers received: %v", headers)

	// Создаем маппинг колонок
	columnMap := make(map[string]int)
	for i, header := range headers {
		header = strings.TrimSpace(header)
		columnMap[header] = i
	}

	// Проверяем наличие обязательных полей
	requiredFields := []string{"id", "title", "description", "price", "category_id"}
	missing := []string{}
	for _, field := range requiredFields {
		if _, ok := columnMap[field]; !ok {
			missing = append(missing, field)
		}
	}

	if len(missing) > 0 {
		errMsg := fmt.Sprintf("Missing required fields: %s", strings.Join(missing, ", "))
		history.Status = "failed"
		history.Log = errMsg
		finishTime := time.Now()
		history.FinishedAt = &finishTime
		s.storage.UpdateImportHistory(ctx, history)
		return history, fmt.Errorf(errMsg)
	}

	// Константа для ID категории "прочее"
	const DefaultCategoryID = 9999

	// Проверяем существование категории "прочее", создаем если нет
	_, err = s.storage.GetCategoryByID(ctx, DefaultCategoryID)
	// Если категория не найдена, логируем это, но продолжаем импорт
	if err != nil {
		log.Printf("Default category (ID: %d) not found. Import will use this ID anyway.", DefaultCategoryID)
	}

	// Обработка строк
	var itemsTotal, itemsImported, itemsFailed int
	var errorLog strings.Builder

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			itemsFailed++
			errorLog.WriteString(fmt.Sprintf("Error reading row: %v\n", err))
			continue
		}

		itemsTotal++

		// Извлекаем данные из строки
		var listingData models.MarketplaceListing

		// Проверяем, что индексы не выходят за пределы массива
		idIdx, ok := columnMap["id"]
		if !ok || idIdx >= len(row) {
			itemsFailed++
			errorLog.WriteString("Row missing required 'id' field\n")
			continue
		}

		// Получаем title
		titleIdx, ok := columnMap["title"]
		if !ok || titleIdx >= len(row) {
			itemsFailed++
			errorLog.WriteString("Row missing required 'title' field\n")
			continue
		}
		listingData.Title = strings.TrimSpace(row[titleIdx])

		// Получаем description
		descIdx, ok := columnMap["description"]
		if !ok || descIdx >= len(row) {
			itemsFailed++
			errorLog.WriteString("Row missing required 'description' field\n")
			continue
		}
		listingData.Description = strings.TrimSpace(row[descIdx])

		// Получаем price
		priceIdx, ok := columnMap["price"]
		if !ok || priceIdx >= len(row) {
			itemsFailed++
			errorLog.WriteString("Row missing required 'price' field\n")
			continue
		}
		price, err := strconv.ParseFloat(strings.TrimSpace(row[priceIdx]), 64)
		if err != nil {
			itemsFailed++
			errorLog.WriteString(fmt.Sprintf("Invalid price value '%s': %v\n", row[priceIdx], err))
			continue
		}
		listingData.Price = price

		// Получаем category_id
		catIdx, ok := columnMap["category_id"]
		if !ok || catIdx >= len(row) {
			itemsFailed++
			errorLog.WriteString("Row missing required 'category_id' field\n")
			continue
		}
		categoryID, err := strconv.Atoi(strings.TrimSpace(row[catIdx]))
		if err != nil {
			// Если категория не является числом, используем категорию "прочее"
			errorLog.WriteString(fmt.Sprintf("Warning: Invalid category_id value '%s': %v. Using default category (ID: %d)\n",
				row[catIdx], err, DefaultCategoryID))
			categoryID = DefaultCategoryID
		} else {
			// Проверяем существование категории
			_, err = s.storage.GetCategoryByID(ctx, categoryID)
			if err != nil {
				// Если категория не найдена, используем категорию "прочее"
				errorLog.WriteString(fmt.Sprintf("Warning: Category with ID '%d' not found. Using default category (ID: %d)\n",
					categoryID, DefaultCategoryID))
				categoryID = DefaultCategoryID
			}
		}
		listingData.CategoryID = categoryID

		// Получаем condition
		if condIdx, ok := columnMap["condition"]; ok && condIdx < len(row) {
			condition := strings.TrimSpace(row[condIdx])
			if condition != "new" && condition != "used" {
				condition = "new" // По умолчанию новый товар
				errorLog.WriteString(fmt.Sprintf("Warning: Invalid condition value '%s', using 'new' as default\n", row[condIdx]))
			}
			listingData.Condition = condition
		} else {
			listingData.Condition = "new" // По умолчанию новый товар
		}

		// Получаем status
		if statusIdx, ok := columnMap["status"]; ok && statusIdx < len(row) {
			status := strings.TrimSpace(row[statusIdx])
			if status != "active" && status != "inactive" {
				status = "active" // По умолчанию активный товар
				errorLog.WriteString(fmt.Sprintf("Warning: Invalid status value '%s', using 'active' as default\n", row[statusIdx]))
			}
			listingData.Status = status
		} else {
			listingData.Status = "active" // По умолчанию активный товар
		}

		// Получаем location
		if locIdx, ok := columnMap["location"]; ok && locIdx < len(row) {
			listingData.Location = strings.TrimSpace(row[locIdx])
		}

		// Получаем latitude
		if latIdx, ok := columnMap["latitude"]; ok && latIdx < len(row) {
			latStr := strings.TrimSpace(row[latIdx])
			if latStr != "" {
				lat, err := strconv.ParseFloat(latStr, 64)
				if err == nil {
					listingData.Latitude = &lat
				} else {
					errorLog.WriteString(fmt.Sprintf("Warning: Invalid latitude value '%s': %v, ignoring\n", latStr, err))
				}
			}
		}

		// Получаем longitude
		if lngIdx, ok := columnMap["longitude"]; ok && lngIdx < len(row) {
			lngStr := strings.TrimSpace(row[lngIdx])
			if lngStr != "" {
				lng, err := strconv.ParseFloat(lngStr, 64)
				if err == nil {
					listingData.Longitude = &lng
				} else {
					errorLog.WriteString(fmt.Sprintf("Warning: Invalid longitude value '%s': %v, ignoring\n", lngStr, err))
				}
			}
		}

		// Получаем город
		if cityIdx, ok := columnMap["address_city"]; ok && cityIdx < len(row) {
			listingData.City = strings.TrimSpace(row[cityIdx])
		}

		// Получаем страну
		if countryIdx, ok := columnMap["address_country"]; ok && countryIdx < len(row) {
			listingData.Country = strings.TrimSpace(row[countryIdx])
		}

		// Получаем show_on_map
		if showOnMapIdx, ok := columnMap["show_on_map"]; ok && showOnMapIdx < len(row) {
			showOnMapStr := strings.TrimSpace(row[showOnMapIdx])
			if showOnMapStr == "true" || showOnMapStr == "1" {
				listingData.ShowOnMap = true
			} else {
				listingData.ShowOnMap = false
			}
		} else {
			listingData.ShowOnMap = true // По умолчанию показываем на карте
		}

		// Получаем original_language
		if langIdx, ok := columnMap["original_language"]; ok && langIdx < len(row) {
			listingData.OriginalLanguage = strings.TrimSpace(row[langIdx])
		} else {
			listingData.OriginalLanguage = "sr" // По умолчанию сербский язык
		}

		// Устанавливаем связь с витриной
		listingData.UserID = userID
		listingData.StorefrontID = &storefront.ID

		// Создание объявления
		listingID, err := s.storage.CreateListing(ctx, &listingData)
		if err != nil {
			itemsFailed++
			errorLog.WriteString(fmt.Sprintf("Error creating listing: %v\n", err))
			continue
		}

		// Переменная для отслеживания, добавлены ли изображения
		imagesAdded := false

		// Если есть колонка с изображениями, обрабатываем их с новым подходом
		if imagesIdx, ok := columnMap["images"]; ok && imagesIdx < len(row) && row[imagesIdx] != "" {
			imagesStr := row[imagesIdx]
		
			// Используем асинхронную версию метода
			s.ProcessImportImagesAsync(ctx, listingID, imagesStr, zipReader)
			
			// Удаляем эту проверку, так как мы использовали асинхронную загрузку
			imagesAdded = true
			log.Printf("Запущена асинхронная обработка изображений для листинга %d", listingID)
		}

		// Получаем созданное объявление для индексации ПОСЛЕ добавления изображений
		if imagesAdded {
			// Небольшая задержка для гарантии, что изображения сохранились в БД
			time.Sleep(200 * time.Millisecond)
		}

		createdListing, err := s.storage.GetListingByID(ctx, listingID)
		if err != nil {
			errorLog.WriteString(fmt.Sprintf("Warning: Listing created but failed to retrieve for indexing: %v\n", err))
		} else {
			// Индексируем объявление в поисковом движке
			err = s.storage.IndexListing(ctx, createdListing)
			if err != nil {
				errorLog.WriteString(fmt.Sprintf("Warning: Listing created but failed to index: %v\n", err))
			}
		}

		itemsImported++
	}

	// Обновляем историю импорта
	finishTime := time.Now()
	history.FinishedAt = &finishTime
	history.ItemsTotal = itemsTotal
	history.ItemsImported = itemsImported
	history.ItemsFailed = itemsFailed
	history.Log = errorLog.String()

	if itemsFailed > 0 {
		if itemsImported > 0 {
			history.Status = "partial"
		} else {
			history.Status = "failed"
		}
	} else {
		history.Status = "success"
	}

	err = s.storage.UpdateImportHistory(ctx, history)
	if err != nil {
		return nil, fmt.Errorf("failed to update import history: %w", err)
	}

	// Обновляем информацию об источнике
	source.LastImportAt = &finishTime
	source.LastImportStatus = history.Status
	source.LastImportLog = errorLog.String()
	s.storage.UpdateImportSource(ctx, source)

	return history, nil
}

// GetImportHistory возвращает историю импорта
func (s *StorefrontService) GetImportHistory(ctx context.Context, sourceID int, userID int, limit, offset int) ([]models.ImportHistory, error) {
	// Проверяем права доступа
	source, err := s.storage.GetImportSourceByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}

	// Получаем информацию о витрине
	storefront, err := s.storage.GetStorefrontByID(ctx, source.StorefrontID)
	if err != nil {
		return nil, err
	}

	if storefront.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return s.storage.GetImportHistory(ctx, sourceID, limit, offset)
}

// generateSlug создает уникальный slug на основе имени
func generateSlug(name string) string {
	// Очищаем строку от специальных символов
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")

	// Добавляем случайное число для уникальности
	rand.Seed(time.Now().UnixNano())
	randomSuffix := rand.Intn(10000)

	return fmt.Sprintf("%s-%d", slug, randomSuffix)
}

// GetImportSourceByID возвращает источник импорта по ID с проверкой прав доступа
func (s *StorefrontService) GetImportSourceByID(ctx context.Context, id int, userID int) (*models.ImportSource, error) {
	// Отладочный лог
	log.Printf("Getting import source ID: %d for user: %d", id, userID)

	// Получаем информацию об источнике
	source, err := s.storage.GetImportSourceByID(ctx, id)
	if err != nil {
		log.Printf("Error getting import source: %v", err)
		return nil, fmt.Errorf("error getting import source: %w", err)
	}

	// Отладочный лог
	log.Printf("Found import source: %+v", source)

	// Проверяем права доступа
	storefront, err := s.storage.GetStorefrontByID(ctx, source.StorefrontID)
	if err != nil {
		log.Printf("Error getting storefront: %v", err)
		return nil, fmt.Errorf("error getting storefront: %w", err)
	}

	// Отладочный лог
	log.Printf("Found storefront: %+v", storefront)

	if storefront.UserID != userID {
		log.Printf("Access denied - storefront owner: %d, requesting user: %d", storefront.UserID, userID)
		return nil, fmt.Errorf("access denied")
	}

	return source, nil
}

// backend/internal/proj/storefront/service/storefront.go

// ImportXMLFromZip выполняет импорт данных из XML файла внутри ZIP-архива
func (s *StorefrontService) ImportXMLFromZip(ctx context.Context, sourceID int, reader io.Reader, userID int) (*models.ImportHistory, error) {
    // Проверяем права доступа
    source, err := s.storage.GetImportSourceByID(ctx, sourceID)
    if err != nil {
        return nil, fmt.Errorf("error getting import source: %w", err)
    }

    // Получаем информацию о витрине
    storefront, err := s.storage.GetStorefrontByID(ctx, source.StorefrontID)
    if err != nil {
        return nil, fmt.Errorf("error getting storefront: %w", err)
    }

    if storefront.UserID != userID {
        return nil, fmt.Errorf("access denied")
    }

    // Создаем запись в истории импорта
    history := &models.ImportHistory{
        SourceID:  sourceID,
        Status:    "in_progress",
        StartedAt: time.Now(),
    }

    historyID, err := s.storage.CreateImportHistory(ctx, history)
    if err != nil {
        return nil, fmt.Errorf("error creating import history: %w", err)
    }
    history.ID = historyID

    // Читаем ZIP-архив
    log.Printf("Reading ZIP archive from source ID %d", sourceID)
    zipData, err := io.ReadAll(reader)
    if err != nil {
        history.Status = "failed"
        history.Log = fmt.Sprintf("Failed to read ZIP archive: %v", err)
        finishTime := time.Now()
        history.FinishedAt = &finishTime
        s.storage.UpdateImportHistory(ctx, history)
        return history, fmt.Errorf("failed to read ZIP archive: %w", err)
    }

    log.Printf("Read %d bytes from ZIP archive", len(zipData))

    // Создаем zip.Reader из буфера
    zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
    if err != nil {
        history.Status = "failed"
        history.Log = fmt.Sprintf("Failed to parse ZIP archive: %v", err)
        finishTime := time.Now()
        history.FinishedAt = &finishTime
        s.storage.UpdateImportHistory(ctx, history)
        return history, fmt.Errorf("failed to parse ZIP archive: %w", err)
    }

    log.Printf("ZIP archive parsed successfully, contains %d files", len(zipReader.File))

    // Поиск XML файла в архиве
    var xmlFile *zip.File
    for _, file := range zipReader.File {
        log.Printf("Found file in ZIP: %s", file.Name)
        if strings.HasSuffix(strings.ToLower(file.Name), ".xml") {
            xmlFile = file
            log.Printf("Selected as XML file: %s", file.Name)
            break
        }
    }

    if xmlFile == nil {
        history.Status = "failed"
        history.Log = "No XML file found in the ZIP archive"
        finishTime := time.Now()
        history.FinishedAt = &finishTime
        s.storage.UpdateImportHistory(ctx, history)
        return history, fmt.Errorf("no XML file found in the ZIP archive")
    }

    // Открываем XML файл
    rc, err := xmlFile.Open()
    if err != nil {
        history.Status = "failed"
        history.Log = fmt.Sprintf("Failed to open XML file: %v", err)
        finishTime := time.Now()
        history.FinishedAt = &finishTime
        s.storage.UpdateImportHistory(ctx, history)
        return history, fmt.Errorf("failed to open XML file: %w", err)
    }
    defer rc.Close()

    // Парсим XML
    xmlContent, err := io.ReadAll(rc)
    if err != nil {
        history.Status = "failed"
        history.Log = fmt.Sprintf("Failed to read XML content: %v", err)
        finishTime := time.Now()
        history.FinishedAt = &finishTime
        s.storage.UpdateImportHistory(ctx, history)
        return history, fmt.Errorf("failed to read XML content: %w", err)
    }

    log.Printf("Read %d bytes of XML content", len(xmlContent))

    // Парсим содержимое XML
    var itemsTotal, itemsImported, itemsFailed int
    var errorLog strings.Builder

    // Используем потоковый парсер XML вместо регулярных выражений
    itemsTotal, itemsImported, itemsFailed, err = s.processXMLContentStream(ctx, bytes.NewReader(xmlContent), storefront.ID, userID, &errorLog)
    if err != nil {
        history.Status = "failed"
        history.Log = fmt.Sprintf("Failed to process XML content: %v\n%s", err, errorLog.String())
        finishTime := time.Now()
        history.FinishedAt = &finishTime
        s.storage.UpdateImportHistory(ctx, history)
        return history, fmt.Errorf("failed to process XML content: %w", err)
    }

    // Обновляем историю импорта
    finishTime := time.Now()
    history.FinishedAt = &finishTime
    history.ItemsTotal = itemsTotal
    history.ItemsImported = itemsImported
    history.ItemsFailed = itemsFailed
    history.Log = errorLog.String()

    if itemsFailed > 0 {
        if itemsImported > 0 {
            history.Status = "partial"
        } else {
            history.Status = "failed"
        }
    } else {
        history.Status = "success"
    }

    log.Printf("Updating import history: Total=%d, Imported=%d, Failed=%d, Status=%s", 
               history.ItemsTotal, history.ItemsImported, history.ItemsFailed, history.Status)

    err = s.storage.UpdateImportHistory(ctx, history)
    if err != nil {
        return nil, fmt.Errorf("failed to update import history: %w", err)
    }

    // Обновляем информацию об источнике
    source.LastImportAt = &finishTime
    source.LastImportStatus = history.Status
    source.LastImportLog = errorLog.String()
    s.storage.UpdateImportSource(ctx, source)

    return history, nil
}
// Функция для сопоставления атрибутов из импорта с атрибутами в системе
func (s *StorefrontService) mapImportAttributes(ctx context.Context, categoryID int, attrMap map[string]string) ([]models.ListingAttributeValue, error) {
    // Получаем атрибуты категории
    categoryAttributes, err := s.storage.GetCategoryAttributes(ctx, categoryID)
    if err != nil {
        return nil, fmt.Errorf("error fetching category attributes: %w", err)
    }

    var result []models.ListingAttributeValue

    // Сопоставляем входящие атрибуты с существующими
    for attrName, attrValue := range attrMap {
        for _, catAttr := range categoryAttributes {
            // Проверяем совпадение по имени или похожие имена
            if strings.EqualFold(catAttr.Name, attrName) || 
               strings.EqualFold(catAttr.DisplayName, attrName) ||
               isSimilarAttributeName(catAttr.Name, attrName) {
                
                // Создаём атрибут с соответствующим типом
                attr := models.ListingAttributeValue{
                    AttributeID:   catAttr.ID,
                    AttributeName: catAttr.Name,
                    AttributeType: catAttr.AttributeType,
                    DisplayName:   catAttr.DisplayName,
                }

                // Заполняем значение в зависимости от типа
                switch catAttr.AttributeType {
                case "number":
                    if numVal, err := strconv.ParseFloat(attrValue, 64); err == nil {
                        attr.NumericValue = &numVal
                        attr.DisplayValue = fmt.Sprintf("%g", numVal)
                    }
                case "boolean":
                    boolVal := attrValue == "true" || attrValue == "1" || 
                               strings.EqualFold(attrValue, "да") || strings.EqualFold(attrValue, "yes")
                    attr.BooleanValue = &boolVal
                    if boolVal {
						attr.DisplayValue = "Да"
					} else {
						attr.DisplayValue = "Нет"
					}
                default: // text, select и другие текстовые типы
                    attr.TextValue = &attrValue
                    attr.DisplayValue = attrValue
                }

                result = append(result, attr)
                break
            }
        }
    }

    return result, nil
}

// Функция для определения похожих имен атрибутов
func isSimilarAttributeName(attrName, importName string) bool {
    // Нормализуем строки
    attrName = strings.ToLower(attrName)
    importName = strings.ToLower(importName)
    
    // Удаляем пробелы и специальные символы
    attrName = regexp.MustCompile(`[^a-z0-9]`).ReplaceAllString(attrName, "")
    importName = regexp.MustCompile(`[^a-z0-9]`).ReplaceAllString(importName, "")
    
    // Проверяем на полное соответствие или вхождение одной строки в другую
    return attrName == importName || 
           strings.Contains(attrName, importName) || 
           strings.Contains(importName, attrName)
}
// processXMLContentStream обрабатывает содержимое XML с использованием потокового парсера
func (s *StorefrontService) processXMLContentStream(ctx context.Context, reader io.Reader, storefrontID int, userID int, errorLog *strings.Builder) (int, int, int, error) {
    var itemsTotal, itemsImported, itemsFailed int

    log.Printf("Starting streaming XML processing for storefront ID %d", storefrontID)

    // Константа для ID категории "прочее"
    const DefaultCategoryID = 9999

    // Создаем XML декодер
    decoder := xml.NewDecoder(reader)
    
    // Переменные для сохранения текущего артикула и его полей
    var (
        inArtikal bool
        inField   string
        id        string
        naziv     string
        kategorija1 string
        kategorija2 string
        kategorija3 string
        opis      string
        mpCena    string
		vpCena	  string
        dostupan  string
        naAkciji  string
        slike     []string
        inSlike   bool
    )

    // Обрабатываем каждый XML токен
    for {
        token, err := decoder.Token()
        if err == io.EOF {
            break
        }
        if err != nil {
            return itemsTotal, itemsImported, itemsFailed, fmt.Errorf("error decoding XML: %w", err)
        }

        switch t := token.(type) {
        case xml.StartElement:
            // Начало элемента
            if t.Name.Local == "artikal" {
                inArtikal = true
                // Сбрасываем переменные для нового артикула
                id = ""
                naziv = ""
                kategorija1 = ""
                kategorija2 = ""
                kategorija3 = ""
                opis = ""
                mpCena = ""
				vpCena = ""
                dostupan = ""
                naAkciji = ""
                slike = nil
            } else if inArtikal {
                if t.Name.Local == "slike" {
                    inSlike = true
                } else if inSlike && t.Name.Local == "slika" {
                    // Будем собирать данные изображения внутри slika
                    inField = "slika"
                } else {
                    inField = t.Name.Local
                }
            }
        case xml.EndElement:
            // Конец элемента
            if t.Name.Local == "artikal" && inArtikal {
                inArtikal = false
                itemsTotal++
                
                // Обрабатываем собранные данные артикула
                if naziv == "" {
                    itemsFailed++
                    errorLog.WriteString(fmt.Sprintf("Item with ID %s skipped: no title\n", id))
                    continue
                }

				
                // Преобразуем цену в число
				price := 0.0

				// Если розничная цена указана и не равна нулю или ".0000"
				mpCenaClean := strings.TrimSpace(mpCena)
				if mpCenaClean != "" && mpCenaClean != ".0000" && mpCenaClean != "0.0000" {
					// Используем розничную цену
					price, err = parsePrice(mpCena)
					if err != nil {
						// Если ошибка парсинга, пытаемся использовать оптовую цену
						price, err = parsePrice(vpCena)
						if err != nil {
							itemsFailed++
							errorLog.WriteString(fmt.Sprintf("Item with ID %s skipped: invalid price: %v\n", id, err))
							continue
						}
					}
				} else {
					// Если розничная цена не указана или равна нулю, используем оптовую цену с наценкой
					price, err = parsePrice(vpCena)
					if err != nil {
						itemsFailed++
						errorLog.WriteString(fmt.Sprintf("Item with ID %s skipped: invalid price: %v\n", id, err))
						continue
					}
					
					// Применяем наценку к оптовой цене (например, 30%)
					price = price * 1.3
				}

                // Находим или создаем категорию
                categoryID := DefaultCategoryID
                if kategorija1 != "" {
                    catID, err := s.findOrCreateCategory(ctx, kategorija1, kategorija2, kategorija3)
                    if err == nil {
                        categoryID = catID
                    } else {
                        errorLog.WriteString(fmt.Sprintf("Warning for item %s: %v. Using default category.\n", id, err))
                    }
                }

                // Создаем объявление
                listing := &models.MarketplaceListing{
                    UserID:       userID,
                    CategoryID:   categoryID,
                    StorefrontID: &storefrontID,
                    Title:        naziv,
                    Description:  opis,
                    Price:        price,
                    Condition:    "new", // По умолчанию новый товар
                    Status: func() string {
                        if dostupan == "1" {
                            return "active"
                        }
                        return "inactive"
                    }(),
                    ShowOnMap:        false,
                    OriginalLanguage: "ru", // Предполагаем русский язык по умолчанию
                }

                // Если товар на акции, отмечаем это в описании
                if naAkciji == "1" {
                    listing.Description = "🔥 sale! 🔥\n\n" + listing.Description
                }

                // Создание объявления
                listingID, err := s.storage.CreateListing(ctx, listing)
                if err != nil {
                    itemsFailed++
                    errorLog.WriteString(fmt.Sprintf("Error creating listing for item %s: %v\n", id, err))
                    continue
                }

                // Если есть изображения, обрабатываем их
				if len(slike) > 0 {
					imagesStr := strings.Join(slike, ",")
					// Используем асинхронную обработку изображений
					s.ProcessImportImagesAsync(ctx, listingID, imagesStr, nil)
				}

                // Получаем созданное объявление для индексации
                createdListing, err := s.storage.GetListingByID(ctx, listingID)
                if err != nil {
                    errorLog.WriteString(fmt.Sprintf("Warning: Listing created but failed to retrieve for indexing: %v\n", err))
                } else {
                    // Индексируем объявление в поисковом движке
                    err = s.storage.IndexListing(ctx, createdListing)
                    if err != nil {
                        errorLog.WriteString(fmt.Sprintf("Warning: Listing created but failed to index: %v\n", err))
                    }
                }

                itemsImported++
                log.Printf("Successfully imported item %s (ID: %s) with DB ID %d", naziv, id, listingID)
            } else if t.Name.Local == "slike" {
                inSlike = false
            } else {
                inField = ""
            }
        case xml.CharData:
            // Текстовые данные
            if inArtikal && inField != "" {
                text := string(t)
                switch inField {
                case "id":
                    id = strings.TrimSpace(text)
                case "naziv":
                    naziv = strings.TrimSpace(text)
                case "kategorija1":
                    kategorija1 = strings.TrimSpace(text)
                case "kategorija2":
                    kategorija2 = strings.TrimSpace(text)
                case "kategorija3":
                    kategorija3 = strings.TrimSpace(text)
                case "opis":
                    opis = strings.TrimSpace(text)
                case "mpCena":
                    mpCena = strings.TrimSpace(text)
                case "dostupan":
                    dostupan = strings.TrimSpace(text)
                case "naAkciji":
                    naAkciji = strings.TrimSpace(text)
                case "slika":
                    if text = strings.TrimSpace(text); text != "" {
                        slike = append(slike, text)
                    }
                }
            }
        case xml.Comment:
            // Пропускаем комментарии
        case xml.ProcInst:
            // Пропускаем инструкции процессора
        case xml.Directive:
            // Пропускаем директивы
        }
    }

    log.Printf("Streaming XML processing completed. Total: %d, Imported: %d, Failed: %d", itemsTotal, itemsImported, itemsFailed)
    return itemsTotal, itemsImported, itemsFailed, nil
}
// processXMLContent обрабатывает содержимое XML и создает товары
func (s *StorefrontService) processXMLContent(ctx context.Context, xmlContent string, storefrontID int, userID int, errorLog *strings.Builder) (int, int, int, error) {
    var itemsTotal, itemsImported, itemsFailed int

    // Добавим логирование для отладки
    log.Printf("Starting XML processing for storefront ID %d, content length: %d bytes", storefrontID, len(xmlContent))

    // Константа для ID категории "прочее"
    const DefaultCategoryID = 9999

    // Используем regexp для поиска всех <artikal> элементов
    re := regexp.MustCompile(`<artikal>(.*?)</artikal>`)
    matches := re.FindAllStringSubmatch(xmlContent, -1)
    
    // Добавим логирование количества найденных товаров
    log.Printf("Found %d <artikal> elements in XML", len(matches))

    for _, match := range matches {
        if len(match) < 2 {
            continue
        }

        itemsTotal++
        artikal := match[1]

        // Извлекаем необходимые поля из элемента artikal
        id := extractField(artikal, "id")
        naziv := cleanXMLContent(extractField(artikal, "naziv"))
        kategorija1 := cleanXMLContent(extractField(artikal, "kategorija1"))
        kategorija2 := cleanXMLContent(extractField(artikal, "kategorija2"))
        kategorija3 := cleanXMLContent(extractField(artikal, "kategorija3"))
        opis := cleanXMLContent(extractField(artikal, "opis"))
        mpCena := extractField(artikal, "mpCena")
        dostupan := extractField(artikal, "dostupan")
        naAkciji := extractField(artikal, "naAkciji")

        // Извлекаем ссылки на изображения
        slike := extractImages(artikal)

        // Добавим логирование для отладки отдельных товаров
        log.Printf("Processing item: ID=%s, Title=%s, Images=%d", id, naziv, len(slike))

        // Если нет названия, пропускаем этот товар
        if naziv == "" {
            itemsFailed++
            errorLog.WriteString(fmt.Sprintf("Item with ID %s skipped: no title\n", id))
            continue
        }

        // Преобразуем цену в число
        price, err := parsePrice(mpCena)
        if err != nil {
            itemsFailed++
            errorLog.WriteString(fmt.Sprintf("Item with ID %s skipped: invalid price %s: %v\n", id, mpCena, err))
            continue
        }

        // Находим или создаем категорию
        categoryID := DefaultCategoryID
        if kategorija1 != "" {
            catID, err := s.findOrCreateCategory(ctx, kategorija1, kategorija2, kategorija3)
            if err == nil {
                categoryID = catID
            } else {
                errorLog.WriteString(fmt.Sprintf("Warning for item %s: %v. Using default category.\n", id, err))
            }
        }

        // Создаем объявление
        listing := &models.MarketplaceListing{
            UserID:       userID,
            CategoryID:   categoryID,
            StorefrontID: &storefrontID,
            Title:        naziv,
            Description:  opis,
            Price:        price,
            Condition:    "new", // По умолчанию новый товар
            Status: func() string {
                if dostupan == "1" {
                    return "active"
                }
                return "inactive"
            }(),
            ShowOnMap:        false,
            OriginalLanguage: "ru", // Предполагаем русский язык по умолчанию
        }

        // Если товар на акции, отмечаем это в описании
        if naAkciji == "1" {
            listing.Description = "🔥 sale! 🔥\n\n" + listing.Description
        }

        // Создание объявления
        listingID, err := s.storage.CreateListing(ctx, listing)
        if err != nil {
            itemsFailed++
            errorLog.WriteString(fmt.Sprintf("Error creating listing for item %s: %v\n", id, err))
            continue
        }

        // Если есть изображения, обрабатываем их
		if len(slike) > 0 {
			imagesStr := strings.Join(slike, ",")
			// Используем асинхронную обработку изображений
			s.ProcessImportImagesAsync(ctx, listingID, imagesStr, nil)
		}

        // Получаем созданное объявление для индексации
        createdListing, err := s.storage.GetListingByID(ctx, listingID)
        if err != nil {
            errorLog.WriteString(fmt.Sprintf("Warning: Listing created but failed to retrieve for indexing: %v\n", err))
        } else {
            // Индексируем объявление в поисковом движке
            err = s.storage.IndexListing(ctx, createdListing)
            if err != nil {
                errorLog.WriteString(fmt.Sprintf("Warning: Listing created but failed to index: %v\n", err))
            }
        }

        itemsImported++
        // Добавляем лог об успешном импорте
        log.Printf("Successfully imported item %s with ID %d", naziv, listingID)
    }

    // Итоговый лог
    log.Printf("Import completed. Total: %d, Imported: %d, Failed: %d", itemsTotal, itemsImported, itemsFailed)

    return itemsTotal, itemsImported, itemsFailed, nil
}

// extractField извлекает значение поля из XML-элемента
func extractField(xml string, field string) string {
    // Пробуем найти поле с CDATA
    reCDATA := regexp.MustCompile(`<` + field + `><!\[CDATA\[(.*?)\]\]></` + field + `>`)
    matchCDATA := reCDATA.FindStringSubmatch(xml)
    if len(matchCDATA) >= 2 {
        return matchCDATA[1]
    }
    
    // Если не найдено с CDATA, ищем обычное поле
    re := regexp.MustCompile(`<` + field + `>(.*?)</` + field + `>`)
    match := re.FindStringSubmatch(xml)
    if len(match) >= 2 {
        return match[1]
    }
    
    return ""
}
// Улучшенная функция cleanXMLContent с поддержкой безопасных HTML тегов
func cleanXMLContent(content string) string {
    // Удаляем CDATA
    content = regexp.MustCompile(`<!\[CDATA\[(.*?)\]\]>`).ReplaceAllString(content, "$1")

    // Создаем политику безопасных HTML тегов
    p := bluemonday.UGCPolicy()
    
    // Разрешаем базовые теги форматирования текста
    p.AllowElements("b", "i", "u", "strong", "em", "p", "br", "ul", "ol", "li")
    
    // Разрешаем атрибут style для параграфов
    p.AllowAttrs("style").OnElements("p")
    
    // Очищаем HTML от небезопасных тегов и атрибутов
    content = p.Sanitize(content)
    
    // Заменяем множественные пробелы на один
    content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ")

    return strings.TrimSpace(content)
}
// extractImages извлекает ссылки на изображения из элемента artikal
func extractImages(xml string) []string {
    var images []string
    
    // Ищем тег <slike>
    slikeRe := regexp.MustCompile(`<slike>(.*?)</slike>`)
    slikeMatch := slikeRe.FindStringSubmatch(xml)
    
    if len(slikeMatch) >= 2 {
        // Нашли тег <slike>, теперь извлекаем все вложенные теги <slika>
        slikaRe := regexp.MustCompile(`<slika><!\[CDATA\[(.*?)\]\]></slika>`)
        slikaMatches := slikaRe.FindAllStringSubmatch(slikeMatch[1], -1)
        
        // Также пробуем найти теги <slika> без CDATA
        simpleSlikaRe := regexp.MustCompile(`<slika>(.*?)</slika>`)
        simpleSlikaMatches := simpleSlikaRe.FindAllStringSubmatch(slikeMatch[1], -1)
        
        // Добавляем все найденные изображения
        for _, match := range slikaMatches {
            if len(match) >= 2 && match[1] != "" {
                images = append(images, match[1])
            }
        }
        
        for _, match := range simpleSlikaMatches {
            if len(match) >= 2 && match[1] != "" {
                images = append(images, match[1])
            }
        }
    }
    
    // Добавим логирование
    log.Printf("Extracted %d images from XML", len(images))
    
    return images
}
// parsePrice преобразует строку с ценой в число
func parsePrice(priceStr string) (float64, error) {
	// Удаляем все нечисловые символы, кроме точки
	priceStr = regexp.MustCompile(`[^0-9.]`).ReplaceAllString(priceStr, "")
	if priceStr == "" {
		return 0, nil
	}
	return strconv.ParseFloat(priceStr, 64)
}

// findOrCreateCategory находит или создает категорию по имени
func (s *StorefrontService) findOrCreateCategory(ctx context.Context, cat1, cat2, cat3 string) (int, error) {
	// Этот метод должен быть реализован для поиска или создания категорий
	// Для упрощения примера просто возвращаем категорию "Прочее"
	return 9999, nil
}
