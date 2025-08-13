#!/usr/bin/env node

/**
 * Скрипт для синхронизации ключей переводов между языками
 * Находит отсутствующие ключи и добавляет их с placeholder значениями
 */

const fs = require('fs');
const path = require('path');
const { glob } = require('glob');

const MESSAGES_DIR = path.join(__dirname, '../src/messages');
const LANGUAGES = ['en', 'ru', 'sr'];

// Цвета для консоли
const colors = {
  reset: '\x1b[0m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m',
};

// Получить все ключи из объекта рекурсивно
function getAllKeys(obj, prefix = '') {
  let keys = [];

  for (const key in obj) {
    const fullKey = prefix ? `${prefix}.${key}` : key;

    if (
      typeof obj[key] === 'object' &&
      obj[key] !== null &&
      !Array.isArray(obj[key])
    ) {
      keys = keys.concat(getAllKeys(obj[key], fullKey));
    } else {
      keys.push(fullKey);
    }
  }

  return keys;
}

// Установить значение по пути ключа
function setNestedValue(obj, path, value) {
  const keys = path.split('.');
  let current = obj;

  for (let i = 0; i < keys.length - 1; i++) {
    const key = keys[i];
    if (!(key in current) || typeof current[key] !== 'object') {
      current[key] = {};
    }
    current = current[key];
  }

  current[keys[keys.length - 1]] = value;
}

// Получить значение по пути ключа
function getNestedValue(obj, path) {
  const keys = path.split('.');
  let current = obj;

  for (const key of keys) {
    if (current && typeof current === 'object' && key in current) {
      current = current[key];
    } else {
      return undefined;
    }
  }

  return current;
}

// Создать placeholder значение
function createPlaceholder(key, targetLang, sourceValue) {
  if (typeof sourceValue === 'string') {
    // Если в оригинале есть переменные {var}, сохраняем их
    const variables = sourceValue.match(/\{[^}]+\}/g) || [];
    const varsString = variables.length > 0 ? ' ' + variables.join(' ') : '';

    return `[MISSING_${targetLang.toUpperCase()}] ${key.split('.').pop()}${varsString}`;
  }
  return sourceValue; // Для не-строковых значений возвращаем как есть
}

// Главная функция синхронизации
async function syncTranslations() {
  console.log(
    `${colors.cyan}🔄 Начинаю синхронизацию ключей переводов...${colors.reset}\n`
  );

  // Получаем список всех модулей
  const enFiles = await glob(`${MESSAGES_DIR}/en/*.json`);
  const modules = enFiles.map((file) => path.basename(file, '.json'));

  console.log(
    `${colors.blue}📦 Найдено модулей: ${modules.length}${colors.reset}`
  );
  console.log(
    `${colors.blue}🌍 Языки: ${LANGUAGES.join(', ')}${colors.reset}\n`
  );

  let totalMissing = 0;
  let totalAdded = 0;
  let totalModulesWithIssues = 0;

  const report = {
    modules: {},
    summary: {
      totalMissing: 0,
      totalAdded: 0,
      byLanguage: {},
    },
  };

  // Инициализация счетчиков по языкам
  LANGUAGES.forEach((lang) => {
    report.summary.byLanguage[lang] = {
      missing: 0,
      added: 0,
    };
  });

  // Обрабатываем каждый модуль
  for (const module of modules) {
    let moduleHasIssues = false;
    const moduleReport = {
      missing: {},
      added: {},
    };

    // Загружаем все файлы для модуля
    const translations = {};
    const allKeys = new Set();

    for (const lang of LANGUAGES) {
      const filePath = path.join(MESSAGES_DIR, lang, `${module}.json`);

      try {
        const content = fs.readFileSync(filePath, 'utf8');
        translations[lang] = JSON.parse(content);

        // Собираем все ключи
        const keys = getAllKeys(translations[lang]);
        keys.forEach((key) => allKeys.add(key));
      } catch (error) {
        console.log(
          `${colors.yellow}⚠️  Файл не найден или поврежден: ${filePath}${colors.reset}`
        );
        translations[lang] = {};
      }
    }

    // Проверяем наличие всех ключей в каждом языке
    for (const lang of LANGUAGES) {
      const missingKeys = [];

      for (const key of allKeys) {
        const value = getNestedValue(translations[lang], key);

        if (value === undefined) {
          missingKeys.push(key);

          // Находим значение из другого языка для использования как образец
          let sourceValue = null;
          let sourceLang = null;

          for (const otherLang of LANGUAGES) {
            if (otherLang !== lang) {
              const otherValue = getNestedValue(translations[otherLang], key);
              if (otherValue !== undefined) {
                sourceValue = otherValue;
                sourceLang = otherLang;
                break;
              }
            }
          }

          // Добавляем placeholder
          const placeholder = createPlaceholder(key, lang, sourceValue);
          setNestedValue(translations[lang], key, placeholder);

          totalAdded++;
          report.summary.byLanguage[lang].added++;
        }
      }

      if (missingKeys.length > 0) {
        moduleHasIssues = true;
        totalMissing += missingKeys.length;
        moduleReport.missing[lang] = missingKeys;
        report.summary.byLanguage[lang].missing += missingKeys.length;

        console.log(
          `${colors.red}❌ ${module}.json (${lang}): отсутствует ${missingKeys.length} ключей${colors.reset}`
        );

        if (missingKeys.length <= 5) {
          missingKeys.forEach((key) => {
            console.log(`   ${colors.yellow}└─ ${key}${colors.reset}`);
          });
        } else {
          missingKeys.slice(0, 3).forEach((key) => {
            console.log(`   ${colors.yellow}└─ ${key}${colors.reset}`);
          });
          console.log(
            `   ${colors.yellow}└─ ... и еще ${missingKeys.length - 3} ключей${colors.reset}`
          );
        }
      }
    }

    // Сохраняем обновленные файлы
    if (moduleHasIssues) {
      totalModulesWithIssues++;

      for (const lang of LANGUAGES) {
        const filePath = path.join(MESSAGES_DIR, lang, `${module}.json`);
        const content = JSON.stringify(translations[lang], null, 2) + '\n';

        try {
          fs.writeFileSync(filePath, content);
          console.log(
            `${colors.green}✅ Обновлен: ${lang}/${module}.json${colors.reset}`
          );
        } catch (error) {
          console.error(
            `${colors.red}❌ Ошибка записи: ${filePath}${colors.reset}`
          );
          console.error(error.message);
        }
      }

      report.modules[module] = moduleReport;
    }
  }

  // Финальный отчет
  console.log(`\n${colors.cyan}📊 ИТОГОВЫЙ ОТЧЕТ${colors.reset}`);
  console.log(`${colors.cyan}${'='.repeat(50)}${colors.reset}`);

  if (totalMissing === 0) {
    console.log(
      `${colors.green}✨ Все ключи синхронизированы! Проблем не найдено.${colors.reset}`
    );
  } else {
    console.log(
      `${colors.yellow}📝 Модулей с проблемами: ${totalModulesWithIssues}/${modules.length}${colors.reset}`
    );
    console.log(
      `${colors.yellow}🔍 Всего отсутствующих ключей: ${totalMissing}${colors.reset}`
    );
    console.log(
      `${colors.green}✅ Добавлено placeholder'ов: ${totalAdded}${colors.reset}`
    );

    console.log(`\n${colors.blue}По языкам:${colors.reset}`);
    for (const lang of LANGUAGES) {
      const stats = report.summary.byLanguage[lang];
      if (stats.missing > 0) {
        console.log(
          `  ${lang}: ${colors.red}${stats.missing} отсутствовало${colors.reset}, ${colors.green}${stats.added} добавлено${colors.reset}`
        );
      }
    }
  }

  // Сохраняем детальный отчет
  const reportPath = path.join(__dirname, 'translation-sync-report.json');
  fs.writeFileSync(reportPath, JSON.stringify(report, null, 2));
  console.log(
    `\n${colors.magenta}📄 Детальный отчет сохранен: ${reportPath}${colors.reset}`
  );

  // Проверка на placeholder'ы
  console.log(
    `\n${colors.cyan}🔍 Проверка существующих placeholder'ов...${colors.reset}`
  );

  let placeholderCount = 0;
  const placeholderModules = new Set();

  for (const module of modules) {
    for (const lang of LANGUAGES) {
      const filePath = path.join(MESSAGES_DIR, lang, `${module}.json`);

      try {
        const content = fs.readFileSync(filePath, 'utf8');
        const placeholderMatches = content.match(
          /\[MISSING_[A-Z]+\]|\[RU\]|\[EN\]|\[SR\]/g
        );

        if (placeholderMatches) {
          placeholderCount += placeholderMatches.length;
          placeholderModules.add(`${lang}/${module}`);
        }
      } catch (error) {
        // Игнорируем ошибки чтения
      }
    }
  }

  if (placeholderCount > 0) {
    console.log(
      `${colors.yellow}⚠️  Найдено ${placeholderCount} placeholder'ов в ${placeholderModules.size} файлах${colors.reset}`
    );
    console.log(
      `${colors.yellow}   Используйте 'npm run clean-placeholders' для их поиска${colors.reset}`
    );
  } else {
    console.log(`${colors.green}✅ Placeholder'ы не найдены${colors.reset}`);
  }

  console.log(`\n${colors.green}✨ Синхронизация завершена!${colors.reset}`);

  return {
    success: true,
    totalMissing,
    totalAdded,
    report,
  };
}

// Запуск если вызван напрямую
if (require.main === module) {
  syncTranslations()
    .then((result) => {
      process.exit(result.totalMissing > 0 ? 1 : 0);
    })
    .catch((error) => {
      console.error(
        `${colors.red}❌ Критическая ошибка:${colors.reset}`,
        error
      );
      process.exit(1);
    });
}

module.exports = { syncTranslations };
