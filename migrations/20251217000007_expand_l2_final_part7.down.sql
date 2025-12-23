-- Migration: Rollback expand L2 categories (Part 7)
-- Date: 2025-12-17
-- Purpose: Remove 60 L2 categories from Hrana, Umetnost, Industrija

DELETE FROM categories WHERE slug IN (
    -- Hrana i pice (22 categories)
    'organsko-voce', 'organsko-povrce', 'mlecni-proizvodi-organik', 'veganska-hrana',
    'vegetarijanska-hrana', 'bezglutenska-hrana', 'hrana-bez-secera', 'bez-laktoze',
    'superfood', 'proteini-u-prahu', 'protein-plocice', 'sportska-ishrana-dodaci',
    'zitarice-organik', 'testenine-integralne', 'keto-dijeta', 'paleo-dijeta',
    'zdrave-grickalice', 'organski-sokovi', 'kafa-premium', 'caj-premium',
    'med-organski', 'maslinovo-ulje-premium',

    -- Umetnost i rukotvorine (15 categories)
    'boje-akrilne', 'boje-uljane', 'boje-vodene', 'cetkice-za-slikanje',
    'platna-za-slikanje', 'molberti', 'olovke-za-crtanje', 'boje-u-olovci',
    'konac-za-vez', 'platno-za-vez', 'perle-za-nakit', 'biser',
    'glina-za-modelovanje', 'gips', 'decoupage-setovi',

    -- Industrija i alati (23 categories)
    'busilice-profesionalne', 'brusilice-ugaone', 'brusilice-vibracione', 'testeri-kruzne',
    'testeri-lancane', 'cirkular', 'rende-elektricne', 'odvijaci-akumulatorski',
    'udarni-odvijaci', 'akumulatori-profesionalni', 'kompresori-klipni', 'kompresori-ventilatorski',
    'generatori-benzin', 'generatori-dizel', 'zavarivaci-inverter', 'zavarivaci-tig',
    'zavarivaci-mig', 'meraci-laserni', 'vodopasi', 'niveliri',
    'merila-digitalna', 'stege', 'tesle'
)
AND level = 2;
