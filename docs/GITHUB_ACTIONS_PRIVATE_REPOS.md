# GitHub Actions: Настройка доступа к приватным репозиториям

## Проблема

При использовании приватных Go модулей (например, `github.com/sveturs/auth`) в GitHub Actions возникает ошибка:

```
fatal: could not read Username for 'https://github.com': terminal prompts disabled
```

Это происходит потому, что Go пытается скачать приватный модуль, но не имеет прав доступа.

## Решение

Для доступа к приватным репозиториям нужно настроить Git authentication и GOPRIVATE переменную.

### 1. Использование встроенного GITHUB_TOKEN

GitHub Actions автоматически предоставляет `GITHUB_TOKEN` с доступом к репозиториям организации. Этот токен НЕ нужно создавать вручную.

### 2. Конфигурация в workflow файлах

Добавьте следующие шаги **ПОСЛЕ** `actions/setup-go` и **ПЕРЕД** `go mod download`:

```yaml
- name: Set up Go
  uses: actions/setup-go@v5
  with:
    go-version: '1.21'

- name: Configure Git for private modules
  run: |
    git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"

- name: Set GOPRIVATE
  run: echo "GOPRIVATE=github.com/sveturs/*" >> $GITHUB_ENV

- name: Install dependencies
  working-directory: backend
  run: go mod download
```

### 3. Что делают эти настройки

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
**Решение:** Убедитесь, что шаг "Configure Git for private modules" выполнен перед `go mod download`.

### Ошибка: "repository not found" или "404"
**Причина:** GITHUB_TOKEN не имеет доступа к приватному репозиторию.
**Решение:**
- Убедитесь, что репозиторий находится в той же организации (`sveturs`)
- Проверьте настройки организации: Settings → Actions → General → Workflow permissions
- Должно быть включено "Read and write permissions" или минимум "Read repository contents"

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
