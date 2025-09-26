-- Откат: удаляем добавленные ключевые слова брендов и моделей

DELETE FROM category_keywords
WHERE category_id = 1301
AND keyword IN (
    'mercedes-benz', 'mercedes benz',
    'alfa-romeo', 'alfa romeo',
    'land-rover', 'land rover',
    'rolls-royce', 'rolls royce',
    'e-class', 'c-class', 's-class', 'a-class', 'g-class',
    'gle', 'glc', 'gla', 'glb', 'gls',
    'x1', 'x2', 'x3', 'x4', 'x5', 'x6', 'x7',
    'a1', 'a3', 'a4', 'a5', 'a6', 'a7', 'a8',
    'q2', 'q3', 'q5', 'q7', 'q8',
    'golf', 'passat', 'tiguan', 'touran', 'touareg',
    'polo', 'arteon', 't-roc', 't-cross',
    'tdi', 'tfsi', 'tsi', 'cdi', 'bluemotion',
    'xdrive', 'quattro', '4matic', '4motion'
);