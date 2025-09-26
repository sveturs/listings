-- Удаляем добавленные модели SEAT
DELETE FROM car_models
WHERE make_id = 22
  AND slug IN (
    'arona', 'ateca', 'ibiza', 'leon', 'formentor',
    'tarraco', 'altea', 'toledo', 'alhambra', 'mii',
    'exeo', 'cordoba'
  );