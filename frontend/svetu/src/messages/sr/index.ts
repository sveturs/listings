// Автосгенерированный файл для модульной загрузки переводов
// Сгенерирован: 2025-08-04T10:32:06.116Z
// НЕ РЕДАКТИРУЙТЕ ВРУЧНУЮ!

// Базовые переводы (загружаются всегда)
import common from './common.json';

// Типы модулей
export type TranslationModule =
  | 'admin'
  | 'ar'
  | 'auth'
  | 'cars'
  | 'cart'
  | 'chat'
  | 'checkout'
  | 'common'
  | 'condition'
  | 'dashboard'
  | 'delivery'
  | 'gis'
  | 'home'
  | 'map'
  | 'marketplace'
  | 'misc'
  | 'notifications'
  | 'orders'
  | 'payment'
  | 'products'
  | 'profile'
  | 'realEstate'
  | 'reviews'
  | 'scanner'
  | 'search'
  | 'services'
  | 'storefronts'
  | 'trust';

// Карта модулей для динамической загрузки
export const moduleLoaders = {
  admin: () => import('./admin.json'),
  ar: () => import('./ar.json'),
  auth: () => import('./auth.json'),
  cars: () => import('./cars.json'),
  cart: () => import('./cart.json'),
  chat: () => import('./chat.json'),
  checkout: () => import('./checkout.json'),
  condition: () => import('./condition.json'),
  dashboard: () => import('./dashboard.json'),
  delivery: () => import('./delivery.json'),
  gis: () => import('./gis.json'),
  home: () => import('./home.json'),
  map: () => import('./map.json'),
  marketplace: () => import('./marketplace.json'),
  misc: () => import('./misc.json'),
  notifications: () => import('./notifications.json'),
  orders: () => import('./orders.json'),
  payment: () => import('./payment.json'),
  products: () => import('./products.json'),
  profile: () => import('./profile.json'),
  realEstate: () => import('./realEstate.json'),
  reviews: () => import('./reviews.json'),
  scanner: () => import('./scanner.json'),
  search: () => import('./search.json'),
  services: () => import('./services.json'),
  storefronts: () => import('./storefronts.json'),
  trust: () => import('./trust.json'),
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
