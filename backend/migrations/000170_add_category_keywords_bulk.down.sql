-- Удаление ключевых слов для 20 категорий

-- Удаляем ключевые слова по category_id
DELETE FROM category_keywords 
WHERE category_id IN (
    1003, -- automotive
    1401, -- apartments
    1001, -- electronics
    1007, -- industrial
    1101, -- smartphones
    1601, -- farm-machinery
    1701, -- industrial-machinery
    1008, -- food-beverages
    1304, -- tires-and-wheels
    1201, -- mens-clothing
    1801, -- organic-food
    1108, -- electronics-accessories
    1508, -- plumbing
    1205, -- kids-clothing
    1020, -- events-tickets
    1208, -- bags
    1102, -- computers
    1315, -- winter-tires
    1337, -- hoods
    2013  -- hunting-fishing
) 
AND source = 'manual';