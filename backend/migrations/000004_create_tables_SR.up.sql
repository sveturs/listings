-- Srpski prevodi kategorija za tabelu translations
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES

('category', 1000, 'sr', 'name', 'Nekretnine', true, true, NOW(), NOW()),
    -- Podkategorije nekretnina 1000
    ('category', 1100, 'sr', 'name', 'Stan', true, true, NOW(), NOW()),
    ('category', 1200, 'sr', 'name', 'Soba', true, true, NOW(), NOW()),
    ('category', 1300, 'sr', 'name', 'Kuća, vikendica, vila', true, true, NOW(), NOW()),
    ('category', 1400, 'sr', 'name', 'Zemljišna parcela', true, true, NOW(), NOW()),
    ('category', 1500, 'sr', 'name', 'Garaža i parking mesto', true, true, NOW(), NOW()),
    ('category', 1600, 'sr', 'name', 'Poslovni prostor', true, true, NOW(), NOW()),
    ('category', 1700, 'sr', 'name', 'Nekretnine u inostranstvu', true, true, NOW(), NOW()),
    ('category', 1800, 'sr', 'name', 'Hotel', true, true, NOW(), NOW()),
    ('category', 1900, 'sr', 'name', 'Apartmani', true, true, NOW(), NOW()),

('category', 2000, 'sr', 'name', 'Vozila', true, true, NOW(), NOW()),

    -- Auto kategorije 2000
    ('category', 2100, 'sr', 'name', 'Putnički automobili', true, true, NOW(), NOW()),
    ('category', 2200, 'sr', 'name', 'Komercijalna vozila', true, true, NOW(), NOW()),
        -- Podkategorije komercijalnih vozila 2200
        ('category', 2210, 'sr', 'name', 'Kamioni', true, true, NOW(), NOW()),
        ('category', 2220, 'sr', 'name', 'Poluprikolice', true, true, NOW(), NOW()),
        ('category', 2230, 'sr', 'name', 'Laka komercijalna vozila', true, true, NOW(), NOW()),
        ('category', 2240, 'sr', 'name', 'Autobusi', true, true, NOW(), NOW()),
        
    ('category', 2300, 'sr', 'name', 'Specijalna mehanizacija', true, true, NOW(), NOW()),
        -- Podkategorije specijalne mehanizacije 2300
        ('category', 2310, 'sr', 'name', 'Bageri', true, true, NOW(), NOW()),
        ('category', 2315, 'sr', 'name', 'Utovarivači', true, true, NOW(), NOW()),
        ('category', 2320, 'sr', 'name', 'Kombinovani bageri', true, true, NOW(), NOW()),
        ('category', 2325, 'sr', 'name', 'Auto-dizalice', true, true, NOW(), NOW()),
        ('category', 2330, 'sr', 'name', 'Automikseri', true, true, NOW(), NOW()),
        ('category', 2335, 'sr', 'name', 'Drumski valjci', true, true, NOW(), NOW()),
        ('category', 2340, 'sr', 'name', 'Čistilice', true, true, NOW(), NOW()),
        ('category', 2345, 'sr', 'name', 'Kamioni za smeće', true, true, NOW(), NOW()),
        ('category', 2350, 'sr', 'name', 'Autokorpe', true, true, NOW(), NOW()),
        ('category', 2355, 'sr', 'name', 'Buldožeri', true, true, NOW(), NOW()),
        ('category', 2360, 'sr', 'name', 'Grejderi', true, true, NOW(), NOW()),
        ('category', 2365, 'sr', 'name', 'Bušilice', true, true, NOW(), NOW()),

    ('category', 2400, 'sr', 'name', 'Poljoprivredna mehanizacija', true, true, NOW(), NOW()),
        -- Podkategorije poljoprivredne mehanizacije 2400
        ('category', 2410, 'sr', 'name', 'Traktori', true, true, NOW(), NOW()),
        ('category', 2415, 'sr', 'name', 'Mini traktori', true, true, NOW(), NOW()),
        ('category', 2420, 'sr', 'name', 'Balirke', true, true, NOW(), NOW()),
        ('category', 2425, 'sr', 'name', 'Drljače', true, true, NOW(), NOW()),
        ('category', 2430, 'sr', 'name', 'Kosačice', true, true, NOW(), NOW()),
        ('category', 2435, 'sr', 'name', 'Kombajni', true, true, NOW(), NOW()),
        ('category', 2440, 'sr', 'name', 'Teleskopski utovarivači', true, true, NOW(), NOW()),
        ('category', 2445, 'sr', 'name', 'Sejačice', true, true, NOW(), NOW()),
        ('category', 2450, 'sr', 'name', 'Kultivatori', true, true, NOW(), NOW()),
        ('category', 2455, 'sr', 'name', 'Plugovi', true, true, NOW(), NOW()),
        ('category', 2460, 'sr', 'name', 'Prskalice', true, true, NOW(), NOW()),

    ('category', 2500, 'sr', 'name', 'Iznajmljivanje vozila i mehanizacije', true, true, NOW(), NOW()),
            -- Iznajmljivanje vozila i mehanizacije (nadređena id=2500)
        ('category', 2510, 'sr', 'name', 'Automobili', true, true, NOW(), NOW()),
        ('category', 2515, 'sr', 'name', 'Oprema za podizanje', true, true, NOW(), NOW()),
        ('category', 2520, 'sr', 'name', 'Oprema za zemljane radove', true, true, NOW(), NOW()),
        ('category', 2525, 'sr', 'name', 'Komunalna mehanizacija', true, true, NOW(), NOW()),
        ('category', 2530, 'sr', 'name', 'Oprema za izgradnju puteva', true, true, NOW(), NOW()),
        ('category', 2535, 'sr', 'name', 'Teretni transport', true, true, NOW(), NOW()),
        ('category', 2540, 'sr', 'name', 'Oprema za utovar', true, true, NOW(), NOW()),
        ('category', 2545, 'sr', 'name', 'Priključci', true, true, NOW(), NOW()),
        ('category', 2550, 'sr', 'name', 'Prikolice', true, true, NOW(), NOW()),
        ('category', 2555, 'sr', 'name', 'Poljoprivredna oprema', true, true, NOW(), NOW()),
        ('category', 2560, 'sr', 'name', 'Kamp kućice', true, true, NOW(), NOW()),

    ('category', 2600, 'sr', 'name', 'Motocikli i motorna vozila', true, true, NOW(), NOW()),
        -- Motocikli i motorna vozila (nadređena id=2600)
        ('category', 2610, 'sr', 'name', 'Terenska vozila', true, true, NOW(), NOW()),
        ('category', 2615, 'sr', 'name', 'Karting', true, true, NOW(), NOW()),
        ('category', 2620, 'sr', 'name', 'Kvadovi i bagiji', true, true, NOW(), NOW()),
        ('category', 2625, 'sr', 'name', 'Mopedi', true, true, NOW(), NOW()),
        ('category', 2630, 'sr', 'name', 'Skuteri', true, true, NOW(), NOW()),
        ('category', 2635, 'sr', 'name', 'Motocikli', true, true, NOW(), NOW()),
        ('category', 2640, 'sr', 'name', 'Motorne sanke', true, true, NOW(), NOW()),

    ('category', 2700, 'sr', 'name', 'Plovila', true, true, NOW(), NOW()),
        -- Plovila (nadređena id=2700)
        ('category', 2710, 'sr', 'name', 'Čamci na vesla', true, true, NOW(), NOW()),
        ('category', 2720, 'sr', 'name', 'Kajaci', true, true, NOW(), NOW()),
        ('category', 2730, 'sr', 'name', 'Džet-ski', true, true, NOW(), NOW()),
        ('category', 2740, 'sr', 'name', 'Čamci i jahte', true, true, NOW(), NOW()),
        ('category', 2750, 'sr', 'name', 'Motorni čamci i motori', true, true, NOW(), NOW()),

    ('category', 2800, 'sr', 'name', 'Delovi i dodaci', true, true, NOW(), NOW()),
        -- Delovi i dodaci (nadređena id=2800)
        ('category', 2810, 'sr', 'name', 'Rezervni delovi', true, true, NOW(), NOW()),
        ('category', 2815, 'sr', 'name', 'Gume, felne i točkovi', true, true, NOW(), NOW()),
        ('category', 2820, 'sr', 'name', 'Audio i video oprema', true, true, NOW(), NOW()),
        ('category', 2825, 'sr', 'name', 'Dodaci', true, true, NOW(), NOW()),
        ('category', 2830, 'sr', 'name', 'Ulja i auto-hemija', true, true, NOW(), NOW()),
        ('category', 2835, 'sr', 'name', 'Alati', true, true, NOW(), NOW()),
        ('category', 2840, 'sr', 'name', 'Krovni nosači i kuke', true, true, NOW(), NOW()),
        ('category', 2845, 'sr', 'name', 'Prikolice', true, true, NOW(), NOW()),
        ('category', 2850, 'sr', 'name', 'Oprema', true, true, NOW(), NOW()),
        ('category', 2855, 'sr', 'name', 'Uređaji protiv krađe', true, true, NOW(), NOW()),
        ('category', 2860, 'sr', 'name', 'GPS navigacija', true, true, NOW(), NOW()),

('category', 3000, 'sr', 'name', 'Elektronika', true, true, NOW(), NOW()),
    -- Elektronika (nadređena id=3000)
    ('category', 3100, 'sr', 'name', 'Telefoni', true, true, NOW(), NOW()),
            -- Podkategorije telefona (nadređena id=3100)
        ('category', 3110, 'sr', 'name', 'Mobilni telefoni', true, true, NOW(), NOW()),
        ('category', 3120, 'sr', 'name', 'Dodaci', true, true, NOW(), NOW()),
                    -- Podkategorije dodataka za telefone (nadređena id=3120)
            ('category', 3121, 'sr', 'name', 'Baterije', true, true, NOW(), NOW()),
            ('category', 3122, 'sr', 'name', 'Slušalice', true, true, NOW(), NOW()),
            ('category', 3123, 'sr', 'name', 'Punjači', true, true, NOW(), NOW()),
            ('category', 3124, 'sr', 'name', 'Kablovi i adapteri', true, true, NOW(), NOW()),
            ('category', 3125, 'sr', 'name', 'Modemi i ruteri', true, true, NOW(), NOW()),
            ('category', 3126, 'sr', 'name', 'Maske i zaštitne folije', true, true, NOW(), NOW()),
            ('category', 3127, 'sr', 'name', 'Rezervni delovi', true, true, NOW(), NOW()),

        ('category', 3130, 'sr', 'name', 'Radio stanice', true, true, NOW(), NOW()),
        ('category', 3140, 'sr', 'name', 'Fiksni telefoni', true, true, NOW(), NOW()),

    ('category', 3200, 'sr', 'name', 'Audio i video', true, true, NOW(), NOW()),
            -- Podkategorije audio i video (nadređena id=3200)
        ('category', 3210, 'sr', 'name', 'Televizori i projektori', true, true, NOW(), NOW()),
        ('category', 3215, 'sr', 'name', 'Slušalice', true, true, NOW(), NOW()),
        ('category', 3220, 'sr', 'name', 'Zvučnici, sistemi, sabvuferi', true, true, NOW(), NOW()),
        ('category', 3225, 'sr', 'name', 'Dodaci', true, true, NOW(), NOW()),
        ('category', 3230, 'sr', 'name', 'Muzički centri, kasetofoni', true, true, NOW(), NOW()),
        ('category', 3235, 'sr', 'name', 'Pojačala i risiveri', true, true, NOW(), NOW()),
        ('category', 3240, 'sr', 'name', 'Video kamere', true, true, NOW(), NOW()),
        ('category', 3245, 'sr', 'name', 'Video, DVD i Blu-rej plejeri', true, true, NOW(), NOW()),
        ('category', 3250, 'sr', 'name', 'Kablovi i adapteri', true, true, NOW(), NOW()),
        ('category', 3255, 'sr', 'name', 'Muzika i filmovi', true, true, NOW(), NOW()),
        ('category', 3260, 'sr', 'name', 'Mikrofoni', true, true, NOW(), NOW()),
        ('category', 3265, 'sr', 'name', 'MP3 plejeri', true, true, NOW(), NOW()),

    ('category', 3300, 'sr', 'name', 'Računarska oprema', true, true, NOW(), NOW()),
    -- Podkategorije računarske opreme (nadređena id=3300)
        ('category', 3310, 'sr', 'name', 'Desktop računari', true, true, NOW(), NOW()),
        ('category', 3320, 'sr', 'name', 'Monitori', true, true, NOW(), NOW()),
        ('category', 3330, 'sr', 'name', 'Komponente', true, true, NOW(), NOW()),
        -- Podkategorije komponenti (nadređena id=3330)
            ('category', 3331, 'sr', 'name', 'CD, DVD i Blu-rej uređaji', true, true, NOW(), NOW()),
            ('category', 3332, 'sr', 'name', 'Napajanja', true, true, NOW(), NOW()),
            ('category', 3333, 'sr', 'name', 'Grafičke karte', true, true, NOW(), NOW()),
            ('category', 3334, 'sr', 'name', 'Hard diskovi', true, true, NOW(), NOW()),
            ('category', 3335, 'sr', 'name', 'Zvučne karte', true, true, NOW(), NOW()),
            ('category', 3336, 'sr', 'name', 'Kontroleri', true, true, NOW(), NOW()),
            ('category', 3337, 'sr', 'name', 'Kućišta', true, true, NOW(), NOW()),
            ('category', 3338, 'sr', 'name', 'Matične ploče', true, true, NOW(), NOW()),
            ('category', 3339, 'sr', 'name', 'RAM memorija', true, true, NOW(), NOW()),
            ('category', 3340, 'sr', 'name', 'Procesori', true, true, NOW(), NOW()),
            ('category', 3341, 'sr', 'name', 'Sistemi za hlađenje', true, true, NOW(), NOW()),

        ('category', 3360, 'sr', 'name', 'Monitori i delovi', true, true, NOW(), NOW()),
        ('category', 3365, 'sr', 'name', 'Mrežna oprema', true, true, NOW(), NOW()),
        ('category', 3370, 'sr', 'name', 'Tastature i miševi', true, true, NOW(), NOW()),
        ('category', 3375, 'sr', 'name', 'Džojstici i volani', true, true, NOW(), NOW()),
        ('category', 3380, 'sr', 'name', 'Fleš memorije i memorijske kartice', true, true, NOW(), NOW()),
        ('category', 3385, 'sr', 'name', 'Veb kamere', true, true, NOW(), NOW()),
        ('category', 3390, 'sr', 'name', 'TV tjuneri', true, true, NOW(), NOW()),

    ('category', 3500, 'sr', 'name', 'Igre, konzole i programi', true, true, NOW(), NOW()),
            -- Podkategorije igara, konzola i programa (nadređena id=3500)
        ('category', 3510, 'sr', 'name', 'Konzole za igre', true, true, NOW(), NOW()),
        ('category', 3520, 'sr', 'name', 'Igre za konzole', true, true, NOW(), NOW()),
        ('category', 3530, 'sr', 'name', 'Računarske igre', true, true, NOW(), NOW()),
        ('category', 3540, 'sr', 'name', 'Programi', true, true, NOW(), NOW()),

    ('category', 3600, 'sr', 'name', 'Laptopovi', true, true, NOW(), NOW()),
    ('category', 3700, 'sr', 'name', 'Foto oprema', true, true, NOW(), NOW()),
            -- Podkategorije foto opreme (nadređena id=3700)
        ('category', 3710, 'sr', 'name', 'Oprema i dodaci', true, true, NOW(), NOW()),
        ('category', 3720, 'sr', 'name', 'Objektivi', true, true, NOW(), NOW()),
        ('category', 3730, 'sr', 'name', 'Kompaktni fotoaparati', true, true, NOW(), NOW()),
        ('category', 3740, 'sr', 'name', 'Filmski fotoaparati', true, true, NOW(), NOW()),
        ('category', 3750, 'sr', 'name', 'DSLR fotoaparati', true, true, NOW(), NOW()),

    ('category', 3800, 'sr', 'name', 'Tableti i e-čitači', true, true, NOW(), NOW()),
            -- Podkategorije tableta i e-čitača (nadređena id=3800)
        ('category', 3810, 'sr', 'name', 'Tableti', true, true, NOW(), NOW()),
        ('category', 3820, 'sr', 'name', 'E-čitači', true, true, NOW(), NOW()),
        ('category', 3830, 'sr', 'name', 'Dodaci', true, true, NOW(), NOW()),

    ('category', 3900, 'sr', 'name', 'Kancelarijska oprema i potrošni materijal', true, true, NOW(), NOW()),
            -- Podkategorije kancelarijske opreme (nadređena id=3900)
        ('category', 3910, 'sr', 'name', 'MFP, kopir i skener uređaji', true, true, NOW(), NOW()),
        ('category', 3920, 'sr', 'name', 'Štampači', true, true, NOW(), NOW()),
        ('category', 3930, 'sr', 'name', 'Kancelarijski materijal', true, true, NOW(), NOW()),
        ('category', 3940, 'sr', 'name', 'UPS, filteri napajanja', true, true, NOW(), NOW()),
        ('category', 3950, 'sr', 'name', 'Telefonija', true, true, NOW(), NOW()),
        ('category', 3960, 'sr', 'name', 'Uništavači papira', true, true, NOW(), NOW()),
        ('category', 3970, 'sr', 'name', 'Potrošni materijal', true, true, NOW(), NOW()),

    ('category', 4100, 'sr', 'name', 'Kućni aparati', true, true, NOW(), NOW()),
        -- Podkategorije kućnih aparata (nadređena id=4100)
        ('category', 4110, 'sr', 'name', 'Kuhinjski aparati', true, true, NOW(), NOW()),
                    -- Podkategorije kuhinjskih aparata (nadređena id=4110)
            ('category', 4111, 'sr', 'name', 'Aspiratori', true, true, NOW(), NOW()),
            ('category', 4112, 'sr', 'name', 'Mali kuhinjski aparati', true, true, NOW(), NOW()),
            ('category', 4113, 'sr', 'name', 'Mikrotalasne pećnice', true, true, NOW(), NOW()),
            ('category', 4114, 'sr', 'name', 'Šporeti i rerne', true, true, NOW(), NOW()),
            ('category', 4115, 'sr', 'name', 'Mašine za pranje sudova', true, true, NOW(), NOW()),
            ('category', 4116, 'sr', 'name', 'Frižideri i zamrzivači', true, true, NOW(), NOW()),

        ('category', 4120, 'sr', 'name', 'Kućni aparati', true, true, NOW(), NOW()),
                    -- Podkategorije kućnih aparata (nadređena id=4120)
            ('category', 4121, 'sr', 'name', 'Usisivači i delovi', true, true, NOW(), NOW()),
            ('category', 4122, 'sr', 'name', 'Mašine za pranje i sušenje veša', true, true, NOW(), NOW()),
            ('category', 4123, 'sr', 'name', 'Pegle', true, true, NOW(), NOW()),
            ('category', 4124, 'sr', 'name', 'Oprema za šivenje', true, true, NOW(), NOW()),

        ('category', 4130, 'sr', 'name', 'Klimatizaciona oprema', true, true, NOW(), NOW()),
                    -- Podkategorije klimatizacione opreme (nadređena id=4130)
            ('category', 4131, 'sr', 'name', 'Ventilatori', true, true, NOW(), NOW()),
            ('category', 4132, 'sr', 'name', 'Klima uređaji i delovi', true, true, NOW(), NOW()),
            ('category', 4133, 'sr', 'name', 'Grejalice', true, true, NOW(), NOW()),
            ('category', 4134, 'sr', 'name', 'Prečišćivači vazduha', true, true, NOW(), NOW()),
            ('category', 4135, 'sr', 'name', 'Termometri i meteo stanice', true, true, NOW(), NOW()),

        ('category', 4140, 'sr', 'name', 'Aparati za ličnu negu', true, true, NOW(), NOW()),
            -- Podkategorije aparata za ličnu negu (nadređena id=4140)
            ('category', 4141, 'sr', 'name', 'Brijači i trimeri', true, true, NOW(), NOW()),
            ('category', 4142, 'sr', 'name', 'Mašinice za šišanje', true, true, NOW(), NOW()),
            ('category', 4143, 'sr', 'name', 'Fenovi i uređaji za oblikovanje', true, true, NOW(), NOW()),
            ('category', 4144, 'sr', 'name', 'Epilatori', true, true, NOW(), NOW()),
            
('category', 5000, 'sr', 'name', 'Sve za kuću i stan', true, true, NOW(), NOW()),
    -- Sve za kuću i stan (nadređena id=5000)
    -- Renoviranje i izgradnja
    ('category', 5100, 'sr', 'name', 'Renoviranje i izgradnja', true, true, NOW(), NOW()),
        -- Podkategorije renoviranja i izgradnje (nadređena id=5100)
        ('category', 5110, 'sr', 'name', 'Vrata', true, true, NOW(), NOW()),
        ('category', 5115, 'sr', 'name', 'Alati', true, true, NOW(), NOW()),
        ('category', 5120, 'sr', 'name', 'Kamini i grejalice', true, true, NOW(), NOW()),
        ('category', 5125, 'sr', 'name', 'Prozori i balkoni', true, true, NOW(), NOW()),
        ('category', 5130, 'sr', 'name', 'Plafoni', true, true, NOW(), NOW()),
        ('category', 5135, 'sr', 'name', 'Za baštu i vikendicu', true, true, NOW(), NOW()),
        ('category', 5140, 'sr', 'name', 'Vodovod, vodosnabdevanje i sauna', true, true, NOW(), NOW()),
        ('category', 5145, 'sr', 'name', 'Gotove konstrukcije i brvnare', true, true, NOW(), NOW()),
        ('category', 5150, 'sr', 'name', 'Kapije, ograde i barijere', true, true, NOW(), NOW()),
        ('category', 5155, 'sr', 'name', 'Sigurnosni i alarmni sistemi', true, true, NOW(), NOW()),

    ('category', 5200, 'sr', 'name', 'Nameštaj i enterijer', true, true, NOW(), NOW()),
        -- Podkategorije nameštaja i enterijera (nadređena id=5200)
        ('category', 5210, 'sr', 'name', 'Kreveti, kauči i fotelje', true, true, NOW(), NOW()),
        ('category', 5215, 'sr', 'name', 'Tekstil i tepisi', true, true, NOW(), NOW()),
        ('category', 5220, 'sr', 'name', 'Osvetljenje', true, true, NOW(), NOW()),
        ('category', 5225, 'sr', 'name', 'Kompjuterski stolovi i stolice', true, true, NOW(), NOW()),
        ('category', 5230, 'sr', 'name', 'Ormari, komode i police', true, true, NOW(), NOW()),
        ('category', 5235, 'sr', 'name', 'Kuhinjske garniture', true, true, NOW(), NOW()),
        ('category', 5240, 'sr', 'name', 'Stolovi i stolice', true, true, NOW(), NOW()),
        ('category', 5250, 'sr', 'name', 'Sobne biljke', true, true, NOW(), NOW()),
            -- Podkategorije sobnih biljaka (nadređena id=5250)
            ('category', 5251, 'sr', 'name', 'Dekorativno lišće biljaka', true, true, NOW(), NOW()),
            ('category', 5252, 'sr', 'name', 'Cvetne biljke', true, true, NOW(), NOW()),
            ('category', 5253, 'sr', 'name', 'Palme i fikusi', true, true, NOW(), NOW()),
            ('category', 5254, 'sr', 'name', 'Kaktusi i sukulenti', true, true, NOW(), NOW()),

('category', 5300, 'sr', 'name', 'Prehrambeni proizvodi', true, true, NOW(), NOW()),
            -- Podkategorije prehrambenih proizvoda (nadređena id=5300)
        ('category', 5310, 'sr', 'name', 'Čaj, kafa, kakao', true, true, NOW(), NOW()),
        ('category', 5315, 'sr', 'name', 'Pića', true, true, NOW(), NOW()),
        ('category', 5320, 'sr', 'name', 'Riba, morski plodovi, kavijar', true, true, NOW(), NOW()),
        ('category', 5325, 'sr', 'name', 'Meso, živina, iznutrice', true, true, NOW(), NOW()),
        ('category', 5330, 'sr', 'name', 'Slatkiši', true, true, NOW(), NOW()),
        ('category', 5340, 'sr', 'name', 'Rakija i vino', true, true, NOW(), NOW()),
                    -- Podkategorije rakije i vina (nadređena id=5340)
            ('category', 5341, 'sr', 'name', 'Šljivovica', true, true, NOW(), NOW()),
            ('category', 5342, 'sr', 'name', 'Lozovača', true, true, NOW(), NOW()),
            ('category', 5343, 'sr', 'name', 'Voćna rakija', true, true, NOW(), NOW()),
            ('category', 5344, 'sr', 'name', 'Domaće vino', true, true, NOW(), NOW()),

        ('category', 5350, 'sr', 'name', 'Domaći sirevi', true, true, NOW(), NOW()),
        ('category', 5360, 'sr', 'name', 'Kajmak', true, true, NOW(), NOW()),
        ('category', 5370, 'sr', 'name', 'Ajvar', true, true, NOW(), NOW()),

    ('category', 5400, 'sr', 'name', 'Posuđe i kuhinjski pribor', true, true, NOW(), NOW()),
        -- Podkategorije posuđa i kuhinjskog pribora (nadređena id=5400)
        ('category', 5405, 'sr', 'name', 'Posuđe', true, true, NOW(), NOW()),
        ('category', 5410, 'sr', 'name', 'Kuhinjski pribor', true, true, NOW(), NOW()),
        ('category', 5415, 'sr', 'name', 'Pribor za postavljanje stola', true, true, NOW(), NOW()),
        ('category', 5420, 'sr', 'name', 'Pribor za kuvanje', true, true, NOW(), NOW()),
        ('category', 5425, 'sr', 'name', 'Čuvanje hrane', true, true, NOW(), NOW()),
        ('category', 5430, 'sr', 'name', 'Priprema napitaka', true, true, NOW(), NOW()),
        ('category', 5435, 'sr', 'name', 'Kućna hemija', true, true, NOW(), NOW()),

('category', 6000, 'sr', 'name', 'Sve za baštu', true, true, NOW(), NOW()),
    -- Sve za baštu (nadređena id=6000, podkategorije 51-60)
    ('category', 6050, 'sr', 'name', 'Baštenski nameštaj', true, true, NOW(), NOW()),
    ('category', 6100, 'sr', 'name', 'Baštenski alati', true, true, NOW(), NOW()),
    ('category', 6200, 'sr', 'name', 'Seme i rasad', true, true, NOW(), NOW()),
    ('category', 6250, 'sr', 'name', 'Roštilj i dodaci', true, true, NOW(), NOW()),
    ('category', 6300, 'sr', 'name', 'Bazeni i oprema', true, true, NOW(), NOW()),
    ('category', 6350, 'sr', 'name', 'Sistemi za navodnjavanje', true, true, NOW(), NOW()),
    ('category', 6400, 'sr', 'name', 'Kompostiranje', true, true, NOW(), NOW()),
    ('category', 6450, 'sr', 'name', 'Staklenici i rasadnici', true, true, NOW(), NOW()),
    ('category', 6500, 'sr', 'name', 'Đubriva i zemlja', true, true, NOW(), NOW()),
    ('category', 6550, 'sr', 'name', 'Osvetljenje', true, true, NOW(), NOW()),
    ('category', 6600, 'sr', 'name', 'Uređenje enterijera', true, true, NOW(), NOW()),
    ('category', 6650, 'sr', 'name', 'Biljke i semena', true, true, NOW(), NOW()),
    ('category', 6700, 'sr', 'name', 'Bašta i povrtnjak', true, true, NOW(), NOW()),
    ('category', 6750, 'sr', 'name', 'Baštenske biljke', true, true, NOW(), NOW()),
            -- Podkategorije baštenskih biljaka (nadređena id=6750)
        ('category', 6751, 'sr', 'name', 'Dekorativno grmlje i drveće', true, true, NOW(), NOW()),
        ('category', 6752, 'sr', 'name', 'Četinari', true, true, NOW(), NOW()),
        ('category', 6753, 'sr', 'name', 'Višegodišnje biljke', true, true, NOW(), NOW()),
        ('category', 6754, 'sr', 'name', 'Voćke', true, true, NOW(), NOW()),
        ('category', 6755, 'sr', 'name', 'Travnjak', true, true, NOW(), NOW()),
        ('category', 6756, 'sr', 'name', 'Zeleniš i začini', true, true, NOW(), NOW()),

    ('category', 6850, 'sr', 'name', 'Semena, lukovice, gomolji', true, true, NOW(), NOW()),
    ('category', 6900, 'sr', 'name', 'Sredstva za negu biljaka', true, true, NOW(), NOW()),
            -- Podkategorije sredstava za negu biljaka (nadređena id=6900)
        ('category', 6901, 'sr', 'name', 'Zemlje i supstrati', true, true, NOW(), NOW()),
        ('category', 6902, 'sr', 'name', 'Đubriva', true, true, NOW(), NOW()),
        ('category', 6903, 'sr', 'name', 'Sredstva protiv štetočina i korova', true, true, NOW(), NOW()),
        ('category', 6904, 'sr', 'name', 'Saksije i žardinjere', true, true, NOW(), NOW()),
        ('category', 6905, 'sr', 'name', 'Fito lampe', true, true, NOW(), NOW()),
        ('category', 6906, 'sr', 'name', 'Merači vlage', true, true, NOW(), NOW()),

    ('category', 6950, 'sr', 'name', 'Panoi i veštačke biljke', true, true, NOW(), NOW()),

('category', 7000, 'sr', 'name', 'Hobi i rekreacija', true, true, NOW(), NOW()),
        -- Hobi i rekreacija (nadređena id=7000)
    ('category', 7050, 'sr', 'name', 'Muzički instrumenti', true, true, NOW(), NOW()),
        -- Muzički instrumenti (nadređena id=7050)
        ('category', 7055, 'sr', 'name', 'Žičani instrumenti', true, true, NOW(), NOW()),
        ('category', 7060, 'sr', 'name', 'Klaviri i klavijature', true, true, NOW(), NOW()),
        ('category', 7065, 'sr', 'name', 'Udaraljke', true, true, NOW(), NOW()),
        ('category', 7070, 'sr', 'name', 'Duvački instrumenti', true, true, NOW(), NOW()),
        ('category', 7075, 'sr', 'name', 'Harmonike', true, true, NOW(), NOW()),
        ('category', 7080, 'sr', 'name', 'Audio oprema', true, true, NOW(), NOW()),
        ('category', 7085, 'sr', 'name', 'Dodaci za instrumente', true, true, NOW(), NOW()),

    ('category', 7100, 'sr', 'name', 'Knjige i časopisi', true, true, NOW(), NOW()),
            -- Podkategorije knjiga i časopisa (nadređena id=7100)
        ('category', 7105, 'sr', 'name', 'Časopisi, novine, brošure', true, true, NOW(), NOW()),
        ('category', 7115, 'sr', 'name', 'Knjige', true, true, NOW(), NOW()),
        ('category', 7130, 'sr', 'name', 'Udžbenici', true, true, NOW(), NOW()),
    ('category', 7150, 'sr', 'name', 'Sportski rekviziti', true, true, NOW(), NOW()),
    ('category', 7250, 'sr', 'name', 'Kolekcionarstvo', true, true, NOW(), NOW()),
            -- Podkategorije kolekcionarstva (nadređena id=7250)
        ('category', 7251, 'sr', 'name', 'Novčanice', true, true, NOW(), NOW()),
        ('category', 7252, 'sr', 'name', 'Karte', true, true, NOW(), NOW()),
        ('category', 7253, 'sr', 'name', 'Stvari poznatih, autogrami', true, true, NOW(), NOW()),
        ('category', 7254, 'sr', 'name', 'Vojne stvari', true, true, NOW(), NOW()),
        ('category', 7255, 'sr', 'name', 'Gramofonske ploče', true, true, NOW(), NOW()),
        ('category', 7256, 'sr', 'name', 'Dokumenta', true, true, NOW(), NOW()),
        ('category', 7257, 'sr', 'name', 'Žetoni, medalje, značke', true, true, NOW(), NOW()),
        ('category', 7258, 'sr', 'name', 'Igre', true, true, NOW(), NOW()),
        ('category', 7259, 'sr', 'name', 'Kalendari', true, true, NOW(), NOW()),
        ('category', 7261, 'sr', 'name', 'Slike', true, true, NOW(), NOW()),
        ('category', 7262, 'sr', 'name', 'Poštanske marke', true, true, NOW(), NOW()),
        ('category', 7263, 'sr', 'name', 'Modeli', true, true, NOW(), NOW()),
        ('category', 7264, 'sr', 'name', 'Novčići', true, true, NOW(), NOW()),

    ('category', 7300, 'sr', 'name', 'Umetnički predmeti', true, true, NOW(), NOW()),
    ('category', 7400, 'sr', 'name', 'Bicikli', true, true, NOW(), NOW()),
            -- Podkategorije bicikala (nadređena id=7400)
        ('category', 7410, 'sr', 'name', 'BMX', true, true, NOW(), NOW()),
        ('category', 7415, 'sr', 'name', 'Gradski bicikli', true, true, NOW(), NOW()),
        ('category', 7420, 'sr', 'name', 'Bicikli za drum', true, true, NOW(), NOW()),
        ('category', 7425, 'sr', 'name', 'Dečiji bicikli', true, true, NOW(), NOW()),
        ('category', 7430, 'sr', 'name', 'Brdski bicikli', true, true, NOW(), NOW()),
        ('category', 7435, 'sr', 'name', 'Delovi i dodaci', true, true, NOW(), NOW()),

    ('category', 7500, 'sr', 'name', 'Lov i ribolov', true, true, NOW(), NOW()),
        -- Podkategorije lova i ribolova (nadređena id=7500)
        ('category', 7510, 'sr', 'name', 'Noževi, multialati, sekire', true, true, NOW(), NOW()),
        ('category', 7520, 'sr', 'name', 'Lov', true, true, NOW(), NOW()),
                    -- Podkategorije lova (nadređena id=7520)
            ('category', 7521, 'sr', 'name', 'Optički nišani', true, true, NOW(), NOW()),
            ('category', 7522, 'sr', 'name', 'Dodaci za optičke nišane', true, true, NOW(), NOW()),
            ('category', 7523, 'sr', 'name', 'Monokulari, dvogledi, daljinomeri', true, true, NOW(), NOW()),

        ('category', 7530, 'sr', 'name', 'Ribolov', true, true, NOW(), NOW()),
            -- Podkategorije ribolova (nadređena id=7530)
            ('category', 7531, 'sr', 'name', 'Štapovi, mašinice', true, true, NOW(), NOW()),
            ('category', 7551, 'sr', 'name', 'Mamci i pribor', true, true, NOW(), NOW()),
            ('category', 7571, 'sr', 'name', 'Sonari i oprema', true, true, NOW(), NOW()),

    ('category', 7650, 'sr', 'name', 'Kampovanje', true, true, NOW(), NOW()),
    ('category', 7700, 'sr', 'name', 'Antikviteti', true, true, NOW(), NOW()),
    ('category', 7750, 'sr', 'name', 'Karte, događaji i putovanja', true, true, NOW(), NOW()),
        -- Podkategorije karata, događaja i putovanja (nadređena id=7750)
        ('category', 7751, 'sr', 'name', 'Karte, kuponi', true, true, NOW(), NOW()),
        ('category', 7752, 'sr', 'name', 'Koncerti', true, true, NOW(), NOW()),
        ('category', 7753, 'sr', 'name', 'Putovanja', true, true, NOW(), NOW()),
        ('category', 7754, 'sr', 'name', 'Sport', true, true, NOW(), NOW()),
        ('category', 7755, 'sr', 'name', 'Pozorište, opera, balet', true, true, NOW(), NOW()),
        ('category', 7756, 'sr', 'name', 'Cirkus, bioskop', true, true, NOW(), NOW()),
        ('category', 7758, 'sr', 'name', 'Šou, mjuzikl', true, true, NOW(), NOW()),

    ('category', 7800, 'sr', 'name', 'Sport', true, true, NOW(), NOW()),
            -- Podkategorije sporta (nadređena id=7800)
        ('category', 7805, 'sr', 'name', 'Bilijar i kuglanje', true, true, NOW(), NOW()),
        ('category', 7810, 'sr', 'name', 'Ronjenje i vodeni sportovi', true, true, NOW(), NOW()),
        ('category', 7815, 'sr', 'name', 'Borilački sportovi', true, true, NOW(), NOW()),
        ('category', 7820, 'sr', 'name', 'Zimski sportovi', true, true, NOW(), NOW()),
        ('category', 7825, 'sr', 'name', 'Igre sa loptom', true, true, NOW(), NOW()),
        ('category', 7830, 'sr', 'name', 'Društvene igre', true, true, NOW(), NOW()),
        ('category', 7835, 'sr', 'name', 'Pejntbol i ersoft', true, true, NOW(), NOW()),
        ('category', 7840, 'sr', 'name', 'Roleri i skejtbord', true, true, NOW(), NOW()),
        ('category', 7845, 'sr', 'name', 'Tenis, badminton, ping-pong', true, true, NOW(), NOW()),
        ('category', 7850, 'sr', 'name', 'Turizam i boravak u prirodi', true, true, NOW(), NOW()),
        ('category', 7855, 'sr', 'name', 'Fitnes i trenažeri', true, true, NOW(), NOW()),
        ('category', 7860, 'sr', 'name', 'Sportska ishrana', true, true, NOW(), NOW()),

    -- Tradicionalni zanati i suveniri
    ('category', 7865, 'sr', 'name', 'Narodni zanati i rukotvorine', true, true, NOW(), NOW()),
        -- Podkategorije narodnih zanata (nadređena id=7865)
        ('category', 7866, 'sr', 'name', 'Opanci', true, true, NOW(), NOW()),
        ('category', 7867, 'sr', 'name', 'Keramika', true, true, NOW(), NOW()),
        ('category', 7868, 'sr', 'name', 'Vez', true, true, NOW(), NOW()),
        ('category', 7869, 'sr', 'name', 'Tkanje', true, true, NOW(), NOW()),
        ('category', 7870, 'sr', 'name', 'Narodni instrumenti', true, true, NOW(), NOW()),
        ('category', 7871, 'sr', 'name', 'Drvodelja', true, true, NOW(), NOW()),

    -- Poljoprivredne kategorije
    ('category', 7900, 'sr', 'name', 'Pčelarstvo', true, true, NOW(), NOW()),
        -- Podkategorije pčelarstva (nadređena id=7900)
        ('category', 7910, 'sr', 'name', 'Med', true, true, NOW(), NOW()),
        ('category', 7920, 'sr', 'name', 'Pčelinji vosak', true, true, NOW(), NOW()),
        ('category', 7930, 'sr', 'name', 'Propolis', true, true, NOW(), NOW()),
        ('category', 7935, 'sr', 'name', 'Pčelarski pribor', true, true, NOW(), NOW()),
        ('category', 7945, 'sr', 'name', 'Pčele', true, true, NOW(), NOW()),

    -- Turističke usluge
    ('category', 7950, 'sr', 'name', 'Seoski turizam', true, true, NOW(), NOW()),
            -- Podkategorije seoskog turizma (nadređena id=7950)
        ('category', 7951, 'sr', 'name', 'Etno-sela', true, true, NOW(), NOW()),
        ('category', 7952, 'sr', 'name', 'Vinske ture', true, true, NOW(), NOW()),
        ('category', 7953, 'sr', 'name', 'Agro turizam', true, true, NOW(), NOW()),
        ('category', 7954, 'sr', 'name', 'Planinski turizam', true, true, NOW(), NOW()),

('category', 8000, 'sr', 'name', 'Životinje', true, true, NOW(), NOW()),
    -- Kategorije životinja (nadređena id=8000)
    ('category', 8050, 'sr', 'name', 'Psi', true, true, NOW(), NOW()),
    ('category', 8100, 'sr', 'name', 'Mačke', true, true, NOW(), NOW()),
    ('category', 8150, 'sr', 'name', 'Ptice', true, true, NOW(), NOW()),
    ('category', 8200, 'sr', 'name', 'Akvarijum', true, true, NOW(), NOW()),
            -- Podkategorije akvarijuma (nadređena id=8200)
        ('category', 8205, 'sr', 'name', 'Akvarijumi', true, true, NOW(), NOW()),
        ('category', 8210, 'sr', 'name', 'Ribe', true, true, NOW(), NOW()),
        ('category', 8215, 'sr', 'name', 'Druge akvarijumske životinje', true, true, NOW(), NOW()),
        ('category', 8220, 'sr', 'name', 'Oprema', true, true, NOW(), NOW()),
        ('category', 8225, 'sr', 'name', 'Biljke', true, true, NOW(), NOW()),
        ('category', 8230, 'sr', 'name', 'Akvarijumski nameštaj', true, true, NOW(), NOW()),
        ('category', 8235, 'sr', 'name', 'Morska akvaristika', true, true, NOW(), NOW()),

    ('category', 8250, 'sr', 'name', 'Druge životinje', true, true, NOW(), NOW()),
        -- Podkategorije drugih životinja (nadređena id=8250)
        ('category', 8251, 'sr', 'name', 'Vodozemci', true, true, NOW(), NOW()),
        ('category', 8252, 'sr', 'name', 'Glodari', true, true, NOW(), NOW()),
        ('category', 8253, 'sr', 'name', 'Zečevi', true, true, NOW(), NOW()),
        ('category', 8254, 'sr', 'name', 'Konji', true, true, NOW(), NOW()),
        ('category', 8255, 'sr', 'name', 'Gmizavci', true, true, NOW(), NOW()),
        ('category', 8256, 'sr', 'name', 'Domaće životinje', true, true, NOW(), NOW()),
        ('category', 8257, 'sr', 'name', 'Živina', true, true, NOW(), NOW()),
        ('category', 8258, 'sr', 'name', 'Proizvodi za životinje', true, true, NOW(), NOW()),

('category', 8500, 'sr', 'name', 'Gotov biznis i oprema', true, true, NOW(), NOW()),


('category', 9000, 'sr', 'name', 'Posao', true, true, NOW(), NOW()),
    ('category', 9050, 'sr', 'name', 'Slobodna radna mesta', true, true, NOW(), NOW()),
    ('category', 9100, 'sr', 'name', 'Biografije', true, true, NOW(), NOW()),
    ('category', 9150, 'sr', 'name', 'Rad na daljinu', true, true, NOW(), NOW()),
    ('category', 9200, 'sr', 'name', 'Partnerstvo i saradnja', true, true, NOW(), NOW()),
    ('category', 9250, 'sr', 'name', 'Obuka i praksa', true, true, NOW(), NOW()),
    ('category', 9300, 'sr', 'name', 'Sezonski posao', true, true, NOW(), NOW()),
        -- Podkategorije sezonskog posla (nadređena id=9300)
        ('category', 9310, 'sr', 'name', 'Berba', true, true, NOW(), NOW()),
        ('category', 9315, 'sr', 'name', 'Rad u vinogradu', true, true, NOW(), NOW()),
        ('category', 9320, 'sr', 'name', 'Sezonski građevinski radovi', true, true, NOW(), NOW()),

('category', 9500, 'sr', 'name', 'Odeća, obuća, dodaci', true, true, NOW(), NOW()),

-- Dečije stvari i igračke (nova kategorija)
('category', 9700, 'sr', 'name', 'Dečije stvari i igračke', true, true, NOW(), NOW()),
    -- Podkategorije dečijih stvari i igračaka (nadređena id=9700)
    ('category', 9705, 'sr', 'name', 'Kolica za bebe', true, true, NOW(), NOW()),
    ('category', 9710, 'sr', 'name', 'Dečiji nameštaj', true, true, NOW(), NOW()),
    ('category', 9715, 'sr', 'name', 'Bicikli i trotineti', true, true, NOW(), NOW()),
    ('category', 9720, 'sr', 'name', 'Proizvodi za hranjenje', true, true, NOW(), NOW()),
    ('category', 9725, 'sr', 'name', 'Auto-sedišta', true, true, NOW(), NOW()),
    ('category', 9730, 'sr', 'name', 'Igračke', true, true, NOW(), NOW()),
    ('category', 9735, 'sr', 'name', 'Posteljina', true, true, NOW(), NOW()),
    ('category', 9740, 'sr', 'name', 'Proizvodi za kupanje', true, true, NOW(), NOW()),
    ('category', 9745, 'sr', 'name', 'Školski pribor', true, true, NOW(), NOW()),
    ('category', 9750, 'sr', 'name', 'Dečija odeća i obuća, dodaci', true, true, NOW(), NOW()),

('category', 9999, 'sr', 'name', 'Ostalo', true, true, NOW(), NOW());

-- Ажурирање sequence-а за translations
SELECT setval('translations_id_seq', (SELECT MAX(id) FROM translations), true);


-- Insert marketplace listings
INSERT INTO marketplace_listings (id, user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, views_count, created_at, updated_at, show_on_map, original_language) VALUES
(8, 2, 2100, 'Toyota Corolla 2018', 'Продајем Toyota Corolla 2018 годиште, 80.000 км, одлично стање. Први власник, редовно одржавање, сва документација доступна.', 1150000.00, 'used', 'active', 'Нови Сад, Србија', 45.26710000, 19.83350000, 'Нови Сад', 'Србија', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'sr'),
(9, 3, 3100, 'mobile Samsung Galaxy S21', 'Selling Samsung Galaxy S21, 256GB, Deep Purple. Perfect condition, complete set with original box and accessories. AppleCare+ until 2024.', 120000.00, 'used', 'active', 'Novi Sad, Serbia', 45.25510000, 19.84520000, 'Novi Sad', 'Serbia', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'en'),
(10, 4, 3310, 'Игровой компьютер RTX 4080', 'Продаю мощный игровой ПК: Intel Core i9-13900K, RTX 4080, 64GB RAM, 2TB NVMe SSD. Идеален для любых игр и тяжелых задач.', 350000.00, 'used', 'active', 'Нови-Сад, Сербия', 45.25410000, 19.84010000, 'Нови-Сад', 'Сербия', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'ru'),
(12, 2, 2100, 'автомобиль Toyota Corolla 2018', 'Продаю Toyota Corolla 2018 года, 80.000 км, отличное состояние. Первый владелец, регулярное обслуживание, вся документация в наличии.', 1475000.00, 'used', 'active', 'Косте Мајинског 4, Ветерник, Сербия', 45.24755670, 19.76878366, 'Ветерник', 'Сербия', 0, '2025-02-07 17:33:27.680035', '2025-02-07 17:40:23.957971', true, 'ru');

SELECT setval('marketplace_listings_id_seq', 12, true);

-- Insert marketplace images
INSERT INTO marketplace_images (id, listing_id, file_path, file_name, file_size, content_type, is_main, created_at) VALUES
(15, 8, 'toyota_1.jpg', 'toyota_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(16, 8, 'toyota_2.jpg', 'toyota_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(17, 9, 'galaxy_s21_1.jpg', 'galaxy_s21_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(18, 9, 'galaxy_s21_2.jpg', 'galaxy_s21_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(19, 10, 'gaming_pc_1.jpg', 'gaming_pc_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(20, 10, 'gaming_pc_2.jpg', 'gaming_pc_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(21, 12, 'toyota_1.jpg', 'toyota_1.jpg', 454842, 'image/jpeg', true, '2025-02-07 17:35:09.579393'),
(22, 12, 'toyota_2.jpg', 'toyota_2.jpg', 398035, 'image/jpeg', true, '2025-02-07 17:40:24.397595');

SELECT setval('marketplace_images_id_seq', 22, true);

-- Insert reviews and related data
INSERT INTO reviews (id, user_id, entity_type, entity_id, rating, comment, pros, cons, photos, likes_count, is_verified_purchase, status, created_at, updated_at, helpful_votes, not_helpful_votes, original_language) VALUES
(1, 2, 'listing', 8, 5, 'норм', NULL, NULL, NULL, 0, true, 'published', '2025-02-07 07:47:17.001726', '2025-02-07 14:25:23.586871', 0, 1, 'ru');

SELECT setval('reviews_id_seq', 1, true);

INSERT INTO review_responses (id, review_id, user_id, response, created_at, updated_at) VALUES
(1, 1, 3, 'ok', '2025-02-07 07:49:14.935308', '2025-02-07 07:49:14.935308');

SELECT setval('review_responses_id_seq', 1, true);

INSERT INTO review_votes (review_id, user_id, vote_type, created_at) VALUES
(1, 3, 'not_helpful', '2025-02-07 07:48:11.709016');
