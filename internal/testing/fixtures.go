package testing

import (
	"time"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
)

// TestFixtures holds pre-configured test data for protobuf messages.
// This provides consistent, reusable test data across all integration tests.
type TestFixtures struct {
	// Listings
	BasicListing      *pb.Listing
	PremiumListing    *pb.Listing
	InactiveListing   *pb.Listing
	DraftListing      *pb.Listing
	DeletedListing    *pb.Listing
	ListingWithImages *pb.Listing

	// Categories
	RootCategory     *pb.Category
	ChildCategory    *pb.Category
	CategoryTreeNode *pb.CategoryTreeNode

	// Images
	ListingImage        *pb.ListingImage
	ListingImageRequest *pb.AddImageRequest

	// Products (B2C)
	SimpleProduct       *pb.Product
	ProductWithVariants *pb.Product
	OutOfStockProduct   *pb.Product

	// Product Variants
	SizeVariant  *pb.ProductVariant
	ColorVariant *pb.ProductVariant

	// Favorites
	FavoriteUserIDs []int64

	// Search/List requests
	BasicSearchRequest     *pb.SearchListingsRequest
	PaginatedSearchRequest *pb.SearchListingsRequest
	ListRequest            *pb.ListListingsRequest

	// Common timestamps
	Now       time.Time
	Yesterday time.Time
	Tomorrow  time.Time
}

// NewTestFixtures creates a new set of test fixtures with sensible defaults.
// All fixtures are ready to use and can be modified as needed.
//
// Example:
//
//	fixtures := testing.NewTestFixtures()
//	listing := fixtures.BasicListing
//	listing.Title = "Custom Title"  // Modify as needed
func NewTestFixtures() *TestFixtures {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	fixtures := &TestFixtures{
		Now:       now,
		Yesterday: yesterday,
		Tomorrow:  tomorrow,
	}

	// Initialize all fixtures
	fixtures.initializeListings()
	fixtures.initializeCategories()
	fixtures.initializeImages()
	fixtures.initializeProducts()
	fixtures.initializeVariants()
	fixtures.initializeFavorites()
	fixtures.initializeRequests()

	return fixtures
}

// initializeListings creates test listing fixtures
func (f *TestFixtures) initializeListings() {
	description := "This is a test listing for integration tests"

	// Basic C2C listing
	f.BasicListing = &pb.Listing{
		Id:             1001,
		Uuid:           "uuid-1001",
		UserId:         100,
		Title:          "Test Listing - Basic",
		Description:    &description,
		Price:          99.99,
		Currency:       "USD",
		CategoryId: "1",
		Status:         "active",
		Visibility:     "public",
		Quantity:       1,
		ViewsCount:     10,
		FavoritesCount: 2,
		CreatedAt:      f.Yesterday.Format(time.RFC3339),
		UpdatedAt:      f.Now.Format(time.RFC3339),
		IsDeleted:      false,
		Location: &pb.ListingLocation{
			City:    stringPtr("New York"),
			Country: stringPtr("US"),
		},
		Images: []*pb.ListingImage{
			{
				Id:           2001,
				ListingId:    1001,
				Url:          "https://storage.example.com/listings/1001/image1.jpg",
				DisplayOrder: 1,
				IsPrimary:    true,
				CreatedAt:    f.Yesterday.Format(time.RFC3339),
				UpdatedAt:    f.Now.Format(time.RFC3339),
			},
		},
	}

	premiumDesc := "Premium listing with enhanced features"
	// Premium listing
	f.PremiumListing = &pb.Listing{
		Id:             1002,
		Uuid:           "uuid-1002",
		UserId:         100,
		Title:          "Test Listing - Premium",
		Description:    &premiumDesc,
		Price:          299.99,
		Currency:       "USD",
		CategoryId: "2",
		Status:         "active",
		Visibility:     "public",
		Quantity:       5,
		ViewsCount:     100,
		FavoritesCount: 15,
		CreatedAt:      f.Yesterday.Format(time.RFC3339),
		UpdatedAt:      f.Now.Format(time.RFC3339),
		IsDeleted:      false,
		Location: &pb.ListingLocation{
			City:       stringPtr("San Francisco"),
			Country:    stringPtr("US"),
			PostalCode: stringPtr("94102"),
		},
		Images: []*pb.ListingImage{
			{
				Id:           2002,
				ListingId:    1002,
				Url:          "https://storage.example.com/listings/1002/image1.jpg",
				DisplayOrder: 1,
				IsPrimary:    true,
				CreatedAt:    f.Yesterday.Format(time.RFC3339),
				UpdatedAt:    f.Now.Format(time.RFC3339),
			},
			{
				Id:           2003,
				ListingId:    1002,
				Url:          "https://storage.example.com/listings/1002/image2.jpg",
				DisplayOrder: 2,
				IsPrimary:    false,
				CreatedAt:    f.Yesterday.Format(time.RFC3339),
				UpdatedAt:    f.Now.Format(time.RFC3339),
			},
		},
	}

	inactiveDesc := "This listing is inactive"
	// Inactive listing
	f.InactiveListing = &pb.Listing{
		Id:             1003,
		Uuid:           "uuid-1003",
		UserId:         101,
		Title:          "Test Listing - Inactive",
		Description:    &inactiveDesc,
		Price:          49.99,
		Currency:       "USD",
		CategoryId: "3",
		Status:         "inactive",
		Visibility:     "public",
		Quantity:       1,
		ViewsCount:     5,
		FavoritesCount: 1,
		CreatedAt:      f.Yesterday.Format(time.RFC3339),
		UpdatedAt:      f.Now.Format(time.RFC3339),
		IsDeleted:      false,
	}

	draftDesc := "This listing is in draft status"
	// Draft listing
	f.DraftListing = &pb.Listing{
		Id:             1004,
		Uuid:           "uuid-1004",
		UserId:         100,
		Title:          "Test Listing - Draft",
		Description:    &draftDesc,
		Price:          79.99,
		Currency:       "USD",
		CategoryId: "1",
		Status:         "draft",
		Visibility:     "private",
		Quantity:       1,
		ViewsCount:     0,
		FavoritesCount: 0,
		CreatedAt:      f.Now.Format(time.RFC3339),
		UpdatedAt:      f.Now.Format(time.RFC3339),
		IsDeleted:      false,
	}

	deletedDesc := "This listing has been soft-deleted"
	deletedAt := f.Now.Format(time.RFC3339)
	// Deleted listing (soft delete)
	f.DeletedListing = &pb.Listing{
		Id:             1005,
		Uuid:           "uuid-1005",
		UserId:         102,
		Title:          "Test Listing - Deleted",
		Description:    &deletedDesc,
		Price:          29.99,
		Currency:       "USD",
		CategoryId: "4",
		Status:         "deleted",
		Visibility:     "private",
		Quantity:       0,
		ViewsCount:     3,
		FavoritesCount: 0,
		CreatedAt:      f.Yesterday.Format(time.RFC3339),
		UpdatedAt:      f.Now.Format(time.RFC3339),
		DeletedAt:      &deletedAt,
		IsDeleted:      true,
	}

	multiImgDesc := "This listing has multiple images for testing"
	// Listing with multiple images
	f.ListingWithImages = &pb.Listing{
		Id:             1006,
		Uuid:           "uuid-1006",
		UserId:         100,
		Title:          "Test Listing - Multiple Images",
		Description:    &multiImgDesc,
		Price:          149.99,
		Currency:       "USD",
		CategoryId: "2",
		Status:         "active",
		Visibility:     "public",
		Quantity:       3,
		ViewsCount:     25,
		FavoritesCount: 8,
		CreatedAt:      f.Yesterday.Format(time.RFC3339),
		UpdatedAt:      f.Now.Format(time.RFC3339),
		IsDeleted:      false,
		Images: []*pb.ListingImage{
			{Id: 2010, ListingId: 1006, Url: "https://storage.example.com/listings/1006/img1.jpg", DisplayOrder: 1, IsPrimary: true, CreatedAt: f.Yesterday.Format(time.RFC3339), UpdatedAt: f.Now.Format(time.RFC3339)},
			{Id: 2011, ListingId: 1006, Url: "https://storage.example.com/listings/1006/img2.jpg", DisplayOrder: 2, IsPrimary: false, CreatedAt: f.Yesterday.Format(time.RFC3339), UpdatedAt: f.Now.Format(time.RFC3339)},
			{Id: 2012, ListingId: 1006, Url: "https://storage.example.com/listings/1006/img3.jpg", DisplayOrder: 3, IsPrimary: false, CreatedAt: f.Yesterday.Format(time.RFC3339), UpdatedAt: f.Now.Format(time.RFC3339)},
			{Id: 2013, ListingId: 1006, Url: "https://storage.example.com/listings/1006/img4.jpg", DisplayOrder: 4, IsPrimary: false, CreatedAt: f.Yesterday.Format(time.RFC3339), UpdatedAt: f.Now.Format(time.RFC3339)},
		},
	}
}

// initializeCategories creates test category fixtures
func (f *TestFixtures) initializeCategories() {
	f.RootCategory = &pb.Category{
		Id:           "1",
		Name:         "Electronics",
		Slug:         "electronics",
		Description:  stringPtr("Electronic devices and accessories"),
		Level:        1,
		SortOrder:    1,
		IsActive:     true,
		ListingCount: 100,
		CreatedAt:    f.Yesterday.Format(time.RFC3339),
	}

	parentID := "1"
	f.ChildCategory = &pb.Category{
		Id:           "10",
		Name:         "Smartphones",
		Slug:         "smartphones",
		Description:  stringPtr("Mobile phones and smartphones"),
		ParentId:     &parentID,
		Level:        2,
		SortOrder:    1,
		IsActive:     true,
		ListingCount: 50,
		CreatedAt:    f.Yesterday.Format(time.RFC3339),
	}

	f.CategoryTreeNode = &pb.CategoryTreeNode{
		Id:            "1",
		Name:          "Electronics",
		Slug:          "electronics",
		Level:         1,
		Path:          "1",
		ListingCount:  100,
		ChildrenCount: 5,
		CreatedAt:     f.Yesterday.Format(time.RFC3339),
		HasCustomUi:   false,
		Children: []*pb.CategoryTreeNode{
			{
				Id:            "10",
				Name:          "Smartphones",
				Slug:          "smartphones",
				ParentId:      &parentID,
				Level:         2,
				Path:          "1.10",
				ListingCount:  50,
				ChildrenCount: 0,
				CreatedAt:     f.Yesterday.Format(time.RFC3339),
				HasCustomUi:   false,
			},
		},
	}
}

// initializeImages creates test image fixtures
func (f *TestFixtures) initializeImages() {
	f.ListingImage = &pb.ListingImage{
		Id:           2001,
		ListingId:    1001,
		Url:          "https://storage.example.com/listings/1001/image1.jpg",
		ThumbnailUrl: stringPtr("https://storage.example.com/listings/1001/thumb1.jpg"),
		DisplayOrder: 1,
		IsPrimary:    true,
		CreatedAt:    f.Yesterday.Format(time.RFC3339),
		UpdatedAt:    f.Now.Format(time.RFC3339),
	}

	f.ListingImageRequest = &pb.AddImageRequest{
		ListingId: 1001,
		Url:       "https://storage.example.com/listings/1001/new-image.jpg",
		IsPrimary: false,
	}
}

// initializeProducts creates test product fixtures
func (f *TestFixtures) initializeProducts() {
	// Simple product without variants
	desc1 := "A simple product without variants"
	sku1 := "TEST-SIMPLE-001"
	f.SimpleProduct = &pb.Product{
		Id:            5001,
		StorefrontId:  3001,
		Name:          "Test Product - Simple",
		Description:   desc1,
		Price:         29.99,
		Currency:      "USD",
		CategoryId: "1",
		Sku:           &sku1,
		StockQuantity: 100,
		StockStatus:   "in_stock",
		IsActive:      true,
		HasVariants:   false,
		ViewCount:     50,
		SoldCount:     10,
		CreatedAt:     timestamppb.New(f.Yesterday),
		UpdatedAt:     timestamppb.New(f.Now),
	}

	// Product with variants
	desc2 := "A product with multiple variants"
	sku2 := "TEST-VAR-001"
	attrs, _ := structpb.NewStruct(map[string]interface{}{
		"brand":    "TestBrand",
		"material": "Cotton",
	})

	f.ProductWithVariants = &pb.Product{
		Id:            5002,
		StorefrontId:  3001,
		Name:          "Test Product - With Variants",
		Description:   desc2,
		Price:         49.99,
		Currency:      "USD",
		CategoryId: "2",
		Sku:           &sku2,
		StockQuantity: 200,
		StockStatus:   "in_stock",
		IsActive:      true,
		HasVariants:   true,
		Attributes:    attrs,
		ViewCount:     150,
		SoldCount:     45,
		CreatedAt:     timestamppb.New(f.Yesterday),
		UpdatedAt:     timestamppb.New(f.Now),
	}

	// Out of stock product
	desc3 := "A product that is out of stock"
	sku3 := "TEST-OOS-001"
	f.OutOfStockProduct = &pb.Product{
		Id:            5003,
		StorefrontId:  3001,
		Name:          "Test Product - Out of Stock",
		Description:   desc3,
		Price:         19.99,
		Currency:      "USD",
		CategoryId: "3",
		Sku:           &sku3,
		StockQuantity: 0,
		StockStatus:   "out_of_stock",
		IsActive:      true,
		HasVariants:   false,
		ViewCount:     75,
		SoldCount:     100,
		CreatedAt:     timestamppb.New(f.Yesterday),
		UpdatedAt:     timestamppb.New(f.Now),
	}
}

// initializeVariants creates test product variant fixtures
func (f *TestFixtures) initializeVariants() {
	sizeAttrs, _ := structpb.NewStruct(map[string]interface{}{
		"size": "Large",
	})
	sku1 := "TEST-VAR-001-L"
	price1 := 49.99

	f.SizeVariant = &pb.ProductVariant{
		Id:                6001,
		ProductId:         5002,
		Sku:               &sku1,
		Price:             &price1,
		StockQuantity:     50,
		StockStatus:       "in_stock",
		VariantAttributes: sizeAttrs,
		IsActive:          true,
		IsDefault:         false,
		ViewCount:         25,
		SoldCount:         10,
		CreatedAt:         timestamppb.New(f.Yesterday),
		UpdatedAt:         timestamppb.New(f.Now),
	}

	colorAttrs, _ := structpb.NewStruct(map[string]interface{}{
		"color": "Blue",
		"size":  "Medium",
	})
	sku2 := "TEST-VAR-001-BL-M"
	price2 := 44.99

	f.ColorVariant = &pb.ProductVariant{
		Id:                6002,
		ProductId:         5002,
		Sku:               &sku2,
		Price:             &price2,
		StockQuantity:     75,
		StockStatus:       "in_stock",
		VariantAttributes: colorAttrs,
		IsActive:          true,
		IsDefault:         false,
		ViewCount:         30,
		SoldCount:         15,
		CreatedAt:         timestamppb.New(f.Yesterday),
		UpdatedAt:         timestamppb.New(f.Now),
	}
}

// initializeFavorites creates test favorite fixtures
func (f *TestFixtures) initializeFavorites() {
	f.FavoriteUserIDs = []int64{100, 101, 102, 103, 104}
}

// initializeRequests creates test request fixtures
func (f *TestFixtures) initializeRequests() {
	// Basic search request
	f.BasicSearchRequest = &pb.SearchListingsRequest{
		Query:  "test",
		Limit:  20,
		Offset: 0,
	}

	// Paginated search request with filters
	categoryID := "1"
	minPrice := 10.0
	maxPrice := 100.0

	f.PaginatedSearchRequest = &pb.SearchListingsRequest{
		Query:      "electronics",
		CategoryId: &categoryID,
		MinPrice:   &minPrice,
		MaxPrice:   &maxPrice,
		Limit:      10,
		Offset:     20,
	}

	// List request with filters
	userID := int64(100)
	status := "active"

	f.ListRequest = &pb.ListListingsRequest{
		UserId: &userID,
		Status: &status,
		Limit:  50,
		Offset: 0,
	}
}

// Helper functions to create pointers for optional fields

func stringPtr(s string) *string {
	return &s
}

// CloneFixtures creates a deep copy of fixtures to avoid mutation in tests.
// Use this when you need to modify fixtures without affecting other tests.
//
// Example:
//
//	original := testing.NewTestFixtures()
//	copy := testing.CloneFixtures(original)
//	copy.BasicListing.Title = "Modified"  // Won't affect original
func CloneFixtures(f *TestFixtures) *TestFixtures {
	// Create a simple clone by calling NewTestFixtures again
	// For more sophisticated cloning, implement proto.Clone() for each message
	return NewTestFixtures()
}
