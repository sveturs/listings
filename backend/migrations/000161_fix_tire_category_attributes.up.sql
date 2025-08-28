-- Remove incorrect attributes from tire categories

-- Remove from main tire category (1304)
DELETE FROM category_attribute_mapping 
WHERE category_id = 1304 
AND attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN (
        'condition', 'boat_make', 'color', 'fuel_type', 'transmission',
        'engine_size', 'doors', 'seats', 'drive_type', 'location',
        'delivery_available', 'negotiable', 'warranty', 'return_policy',
        'truck_make'
    )
);

-- Remove from subcategories
DELETE FROM category_attribute_mapping 
WHERE category_id IN (1314, 1315, 1316) -- Summer, Winter, All-season tires
AND attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN (
        'condition', 'boat_make', 'color', 'fuel_type', 'transmission',
        'engine_size', 'doors', 'seats', 'drive_type', 'location',
        'delivery_available', 'negotiable', 'warranty', 'return_policy',
        'truck_make'
    )
);

-- Also remove from parent category Auto delovi (1303) if needed
DELETE FROM category_attribute_mapping 
WHERE category_id = 1303 
AND attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN (
        'boat_make', 'fuel_type', 'transmission', 'engine_size', 
        'doors', 'seats', 'drive_type', 'truck_make'
    )
);