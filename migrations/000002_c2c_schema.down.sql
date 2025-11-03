-- Phase 7.4: C2C Schema Rollback
-- Removes all c2c_* tables from listings microservice
-- Order is important to respect foreign key constraints

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS c2c_orders CASCADE;
DROP TABLE IF EXISTS c2c_messages CASCADE;
DROP TABLE IF EXISTS c2c_listing_variants CASCADE;
DROP TABLE IF EXISTS c2c_images CASCADE;
DROP TABLE IF EXISTS c2c_favorites CASCADE;
DROP TABLE IF EXISTS c2c_chats CASCADE;
DROP TABLE IF EXISTS c2c_listings CASCADE;
DROP TABLE IF EXISTS c2c_categories CASCADE;

-- Drop sequences
DROP SEQUENCE IF EXISTS c2c_orders_id_seq CASCADE;
DROP SEQUENCE IF EXISTS c2c_messages_id_seq CASCADE;
DROP SEQUENCE IF EXISTS c2c_listing_variants_id_seq CASCADE;
DROP SEQUENCE IF EXISTS c2c_images_id_seq CASCADE;
DROP SEQUENCE IF EXISTS c2c_chats_id_seq CASCADE;
DROP SEQUENCE IF EXISTS c2c_listings_id_seq CASCADE;
DROP SEQUENCE IF EXISTS c2c_categories_id_seq CASCADE;
