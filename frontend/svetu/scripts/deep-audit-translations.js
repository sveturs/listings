#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const glob = require('glob');

// Цвета для вывода
const colors = {
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  reset: '\x1b[0m',
};

// Счетчики проблем
let totalProblems = 0;
const problems = {
  missingKeys: [],
  wrongUsage: [],
  emptyTranslations: [],
  unusedKeys: [],
  duplicateKeys: [],
  snakeCaseKeys: [],
  invalidNesting: [],
};

// Загрузка всех переводов
function loadAllTranslations() {
  const translations = {};
  const locales = ['en', 'ru', 'sr'];

  locales.forEach((locale) => {
    translations[locale] = {};
    const localeDir = path.join(__dirname, '..', 'src', 'messages', locale);

    if (fs.existsSync(localeDir)) {
      const files = fs
        .readdirSync(localeDir)
        .filter((f) => f.endsWith('.json'));
      files.forEach((file) => {
        const module = file.replace('.json', '');
        try {
          const content = JSON.parse(
            fs.readFileSync(path.join(localeDir, file), 'utf8')
          );
          translations[locale][module] = content;
        } catch (e) {
          console.error(
            `${colors.red}Ошибка загрузки ${locale}/${file}: ${e.message}${colors.reset}`
          );
        }
      });
    }
  });

  return translations;
}

// Получение всех ключей из объекта
function getAllKeys(obj, prefix = '') {
  const keys = [];

  for (const key in obj) {
    if (obj.hasOwnProperty(key)) {
      const fullKey = prefix ? `${prefix}.${key}` : key;

      if (
        typeof obj[key] === 'object' &&
        obj[key] !== null &&
        !Array.isArray(obj[key])
      ) {
        keys.push(...getAllKeys(obj[key], fullKey));
      } else {
        keys.push(fullKey);
      }
    }
  }

  return keys;
}

// Проверка согласованности ключей между языками
function checkKeyConsistency(translations) {
  console.log(
    `\n${colors.blue}=== Проверка согласованности ключей ===${colors.reset}`
  );

  const locales = Object.keys(translations);
  const modules = [
    ...new Set(locales.flatMap((l) => Object.keys(translations[l]))),
  ];

  modules.forEach((module) => {
    const keysByLocale = {};

    locales.forEach((locale) => {
      if (translations[locale][module]) {
        keysByLocale[locale] = getAllKeys(translations[locale][module]);
      } else {
        problems.missingKeys.push({
          module,
          locale,
          message: `Модуль ${module} отсутствует в локали ${locale}`,
        });
      }
    });

    // Сравнение ключей между локалями
    if (keysByLocale.en) {
      const enKeys = new Set(keysByLocale.en);

      ['ru', 'sr'].forEach((locale) => {
        if (keysByLocale[locale]) {
          const localeKeys = new Set(keysByLocale[locale]);

          // Ключи есть в en, но нет в locale
          enKeys.forEach((key) => {
            if (!localeKeys.has(key)) {
              problems.missingKeys.push({
                module,
                locale,
                key,
                message: `Ключ "${key}" есть в en, но отсутствует в ${locale}`,
              });
            }
          });

          // Ключи есть в locale, но нет в en
          localeKeys.forEach((key) => {
            if (!enKeys.has(key)) {
              problems.unusedKeys.push({
                module,
                locale,
                key,
                message: `Ключ "${key}" есть в ${locale}, но отсутствует в en`,
              });
            }
          });
        }
      });
    }
  });

  console.log(`Найдено ${problems.missingKeys.length} отсутствующих ключей`);
  console.log(`Найдено ${problems.unusedKeys.length} лишних ключей`);
}

// Проверка использования в компонентах
function checkComponentUsage() {
  console.log(
    `\n${colors.blue}=== Проверка использования в компонентах ===${colors.reset}`
  );

  const componentFiles = glob.sync(
    path.join(__dirname, '..', 'src', '**', '*.{tsx,ts,jsx,js}'),
    {
      ignore: [
        '**/node_modules/**',
        '**/scripts/**',
        '**/*.test.*',
        '**/*.spec.*',
      ],
    }
  );

  const translations = loadAllTranslations();
  const allAvailableKeys = new Set();

  // Собираем все доступные ключи
  Object.values(translations).forEach((locale) => {
    Object.entries(locale).forEach(([module, content]) => {
      getAllKeys(content).forEach((key) => {
        allAvailableKeys.add(`${module}.${key}`);
      });
    });
  });

  componentFiles.forEach((file) => {
    const content = fs.readFileSync(file, 'utf8');

    // Проверка на useTranslations без аргумента
    const emptyUseTranslations = content.match(/useTranslations\s*\(\s*\)/g);
    if (emptyUseTranslations) {
      problems.emptyTranslations.push({
        file: path.relative(process.cwd(), file),
        count: emptyUseTranslations.length,
      });
    }

    // Проверка на неправильное использование вложенных namespace
    const nestedNamespace = content.match(
      /useTranslations\s*\(\s*['"`]([^'"`]+\.[^'"`]+)['"`]\s*\)/g
    );
    if (nestedNamespace) {
      nestedNamespace.forEach((match) => {
        const namespace = match.match(/['"`]([^'"`]+)['"`]/)[1];
        problems.wrongUsage.push({
          file: path.relative(process.cwd(), file),
          usage: match,
          namespace,
        });
      });
    }

    // Проверка использования ключей
    const useTranslationsMatch = content.match(
      /const\s+(\w+)\s*=\s*useTranslations\s*\(\s*['"`](\w+)['"`]\s*\)/g
    );
    if (useTranslationsMatch) {
      useTranslationsMatch.forEach((match) => {
        const [, varName, module] = match.match(
          /const\s+(\w+)\s*=\s*useTranslations\s*\(\s*['"`](\w+)['"`]\s*\)/
        );

        // Ищем использование переводов
        const usageRegex = new RegExp(
          `${varName}\\s*\\(\\s*['"\`]([^'"\`]+)['"\`]`,
          'g'
        );
        let usage;
        while ((usage = usageRegex.exec(content)) !== null) {
          const key = usage[1];
          const fullKey = `${module}.${key}`;

          // Проверяем существование ключа
          if (!allAvailableKeys.has(fullKey)) {
            problems.missingKeys.push({
              file: path.relative(process.cwd(), file),
              key: fullKey,
              line: content.substring(0, usage.index).split('\n').length,
              message: `Использование несуществующего ключа "${fullKey}"`,
            });
          }
        }
      });
    }
  });

  console.log(
    `Найдено ${problems.emptyTranslations.length} файлов с пустым useTranslations()`
  );
  console.log(
    `Найдено ${problems.wrongUsage.length} неправильных использований namespace`
  );
}

// Проверка на snake_case ключи
function checkSnakeCaseKeys(translations) {
  console.log(
    `\n${colors.blue}=== Проверка на snake_case ключи ===${colors.reset}`
  );

  Object.entries(translations).forEach(([locale, modules]) => {
    Object.entries(modules).forEach(([module, content]) => {
      const keys = getAllKeys(content);
      keys.forEach((key) => {
        if (key.includes('_')) {
          problems.snakeCaseKeys.push({
            locale,
            module,
            key,
            suggestion: key.replace(/_([a-z])/g, (_, letter) =>
              letter.toUpperCase()
            ),
          });
        }
      });
    });
  });

  console.log(`Найдено ${problems.snakeCaseKeys.length} snake_case ключей`);
}

// Проверка дубликатов
function checkDuplicates(translations) {
  console.log(`\n${colors.blue}=== Проверка дубликатов ===${colors.reset}`);

  Object.entries(translations).forEach(([locale, modules]) => {
    Object.entries(modules).forEach(([module, content]) => {
      const jsonStr = JSON.stringify(content, null, 2);
      const lines = jsonStr.split('\n');
      const seenKeys = new Map();

      lines.forEach((line, index) => {
        const keyMatch = line.match(/^\s*"([^"]+)":/);
        if (keyMatch) {
          const key = keyMatch[1];
          if (seenKeys.has(key)) {
            problems.duplicateKeys.push({
              locale,
              module,
              key,
              firstLine: seenKeys.get(key) + 1,
              duplicateLine: index + 1,
            });
          } else {
            seenKeys.set(key, index);
          }
        }
      });
    });
  });

  console.log(`Найдено ${problems.duplicateKeys.length} дублирующихся ключей`);
}

// Вывод отчета
function printReport() {
  console.log(`\n${colors.yellow}=== ОТЧЕТ О ПРОБЛЕМАХ ===${colors.reset}\n`);

  totalProblems = Object.values(problems).reduce(
    (sum, arr) => sum + arr.length,
    0
  );

  if (totalProblems === 0) {
    console.log(
      `${colors.green}✅ Проблем не найдено! Система переводов в порядке.${colors.reset}`
    );
    return;
  }

  console.log(
    `${colors.red}Всего найдено проблем: ${totalProblems}${colors.reset}\n`
  );

  // Вывод отсутствующих ключей
  if (problems.missingKeys.length > 0) {
    console.log(
      `${colors.red}❌ Отсутствующие ключи (${problems.missingKeys.length})${colors.reset}`
    );
    problems.missingKeys.slice(0, 10).forEach((p) => {
      if (p.file) {
        console.log(`  - ${p.file}:${p.line || '?'} - ${p.message}`);
      } else {
        console.log(`  - ${p.module}/${p.locale}: ${p.message || p.key}`);
      }
    });
    if (problems.missingKeys.length > 10) {
      console.log(`  ... и еще ${problems.missingKeys.length - 10} проблем`);
    }
    console.log();
  }

  // Вывод неправильного использования
  if (problems.wrongUsage.length > 0) {
    console.log(
      `${colors.red}❌ Неправильное использование namespace (${problems.wrongUsage.length})${colors.reset}`
    );
    problems.wrongUsage.slice(0, 5).forEach((p) => {
      console.log(`  - ${p.file}: ${p.usage}`);
    });
    if (problems.wrongUsage.length > 5) {
      console.log(`  ... и еще ${problems.wrongUsage.length - 5} проблем`);
    }
    console.log();
  }

  // Вывод пустых useTranslations
  if (problems.emptyTranslations.length > 0) {
    console.log(
      `${colors.red}❌ Пустые useTranslations() (${problems.emptyTranslations.length} файлов)${colors.reset}`
    );
    problems.emptyTranslations.slice(0, 5).forEach((p) => {
      console.log(`  - ${p.file} (${p.count} раз)`);
    });
    if (problems.emptyTranslations.length > 5) {
      console.log(
        `  ... и еще ${problems.emptyTranslations.length - 5} файлов`
      );
    }
    console.log();
  }

  // Вывод snake_case ключей
  if (problems.snakeCaseKeys.length > 0) {
    console.log(
      `${colors.yellow}⚠️  Snake_case ключи (${problems.snakeCaseKeys.length})${colors.reset}`
    );
    problems.snakeCaseKeys.slice(0, 5).forEach((p) => {
      console.log(
        `  - ${p.locale}/${p.module}: "${p.key}" → "${p.suggestion}"`
      );
    });
    if (problems.snakeCaseKeys.length > 5) {
      console.log(`  ... и еще ${problems.snakeCaseKeys.length - 5} ключей`);
    }
    console.log();
  }

  // Сохранение полного отчета
  const fullReport = {
    timestamp: new Date().toISOString(),
    totalProblems,
    problems,
  };

  fs.writeFileSync(
    path.join(__dirname, 'deep-audit-report.json'),
    JSON.stringify(fullReport, null, 2)
  );

  console.log(
    `\n${colors.blue}Полный отчет сохранен в scripts/deep-audit-report.json${colors.reset}`
  );
}

// Проверка браузера
async function checkBrowserErrors() {
  console.log(
    `\n${colors.blue}=== Проверка ошибок в браузере ===${colors.reset}`
  );
  console.log(
    'Откройте http://localhost:3001 и проверьте консоль браузера на наличие ошибок переводов'
  );
  console.log('Особенно проверьте страницы:');
  console.log('  - / (главная)');
  console.log('  - /map (карта)');
  console.log('  - /profile/storefronts (витрины)');
  console.log('  - /admin (админка)');
}

// Основная функция
async function main() {
  console.log(
    `${colors.blue}Запуск глубокого аудита системы переводов...${colors.reset}`
  );

  const translations = loadAllTranslations();

  checkKeyConsistency(translations);
  checkComponentUsage();
  checkSnakeCaseKeys(translations);
  checkDuplicates(translations);

  await checkBrowserErrors();

  printReport();
}

main().catch(console.error);
