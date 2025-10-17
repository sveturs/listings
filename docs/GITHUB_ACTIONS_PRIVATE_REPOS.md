# GitHub Actions: Настройка доступа к приватным репозиториям

## Проблема

При использовании приватных Go модулей (например, `github.com/sveturs/auth`) в GitHub Actions возникает ошибка:

```
fatal: could not read Username for 'https://github.com': terminal prompts disabled
```

Это происходит потому, что Go пытается скачать приватный модуль, но не имеет прав доступа.

## Решение

Для доступа к приватным репозиториям нужно настроить Git authentication, GOPRIVATE переменную и permissions.

### 1. Добавить permissions в workflow

**КРИТИЧЕСКИ ВАЖНО:** Workflow должен иметь явные permissions для доступа к другим репозиториям:

```yaml
permissions:
  contents: read
  packages: read
```

Это дает GITHUB_TOKEN права на чтение других репозиториев в организации.

### 2. Использование встроенного GITHUB_TOKEN

GitHub Actions автоматически предоставляет `GITHUB_TOKEN` с доступом к репозиториям организации. Этот токен НЕ нужно создавать вручную, но нужно настроить permissions (см. выше).

### 3. Конфигурация в workflow файлах

**ВАЖНО:** Git конфигурация должна быть **ДО** `actions/setup-go`, потому что setup-go с `cache: true` пытается восстановить кэш зависимостей сразу при выполнении.

Полный пример workflow файла:

```yaml
name: My Workflow

on:
  pull_request:

permissions:
  contents: read
  packages: read  # КРИТИЧЕСКИ ВАЖНО!

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Configure Git for private modules
        run: |
          git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"

      - name: Set GOPRIVATE
        run: echo "GOPRIVATE=github.com/sveturs/*" >> $GITHUB_ENV

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true
          cache-dependency-path: backend/go.sum

      - name: Install dependencies
        working-directory: backend
        run: go mod download
```

### 4. Что делают эти настройки

#### permissions: packages: read
Дает GITHUB_TOKEN права на чтение packages (включая приватные репозитории) в организации. Без этого токен не сможет клонировать приватные модули.

#### Configure Git for private modules
```bash
git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"
```
- Настраивает Git использовать GITHUB_TOKEN для HTTPS запросов к GitHub
- Автоматически добавляет токен в URL при клонировании/загрузке

#### Set GOPRIVATE
```bash
echo "GOPRIVATE=github.com/sveturs/*" >> $GITHUB_ENV
```
- Указывает Go, что модули из `github.com/sveturs/*` приватные
- Go не будет пытаться использовать публичный прокси (proxy.golang.org)
- Go будет обращаться напрямую к GitHub используя Git

## Обновленные workflow файлы

Следующие workflow файлы обновлены:

1. `.github/workflows/go-format.yml` - проверка форматирования Go кода
2. `.github/workflows/unified-attributes-test.yml` - тесты unified attributes (4 джоба)
3. `.github/workflows/integration-tests.yml` - интеграционные тесты

## Проверка

После применения изменений:

1. Создайте PR или запустите workflow вручную
2. Проверьте логи GitHub Actions
3. Шаг "Install dependencies" должен успешно скачать приватные модули

## Альтернативный вариант (НЕ рекомендуется)

Можно создать Personal Access Token (PAT) вручную, но это:
- Менее безопасно (токен с широкими правами)
- Требует ручного управления
- Привязан к конкретному пользователю

Используйте встроенный `GITHUB_TOKEN` - он автоматический, безопасный и работает для всех в организации.

## Troubleshooting

### Ошибка: "terminal prompts disabled"
**Причина:** Git не может запросить credentials интерактивно.
**Решение:**
1. Убедитесь, что шаг "Configure Git for private modules" выполнен **ДО** `actions/setup-go`
2. Если используете `cache: true` в setup-go, Git ОБЯЗАТЕЛЬНО должен быть настроен до этого шага
3. Проверьте порядок: Checkout → Configure Git → Set GOPRIVATE → Setup Go

### Ошибка: "repository not found" или "404"
**Причина:** GITHUB_TOKEN не имеет доступа к приватному репозиторию.
**Решение:**
1. **Проверьте permissions в workflow** - ОБЯЗАТЕЛЬНО должно быть:
   ```yaml
   permissions:
     contents: read
     packages: read
   ```
2. Убедитесь, что репозиторий находится в той же организации (`sveturs`)
3. Проверьте настройки организации: Settings → Actions → General → Workflow permissions
4. Должно быть включено "Read and write permissions" или минимум "Read repository contents"

### Ошибка: "invalid version: git ls-remote"
**Причина:** Go не смог найти тег/версию модуля.
**Решение:**
- Проверьте, что в `go.mod` указана существующая версия/тег
- Убедитесь, что GOPRIVATE настроен правильно

## Дополнительная информация

- [GitHub Actions: Automatic token authentication](https://docs.github.com/en/actions/security-guides/automatic-token-authentication)
- [Go Modules: Private modules](https://go.dev/ref/mod#private-modules)
- [Git config: URL insteadOf](https://git-scm.com/docs/git-config#Documentation/git-config.txt-urlltbasegtinsteadOf)

---

**Дата последнего обновления:** 2025-01-17
