// Автосгенерированный файл для модульной загрузки переводов
// Сгенерирован: 2025-08-04T10:32:06.116Z
// НЕ РЕДАКТИРУЙТЕ ВРУЧНУЮ!

// Базовые переводы (загружаются всегда)
import common from './common.json';

// Типы модулей
export type TranslationModule =
  | 'admin'
  | 'auth'
  | 'cars'
  | 'cart'
  | 'chat'
  | 'common'
  | 'map'
  | 'marketplace'
  | 'misc'
  | 'realEstate'
  | 'services'
  | 'storefronts';

// Карта модулей для динамической загрузки
export const moduleLoaders = {
  admin: () => import('./admin.json'),
  auth: () => import('./auth.json'),
  cars: () => import('./cars.json'),
  cart: () => import('./cart.json'),
  chat: () => import('./chat.json'),
  map: () => import('./map.json'),
  marketplace: () => import('./marketplace.json'),
  misc: () => import('./misc.json'),
  realEstate: () => import('./realEstate.json'),
  services: () => import('./services.json'),
  storefronts: () => import('./storefronts.json'),
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
