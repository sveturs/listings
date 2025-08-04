#!/usr/bin/env node

/**
 * Скрипт миграции переводов из misc.json в правильные модули
 * Анализирует структуру и автоматически распределяет переводы
 */

const fs = require('fs');
const path = require('path');

// Маппинг корневых ключей к модулям
const KEY_TO_MODULE_MAP = {
  // Common module
  'common': 'common',
  'countries': 'common',
  'languages': 'common',
  'errors': 'common',
  'validation': 'common',
  'dateTime': 'common',
  'toast': 'common',
  'filters': 'common',
  'sort': 'common',
  'pagination': 'common',
  'search': 'common',
  'navigation': 'common',
  'footer': 'common',
  
  // Auth module
  'auth': 'auth',
  'login': 'auth',
  'register': 'auth',
  'profile': 'auth',
  'resetPassword': 'auth',
  'userMenu': 'auth',
  'settings': 'auth',
  
  // Marketplace module
  'marketplace': 'marketplace',
  'listing': 'marketplace',
  'listingDetails': 'marketplace',
  'listingActions': 'marketplace',
  'listingForm': 'marketplace',
  'categories': 'marketplace',
  'comparison': 'marketplace',
  'favorites': 'marketplace',
  'priceHistory': 'marketplace',
  'seller': 'marketplace',
  'buyer': 'marketplace',
  'map': 'marketplace',
  'home': 'marketplace', // Главная страница маркетплейса
  
  // Admin module
  'admin': 'admin',
  'attributes': 'admin',
  'attributeGroups': 'admin',
  'variantAttributes': 'admin',
  'analytics': 'admin',
  'moderation': 'admin',
  
  // Storefront module
  'storefront': 'storefront',
  'store': 'storefront',
  'products': 'storefront',
  'dashboard': 'storefront',
  'orders': 'storefront',
  'inventory': 'storefront',
  'shopSettings': 'storefront',
  
  // Cars module
  'cars': 'cars',
  'automotive': 'cars',
  'carDetails': 'cars',
  'carFilters': 'cars',
  'vinDecoder': 'cars',
  
  // Chat module
  'chat': 'chat',
  'messages': 'chat',
  'conversations': 'chat',
  'notifications': 'chat',
  
  // Cart module
  'cart': 'cart',
  'checkout': 'cart',
  'payment': 'cart',
  'shipping': 'cart',
  'orderHistory': 'cart',
  
  // Real Estate module
  'realEstate': 'realEstate',
  'property': 'realEstate',
  'propertyDetails': 'realEstate',
  'propertyFilters': 'realEstate',
  
  // Services module
  'services': 'services',
  'booking': 'services',
  'serviceProviders': 'services',
  'serviceCategories': 'services',
  
  // Special cases
  'header': 'common', // Header в common для доступности везде
  'blog': 'misc', // Блог остается в misc
  'news': 'misc', // Новости остаются в misc
};

// Функция для глубокого слияния объектов
function deepMerge(target, source) {
  const output = Object.assign({}, target);
  if (isObject(target) && isObject(source)) {
    Object.keys(source).forEach(key => {
      if (isObject(source[key])) {
        if (!(key in target))
          Object.assign(output, { [key]: source[key] });
        else
          output[key] = deepMerge(target[key], source[key]);
      } else {
        Object.assign(output, { [key]: source[key] });
      }
    });
  }
  return output;
}

function isObject(item) {
  return item && typeof item === 'object' && !Array.isArray(item);
}

// Функция миграции для одного языка
async function migrateLanguage(lang) {
  console.log(`\n🔄 Миграция языка: ${lang}`);
  console.log('─'.repeat(50));
  
  const langDir = path.join(__dirname, '..', 'src', 'messages', lang);
  const miscPath = path.join(langDir, 'misc.json');
  
  // Проверяем существование misc.json
  if (!fs.existsSync(miscPath)) {
    console.log('  ℹ️  misc.json не найден, пропускаем');
    return;
  }
  
  // Читаем misc.json
  const miscData = JSON.parse(fs.readFileSync(miscPath, 'utf8'));
  const totalKeys = Object.keys(miscData).length;
  console.log(`  📊 Найдено ${totalKeys} ключей в misc.json`);
  
  // Распределяем переводы по модулям
  const moduleUpdates = {};
  const remainingMisc = {};
  let migratedCount = 0;
  
  Object.entries(miscData).forEach(([key, value]) => {
    const targetModule = KEY_TO_MODULE_MAP[key];
    
    if (targetModule && targetModule !== 'misc') {
      if (!moduleUpdates[targetModule]) {
        moduleUpdates[targetModule] = {};
      }
      moduleUpdates[targetModule][key] = value;
      migratedCount++;
    } else {
      remainingMisc[key] = value;
    }
  });
  
  // Обновляем модули
  for (const [moduleName, updates] of Object.entries(moduleUpdates)) {
    const modulePath = path.join(langDir, `${moduleName}.json`);
    let moduleData = {};
    
    // Читаем существующий модуль
    if (fs.existsSync(modulePath)) {
      moduleData = JSON.parse(fs.readFileSync(modulePath, 'utf8'));
    }
    
    // Сливаем с новыми данными
    const mergedData = deepMerge(moduleData, updates);
    
    // Сохраняем обновленный модуль
    fs.writeFileSync(modulePath, JSON.stringify(mergedData, null, 2));
    console.log(`  ✅ Обновлен ${moduleName}.json (+${Object.keys(updates).length} ключей)`);
  }
  
  // Обновляем misc.json с оставшимися ключами
  if (Object.keys(remainingMisc).length > 0) {
    fs.writeFileSync(miscPath, JSON.stringify(remainingMisc, null, 2));
    console.log(`  ⚠️  Осталось в misc.json: ${Object.keys(remainingMisc).length} ключей`);
  } else {
    fs.unlinkSync(miscPath);
    console.log(`  🗑️  misc.json удален (все ключи мигрированы)`);
  }
  
  console.log(`  ✨ Мигрировано: ${migratedCount} из ${totalKeys} ключей`);
}

// Функция для обновления index.ts файлов
function updateIndexFiles() {
  const languages = ['ru', 'en', 'sr'];
  
  languages.forEach(lang => {
    const langDir = path.join(__dirname, '..', 'src', 'messages', lang);
    const indexPath = path.join(langDir, 'index.ts');
    
    // Получаем список модулей
    const modules = fs.readdirSync(langDir)
      .filter(file => file.endsWith('.json'))
      .map(file => file.replace('.json', ''));
    
    // Генерируем новый index.ts
    const indexContent = `// Автосгенерированный файл для модульной загрузки переводов
// Сгенерирован: ${new Date().toISOString()}
// НЕ РЕДАКТИРУЙТЕ ВРУЧНУЮ!

// Базовые переводы (загружаются всегда)
import common from './common.json';

// Типы модулей
export type TranslationModule = 
${modules.map(m => `  | '${m}'`).join('\n')};

// Карта модулей для динамической загрузки
export const moduleLoaders = {
${modules.filter(m => m !== 'common').map(m => 
  `  '${m}': () => import('./${m}.json')`
).join(',\n')}
};

// Функция загрузки модуля
export async function loadModule(moduleName: TranslationModule) {
  if (moduleName === 'common') return common;
  
  const loader = moduleLoaders[moduleName];
  if (!loader) {
    throw new Error(\`Unknown module: \${moduleName}\`);
  }
  
  const module = await loader();
  return module.default || module;
}

// Экспорт базовых переводов
export default common;
`;
    
    fs.writeFileSync(indexPath, indexContent);
    console.log(`\n✅ Обновлен ${lang}/index.ts`);
  });
}

// Главная функция
async function main() {
  console.log('🚀 Начинаем миграцию переводов из misc.json\n');
  
  const languages = ['ru', 'en', 'sr'];
  
  // Мигрируем каждый язык
  for (const lang of languages) {
    await migrateLanguage(lang);
  }
  
  // Обновляем index файлы
  updateIndexFiles();
  
  console.log('\n✨ Миграция завершена!');
  console.log('\n📋 Следующие шаги:');
  console.log('1. Проверьте оставшиеся ключи в misc.json файлах');
  console.log('2. Обновите KEY_TO_MODULE_MAP для нераспределенных ключей');
  console.log('3. Запустите скрипт повторно при необходимости');
}

// Запуск
main().catch(console.error);