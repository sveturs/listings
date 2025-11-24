# Listings Microservice Architecture Analysis: Attributes Handling

## Executive Summary

The listings microservice has a **hybrid attributes system**:
1. **Legacy system**: `listing_attributes` table (key-value storage) - used for C2C listings
2. **Modern system**: JSONB columns (`attributes` field) - used for B2C products (in `listings` table)
3. **No metadata system**: There is NO table for attribute definitions, types, or validation rules

Categories exist as a separate table with basic metadata, but no category-specific attribute definitions.

---

## Database Schema Analysis

### 1. Attributes Storage - TWO APPROACHES

#### Approach A: Legacy Key-Value Storage (`listing_attributes` table)
**Location**: Migration `000001_initial_schema.up.sql`

```sql
CREATE TABLE listing_attributes (
    id BIGSERIAL PRIMARY KEY,
    listing_id BIGINT NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    attribute_key VARCHAR(100) NOT NULL,
    attribute_value TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(listing_id, attribute_key)
);

CREATE INDEX idx_listing_attributes_listing_id ON listing_attributes(listing_id);
CREATE INDEX idx_listing_attributes_key ON listing_attributes(attribute_key);
```

**Status**: Still exists but **NOT USED** for B2C products
**Purpose**: Originally designed for C2C marketplace attributes
**Limitation**: One-key-one-value only, no type information

#### Approach B: Modern JSONB Storage (in `listings` table)
**Location**: Migration `000012_add_attributes_to_listings.up.sql`

```sql
ALTER TABLE listings
ADD COLUMN IF NOT EXISTS attributes JSONB DEFAULT '{}'::jsonb;

CREATE INDEX IF NOT EXISTS idx_listings_attributes ON listings USING GIN (attributes);
```

**Also in B2C products**:
```sql
-- From migration 000004_add_b2c_products.up.sql
attributes JSONB DEFAULT '{}'::jsonb,  -- Product attributes & metadata
```

**And product variants**:
```sql
variant_attributes JSONB DEFAULT '{}'::jsonb,  -- e.g., {"size": "L", "color": "Red"}
dimensions JSONB,  -- {"length": 10, "width": 5, "height": 3, "unit": "cm"}
```

**Status**: **ACTIVELY USED** for B2C products
**Advantages**: 
- Flexible nested structure
- GIN indexing for efficient queries
- Native support for complex types

### 2. Categories Table (Schema & Capabilities)
**Location**: Migration `000011_restore_categories_table.up.sql`

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
    seo_keywords TEXT,
    CONSTRAINT check_root_categories_level CHECK (...)
);
```

**Renamed in migration 20251111230000_unify_table_names.up.sql to**: `categories`

**What's Missing**: NO attribute definitions per category
- No table like `category_attributes` or `attribute_definitions`
- No validation rules storage
- No attribute type information (string, number, enum, etc.)

### 3. Products Table (Main Listings Table)
**Location**: Unified in `listings` table with `source_type = 'b2c'`

**Key columns for B2C products**:
```sql
id BIGSERIAL PRIMARY KEY,
storefront_id BIGINT,  -- For B2C products only
title VARCHAR(255),
description TEXT,
price DECIMAL(15,2),
currency VARCHAR(3),
category_id BIGINT,
status VARCHAR(50),  -- 'draft', 'active', 'inactive', 'sold', 'archived'
quantity INTEGER,
sku VARCHAR(100),
attributes JSONB DEFAULT '{}'::jsonb,  -- FLEXIBLE ATTRIBUTES HERE
source_type VARCHAR(20),  -- CHECK (source_type IN ('c2c', 'b2c', 'storefront'))
stock_status VARCHAR(20),  -- 'in_stock', 'out_of_stock', 'low_stock'
has_variants BOOLEAN,
created_at, updated_at, deleted_at
```

---

## Domain Model Analysis

### Product Model (`internal/domain/product.go`)
```go
type Product struct {
    ID                int64                  // Database ID
    StorefrontID      int64                  // Owner's storefront
    Name              string
    Description       string
    Price             float64
    Currency          string
    CategoryID        int64
    SKU               *string
    Barcode           *string
    StockQuantity     int32
    StockStatus       string                 // enum: in_stock, low_stock, out_of_stock, pre_order
    IsActive          bool
    Attributes        map[string]interface{} // ← JSONB STORED AS MAP
    ViewCount         int32
    SoldCount         int32
    HasIndividualLocation bool
    IndividualAddress *string
    IndividualLatitude *float64
    IndividualLongitude *float64
    LocationPrivacy   *string               // exact, approximate, hidden
    ShowOnMap         bool
    HasVariants       bool
    CreatedAt         time.Time
    UpdatedAt         time.Time
    
    // Relations
    Variants []ProductVariant
    Images   []*ProductImage
}
```

### ProductVariant Model
```go
type ProductVariant struct {
    ID                int64                  
    ProductID         int64                  
    SKU               *string                
    Barcode           *string                
    Price             *float64               
    CompareAtPrice    *float64               // "Was" price
    CostPrice         *float64               
    StockQuantity     int32                  
    StockStatus       string                 
    LowStockThreshold *int32                 
    VariantAttributes map[string]interface{} // ← JSONB FOR VARIANT-SPECIFIC ATTRIBUTES
    Weight            *float64               
    Dimensions        map[string]interface{} // ← JSONB FOR DIMENSIONS
    IsActive          bool                   
    IsDefault         bool                   
    ViewCount         int32                  
    SoldCount         int32                  
    CreatedAt         time.Time              
    UpdatedAt         time.Time              
    
    // Relations
    Images []*ProductImage
}
```

### Listing Model (C2C) - `internal/domain/listing.go`
```go
type Listing struct {
    ID             int64
    UUID           string
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
    StockStatus    *string
    AttributesJSON *string                  // ← JSONB (not mapped to objects)
    
    // Relations
    Attributes []*ListingAttribute         // ← KEY-VALUE TABLE USED HERE
    Images     []*ListingImage
    Tags       []string
    Location   *ListingLocation
}

type ListingAttribute struct {
    ID             int64
    ListingID      int64
    AttributeKey   string
    AttributeValue string
    CreatedAt      time.Time
}
```

---

## Repository Layer Analysis

### Products Repository (`internal/repository/postgres/products_repository.go`)

**Attributes Handling**:
```go
// Reading attributes from JSONB
var attributesJSON []byte
// ... in SELECT query: p.attributes
// After scanning:
if len(attributesJSON) > 0 {
    if err := json.Unmarshal(attributesJSON, &product.Attributes); err != nil {
        // Handle error
    }
}

// Writing attributes to JSONB
var attributesJSON []byte
if len(input.Attributes) > 0 {
    attributesJSON, err = json.Marshal(input.Attributes)
}
// ... used in INSERT/UPDATE query: $N
```

**No validation**: Attributes are accepted as-is without:
- Type checking
- Required fields validation
- Allowed values validation
- Size constraints

### Categories Repository (`internal/repository/postgres/categories_repository.go`)

**No attribute metadata queries**:
- Only fetches category info (id, name, slug, level, etc.)
- No support for category-specific attribute definitions

### Product Variants Repository
**Status**: DEPRECATED (as of migration 000010)
- File header states: "⚠️ DEPRECATED: Product variants functionality is deprecated"
- Still attempts to use `b2c_product_variants` table (dropped in Phase 11.5)
- Methods return errors
- File kept for "API compatibility" only

---

## gRPC Proto Definition Analysis

### File: `api/proto/listings/v1/listings.proto`

**Attributes in Proto Messages**:

#### 1. ListingAttribute Message (line 420-427)
```protobuf
message ListingAttribute {
  int64 id = 1;
  int64 listing_id = 2;
  string attribute_key = 3;
  string attribute_value = 4;
  string created_at = 5;
}
```
Status: **KEY-VALUE ONLY** - reflects legacy `listing_attributes` table

#### 2. Listing Message (line 369-401)
```protobuf
message Listing {
  // ... other fields
  repeated ListingAttribute attributes = 22;  // Uses key-value model
  repeated ListingVariant variants = 25;
  // ...
}
```

#### 3. Product Message (line 494-521)
```protobuf
message Product {
  int64 id = 1;
  int64 storefront_id = 2;
  string name = 3;
  // ... other fields
  google.protobuf.Struct attributes = 13;  // ← JSONB STRUCT
  // ...
  repeated ProductVariant variants = 25;
}
```
Status: **USES protobuf.Struct** - supports arbitrary JSON/JSONB

#### 4. ProductVariant Message (line 523-544)
```protobuf
message ProductVariant {
  int64 id = 1;
  int64 product_id = 2;
  // ... other fields
  google.protobuf.Struct variant_attributes = 11;  // ← JSONB
  google.protobuf.Struct dimensions = 13;          // ← JSONB
  // ...
}
```
Status: **USES protobuf.Struct** - allows arbitrary JSON

#### 5. Category Message (line 458-474)
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
  map<string, string> translations = 11;  // Simple translations, not attributes
  bool has_custom_ui = 12;
  optional string custom_ui_component = 13;
  string created_at = 14;
}
```
Status: **NO ATTRIBUTES SUPPORT** - only basic fields

### Service Methods Related to Attributes

**For Products** (lines 186-246):
- CreateProduct, UpdateProduct, DeleteProduct
- BulkCreateProducts, BulkUpdateProducts, BulkDeleteProducts
- CreateProductVariant, UpdateProductVariant, DeleteProductVariant
- All accept `google.protobuf.Struct attributes` parameter
- **No validation or schema definition**

**For Listings** (lines 78-98):
- GetListing, CreateListing, UpdateListing, DeleteListing
- SearchListings, ListListings
- **Uses ListingAttribute** (key-value model)

**For Categories** (lines 116-132):
- GetCategory, GetCategoryTree, GetPopularCategories
- **No attribute management methods**

---

## Current Architecture Assessment

### STRENGTHS ✅
1. **Flexible attributes system** - JSONB allows any structure without schema changes
2. **Backward compatible** - Key-value table still available for C2C listings
3. **Indexed for performance** - GIN indexes on JSONB for efficient queries
4. **Version migration path** - Gradual migration from key-value to JSONB
5. **Multiple data types supported** - Nested objects, arrays, primitives

### WEAKNESSES ❌
1. **No metadata table** - No way to define allowed attributes per category
2. **No validation** - Any attribute name/value accepted without checking
3. **No type information** - Can't enforce string vs number vs enum
4. **No documentation** - No schema definition system for attributes
5. **Variant support broken** - `b2c_product_variants` table dropped but code references it
6. **No defaults** - No way to set default attributes per category
7. **No required fields** - Can't mark certain attributes as mandatory
8. **Inconsistent model** - C2C uses key-value, B2C uses JSONB, confusion possible

---

## Key Questions ANSWERED

### Q1: Is there a table for attribute metadata?
**Answer**: NO - there is no table for:
- Attribute definitions
- Type information
- Validation rules
- Category-specific attribute requirements
- Default values

### Q2: How are attributes currently handled?
**Answer**: TWO SYSTEMS:
- **C2C Listings**: `listing_attributes` key-value table
- **B2C Products**: `attributes` JSONB column in `listings` table
- **Variants**: `variant_attributes` JSONB column (in dropped table, but still referenced in code)

### Q3: What's the gRPC API for attributes?
**Answer**:
- **Products**: Accept `google.protobuf.Struct attributes` (flexible, no schema)
- **Listings**: Use `repeated ListingAttribute` (key-value pairs)
- **Variants**: Accept `google.protobuf.Struct variant_attributes` (flexible)
- **Categories**: NO attribute support in API

### Q4: Are categories in the microservice?
**Answer**: YES - fully present:
- Table: `categories` (renamed from `c2c_categories`)
- Repository: `categories_repository.go`
- Proto: `Category` message (lines 458-474)
- Service methods: GetCategory, GetCategoryTree, etc.
- Features: Hierarchical (parent_id), level tracking, custom UI support

---

## File Locations Reference

### Database Migrations
- **Attributes (key-value)**: `/p/github.com/sveturs/listings/migrations/000001_initial_schema.up.sql`
- **Attributes (JSONB)**: `/p/github.com/sveturs/listings/migrations/000012_add_attributes_to_listings.up.sql`
- **B2C Products**: `/p/github.com/sveturs/listings/migrations/000004_add_b2c_products.up.sql`
- **Categories**: `/p/github.com/sveturs/listings/migrations/000011_restore_categories_table.up.sql`
- **Latest unification**: `/p/github.com/sveturs/listings/migrations/20251111230000_unify_table_names.up.sql`

### Domain Models
- **Product**: `/p/github.com/sveturs/listings/internal/domain/product.go` (lines 10-39)
- **ProductVariant**: `/p/github.com/sveturs/listings/internal/domain/product.go` (lines 41-65)
- **Listing**: `/p/github.com/sveturs/listings/internal/domain/listing.go` (lines 10-50)
- **ListingAttribute**: `/p/github.com/sveturs/listings/internal/domain/listing.go` (lines 52-59)

### Repository Layer
- **Products**: `/p/github.com/sveturs/listings/internal/repository/postgres/products_repository.go`
- **Variants**: `/p/github.com/sveturs/listings/internal/repository/postgres/product_variants_repository.go` (DEPRECATED)
- **Categories**: `/p/github.com/sveturs/listings/internal/repository/postgres/categories_repository.go`

### gRPC Definitions
- **Proto file**: `/p/github.com/sveturs/listings/api/proto/listings/v1/listings.proto`
- **Product message**: Line 494-521
- **ProductVariant message**: Line 523-544
- **Listing message**: Line 369-401
- **ListingAttribute message**: Line 420-427
- **Category message**: Line 458-474
- **Service definitions**: Line 78-353

### Service Layer
- **Listings service**: `/p/github.com/sveturs/listings/internal/service/listings/service.go`
- **Storefront service**: `/p/github.com/sveturs/listings/internal/service/listings/storefront_service.go`
- **Stock service**: `/p/github.com/sveturs/listings/internal/service/listings/stock_service.go`

---

## Recommendations for Future Development

If you need to add attribute validation/definitions:

1. **Option A: Lightweight** - Add validation in service layer
   - Maintain current JSONB flexibility
   - Add whitelist of allowed attributes per category
   - Document in code comments or README

2. **Option B: Structured** - Create metadata tables
   ```sql
   CREATE TABLE attribute_definitions (
       id SERIAL PRIMARY KEY,
       category_id INTEGER REFERENCES categories(id),
       attribute_name VARCHAR(255) NOT NULL,
       attribute_type VARCHAR(50),  -- string, number, enum, boolean
       is_required BOOLEAN,
       allowed_values JSONB,  -- For enums
       default_value JSONB,
       validation_regex TEXT,
       UNIQUE(category_id, attribute_name)
   );
   ```

3. **Option C: Hybrid** - Use OpenSearch for attribute queries
   - Store attributes in ES for full-text search
   - Keep JSONB in PostgreSQL for storage
   - Use materialized views for validation

---

## Summary Table

| Aspect | Status | Location |
|--------|--------|----------|
| **Attribute Storage** | JSONB + Legacy key-value | `listings.attributes`, `listing_attributes` |
| **Attribute Metadata** | ❌ NOT IMPLEMENTED | N/A |
| **Categories Table** | ✅ PRESENT | `categories` table |
| **Category Attributes** | ❌ NOT SUPPORTED | N/A |
| **B2C Attributes** | ✅ JSONB | `listings.attributes` column |
| **Variant Attributes** | ⚠️ BROKEN | Dropped table, deprecated code |
| **gRPC for Products** | ✅ Flexible | `google.protobuf.Struct` |
| **gRPC for Listings** | ✅ Key-value | `repeated ListingAttribute` |
| **gRPC for Categories** | ❌ NO ATTRIBUTES | Basic category fields only |
| **Validation** | ❌ NONE | No constraints |
| **Documentation** | ❌ NONE | No schema definitions |

