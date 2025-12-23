-- Migration: Fix category_id column type from bigint to UUID
-- This migration changes listings.category_id to match categories.id (UUID type)
-- Context: Categories table uses UUID primary keys, but listings.category_id was bigint

BEGIN;

-- Step 1: Drop existing foreign key constraint (if exists)
ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_category_id_fkey;

-- Step 2: Drop existing check constraint
ALTER TABLE listings DROP CONSTRAINT IF EXISTS listings_category_id_check;

-- Step 3: Drop existing index
DROP INDEX IF EXISTS idx_listings_category_id;

-- Step 4: Add new UUID column
ALTER TABLE listings ADD COLUMN category_uuid UUID;

-- Step 5: Set all existing values to NULL (no mapping possible from bigint IDs)
-- In a production scenario with data migration, we would populate this based on a mapping table
UPDATE listings SET category_uuid = NULL;

-- Step 6: Drop old bigint column
ALTER TABLE listings DROP COLUMN category_id;

-- Step 7: Rename new column to category_id
ALTER TABLE listings RENAME COLUMN category_uuid TO category_id;

-- Step 8: Add foreign key constraint to categories table
ALTER TABLE listings
  ADD CONSTRAINT listings_category_id_fkey
  FOREIGN KEY (category_id)
  REFERENCES categories(id)
  ON DELETE SET NULL
  ON UPDATE CASCADE;

-- Step 9: Create index for performance
CREATE INDEX idx_listings_category_id ON listings(category_id)
WHERE category_id IS NOT NULL AND is_deleted = false;

-- Step 10: Add comment
COMMENT ON COLUMN listings.category_id IS 'UUID reference to categories.id (changed from bigint in migration 20251217150000)';

COMMIT;
