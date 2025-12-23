#!/usr/bin/env python3
"""
Создание связей category_attributes для мигрированных атрибутов.

Логика назначения:
- Атрибуты с префиксом book_* → категории книг
- Атрибуты с префиксом pet_* → категории животных
- Атрибуты с префиксом auto_*, car_* → автомобильные категории
- Атрибуты с префиксом job_* → категории работы
- И т.д.
"""

import sys
import psycopg2
from psycopg2.extras import RealDictCursor

class Colors:
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'

LISTINGS_DB = "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db"

# Маппинг префикса атрибута → slug категорий (L1/L2/L3)
ATTRIBUTE_CATEGORY_MAPPING = {
    # Книги
    'book_': ['knjige-i-mediji'],

    # Животные
    'pet_': ['kucni-ljubimci'],

    # Автомобили и автозапчасти
    'auto_': ['automobilizam', 'auto-delovi'],
    'car_': ['automobilizam', 'automobili'],
    'tire_': ['automobilizam', 'gume-i-felne'],
    'rim_': ['automobilizam', 'gume-i-felne'],
    'engine_': ['automobilizam'],
    'vehicle_': ['automobilizam'],
    'motorcycle_': ['automobilizam', 'motocikli'],
    'truck_': ['automobilizam', 'kamioni'],
    'boat_': ['automobilizam', 'plovila'],

    # Работа
    'job_': ['poslovi'],
    'salary_': ['poslovi'],

    # Здоровье и красота
    'health_': ['zdravlje-i-lepota'],

    # Дети
    'kids_': ['deca-i-bebe'],

    # Музыка
    'music_': ['hobiji-i-zabava', 'muzicki-instrumenti'],

    # События
    'event_': ['dogadjaji-i-karte'],

    # Образование
    'edu_': ['usluge', 'obrazovanje'],

    # Искусство
    'art_': ['hobiji-i-zabava', 'umetnost'],

    # Хобби
    'hobby_': ['hobiji-i-zabava'],

    # Недвижимость
    'floor': ['nekretnine'],
    'rooms': ['nekretnine'],
    'furnished': ['nekretnine'],
    'parking': ['nekretnine'],
    'balcony': ['nekretnine'],
    'house_area': ['nekretnine'],
    'land_area': ['nekretnine'],
    'bathrooms': ['nekretnine'],
    'garden': ['nekretnine'],
    'garage': ['nekretnine'],
    'heating_type': ['nekretnine'],
    'construction_year': ['nekretnine'],
    'elevator': ['nekretnine'],
    'security': ['nekretnine'],
    'area': ['nekretnine'],

    # Услуги
    'working_hours': ['usluge'],
    'experience_years': ['usluge'],
    'service_': ['usluge'],
    'availability': ['usluge'],
    'portfolio_url': ['usluge'],

    # Общие автомобильные
    'mileage': ['automobilizam'],
    'year': ['automobilizam', 'elektronika'],
    'power_hp': ['automobilizam'],
    'body_type': ['automobilizam'],
    'drive_type': ['automobilizam'],
    'transmission': ['automobilizam'],
    'fuel_type': ['automobilizam'],
    'doors': ['automobilizam'],
    'seats': ['automobilizam'],
    'vin_number': ['automobilizam'],
    'owner_count': ['automobilizam'],
    'service_book': ['automobilizam'],
    'first_registration': ['automobilizam'],
    'inspection_valid_until': ['automobilizam'],
    'registration_valid_until': ['automobilizam'],
    'equipment_features': ['automobilizam'],
    'additional_equipment': ['automobilizam'],
    'exchange_possible': ['automobilizam'],
    'financing_available': ['automobilizam'],

    # Общие
    'delivery_available': [],  # Применимо ко всем
    'negotiable': [],
    'return_policy': [],
    'warranty': [],
    'warranty_period': [],
    'certification': [],
    'language': [],
    'location': [],
    'price_type': [],
    'country_origin': [],
    'storage': [],
    'storage_type': [],
    'battery_life': ['elektronika'],
    'size': ['moda', 'odeca'],
}

def log_info(msg: str):
    print(f"{Colors.OKBLUE}[INFO]{Colors.ENDC} {msg}")

def log_success(msg: str):
    print(f"{Colors.OKGREEN}[OK]{Colors.ENDC} {msg}")

def log_warning(msg: str):
    print(f"{Colors.WARNING}[WARN]{Colors.ENDC} {msg}")

def log_error(msg: str):
    print(f"{Colors.FAIL}[ERROR]{Colors.ENDC} {msg}", file=sys.stderr)

def find_matching_categories(attr_code: str, all_categories: dict) -> list:
    """Найти подходящие категории для атрибута по его коду"""
    matching_slugs = []

    for prefix, slugs in ATTRIBUTE_CATEGORY_MAPPING.items():
        if attr_code.startswith(prefix) or attr_code == prefix.rstrip('_'):
            matching_slugs.extend(slugs)
            break

    # Находим UUID по slug
    category_ids = []
    for slug in matching_slugs:
        if slug in all_categories:
            category_ids.append(all_categories[slug])

    return category_ids

def main():
    print(f"\n{Colors.BOLD}=== СОЗДАНИЕ СВЯЗЕЙ CATEGORY_ATTRIBUTES ==={Colors.ENDC}\n")

    try:
        log_info("Подключение к listings БД...")
        conn = psycopg2.connect(LISTINGS_DB)
        cur = conn.cursor(cursor_factory=RealDictCursor)
        log_success("Подключено")

        # Загружаем все категории
        log_info("Загрузка категорий...")
        cur.execute("SELECT id, slug FROM categories")
        all_categories = {row['slug']: row['id'] for row in cur.fetchall()}
        log_info(f"Загружено {len(all_categories)} категорий")

        # Находим атрибуты без связей
        log_info("Поиск атрибутов без связей...")
        cur.execute("""
            SELECT a.id, a.code
            FROM attributes a
            LEFT JOIN category_attributes ca ON a.id = ca.attribute_id
            WHERE ca.id IS NULL
            ORDER BY a.code
        """)
        orphan_attrs = cur.fetchall()
        log_info(f"Найдено {len(orphan_attrs)} атрибутов без связей")

        if not orphan_attrs:
            log_success("Все атрибуты уже имеют связи!")
            return

        # Создаём связи
        created = 0
        skipped = 0
        no_category = []

        for attr in orphan_attrs:
            attr_id = attr['id']
            attr_code = attr['code']

            # Пропускаем тестовые атрибуты
            if 'test_' in attr_code:
                skipped += 1
                continue

            category_ids = find_matching_categories(attr_code, all_categories)

            if not category_ids:
                no_category.append(attr_code)
                continue

            for cat_id in category_ids:
                try:
                    cur.execute("""
                        INSERT INTO category_attributes (category_id, attribute_id, is_enabled, is_required, sort_order)
                        VALUES (%s, %s, true, false, 0)
                        ON CONFLICT (category_id, attribute_id) DO NOTHING
                    """, (cat_id, attr_id))
                    created += 1
                except Exception as e:
                    log_warning(f"Не удалось создать связь {attr_code} → {cat_id}: {e}")

        conn.commit()

        # Итоги
        print(f"\n{Colors.BOLD}=== ИТОГИ ==={Colors.ENDC}")
        print(f"Создано связей: {Colors.OKGREEN}{created}{Colors.ENDC}")
        print(f"Пропущено (тестовые): {skipped}")
        print(f"Без категории: {len(no_category)}")

        if no_category:
            print(f"\n{Colors.WARNING}Атрибуты без категории:{Colors.ENDC}")
            for code in no_category[:20]:
                print(f"  - {code}")
            if len(no_category) > 20:
                print(f"  ... и ещё {len(no_category) - 20}")

        log_success("Готово!")

    except Exception as e:
        log_error(f"Ошибка: {e}")
        sys.exit(1)
    finally:
        if 'conn' in dir():
            conn.close()

if __name__ == "__main__":
    main()
