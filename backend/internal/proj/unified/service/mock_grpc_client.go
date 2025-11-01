// backend/internal/proj/unified/service/mock_grpc_client.go
package service

import (
	"context"
	"fmt"
	"sync"

	"backend/internal/domain/models"
)

// MockListingsGRPCClient - mock implementation для тестирования
// Используется для симуляции listings микросервиса в тестах
type MockListingsGRPCClient struct {
	mu       sync.RWMutex
	listings map[int64]*models.UnifiedListing
	nextID   int64

	// Flags для тестирования различных сценариев
	ShouldFailGet    bool
	ShouldFailCreate bool
	ShouldFailUpdate bool
	ShouldFailDelete bool

	// Счётчики вызовов для проверки в тестах
	GetCallCount    int
	CreateCallCount int
	UpdateCallCount int
	DeleteCallCount int
}

// NewMockListingsGRPCClient создаёт новый mock client
func NewMockListingsGRPCClient() *MockListingsGRPCClient {
	return &MockListingsGRPCClient{
		listings: make(map[int64]*models.UnifiedListing),
		nextID:   1,
	}
}

// GetListing получает listing по ID (mock)
func (m *MockListingsGRPCClient) GetListing(ctx context.Context, id int64) (*models.UnifiedListing, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.GetCallCount++

	if m.ShouldFailGet {
		return nil, fmt.Errorf("mock: get listing failed")
	}

	listing, exists := m.listings[id]
	if !exists {
		return nil, fmt.Errorf("listing not found: id=%d", id)
	}

	// Возвращаем копию чтобы избежать race conditions
	return copyUnifiedListing(listing), nil
}

// CreateListing создаёт новый listing (mock)
func (m *MockListingsGRPCClient) CreateListing(ctx context.Context, unified *models.UnifiedListing) (*models.UnifiedListing, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CreateCallCount++

	if m.ShouldFailCreate {
		return nil, fmt.Errorf("mock: create listing failed")
	}

	// Присваиваем ID
	created := copyUnifiedListing(unified)
	created.ID = int(m.nextID)
	m.nextID++

	// Сохраняем в памяти
	m.listings[int64(created.ID)] = created

	return created, nil
}

// UpdateListing обновляет существующий listing (mock)
func (m *MockListingsGRPCClient) UpdateListing(ctx context.Context, unified *models.UnifiedListing) (*models.UnifiedListing, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.UpdateCallCount++

	if m.ShouldFailUpdate {
		return nil, fmt.Errorf("mock: update listing failed")
	}

	id := int64(unified.ID)
	_, exists := m.listings[id]
	if !exists {
		return nil, fmt.Errorf("listing not found: id=%d", id)
	}

	// Обновляем
	updated := copyUnifiedListing(unified)
	m.listings[id] = updated

	return updated, nil
}

// DeleteListing удаляет listing (mock)
func (m *MockListingsGRPCClient) DeleteListing(ctx context.Context, id int64, userID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DeleteCallCount++

	if m.ShouldFailDelete {
		return fmt.Errorf("mock: delete listing failed")
	}

	listing, exists := m.listings[id]
	if !exists {
		return fmt.Errorf("listing not found: id=%d", id)
	}

	// Проверяем ownership
	if listing.UserID != int(userID) {
		return fmt.Errorf("permission denied: user %d cannot delete listing %d (owner: %d)", userID, id, listing.UserID)
	}

	// Удаляем
	delete(m.listings, id)

	return nil
}

// Reset сбрасывает состояние mock client
func (m *MockListingsGRPCClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.listings = make(map[int64]*models.UnifiedListing)
	m.nextID = 1
	m.ShouldFailGet = false
	m.ShouldFailCreate = false
	m.ShouldFailUpdate = false
	m.ShouldFailDelete = false
	m.GetCallCount = 0
	m.CreateCallCount = 0
	m.UpdateCallCount = 0
	m.DeleteCallCount = 0
}

// GetAllListings возвращает все сохранённые listings (для проверки в тестах)
func (m *MockListingsGRPCClient) GetAllListings() []*models.UnifiedListing {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*models.UnifiedListing, 0, len(m.listings))
	for _, listing := range m.listings {
		result = append(result, copyUnifiedListing(listing))
	}

	return result
}

// copyUnifiedListing создаёт копию UnifiedListing для thread safety
func copyUnifiedListing(src *models.UnifiedListing) *models.UnifiedListing {
	if src == nil {
		return nil
	}

	dst := &models.UnifiedListing{
		ID:           src.ID,
		UserID:       src.UserID,
		SourceType:   src.SourceType,
		Title:        src.Title,
		Description:  src.Description,
		Price:        src.Price,
		CategoryID:   src.CategoryID,
		Status:       src.Status,
		Condition:    src.Condition,
		Location:     src.Location,
		City:         src.City,
		Country:      src.Country,
		ShowOnMap:    src.ShowOnMap,
		ViewsCount:   src.ViewsCount,
		OriginalLang: src.OriginalLang,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
	}

	// Copy optional fields
	if src.StorefrontID != nil {
		storefrontID := *src.StorefrontID
		dst.StorefrontID = &storefrontID
	}

	if src.Latitude != nil {
		lat := *src.Latitude
		dst.Latitude = &lat
	}

	if src.Longitude != nil {
		lon := *src.Longitude
		dst.Longitude = &lon
	}

	// Copy slices
	if src.Images != nil {
		dst.Images = make([]models.UnifiedImage, len(src.Images))
		copy(dst.Images, src.Images)
	}

	// Copy metadata
	if src.Metadata != nil {
		dst.Metadata = make(map[string]interface{})
		for k, v := range src.Metadata {
			dst.Metadata[k] = v
		}
	}

	return dst
}
