package listings

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/service/listings/mocks"
)

func TestSlugGenerator_Generate(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	generator := NewSlugGenerator(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name      string
		title     string
		setupMock func()
		wantSlug  string
		wantErr   bool
	}{
		{
			name:  "unique slug - first attempt",
			title: "My Product Title",
			setupMock: func() {
				// Slug "my-product-title" doesn't exist
				mockRepo.On("GetListingBySlug", ctx, "my-product-title").
					Return(nil, errors.New("not found")).Once()
			},
			wantSlug: "my-product-title",
			wantErr:  false,
		},
		{
			name:  "slug exists - use counter",
			title: "Existing Product",
			setupMock: func() {
				// Base slug exists
				mockRepo.On("GetListingBySlug", ctx, "existing-product").
					Return(&domain.Listing{ID: 1}, nil).Once()
				// First counter attempt also exists
				mockRepo.On("GetListingBySlug", ctx, "existing-product-1").
					Return(&domain.Listing{ID: 2}, nil).Once()
				// Second counter attempt is available
				mockRepo.On("GetListingBySlug", ctx, "existing-product-2").
					Return(nil, errors.New("not found")).Once()
			},
			wantSlug: "existing-product-2",
			wantErr:  false,
		},
		{
			name:      "empty title",
			title:     "",
			setupMock: func() {},
			wantSlug:  "",
			wantErr:   true,
		},
		{
			name:  "cyrillic title",
			title: "Мой продукт",
			setupMock: func() {
				mockRepo.On("GetListingBySlug", ctx, "moi-produkt").
					Return(nil, errors.New("not found")).Once()
			},
			wantSlug: "moi-produkt",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil // Reset mock
			tt.setupMock()

			gotSlug, err := generator.Generate(ctx, tt.title)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantSlug, gotSlug)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSlugGenerator_GenerateWithExclusion(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	generator := NewSlugGenerator(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name             string
		title            string
		excludeListingID int64
		setupMock        func()
		wantSlug         string
		wantErr          bool
	}{
		{
			name:             "slug belongs to excluded listing - can reuse",
			title:            "My Product",
			excludeListingID: 42,
			setupMock: func() {
				// Slug exists but belongs to the listing we're updating
				mockRepo.On("GetListingBySlug", ctx, "my-product").
					Return(&domain.Listing{ID: 42}, nil).Once()
			},
			wantSlug: "my-product",
			wantErr:  false,
		},
		{
			name:             "slug belongs to different listing - use counter",
			title:            "My Product",
			excludeListingID: 42,
			setupMock: func() {
				// Base slug belongs to different listing
				mockRepo.On("GetListingBySlug", ctx, "my-product").
					Return(&domain.Listing{ID: 99}, nil).Once()
				// Counter slug available
				mockRepo.On("GetListingBySlug", ctx, "my-product-1").
					Return(nil, errors.New("not found")).Once()
			},
			wantSlug: "my-product-1",
			wantErr:  false,
		},
		{
			name:             "slug doesn't exist - can use",
			title:            "Unique Product",
			excludeListingID: 42,
			setupMock: func() {
				mockRepo.On("GetListingBySlug", ctx, "unique-product").
					Return(nil, errors.New("not found")).Once()
			},
			wantSlug: "unique-product",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil // Reset mock
			tt.setupMock()

			gotSlug, err := generator.GenerateWithExclusion(ctx, tt.title, tt.excludeListingID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantSlug, gotSlug)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSlugGenerator_ValidateSlug(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	generator := NewSlugGenerator(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name      string
		slug      string
		setupMock func()
		wantErr   bool
	}{
		{
			name: "valid unique slug",
			slug: "valid-slug-123",
			setupMock: func() {
				mockRepo.On("GetListingBySlug", ctx, "valid-slug-123").
					Return(nil, errors.New("not found")).Once()
			},
			wantErr: false,
		},
		{
			name: "slug already exists",
			slug: "existing-slug",
			setupMock: func() {
				mockRepo.On("GetListingBySlug", ctx, "existing-slug").
					Return(&domain.Listing{ID: 1}, nil).Once()
			},
			wantErr: true,
		},
		{
			name:      "invalid slug format - uppercase",
			slug:      "Invalid-Slug",
			setupMock: func() {},
			wantErr:   true,
		},
		{
			name:      "invalid slug format - spaces",
			slug:      "invalid slug",
			setupMock: func() {},
			wantErr:   true,
		},
		{
			name:      "empty slug",
			slug:      "",
			setupMock: func() {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil // Reset mock
			tt.setupMock()

			err := generator.ValidateSlug(ctx, tt.slug)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if len(mockRepo.ExpectedCalls) > 0 {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
