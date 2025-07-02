#!/bin/bash
# Скрипт для запуска Next.js dev сервера с обходом проблемы сетевых интерфейсов

cd /data/hostel-booking-system/frontend/svetu

# Очистка старых процессов
/home/dim/.local/bin/kill-port-3001.sh

# Установка переменных окружения для обхода проблемы с сетевыми интерфейсами
export NODE_ENV=development
export HOSTNAME=localhost
export HOST=0.0.0.0

# Запуск с явным указанием хоста
exec yarn dev --hostname 0.0.0.0 -p 3001