#!/usr/bin/env python3
"""
Generate meta_keywords for all L2/L3 categories based on category names.
Uses translation mappings and category-specific keyword expansions.
"""

import psycopg2
import json
from typing import Dict, List, Tuple

# Database connection
DB_CONFIG = {
    "host": "localhost",
    "port": 35434,
    "database": "listings_dev_db",
    "user": "listings_user",
    "password": "listings_password"
}

# Serbian to English/Russian translations for common words
TRANSLATIONS = {
    # Electronics
    "telefon": {"en": "phone, mobile", "ru": "телефон, мобильный"},
    "pametni": {"en": "smart", "ru": "умный, смарт"},
    "laptop": {"en": "laptop, notebook", "ru": "ноутбук, лэптоп"},
    "racunar": {"en": "computer, PC", "ru": "компьютер, ПК"},
    "tablet": {"en": "tablet", "ru": "планшет"},
    "kamera": {"en": "camera", "ru": "камера"},
    "slusalice": {"en": "headphones, earphones", "ru": "наушники"},
    "zvucnik": {"en": "speaker", "ru": "колонка, динамик"},
    "monitor": {"en": "monitor, display", "ru": "монитор, дисплей"},
    "stampac": {"en": "printer", "ru": "принтер"},
    "tastatura": {"en": "keyboard", "ru": "клавиатура"},
    "mis": {"en": "mouse", "ru": "мышь, мышка"},
    "punjac": {"en": "charger", "ru": "зарядка, зарядное"},
    "baterija": {"en": "battery", "ru": "батарея, аккумулятор"},
    "kabl": {"en": "cable, cord", "ru": "кабель, провод"},
    "adapter": {"en": "adapter", "ru": "адаптер"},
    "gaming": {"en": "gaming", "ru": "игровой, геймерский"},
    "konzola": {"en": "console", "ru": "консоль, приставка"},

    # Clothing
    "majica": {"en": "t-shirt, tee, top", "ru": "футболка, майка"},
    "kosulja": {"en": "shirt", "ru": "рубашка"},
    "pantalone": {"en": "pants, trousers", "ru": "брюки, штаны"},
    "farmerke": {"en": "jeans, denim", "ru": "джинсы"},
    "jakna": {"en": "jacket", "ru": "куртка"},
    "kaput": {"en": "coat", "ru": "пальто"},
    "dzemperi": {"en": "sweater, jumper", "ru": "свитер, джемпер"},
    "haljina": {"en": "dress", "ru": "платье"},
    "suknja": {"en": "skirt", "ru": "юбка"},
    "sorts": {"en": "shorts", "ru": "шорты"},
    "patike": {"en": "sneakers, trainers", "ru": "кроссовки"},
    "cipele": {"en": "shoes", "ru": "туфли, ботинки"},
    "cizme": {"en": "boots", "ru": "сапоги"},
    "sandale": {"en": "sandals", "ru": "сандалии"},
    "torba": {"en": "bag", "ru": "сумка"},
    "ranac": {"en": "backpack", "ru": "рюкзак"},
    "kais": {"en": "belt", "ru": "ремень"},
    "muski": {"en": "mens, male", "ru": "мужской"},
    "zenski": {"en": "womens, female", "ru": "женский"},
    "deciji": {"en": "kids, children", "ru": "детский"},

    # Home & Garden
    "namestaj": {"en": "furniture", "ru": "мебель"},
    "sofa": {"en": "sofa, couch", "ru": "диван, софа"},
    "krevet": {"en": "bed", "ru": "кровать"},
    "sto": {"en": "table, desk", "ru": "стол"},
    "stolica": {"en": "chair", "ru": "стул"},
    "ormar": {"en": "wardrobe, closet", "ru": "шкаф"},
    "polica": {"en": "shelf", "ru": "полка"},
    "lampa": {"en": "lamp, light", "ru": "лампа, светильник"},
    "tepih": {"en": "carpet, rug", "ru": "ковёр"},
    "zavesa": {"en": "curtain", "ru": "штора, занавеска"},
    "posteljina": {"en": "bedding, linens", "ru": "постельное бельё"},
    "jastuk": {"en": "pillow", "ru": "подушка"},
    "dusek": {"en": "mattress", "ru": "матрас"},
    "basta": {"en": "garden", "ru": "сад, огород"},
    "biljka": {"en": "plant", "ru": "растение"},
    "cvece": {"en": "flowers", "ru": "цветы"},
    "alat": {"en": "tools", "ru": "инструменты"},
    "kuhinja": {"en": "kitchen", "ru": "кухня"},
    "kupatilo": {"en": "bathroom", "ru": "ванная"},

    # Appliances
    "ves-masina": {"en": "washing machine, washer", "ru": "стиральная машина"},
    "frizider": {"en": "refrigerator, fridge", "ru": "холодильник"},
    "sudoper": {"en": "dishwasher", "ru": "посудомоечная машина"},
    "rerna": {"en": "oven", "ru": "духовка"},
    "mikrotalasna": {"en": "microwave", "ru": "микроволновка"},
    "usisivac": {"en": "vacuum cleaner", "ru": "пылесос"},
    "klima": {"en": "air conditioner, AC", "ru": "кондиционер"},
    "grejalica": {"en": "heater", "ru": "обогреватель"},
    "fen": {"en": "hair dryer", "ru": "фен"},
    "pegla": {"en": "iron", "ru": "утюг"},
    "blender": {"en": "blender", "ru": "блендер"},
    "tosteri": {"en": "toaster", "ru": "тостер"},
    "aparat-kafa": {"en": "coffee maker, coffee machine", "ru": "кофеварка, кофемашина"},

    # Auto
    "auto": {"en": "car, auto, vehicle", "ru": "авто, автомобиль, машина"},
    "delovi": {"en": "parts, spare parts", "ru": "запчасти, детали"},
    "gume": {"en": "tires, tyres", "ru": "шины, резина"},
    "ulje": {"en": "oil", "ru": "масло"},
    "kocnice": {"en": "brakes", "ru": "тормоза"},
    "filter": {"en": "filter", "ru": "фильтр"},
    "akumulator": {"en": "battery", "ru": "аккумулятор"},
    "motor": {"en": "engine, motor", "ru": "двигатель, мотор"},
    "karoserija": {"en": "body parts", "ru": "кузов"},
    "motocikl": {"en": "motorcycle, motorbike", "ru": "мотоцикл"},
    "bicikl": {"en": "bicycle, bike", "ru": "велосипед"},

    # Sports
    "sport": {"en": "sports", "ru": "спорт"},
    "fitnes": {"en": "fitness, gym", "ru": "фитнес"},
    "trening": {"en": "training, workout", "ru": "тренировка"},
    "oprema": {"en": "equipment, gear", "ru": "оборудование, снаряжение"},
    "lopta": {"en": "ball", "ru": "мяч"},
    "fudbal": {"en": "football, soccer", "ru": "футбол"},
    "kosarka": {"en": "basketball", "ru": "баскетбол"},
    "tenis": {"en": "tennis", "ru": "теннис"},
    "plivanje": {"en": "swimming", "ru": "плавание"},
    "planinarenje": {"en": "hiking, trekking", "ru": "походы, треккинг"},
    "kampovanje": {"en": "camping", "ru": "кемпинг"},
    "biciklizam": {"en": "cycling", "ru": "велоспорт"},

    # Beauty
    "kozmetika": {"en": "cosmetics, makeup", "ru": "косметика"},
    "sminka": {"en": "makeup", "ru": "макияж"},
    "parfem": {"en": "perfume, fragrance", "ru": "духи, парфюм"},
    "nega": {"en": "care", "ru": "уход"},
    "krema": {"en": "cream", "ru": "крем"},
    "sampon": {"en": "shampoo", "ru": "шампунь"},
    "kosa": {"en": "hair", "ru": "волосы"},
    "lice": {"en": "face", "ru": "лицо"},
    "telo": {"en": "body", "ru": "тело"},
    "nokti": {"en": "nails", "ru": "ногти"},

    # Kids
    "beba": {"en": "baby", "ru": "малыш, ребёнок"},
    "deca": {"en": "kids, children", "ru": "дети"},
    "igracka": {"en": "toy", "ru": "игрушка"},
    "kolica": {"en": "stroller, pram", "ru": "коляска"},
    "krevetac": {"en": "crib, cot", "ru": "кроватка"},
    "pelene": {"en": "diapers", "ru": "подгузники"},
    "hrana-beba": {"en": "baby food", "ru": "детское питание"},

    # Jewelry
    "nakit": {"en": "jewelry", "ru": "украшения, ювелирные изделия"},
    "sat": {"en": "watch", "ru": "часы"},
    "prsten": {"en": "ring", "ru": "кольцо"},
    "ogrlica": {"en": "necklace", "ru": "ожерелье, колье"},
    "narukvica": {"en": "bracelet", "ru": "браслет"},
    "mindjuse": {"en": "earrings", "ru": "серьги"},
    "zlato": {"en": "gold", "ru": "золото"},
    "srebro": {"en": "silver", "ru": "серебро"},

    # General
    "nov": {"en": "new", "ru": "новый"},
    "polovan": {"en": "used, second-hand", "ru": "б/у, подержанный"},
    "originalni": {"en": "original, genuine", "ru": "оригинальный"},
    "profesionalni": {"en": "professional", "ru": "профессиональный"},
    "bezicni": {"en": "wireless", "ru": "беспроводной"},
    "električni": {"en": "electric", "ru": "электрический"},
}

# Category-specific keyword expansions
CATEGORY_KEYWORDS = {
    "pametni-telefoni": {
        "sr": "android, ios, 5g, dual sim, smartphone, mobilni, ekran, memorija, ram",
        "en": "android, ios, 5g, dual sim, smartphone, mobile, screen, memory, ram, unlocked",
        "ru": "андроид, ios, 5g, две сим, смартфон, экран, память, разблокированный"
    },
    "laptopi": {
        "sr": "intel, amd, ryzen, ssd, nvidia, gejming laptop, radni laptop, ultrabook",
        "en": "intel, amd, ryzen, ssd, nvidia, gaming laptop, work laptop, ultrabook, lightweight",
        "ru": "интел, амд, ssd, nvidia, игровой ноутбук, рабочий ноутбук, ультрабук"
    },
    "televizori": {
        "sr": "led, oled, qled, smart tv, 4k, 8k, uhd, hdr, dijagonala, inch",
        "en": "led, oled, qled, smart tv, 4k, 8k, uhd, hdr, screen size, inch",
        "ru": "led, oled, qled, смарт тв, 4k, 8k, uhd, hdr, диагональ, дюйм"
    },
    "frizideri": {
        "sr": "kombinovani, side by side, no frost, zamrzivac, ugradni, mini frizider",
        "en": "combo, side by side, no frost, freezer, built-in, mini fridge, french door",
        "ru": "двухкамерный, side by side, no frost, морозильник, встраиваемый, мини"
    },
    "ves-masine": {
        "sr": "pranje, centrifuga, kg, energetska klasa, frontalno punjenje, gornje punjenje",
        "en": "wash, spin, kg, energy class, front load, top load, capacity",
        "ru": "стирка, отжим, кг, энергокласс, фронтальная, вертикальная"
    }
}


def get_categories() -> List[Tuple]:
    """Get all L2/L3 categories without meta_keywords."""
    conn = psycopg2.connect(**DB_CONFIG)
    cur = conn.cursor()

    cur.execute("""
        SELECT id, slug, name->>'sr' as name_sr, name->>'en' as name_en
        FROM categories
        WHERE parent_id IS NOT NULL
        AND (meta_keywords IS NULL OR meta_keywords::text = '{}' OR meta_keywords::text = 'null')
        ORDER BY slug
    """)

    results = cur.fetchall()
    cur.close()
    conn.close()
    return results


def generate_keywords(slug: str, name_sr: str, name_en: str) -> Dict[str, str]:
    """Generate keywords for a category based on its name."""

    # Start with specific keywords if available
    if slug in CATEGORY_KEYWORDS:
        return CATEGORY_KEYWORDS[slug]

    # Generate from name
    words_sr = name_sr.lower().replace("-", " ").split() if name_sr else []
    words_en = name_en.lower().replace("-", " ").split() if name_en else []

    keywords_sr = set()
    keywords_en = set()
    keywords_ru = set()

    # Add original words
    keywords_sr.update(words_sr)
    keywords_en.update(words_en)

    # Translate and expand
    for word in words_sr:
        # Normalize Serbian word
        word_clean = word.replace("š", "s").replace("č", "c").replace("ć", "c").replace("ž", "z").replace("đ", "dj")

        # Check translations
        for sr_key, translations in TRANSLATIONS.items():
            if word in sr_key or sr_key in word or word_clean in sr_key.replace("š", "s"):
                keywords_en.update(translations["en"].split(", "))
                keywords_ru.update(translations["ru"].split(", "))

    # Add common variations
    keywords_sr.add(name_sr.lower() if name_sr else slug.replace("-", " "))
    keywords_en.add(name_en.lower() if name_en else slug.replace("-", " "))

    # Add kupiti/prodati variations for Serbian
    if name_sr:
        keywords_sr.add(f"kupiti {name_sr.lower()}")
        keywords_sr.add(f"prodaja {name_sr.lower()}")

    # Add buy/sell for English
    if name_en:
        keywords_en.add(f"buy {name_en.lower()}")
        keywords_en.add(f"sell {name_en.lower()}")

    # Add Russian variations
    keywords_ru.add("купить")
    keywords_ru.add("продать")

    return {
        "sr": ", ".join(sorted(keywords_sr)[:15]),
        "en": ", ".join(sorted(keywords_en)[:15]),
        "ru": ", ".join(sorted(keywords_ru)[:15]) if keywords_ru else ", ".join(sorted(keywords_en)[:10])
    }


def generate_sql() -> str:
    """Generate SQL migration for all categories."""
    categories = get_categories()

    sql_lines = [
        "-- Auto-generated meta_keywords for L2/L3 categories",
        "-- Generated by generate_keywords.py",
        "",
        "BEGIN;",
        ""
    ]

    for cat_id, slug, name_sr, name_en in categories:
        keywords = generate_keywords(slug, name_sr, name_en)

        # Escape single quotes
        kw_sr = keywords["sr"].replace("'", "''")
        kw_en = keywords["en"].replace("'", "''")
        kw_ru = keywords["ru"].replace("'", "''")

        sql = f"""UPDATE categories SET meta_keywords = jsonb_build_object(
    'sr', '{kw_sr}',
    'en', '{kw_en}',
    'ru', '{kw_ru}'
) WHERE slug = '{slug}';"""

        sql_lines.append(sql)
        sql_lines.append("")

    sql_lines.append("COMMIT;")

    return "\n".join(sql_lines)


def main():
    print("Generating keywords for L2/L3 categories...")

    categories = get_categories()
    print(f"Found {len(categories)} categories without keywords")

    sql = generate_sql()

    # Write to file
    output_file = "/p/github.com/vondi-global/listings/migrations/20251222200001_l2l3_keywords.up.sql"
    with open(output_file, "w") as f:
        f.write(sql)

    print(f"Generated SQL saved to: {output_file}")
    print(f"Total UPDATE statements: {len(categories)}")


if __name__ == "__main__":
    main()
