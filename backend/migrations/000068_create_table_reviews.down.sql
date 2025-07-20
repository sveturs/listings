-- Drop table: reviews
DROP SEQUENCE IF EXISTS public.reviews_id_seq;
DROP TABLE IF EXISTS public.reviews;
DROP INDEX IF EXISTS public.idx_reviews_entity;
DROP INDEX IF EXISTS public.idx_reviews_entity_origin;
DROP INDEX IF EXISTS public.idx_reviews_has_dispute;
DROP INDEX IF EXISTS public.idx_reviews_rating;
DROP INDEX IF EXISTS public.idx_reviews_seller_confirmed;
DROP INDEX IF EXISTS public.idx_reviews_status;
DROP INDEX IF EXISTS public.idx_reviews_user;
DROP INDEX IF EXISTS public.idx_reviews_user_entity_unique;
DROP TRIGGER IF EXISTS refresh_rating_summaries_trigger ON public.reviews;
DROP TRIGGER IF EXISTS trigger_refresh_rating_distributions ON public.reviews;
DROP TRIGGER IF EXISTS update_ratings_after_review_change ON public.reviews;
DROP TRIGGER IF EXISTS update_reviews_updated_at ON public.reviews;