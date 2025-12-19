-- Rollback: L3 Ostalo Categories

DELETE FROM categories WHERE slug IN (
  -- Kosa i frizura (10)
  'fen-za-kosu',
  'peglazakosu',
  'figaro',
  'trimer-za-kosu',
  'masnica-za-kosu',
  'sampon',
  'regenerator-za-kosu',
  'maska-za-kosu',
  'farba-za-kosu',
  'sprej-za-kosu',

  -- Nega koze (10)
  'mleko-za-ciscenje',
  'tonik-za-lice',
  'krema-za-lice',
  'serum-za-lice',
  'maska-za-lice',
  'krema-za-suncanje',
  'anti-age-krema',
  'krema-za-telo',
  'gel-za-tusiranje',
  'krema-za-ruke',

  -- Bebe oprema (12)
  'kolica-za-bebe',
  'nosiljka-za-bebe',
  'auto-sediste-za-bebe',
  'krevetac-za-bebe',
  'hranilica',
  'baby-alarm',
  'bocice-za-bebe',
  'pumpa-za-mleko',
  'pelene',
  'kadica-za-kupanje',
  'bebe-igracke',
  'ogradica-za-bebe',

  -- Auto delovi (10)
  'gume-ljetne',
  'gume-zimske',
  'aluminijske-felne',
  'akumulator',
  'auto-farovi',
  'led-auto-sijalice',
  'kocione-plocice',
  'uljni-filter',
  'vazdusni-filter',
  'brisaci',

  -- Aparati za kucu (10)
  'masina-za-pranje-vesa',
  'masina-za-susenje',
  'usisivac',
  'robotski-usisivac',
  'pegla',
  'pegla-sa-parom',
  'klima-uredjaj',
  'grejalica',
  'preciscivac-vazduha',
  'razvlazivac',

  -- Knjige (6)
  'romani',
  'price',
  'biografije',
  'popularna-nauka',
  'decije-knjige',
  'stripovi',

  -- Hrana za ljubimce (6)
  'hrana-za-pse-suva',
  'hrana-za-pse-konzerve',
  'hrana-za-macke-suva',
  'hrana-za-macke-konzerve',
  'hrana-za-ptice',
  'hrana-za-ribe'
) AND level = 3;
