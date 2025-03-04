package service

import (
	"backend/internal/domain/models"
	"backend/internal/storage"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
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

// RunImport запускает импорт данных
func (s *StorefrontService) RunImport(ctx context.Context, sourceID int, userID int) (*models.ImportHistory, error) {
	// Получаем информацию об источнике
	source, err := s.storage.GetImportSourceByID(ctx, sourceID)
	if err != nil {
		return nil, err
	}

	// Проверяем права доступа
	storefront, err := s.storage.GetStorefrontByID(ctx, source.StorefrontID)
	if err != nil {
		return nil, err
	}

	if storefront.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// В зависимости от типа источника, запускаем соответствующий импорт
	// TODO: Реализовать полноценный импорт

	// Создаем запись в истории импорта
	history := &models.ImportHistory{
		SourceID:  sourceID,
		Status:    "pending",
		StartedAt: time.Now(),
	}

	historyID, err := s.storage.CreateImportHistory(ctx, history)
	if err != nil {
		return nil, err
	}

	history.ID = historyID
	return history, nil
}

// ImportCSV импортирует данные из CSV
func (s *StorefrontService) ImportCSV(ctx context.Context, sourceID int, reader io.Reader, userID int) (*models.ImportHistory, error) {
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
			itemsFailed++
			errorLog.WriteString(fmt.Sprintf("Invalid category_id value '%s': %v\n", row[catIdx], err))
			continue
		}
		listingData.CategoryID = categoryID

		// Получаем condition
		if condIdx, ok := columnMap["condition"]; ok && condIdx < len(row) {
			condition := strings.TrimSpace(row[condIdx])
			if condition != "new" && condition != "used" {
				condition = "new" // По умолчанию новый товар
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
			}
			listingData.Status = status
		} else {
			listingData.Status = "active" // По умолчанию активный товар
		}

		// ДОБАВЛЯЕМ ОБРАБОТКУ МЕСТОПОЛОЖЕНИЯ
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
					errorLog.WriteString(fmt.Sprintf("Invalid latitude value '%s': %v\n", latStr, err))
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
					errorLog.WriteString(fmt.Sprintf("Invalid longitude value '%s': %v\n", lngStr, err))
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
			listingData.OriginalLanguage = "ru" // По умолчанию русский язык
		}

		// Устанавливаем связь с витриной
		listingData.UserID = userID
		listingData.StorefrontID = &storefront.ID
		//	listingData.ShowOnMap = false // Товары из витрины не показываем на карте

		// Создание объявления
		listingID, err := s.storage.CreateListing(ctx, &listingData)
		if err != nil {
			itemsFailed++
			errorLog.WriteString(fmt.Sprintf("Error creating listing: %v\n", err))
			continue
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

		// Если есть колонка с изображениями, обрабатываем их
		if imagesIdx, ok := columnMap["images"]; ok && imagesIdx < len(row) && row[imagesIdx] != "" {
			imagesList := strings.Split(row[imagesIdx], ",")
			for i, imagePath := range imagesList {
				image := &models.MarketplaceImage{
					ListingID:   listingID,
					FilePath:    strings.TrimSpace(imagePath),
					FileName:    strings.TrimSpace(imagePath),
					FileSize:    0,            // Неизвестно
					ContentType: "image/jpeg", // Предполагаем jpeg
					IsMain:      i == 0,       // Первое изображение - основное
				}

				_, err := s.storage.AddListingImage(ctx, image)
				if err != nil {
					errorLog.WriteString(fmt.Sprintf("Error adding image %s to listing %d: %v\n", imagePath, listingID, err))
					// Не увеличиваем itemsFailed, так как само объявление создалось успешно
				}
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
