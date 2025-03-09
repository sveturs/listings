-- Обновленная схема с каскадным удалением
ALTER TABLE marketplace_favorites DROP CONSTRAINT marketplace_favorites_listing_id_fkey;
ALTER TABLE marketplace_favorites ADD CONSTRAINT marketplace_favorites_listing_id_fkey 
    FOREIGN KEY (listing_id) REFERENCES marketplace_listings(id) ON DELETE CASCADE;