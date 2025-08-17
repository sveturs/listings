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
  | 'checkout'
  | 'common'
  | 'delivery'
  | 'map'
  | 'marketplace'
  | 'misc'
  | 'notifications'
  | 'orders'
  | 'products'
  | 'profile'
  | 'realEstate'
  | 'reviews'
  | 'search'
  | 'services'
  | 'storefronts';

// Карта модулей для динамической загрузки
export const moduleLoaders = {
  admin: () => import('./admin.json'),
  auth: () => import('./auth.json'),
  cars: () => import('./cars.json'),
  cart: () => import('./cart.json'),
  chat: () => import('./chat.json'),
  checkout: () => import('./checkout.json'),
  delivery: () => import('./delivery.json'),
  map: () => import('./map.json'),
  marketplace: () => import('./marketplace.json'),
  misc: () => import('./misc.json'),
  notifications: () => import('./notifications.json'),
  orders: () => import('./orders.json'),
  products: () => import('./products.json'),
  profile: () => import('./profile.json'),
  realEstate: () => import('./realEstate.json'),
  reviews: () => import('./reviews.json'),
  search: () => import('./search.json'),
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

  const moduleData = await loader();
  return moduleData.default || moduleData;
}

// Экспорт базовых переводов
export default common;
