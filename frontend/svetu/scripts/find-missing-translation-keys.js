#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const glob = require('glob');

console.log('🔍 Поиск отсутствующих ключей переводов...\n');

// Загружаем все JSON файлы переводов
const locales = ['ru', 'en', 'sr'];
const modules = [
  'common',
  'auth',
  'marketplace',
  'admin',
  'storefront',
  'cars',
  'chat',
  'cart',
  'realEstate',
  'services',
  'map',
  'misc',
];

const translations = {};

// Загружаем существующие переводы
locales.forEach((locale) => {
  translations[locale] = {};
  modules.forEach((module) => {
    try {
      const filePath = path.join(
        __dirname,
        `../src/messages/${locale}/${module}.json`
      );
      if (fs.existsSync(filePath)) {
        translations[locale][module] = JSON.parse(
          fs.readFileSync(filePath, 'utf8')
        );
      }
    } catch (e) {
      console.error(`Ошибка загрузки ${locale}/${module}.json:`, e.message);
    }
  });
});

// Находим все использования переводов в коде
const usedKeys = new Map(); // module -> Set of keys
const missingKeys = new Map(); // module -> Set of missing keys

const patterns = ['src/**/*.tsx', 'src/**/*.ts'];

patterns.forEach((pattern) => {
  const files = glob.sync(pattern, {
    cwd: path.join(__dirname, '..'),
    ignore: [
      '**/node_modules/**',
      '**/*.test.*',
      '**/*.spec.*',
      '**/messages/**',
    ],
  });

  files.forEach((file) => {
    const filePath = path.join(__dirname, '..', file);
    const content = fs.readFileSync(filePath, 'utf8');

    // Паттерны для поиска использования переводов
    const patterns = [
      // useTranslations('module')
      /useTranslations\(['"]([^'"]+)['"]\)/g,
      // t('key') или t('key.subkey')
      /(?:t|tCommon|tAdmin|tAuth|tMarketplace|tStorefront|tCars|tChat|tCart|tRealEstate|tServices|tMap|tMisc)\(['"]([^'"]+)['"]\)/g,
    ];

    // Анализируем какой модуль используется
    let currentModule = null;
    const moduleMatches = content.matchAll(
      /useTranslations\(['"]([^'"]+)['"]\)/g
    );
    for (const match of moduleMatches) {
      currentModule = match[1];
    }

    // Находим все использования ключей
    const keyMatches = content.matchAll(
      /(?:t|tCommon|tAdmin|tAuth|tMarketplace|tStorefront|tCars|tChat|tCart|tRealEstate|tServices|tMap|tMisc)\(['"]([^'"]+)['"]\)/g
    );

    for (const match of keyMatches) {
      const key = match[1];

      // Определяем модуль из имени функции
      let module = currentModule;
      const funcName = match[0].split('(')[0];

      if (funcName === 'tCommon') module = 'common';
      else if (funcName === 'tAdmin') module = 'admin';
      else if (funcName === 'tAuth') module = 'auth';
      else if (funcName === 'tMarketplace') module = 'marketplace';
      else if (funcName === 'tStorefront') module = 'storefront';
      else if (funcName === 'tCars') module = 'cars';
      else if (funcName === 'tChat') module = 'chat';
      else if (funcName === 'tCart') module = 'cart';
      else if (funcName === 'tRealEstate') module = 'realEstate';
      else if (funcName === 'tServices') module = 'services';
      else if (funcName === 'tMap') module = 'map';
      else if (funcName === 'tMisc') module = 'misc';

      if (module) {
        if (!usedKeys.has(module)) {
          usedKeys.set(module, new Set());
        }
        usedKeys.get(module).add(key);
      }
    }
  });
});

// Проверяем какие ключи отсутствуют
usedKeys.forEach((keys, module) => {
  keys.forEach((key) => {
    // Проверяем наличие ключа во всех локалях
    let missing = false;
    locales.forEach((locale) => {
      if (!translations[locale][module]) {
        missing = true;
        return;
      }

      // Проверяем вложенные ключи
      const keyParts = key.split('.');
      let current = translations[locale][module];

      for (const part of keyParts) {
        if (!current || !current[part]) {
          missing = true;
          break;
        }
        current = current[part];
      }
    });

    if (missing) {
      if (!missingKeys.has(module)) {
        missingKeys.set(module, new Set());
      }
      missingKeys.get(module).add(key);
    }
  });
});

// Выводим результаты
console.log('📊 Результаты анализа:\n');

if (missingKeys.size === 0) {
  console.log('✅ Все ключи переводов присутствуют!');
} else {
  console.log('❌ Обнаружены отсутствующие ключи:\n');

  let totalMissing = 0;
  missingKeys.forEach((keys, module) => {
    console.log(`📦 Модуль: ${module}`);
    const sortedKeys = Array.from(keys).sort();
    sortedKeys.forEach((key) => {
      console.log(`   - ${key}`);
      totalMissing++;
    });
    console.log('');
  });

  console.log(`\n🔢 Всего отсутствующих ключей: ${totalMissing}`);

  // Предлагаем автоматическое добавление
  console.log(
    '\n💡 Для автоматического добавления отсутствующих ключей запустите:'
  );
  console.log('   node scripts/add-missing-keys.js');
}

// Сохраняем отчёт для скрипта добавления
const report = {
  missingKeys: Object.fromEntries(
    Array.from(missingKeys.entries()).map(([module, keys]) => [
      module,
      Array.from(keys),
    ])
  ),
  timestamp: new Date().toISOString(),
};

fs.writeFileSync(
  path.join(__dirname, 'missing-keys-report.json'),
  JSON.stringify(report, null, 2)
);

console.log('\n📄 Отчёт сохранён в scripts/missing-keys-report.json');
