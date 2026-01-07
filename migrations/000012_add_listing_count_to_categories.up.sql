-- Add listing_count column to categories table
ALTER TABLE public.categories ADD COLUMN IF NOT EXISTS listing_count INTEGER DEFAULT 0 NOT NULL;

COMMENT ON COLUMN public.categories.listing_count IS 'Number of active listings in this category';
