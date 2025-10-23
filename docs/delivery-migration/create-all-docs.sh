#!/bin/bash

SRC="/data/hostel-booking-system/docs/DELIVERY_MICROSERVICE_MIGRATION_CLEAN_CUT.md"
DEST_DIR="/data/hostel-booking-system/docs/delivery-migration"

# 01-OVERVIEW.md (строки 1-169)
sed -n '1,169p' "$SRC" > "$DEST_DIR/01-OVERVIEW.md"

# 02-ARCHITECTURE.md (строки 50-302)
echo "# Архитектура: Текущая vs Целевая" > "$DEST_DIR/02-ARCHITECTURE.md"
sed -n '50,302p' "$SRC" >> "$DEST_DIR/02-ARCHITECTURE.md"

# 03-TECHNICAL-SPEC.md (строки 303-433)
echo "# Техническая спецификация микросервиса" > "$DEST_DIR/03-TECHNICAL-SPEC.md"
sed -n '303,433p' "$SRC" >> "$DEST_DIR/03-TECHNICAL-SPEC.md"

# 04-DATABASE.md (строки 210-266)
echo "# База данных: Схема и миграция" > "$DEST_DIR/04-DATABASE.md"
sed -n '210,266p' "$SRC" >> "$DEST_DIR/04-DATABASE.md"

# 05-PHASE-0 (строки 3271-3278 из обновленного чеклиста)
echo "# Фаза 0: Подготовка инфраструктуры" > "$DEST_DIR/05-PHASE-0-INFRASTRUCTURE.md"
echo "" >> "$DEST_DIR/05-PHASE-0-INFRASTRUCTURE.md"
echo "**Срок**: Week 0" >> "$DEST_DIR/05-PHASE-0-INFRASTRUCTURE.md"
echo "" >> "$DEST_DIR/05-PHASE-0-INFRASTRUCTURE.md"
sed -n '3271,3278p' "$SRC" >> "$DEST_DIR/05-PHASE-0-INFRASTRUCTURE.md"

# 06-PHASE-1 (строки 337-1300 примерно - план миграции фаза 1)
echo "# Фаза 1: Разработка микросервиса" > "$DEST_DIR/06-PHASE-1-DEVELOPMENT.md"
echo "" >> "$DEST_DIR/06-PHASE-1-DEVELOPMENT.md"
echo "**Срок**: Week 1-2" >> "$DEST_DIR/06-PHASE-1-DEVELOPMENT.md"
echo "" >> "$DEST_DIR/06-PHASE-1-DEVELOPMENT.md"
sed -n '337,1300p' "$SRC" >> "$DEST_DIR/06-PHASE-1-DEVELOPMENT.md"

# 07-PHASE-2 (строки 1301-2100 примерно - тестирование)
echo "# Фаза 2: Тестирование" > "$DEST_DIR/07-PHASE-2-TESTING.md"
echo "" >> "$DEST_DIR/07-PHASE-2-TESTING.md"
echo "**Срок**: Week 3" >> "$DEST_DIR/07-PHASE-2-TESTING.md"
echo "" >> "$DEST_DIR/07-PHASE-2-TESTING.md"
sed -n '1301,2100p' "$SRC" >> "$DEST_DIR/07-PHASE-2-TESTING.md"

# 08-PHASE-3 (строки 2101-2433 - deploy)
echo "# Фаза 3: Развертывание и миграция монолита" > "$DEST_DIR/08-PHASE-3-DEPLOYMENT.md"
echo "" >> "$DEST_DIR/08-PHASE-3-DEPLOYMENT.md"
echo "**Срок**: Week 4" >> "$DEST_DIR/08-PHASE-3-DEPLOYMENT.md"
echo "" >> "$DEST_DIR/08-PHASE-3-DEPLOYMENT.md"
sed -n '2101,2433p' "$SRC" >> "$DEST_DIR/08-PHASE-3-DEPLOYMENT.md"

# 10-DEPLOYMENT-GUIDE (строки 3011-3138 - инструкции развертывания)
echo "# Пошаговая инструкция развертывания" > "$DEST_DIR/10-DEPLOYMENT-GUIDE.md"
sed -n '3011,3138p' "$SRC" >> "$DEST_DIR/10-DEPLOYMENT-GUIDE.md"

# 11-MONITORING-ROLLBACK (строки 2547-2614)
echo "# Мониторинг и Rollback" > "$DEST_DIR/11-MONITORING-ROLLBACK.md"
sed -n '2547,2614p' "$SRC" >> "$DEST_DIR/11-MONITORING-ROLLBACK.md"

# 12-CHECKLISTS (строки 2617-2658 + 3269-3332)
echo "# Чеклисты выполнения" > "$DEST_DIR/12-CHECKLISTS.md"
sed -n '2617,2658p' "$SRC" >> "$DEST_DIR/12-CHECKLISTS.md"
echo "" >> "$DEST_DIR/12-CHECKLISTS.md"
echo "---" >> "$DEST_DIR/12-CHECKLISTS.md"
echo "" >> "$DEST_DIR/12-CHECKLISTS.md"
sed -n '3269,3332p' "$SRC" >> "$DEST_DIR/12-CHECKLISTS.md"

# 13-TROUBLESHOOTING (строки 3140-3265)
echo "# Troubleshooting" > "$DEST_DIR/13-TROUBLESHOOTING.md"
sed -n '3140,3265p' "$SRC" >> "$DEST_DIR/13-TROUBLESHOOTING.md"

echo "All files created successfully!"
ls -lh "$DEST_DIR"/*.md | wc -l
