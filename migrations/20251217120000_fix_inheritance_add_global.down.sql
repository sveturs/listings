-- Rollback: restore original function without global attributes

DROP FUNCTION IF EXISTS get_category_attributes_with_inheritance(uuid);

-- Recreate original version (without global attributes auto-inclusion)
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
        SELECT
            c.id,
            c.parent_id,
            c.level,
            1 as distance
        FROM categories c
        WHERE c.id = p_category_id

        UNION ALL

        SELECT
            c.id,
            c.parent_id,
            c.level,
            cp.distance + 1
        FROM categories c
        INNER JOIN category_path cp ON c.id = cp.parent_id
    )
    SELECT DISTINCT ON (a.id)
        a.id as attribute_id,
        ca.category_id,
        ca.is_enabled,
        COALESCE(ca.is_required, false) as is_required,
        COALESCE(ca.is_searchable, a.is_searchable) as is_searchable,
        COALESCE(ca.is_filterable, a.is_filterable) as is_filterable,
        COALESCE(ca.sort_order, a.sort_order) as sort_order
    FROM category_path cp
    INNER JOIN category_attributes ca ON cp.id = ca.category_id
    INNER JOIN attributes a ON ca.attribute_id = a.id
    WHERE ca.is_enabled = true
      AND a.is_active = true
    ORDER BY a.id, cp.distance ASC;
END;
$$ LANGUAGE plpgsql STABLE;
