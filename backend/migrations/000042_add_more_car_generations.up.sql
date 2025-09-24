-- Добавляем популярные поколения автомобилей
-- Эта миграция расширяет данные для автомобильного раздела

-- BMW X5
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
((SELECT id FROM car_models WHERE name = 'X5' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'E53 (1999-2006)', 1999, 2006, 1, 'e53-1999-2006'),
((SELECT id FROM car_models WHERE name = 'X5' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'E70 (2007-2013)', 2007, 2013, 2, 'e70-2007-2013'),
((SELECT id FROM car_models WHERE name = 'X5' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'F15 (2013-2018)', 2013, 2018, 3, 'f15-2013-2018'),
((SELECT id FROM car_models WHERE name = 'X5' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'G05 (2018-)', 2018, NULL, 4, 'g05-2018-current');

-- BMW X3
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
((SELECT id FROM car_models WHERE name = 'X3' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'E83 (2003-2010)', 2003, 2010, 1, 'e83-2003-2010'),
((SELECT id FROM car_models WHERE name = 'X3' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'F25 (2010-2017)', 2010, 2017, 2, 'f25-2010-2017'),
((SELECT id FROM car_models WHERE name = 'X3' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'G01 (2017-)', 2017, NULL, 3, 'g01-2017-current');

-- BMW 5 Series
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
((SELECT id FROM car_models WHERE name = '5 Series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'E39 (1995-2003)', 1995, 2003, 1, 'e39-1995-2003'),
((SELECT id FROM car_models WHERE name = '5 Series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'E60/E61 (2003-2010)', 2003, 2010, 2, 'e60-e61-2003-2010'),
((SELECT id FROM car_models WHERE name = '5 Series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'F10/F11 (2010-2017)', 2010, 2017, 3, 'f10-f11-2010-2017'),
((SELECT id FROM car_models WHERE name = '5 Series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'G30/G31 (2016-)', 2016, NULL, 4, 'g30-g31-2016-current');

-- BMW 3 Series
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
((SELECT id FROM car_models WHERE name = '3 Series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'E46 (1998-2006)', 1998, 2006, 1, 'e46-1998-2006'),
((SELECT id FROM car_models WHERE name = '3 Series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'E90/E91/E92/E93 (2005-2012)', 2005, 2012, 2, 'e90-series-2005-2012'),
((SELECT id FROM car_models WHERE name = '3 Series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'F30/F31/F34/F35 (2012-2019)', 2012, 2019, 3, 'f30-series-2012-2019'),
((SELECT id FROM car_models WHERE name = '3 Series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1), 'G20/G21 (2019-)', 2019, NULL, 4, 'g20-g21-2019-current');

-- Audi A6
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
((SELECT id FROM car_models WHERE name = 'A6' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1), 'C5 (1997-2004)', 1997, 2004, 1, 'c5-1997-2004'),
((SELECT id FROM car_models WHERE name = 'A6' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1), 'C6 (2004-2011)', 2004, 2011, 2, 'c6-2004-2011'),
((SELECT id FROM car_models WHERE name = 'A6' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1), 'C7 (2011-2018)', 2011, 2018, 3, 'c7-2011-2018'),
((SELECT id FROM car_models WHERE name = 'A6' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1), 'C8 (2018-)', 2018, NULL, 4, 'c8-2018-current');

-- Audi A4
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
((SELECT id FROM car_models WHERE name = 'A4' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1), 'B5 (1994-2001)', 1994, 2001, 1, 'b5-1994-2001'),
((SELECT id FROM car_models WHERE name = 'A4' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1), 'B6 (2000-2005)', 2000, 2005, 2, 'b6-2000-2005'),
((SELECT id FROM car_models WHERE name = 'A4' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1), 'B7 (2004-2008)', 2004, 2008, 3, 'b7-2004-2008'),
((SELECT id FROM car_models WHERE name = 'A4' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1), 'B8 (2007-2015)', 2007, 2015, 4, 'b8-2007-2015'),
((SELECT id FROM car_models WHERE name = 'A4' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1), 'B9 (2015-)', 2015, NULL, 5, 'b9-2015-current');

-- Mercedes-Benz E-Class
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
((SELECT id FROM car_models WHERE name = 'E-Class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1), 'W210 (1995-2002)', 1995, 2002, 1, 'w210-1995-2002'),
((SELECT id FROM car_models WHERE name = 'E-Class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1), 'W211 (2002-2009)', 2002, 2009, 2, 'w211-2002-2009'),
((SELECT id FROM car_models WHERE name = 'E-Class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1), 'W212 (2009-2016)', 2009, 2016, 3, 'w212-2009-2016'),
((SELECT id FROM car_models WHERE name = 'E-Class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1), 'W213 (2016-)', 2016, NULL, 4, 'w213-2016-current');

-- Mercedes-Benz C-Class
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
((SELECT id FROM car_models WHERE name = 'C-Class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1), 'W202 (1993-2000)', 1993, 2000, 1, 'w202-1993-2000'),
((SELECT id FROM car_models WHERE name = 'C-Class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1), 'W203 (2000-2007)', 2000, 2007, 2, 'w203-2000-2007'),
((SELECT id FROM car_models WHERE name = 'C-Class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1), 'W204 (2007-2014)', 2007, 2014, 3, 'w204-2007-2014'),
((SELECT id FROM car_models WHERE name = 'C-Class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1), 'W205 (2014-2021)', 2014, 2021, 4, 'w205-2014-2021'),
((SELECT id FROM car_models WHERE name = 'C-Class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1), 'W206 (2021-)', 2021, NULL, 5, 'w206-2021-current');

-- Volkswagen Golf
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
((SELECT id FROM car_models WHERE name = 'Golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1), 'Golf III (1991-1997)', 1991, 1997, 1, 'golf-iii-1991-1997'),
((SELECT id FROM car_models WHERE name = 'Golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1), 'Golf IV (1997-2003)', 1997, 2003, 2, 'golf-iv-1997-2003'),
((SELECT id FROM car_models WHERE name = 'Golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1), 'Golf V (2003-2008)', 2003, 2008, 3, 'golf-v-2003-2008'),
((SELECT id FROM car_models WHERE name = 'Golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1), 'Golf VI (2008-2012)', 2008, 2012, 4, 'golf-vi-2008-2012'),
((SELECT id FROM car_models WHERE name = 'Golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1), 'Golf VII (2012-2019)', 2012, 2019, 5, 'golf-vii-2012-2019'),
((SELECT id FROM car_models WHERE name = 'Golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1), 'Golf VIII (2019-)', 2019, NULL, 6, 'golf-viii-2019-current');

-- Toyota Camry
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
((SELECT id FROM car_models WHERE name = 'Camry' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1), 'XV30 (2001-2006)', 2001, 2006, 1, 'xv30-2001-2006'),
((SELECT id FROM car_models WHERE name = 'Camry' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1), 'XV40 (2006-2011)', 2006, 2011, 2, 'xv40-2006-2011'),
((SELECT id FROM car_models WHERE name = 'Camry' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1), 'XV50 (2011-2017)', 2011, 2017, 3, 'xv50-2011-2017'),
((SELECT id FROM car_models WHERE name = 'Camry' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1), 'XV70 (2017-)', 2017, NULL, 4, 'xv70-2017-current');

-- Honda Civic
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
((SELECT id FROM car_models WHERE name = 'Civic' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1), '7th Gen (2001-2005)', 2001, 2005, 1, '7th-gen-2001-2005'),
((SELECT id FROM car_models WHERE name = 'Civic' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1), '8th Gen (2006-2011)', 2006, 2011, 2, '8th-gen-2006-2011'),
((SELECT id FROM car_models WHERE name = 'Civic' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1), '9th Gen (2012-2015)', 2012, 2015, 3, '9th-gen-2012-2015'),
((SELECT id FROM car_models WHERE name = 'Civic' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1), '10th Gen (2016-2021)', 2016, 2021, 4, '10th-gen-2016-2021'),
((SELECT id FROM car_models WHERE name = 'Civic' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1), '11th Gen (2021-)', 2021, NULL, 5, '11th-gen-2021-current');

-- Ford Focus
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
((SELECT id FROM car_models WHERE name = 'Focus' AND make_id = (SELECT id FROM car_makes WHERE slug = 'ford') LIMIT 1), '1st Gen (1998-2004)', 1998, 2004, 1, '1st-gen-1998-2004'),
((SELECT id FROM car_models WHERE name = 'Focus' AND make_id = (SELECT id FROM car_makes WHERE slug = 'ford') LIMIT 1), '2nd Gen (2004-2011)', 2004, 2011, 2, '2nd-gen-2004-2011'),
((SELECT id FROM car_models WHERE name = 'Focus' AND make_id = (SELECT id FROM car_makes WHERE slug = 'ford') LIMIT 1), '3rd Gen (2011-2018)', 2011, 2018, 3, '3rd-gen-2011-2018'),
((SELECT id FROM car_models WHERE name = 'Focus' AND make_id = (SELECT id FROM car_makes WHERE slug = 'ford') LIMIT 1), '4th Gen (2018-)', 2018, NULL, 4, '4th-gen-2018-current');