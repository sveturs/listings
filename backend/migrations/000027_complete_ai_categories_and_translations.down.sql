-- Откат миграции 000027_complete_ai_categories_and_translations

-- Удаление AI маппингов для новых категорий
DELETE FROM category_ai_mappings
WHERE ai_domain IN ('entertainment', 'nature', 'construction', 'aviation', 'military', 'other')
  AND product_type IN (
    'toy', 'doll', 'action_figure', 'puzzle', 'jigsaw', 'game', 'board_game', 'card_game',
    'acorn', 'wood', 'stone', 'mineral', 'seed', 'plant',
    'sand', 'cement', 'brick', 'tool', 'paint',
    'airplane', 'aircraft', 'helicopter',
    'uniform', 'equipment', 'medal',
    'unknown', 'misc'
  );

-- Удаление ключевых слов для новых категорий
DELETE FROM category_keyword_weights
WHERE category_id IN (
    SELECT id FROM marketplace_categories
    WHERE slug IN (
        'puzzles', 'toys', 'natural-materials', 'construction-materials',
        'aviation', 'military'
    )
);

-- Удаление переводов для новых категорий
DELETE FROM translations
WHERE entity_type = 'category'
  AND entity_id IN (
    SELECT id FROM marketplace_categories
    WHERE slug IN (
        'natural-materials', 'wood-materials', 'stones-minerals', 'plants-seeds', 'natural-decor',
        'construction-materials', 'bulk-materials', 'tools', 'paints', 'plumbing', 'electrical',
        'crafts', 'craft-materials', 'handmade', 'sewing', 'knitting',
        'toys', 'puzzles', 'board-games', 'collectibles', 'constructors', 'educational-games', 'models',
        'antiques', 'aviation', 'military', 'miscellaneous'
    )
  );

-- Удаление новых категорий (сначала подкатегории, потом родительские)
DELETE FROM marketplace_categories
WHERE slug IN (
    -- Подкатегории природных материалов
    'wood-materials', 'stones-minerals', 'plants-seeds', 'natural-decor',
    -- Подкатегории строительных материалов
    'bulk-materials', 'tools', 'paints', 'plumbing', 'electrical',
    -- Подкатегории декора и рукоделия
    'craft-materials', 'handmade', 'sewing', 'knitting',
    -- Подкатегории игрушек и хобби
    'board-games', 'constructors', 'educational-games', 'models'
);

-- Удаление основных категорий (если они были созданы этой миграцией)
DELETE FROM marketplace_categories
WHERE slug IN ('crafts', 'antiques', 'aviation', 'military', 'miscellaneous')
  AND NOT EXISTS (
    SELECT 1 FROM marketplace_listings WHERE category_id = marketplace_categories.id
  );