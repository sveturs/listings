#!/usr/bin/env python3
"""
OpenSearch Index Migration Script
Миграция данных из старых индексов (marketplace_listings, storefront_products)
в новые индексы (c2c_listings, b2c_products)
"""

import json
import sys
from typing import Any, Dict, List
import requests
from requests.auth import HTTPBasicAuth

# OpenSearch configuration
OPENSEARCH_HOST = "localhost"
OPENSEARCH_PORT = 9200
OPENSEARCH_URL = f"http://{OPENSEARCH_HOST}:{OPENSEARCH_PORT}"

# Index mapping
INDEX_MIGRATIONS = [
    {
        "source": "marketplace_listings",
        "target": "c2c_listings",
        "description": "C2C Listings (Customer-to-Customer marketplace)"
    },
    {
        "source": "storefront_products",
        "target": "b2c_products",
        "description": "B2C Products (Business-to-Customer stores)"
    },
]


def get_index_mapping(index_name: str) -> Dict[str, Any]:
    """Получить mapping существующего индекса"""
    url = f"{OPENSEARCH_URL}/{index_name}/_mapping"
    response = requests.get(url)

    if response.status_code == 404:
        return {}

    response.raise_for_status()
    return response.json().get(index_name, {}).get("mappings", {})


def get_index_settings(index_name: str) -> Dict[str, Any]:
    """Получить settings существующего индекса"""
    url = f"{OPENSEARCH_URL}/{index_name}/_settings"
    response = requests.get(url)

    if response.status_code == 404:
        return {}

    response.raise_for_status()
    settings = response.json().get(index_name, {}).get("settings", {})

    # Удаляем системные настройки, которые нельзя переносить
    if "index" in settings:
        settings["index"].pop("creation_date", None)
        settings["index"].pop("uuid", None)
        settings["index"].pop("version", None)
        settings["index"].pop("provided_name", None)

    return settings


def create_index(index_name: str, mappings: Dict[str, Any], settings: Dict[str, Any]) -> bool:
    """Создать новый индекс с заданными mappings и settings"""
    url = f"{OPENSEARCH_URL}/{index_name}"

    # Проверяем, существует ли индекс
    check_response = requests.head(url)
    if check_response.status_code == 200:
        print(f"⚠️  Индекс '{index_name}' уже существует")
        user_input = input(f"Удалить существующий индекс '{index_name}'? (yes/no): ")
        if user_input.lower() != 'yes':
            print(f"Пропускаем создание индекса '{index_name}'")
            return False

        # Удаляем существующий индекс
        delete_response = requests.delete(url)
        delete_response.raise_for_status()
        print(f"✅ Индекс '{index_name}' удален")

    # Создаем индекс
    body = {
        "settings": settings,
        "mappings": mappings
    }

    response = requests.put(url, json=body, headers={"Content-Type": "application/json"})
    response.raise_for_status()
    print(f"✅ Индекс '{index_name}' создан успешно")
    return True


def reindex_data(source_index: str, target_index: str) -> Dict[str, Any]:
    """Переиндексировать данные из source в target"""
    url = f"{OPENSEARCH_URL}/_reindex"

    body = {
        "source": {
            "index": source_index
        },
        "dest": {
            "index": target_index
        }
    }

    response = requests.post(url, json=body, headers={"Content-Type": "application/json"})
    response.raise_for_status()
    result = response.json()

    print(f"✅ Переиндексация завершена:")
    print(f"   • Всего документов: {result.get('total', 0)}")
    print(f"   • Создано: {result.get('created', 0)}")
    print(f"   • Обновлено: {result.get('updated', 0)}")
    print(f"   • Ошибок: {len(result.get('failures', []))}")

    if result.get('failures'):
        print(f"⚠️  Ошибки переиндексации:")
        for failure in result['failures']:
            print(f"   • {failure}")

    return result


def get_document_count(index_name: str) -> int:
    """Получить количество документов в индексе"""
    url = f"{OPENSEARCH_URL}/{index_name}/_count"
    response = requests.get(url)

    if response.status_code == 404:
        return 0

    response.raise_for_status()
    return response.json().get("count", 0)


def verify_migration(source_index: str, target_index: str) -> bool:
    """Проверить, что миграция прошла успешно"""
    source_count = get_document_count(source_index)
    target_count = get_document_count(target_index)

    print(f"\n📊 Проверка миграции:")
    print(f"   • Исходный индекс '{source_index}': {source_count} документов")
    print(f"   • Целевой индекс '{target_index}': {target_count} документов")

    if source_count == target_count:
        print(f"✅ Миграция успешна! Все документы перенесены.")
        return True
    else:
        print(f"⚠️  Внимание! Количество документов не совпадает.")
        return False


def main():
    """Главная функция миграции"""
    print("=" * 80)
    print("🔄 Миграция OpenSearch индексов: marketplace/storefronts → c2c/b2c")
    print("=" * 80)

    # Проверяем доступность OpenSearch
    try:
        response = requests.get(OPENSEARCH_URL)
        response.raise_for_status()
        print(f"✅ OpenSearch доступен на {OPENSEARCH_URL}")
    except Exception as e:
        print(f"❌ Ошибка подключения к OpenSearch: {e}")
        sys.exit(1)

    # Выполняем миграцию для каждой пары индексов
    for migration in INDEX_MIGRATIONS:
        source = migration["source"]
        target = migration["target"]
        description = migration["description"]

        print(f"\n{'=' * 80}")
        print(f"📦 Миграция: {source} → {target}")
        print(f"   Описание: {description}")
        print(f"{'=' * 80}")

        # Проверяем существование исходного индекса
        source_count = get_document_count(source)
        if source_count == 0:
            print(f"⚠️  Исходный индекс '{source}' пуст или не существует. Пропускаем.")

            # Все равно создаем целевой индекс с правильной структурой
            print(f"📝 Создаем пустой индекс '{target}' с правильной структурой...")

            # Получаем mapping и settings из исходного индекса (если существует)
            mappings = get_index_mapping(source)
            settings = get_index_settings(source)

            if mappings or settings:
                create_index(target, mappings, settings)
            else:
                print(f"⚠️  Не удалось получить структуру индекса '{source}'. Пропускаем.")

            continue

        print(f"📊 Исходный индекс содержит {source_count} документов")

        # Получаем mapping и settings
        print(f"📝 Получаем структуру индекса '{source}'...")
        mappings = get_index_mapping(source)
        settings = get_index_settings(source)

        # Создаем новый индекс
        print(f"🔨 Создаем новый индекс '{target}'...")
        if not create_index(target, mappings, settings):
            print(f"⏭️  Пропускаем переиндексацию для '{target}'")
            continue

        # Переиндексируем данные
        print(f"🔄 Переиндексация данных из '{source}' в '{target}'...")
        reindex_data(source, target)

        # Проверяем результат
        verify_migration(source, target)

    print(f"\n{'=' * 80}")
    print("✅ Миграция OpenSearch индексов завершена!")
    print("=" * 80)

    # Финальная сводка
    print("\n📊 Итоговая статистика:")
    for migration in INDEX_MIGRATIONS:
        source = migration["source"]
        target = migration["target"]
        source_count = get_document_count(source)
        target_count = get_document_count(target)
        status = "✅" if source_count == target_count else "⚠️"
        print(f"   {status} {source} ({source_count}) → {target} ({target_count})")

    print("\n💡 Следующие шаги:")
    print("   1. Обновите код backend для использования новых индексов")
    print("   2. Протестируйте поиск в новых индексах")
    print("   3. После проверки можете удалить старые индексы:")
    for migration in INDEX_MIGRATIONS:
        print(f"      curl -X DELETE {OPENSEARCH_URL}/{migration['source']}")


if __name__ == "__main__":
    main()
