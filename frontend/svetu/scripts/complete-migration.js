#!/usr/bin/env node

/**
 * Полная миграция всех переводов из misc.json в соответствующие модули
 * с анализом структуры и автоматическим распределением
 */

const fs = require('fs');
const path = require('path');

// Расширенный маппинг с учетом вложенных структур
const COMPLETE_MAPPING = {
  // Common module - базовые элементы UI
  common: {
    module: 'common',
    subkeys: ['*'], // Все подключи
  },
  countries: {
    module: 'common',
    subkeys: ['*'],
  },
  languages: {
    module: 'common',
    subkeys: ['*'],
  },
  errors: {
    module: 'common',
    subkeys: ['*'],
  },
  validation: {
    module: 'common',
    subkeys: ['*'],
  },
  dateTime: {
    module: 'common',
    subkeys: ['*'],
  },
  toast: {
    module: 'common',
    subkeys: ['*'],
  },
  footer: {
    module: 'common',
    subkeys: ['*'],
  },
  header: {
    module: 'common',
    subkeys: ['*'],
  },
  navigation: {
    module: 'common',
    subkeys: ['*'],
  },

  // Auth module
  auth: {
    module: 'auth',
    subkeys: ['*'],
  },
  login: {
    module: 'auth',
    subkeys: ['*'],
  },
  register: {
    module: 'auth',
    subkeys: ['*'],
  },
  profile: {
    module: 'auth',
    subkeys: ['*'],
  },
  resetPassword: {
    module: 'auth',
    subkeys: ['*'],
  },
  userMenu: {
    module: 'auth',
    subkeys: ['*'],
  },
  settings: {
    module: 'auth',
    subkeys: ['*'],
  },

  // Marketplace module
  home: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  marketplace: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  listing: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  listingDetails: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  listingActions: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  listingForm: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  categories: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  comparison: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  favorites: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  priceHistory: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  seller: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  buyer: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  map: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  filters: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  sort: {
    module: 'marketplace',
    subkeys: ['*'],
  },
  search: {
    module: 'marketplace',
    subkeys: ['*'],
  },

  // Admin module
  admin: {
    module: 'admin',
    subkeys: ['*'],
  },
  attributes: {
    module: 'admin',
    subkeys: ['*'],
  },
  attributeGroups: {
    module: 'admin',
    subkeys: ['*'],
  },
  variantAttributes: {
    module: 'admin',
    subkeys: ['*'],
  },
  analytics: {
    module: 'admin',
    subkeys: ['*'],
  },
  moderation: {
    module: 'admin',
    subkeys: ['*'],
  },

  // Storefront module
  storefront: {
    module: 'storefront',
    subkeys: ['*'],
  },
  store: {
    module: 'storefront',
    subkeys: ['*'],
  },
  products: {
    module: 'storefront',
    subkeys: ['*'],
  },
  dashboard: {
    module: 'storefront',
    subkeys: ['*'],
  },
  orders: {
    module: 'storefront',
    subkeys: ['*'],
  },
  inventory: {
    module: 'storefront',
    subkeys: ['*'],
  },
  shopSettings: {
    module: 'storefront',
    subkeys: ['*'],
  },

  // Cars module
  cars: {
    module: 'cars',
    subkeys: ['*'],
  },
  automotive: {
    module: 'cars',
    subkeys: ['*'],
  },
  carDetails: {
    module: 'cars',
    subkeys: ['*'],
  },
  carFilters: {
    module: 'cars',
    subkeys: ['*'],
  },
  vinDecoder: {
    module: 'cars',
    subkeys: ['*'],
  },

  // Chat module
  chat: {
    module: 'chat',
    subkeys: ['*'],
  },
  messages: {
    module: 'chat',
    subkeys: ['*'],
  },
  conversations: {
    module: 'chat',
    subkeys: ['*'],
  },
  notifications: {
    module: 'chat',
    subkeys: ['*'],
  },

  // Cart module
  cart: {
    module: 'cart',
    subkeys: ['*'],
  },
  checkout: {
    module: 'cart',
    subkeys: ['*'],
  },
  payment: {
    module: 'cart',
    subkeys: ['*'],
  },
  shipping: {
    module: 'cart',
    subkeys: ['*'],
  },
  orderHistory: {
    module: 'cart',
    subkeys: ['*'],
  },

  // Real Estate module
  realEstate: {
    module: 'realEstate',
    subkeys: ['*'],
  },
  property: {
    module: 'realEstate',
    subkeys: ['*'],
  },
  propertyDetails: {
    module: 'realEstate',
    subkeys: ['*'],
  },
  propertyFilters: {
    module: 'realEstate',
    subkeys: ['*'],
  },

  // Services module
  services: {
    module: 'services',
    subkeys: ['*'],
  },
  booking: {
    module: 'services',
    subkeys: ['*'],
  },
  serviceProviders: {
    module: 'services',
    subkeys: ['*'],
  },
  serviceCategories: {
    module: 'services',
    subkeys: ['*'],
  },

  // Special cases - остаются в misc
  blog: {
    module: 'misc',
    subkeys: ['*'],
  },
  news: {
    module: 'misc',
    subkeys: ['*'],
  },
  metadata: {
    module: 'misc',
    subkeys: ['*'],
  },
};

// Функция для глубокого слияния объектов
function deepMerge(target, source) {
  const output = Object.assign({}, target);
  if (isObject(target) && isObject(source)) {
    Object.keys(source).forEach((key) => {
      if (isObject(source[key])) {
        if (!(key in target)) Object.assign(output, { [key]: source[key] });
        else output[key] = deepMerge(target[key], source[key]);
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

// Функция для создания или обновления модуля
function ensureModule(langDir, moduleName) {
  const modulePath = path.join(langDir, `${moduleName}.json`);

  if (!fs.existsSync(modulePath)) {
    fs.writeFileSync(modulePath, '{}');
  }
}

// Главная функция миграции
async function completeMigration() {
  console.log('🚀 Начинаем полную миграцию переводов\n');

  const languages = ['ru', 'en', 'sr'];
  const globalStats = {
    totalMigrated: 0,
    totalRemaining: 0,
    byModule: {},
  };

  for (const lang of languages) {
    console.log(`\n📦 Обработка языка: ${lang}`);
    console.log('─'.repeat(60));

    const langDir = path.join(__dirname, '..', 'src', 'messages', lang);
    const miscPath = path.join(langDir, 'misc.json');

    if (!fs.existsSync(miscPath)) {
      console.log('  ℹ️  misc.json не найден');
      continue;
    }

    // Читаем misc.json
    const miscData = JSON.parse(fs.readFileSync(miscPath, 'utf8'));
    const totalKeys = Object.keys(miscData).length;
    console.log(`  📊 Найдено ${totalKeys} ключей в misc.json`);

    // Создаем недостающие модули
    const requiredModules = new Set(
      Object.values(COMPLETE_MAPPING).map((m) => m.module)
    );
    requiredModules.forEach((module) => {
      if (module !== 'misc') {
        ensureModule(langDir, module);
      }
    });

    // Распределяем переводы
    const moduleUpdates = {};
    const remainingMisc = {};
    let migratedCount = 0;

    // Обрабатываем каждый ключ
    Object.entries(miscData).forEach(([key, value]) => {
      const mapping = COMPLETE_MAPPING[key];

      if (mapping && mapping.module !== 'misc') {
        const targetModule = mapping.module;

        if (!moduleUpdates[targetModule]) {
          moduleUpdates[targetModule] = {};
        }

        // Добавляем весь объект с подключами
        moduleUpdates[targetModule][key] = value;
        migratedCount++;

        // Обновляем статистику
        if (!globalStats.byModule[targetModule]) {
          globalStats.byModule[targetModule] = 0;
        }
        globalStats.byModule[targetModule]++;
      } else {
        remainingMisc[key] = value;
      }
    });

    // Применяем обновления к модулям
    for (const [moduleName, updates] of Object.entries(moduleUpdates)) {
      const modulePath = path.join(langDir, `${moduleName}.json`);

      // Читаем существующий модуль
      let moduleData = {};
      if (fs.existsSync(modulePath)) {
        moduleData = JSON.parse(fs.readFileSync(modulePath, 'utf8'));
      }

      // Сливаем с новыми данными
      const mergedData = deepMerge(moduleData, updates);

      // Сохраняем
      fs.writeFileSync(modulePath, JSON.stringify(mergedData, null, 2));
      console.log(
        `  ✅ Обновлен ${moduleName}.json (+${Object.keys(updates).length} ключей)`
      );
    }

    // Обновляем misc.json
    if (Object.keys(remainingMisc).length > 0) {
      fs.writeFileSync(miscPath, JSON.stringify(remainingMisc, null, 2));
      console.log(
        `  ⚠️  Осталось в misc.json: ${Object.keys(remainingMisc).length} ключей`
      );
    } else {
      // Удаляем пустой misc.json
      fs.unlinkSync(miscPath);
      console.log(`  🗑️  misc.json удален (все ключи мигрированы)`);
    }

    console.log(`  ✨ Мигрировано: ${migratedCount} из ${totalKeys} ключей`);

    globalStats.totalMigrated += migratedCount;
    globalStats.totalRemaining += Object.keys(remainingMisc).length;
  }

  // Обновляем index файлы
  updateIndexFiles();

  // Выводим финальную статистику
  console.log('\n' + '='.repeat(60));
  console.log('📊 ФИНАЛЬНАЯ СТАТИСТИКА');
  console.log('='.repeat(60));
  console.log(`\n✅ Всего мигрировано: ${globalStats.totalMigrated} ключей`);
  console.log(`⚠️  Осталось в misc: ${globalStats.totalRemaining} ключей`);
  console.log(`\n📦 Распределение по модулям:`);

  Object.entries(globalStats.byModule)
    .sort((a, b) => b[1] - a[1])
    .forEach(([module, count]) => {
      console.log(`   ${module.padEnd(15)} ${count} ключей`);
    });

  console.log('\n✨ Миграция завершена!');

  if (globalStats.totalRemaining > 0) {
    console.log(
      '\n⚠️  Оставшиеся ключи в misc.json требуют ручного распределения'
    );
  }
}

// Функция обновления index файлов
function updateIndexFiles() {
  const languages = ['ru', 'en', 'sr'];

  languages.forEach((lang) => {
    const langDir = path.join(__dirname, '..', 'src', 'messages', lang);
    const indexPath = path.join(langDir, 'index.ts');

    // Получаем список модулей
    const modules = fs
      .readdirSync(langDir)
      .filter((file) => file.endsWith('.json'))
      .map((file) => file.replace('.json', ''));

    // Генерируем index.ts
    const indexContent = `// Автосгенерированный файл для модульной загрузки переводов
// Сгенерирован: ${new Date().toISOString()}
// НЕ РЕДАКТИРУЙТЕ ВРУЧНУЮ!

// Базовые переводы (загружаются всегда)
import common from './common.json';

// Типы модулей
export type TranslationModule = 
${modules.map((m) => `  | '${m}'`).join('\n')};

// Карта модулей для динамической загрузки
export const moduleLoaders = {
${modules
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

    fs.writeFileSync(indexPath, indexContent);
  });

  console.log('\n✅ Обновлены index.ts файлы для всех языков');
}

// Запуск
completeMigration().catch(console.error);
