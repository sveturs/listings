#!/bin/bash

LOG_FILE="test_log.txt"

# Очистка файла логов
> $LOG_FILE

# Функция для логирования команд и их результатов
log_command() {
    echo "Executing: $@" >> $LOG_FILE
    eval "$@" >> $LOG_FILE 2>&1
    echo "---------------------------------------" >> $LOG_FILE
}

# Создание пользователей
log_command 'curl -X POST http://localhost:3000/users -H "Content-Type: application/json" -d '\''{"name":"Иван", "email":"ivan@example.com"}'\'''
log_command 'curl -X POST http://localhost:3000/users -H "Content-Type: application/json" -d '\''{"name":"Анна", "email":"anna@example.com"}'\'''

# Создание комнат
log_command 'curl -X POST http://localhost:3000/rooms -H "Content-Type: application/json" -d '\''{"name":"Luxury Suite", "capacity":2, "price_per_night":150}'\'''
log_command 'curl -X POST http://localhost:3000/rooms -H "Content-Type: application/json" -d '\''{"name":"Economy Room", "capacity":4, "price_per_night":100}'\'''
log_command 'curl -X POST http://localhost:3000/rooms -H "Content-Type: application/json" -d '\''{"name":"Standard Room", "capacity":3, "price_per_night":120}'\'''

# Создание бронирований
log_command 'curl -X POST http://localhost:3000/bookings -H "Content-Type: application/json" -d '\''{"user_id":1, "room_id":1, "start_date":"2024-11-20", "end_date":"2024-11-25"}'\'''
log_command 'curl -X POST http://localhost:3000/bookings -H "Content-Type: application/json" -d '\''{"user_id":2, "room_id":2, "start_date":"2024-11-22", "end_date":"2024-11-27"}'\'''

# Проверка фильтров
log_command 'curl "http://localhost:3000/rooms?capacity=2"'
log_command 'curl "http://localhost:3000/rooms?min_price=100&max_price=150"'
log_command 'curl "http://localhost:3000/rooms?start_date=2024-11-21&end_date=2024-11-23"'
log_command 'curl "http://localhost:3000/rooms?capacity=4&min_price=100&max_price=200&start_date=2024-11-21&end_date=2024-11-23"'

# Удаление бронирований
log_command 'curl -X DELETE http://localhost:3000/bookings/1'
log_command 'curl -X DELETE http://localhost:3000/bookings/2'

# Удаление комнат
log_command 'curl -X DELETE http://localhost:3000/rooms/1'
log_command 'curl -X DELETE http://localhost:3000/rooms/2'
log_command 'curl -X DELETE http://localhost:3000/rooms/3'

# Финальная проверка списка комнат и бронирований
log_command 'curl http://localhost:3000/rooms'
log_command 'curl http://localhost:3000/bookings'

echo "Тест завершён. Логи записаны в $LOG_FILE."

