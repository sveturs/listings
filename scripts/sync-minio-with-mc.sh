#!/bin/bash

# Скрипт синхронизации MinIO с использованием MinIO Client (mc)
# Более элегантное решение через S3 API

set -e

echo "🔄 Синхронизация MinIO через MinIO Client"

# Установка MinIO Client если не установлен
if ! command -v mc &> /dev/null; then
    echo "📥 Устанавливаем MinIO Client..."
    wget https://dl.min.io/client/mc/release/linux-amd64/mc -O /tmp/mc
    chmod +x /tmp/mc
    sudo mv /tmp/mc /usr/local/bin/
fi

echo "🔧 Настройка подключений к MinIO..."

# Настройка локального MinIO
mc alias set local http://localhost:9000 minioadmin zhmEsJZZNFN0vrCO7Hya

# Настройка dev MinIO (через Tailscale VPN)
mc alias set dev http://100.88.44.15:9002 minioadmin minioadmin

echo "📋 Проверка подключений..."
mc admin info local
mc admin info dev

echo "🔄 Синхронизация buckets..."

# Синхронизация listings
echo "📁 Синхронизация listings..."
mc mirror local/listings dev/listings --overwrite

# Синхронизация chat-files
echo "💬 Синхронизация chat-files..."
mc mirror local/chat-files dev/chat-files --overwrite

# Синхронизация review-photos
echo "📸 Синхронизация review-photos..."
mc mirror local/review-photos dev/review-photos --overwrite

# Синхронизация storefronts (если есть)
if mc ls local/storefronts >/dev/null 2>&1; then
    echo "🏪 Синхронизация storefronts..."
    mc mirror local/storefronts dev/storefronts --overwrite
fi

# Синхронизация products (если есть)
if mc ls local/products >/dev/null 2>&1; then
    echo "🛍️ Синхронизация products..."
    mc mirror local/products dev/products --overwrite
fi

echo "📊 Статистика после синхронизации:"
echo "=== Локальное хранилище ==="
mc du local --depth 2

echo "=== Dev хранилище ==="
mc du dev --depth 2

echo "✅ Синхронизация завершена!"
echo "🔗 Проверить можно по адресу: https://devs3.svetu.rs"