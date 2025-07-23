-- Fix marketplace_listings table to properly auto-increment ID
ALTER TABLE marketplace_listings 
ALTER COLUMN id SET DEFAULT nextval('marketplace_listings_id_seq'::regclass);

-- Ensure the sequence is owned by the column (should already be set, but making sure)
ALTER SEQUENCE marketplace_listings_id_seq OWNED BY marketplace_listings.id;