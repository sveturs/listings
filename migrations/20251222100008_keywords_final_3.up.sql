-- Final 3 L3 categories missing keywords
-- Generated 2025-12-22

BEGIN;

-- 1. Lomači dokumenta (Document Shredders) - parent: kancelarijski-materijal
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'document shredder, paper shredder, office shredder, cross-cut shredder, micro-cut shredder, strip-cut shredder, security shredder, confidential document disposal',
    'sr', 'lomač dokumenata, uništavač papira, kancelarijski lomač, rezač dokumenata, uništavač dokumenata, poverljivi dokumenti, bezbednosni lomač',
    'ru', 'уничтожитель документов, шредер, измельчитель бумаги, офисный шредер, перекрестная резка, конфиденциальные документы, уничтожение бумаг'
) WHERE slug = 'lomači-dokumenta';

-- 2. Mašine za suđe (Dishwashers) - parent: kucni-aparati
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'dishwasher, dish washing machine, built-in dishwasher, freestanding dishwasher, compact dishwasher, countertop dishwasher, energy efficient dishwasher',
    'sr', 'mašina za sudove, mašina za pranje sudova, ugradna mašina za sudove, sudopera, sudomašina, kompaktna mašina za sudove',
    'ru', 'посудомоечная машина, посудомойка, встраиваемая посудомойка, отдельностоящая посудомойка, компактная посудомойка, настольная посудомойка'
) WHERE slug = 'masine-sudjе';

-- 3. Tesle (Adzes/Axes) - parent: industrija-i-alati
UPDATE categories SET meta_keywords = jsonb_build_object(
    'en', 'adze, hand adze, carpenters adze, woodworking adze, timber framing tool, hewing tool, shaping tool, traditional woodworking',
    'sr', 'tesla, tesle, stolarska tesla, drvodeljska tesla, ručna tesla, alat za tesanje, tradicionalni alat, obrada drveta',
    'ru', 'тесло, плотницкое тесло, столярное тесло, ручное тесло, инструмент для тесания, деревообрабатывающий инструмент, традиционный инструмент'
) WHERE slug = 'tesle';

COMMIT;
