# GetListing Images - Verification Guide

## Quick Verification

### 1. Run Integration Tests

```bash
cd /p/github.com/sveturs/listings
go test -v ./test/integration -run TestGetListing_WithImages -count=1
```

**Expected Result**: All 3 sub-tests should PASS ✅

### 2. Manual Testing via Database

#### Step 1: Check Docker Containers
```bash
docker ps | grep listings
```

Expected containers:
- `listings_app` - Main application (port 50051 gRPC, 8086 HTTP)
- `listings_postgres` - Database (port 35434)
- `listings_redis` - Cache (port 36380)

#### Step 2: Insert Test Data

```bash
# Connect to database
docker exec -it listings_postgres psql -U listings_user -d listings_dev_db

# Create category
INSERT INTO c2c_categories (id, name, slug, sort_order, level, is_active, count)
VALUES (9999, 'Test Category', 'test-category', 1, 0, true, 0);

# Create listing
INSERT INTO listings (id, user_id, category_id, title, description, price, currency, quantity, status, visibility, source_type, uuid, slug)
VALUES (9999, 1, 9999, 'Test Product with Images', 'Product description', 99.99, 'USD', 10, 'active', 'public', 'b2c', gen_random_uuid(), 'test-product-images');

# Create images
INSERT INTO listing_images (listing_id, url, storage_path, thumbnail_url, display_order, is_primary, width, height, file_size, mime_type)
VALUES
  (9999, 'https://example.com/product-main.jpg', '/storage/9999/main.jpg', 'https://example.com/thumbnails/main-thumb.jpg', 1, true, 1920, 1080, 256000, 'image/jpeg'),
  (9999, 'https://example.com/product-side.jpg', '/storage/9999/side.jpg', 'https://example.com/thumbnails/side-thumb.jpg', 2, false, 1920, 1080, 198000, 'image/jpeg'),
  (9999, 'https://example.com/product-back.jpg', '/storage/9999/back.jpg', NULL, 3, false, 1024, 768, 165000, 'image/jpeg');

# Verify images were created
SELECT id, listing_id, url, is_primary, display_order FROM listing_images WHERE listing_id = 9999 ORDER BY display_order;

\q
```

#### Step 3: Test via gRPC (using grpcurl)

If you have `grpcurl` installed:

```bash
grpcurl -plaintext \
  -d '{"id": 9999}' \
  localhost:50051 \
  listings.v1.ListingsService/GetListing
```

**Expected Response** (JSON format):
```json
{
  "listing": {
    "id": "9999",
    "title": "Test Product with Images",
    "images": [
      {
        "id": "...",
        "listingId": "9999",
        "url": "https://example.com/product-main.jpg",
        "storagePath": "/storage/9999/main.jpg",
        "thumbnailUrl": "https://example.com/thumbnails/main-thumb.jpg",
        "displayOrder": 1,
        "isPrimary": true,
        "width": 1920,
        "height": 1080,
        "fileSize": "256000",
        "mimeType": "image/jpeg"
      },
      {
        "id": "...",
        "listingId": "9999",
        "url": "https://example.com/product-side.jpg",
        "displayOrder": 2,
        "isPrimary": false
      },
      {
        "id": "...",
        "listingId": "9999",
        "url": "https://example.com/product-back.jpg",
        "displayOrder": 3,
        "isPrimary": false
      }
    ]
  }
}
```

### 3. Check Logs for Image Loading

```bash
# Watch real-time logs
docker logs -f listings_app

# Look for these log entries when calling GetListing:
# - "GetListing called" (debug level)
# - "failed to load images" (warning level, only if images fail to load)
```

## Verification Checklist

- [ ] Integration tests pass (all 3 scenarios)
- [ ] Images are returned in correct order (by display_order)
- [ ] Primary image has `is_primary: true`
- [ ] Optional fields are properly handled (null vs. filled)
- [ ] GetListing succeeds even without images
- [ ] GetListing doesn't fail if images fail to load (defensive error handling)

## Database Schema Reference

### listing_images table structure

```sql
Column         | Type                     | Nullable | Default
---------------|--------------------------|----------|------------------
id             | bigint                   | NOT NULL | nextval(...)
listing_id     | bigint                   | NOT NULL |
url            | text                     | NOT NULL |
storage_path   | text                     | NULL     |
thumbnail_url  | text                     | NULL     |
display_order  | integer                  | NOT NULL | 0
is_primary     | boolean                  | NOT NULL | false
width          | integer                  | NULL     |
height         | integer                  | NULL     |
file_size      | bigint                   | NULL     |
mime_type      | varchar(100)             | NULL     |
created_at     | timestamp with time zone | NOT NULL | CURRENT_TIMESTAMP
updated_at     | timestamp with time zone | NOT NULL | CURRENT_TIMESTAMP
```

**Foreign Key**: `listing_id` → `listings(id)` ON DELETE CASCADE

## Troubleshooting

### Images Not Returned

**Symptoms**: GetListing returns listing but `images` field is null or empty

**Possible Causes:**
1. No images in database for this listing → Expected behavior
2. Foreign key constraint violation → Check `listing_id` exists in `listings` table
3. Repository method error → Check logs for "failed to load images"

**Debug Steps:**
```bash
# Check if images exist in DB
docker exec listings_postgres psql -U listings_user -d listings_dev_db -c \
  "SELECT COUNT(*) FROM listing_images WHERE listing_id = YOUR_LISTING_ID;"

# Check container logs
docker logs listings_app --tail 100 | grep -i "image"
```

### Service Won't Start

**Symptoms**: Container exits or restarts repeatedly

**Debug Steps:**
```bash
# Check container status
docker ps -a | grep listings_app

# Check logs for errors
docker logs listings_app

# Rebuild and restart
cd /p/github.com/sveturs/listings
docker compose build app && docker compose up -d app
```

### Tests Fail

**Symptoms**: Integration tests fail with errors

**Common Issues:**
- Wrong table name (`categories` vs `c2c_categories`)
- Database connection issues
- Missing Docker containers

**Fix:**
```bash
# Restart test database
docker compose -f docker-compose.yml restart postgres

# Clear and rebuild
docker compose down -v
docker compose up -d
```

## Cleanup Test Data

```bash
docker exec -it listings_postgres psql -U listings_user -d listings_dev_db

DELETE FROM listing_images WHERE listing_id = 9999;
DELETE FROM listings WHERE id = 9999;
DELETE FROM c2c_categories WHERE id = 9999;

\q
```

---

**Last Updated**: 2025-11-11
**Service Version**: e04e8287-dirty
