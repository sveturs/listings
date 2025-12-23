-- ============================================================================
-- Migration: Fix inheritance function to include global attributes
-- Date: 2025-12-17
-- Description: Updates get_category_attributes_with_inheritance() to automatically
--              include is_global = true attributes for all categories
-- ============================================================================

-- Drop existing UUID-based function
DROP FUNCTION IF EXISTS get_category_attributes_with_inheritance(uuid);

-- Create improved function with global attributes support
CREATE OR REPLACE FUNCTION get_category_attributes_with_inheritance(
    p_category_id UUID
)
RETURNS TABLE (
    attribute_id INTEGER,
    category_id UUID,
    is_enabled BOOLEAN,
    is_required BOOLEAN,
    is_searchable BOOLEAN,
    is_filterable BOOLEAN,
    sort_order INTEGER
) AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE category_path AS (
        -- Start with the given category
        SELECT
            c.id,
            c.parent_id,
            c.level,
            1 as distance
        FROM categories c
        WHERE c.id = p_category_id

        UNION ALL

        -- Recursively get parent categories (up the tree)
        SELECT
            c.id,
            c.parent_id,
            c.level,
            cp.distance + 1
        FROM categories c
        INNER JOIN category_path cp ON c.id = cp.parent_id
    ),
    inherited_attributes AS (
        -- Get attributes from category chain (closest category wins)
        SELECT DISTINCT ON (a.id)
            a.id as attribute_id,
            ca.category_id,
            ca.is_enabled,
            COALESCE(ca.is_required, false) as is_required,
            COALESCE(ca.is_searchable, a.is_searchable) as is_searchable,
            COALESCE(ca.is_filterable, a.is_filterable) as is_filterable,
            COALESCE(ca.sort_order, a.sort_order) as sort_order,
            cp.distance
        FROM category_path cp
        INNER JOIN category_attributes ca ON cp.id = ca.category_id
        INNER JOIN attributes a ON ca.attribute_id = a.id
        WHERE ca.is_enabled = true
          AND a.is_active = true
        ORDER BY a.id, cp.distance ASC  -- Closest category wins
    ),
    global_attributes AS (
        -- Add global attributes (is_global = true) if not already included
        SELECT
            a.id as attribute_id,
            p_category_id as category_id,
            true as is_enabled,
            false as is_required,
            a.is_searchable,
            a.is_filterable,
            a.sort_order,
            999 as distance  -- Low priority (globals added last)
        FROM attributes a
        WHERE a.is_active = true
          AND a.is_global = true
          AND a.id NOT IN (SELECT attribute_id FROM inherited_attributes)
    ),
    combined AS (
        SELECT * FROM inherited_attributes
        UNION ALL
        SELECT * FROM global_attributes
    )
    SELECT
        c.attribute_id,
        c.category_id,
        c.is_enabled,
        c.is_required,
        c.is_searchable,
        c.is_filterable,
        c.sort_order
    FROM combined c
    ORDER BY c.sort_order, c.attribute_id;
END;
$$ LANGUAGE plpgsql STABLE;

-- ============================================================================
-- COMMENT
-- ============================================================================

COMMENT ON FUNCTION get_category_attributes_with_inheritance(UUID) IS
'Returns all attributes for a category including:
1. Direct attributes from category_attributes
2. Inherited attributes from parent categories (closest wins)
3. Global attributes (is_global = true) if not overridden
Used by AttributeRepository.GetByCategoryID() for seamless attribute inheritance.';
