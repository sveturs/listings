package opensearch

// GetListingsIndexMapping returns the OpenSearch mapping for marketplace_listings index
// This mapping supports both C2C and B2C listings with full attribute support
func GetListingsIndexMapping() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 1,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"standard": map[string]interface{}{
						"type": "standard",
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				// Basic listing fields
				"id": map[string]interface{}{
					"type": "long",
				},
				"uuid": map[string]interface{}{
					"type": "keyword",
				},
				"user_id": map[string]interface{}{
					"type": "long",
				},
				"storefront_id": map[string]interface{}{
					"type": "long",
				},
				"title": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"description": map[string]interface{}{
					"type": "text",
				},
				"price": map[string]interface{}{
					"type": "double",
				},
				"currency": map[string]interface{}{
					"type": "keyword",
				},
				"category_id": map[string]interface{}{
					"type": "long",
				},
				"status": map[string]interface{}{
					"type": "keyword",
				},
				"visibility": map[string]interface{}{
					"type": "keyword",
				},
				"quantity": map[string]interface{}{
					"type": "integer",
				},
				"sku": map[string]interface{}{
					"type": "keyword",
				},
				"source_type": map[string]interface{}{
					"type": "keyword",
				},
				"document_type": map[string]interface{}{
					"type": "keyword",
				},
				"stock_status": map[string]interface{}{
					"type": "keyword",
				},
				"views_count": map[string]interface{}{
					"type": "integer",
				},
				"favorites_count": map[string]interface{}{
					"type": "integer",
				},
				"created_at": map[string]interface{}{
					"type": "date",
				},
				"updated_at": map[string]interface{}{
					"type": "date",
				},
				"published_at": map[string]interface{}{
					"type": "date",
				},

				// Images
				"images": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type": "long",
						},
						"public_url": map[string]interface{}{
							"type": "keyword",
						},
						"file_path": map[string]interface{}{
							"type": "keyword",
						},
						"is_main": map[string]interface{}{
							"type": "boolean",
						},
					},
				},

				// Location fields
				"location": map[string]interface{}{
					"type": "geo_point",
				},
				"has_individual_location": map[string]interface{}{
					"type": "boolean",
				},
				"individual_latitude": map[string]interface{}{
					"type": "double",
				},
				"individual_longitude": map[string]interface{}{
					"type": "double",
				},
				"country": map[string]interface{}{
					"type": "keyword",
				},
				"city": map[string]interface{}{
					"type": "keyword",
				},

				// Attributes - Nested structure for advanced filtering
				"attributes": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type": "integer",
						},
						"code": map[string]interface{}{
							"type": "keyword",
						},
						"name": map[string]interface{}{
							"type": "text",
							"fields": map[string]interface{}{
								"keyword": map[string]interface{}{
									"type": "keyword",
								},
							},
						},
						"value_text": map[string]interface{}{
							"type": "text",
							"fields": map[string]interface{}{
								"keyword": map[string]interface{}{
									"type": "keyword",
								},
							},
						},
						"value_number": map[string]interface{}{
							"type": "double",
						},
						"value_boolean": map[string]interface{}{
							"type": "boolean",
						},
						"is_searchable": map[string]interface{}{
							"type": "boolean",
						},
						"is_filterable": map[string]interface{}{
							"type": "boolean",
						},
					},
				},

				// Flattened attributes for simple queries (denormalized)
				"attributes_flat": map[string]interface{}{
					"type":    "object",
					"enabled": false, // Store only, don't index
				},

				// Searchable text from all searchable attributes
				"attributes_searchable_text": map[string]interface{}{
					"type": "text",
				},
			},
		},
	}
}

// GetAttributeNestedQuery builds a nested query for attribute filtering
// Example: Find listings where attribute "brand" equals "Toyota"
func GetAttributeNestedQuery(attributeCode string, valueText *string, valueNumber *float64, valueBool *bool) map[string]interface{} {
	must := []interface{}{
		map[string]interface{}{
			"term": map[string]interface{}{
				"attributes.code": attributeCode,
			},
		},
	}

	if valueText != nil {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"attributes.value_text.keyword": *valueText,
			},
		})
	}

	if valueNumber != nil {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"attributes.value_number": *valueNumber,
			},
		})
	}

	if valueBool != nil {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"attributes.value_boolean": *valueBool,
			},
		})
	}

	return map[string]interface{}{
		"nested": map[string]interface{}{
			"path": "attributes",
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": must,
				},
			},
		},
	}
}

// GetAttributeRangeQuery builds a range query for numeric attributes
// Example: Find listings where attribute "year" >= 2020
func GetAttributeRangeQuery(attributeCode string, gte *float64, lte *float64) map[string]interface{} {
	rangeFilter := make(map[string]interface{})
	if gte != nil {
		rangeFilter["gte"] = *gte
	}
	if lte != nil {
		rangeFilter["lte"] = *lte
	}

	return map[string]interface{}{
		"nested": map[string]interface{}{
			"path": "attributes",
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []interface{}{
						map[string]interface{}{
							"term": map[string]interface{}{
								"attributes.code": attributeCode,
							},
						},
						map[string]interface{}{
							"range": map[string]interface{}{
								"attributes.value_number": rangeFilter,
							},
						},
					},
				},
			},
		},
	}
}
