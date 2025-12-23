-- Rollback: Remove meta_keywords for Odeća i obuća category tree

UPDATE categories SET meta_keywords = NULL WHERE slug IN (
    -- Level 2
    'bade-mantili', 'decija-obuca', 'decija-odeca', 'dodaci-i-aksesoari', 'donji-ves',
    'elegantna-odeca', 'esarpe-i-salovi', 'kosmeticke-uniforme', 'kosulјe-kratkih-rukava',
    'kucne-halјine', 'kupaci-kostimi', 'muska-obuca', 'muska-odeca', 'muskarci-veliki-brojevi',
    'muske-velike-velicine', 'naocari-i-dodaci', 'negl ize-i-spavacice', 'odeca-ronjenje',
    'odeca-trudnice-svecana', 'odela-i-smokingzi', 'pidzhame-svila', 'plus-size-odeca',
    'posteljina-i-peskiri', 'radna-odeca', 'sportska-odeca', 'termo-ves', 'torbice-i-novcanici',
    'trudnicka-odeca', 'vencana-odeca', 'ves-masine-dodaci', 'zene-veliki-brojevi',
    'zenska-obuca', 'zenska-odeca', 'zenske-velike-velicine', 'zimska-garderoba',

    -- Level 3 - Dečija obuća
    'bebe-obuca-0-6-meseci', 'bebe-obuca-6-12-meseci', 'decija-obuca-prve-korake',
    'decije-baletanke', 'decije-cipele-skolske', 'decije-cizme-duboke', 'decije-cizme-gumene',
    'decije-cizme-zimske', 'decije-fudbalske-kopacke', 'decije-papuce', 'decije-patike-1-3-godine',
    'decije-patike-4-7-godina', 'decije-patike-8-12-godina', 'decije-sandale', 'decije-sportske-patike',

    -- Level 3 - Dečija odeća
    'bebe-odeca-0-3-meseca', 'bebe-odeca-3-6-meseci', 'bebe-odeca-6-12-meseci',
    'decaci-odeca-1-3-godine', 'decaci-odeca-4-7-godina', 'decaci-odeca-8-12-godina',
    'devojcice-odeca-1-3-godine', 'devojcice-odeca-4-7-godina', 'devojcice-odeca-8-12-godina',
    'decije-jakne', 'decije-kaputi', 'decije-trenerke', 'deciji-duksevi', 'deciji-kupaci-kostimi',
    'skolska-uniforma', 'svecana-decija-odeca', 'kupaći-majice-uv-zaštita',

    -- Level 3 - Muška obuća
    'muske-cipele-brodske', 'muske-cipele-derby', 'muske-cipele-koza', 'muske-cipele-oxford',
    'muske-cizme-chelsea', 'muske-cizme-duboke', 'muske-cizme-radne', 'muske-espadrile',
    'muske-mokasine', 'muske-papuce', 'muske-patike-basketball', 'muske-patike-casual',
    'muske-patike-football', 'muske-patike-running', 'muske-sandale',

    -- Level 3 - Muška odeća
    'muska-poslovna-odeca', 'muska-sportska-odeca-l3', 'muske-jakne', 'muske-kosulje',
    'muske-majice', 'muske-pantalone', 'muski-dž emperi', 'muski-kupaci-slip',
    'muski-kupaci-sorcevi', 'muski-sorcevi', 'pareo-tunike',

    -- Level 3 - Ženska obuća
    'zenske-balerinke', 'zenske-cipele-platforme', 'zenske-cipele-potpetica',
    'zenske-cizme-duboke', 'zenske-cizme-gležnjače', 'zenske-cizme-preko-kolena',
    'zenske-espadrile', 'zenske-mokasine', 'zenske-natikace', 'zenske-papuce',
    'zenske-patike-casual', 'zenske-patike-fitness', 'zenske-patike-running',
    'zenske-sandale', 'zenske-stikle',

    -- Level 3 - Ženska odeća
    'zenska-poslovna-garderoba', 'zenska-poslovna-odeca', 'zenska-sportska-odeca',
    'zenska-vecernja-garderoba', 'zenske-bluze', 'zenske-haljine', 'zenske-jakne',
    'zenske-majice', 'zenske-mantile', 'zenske-pantalone', 'zenske-pantalone-elegantne',
    'zenske-pantalone-jeans', 'zenske-suknje', 'zenske-trenerke', 'zenski-duksevi',
    'zenski-dzemperi', 'zenski-džemperi', 'zenski-kaputi', 'zenski-monokini',
    'zenski-sorcevi', 'zenski-tankini'
);
