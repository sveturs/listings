-- Rollback Migration: Expand L2 categories (Part 6)
-- Date: 2025-12-17
-- Purpose: Remove L2 categories added in part 6 (partial)

-- Note: This migration was partially completed
-- Created: Kancelarija (15) + Muzički instrumenti (15) = 30 L2
-- Pending: Hrana (20) + Igračke (18) + Umetnost (12) = 50 L2

DELETE FROM categories WHERE slug IN (
  -- Kancelarijski materijal (15)
  'kancelarijska-hartija', 'stampaci-toneri', 'registratori-fascikle', 'olovke-i-markeri', 'beleznice-blokovi',
  'kancelarijski-pribor', 'organizacija-stola', 'kalkulator-i-oprema', 'bele-table-i-prezentacije', 'lomaci-dokumenta',
  'kancelarijski-namestaj', 'arhiviranje', 'laminating', 'korporativni-pokloni', 'skener-i-kopir',
  'projektna-oprema',

  -- Muzički instrumenti (15)
  'gitare-akusticne', 'elektricne-gitare', 'bas-gitare', 'bubnjevi-setovi', 'klavijature-i-pianos',
  'duvacki-instrumenti', 'violina-i-zicani', 'audio-oprema-za-muziku', 'mikrofoni-za-muziku', 'efekti-i-procesori',
  'dj-oprema', 'snimanje-i-produkcija', 'muzicki-dodaci', 'ukulele-i-mandoline', 'orgulje-i-harmonijum'
);

-- Progress notification
DO $$
BEGIN
  RAISE NOTICE 'Migration 20251217000003 PARTIALLY ROLLED BACK: Removed 30 L2 categories';
END $$;
