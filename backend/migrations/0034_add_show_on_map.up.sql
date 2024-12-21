-- /backend/migrations/0034_add_show_on_map.up.sql
ALTER TABLE marketplace_listings
ADD COLUMN show_on_map BOOLEAN NOT NULL DEFAULT true;