#!/bin/bash

# Ищем процесс на порту 3000
PID=$(netstat -tlnp 2>/dev/null | grep :3000 | grep -oP '\d+(?=/\w+)')

if [ -n "$PID" ]; then
    echo "Найден процесс $PID на порту 3000"
    kill -9 $PID
    if [ $? -eq 0 ]; then
        echo "Процесс $PID успешно остановлен"
        notify-send "Kill Port 3000" "Процесс $PID на порту 3000 успешно остановлен"
    else
        echo "Ошибка при остановке процесса $PID"
        notify-send "Kill Port 3000" "Ошибка при остановке процесса $PID"
    fi
else
    echo "Процесс на порту 3000 не найден"
    notify-send "Kill Port 3000" "Процесс на порту 3000 не найден"
fi