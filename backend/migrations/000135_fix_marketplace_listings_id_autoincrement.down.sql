-- Revert the auto-increment fix
ALTER TABLE marketplace_listings 
ALTER COLUMN id DROP DEFAULT;