#!/bin/bash
set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Функция для логирования с разным уровнем важности
log_info() {
    echo -e "${YELLOW}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${RED}[WARNING]${NC} $1"
}

log_debug() {
    if [ "${DEBUG:-false}" = "true" ]; then
        echo -e "${CYAN}[DEBUG]${NC} $1"
    fi
}

log_warn "Не удалось получить IP нового контейнера! Используем имя контейнера вместо IP."
log_error "Ошибка: Новый контейнер не отвечает на запросы API!"
log_success "MINIO_ROOT_USER=$MINIO_ROOT_USER"
log_info "POSTGRES_USER=$POSTGRES_USER"
