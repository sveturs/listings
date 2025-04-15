// backend/internal/proj/storefront/service/storefront.go
package service

import (
	"archive/zip"
	"backend/internal/domain/models"
	"backend/internal/storage"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
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

	"github.com/microcosm-cc/bluemonday"
)

const (
	StorefrontCreationCost = 15000.0 // стоимость создания витрины
)

type StorefrontService struct {
	storage             storage.Storage
	priceHistoryService PriceHistoryServiceInterface
}

func NewStorefrontService(storage storage.Storage) StorefrontServiceInterface {
	return &StorefrontService{
		storage:             storage,
		priceHistoryService: nil,
	}
}
func (s *StorefrontService) SetPriceHistoryService(priceHistoryService PriceHistoryServiceInterface) {
	s.priceHistoryService = priceHistoryService
}

// GetCategoryMappings возвращает текущие сопоставления категорий для источника импорта
func (s *StorefrontService) GetCategoryMappings(ctx context.Context, sourceID int, userID int) (map[string]int, error) {
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

	// Извлекаем сопоставления из поля Mapping
	var mappings map[string]int
	if len(source.Mapping) > 0 {
		if err := json.Unmarshal(source.Mapping, &mappings); err != nil {
			return nil, fmt.Errorf("error parsing mappings: %w", err)
		}
	} else {
		mappings = make(map[string]int)
	}

	return mappings, nil
}

// UpdateCategoryMappings обновляет сопоставления категорий для источника импорта
func (s *StorefrontService) UpdateCategoryMappings(ctx context.Context, sourceID int, userID int, mappings map[string]int) error {
	// Проверяем права доступа
	source, err := s.storage.GetImportSourceByID(ctx, sourceID)
	if err != nil {
		return err
	}

	// Получаем информацию о витрине
	storefront, err := s.storage.GetStorefrontByID(ctx, source.StorefrontID)
	if err != nil {
		return err
	}

	if storefront.UserID != userID {
		return fmt.Errorf("access denied")
	}

	// Сериализуем сопоставления в JSON
	mappingJSON, err := json.Marshal(mappings)
	if err != nil {
		return fmt.Errorf("error serializing mappings: %w", err)
	}

	// Обновляем поле Mapping
	source.Mapping = mappingJSON

	// Сохраняем изменения
	return s.storage.UpdateImportSource(ctx, source)
}

type PriceHistoryServiceInterface interface {
	AnalyzeDiscount(ctx context.Context, listingID int) (*models.DiscountInfo, error)
	RecordPriceChange(ctx context.Context, listingID int, oldPrice, newPrice float64, source string) error
}
// ApplyCategoryMappings применяет настроенные сопоставления категорий ко всем товарам,
// которые были импортированы из указанного источника
func (s *StorefrontService) ApplyCategoryMappings(ctx context.Context, sourceID int, userID int) (int, error) {
	// Проверяем права доступа
	source, err := s.storage.GetImportSourceByID(ctx, sourceID)
	if err != nil {
		return 0, fmt.Errorf("error getting import source: %w", err)
	}

	// Получаем информацию о витрине
	storefront, err := s.storage.GetStorefrontByID(ctx, source.StorefrontID)
	if err != nil {
		return 0, fmt.Errorf("error getting storefront: %w", err)
	}

	if storefront.UserID != userID {
		return 0, fmt.Errorf("access denied")
	}

	// Получаем сопоставления категорий
	mappings, err := s.GetCategoryMappings(ctx, sourceID, userID)
	if err != nil {
		return 0, fmt.Errorf("error getting category mappings: %w", err)
	}

	if len(mappings) == 0 {
		return 0, fmt.Errorf("no category mappings found")
	}

	// Получаем список всех импортированных товаров для этого источника
	query := `
        WITH source_listings AS (
            SELECT DISTINCT ml.id, ml.title, ml.category_id, ml.external_id
            FROM marketplace_listings ml
            WHERE ml.storefront_id = $1 
              AND ml.external_id IS NOT NULL AND ml.external_id != ''
        ),
        imported_categories_by_source AS (
            SELECT source_category, category_id
            FROM imported_categories
            WHERE source_id = $2 AND category_id > 0
        )
        SELECT sl.id, sl.title, sl.category_id, ic.source_category, ic.category_id as mapped_category_id
        FROM source_listings sl
        JOIN imported_categories ic ON ic.source_category IN (
            SELECT source_category 
            FROM imported_categories 
            WHERE source_id = $2
        )
        WHERE ic.category_id > 0
    `

	rows, err := s.storage.Query(ctx, query, storefront.ID, sourceID)
	if err != nil {
		return 0, fmt.Errorf("error querying listings: %w", err)
	}
	defer rows.Close()

	type listingToUpdate struct {
		ID              int
		Title           string
		CurrentCategory int
		SourceCategory  string
		MappedCategory  int
	}

	var listingsToUpdate []listingToUpdate
	for rows.Next() {
		var listing listingToUpdate
		if err := rows.Scan(&listing.ID, &listing.Title, &listing.CurrentCategory, &listing.SourceCategory, &listing.MappedCategory); err != nil {
			return 0, fmt.Errorf("error scanning listing: %w", err)
		}

		// Проверяем, изменилась ли категория
		if listing.CurrentCategory != listing.MappedCategory {
			listingsToUpdate = append(listingsToUpdate, listing)
		}
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("error iterating listings: %w", err)
	}

	// Обновляем категории товаров
	updatedCount := 0
	for _, listing := range listingsToUpdate {
		categoryID, ok := mappings[listing.SourceCategory]
		if !ok || categoryID == 0 {
			continue // Пропускаем, если нет сопоставления
		}

		// Обновляем категорию товара
		_, err := s.storage.Exec(ctx, `
            UPDATE marketplace_listings
            SET category_id = $1
            WHERE id = $2
        `, categoryID, listing.ID)

		if err != nil {
			log.Printf("Error updating category for listing %d: %v", listing.ID, err)
			continue
		}

		updatedCount++

		// Переиндексируем товар в поисковом движке
		listing, err := s.storage.GetListingByID(ctx, listing.ID)
		if err != nil {
			log.Printf("Error getting listing %d for reindexing: %v", listing.ID, err)
			continue
		}

		if err := s.storage.IndexListing(ctx, listing); err != nil {
			log.Printf("Error reindexing listing %d: %v", listing.ID, err)
		}
	}

	return updatedCount, nil
}
// GetImportedCategories возвращает список категорий, которые были импортированы этим источником
func (s *StorefrontService) GetImportedCategories(ctx context.Context, sourceID int, userID int) ([]string, error) {
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

	// Выполняем запрос к базе данных для получения уникальных категорий из истории импорта
	query := `
        WITH listing_categories AS (
            SELECT DISTINCT h.source_category
            FROM imported_categories h
            WHERE h.source_id = $1
            AND h.source_category IS NOT NULL
            AND h.source_category != ''
        )
        SELECT source_category FROM listing_categories
        ORDER BY source_category ASC
    `

	rows, err := s.storage.Query(ctx, query, sourceID)
	if err != nil {
		return nil, fmt.Errorf("error querying imported categories: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, fmt.Errorf("error scanning category: %w", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
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
	now := time.Now().UTC()
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
func (s *StorefrontService) findOrCreateCategory(ctx context.Context, sourceID int, cat1, cat2, cat3 string) (int, error) {
	var categoryID int = 9999 // По умолчанию категория "Прочее"

	// Получаем источник импорта, чтобы прочитать пользовательские сопоставления
	source, err := s.storage.GetImportSourceByID(ctx, sourceID)
	if err != nil {
		log.Printf("Ошибка получения источника импорта %d: %v", sourceID, err)
		return categoryID, nil // Используем "Прочее" по умолчанию
	}

	// Пытаемся найти сопоставление в пользовательских настройках
	if len(source.Mapping) > 0 {
		var mappings map[string]int
		if err := json.Unmarshal(source.Mapping, &mappings); err == nil {
			// Проверяем все возможные комбинации категорий
			if cat1 != "" && cat2 != "" && cat3 != "" {
				key := fmt.Sprintf("%s|%s|%s", cat1, cat2, cat3)
				if catID, ok := mappings[key]; ok && catID > 0 {
					// Сохраняем использованные категории для будущего использования
					s.saveImportedCategory(ctx, sourceID, key, catID)
					return catID, nil
				}
			}

			if cat1 != "" && cat2 != "" {
				key := fmt.Sprintf("%s|%s", cat1, cat2)
				if catID, ok := mappings[key]; ok && catID > 0 {
					s.saveImportedCategory(ctx, sourceID, key, catID)
					return catID, nil
				}
			}

			if cat1 != "" {
				if catID, ok := mappings[cat1]; ok && catID > 0 {
					s.saveImportedCategory(ctx, sourceID, cat1, catID)
					return catID, nil
				}
			}

			if cat2 != "" {
				if catID, ok := mappings[cat2]; ok && catID > 0 {
					s.saveImportedCategory(ctx, sourceID, cat2, catID)
					return catID, nil
				}
			}

			if cat3 != "" {
				if catID, ok := mappings[cat3]; ok && catID > 0 {
					s.saveImportedCategory(ctx, sourceID, cat3, catID)
					return catID, nil
				}
			}
		} else {
			log.Printf("Ошибка при разборе сопоставлений категорий: %v", err)
		}
	}

	// Если пользовательских сопоставлений нет или они не сработали,
	// используем встроенный маппинг

	// Маппинг для комбинаций категорий (ключ - комбинация категории 1 и 2)
	combinedMapping := map[string]int{
		"OPREMA ZA MOBILNI|ALATI":                  3127, // Запчасти для телефонов
		"OPREMA ZA MOBILNI|SRAFCIGERI":             3127, // Запчасти для телефонов
		"OPREMA ZA MOBILNI|BATERIJE":               3121, // Аккумуляторы
		"OPREMA ZA MOBILNI|PUNJACI ZA TELEFONE":    3123, // Зарядные устройства
		"OPREMA ZA MOBILNI|KUCNI PUNJACI":          3123, // Зарядные устройства
		"OPREMA ZA MOBILNI|ZASTITE ZA EKRANE":      3126, // Чехлы и плёнки
		"OPREMA ZA MOBILNI|ZASTITNA STAKLA OUTLET": 3126, // Чехлы и плёнки
		"OPREMA ZA MOBILNI|DATA KABLOVI":           3124, // Кабели и адаптеры
		"OPREMA ZA MOBILNI|TYPE-C KABLOVI":         3124, // Кабели и адаптеры
		"OPREMA ZA MOBILNI|MASKE OUTLET":           3126, // Чехлы и плёнки
	}

	// Простой маппинг для отдельных категорий
	simpleMapping := map[string]int{
		"OPREMA ZA MOBILNI":        3000,  // Электроника
		"BATERIJE":                 3121,  // Аккумуляторы
		"PUNJACI ZA TELEFONE":      3123,  // Зарядные устройства
		"KUCNI PUNJACI":            3123,  // Зарядные устройства
		"ALATI":                    2835,  // Инструменты (из авто)
		"SRAFCIGERI":               2835,  // Инструменты
		"ZASTITE ZA EKRANE":        3126,  // Чехлы и плёнки
		"ZASTITNA STAKLA OUTLET":   3126,  // Чехлы и плёнки
		"DATA KABLOVI":             3124,  // Кабели и адаптеры
		"TYPE-C KABLOVI":           3124,  // Кабели и адаптеры
		"MASKE OUTLET":             3126,  // Чехлы и плёнки
		"OUTLET":                   -1,    // Игнорируем как категорию
		"AUTO OPREMA OUTLET":       2800,  // Запчасти и аксессуары (авто)
		"SECURITY":                 10000, // Безопасность
		"VIDEO NADZOR":             10100, // Видеонаблюдение
		"KONEKTORI I VIDEO BALUNI": 10410, // Коннекторы и видео баллуны
	}

	// Сохраняем все исходные категории в таблицу imported_categories
	if cat1 != "" {
		s.saveImportedCategory(ctx, sourceID, cat1, 0)
	}
	if cat2 != "" {
		s.saveImportedCategory(ctx, sourceID, cat2, 0)
	}
	if cat3 != "" {
		s.saveImportedCategory(ctx, sourceID, cat3, 0)
	}
	if cat1 != "" && cat2 != "" {
		s.saveImportedCategory(ctx, sourceID, cat1+"|"+cat2, 0)
	}

	// Сначала проверяем комбинацию cat1|cat2
	if cat1 != "" && cat2 != "" {
		combinedKey := cat1 + "|" + cat2
		if catID, ok := combinedMapping[combinedKey]; ok && catID != -1 {
			s.saveImportedCategory(ctx, sourceID, combinedKey, catID)
			return catID, nil
		}
	}

	// Если комбинация не найдена, проверяем kategorija2 как основной уровень
	if cat2 != "" {
		if catID, ok := simpleMapping[cat2]; ok && catID != -1 {
			s.saveImportedCategory(ctx, sourceID, cat2, catID)
			return catID, nil
		}
	}

	// Если нет, проверяем kategorija1
	if cat1 != "" {
		if catID, ok := simpleMapping[cat1]; ok && catID != -1 {
			s.saveImportedCategory(ctx, sourceID, cat1, catID)
			return catID, nil
		}
	}

	// Если ничего не найдено, возвращаем "Прочее"
	log.Printf("Категория не найдена для %s > %s > %s, используется 'Прочее' (9999)", cat1, cat2, cat3)

	// Сохраняем информацию о несопоставленной категории
	combinedCategory := ""
	if cat1 != "" && cat2 != "" && cat3 != "" {
		combinedCategory = fmt.Sprintf("%s|%s|%s", cat1, cat2, cat3)
	} else if cat1 != "" && cat2 != "" {
		combinedCategory = fmt.Sprintf("%s|%s", cat1, cat2)
	} else if cat1 != "" {
		combinedCategory = cat1
	} else if cat2 != "" {
		combinedCategory = cat2
	} else if cat3 != "" {
		combinedCategory = cat3
	}

	if combinedCategory != "" {
		s.saveImportedCategory(ctx, sourceID, combinedCategory, categoryID)
	}

	return categoryID, nil
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

// ReindexStorefrontListings переиндексирует все объявления для конкретной витрины
// Улучшенная функция ReindexStorefrontListings
func (s *StorefrontService) ReindexStorefrontListings(ctx context.Context, storefrontID int) error {
	log.Printf("Переиндексация всех объявлений для витрины %d...", storefrontID)

	// Получаем список всех объявлений для витрины, независимо от статуса
	query := `
        SELECT id FROM marketplace_listings 
        WHERE storefront_id = $1
    `

	rows, err := s.storage.Query(ctx, query, storefrontID)
	if err != nil {
		return fmt.Errorf("ошибка получения списка объявлений: %w", err)
	}
	defer rows.Close()

	var listingIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("Ошибка при сканировании ID объявления: %v", err)
			continue
		}
		listingIDs = append(listingIDs, id)
	}

	log.Printf("Найдено %d объявлений для переиндексации", len(listingIDs))

	// Переиндексируем каждое объявление
	successCount := 0
	failCount := 0

	for i, listingID := range listingIDs {
		if i > 0 && i%100 == 0 {
			// Логируем прогресс каждые 100 объявлений
			log.Printf("Переиндексировано %d из %d объявлений...", i, len(listingIDs))
		}

		listing, err := s.storage.GetListingByID(ctx, listingID)
		if err != nil {
			log.Printf("Ошибка получения объявления %d: %v", listingID, err)
			failCount++
			continue
		}
		if listing.City == "" || listing.Country == "" || listing.Location == "" {
			// Получаем информацию о витрине
			storefront, err := s.storage.GetStorefrontByID(ctx, storefrontID)
			if err == nil && storefront != nil {
				updated := false

				if listing.City == "" && storefront.City != "" {
					listing.City = storefront.City
					updated = true
				}

				if listing.Country == "" && storefront.Country != "" {
					listing.Country = storefront.Country
					updated = true
				}

				if listing.Location == "" && storefront.Address != "" {
					listing.Location = storefront.Address
					updated = true
				}

				if updated {
					log.Printf("Обновлена адресная информация для объявления %d из витрины", listing.ID)

					// Обновляем объявление в базе данных
					if err := s.storage.UpdateListing(ctx, listing); err != nil {
						log.Printf("Ошибка обновления адресной информации для объявления %d: %v", listing.ID, err)
					}
				}
			}
		}
		// Проверка наличия storefront_id
		if listing.StorefrontID == nil || *listing.StorefrontID != storefrontID {
			// Исправляем storefront_id, если он отсутствует или неверен
			listing.StorefrontID = &storefrontID
			log.Printf("Исправлен storefront_id для объявления %d", listingID)

			// Обновляем объявление в базе данных
			if err := s.storage.UpdateListing(ctx, listing); err != nil {
				log.Printf("Ошибка обновления объявления %d в БД: %v", listingID, err)
			}
		}

		if err := s.storage.IndexListing(ctx, listing); err != nil {
			log.Printf("Ошибка индексации объявления %d: %v", listingID, err)
			failCount++
		} else {
			successCount++
		}
	}

	log.Printf("Завершена переиндексация %d объявлений для витрины %d (успешно: %d, ошибок: %d)",
		len(listingIDs), storefrontID, successCount, failCount)
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
	if updatedHistory != nil && (updatedHistory.Status == "success" || updatedHistory.Status == "partial") {
		log.Printf("Запуск переиндексации всех объявлений после импорта...")
		// Запускаем переиндексацию синхронно, а не в горутине
		reindexCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if err := s.ReindexStorefrontListings(reindexCtx, source.StorefrontID); err != nil {
			log.Printf("Ошибка при переиндексации объявлений витрины после импорта: %v", err)
		} else {
			log.Printf("Переиндексация объявлений витрины после импорта успешно завершена")
		}
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
		log.Printf("Ошибка получения информации о витрине для адресов: %v", err)
		// Продолжаем выполнение, даже если не удалось получить витрину
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
		// В методе ImportCSV после обработки категорий
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

		// Сохраняем информацию о категории, если она указана в строке
		if catNameIdx, ok := columnMap["category_name"]; ok && catNameIdx < len(row) {
			categoryName := strings.TrimSpace(row[catNameIdx])
			if categoryName != "" {
				s.saveImportedCategory(ctx, sourceID, categoryName, categoryID)
			}
		}
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
		// После установки связи с витриной
		listingData.UserID = userID
		listingData.StorefrontID = &storefront.ID

		// Применяем данные о местоположении из витрины ВСЕГДА
		if storefront != nil {
			// Для города
			if listingData.City == "" && storefront.City != "" {
				listingData.City = storefront.City
				log.Printf("Применен город витрины для товара: %s", storefront.City)
			}

			// Для страны
			if listingData.Country == "" && storefront.Country != "" {
				listingData.Country = storefront.Country
				log.Printf("Применена страна витрины для товара: %s", storefront.Country)
			}

			// Для адреса
			if listingData.Location == "" && storefront.Address != "" {
				listingData.Location = storefront.Address
				log.Printf("Применен адрес витрины для товара: %s", storefront.Address)
			}

			// Для координат
			if (listingData.Latitude == nil || listingData.Longitude == nil) &&
				storefront.Latitude != nil && storefront.Longitude != nil {
				listingData.Latitude = storefront.Latitude
				listingData.Longitude = storefront.Longitude
				log.Printf("Применены координаты витрины для товара: Lat=%f, Lon=%f",
					*storefront.Latitude, *storefront.Longitude)
			}

			// Включаем показ на карте, если есть координаты
			if listingData.Latitude != nil && listingData.Longitude != nil &&
				*listingData.Latitude != 0 && *listingData.Longitude != 0 {
				listingData.ShowOnMap = true
			}
		}
		// Создание объявления
		listingID, err := s.storage.CreateListing(ctx, &listingData)
		if err != nil {
			itemsFailed++
			errorLog.WriteString(fmt.Sprintf("Error creating listing: %v\n", err))
			continue
		}

		// Если есть колонка с изображениями, обрабатываем их
		if imagesIdx, ok := columnMap["images"]; ok && imagesIdx < len(row) && row[imagesIdx] != "" {
			imagesStr := row[imagesIdx]

			// Запускаем асинхронную обработку изображений без предварительного создания записей
			go func(listID int, imgs string, zipR *zip.Reader) {
				processingCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
				defer cancel()
				s.ProcessImportImagesAsync(processingCtx, listID, imgs, zipR)
				log.Printf("Запущена обработка изображений для листинга %d", listID)
			}(listingID, imagesStr, zipReader)
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
// Обновляем метод ImportXMLFromZip в файле backend/internal/proj/storefront/service/storefront.go

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

	// Используем потоковый парсер XML, передавая sourceID в функцию processXMLContentStream
	itemsTotal, itemsImported, itemsFailed, err = s.processXMLContentStream(ctx, bytes.NewReader(xmlContent), storefront.ID, sourceID, userID, &errorLog)
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

	if history.Status == "success" || history.Status == "partial" {
		log.Printf("Запуск переиндексации всех объявлений после XML импорта...")
		go func() {
			reindexCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
			defer cancel()

			// Запускаем переиндексацию всех объявлений
			if err := s.ReindexStorefrontListings(reindexCtx, source.StorefrontID); err != nil {
				log.Printf("Ошибка при переиндексации объявлений витрины после импорта: %v", err)
			} else {
				log.Printf("Переиндексация объявлений витрины после импорта успешно завершена")
			}
		}()
	}

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
// Модифицированная функция processXMLContentStream, добавляющая поддержку атрибутов недвижимости
func (s *StorefrontService) processXMLContentStream(ctx context.Context, reader io.Reader, storefrontID int, sourceID int, userID int, errorLog *strings.Builder) (int, int, int, error) {
	var itemsTotal, itemsImported, itemsFailed, itemsUpdated int

	log.Printf("Starting streaming XML processing for storefront ID %d, source ID %d", storefrontID, sourceID)
	storefront, err := s.storage.GetStorefrontByID(ctx, storefrontID)
	if err != nil {
		log.Printf("Ошибка получения информации о витрине для адресов: %v", err)
		// Продолжаем выполнение, даже если не удалось получить витрину
	}
	const DefaultCategoryID = 9999

	decoder := xml.NewDecoder(reader)

	var (
		inArtikal   bool
		inField     string
		id          string
		naziv       string
		kategorija1 string
		kategorija2 string
		kategorija3 string
		opis        string
		mpCena      string
		vpCena      string
		dostupan    string
		naAkciji    string
		slike       []string
		inSlike     bool
		
		// Добавляем атрибуты для недвижимости
		propertyType string
		rooms        string
		floor        string
		totalFloors  string
		area         string
		landArea     string
		buildingType string
		hasBalcony   string
		hasElevator  string
		hasParking   string
		yearBuilt    string
		
		// Местоположение
		location      string
		addressCity   string
		addressCountry string
		latitude      string
		longitude     string
		showOnMap     string
	)

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
			if t.Name.Local == "artikal" {
				inArtikal = true
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
				
				// Сбрасываем атрибуты недвижимости
				propertyType = ""
				rooms = ""
				floor = ""
				totalFloors = ""
				area = ""
				landArea = ""
				buildingType = ""
				hasBalcony = ""
				hasElevator = ""
				hasParking = ""
				yearBuilt = ""
				
				// Сбрасываем атрибуты местоположения
				location = ""
				addressCity = ""
				addressCountry = ""
				latitude = ""
				longitude = ""
				showOnMap = ""
			} else if inArtikal {
				if t.Name.Local == "slike" {
					inSlike = true
				} else if inSlike && t.Name.Local == "slika" {
					inField = "slika"
				} else {
					inField = t.Name.Local
				}
			}
		case xml.EndElement:
			if t.Name.Local == "artikal" && inArtikal {
				inArtikal = false
				itemsTotal++

				if naziv == "" {
					itemsFailed++
					errorLog.WriteString(fmt.Sprintf("Item with ID %s skipped: no title\n", id))
					continue
				}

				var price float64 = 0.0
				mpCenaClean := strings.TrimSpace(mpCena)
				if mpCenaClean != "" && mpCenaClean != ".0000" && mpCenaClean != "0.0000" {
					price, err = s.parsePrice(mpCena)
					if err != nil {
						price, err = s.parsePrice(vpCena)
						if err != nil || price == 0 {
							price = 1.00
							log.Printf("For item ID %s: both retail and wholesale prices are invalid, using minimal price: %f", id, price)
						} else {
							price = price * 1.5
							log.Printf("For item ID %s: retail price invalid, using wholesale price with markup: %f", id, price)
						}
					}
				} else {
					wholesalePrice, wpErr := s.parsePrice(vpCena)
					if wpErr == nil && wholesalePrice > 0 {
						price = wholesalePrice * 1.5
						log.Printf("For item ID %s: retail price not set, using wholesale price with markup: %f", id, price)
					} else {
						price = 1.00
						log.Printf("For item ID %s: both retail and wholesale prices are invalid, using minimal price: %f", id, price)
					}
				}

				// Определим категорию на основе категории из XML
				categoryID, err := s.findOrCreateCategory(ctx, sourceID, kategorija1, kategorija2, kategorija3)
				if err != nil {
					errorLog.WriteString(fmt.Sprintf("Warning for item %s: %v. Using default category.\n", id, err))
					categoryID = DefaultCategoryID
				}

				var existingListing *models.MarketplaceListing
				var existingListingID int
				var existingPrice float64

				sqlQuery := `
                    SELECT id, price FROM marketplace_listings 
                    WHERE external_id = $1 AND storefront_id = $2 AND user_id = $3 
                    LIMIT 1
                `
				err = s.storage.QueryRow(ctx, sqlQuery, id, storefrontID, userID).Scan(&existingListingID, &existingPrice)
				if err == nil {
					existingListing, err = s.storage.GetListingByID(ctx, existingListingID)
					if err != nil {
						log.Printf("Failed to fetch existing listing %d: %v", existingListingID, err)
					} else {
						log.Printf("Found existing listing ID %d for external ID %s", existingListingID, id)
					}
				}

				var discountInfo *models.DiscountInfo
				if existingListing != nil && s.priceHistoryService != nil {
					discountInfo, err = s.priceHistoryService.AnalyzeDiscount(ctx, existingListingID)
					if err != nil {
						log.Printf("Ошибка при анализе скидки для товара %s: %v", id, err)
					}
				}

				var discountLabel string
				if existingListing != nil && existingPrice > 0 && price < existingPrice {
					discountPercent := int((existingPrice - price) / existingPrice * 100)
					if discountInfo != nil && discountInfo.IsRealDiscount && discountInfo.DiscountPercent >= 10 {
						discountLabel = fmt.Sprintf("🔥 %d%% SALE! 🔥\nold price: %.2f RSD\n\n",
							discountInfo.DiscountPercent, existingPrice)
						log.Printf("Обнаружена реальная скидка для товара %s: %d%% (old price: %.2f, новая цена: %.2f)",
							id, discountInfo.DiscountPercent, existingPrice, price)
					} else if discountPercent >= 10 {
						discountLabel = fmt.Sprintf("🔥 %d%% SALE! 🔥\nold price: %.2f RSD\n\n",
							discountPercent, existingPrice)
						log.Printf("Обнаружена скидка для товара %s: %d%% (old price: %.2f, новая цена: %.2f)",
							id, discountPercent, existingPrice, price)
					}
				}

				if naAkciji == "1" && discountLabel == "" {
					discountLabel = "🔥 SALE! 🔥\n\n"
				}

				if existingListing != nil && existingPrice != price && s.priceHistoryService != nil {
					err = s.priceHistoryService.RecordPriceChange(ctx, existingListingID, existingPrice, price, "import")
					if err != nil {
						log.Printf("Ошибка при записи изменения цены для товара %s: %v", id, err)
					}
				}
				if kategorija1 != "" || kategorija2 != "" || kategorija3 != "" {
					// Сохраняем все непустые категории
					if kategorija1 != "" {
						s.saveImportedCategory(ctx, sourceID, kategorija1, categoryID)
					}
					if kategorija2 != "" {
						s.saveImportedCategory(ctx, sourceID, kategorija2, categoryID)
					}
					if kategorija3 != "" {
						s.saveImportedCategory(ctx, sourceID, kategorija3, categoryID)
					}

					// Если есть комбинация категорий, сохраняем и ее
					if kategorija1 != "" && kategorija2 != "" {
						combinedKey := kategorija1 + "|" + kategorija2
						s.saveImportedCategory(ctx, sourceID, combinedKey, categoryID)
					}
				}
				descriptionWithDiscount := discountLabel + opis

				if existingListing != nil {
					existingListing.Title = naziv
					existingListing.Description = descriptionWithDiscount
					existingListing.Price = price
					existingListing.CategoryID = categoryID
					if dostupan == "1" {
						existingListing.Status = "active"
					} else {
						existingListing.Status = "inactive"
					}
					
					// Обновляем поля местоположения, если они есть в XML
					if location != "" {
						existingListing.Location = location
					}
					if addressCity != "" {
						existingListing.City = addressCity
					}
					if addressCountry != "" {
						existingListing.Country = addressCountry
					}
					if latitude != "" {
						lat, err := strconv.ParseFloat(latitude, 64)
						if err == nil {
							existingListing.Latitude = &lat
						}
					}
					if longitude != "" {
						lng, err := strconv.ParseFloat(longitude, 64)
						if err == nil {
							existingListing.Longitude = &lng
						}
					}
					if showOnMap == "1" {
						existingListing.ShowOnMap = true
					} else if showOnMap == "0" {
						existingListing.ShowOnMap = false
					}

					if discountInfo != nil && discountInfo.DiscountPercent > 0 {
						discountData := models.DiscountData{
							DiscountPercent: discountInfo.DiscountPercent,
							PreviousPrice:   discountInfo.PreviousPrice,
							EffectiveFrom:   discountInfo.EffectiveFrom,
							HasPriceHistory: true,
						}
						discountJSON, err := json.Marshal(discountData)
						if err == nil {
							if existingListing.Metadata == nil {
								existingListing.Metadata = make(map[string]interface{})
							}
							existingListing.Metadata["discount"] = json.RawMessage(discountJSON)
						}
					}

					if err := s.storage.UpdateListing(ctx, existingListing); err != nil {
						itemsFailed++
						errorLog.WriteString(fmt.Sprintf("Error updating listing %d for item %s: %v\n",
							existingListing.ID, id, err))
						continue
					}

					itemsUpdated++
					log.Printf("Successfully updated listing %d for item %s", existingListing.ID, id)
					
					// Добавляем атрибуты недвижимости, если они есть
					if categoryID >= 1000 && categoryID < 2000 {
						s.savePropertyAttributes(ctx, existingListing.ID, propertyType, rooms, floor, totalFloors, area, landArea, buildingType, hasBalcony, hasElevator, hasParking, yearBuilt)
					}
					
					if len(slike) > 0 {
						imagesStr := strings.Join(slike, ",")
						go func(listID int, imgs string) {
							processingCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
							defer cancel()
							s.ProcessImportImages(processingCtx, listID, imgs, nil)
							log.Printf("Завершена обработка изображений для листинга %d", listID)
						}(existingListing.ID, imagesStr)
					}

					if err := s.storage.IndexListing(ctx, existingListing); err != nil {
						log.Printf("Warning: Failed to reindex updated listing: %v", err)
					}
				} else {
					// Создаем новое объявление
					status := "inactive"
					if dostupan == "1" {
						status = "active"
					}

					// Внутри блока для создания нового объявления
					listing := &models.MarketplaceListing{
						UserID:           userID,
						CategoryID:       categoryID,
						StorefrontID:     &storefrontID,
						Title:            naziv,
						Description:      descriptionWithDiscount,
						Price:            price,
						Condition:        "new",
						Status:           status,
						ShowOnMap:        false,
						OriginalLanguage: "sr",
						ExternalID:       id,
					}

					// Добавляем данные о местоположении из XML
					if location != "" {
						listing.Location = location
					}
					if addressCity != "" {
						listing.City = addressCity
					}
					if addressCountry != "" {
						listing.Country = addressCountry
					}
					if latitude != "" {
						lat, err := strconv.ParseFloat(latitude, 64)
						if err == nil {
							listing.Latitude = &lat
						}
					}
					if longitude != "" {
						lng, err := strconv.ParseFloat(longitude, 64)
						if err == nil {
							listing.Longitude = &lng
						}
					}
					if showOnMap == "1" {
						listing.ShowOnMap = true
					}

					// Применяем данные о местоположении из витрины, если они отсутствуют в XML
					if storefront != nil {
						// Для города
						if listing.City == "" && storefront.City != "" {
							listing.City = storefront.City
							log.Printf("Применен город витрины для товара %s: %s", id, storefront.City)
						}

						// Для страны
						if listing.Country == "" && storefront.Country != "" {
							listing.Country = storefront.Country
							log.Printf("Применена страна витрины для товара %s: %s", id, storefront.Country)
						}

						// Для адреса
						if listing.Location == "" && storefront.Address != "" {
							listing.Location = storefront.Address
							log.Printf("Применен адрес витрины для товара %s: %s", id, storefront.Address)
						}

						// Для координат
						if (listing.Latitude == nil || listing.Longitude == nil) &&
							storefront.Latitude != nil && storefront.Longitude != nil {
							listing.Latitude = storefront.Latitude
							listing.Longitude = storefront.Longitude
							log.Printf("Применены координаты витрины для товара %s: Lat=%f, Lon=%f",
								id, *storefront.Latitude, *storefront.Longitude)
						}

						// Включаем показ на карте, если есть координаты
						if listing.Latitude != nil && listing.Longitude != nil &&
							*listing.Latitude != 0 && *listing.Longitude != 0 {
							listing.ShowOnMap = true
						}
					}
					listingID, err := s.storage.CreateListing(ctx, listing)
					if err != nil {
						itemsFailed++
						errorLog.WriteString(fmt.Sprintf("Error creating listing for item %s: %v\n", id, err))
						continue
					}
					
					// Добавляем атрибуты недвижимости, если категория относится к недвижимости
					if categoryID >= 1000 && categoryID < 2000 {
						s.savePropertyAttributes(ctx, listingID, propertyType, rooms, floor, totalFloors, area, landArea, buildingType, hasBalcony, hasElevator, hasParking, yearBuilt)
					}

					if len(slike) > 0 {
						imagesStr := strings.Join(slike, ",")
						go func(listID int, imgs string) {
							processingCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
							defer cancel()
							s.ProcessImportImagesAsync(processingCtx, listID, imgs, nil)
							log.Printf("Запущена обработка изображений для листинга %d", listID)
						}(listingID, imagesStr)
					}
					createdListing, err := s.storage.GetListingByID(ctx, listingID)
					if err != nil {
						errorLog.WriteString(fmt.Sprintf("Warning: Listing created but failed to retrieve for indexing: %v\n", err))
					} else {
						err = s.storage.IndexListing(ctx, createdListing)
						if err != nil {
							errorLog.WriteString(fmt.Sprintf("Warning: Listing created but failed to index: %v\n", err))
						}
					}

					itemsImported++
					log.Printf("Successfully imported item %s (ID: %s) with DB ID %d", naziv, id, listingID)
					// Получаем обновленное объявление для переиндексации
					updatedListing, err := s.storage.GetListingByID(ctx, listingID)
					if err != nil {
						log.Printf("Warning: Failed to get listing %d for reindexing after translation: %v", listingID, err)
					} else {
						// Переиндексируем объявление
						if err := s.storage.IndexListing(ctx, updatedListing); err != nil {
							log.Printf("Warning: Failed to reindex listing %d after translation: %v", listingID, err)
						} else {
							log.Printf("Successfully reindexed listing %d after translation", listingID)
						}
					}
				}
			} else if t.Name.Local == "slike" {
				inSlike = false
			} else {
				inField = ""
			}
		case xml.CharData:
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
				case "vpCena":
					vpCena = strings.TrimSpace(text)
				case "dostupan":
					dostupan = strings.TrimSpace(text)
				case "naAkciji":
					naAkciji = strings.TrimSpace(text)
				case "slika":
					if text = strings.TrimSpace(text); text != "" {
						slike = append(slike, text)
					}
				
				// Обработка атрибутов недвижимости
				case "property_type":
					propertyType = strings.TrimSpace(text)
				case "rooms":
					rooms = strings.TrimSpace(text)
				case "floor":
					floor = strings.TrimSpace(text)
				case "total_floors":
					totalFloors = strings.TrimSpace(text)
				case "area":
					area = strings.TrimSpace(text)
				case "land_area":
					landArea = strings.TrimSpace(text)
				case "building_type":
					buildingType = strings.TrimSpace(text)
				case "has_balcony":
					hasBalcony = strings.TrimSpace(text)
				case "has_elevator":
					hasElevator = strings.TrimSpace(text)
				case "has_parking":
					hasParking = strings.TrimSpace(text)
				case "year_built":
					yearBuilt = strings.TrimSpace(text)
				
				// Обработка данных о местоположении
				case "location":
					location = strings.TrimSpace(text)
				case "address_city":
					addressCity = strings.TrimSpace(text)
				case "address_country":
					addressCountry = strings.TrimSpace(text)
				case "latitude":
					latitude = strings.TrimSpace(text)
				case "longitude":
					longitude = strings.TrimSpace(text)
				case "show_on_map":
					showOnMap = strings.TrimSpace(text)
				}
			}
		}
	}

	log.Printf("Streaming XML processing completed. Total: %d, Imported: %d, Updated: %d, Failed: %d",
		itemsTotal, itemsImported, itemsUpdated, itemsFailed)

	return itemsTotal, itemsImported + itemsUpdated, itemsFailed, nil
}

// Новая функция для сохранения атрибутов недвижимости

// Улучшенная функция для сохранения атрибутов недвижимости
func (s *StorefrontService) savePropertyAttributes(ctx context.Context, listingID int, propertyType, rooms, floor, totalFloors, area, landArea, buildingType, hasBalcony, hasElevator, hasParking, yearBuilt string) {
    log.Printf("Сохранение атрибутов недвижимости для листинга ID=%d", listingID)
    
    // Получаем атрибуты категории "Недвижимость"
    var attributeValues []models.ListingAttributeValue
    
    // Обрабатываем тип недвижимости
    if propertyType != "" {
        var propTypeAttrID int
        err := s.storage.QueryRow(ctx, `
            SELECT id FROM category_attributes WHERE name = 'property_type' LIMIT 1
        `).Scan(&propTypeAttrID)
        
        if err == nil {
            textVal := propertyType
            attributeValues = append(attributeValues, models.ListingAttributeValue{
                ListingID:    listingID,
                AttributeID:  propTypeAttrID,
                TextValue:    &textVal,
                DisplayValue: textVal,
                AttributeName: "property_type",
                AttributeType: "select",
            })
        }
    }
    
    // Обрабатываем количество комнат
    if rooms != "" {
        var roomsAttrID int
        err := s.storage.QueryRow(ctx, `
            SELECT id FROM category_attributes WHERE name = 'rooms' LIMIT 1
        `).Scan(&roomsAttrID)
        
        if err == nil {
            // Улучшенная валидация
            roomsInt, err := strconv.Atoi(strings.TrimSpace(rooms))
            if err == nil && roomsInt >= 0 && roomsInt < 100 { // Разумное ограничение
                roomsFloat := float64(roomsInt)
                textVal := rooms
                attributeValues = append(attributeValues, models.ListingAttributeValue{
                    ListingID:     listingID,
                    AttributeID:   roomsAttrID,
                    NumericValue:  &roomsFloat,
                    TextValue:     &textVal,
                    DisplayValue:  textVal,
                    AttributeName: "rooms",
                    AttributeType: "number",
                    Unit:          "soba",
                })
            } else {
                log.Printf("Invalid rooms value ignored: %s", rooms)
            }
        }
    }
    
    // Обрабатываем этаж
    if floor != "" {
        var floorAttrID int
        err := s.storage.QueryRow(ctx, `
            SELECT id FROM category_attributes WHERE name = 'floor' LIMIT 1
        `).Scan(&floorAttrID)
        
        if err == nil {
            // Улучшенная валидация
            floorInt, err := strconv.Atoi(strings.TrimSpace(floor))
            if err == nil && floorInt >= -2 && floorInt < 200 { // Разумное ограничение
                floorFloat := float64(floorInt)
                textVal := floor
                attributeValues = append(attributeValues, models.ListingAttributeValue{
                    ListingID:     listingID,
                    AttributeID:   floorAttrID,
                    NumericValue:  &floorFloat,
                    TextValue:     &textVal,
                    DisplayValue:  textVal,
                    AttributeName: "floor",
                    AttributeType: "number",
                    Unit:          "sprat",
                })
            } else {
                log.Printf("Invalid floor value ignored: %s", floor)
            }
        }
    }
    
    // Обрабатываем всего этажей в доме
    if totalFloors != "" {
        var totalFloorsAttrID int
        err := s.storage.QueryRow(ctx, `
            SELECT id FROM category_attributes WHERE name = 'total_floors' LIMIT 1
        `).Scan(&totalFloorsAttrID)
        
        if err == nil {
            // Улучшенная валидация
            floorsInt, err := strconv.Atoi(strings.TrimSpace(totalFloors))
            if err == nil && floorsInt > 0 && floorsInt < 200 { // Разумное ограничение
                floorsFloat := float64(floorsInt)
                textVal := totalFloors
                attributeValues = append(attributeValues, models.ListingAttributeValue{
                    ListingID:     listingID,
                    AttributeID:   totalFloorsAttrID,
                    NumericValue:  &floorsFloat,
                    TextValue:     &textVal,
                    DisplayValue:  textVal,
                    AttributeName: "total_floors",
                    AttributeType: "number",
                    Unit:          "sprat",
                })
            } else {
                log.Printf("Invalid total_floors value ignored: %s", totalFloors)
            }
        }
    }
    
    // Обрабатываем площадь
    if area != "" {
        var areaAttrID int
        err := s.storage.QueryRow(ctx, `
            SELECT id FROM category_attributes WHERE name = 'area' LIMIT 1
        `).Scan(&areaAttrID)
        
        if err == nil {
            // Улучшенная валидация
            areaFloat, err := strconv.ParseFloat(strings.TrimSpace(area), 64)
            if err == nil && areaFloat >= 0 && areaFloat < 10000 { // Разумное ограничение
                textVal := area
                displayVal := fmt.Sprintf("%s м²", area)
                attributeValues = append(attributeValues, models.ListingAttributeValue{
                    ListingID:     listingID,
                    AttributeID:   areaAttrID,
                    NumericValue:  &areaFloat,
                    TextValue:     &textVal,
                    DisplayValue:  displayVal,
                    AttributeName: "area",
                    AttributeType: "number",
                    Unit:          "m²",
                })
            } else {
                log.Printf("Invalid area value ignored: %s", area)
            }
        }
    }
    
    // Обрабатываем площадь участка
    if landArea != "" {
        var landAreaAttrID int
        err := s.storage.QueryRow(ctx, `
            SELECT id FROM category_attributes WHERE name = 'land_area' LIMIT 1
        `).Scan(&landAreaAttrID)
        
        if err == nil {
            // Улучшенная валидация
            landAreaFloat, err := strconv.ParseFloat(strings.TrimSpace(landArea), 64)
            if err == nil && landAreaFloat >= 0 && landAreaFloat < 10000 { // Разумное ограничение
                textVal := landArea
                displayVal := fmt.Sprintf("%s сот.", landArea)
                attributeValues = append(attributeValues, models.ListingAttributeValue{
                    ListingID:     listingID,
                    AttributeID:   landAreaAttrID,
                    NumericValue:  &landAreaFloat,
                    TextValue:     &textVal,
                    DisplayValue:  displayVal,
                    AttributeName: "land_area",
                    AttributeType: "number",
                    Unit:          "ar",
                })
            } else {
                log.Printf("Invalid land_area value ignored: %s", landArea)
            }
        }
    }
    
    // Обрабатываем тип здания
    if buildingType != "" {
        var buildingTypeAttrID int
        err := s.storage.QueryRow(ctx, `
            SELECT id FROM category_attributes WHERE name = 'building_type' LIMIT 1
        `).Scan(&buildingTypeAttrID)
        
        if err == nil {
            // Ограничим длину текста для безопасности
            textVal := buildingType
            if len(textVal) > 100 {
                textVal = textVal[:100]
                log.Printf("Building type value truncated to 100 chars")
            }
            attributeValues = append(attributeValues, models.ListingAttributeValue{
                ListingID:    listingID,
                AttributeID:  buildingTypeAttrID,
                TextValue:    &textVal,
                DisplayValue: textVal,
                AttributeName: "building_type",
                AttributeType: "select",
            })
        }
    }
    
    // Обрабатываем булевы атрибуты (с переработкой отображения и обработки)
    // Добавляем наличие балкона
    if hasBalcony != "" {
        var balconyAttrID int
        err := s.storage.QueryRow(ctx, `
            SELECT id FROM category_attributes WHERE name = 'has_balcony' LIMIT 1
        `).Scan(&balconyAttrID)
        
        if err == nil {
            boolVal := hasBalcony == "1" || strings.ToLower(hasBalcony) == "true"
            var displayVal string
            if boolVal {
                displayVal = "Да"
            } else {
                displayVal = "Нет"
            }
            attributeValues = append(attributeValues, models.ListingAttributeValue{
                ListingID:     listingID,
                AttributeID:   balconyAttrID,
                BooleanValue:  &boolVal,
                DisplayValue:  displayVal,
                AttributeName: "has_balcony",
                AttributeType: "boolean",
            })
        }
    }
    
    // Обрабатываем наличие лифта
    if hasElevator != "" {
        var elevatorAttrID int
        err := s.storage.QueryRow(ctx, `
            SELECT id FROM category_attributes WHERE name = 'has_elevator' LIMIT 1
        `).Scan(&elevatorAttrID)
        
        if err == nil {
            boolVal := hasElevator == "1" || strings.ToLower(hasElevator) == "true"
            textVal := hasElevator
            var displayVal string
            if boolVal {
                displayVal = "Да"
            } else {
                displayVal = "Нет"
            }
            attributeValues = append(attributeValues, models.ListingAttributeValue{
                ListingID:     listingID,
                AttributeID:   elevatorAttrID,
                BooleanValue:  &boolVal,
                TextValue:     &textVal,
                DisplayValue:  displayVal,
                AttributeName: "has_elevator",
                AttributeType: "boolean",
            })
        }
    }
    
    // Обрабатываем наличие парковки
    if hasParking != "" {
        var parkingAttrID int
        err := s.storage.QueryRow(ctx, `
            SELECT id FROM category_attributes WHERE name = 'has_parking' LIMIT 1
        `).Scan(&parkingAttrID)
        
        if err == nil {
            boolVal := hasParking == "1" || strings.ToLower(hasParking) == "true"
            textVal := hasParking
            var displayVal string
            if boolVal {
                displayVal = "Да"
            } else {
                displayVal = "Нет"
            }
            attributeValues = append(attributeValues, models.ListingAttributeValue{
                ListingID:     listingID,
                AttributeID:   parkingAttrID,
                BooleanValue:  &boolVal,
                TextValue:     &textVal,
                DisplayValue:  displayVal,
                AttributeName: "has_parking",
                AttributeType: "boolean",
            })
        }
    }
    
    if yearBuilt != "" {
        var yearBuiltAttrID int
        err := s.storage.QueryRow(ctx, `
            SELECT id FROM category_attributes WHERE name = 'year_built' LIMIT 1
        `).Scan(&yearBuiltAttrID)
        
        if err == nil {
            // Улучшенная валидация года
            yearInt, err := strconv.Atoi(strings.TrimSpace(yearBuilt))
            if err == nil && yearInt >= 1800 && yearInt <= time.Now().Year() + 5 { // Разумное ограничение
                yearFloat := float64(yearInt)
                textVal := yearBuilt
                attributeValues = append(attributeValues, models.ListingAttributeValue{
                    ListingID:     listingID,
                    AttributeID:   yearBuiltAttrID,
                    NumericValue:  &yearFloat,
                    TextValue:     &textVal,
                    DisplayValue:  yearBuilt,
                    AttributeName: "year_built",
                    AttributeType: "number",
                })
            } else {
                log.Printf("Invalid year_built value ignored: %s", yearBuilt)
            }
        }
    }
    
    // Сохраняем все атрибуты в базу данных
    if len(attributeValues) > 0 {
        log.Printf("Сохранение %d атрибутов для объявления ID=%d", len(attributeValues), listingID)
        
        // Очищаем существующие атрибуты для данного объявления
        _, err := s.storage.Exec(ctx, `
            DELETE FROM listing_attribute_values WHERE listing_id = $1
        `, listingID)
        
        if err != nil {
            log.Printf("Ошибка при очистке существующих атрибутов для объявления ID=%d: %v", listingID, err)
        }
        
        // Вставляем новые атрибуты с учетом unit
        for _, attr := range attributeValues {
            var valueText, valueNum, valueBool interface{}
            
            if attr.TextValue != nil {
                valueText = *attr.TextValue
            }
            
            if attr.NumericValue != nil {
                valueNum = *attr.NumericValue
            }
            
            if attr.BooleanValue != nil {
                valueBool = *attr.BooleanValue
            }
            
            _, err := s.storage.Exec(ctx, `
                INSERT INTO listing_attribute_values (
                    listing_id, attribute_id, text_value, numeric_value, boolean_value, json_value, unit
                ) VALUES ($1, $2, $3, $4, $5, NULL, $6)
                ON CONFLICT (listing_id, attribute_id) DO UPDATE SET
                    text_value = $3,
                    numeric_value = $4,
                    boolean_value = $5,
                    unit = $6
            `, listingID, attr.AttributeID, valueText, valueNum, valueBool, attr.Unit)
            
            if err != nil {
                log.Printf("Ошибка при сохранении атрибута ID=%d для объявления ID=%d: %v", 
                    attr.AttributeID, listingID, err)
            } else {
                log.Printf("Успешно сохранен атрибут ID=%d для объявления ID=%d", 
                    attr.AttributeID, listingID)
            }
        }
    } else {
        log.Printf("Нет атрибутов для сохранения для объявления ID=%d", listingID)
    }
}



func (s *StorefrontService) saveImportedCategory(ctx context.Context, sourceID int, sourceCategory string, categoryID int) {
	// Пропускаем пустые категории
	if sourceCategory == "" {
		return
	}

	// Используем UPSERT для обновления или вставки записи
	_, err := s.storage.Exec(ctx, `
        INSERT INTO imported_categories (source_id, source_category, category_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (source_id, source_category) 
        DO UPDATE SET 
            category_id = $3,
            created_at = CURRENT_TIMESTAMP
    `, sourceID, sourceCategory, categoryID)

	if err != nil {
		log.Printf("Ошибка при сохранении информации о категории: %v", err)
	}
}
// processXMLContent обрабатывает содержимое XML и создает товары
func (s *StorefrontService) processXMLContent(ctx context.Context, xmlContent string, storefrontID int, sourceID int, userID int, errorLog *strings.Builder) (int, int, int, error) {
	var itemsTotal, itemsImported, itemsFailed int

	// Добавим логирование для отладки
	log.Printf("Starting XML processing for storefront ID %d, content length: %d bytes", storefrontID, len(xmlContent))

	// Получаем информацию о витрине для использования в данных о местоположении
	storefront, err := s.storage.GetStorefrontByID(ctx, storefrontID)
	if err != nil {
		log.Printf("Ошибка получения информации о витрине для адресов: %v", err)
		// Продолжаем выполнение, даже если не удалось получить витрину
	}

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
		vpCena := extractField(artikal, "vpCena")
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
		price, err := s.parsePrice(mpCena)
		if err != nil {
			// Если не удалось разобрать цену, логируем и устанавливаем 0
			errorLog.WriteString(fmt.Sprintf("Warning for item %s: invalid price %s: %v. Using 0.\n", id, mpCena, err))
			price = 0
		}

		// Пробуем получить оптовую цену, если розничная равна 0
		if price == 0 && vpCena != "" {
			wholesalePrice, _ := s.parsePrice(vpCena)
			if wholesalePrice > 0 {
				// Если есть оптовая цена, используем её с наценкой 100%
				price = wholesalePrice * 2
				errorLog.WriteString(fmt.Sprintf("Warning for item %s: retail price is 0, using wholesale price with markup: %f.\n", id, price))
			}
		}

		// Если после всех попыток цена все равно 0, устанавливаем минимальную цену
		if price == 0 {
			price = 1.00 // Минимальная цена в 1.00
			errorLog.WriteString(fmt.Sprintf("Warning for item %s: both retail and wholesale prices are invalid. Using minimal price of %f.\n", id, price))
		}

		// Находим или создаем категорию
		categoryID := DefaultCategoryID
		if kategorija1 != "" {
			catID, err := s.findOrCreateCategory(ctx, sourceID, kategorija1, kategorija2, kategorija3)
			if err == nil {
				categoryID = catID
			} else {
				errorLog.WriteString(fmt.Sprintf("Warning for item %s: %v. Using default category.\n", id, err))
			}
		}

		// Определяем метку скидки, если товар на акции
		var discountLabel string = ""
		if naAkciji == "1" {
			discountLabel = "🔥 SALE! 🔥\n\n"
		}

		// Определяем статус в зависимости от значения dostupan
		var status string
		if dostupan == "1" {
			status = "active"
		} else {
			status = "inactive"
		}

		// Создаем объявление
		listing := &models.MarketplaceListing{
			UserID:           userID,
			CategoryID:       categoryID,
			StorefrontID:     &storefrontID,
			Title:            naziv,
			Description:      discountLabel + opis,
			Price:            price,
			Condition:        "new", // По умолчанию новый товар
			Status:           status,
			ShowOnMap:        false,
			OriginalLanguage: "sr", // По умолчанию сербский язык
			ExternalID:       id,   // НОВОЕ ПОЛЕ: добавляем внешний ID
		}

		// Применяем данные о местоположении из витрины ВСЕГДА, но не перезаписываем существующие значения
		if storefront != nil {
			// Для города
			if listing.City == "" && storefront.City != "" {
				listing.City = storefront.City
				log.Printf("Применен город витрины для товара %s: %s", id, storefront.City)
			}

			// Для страны
			if listing.Country == "" && storefront.Country != "" {
				listing.Country = storefront.Country
				log.Printf("Применена страна витрины для товара %s: %s", id, storefront.Country)
			}

			// Для адреса
			if listing.Location == "" && storefront.Address != "" {
				listing.Location = storefront.Address
				log.Printf("Применен адрес витрины для товара %s: %s", id, storefront.Address)
			}

			// Для координат
			if (listing.Latitude == nil || listing.Longitude == nil) &&
				storefront.Latitude != nil && storefront.Longitude != nil {
				listing.Latitude = storefront.Latitude
				listing.Longitude = storefront.Longitude
				log.Printf("Применены координаты витрины для товара %s: Lat=%f, Lon=%f",
					id, *storefront.Latitude, *storefront.Longitude)
			}

			// Включаем показ на карте, если есть координаты
			if listing.Latitude != nil && listing.Longitude != nil &&
				*listing.Latitude != 0 && *listing.Longitude != 0 {
				listing.ShowOnMap = true
			}
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
func (s *StorefrontService) parsePrice(priceStr string) (float64, error) {
	priceStr = strings.Replace(priceStr, ",", ".", -1)
	priceStr = strings.TrimSpace(priceStr)
	if priceStr == "" || priceStr == "." {
		return 0, nil
	}
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}
