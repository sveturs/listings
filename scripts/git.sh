#!/bin/bash

# Скрипт для слияния и отправки изменений из текущей ветки в main
# Автор: Claude для dim
# Дата: 16-05-2025

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Функция для вывода с цветом
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Функция для показа информации о текущем состоянии
show_status() {
    print_message "$BLUE" "\n--- Текущее состояние ---"
    print_message "$BLUE" "Текущая ветка: $(git branch --show-current)"
    git status -s
    print_message "$BLUE" "------------------------\n"
}

# Функция для обработки конфликтов слияния
handle_merge_conflicts() {
    print_message "$YELLOW" "Обнаружены конфликты слияния!"

    PS3="Выберите действие: "
    options=("Открыть редактор для решения конфликтов" "Прервать слияние и вернуться в исходную ветку" "Запустить команду git mergetool" "Показать конфликтующие файлы")
    select opt in "${options[@]}"; do
        case $opt in
            "Открыть редактор для решения конфликтов")
                ${EDITOR:-vim}
                print_message "$YELLOW" "После решения конфликтов выполните:"
                print_message "$YELLOW" "1. git add <файлы с конфликтами>"
                print_message "$YELLOW" "2. git commit"
                print_message "$YELLOW" "3. Запустите скрипт снова с параметром --continue"
                exit 1
                ;;
            "Прервать слияние и вернуться в исходную ветку")
                git merge --abort
                git checkout "$source_branch"
                print_message "$RED" "Слияние отменено. Вы вернулись в ветку $source_branch."
                exit 1
                ;;
            "Запустить команду git mergetool")
                git mergetool
                print_message "$YELLOW" "После решения конфликтов выполните:"
                print_message "$YELLOW" "1. git commit"
                print_message "$YELLOW" "2. Запустите скрипт снова с параметром --continue"
                exit 1
                ;;
            "Показать конфликтующие файлы")
                print_message "$YELLOW" "Конфликтующие файлы:"
                git diff --name-only --diff-filter=U
                echo ""
                # Возвращаемся к выбору
                ;;
            *)
                print_message "$RED" "Неверный выбор. Пожалуйста, выберите число от 1 до ${#options[@]}."
                ;;
        esac
    done
}

# Функция для сохранения незафиксированных изменений
handle_uncommitted_changes() {
    print_message "$YELLOW" "Обнаружены незафиксированные изменения!"

    PS3="Выберите действие: "
    options=("Автоматически зафиксировать все изменения" "Выборочно добавить файлы" "Отложить изменения (stash)" "Отменить операцию")
    select opt in "${options[@]}"; do
        case $opt in
            "Автоматически зафиксировать все изменения")
                read -p "Введите сообщение коммита: " commit_message
                git add .
                git commit -m "$commit_message"
                print_message "$GREEN" "Изменения зафиксированы."
                return 0
                ;;
            "Выборочно добавить файлы")
                git add -i
                if git diff --cached --quiet; then
                    print_message "$YELLOW" "Файлы не были добавлены в индекс."
                    return 1
                else
                    read -p "Введите сообщение коммита: " commit_message
                    git commit -m "$commit_message"
                    print_message "$GREEN" "Выбранные изменения зафиксированы."
                    return 0
                fi
                ;;
            "Отложить изменения (stash)")
                read -p "Введите описание для stash (или оставьте пустым): " stash_message
                if [ -z "$stash_message" ]; then
                    git stash
                else
                    git stash push -m "$stash_message"
                fi
                print_message "$GREEN" "Изменения отложены (stashed)."
                return 0
                ;;
            "Отменить операцию")
                print_message "$RED" "Операция отменена."
                return 1
                ;;
            *)
                print_message "$RED" "Неверный выбор. Пожалуйста, выберите число от 1 до ${#options[@]}."
                ;;
        esac
    done
}

# Функция для проверки, существует ли удаленная ветка
check_remote_branch() {
    local branch_name=$1
    if git ls-remote --heads origin $branch_name | grep -q $branch_name; then
        return 0
    else
        return 1
    fi
}

# Функция для обработки ситуации, когда локальная ветка отстает от удаленной
handle_behind_remote() {
    print_message "$YELLOW" "Локальная ветка отстает от удаленной!"

    PS3="Выберите действие: "
    options=("Выполнить git pull" "Выполнить git pull --rebase" "Отменить операцию")
    select opt in "${options[@]}"; do
        case $opt in
            "Выполнить git pull")
                git pull origin $(git branch --show-current)
                if [ $? -eq 0 ]; then
                    print_message "$GREEN" "Локальная ветка обновлена."
                    return 0
                else
                    print_message "$RED" "Произошла ошибка при обновлении ветки."
                    return 1
                fi
                ;;
            "Выполнить git pull --rebase")
                git pull --rebase origin $(git branch --show-current)
                if [ $? -eq 0 ]; then
                    print_message "$GREEN" "Локальная ветка обновлена с перебазированием."
                    return 0
                else
                    print_message "$RED" "Произошла ошибка при обновлении ветки с перебазированием."
                    handle_rebase_conflicts
                    return 1
                fi
                ;;
            "Отменить операцию")
                print_message "$RED" "Операция отменена."
                return 1
                ;;
            *)
                print_message "$RED" "Неверный выбор. Пожалуйста, выберите число от 1 до ${#options[@]}."
                ;;
        esac
    done
}

# Функция для обработки конфликтов перебазирования
handle_rebase_conflicts() {
    print_message "$YELLOW" "Обнаружены конфликты перебазирования!"

    PS3="Выберите действие: "
    options=("Открыть редактор для решения конфликтов" "Прервать перебазирование" "Показать конфликтующие файлы")
    select opt in "${options[@]}"; do
        case $opt in
            "Открыть редактор для решения конфликтов")
                ${EDITOR:-vim}
                print_message "$YELLOW" "После решения конфликтов выполните:"
                print_message "$YELLOW" "1. git add <файлы с конфликтами>"
                print_message "$YELLOW" "2. git rebase --continue"
                print_message "$YELLOW" "3. Запустите скрипт снова с параметром --continue"
                exit 1
                ;;
            "Прервать перебазирование")
                git rebase --abort
                print_message "$RED" "Перебазирование отменено."
                exit 1
                ;;
            "Показать конфликтующие файлы")
                print_message "$YELLOW" "Конфликтующие файлы:"
                git diff --name-only --diff-filter=U
                echo ""
                # Возвращаемся к выбору
                ;;
            *)
                print_message "$RED" "Неверный выбор. Пожалуйста, выберите число от 1 до ${#options[@]}."
                ;;
        esac
    done
}

# Функция для проверки и установки источника и цели
set_branches() {
    # Если не указаны параметры, используем интерактивный режим
    if [ -z "$source_branch" ] || [ -z "$target_branch" ]; then
        # Получаем список локальных веток
        branches=($(git branch --format='%(refname:short)'))

        if [ -z "$source_branch" ]; then
            current_branch=$(git branch --show-current)
            print_message "$BLUE" "Текущая ветка: $current_branch"

            PS3="Выберите исходную ветку (или 0 для использования текущей): "
            select branch in "Использовать текущую ветку ($current_branch)" "${branches[@]}"; do
                if [ "$REPLY" = "0" ] || [ "$branch" = "Использовать текущую ветку ($current_branch)" ]; then
                    source_branch=$current_branch
                    break
                elif [ -n "$branch" ]; then
                    source_branch=$branch
                    break
                else
                    print_message "$RED" "Неверный выбор."
                fi
            done
        fi

        if [ -z "$target_branch" ]; then
            PS3="Выберите целевую ветку (или введите 'm' для main): "
            select branch in "main" "${branches[@]}"; do
                if [ "$REPLY" = "m" ] || [ "$branch" = "main" ]; then
                    target_branch="main"
                    break
                elif [ -n "$branch" ]; then
                    target_branch=$branch
                    break
                else
                    print_message "$RED" "Неверный выбор."
                fi
            done
        fi
    fi

    print_message "$GREEN" "Исходная ветка: $source_branch"
    print_message "$GREEN" "Целевая ветка: $target_branch"
}

# Функция для выбора стратегии слияния
select_merge_strategy() {
    PS3="Выберите стратегию слияния: "
    options=("Обычное слияние (merge)" "Перебазирование (rebase)" "Отмена")
    select opt in "${options[@]}"; do
        case $opt in
            "Обычное слияние (merge)")
                merge_strategy="merge"
                break
                ;;
            "Перебазирование (rebase)")
                merge_strategy="rebase"
                break
                ;;
            "Отмена")
                print_message "$RED" "Операция отменена."
                exit 1
                ;;
            *)
                print_message "$RED" "Неверный выбор. Пожалуйста, выберите число от 1 до ${#options[@]}."
                ;;
        esac
    done

    print_message "$GREEN" "Выбрана стратегия: $merge_strategy"
}

# Функция для проверки, нужно ли деплоить изменения
ask_deploy() {
    read -p "Хотите запустить деплой после слияния? (y/n): " deploy_answer
    if [[ "$deploy_answer" =~ ^[Yy]$ ]]; then
        PS3="Выберите тип деплоя: "
        options=("Обычный деплой" "Деплой с параметром -m (все)" "Отмена деплоя")
        select opt in "${options[@]}"; do
            case $opt in
                "Обычный деплой")
                    deploy_type="normal"
                    break
                    ;;
                "Деплой с параметром -m (все)")
                    deploy_type="all"
                    break
                    ;;
                "Отмена деплоя")
                    deploy_type="none"
                    break
                    ;;
                *)
                    print_message "$RED" "Неверный выбор. Пожалуйста, выберите число от 1 до ${#options[@]}."
                    ;;
            esac
        done
    else
        deploy_type="none"
    fi
}

# Основная логика скрипта
main() {
    clear
    print_message "$PURPLE" "=========================================="
    print_message "$PURPLE" "     Скрипт отправки кода в main          "
    print_message "$PURPLE" "=========================================="

    # Проверка, находимся ли мы в git репозитории
    if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
        print_message "$RED" "Текущая директория не является git репозиторием."
        exit 1
    fi

    # Инициализация переменных
    source_branch=""
    target_branch=""
    continue_mode=false
    merge_strategy=""
    deploy_type="none"

    # Обработка параметров командной строки
    while [[ $# -gt 0 ]]; do
        case $1 in
            --source|-s)
                source_branch="$2"
                shift 2
                ;;
            --target|-t)
                target_branch="$2"
                shift 2
                ;;
            --continue)
                continue_mode=true
                shift
                ;;
            --strategy)
                if [[ "$2" == "merge" || "$2" == "rebase" ]]; then
                    merge_strategy="$2"
                    shift 2
                else
                    print_message "$RED" "Неверная стратегия. Доступные варианты: merge, rebase"
                    exit 1
                fi
                ;;
            --deploy)
                case "$2" in
                    normal)
                        deploy_type="normal"
                        ;;
                    all)
                        deploy_type="all"
                        ;;
                    none)
                        deploy_type="none"
                        ;;
                    *)
                        print_message "$RED" "Неверный тип деплоя. Доступные варианты: normal, all, none"
                        exit 1
                        ;;
                esac
                shift 2
                ;;
            --help|-h)
                print_message "$BLUE" "Использование: $0 [опции]"
                print_message "$BLUE" "Опции:"
                print_message "$BLUE" "  --source, -s BRANCH     Исходная ветка (по умолчанию: текущая)"
                print_message "$BLUE" "  --target, -t BRANCH     Целевая ветка (по умолчанию: main)"
                print_message "$BLUE" "  --continue              Продолжить после разрешения конфликтов"
                print_message "$BLUE" "  --strategy STRATEGY     Стратегия слияния (merge или rebase)"
                print_message "$BLUE" "  --deploy TYPE           Тип деплоя (normal, all, none)"
                print_message "$BLUE" "  --help, -h              Показать эту справку"
                print_message "$BLUE" "\nПримеры использования:"
                print_message "$BLUE" "  $0                      Интерактивный режим"
                print_message "$BLUE" "  $0 -s ts -t main        Слияние из ветки ts в main"
                print_message "$BLUE" "  $0 --continue           Продолжить после разрешения конфликтов"
                exit 0
                ;;
            *)
                print_message "$RED" "Неизвестный параметр: $1"
                print_message "$RED" "Используйте --help для просмотра справки."
                exit 1
                ;;
        esac
    done

    show_status

    # Если не в режиме продолжения, начинаем процесс с начала
    if [ "$continue_mode" = false ]; then
        # Установка веток
        set_branches

        # Проверяем, есть ли незафиксированные изменения
        if ! git diff --quiet || ! git diff --cached --quiet; then
            handle_uncommitted_changes
            if [ $? -ne 0 ]; then
                exit 1
            fi
        fi

        # Сохраняем текущую ветку
        current_branch=$(git branch --show-current)

        # Если текущая ветка не совпадает с исходной, переключаемся
        if [ "$current_branch" != "$source_branch" ]; then
            print_message "$BLUE" "Переключение на исходную ветку $source_branch..."
            git checkout "$source_branch"
            if [ $? -ne 0 ]; then
                print_message "$RED" "Ошибка при переключении на ветку $source_branch."
                exit 1
            fi
        fi

        # Проверяем, существует ли удаленная ветка
        if check_remote_branch "$source_branch"; then
            # Проверяем, отстает ли локальная ветка от удаленной
            git fetch origin "$source_branch"
            if [ $(git rev-list HEAD..origin/"$source_branch" --count) -gt 0 ]; then
                print_message "$YELLOW" "Локальная ветка $source_branch отстает от удаленной."
                handle_behind_remote
                if [ $? -ne 0 ]; then
                    exit 1
                fi
            fi

            # Отправляем локальные изменения в удаленную ветку
            print_message "$BLUE" "Отправка изменений из $source_branch в удаленный репозиторий..."
            git push origin "$source_branch"
            if [ $? -ne 0 ]; then
                print_message "$RED" "Ошибка при отправке изменений в удаленный репозиторий."
                exit 1
            fi
        else
            print_message "$YELLOW" "Удаленная ветка $source_branch не существует. Создаем..."
            git push -u origin "$source_branch"
            if [ $? -ne 0 ]; then
                print_message "$RED" "Ошибка при создании удаленной ветки $source_branch."
                exit 1
            fi
        fi

        # Переключаемся на целевую ветку
        print_message "$BLUE" "Переключение на целевую ветку $target_branch..."
        git checkout "$target_branch"
        if [ $? -ne 0 ]; then
            print_message "$RED" "Ошибка при переключении на ветку $target_branch."
            git checkout "$source_branch" # Возвращаемся на исходную ветку
            exit 1
        fi

        # Подтягиваем изменения целевой ветки
        print_message "$BLUE" "Обновление $target_branch из удаленного репозитория..."
        git pull origin "$target_branch"
        if [ $? -ne 0 ]; then
            print_message "$RED" "Ошибка при обновлении ветки $target_branch."
            git checkout "$source_branch" # Возвращаемся на исходную ветку
            exit 1
        fi

        # Выбираем стратегию слияния, если не указана
        if [ -z "$merge_strategy" ]; then
            select_merge_strategy
        fi

        # Выполняем слияние или перебазирование
        if [ "$merge_strategy" = "merge" ]; then
            print_message "$BLUE" "Выполнение слияния $source_branch в $target_branch..."
            git merge "$source_branch"
            merge_result=$?
        else
            print_message "$BLUE" "Выполнение перебазирования $target_branch на $source_branch..."
            git checkout "$source_branch"
            git rebase "$target_branch"
            merge_result=$?
            if [ $merge_result -eq 0 ]; then
                git checkout "$target_branch"
                git merge "$source_branch"
                merge_result=$?
            fi
        fi

        # Обработка результата слияния
        if [ $merge_result -ne 0 ]; then
            print_message "$RED" "Возникли конфликты при слиянии."
            if [ "$merge_strategy" = "merge" ]; then
                handle_merge_conflicts
            else
                handle_rebase_conflicts
            fi
            exit 1
        fi

        # Отправляем изменения в удаленный репозиторий
        print_message "$BLUE" "Отправка изменений в удаленный репозиторий..."
        git push origin "$target_branch"
        if [ $? -ne 0 ]; then
            print_message "$RED" "Ошибка при отправке изменений в удаленный репозиторий."
            exit 1
        fi

        # Возвращаемся на исходную ветку
        print_message "$BLUE" "Возврат на исходную ветку $source_branch..."
        git checkout "$source_branch"
    else
        # Режим продолжения - определяем текущую ветку
        current_branch=$(git branch --show-current)
        print_message "$BLUE" "Продолжение в ветке $current_branch..."

        # Проверяем, находимся ли мы в состоянии слияния или перебазирования
        if [ -d ".git/rebase-apply" ] || [ -d ".git/rebase-merge" ]; then
            print_message "$BLUE" "Продолжение перебазирования..."
            git rebase --continue
            if [ $? -ne 0 ]; then
                print_message "$RED" "Всё еще есть конфликты. Разрешите их и запустите скрипт снова с --continue."
                exit 1
            fi

            # Переключаемся на целевую ветку и сливаем
            target_branch="main" # Предполагаем, что целевая ветка - main
            source_branch=$(git branch --show-current)

            git checkout "$target_branch"
            git merge "$source_branch"
            if [ $? -ne 0 ]; then
                print_message "$RED" "Возникли конфликты при слиянии после перебазирования."
                handle_merge_conflicts
                exit 1
            fi
        elif [ -f ".git/MERGE_HEAD" ]; then
            print_message "$BLUE" "Завершение слияния..."
            git commit
            if [ $? -ne 0 ]; then
                print_message "$RED" "Ошибка при фиксации слияния."
                exit 1
            fi
        else
            print_message "$RED" "Нет активного слияния или перебазирования для продолжения."
            exit 1
        fi

        # Отправка изменений
        print_message "$BLUE" "Отправка изменений в удаленный репозиторий..."
        git push origin "$current_branch"
        if [ $? -ne 0 ]; then
            print_message "$RED" "Ошибка при отправке изменений в удаленный репозиторий."
            exit 1
        fi
    fi

    # Спрашиваем о необходимости деплоя, если не указано
    if [ "$deploy_type" = "none" ] && [ "$continue_mode" = false ]; then
        ask_deploy
    fi

    # Выполняем деплой, если требуется
    if [ "$deploy_type" != "none" ]; then
        print_message "$BLUE" "Выполнение деплоя..."
        deploy_script="/data/hostel-booking-system/scripts/blue_green_deploy_on_svetu.rs.sh.improved"

        if [ -f "$deploy_script" ]; then
            if [ "$deploy_type" = "all" ]; then
                $deploy_script all -m
            else
                $deploy_script
            fi

            if [ $? -ne 0 ]; then
                print_message "$RED" "Ошибка при выполнении деплоя."
                exit 1
            fi
        else
            print_message "$RED" "Скрипт деплоя не найден: $deploy_script"
            exit 1
        fi
    fi

    print_message "$GREEN" "Процесс успешно завершен!"
    print_message "$GREEN" "Изменения из ветки $source_branch успешно отправлены в $target_branch."
}

# Запуск основной функции
main "$@"