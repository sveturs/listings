-- Remove added categories in reverse order (children first, then parents)

-- Remove Electronics subcategories
DELETE FROM marketplace_categories WHERE id IN (2400, 2401, 2402);

-- Remove Fashion subcategories  
DELETE FROM marketplace_categories WHERE id IN (2500, 2501, 2502);

-- Remove Business & Industrial subcategories
DELETE FROM marketplace_categories WHERE id IN (2101, 2102, 2103);

-- Remove Collectibles & Hobby subcategories
DELETE FROM marketplace_categories WHERE id IN (2201, 2202, 2203);

-- Remove Travel & Tourism subcategories
DELETE FROM marketplace_categories WHERE id IN (2301, 2302, 2303);

-- Remove parent categories
DELETE FROM marketplace_categories WHERE id IN (2100, 2200, 2300);