#!/usr/bin/env node

/**
 * Скрипт для поиска и анализа placeholder'ов в переводах
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

// Паттерны placeholder'ов
const PLACEHOLDER_PATTERNS = [
  /\[MISSING_[A-Z]+\]/g,
  /\[RU\]/g,
  /\[EN\]/g,
  /\[SR\]/g,
  /\[TODO\]/g,
  /\[TRANSLATE\]/g,
  /XXX/g,
  /TODO:/g,
];

// Получить все значения из объекта рекурсивно
function getAllValues(obj, currentPath = '') {
  let values = [];
  
  for (const key in obj) {
    const fullPath = currentPath ? `${currentPath}.${key}` : key;
    
    if (typeof obj[key] === 'object' && obj[key] !== null && !Array.isArray(obj[key])) {
      values = values.concat(getAllValues(obj[key], fullPath));
    } else if (typeof obj[key] === 'string') {
      values.push({
        path: fullPath,
        value: obj[key]
      });
    }
  }
  
  return values;
}

// Найти placeholder'ы в строке
function findPlaceholders(text) {
  const found = [];
  
  for (const pattern of PLACEHOLDER_PATTERNS) {
    const matches = text.match(pattern);
    if (matches) {
      matches.forEach(match => {
        if (!found.includes(match)) {
          found.push(match);
        }
      });
    }
  }
  
  return found;
}

// Главная функция
async function findAllPlaceholders() {
  console.log(`${colors.cyan}🔍 Поиск placeholder'ов в переводах...${colors.reset}\n`);
  
  const report = {
    totalPlaceholders: 0,
    byLanguage: {},
    byModule: {},
    byType: {},
    details: []
  };
  
  // Инициализация счетчиков
  LANGUAGES.forEach(lang => {
    report.byLanguage[lang] = 0;
  });
  
  // Получаем список всех модулей
  const enFiles = await glob(`${MESSAGES_DIR}/en/*.json`);
  const modules = enFiles.map(file => path.basename(file, '.json'));
  
  console.log(`${colors.blue}📦 Анализирую ${modules.length} модулей...${colors.reset}\n`);
  
  // Обрабатываем каждый модуль
  for (const module of modules) {
    let modulePlaceholders = 0;
    
    for (const lang of LANGUAGES) {
      const filePath = path.join(MESSAGES_DIR, lang, `${module}.json`);
      
      try {
        const content = fs.readFileSync(filePath, 'utf8');
        const translations = JSON.parse(content);
        const values = getAllValues(translations);
        
        for (const { path: keyPath, value } of values) {
          const placeholders = findPlaceholders(value);
          
          if (placeholders.length > 0) {
            modulePlaceholders++;
            report.totalPlaceholders++;
            report.byLanguage[lang]++;
            
            // Считаем по типам
            placeholders.forEach(ph => {
              if (!report.byType[ph]) {
                report.byType[ph] = 0;
              }
              report.byType[ph]++;
            });
            
            // Сохраняем детали
            report.details.push({
              module,
              language: lang,
              key: keyPath,
              value,
              placeholders
            });
          }
        }
      } catch (error) {
        console.log(`${colors.yellow}⚠️  Ошибка чтения: ${filePath}${colors.reset}`);
      }
    }
    
    if (modulePlaceholders > 0) {
      report.byModule[module] = modulePlaceholders;
    }
  }
  
  // Вывод результатов
  console.log(`${colors.cyan}📊 РЕЗУЛЬТАТЫ АНАЛИЗА${colors.reset}`);
  console.log(`${colors.cyan}${'='.repeat(50)}${colors.reset}\n`);
  
  if (report.totalPlaceholders === 0) {
    console.log(`${colors.green}✨ Placeholder'ы не найдены! Все переводы заполнены.${colors.reset}`);
  } else {
    console.log(`${colors.yellow}⚠️  Найдено placeholder'ов: ${report.totalPlaceholders}${colors.reset}\n`);
    
    // По языкам
    console.log(`${colors.blue}По языкам:${colors.reset}`);
    Object.entries(report.byLanguage)
      .sort((a, b) => b[1] - a[1])
      .forEach(([lang, count]) => {
        if (count > 0) {
          const percent = ((count / report.totalPlaceholders) * 100).toFixed(1);
          console.log(`  ${lang}: ${colors.yellow}${count}${colors.reset} (${percent}%)`);
        }
      });
    
    // По типам placeholder'ов
    console.log(`\n${colors.blue}По типам:${colors.reset}`);
    Object.entries(report.byType)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 10)
      .forEach(([type, count]) => {
        console.log(`  ${colors.yellow}${type}${colors.reset}: ${count}`);
      });
    
    // По модулям
    console.log(`\n${colors.blue}Топ модулей с placeholder'ами:${colors.reset}`);
    Object.entries(report.byModule)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 10)
      .forEach(([module, count]) => {
        console.log(`  ${module}: ${colors.yellow}${count}${colors.reset}`);
      });
    
    // Примеры
    console.log(`\n${colors.blue}Примеры найденных placeholder'ов:${colors.reset}`);
    const examples = report.details.slice(0, 5);
    examples.forEach(({ module, language, key, value }) => {
      console.log(`${colors.yellow}📍 ${module}/${language}/${key}:${colors.reset}`);
      console.log(`   "${value.substring(0, 100)}${value.length > 100 ? '...' : ''}"`);
    });
    
    if (report.details.length > 5) {
      console.log(`\n${colors.cyan}... и еще ${report.details.length - 5} placeholder'ов${colors.reset}`);
    }
  }
  
  // Сохраняем детальный отчет
  const reportPath = path.join(__dirname, 'placeholders-report.json');
  fs.writeFileSync(reportPath, JSON.stringify(report, null, 2));
  console.log(`\n${colors.magenta}📄 Детальный отчет сохранен: ${reportPath}${colors.reset}`);
  
  // CSV экспорт для удобства
  if (report.totalPlaceholders > 0) {
    const csvPath = path.join(__dirname, 'placeholders-report.csv');
    const csvHeader = 'Module,Language,Key,Value,Placeholders\n';
    const csvContent = report.details.map(d => 
      `"${d.module}","${d.language}","${d.key}","${d.value.replace(/"/g, '""')}","${d.placeholders.join(', ')}"`
    ).join('\n');
    
    fs.writeFileSync(csvPath, csvHeader + csvContent);
    console.log(`${colors.magenta}📄 CSV отчет сохранен: ${csvPath}${colors.reset}`);
  }
  
  // Рекомендации
  if (report.totalPlaceholders > 0) {
    console.log(`\n${colors.cyan}💡 РЕКОМЕНДАЦИИ${colors.reset}`);
    console.log(`${colors.cyan}${'='.repeat(50)}${colors.reset}`);
    
    if (report.byType['[MISSING_EN]'] || report.byType['[MISSING_RU]'] || report.byType['[MISSING_SR]']) {
      console.log(`\n1. ${colors.yellow}Заполните недостающие переводы:${colors.reset}`);
      console.log(`   - Используйте массовый перевод через AI в админ-панели`);
      console.log(`   - Или заполните вручную критические переводы`);
    }
    
    if (report.byType['[RU]'] || report.byType['[EN]'] || report.byType['[SR]']) {
      console.log(`\n2. ${colors.yellow}Удалите старые placeholder'ы [RU], [EN], [SR]:${colors.reset}`);
      console.log(`   - Это устаревший формат, замените на реальные переводы`);
    }
    
    console.log(`\n3. ${colors.yellow}Приоритетные модули для исправления:${colors.reset}`);
    Object.entries(report.byModule)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 3)
      .forEach(([module, count]) => {
        console.log(`   - ${module} (${count} placeholder'ов)`);
      });
  }
  
  console.log(`\n${colors.green}✨ Анализ завершен!${colors.reset}`);
  
  return report;
}

// Запуск если вызван напрямую
if (require.main === module) {
  findAllPlaceholders()
    .then(report => {
      process.exit(report.totalPlaceholders > 0 ? 1 : 0);
    })
    .catch(error => {
      console.error(`${colors.red}❌ Ошибка:${colors.reset}`, error);
      process.exit(1);
    });
}

module.exports = { findAllPlaceholders };