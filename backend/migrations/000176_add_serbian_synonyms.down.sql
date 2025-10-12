-- Откат добавления сербских синонимов
-- Удаляем только те синонимы, которые были добавлены этой миграцией

DELETE FROM search_synonyms WHERE language = 'sr' AND term IN (
    'stan', 'kuća', 'smeštaj', 'nekretnina', 'iznajmljivanje',
    'automobil', 'kola', 'auto', 'prevoz', 'putničko vozilo',
    'telefon', 'mobilni', 'laptop', 'računar', 'tablet',
    'posao', 'rad', 'zaposlenje', 'usluga', 'uslužni posao',
    'obrazovanje', 'kurs', 'instrukcije',
    'teretana', 'vežbanje',
    'ljubimac', 'pas', 'mačka',
    'odeća', 'obuća',
    'nameštaj', 'stolica', 'sto',
    'lepota', 'zdravlje',
    'hrana', 'piće', 'restoran',
    'putovanje', 'hotel',
    'zabava', 'muzika', 'film'
) AND id > 27;  -- Сохраняем старые синонимы (id <= 27)
