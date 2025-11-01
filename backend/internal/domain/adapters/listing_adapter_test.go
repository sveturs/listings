package adapters

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/domain/models"
)

func TestC2CToUnified(t *testing.T) {
	now := time.Now()
	lat := 43.3209
	lon := 21.8954
	externalID := "external-123"

	tests := []struct {
		name     string
		input    *models.MarketplaceListing
		expected *models.UnifiedListing
		wantErr  bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name: "complete c2c listing",
			input: &models.MarketplaceListing{
				ID:               1,
				UserID:           100,
				CategoryID:       10,
				Title:            "Test Listing",
				Description:      "Test Description",
				Price:            99.99,
				Condition:        "used",
				Status:           "active",
				Location:         "Niš, Serbia",
				Latitude:         &lat,
				Longitude:        &lon,
				City:             "Niš",
				Country:          "Serbia",
				ViewsCount:       50,
				ShowOnMap:        true,
				OriginalLanguage: "sr",
				CreatedAt:        now,
				UpdatedAt:        now,
				ExternalID:       externalID,
				Images: []models.MarketplaceImage{
					{
						ID:           1,
						ListingID:    1,
						PublicURL:    "https://example.com/image1.jpg",
						ThumbnailURL: "https://example.com/thumb1.jpg",
						IsMain:       true,
						DisplayOrder: 0,
					},
				},
			},
			expected: &models.UnifiedListing{
				ID:           1,
				SourceType:   "c2c",
				UserID:       100,
				CategoryID:   10,
				Title:        "Test Listing",
				Description:  "Test Description",
				Price:        99.99,
				Condition:    "used",
				Status:       "active",
				Location:     "Niš, Serbia",
				Latitude:     &lat,
				Longitude:    &lon,
				City:         "Niš",
				Country:      "Serbia",
				ViewsCount:   50,
				ShowOnMap:    true,
				OriginalLang: "sr",
				CreatedAt:    now,
				UpdatedAt:    now,
				StorefrontID: nil,
				Metadata: map[string]interface{}{
					"external_id": externalID,
				},
				Images: []models.UnifiedImage{
					{
						ID:           1,
						URL:          "https://example.com/image1.jpg",
						ThumbnailURL: "https://example.com/thumb1.jpg",
						IsMain:       true,
						DisplayOrder: 0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "minimal c2c listing",
			input: &models.MarketplaceListing{
				ID:               2,
				UserID:           200,
				CategoryID:       20,
				Title:            "Minimal",
				Description:      "Desc",
				Price:            10.00,
				Status:           "active",
				ShowOnMap:        false,
				OriginalLanguage: "en",
				CreatedAt:        now,
				UpdatedAt:        now,
			},
			expected: &models.UnifiedListing{
				ID:           2,
				SourceType:   "c2c",
				UserID:       200,
				CategoryID:   20,
				Title:        "Minimal",
				Description:  "Desc",
				Price:        10.00,
				Status:       "active",
				ShowOnMap:    false,
				OriginalLang: "en",
				CreatedAt:    now,
				UpdatedAt:    now,
				StorefrontID: nil,
				Images:       nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := C2CToUnified(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.SourceType, result.SourceType)
			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.Equal(t, tt.expected.CategoryID, result.CategoryID)
			assert.Equal(t, tt.expected.Title, result.Title)
			assert.Equal(t, tt.expected.Description, result.Description)
			assert.Equal(t, tt.expected.Price, result.Price)
			assert.Equal(t, tt.expected.Condition, result.Condition)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.Location, result.Location)
			assert.Equal(t, tt.expected.ShowOnMap, result.ShowOnMap)
			assert.Equal(t, tt.expected.StorefrontID, result.StorefrontID)

			if len(tt.expected.Images) > 0 {
				assert.Len(t, result.Images, len(tt.expected.Images))
				assert.Equal(t, tt.expected.Images[0].URL, result.Images[0].URL)
			}
		})
	}
}

func TestB2CToUnified(t *testing.T) {
	now := time.Now()
	lat := 43.3209
	lon := 21.8954
	sku := "SKU-123"
	barcode := "1234567890"
	individualAddr := "Individual Street 1"

	storefrontLat := 45.2671
	storefrontLon := 19.8335

	storefrontAddr := "Storefront Street 10"
	storefrontCity := "Novi Sad"
	storefrontCountry := "Serbia"
	storefront := &models.Storefront{
		ID:        1,
		UserID:    500,
		Address:   &storefrontAddr,
		City:      &storefrontCity,
		Country:   &storefrontCountry,
		Latitude:  &storefrontLat,
		Longitude: &storefrontLon,
	}

	tests := []struct {
		name       string
		input      *models.StorefrontProduct
		storefront *models.Storefront
		expected   *models.UnifiedListing
		wantErr    bool
	}{
		{
			name:       "nil input",
			input:      nil,
			storefront: nil,
			expected:   nil,
			wantErr:    false,
		},
		{
			name: "b2c product with storefront location",
			input: &models.StorefrontProduct{
				ID:            1,
				StorefrontID:  1,
				Name:          "Test Product",
				Description:   "Product Description",
				Price:         149.99,
				Currency:      "EUR",
				CategoryID:    15,
				SKU:           &sku,
				Barcode:       &barcode,
				StockQuantity: 10,
				StockStatus:   "in_stock",
				IsActive:      true,
				ViewCount:     100,
				SoldCount:     5,
				CreatedAt:     now,
				UpdatedAt:     now,
				ShowOnMap:     true,
				HasVariants:   false,
			},
			storefront: storefront,
			expected: &models.UnifiedListing{
				ID:           1,
				SourceType:   "b2c",
				UserID:       500,
				CategoryID:   15,
				Title:        "Test Product",
				Description:  "Product Description",
				Price:        149.99,
				Condition:    "new",
				Status:       "active",
				Location:     "Storefront Street 10",
				Latitude:     &storefrontLat,
				Longitude:    &storefrontLon,
				City:         "Novi Sad",
				Country:      "Serbia",
				ViewsCount:   100,
				ShowOnMap:    true,
				OriginalLang: "sr",
				CreatedAt:    now,
				UpdatedAt:    now,
				StorefrontID: ptrInt(1),
				Metadata: map[string]interface{}{
					"source":         "storefront",
					"storefront_id":  1,
					"stock_quantity": 10,
					"stock_status":   "in_stock",
					"currency":       "EUR",
					"has_variants":   false,
					"sold_count":     5,
					"sku":            "SKU-123",
					"barcode":        "1234567890",
				},
			},
			wantErr: false,
		},
		{
			name: "b2c product with individual location",
			input: &models.StorefrontProduct{
				ID:                    2,
				StorefrontID:          1,
				Name:                  "Individual Product",
				Description:           "Has its own location",
				Price:                 99.99,
				Currency:              "USD",
				CategoryID:            20,
				StockQuantity:         5,
				StockStatus:           "low_stock",
				IsActive:              true,
				ViewCount:             50,
				CreatedAt:             now,
				UpdatedAt:             now,
				HasIndividualLocation: true,
				IndividualAddress:     &individualAddr,
				IndividualLatitude:    &lat,
				IndividualLongitude:   &lon,
				ShowOnMap:             true,
			},
			storefront: storefront,
			expected: &models.UnifiedListing{
				ID:           2,
				SourceType:   "b2c",
				UserID:       500,
				CategoryID:   20,
				Title:        "Individual Product",
				Description:  "Has its own location",
				Price:        99.99,
				Condition:    "new",
				Status:       "active",
				Location:     "Individual Street 1",
				Latitude:     &lat,
				Longitude:    &lon,
				ViewsCount:   50,
				ShowOnMap:    true,
				OriginalLang: "sr",
				CreatedAt:    now,
				UpdatedAt:    now,
				StorefrontID: ptrInt(1),
				Metadata: map[string]interface{}{
					"source":         "storefront",
					"storefront_id":  1,
					"stock_quantity": 5,
					"stock_status":   "low_stock",
					"currency":       "USD",
					"has_variants":   false,
					"sold_count":     0,
				},
			},
			wantErr: false,
		},
		{
			name: "inactive b2c product",
			input: &models.StorefrontProduct{
				ID:            3,
				StorefrontID:  1,
				Name:          "Inactive Product",
				Description:   "Not active",
				Price:         10.00,
				Currency:      "USD",
				CategoryID:    25,
				StockQuantity: 0,
				StockStatus:   "out_of_stock",
				IsActive:      false,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			storefront: storefront,
			expected: &models.UnifiedListing{
				ID:           3,
				SourceType:   "b2c",
				UserID:       500,
				CategoryID:   25,
				Title:        "Inactive Product",
				Description:  "Not active",
				Price:        10.00,
				Condition:    "new",
				Status:       "inactive",
				Location:     "Storefront Street 10",
				Latitude:     &storefrontLat,
				Longitude:    &storefrontLon,
				City:         "Novi Sad",
				Country:      "Serbia",
				ShowOnMap:    false,
				OriginalLang: "sr",
				CreatedAt:    now,
				UpdatedAt:    now,
				StorefrontID: ptrInt(1),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := B2CToUnified(tt.input, tt.storefront)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.SourceType, result.SourceType)
			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.Equal(t, tt.expected.Title, result.Title)
			assert.Equal(t, tt.expected.Condition, result.Condition)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.Location, result.Location)
			assert.NotNil(t, result.StorefrontID)
			assert.Equal(t, *tt.expected.StorefrontID, *result.StorefrontID)
		})
	}
}

func TestUnifiedToC2C(t *testing.T) {
	now := time.Now()
	lat := 43.3209
	lon := 21.8954

	tests := []struct {
		name     string
		input    *models.UnifiedListing
		expected *models.MarketplaceListing
		wantErr  bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name: "b2c source type returns nil",
			input: &models.UnifiedListing{
				ID:         1,
				SourceType: "b2c",
			},
			expected: nil,
			wantErr:  false,
		},
		{
			name: "complete unified to c2c",
			input: &models.UnifiedListing{
				ID:           1,
				SourceType:   "c2c",
				UserID:       100,
				CategoryID:   10,
				Title:        "Test Listing",
				Description:  "Test Description",
				Price:        99.99,
				Condition:    "used",
				Status:       "active",
				Location:     "Niš, Serbia",
				Latitude:     &lat,
				Longitude:    &lon,
				City:         "Niš",
				Country:      "Serbia",
				ViewsCount:   50,
				ShowOnMap:    true,
				OriginalLang: "sr",
				CreatedAt:    now,
				UpdatedAt:    now,
				Metadata: map[string]interface{}{
					"external_id": "ext-123",
				},
			},
			expected: &models.MarketplaceListing{
				ID:               1,
				UserID:           100,
				CategoryID:       10,
				Title:            "Test Listing",
				Description:      "Test Description",
				Price:            99.99,
				Condition:        "used",
				Status:           "active",
				Location:         "Niš, Serbia",
				Latitude:         &lat,
				Longitude:        &lon,
				City:             "Niš",
				Country:          "Serbia",
				ViewsCount:       50,
				ShowOnMap:        true,
				OriginalLanguage: "sr",
				CreatedAt:        now,
				UpdatedAt:        now,
				ExternalID:       "ext-123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnifiedToC2C(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.Equal(t, tt.expected.Title, result.Title)
			assert.Equal(t, tt.expected.Status, result.Status)
		})
	}
}

func TestUnifiedToB2C(t *testing.T) {
	now := time.Now()
	storefrontID := 1

	tests := []struct {
		name     string
		input    *models.UnifiedListing
		expected *models.StorefrontProduct
		wantErr  bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
			wantErr:  false,
		},
		{
			name: "c2c source type returns nil",
			input: &models.UnifiedListing{
				ID:         1,
				SourceType: "c2c",
			},
			expected: nil,
			wantErr:  false,
		},
		{
			name: "b2c without storefront_id returns nil",
			input: &models.UnifiedListing{
				ID:         1,
				SourceType: "b2c",
			},
			expected: nil,
			wantErr:  false,
		},
		{
			name: "complete unified to b2c",
			input: &models.UnifiedListing{
				ID:           1,
				SourceType:   "b2c",
				UserID:       500,
				CategoryID:   15,
				Title:        "Test Product",
				Description:  "Product Description",
				Price:        149.99,
				Condition:    "new",
				Status:       "active",
				ViewsCount:   100,
				ShowOnMap:    true,
				OriginalLang: "sr",
				CreatedAt:    now,
				UpdatedAt:    now,
				StorefrontID: &storefrontID,
				Metadata: map[string]interface{}{
					"stock_quantity": float64(10),
					"stock_status":   "in_stock",
					"currency":       "EUR",
					"has_variants":   false,
					"sold_count":     float64(5),
					"sku":            "SKU-123",
				},
			},
			expected: &models.StorefrontProduct{
				ID:            1,
				StorefrontID:  1,
				Name:          "Test Product",
				Description:   "Product Description",
				Price:         149.99,
				CategoryID:    15,
				StockQuantity: 10,
				StockStatus:   "in_stock",
				IsActive:      true,
				ViewCount:     100,
				SoldCount:     5,
				CreatedAt:     now,
				UpdatedAt:     now,
				ShowOnMap:     true,
				Currency:      "EUR",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnifiedToB2C(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.expected == nil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.StorefrontID, result.StorefrontID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.IsActive, result.IsActive)
			assert.Equal(t, tt.expected.StockQuantity, result.StockQuantity)
		})
	}
}

func TestBatchC2CToUnified(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    []*models.MarketplaceListing
		expected int
		wantErr  bool
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "empty input",
			input:    []*models.MarketplaceListing{},
			expected: 0,
			wantErr:  false,
		},
		{
			name: "batch conversion",
			input: []*models.MarketplaceListing{
				{
					ID:               1,
					UserID:           100,
					CategoryID:       10,
					Title:            "Listing 1",
					Description:      "Desc 1",
					Price:            10.00,
					Status:           "active",
					OriginalLanguage: "en",
					CreatedAt:        now,
					UpdatedAt:        now,
				},
				{
					ID:               2,
					UserID:           200,
					CategoryID:       20,
					Title:            "Listing 2",
					Description:      "Desc 2",
					Price:            20.00,
					Status:           "active",
					OriginalLanguage: "sr",
					CreatedAt:        now,
					UpdatedAt:        now,
				},
			},
			expected: 2,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BatchC2CToUnified(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.expected == 0 {
				assert.Nil(t, result)
				return
			}

			assert.Len(t, result, tt.expected)
			for i, unified := range result {
				assert.Equal(t, "c2c", unified.SourceType)
				assert.Equal(t, tt.input[i].ID, unified.ID)
			}
		})
	}
}

func TestBatchB2CToUnified(t *testing.T) {
	now := time.Now()
	lat := 45.2671
	lon := 19.8335

	addr := "Street 10"
	city := "Novi Sad"
	country := "Serbia"
	storefronts := map[int]*models.Storefront{
		1: {
			ID:        1,
			UserID:    500,
			Address:   &addr,
			City:      &city,
			Country:   &country,
			Latitude:  &lat,
			Longitude: &lon,
		},
	}

	tests := []struct {
		name        string
		input       []*models.StorefrontProduct
		storefronts map[int]*models.Storefront
		expected    int
		wantErr     bool
	}{
		{
			name:        "nil input",
			input:       nil,
			storefronts: nil,
			expected:    0,
			wantErr:     false,
		},
		{
			name:        "empty input",
			input:       []*models.StorefrontProduct{},
			storefronts: nil,
			expected:    0,
			wantErr:     false,
		},
		{
			name: "batch conversion",
			input: []*models.StorefrontProduct{
				{
					ID:            1,
					StorefrontID:  1,
					Name:          "Product 1",
					Description:   "Desc 1",
					Price:         10.00,
					Currency:      "USD",
					CategoryID:    10,
					StockQuantity: 5,
					IsActive:      true,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
				{
					ID:            2,
					StorefrontID:  1,
					Name:          "Product 2",
					Description:   "Desc 2",
					Price:         20.00,
					Currency:      "EUR",
					CategoryID:    20,
					StockQuantity: 10,
					IsActive:      true,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
			},
			storefronts: storefronts,
			expected:    2,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BatchB2CToUnified(tt.input, tt.storefronts)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			if tt.expected == 0 {
				assert.Nil(t, result)
				return
			}

			assert.Len(t, result, tt.expected)
			for i, unified := range result {
				assert.Equal(t, "b2c", unified.SourceType)
				assert.Equal(t, tt.input[i].ID, unified.ID)
				assert.NotNil(t, unified.StorefrontID)
			}
		})
	}
}

func TestImageConversion(t *testing.T) {
	c2cListing := &models.MarketplaceListing{
		ID:               1,
		UserID:           100,
		CategoryID:       10,
		Title:            "Test",
		Description:      "Desc",
		Price:            10.00,
		Status:           "active",
		OriginalLanguage: "en",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Images: []models.MarketplaceImage{
			{
				ID:           1,
				ListingID:    1,
				PublicURL:    "https://example.com/img1.jpg",
				ThumbnailURL: "https://example.com/thumb1.jpg",
				IsMain:       true,
				DisplayOrder: 0,
			},
			{
				ID:           2,
				ListingID:    1,
				PublicURL:    "https://example.com/img2.jpg",
				ThumbnailURL: "https://example.com/thumb2.jpg",
				IsMain:       false,
				DisplayOrder: 1,
			},
		},
	}

	unified, err := C2CToUnified(c2cListing)
	require.NoError(t, err)
	require.NotNil(t, unified)

	// Check images converted
	assert.Len(t, unified.Images, 2)
	assert.Equal(t, "https://example.com/img1.jpg", unified.Images[0].URL)
	assert.True(t, unified.Images[0].IsMain)

	// Check ImagesJSON populated
	assert.NotNil(t, unified.ImagesJSON)

	var imagesFromJSON []models.UnifiedImage
	err = json.Unmarshal(unified.ImagesJSON, &imagesFromJSON)
	require.NoError(t, err)
	assert.Len(t, imagesFromJSON, 2)
	assert.Equal(t, unified.Images[0].URL, imagesFromJSON[0].URL)
}

// Helper function
func ptrInt(i int) *int {
	return &i
}
