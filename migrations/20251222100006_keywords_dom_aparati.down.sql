-- Откат миграции: Очистка meta_keywords для категорий Dom i bašta и Kućni aparati
-- Создано: 2025-12-22

-- Очистка meta_keywords для всех категорий Dom i bašta и Kućni aparati
UPDATE categories SET meta_keywords = NULL
WHERE slug IN (
    -- Kućni aparati
    'frizideri', 'zamrzivaci-vertikalni', 'vinske-vitrine',
    'masine-za-pranje', 'masine-pranje-susenje', 'masine-sudje', 'sudopere-i-masine',
    'sporet-i-rerna', 'indukcionе-ploce', 'mikotalasne-rerne', 'parne-rerne',
    'usisivaci', 'aspiratori-bez-kesice', 'roboti-usisivaci',
    'ventilacija-i-klimatizacija', 'ventilatori-i-grejalice', 'precistaci-vazduha', 'bojleri',
    'mali-kucni-aparati', 'masine-za-kafu', 'blenderi-profesionalni',
    'friteze', 'friteze-sa-vrucim-vazduhom', 'multi-kukeri', 'instant-lonci',
    'kuhinjski-pribor', 'kuhinjski-set-nozeva', 'drveno-posude-kuhinja', 'keramika-posude-rucno',
    'staklo-kristal-posude', 'kozarci-case-luksuz', 'kafe-servisi', 'caj-servisi',
    'tacne-posluzavnici', 'bar-oprema-dom',

    -- Dom i bašta - Nameštaj
    'namestaj-dnevna-soba', 'namestaj-spavaca-soba', 'namestaj-kuhinja', 'namestaj-kancelarija',
    'organizacija-i-skladistenje',

    -- Dom i bašta - Tekstil
    'tekstil-za-dom', 'tepihi-i-prostirke',

    -- Dom i bašta - Rasveta
    'rasveta', 'luster', 'plafonjera', 'zidna-lampa', 'stojeća-lampa', 'stonalampa',
    'led-sijalice', 'ugradna-led-rasveta', 'spoljna-rasveta', 'dekorativna-rasveta', 'pametna-rasveta',

    -- Dom i bašta - Dekoracije
    'dekoracije', 'vaze-i-dekor', 'ogledala', 'sat i-za-zid', 'pregradni-zidovi',

    -- Dom i bašta - Kupatilo
    'kupatilo', 'wc-solja', 'bidе', 'lavabo', 'ugradni-lavabo', 'slavina-za-lavabo',
    'kada', 'hidromasazna-kada', 'slavina-za-kadu', 'tus-kabina', 'tus-set',
    'kupatilski-namestaj', 'ogledalo-sa-ormarićem',

    -- Dom i bašta - Bašta
    'bastenska-garnitura', 'bastenska-oprema', 'alati-za-basta', 'alati-i-popravke',
    'bastenska-rasveta', 'bastenske-ukrase', 'grncari ja-i-biljke', 'kompostiranje', 'bazeni-i-spa'
);
