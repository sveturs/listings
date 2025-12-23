-- +migrate Down
-- Откат meta_keywords для категорий Automobilizam и Sport i turizam
-- Дата: 2025-12-22

-- Удаляем meta_keywords для всех указанных категорий
UPDATE categories SET meta_keywords = NULL WHERE slug IN (
    -- Automobilizam L2
    'akumulatori',
    'alati-za-automobile',
    'ambijentalno-osvetljenje',
    'audio-i-navigacija',
    'auto-aspiratori',
    'auto-dodaci',
    'auto-kozmetika',
    'auto-organizatori',
    'auto-pokrivaci',
    'auto-sedista-bebe',
    'dash-kamere',
    'delovi-za-automobile',
    'delovi-za-motocikle',
    'drzaci-telefona-auto',
    'tuniranje',
    'gps-trekeri-auto',
    'gume-i-felne',
    'klime-uređaji-prenosni',
    'led-trake-auto',
    'moto-oprema',
    'parking-senzori',
    'punjaci-elektricni-auto',

    -- Sport i turizam L2
    'badminton-oprema',
    'bicikli-i-trotineti',
    'borilacke-vestine-zastita',
    'dzonovanje',
    'fitnes-i-teretana',
    'fudbal',
    'joga-blokovi',
    'kajak-kanu-oprema',
    'kampovanje',
    'kosarka',
    'lov',
    'pilates-lopte',
    'planinarenje',
    'plivanje',
    'ribolov',
    'ronilacka-oprema-profesionalna',
    'stoni-tenis-reketi',
    'surfovanje-daske',
    'tenis',
    'tenis-reketi-pro',
    'trake-istezanje',
    'zimski-sportovi',

    -- Kampovanje L3
    'gas-reaud',
    'kamp-lampa',
    'kamp-ranac',
    'kamp-roštilj',
    'kamp-stolica',
    'prenosivi-frižider',
    'samoduvajuci-dušek',
    'sator-2-3-osobe',
    'sator-4-6-osoba',
    'spavaca-vreća'
);
