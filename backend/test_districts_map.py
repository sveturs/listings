#!/usr/bin/env python3
import subprocess
import time
import os

def run_test_step(step_name, command):
    """Выполнить шаг теста через claude"""
    print(f"\n{'='*60}")
    print(f"Шаг: {step_name}")
    print('='*60)
    
    try:
        result = subprocess.run(
            ['claude', '-p', '--dangerously-skip-permissions', command],
            capture_output=True,
            text=True,
            timeout=45
        )
        
        if result.returncode == 0:
            print("✓ Успешно выполнено")
            if result.stdout:
                print(f"Результат: {result.stdout[:500]}...")
        else:
            print(f"✗ Ошибка: {result.stderr}")
            
        return result.returncode == 0
        
    except subprocess.TimeoutExpired:
        print("✗ Превышен таймаут")
        return False
    except Exception as e:
        print(f"✗ Исключение: {str(e)}")
        return False

def main():
    """Основная функция тестирования"""
    
    # Создаем файл логов
    log_file = "/tmp/districts-test.log"
    
    with open(log_file, "w") as f:
        f.write("Тестирование фильтра районов на карте\n")
        f.write(f"Дата: {time.strftime('%Y-%m-%d %H:%M:%S')}\n")
        f.write("="*60 + "\n\n")
    
    # Шаг 1: Проверка страницы карты
    success = run_test_step(
        "Проверка доступности карты",
        "Проверь что страница http://localhost:3001/ru/map доступна используя curl"
    )
    
    if not success:
        print("\n❌ Страница карты недоступна!")
        return
    
    # Шаг 2: Открытие карты в браузере
    success = run_test_step(
        "Открытие карты",
        """Используй MCP Playwright:
        1. Открой Chrome браузер
        2. Перейди на http://localhost:3001/ru/map
        3. Подожди 3 секунды
        4. Сделай скриншот /tmp/map-initial.png
        5. Найди кнопку или переключатель 'Поиск по районам'"""
    )
    
    if not success:
        print("\n❌ Не удалось открыть карту!")
        return
        
    # Шаг 3: Переход к Нови Саду
    time.sleep(2)
    success = run_test_step(
        "Переход к Нови Саду",
        """Используй MCP Playwright:
        1. Перейди на http://localhost:3001/ru/map?lat=45.2671&lng=19.8335
        2. Подожди 3 секунды для загрузки
        3. Сделай скриншот /tmp/map-novi-sad.png
        4. Проверь есть ли селектор районов на странице"""
    )
    
    # Финальный отчет
    with open(log_file, "a") as f:
        f.write(f"\n\nЗавершено: {time.strftime('%Y-%m-%d %H:%M:%S')}\n")
        f.write("Созданные файлы:\n")
        for file in ["/tmp/map-initial.png", "/tmp/map-novi-sad.png"]:
            if os.path.exists(file):
                f.write(f"✓ {file}\n")
            else:
                f.write(f"✗ {file} - не создан\n")
    
    print(f"\n\n📄 Отчет сохранен в: {log_file}")
    print("🖼️  Скриншоты:")
    os.system("ls -la /tmp/map-*.png 2>/dev/null | tail -5")

if __name__ == "__main__":
    main()