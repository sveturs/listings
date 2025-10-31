# OpenSearch Setup Guide

This guide covers the OpenSearch indexing setup for the listings microservice.

## Overview

The listings microservice uses OpenSearch for full-text search, faceted filtering, and geo-location queries. This document describes the index structure, reindexing process, and validation procedures.

## Architecture

```
┌─────────────────┐         ┌──────────────────┐
│   PostgreSQL    │────────▶│   Reindex Script │
│   (listings_db) │         │   (Python 3)     │
└─────────────────┘         └──────────┬───────┘
                                       │
                                       ▼
                           ┌───────────────────────┐
                           │     OpenSearch        │
                           │ listings_microservice │
                           └───────────────────────┘
                                       │
                                       ▼
                           ┌───────────────────────┐
                           │  Listings Service     │
                           │  (Go microservice)    │
                           └───────────────────────┘
```

## Index Schema

**Index Name:** `listings_microservice`

**Key Features:**
- Full-text search on title and description (Russian stopwords)
- Autocomplete support (edge n-grams)
- Price range filtering (scaled_float)
- Category and status faceting
- Geo-location queries (geo_point)
- Nested objects for images and attributes

**Schema File:** `scripts/opensearch_schema.json`

### Mapped Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | long | Primary key |
| `uuid` | keyword | Unique identifier |
| `user_id` | long | Owner user ID |
| `storefront_id` | long | Associated storefront |
| `title` | text | Listing title (analyzed, autocomplete) |
| `description` | text | Listing description (analyzed) |
| `price` | scaled_float | Price (scaling factor: 100) |
| `currency` | keyword | Currency code |
| `category_id` | long | Category ID |
| `status` | keyword | Status (draft, active, inactive, sold, archived) |
| `visibility` | keyword | Visibility (public, private, unlisted) |
| `quantity` | integer | Available quantity |
| `sku` | keyword | Stock keeping unit |
| `views_count` | integer | Number of views |
| `favorites_count` | integer | Number of favorites |
| `location` | geo_point | Geographic coordinates |
| `country` | keyword | Country name |
| `city` | keyword | City name |
| `postal_code` | keyword | Postal/ZIP code |
| `address_line1` | text | Address line 1 |
| `address_line2` | text | Address line 2 |
| `images` | nested | Array of image objects |
| `tags` | keyword | Array of tags |
| `attributes` | nested | Array of key-value attributes |
| `created_at` | date | Creation timestamp |
| `updated_at` | date | Last update timestamp |
| `published_at` | date | Publication timestamp |

## Setup Process

### Prerequisites

1. **OpenSearch Running**
   ```bash
   # Check OpenSearch is accessible
   curl -X GET "http://localhost:9200/"
   ```

2. **Python 3 with Dependencies**
   ```bash
   pip3 install psycopg2-binary requests rich
   ```

3. **PostgreSQL Database Accessible**
   - Host: localhost (or remote)
   - Port: 35433 (microservice DB)
   - Database: listings_db
   - User: listings_user
   - Password: (from .env)

### Step 1: Create Index

Create the OpenSearch index with proper mapping:

```bash
cd /p/github.com/sveturs/listings
python3 scripts/create_opensearch_index.py
```

**Options:**
- `--force` - Delete existing index without confirmation

**Output:**
```
✓ OpenSearch connection successful
✓ Index 'listings_microservice' created successfully
✓ Index verification passed
```

### Step 2: Reindex Data

Migrate all listings from PostgreSQL to OpenSearch:

```bash
python3 scripts/reindex_listings.py \
  --target-password <password> \
  --verbose
```

**Options:**
- `--target-host HOST` - Database host (default: localhost)
- `--target-port PORT` - Database port (default: 35433)
- `--target-user USER` - Database user (default: listings_user)
- `--target-password PASS` - Database password (required)
- `--target-db DB` - Database name (default: listings_db)
- `--opensearch-url URL` - OpenSearch URL (default: http://localhost:9200)
- `--opensearch-index NAME` - Index name (default: listings_microservice)
- `--batch-size N` - Batch size (default: 500)
- `--verbose` - Enable verbose output
- `--dry-run` - Show what would be done without executing

**Example Output:**
```
✓ Connected to PostgreSQL
✓ Index 'listings_microservice' exists
Total listings to index: 8

Indexing listings... ━━━━━━━━━━━━━━━━━━━━ 100% (8/8)

┏━━━━━━━━━━━━━━━┳━━━━━━━┓
┃ Metric        ┃ Value ┃
┡━━━━━━━━━━━━━━━╇━━━━━━━┩
│ Total Listings│ 8     │
│ Indexed       │ 8     │
│ Failed        │ 0     │
│ Duration      │ 1.23s │
│ Rate          │ 6.50  │
└───────────────┴───────┘

✓ Reindex completed successfully!
```

### Step 3: Validate Index

Verify the reindex was successful:

```bash
python3 scripts/validate_opensearch.py \
  --target-password <password>
```

**Validations Performed:**
1. ✓ Document count matches PostgreSQL
2. ✓ All required fields present
3. ✓ Search functionality works
4. ✓ Aggregations work
5. ✓ Geo-location search works (if data available)
6. ✓ Index statistics healthy

**Example Output:**
```
======================================================================
                        Validation Results
======================================================================

✓ Passed (6)
  • Document Count
    Matched: 8 documents
  • Required Fields
    All 11 required fields present
  • Search Functionality
    Search works (45.23ms)
  • Aggregations
    Aggregations work (78.45ms)
  • Geo Search
    Geo search works (2 hits, 32.10ms)
  • Index Stats
    8 docs, 0.02 MB, 1 segments

======================================================================
✓ All validations passed!
======================================================================
```

## Microservice Configuration

Update the microservice environment variables:

```bash
# .env or deployment config
SVETULISTINGS_OPENSEARCH_ADDRESSES=http://localhost:9200
SVETULISTINGS_OPENSEARCH_USERNAME=admin
SVETULISTINGS_OPENSEARCH_PASSWORD=admin
SVETULISTINGS_OPENSEARCH_INDEX=listings_microservice
```

**Important:** Change the index name from `marketplace_listings` (monolith) to `listings_microservice` (microservice).

## Search API Examples

### Full-Text Search

```go
query := &domain.SearchListingsQuery{
    Query:  "phone case",
    Limit:  10,
    Offset: 0,
}

listings, total, err := openSearchClient.SearchListings(ctx, query)
```

### Search with Filters

```go
minPrice := 100.0
maxPrice := 5000.0
categoryID := int64(1301)

query := &domain.SearchListingsQuery{
    Query:      "phone",
    CategoryID: &categoryID,
    MinPrice:   &minPrice,
    MaxPrice:   &maxPrice,
    Limit:      20,
    Offset:     0,
}

listings, total, err := openSearchClient.SearchListings(ctx, query)
```

### Get Listing by ID

```go
listing, err := openSearchClient.GetListingByID(ctx, 328)
```

## Ongoing Maintenance

### Reindexing

Full reindex should be performed:
- After schema changes
- When data inconsistencies detected
- During major migrations
- Periodically (e.g., monthly) for optimization

```bash
# Full reindex
python3 scripts/reindex_listings.py --target-password <pwd>

# Validate after reindex
python3 scripts/validate_opensearch.py --target-password <pwd>
```

### Async Indexing (Production)

In production, the microservice uses an async indexing queue:

1. **On listing create/update:** Enqueue indexing job
2. **Worker process:** Consumes queue and updates OpenSearch
3. **Retry logic:** Failed operations retried up to 3 times

**Queue Table:** `indexing_queue`
```sql
SELECT * FROM indexing_queue WHERE status = 'failed';
```

### Monitoring

**Check Index Health:**
```bash
curl -X GET "http://localhost:9200/_cat/indices/listings_microservice?v"
```

**Check Document Count:**
```bash
curl -X GET "http://localhost:9200/listings_microservice/_count" | jq '.count'
```

**Check Index Stats:**
```bash
curl -X GET "http://localhost:9200/listings_microservice/_stats" | jq '.indices.listings_microservice.total'
```

**Search Performance:**
```bash
# Slow queries logged by microservice
grep "search completed" /var/log/listings-service.log
```

## Troubleshooting

### Problem: Index creation fails

**Symptoms:**
```
✗ Failed to create index: Connection refused
```

**Solutions:**
1. Check OpenSearch is running: `systemctl status opensearch` or `docker ps | grep opensearch`
2. Verify URL: `curl http://localhost:9200/`
3. Check credentials (username/password)

---

### Problem: Document count mismatch

**Symptoms:**
```
✗ Document Count
  Mismatch: PostgreSQL=100, OpenSearch=95 (diff=5)
```

**Solutions:**
1. Check for failed indexing operations:
   ```sql
   SELECT * FROM indexing_queue WHERE status = 'failed';
   ```
2. Rerun reindex script
3. Check OpenSearch logs for errors

---

### Problem: Search returns no results

**Symptoms:**
- API returns empty results for known data

**Solutions:**
1. Verify index exists: `curl http://localhost:9200/_cat/indices`
2. Check document count: `curl http://localhost:9200/listings_microservice/_count`
3. Test direct OpenSearch query:
   ```bash
   curl -X POST "http://localhost:9200/listings_microservice/_search" \
     -H 'Content-Type: application/json' \
     -d '{"query": {"match_all": {}}, "size": 1}'
   ```
4. Verify microservice uses correct index name (check SVETULISTINGS_OPENSEARCH_INDEX)

---

### Problem: Slow search performance

**Symptoms:**
- Search takes >500ms
- Validation warnings about performance

**Solutions:**
1. Check segment count: `curl http://localhost:9200/listings_microservice/_stats | jq '.indices.listings_microservice.total.segments.count'`
2. Force merge if segments >10:
   ```bash
   curl -X POST "http://localhost:9200/listings_microservice/_forcemerge?max_num_segments=1"
   ```
3. Increase `refresh_interval` in schema (default: 5s)
4. Add more replicas for read performance
5. Review query complexity

---

### Problem: Geo queries fail

**Symptoms:**
```
✗ Geo Search failed: field [location] not found
```

**Solutions:**
1. Check if listings have location data:
   ```sql
   SELECT COUNT(*) FROM listing_locations WHERE latitude IS NOT NULL;
   ```
2. Verify location field in OpenSearch:
   ```bash
   curl "http://localhost:9200/listings_microservice/_mapping" | jq '.listings_microservice.mappings.properties.location'
   ```
3. If schema missing location, recreate index and reindex

---

## Performance Tuning

### Index Settings

**For bulk indexing (reindex):**
```json
{
  "refresh_interval": "30s",
  "number_of_replicas": 0
}
```

**For production (search):**
```json
{
  "refresh_interval": "5s",
  "number_of_replicas": 1
}
```

### Query Optimization

1. **Use filters instead of queries** for exact matches (status, category)
2. **Limit result size** - Don't fetch more than needed
3. **Use pagination** - Offset + Limit pattern
4. **Cache frequent searches** - Redis cache for popular queries
5. **Use aggregations wisely** - Can be expensive

### Hardware Recommendations

**Development:**
- 2 GB RAM for OpenSearch
- 1 CPU core
- 10 GB disk

**Production (1M listings):**
- 8 GB RAM for OpenSearch
- 4 CPU cores
- 50 GB disk (SSD)
- 1 replica (for HA)

## References

- [OpenSearch Documentation](https://opensearch.org/docs/latest/)
- [Query DSL](https://opensearch.org/docs/latest/query-dsl/)
- [Mapping Types](https://opensearch.org/docs/latest/field-types/)
- [Geo Queries](https://opensearch.org/docs/latest/query-dsl/geo-and-xy/)

## Change Log

| Date | Version | Changes |
|------|---------|---------|
| 2025-10-31 | 1.0.0 | Initial OpenSearch setup for listings microservice |

---

**Last Updated:** 2025-10-31
**Maintained By:** Listings Microservice Team
