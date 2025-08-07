// backend/internal/storage/opensearch/mappings.go
package opensearch

const ListingMapping = `
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1,
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
      "car_keywords": {
          "type": "text",
          "analyzer": "default_analyzer",
          "fields": {
              "keyword": {
                  "type": "keyword",
                  "ignore_above": 256
              },
              "autocomplete": {
                  "type": "text",
                  "analyzer": "autocomplete"
              }
          }
      },
      "id": {
        "type": "integer"
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
      "average_rating": {
        "type": "double"
      },
      "review_count": {
        "type": "integer"
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
      "old_price": {
        "type": "double"
      },
      "has_discount": {
        "type": "boolean"
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
      "metadata": {
        "type": "object",
        "properties": {
          "discount": {
            "type": "object",
            "properties": {
              "discount_percent": {
                "type": "integer"
              },
              "previous_price": {
                "type": "double"
              },
              "effective_from": {
                "type": "date"
              },
              "has_price_history": {
                "type": "boolean"
              }
            }
          }
        }
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
      },
      "attributes": {
          "type": "nested",
          "properties": {
              "attribute_id": {
                  "type": "integer"
              },
              "attribute_name": {
                  "type": "keyword"
              },
              "display_name": {
                  "type": "text",
                  "fields": {
                      "keyword": {
                          "type": "keyword"
                      }
                  }
              },
              "attribute_type": {
                  "type": "keyword"
              },
              "text_value": {
                  "type": "text",
                  "analyzer": "default_analyzer",
                  "fields": {
                      "keyword": {
                          "type": "keyword"
                      },
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
              "text_value_lowercase": {
                  "type": "keyword"
              },
              "text_value_variations": {
                  "type": "text",
                  "analyzer": "default_analyzer"
              },
              "numeric_value": {
                  "type": "double"
              },
              "boolean_value": {
                  "type": "boolean"
              },
              "boolean_text": {
                  "type": "object"
              },
              "json_value": {
                  "type": "text",
                  "fields": {
                      "keyword": {
                          "type": "keyword"
                      }
                  }
              },
              "json_array": {
                  "type": "text",
                  "fields": {
                      "keyword": {
                          "type": "keyword"
                      }
                  }
              },
              "display_value": {
                  "type": "text",
                  "analyzer": "default_analyzer",
                  "fields": {
                      "keyword": {
                          "type": "keyword"
                      },
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
              "translations": {
                  "type": "object",
                  "properties": {
                      "en": { "type": "text", "analyzer": "english_analyzer" },
                      "sr": { "type": "text", "analyzer": "serbian_analyzer" },
                      "ru": { "type": "text", "analyzer": "russian_analyzer" }
                  }
              }
          }
      },
      "all_attributes_text": {
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
      "db_translations": {
          "type": "nested",
          "properties": {
              "field_name": {
                  "type": "keyword"
              },
              "language": {
                  "type": "keyword"
              },
              "translation": {
                  "type": "text"
              }
          }
      },
      "popularity_score": {
          "type": "float"
      }
    }
  }
}
  `
