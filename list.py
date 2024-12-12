import os
import re

# Путь к файлу проекта
PROJECT_FILE = "/data/proj/hostel-booking-system/collected_content.txt"

# Оглавление проекта, пример структуры
SECTIONS = {
    1: {"name": "Общее", "keywords": ["server", "config"]},
    2: {"name": "Автопрокат", "keywords": ["car", "booking"]},
    3: {"name": "Бронирование жилья", "keywords": ["room", "bed"]},
    4: {"name": "Пользовательская система", "keywords": ["auth", "user"]},
}

def parse_project(file_path, keywords):
    """
    Извлекает части кода, содержащие указанные ключевые слова.
    """
    with open(file_path, "r", encoding="utf-8") as f:
        lines = f.readlines()

    result = []
    inside_section = False
    for line in lines:
        if any(keyword in line for keyword in keywords):
            inside_section = True
        if inside_section:
            result.append(line)
            if line.strip() == "":  # Конец секции по пустой строке
                inside_section = False

    return result

def save_to_file(content, output_path):
    """Сохраняет содержимое в файл."""
    with open(output_path, "w", encoding="utf-8") as f:
        f.writelines(content)

def main():
    print("Оглавление проекта:")
    for key, section in SECTIONS.items():
        print(f"{key}. {section['name']}")

    try:
        choice = int(input("\nВыберите номер раздела: "))
        if choice not in SECTIONS:
            print("Неверный выбор! Попробуйте снова.")
            return

        section = SECTIONS[choice]
        print(f"\nИзвлекаем раздел: {section['name']}")

        # Извлекаем код из проекта
        extracted_code = parse_project(PROJECT_FILE, section["keywords"])

        # Сохраняем в файл
        output_file = f"{section['name']}_code.txt"
        save_to_file(extracted_code, output_file)

        print(f"Код сохранен в файл: {output_file}")

    except ValueError:
        print("Ошибка ввода! Введите номер раздела.")

if __name__ == "__main__":
    main()
