-- Fix autoincrement for refresh_tokens table
ALTER TABLE refresh_tokens 
ALTER COLUMN id SET DEFAULT nextval('refresh_tokens_id_seq');

-- Ensure the sequence is owned by the table
ALTER SEQUENCE refresh_tokens_id_seq OWNED BY refresh_tokens.id;

-- Set the sequence to the correct value
SELECT setval('refresh_tokens_id_seq', COALESCE((SELECT MAX(id) FROM refresh_tokens), 0) + 1, false);