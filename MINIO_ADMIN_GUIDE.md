# 📚 РУКОВОДСТВО АДМИНИСТРАТОРА MinIO S3 Storage
## Централизованное хранилище для проекта SveTu.rs

---

## 🏗️ АРХИТЕКТУРА СИСТЕМЫ

### Характеристики инфраструктуры:
- **Сервер:** Contabo Storage VPS
  - CPU: 6 ядер
  - RAM: 18 GB  
  - SSD: 1 TB (локальное горячее хранилище)
  - Object Storage: 500 GB (Contabo, холодное хранилище)
  - **Общая емкость: 1.5 TB**
  - Сеть: 1 Gbps
  - Локация: Европа (низкая латентность для Балкан)

### Возможности системы:
- ✅ S3-совместимое API (100% совместимость с AWS S3)
- ✅ Гибридное хранилище (SSD + Object Storage)
- ✅ Автоматический tiering (горячие → холодные данные)
- ✅ Публичный доступ для изображений через CDN
- ✅ Многопользовательская система с изоляцией
- ✅ Версионирование объектов (опционально)
- ✅ Шифрование данных в транзите (TLS)
- ✅ Автоматические бэкапы на Contabo Storage
- ✅ Поддержка больших файлов (до 5TB на объект)
- ✅ Multipart upload для оптимизации скорости

---

## 📁 СТРУКТУРА ДАННЫХ

```
MinIO Local (1TB SSD) - Горячие данные
├── production-listings/      # Фото объявлений (публичный доступ)
├── production-chat-files/    # Файлы из чатов
├── production-storefronts/   # Витрины магазинов (публичный доступ)  
├── production-user-documents/ # Документы пользователей
├── development-listings/     # Dev фото (публичный доступ)
├── development-chat-files/   # Dev файлы чатов
└── development-storefronts/  # Dev витрины

Contabo Storage (500GB) - Холодные данные
├── svetu-backup/   # Ежедневные бэкапы
├── svetu-archive/  # Архив старых данных (>30 дней)
└── svetu-cold/     # Долгосрочное хранение
```

---

## 🔐 ДОСТУПЫ И CREDENTIALS

### Административный доступ:
```bash
# MinIO Console
URL: https://console.s3.svetu.rs (после добавления DNS)
Временный URL: http://194.163.132.116:9001
Login: svetu_admin_s3
Password: BLcLlznxtWzb6j5vdRUumFA1t

# MinIO API Endpoint
URL: https://s3.svetu.rs
Прямой доступ: http://194.163.132.116:9000
```

### Сервисные аккаунты:
```bash
# Production Backend
Access Key: production_backend
Secret Key: xK9mNjR3tP5wQ2aLbV7cH8dS

# Development Backend  
Access Key: dev_backend
Secret Key: pL4kJ8nM2qR6tY9wX5zC3vB7
```

### Contabo Storage:
```bash
Endpoint: https://eu2.contabostorage.com
Access Key: 39e2e4987c6c4c9926c9b24bca119cd0
Secret Key: a479932e4af2c29b16049223b3e54d42
```

---

## 🛠️ ЕЖЕДНЕВНОЕ ОБСЛУЖИВАНИЕ

### 1. Проверка статуса системы
```bash
# Статус контейнера
docker ps | grep minio

# Логи MinIO
docker logs -f minio-hybrid --tail 100

# Статус хранилища
mc admin info local

# Использование диска
df -h /opt/minio/data
mc du --depth 1 local/
mc du --depth 1 contabo/

# Проверка health
curl -I http://localhost:9000/minio/health/live
```

### 2. Мониторинг производительности
```bash
# Статистика по бакетам
mc stat local/production-listings

# Активные сессии
mc admin trace local

# Топ объектов по размеру
mc find local --larger 100MB --maxdepth 2

# Скорость сети
mc admin speedtest local
```

### 3. Управление пользователями
```bash
# Список пользователей
mc admin user list local

# Добавить пользователя
mc admin user add local NEW_USER NEW_PASSWORD

# Назначить политику
mc admin policy attach local readwrite --user NEW_USER

# Удалить пользователя
mc admin user remove local USER_NAME

# Сбросить пароль
mc admin user password local USER_NAME NEW_PASSWORD
```

---

## 💾 РЕЗЕРВНОЕ КОПИРОВАНИЕ

### Автоматические бэкапы (настроены в cron):
```bash
# Просмотр расписания
crontab -l

# Ручной запуск бэкапа
/opt/minio/scripts/backup.sh

# Проверка последних бэкапов
mc ls contabo/svetu-backup/ --recursive | tail -20
```

### Создание полного бэкапа:
```bash
# Бэкап критичных production данных
DATE=$(date +%Y%m%d_%H%M%S)
mc mirror local/production-listings contabo/svetu-backup/$DATE/listings
mc mirror local/production-user-documents contabo/svetu-backup/$DATE/documents
```

### Восстановление из бэкапа:
```bash
# Список доступных бэкапов
mc ls contabo/svetu-backup/

# Восстановление
mc mirror contabo/svetu-backup/20250901/listings local/production-listings-restored
```

---

## 🚨 РЕШЕНИЕ ПРОБЛЕМ

### Проблема: MinIO не запускается
```bash
# Проверка логов
docker logs minio-hybrid

# Перезапуск
cd /opt/minio
docker-compose restart

# Полная переустановка
docker-compose down
docker-compose up -d
```

### Проблема: Нехватка места
```bash
# Анализ использования
mc du --depth 2 local/ | sort -rh | head -20

# Очистка старых версий
mc rm --recursive --force --older-than 30d local/development-listings

# Принудительный запуск lifecycle policies
mc ilm list local/production-listings
```

### Проблема: Медленная загрузка
```bash
# Проверка кэша
ls -lah /opt/minio/cache/

# Очистка кэша
docker exec minio-hybrid rm -rf /cache/*

# Проверка сети
speedtest-cli
mc admin speedtest local
```

### Проблема: Не работает Contabo Storage
```bash
# Проверка подключения
mc ls contabo --debug

# Тест доступа
curl -I https://eu2.contabostorage.com

# Переконфигурация
mc alias set contabo https://eu2.contabostorage.com ACCESS_KEY SECRET_KEY
```

---

## 📊 ОПТИМИЗАЦИЯ ПРОИЗВОДИТЕЛЬНОСТИ

### 1. Настройка кэширования
```bash
# Проверка текущих настроек
docker exec minio-hybrid printenv | grep CACHE

# Изменение параметров кэша (в /opt/minio/.env)
MINIO_CACHE_QUOTA=90  # Использовать 90% диска для кэша
MINIO_CACHE_AFTER=0   # Кэшировать сразу после первого доступа
```

### 2. Настройка lifecycle policies
```bash
# Автоудаление старых dev данных
mc ilm add --expiry-days 7 local/development-chat-files

# Перемещение в архив
cat > lifecycle.json <<EOF
{
  "Rules": [{
    "ID": "MoveToArchive",
    "Status": "Enabled",
    "Transition": {
      "Days": 30,
      "StorageClass": "GLACIER"
    }
  }]
}
EOF
mc ilm import local/production-listings < lifecycle.json
```

### 3. Оптимизация изображений
```bash
# Найти большие изображения
mc find local --name "*.jpg" --larger 5MB

# Установить квоты на бакеты
mc admin bucket quota local/development-listings --hard 100GB
```

---

## 🔄 ОБНОВЛЕНИЕ СИСТЕМЫ

### Обновление MinIO:
```bash
cd /opt/minio

# Бэкап конфигурации
cp docker-compose.yml docker-compose.yml.backup
cp .env .env.backup

# Обновление образа
docker-compose pull
docker-compose down
docker-compose up -d

# Проверка версии
mc admin info local | grep Version
```

### Обновление MinIO Client:
```bash
# Скачать последнюю версию
sudo curl https://dl.min.io/client/mc/release/linux-amd64/mc \
  -o /usr/local/bin/mc.new

# Заменить
sudo mv /usr/local/bin/mc /usr/local/bin/mc.old
sudo mv /usr/local/bin/mc.new /usr/local/bin/mc
sudo chmod +x /usr/local/bin/mc

# Проверка
mc --version
```

---

## 📈 МАСШТАБИРОВАНИЕ

### Когда нужно масштабировать:
- Использование диска > 80%
- Latency > 500ms
- Количество объектов > 10 миллионов

### План масштабирования:
1. **Краткосрочный (1-3 месяца):**
   - Добавить Contabo Storage блоки (+500GB = €2.49/месяц)
   - Настроить агрессивное архивирование

2. **Среднесрочный (3-6 месяцев):**
   - Миграция на больший VPS (4TB)
   - Добавить второй MinIO узел

3. **Долгосрочный (6-12 месяцев):**
   - Distributed MinIO (4+ узла)
   - CDN интеграция (CloudFlare)
   - Георепликация

---

## 🔒 БЕЗОПАСНОСТЬ

### Регулярные задачи:
```bash
# Аудит доступов (еженедельно)
mc admin user list local
mc admin policy list local

# Проверка публичных политик
mc anonymous list local --recursive

# Ротация паролей (каждые 90 дней)
mc admin user password local svetu_admin_s3 NEW_PASSWORD

# Проверка логов доступа
docker logs minio-hybrid | grep -i "error\|warn\|fail"
```

### Настройка firewall:
```bash
# Текущие правила
sudo ufw status

# Ограничение доступа к консоли
sudo ufw allow from YOUR_IP to any port 9001
```

---

## 📞 КОНТАКТЫ И ПОДДЕРЖКА

### При критических проблемах:
1. Проверьте логи: `docker logs minio-hybrid`
2. Проверьте статус: `mc admin info local`
3. Перезапустите: `docker-compose restart`

### Полезные ресурсы:
- MinIO Docs: https://min.io/docs/
- Contabo Support: https://contabo.com/support
- S3 API Reference: https://docs.aws.amazon.com/s3/

### Расположение файлов:
```
/opt/minio/              # Корневая директория MinIO
├── docker-compose.yml   # Конфигурация Docker
├── .env                 # Переменные окружения и пароли
├── data/                # Данные MinIO
├── cache/               # Кэш для горячих данных
├── config/              # Конфигурация MinIO
├── scripts/             # Скрипты обслуживания
│   ├── backup.sh        # Скрипт бэкапа
│   └── monitor.sh       # Скрипт мониторинга
└── backup/              # Локальные бэкапы
```

---

## ⚡ БЫСТРЫЕ КОМАНДЫ

```bash
# Рестарт MinIO
docker-compose -f /opt/minio/docker-compose.yml restart

# Статус системы одной командой
echo "=== Docker ===" && docker ps | grep minio && \
echo "=== Storage ===" && mc admin info local && \
echo "=== Disk ===" && df -h /opt/minio/data

# Быстрый бэкап
mc mirror local/production contabo/svetu-backup/quick-$(date +%Y%m%d)

# Очистка кэша
docker exec minio-hybrid sh -c "rm -rf /cache/*"

# Топ-10 больших файлов
mc find local --larger 10MB | head -10
```

---

*Последнее обновление: Сентябрь 2025*
*Версия системы: MinIO RELEASE.2025-07-23*