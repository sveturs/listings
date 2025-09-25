-- Откат изменений (возврат к предыдущему состоянию)
-- Внимание: это вернет неправильные категории обратно!

UPDATE ai_category_decisions
SET category_id = 1002,
    reasoning = REPLACE(reasoning, ' [FIXED: Corrected from Fashion(1002) to Cars(1301)]', '')
WHERE category_id = 1301
  AND reasoning LIKE '%[FIXED: Corrected from Fashion(1002) to Cars(1301)]%';

-- Восстановление альтернативных категорий не требуется, так как
-- мы не можем точно определить какие из них были изменены