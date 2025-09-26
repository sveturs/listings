-- Добавляем поколения для реально существующих моделей
-- Эта миграция использует точные ID моделей из базы данных

-- BMW 3 Series (ID: 2752)
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
(2752, 'E46 (1998-2006)', 1998, 2006, 1, 'e46-1998-2006'),
(2752, 'E90/E91/E92/E93 (2005-2012)', 2005, 2012, 2, 'e90-series-2005-2012'),
(2752, 'F30/F31/F34/F35 (2012-2019)', 2012, 2019, 3, 'f30-series-2012-2019'),
(2752, 'G20/G21 (2019-)', 2019, NULL, 4, 'g20-g21-2019-current');

-- BMW 5 Series (ID: 2754)
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
(2754, 'E39 (1995-2003)', 1995, 2003, 1, 'e39-1995-2003'),
(2754, 'E60/E61 (2003-2010)', 2003, 2010, 2, 'e60-e61-2003-2010'),
(2754, 'F10/F11 (2010-2017)', 2010, 2017, 3, 'f10-f11-2010-2017'),
(2754, 'G30/G31 (2016-)', 2016, NULL, 4, 'g30-g31-2016-current');

-- Audi A4 (ID: 944)
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
(944, 'B5 (1994-2001)', 1994, 2001, 1, 'b5-1994-2001'),
(944, 'B6 (2000-2005)', 2000, 2005, 2, 'b6-2000-2005'),
(944, 'B7 (2004-2008)', 2004, 2008, 3, 'b7-2004-2008'),
(944, 'B8 (2007-2015)', 2007, 2015, 4, 'b8-2007-2015'),
(944, 'B9 (2015-)', 2015, NULL, 5, 'b9-2015-current');

-- Audi A6 (ID: 950)
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug) VALUES
(950, 'C5 (1997-2004)', 1997, 2004, 1, 'c5-1997-2004'),
(950, 'C6 (2004-2011)', 2004, 2011, 2, 'c6-2004-2011'),
(950, 'C7 (2011-2018)', 2011, 2018, 3, 'c7-2011-2018'),
(950, 'C8 (2018-)', 2018, NULL, 4, 'c8-2018-current');

-- Найдем и добавим больше моделей
-- Mercedes-Benz C-Class
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'W202 (1993-2000)', 1993, 2000, 1, 'w202-1993-2000'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'mercedes-benz' AND cm.name = 'C-Class' LIMIT 1;

INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'W203 (2000-2007)', 2000, 2007, 2, 'w203-2000-2007'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'mercedes-benz' AND cm.name = 'C-Class' LIMIT 1;

INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'W204 (2007-2014)', 2007, 2014, 3, 'w204-2007-2014'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'mercedes-benz' AND cm.name = 'C-Class' LIMIT 1;

INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'W205 (2014-2021)', 2014, 2021, 4, 'w205-2014-2021'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'mercedes-benz' AND cm.name = 'C-Class' LIMIT 1;

INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'W206 (2021-)', 2021, NULL, 5, 'w206-2021-current'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'mercedes-benz' AND cm.name = 'C-Class' LIMIT 1;

-- Mercedes-Benz E-Class
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'W210 (1995-2002)', 1995, 2002, 1, 'w210-1995-2002'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'mercedes-benz' AND cm.name = 'E-Class' LIMIT 1;

INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'W211 (2002-2009)', 2002, 2009, 2, 'w211-2002-2009'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'mercedes-benz' AND cm.name = 'E-Class' LIMIT 1;

INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'W212 (2009-2016)', 2009, 2016, 3, 'w212-2009-2016'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'mercedes-benz' AND cm.name = 'E-Class' LIMIT 1;

INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'W213 (2016-)', 2016, NULL, 4, 'w213-2016-current'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'mercedes-benz' AND cm.name = 'E-Class' LIMIT 1;

-- BMW X5
INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'E53 (1999-2006)', 1999, 2006, 1, 'e53-1999-2006'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'bmw' AND cm.name LIKE 'X5%' LIMIT 1;

INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'E70 (2007-2013)', 2007, 2013, 2, 'e70-2007-2013'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'bmw' AND cm.name LIKE 'X5%' LIMIT 1;

INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'F15 (2013-2018)', 2013, 2018, 3, 'f15-2013-2018'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'bmw' AND cm.name LIKE 'X5%' LIMIT 1;

INSERT INTO car_generations (model_id, name, year_start, year_end, sort_order, slug)
SELECT cm.id, 'G05 (2018-)', 2018, NULL, 4, 'g05-2018-current'
FROM car_models cm
JOIN car_makes mk ON cm.make_id = mk.id
WHERE mk.slug = 'bmw' AND cm.name LIKE 'X5%' LIMIT 1;