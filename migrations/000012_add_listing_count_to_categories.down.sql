-- Remove listing_count column from categories table
ALTER TABLE public.categories DROP COLUMN IF EXISTS listing_count;
