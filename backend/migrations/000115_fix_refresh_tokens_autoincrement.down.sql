-- Remove autoincrement from refresh_tokens table
ALTER TABLE refresh_tokens 
ALTER COLUMN id DROP DEFAULT;