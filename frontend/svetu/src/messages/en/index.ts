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
  | 'calculator'
  | 'cars'
  | 'cart'
  | 'chat'
  | 'checkout'
  | 'common'
  | 'condition'
  | 'dashboard'
  | 'delivery'
  | 'filters'
  | 'gis'
  | 'home'
  | 'map'
  | 'marketplace'
  | 'marketplace.home'
  | 'misc'
  | 'notifications'
  | 'orders'
  | 'payment'
  | 'products'
  | 'profile'
  | 'realEstate'
  | 'recommendations'
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
  calculator: () => import('./calculator.json'),
  cars: () => import('./cars.json'),
  cart: () => import('./cart.json'),
  chat: () => import('./chat.json'),
  checkout: () => import('./checkout.json'),
  condition: () => import('./condition.json'),
  dashboard: () => import('./dashboard.json'),
  delivery: () => import('./delivery.json'),
  filters: () => import('./filters.json'),
  gis: () => import('./gis.json'),
  home: () => import('./home.json'),
  map: () => import('./map.json'),
  marketplace: () => import('./c2c.json'),
  'marketplace.home': () => import('./c2c.home.json'),
  misc: () => import('./misc.json'),
  notifications: () => import('./notifications.json'),
  orders: () => import('./orders.json'),
  payment: () => import('./payment.json'),
  products: () => import('./products.json'),
  profile: () => import('./profile.json'),
  realEstate: () => import('./realEstate.json'),
  recommendations: () => import('./recommendations.json'),
  reviews: () => import('./reviews.json'),
  scanner: () => import('./scanner.json'),
  search: () => import('./search.json'),
  services: () => import('./services.json'),
  storefronts: () => import('./b2c.json'),
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
