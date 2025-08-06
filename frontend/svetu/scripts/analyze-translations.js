#!/usr/bin/env node

/**
 * Скрипт анализа размера переводов и потенциальной экономии
 */

const fs = require('fs');
const path = require('path');
const { gzipSync } = require('zlib');

// Цвета для консоли
const colors = {
  reset: '\x1b[0m',
  bright: '\x1b[1m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[36m',
  red: '\x1b[31m',
};

// Форматирование размера
function formatSize(bytes) {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
  return (bytes / 1024 / 1024).toFixed(2) + ' MB';
}

// Получение размера файла и gzip размера
function getFileSizes(filePath) {
  const content = fs.readFileSync(filePath, 'utf8');
  const raw = Buffer.byteLength(content, 'utf8');
  const gzipped = gzipSync(content).length;

  return { raw, gzipped, content };
}

// Анализ модульной системы
function analyzeModularSystem(lang) {
  const langDir = path.join(__dirname, '..', 'src', 'messages', lang);
  const modules = {};
  let totalRaw = 0;
  let totalGzipped = 0;

  // Анализируем каждый модуль
  const files = fs.readdirSync(langDir).filter((f) => f.endsWith('.json'));

  files.forEach((file) => {
    const filePath = path.join(langDir, file);
    const sizes = getFileSizes(filePath);
    const keyCount = Object.keys(JSON.parse(sizes.content)).length;

    modules[file] = {
      ...sizes,
      keyCount,
    };

    totalRaw += sizes.raw;
    totalGzipped += sizes.gzipped;
  });

  return { modules, totalRaw, totalGzipped };
}

// Анализ старой системы
function analyzeOldSystem(lang) {
  const filePath = path.join(
    __dirname,
    '..',
    'src',
    'messages',
    `${lang}.json`
  );

  if (!fs.existsSync(filePath)) {
    return null;
  }

  const sizes = getFileSizes(filePath);
  const keyCount = Object.keys(JSON.parse(sizes.content)).length;

  return { ...sizes, keyCount };
}

// Симуляция загрузки страниц
function simulatePageLoads() {
  const scenarios = [
    { name: 'Главная страница', modules: ['common', 'marketplace'] },
    { name: 'Админ панель', modules: ['common', 'admin'] },
    { name: 'Корзина', modules: ['common', 'cart'] },
    { name: 'Витрина магазина', modules: ['common', 'storefront'] },
    { name: 'Автомобили', modules: ['common', 'marketplace', 'cars'] },
  ];

  return scenarios;
}

// Главная функция
async function main() {
  console.log(`${colors.bright}📊 Анализ размера переводов${colors.reset}\n`);

  const languages = ['ru', 'en', 'sr'];
  const allStats = {};

  // Анализируем каждый язык
  for (const lang of languages) {
    console.log(`${colors.blue}Язык: ${lang}${colors.reset}`);
    console.log('─'.repeat(60));

    // Старая система
    const oldSystem = analyzeOldSystem(lang);
    if (oldSystem) {
      console.log(
        `\n${colors.yellow}Старая система (один файл):${colors.reset}`
      );
      console.log(
        `  Размер: ${formatSize(oldSystem.raw)} (gzip: ${formatSize(oldSystem.gzipped)})`
      );
      console.log(`  Ключей: ${oldSystem.keyCount}`);
    }

    // Новая модульная система
    const modularSystem = analyzeModularSystem(lang);
    console.log(`\n${colors.green}Модульная система:${colors.reset}`);

    // Сортируем модули по размеру
    const sortedModules = Object.entries(modularSystem.modules).sort(
      (a, b) => b[1].raw - a[1].raw
    );

    sortedModules.forEach(([name, data]) => {
      const percentage = ((data.raw / modularSystem.totalRaw) * 100).toFixed(1);
      console.log(
        `  ${name.padEnd(20)} ${formatSize(data.raw).padEnd(10)} (${percentage}%) - ${data.keyCount} ключей`
      );
    });

    console.log(`  ${'─'.repeat(50)}`);
    console.log(
      `  ${'ИТОГО:'.padEnd(20)} ${formatSize(modularSystem.totalRaw).padEnd(10)} (gzip: ${formatSize(modularSystem.totalGzipped)})`
    );

    // Сравнение
    if (oldSystem) {
      const savings = oldSystem.raw - modularSystem.totalRaw;
      const savingsPercent = ((savings / oldSystem.raw) * 100).toFixed(1);
      console.log(
        `\n${colors.bright}Изменение размера:${colors.reset} ${savings > 0 ? colors.green : colors.red}${savingsPercent}%${colors.reset}`
      );
    }

    allStats[lang] = { oldSystem, modularSystem };
    console.log('\n');
  }

  // Симуляция загрузки страниц
  console.log(`${colors.bright}📱 Симуляция загрузки страниц${colors.reset}\n`);

  const scenarios = simulatePageLoads();
  const exampleLang = 'ru';
  const langStats = allStats[exampleLang];

  if (langStats.modularSystem) {
    scenarios.forEach((scenario) => {
      let pageSize = 0;
      let pageGzipSize = 0;

      scenario.modules.forEach((moduleName) => {
        const module = langStats.modularSystem.modules[`${moduleName}.json`];
        if (module) {
          pageSize += module.raw;
          pageGzipSize += module.gzipped;
        }
      });

      const oldSize = langStats.oldSystem ? langStats.oldSystem.gzipped : 0;
      const savings = oldSize - pageGzipSize;
      const savingsPercent = oldSize
        ? ((savings / oldSize) * 100).toFixed(1)
        : 0;

      console.log(`${scenario.name}:`);
      console.log(`  Модули: ${scenario.modules.join(', ')}`);
      console.log(`  Размер: ${formatSize(pageGzipSize)} (gzip)`);
      if (oldSize) {
        console.log(
          `  Экономия: ${colors.green}${formatSize(savings)} (${savingsPercent}%)${colors.reset}`
        );
      }
      console.log('');
    });
  }

  // Рекомендации
  console.log(`${colors.bright}💡 Рекомендации:${colors.reset}\n`);
  console.log(
    '1. Перенесите оставшиеся переводы из misc.json в соответствующие модули'
  );
  console.log('2. Используйте предзагрузку для критичных модулей');
  console.log('3. Настройте HTTP/2 Server Push для common.json');
  console.log('4. Добавьте долгосрочное кэширование для модулей переводов');
  console.log(
    '5. Рассмотрите использование Service Worker для offline доступа'
  );
}

// Запуск
main().catch(console.error);
