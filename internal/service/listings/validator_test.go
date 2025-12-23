package listings

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/service/listings/mocks"
)

func TestValidator_ValidatePrice(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	validator := NewValidator(mockRepo)

	tests := []struct {
		name    string
		price   float64
		wantErr bool
	}{
		{"valid price", 100.50, false},
		{"zero price", 0.0, false},
		{"negative price", -10.0, true},
		{"max price", MaxPrice, false},
		{"exceeds max price", MaxPrice + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePrice(tt.price)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateTitle(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	validator := NewValidator(mockRepo)

	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"valid title", "Valid Product Title", false},
		{"min length title", "abc", false},
		{"too short title", "ab", true},
		{"empty title", "", true},
		{"whitespace only", "   ", true},
		{"max length title", string(make([]rune, MaxTitleLength)), false},
		{"exceeds max length", string(make([]rune, MaxTitleLength+1)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTitle(tt.title)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateDescription(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	validator := NewValidator(mockRepo)

	validDesc := "This is a valid description"
	longDesc := string(make([]rune, MaxDescriptionLength+1))

	tests := []struct {
		name    string
		desc    *string
		wantErr bool
	}{
		{"nil description", nil, false},
		{"valid description", &validDesc, false},
		{"exceeds max length", &longDesc, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateDescription(tt.desc)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateQuantity(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	validator := NewValidator(mockRepo)

	tests := []struct {
		name     string
		quantity int32
		wantErr  bool
	}{
		{"valid quantity", 10, false},
		{"zero quantity", 0, false},
		{"negative quantity", -5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateQuantity(tt.quantity)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateCurrency(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	validator := NewValidator(mockRepo)

	tests := []struct {
		name     string
		currency string
		wantErr  bool
	}{
		{"valid RSD", "RSD", false},
		{"valid EUR", "EUR", false},
		{"valid USD", "USD", false},
		{"lowercase", "usd", true},
		{"too long", "USDD", true},
		{"too short", "US", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCurrency(tt.currency)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateCategory(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	validator := NewValidator(mockRepo)

	ctx := context.Background()

	// Setup mock for active category
	activeCategory := &domain.Category{
		ID:       "1",
		Name:     "Electronics",
		IsActive: true,
	}
	mockRepo.On("GetCategoryByID", ctx, "1").Return(activeCategory, nil)

	// Setup mock for inactive category
	inactiveCategory := &domain.Category{
		ID:       "2",
		Name:     "Inactive Category",
		IsActive: false,
	}
	mockRepo.On("GetCategoryByID", ctx, "2").Return(inactiveCategory, nil)

	tests := []struct {
		name       string
		categoryID string
		wantErr    bool
	}{
		{"valid active category", "1", false},
		{"inactive category", "2", true},
		{"invalid category ID", "", true},
		{"negative category ID", "-1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCategory(ctx, tt.categoryID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateStatusTransition(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	validator := NewValidator(mockRepo)

	tests := []struct {
		name    string
		from    string
		to      string
		wantErr bool
	}{
		{"draft to active", domain.StatusDraft, domain.StatusActive, false},
		{"draft to inactive", domain.StatusDraft, domain.StatusInactive, false},
		{"active to sold", domain.StatusActive, domain.StatusSold, false},
		{"active to inactive", domain.StatusActive, domain.StatusInactive, false},
		{"inactive to active", domain.StatusInactive, domain.StatusActive, false},
		{"sold to active", domain.StatusSold, domain.StatusActive, false},
		{"invalid: draft to sold", domain.StatusDraft, domain.StatusSold, true},
		{"invalid: sold to inactive", domain.StatusSold, domain.StatusInactive, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateStatusTransition(tt.from, tt.to)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidateImages(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	validator := NewValidator(mockRepo)

	validFileSize := int64(1024 * 1024) // 1MB
	largeFileSize := int64(MaxImageSize + 1)
	validWidth := int32(800)
	validHeight := int32(600)
	smallWidth := int32(50)
	largeWidth := int32(MaxImageWidth + 1)
	validMimeType := "image/jpeg"
	invalidMimeType := "image/gif"

	tests := []struct {
		name    string
		images  []*domain.ListingImage
		wantErr bool
	}{
		{
			name: "valid images",
			images: []*domain.ListingImage{
				{FileSize: &validFileSize, Width: &validWidth, Height: &validHeight, MimeType: &validMimeType},
			},
			wantErr: false,
		},
		{
			name:    "too many images",
			images:  make([]*domain.ListingImage, MaxImagesPerListing+1),
			wantErr: true,
		},
		{
			name: "file size exceeds limit",
			images: []*domain.ListingImage{
				{FileSize: &largeFileSize},
			},
			wantErr: true,
		},
		{
			name: "dimensions too small",
			images: []*domain.ListingImage{
				{Width: &smallWidth, Height: &validHeight},
			},
			wantErr: true,
		},
		{
			name: "dimensions too large",
			images: []*domain.ListingImage{
				{Width: &largeWidth, Height: &validHeight},
			},
			wantErr: true,
		},
		{
			name: "invalid mime type",
			images: []*domain.ListingImage{
				{MimeType: &invalidMimeType},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateImages(tt.images)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSlug(t *testing.T) {
	tests := []struct {
		name string
		slug string
		want bool
	}{
		{"valid slug", "valid-product-slug", true},
		{"valid with numbers", "product-123", true},
		{"uppercase not allowed", "Product-Slug", false},
		{"spaces not allowed", "product slug", false},
		{"empty slug", "", false},
		{"too long", string(make([]rune, 251)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateSlug(tt.slug)
			assert.Equal(t, tt.want, got)
		})
	}
}
