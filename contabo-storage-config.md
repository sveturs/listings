# Конфигурация Contabo Object Storage с MinIO

## 📊 Архитектура гибридного хранилища

```
┌─────────────────────────────────────────────────────────┐
│                    ПОЛЬЗОВАТЕЛИ                          │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
        ┌──────────────────────────┐
        │   CloudFlare CDN (Free)  │ ← Кэширование изображений
        └────────────┬─────────────┘
                     │
                     ▼
        ┌──────────────────────────┐
        │     Nginx + SSL          │ ← https://s3.svetu.rs
        └────────────┬─────────────┘
                     │
        ┌────────────┴─────────────┐
        │                          │
        ▼                          ▼
┌──────────────────┐    ┌──────────────────┐
│   MinIO Local    │    │  MinIO Gateway   │
│   (Hot Storage)  │    │  (Cold Storage)  │
│     1TB SSD      │    │  Contabo 500GB   │
└──────────────────┘    └──────────────────┘
        │                          │
        │                          │
        ▼                          ▼
  Новые файлы,            Архивы, бэкапы,
  частый доступ           старые данные
```

## 🔧 Получение доступов Contabo Object Storage

1. Войдите в панель управления Contabo
2. Перейдите в **Storage → Object Storage**
3. Нажмите на **Object Storage European Union 9967**
4. Найдите раздел **S3 Credentials**
5. Создайте новые Access Keys:
   - Нажмите **Create Access Key**
   - Сохраните **Access Key ID** и **Secret Access Key**

## 🚀 Быстрая установка

```bash
# Скачайте и запустите скрипт установки
wget https://raw.githubusercontent.com/your-repo/minio-hybrid-setup.sh
chmod +x minio-hybrid-setup.sh
sudo ./minio-hybrid-setup.sh
```

## 📝 Ручная конфигурация

### 1. Настройка переменных окружения

```bash
# /opt/minio/.env
CONTABO_ACCESS_KEY=your_access_key_here
CONTABO_SECRET_KEY=your_secret_key_here
CONTABO_ENDPOINT=eu2.contabostorage.com
CONTABO_REGION=EU
```

### 2. Стратегия хранения данных

| Тип данных | Локальное хранилище | Contabo Storage | Lifecycle |
|------------|-------------------|-----------------|-----------|
| Новые фото объявлений | ✅ Первые 30 дней | ➡️ После 30 дней | Auto-transition |
| Витрины магазинов | ✅ Всегда | 📋 Backup daily | Mirror |
| Чат файлы | ✅ 7 дней | ❌ | Auto-delete |
| Dev данные | ✅ 30 дней | ❌ | Auto-delete |
| Бэкапы | ❌ | ✅ Всегда | Direct upload |
| Архивы | ❌ | ✅ Всегда | Direct upload |

### 3. Команды MinIO Client для работы с Contabo

```bash
# Настройка алиасов
mc alias set local https://s3.svetu.rs svetu_admin_s3 BLcLlznxtWzb6j5vdRUumFA1t
mc alias set contabo https://eu2.contabostorage.com 39e2e4987c6c4c9926c9b24bca119cd0 a479932e4af2c29b16049223b3e54d42

# Создание бакетов на Contabo
mc mb contabo/svetu-production-archive
mc mb contabo/svetu-backups
mc mb contabo/svetu-cold-storage

# Копирование данных на Contabo
mc cp local/production-listings/old-data/* contabo/svetu-production-archive/

# Зеркалирование для бэкапа
mc mirror local/production contabo/svetu-backups/$(date +%Y%m%d)/

# Проверка использования
mc du contabo/svetu-backups
```

## 🔄 Автоматизация lifecycle

### Политика для production данных:
```json
{
  "Rules": [{
    "ID": "MoveOldToContabo",
    "Status": "Enabled",
    "Prefix": "",
    "Transitions": [{
      "Days": 30,
      "StorageClass": "GLACIER"
    }],
    "NoncurrentVersionTransitions": [{
      "NoncurrentDays": 7,
      "StorageClass": "GLACIER"
    }]
  }]
}
```

### Применение политики:
```bash
mc ilm import local/production-listings < lifecycle-policy.json
```

## 📊 Мониторинг использования

```bash
# Проверка локального хранилища
df -h /opt/minio/data
mc admin info local

# Проверка Contabo Storage
mc du --depth 1 contabo/
mc stat contabo/svetu-backups

# Общая статистика
docker exec minio-hybrid mc admin info local
```

## 🔐 Безопасность

1. **Разделение доступов:**
   - Production приложение → только локальный MinIO
   - Backup scripts → доступ к Contabo
   - Dev окружение → изолированные бакеты

2. **Шифрование:**
   - TLS для всех соединений
   - Опционально: шифрование на Contabo стороне

3. **Ротация ключей:**
   - Каждые 90 дней для Contabo
   - Каждые 30 дней для приложений

## 💰 Оптимизация стоимости

### Текущие расходы:
- VPS Storage (1TB): ~€30/месяц
- Contabo Object Storage (500GB): €2.49/месяц
- **Итого: ~€32.49/месяц**

### При росте до 5TB:
- Опция 1: Добавить Contabo Storage (€2.49/500GB)
  - 5TB = €24.90/месяц
- Опция 2: Cloudflare R2
  - 5TB = $75/месяц (хранение) + $0 egress
- Опция 3: Больший VPS
  - Storage VPS 4TB = €55/месяц

**Рекомендация:** Гибрид Contabo Object Storage + Cloudflare R2

## 🚨 Troubleshooting

### Проблема: Не подключается к Contabo
```bash
# Проверка подключения
mc ls contabo --debug

# Проверка credentials
curl -I https://eu2.contabostorage.com
```

### Проблема: Медленная загрузка
```bash
# Включить кэширование на MinIO
docker exec minio-hybrid mc admin config set local cache drives="/cache"
docker exec minio-hybrid mc admin service restart local
```

### Проблема: Нехватка места
```bash
# Очистка кэша
docker exec minio-hybrid rm -rf /cache/*

# Принудительный запуск lifecycle
mc ilm rule apply local/production-listings
```

## 📈 План масштабирования

### Фаза 1 (сейчас - 3 месяца):
- ✅ 1TB local + 500GB Contabo
- Ручное управление lifecycle
- Базовый мониторинг

### Фаза 2 (3-6 месяцев):
- Увеличение Contabo до 2TB
- Автоматизация tiering
- CloudFlare CDN integration

### Фаза 3 (6-12 месяцев):
- Migration to distributed MinIO
- Multi-region replication
- Advanced caching strategies

## 📞 Поддержка

- **Contabo Support:** https://contabo.com/support
- **MinIO Docs:** https://min.io/docs/
- **Наш Telegram:** @svetu_devops