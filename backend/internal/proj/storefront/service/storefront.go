package service

import (
	"backend/internal/domain/models"
	"backend/internal/storage"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
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
		UserID:              userID,
		Name:                create.Name,
		Description:         create.Description,
		Slug:                slug,
		Status:              "active",
		CreationTransactionID: &transactionID,
		CreatedAt:           now,
		UpdatedAt:           now,
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

    // Создаем историю импорта
    history := &models.ImportHistory{
        SourceID:  sourceID,
        Status:    "in_progress",
        StartedAt: time.Now(),
    }

    historyID, err := s.storage.CreateImportHistory(ctx, history)
    if err != nil {
        return nil, err
    }
    history.ID = historyID

    // Читаем CSV файл
    csvReader := csv.NewReader(reader)
    csvReader.Comma = ';'  // Используем точку с запятой как разделитель

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

    // В реальной реализации здесь бы использовались headers для маппинга полей
    _ = headers // явно отмечаем, что мы знаем о неиспользуемой переменной

    // TODO: Реализовать полноценный импорт CSV с учетом маппинга полей
    // В этом примере будем просто считать строки
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
        // TODO: Обработка строки и создание объявления
        // Здесь будет код для преобразования строки CSV в объявление
        _ = row // явно отмечаем, что мы знаем о неиспользуемой переменной
        
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