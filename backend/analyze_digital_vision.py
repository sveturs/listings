#!/usr/bin/env python3
"""
Скрипт анализа Digital Vision прайса

Анализирует XML прайс Digital Vision и выдает:
- Статистику категорий
- Статистику атрибутов
- Потенциальные группы вариантов
- Статистику изображений

Usage:
    python3 analyze_digital_vision.py --file DigitalVision.xml --output analysis.json
"""

import argparse
import json
import re
import xml.etree.ElementTree as ET
from collections import defaultdict, Counter
from dataclasses import dataclass, field, asdict
from typing import List, Dict, Set, Optional
from pathlib import Path


@dataclass
class CategoryStats:
    """Статистика по категориям"""
    total: int = 0
    level1: Set[str] = field(default_factory=set)
    level2: Set[str] = field(default_factory=set)
    level3: Set[str] = field(default_factory=set)
    category_product_count: Dict[str, int] = field(default_factory=lambda: defaultdict(int))
    top_categories: List[Dict[str, any]] = field(default_factory=list)


@dataclass
class AttributeStats:
    """Статистика по атрибутам"""
    detected: List[str] = field(default_factory=list)
    unique_values: Dict[str, Set[str]] = field(default_factory=lambda: defaultdict(set))
    value_counts: Dict[str, Dict[str, int]] = field(default_factory=lambda: defaultdict(lambda: defaultdict(int)))


@dataclass
class VariantGroup:
    """Группа потенциальных вариантов"""
    base_name: str
    products: List[Dict[str, str]]
    variant_count: int
    variant_attributes: Set[str]
    confidence: float


@dataclass
class ImageStats:
    """Статистика изображений"""
    total_products_with_images: int = 0
    total_images: int = 0
    avg_images_per_product: float = 0.0
    max_images_per_product: int = 0
    percentage: float = 0.0


@dataclass
class AnalysisResult:
    """Результат полного анализа"""
    total_products: int
    categories: Dict
    attributes: Dict
    variants: Dict
    images: Dict


class DigitalVisionAnalyzer:
    """Анализатор Digital Vision прайса"""

    # Паттерны для определения вариантов
    COLOR_PATTERNS = [
        # Основные цвета (английский, сербский)
        r'\b(crn[ai]|bel[ai]|crven[ai]?|zelen[ai]?|plav[ai]|pink|black|white|red|blue|green|yellow|grey|gray|silver|gold)\b',
        # Модификаторы цвета + цвет (dark blue, light yellow)
        r'\b(dark|light|bright|deep)\s+(blue|red|green|yellow|black|white|grey|gray)\b',
        # Одиночные модификаторы цвета
        r'\b(light|dark|bright|deep)\b',
    ]

    SIZE_PATTERNS = [
        r'\b\d+\/\d+\/\d+\/\d+\s*mm\b',  # 42/44/45/49mm - четыре размера
        r'\b\d+\/\d+\/\d+\s*mm\b',  # 38/40/41mm - три размера
        r'\b\d+\/\d+\s*mm\b',  # 42/44mm - два размера
        r'\/\d+mm\b',  # /49mm - оставшиеся части после слэша
        r'\b\d+mm\b',  # 40mm - одиночные размеры
        r'\b[SML]\/[ML]\b',  # S/M, M/L
        r'\b(small|medium|large|xs|s|m|l|xl|xxl)\b',
    ]

    MODEL_PATTERNS = [
        # Модели телефонов - порядок важен!
        r'\bSamsung\s+Galaxy\s+[A-Z]\d+\+?\b',  # Samsung Galaxy S21 - полная форма сначала
        r'\bGalaxy\s+[A-Z]\d+\+?\b',  # Galaxy S21
        r'\biPhone\s+\d+\s*(Pro|Plus|Max|Mini)?\b',  # iPhone 12, iPhone 13 Pro
        r'\b(Samsung|Apple|Xiaomi|Huawei)\b',  # Производители отдельно
        # Общие паттерны моделей
        r'\b\d{2,4}[A-Z]+\b',  # 2021G, KB-UM-104
    ]

    def __init__(self, xml_file: str):
        self.xml_file = xml_file
        self.tree = None
        self.root = None

    def load_xml(self) -> bool:
        """Загрузить XML файл"""
        try:
            self.tree = ET.parse(self.xml_file)
            self.root = self.tree.getroot()
            return True
        except Exception as e:
            print(f"Ошибка загрузки XML: {e}")
            return False

    def analyze_categories(self) -> CategoryStats:
        """Анализ категорий"""
        stats = CategoryStats()

        if self.root is None:
            return stats

        for product in self.root.findall('artikal'):
            kat1 = product.findtext('kategorija1', '')
            kat2 = product.findtext('kategorija2', '')
            kat3 = product.findtext('kategorija3', '')

            # Собираем уникальные категории
            if kat1:
                stats.level1.add(kat1)
                stats.category_product_count[kat1] += 1

            if kat2:
                full_kat2 = f"{kat1} > {kat2}" if kat1 else kat2
                stats.level2.add(full_kat2)
                stats.category_product_count[full_kat2] += 1

            if kat3:
                full_kat3 = f"{kat1} > {kat2} > {kat3}" if kat1 and kat2 else kat3
                stats.level3.add(full_kat3)
                stats.category_product_count[full_kat3] += 1
                stats.total += 1

        # Топ категорий по количеству товаров
        stats.top_categories = [
            {"category": cat, "product_count": count}
            for cat, count in sorted(stats.category_product_count.items(), key=lambda x: x[1], reverse=True)[:20]
        ]

        return stats

    def analyze_attributes(self) -> AttributeStats:
        """Анализ атрибутов"""
        stats = AttributeStats()

        # Атрибуты для анализа
        attr_fields = [
            'uvoznik', 'godinaUvoza', 'zemljaPorekla',
            'dostupan', 'naAkciji', 'barKod'
        ]

        if self.root is None:
            return stats

        for product in self.root.findall('artikal'):
            for attr in attr_fields:
                value = product.findtext(attr, '').strip()
                if value:
                    stats.unique_values[attr].add(value)
                    stats.value_counts[attr][value] += 1

        stats.detected = attr_fields

        return stats

    def extract_variant_attributes(self, product_name: str) -> Dict[str, Optional[str]]:
        """Извлечь атрибуты варианта из названия"""
        attributes = {}

        # Цвет
        for pattern in self.COLOR_PATTERNS:
            match = re.search(pattern, product_name, re.IGNORECASE)
            if match:
                attributes['color'] = match.group(0)
                break

        # Размер
        for pattern in self.SIZE_PATTERNS:
            match = re.search(pattern, product_name, re.IGNORECASE)
            if match:
                attributes['size'] = match.group(0)
                break

        # Модель
        for pattern in self.MODEL_PATTERNS:
            match = re.search(pattern, product_name, re.IGNORECASE)
            if match:
                attributes['model'] = match.group(0)
                break

        return attributes

    def extract_base_name(self, product_name: str) -> str:
        """Извлечь базовое название без вариантов"""
        name = product_name

        # Убираем цвета
        for pattern in self.COLOR_PATTERNS:
            name = re.sub(pattern, '', name, flags=re.IGNORECASE)

        # Убираем размеры
        for pattern in self.SIZE_PATTERNS:
            name = re.sub(pattern, '', name, flags=re.IGNORECASE)

        # Убираем модели
        for pattern in self.MODEL_PATTERNS:
            name = re.sub(pattern, '', name, flags=re.IGNORECASE)

        # Очистка пробелов
        name = re.sub(r'\s+', ' ', name).strip()

        return name

    def detect_variants(self, min_group_size: int = 2, min_confidence: float = 0.7) -> List[VariantGroup]:
        """Детектировать потенциальные группы вариантов"""
        groups = defaultdict(list)

        if self.root is None:
            return []

        # Группируем по базовому названию
        for product in self.root.findall('artikal'):
            name = product.findtext('naziv', '').strip()
            if not name:
                continue

            base_name = self.extract_base_name(name)
            variant_attrs = self.extract_variant_attributes(name)

            # Только товары с вариантными атрибутами
            if variant_attrs:
                groups[base_name].append({
                    'name': name,
                    'sku': product.findtext('sifra', ''),
                    'attributes': variant_attrs
                })

        # Фильтруем группы
        variant_groups = []
        for base_name, products in groups.items():
            if len(products) < min_group_size:
                continue

            # Собираем уникальные атрибуты
            all_attrs = set()
            for p in products:
                all_attrs.update(p['attributes'].keys())

            # Confidence: процент товаров с вариантными атрибутами
            confidence = len(products) / max(len(products), 1)

            if confidence >= min_confidence:
                variant_groups.append(VariantGroup(
                    base_name=base_name,
                    products=products,
                    variant_count=len(products),
                    variant_attributes=all_attrs,
                    confidence=confidence
                ))

        # Сортировка по количеству вариантов
        variant_groups.sort(key=lambda x: x.variant_count, reverse=True)

        return variant_groups

    def analyze_images(self) -> ImageStats:
        """Анализ изображений"""
        stats = ImageStats()

        if self.root is None:
            return stats

        total_products = len(self.root.findall('artikal'))
        image_counts = []

        for product in self.root.findall('artikal'):
            slike = product.find('slike')
            if slike is not None:
                images = slike.findall('slika')
                if images:
                    stats.total_products_with_images += 1
                    image_count = len(images)
                    stats.total_images += image_count
                    image_counts.append(image_count)

        if image_counts:
            stats.avg_images_per_product = sum(image_counts) / len(image_counts)
            stats.max_images_per_product = max(image_counts)

        if total_products > 0:
            stats.percentage = (stats.total_products_with_images / total_products) * 100

        return stats

    def analyze(self) -> AnalysisResult:
        """Полный анализ прайса"""
        if not self.load_xml():
            return None

        total_products = len(self.root.findall('artikal'))

        print(f"📊 Анализирую {total_products} товаров...")

        # Категории
        print("  🏷️  Анализ категорий...")
        category_stats = self.analyze_categories()

        # Атрибуты
        print("  🔧 Анализ атрибутов...")
        attribute_stats = self.analyze_attributes()

        # Варианты
        print("  🎨 Детекция вариантов...")
        variant_groups = self.detect_variants()

        # Изображения
        print("  📸 Анализ изображений...")
        image_stats = self.analyze_images()

        return AnalysisResult(
            total_products=total_products,
            categories={
                "total": category_stats.total,
                "level1": len(category_stats.level1),
                "level2": len(category_stats.level2),
                "level3": len(category_stats.level3),
                "unique_level1": sorted(list(category_stats.level1)),
                "unique_level2": sorted(list(category_stats.level2))[:50],  # Top 50
                "unique_level3": sorted(list(category_stats.level3))[:100],  # Top 100
                "top_categories": category_stats.top_categories
            },
            attributes={
                "detected": attribute_stats.detected,
                "unique_values": {
                    k: sorted(list(v))[:50]  # Top 50 values per attribute
                    for k, v in attribute_stats.unique_values.items()
                },
                "value_distribution": {
                    k: sorted(v.items(), key=lambda x: x[1], reverse=True)[:20]
                    for k, v in attribute_stats.value_counts.items()
                }
            },
            variants={
                "potential_groups": len(variant_groups),
                "products_affected": sum(g.variant_count for g in variant_groups),
                "top_variant_groups": [
                    {
                        "base_name": g.base_name,
                        "variant_count": g.variant_count,
                        "attributes": list(g.variant_attributes),
                        "confidence": g.confidence,
                        "examples": g.products[:5]  # Top 5 examples
                    }
                    for g in variant_groups[:20]  # Top 20 groups
                ]
            },
            images={
                "total_products_with_images": image_stats.total_products_with_images,
                "total_images": image_stats.total_images,
                "avg_images_per_product": round(image_stats.avg_images_per_product, 2),
                "max_images_per_product": image_stats.max_images_per_product,
                "percentage": round(image_stats.percentage, 1)
            }
        )


def main():
    parser = argparse.ArgumentParser(description='Analyze Digital Vision XML price list')
    parser.add_argument('--file', required=True, help='Path to Digital Vision XML file')
    parser.add_argument('--output', required=True, help='Output JSON file path')
    parser.add_argument('--min-variants', type=int, default=2, help='Minimum variants in group (default: 2)')
    parser.add_argument('--confidence', type=float, default=0.7, help='Minimum confidence (default: 0.7)')

    args = parser.parse_args()

    # Проверка файла
    if not Path(args.file).exists():
        print(f"❌ Файл не найден: {args.file}")
        return 1

    print(f"🚀 Запуск анализа Digital Vision прайса")
    print(f"📁 Файл: {args.file}")
    print()

    # Анализ
    analyzer = DigitalVisionAnalyzer(args.file)
    result = analyzer.analyze()

    if result is None:
        print("❌ Ошибка анализа")
        return 1

    # Сохранение результатов
    output_data = asdict(result)
    with open(args.output, 'w', encoding='utf-8') as f:
        json.dump(output_data, f, indent=2, ensure_ascii=False)

    print()
    print("✅ Анализ завершен!")
    print(f"📄 Результаты сохранены в: {args.output}")
    print()
    print("📊 Краткая статистика:")
    print(f"  📦 Всего товаров: {result.total_products}")
    print(f"  🏷️  Категорий (level 1/2/3): {result.categories['level1']}/{result.categories['level2']}/{result.categories['level3']}")
    print(f"  🎨 Потенциальных групп вариантов: {result.variants['potential_groups']}")
    print(f"  📸 Товаров с изображениями: {result.images['total_products_with_images']} ({result.images['percentage']}%)")
    print()

    return 0


if __name__ == '__main__':
    exit(main())
