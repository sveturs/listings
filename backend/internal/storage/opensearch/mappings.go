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
        "type": "keyword"
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
        "type": "keyword"
      },
      "country": {
        "type": "keyword"
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