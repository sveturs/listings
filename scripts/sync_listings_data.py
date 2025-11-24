#!/usr/bin/env python3
"""
Синхронизация справочных данных между монолитом (svetubd) и микросервисом listings (listings_dev_db).

Основная сложность: трансформация VARCHAR полей в JSONB для многоязычности.
"""

import sys
import json
from typing import Dict, List, Any, Optional
from dataclasses import dataclass
import psycopg2
from psycopg2.extensions import ISOLATION_LEVEL_AUTOCOMMIT
from psycopg2.extras import RealDictCursor

# ANSI цвета для красивого вывода
class Colors:
    HEADER = '\033[95m'
    OKBLUE = '\033[94m'
    OKCYAN = '\033[96m'
    OKGREEN = '\033[92m'
    WARNING = '\033[93m'
    FAIL = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

@dataclass
class SyncStats:
    """Статистика синхронизации"""
    categories_before: int = 0
    categories_after: int = 0
    attributes_before: int = 0
    attributes_after: int = 0
    category_attributes_before: int = 0
    category_attributes_after: int = 0

    def print_summary(self):
        """Печать итоговой статистики"""
        print(f"\n{Colors.HEADER}{Colors.BOLD}=== ИТОГОВАЯ СТАТИСТИКА ==={Colors.ENDC}")

        print(f"\n{Colors.OKBLUE}Categories:{Colors.ENDC}")
        print(f"  До:    {self.categories_before}")
        print(f"  После: {Colors.OKGREEN}{self.categories_after}{Colors.ENDC}")
        print(f"  Добавлено: {Colors.OKCYAN}{self.categories_after - self.categories_before}{Colors.ENDC}")

        print(f"\n{Colors.OKBLUE}Attributes:{Colors.ENDC}")
        print(f"  До:    {self.attributes_before}")
        print(f"  После: {Colors.OKGREEN}{self.attributes_after}{Colors.ENDC}")
        print(f"  Добавлено: {Colors.OKCYAN}{self.attributes_after - self.attributes_before}{Colors.ENDC}")

        print(f"\n{Colors.OKBLUE}Category Attributes:{Colors.ENDC}")
        print(f"  До:    {self.category_attributes_before}")
        print(f"  После: {Colors.OKGREEN}{self.category_attributes_after}{Colors.ENDC}")
        print(f"  Добавлено: {Colors.OKCYAN}{self.category_attributes_after - self.category_attributes_before}{Colors.ENDC}")

class DatabaseSync:
    """Класс для синхронизации баз данных"""

    MONOLITH_DB = "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd"
    MICROSERVICE_DB = "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db"

    def __init__(self):
        self.monolith_conn = None
        self.microservice_conn = None
        self.stats = SyncStats()

    def log_info(self, message: str):
        """Логирование информационных сообщений"""
        print(f"{Colors.OKBLUE}[INFO]{Colors.ENDC} {message}")

    def log_success(self, message: str):
        """Логирование успешных операций"""
        print(f"{Colors.OKGREEN}[SUCCESS]{Colors.ENDC} {message}")

    def log_warning(self, message: str):
        """Логирование предупреждений"""
        print(f"{Colors.WARNING}[WARNING]{Colors.ENDC} {message}")

    def log_error(self, message: str):
        """Логирование ошибок"""
        print(f"{Colors.FAIL}[ERROR]{Colors.ENDC} {message}", file=sys.stderr)

    def connect_databases(self):
        """Подключение к обеим базам данных"""
        self.log_info("Подключение к базам данных...")

        try:
            # Подключение к монолиту
            self.log_info(f"Подключение к монолиту (порт 5433)...")
            self.monolith_conn = psycopg2.connect(self.MONOLITH_DB)
            self.monolith_conn.set_isolation_level(ISOLATION_LEVEL_AUTOCOMMIT)
            self.log_success("Подключено к монолиту")

            # Подключение к микросервису
            self.log_info(f"Подключение к микросервису listings (порт 35434)...")
            self.microservice_conn = psycopg2.connect(self.MICROSERVICE_DB)
            self.log_success("Подключено к микросервису")

        except Exception as e:
            self.log_error(f"Ошибка подключения к БД: {e}")
            raise

    def get_table_count(self, conn, table: str) -> int:
        """Получить количество записей в таблице"""
        with conn.cursor() as cur:
            # Все таблицы в схеме public
            cur.execute(f"SELECT COUNT(*) FROM {table}")
            return cur.fetchone()[0]

    def collect_stats_before(self):
        """Собрать статистику перед синхронизацией"""
        self.log_info("Сбор статистики ДО синхронизации...")

        self.stats.categories_before = self.get_table_count(self.microservice_conn, "categories")
        self.stats.attributes_before = self.get_table_count(self.microservice_conn, "attributes")
        self.stats.category_attributes_before = self.get_table_count(self.microservice_conn, "category_attributes")

        self.log_info(f"Categories: {self.stats.categories_before}")
        self.log_info(f"Attributes: {self.stats.attributes_before}")
        self.log_info(f"Category Attributes: {self.stats.category_attributes_before}")

    def collect_stats_after(self):
        """Собрать статистику после синхронизации"""
        self.log_info("Сбор статистики ПОСЛЕ синхронизации...")

        self.stats.categories_after = self.get_table_count(self.microservice_conn, "categories")
        self.stats.attributes_after = self.get_table_count(self.microservice_conn, "attributes")
        self.stats.category_attributes_after = self.get_table_count(self.microservice_conn, "category_attributes")

        self.log_info(f"Categories: {self.stats.categories_after}")
        self.log_info(f"Attributes: {self.stats.attributes_after}")
        self.log_info(f"Category Attributes: {self.stats.category_attributes_after}")

    def varchar_to_jsonb(self, value: str) -> str:
        """Преобразовать VARCHAR в JSONB для многоязычности"""
        if not value:
            return json.dumps({"en": "", "sr": "", "ru": ""})

        # Создаем JSON объект с одинаковым значением для всех языков
        multilang = {
            "en": value,
            "sr": value,
            "ru": value
        }
        return json.dumps(multilang)

    def clear_target_tables(self):
        """Очистить таблицы микросервиса в правильном порядке (учет FK)"""
        self.log_info("Очистка таблиц микросервиса (с учетом FK constraints)...")

        cur = self.microservice_conn.cursor()

        try:
            # Порядок важен! Сначала зависимые таблицы, потом родительские
            tables = [
                "category_attributes",  # Зависит от categories и attributes
                "attributes",           # Независимая
                "categories"            # Независимая
            ]

            for table in tables:
                self.log_info(f"Очистка {table}...")
                cur.execute(f"TRUNCATE TABLE {table} RESTART IDENTITY CASCADE")
                self.log_success(f"Очищена таблица {table}")

            self.microservice_conn.commit()
            self.log_success("Все таблицы очищены")

        except Exception as e:
            self.log_error(f"Ошибка очистки таблиц: {e}")
            self.microservice_conn.rollback()
            raise
        finally:
            cur.close()

    def sync_categories(self):
        """Синхронизация категорий (прямое копирование)"""
        self.log_info("Синхронизация categories (c2c_categories → categories)...")

        # Читаем из монолита
        monolith_cur = self.monolith_conn.cursor(cursor_factory=RealDictCursor)
        monolith_cur.execute("""
            SELECT id, name, slug, parent_id, level, is_active, icon, created_at,
                   has_custom_ui, custom_ui_component, sort_order, count,
                   external_id, description, seo_title, seo_description, seo_keywords
            FROM c2c_categories
            ORDER BY id
        """)
        categories = monolith_cur.fetchall()
        monolith_cur.close()

        self.log_info(f"Найдено {len(categories)} категорий в монолите")

        # Пишем в микросервис
        micro_cur = self.microservice_conn.cursor()

        try:
            for cat in categories:
                micro_cur.execute("""
                    INSERT INTO categories
                    (id, name, slug, parent_id, level, is_active, icon, created_at,
                     has_custom_ui, custom_ui_component, sort_order, count,
                     external_id, description, seo_title, seo_description, seo_keywords)
                    VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                """, (
                    cat['id'],
                    cat['name'],
                    cat['slug'],
                    cat['parent_id'],
                    cat['level'],
                    cat['is_active'],
                    cat['icon'],
                    cat['created_at'],
                    cat['has_custom_ui'],
                    cat['custom_ui_component'],
                    cat['sort_order'],
                    cat['count'],
                    cat['external_id'],
                    cat['description'],
                    cat['seo_title'],
                    cat['seo_description'],
                    cat['seo_keywords']
                ))

            self.microservice_conn.commit()
            self.log_success(f"Синхронизировано {len(categories)} категорий")

        except Exception as e:
            self.log_error(f"Ошибка синхронизации категорий: {e}")
            self.microservice_conn.rollback()
            raise
        finally:
            micro_cur.close()

    def sync_attributes(self):
        """Синхронизация атрибутов с трансформацией VARCHAR → JSONB"""
        self.log_info("Синхронизация attributes (unified_attributes → attributes) с трансформацией VARCHAR → JSONB...")

        # Читаем из монолита
        monolith_cur = self.monolith_conn.cursor(cursor_factory=RealDictCursor)
        monolith_cur.execute("""
            SELECT id, code, name, display_name, attribute_type, purpose,
                   options, validation_rules, ui_settings,
                   is_searchable, is_filterable, is_required,
                   is_variant_compatible, affects_stock, affects_price,
                   show_in_card, is_active, sort_order,
                   legacy_category_attribute_id, legacy_product_variant_attribute_id,
                   icon, search_vector, created_at, updated_at
            FROM unified_attributes
            ORDER BY id
        """)
        attributes = monolith_cur.fetchall()
        monolith_cur.close()

        self.log_info(f"Найдено {len(attributes)} атрибутов в монолите")

        # Пишем в микросервис с трансформацией
        micro_cur = self.microservice_conn.cursor()

        try:
            for attr in attributes:
                # Трансформация VARCHAR → JSONB для name и display_name
                name_jsonb = self.varchar_to_jsonb(attr['name'])
                display_name_jsonb = self.varchar_to_jsonb(attr['display_name'])

                # Сериализация JSONB полей из dict → string
                options_json = json.dumps(attr['options']) if attr['options'] else '{}'
                validation_rules_json = json.dumps(attr['validation_rules']) if attr['validation_rules'] else '{}'
                ui_settings_json = json.dumps(attr['ui_settings']) if attr['ui_settings'] else '{}'

                micro_cur.execute("""
                    INSERT INTO attributes
                    (id, code, name, display_name, attribute_type, purpose,
                     options, validation_rules, ui_settings,
                     is_searchable, is_filterable, is_required,
                     is_variant_compatible, affects_stock, affects_price,
                     show_in_card, is_active, sort_order,
                     legacy_category_attribute_id, legacy_product_variant_attribute_id,
                     icon, search_vector, created_at, updated_at)
                    VALUES (%s, %s, %s::jsonb, %s::jsonb, %s, %s, %s::jsonb, %s::jsonb, %s::jsonb,
                            %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
                """, (
                    attr['id'],
                    attr['code'],
                    name_jsonb,
                    display_name_jsonb,
                    attr['attribute_type'],
                    attr['purpose'],
                    options_json,
                    validation_rules_json,
                    ui_settings_json,
                    attr['is_searchable'],
                    attr['is_filterable'],
                    attr['is_required'],
                    attr['is_variant_compatible'],
                    attr['affects_stock'],
                    attr['affects_price'],
                    attr['show_in_card'],
                    attr['is_active'],
                    attr['sort_order'],
                    attr['legacy_category_attribute_id'],
                    attr['legacy_product_variant_attribute_id'],
                    attr['icon'],
                    attr['search_vector'],
                    attr['created_at'],
                    attr['updated_at']
                ))

            self.microservice_conn.commit()
            self.log_success(f"Синхронизировано {len(attributes)} атрибутов (с трансформацией VARCHAR → JSONB)")

        except Exception as e:
            self.log_error(f"Ошибка синхронизации атрибутов: {e}")
            self.microservice_conn.rollback()
            raise
        finally:
            micro_cur.close()

    def sync_category_attributes(self):
        """Синхронизация связей категорий и атрибутов (прямое копирование)"""
        self.log_info("Синхронизация category_attributes (unified_category_attributes → category_attributes)...")

        # Читаем из монолита
        monolith_cur = self.monolith_conn.cursor(cursor_factory=RealDictCursor)
        monolith_cur.execute("""
            SELECT category_id, attribute_id, is_required, sort_order, created_at
            FROM unified_category_attributes
            ORDER BY category_id, attribute_id
        """)
        cat_attrs = monolith_cur.fetchall()
        monolith_cur.close()

        self.log_info(f"Найдено {len(cat_attrs)} связей категория-атрибут в монолите")

        # Пишем в микросервис
        micro_cur = self.microservice_conn.cursor()

        try:
            for ca in cat_attrs:
                micro_cur.execute("""
                    INSERT INTO category_attributes
                    (category_id, attribute_id, is_required, sort_order, created_at)
                    VALUES (%s, %s, %s, %s, %s)
                """, (
                    ca['category_id'],
                    ca['attribute_id'],
                    ca['is_required'],
                    ca['sort_order'],
                    ca['created_at']
                ))

            self.microservice_conn.commit()
            self.log_success(f"Синхронизировано {len(cat_attrs)} связей категория-атрибут")

        except Exception as e:
            self.log_error(f"Ошибка синхронизации связей: {e}")
            self.microservice_conn.rollback()
            raise
        finally:
            micro_cur.close()

    def reset_sequences(self):
        """Сброс sequences для auto-increment полей"""
        self.log_info("Сброс sequences для auto-increment...")

        cur = self.microservice_conn.cursor()

        try:
            sequences = [
                ("categories", "id"),
                ("attributes", "id")
            ]

            for table, id_col in sequences:
                # Получаем максимальный ID
                cur.execute(f"SELECT MAX({id_col}) FROM {table}")
                max_id = cur.fetchone()[0]

                if max_id:
                    # Устанавливаем sequence на max_id + 1
                    cur.execute(f"SELECT setval('{table}_{id_col}_seq', %s, true)", (max_id,))
                    self.log_success(f"Sequence {table}_{id_col}_seq установлен на {max_id}")

            self.microservice_conn.commit()

        except Exception as e:
            self.log_error(f"Ошибка сброса sequences: {e}")
            self.microservice_conn.rollback()
            raise
        finally:
            cur.close()

    def verify_integrity(self):
        """Проверка целостности данных после синхронизации"""
        self.log_info("Проверка целостности данных...")

        cur = self.microservice_conn.cursor(cursor_factory=RealDictCursor)

        try:
            # 1. Проверка FK constraints в category_attributes
            cur.execute("""
                SELECT COUNT(*) as invalid_count
                FROM category_attributes ca
                LEFT JOIN categories c ON ca.category_id = c.id
                WHERE c.id IS NULL
            """)
            invalid_cats = cur.fetchone()['invalid_count']

            if invalid_cats > 0:
                self.log_error(f"Найдено {invalid_cats} связей с несуществующими категориями!")
            else:
                self.log_success("Все связи категорий валидны")

            cur.execute("""
                SELECT COUNT(*) as invalid_count
                FROM category_attributes ca
                LEFT JOIN attributes a ON ca.attribute_id = a.id
                WHERE a.id IS NULL
            """)
            invalid_attrs = cur.fetchone()['invalid_count']

            if invalid_attrs > 0:
                self.log_error(f"Найдено {invalid_attrs} связей с несуществующими атрибутами!")
            else:
                self.log_success("Все связи атрибутов валидны")

            # 2. Проверка наличия критичной категории 1008
            cur.execute("SELECT id, name FROM categories WHERE id = 1008")
            cat_1008 = cur.fetchone()

            if cat_1008:
                self.log_success(f"Категория 1008 найдена: {cat_1008['name']}")
            else:
                self.log_error("Категория 1008 НЕ НАЙДЕНА!")

            # 3. Проверка трансформации JSONB (на примере первых 3 атрибутов)
            cur.execute("""
                SELECT id, name, display_name
                FROM attributes
                ORDER BY id
                LIMIT 3
            """)
            sample_attrs = cur.fetchall()

            self.log_info("Проверка трансформации VARCHAR → JSONB (примеры):")
            for attr in sample_attrs:
                name_obj = attr['name']
                display_obj = attr['display_name']

                # Проверяем что это валидный JSON с нужными ключами
                if isinstance(name_obj, dict) and all(k in name_obj for k in ['en', 'sr', 'ru']):
                    self.log_success(f"  Атрибут {attr['id']}: name = {name_obj}")
                else:
                    self.log_error(f"  Атрибут {attr['id']}: name НЕ JSONB или некорректный формат!")

                if isinstance(display_obj, dict) and all(k in display_obj for k in ['en', 'sr', 'ru']):
                    self.log_success(f"  Атрибут {attr['id']}: display_name = {display_obj}")
                else:
                    self.log_error(f"  Атрибут {attr['id']}: display_name НЕ JSONB или некорректный формат!")

        except Exception as e:
            self.log_error(f"Ошибка проверки целостности: {e}")
            raise
        finally:
            cur.close()

    def run_sync(self):
        """Основной процесс синхронизации"""
        try:
            print(f"\n{Colors.HEADER}{Colors.BOLD}=== СИНХРОНИЗАЦИЯ СПРАВОЧНЫХ ДАННЫХ ==={Colors.ENDC}")
            print(f"{Colors.OKCYAN}Монолит БД (svetubd, порт 5433) → Микросервис БД (listings_dev_db, порт 35434){Colors.ENDC}\n")

            # 1. Подключение
            self.connect_databases()

            # 2. Статистика ДО
            self.collect_stats_before()

            # 3. Очистка таблиц
            self.clear_target_tables()

            # 4. Синхронизация в правильном порядке
            self.sync_categories()
            self.sync_attributes()
            self.sync_category_attributes()

            # 5. Сброс sequences
            self.reset_sequences()

            # 6. Проверка целостности
            self.verify_integrity()

            # 7. Статистика ПОСЛЕ
            self.collect_stats_after()

            # 8. Итоговая статистика
            self.stats.print_summary()

            print(f"\n{Colors.OKGREEN}{Colors.BOLD}=== СИНХРОНИЗАЦИЯ ЗАВЕРШЕНА УСПЕШНО ==={Colors.ENDC}\n")

        except Exception as e:
            self.log_error(f"Критическая ошибка синхронизации: {e}")
            print(f"\n{Colors.FAIL}{Colors.BOLD}=== СИНХРОНИЗАЦИЯ ПРЕРВАНА ==={Colors.ENDC}\n")
            sys.exit(1)

        finally:
            # Закрытие подключений
            if self.monolith_conn:
                self.monolith_conn.close()
                self.log_info("Подключение к монолиту закрыто")

            if self.microservice_conn:
                self.microservice_conn.close()
                self.log_info("Подключение к микросервису закрыто")

def main():
    """Точка входа"""
    sync = DatabaseSync()
    sync.run_sync()

if __name__ == "__main__":
    main()
