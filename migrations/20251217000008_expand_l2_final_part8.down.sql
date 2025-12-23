-- Migration: Rollback expand L2 categories (Part 8 - FINAL)
-- Date: 2025-12-17
-- Purpose: Remove final 62 L2 categories from Usluge, Ljubimci, Knjige, Ostalo

DELETE FROM categories WHERE slug IN (
    -- Usluge (23 categories)
    'elektricar', 'vodoinstalater', 'moler-farbanje', 'stolar', 'bravar',
    'krovopokrivac', 'keramicar', 'parketar', 'majstor-general', 'servis-bele-tehnike',
    'servis-laptopa', 'servis-telefona', 'servis-tv', 'fotograf-event', 'fotograf-portret',
    'snimanje-video', 'montaza-video', 'web-dizajn', 'graficki-dizajn', 'prevod-jezik',
    'prepis-teksta', 'korepetitor-matematika', 'korepetitor-engleski',

    -- Kucni ljubimci (12 categories)
    'hrana-za-pse', 'hrana-za-macke', 'hrana-za-ptice', 'hrana-za-glodare',
    'igracke-za-pse', 'igracke-za-macke', 'oprema-za-setnju', 'kreveti-za-ljubimce',
    'toaleta-za-ljubimce', 'akvarijumi', 'terarijumi', 'nosilje-za-ljubimce',

    -- Knjige i mediji (10 categories)
    'biografije', 'istorijske-knjige', 'fantastika', 'sf-knjige', 'kriminalisticki-romani',
    'ljubavni-romani', 'horor', 'poezija', 'eseji', 'knjige-za-roditelje',

    -- Ostalo (17 categories)
    'antikviteti-namestaj', 'antikviteti-satovi', 'antikviteti-slike', 'vintage-odeca-muska',
    'vintage-odeca-zenska', 'retro-elektronika', 'kolekcionarske-karte', 'kolekcionarski-novcici',
    'markice-filatelija', 'minerali-i-kamenje', 'fosili', 'rucni-rad',
    'custom-proizvodi', 'personalizovani-pokloni', 'korporativni-pokloni-misc', 'razno',
    'nerazvrstavano'
)
AND level = 2;
