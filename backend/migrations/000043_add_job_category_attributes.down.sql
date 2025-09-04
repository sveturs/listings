-- Откат миграции 000043: Удаление атрибутов для категории Jobs

-- Удаляем связи категорий с атрибутами
DELETE FROM unified_category_attributes
WHERE attribute_id IN (
    SELECT id FROM unified_attributes 
    WHERE code IN ('job_type', 'job_category', 'salary_min', 'salary_max', 'salary_currency', 
                   'experience_years', 'education_level', 'remote_work', 'work_schedule', 'company_name')
);

-- Удаляем переводы атрибутов
DELETE FROM translations
WHERE entity_type = 'unified_attribute'
  AND entity_id IN (
    SELECT id FROM unified_attributes 
    WHERE code IN ('job_type', 'job_category', 'salary_min', 'salary_max', 'salary_currency', 
                   'experience_years', 'education_level', 'remote_work', 'work_schedule', 'company_name')
  );

-- Удаляем сами атрибуты
DELETE FROM unified_attributes 
WHERE code IN ('job_type', 'job_category', 'salary_min', 'salary_max', 'salary_currency', 
               'experience_years', 'education_level', 'remote_work', 'work_schedule', 'company_name');