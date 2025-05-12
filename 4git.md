# Инструкция по работе с Git

## Настройка токена GitHub

Для безопасной работы с удаленным репозиторием рекомендуется использовать персональный токен доступа (PAT):

1. Создайте токен GitHub:
   - Перейдите в Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Нажмите "Generate new token" → "Generate new token (classic)"
   - Укажите название и выберите scopes (минимально необходимо `repo`)
   - Нажмите "Generate token" и скопируйте полученный токен

2. Сохраните токен локально (единожды):
   ```bash
   echo "github_pat_your_token_here" > ~/.github_token
   chmod 600 ~/.github_token
   ```

3. Настройте Git для использования токена:
   ```bash
   git remote set-url origin https://username:$(cat ~/.github_token)@github.com/DmitruNS/hostel-booking-system.git
   ```
   где `username` - ваше имя пользователя GitHub.

## Отправка изменений из локальной ветки в main

### 1. Подготовка к отправке

1. Убедитесь, что вы находитесь в нужной ветке:
   ```bash
   git branch
   ```

2. Если вы в другой ветке, переключитесь на нужную:
   ```bash
   git checkout your-branch-name
   ```

3. Подготовьте файлы для коммита:
   ```bash
   git add path/to/file1 path/to/file2
   ```
   или добавьте все изменения:
   ```bash
   git add .
   ```

4. Создайте коммит с описательным сообщением:
   ```bash
   git commit -m "Описательное сообщение о внесенных изменениях"
   ```

### 2. Получение последних изменений из удаленного репозитория

1. Получите последние изменения с удаленного репозитория:
   ```bash
   git fetch origin
   ```

2. Просмотрите, какие изменения есть в удаленной ветке main, которых нет в вашей ветке:
   ```bash
   git log --oneline origin/main ^your-branch-name
   ```

3. И наоборот, какие изменения есть в вашей ветке, которых нет в main:
   ```bash
   git log --oneline your-branch-name ^origin/main
   ```

### 3. Слияние изменений

#### Вариант 1: Слияние в локальную ветку main

1. Переключитесь на ветку main:
   ```bash
   git checkout main
   ```

2. Получите последние изменения:
   ```bash
   git pull origin main
   ```

3. Выполните слияние вашей ветки:
   ```bash
   git merge your-branch-name
   ```

4. Отправьте результат слияния в удаленный репозиторий:
   ```bash
   git push origin main
   ```

#### Вариант 2: Слияние в вашу рабочую ветку и затем в main

1. Находясь в вашей ветке, получите изменения из main:
   ```bash
   git merge origin/main
   ```

2. Разрешите конфликты слияния, если они возникли:
   ```bash
   # После ручного разрешения конфликтов
   git add .
   git commit -m "Resolve merge conflicts"
   ```

3. Переключитесь на ветку main:
   ```bash
   git checkout main
   ```

4. Получите последние изменения:
   ```bash
   git pull origin main
   ```

5. Выполните слияние вашей ветки (уже включающей изменения из main):
   ```bash
   git merge your-branch-name
   ```

6. Отправьте изменения:
   ```bash
   git push origin main
   ```

#### Вариант 3: Создание Pull Request (рекомендуется)

1. Отправьте вашу ветку в удаленный репозиторий:
   ```bash
   git push origin your-branch-name
   ```

2. На GitHub создайте Pull Request из вашей ветки в main
3. Проверьте изменения и при необходимости обсудите с коллегами
4. После получения одобрения выполните слияние через интерфейс GitHub

### 4. Решение конфликтов слияния

Если при слиянии возникают конфликты:

1. Откройте файлы с конфликтами и отредактируйте их (найдите маркеры `<<<<<<<`, `=======`, `>>>>>>>`)
2. После редактирования добавьте файлы:
   ```bash
   git add .
   ```
3. Завершите процесс слияния:
   ```bash
   git commit
   ```

### 5. Возврат к предыдущему состоянию (если что-то пошло не так)

Если при слиянии возникли проблемы, вы можете отменить слияние:

```bash
git merge --abort
```

Или вернуться к предыдущему коммиту:

```bash
git reset --hard HEAD~1
```

## Работа с миграциями базы данных

При отправке изменений, содержащих миграции базы данных:

1. Убедитесь, что файлы миграций добавлены в коммит:
   ```bash
   git add backend/migrations/your_migration_file.up.sql backend/migrations/your_migration_file.down.sql
   ```

2. Создайте информативный коммит:
   ```bash
   git commit -m "Добавлена миграция для [описание функциональности]"
   ```

3. Для применения миграций на сервере используйте параметр `-m` при запуске скрипта деплоя:
   ```bash
   ./scripts/blue_green_deploy_on_svetu.rs.sh backend -m
   ```

## Полезные команды Git

### Просмотр статуса и истории

```bash
# Просмотр статуса репозитория
git status

# Просмотр лога коммитов
git log --oneline

# Просмотр изменений в файле
git diff path/to/file

# Просмотр изменений в репозитории между ветками
git diff branch1..branch2

# Просмотр изменений в конкретной директории между ветками
git diff branch1..branch2 -- path/to/directory/
```

### Работа с ветками

```bash
# Создание новой ветки
git checkout -b new-branch-name

# Переключение между ветками
git checkout branch-name

# Просмотр всех веток (локальных и удалённых)
git branch -a

# Удаление ветки (после слияния)
git branch -d branch-name

# Принудительное удаление ветки
git branch -D branch-name
```

### Отмена изменений

```bash
# Отмена изменений в файле (не добавленных в stage)
git checkout -- path/to/file

# Отмена последнего коммита, сохраняя изменения
git reset --soft HEAD~1

# Полная отмена последнего коммита (с удалением изменений)
git reset --hard HEAD~1

# Создание нового коммита, отменяющего изменения предыдущего
git revert HEAD
```