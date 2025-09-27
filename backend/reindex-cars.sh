#!/bin/bash

# Переиндексация автомобильных объявлений в OpenSearch

echo "=== Переиндексация автомобильных объявлений ==="
echo "Запуск времени: $(date)"

# Переходим в директорию backend
cd /data/hostel-booking-system/backend

# Проверяем наличие .env файла
if [ ! -f .env ]; then
    echo "❌ Файл .env не найден!"
    exit 1
fi

# Загружаем переменные окружения
export $(cat .env | xargs)

# Компилируем и запускаем программу переиндексации
echo "Компиляция программы переиндексации..."
if go build -o /tmp/reindex cmd/reindex/main.go; then
    echo "✅ Компиляция успешна"

    echo "Запуск переиндексации автомобилей..."
    /tmp/reindex cars

    # Удаляем временный файл
    rm /tmp/reindex
else
    echo "❌ Ошибка компиляции"
    exit 1
fi

echo "Завершено: $(date)"
