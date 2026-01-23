# Phase 1 - Category V2 Integration Progress

**Started:** 2025-12-16
**Goal:** Integrate UUID-based categories with existing int32-based code

---

## Database Status (✅ DONE)

- ✅ BE-1.1: UUID primary key migration applied
- ✅ BE-1.2: JSONB fields (name, description, meta_*) migration applied
- ✅ BE-1.3: GIN indexes on JSONB fields created
- ✅ BE-1.4: 18 L1 categories seeded

**Verified:**
```sql
-- Table has UUID id, JSONB name/description/meta fields
SELECT id, slug, name->>'sr', level FROM categories LIMIT 5;
```

---

## Code Integration Tasks

### STAGE 1: Domain Layer (BE-1.5-1.6) - ✅ DONE

**File:** `/p/github.com/vondi-global/listings/internal/domain/category.go`

**Status:** Completed 2025-12-16 14:45 UTC

**Tasks:**
- [x] Add `CategoryV2` struct with UUID and JSONB support
- [x] Add `LocalizedCategory` struct for API responses
- [x] Add `CategoryTreeV2` struct
- [x] Add `CategoryBreadcrumb` struct
- [x] Add `GetCategoryTreeFilterV2` with locale support
- [x] Add `Localize(locale string)` method for JSONB extraction
- [x] Add `getLocalized()` helper with fallback logic (sr -> en -> first available)
- [x] Keep existing `CategoryDetail` struct (backward compatibility)
- [x] Verified compilation: `go build ./...` passes

**Code to add:**
```go
type CategoryV2 struct {
    ID              uuid.UUID              `json:"id" db:"id"`
    Slug            string                 `json:"slug" db:"slug"`
    ParentID        *uuid.UUID             `json:"parent_id" db:"parent_id"`
    Level           int32                  `json:"level" db:"level"`
    Path            string                 `json:"path" db:"path"`
    SortOrder       int32                  `json:"sort_order" db:"sort_order"`
    Name            map[string]string      `json:"name" db:"name"`
    Description     map[string]string      `json:"description" db:"description"`
    MetaTitle       map[string]string      `json:"meta_title" db:"meta_title"`
    MetaDescription map[string]string      `json:"meta_description" db:"meta_description"`
    MetaKeywords    map[string]string      `json:"meta_keywords" db:"meta_keywords"`
    Icon            *string                `json:"icon" db:"icon"`
    ImageURL        *string                `json:"image_url" db:"image_url"`
    IsActive        bool                   `json:"is_active" db:"is_active"`
    CreatedAt       time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
}

type CategoryBreadcrumb struct {
    ID    uuid.UUID `json:"id"`
    Slug  string    `json:"slug"`
    Name  string    `json:"name"`  // Localized
    Level int32     `json:"level"`
}

// Localize extracts localized values from JSONB fields
func (c *CategoryV2) Localize(locale string) LocalizedCategory {
    return LocalizedCategory{
        ID:              c.ID,
        Slug:            c.Slug,
        Name:            getLocalized(c.Name, locale),
        Description:     getLocalized(c.Description, locale),
        MetaTitle:       getLocalized(c.MetaTitle, locale),
        MetaDescription: getLocalized(c.MetaDescription, locale),
        // ... other fields
    }
}
```

---

### STAGE 2: Repository Layer (BE-1.7-1.9) - ✅ DONE

**Files:**
- Updated `/p/github.com/vondi-global/listings/internal/repository/category_repository.go`
- Created `/p/github.com/vondi-global/listings/internal/repository/postgres/category_repository_v2.go`

**Status:** Completed 2025-12-16 15:00 UTC

**Tasks:**
- [x] Add `CategoryRepositoryV2` interface with UUID methods
- [x] Implement `GetByUUID(ctx, id) (*domain.CategoryV2, error)`
- [x] Implement `GetBySlugV2(ctx, slug) (*domain.CategoryV2, error)`
- [x] Implement `GetTreeV2(ctx, filter) ([]*domain.CategoryTreeV2, error)` with recursive loading
- [x] Implement `GetBreadcrumb(ctx, categoryID, locale) ([]*domain.CategoryBreadcrumb, error)` with CTE
- [x] Implement `ListV2(ctx, parentID, activeOnly, page, pageSize) ([]*domain.CategoryV2, int64, error)`
- [x] Add helper `scanCategoriesV2()` for row scanning with JSONB parsing
- [x] Add helper `getLocalizedFromMap()` for locale fallback
- [x] Keep existing int32 methods (backward compatibility)
- [x] Verified compilation: `go build ./...` passes

---

### STAGE 3: Cache Layer (BE-1.10) - ✅ DONE

**File:** Created `/p/github.com/vondi-global/listings/internal/cache/category_cache.go`

**Status:** Completed 2025-12-16 15:10 UTC

**Tasks:**
- [x] Create `CategoryCache` struct with Redis client
- [x] Implement `GetCategoryTree(ctx, key)` - retrieve cached tree
- [x] Implement `SetCategoryTree(ctx, key, tree, ttl)` - cache tree
- [x] Implement `GetCategoryBySlug(ctx, slug)` - retrieve cached category by slug
- [x] Implement `SetCategoryBySlug(ctx, slug, cat, ttl)` - cache category by slug
- [x] Implement `GetCategoryByUUID(ctx, uuid)` - retrieve cached category by UUID
- [x] Implement `SetCategoryByUUID(ctx, uuid, cat, ttl)` - cache category by UUID
- [x] Implement `GetBreadcrumb(ctx, categoryID, locale)` - retrieve cached breadcrumb
- [x] Implement `SetBreadcrumb(ctx, categoryID, locale, breadcrumbs, ttl)` - cache breadcrumb
- [x] Implement `InvalidateCategoryCache(ctx, pattern)` - invalidate by pattern (uses SCAN not KEYS)
- [x] Implement `InvalidateAll(ctx)` - invalidate all category caches
- [x] Implement `InvalidateCategory(ctx, categoryID, slug)` - invalidate specific category
- [x] Add cache key prefixes as constants
- [x] Add default TTL (1 hour)
- [x] Add proper logging with zerolog
- [x] Verified compilation: `go build ./...` passes

---

### STAGE 4: Service Layer - ⚪ TODO

**File:** Update `/p/github.com/vondi-global/listings/internal/service/category_service.go`

**Tasks:**
- [ ] Integrate repository V2 methods
- [ ] Add caching layer
- [ ] Add locale support

---

### STAGE 5: gRPC Layer (BE-1.11-1.12) - ⚪ TODO

**File 1:** Update `/p/github.com/vondi-global/listings/api/proto/categories/v1/categories.proto`

**Tasks:**
- [ ] Add `locale` field to requests
- [ ] Add `GetBreadcrumb` RPC method
- [ ] Update `Category` message to support localized fields

**File 2:** Create gRPC handler (reuse existing or update)

---

### STAGE 6: Frontend Components (FE-1.1-1.8) - ⚪ TODO

**Directory:** `/p/github.com/vondi-global/vondi/frontend/src/components/categories/`

**Tasks:**
- [ ] `CategoryTree.tsx` - Recursive tree component
- [ ] `CategoryBreadcrumb.tsx` - Navigation breadcrumb
- [ ] `CategoryCard.tsx` - Category card display
- [ ] Category page with SEO metadata

---

## Progress Log

### 2025-12-16 14:30 UTC
- Started integration phase
- Verified DB structure (UUID + JSONB confirmed)
- Created progress tracking document
- Beginning STAGE 1: Domain Layer

### 2025-12-16 14:45 UTC
- ✅ STAGE 1 COMPLETED: Domain Layer
- Added CategoryV2, LocalizedCategory, CategoryTreeV2, CategoryBreadcrumb
- Added Localize() method with fallback logic
- Verified compilation

### 2025-12-16 15:00 UTC
- ✅ STAGE 2 COMPLETED: Repository Layer
- Created category_repository_v2.go (544 lines)
- Implemented all V2 repository methods
- Added JSONB parsing, breadcrumb CTE, tree recursion
- Verified compilation

### 2025-12-16 15:10 UTC
- ✅ STAGE 3 COMPLETED: Cache Layer
- Created category_cache.go (316 lines)
- Implemented get/set/invalidate methods for all cache types
- Used SCAN for safe cache invalidation
- Verified compilation

### 2025-12-16 15:15 UTC
- 📊 Created PHASE1_BACKEND_SUMMARY.md
- **Status: 3 out of 6 stages completed (50% backend)**
- Total code written: ~990 lines
- Files created: 2, Files updated: 2
- Zero breaking changes - backward compatible

---

## Next Actions

1. Update `internal/domain/category.go` with V2 types
2. Test compilation: `go build ./...`
3. Update this document
4. Move to STAGE 2

---

## Notes

- **Backward Compatibility:** Keep existing `CategoryDetail` (int32 ID) for old code
- **Parallel Operation:** V1 and V2 code coexist during migration
- **No File Duplication:** Update existing files, don't create `*_new.go`
