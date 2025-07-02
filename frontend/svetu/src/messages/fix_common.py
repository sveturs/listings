#!/usr/bin/env python3
import json
import sys

def fix_common_sections(file_path, unified_common):
    """Заменяет все common секции на объединенную версию"""
    
    # Читаем файл построчно
    with open(file_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()
    
    # Найдем все позиции common секций
    common_positions = []
    for i, line in enumerate(lines):
        if '"common":' in line:
            common_positions.append(i)
    
    print(f"Found {len(common_positions)} common sections in {file_path}")
    
    # Удаляем все common секции кроме первой
    result_lines = []
    skip_until = -1
    
    for i, line in enumerate(lines):
        if skip_until >= i:
            continue
            
        # Если это common секция (не первая)
        if '"common":' in line and i not in [common_positions[0]]:
            # Найдем конец этой секции
            brace_count = 0
            found_opening = False
            
            for j in range(i, len(lines)):
                test_line = lines[j]
                for char in test_line:
                    if char == '{':
                        brace_count += 1
                        found_opening = True
                    elif char == '}':
                        brace_count -= 1
                
                if found_opening and brace_count == 0:
                    skip_until = j
                    print(f"Removing common section from line {i+1} to {j+1}")
                    break
        else:
            result_lines.append(line)
    
    # Теперь заменим первую common секцию
    final_lines = []
    replaced_first = False
    
    for i, line in enumerate(result_lines):
        if '"common":' in line and not replaced_first:
            # Найдем конец первой секции
            brace_count = 0
            found_opening = False
            
            for j in range(i, len(result_lines)):
                test_line = result_lines[j]
                for char in test_line:
                    if char == '{':
                        brace_count += 1
                        found_opening = True
                    elif char == '}':
                        brace_count -= 1
                
                if found_opening and brace_count == 0:
                    # Заменяем секцию
                    indent = len(line) - len(line.lstrip())
                    
                    # Создаем новую common секцию
                    final_lines.append(' ' * indent + '"common": {\n')
                    
                    for key, value in unified_common.items():
                        if isinstance(value, dict):
                            final_lines.append(' ' * (indent + 2) + f'"{key}": {{\n')
                            for subkey, subvalue in value.items():
                                final_lines.append(' ' * (indent + 4) + f'"{subkey}": "{subvalue}",\n')
                            # Удаляем последнюю запятую
                            if final_lines[-1].endswith(',\n'):
                                final_lines[-1] = final_lines[-1][:-2] + '\n'
                            final_lines.append(' ' * (indent + 2) + '},\n')
                        else:
                            final_lines.append(' ' * (indent + 2) + f'"{key}": "{value}",\n')
                    
                    # Удаляем последнюю запятую
                    if final_lines[-1].endswith(',\n'):
                        final_lines[-1] = final_lines[-1][:-2] + '\n'
                    
                    final_lines.append(' ' * indent + '},\n')
                    replaced_first = True
                    
                    # Пропускаем оригинальные строки этой секции
                    skip_next = j - i
                    for _ in range(skip_next):
                        next(enumerate(result_lines[i+1:]), None)
                    break
        else:
            if not replaced_first or '"common":' not in line:
                final_lines.append(line)
    
    # Записываем результат
    with open(file_path, 'w', encoding='utf-8') as f:
        f.writelines(final_lines)
    
    print(f"Updated {file_path}")

# Объединенные common секции
ru_unified_common = {
    "next": "Далее",
    "back": "Назад", 
    "save": "Сохранить",
    "submit": "Отправить",
    "edit": "Редактировать",
    "delete": "Удалить",
    "cancel": "Отмена",
    "loading": "Загрузка...",
    "error": "Ошибка",
    "success": "Успешно",
    "saved": "Сохранено",
    "updated": "Обновлено",
    "noData": "Нет данных",
    "tryAgain": "Попробуйте снова",
    "status": "Статус",
    "view": "Просмотр",
    "new": "Новый",
    "minutesAgo": "{minutes} минут назад",
    "hourAgo": "{hours} час назад",
    "loadMore": "Загрузить еще",
    "justNow": "Только что",
    "hoursAgo": "ч. назад",
    "daysAgo": "д. назад",
    "daysAgoWithCount": "{count} д. назад",
    "chat": "Чат",
    "sendMessage": "Написать сообщение",
    "viewAll": "Посмотреть все",
    "documentation": "Документация",
    "close": "Закрыть",
    "add": "Добавить",
    "search": "Поиск",
    "filter": "Фильтр",
    "actions": "Действия",
    "active": "Активный",
    "inactive": "Неактивный",
    "yes": "Да",
    "no": "Нет",
    "confirmDelete": "Вы уверены, что хотите удалить этот элемент?",
    "saveSuccess": "Успешно сохранено",
    "deleteSuccess": "Успешно удалено",
    "total": "Всего",
    "showing": "Показано",
    "of": "из",
    "items": "элементов",
    "translate": "Перевести",
    "continue": "Продолжить",
    "select": "Выбрать",
    "accessDenied": "Доступ запрещен",
    "comingSoon": "Скоро",
    "popular": "Популярное",
    "optional": "Необязательно",
    "closed": "Закрыто",
    "open": "Открыто",
    "searching": "Поиск...",
    "publishing": "Публикация...",
    "publish": "Опубликовать",
    "free_above": "Бесплатно от",
    "viewGrid": "Сеткой",
    "viewList": "Списком",
    "days": {
        "monday": "Понедельник",
        "tuesday": "Вторник",
        "wednesday": "Среда",
        "thursday": "Четверг",
        "friday": "Пятница",
        "saturday": "Суббота",
        "sunday": "Воскресенье"
    }
}

en_unified_common = {
    "next": "Next",
    "back": "Back",
    "save": "Save",
    "submit": "Submit",
    "edit": "Edit",
    "delete": "Delete",
    "cancel": "Cancel",
    "loading": "Loading...",
    "error": "Error",
    "success": "Success",
    "saved": "Saved",
    "updated": "Updated",
    "noData": "No data available",
    "tryAgain": "Try again",
    "status": "Status",
    "view": "View",
    "new": "New",
    "minutesAgo": "{minutes} minutes ago",
    "hourAgo": "{hours} hour ago",
    "loadMore": "Load more",
    "justNow": "Just now",
    "hoursAgo": "h ago",
    "daysAgo": "d ago",
    "daysAgoWithCount": "{count} days ago",
    "chat": "Chat",
    "sendMessage": "Send message",
    "viewAll": "View all",
    "documentation": "Documentation",
    "close": "Close",
    "add": "Add",
    "search": "Search",
    "filter": "Filter",
    "actions": "Actions",
    "active": "Active",
    "inactive": "Inactive",
    "yes": "Yes",
    "no": "No",
    "confirmDelete": "Are you sure you want to delete this item?",
    "saveSuccess": "Successfully saved",
    "deleteSuccess": "Successfully deleted",
    "total": "Total",
    "showing": "Showing",
    "of": "of",
    "items": "items",
    "translate": "Translate",
    "continue": "Continue",
    "select": "Select",
    "accessDenied": "Access Denied",
    "comingSoon": "Coming Soon",
    "popular": "Popular",
    "optional": "Optional",
    "closed": "Closed",
    "open": "Open",
    "searching": "Searching...",
    "publishing": "Publishing...",
    "publish": "Publish",
    "free_above": "Free above",
    "all": "All",
    "reviews": "reviews",
    "notSpecified": "Not specified",
    "rating": "Rating",
    "share": "Share",
    "viewGrid": "Grid view",
    "viewList": "List view",
    "phone": "Phone",
    "email": "Email",
    "days": {
        "monday": "Monday",
        "tuesday": "Tuesday",
        "wednesday": "Wednesday",
        "thursday": "Thursday",
        "friday": "Friday",
        "saturday": "Saturday",
        "sunday": "Sunday"
    }
}

if __name__ == "__main__":
    print("Fixing common sections...")
    
    # Фиксим RU файл
    fix_common_sections("ru.json", ru_unified_common)
    
    # Фиксим EN файл  
    fix_common_sections("en.json", en_unified_common)
    
    print("Done!")