-- Restore original refresh_rating_views function
CREATE OR REPLACE FUNCTION public.refresh_rating_views()
RETURNS trigger
LANGUAGE plpgsql
AS $function$
BEGIN
    -- Обновляем только затронутые строки, а не всё представление
    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        -- Для пользователей
        IF NEW.entity_origin_type = 'user' THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
        END IF;

        -- Для магазинов
        IF NEW.entity_origin_type = 'storefront' THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        -- При удалении также обновляем
        IF OLD.entity_origin_type = 'user' THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY user_ratings;
        END IF;

        IF OLD.entity_origin_type = 'storefront' THEN
            REFRESH MATERIALIZED VIEW CONCURRENTLY storefront_ratings;
        END IF;
    END IF;

    RETURN NULL;
END;
$function$;

-- Drop materialized views
DROP MATERIALIZED VIEW IF EXISTS storefront_ratings CASCADE;
DROP MATERIALIZED VIEW IF EXISTS user_ratings CASCADE;
