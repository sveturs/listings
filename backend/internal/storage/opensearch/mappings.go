// backend/internal/storage/opensearch/mappings.go
package opensearch

// ListingMapping содержит схему для индекса объявлений
const ListingMapping = `
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
    "analysis": {
      "analyzer": {
        "serbian_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": [
            "lowercase",
            "serbian_latin_stemmer"
          ]
        },
        "russian_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": [
            "lowercase",
            "russian_stemmer"
          ]
        },
        "english_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": [
            "lowercase",
            "english_stemmer"
          ]
        },
        "default_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": [
            "lowercase"
          ],
          "char_filter": [
            "html_strip"
          ]
        },
        "autocomplete": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": [
            "lowercase",
            "autocomplete_filter"
          ]
        },
        "shingle_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": [
            "lowercase",
            "shingle_filter"
          ]
        }
      },
      "filter": {
        "serbian_latin_stemmer": {
          "type": "stemmer",
          "language": "serbian"
        },
        "russian_stemmer": {
          "type": "stemmer",
          "language": "russian"
        },
        "english_stemmer": {
          "type": "stemmer",
          "language": "english"
        },
        "autocomplete_filter": {
          "type": "edge_ngram",
          "min_gram": 1,
          "max_gram": 20
        },
        "shingle_filter": {
          "type": "shingle",
          "min_shingle_size": 2,
          "max_shingle_size": 3
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "title": {
        "type": "text",
        "analyzer": "default_analyzer",
        "fields": {
          "serbian": {
            "type": "text",
            "analyzer": "serbian_analyzer"
          },
          "russian": {
            "type": "text",
            "analyzer": "russian_analyzer"
          },
          "english": {
            "type": "text",
            "analyzer": "english_analyzer"
          },
          "autocomplete": {
            "type": "text",
            "analyzer": "autocomplete"
          },
          "keyword": {
            "type": "keyword"
          },
          "shingles": {
            "type": "text",
            "analyzer": "shingle_analyzer"
          }
        }
      },
      "title_suggest": {
        "type": "completion",
        "analyzer": "autocomplete"
      },
      "title_variations": {
        "type": "text",
        "analyzer": "default_analyzer",
        "fields": {
          "keyword": {
            "type": "keyword",
            "ignore_above": 256
          }
        }
      },
      "description": {
        "type": "text",
        "analyzer": "default_analyzer",
        "fields": {
          "serbian": {
            "type": "text",
            "analyzer": "serbian_analyzer"
          },
          "russian": {
            "type": "text",
            "analyzer": "russian_analyzer"
          },
          "english": {
            "type": "text",
            "analyzer": "english_analyzer"
          }
        }
      },
      "price": {
        "type": "double"
      },
      "condition": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "status": {
        "type": "keyword"
      },
      "location": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "coordinates": {
        "type": "geo_point"
      },
      "city": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "country": {
        "type": "text",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "auto_brand": {
        "type": "keyword",
        "fields": {
          "text": {
            "type": "text",
            "analyzer": "default_analyzer"
          }
        }
      },
      "auto_model": {
        "type": "keyword",
        "fields": {
          "text": {
            "type": "text",
            "analyzer": "default_analyzer"
          }
        }
      },
      "auto_year": {
        "type": "integer"
      },
      "auto_mileage": {
        "type": "integer"
      },
      "auto_fuel_type": {
        "type": "keyword"
      },
      "auto_transmission": {
        "type": "keyword"
      },
      "auto_body_type": {
        "type": "keyword"
      },
      "auto_drive_type": {
        "type": "keyword"
      },
      "auto_engine_capacity": {
        "type": "float"
      },
      "auto_power": {
        "type": "integer"
      },
      "auto_color": {
        "type": "keyword"
      },
      "auto_number_of_doors": {
        "type": "integer"
      },
      "auto_number_of_seats": {
        "type": "integer"
      },
      "views_count": {
        "type": "integer"
      },
      "created_at": {
        "type": "date"
      },
      "updated_at": {
        "type": "date"
      },
      "show_on_map": {
        "type": "boolean"
      },
      "original_language": {
        "type": "keyword"
      },
      "category_id": {
        "type": "integer"
      },
      "user_id": {
        "type": "integer"
      },
      "storefront_id": {
        "type": "integer"
      },
      "category_path_ids": {
        "type": "integer"
      },
      "translations": {
        "properties": {
          "sr": {
            "properties": {
              "title": {
                "type": "text",
                "analyzer": "serbian_analyzer",
                "fields": {
                  "autocomplete": {
                    "type": "text",
                    "analyzer": "autocomplete"
                  }
                }
              },
              "description": {
                "type": "text",
                "analyzer": "serbian_analyzer"
              }
            }
          },
          "ru": {
            "properties": {
              "title": {
                "type": "text",
                "analyzer": "russian_analyzer",
                "fields": {
                  "autocomplete": {
                    "type": "text",
                    "analyzer": "autocomplete"
                  }
                }
              },
              "description": {
                "type": "text",
                "analyzer": "russian_analyzer"
              }
            }
          },
          "en": {
            "properties": {
              "title": {
                "type": "text",
                "analyzer": "english_analyzer",
                "fields": {
                  "autocomplete": {
                    "type": "text",
                    "analyzer": "autocomplete"
                  }
                }
              },
              "description": {
                "type": "text",
                "analyzer": "english_analyzer"
              }
            }
          }
        }
      },
      "images": {
        "type": "nested",
        "properties": {
          "id": {
            "type": "integer"
          },
          "file_path": {
            "type": "keyword"
          },
          "is_main": {
            "type": "boolean"
          }
        }
      },
      "category": {
        "type": "object",
        "properties": {
          "id": {
            "type": "integer"
          },
          "name": {
            "type": "text",
            "fields": {
              "keyword": {
                "type": "keyword"
              }
            }
          },
          "slug": {
            "type": "keyword"
          }
        }
      },
      "user": {
        "type": "object",
        "properties": {
          "id": {
            "type": "integer"
          },
          "name": {
            "type": "text",
            "fields": {
              "keyword": {
                "type": "keyword"
              }
            }
          },
          "email": {
            "type": "keyword"
          }
        }
      }
    }
  }
}
`