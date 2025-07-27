// Утилиты для работы с адресами и их приватностью

export type LocationPrivacyLevel = 'exact' | 'street' | 'district' | 'city';

/**
 * Форматирует адрес с учетом уровня приватности
 */
export function formatAddressWithPrivacy(
  address: string | undefined,
  privacyLevel: LocationPrivacyLevel | undefined
): string {
  if (!address) return '';
  if (!privacyLevel || privacyLevel === 'exact') return address;

  // Разбиваем адрес на части
  const parts = address.split(',').map((part) => part.trim());

  switch (privacyLevel) {
    case 'street':
      // Показываем улицу и город, скрываем номер дома
      if (parts.length > 1) {
        const street = parts[0].replace(/\d+/g, '').trim();
        return [street, ...parts.slice(1)].join(', ');
      }
      return address;

    case 'district':
      // Показываем только район и город
      if (parts.length > 2) {
        return parts.slice(-2).join(', ');
      }
      return address;

    case 'city':
      // Показываем только город
      if (parts.length > 0) {
        return parts[parts.length - 1];
      }
      return address;

    default:
      return address;
  }
}

/**
 * Получает локализованный адрес из объекта с переводами
 */
export function getLocalizedAddress(
  defaultAddress: string | undefined,
  translations:
    | {
        sr?: string;
        en?: string;
        ru?: string;
      }
    | undefined,
  locale: string
): string {
  if (!translations) return defaultAddress || '';

  // Пытаемся получить адрес на нужном языке
  const localizedAddress = translations[locale as keyof typeof translations];

  // Если локализованный адрес есть, возвращаем его
  if (localizedAddress) return localizedAddress;

  // Иначе возвращаем адрес по умолчанию
  return defaultAddress || '';
}

/**
 * Получает полный локализованный адрес из компонентов
 */
export function getFullLocalizedAddress(
  item: {
    location?: string;
    city?: string;
    country?: string;
    translations?: any; // Принимаем любой формат translations
  },
  locale: string
): string {
  let location = item.location || '';
  let city = item.city || '';
  let country = item.country || '';

  if (item.translations) {
    // Проверяем формат 1: { [locale]: { location, city, country } }
    if (item.translations[locale] && typeof item.translations[locale] === 'object' && 'location' in item.translations[locale]) {
      location = item.translations[locale].location || location;
      city = item.translations[locale].city || city;
      country = item.translations[locale].country || country;
    }
    // Проверяем формат 2: { location: { [locale]: string }, ... }
    else {
      if (item.translations.location && item.translations.location[locale]) {
        location = item.translations.location[locale] || location;
      }
      if (item.translations.city && item.translations.city[locale]) {
        city = item.translations.city[locale] || city;
      }
      if (item.translations.country && item.translations.country[locale]) {
        country = item.translations.country[locale] || country;
      }
    }
  }

  // Собираем полный адрес из непустых компонентов
  const parts = [location, city, country].filter(Boolean);
  return parts.join(', ');
}

/**
 * Получает иконку для уровня приватности
 */
export function getPrivacyIcon(
  privacyLevel: LocationPrivacyLevel | undefined
): string {
  switch (privacyLevel) {
    case 'exact':
      return '📍'; // Точный адрес
    case 'street':
      return '🏘️'; // Улица
    case 'district':
      return '🏙️'; // Район
    case 'city':
      return '🌆'; // Город
    default:
      return '📍';
  }
}

/**
 * Получает текстовое описание уровня приватности
 */
export function getPrivacyDescription(
  privacyLevel: LocationPrivacyLevel | undefined,
  t: (key: string) => string
): string {
  switch (privacyLevel) {
    case 'exact':
      return t('products.privacy.exact');
    case 'street':
      return t('products.privacy.street');
    case 'district':
      return t('products.privacy.district');
    case 'city':
      return t('products.privacy.city');
    default:
      return t('products.privacy.exact');
  }
}
