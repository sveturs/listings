#!/usr/bin/env python3
"""
Скрипт для генерации тестовых CSV файлов с товарами для тестирования импорта.
Использование:
    python3 generate_test_products.py --count 1000 --output test_products_1000.csv
"""

import argparse
import csv
import random
from typing import List

# Списки для генерации случайных товаров
CATEGORIES = [1000, 1001, 1002, 1003, 1004]
BRANDS = ["Samsung", "Apple", "Sony", "LG", "Huawei", "Xiaomi", "OnePlus"]
MODELS = ["Pro", "Max", "Ultra", "Lite", "Plus", "SE", "Mini"]
COLORS = ["Black", "White", "Blue", "Red", "Green", "Gold", "Silver"]
CURRENCIES = ["USD", "EUR", "RSD"]

def generate_product(index: int) -> dict:
    """Генерирует случайный товар"""
    brand = random.choice(BRANDS)
    model = random.choice(MODELS)
    color = random.choice(COLORS)

    return {
        "name": f"{brand} {model} {color}",
        "description": f"High-quality product from {brand}. {model} series with amazing features.",
        "price": round(random.uniform(100, 2000), 2),
        "currency": random.choice(CURRENCIES),
        "category_id": random.choice(CATEGORIES),
        "sku": f"SKU-{index:06d}",
        "barcode": f"{random.randint(1000000000, 9999999999)}",
        "stock_quantity": random.randint(0, 100),
        "is_active": random.choice([True, False]),
    }

def generate_csv(count: int, output_file: str):
    """Генерирует CSV файл с заданным количеством товаров"""
    print(f"Генерация {count} товаров в файл {output_file}...")

    fieldnames = [
        "name",
        "description",
        "price",
        "currency",
        "category_id",
        "sku",
        "barcode",
        "stock_quantity",
        "is_active",
    ]

    with open(output_file, 'w', newline='', encoding='utf-8') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()

        for i in range(1, count + 1):
            product = generate_product(i)
            writer.writerow(product)

            if i % 100 == 0:
                print(f"Сгенерировано {i}/{count} товаров...")

    print(f"✅ Готово! Файл {output_file} создан с {count} товарами.")

def main():
    parser = argparse.ArgumentParser(
        description="Генератор тестовых CSV файлов с товарами"
    )
    parser.add_argument(
        "--count",
        type=int,
        default=1000,
        help="Количество товаров для генерации (по умолчанию: 1000)"
    )
    parser.add_argument(
        "--output",
        type=str,
        default="test_products.csv",
        help="Имя выходного файла (по умолчанию: test_products.csv)"
    )

    args = parser.parse_args()

    if args.count <= 0:
        print("❌ Ошибка: количество товаров должно быть больше 0")
        return

    generate_csv(args.count, args.output)

if __name__ == "__main__":
    main()
