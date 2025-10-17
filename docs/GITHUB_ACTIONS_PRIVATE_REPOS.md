# GitHub Actions: Настройка доступа к приватным репозиториям

## Проблема

При использовании приватных Go модулей (например, `github.com/sveturs/auth`) в GitHub Actions возникает ошибка:

```
fatal: could not read Username for 'https://github.com': terminal prompts disabled
```

Это происходит потому, что Go пытается скачать приватный модуль, но не имеет прав доступа.

## Решение

**ВАЖНО:** `GITHUB_TOKEN` НЕ РАБОТАЕТ для доступа к приватным Go модулям в других репозиториях!

Нужно создать Personal Access Token (PAT) с правами `repo` и добавить его в secrets репозитория.

### 1. Создать Personal Access Token (Classic)

1. Перейдите на https://github.com/settings/tokens
2. Нажмите "Generate new token" → "Generate new token (classic)"
3. Настройте токен:
   - **Note**: `SVETU_PRIVATE_MODULES_TOKEN`
   - **Expiration**: 90 days (или No expiration)
   - **Scopes**:
     - ✅ `repo` (Full control of private repositories)
     - ✅ `read:packages`
4. Нажмите "Generate token" и **скопируйте токен**

### 2. Добавить токен в GitHub Secrets

1. Перейдите в: https://github.com/sveturs/svetu/settings/secrets/actions
2. Нажмите "New repository secret"
3. Настройте:
   - **Name**: `PRIVATE_MODULES_TOKEN`
   - **Secret**: вставьте скопированный токен
4. Нажмите "Add secret"

### 3. Конфигурация в workflow файлах

**ВАЖНО:** Git конфигурация должна быть **ДО** `actions/setup-go`, потому что setup-go с `cache: true` пытается восстановить кэш зависимостей сразу при выполнении.

Полный пример workflow файла:

```yaml
name: My Workflow

on:
  pull_request:

permissions:
  contents: read

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Configure Git for private modules
        run: |
          git config --global url."https://${{ secrets.PRIVATE_MODULES_TOKEN }}@github.com/".insteadOf "https://github.com/"

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

#### Configure Git for private modules
```bash
git config --global url."https://${{ secrets.PRIVATE_MODULES_TOKEN }}@github.com/".insteadOf "https://github.com/"
```
- Настраивает Git использовать PRIVATE_MODULES_TOKEN для HTTPS запросов к GitHub
- Автоматически добавляет токен в URL при клонировании/загрузке
- PAT имеет полные права на приватные репозитории (в отличие от GITHUB_TOKEN)

#### Set GOPRIVATE
```bash
echo "GOPRIVATE=github.com/sveturs/*" >> $GITHUB_ENV
```
- Указывает Go, что модули из `github.com/sveturs/*` приватные
- Go не будет пытаться использовать публичный прокси (proxy.golang.org)
- Go будет обращаться напрямую к GitHub используя Git

## Обновленные workflow файлы

Следующие workflow файлы обновлены для использования `PRIVATE_MODULES_TOKEN`:

1. `.github/workflows/go-format.yml` - проверка форматирования Go кода
2. `.github/workflows/unified-attributes-test.yml` - тесты unified attributes (4 джоба)
3. `.github/workflows/integration-tests.yml` - интеграционные тесты

## Проверка

После применения изменений:

1. Убедитесь, что `PRIVATE_MODULES_TOKEN` добавлен в GitHub Secrets
2. Создайте PR или запустите workflow вручную
3. Проверьте логи GitHub Actions
4. Шаг "Install dependencies" должен успешно скачать приватные модули

## Почему не GITHUB_TOKEN?

`GITHUB_TOKEN` имеет ограничения безопасности и **НЕ МОЖЕТ** клонировать другие приватные репозитории организации для целей Go modules. Это документированное ограничение GitHub Actions.

Personal Access Token (PAT) с правами `repo` - единственное рабочее решение для приватных Go модулей.

## Troubleshooting

### Ошибка: "terminal prompts disabled"
**Причина:** Git не может запросить credentials интерактивно.
**Решение:**
1. Убедитесь, что шаг "Configure Git for private modules" выполнен **ДО** `actions/setup-go`
2. Если используете `cache: true` в setup-go, Git ОБЯЗАТЕЛЬНО должен быть настроен до этого шага
3. Проверьте порядок: Checkout → Configure Git → Set GOPRIVATE → Setup Go

### Ошибка: "repository not found" или "404"
**Причина:** Токен не имеет доступа к приватному репозиторию.
**Решение:**
1. **Проверьте, что токен добавлен в Secrets:**
   - Название: `PRIVATE_MODULES_TOKEN`
   - Значение: Personal Access Token с правами `repo`
2. Убедитесь, что токен создан с правильными scopes:
   - ✅ `repo` - обязательно!
   - ✅ `read:packages` - для пакетов
3. Проверьте, что токен не истек (GitHub показывает дату истечения)

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
