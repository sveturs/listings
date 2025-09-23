-- Миграция для заполнения таблицы car_generations
-- Добавляем поколения для популярных моделей автомобилей

-- BMW 3 Series (model_id из таблицы car_models)
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, external_id) VALUES
((SELECT id FROM car_models WHERE slug = '3-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'E21', 'e21', 1975, 1983, 'bmw_3_e21'),
((SELECT id FROM car_models WHERE slug = '3-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'E30', 'e30', 1982, 1994, 'bmw_3_e30'),
((SELECT id FROM car_models WHERE slug = '3-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'E36', 'e36', 1990, 2000, 'bmw_3_e36'),
((SELECT id FROM car_models WHERE slug = '3-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'E46', 'e46', 1998, 2006, 'bmw_3_e46'),
((SELECT id FROM car_models WHERE slug = '3-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'E90/E91/E92/E93', 'e90', 2005, 2013, 'bmw_3_e90'),
((SELECT id FROM car_models WHERE slug = '3-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'F30/F31/F34', 'f30', 2011, 2019, 'bmw_3_f30'),
((SELECT id FROM car_models WHERE slug = '3-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'G20/G21', 'g20', 2018, NULL, 'bmw_3_g20');

-- BMW 5 Series
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, external_id) VALUES
((SELECT id FROM car_models WHERE slug = '5-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'E28', 'e28', 1981, 1988, 'bmw_5_e28'),
((SELECT id FROM car_models WHERE slug = '5-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'E34', 'e34', 1988, 1996, 'bmw_5_e34'),
((SELECT id FROM car_models WHERE slug = '5-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'E39', 'e39', 1995, 2003, 'bmw_5_e39'),
((SELECT id FROM car_models WHERE slug = '5-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'E60/E61', 'e60', 2003, 2010, 'bmw_5_e60'),
((SELECT id FROM car_models WHERE slug = '5-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'F10/F11', 'f10', 2010, 2017, 'bmw_5_f10'),
((SELECT id FROM car_models WHERE slug = '5-series' AND make_id = (SELECT id FROM car_makes WHERE slug = 'bmw') LIMIT 1),
 'G30/G31', 'g30', 2016, NULL, 'bmw_5_g30');

-- Mercedes-Benz C-Class
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, external_id) VALUES
((SELECT id FROM car_models WHERE slug = 'c-class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1),
 'W202', 'w202', 1993, 2000, 'mb_c_w202'),
((SELECT id FROM car_models WHERE slug = 'c-class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1),
 'W203', 'w203', 2000, 2007, 'mb_c_w203'),
((SELECT id FROM car_models WHERE slug = 'c-class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1),
 'W204', 'w204', 2007, 2014, 'mb_c_w204'),
((SELECT id FROM car_models WHERE slug = 'c-class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1),
 'W205', 'w205', 2014, 2021, 'mb_c_w205'),
((SELECT id FROM car_models WHERE slug = 'c-class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1),
 'W206', 'w206', 2021, NULL, 'mb_c_w206');

-- Mercedes-Benz E-Class
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, external_id) VALUES
((SELECT id FROM car_models WHERE slug = 'e-class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1),
 'W124', 'w124', 1984, 1995, 'mb_e_w124'),
((SELECT id FROM car_models WHERE slug = 'e-class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1),
 'W210', 'w210', 1995, 2002, 'mb_e_w210'),
((SELECT id FROM car_models WHERE slug = 'e-class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1),
 'W211', 'w211', 2002, 2009, 'mb_e_w211'),
((SELECT id FROM car_models WHERE slug = 'e-class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1),
 'W212', 'w212', 2009, 2016, 'mb_e_w212'),
((SELECT id FROM car_models WHERE slug = 'e-class' AND make_id = (SELECT id FROM car_makes WHERE slug = 'mercedes-benz') LIMIT 1),
 'W213', 'w213', 2016, NULL, 'mb_e_w213');

-- Audi A4
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, external_id) VALUES
((SELECT id FROM car_models WHERE slug = 'a4' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1),
 'B5', 'b5', 1994, 2001, 'audi_a4_b5'),
((SELECT id FROM car_models WHERE slug = 'a4' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1),
 'B6', 'b6', 2000, 2006, 'audi_a4_b6'),
((SELECT id FROM car_models WHERE slug = 'a4' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1),
 'B7', 'b7', 2004, 2009, 'audi_a4_b7'),
((SELECT id FROM car_models WHERE slug = 'a4' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1),
 'B8', 'b8', 2007, 2015, 'audi_a4_b8'),
((SELECT id FROM car_models WHERE slug = 'a4' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1),
 'B9', 'b9', 2015, NULL, 'audi_a4_b9');

-- Audi A6
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, external_id) VALUES
((SELECT id FROM car_models WHERE slug = 'a6' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1),
 'C5', 'c5', 1997, 2005, 'audi_a6_c5'),
((SELECT id FROM car_models WHERE slug = 'a6' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1),
 'C6', 'c6', 2004, 2011, 'audi_a6_c6'),
((SELECT id FROM car_models WHERE slug = 'a6' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1),
 'C7', 'c7', 2011, 2018, 'audi_a6_c7'),
((SELECT id FROM car_models WHERE slug = 'a6' AND make_id = (SELECT id FROM car_makes WHERE slug = 'audi') LIMIT 1),
 'C8', 'c8', 2018, NULL, 'audi_a6_c8');

-- Volkswagen Golf
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, external_id) VALUES
((SELECT id FROM car_models WHERE slug = 'golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'Mk1', 'mk1', 1974, 1983, 'vw_golf_mk1'),
((SELECT id FROM car_models WHERE slug = 'golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'Mk2', 'mk2', 1983, 1991, 'vw_golf_mk2'),
((SELECT id FROM car_models WHERE slug = 'golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'Mk3', 'mk3', 1991, 1997, 'vw_golf_mk3'),
((SELECT id FROM car_models WHERE slug = 'golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'Mk4', 'mk4', 1997, 2003, 'vw_golf_mk4'),
((SELECT id FROM car_models WHERE slug = 'golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'Mk5', 'mk5', 2003, 2008, 'vw_golf_mk5'),
((SELECT id FROM car_models WHERE slug = 'golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'Mk6', 'mk6', 2008, 2013, 'vw_golf_mk6'),
((SELECT id FROM car_models WHERE slug = 'golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'Mk7', 'mk7', 2012, 2019, 'vw_golf_mk7'),
((SELECT id FROM car_models WHERE slug = 'golf' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'Mk8', 'mk8', 2019, NULL, 'vw_golf_mk8');

-- Volkswagen Passat
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, external_id) VALUES
((SELECT id FROM car_models WHERE slug = 'passat' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'B3', 'b3', 1988, 1993, 'vw_passat_b3'),
((SELECT id FROM car_models WHERE slug = 'passat' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'B4', 'b4', 1993, 1997, 'vw_passat_b4'),
((SELECT id FROM car_models WHERE slug = 'passat' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'B5', 'b5', 1996, 2005, 'vw_passat_b5'),
((SELECT id FROM car_models WHERE slug = 'passat' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'B6', 'b6', 2005, 2010, 'vw_passat_b6'),
((SELECT id FROM car_models WHERE slug = 'passat' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'B7', 'b7', 2010, 2015, 'vw_passat_b7'),
((SELECT id FROM car_models WHERE slug = 'passat' AND make_id = (SELECT id FROM car_makes WHERE slug = 'volkswagen') LIMIT 1),
 'B8', 'b8', 2014, NULL, 'vw_passat_b8');

-- Toyota Camry
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, external_id) VALUES
((SELECT id FROM car_models WHERE slug = 'camry' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1),
 'XV20', 'xv20', 1996, 2001, 'toyota_camry_xv20'),
((SELECT id FROM car_models WHERE slug = 'camry' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1),
 'XV30', 'xv30', 2001, 2006, 'toyota_camry_xv30'),
((SELECT id FROM car_models WHERE slug = 'camry' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1),
 'XV40', 'xv40', 2006, 2011, 'toyota_camry_xv40'),
((SELECT id FROM car_models WHERE slug = 'camry' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1),
 'XV50', 'xv50', 2011, 2017, 'toyota_camry_xv50'),
((SELECT id FROM car_models WHERE slug = 'camry' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1),
 'XV70', 'xv70', 2017, NULL, 'toyota_camry_xv70');

-- Toyota Corolla
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, external_id) VALUES
((SELECT id FROM car_models WHERE slug = 'corolla' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1),
 'E110', 'e110', 1995, 2002, 'toyota_corolla_e110'),
((SELECT id FROM car_models WHERE slug = 'corolla' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1),
 'E120', 'e120', 2000, 2007, 'toyota_corolla_e120'),
((SELECT id FROM car_models WHERE slug = 'corolla' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1),
 'E140/E150', 'e140', 2006, 2013, 'toyota_corolla_e140'),
((SELECT id FROM car_models WHERE slug = 'corolla' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1),
 'E170', 'e170', 2013, 2018, 'toyota_corolla_e170'),
((SELECT id FROM car_models WHERE slug = 'corolla' AND make_id = (SELECT id FROM car_makes WHERE slug = 'toyota') LIMIT 1),
 'E210', 'e210', 2018, NULL, 'toyota_corolla_e210');

-- Honda Civic
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, external_id) VALUES
((SELECT id FROM car_models WHERE slug = 'civic' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1),
 '6th Gen', 'ek', 1995, 2000, 'honda_civic_6'),
((SELECT id FROM car_models WHERE slug = 'civic' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1),
 '7th Gen', 'em', 2000, 2005, 'honda_civic_7'),
((SELECT id FROM car_models WHERE slug = 'civic' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1),
 '8th Gen', 'fd', 2005, 2011, 'honda_civic_8'),
((SELECT id FROM car_models WHERE slug = 'civic' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1),
 '9th Gen', 'fb', 2011, 2016, 'honda_civic_9'),
((SELECT id FROM car_models WHERE slug = 'civic' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1),
 '10th Gen', 'fc', 2015, 2021, 'honda_civic_10'),
((SELECT id FROM car_models WHERE slug = 'civic' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1),
 '11th Gen', 'fe', 2021, NULL, 'honda_civic_11');

-- Honda Accord
INSERT INTO car_generations (model_id, name, slug, year_start, year_end, external_id) VALUES
((SELECT id FROM car_models WHERE slug = 'accord' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1),
 '6th Gen', 'cf', 1997, 2002, 'honda_accord_6'),
((SELECT id FROM car_models WHERE slug = 'accord' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1),
 '7th Gen', 'cl', 2002, 2007, 'honda_accord_7'),
((SELECT id FROM car_models WHERE slug = 'accord' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1),
 '8th Gen', 'cu', 2007, 2012, 'honda_accord_8'),
((SELECT id FROM car_models WHERE slug = 'accord' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1),
 '9th Gen', 'cr', 2012, 2017, 'honda_accord_9'),
((SELECT id FROM car_models WHERE slug = 'accord' AND make_id = (SELECT id FROM car_makes WHERE slug = 'honda') LIMIT 1),
 '10th Gen', 'cv', 2017, NULL, 'honda_accord_10');