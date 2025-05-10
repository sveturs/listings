#!/bin/bash

# Скрипт для проверки TS-файлов в GIS-компонентах с игнорированием проблемных типов i18next
# Использование: ./gis-typecheck.sh [файлы...]

# Переход в директорию скрипта
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

# Установка переменной окружения для старых версий OpenSSL
export NODE_OPTIONS="--openssl-legacy-provider"

# Опции TypeScript для игнорирования проблемных библиотек
TS_OPTIONS="--noEmit --jsx react --skipLibCheck --esModuleInterop --allowSyntheticDefaultImports --moduleResolution node"

# Добавляем опцию для игнорирования директории node_modules/i18next
TS_OPTIONS="$TS_OPTIONS --typeRoots ./node_modules/@types"

# Проверка наличия аргументов
if [ $# -gt 0 ]; then
  echo "Проверка TypeScript файлов: $@"
  
  # Компиляция указанных файлов
  npx tsc $TS_OPTIONS $@
else
  echo "Необходимо указать файлы для проверки"
  echo "Использование: ./gis-typecheck.sh путь/к/файлу.tsx [путь/к/другому/файлу.tsx ...]"
  exit 1
fi

# Проверка результата выполнения
if [ $? -eq 0 ]; then
  echo "✅ TypeScript проверка пройдена"
else
  echo "❌ TypeScript проверка не пройдена"
  exit 1
fi