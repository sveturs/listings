-- Добавление сербских синонимов для улучшения поиска

-- Недвижимость и жилье
INSERT INTO search_synonyms (term, synonym, language, is_active) VALUES
('stan', 'garsonjera', 'sr', true),
('stan', 'dvosoban', 'sr', true),
('stan', 'trosoban', 'sr', true),
('stan', 'četvorosoban', 'sr', true),
('kuća', 'kuća', 'sr', true),
('kuća', 'vikendica', 'sr', true),
('kuća', 'rezidencija', 'sr', true),
('smeštaj', 'stan', 'sr', true),
('smeštaj', 'apartman', 'sr', true),
('smeštaj', 'boravište', 'sr', true),
('nekretnina', 'imovina', 'sr', true),
('nekretnina', 'stan', 'sr', true),
('iznajmljivanje', 'rentiranje', 'sr', true),
('iznajmljivanje', 'najam', 'sr', true),

-- Автомобили
('automobil', 'vozilo', 'sr', true),
('automobil', 'motorno vozilo', 'sr', true),
('kola', 'automobil', 'sr', true),
('kola', 'vozilo', 'sr', true),
('auto', 'vozilo', 'sr', true),
('auto', 'kola', 'sr', true),
('prevoz', 'transport', 'sr', true),
('putničko vozilo', 'automobil', 'sr', true),

-- Электроника
('telefon', 'pametni telefon', 'sr', true),
('telefon', 'mobilni uređaj', 'sr', true),
('mobilni', 'telefon', 'sr', true),
('mobilni', 'smart telefon', 'sr', true),
('laptop', 'prenosni računar', 'sr', true),
('laptop', 'notebook', 'sr', true),
('računar', 'kompjuter', 'sr', true),
('računar', 'PC', 'sr', true),
('tablet', 'tablet računar', 'sr', true),

-- Работа и услуги
('posao', 'zaposlenje', 'sr', true),
('posao', 'poslovanje', 'sr', true),
('rad', 'posao', 'sr', true),
('rad', 'zaposlenje', 'sr', true),
('zaposlenje', 'posao', 'sr', true),
('zaposlenje', 'rad', 'sr', true),
('usluga', 'servis', 'sr', true),
('usluga', 'pomoć', 'sr', true),
('uslužni posao', 'usluga', 'sr', true),

-- Образование
('obrazovanje', 'edukacija', 'sr', true),
('obrazovanje', 'školovanje', 'sr', true),
('kurs', 'obuka', 'sr', true),
('kurs', 'časovi', 'sr', true),
('instrukcije', 'časovi', 'sr', true),
('instrukcije', 'nastava', 'sr', true),

-- Спорт и фитнес
('teretana', 'fitnes', 'sr', true),
('teretana', 'gym', 'sr', true),
('vežbanje', 'trening', 'sr', true),
('vežbanje', 'fitnes', 'sr', true),

-- Животные
('ljubimac', 'kućni ljubimac', 'sr', true),
('ljubimac', 'životinja', 'sr', true),
('pas', 'kućni pas', 'sr', true),
('mačka', 'kućna mačka', 'sr', true),

-- Мода и одежда
('odeća', 'garderoboa', 'sr', true),
('odeća', 'odela', 'sr', true),
('obuća', 'cipele', 'sr', true),
('obuća', 'obućа', 'sr', true),

-- Мебель и интерьер
('nameštaj', 'mobilija', 'sr', true),
('nameštaj', 'enterijer', 'sr', true),
('stolica', 'sedište', 'sr', true),
('sto', 'radni sto', 'sr', true),

-- Красота и здоровье
('lepota', 'kozmetika', 'sr', true),
('lepota', 'nега', 'sr', true),
('zdravlje', 'medicina', 'sr', true),
('zdravlje', 'medicinska nega', 'sr', true),

-- Еда и напитки
('hrana', 'jelo', 'sr', true),
('hrana', 'obrok', 'sr', true),
('piće', 'napitak', 'sr', true),
('restoran', 'lokal', 'sr', true),
('restoran', 'restoracija', 'sr', true),

-- Путешествия и туризм
('putovanje', 'turizam', 'sr', true),
('putovanje', 'ekskurzija', 'sr', true),
('hotel', 'smeštaj', 'sr', true),
('hotel', 'prenoćište', 'sr', true),

-- Развлечения
('zabava', 'entertainment', 'sr', true),
('zabava', 'razonoda', 'sr', true),
('muzika', 'muzički sadržaj', 'sr', true),
('film', 'bioskop', 'sr', true),
('film', 'filmski sadržaj', 'sr', true)

ON CONFLICT (term, synonym, language) WHERE is_active = true DO NOTHING;
