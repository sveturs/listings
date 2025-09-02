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
# MinIO Console (✅ ЗАЩИЩЕНО HTTPS с заголовками безопасности)
URL: https://console.s3.svetu.rs (✅ РАБОТАЕТ)
🔒 Безопасность: HSTS, X-Frame-Options, X-Content-Type-Options
Login: Хранится в /opt/minio/secrets/minio_root_user.txt
Password: Хранится в /opt/minio/secrets/minio_root_password.txt

# Просмотр текущих учетных данных:
cat /opt/minio/secrets/credentials.info

# MinIO API Endpoint (✅ ЗАЩИЩЕНО HTTPS)
URL: https://s3.svetu.rs (✅ РАБОТАЕТ)
❌ ПРЯМОЙ ДОСТУП БЛОКИРОВАН: ~~http://194.163.132.116:9000~~ (заблокирован)
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

# Проверка health (ОБНОВЛЕНО для безопасности)
# ✅ ЛОКАЛЬНО (работает):
curl -I http://localhost:9000/minio/health/live

# ✅ ЧЕРЕЗ REVERSE PROXY (рекомендуемо):
curl -I https://s3.svetu.rs/

# ❌ ПРЯМОЙ ДОСТУП БЛОКИРОВАН:
# curl -I http://194.163.132.116:9000/minio/health/live
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

### ✅ ОБНОВЛЕННЫЕ НАСТРОЙКИ БЕЗОПАСНОСТИ (2025-09-01):
- **✅ Network Security:** MinIO порты 9000/9001 заблокированы для внешнего доступа
- **✅ Private Buckets:** Все production bucket'ы установлены в режим PRIVATE
- **✅ HTTPS Only:** Доступ только через nginx reverse proxy
- **✅ Security Headers:** HSTS, X-Frame-Options, X-Content-Type-Options, Referrer-Policy
- **✅ Docker Secrets:** Учетные данные защищены через Docker secrets
- **✅ UFW Firewall:** DENY правила для MinIO портов + ALLOW только для 80/443/22
- **✅ SSL/TLS:** Let's Encrypt сертификаты для всех endpoints
- **✅ Версионирование:** Включено для production bucket'ов
- **✅ Ротация логов:** Настроена через logrotate

### Статус систем безопасности:
```bash
# Проверка Fail2ban
sudo systemctl status fail2ban
sudo fail2ban-client status

# Проверка UFW
sudo ufw status verbose

# Проверка SSL сертификатов
sudo certbot certificates

# Проверка Docker secrets
docker secret ls
```

### Регулярные задачи:
```bash
# Аудит доступов (еженедельно)
mc admin user list local
mc admin policy list local

# Проверка приватности bucket'ов (ОБНОВЛЕНО)
mc anonymous list local --recursive  # Должно показать "private" для всех production bucket'ов

# Ротация паролей (каждые 90 дней) - теперь безопасно
openssl rand -base64 32 > /opt/minio/secrets/minio_root_password.txt
cd /opt/minio && docker-compose restart

# Проверка логов доступа
docker logs minio-hybrid | grep -i "error\|warn\|fail"

# Мониторинг попыток взлома
sudo fail2ban-client status sshd
sudo fail2ban-client status minio
```

### Управление Fail2ban:
```bash
# Статус всех jail'ов
sudo fail2ban-client status

# Разблокировка IP
sudo fail2ban-client set sshd unbanip <IP_ADDRESS>

# Просмотр заблокированных IP
sudo fail2ban-client get sshd banned
```

### 🔐 Проверка новых мер безопасности:
```bash
# Проверка блокировки прямого доступа к MinIO портам
curl -m 5 http://194.163.132.116:9000 || echo "✅ Порт 9000 заблокирован"
curl -m 5 http://194.163.132.116:9001 || echo "✅ Порт 9001 заблокирован"

# Проверка работы через reverse proxy
curl -s -w "%{http_code}" https://s3.svetu.rs/ | head -1  # Должно быть 403
curl -s -w "%{http_code}" https://console.s3.svetu.rs/ | head -1  # Должно быть 200

# Проверка приватности production bucket'ов
curl -s -w "%{http_code}" https://s3.svetu.rs/production-listings/test.txt | head -1  # Должно быть 403

# Проверка security headers
curl -I https://s3.svetu.rs/ | grep -E "(Strict-Transport|X-Content)"
curl -I https://console.s3.svetu.rs/ | grep -E "(X-Frame|Strict-Transport)"

# Проверка UFW правил
sudo ufw status | grep -E "(9000|9001)"  # Должно показать DENY
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
/opt/minio/                    # Корневая директория MinIO
├── docker-compose.yml         # Текущая конфигурация Docker
├── docker-compose-secure.yml  # Безопасная конфигурация
├── .env                       # Переменные окружения (очищены)
├── secrets/                   # Защищенные учетные данные
│   ├── minio_root_user.txt    # Имя администратора
│   ├── minio_root_password.txt # Пароль администратора
│   └── credentials.info       # Сводка учетных данных
├── data/                      # Данные MinIO
├── cache/                     # Кэш для горячих данных
├── config/                    # Конфигурация MinIO
├── scripts/                   # Скрипты обслуживания
│   ├── backup.sh              # Автоматический бэкап
│   ├── secure_migration.sh    # Миграция на безопасную конфигурацию
│   ├── enable_versioning.sh   # Включение версионирования
│   └── setup_lifecycle.sh     # Настройка lifecycle политик
├── backup/                    # Локальные бэкапы
└── DISASTER_RECOVERY.md       # План аварийного восстановления

/etc/fail2ban/jail.local       # Конфигурация Fail2ban
/etc/logrotate.d/minio         # Настройки ротации логов
/var/log/minio-backup.log      # Лог файлы бэкапов
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