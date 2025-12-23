-- Rollback: Complete L2 Categories (Part 1/2)
-- Delete 18 newly added L2 categories

DELETE FROM categories WHERE slug IN (
  'bluetooth-zvucnici-premium',
  'bluetooth-slusalice-profesionalne',
  'wi-fi-routeri-mesh',
  'powerline-adapteri',
  'usb-hub-powered',
  'eksterni-hard-diskovi-ssd',
  'memorijske-kartice-pro',
  'citaci-e-ink',
  'elektronske-knjige-citalice',
  'bluetooth-tastature',
  'graficke-tablice',
  '3d-stampaci',
  'streaming-oprema',
  'pametne-sijalice',
  'punjaci-bezicni',
  'stabilizatori-gimbal',
  'mikroskopi-digitalni',
  'solarni-punjaci'
) AND level = 2;
