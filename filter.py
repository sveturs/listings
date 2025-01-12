import re

# Функция для получения списка блоков для удаления
def get_blocks_to_remove():
    blocks = [
        "=== /data/proj/hostel-booking-system/backend/internal/proj/accommodation",
        "=== /data/proj/hostel-booking-system/backend/internal/proj/car",
        "=== /data/proj/hostel-booking-system/frontend/hostel-frontend/src/components/accommodation",
        "=== /data/proj/hostel-booking-system/frontend/hostel-frontend/src/components/car",
        "=== /data/proj/hostel-booking-system/frontend/hostel-frontend/src/pages/car",
        "=== /data/proj/hostel-booking-system/frontend/hostel-frontend/src/pages/accommodation",
    ]
    print("Доступные блоки для удаления:")
    for i, block in enumerate(blocks, 1):
        print(f"{i}. {block}")
    
    choices = input("Введите номера блоков для удаления (через запятую): ")
    selected = [int(choice.strip()) for choice in choices.split(",")]
    return [blocks[i - 1] for i in selected if 0 < i <= len(blocks)]

# Получаем список блоков для удаления
blocks_to_remove = get_blocks_to_remove()

# Открываем файл для чтения и создаем новый файл для записи
input_file = "project_code.txt"
output_file = "filtered_project_code.txt"

# Функция для проверки начала блока
def is_block_start(line):
    return line.strip().startswith("===")

# Читаем исходный файл и записываем только нужные блоки в новый файл
with open(input_file, "r", encoding="utf-8") as infile, open(output_file, "w", encoding="utf-8") as outfile:
    remove_block = False

    for line in infile:
        if is_block_start(line):
            remove_block = any(line.startswith(block) for block in blocks_to_remove)

        if not remove_block:
            outfile.write(line)

print(f"Фильтрация завершена. Результат сохранен в {output_file}")
