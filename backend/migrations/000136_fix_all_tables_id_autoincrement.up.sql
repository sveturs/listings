-- Fix auto-increment for all tables that have sequences but missing DEFAULT
-- This migration fixes 56 tables that were created with sequences but without proper DEFAULT values

-- 1. address_change_log
ALTER TABLE address_change_log 
ALTER COLUMN id SET DEFAULT nextval('address_change_log_id_seq'::regclass);

-- 2. admin_users
ALTER TABLE admin_users 
ALTER COLUMN id SET DEFAULT nextval('admin_users_id_seq'::regclass);

-- 3. attribute_group_items
ALTER TABLE attribute_group_items 
ALTER COLUMN id SET DEFAULT nextval('attribute_group_items_id_seq'::regclass);

-- 4. attribute_groups
ALTER TABLE attribute_groups 
ALTER COLUMN id SET DEFAULT nextval('attribute_groups_id_seq'::regclass);

-- 5. attribute_option_translations
ALTER TABLE attribute_option_translations 
ALTER COLUMN id SET DEFAULT nextval('attribute_option_translations_id_seq'::regclass);

-- 6. balance_transactions
ALTER TABLE balance_transactions 
ALTER COLUMN id SET DEFAULT nextval('balance_transactions_id_seq'::regclass);

-- 7. category_attribute_groups
ALTER TABLE category_attribute_groups 
ALTER COLUMN id SET DEFAULT nextval('category_attribute_groups_id_seq'::regclass);

-- 8. category_attributes
ALTER TABLE category_attributes 
ALTER COLUMN id SET DEFAULT nextval('category_attributes_id_seq'::regclass);

-- 9. chat_attachments
ALTER TABLE chat_attachments 
ALTER COLUMN id SET DEFAULT nextval('chat_attachments_id_seq'::regclass);

-- 10. component_templates
ALTER TABLE component_templates 
ALTER COLUMN id SET DEFAULT nextval('component_templates_id_seq'::regclass);

-- 11. custom_ui_component_usage
ALTER TABLE custom_ui_component_usage 
ALTER COLUMN id SET DEFAULT nextval('custom_ui_component_usage_id_seq'::regclass);

-- 12. custom_ui_components
ALTER TABLE custom_ui_components 
ALTER COLUMN id SET DEFAULT nextval('custom_ui_components_id_seq'::regclass);

-- 13. custom_ui_templates
ALTER TABLE custom_ui_templates 
ALTER COLUMN id SET DEFAULT nextval('custom_ui_templates_id_seq'::regclass);

-- 14. escrow_payments
ALTER TABLE escrow_payments 
ALTER COLUMN id SET DEFAULT nextval('escrow_payments_id_seq'::regclass);

-- 15. geocoding_cache
ALTER TABLE geocoding_cache 
ALTER COLUMN id SET DEFAULT nextval('geocoding_cache_id_seq'::regclass);

-- 16. gis_filter_analytics
ALTER TABLE gis_filter_analytics 
ALTER COLUMN id SET DEFAULT nextval('gis_filter_analytics_id_seq'::regclass);

-- 17. gis_isochrone_cache
ALTER TABLE gis_isochrone_cache 
ALTER COLUMN id SET DEFAULT nextval('gis_isochrone_cache_id_seq'::regclass);

-- 18. gis_poi_cache
ALTER TABLE gis_poi_cache 
ALTER COLUMN id SET DEFAULT nextval('gis_poi_cache_id_seq'::regclass);

-- 19. import_history
ALTER TABLE import_history 
ALTER COLUMN id SET DEFAULT nextval('import_history_id_seq'::regclass);

-- 20. import_sources
ALTER TABLE import_sources 
ALTER COLUMN id SET DEFAULT nextval('import_sources_id_seq'::regclass);

-- 21. imported_categories
ALTER TABLE imported_categories 
ALTER COLUMN id SET DEFAULT nextval('imported_categories_id_seq'::regclass);

-- 22. item_performance_metrics
ALTER TABLE item_performance_metrics 
ALTER COLUMN id SET DEFAULT nextval('item_performance_metrics_id_seq'::regclass);

-- 23. listing_attribute_values
ALTER TABLE listing_attribute_values 
ALTER COLUMN id SET DEFAULT nextval('listing_attribute_values_id_seq'::regclass);

-- 24. listing_views
ALTER TABLE listing_views 
ALTER COLUMN id SET DEFAULT nextval('listing_views_id_seq'::regclass);

-- 25. listings_geo
ALTER TABLE listings_geo 
ALTER COLUMN id SET DEFAULT nextval('listings_geo_id_seq'::regclass);

-- 26. marketplace_chats
ALTER TABLE marketplace_chats 
ALTER COLUMN id SET DEFAULT nextval('marketplace_chats_id_seq'::regclass);

-- 27. marketplace_images
ALTER TABLE marketplace_images 
ALTER COLUMN id SET DEFAULT nextval('marketplace_images_id_seq'::regclass);

-- 28. marketplace_messages
ALTER TABLE marketplace_messages 
ALTER COLUMN id SET DEFAULT nextval('marketplace_messages_id_seq'::regclass);

-- 29. marketplace_orders
ALTER TABLE marketplace_orders 
ALTER COLUMN id SET DEFAULT nextval('marketplace_orders_id_seq'::regclass);

-- 30. merchant_payouts
ALTER TABLE merchant_payouts 
ALTER COLUMN id SET DEFAULT nextval('merchant_payouts_id_seq'::regclass);

-- 31. payment_gateways
ALTER TABLE payment_gateways 
ALTER COLUMN id SET DEFAULT nextval('payment_gateways_id_seq'::regclass);

-- 32. payment_methods
ALTER TABLE payment_methods 
ALTER COLUMN id SET DEFAULT nextval('payment_methods_id_seq'::regclass);

-- 33. payment_transactions
ALTER TABLE payment_transactions 
ALTER COLUMN id SET DEFAULT nextval('payment_transactions_id_seq'::regclass);

-- 34. price_history
ALTER TABLE price_history 
ALTER COLUMN id SET DEFAULT nextval('price_history_id_seq'::regclass);

-- 35. product_variant_attribute_values
ALTER TABLE product_variant_attribute_values 
ALTER COLUMN id SET DEFAULT nextval('product_variant_attribute_values_id_seq'::regclass);

-- 36. product_variant_attributes
ALTER TABLE product_variant_attributes 
ALTER COLUMN id SET DEFAULT nextval('product_variant_attributes_id_seq'::regclass);

-- 37. review_confirmations
ALTER TABLE review_confirmations 
ALTER COLUMN id SET DEFAULT nextval('review_confirmations_id_seq'::regclass);

-- 38. review_disputes
ALTER TABLE review_disputes 
ALTER COLUMN id SET DEFAULT nextval('review_disputes_id_seq'::regclass);

-- 39. review_responses
ALTER TABLE review_responses 
ALTER COLUMN id SET DEFAULT nextval('review_responses_id_seq'::regclass);

-- 40. reviews
ALTER TABLE reviews 
ALTER COLUMN id SET DEFAULT nextval('reviews_id_seq'::regclass);

-- 41. search_behavior_metrics
ALTER TABLE search_behavior_metrics 
ALTER COLUMN id SET DEFAULT nextval('search_behavior_metrics_id_seq'::regclass);

-- 42. search_config
ALTER TABLE search_config 
ALTER COLUMN id SET DEFAULT nextval('search_config_id_seq'::regclass);

-- 43. search_optimization_sessions
ALTER TABLE search_optimization_sessions 
ALTER COLUMN id SET DEFAULT nextval('search_optimization_sessions_id_seq'::regclass);

-- 44. search_queries
ALTER TABLE search_queries 
ALTER COLUMN id SET DEFAULT nextval('search_queries_id_seq'::regclass);

-- 45. search_statistics
ALTER TABLE search_statistics 
ALTER COLUMN id SET DEFAULT nextval('search_statistics_id_seq'::regclass);

-- 46. search_synonyms
ALTER TABLE search_synonyms 
ALTER COLUMN id SET DEFAULT nextval('search_synonyms_id_seq'::regclass);

-- 47. search_synonyms_config
ALTER TABLE search_synonyms_config 
ALTER COLUMN id SET DEFAULT nextval('search_synonyms_config_id_seq'::regclass);

-- 48. search_weights
ALTER TABLE search_weights 
ALTER COLUMN id SET DEFAULT nextval('search_weights_id_seq'::regclass);

-- 49. search_weights_history
ALTER TABLE search_weights_history 
ALTER COLUMN id SET DEFAULT nextval('search_weights_history_id_seq'::regclass);

-- 50. shopping_cart_items
ALTER TABLE shopping_cart_items 
ALTER COLUMN id SET DEFAULT nextval('shopping_cart_items_id_seq'::regclass);

-- 51. shopping_carts
ALTER TABLE shopping_carts 
ALTER COLUMN id SET DEFAULT nextval('shopping_carts_id_seq'::regclass);

-- 52. translations
ALTER TABLE translations 
ALTER COLUMN id SET DEFAULT nextval('translations_id_seq'::regclass);

-- 53. transliteration_rules
ALTER TABLE transliteration_rules 
ALTER COLUMN id SET DEFAULT nextval('transliteration_rules_id_seq'::regclass);

-- 54. user_behavior_events
ALTER TABLE user_behavior_events 
ALTER COLUMN id SET DEFAULT nextval('user_behavior_events_id_seq'::regclass);

-- 55. user_contacts
ALTER TABLE user_contacts 
ALTER COLUMN id SET DEFAULT nextval('user_contacts_id_seq'::regclass);

-- 56. user_storefronts
ALTER TABLE user_storefronts 
ALTER COLUMN id SET DEFAULT nextval('user_storefronts_id_seq'::regclass);

-- Ensure all sequences are properly owned by their respective columns
ALTER SEQUENCE address_change_log_id_seq OWNED BY address_change_log.id;
ALTER SEQUENCE admin_users_id_seq OWNED BY admin_users.id;
ALTER SEQUENCE attribute_group_items_id_seq OWNED BY attribute_group_items.id;
ALTER SEQUENCE attribute_groups_id_seq OWNED BY attribute_groups.id;
ALTER SEQUENCE attribute_option_translations_id_seq OWNED BY attribute_option_translations.id;
ALTER SEQUENCE balance_transactions_id_seq OWNED BY balance_transactions.id;
ALTER SEQUENCE category_attribute_groups_id_seq OWNED BY category_attribute_groups.id;
ALTER SEQUENCE category_attributes_id_seq OWNED BY category_attributes.id;
ALTER SEQUENCE chat_attachments_id_seq OWNED BY chat_attachments.id;
ALTER SEQUENCE component_templates_id_seq OWNED BY component_templates.id;
ALTER SEQUENCE custom_ui_component_usage_id_seq OWNED BY custom_ui_component_usage.id;
ALTER SEQUENCE custom_ui_components_id_seq OWNED BY custom_ui_components.id;
ALTER SEQUENCE custom_ui_templates_id_seq OWNED BY custom_ui_templates.id;
ALTER SEQUENCE escrow_payments_id_seq OWNED BY escrow_payments.id;
ALTER SEQUENCE geocoding_cache_id_seq OWNED BY geocoding_cache.id;
ALTER SEQUENCE gis_filter_analytics_id_seq OWNED BY gis_filter_analytics.id;
ALTER SEQUENCE gis_isochrone_cache_id_seq OWNED BY gis_isochrone_cache.id;
ALTER SEQUENCE gis_poi_cache_id_seq OWNED BY gis_poi_cache.id;
ALTER SEQUENCE import_history_id_seq OWNED BY import_history.id;
ALTER SEQUENCE import_sources_id_seq OWNED BY import_sources.id;
ALTER SEQUENCE imported_categories_id_seq OWNED BY imported_categories.id;
ALTER SEQUENCE item_performance_metrics_id_seq OWNED BY item_performance_metrics.id;
ALTER SEQUENCE listing_attribute_values_id_seq OWNED BY listing_attribute_values.id;
ALTER SEQUENCE listing_views_id_seq OWNED BY listing_views.id;
ALTER SEQUENCE listings_geo_id_seq OWNED BY listings_geo.id;
ALTER SEQUENCE marketplace_chats_id_seq OWNED BY marketplace_chats.id;
ALTER SEQUENCE marketplace_images_id_seq OWNED BY marketplace_images.id;
ALTER SEQUENCE marketplace_messages_id_seq OWNED BY marketplace_messages.id;
ALTER SEQUENCE marketplace_orders_id_seq OWNED BY marketplace_orders.id;
ALTER SEQUENCE merchant_payouts_id_seq OWNED BY merchant_payouts.id;
ALTER SEQUENCE payment_gateways_id_seq OWNED BY payment_gateways.id;
ALTER SEQUENCE payment_methods_id_seq OWNED BY payment_methods.id;
ALTER SEQUENCE payment_transactions_id_seq OWNED BY payment_transactions.id;
ALTER SEQUENCE price_history_id_seq OWNED BY price_history.id;
ALTER SEQUENCE product_variant_attribute_values_id_seq OWNED BY product_variant_attribute_values.id;
ALTER SEQUENCE product_variant_attributes_id_seq OWNED BY product_variant_attributes.id;
ALTER SEQUENCE review_confirmations_id_seq OWNED BY review_confirmations.id;
ALTER SEQUENCE review_disputes_id_seq OWNED BY review_disputes.id;
ALTER SEQUENCE review_responses_id_seq OWNED BY review_responses.id;
ALTER SEQUENCE reviews_id_seq OWNED BY reviews.id;
ALTER SEQUENCE search_behavior_metrics_id_seq OWNED BY search_behavior_metrics.id;
ALTER SEQUENCE search_config_id_seq OWNED BY search_config.id;
ALTER SEQUENCE search_optimization_sessions_id_seq OWNED BY search_optimization_sessions.id;
ALTER SEQUENCE search_queries_id_seq OWNED BY search_queries.id;
ALTER SEQUENCE search_statistics_id_seq OWNED BY search_statistics.id;
ALTER SEQUENCE search_synonyms_id_seq OWNED BY search_synonyms.id;
ALTER SEQUENCE search_synonyms_config_id_seq OWNED BY search_synonyms_config.id;
ALTER SEQUENCE search_weights_id_seq OWNED BY search_weights.id;
ALTER SEQUENCE search_weights_history_id_seq OWNED BY search_weights_history.id;
ALTER SEQUENCE shopping_cart_items_id_seq OWNED BY shopping_cart_items.id;
ALTER SEQUENCE shopping_carts_id_seq OWNED BY shopping_carts.id;
ALTER SEQUENCE translations_id_seq OWNED BY translations.id;
ALTER SEQUENCE transliteration_rules_id_seq OWNED BY transliteration_rules.id;
ALTER SEQUENCE user_behavior_events_id_seq OWNED BY user_behavior_events.id;
ALTER SEQUENCE user_contacts_id_seq OWNED BY user_contacts.id;
ALTER SEQUENCE user_storefronts_id_seq OWNED BY user_storefronts.id;