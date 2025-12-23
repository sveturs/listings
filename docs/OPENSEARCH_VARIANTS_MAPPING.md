# OpenSearch Mapping для Product Variants (Phase 3)

## Обновленный маппинг для products index

Добавить в существующий маппинг products:

```json
{
  "mappings": {
    "properties": {
      // ... existing fields ...

      "has_variants": {
        "type": "boolean"
      },

      "variants": {
        "type": "nested",
        "properties": {
          "variant_id": {
            "type": "keyword"
          },
          "sku": {
            "type": "keyword"
          },
          "price": {
            "type": "double"
          },
          "available_quantity": {
            "type": "integer"
          },
          "status": {
            "type": "keyword"
          },
          "attributes": {
            "type": "object",
            "dynamic": true,
            "properties": {
              "size": {"type": "keyword"},
              "color": {"type": "keyword"},
              "storage_capacity": {"type": "keyword"},
              "material": {"type": "keyword"}
            }
          }
        }
      },

      "min_price": {
        "type": "double"
      },

      "max_price": {
        "type": "double"
      },

      "available_sizes": {
        "type": "keyword"
      },

      "available_colors": {
        "type": "keyword"
      },

      "total_stock": {
        "type": "integer"
      },

      "has_stock": {
        "type": "boolean"
      }
    }
  }
}
```

## Пример документа с вариантами

```json
{
  "product_id": "a1b2c3d4-e5f6-4789-a0b1-c2d3e4f5g6h7",
  "name": {
    "en": "Nike Air Max 90",
    "ru": "Nike Air Max 90",
    "sr": "Nike Air Max 90"
  },
  "category_id": "shoes/sneakers",
  "storefront_id": "store-uuid",

  "base_price": 12990,

  "has_variants": true,

  "variants": [
    {
      "variant_id": "var-uuid-1",
      "sku": "SHO-A1B2C3-42-BLK",
      "price": 12990,
      "available_quantity": 5,
      "status": "active",
      "attributes": {
        "size": "42",
        "color": "black"
      }
    },
    {
      "variant_id": "var-uuid-2",
      "sku": "SHO-A1B2C3-43-WHT",
      "price": 12990,
      "available_quantity": 3,
      "status": "active",
      "attributes": {
        "size": "43",
        "color": "white"
      }
    },
    {
      "variant_id": "var-uuid-3",
      "sku": "SHO-A1B2C3-44-BLK",
      "price": 13490,
      "available_quantity": 2,
      "status": "active",
      "attributes": {
        "size": "44",
        "color": "black"
      }
    }
  ],

  "min_price": 12990,
  "max_price": 13490,

  "available_sizes": ["42", "43", "44"],
  "available_colors": ["black", "white"],

  "total_stock": 10,
  "has_stock": true,

  "created_at": "2025-12-17T10:00:00Z",
  "updated_at": "2025-12-17T10:00:00Z"
}
```

## Примеры поисковых запросов

### 1. Поиск товаров с конкретным размером в наличии

```json
{
  "query": {
    "bool": {
      "must": [
        {"term": {"category_id": "shoes/sneakers"}}
      ],
      "filter": [
        {
          "nested": {
            "path": "variants",
            "query": {
              "bool": {
                "must": [
                  {"term": {"variants.attributes.size": "42"}},
                  {"range": {"variants.available_quantity": {"gt": 0}}},
                  {"term": {"variants.status": "active"}}
                ]
              }
            }
          }
        }
      ]
    }
  }
}
```

### 2. Фильтр по цвету и размеру

```json
{
  "query": {
    "bool": {
      "filter": [
        {
          "nested": {
            "path": "variants",
            "query": {
              "bool": {
                "must": [
                  {"term": {"variants.attributes.size": "M"}},
                  {"term": {"variants.attributes.color": "black"}},
                  {"range": {"variants.available_quantity": {"gt": 0}}}
                ]
              }
            }
          }
        }
      ]
    }
  }
}
```

### 3. Агрегация: доступные размеры

```json
{
  "size": 0,
  "query": {
    "term": {"product_id": "a1b2c3d4-e5f6-4789-a0b1-c2d3e4f5g6h7"}
  },
  "aggs": {
    "available_sizes": {
      "nested": {"path": "variants"},
      "aggs": {
        "in_stock": {
          "filter": {
            "bool": {
              "must": [
                {"term": {"variants.status": "active"}},
                {"range": {"variants.available_quantity": {"gt": 0}}}
              ]
            }
          },
          "aggs": {
            "sizes": {
              "terms": {"field": "variants.attributes.size"}
            }
          }
        }
      }
    }
  }
}
```

### 4. Сортировка по min_price

```json
{
  "query": {
    "bool": {
      "must": [{"match": {"name.en": "sneakers"}}],
      "filter": [{"term": {"has_stock": true}}]
    }
  },
  "sort": [
    {"min_price": {"order": "asc"}},
    {"_score": {"order": "desc"}}
  ]
}
```

## Go Code: Индексация товара с вариантами

Файл: `internal/service/search/product_indexer.go` (обновить существующий)

```go
type ProductDocument struct {
    ProductID   string            `json:"product_id"`
    Name        map[string]string `json:"name"`
    CategoryID  string            `json:"category_id"`
    BasePrice   float64           `json:"base_price"`

    // НОВЫЕ ПОЛЯ
    HasVariants     bool                  `json:"has_variants"`
    Variants        []VariantDocument     `json:"variants,omitempty"`
    MinPrice        float64               `json:"min_price"`
    MaxPrice        float64               `json:"max_price"`
    AvailableSizes  []string              `json:"available_sizes,omitempty"`
    AvailableColors []string              `json:"available_colors,omitempty"`
    TotalStock      int32                 `json:"total_stock"`
    HasStock        bool                  `json:"has_stock"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type VariantDocument struct {
    VariantID         string                 `json:"variant_id"`
    SKU               string                 `json:"sku"`
    Price             float64                `json:"price"`
    AvailableQuantity int32                  `json:"available_quantity"`
    Status            string                 `json:"status"`
    Attributes        map[string]interface{} `json:"attributes"`
}

func (idx *ProductIndexer) IndexProductWithVariants(ctx context.Context, product *domain.Product, variants []*domain.ProductVariantV2) error {
    doc := &ProductDocument{
        ProductID:  product.ID.String(),
        Name:       product.NameTranslations,
        CategoryID: product.CategoryID.String(),
        BasePrice:  product.Price,
        HasVariants: len(variants) > 0,
    }

    if len(variants) > 0 {
        doc.Variants = make([]VariantDocument, 0, len(variants))

        minPrice := variants[0].GetEffectivePrice()
        maxPrice := minPrice
        totalStock := int32(0)
        sizesMap := make(map[string]bool)
        colorsMap := make(map[string]bool)

        for _, v := range variants {
            if v.Status != "active" {
                continue
            }

            // Build variant document
            varDoc := VariantDocument{
                VariantID:         v.ID.String(),
                SKU:               v.SKU,
                Price:             *v.Price,
                AvailableQuantity: v.GetAvailableQuantity(),
                Status:            v.Status,
                Attributes:        make(map[string]interface{}),
            }

            // Extract attributes
            for _, attr := range v.Attributes {
                if attr.ValueText != nil {
                    varDoc.Attributes[getAttributeCode(attr.AttributeID)] = *attr.ValueText

                    // Collect sizes/colors
                    if isSize(attr.AttributeID) {
                        sizesMap[*attr.ValueText] = true
                    }
                    if isColor(attr.AttributeID) {
                        colorsMap[*attr.ValueText] = true
                    }
                }
            }

            doc.Variants = append(doc.Variants, varDoc)

            // Update aggregates
            if v.Price != nil {
                if minPrice == nil || *v.Price < *minPrice {
                    minPrice = v.Price
                }
                if maxPrice == nil || *v.Price > *maxPrice {
                    maxPrice = v.Price
                }
            }

            totalStock += v.GetAvailableQuantity()
        }

        if minPrice != nil {
            doc.MinPrice = *minPrice
        }
        if maxPrice != nil {
            doc.MaxPrice = *maxPrice
        }

        doc.TotalStock = totalStock
        doc.HasStock = totalStock > 0

        // Convert maps to arrays
        for size := range sizesMap {
            doc.AvailableSizes = append(doc.AvailableSizes, size)
        }
        for color := range colorsMap {
            doc.AvailableColors = append(doc.AvailableColors, color)
        }
    }

    // Index document
    return idx.client.Index(ctx, "products", product.ID.String(), doc)
}
```

## Migration Plan

1. **Update mapping:** Add new fields to products index
2. **Update indexer:** Modify product indexer to include variants
3. **Reindex:** Run full reindex of products with variants
4. **Update search queries:** Add variant filtering support

## Notes

- Use `nested` type for variants to support independent filtering
- Pre-aggregate common filters (sizes, colors) for performance
- Update index on variant stock changes
- Consider partial updates for stock-only changes
