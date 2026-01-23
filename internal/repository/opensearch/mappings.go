package opensearch

// GetListingsIndexMapping returns the OpenSearch mapping for marketplace_listings index
// This mapping supports both C2C and B2C listings with full attribute support
// ФАЗА 2: Advanced multilingual search with Serbian/Russian/English analysis, synonyms, autocomplete
func GetListingsIndexMapping() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 1,
			"analysis": map[string]interface{}{
				// ====================
				// CHAR FILTERS
				// ====================
				"char_filter": map[string]interface{}{
					"serbian_cyrillic_to_latin": map[string]interface{}{
						"type": "mapping",
						"mappings": []string{
							// Lowercase cyrillic to latin
							"а=>a", "б=>b", "в=>v", "г=>g", "д=>d",
							"ђ=>đ", "е=>e", "ж=>ž", "з=>z", "и=>i",
							"ј=>j", "к=>k", "л=>l", "љ=>lj", "м=>m",
							"н=>n", "њ=>nj", "о=>o", "п=>p", "р=>r",
							"с=>s", "т=>t", "ћ=>ć", "у=>u", "ф=>f",
							"х=>h", "ц=>c", "ч=>č", "џ=>dž", "ш=>š",
							// Uppercase cyrillic to latin
							"А=>A", "Б=>B", "В=>V", "Г=>G", "Д=>D",
							"Ђ=>Đ", "Е=>E", "Ж=>Ž", "З=>Z", "И=>I",
							"Ј=>J", "К=>K", "Л=>L", "Љ=>Lj", "М=>M",
							"Н=>N", "Њ=>Nj", "О=>O", "П=>P", "Р=>R",
							"С=>S", "Т=>T", "Ћ=>Ć", "У=>U", "Ф=>F",
							"Х=>H", "Ц=>C", "Ч=>Č", "Џ=>Dž", "Ш=>Š",
						},
					},
				},

				// ====================
				// TOKEN FILTERS
				// ====================
				"filter": map[string]interface{}{
					// Serbian synonyms
					"serbian_synonyms": map[string]interface{}{
						"type": "synonym",
						"synonyms": []string{
							// Electronics
							"telefon,mobilni,smartphone,mob",
							"laptop,notebook,računar",
							"tablet,tablet računar",
							"patike,sneakers,cipele",
							"auto,automobil,kola",
							"bicikl,bike,bajs",
							// Clothing
							"majica,shirt,košulja",
							"farmerke,jeans,džins",
							"jakna,jacket,jaketa",
							// Home
							"nameštaj,furniture",
							"lampa,svetlo,light",
							"tepih,ćilim,carpet",
						},
					},

					// Russian synonyms
					"russian_synonyms": map[string]interface{}{
						"type": "synonym",
						"synonyms": []string{
							// Electronics
							"телефон,смартфон,мобильный,моб",
							"ноутбук,лаптоп,компьютер",
							"планшет,таблет",
							// Transport
							"автомобиль,машина,авто,тачка",
							"велосипед,байк,вело",
							// Clothing
							"кроссовки,кеды,обувь",
							"футболка,майка",
							"джинсы,штаны",
							"куртка,пальто",
						},
					},

					// English synonyms
					"english_synonyms": map[string]interface{}{
						"type": "synonym",
						"synonyms": []string{
							// Electronics
							"phone,smartphone,mobile,cell",
							"laptop,notebook,computer",
							"tablet,ipad",
							// Transport
							"car,automobile,vehicle",
							"bike,bicycle",
							// Clothing
							"sneakers,shoes,trainers",
							"tshirt,shirt",
							"jeans,pants,trousers",
							"jacket,coat",
						},
					},

					// Universal brand synonyms (cross-language)
					"universal_synonyms": map[string]interface{}{
						"type": "synonym",
						"synonyms": []string{
							// Brands - phone
							"iphone,ajfon,айфон",
							"samsung,самсунг",
							"xiaomi,сяоми,ксяоми",
							"huawei,хуавей",
							// Brands - clothing
							"nike,найк,најк",
							"adidas,адидас,адидас",
							"puma,пума",
							// Brands - cars
							"bmw,бмв",
							"mercedes,мерцедес,мерседес",
							"audi,ауди",
						},
					},

					// Autocomplete edge ngram
					"autocomplete_filter": map[string]interface{}{
						"type":     "edge_ngram",
						"min_gram": 2,
						"max_gram": 20,
					},

					// Trigram for did-you-mean
					"trigram_filter": map[string]interface{}{
						"type":     "ngram",
						"min_gram": 3,
						"max_gram": 3,
					},

					// Serbian stemmer
					"serbian_stemmer": map[string]interface{}{
						"type":     "stemmer",
						"language": "light_german", // Closest to Serbian
					},

					// Russian stemmer
					"russian_stemmer": map[string]interface{}{
						"type":     "stemmer",
						"language": "russian",
					},

					// English stemmer
					"english_stemmer": map[string]interface{}{
						"type":     "stemmer",
						"language": "english",
					},
				},

				// ====================
				// ANALYZERS
				// ====================
				"analyzer": map[string]interface{}{
					// Serbian analyzer (with cyrillic→latin transliteration)
					"serbian_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"char_filter": []string{
							"serbian_cyrillic_to_latin",
						},
						"filter": []string{
							"lowercase",
							"serbian_synonyms",
							"universal_synonyms",
							"serbian_stemmer",
						},
					},

					// Russian analyzer
					"russian_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"russian_synonyms",
							"universal_synonyms",
							"russian_stemmer",
						},
					},

					// English analyzer
					"english_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"english_synonyms",
							"universal_synonyms",
							"english_stemmer",
						},
					},

					// Universal analyzer (ICU folding + all synonyms)
					"universal_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"char_filter": []string{
							"serbian_cyrillic_to_latin",
						},
						"filter": []string{
							"lowercase",
							"icu_folding", // Normalize diacritics (ć→c, đ→d, etc.)
							"serbian_synonyms",
							"russian_synonyms",
							"english_synonyms",
							"universal_synonyms",
						},
					},

					// Autocomplete analyzer
					"autocomplete_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"char_filter": []string{
							"serbian_cyrillic_to_latin",
						},
						"filter": []string{
							"lowercase",
							"icu_folding",
							"autocomplete_filter",
						},
					},

					// Autocomplete search analyzer (no edge ngram)
					"autocomplete_search_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"char_filter": []string{
							"serbian_cyrillic_to_latin",
						},
						"filter": []string{
							"lowercase",
							"icu_folding",
						},
					},

					// Trigram analyzer for did-you-mean
					"trigram_analyzer": map[string]interface{}{
						"type":      "custom",
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"icu_folding",
							"trigram_filter",
						},
					},

					// Standard (unchanged)
					"standard": map[string]interface{}{
						"type": "standard",
					},
				},
			},
		},

		// ====================
		// MAPPINGS
		// ====================
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				// ==================== BASIC LISTING FIELDS ====================
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

				// ==================== TITLE (multilingual with subfields) ====================
				"title": map[string]interface{}{
					"type":     "text",
					"analyzer": "universal_analyzer",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
						"autocomplete": map[string]interface{}{
							"type":            "text",
							"analyzer":        "autocomplete_analyzer",
							"search_analyzer": "autocomplete_search_analyzer",
						},
						"trigram": map[string]interface{}{
							"type":     "text",
							"analyzer": "trigram_analyzer",
						},
					},
				},
				"description": map[string]interface{}{
					"type":     "text",
					"analyzer": "universal_analyzer",
				},

				// ==================== MULTILINGUAL FIELDS ====================
				// Serbian
				"title_sr": map[string]interface{}{
					"type":     "text",
					"analyzer": "serbian_analyzer",
				},
				"description_sr": map[string]interface{}{
					"type":     "text",
					"analyzer": "serbian_analyzer",
				},
				// English
				"title_en": map[string]interface{}{
					"type":     "text",
					"analyzer": "english_analyzer",
				},
				"description_en": map[string]interface{}{
					"type":     "text",
					"analyzer": "english_analyzer",
				},
				// Russian
				"title_ru": map[string]interface{}{
					"type":     "text",
					"analyzer": "russian_analyzer",
				},
				"description_ru": map[string]interface{}{
					"type":     "text",
					"analyzer": "russian_analyzer",
				},
				"original_language": map[string]interface{}{
					"type": "keyword",
				},

				// ==================== PRICE & CURRENCY ====================
				"price": map[string]interface{}{
					"type": "double",
				},
				"old_price": map[string]interface{}{
					"type": "double",
				},
				"discount_percent": map[string]interface{}{
					"type": "float",
				},
				"has_discount": map[string]interface{}{
					"type": "boolean",
				},
				"currency": map[string]interface{}{
					"type": "keyword",
				},

				// ==================== BRAND & CONDITION ====================
				"brand": map[string]interface{}{
					"type": "keyword",
					"fields": map[string]interface{}{
						"text": map[string]interface{}{
							"type":     "text",
							"analyzer": "universal_analyzer",
						},
					},
				},
				"condition": map[string]interface{}{
					"type": "keyword", // new, used, refurbished
				},

				// ==================== RATING & REVIEWS ====================
				"rating": map[string]interface{}{
					"type": "float",
				},
				"review_count": map[string]interface{}{
					"type": "integer",
				},

				// ==================== TAGS ====================
				"tags": map[string]interface{}{
					"type": "keyword",
				},

				// ==================== CATEGORY ====================
				"category_id": map[string]interface{}{
					"type": "keyword", // UUID as string
				},
				"category_slug": map[string]interface{}{
					"type": "keyword",
				},

				// ==================== STATUS & VISIBILITY ====================
				"status": map[string]interface{}{
					"type": "keyword",
				},
				"visibility": map[string]interface{}{
					"type": "keyword",
				},
				"stock_status": map[string]interface{}{
					"type": "keyword",
				},

				// ==================== QUANTITY & SKU ====================
				"quantity": map[string]interface{}{
					"type": "integer",
				},
				"sku": map[string]interface{}{
					"type": "keyword",
				},

				// ==================== SOURCE TYPE ====================
				"source_type": map[string]interface{}{
					"type": "keyword", // c2c, b2c
				},
				"document_type": map[string]interface{}{
					"type": "keyword",
				},

				// ==================== STOREFRONT INFO ====================
				"storefront_name": map[string]interface{}{
					"type": "keyword",
					"fields": map[string]interface{}{
						"text": map[string]interface{}{
							"type":     "text",
							"analyzer": "universal_analyzer",
						},
					},
				},
				"storefront_slug": map[string]interface{}{
					"type": "keyword",
				},
				"storefront_rating": map[string]interface{}{
					"type": "float",
				},
				"seller_verified": map[string]interface{}{
					"type": "boolean",
				},

				// ==================== SHIPPING ====================
				"shipping_available": map[string]interface{}{
					"type": "boolean",
				},
				"shipping_free": map[string]interface{}{
					"type": "boolean",
				},
				"shipping_price": map[string]interface{}{
					"type": "double",
				},

				// ==================== FLAGS ====================
				"is_promoted": map[string]interface{}{
					"type": "boolean",
				},
				"is_featured": map[string]interface{}{
					"type": "boolean",
				},
				"is_new_arrival": map[string]interface{}{
					"type": "boolean",
				},

				// ==================== POPULARITY SCORE ====================
				"popularity_score": map[string]interface{}{
					"type": "float",
				},
				"views_count": map[string]interface{}{
					"type": "integer",
				},
				"favorites_count": map[string]interface{}{
					"type": "integer",
				},

				// ==================== TIMESTAMPS ====================
				"created_at": map[string]interface{}{
					"type": "date",
				},
				"updated_at": map[string]interface{}{
					"type": "date",
				},
				"published_at": map[string]interface{}{
					"type": "date",
				},

				// ==================== IMAGES ====================
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

				// ==================== LOCATION ====================
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

				// ==================== ATTRIBUTES (nested for filtering) ====================
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

				// ==================== FLATTENED ATTRIBUTES ====================
				"attributes_flat": map[string]interface{}{
					"type":    "object",
					"enabled": false, // Store only, don't index
				},
				"attributes_searchable_text": map[string]interface{}{
					"type":     "text",
					"analyzer": "universal_analyzer",
				},

				// ==================== VARIANTS (nested) ====================
				"variants": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{
							"type": "long",
						},
						"sku": map[string]interface{}{
							"type": "keyword",
						},
						"price": map[string]interface{}{
							"type": "double",
						},
						"stock": map[string]interface{}{
							"type": "integer",
						},
						"attributes": map[string]interface{}{
							"type": "nested",
							"properties": map[string]interface{}{
								"code": map[string]interface{}{
									"type": "keyword",
								},
								"value": map[string]interface{}{
									"type": "keyword",
								},
							},
						},
					},
				},

				// ==================== VARIANT AGGREGATION FIELDS ====================
				"variant_skus": map[string]interface{}{
					"type": "keyword",
				},
				"variant_colors": map[string]interface{}{
					"type": "keyword",
				},
				"variant_sizes": map[string]interface{}{
					"type": "keyword",
				},
				"min_price": map[string]interface{}{
					"type": "double",
				},
				"max_price": map[string]interface{}{
					"type": "double",
				},
				"total_stock": map[string]interface{}{
					"type": "integer",
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

// GetVariantNestedQuery builds a nested query for variant filtering
// Example: Find listings with variants where color="red" AND size="M"
func GetVariantNestedQuery(attributeFilters map[string]string) map[string]interface{} {
	must := []interface{}{}

	for code, value := range attributeFilters {
		must = append(must, map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "variants.attributes",
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"must": []interface{}{
							map[string]interface{}{
								"term": map[string]interface{}{
									"variants.attributes.code": code,
								},
							},
							map[string]interface{}{
								"term": map[string]interface{}{
									"variants.attributes.value": value,
								},
							},
						},
					},
				},
			},
		})
	}

	return map[string]interface{}{
		"nested": map[string]interface{}{
			"path": "variants",
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": must,
				},
			},
		},
	}
}

// GetAutocompleteQuery builds autocomplete query using edge ngram
// Example: User types "sam" → matches "samsung"
func GetAutocompleteQuery(query string) map[string]interface{} {
	return map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query": query,
			"fields": []string{
				"title.autocomplete^3",
				"brand.text^2",
				"storefront_name.text",
			},
			"type": "phrase_prefix",
		},
	}
}

// GetDidYouMeanQuery builds fuzzy query using trigrams
// Example: "samsng" → suggests "samsung"
func GetDidYouMeanQuery(query string) map[string]interface{} {
	return map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query":     query,
			"fields":    []string{"title.trigram"},
			"fuzziness": "AUTO",
		},
	}
}
