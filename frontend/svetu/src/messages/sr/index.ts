// Автосгенерированный файл для модульной загрузки переводов
// Сгенерирован: 2025-08-04T08:39:33.520Z
// НЕ РЕДАКТИРУЙТЕ ВРУЧНУЮ!

// Базовые переводы (загружаются всегда)
import common from './common.json';

// Типы модулей
export type TranslationModule = 
  | 'common'
  | 'auth'
  | 'marketplace'
  | 'admin'
  | 'storefront'
  | 'cars'
  | 'cart'
  | 'misc';

// Карта модулей для динамической загрузки
export const moduleLoaders = {
  'auth': () => import('./auth.json'),
  'marketplace': () => import('./marketplace.json'),
  'admin': () => import('./admin.json'),
  'storefront': () => import('./storefront.json'),
  'cars': () => import('./cars.json'),
  'cart': () => import('./cart.json'),
  'misc': () => import('./misc.json')
};

// Функция загрузки модуля
export async function loadModule(moduleName: TranslationModule) {
  if (moduleName === 'common') return common;
  
  const loader = moduleLoaders[moduleName];
  if (!loader) {
    throw new Error(`Unknown module: ${moduleName}`);
  }
  
  const module = await loader();
  return module.default || module;
}

// Экспорт базовых переводов
export default common;
