-- Migration: 20251217_phase2_inheritance_function
-- Description: Create function get_category_attributes_with_inheritance() for Phase 2
-- Author: Elite Full-Stack Architect (Claude AI)
-- Date: 2025-12-17
-- Task: BE-2.9 - Inheritance function for attributes

-- ============================================================================
-- FUNCTION: get_category_attributes_with_inheritance
-- Description: Returns all attributes for a category (inherited + own)
-- ============================================================================

CREATE OR REPLACE FUNCTION get_category_attributes_with_inheritance(p_category_id UUID)
RETURNS TABLE (
    attribute_id INT,
    category_id UUID,
    is_enabled BOOLEAN,
    is_required BOOLEAN,
    is_searchable BOOLEAN,
    is_filterable BOOLEAN,
    sort_order INT
) AS $$
BEGIN
    -- Return attributes from the entire hierarchy
    -- Priority: closer categories override parent settings
    RETURN QUERY
    WITH RECURSIVE category_hierarchy AS (
        -- Start with the target category
        SELECT
            c.id,
            c.parent_id,
            c.path,
            c.level,
            0 AS depth  -- 0 = target category
        FROM categories c
        WHERE c.id = p_category_id

        UNION ALL

        -- Recursively get parent categories
        SELECT
            c.id,
            c.parent_id,
            c.path,
            c.level,
            ch.depth + 1 AS depth
        FROM categories c
        INNER JOIN category_hierarchy ch ON c.id = ch.parent_id
        WHERE c.parent_id IS NOT NULL
    ),
    -- Collect attributes from all levels (with priority)
    collected_attributes AS (
        SELECT DISTINCT ON (ca.attribute_id)
            ca.attribute_id,
            ca.category_id,
            ca.is_enabled,
            ca.is_required,
            ca.is_searchable,
            ca.is_filterable,
            ca.sort_order,
            ch.depth  -- Lower depth = higher priority
        FROM category_hierarchy ch
        INNER JOIN category_attributes ca ON ca.category_id = ch.id
        WHERE ca.is_active = true
          AND ca.is_enabled = true
        ORDER BY ca.attribute_id, ch.depth ASC  -- Closer categories win
    )
    SELECT
        ca.attribute_id,
        ca.category_id,
        ca.is_enabled,
        ca.is_required,
        ca.is_searchable,
        ca.is_filterable,
        ca.sort_order
    FROM collected_attributes ca
    ORDER BY ca.sort_order, ca.attribute_id;
END;
$$ LANGUAGE plpgsql STABLE;

-- Add comment
COMMENT ON FUNCTION get_category_attributes_with_inheritance(UUID) IS 'Returns all attributes for a category including inherited from parents. Closer categories override parent settings.';
