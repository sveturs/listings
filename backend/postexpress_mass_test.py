#!/usr/bin/env python3
"""
Post Express Mass Shipment Testing Script
Создает множество тестовых посылок для демонстрации интеграции
"""

import requests
import json
import time
from datetime import datetime
from typing import List, Dict, Any
import sys

# Конфигурация
BACKEND_URL = "http://localhost:3000"
OUTPUT_FILE = "/tmp/postexpress_mass_test_results.json"
TOKEN_FILE = "/tmp/token"

# Тестовые данные (соответствуют TestShipmentRequest структуре)
TEST_DATA = {
    "standard": {
        "count": 10,
        "template": {
            "recipient_name": "Test Recipient",
            "recipient_phone": "0641234567",
            "recipient_email": "recipient@test.rs",
            "recipient_city": "Novi Sad",
            "recipient_address": "Somborska 32",
            "recipient_zip": "21000",
            "sender_name": "Sve Tu d.o.o.",
            "sender_phone": "0641234567",
            "sender_email": "b2b@svetu.rs",
            "sender_city": "Beograd",
            "sender_address": "Bulevar kralja Aleksandra 73",
            "sender_zip": "11000",
            "weight": 500,  # граммы
            "content": "Test package from SVETU platform",
            "cod_amount": 0,
            "insured_value": 10000,  # RSD
            "services": "PNA",
            "delivery_method": "K",
            "delivery_type": "standard",
            "payment_method": "POF",
            "id_rukovanje": 71
        }
    },
    "cod": {
        "count": 10,
        "template": {
            "recipient_name": "COD Test Recipient",
            "recipient_phone": "0641234567",
            "recipient_email": "recipient@test.rs",
            "recipient_city": "Novi Sad",
            "recipient_address": "Somborska 32",
            "recipient_zip": "21000",
            "sender_name": "Sve Tu d.o.o.",
            "sender_phone": "0641234567",
            "sender_email": "b2b@svetu.rs",
            "sender_city": "Beograd",
            "sender_address": "Bulevar kralja Aleksandra 73",
            "sender_zip": "11000",
            "weight": 500,
            "content": "COD package from SVETU",
            "cod_amount": 5000,  # RSD
            "insured_value": 15000,
            "services": "PNA",
            "delivery_method": "K",
            "delivery_type": "cod",
            "payment_method": "POF",
            "id_rukovanje": 71
        }
    },
    "parcel_locker": {
        "count": 10,
        "template": {
            "recipient_name": "Parcel Locker Recipient",
            "recipient_phone": "0641234567",
            "recipient_email": "recipient@test.rs",
            "parcel_locker_code": "PAK001",
            "sender_name": "Sve Tu d.o.o.",
            "sender_phone": "0641234567",
            "sender_email": "b2b@svetu.rs",
            "sender_city": "Beograd",
            "sender_address": "Bulevar kralja Aleksandra 73",
            "sender_zip": "11000",
            "weight": 500,
            "content": "Parcel locker package",
            "cod_amount": 0,
            "insured_value": 8000,
            "services": "PNA",
            "delivery_method": "PAK",
            "delivery_type": "parcel_locker",
            "payment_method": "POF",
            "id_rukovanje": 85
        }
    }
}

# Различные суммы для COD (RSD)
COD_AMOUNTS = [3000, 5000, 7500, 10000, 15000, 20000, 25000, 30000, 40000, 50000]

# Различные коды паккетоматов
PARCEL_LOCKER_CODES = [
    "PAK001", "PAK002", "PAK003", "PAK004", "PAK005",
    "PAK006", "PAK007", "PAK008", "PAK009", "PAK010"
]


def load_token() -> str:
    """Загрузить JWT токен из файла"""
    try:
        with open(TOKEN_FILE, 'r') as f:
            return f.read().strip()
    except FileNotFoundError:
        print(f"❌ Файл токена не найден: {TOKEN_FILE}")
        sys.exit(1)


def create_shipment(shipment_type: str, data: Dict[str, Any], index: int, token: str) -> Dict[str, Any]:
    """Создать одну посылку"""

    # Модифицируем данные для уникальности
    modified_data = data.copy()

    if shipment_type == "cod":
        # Используем разные суммы COD
        modified_data["cod_amount"] = COD_AMOUNTS[index % len(COD_AMOUNTS)]
        modified_data["content"] = f"COD test #{index+1} - Amount: {modified_data['cod_amount']} RSD"

    elif shipment_type == "parcel_locker":
        # Используем разные коды паккетоматов
        modified_data["parcel_locker_code"] = PARCEL_LOCKER_CODES[index % len(PARCEL_LOCKER_CODES)]
        modified_data["content"] = f"Parcel locker test #{index+1} - Code: {modified_data['parcel_locker_code']}"

    else:  # standard
        modified_data["content"] = f"Standard delivery test #{index+1}"

    # Отправка запроса
    url = f"{BACKEND_URL}/api/v1/postexpress/test/shipment"
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}"
    }

    start_time = time.time()
    try:
        response = requests.post(url, json=modified_data, headers=headers, timeout=30)
        elapsed_time = (time.time() - start_time) * 1000  # ms

        result = {
            "type": shipment_type,
            "index": index + 1,
            "success": response.status_code == 200,
            "status_code": response.status_code,
            "elapsed_ms": round(elapsed_time, 2),
            "request": modified_data,
            "response": response.json() if response.status_code == 200 else response.text,
            "timestamp": datetime.now().isoformat()
        }

        if response.status_code == 200:
            resp_data = response.json()
            # Результат в data.success и data.tracking_number
            success = resp_data.get("success", False)
            data_obj = resp_data.get("data", {})

            if isinstance(data_obj, dict) and data_obj.get("success"):
                tracking = data_obj.get("tracking_number", "N/A")
                cost = data_obj.get("cost", "N/A")
                print(f"  ✅ {shipment_type.upper()} #{index+1}: {tracking} | Cost: {cost} RSD | {elapsed_time:.0f}ms")
            else:
                errors = data_obj.get("errors", [])
                print(f"  ❌ {shipment_type.upper()} #{index+1}: API Error: {errors} | {elapsed_time:.0f}ms")
                result["success"] = False
        else:
            print(f"  ❌ {shipment_type.upper()} #{index+1}: HTTP {response.status_code} | {elapsed_time:.0f}ms")

        return result

    except requests.exceptions.Timeout:
        print(f"  ⏱️  {shipment_type.upper()} #{index+1}: Timeout после 30 секунд")
        return {
            "type": shipment_type,
            "index": index + 1,
            "success": False,
            "error": "Timeout",
            "timestamp": datetime.now().isoformat()
        }
    except Exception as e:
        print(f"  ❌ {shipment_type.upper()} #{index+1}: Ошибка - {str(e)}")
        return {
            "type": shipment_type,
            "index": index + 1,
            "success": False,
            "error": str(e),
            "timestamp": datetime.now().isoformat()
        }


def main():
    print("=" * 80)
    print("🚀 POST EXPRESS MASS TESTING SCRIPT")
    print("=" * 80)
    print()

    # Загрузка токена
    print("📝 Загрузка JWT токена...")
    token = load_token()
    print(f"✅ Токен загружен: {token[:20]}...")
    print()

    results = []

    # Standard Delivery
    print("📦 STANDARD DELIVERY")
    print("-" * 80)
    for i in range(TEST_DATA["standard"]["count"]):
        result = create_shipment("standard", TEST_DATA["standard"]["template"], i, token)
        results.append(result)
        time.sleep(1)  # Пауза между запросами
    print()

    # COD Delivery
    print("💰 COD DELIVERY (Cash on Delivery)")
    print("-" * 80)
    for i in range(TEST_DATA["cod"]["count"]):
        result = create_shipment("cod", TEST_DATA["cod"]["template"], i, token)
        results.append(result)
        time.sleep(1)
    print()

    # Parcel Locker
    print("🏪 PARCEL LOCKER DELIVERY")
    print("-" * 80)
    for i in range(TEST_DATA["parcel_locker"]["count"]):
        result = create_shipment("parcel_locker", TEST_DATA["parcel_locker"]["template"], i, token)
        results.append(result)
        time.sleep(1)
    print()

    # Статистика
    print("=" * 80)
    print("📊 СТАТИСТИКА")
    print("=" * 80)

    total = len(results)
    successful = len([r for r in results if r.get("success", False)])
    failed = total - successful

    print(f"Всего попыток: {total}")
    print(f"✅ Успешных: {successful} ({successful/total*100:.1f}%)")
    print(f"❌ Неудачных: {failed} ({failed/total*100:.1f}%)")
    print()

    # Tracking numbers
    tracking_numbers = []
    for r in results:
        if r.get("success") and isinstance(r.get("response"), dict):
            data_obj = r["response"].get("data", {})
            if isinstance(data_obj, dict):
                tn = data_obj.get("tracking_number")
                if tn:
                    tracking_numbers.append({
                        "type": r["type"],
                        "tracking": tn,
                        "cost": data_obj.get("cost")
                    })

    print(f"📋 Tracking Numbers: {len(tracking_numbers)}")
    for tn in tracking_numbers:
        cost_str = str(tn['cost']) if tn['cost'] is not None else "N/A"
        print(f"  {tn['type']:15} | {tn['tracking']:15} | {cost_str:>6} RSD")
    print()

    # Сохранение результатов
    output = {
        "metadata": {
            "timestamp": datetime.now().isoformat(),
            "total_tests": total,
            "successful": successful,
            "failed": failed,
            "backend_url": BACKEND_URL
        },
        "tracking_numbers": tracking_numbers,
        "detailed_results": results
    }

    with open(OUTPUT_FILE, 'w', encoding='utf-8') as f:
        json.dump(output, f, ensure_ascii=False, indent=2)

    print(f"💾 Результаты сохранены: {OUTPUT_FILE}")
    print()
    print("=" * 80)
    print("✅ MASS TESTING ЗАВЕРШЁН!")
    print("=" * 80)


if __name__ == "__main__":
    main()
