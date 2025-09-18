/**
 * Статические импорты всех переводов для корректной работы SSG
 */

// Английские переводы
import enCommon from '@/messages/en/common.json';
import enMarketplace from '@/messages/en/marketplace.json';
import enAuth from '@/messages/en/auth.json';
import enMisc from '@/messages/en/misc.json';
import enCart from '@/messages/en/cart.json';
import enMap from '@/messages/en/map.json';
import enStorefronts from '@/messages/en/storefronts.json';
import enAdmin from '@/messages/en/admin.json';
import enCars from '@/messages/en/cars.json';
import enReviews from '@/messages/en/reviews.json';
import enAuthShared from '@/messages/en/auth-shared.json';
import enBalance from '@/messages/en/balance.json';
import enProfile from '@/messages/en/profile.json';
import enChat from '@/messages/en/chat.json';
import enSearch from '@/messages/en/search.json';
import enCheckout from '@/messages/en/checkout.json';
import enOrders from '@/messages/en/orders.json';
import enSubscription from '@/messages/en/subscription.json';
import enFavorites from '@/messages/en/favorites.json';
import enTracking from '@/messages/en/tracking.json';

// Русские переводы
import ruCommon from '@/messages/ru/common.json';
import ruMarketplace from '@/messages/ru/marketplace.json';
import ruAuth from '@/messages/ru/auth.json';
import ruMisc from '@/messages/ru/misc.json';
import ruCart from '@/messages/ru/cart.json';
import ruMap from '@/messages/ru/map.json';
import ruStorefronts from '@/messages/ru/storefronts.json';
import ruAdmin from '@/messages/ru/admin.json';
import ruCars from '@/messages/ru/cars.json';
import ruReviews from '@/messages/ru/reviews.json';
import ruAuthShared from '@/messages/ru/auth-shared.json';
import ruBalance from '@/messages/ru/balance.json';
import ruProfile from '@/messages/ru/profile.json';
import ruChat from '@/messages/ru/chat.json';
import ruSearch from '@/messages/ru/search.json';
import ruCheckout from '@/messages/ru/checkout.json';
import ruOrders from '@/messages/ru/orders.json';
import ruSubscription from '@/messages/ru/subscription.json';
import ruFavorites from '@/messages/ru/favorites.json';
import ruTracking from '@/messages/ru/tracking.json';

// Сербские переводы
import srCommon from '@/messages/sr/common.json';
import srMarketplace from '@/messages/sr/marketplace.json';
import srAuth from '@/messages/sr/auth.json';
import srMisc from '@/messages/sr/misc.json';
import srCart from '@/messages/sr/cart.json';
import srMap from '@/messages/sr/map.json';
import srStorefronts from '@/messages/sr/storefronts.json';
import srAdmin from '@/messages/sr/admin.json';
import srCars from '@/messages/sr/cars.json';
import srReviews from '@/messages/sr/reviews.json';
import srAuthShared from '@/messages/sr/auth-shared.json';
import srBalance from '@/messages/sr/balance.json';
import srProfile from '@/messages/sr/profile.json';
import srChat from '@/messages/sr/chat.json';
import srSearch from '@/messages/sr/search.json';
import srCheckout from '@/messages/sr/checkout.json';
import srOrders from '@/messages/sr/orders.json';
import srSubscription from '@/messages/sr/subscription.json';
import srFavorites from '@/messages/sr/favorites.json';
import srTracking from '@/messages/sr/tracking.json';

// Объединённые сообщения по языкам
export const messages = {
  en: {
    common: enCommon,
    marketplace: enMarketplace,
    auth: enAuth,
    misc: enMisc,
    cart: enCart,
    map: enMap,
    storefronts: enStorefronts,
    admin: enAdmin,
    cars: enCars,
    reviews: enReviews,
    'auth-shared': enAuthShared,
    balance: enBalance,
    profile: enProfile,
    chat: enChat,
    search: enSearch,
    checkout: enCheckout,
    orders: enOrders,
    subscription: enSubscription,
    favorites: enFavorites,
    tracking: enTracking,
  },
  ru: {
    common: ruCommon,
    marketplace: ruMarketplace,
    auth: ruAuth,
    misc: ruMisc,
    cart: ruCart,
    map: ruMap,
    storefronts: ruStorefronts,
    admin: ruAdmin,
    cars: ruCars,
    reviews: ruReviews,
    'auth-shared': ruAuthShared,
    balance: ruBalance,
    profile: ruProfile,
    chat: ruChat,
    search: ruSearch,
    checkout: ruCheckout,
    orders: ruOrders,
    subscription: ruSubscription,
    favorites: ruFavorites,
    tracking: ruTracking,
  },
  sr: {
    common: srCommon,
    marketplace: srMarketplace,
    auth: srAuth,
    misc: srMisc,
    cart: srCart,
    map: srMap,
    storefronts: srStorefronts,
    admin: srAdmin,
    cars: srCars,
    reviews: srReviews,
    'auth-shared': srAuthShared,
    balance: srBalance,
    profile: srProfile,
    chat: srChat,
    search: srSearch,
    checkout: srCheckout,
    orders: srOrders,
    subscription: srSubscription,
    favorites: srFavorites,
    tracking: srTracking,
  },
};

/**
 * Функция для получения сообщений с правильной структурой для next-intl
 */
export function getMessages(locale: 'en' | 'ru' | 'sr') {
  const localeMessages = messages[locale];
  const result: Record<string, any> = {};

  // Просто возвращаем вложенную структуру модулей
  // Next-intl сам обрабатывает вложенность через точечную нотацию
  Object.entries(localeMessages).forEach(([module, data]) => {
    // Добавляем модуль как namespace
    result[module] = data;

    // Также добавляем ключи верхнего уровня для обратной совместимости
    // Это позволяет обращаться напрямую к ключам без указания модуля
    Object.entries(data).forEach(([key, value]) => {
      if (!result[key]) {
        result[key] = value;
      }
    });
  });

  return result;
}
