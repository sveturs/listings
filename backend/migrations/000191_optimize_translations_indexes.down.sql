-- Rollback: Restore dropped indexes on translations table

-- Restore metadata GIN index
CREATE INDEX idx_translations_metadata ON translations USING GIN (metadata);

-- Restore entity_field partial index (only for listings)
CREATE INDEX idx_translations_entity_field
ON translations (entity_type, entity_id, field_name, language)
WHERE entity_type = 'listing';
