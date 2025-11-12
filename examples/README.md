# Listings gRPC Client Examples

## ReindexAll Client

Простой Go клиент для тестирования метода `ReindexAll` gRPC.

### Компиляция

```bash
# Скомпилировать клиент
go build -o reindex_client ./examples/reindex_client.go
```

### Использование

```bash
# Переиндексировать ВСЕ продукты (B2C + C2C), batch size = 1000
./reindex_client

# Переиндексировать только B2C продукты
./reindex_client b2c

# Переиндексировать только C2C листинги
./reindex_client c2c

# Переиндексировать B2C с custom batch size = 500
./reindex_client b2c 500

# Использовать другой gRPC сервер
GRPC_HOST=listings-service:50051 ./reindex_client b2c
```

### Пример вывода

```
🔌 Connecting to gRPC server at localhost:50051
📦 Source Type: b2c
📦 Batch Size: 1000

🚀 Starting reindexing...

✅ Reindexing completed successfully!
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 Total Indexed:    5432 products
❌ Total Failed:     0 products
⏱️  Duration:         45 seconds (0.75 minutes)
🕐 Client Elapsed:   46.23 seconds
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

📈 Success Rate: 100.00%

✨ Done!
```

### Troubleshooting

#### Connection refused
```bash
# Проверить что gRPC сервер запущен
netstat -tlnp | grep 50051

# Или используй правильный адрес
GRPC_HOST=localhost:50051 ./reindex_client
```

#### Timeout errors
```bash
# Увеличь timeout в коде (по умолчанию 10 минут)
# Или используй меньший batch size
./reindex_client b2c 500
```

## Shell Script для Batch Testing

Используй `test_reindex.sh` для автоматического тестирования:

```bash
# Запустить все тесты
cd /p/github.com/sveturs/listings
./test_reindex.sh

# Или с custom gRPC host
GRPC_HOST=listings-service:50051 ./test_reindex.sh
```
