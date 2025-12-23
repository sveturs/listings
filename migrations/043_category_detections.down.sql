DROP INDEX IF EXISTS idx_category_detections_created_at;
DROP INDEX IF EXISTS idx_category_detections_method;
DROP INDEX IF EXISTS idx_category_detections_detected_category;
DROP INDEX IF EXISTS idx_category_detections_user_confirmed;
DROP TABLE IF EXISTS category_detections;
-- НЕ удаляем индексы на categories - они могут использоваться другими
