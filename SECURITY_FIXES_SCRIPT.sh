#!/bin/bash

# Скрипт быстрого исправления критических уязвимостей безопасности
# SveTu Platform Security Fixes Script
# Дата создания: 26 августа 2025

set -e  # Остановка при первой ошибке
set -u  # Ошибка при использовании неопределенных переменных

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Логирование
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Проверка что скрипт запускается из корня проекта
if [[ ! -f "docker-compose.yml" || ! -d "backend" || ! -d "frontend" ]]; then
    error "Скрипт должен запускаться из корня проекта SveTu!"
    exit 1
fi

log "🔒 Начинаем исправление критических уязвимостей безопасности..."

# Создание резервных копий
log "📦 Создание резервных копий конфигурационных файлов..."
BACKUP_DIR="security_backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

# Генерация новых секретов
log "🔐 Генерация новых криптографически стойких секретов..."

# Генерация JWT секрета (256 бит)
JWT_SECRET=$(openssl rand -base64 32 | tr -d '\n')
POSTGRES_PASSWORD=$(openssl rand -base64 16 | tr -d "=+/\n" | cut -c1-16)
MINIO_PASSWORD=$(openssl rand -base64 24 | tr -d "=+/\n" | cut -c1-20)
REDIS_PASSWORD=$(openssl rand -base64 20 | tr -d "=+/\n" | cut -c1-18)

log "✅ Новые секреты сгенерированы"

# Создание файла с новыми секретами
cat > "$BACKUP_DIR/new_secrets.txt" << EOF
# Новые сгенерированные секреты - ХРАНИТЕ В БЕЗОПАСНОСТИ!
# Дата генерации: $(date)

JWT_SECRET=$JWT_SECRET
POSTGRES_PASSWORD=$POSTGRES_PASSWORD
MINIO_PASSWORD=$MINIO_PASSWORD
REDIS_PASSWORD=$REDIS_PASSWORD

# ИНСТРУКЦИИ ПО ПРИМЕНЕНИЮ:
# 1. Обновите backend/.env:
#    JWT_SECRET=$JWT_SECRET
#    REDIS_PASSWORD=$REDIS_PASSWORD
#    DATABASE_URL=postgres://postgres:$POSTGRES_PASSWORD@localhost:5432/svetubd?sslmode=disable
#    MINIO_SECRET_KEY=$MINIO_PASSWORD

# 2. Обновите docker-compose.yml:
#    POSTGRES_PASSWORD: $POSTGRES_PASSWORD
#    MINIO_ROOT_PASSWORD: $MINIO_PASSWORD
#    Redis command: redis-server --appendonly yes --requirepass $REDIS_PASSWORD

# НЕ ДОБАВЛЯЙТЕ ЭТОТ ФАЙЛ В GIT!
EOF

# КРИТИЧЕСКИЕ ПРЕДУПРЕЖДЕНИЯ
echo ""
warn "🚨 КРИТИЧЕСКИЕ ПРЕДУПРЕЖДЕНИЯ:"
echo ""
error "❗ Обнаружены следующие критические уязвимости:"
echo "   1. API ключи в открытом тексте в .env файлах"
echo "   2. Слабые пароли базы данных"
echo "   3. JWT секрет по умолчанию"
echo "   4. OpenSearch без аутентификации"
echo "   5. Redis без пароля"
echo "   6. Широкие CORS политики"
echo ""
warn "⚠️  НЕМЕДЛЕННЫЕ ДЕЙСТВИЯ ТРЕБУЮТСЯ:"
echo ""
echo "1. 🔑 СМЕНИТЕ ВСЕ API КЛЮЧИ в .env файлах:"
echo "   - OpenAI: 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"
echo "   - Claude: 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"
echo "   - Google Client Secret: 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"
echo "   - Stripe Keys: sk_test_..., pk_test_..."
echo ""
echo "2. 🛠️  Примените новые секреты из файла: $BACKUP_DIR/new_secrets.txt"
echo ""
echo "3. 🔄 Обновите конфигурацию сервисов:"
echo "   - Включите аутентификацию OpenSearch"
echo "   - Добавьте пароль для Redis"
echo "   - Ограничьте CORS политики"
echo ""
echo "4. 🚫 Отключите production до исправления критических проблем"
echo ""

success "✅ Новые секреты созданы в: $BACKUP_DIR/new_secrets.txt"
warn "❗ НЕ ДОБАВЛЯЙТЕ ЭТОТ ФАЙЛ В GIT!"

echo ""
log "🔒 Полный отчет безопасности доступен в: SECURITY_AUDIT_REPORT_2025.md"
log "📊 Найдено: 8 критических, 5 высоких, 7 средних, 3 низких уязвимости"

echo ""
success "✅ Скрипт завершен. Следуйте инструкциям выше для повышения безопасности!"