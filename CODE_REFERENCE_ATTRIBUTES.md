# Listings Microservice - Attributes Code Reference

## Database Schemas - Quick Reference

### 1. Legacy Listing Attributes (Key-Value)
**File**: `/p/github.com/sveturs/listings/migrations/000001_initial_schema.up.sql`

```sql
CREATE TABLE listing_attributes (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    attribute_key VARCHAR(100) NOT NULL,
    attribute_value TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(listing_id, attribute_key)
);
```

### 2. B2C Products with JSONB Attributes
**File**: `/p/github.com/sveturs/listings/migrations/000004_add_b2c_products.up.sql`

Key columns in b2c_products (merged into listings table):
```sql
attributes JSONB DEFAULT '{}'::jsonb,  -- Main product attributes
```

### 3. Product Variants (DEPRECATED)
**File**: `/p/github.com/sveturs/listings/migrations/000004_add_b2c_products.up.sql`

```sql
CREATE TABLE b2c_product_variants (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES b2c_products(id),
    sku VARCHAR(100) UNIQUE,
    barcode VARCHAR(100),
    price NUMERIC(10, 2),
    compare_at_price NUMERIC(10, 2),
    cost_price NUMERIC(10, 2),
    stock_quantity INTEGER DEFAULT 0,
    stock_status VARCHAR(20),
    low_stock_threshold INTEGER,
    
    -- ATTRIBUTES FIELDS (NO LONGER IN USE)
    variant_attributes JSONB DEFAULT '{}'::jsonb,  -- {"size": "L", "color": "Red"}
    dimensions JSONB,                              -- {"length": 10, "width": 5, "height": 3}
    
    is_active BOOLEAN DEFAULT TRUE,
    is_default BOOLEAN DEFAULT FALSE,
    view_count INTEGER DEFAULT 0,
    sold_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**Status**: TABLE DROPPED in migration 000010 (Phase 11.5)

### 4. Categories Table
**File**: `/p/github.com/sveturs/listings/migrations/000011_restore_categories_table.up.sql`

```sql
CREATE TABLE c2c_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    parent_id INTEGER REFERENCES c2c_categories(id),
    icon VARCHAR(50),
    created_at TIMESTAMP,
    has_custom_ui BOOLEAN DEFAULT FALSE,
    custom_ui_component VARCHAR(255),
    sort_order INTEGER DEFAULT 0,
    level INTEGER DEFAULT 0,
    count INTEGER DEFAULT 0,
    external_id VARCHAR(255),
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    seo_title VARCHAR(255),
    seo_description TEXT,
    seo_keywords TEXT
);
```

**Renamed to**: `categories` (migration 20251111230000)

---

## Domain Models - Quick Reference

### Product Model
**File**: `/p/github.com/sveturs/listings/internal/domain/product.go:10`

```go
type Product struct {
    ID                    int64                  `json:"id"`
    StorefrontID          int64                  `json:"storefront_id"`
    Name                  string                 `json:"name"`
    Description           string                 `json:"description"`
    Price                 float64                `json:"price"`
    Currency              string                 `json:"currency"`
    CategoryID            int64                  `json:"category_id"`
    SKU                   *string                `json:"sku,omitempty"`
    Barcode               *string                `json:"barcode,omitempty"`
    StockQuantity         int32                  `json:"stock_quantity"`
    StockStatus           string                 `json:"stock_status"` // in_stock, low_stock, out_of_stock
    IsActive              bool                   `json:"is_active"`
    Attributes            map[string]interface{} `json:"attributes,omitempty"` // ← JSONB
    ViewCount             int32                  `json:"view_count"`
    SoldCount             int32                  `json:"sold_count"`
    CreatedAt             time.Time              `json:"created_at"`
    UpdatedAt             time.Time              `json:"updated_at"`
    HasIndividualLocation bool                   `json:"has_individual_location"`
    IndividualAddress     *string                `json:"individual_address,omitempty"`
    IndividualLatitude    *float64               `json:"individual_latitude,omitempty"`
    IndividualLongitude   *float64               `json:"individual_longitude,omitempty"`
    LocationPrivacy       *string                `json:"location_privacy,omitempty"` // exact, approximate, hidden
    ShowOnMap             bool                   `json:"show_on_map"`
    HasVariants           bool                   `json:"has_variants"`
    
    // Relations
    Variants []ProductVariant `json:"variants,omitempty"`
    Images   []*ProductImage  `json:"images,omitempty"`
}
```

### ProductVariant Model (Deprecated)
**File**: `/p/github.com/sveturs/listings/internal/domain/product.go:41`

```go
type ProductVariant struct {
    ID                int64                  `json:"id"`
    ProductID         int64                  `json:"product_id"`
    SKU               *string                `json:"sku,omitempty"`
    Barcode           *string                `json:"barcode,omitempty"`
    Price             *float64               `json:"price,omitempty"`
    CompareAtPrice    *float64               `json:"compare_at_price,omitempty"`
    CostPrice         *float64               `json:"cost_price,omitempty"`
    StockQuantity     int32                  `json:"stock_quantity"`
    StockStatus       string                 `json:"stock_status"`
    LowStockThreshold *int32                 `json:"low_stock_threshold,omitempty"`
    VariantAttributes map[string]interface{} `json:"variant_attributes,omitempty"` // ← JSONB (NO TABLE)
    Weight            *float64               `json:"weight,omitempty"`
    Dimensions        map[string]interface{} `json:"dimensions,omitempty"` // ← JSONB (NO TABLE)
    IsActive          bool                   `json:"is_active"`
    IsDefault         bool                   `json:"is_default"`
    ViewCount         int32                  `json:"view_count"`
    SoldCount         int32                  `json:"sold_count"`
    CreatedAt         time.Time              `json:"created_at"`
    UpdatedAt         time.Time              `json:"updated_at"`
    
    // Relations
    Images []*ProductImage `json:"images,omitempty"`
}
```

### Listing Model (C2C - Legacy)
**File**: `/p/github.com/sveturs/listings/internal/domain/listing.go:10`

```go
type Listing struct {
    ID             int64
    UUID           string
    Slug           string
    UserID         int64
    StorefrontID   *int64
    Title          string
    Description    *string
    Price          float64
    Currency       string
    CategoryID     int64
    Status         string
    Visibility     string
    Quantity       int32
    SKU            *string
    SourceType     string                   // 'c2c' or 'b2c'
    StockStatus    *string                  // in_stock, out_of_stock, low_stock
    AttributesJSON *string                  // ← JSONB (not mapped)
    ViewsCount     int32
    FavoritesCount int32
    CreatedAt      time.Time
    UpdatedAt      time.Time
    PublishedAt    *time.Time
    DeletedAt      *time.Time
    IsDeleted      bool
    
    // Relations
    Attributes []*ListingAttribute         // ← KEY-VALUE TABLE
    Images     []*ListingImage
    Tags       []string
    Location   *ListingLocation
}

type ListingAttribute struct {
    ID             int64     // Table ID
    ListingID      int64     // Reference to listing
    AttributeKey   string    // Key name (e.g., "brand", "size")
    AttributeValue string    // Value (e.g., "Nike", "Large")
    CreatedAt      time.Time
}
```

---

## Repository Layer - Code Samples

### Reading Product Attributes (JSONB)
**File**: `/p/github.com/sveturs/listings/internal/repository/postgres/products_repository.go`

```go
func (r *Repository) GetProductByID(ctx context.Context, productID int64, storefrontID *int64) (*domain.Product, error) {
    query := `
        SELECT
            p.id, p.storefront_id, p.title, p.description, p.price, p.currency,
            p.category_id, p.sku, p.quantity, p.stock_status,
            p.status, p.attributes, p.view_count, p.sold_count,
            p.created_at, p.updated_at,
            p.has_individual_location, p.individual_address,
            p.individual_latitude, p.individual_longitude,
            p.location_privacy, p.show_on_map, p.has_variants
        FROM listings p
        WHERE p.id = $1
        AND p.source_type = 'b2c'
        AND ($2::bigint IS NULL OR p.storefront_id = $2)
        AND p.deleted_at IS NULL
    `

    var product domain.Product
    var attributesJSON []byte  // ← Will hold JSONB data
    
    err := r.db.QueryRowContext(ctx, query, productID, storefrontID).Scan(
        &product.ID,
        // ... other fields ...
        &attributesJSON,  // ← Scanned here
        // ... rest of fields ...
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to get product by ID: %w", err)
    }
    
    // Parse JSONB attributes
    if len(attributesJSON) > 0 {
        if err := json.Unmarshal(attributesJSON, &product.Attributes); err != nil {
            return nil, fmt.Errorf("failed to unmarshal product attributes: %w", err)
        }
    }
    
    return &product, nil
}
```

### Writing Product Attributes (JSONB)
```go
func (r *Repository) CreateProduct(ctx context.Context, input *domain.CreateProductInput) (*domain.Product, error) {
    // Marshal attributes to JSONB
    var attributesJSON []byte
    if len(input.Attributes) > 0 {
        var err error
        attributesJSON, err = json.Marshal(input.Attributes)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal attributes: %w", err)
        }
    }
    
    query := `
        INSERT INTO listings (
            storefront_id, title, description, price, currency,
            category_id, sku, quantity, status, attributes,
            source_type, has_individual_location,
            individual_address, individual_latitude, individual_longitude,
            location_privacy, show_on_map, has_variants
        ) VALUES (
            $1, $2, $3, $4, $5,
            $6, $7, $8, $9, $10::jsonb,
            $11, $12,
            $13, $14, $15,
            $16, $17, $18
        )
        RETURNING id, created_at, updated_at
    `
    
    var product domain.Product
    product.Attributes = input.Attributes  // ← Store the map
    
    err := r.db.QueryRowContext(ctx, query,
        input.StorefrontID, input.Name, input.Description, input.Price, input.Currency,
        input.CategoryID, input.SKU, input.StockQuantity, "draft", attributesJSON,
        "b2c", input.HasIndividualLocation,
        input.IndividualAddress, input.IndividualLatitude, input.IndividualLongitude,
        input.LocationPrivacy, input.ShowOnMap, input.HasVariants,
    ).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
    
    return &product, err
}
```

### Deprecated Variants Code (DO NOT USE)
**File**: `/p/github.com/sveturs/listings/internal/repository/postgres/product_variants_repository.go:1-6`

```go
package postgres

// ⚠️ DEPRECATED: Product variants functionality is deprecated.
// The b2c_product_variants table was removed in Phase 11.5 (migration 000010).
// This file is kept for API compatibility but methods will return errors.
// Consider creating a new unified product_variants table if variants are needed.
```

---

## gRPC Definitions - Proto Messages

### Product Message (Accepts Attributes)
**File**: `/p/github.com/sveturs/listings/api/proto/listings/v1/listings.proto:494-521`

```protobuf
message Product {
  int64 id = 1;
  int64 storefront_id = 2;
  string name = 3;
  string description = 4;
  double price = 5;
  string currency = 6;
  int64 category_id = 7;
  optional string sku = 8;
  optional string barcode = 9;
  int32 stock_quantity = 10;
  string stock_status = 11;
  bool is_active = 12;
  google.protobuf.Struct attributes = 13;  // ← FLEXIBLE ATTRIBUTES
  int32 view_count = 14;
  int32 sold_count = 15;
  google.protobuf.Timestamp created_at = 16;
  google.protobuf.Timestamp updated_at = 17;
  
  bool has_individual_location = 18;
  optional string individual_address = 19;
  optional double individual_latitude = 20;
  optional double individual_longitude = 21;
  optional string location_privacy = 22;
  bool show_on_map = 23;
  bool has_variants = 24;
  repeated ProductVariant variants = 25;
}
```

### ProductVariant Message (ATTRIBUTES STORED BUT TABLE DROPPED)
**File**: `/p/github.com/sveturs/listings/api/proto/listings/v1/listings.proto:523-544`

```protobuf
message ProductVariant {
  int64 id = 1;
  int64 product_id = 2;
  optional string sku = 3;
  optional string barcode = 4;
  optional double price = 5;
  optional double compare_at_price = 6;
  optional double cost_price = 7;
  int32 stock_quantity = 8;
  string stock_status = 9;
  optional int32 low_stock_threshold = 10;
  google.protobuf.Struct variant_attributes = 11;  // ← NO TABLE!
  optional double weight = 12;
  google.protobuf.Struct dimensions = 13;          // ← NO TABLE!
  bool is_active = 14;
  bool is_default = 15;
  int32 view_count = 16;
  int32 sold_count = 17;
  google.protobuf.Timestamp created_at = 18;
  google.protobuf.Timestamp updated_at = 19;
}
```

### ListingAttribute Message (Key-Value)
**File**: `/p/github.com/sveturs/listings/api/proto/listings/v1/listings.proto:420-427`

```protobuf
message ListingAttribute {
  int64 id = 1;
  int64 listing_id = 2;
  string attribute_key = 3;
  string attribute_value = 4;
  string created_at = 5;
}
```

### Category Message (NO ATTRIBUTES)
**File**: `/p/github.com/sveturs/listings/api/proto/listings/v1/listings.proto:458-474`

```protobuf
message Category {
  int64 id = 1;
  string name = 2;
  string slug = 3;
  optional int64 parent_id = 4;
  optional string icon = 5;
  optional string description = 6;
  bool is_active = 7;
  int32 listing_count = 8;
  int32 sort_order = 9;
  int32 level = 10;
  map<string, string> translations = 11;  // Simple map, not attributes
  bool has_custom_ui = 12;
  optional string custom_ui_component = 13;
  string created_at = 14;
}
```

### CreateProductRequest (Accepts Attributes)
**File**: `/p/github.com/sveturs/listings/api/proto/listings/v1/listings.proto:1051-1074`

```protobuf
message CreateProductRequest {
  int64 storefront_id = 1;
  string name = 2;
  string description = 3;
  double price = 4;
  string currency = 5;
  int64 category_id = 6;
  optional string sku = 7;
  optional string barcode = 8;
  int32 stock_quantity = 9;
  bool is_active = 10;
  google.protobuf.Struct attributes = 11;  // ← NO VALIDATION
  
  bool has_individual_location = 12;
  optional string individual_address = 13;
  optional double individual_latitude = 14;
  optional double individual_longitude = 15;
  optional string location_privacy = 16;
  bool show_on_map = 17;
  bool has_variants = 18;
}
```

---

## Important Notes

1. **JSONB vs Key-Value Hybrid**:
   - Legacy listings use `listing_attributes` table (key-value)
   - New products use `attributes` JSONB column
   - No validation in either approach

2. **Variant Issue**:
   - Code still defines variant_attributes in proto
   - But b2c_product_variants table was DROPPED
   - Storing variant data will FAIL

3. **No Metadata Support**:
   - Can't define attribute schemas per category
   - No type information (string, number, enum)
   - No validation rules
   - No required/optional marking

4. **Categories Are Present**:
   - Full hierarchy support
   - Custom UI component support
   - But NO link to attribute definitions
