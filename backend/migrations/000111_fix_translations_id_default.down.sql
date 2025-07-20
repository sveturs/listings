-- Revert fix for translations table ID auto-increment
ALTER TABLE ONLY public.translations ALTER COLUMN id DROP DEFAULT;