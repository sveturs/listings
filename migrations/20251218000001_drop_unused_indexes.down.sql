-- Rollback: Recreate indexes if needed

CREATE INDEX IF NOT EXISTS idx_categories_tree
  ON categories (parent_id, level, sort_order)
  WHERE is_active = true;

CREATE INDEX IF NOT EXISTS idx_categories_is_active
  ON categories (is_active)
  WHERE is_active = true;

CREATE INDEX IF NOT EXISTS idx_categories_path
  ON categories (path varchar_pattern_ops);

CREATE INDEX IF NOT EXISTS idx_categories_active_l1
  ON categories (sort_order)
  WHERE level = 1 AND is_active = true;
