#!/usr/bin/env node

/**
 * Скрипт для разбиения больших JSON файлов переводов на модули
 * Запуск: node scripts/split-translations.js
 */

const fs = require('fs');
const path = require('path');

// Конфигурация модулей и их ключей
const MODULES = {
  common: [
    'common',
    'countries',
    'navigation',
    'footer',
    'toast',
    'errors',
    'filters',
    'sort',
    'pagination',
    'search',
    'dateTime',
    'validation'
  ],
  auth: [
    'auth',
    'login',
    'register',
    'profile',
    'resetPassword',
    'userMenu',
    'settings'
  ],
  marketplace: [
    'marketplace',
    'listing',
    'listingDetails',
    'listingActions',
    'listingForm',
    'categories',
    'comparison',
    'favorites',
    'priceHistory',
    'seller',
    'buyer',
    'map'
  ],
  admin: [
    'admin',
    'attributes',
    'attributeGroups',
    'variantAttributes',
    'analytics',
    'users',
    'system'
  ],
  storefront: [
    'storefront',
    'store',
    'products',
    'dashboard',
    'orders',
    'inventory',
    'shopSettings'
  ],
  cars: [
    'cars',
    'automotive',
    'carDetails',
    'carFilters',
    'vinDecoder'
  ],
  chat: [
    'chat',
    'messages',
    'conversations',
    'notifications'
  ],
  cart: [
    'cart',
    'checkout',
    'payment',
    'shipping',
    'orderHistory'
  ],
  realEstate: [
    'realEstate',
    'property',
    'propertyDetails',
    'propertyFilters'
  ],
  services: [
    'services',
    'booking',
    'serviceProviders',
    'serviceCategories'
  ]
};

// Функция для извлечения ключей из большого JSON
function extractModuleKeys(fullTranslations, moduleKeys) {
  const moduleData = {};
  
  moduleKeys.forEach(key => {
    if (fullTranslations[key]) {
      moduleData[key] = fullTranslations[key];
    }
  });
  
  return moduleData;
}

// Функция для создания оставшихся переводов (не вошедших в модули)
function getUnusedKeys(fullTranslations, allUsedKeys) {
  const unused = {};
  const usedKeysSet = new Set(allUsedKeys);
  
  Object.keys(fullTranslations).forEach(key => {
    if (!usedKeysSet.has(key)) {
      unused[key] = fullTranslations[key];
    }
  });
  
  return unused;
}

// Главная функция
async function splitTranslations() {
  const languages = ['ru', 'en', 'sr'];
  
  for (const lang of languages) {
    console.log(`\n📦 Обработка языка: ${lang}`);
    
    // Читаем исходный файл
    const sourcePath = path.join(__dirname, '..', 'src', 'messages', `${lang}.json`);
    const fullTranslations = JSON.parse(fs.readFileSync(sourcePath, 'utf8'));
    
    // Создаем директорию для языка
    const langDir = path.join(__dirname, '..', 'src', 'messages', lang);
    if (!fs.existsSync(langDir)) {
      fs.mkdirSync(langDir, { recursive: true });
    }
    
    // Собираем все использованные ключи
    const allUsedKeys = [];
    
    // Создаем модули
    for (const [moduleName, moduleKeys] of Object.entries(MODULES)) {
      const moduleData = extractModuleKeys(fullTranslations, moduleKeys);
      
      if (Object.keys(moduleData).length > 0) {
        const modulePath = path.join(langDir, `${moduleName}.json`);
        fs.writeFileSync(modulePath, JSON.stringify(moduleData, null, 2));
        console.log(`  ✅ Создан модуль ${moduleName}.json (${Object.keys(moduleData).length} ключей)`);
        
        allUsedKeys.push(...moduleKeys);
      }
    }
    
    // Сохраняем неиспользованные ключи в misc.json
    const unusedKeys = getUnusedKeys(fullTranslations, allUsedKeys);
    if (Object.keys(unusedKeys).length > 0) {
      const miscPath = path.join(langDir, 'misc.json');
      fs.writeFileSync(miscPath, JSON.stringify(unusedKeys, null, 2));
      console.log(`  ⚠️  Создан модуль misc.json для неопределенных ключей (${Object.keys(unusedKeys).length} ключей)`);
    }
    
    // Создаем index.ts для динамического импорта
    const indexContent = `// Автосгенерированный файл для lazy loading переводов
// НЕ РЕДАКТИРУЙТЕ ВРУЧНУЮ!

export type TranslationModule = 
  | 'common'
  | 'auth'
  | 'marketplace'
  | 'admin'
  | 'storefront'
  | 'cars'
  | 'chat'
  | 'cart'
  | 'realEstate'
  | 'services'
  | 'misc';

// Базовые переводы, которые загружаются всегда
import common from './common.json';

// Функция для динамической загрузки модулей
export async function loadTranslationModule(module: TranslationModule) {
  switch (module) {
    case 'common':
      return common;
    case 'auth':
      return import('./auth.json');
    case 'marketplace':
      return import('./marketplace.json');
    case 'admin':
      return import('./admin.json');
    case 'storefront':
      return import('./storefront.json');
    case 'cars':
      return import('./cars.json');
    case 'chat':
      return import('./chat.json');
    case 'cart':
      return import('./cart.json');
    case 'realEstate':
      return import('./realEstate.json');
    case 'services':
      return import('./services.json');
    case 'misc':
      return import('./misc.json');
    default:
      throw new Error(\`Unknown translation module: \${module}\`);
  }
}

// Экспорт базовых переводов
export default common;
`;
    
    const indexPath = path.join(langDir, 'index.ts');
    fs.writeFileSync(indexPath, indexContent);
    console.log(`  ✅ Создан index.ts для динамической загрузки`);
  }
  
  console.log('\n✨ Разбиение переводов завершено!');
  console.log('\n📋 Следующие шаги:');
  console.log('1. Обновите i18n.ts для использования модульных переводов');
  console.log('2. Обновите компоненты для использования правильных namespace');
  console.log('3. Протестируйте lazy loading в разных разделах');
}

// Запуск
splitTranslations().catch(console.error);