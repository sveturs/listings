-- Revert auto-increment fix for all tables

-- 1. address_change_log
ALTER TABLE address_change_log ALTER COLUMN id DROP DEFAULT;

-- 2. admin_users
ALTER TABLE admin_users ALTER COLUMN id DROP DEFAULT;

-- 3. attribute_group_items
ALTER TABLE attribute_group_items ALTER COLUMN id DROP DEFAULT;

-- 4. attribute_groups
ALTER TABLE attribute_groups ALTER COLUMN id DROP DEFAULT;

-- 5. attribute_option_translations
ALTER TABLE attribute_option_translations ALTER COLUMN id DROP DEFAULT;

-- 6. balance_transactions
ALTER TABLE balance_transactions ALTER COLUMN id DROP DEFAULT;

-- 7. category_attribute_groups
ALTER TABLE category_attribute_groups ALTER COLUMN id DROP DEFAULT;

-- 8. category_attributes
ALTER TABLE category_attributes ALTER COLUMN id DROP DEFAULT;

-- 9. chat_attachments
ALTER TABLE chat_attachments ALTER COLUMN id DROP DEFAULT;

-- 10. component_templates
ALTER TABLE component_templates ALTER COLUMN id DROP DEFAULT;

-- 11. custom_ui_component_usage
ALTER TABLE custom_ui_component_usage ALTER COLUMN id DROP DEFAULT;

-- 12. custom_ui_components
ALTER TABLE custom_ui_components ALTER COLUMN id DROP DEFAULT;

-- 13. custom_ui_templates
ALTER TABLE custom_ui_templates ALTER COLUMN id DROP DEFAULT;

-- 14. escrow_payments
ALTER TABLE escrow_payments ALTER COLUMN id DROP DEFAULT;

-- 15. geocoding_cache
ALTER TABLE geocoding_cache ALTER COLUMN id DROP DEFAULT;

-- 16. gis_filter_analytics
ALTER TABLE gis_filter_analytics ALTER COLUMN id DROP DEFAULT;

-- 17. gis_isochrone_cache
ALTER TABLE gis_isochrone_cache ALTER COLUMN id DROP DEFAULT;

-- 18. gis_poi_cache
ALTER TABLE gis_poi_cache ALTER COLUMN id DROP DEFAULT;

-- 19. import_history
ALTER TABLE import_history ALTER COLUMN id DROP DEFAULT;

-- 20. import_sources
ALTER TABLE import_sources ALTER COLUMN id DROP DEFAULT;

-- 21. imported_categories
ALTER TABLE imported_categories ALTER COLUMN id DROP DEFAULT;

-- 22. item_performance_metrics
ALTER TABLE item_performance_metrics ALTER COLUMN id DROP DEFAULT;

-- 23. listing_attribute_values
ALTER TABLE listing_attribute_values ALTER COLUMN id DROP DEFAULT;

-- 24. listing_views
ALTER TABLE listing_views ALTER COLUMN id DROP DEFAULT;

-- 25. listings_geo
ALTER TABLE listings_geo ALTER COLUMN id DROP DEFAULT;

-- 26. marketplace_chats
ALTER TABLE marketplace_chats ALTER COLUMN id DROP DEFAULT;

-- 27. marketplace_images
ALTER TABLE marketplace_images ALTER COLUMN id DROP DEFAULT;

-- 28. marketplace_messages
ALTER TABLE marketplace_messages ALTER COLUMN id DROP DEFAULT;

-- 29. marketplace_orders
ALTER TABLE marketplace_orders ALTER COLUMN id DROP DEFAULT;

-- 30. merchant_payouts
ALTER TABLE merchant_payouts ALTER COLUMN id DROP DEFAULT;

-- 31. payment_gateways
ALTER TABLE payment_gateways ALTER COLUMN id DROP DEFAULT;

-- 32. payment_methods
ALTER TABLE payment_methods ALTER COLUMN id DROP DEFAULT;

-- 33. payment_transactions
ALTER TABLE payment_transactions ALTER COLUMN id DROP DEFAULT;

-- 34. price_history
ALTER TABLE price_history ALTER COLUMN id DROP DEFAULT;

-- 35. product_variant_attribute_values
ALTER TABLE product_variant_attribute_values ALTER COLUMN id DROP DEFAULT;

-- 36. product_variant_attributes
ALTER TABLE product_variant_attributes ALTER COLUMN id DROP DEFAULT;

-- 37. review_confirmations
ALTER TABLE review_confirmations ALTER COLUMN id DROP DEFAULT;

-- 38. review_disputes
ALTER TABLE review_disputes ALTER COLUMN id DROP DEFAULT;

-- 39. review_responses
ALTER TABLE review_responses ALTER COLUMN id DROP DEFAULT;

-- 40. reviews
ALTER TABLE reviews ALTER COLUMN id DROP DEFAULT;

-- 41. search_behavior_metrics
ALTER TABLE search_behavior_metrics ALTER COLUMN id DROP DEFAULT;

-- 42. search_config
ALTER TABLE search_config ALTER COLUMN id DROP DEFAULT;

-- 43. search_optimization_sessions
ALTER TABLE search_optimization_sessions ALTER COLUMN id DROP DEFAULT;

-- 44. search_queries
ALTER TABLE search_queries ALTER COLUMN id DROP DEFAULT;

-- 45. search_statistics
ALTER TABLE search_statistics ALTER COLUMN id DROP DEFAULT;

-- 46. search_synonyms
ALTER TABLE search_synonyms ALTER COLUMN id DROP DEFAULT;

-- 47. search_synonyms_config
ALTER TABLE search_synonyms_config ALTER COLUMN id DROP DEFAULT;

-- 48. search_weights
ALTER TABLE search_weights ALTER COLUMN id DROP DEFAULT;

-- 49. search_weights_history
ALTER TABLE search_weights_history ALTER COLUMN id DROP DEFAULT;

-- 50. shopping_cart_items
ALTER TABLE shopping_cart_items ALTER COLUMN id DROP DEFAULT;

-- 51. shopping_carts
ALTER TABLE shopping_carts ALTER COLUMN id DROP DEFAULT;

-- 52. translations
ALTER TABLE translations ALTER COLUMN id DROP DEFAULT;

-- 53. transliteration_rules
ALTER TABLE transliteration_rules ALTER COLUMN id DROP DEFAULT;

-- 54. user_behavior_events
ALTER TABLE user_behavior_events ALTER COLUMN id DROP DEFAULT;

-- 55. user_contacts
ALTER TABLE user_contacts ALTER COLUMN id DROP DEFAULT;

-- 56. user_storefronts
ALTER TABLE user_storefronts ALTER COLUMN id DROP DEFAULT;