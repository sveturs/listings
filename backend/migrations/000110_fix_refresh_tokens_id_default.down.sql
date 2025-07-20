-- Revert fix for refresh_tokens table ID auto-increment
ALTER TABLE ONLY public.refresh_tokens ALTER COLUMN id DROP DEFAULT;