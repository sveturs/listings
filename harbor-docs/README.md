# Документация по интеграции с Harbor

Эта директория содержит документацию, связанную с интеграцией Sve Tu Platform с Harbor - приватным реестром Docker-образов.

## Содержимое

- **harbor_integration_status.md** - текущий статус интеграции с Harbor
- **harbor_migration_plan.md** - детальный план миграции на Harbor
- **harbor_installation_instructions.md** - инструкции по установке Harbor
- **harbor_config.yml** - пример конфигурации Harbor
- **harbor_final_report.md** - итоговый отчет о миграции
- **docker_image_list.md** - список образов Docker для миграции

## Ключевая информация

- **URL Harbor**: http://207.180.197.172
- **Логин**: admin
- **Пароль**: SveTu2025
- **Проект**: svetu

## Использование Harbor

Для использования образов из Harbor в проекте необходимо:

1. Авторизоваться в реестре: 
   ```
   docker login -u admin -p SveTu2025 207.180.197.172
   ```

2. Загрузить образ:
   ```
   docker pull 207.180.197.172/svetu/backend/api:latest
   ```

3. Отправить образ:
   ```
   docker push 207.180.197.172/svetu/backend/api:latest
   ```

## Важные замечания

- Необходимо настроить HTTPS для продакшн использования
- Рекомендуется изменить пароль admin после завершения миграции
- Образы в docker-compose.yml и docker-compose.prod.yml обновлены для использования Harbor