#!/bin/bash

# Скрипт для запуска загрузки mailserver в фоновом режиме с nohup
echo "Запуск загрузки mailserver в фоновом режиме..."
cd /data/hostel-booking-system/harbor-scripts
nohup ./upload_mailserver.sh > mailserver_upload.log 2>&1 &
echo "Процесс запущен в фоне. Для проверки статуса используйте:"
echo "tail -f /data/hostel-booking-system/harbor-scripts/mailserver_upload.log"