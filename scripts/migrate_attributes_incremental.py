#!/usr/bin/env python3
"""
Инкрементальная миграция атрибутов из монолита (unified_attributes) в listings (attributes).

Особенности:
- Пропускает атрибуты, которые уже есть в listings по code
- Добавляет только новые с авто-генерацией ID
- Трансформирует VARCHAR → JSONB для name/display_name
- Сохраняет маппинг old_id → new_id для последующей миграции связей
"""

import sys
import json
from typing import Dict, Set
import psycopg2
from psycopg2.extras import RealDictCursor

# ANSI colors
class Colors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'

MONOLITH_DB = "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd"
MICROSERVICE_DB = "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db"

def log_info(msg: str):
    print(f"{Colors.OKBLUE}[INFO]{Colors.ENDC} {msg}")

def log_success(msg: str):
    print(f"{Colors.OKGREEN}[OK]{Colors.ENDC} {msg}")

def log_warning(msg: str):
    print(f"{Colors.WARNING}[WARN]{Colors.ENDC} {msg}")

def log_error(msg: str):
    print(f"{Colors.FAIL}[ERROR]{Colors.ENDC} {msg}", file=sys.stderr)

def varchar_to_jsonb(value: str) -> str:
    """Трансформация VARCHAR в JSONB для многоязычности"""
    if not value:
        return json.dumps({"en": "", "sr": "", "ru": ""})
    return json.dumps({"en": value, "sr": value, "ru": value})

def main():
    print(f"\n{Colors.HEADER}{Colors.BOLD}=== ИНКРЕМЕНТАЛЬНАЯ МИГРАЦИЯ АТРИБУТОВ ==={Colors.ENDC}\n")

    try:
        # Подключение к базам
        log_info("Подключение к монолиту (порт 5433)...")
        monolith_conn = psycopg2.connect(MONOLITH_DB)
        log_success("Подключено к монолиту")

        log_info("Подключение к listings (порт 35434)...")
        micro_conn = psycopg2.connect(MICROSERVICE_DB)
        log_success("Подключено к listings")

        # Получаем существующие коды в listings
        log_info("Чтение существующих атрибутов в listings...")
        micro_cur = micro_conn.cursor(cursor_factory=RealDictCursor)
        micro_cur.execute("SELECT id, code FROM attributes")
        existing_attrs = {row['code']: row['id'] for row in micro_cur.fetchall()}
        log_info(f"Найдено {len(existing_attrs)} атрибутов в listings")

        # Получаем все атрибуты из монолита
        log_info("Чтение атрибутов из монолита...")
        mono_cur = monolith_conn.cursor(cursor_factory=RealDictCursor)
        mono_cur.execute("""
            SELECT id, code, name, display_name, attribute_type, purpose,
                   options, validation_rules, ui_settings,
                   is_searchable, is_filterable, is_required,
                   is_variant_compatible, affects_stock, affects_price,
                   show_in_card, is_active, sort_order,
                   legacy_category_attribute_id, icon,
                   created_at, updated_at
            FROM unified_attributes
            ORDER BY id
        """)
        monolith_attrs = mono_cur.fetchall()
        log_info(f"Найдено {len(monolith_attrs)} атрибутов в монолите")

        # Определяем что нужно мигрировать
        to_migrate = []
        skipped_codes = []

        for attr in monolith_attrs:
            if attr['code'] in existing_attrs:
                skipped_codes.append(attr['code'])
            else:
                to_migrate.append(attr)

        log_info(f"Пропущено (уже есть): {len(skipped_codes)}")
        log_info(f"К миграции: {len(to_migrate)}")

        if not to_migrate:
            log_warning("Нет атрибутов для миграции. Все уже существуют в listings.")
            return

        # Миграция
        log_info(f"\nМиграция {len(to_migrate)} атрибутов...")

        id_mapping: Dict[int, int] = {}  # old_id → new_id
        migrated = 0

        for attr in to_migrate:
            try:
                # Трансформация VARCHAR → JSONB
                name_jsonb = varchar_to_jsonb(attr['name'])
                display_name_jsonb = varchar_to_jsonb(attr['display_name'])

                # Сериализация JSONB полей
                options_json = json.dumps(attr['options']) if attr['options'] else '{}'
                validation_rules_json = json.dumps(attr['validation_rules']) if attr['validation_rules'] else '{}'
                ui_settings_json = json.dumps(attr['ui_settings']) if attr['ui_settings'] else '{}'

                # INSERT с авто-генерацией ID
                micro_cur.execute("""
                    INSERT INTO attributes
                    (code, name, display_name, attribute_type, purpose,
                     options, validation_rules, ui_settings,
                     is_searchable, is_filterable, is_required,
                     is_variant_compatible, affects_stock, affects_price,
                     show_in_card, is_active, sort_order,
                     legacy_category_attribute_id, icon,
                     created_at, updated_at)
                    VALUES (%s, %s::jsonb, %s::jsonb, %s, %s, %s::jsonb, %s::jsonb, %s::jsonb,
                            %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                    RETURNING id
                """, (
                    attr['code'],
                    name_jsonb,
                    display_name_jsonb,
                    attr['attribute_type'],
                    attr['purpose'],
                    options_json,
                    validation_rules_json,
                    ui_settings_json,
                    attr['is_searchable'] or False,
                    attr['is_filterable'] or False,
                    attr['is_required'] or False,
                    attr['is_variant_compatible'] or False,
                    attr['affects_stock'] or False,
                    attr['affects_price'] or False,
                    attr['show_in_card'] or False,
                    attr['is_active'] if attr['is_active'] is not None else True,
                    attr['sort_order'] or 0,
                    attr['legacy_category_attribute_id'],
                    attr['icon'] or '',
                    attr['created_at'],
                    attr['updated_at']
                ))

                result = micro_cur.fetchone()
                new_id = result['id'] if isinstance(result, dict) else result[0]
                id_mapping[attr['id']] = new_id
                migrated += 1

                if migrated % 50 == 0:
                    log_info(f"  Мигрировано {migrated}/{len(to_migrate)}...")

            except Exception as e:
                import traceback
                log_error(f"Ошибка миграции атрибута {attr['code']}: {e}")
                log_error(f"Traceback: {traceback.format_exc()}")
                log_error(f"Данные атрибута: id={attr['id']}, type={attr['attribute_type']}, purpose={attr['purpose']}")
                micro_conn.rollback()
                raise

        micro_conn.commit()
        log_success(f"Мигрировано {migrated} атрибутов")

        # Сохраняем маппинг в файл
        mapping_file = "/tmp/attribute_id_mapping.json"
        with open(mapping_file, 'w') as f:
            json.dump(id_mapping, f, indent=2)
        log_success(f"Маппинг ID сохранён в {mapping_file}")

        # Обновляем маппинг для существующих (code → new_id)
        code_to_new_id = {}
        micro_cur.execute("SELECT id, code FROM attributes")
        for row in micro_cur.fetchall():
            code_to_new_id[row['code']] = row['id']

        # Добавляем пропущенные в маппинг (старый id → новый id по code)
        for attr in monolith_attrs:
            if attr['code'] in existing_attrs and attr['id'] not in id_mapping:
                id_mapping[attr['id']] = code_to_new_id[attr['code']]

        # Сохраняем полный маппинг
        full_mapping_file = "/tmp/attribute_id_mapping_full.json"
        with open(full_mapping_file, 'w') as f:
            json.dump(id_mapping, f, indent=2)
        log_success(f"Полный маппинг (включая существующие) сохранён в {full_mapping_file}")

        # Итоговая статистика
        print(f"\n{Colors.HEADER}{Colors.BOLD}=== ИТОГИ ==={Colors.ENDC}")

        # Проверяем итоговое количество
        micro_cur.execute("SELECT COUNT(*) FROM attributes")
        total = micro_cur.fetchone()[0]

        print(f"\n{Colors.OKBLUE}Атрибуты в listings:{Colors.ENDC}")
        print(f"  Было:    30")
        print(f"  Добавлено: {Colors.OKGREEN}{migrated}{Colors.ENDC}")
        print(f"  Итого: {Colors.OKGREEN}{total}{Colors.ENDC}")

        print(f"\n{Colors.OKBLUE}Пропущенные (пересечение):{Colors.ENDC}")
        for code in sorted(skipped_codes):
            print(f"  - {code}")

        print(f"\n{Colors.OKGREEN}Миграция завершена успешно!{Colors.ENDC}\n")

    except Exception as e:
        log_error(f"Критическая ошибка: {e}")
        sys.exit(1)
    finally:
        if 'monolith_conn' in dir():
            monolith_conn.close()
        if 'micro_conn' in dir():
            micro_conn.close()

if __name__ == "__main__":
    main()
