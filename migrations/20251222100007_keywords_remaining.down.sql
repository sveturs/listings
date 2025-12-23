-- Migration rollback: Remove meta_keywords for remaining categories
-- Generated: 2025-12-22

-- Reset meta_keywords to NULL for all affected categories
UPDATE categories SET meta_keywords = NULL WHERE slug IN (
    -- Knjige i mediji
    'knjige-beletristika', 'decije-knjige', 'biografije', 'autobiografije', 'poezija',
    'eseji', 'kriminalisticki-romani', 'ljubavni-romani', 'fantastika', 'sf-knjige',
    'horor', 'istorijske-knjige', 'strucne-knjige', 'kuvari-recepti', 'knjige-za-roditelje',
    'audio-knjige', 'e-knjige', 'casopisi', 'stripovi', 'filmovi', 'muzika', 'retke-knjige',

    -- Hrana i piće
    'organsko-voce', 'organsko-povrce', 'mlecni-proizvodi-organik', 'zitarice-organik',
    'med-organski', 'organski-sokovi', 'maslinovo-ulje-premium', 'veganska-hrana',
    'vegetarijanska-hrana', 'bezglutenska-hrana', 'bez-laktoze', 'hrana-bez-secera',
    'keto-dijeta', 'paleo-dijeta', 'kafa-premium', 'caj-premium', 'proteini-u-prahu',
    'protein-pločice', 'sportska-ishrana-dodaci', 'superfood', 'zdrave-grickalice',
    'testenine-integralne',

    -- Kućni ljubimci
    'hrana-za-pse', 'hrana-za-macke', 'hrana-za-ptice', 'hrana-za-glodare',
    'igracke-za-pse', 'igracke-za-macke', 'oprema-za-setnju', 'kreveti-za-ljubimce',
    'toaleta-za-ljubimce', 'nosilje-za-ljubimce', 'akvarijumi', 'terarijumi', 'vodopasi',

    -- Industrija i alati
    'busilice-profesionalne', 'odvijaci-akumulatorski', 'udarni-odvijaci', 'brusilice-ugaone',
    'brusilice-vibracione', 'testeri-kruzne', 'testeri-lancane', 'rende-elektricne', 'cirkular',
    'generatori-benzin', 'generatori-dizel', 'kompresori-klipni', 'kompresori-ventilatorski',
    'zavarivaci-inverter', 'zavarivaci-mig', 'zavarivaci-tig', 'niveliri', 'meraci-laserni',
    'merila-digitalna', 'stege', 'akumulatori-profesionalni',

    -- Kancelarijski materijal
    'kancelarijski-pribor', 'olovke-i-markeri', 'beležnice-blokovi', 'kancelarijska-hartija',
    'registratori-fascikle', 'organizacija-stola', 'kalkulator-i-oprema', 'bele-table-i-prezentacije',
    'projektna-oprema', 'laminating', 'lomaci-dokumenta', 'skener-i-kopir', 'stampaci-toneri',
    'arhiviranje', 'kancelarijski-nameštaj',

    -- Muzički instrumenti
    'gitare-akusticne', 'elektricne-gitare', 'bas-gitare', 'ukulele-i-mandoline',
    'klavijature-i-pianos', 'orgulje-i-harmonijum', 'bubnjevi-setovi', 'duvacki-instrumenti',
    'violina-i-zicani', 'mikrofoni-za-muziku', 'audio-oprema-za-muziku', 'efekti-i-procesori',
    'dj-oprema', 'muzicki-dodaci', 'snimanje-i-produkcija',

    -- Umetnost i rukotvorine
    'boje-uljane', 'boje-akrilne', 'boje-vodene', 'boje-u-olovci', 'olovke-za-crtanje',
    'cetkice-za-slikanje', 'platna-za-slikanje', 'molberti', 'glina-za-modelovanje',
    'gips', 'konac-za-vez', 'platno-za-vez', 'perle-za-nakit', 'biser', 'decoupage-setovi',
    'rucni-rad',

    -- Usluge
    'majstor-general', 'elektricar', 'vodoinstalater', 'moler-farbanje', 'bravar',
    'stolar', 'keramicar', 'parketar', 'krovopokrivac', 'servis-bele-tehnike',
    'servis-tv', 'servis-laptopa', 'servis-telefona', 'fotograf-event', 'fotograf-portret',
    'snimanje-video', 'montaza-video', 'graficki-dizajn', 'web-dizajn', 'prevod-jezik',
    'prepis-teksta', 'korepetitor-matematika', 'korepetitor-engleski',

    -- Ostalo
    'kolekcionarske-karte', 'kolekcionarski-novcici', 'markice-filatelija', 'fosili',
    'minerali-i-kamenje', 'antikviteti-namestaj', 'antikviteti-satovi', 'antikviteti-slike',
    'vintage-odeca-muska', 'vintage-odeca-zenska', 'retro-elektronika', 'korporativni-pokloni',
    'korporativni-pokloni-misc', 'personalizovani-pokloni', 'custom-proizvodi', 'razno',
    'nerazvrstavano'
);
