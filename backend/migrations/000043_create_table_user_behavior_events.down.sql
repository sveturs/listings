-- Drop table: user_behavior_events
DROP SEQUENCE IF EXISTS public.user_behavior_events_id_seq;
DROP TABLE IF EXISTS public.user_behavior_events;
DROP INDEX IF EXISTS public.idx_user_behavior_events_created_at;
DROP INDEX IF EXISTS public.idx_user_behavior_events_event_type;
DROP INDEX IF EXISTS public.idx_user_behavior_events_item;
DROP INDEX IF EXISTS public.idx_user_behavior_events_search_query_type;
DROP INDEX IF EXISTS public.idx_user_behavior_events_session_id;
DROP INDEX IF EXISTS public.idx_user_behavior_events_user_id;