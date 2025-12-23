# Meta Keywords для категорий Automobilizam и Sport i turizam

## Статус: ✅ Завершено

**Дата:** 2025-12-22
**Миграция:** `20251222100004_keywords_auto_sport.up.sql`

---

## Охват

| Категория | Уровень | Количество | Покрытие |
|-----------|---------|------------|----------|
| Automobilizam L2 | 2 | 22 | 100% |
| Sport i turizam L2 | 2 | 22 | 100% |
| Kampovanje L3 | 3 | 10 | 100% |
| **ИТОГО** | - | **54** | **100%** |

---

## Примеры keywords

### Automobilizam

#### Gume i felne
- **SR:** gume, auto gume, letnje gume, zimske gume, all season, run flat, felne, alu felge
- **EN:** tires, car tires, summer tires, winter tires, all season, run flat, alloy wheels, rims
- **RU:** шины, автошины, летние, зимние, всесезонные, run flat, диски, литые диски

#### Akumulatori
- **SR:** akumulatori za auto, akumulator 12V, start stop, bezodržavajući, punjivi, AGM, EFB
- **EN:** car batteries, accumulator, 12V battery, start stop battery, AGM, EFB, maintenance free
- **RU:** автомобильные аккумуляторы, 12V батарея, start stop, необслуживаемый, AGM, EFB

### Sport i turizam

#### Fudbal
- **SR:** fudbal, lopta, kopačke, štitnici, golmanske rukavice, dres, trening
- **EN:** football, soccer ball, cleats, shin guards, goalkeeper gloves, football jersey, training
- **RU:** футбол, мяч, бутсы, щитки, вратарские перчатки, форма, тренировка

#### Kampovanje (L3 пример: šator 2-3 osobe)
- **SR:** šator 2-3 osobe, mali šator, kamp šator, lagan šator, vodootporan, kupola šator
- **EN:** tent 2-3 person, small tent, camping tent, lightweight tent, waterproof, dome tent
- **RU:** палатка 2-3 человека, маленькая палатка, лёгкая палатка, водонепроницаемая

---

## Категории Automobilizam (L2)

1. akumulatori
2. alati-za-automobile
3. ambijentalno-osvetljenje
4. audio-i-navigacija
5. auto-aspiratori
6. auto-dodaci
7. auto-kozmetika
8. auto-organizatori
9. auto-pokrivaci
10. auto-sedista-bebe
11. dash-kamere
12. delovi-za-automobile
13. delovi-za-motocikle
14. drzaci-telefona-auto
15. tuniranje
16. gps-trekeri-auto
17. gume-i-felne
18. klime-uređaji-prenosni
19. led-trake-auto
20. moto-oprema
21. parking-senzori
22. punjaci-elektricni-auto

---

## Категории Sport i turizam (L2)

1. badminton-oprema
2. bicikli-i-trotineti
3. borilacke-vestine-zastita
4. dzonovanje
5. fitnes-i-teretana
6. fudbal
7. joga-blokovi
8. kajak-kanu-oprema
9. kampovanje
10. kosarka
11. lov
12. pilates-lopte
13. planinarenje
14. plivanje
15. ribolov
16. ronilacka-oprema-profesionalna
17. stoni-tenis-reketi
18. surfovanje-daske
19. tenis
20. tenis-reketi-pro
21. trake-istezanje
22. zimski-sportovi

---

## Категории Kampovanje (L3)

1. gas-reaud
2. kamp-lampa
3. kamp-ranac
4. kamp-roštilj
5. kamp-stolica
6. prenosivi-frižider
7. samoduvajuci-dušek
8. sator-2-3-osobe
9. sator-4-6-osoba
10. spavaca-vreća

---

## Применение миграции

```bash
# Миграция применена напрямую через psql
docker exec -i listings_postgres psql -U listings_user -d listings_dev_db < migrations/20251222100004_keywords_auto_sport.up.sql

# Зарегистрирована в schema_migrations
INSERT INTO schema_migrations (version, dirty) VALUES (20251222100004, false);
```

---

## Проверка результата

```sql
-- Проверить keywords для конкретной категории
SELECT slug,
       meta_keywords->>'sr' as keywords_sr,
       meta_keywords->>'en' as keywords_en,
       meta_keywords->>'ru' as keywords_ru
FROM categories
WHERE slug = 'gume-i-felne';

-- Проверить покрытие
WITH RECURSIVE cat_tree AS (
    SELECT id, slug, parent_id, meta_keywords, 1 as level
    FROM categories WHERE slug IN ('automobilizam', 'sport-i-turizam')
    UNION ALL
    SELECT c.id, c.slug, c.parent_id, c.meta_keywords, ct.level + 1
    FROM categories c
    JOIN cat_tree ct ON c.parent_id = ct.id
)
SELECT
    level,
    COUNT(*) as total,
    COUNT(CASE WHEN meta_keywords::text != '{}'::jsonb::text THEN 1 END) as with_keywords
FROM cat_tree
WHERE level > 1
GROUP BY level;
```

---

## Стратегия keywords

### Automobilizam
- **Типы запчастей:** motor, kočnice, karoserija, menjač, filteri
- **Автомобили:** BMW, Mercedes, Audi, delovi za auto
- **Расходники:** ulje, filter, akumulator, gume
- **Аксессуары:** držač telefona, dash kamera, parking senzori

### Sport i turizam
- **Виды спорта:** fudbal, košarka, tenis, plivanje, bicikl
- **Outdoor:** planinarenje, kampovanje, šator, spavaća vreća
- **Фитнес:** tegovi, bučice, traka za trčanje, oprema
- **Оборудование:** lopta, reket, patike, dres

---

## SEO эффект

Добавленные keywords улучшают:
1. ✅ Внутренний поиск по категориям
2. ✅ SEO метатеги страниц категорий
3. ✅ Релевантность в поисковых системах
4. ✅ Подсказки при поиске товаров
5. ✅ Фильтрацию и навигацию

---

## Следующие шаги

- [ ] Заполнить keywords для остальных L1 категорий:
  - Elektronika
  - Odeca i obuca
  - Kucni aparati
  - Lepota i zdravlje
  - Za bebe i decu
  - Nakit i satovi
  - ...и другие

- [ ] Обновить SEO метатеги на основе keywords
- [ ] Настроить полнотекстовый поиск через FTS индексы (уже созданы)
- [ ] Интегрировать keywords в автодополнение поиска
