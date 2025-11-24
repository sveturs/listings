package postgres

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/tests"
)

func setupAttributeTestRepo(t *testing.T) (*AttributeRepository, *tests.TestDB) {
	t.Helper()

	// Skip if running in short mode or no Docker
	tests.SkipIfShort(t)
	tests.SkipIfNoDocker(t)

	// Setup test database
	testDB := tests.SetupTestPostgres(t)

	// Run migrations
	tests.RunMigrations(t, testDB.DB, "../../../migrations")

	// Create sqlx.DB wrapper
	db := sqlx.NewDb(testDB.DB, "postgres")

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()

	// Create attribute repository
	repo := NewAttributeRepository(db, logger)

	// Create test categories
	setupTestCategories(t, db)

	return repo, testDB
}

// =============================================================================
// Test: NewAttributeRepository
// =============================================================================

func TestNewAttributeRepository(t *testing.T) {
	db := &sqlx.DB{}
	logger := zerolog.Nop()

	repo := NewAttributeRepository(db, logger)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)
}

// =============================================================================
// Test: Create
// =============================================================================

func TestAttributeRepository_Create(t *testing.T) {
	repo, testDB := setupAttributeTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	testCases := []struct {
		name      string
		input     *domain.CreateAttributeInput
		wantErr   bool
		errString string
	}{
		{
			name: "valid attribute - text type",
			input: &domain.CreateAttributeInput{
				Code: "brand_test",
				Name: map[string]string{
					"en": "Brand",
					"ru": "Бренд",
					"sr": "Бренд",
				},
				DisplayName: map[string]string{
					"en": "Brand Name",
					"ru": "Название бренда",
					"sr": "Име бренда",
				},
				AttributeType: domain.AttributeTypeText,
				Purpose:       domain.AttributePurposeRegular,
				IsSearchable:  true,
				IsFilterable:  true,
				IsRequired:    false,
				SortOrder:     1,
			},
			wantErr: false,
		},
		{
			name: "valid attribute - select type with options",
			input: &domain.CreateAttributeInput{
				Code: "size_test",
				Name: map[string]string{
					"en": "Size",
					"ru": "Размер",
					"sr": "Величина",
				},
				DisplayName: map[string]string{
					"en": "Size",
					"ru": "Размер",
					"sr": "Величина",
				},
				AttributeType: domain.AttributeTypeSelect,
				Purpose:       domain.AttributePurposeBoth,
				Options: []domain.AttributeOption{
					{
						Value: "s",
						Label: map[string]string{"en": "Small", "ru": "Маленький", "sr": "Мали"},
					},
					{
						Value: "m",
						Label: map[string]string{"en": "Medium", "ru": "Средний", "sr": "Средњи"},
					},
					{
						Value: "l",
						Label: map[string]string{"en": "Large", "ru": "Большой", "sr": "Велики"},
					},
				},
				ValidationRules: map[string]interface{}{
					"required": true,
				},
				UISettings: map[string]interface{}{
					"display_as": "buttons",
				},
				IsSearchable:        false,
				IsFilterable:        true,
				IsRequired:          true,
				IsVariantCompatible: true,
				SortOrder:           2,
			},
			wantErr: false,
		},
		{
			name: "valid attribute - number type",
			input: &domain.CreateAttributeInput{
				Code: "weight_test",
				Name: map[string]string{
					"en": "Weight",
					"ru": "Вес",
					"sr": "Тежина",
				},
				DisplayName: map[string]string{
					"en": "Weight (kg)",
					"ru": "Вес (кг)",
					"sr": "Тежина (kg)",
				},
				AttributeType: domain.AttributeTypeNumber,
				Purpose:       domain.AttributePurposeRegular,
				ValidationRules: map[string]interface{}{
					"min": 0.0,
					"max": 999.99,
				},
				IsSearchable: false,
				IsFilterable: true,
				SortOrder:    3,
			},
			wantErr: false,
		},
		{
			name:      "nil input",
			input:     nil,
			wantErr:   true,
			errString: "input cannot be nil",
		},
		{
			name: "empty code",
			input: &domain.CreateAttributeInput{
				Code:          "",
				Name:          map[string]string{"en": "Test"},
				DisplayName:   map[string]string{"en": "Test"},
				AttributeType: domain.AttributeTypeText,
			},
			wantErr:   true,
			errString: "code is required",
		},
		{
			name: "empty name",
			input: &domain.CreateAttributeInput{
				Code:          "test_empty_name",
				Name:          map[string]string{},
				DisplayName:   map[string]string{"en": "Test"},
				AttributeType: domain.AttributeTypeText,
			},
			wantErr:   true,
			errString: "name is required",
		},
		{
			name: "empty attribute_type",
			input: &domain.CreateAttributeInput{
				Code:        "test_empty_type",
				Name:        map[string]string{"en": "Test"},
				DisplayName: map[string]string{"en": "Test"},
			},
			wantErr:   true,
			errString: "attribute_type is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attr, err := repo.Create(ctx, tc.input)

			if tc.wantErr {
				require.Error(t, err)
				if tc.errString != "" {
					assert.Contains(t, err.Error(), tc.errString)
				}
				assert.Nil(t, attr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, attr)
				assert.Greater(t, attr.ID, int32(0))
				assert.Equal(t, tc.input.Code, attr.Code)
				assert.Equal(t, tc.input.Name, attr.Name)
				assert.Equal(t, tc.input.DisplayName, attr.DisplayName)
				assert.Equal(t, tc.input.AttributeType, attr.AttributeType)
				assert.True(t, attr.IsActive) // Default is true
				assert.NotZero(t, attr.CreatedAt)
				assert.NotZero(t, attr.UpdatedAt)

				// Check options if provided
				if len(tc.input.Options) > 0 {
					assert.Equal(t, len(tc.input.Options), len(attr.Options))
				}
			}
		})
	}
}

// =============================================================================
// Test: GetByID
// =============================================================================

func TestAttributeRepository_GetByID(t *testing.T) {
	repo, testDB := setupAttributeTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create test attribute
	input := &domain.CreateAttributeInput{
		Code:          "test_get_by_id",
		Name:          map[string]string{"en": "Test"},
		DisplayName:   map[string]string{"en": "Test Attribute"},
		AttributeType: domain.AttributeTypeText,
	}
	created, err := repo.Create(ctx, input)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		id      int32
		wantErr bool
	}{
		{
			name:    "existing attribute",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "non-existent attribute",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attr, err := repo.GetByID(ctx, tc.id)

			if tc.wantErr {
				require.Error(t, err)
				assert.Nil(t, attr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, attr)
				assert.Equal(t, tc.id, attr.ID)
				assert.Equal(t, created.Code, attr.Code)
			}
		})
	}
}

// =============================================================================
// Test: GetByCode
// =============================================================================

func TestAttributeRepository_GetByCode(t *testing.T) {
	repo, testDB := setupAttributeTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create test attribute
	input := &domain.CreateAttributeInput{
		Code:          "test_get_by_code",
		Name:          map[string]string{"en": "Test"},
		DisplayName:   map[string]string{"en": "Test Attribute"},
		AttributeType: domain.AttributeTypeText,
	}
	created, err := repo.Create(ctx, input)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		code    string
		wantErr bool
	}{
		{
			name:    "existing attribute",
			code:    created.Code,
			wantErr: false,
		},
		{
			name:    "non-existent attribute",
			code:    "non_existent_code",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attr, err := repo.GetByCode(ctx, tc.code)

			if tc.wantErr {
				require.Error(t, err)
				assert.Nil(t, attr)
			} else {
				require.NoError(t, err)
				require.NotNil(t, attr)
				assert.Equal(t, tc.code, attr.Code)
				assert.Equal(t, created.ID, attr.ID)
			}
		})
	}
}

// =============================================================================
// Test: Update
// =============================================================================

func TestAttributeRepository_Update(t *testing.T) {
	repo, testDB := setupAttributeTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create test attribute
	input := &domain.CreateAttributeInput{
		Code:          "test_update",
		Name:          map[string]string{"en": "Original Name"},
		DisplayName:   map[string]string{"en": "Original Display Name"},
		AttributeType: domain.AttributeTypeText,
		IsSearchable:  false,
		IsFilterable:  false,
		SortOrder:     1,
	}
	created, err := repo.Create(ctx, input)
	require.NoError(t, err)

	testCases := []struct {
		name      string
		id        int32
		input     *domain.UpdateAttributeInput
		wantErr   bool
		errString string
	}{
		{
			name: "update name",
			id:   created.ID,
			input: &domain.UpdateAttributeInput{
				Name: &map[string]string{
					"en": "Updated Name",
					"ru": "Обновленное имя",
				},
			},
			wantErr: false,
		},
		{
			name: "update boolean flags",
			id:   created.ID,
			input: &domain.UpdateAttributeInput{
				IsSearchable: boolPtr(true),
				IsFilterable: boolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "update sort order",
			id:   created.ID,
			input: &domain.UpdateAttributeInput{
				SortOrder: int32Ptr(10),
			},
			wantErr: false,
		},
		{
			name: "update non-existent attribute",
			id:   99999,
			input: &domain.UpdateAttributeInput{
				Name: &map[string]string{"en": "Test"},
			},
			wantErr:   true,
			errString: "not found",
		},
		{
			name:    "nil input",
			id:      created.ID,
			input:   nil,
			wantErr: true,
		},
		{
			name:    "empty update (no fields)",
			id:      created.ID,
			input:   &domain.UpdateAttributeInput{},
			wantErr: false, // Should return current state
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attr, err := repo.Update(ctx, tc.id, tc.input)

			if tc.wantErr {
				require.Error(t, err)
				if tc.errString != "" {
					assert.Contains(t, err.Error(), tc.errString)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, attr)
				assert.Equal(t, tc.id, attr.ID)

				// Verify updates
				if tc.input != nil {
					if tc.input.Name != nil {
						assert.Equal(t, *tc.input.Name, attr.Name)
					}
					if tc.input.IsSearchable != nil {
						assert.Equal(t, *tc.input.IsSearchable, attr.IsSearchable)
					}
					if tc.input.IsFilterable != nil {
						assert.Equal(t, *tc.input.IsFilterable, attr.IsFilterable)
					}
					if tc.input.SortOrder != nil {
						assert.Equal(t, *tc.input.SortOrder, attr.SortOrder)
					}
				}
			}
		})
	}
}

// =============================================================================
// Test: Delete (Soft Delete)
// =============================================================================

func TestAttributeRepository_Delete(t *testing.T) {
	repo, testDB := setupAttributeTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create test attribute
	input := &domain.CreateAttributeInput{
		Code:          "test_delete",
		Name:          map[string]string{"en": "Test Delete"},
		DisplayName:   map[string]string{"en": "Test Delete"},
		AttributeType: domain.AttributeTypeText,
	}
	created, err := repo.Create(ctx, input)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		id      int32
		wantErr bool
	}{
		{
			name:    "delete existing attribute",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "delete non-existent attribute",
			id:      99999,
			wantErr: true,
		},
		{
			name:    "delete already deleted attribute",
			id:      created.ID,
			wantErr: true, // Already soft-deleted
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Delete(ctx, tc.id)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify attribute is soft-deleted (not retrievable by GetByID)
				attr, err := repo.GetByID(ctx, tc.id)
				assert.Error(t, err)
				assert.Nil(t, attr)
			}
		})
	}
}

// =============================================================================
// Test: List
// =============================================================================

func TestAttributeRepository_List(t *testing.T) {
	repo, testDB := setupAttributeTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create multiple test attributes
	attributes := []struct {
		code         string
		attrType     domain.AttributeType
		purpose      domain.AttributePurpose
		isSearchable bool
		isFilterable bool
		sortOrder    int32
	}{
		{"list_test_1", domain.AttributeTypeText, domain.AttributePurposeRegular, true, true, 1},
		{"list_test_2", domain.AttributeTypeNumber, domain.AttributePurposeVariant, false, true, 2},
		{"list_test_3", domain.AttributeTypeSelect, domain.AttributePurposeBoth, true, false, 3},
	}

	for _, attr := range attributes {
		_, err := repo.Create(ctx, &domain.CreateAttributeInput{
			Code:          attr.code,
			Name:          map[string]string{"en": attr.code},
			DisplayName:   map[string]string{"en": attr.code},
			AttributeType: attr.attrType,
			Purpose:       attr.purpose,
			IsSearchable:  attr.isSearchable,
			IsFilterable:  attr.isFilterable,
			SortOrder:     attr.sortOrder,
		})
		require.NoError(t, err)
	}

	testCases := []struct {
		name          string
		filter        *domain.ListAttributesFilter
		expectedCount int
		wantErr       bool
	}{
		{
			name: "list all attributes",
			filter: &domain.ListAttributesFilter{
				Limit:  10,
				Offset: 0,
			},
			expectedCount: 3,
			wantErr:       false,
		},
		{
			name: "filter by attribute type",
			filter: &domain.ListAttributesFilter{
				AttributeType: attrTypePtr(domain.AttributeTypeText),
				Limit:         10,
				Offset:        0,
			},
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name: "filter by purpose",
			filter: &domain.ListAttributesFilter{
				Purpose: attrPurposePtr(domain.AttributePurposeRegular),
				Limit:   10,
				Offset:  0,
			},
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name: "filter by is_searchable",
			filter: &domain.ListAttributesFilter{
				IsSearchable: boolPtr(true),
				Limit:        10,
				Offset:       0,
			},
			expectedCount: 2,
			wantErr:       false,
		},
		{
			name: "filter by is_filterable",
			filter: &domain.ListAttributesFilter{
				IsFilterable: boolPtr(true),
				Limit:        10,
				Offset:       0,
			},
			expectedCount: 2,
			wantErr:       false,
		},
		{
			name: "pagination - limit 2",
			filter: &domain.ListAttributesFilter{
				Limit:  2,
				Offset: 0,
			},
			expectedCount: 2,
			wantErr:       false,
		},
		{
			name: "pagination - offset 2",
			filter: &domain.ListAttributesFilter{
				Limit:  10,
				Offset: 2,
			},
			expectedCount: 1, // Only 1 remaining after offset 2
			wantErr:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attrs, total, err := repo.List(ctx, tc.filter)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, attrs, tc.expectedCount)
				assert.GreaterOrEqual(t, int(total), tc.expectedCount)

				// Verify ordering by sort_order
				for i := 1; i < len(attrs); i++ {
					assert.GreaterOrEqual(t, attrs[i].SortOrder, attrs[i-1].SortOrder)
				}
			}
		})
	}
}

// =============================================================================
// Test: LinkToCategory
// =============================================================================

func TestAttributeRepository_LinkToCategory(t *testing.T) {
	repo, testDB := setupAttributeTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create test attribute
	attr, err := repo.Create(ctx, &domain.CreateAttributeInput{
		Code:          "link_test",
		Name:          map[string]string{"en": "Link Test"},
		DisplayName:   map[string]string{"en": "Link Test"},
		AttributeType: domain.AttributeTypeText,
		IsSearchable:  false,
		IsFilterable:  true,
	})
	require.NoError(t, err)

	testCases := []struct {
		name        string
		categoryID  int32
		attributeID int32
		settings    *domain.CategoryAttributeSettings
		wantErr     bool
	}{
		{
			name:        "link attribute to category",
			categoryID:  100,
			attributeID: attr.ID,
			settings: &domain.CategoryAttributeSettings{
				IsEnabled:    true,
				IsRequired:   boolPtr(true),
				IsSearchable: boolPtr(true),
				SortOrder:    1,
			},
			wantErr: false,
		},
		{
			name:        "update existing link (upsert)",
			categoryID:  100,
			attributeID: attr.ID,
			settings: &domain.CategoryAttributeSettings{
				IsEnabled:  true,
				IsRequired: boolPtr(false),
				SortOrder:  2,
			},
			wantErr: false,
		},
		{
			name:        "nil settings",
			categoryID:  100,
			attributeID: attr.ID,
			settings:    nil,
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			catAttr, err := repo.LinkToCategory(ctx, tc.categoryID, tc.attributeID, tc.settings)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, catAttr)
				assert.Greater(t, catAttr.ID, int32(0))
				assert.Equal(t, tc.categoryID, catAttr.CategoryID)
				assert.Equal(t, tc.attributeID, catAttr.AttributeID)
				assert.Equal(t, tc.settings.IsEnabled, catAttr.IsEnabled)

				if tc.settings.IsRequired != nil {
					assert.Equal(t, *tc.settings.IsRequired, *catAttr.IsRequired)
				}
			}
		})
	}
}

// =============================================================================
// Test: GetCategoryAttributes
// =============================================================================

func TestAttributeRepository_GetCategoryAttributes(t *testing.T) {
	repo, testDB := setupAttributeTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create test attributes
	attr1, err := repo.Create(ctx, &domain.CreateAttributeInput{
		Code:          "cat_attr_1",
		Name:          map[string]string{"en": "Category Attr 1"},
		DisplayName:   map[string]string{"en": "Category Attr 1"},
		AttributeType: domain.AttributeTypeText,
		IsSearchable:  true,
	})
	require.NoError(t, err)

	attr2, err := repo.Create(ctx, &domain.CreateAttributeInput{
		Code:          "cat_attr_2",
		Name:          map[string]string{"en": "Category Attr 2"},
		DisplayName:   map[string]string{"en": "Category Attr 2"},
		AttributeType: domain.AttributeTypeNumber,
		IsFilterable:  true,
	})
	require.NoError(t, err)

	// Link attributes to category
	_, err = repo.LinkToCategory(ctx, 100, attr1.ID, &domain.CategoryAttributeSettings{
		IsEnabled:  true,
		IsRequired: boolPtr(true),
		SortOrder:  1,
	})
	require.NoError(t, err)

	_, err = repo.LinkToCategory(ctx, 100, attr2.ID, &domain.CategoryAttributeSettings{
		IsEnabled:  true,
		IsRequired: boolPtr(false),
		SortOrder:  2,
	})
	require.NoError(t, err)

	testCases := []struct {
		name          string
		categoryID    int32
		filter        *domain.GetCategoryAttributesFilter
		expectedCount int
		wantErr       bool
	}{
		{
			name:          "get all category attributes",
			categoryID:    100,
			filter:        nil,
			expectedCount: 2,
			wantErr:       false,
		},
		{
			name:       "filter by is_enabled",
			categoryID: 100,
			filter: &domain.GetCategoryAttributesFilter{
				IsEnabled: boolPtr(true),
			},
			expectedCount: 2,
			wantErr:       false,
		},
		{
			name:       "filter by is_required",
			categoryID: 100,
			filter: &domain.GetCategoryAttributesFilter{
				IsRequired: boolPtr(true),
			},
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name:          "category with no attributes",
			categoryID:    200,
			filter:        nil,
			expectedCount: 0,
			wantErr:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attrs, err := repo.GetCategoryAttributes(ctx, tc.categoryID, tc.filter)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Len(t, attrs, tc.expectedCount)

				// Verify attribute loaded
				for _, attr := range attrs {
					assert.NotNil(t, attr.Attribute)
					assert.Greater(t, attr.Attribute.ID, int32(0))
				}
			}
		})
	}
}

// =============================================================================
// Test: UnlinkFromCategory
// =============================================================================

func TestAttributeRepository_UnlinkFromCategory(t *testing.T) {
	repo, testDB := setupAttributeTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create and link test attribute
	attr, err := repo.Create(ctx, &domain.CreateAttributeInput{
		Code:          "unlink_test",
		Name:          map[string]string{"en": "Unlink Test"},
		DisplayName:   map[string]string{"en": "Unlink Test"},
		AttributeType: domain.AttributeTypeText,
	})
	require.NoError(t, err)

	_, err = repo.LinkToCategory(ctx, 100, attr.ID, &domain.CategoryAttributeSettings{
		IsEnabled: true,
		SortOrder: 1,
	})
	require.NoError(t, err)

	testCases := []struct {
		name        string
		categoryID  int32
		attributeID int32
		wantErr     bool
	}{
		{
			name:        "unlink existing link",
			categoryID:  100,
			attributeID: attr.ID,
			wantErr:     false,
		},
		{
			name:        "unlink non-existent link",
			categoryID:  100,
			attributeID: attr.ID,
			wantErr:     true, // Already unlinked
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.UnlinkFromCategory(ctx, tc.categoryID, tc.attributeID)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Verify link is removed
				attrs, err := repo.GetCategoryAttributes(ctx, tc.categoryID, nil)
				require.NoError(t, err)
				assert.Len(t, attrs, 0)
			}
		})
	}
}

// =============================================================================
// Test: SetListingValues & GetListingValues
// =============================================================================

func TestAttributeRepository_SetAndGetListingValues(t *testing.T) {
	repo, testDB := setupAttributeTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create test listing (via direct SQL for simplicity)
	var listingID int32
	err := testDB.DB.QueryRow(`
		INSERT INTO listings (user_id, title, price, currency, category_id, status, sku, source_type, slug)
		VALUES (1, 'Test Listing', 99.99, 'USD', 100, 'active', 'TEST-SKU', 'c2c', 'test-listing')
		RETURNING id
	`).Scan(&listingID)
	require.NoError(t, err)

	// Create test attributes
	textAttr, err := repo.Create(ctx, &domain.CreateAttributeInput{
		Code:          "brand_for_listing",
		Name:          map[string]string{"en": "Brand"},
		DisplayName:   map[string]string{"en": "Brand"},
		AttributeType: domain.AttributeTypeText,
	})
	require.NoError(t, err)

	numberAttr, err := repo.Create(ctx, &domain.CreateAttributeInput{
		Code:          "weight_for_listing",
		Name:          map[string]string{"en": "Weight"},
		DisplayName:   map[string]string{"en": "Weight"},
		AttributeType: domain.AttributeTypeNumber,
	})
	require.NoError(t, err)

	boolAttr, err := repo.Create(ctx, &domain.CreateAttributeInput{
		Code:          "is_new_for_listing",
		Name:          map[string]string{"en": "Is New"},
		DisplayName:   map[string]string{"en": "Is New"},
		AttributeType: domain.AttributeTypeBoolean,
	})
	require.NoError(t, err)

	// Set listing values
	values := []domain.SetListingAttributeValue{
		{
			AttributeID: textAttr.ID,
			ValueText:   stringPtr("Nike"),
		},
		{
			AttributeID: numberAttr.ID,
			ValueNumber: float64Ptr(1.5),
		},
		{
			AttributeID:  boolAttr.ID,
			ValueBoolean: boolPtr(true),
		},
	}

	err = repo.SetListingValues(ctx, listingID, values)
	require.NoError(t, err)

	// Get listing values
	retrievedValues, err := repo.GetListingValues(ctx, listingID)
	require.NoError(t, err)
	assert.Len(t, retrievedValues, 3)

	// Verify values
	for _, val := range retrievedValues {
		switch val.AttributeID {
		case textAttr.ID:
			assert.NotNil(t, val.ValueText)
			assert.Equal(t, "Nike", *val.ValueText)
		case numberAttr.ID:
			assert.NotNil(t, val.ValueNumber)
			assert.Equal(t, 1.5, *val.ValueNumber)
		case boolAttr.ID:
			assert.NotNil(t, val.ValueBoolean)
			assert.True(t, *val.ValueBoolean)
		}

		// Verify attribute loaded
		assert.NotNil(t, val.Attribute)
		assert.Greater(t, val.Attribute.ID, int32(0))
	}

	// Update existing values (upsert)
	updatedValues := []domain.SetListingAttributeValue{
		{
			AttributeID: textAttr.ID,
			ValueText:   stringPtr("Adidas"), // Changed
		},
	}

	err = repo.SetListingValues(ctx, listingID, updatedValues)
	require.NoError(t, err)

	// Verify update
	retrievedValues, err = repo.GetListingValues(ctx, listingID)
	require.NoError(t, err)
	assert.Len(t, retrievedValues, 3) // Still 3 values

	for _, val := range retrievedValues {
		if val.AttributeID == textAttr.ID {
			assert.Equal(t, "Adidas", *val.ValueText)
		}
	}
}

// =============================================================================
// Test: DeleteListingValues
// =============================================================================

func TestAttributeRepository_DeleteListingValues(t *testing.T) {
	repo, testDB := setupAttributeTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create test listing
	var listingID int32
	err := testDB.DB.QueryRow(`
		INSERT INTO listings (user_id, title, price, currency, category_id, status, sku, source_type, slug)
		VALUES (1, 'Test Listing 2', 99.99, 'USD', 100, 'active', 'TEST-SKU-2', 'c2c', 'test-listing-2')
		RETURNING id
	`).Scan(&listingID)
	require.NoError(t, err)

	// Create and set test values
	attr, err := repo.Create(ctx, &domain.CreateAttributeInput{
		Code:          "delete_test_attr",
		Name:          map[string]string{"en": "Delete Test"},
		DisplayName:   map[string]string{"en": "Delete Test"},
		AttributeType: domain.AttributeTypeText,
	})
	require.NoError(t, err)

	err = repo.SetListingValues(ctx, listingID, []domain.SetListingAttributeValue{
		{AttributeID: attr.ID, ValueText: stringPtr("Test")},
	})
	require.NoError(t, err)

	// Delete values
	err = repo.DeleteListingValues(ctx, listingID)
	require.NoError(t, err)

	// Verify deleted
	values, err := repo.GetListingValues(ctx, listingID)
	require.NoError(t, err)
	assert.Len(t, values, 0)
}

// =============================================================================
// Test: GetCategoryVariantAttributes
// =============================================================================

func TestAttributeRepository_GetCategoryVariantAttributes(t *testing.T) {
	repo, testDB := setupAttributeTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := context.Background()

	// Create variant-compatible attribute
	attr, err := repo.Create(ctx, &domain.CreateAttributeInput{
		Code:                "size_variant",
		Name:                map[string]string{"en": "Size"},
		DisplayName:         map[string]string{"en": "Size"},
		AttributeType:       domain.AttributeTypeSelect,
		Purpose:             domain.AttributePurposeVariant,
		IsVariantCompatible: true,
	})
	require.NoError(t, err)

	// Insert category variant attribute directly (since we don't have the full API yet)
	_, err = testDB.DB.Exec(`
		INSERT INTO category_variant_attributes (category_id, attribute_id, is_required, affects_price, affects_stock, sort_order, display_as)
		VALUES ($1, $2, true, false, true, 1, 'buttons')
	`, 100, attr.ID)
	require.NoError(t, err)

	// Get category variant attributes
	variantAttrs, err := repo.GetCategoryVariantAttributes(ctx, 100)
	require.NoError(t, err)
	assert.Len(t, variantAttrs, 1)

	va := variantAttrs[0]
	assert.Equal(t, int32(100), va.CategoryID)
	assert.Equal(t, attr.ID, va.AttributeID)
	assert.True(t, va.IsRequired)
	assert.True(t, va.AffectsStock)
	assert.Equal(t, "buttons", va.DisplayAs)

	// Verify attribute loaded
	assert.NotNil(t, va.Attribute)
	assert.Equal(t, attr.Code, va.Attribute.Code)
}
