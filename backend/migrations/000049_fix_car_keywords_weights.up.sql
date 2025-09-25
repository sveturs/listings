-- Исправление неправильных весов ключевых слов для категории автомобилей
-- Все автомобильные ключевые слова должны указывать на категорию 1301 (Lični automobili), а не на 1002 (Fashion/Moda)

-- Сначала удаляем все неправильные связи автомобильных ключевых слов с категорией Fashion (1002)
DELETE FROM category_keyword_weights
WHERE category_id = 1002
  AND keyword IN (
    'volkswagen', 'vw', 'touran', 'car', 'cars', 'auto', 'automobile', 'vehicle',
    'bmw', 'mercedes', 'audi', 'toyota', 'honda', 'ford', 'opel', 'peugeot',
    'renault', 'citroen', 'fiat', 'skoda', 'škoda', 'mazda', 'nissan', 'hyundai',
    'kia', 'volvo', 'saab', 'seat', 'dacia', 'chevrolet', 'chrysler', 'dodge',
    'jeep', 'land', 'rover', 'jaguar', 'porsche', 'ferrari', 'lamborghini',
    'minivan', 'suv', 'sedan', 'hatchback', 'coupe', 'convertible', 'wagon',
    'truck', 'van', 'pickup', 'crossover'
  );

-- Вставляем правильные веса для основных автомобильных брендов и ключевых слов
INSERT INTO category_keyword_weights (category_id, keyword, weight, success_rate, language)
VALUES
  -- Основные автомобильные термины
  (1301, 'car', 1.9, 0.95, 'ru'),
  (1301, 'cars', 1.9, 0.95, 'ru'),
  (1301, 'auto', 1.9, 0.95, 'ru'),
  (1301, 'automobile', 1.9, 0.95, 'ru'),
  (1301, 'vehicle', 1.8, 0.90, 'ru'),
  (1301, 'automobili', 1.9, 0.95, 'ru'),

  -- Немецкие бренды
  (1301, 'volkswagen', 1.85, 0.95, 'ru'),
  (1301, 'vw', 1.85, 0.95, 'ru'),
  (1301, 'bmw', 1.85, 0.95, 'ru'),
  (1301, 'mercedes', 1.85, 0.95, 'ru'),
  (1301, 'audi', 1.85, 0.95, 'ru'),
  (1301, 'opel', 1.85, 0.95, 'ru'),
  (1301, 'porsche', 1.85, 0.95, 'ru'),

  -- Японские бренды
  (1301, 'toyota', 1.85, 0.95, 'ru'),
  (1301, 'honda', 1.85, 0.95, 'ru'),
  (1301, 'nissan', 1.85, 0.95, 'ru'),
  (1301, 'mazda', 1.85, 0.95, 'ru'),
  (1301, 'mitsubishi', 1.85, 0.95, 'ru'),
  (1301, 'suzuki', 1.85, 0.95, 'ru'),
  (1301, 'subaru', 1.85, 0.95, 'ru'),

  -- Французские бренды
  (1301, 'peugeot', 1.85, 0.95, 'ru'),
  (1301, 'renault', 1.85, 0.95, 'ru'),
  (1301, 'citroen', 1.85, 0.95, 'ru'),
  (1301, 'citroën', 1.85, 0.95, 'ru'),

  -- Итальянские бренды
  (1301, 'fiat', 1.85, 0.95, 'ru'),
  (1301, 'alfa', 1.85, 0.95, 'ru'),
  (1301, 'ferrari', 1.85, 0.95, 'ru'),
  (1301, 'lamborghini', 1.85, 0.95, 'ru'),

  -- Корейские бренды
  (1301, 'hyundai', 1.85, 0.95, 'ru'),
  (1301, 'kia', 1.85, 0.95, 'ru'),

  -- Американские бренды
  (1301, 'ford', 1.85, 0.95, 'ru'),
  (1301, 'chevrolet', 1.85, 0.95, 'ru'),
  (1301, 'chrysler', 1.85, 0.95, 'ru'),
  (1301, 'dodge', 1.85, 0.95, 'ru'),
  (1301, 'jeep', 1.85, 0.95, 'ru'),
  (1301, 'tesla', 1.85, 0.95, 'ru'),

  -- Чешские и другие бренды
  (1301, 'skoda', 1.85, 0.95, 'ru'),
  (1301, 'škoda', 1.85, 0.95, 'ru'),
  (1301, 'seat', 1.85, 0.95, 'ru'),
  (1301, 'dacia', 1.85, 0.95, 'ru'),
  (1301, 'volvo', 1.85, 0.95, 'ru'),
  (1301, 'saab', 1.85, 0.95, 'ru'),

  -- Модели VW (для нашего случая)
  (1301, 'touran', 1.80, 0.90, 'ru'),
  (1301, 'golf', 1.80, 0.90, 'ru'),
  (1301, 'passat', 1.80, 0.90, 'ru'),
  (1301, 'tiguan', 1.80, 0.90, 'ru'),
  (1301, 'polo', 1.80, 0.90, 'ru'),

  -- Типы кузовов
  (1301, 'minivan', 1.75, 0.85, 'ru'),
  (1301, 'suv', 1.75, 0.85, 'ru'),
  (1301, 'sedan', 1.75, 0.85, 'ru'),
  (1301, 'hatchback', 1.75, 0.85, 'ru'),
  (1301, 'coupe', 1.75, 0.85, 'ru'),
  (1301, 'convertible', 1.75, 0.85, 'ru'),
  (1301, 'wagon', 1.75, 0.85, 'ru'),
  (1301, 'truck', 1.75, 0.85, 'ru'),
  (1301, 'van', 1.75, 0.85, 'ru'),
  (1301, 'pickup', 1.75, 0.85, 'ru'),
  (1301, 'crossover', 1.75, 0.85, 'ru')
ON CONFLICT (category_id, keyword, language) DO UPDATE
SET weight = EXCLUDED.weight,
    success_rate = EXCLUDED.success_rate;