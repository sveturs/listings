-- Откат изменений (удаление добавленных весов ключевых слов)
DELETE FROM category_keyword_weights
WHERE category_id = 1301
  AND keyword IN (
    'volkswagen', 'vw', 'touran', 'car', 'cars', 'auto', 'automobile', 'vehicle',
    'bmw', 'mercedes', 'audi', 'toyota', 'honda', 'ford', 'opel', 'peugeot',
    'renault', 'citroen', 'citroën', 'fiat', 'alfa', 'skoda', 'škoda', 'mazda',
    'nissan', 'hyundai', 'kia', 'volvo', 'saab', 'seat', 'dacia', 'chevrolet',
    'chrysler', 'dodge', 'jeep', 'ferrari', 'lamborghini', 'tesla',
    'minivan', 'suv', 'sedan', 'hatchback', 'coupe', 'convertible', 'wagon',
    'truck', 'van', 'pickup', 'crossover', 'automobili',
    'golf', 'passat', 'tiguan', 'polo', 'mitsubishi', 'suzuki', 'subaru'
  );