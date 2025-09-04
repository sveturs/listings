-- Откат добавления переводов для атрибутов
-- Migration 000050 DOWN: Remove attribute translations

-- Удаление всех переводов для unified_attribute добавленных в этой миграции
-- Сохраняем только переводы, которые были до миграции (атрибуты с ID >= 107)
DELETE FROM translations 
WHERE entity_type = 'unified_attribute'
  AND entity_id < 107;

-- Удаление переводов на русский язык для всех атрибутов (их не было вообще)
DELETE FROM translations 
WHERE entity_type = 'unified_attribute'
  AND language = 'ru';

-- Обновление статистики
UPDATE unified_attributes SET updated_at = NOW() WHERE is_active = true;