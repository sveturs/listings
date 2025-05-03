#!/bin/bash

# Скрипт для запуска загрузки mail-webui в фоновом режиме с nohup
echo "Запуск загрузки mail-webui в фоновом режиме..."
cd /data/hostel-booking-system/harbor-scripts
nohup ./upload_mail_webui.sh > mail_webui_upload.log 2>&1 &
echo "Процесс запущен в фоне. Для проверки статуса используйте:"
echo "tail -f /data/hostel-booking-system/harbor-scripts/mail_webui_upload.log"