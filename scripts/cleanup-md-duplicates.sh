#!/bin/bash

# Скрипт для очистки дублирующихся MD файлов в проекте
# Автор: System Administrator
# Дата: 2025-09-05

set -e

PROJECT_ROOT="/data/hostel-booking-system"
ARCHIVE_DIR="$PROJECT_ROOT/docs/archive/cleanup-$(date +%Y%m%d)"
LOG_FILE="$PROJECT_ROOT/cleanup-md-$(date +%Y%m%d).log"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "📋 Начинаем очистку дублирующихся MD файлов..." | tee "$LOG_FILE"
echo "Дата: $(date)" | tee -a "$LOG_FILE"
echo "----------------------------------------" | tee -a "$LOG_FILE"

# Функция для безопасного перемещения файлов
safe_move() {
    local src="$1"
    local dst="$2"
    if [ -f "$src" ]; then
        mkdir -p "$(dirname "$dst")"
        mv "$src" "$dst"
        echo -e "${GREEN}✓${NC} Перемещен: $src → $dst" | tee -a "$LOG_FILE"
    else
        echo -e "${YELLOW}⚠${NC} Файл не найден: $src" | tee -a "$LOG_FILE"
    fi
}

# Функция для удаления файла с логированием
safe_delete() {
    local file="$1"
    if [ -f "$file" ]; then
        rm "$file"
        echo -e "${RED}✗${NC} Удален: $file" | tee -a "$LOG_FILE"
    fi
}

# 1. Создаем директорию для архива
echo -e "\n${YELLOW}1. Создание архивной директории...${NC}" | tee -a "$LOG_FILE"
mkdir -p "$ARCHIVE_DIR"

# 2. Удаляем директорию /task/ с дубликатами
echo -e "\n${YELLOW}2. Удаление директории /task/ (полные дубликаты)...${NC}" | tee -a "$LOG_FILE"
if [ -d "$PROJECT_ROOT/task" ]; then
    # Сначала архивируем
    mkdir -p "$ARCHIVE_DIR/task-backup"
    cp -r "$PROJECT_ROOT/task/"*.md "$ARCHIVE_DIR/task-backup/" 2>/dev/null || true
    # Затем удаляем
    rm -rf "$PROJECT_ROOT/task"
    echo -e "${GREEN}✓${NC} Директория /task/ удалена (архив создан)" | tee -a "$LOG_FILE"
else
    echo -e "${YELLOW}⚠${NC} Директория /task/ не найдена" | tee -a "$LOG_FILE"
fi

# 3. Архивируем дневные отчеты по категориям
echo -e "\n${YELLOW}3. Архивирование дневных отчетов (DAY_01 - DAY_29)...${NC}" | tee -a "$LOG_FILE"
mkdir -p "$PROJECT_ROOT/docs/categories/archive/daily-reports-2025-09"

for file in $PROJECT_ROOT/docs/categories/ATTRIBUTE_UNIFICATION_PROGRESS_DAY_*.md; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        # Пропускаем финальный отчет DAY_30_FINAL
        if [[ ! "$filename" =~ DAY_30_FINAL ]]; then
            safe_move "$file" "$PROJECT_ROOT/docs/categories/archive/daily-reports-2025-09/$filename"
        fi
    fi
done

for file in $PROJECT_ROOT/docs/categories/ATTRIBUTE_UNIFICATION_HANDOVER_DAY_*.md; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        safe_move "$file" "$PROJECT_ROOT/docs/categories/archive/daily-reports-2025-09/$filename"
    fi
done

# 4. Поиск и удаление точных дубликатов по MD5
echo -e "\n${YELLOW}4. Поиск точных дубликатов по контрольной сумме...${NC}" | tee -a "$LOG_FILE"

# Создаем временный файл для хешей
TEMP_HASHES=$(mktemp)

# Находим все MD файлы и вычисляем их MD5
find "$PROJECT_ROOT" -name "*.md" -type f ! -path "*/node_modules/*" ! -path "*/.git/*" -exec md5sum {} \; | sort > "$TEMP_HASHES"

# Находим дубликаты
awk '{print $1}' "$TEMP_HASHES" | uniq -d | while read hash; do
    echo -e "\n${YELLOW}Найдены дубликаты с хешем $hash:${NC}" | tee -a "$LOG_FILE"
    grep "^$hash" "$TEMP_HASHES" | while read line; do
        file=$(echo "$line" | cut -d' ' -f2-)
        echo "  - $file" | tee -a "$LOG_FILE"
    done
    
    # Оставляем первый файл, остальные архивируем
    first_file=""
    grep "^$hash" "$TEMP_HASHES" | while read line; do
        file=$(echo "$line" | cut -d' ' -f2-)
        if [ -z "$first_file" ]; then
            first_file="$file"
            echo -e "  ${GREEN}Оставляем: $first_file${NC}" | tee -a "$LOG_FILE"
        else
            relative_path="${file#$PROJECT_ROOT/}"
            archive_path="$ARCHIVE_DIR/duplicates/$relative_path"
            safe_move "$file" "$archive_path"
        fi
    done
done

rm "$TEMP_HASHES"

# 5. Очистка UI/UX дубликатов с номерами в скобках
echo -e "\n${YELLOW}5. Очистка файлов с (2) в названии...${NC}" | tee -a "$LOG_FILE"
find "$PROJECT_ROOT" -name "*(*).md" -type f ! -path "*/node_modules/*" | while read file; do
    safe_move "$file" "$ARCHIVE_DIR/numbered-duplicates/$(basename "$file")"
done

# 6. Организация Post Express документации
echo -e "\n${YELLOW}6. Консолидация Post Express документации...${NC}" | tee -a "$LOG_FILE"
mkdir -p "$PROJECT_ROOT/docs/features/logistics/post-express"
mkdir -p "$ARCHIVE_DIR/post-express-old"

# Перемещаем старые версии в архив
for file in $PROJECT_ROOT/docs/POST_EXPRESS_*.md; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        # Оставляем только главные файлы
        if [[ "$filename" == "POST_EXPRESS_INTEGRATION_COMPLETE.md" ]]; then
            safe_move "$file" "$PROJECT_ROOT/docs/features/logistics/post-express/README.md"
        else
            safe_move "$file" "$ARCHIVE_DIR/post-express-old/$filename"
        fi
    fi
done

# 7. Создание структуры директорий
echo -e "\n${YELLOW}7. Создание новой структуры директорий...${NC}" | tee -a "$LOG_FILE"
mkdir -p "$PROJECT_ROOT/docs/"{architecture/{backend,frontend,database,infrastructure},guides/{development,deployment,testing,maintenance},features/{marketplace,storefronts,payments,logistics,translations,categories},api,business/{investor-materials,plans,reports},ui-ux/{current,archive}}

# 8. Статистика
echo -e "\n${GREEN}========================================${NC}" | tee -a "$LOG_FILE"
echo -e "${GREEN}СТАТИСТИКА ОЧИСТКИ:${NC}" | tee -a "$LOG_FILE"
echo -e "${GREEN}========================================${NC}" | tee -a "$LOG_FILE"

# Подсчет файлов
TOTAL_BEFORE=$(find "$PROJECT_ROOT" -name "*.md" -type f ! -path "*/node_modules/*" | wc -l)
ARCHIVED=$(find "$ARCHIVE_DIR" -name "*.md" -type f 2>/dev/null | wc -l || echo 0)

echo "MD файлов до очистки: $TOTAL_BEFORE" | tee -a "$LOG_FILE"
echo "Файлов архивировано: $ARCHIVED" | tee -a "$LOG_FILE"
echo "Архив создан в: $ARCHIVE_DIR" | tee -a "$LOG_FILE"
echo "Лог сохранен в: $LOG_FILE" | tee -a "$LOG_FILE"

echo -e "\n${GREEN}✓ Очистка завершена!${NC}" | tee -a "$LOG_FILE"
echo -e "${YELLOW}⚠ Рекомендуется проверить архив перед окончательным удалением${NC}" | tee -a "$LOG_FILE"

# 9. Создание отчета о дубликатах
echo -e "\n${YELLOW}Создание детального отчета о дубликатах...${NC}"
DUPLICATES_REPORT="$PROJECT_ROOT/docs/DUPLICATES_REPORT_$(date +%Y%m%d).md"

cat > "$DUPLICATES_REPORT" << EOF
# Отчет об обнаруженных дубликатах MD файлов
Дата: $(date)

## Статистика
- Всего MD файлов проверено: $TOTAL_BEFORE
- Файлов архивировано: $ARCHIVED
- Местоположение архива: $ARCHIVE_DIR

## Выполненные действия

### 1. Удаление директории /task/
Полностью дублировала файлы из корневой директории

### 2. Архивация дневных отчетов
Перемещены отчеты DAY_01 - DAY_29 в архив

### 3. Консолидация Post Express
Объединена разрозненная документация

### 4. Очистка дубликатов
Удалены файлы с идентичным содержимым

## Рекомендации

1. Проверить архив: \`$ARCHIVE_DIR\`
2. Если все в порядке, удалить архив через 30 дней
3. Обновить ссылки в CLAUDE.md на новые пути
4. Запустить проверку битых ссылок

EOF

echo -e "${GREEN}✓ Отчет создан: $DUPLICATES_REPORT${NC}"