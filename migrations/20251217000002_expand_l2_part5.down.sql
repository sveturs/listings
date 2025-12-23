-- Rollback Migration: Expand L2 categories (Part 5)
-- Date: 2025-12-17
-- Purpose: Remove 80 L2 categories added in part 5

DELETE FROM categories WHERE slug IN (
  -- Odeća i obuća (10)
  'spec-odeca', 'uniforma', 'tradicionalna-odeca', 'premium-aksesoari', 'vintage-odeca',
  'sportske-uniforme', 'kucni-mantili', 'pletena-odeca', 'funkcionalna-odeca', 'modna-obuca',

  -- Elektronika (9, removed smart-home duplicate)
  'dronovi', 'vr-ar-oprema', 'roboti', '3d-printeri',
  'gaming-stolice', 'streaming-oprema', 'power-stanice', 'projektori-prenosni', 'elektro-transport',

  -- Dom i bašta (10)
  'premium-tekstil', 'premium-posudje', 'dekor-zidova', 'tepisi-i-circi', 'zavese-i-zastori',
  'pametna-rasveta', 'bastenska-premium-oprema', 'gril-i-bbq', 'bazeni-i-dzakuzi', 'sigurnost-za-dom',

  -- Lepota i zdravlje (9, removed organska-kozmetika duplicate)
  'masazeri', 'medicinski-aparati', 'pametne-vage', 'fitnes-trakeri-zdravlje', 'vitamini-premium',
  'ajurveda', 'anti-age', 'muska-nega', 'salonska-oprema',

  -- Za bebe i decu (10)
  'decija-oprema-namestaj', 'skolski-ranaci', 'muzicke-igracke', 'obrazovne-igre', 'decija-elektronika',
  'decija-kozmetika', 'artikli-za-novorodjence', 'kupanje-bebe', 'deciji-tekstil', 'dekoracija-decije-sobe',

  -- Sport i turizam (10)
  'joga-i-pilates', 'boks-i-borilacke-vestine', 'plivanje', 'tenis', 'odbojka',
  'kosarka', 'fudbal-oprema', 'trcanje-i-atletika', 'ekstremni-sportovi', 'lov-i-ribolov',

  -- Automobilizam (10)
  'elektromobil-aksesoari', 'tjuning', 'autozvuk-premium', 'autokozmetika-premium', 'krovni-nosaci',
  'auto-elektronika-premium', 'gps-navigacija', 'video-registratori', 'parking-senzori', 'auto-cehli-premium',

  -- Kućni aparati (10)
  'roboti-usisivaci', 'pametni-frizidera', 'vinski-frizideri', 'sudopere-premium', 'ves-masine-premium',
  'susare-masine', 'parne-stanice', 'multivarka', 'aerofriteze', 'hlebopekare'
);

-- Progress notification
DO $$
BEGIN
  RAISE NOTICE 'Migration 20251217000002 ROLLED BACK: Removed 78 L2 categories (2 duplicates removed from this migration)';
END $$;
