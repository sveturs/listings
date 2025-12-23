#!/bin/bash
# OpenSearch Mapping Migration: Add category_slug field
# Phase 4 - PF-4.4: OpenSearch mapping optimization
# Date: 2025-12-18

set -euo pipefail

OPENSEARCH_URL="${OPENSEARCH_URL:-http://localhost:9200}"
INDEX_NAME="listings_microservice"
TEMP_INDEX="${INDEX_NAME}_temp_$(date +%s)"

echo "============================================"
echo "OpenSearch Mapping Migration: Add category_slug"
echo "============================================"
echo ""
echo "Index: $INDEX_NAME"
echo "OpenSearch: $OPENSEARCH_URL"
echo ""

# Step 1: Check if index exists
echo "[1/7] Checking if index exists..."
if ! curl -s -f -XHEAD "$OPENSEARCH_URL/$INDEX_NAME" > /dev/null 2>&1; then
    echo "❌ Error: Index $INDEX_NAME does not exist"
    exit 1
fi
echo "✅ Index exists"

# Step 2: Get current document count
DOC_COUNT=$(curl -s "$OPENSEARCH_URL/$INDEX_NAME/_count" | jq -r '.count')
echo ""
echo "[2/7] Current document count: $DOC_COUNT"

if [ "$DOC_COUNT" -eq 0 ]; then
    echo "⚠️  Index is empty. We can just update mapping without reindexing."

    # Add category_slug field to existing mapping
    echo ""
    echo "[3/7] Adding category_slug field to mapping..."

    curl -s -X PUT "$OPENSEARCH_URL/$INDEX_NAME/_mapping" \
        -H 'Content-Type: application/json' \
        -d '{
          "properties": {
            "category_slug": {
              "type": "keyword"
            }
          }
        }' | jq '.'

    echo "✅ Mapping updated successfully"
    echo ""
    echo "============================================"
    echo "Migration Complete!"
    echo "============================================"
    echo ""
    echo "Next steps:"
    echo "1. Reindex all listings from database to populate category_slug"
    echo "2. Run: python3 reindex_unified.py"
    exit 0
fi

# Step 3: Get current mapping
echo ""
echo "[3/7] Fetching current mapping..."
curl -s "$OPENSEARCH_URL/$INDEX_NAME/_mapping" > /tmp/current_mapping.json
echo "✅ Mapping saved to /tmp/current_mapping.json"

# Step 4: Create new index with updated mapping
echo ""
echo "[4/7] Creating temporary index with updated mapping..."

# Create new mapping with category_slug added
curl -s -X PUT "$OPENSEARCH_URL/$TEMP_INDEX" \
    -H 'Content-Type: application/json' \
    -d '{
      "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0,
        "analysis": {
          "analyzer": {
            "listing_analyzer": {
              "type": "custom",
              "tokenizer": "standard",
              "filter": ["lowercase", "asciifolding"]
            }
          }
        }
      },
      "mappings": {
        "properties": {
          "id": {"type": "long"},
          "uuid": {"type": "keyword"},
          "title": {"type": "text", "analyzer": "listing_analyzer"},
          "description": {"type": "text", "analyzer": "listing_analyzer"},
          "price": {"type": "scaled_float", "scaling_factor": 100.0},
          "currency": {"type": "keyword"},
          "category_id": {"type": "keyword"},
          "category_slug": {"type": "keyword"},
          "storefront_id": {"type": "keyword"},
          "user_id": {"type": "keyword"},
          "status": {"type": "keyword"},
          "stock_status": {"type": "keyword"},
          "visibility": {"type": "keyword"},
          "source_type": {"type": "keyword"},
          "quantity": {"type": "integer"},
          "sku": {"type": "keyword"},
          "tags": {"type": "keyword"},
          "views_count": {"type": "integer"},
          "favorites_count": {"type": "integer"},
          "created_at": {"type": "date"},
          "updated_at": {"type": "date"},
          "published_at": {"type": "date"},
          "location": {"type": "geo_point"},
          "country": {"type": "keyword"},
          "city": {"type": "keyword"},
          "postal_code": {"type": "keyword"},
          "address_line1": {"type": "text"},
          "address_line2": {"type": "text"},
          "images": {
            "type": "nested",
            "properties": {
              "id": {"type": "long"},
              "url": {"type": "keyword"},
              "thumbnail_url": {"type": "keyword"},
              "is_primary": {"type": "boolean"},
              "display_order": {"type": "integer"}
            }
          },
          "attributes": {
            "type": "nested",
            "properties": {
              "code": {"type": "keyword"},
              "type": {"type": "keyword"},
              "value_text": {
                "type": "text",
                "fields": {"keyword": {"type": "keyword", "ignore_above": 256}}
              },
              "value_number": {"type": "float"},
              "value_boolean": {"type": "boolean"},
              "value_date": {"type": "date"},
              "value_select": {"type": "keyword"},
              "value_multiselect": {"type": "keyword"}
            }
          }
        }
      }
    }' | jq '.'

echo "✅ Temporary index created: $TEMP_INDEX"

# Step 5: Reindex with category_slug population
echo ""
echo "[5/7] Reindexing data with category_slug population..."
echo "⚠️  This step requires database lookup for category slugs"
echo "⚠️  We'll copy data as-is for now, category_slug will be null"

curl -s -X POST "$OPENSEARCH_URL/_reindex?wait_for_completion=false" \
    -H 'Content-Type: application/json' \
    -d "{
      \"source\": {
        \"index\": \"$INDEX_NAME\"
      },
      \"dest\": {
        \"index\": \"$TEMP_INDEX\"
      }
    }" | jq '.'

echo "✅ Reindex task started (running in background)"
echo ""
echo "Wait 10 seconds for reindex to complete..."
sleep 10

# Step 6: Verify document count
NEW_COUNT=$(curl -s "$OPENSEARCH_URL/$TEMP_INDEX/_count" | jq -r '.count')
echo ""
echo "[6/7] Verifying document counts..."
echo "Original index: $DOC_COUNT documents"
echo "New index: $NEW_COUNT documents"

if [ "$NEW_COUNT" -ne "$DOC_COUNT" ]; then
    echo "❌ Error: Document count mismatch!"
    echo "Original: $DOC_COUNT, New: $NEW_COUNT"
    echo "Temporary index preserved for inspection: $TEMP_INDEX"
    exit 1
fi

echo "✅ Document counts match"

# Step 7: Swap indices (delete old, rename new)
echo ""
echo "[7/7] Swapping indices..."

# Delete old index
curl -s -X DELETE "$OPENSEARCH_URL/$INDEX_NAME" | jq '.'
echo "✅ Old index deleted"

# Create alias from temp to original name
curl -s -X POST "$OPENSEARCH_URL/_aliases" \
    -H 'Content-Type: application/json' \
    -d "{
      \"actions\": [
        {\"add\": {\"index\": \"$TEMP_INDEX\", \"alias\": \"$INDEX_NAME\"}}
      ]
    }" | jq '.'

echo "✅ Alias created: $TEMP_INDEX → $INDEX_NAME"

# Verify
FINAL_COUNT=$(curl -s "$OPENSEARCH_URL/$INDEX_NAME/_count" | jq -r '.count')
echo ""
echo "============================================"
echo "Migration Complete!"
echo "============================================"
echo ""
echo "Results:"
echo "- Original documents: $DOC_COUNT"
echo "- Final documents: $FINAL_COUNT"
echo "- Physical index: $TEMP_INDEX"
echo "- Alias: $INDEX_NAME"
echo ""
echo "⚠️  IMPORTANT: category_slug field is now in mapping but values are NULL"
echo ""
echo "Next steps:"
echo "1. Run full reindex to populate category_slug from database:"
echo "   cd /p/github.com/vondi-global/vondi/backend"
echo "   python3 reindex_unified.py"
echo ""
echo "2. Verify category_slug is populated:"
echo "   curl -s '$OPENSEARCH_URL/$INDEX_NAME/_search?q=*&size=1' | jq '.hits.hits[0]._source.category_slug'"
echo ""
