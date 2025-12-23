-- Migration: Remove meta_keywords for Lepota i zdravlje, Za bebe i decu, Nakit i satovi categories
-- Date: 2025-12-22
-- Description: Rollback SEO keywords

BEGIN;

-- Clear meta_keywords for all affected categories
UPDATE categories SET meta_keywords = '{}'::jsonb
WHERE slug IN (
    -- Lepota i zdravlje
    'anti-aging', 'anti-age-krema', 'nega-koze', 'krema-za-lice', 'serum-za-lice',
    'maska-za-lice', 'tonik-za-lice', 'mleko-za-ciscenje', 'krema-za-telo', 'krema-za-ruke',
    'zastita-od-sunca', 'krema-za-suncanje', 'higijena', 'gelovi-za-tusiranje', 'gel-za-tusiranje',
    'intimna-higijena', 'oralna-higijena', 'dekorativna-kozmetika', 'makeup-cetkice',
    'nega-kose', 'fene-profesionalni', 'pegla-kosa-turmalin', 'ionske-cetkice-kosa',
    'muska-nega', 'muski-stil', 'parfemi', 'luksuzna-kozmetika', 'organska-kozmetika',
    'lepota-aparati', 'dermaroller-mikronedling', 'led-maske-lice', 'ultrazvucni-cistaci',
    'laseri-uklanjanje-dlaka', 'terapije-svetlom', 'depilacija', 'manikir-i-pedikir',
    'spa-i-relaksacija', 'vitamini-i-suplementi', 'kolagen-preparati', 'biotin-kompleksi',
    'eterična-ulja', 'medicinski-proizvodi', 'decija-kozmetika',

    -- Za bebe i decu
    'oprema-za-bebe', 'bebi-monitori-video', 'kolevke-automatske', 'grejaci-flasica',
    'sterilizatori', 'bebi-jumperoo', 'kada-bebe-sklopljiva', 'nocnici',
    'nosiljke-ergonomske', 'zaštitne-ograde', 'namestaj-za-bebe', 'deciji-namestaj',
    'nega-i-higijena-beba', 'hrana-za-bebe', 'igracke-za-bebe', 'igracke-za-decu',
    'edukativni-tepisi', 'elektronika-za-decu', 'skolski-pribor', 'decija-odeca-bebe',
    'decija-obuca-bebe',

    -- Nakit i satovi
    'zlatni-nakit', 'zlatne-ogrlice', 'srebrni-nakit', 'srebrne-minđuše',
    'dijamanti', 'verenicko-prstenje', 'prstenje-veridba-dijamant', 'biserni-nakit',
    'perlice-narukvice', 'moderni-nakit', 'broše-vintage', 'privesci-lanci',
    'muski-satovi', 'zenski-satovi', 'pametni-satovi-fitness', 'luksuzni-satovi',
    'luksuzni-satovi-swiss', 'satovski-dodaci'
);

COMMIT;
