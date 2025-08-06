#!/usr/bin/env node

/**
 * Улучшенный скрипт для разбиения больших JSON файлов переводов на модули
 * Версия 2.0 с более точной категоризацией
 */

const fs = require('fs');
const path = require('path');

// Расширенная конфигурация модулей с поддержкой вложенных ключей
const MODULES_CONFIG = {
  common: {
    // Прямые ключи
    directKeys: [
      'common',
      'countries',
      'navigation',
      'footer',
      'errors',
      'filters',
      'sort',
      'pagination',
      'search',
      'dateTime',
      'validation',
      'notifications',
      'languages',
      'cookieConsent',
      'privacyPolicy',
      'termsOfService',
      'support',
      'contact',
      'about',
      'faq',
    ],
    // Паттерны для поиска
    patterns: [
      /^common\./,
      /^errors\./,
      /^validation\./,
      /^dateTime\./,
      /^languages\./,
    ],
  },

  auth: {
    directKeys: [
      'auth',
      'login',
      'register',
      'profile',
      'resetPassword',
      'userMenu',
      'settings',
      'accountSettings',
      'security',
      'twoFactor',
      'sessions',
      'apiKeys',
    ],
    patterns: [
      /^auth\./,
      /^login\./,
      /^register\./,
      /^profile\./,
      /^settings\./,
    ],
  },

  marketplace: {
    directKeys: [
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
      'map',
      'reviews',
      'ratings',
      'questions',
      'offers',
      'negotiation',
      'reportListing',
      'listingStatistics',
    ],
    patterns: [
      /^marketplace\./,
      /^listing\./,
      /^categories\./,
      /^map\./,
      /^seller\./,
      /^buyer\./,
    ],
  },

  admin: {
    directKeys: [
      'admin',
      'adminDashboard',
      'adminUsers',
      'adminListings',
      'adminCategories',
      'adminAttributes',
      'adminAnalytics',
      'adminReports',
      'adminSettings',
      'adminLogs',
      'adminAudit',
      'adminPermissions',
      'adminRoles',
      'moderation',
      'contentModeration',
    ],
    patterns: [
      /^admin\./,
      /^moderation\./,
      /^attributes\./,
      /^attributeGroups\./,
      /^variantAttributes\./,
    ],
  },

  storefront: {
    directKeys: [
      'storefront',
      'store',
      'products',
      'dashboard',
      'orders',
      'inventory',
      'shopSettings',
      'storeAnalytics',
      'productVariants',
      'collections',
      'discounts',
      'coupons',
      'salesReports',
      'customerManagement',
    ],
    patterns: [
      /^storefront\./,
      /^store\./,
      /^products\./,
      /^dashboard\./,
      /^inventory\./,
    ],
  },

  cars: {
    directKeys: [
      'cars',
      'automotive',
      'carDetails',
      'carFilters',
      'vinDecoder',
      'carMakes',
      'carModels',
      'carFeatures',
      'carSpecifications',
      'carHistory',
      'carInspection',
      'carFinancing',
      'carInsurance',
    ],
    patterns: [/^cars\./, /^automotive\./, /^vinDecoder\./, /^carFilters\./],
  },

  chat: {
    directKeys: [
      'chat',
      'messages',
      'conversations',
      'chatNotifications',
      'typing',
      'messageStatus',
      'chatSettings',
      'chatActions',
      'attachments',
      'voiceMessages',
      'videoCall',
    ],
    patterns: [/^chat\./, /^messages\./, /^conversations\./],
  },

  cart: {
    directKeys: [
      'cart',
      'checkout',
      'payment',
      'shipping',
      'orderHistory',
      'orderStatus',
      'orderTracking',
      'paymentMethods',
      'shippingMethods',
      'orderSummary',
      'invoice',
      'refunds',
    ],
    patterns: [
      /^cart\./,
      /^checkout\./,
      /^payment\./,
      /^shipping\./,
      /^order\./,
    ],
  },

  realEstate: {
    directKeys: [
      'realEstate',
      'property',
      'propertyDetails',
      'propertyFilters',
      'propertyTypes',
      'propertyFeatures',
      'propertyLocation',
      'propertyPricing',
      'propertyDocuments',
      'propertyTours',
      'mortgage',
    ],
    patterns: [/^realEstate\./, /^property\./, /^mortgage\./],
  },

  services: {
    directKeys: [
      'services',
      'booking',
      'serviceProviders',
      'serviceCategories',
      'serviceDetails',
      'serviceReviews',
      'serviceSchedule',
      'servicePricing',
      'serviceBooking',
      'serviceCalendar',
    ],
    patterns: [/^services\./, /^booking\./, /^serviceProviders\./],
  },
};

// Функция для проверки, принадлежит ли ключ модулю
function belongsToModule(key, moduleConfig) {
  // Проверяем прямые ключи
  if (moduleConfig.directKeys.includes(key)) {
    return true;
  }

  // Проверяем паттерны
  return moduleConfig.patterns.some((pattern) => pattern.test(key));
}

// Функция для извлечения всех ключей модуля (включая вложенные)
function extractModuleData(fullTranslations, moduleConfig, processedKeys) {
  const moduleData = {};

  Object.entries(fullTranslations).forEach(([key, value]) => {
    if (!processedKeys.has(key) && belongsToModule(key, moduleConfig)) {
      moduleData[key] = value;
      processedKeys.add(key);
    }
  });

  return moduleData;
}

// Функция для создания статистики
function generateStats(modules) {
  const stats = {};
  let totalKeys = 0;

  Object.entries(modules).forEach(([moduleName, moduleData]) => {
    const keyCount = Object.keys(moduleData).length;
    stats[moduleName] = keyCount;
    totalKeys += keyCount;
  });

  return { stats, totalKeys };
}

// Главная функция
async function splitTranslations() {
  const languages = ['ru', 'en', 'sr'];

  for (const lang of languages) {
    console.log(`\n📦 Обработка языка: ${lang}`);
    console.log('─'.repeat(50));

    // Читаем исходный файл
    const sourcePath = path.join(
      __dirname,
      '..',
      'src',
      'messages',
      `${lang}.json`
    );
    const fullTranslations = JSON.parse(fs.readFileSync(sourcePath, 'utf8'));
    const totalOriginalKeys = Object.keys(fullTranslations).length;

    console.log(`📊 Всего ключей в исходном файле: ${totalOriginalKeys}`);

    // Создаем директорию для языка
    const langDir = path.join(__dirname, '..', 'src', 'messages', lang);

    // Очищаем директорию
    if (fs.existsSync(langDir)) {
      fs.rmSync(langDir, { recursive: true });
    }
    fs.mkdirSync(langDir, { recursive: true });

    // Отслеживаем обработанные ключи
    const processedKeys = new Set();
    const modules = {};

    // Создаем модули в определенном порядке приоритета
    for (const [moduleName, moduleConfig] of Object.entries(MODULES_CONFIG)) {
      const moduleData = extractModuleData(
        fullTranslations,
        moduleConfig,
        processedKeys
      );

      if (Object.keys(moduleData).length > 0) {
        modules[moduleName] = moduleData;
        const modulePath = path.join(langDir, `${moduleName}.json`);
        fs.writeFileSync(modulePath, JSON.stringify(moduleData, null, 2));
        console.log(
          `  ✅ ${moduleName}.json - ${Object.keys(moduleData).length} ключей`
        );
      }
    }

    // Собираем неиспользованные ключи
    const miscData = {};
    Object.entries(fullTranslations).forEach(([key, value]) => {
      if (!processedKeys.has(key)) {
        miscData[key] = value;
      }
    });

    if (Object.keys(miscData).length > 0) {
      modules.misc = miscData;
      const miscPath = path.join(langDir, 'misc.json');
      fs.writeFileSync(miscPath, JSON.stringify(miscData, null, 2));
      console.log(
        `  ⚠️  misc.json - ${Object.keys(miscData).length} нераспределенных ключей`
      );
    }

    // Генерируем статистику
    const { stats, totalKeys } = generateStats(modules);
    console.log(`\n📈 Статистика распределения:`);
    Object.entries(stats).forEach(([module, count]) => {
      const percentage = ((count / totalOriginalKeys) * 100).toFixed(1);
      console.log(`   ${module}: ${count} ключей (${percentage}%)`);
    });
    console.log(`   ─────────────────────────`);
    console.log(`   Всего обработано: ${totalKeys} из ${totalOriginalKeys}`);

    // Создаем index.ts для языка
    const indexContent = `// Автосгенерированный файл для модульной загрузки переводов
// Сгенерирован: ${new Date().toISOString()}
// НЕ РЕДАКТИРУЙТЕ ВРУЧНУЮ!

// Базовые переводы (загружаются всегда)
import common from './common.json';

// Типы модулей
export type TranslationModule = 
${Object.keys(modules)
  .map((m) => `  | '${m}'`)
  .join('\n')};

// Карта модулей для динамической загрузки
export const moduleLoaders = {
${Object.keys(modules)
  .filter((m) => m !== 'common')
  .map((m) => `  '${m}': () => import('./${m}.json')`)
  .join(',\n')}
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

    const indexPath = path.join(langDir, 'index.ts');
    fs.writeFileSync(indexPath, indexContent);
    console.log(`\n  ✅ Создан index.ts для динамической загрузки`);
  }

  console.log('\n✨ Разбиение переводов завершено!');
  console.log('\n📋 Рекомендации:');
  console.log('1. Проверьте misc.json файлы и распределите оставшиеся ключи');
  console.log('2. Обновите i18n.ts для использования модульной системы');
  console.log('3. Обновите компоненты для правильных namespace');
  console.log('4. Запустите yarn build для проверки');
}

// Запуск
splitTranslations().catch(console.error);
