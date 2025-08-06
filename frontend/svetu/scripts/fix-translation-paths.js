#!/usr/bin/env node

/**
 * Скрипт для исправления неправильного использования вложенных путей в useTranslations
 * Заменяет useTranslations('module.submodule') на useTranslations('module')
 * и обновляет все вызовы t() чтобы добавить префикс submodule.
 */

const fs = require('fs');
const path = require('path');
const glob = require('glob');

// Паттерн для поиска useTranslations с вложенными путями
const USE_TRANSLATIONS_PATTERN =
  /const\s+(\w+)\s*=\s*useTranslations\(['"](\w+)\.(\w+)['"]\)/g;

// Функция для обработки файла
function processFile(filePath) {
  let fileContent = fs.readFileSync(filePath, 'utf8');
  let hasChanges = false;

  // Сохраняем информацию о найденных заменах
  const replacements = [];

  // Находим все вызовы useTranslations с вложенными путями
  fileContent.replace(
    USE_TRANSLATIONS_PATTERN,
    (match, varName, module, submodule) => {
      replacements.push({ varName, module, submodule });
      hasChanges = true;
      return match;
    }
  );

  if (!hasChanges) {
    return false;
  }

  console.log(`\nОбработка файла: ${filePath}`);

  // Применяем замены
  replacements.forEach(({ varName, module, submodule }) => {
    console.log(
      `  Найдено: const ${varName} = useTranslations('${module}.${submodule}')`
    );

    // Заменяем useTranslations
    fileContent = fileContent.replace(
      new RegExp(
        `const\\s+${varName}\\s*=\\s*useTranslations\\(['"]${module}\\.${submodule}['"]\\)`,
        'g'
      ),
      `const ${varName} = useTranslations('${module}')`
    );

    // Находим все вызовы t() и добавляем префикс
    // Паттерн для поиска вызовов вида t('key') или t("key")
    const tCallPattern = new RegExp(`${varName}\\((['"])([^'"]+)\\1`, 'g');

    fileContent = fileContent.replace(tCallPattern, (match, quote, key) => {
      // Если ключ уже содержит точку, не добавляем префикс
      if (key.includes('.')) {
        return match;
      }
      return `${varName}(${quote}${submodule}.${key}${quote}`;
    });

    console.log(
      `  Заменено на: const ${varName} = useTranslations('${module}')`
    );
    console.log(`  Обновлены вызовы ${varName}() с префиксом '${submodule}.'`);
  });

  // Записываем обновленный файл
  fs.writeFileSync(filePath, fileContent);
  return true;
}

// Главная функция
function main() {
  console.log('Поиск файлов с неправильными путями переводов...\n');

  // Находим все файлы TypeScript и TSX
  const files = glob.sync('src/**/*.{ts,tsx}', {
    cwd: path.join(__dirname, '..'),
    absolute: true,
    ignore: ['**/node_modules/**', '**/build/**', '**/.next/**'],
  });

  let processedCount = 0;

  files.forEach((file) => {
    if (processFile(file)) {
      processedCount++;
    }
  });

  console.log(`\nОбработано файлов: ${processedCount}`);

  if (processedCount > 0) {
    console.log('\n⚠️  ВАЖНО: Проверьте изменения и протестируйте приложение!');
    console.log(
      'Некоторые сложные случаи могут требовать ручной корректировки.'
    );
  }
}

// Запуск скрипта
main();
