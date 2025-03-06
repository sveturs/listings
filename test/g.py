import re

# Файл ввода и вывода
input_file = "input.sql"
output_file = "output.sql"

# Регулярное выражение для поиска строк
pattern = re.compile(r"^(\s*)\('category', (\d+), 'ru', 'name', '([^']+)', '[^']+', (\d+|NULL), '[^']+', '[^']+'\)")

# Открываем файлы
with open(input_file, "r", encoding="utf-8") as f:
    lines = f.readlines()

output_lines = []
for line in lines:
    match = pattern.search(line)
    if match:
        indent, category_id, name, parent_id = match.groups()
        output_line = f"{indent}('category', {category_id}, 'ru', 'name', '{name}', true, true, NOW(), NOW()),\n"
        output_lines.append(output_line)
    else:
        output_lines.append(line)  # Оставляем строки, не попадающие под шаблон

# Записываем преобразованные данные в новый файл
with open(output_file, "w", encoding="utf-8") as f:
    f.writelines(output_lines)

print("Конвертация завершена. Результат сохранён в", output_file)
