/**
 * Утилита для загрузки переводов с поддержкой lazy loading
 * Используется с next-intl для оптимизации загрузки переводов
 */

type Locale = 'ru' | 'en' | 'sr';

// Типы доступных модулей
export type TranslationModule =
  | 'common' // Базовые переводы (всегда загружаются)
  | 'auth' // Авторизация и профиль
  | 'marketplace' // Маркетплейс и объявления
  | 'admin' // Админ панель
  | 'storefront' // Витрины магазинов
  | 'cars' // Автомобильный раздел
  | 'chat' // Чат и сообщения
  | 'cart' // Корзина и заказы
  | 'realEstate' // Недвижимость
  | 'services' // Услуги
  | 'map' // Карта
  | 'misc'; // Разное (metadata, bentoGrid и др.)

// Кэш для загруженных модулей
const moduleCache = new Map<string, any>();

/**
 * Загружает переводы для указанной локали и модулей
 * @param locale - Локаль (ru, en, sr)
 * @param modules - Массив модулей для загрузки
 * @returns Объект с переводами
 */
export async function loadMessages(
  locale: Locale,
  modules: TranslationModule[] = ['common']
): Promise<Record<string, any>> {
  const messages: Record<string, any> = {};

  // Всегда загружаем common модуль
  if (!modules.includes('common')) {
    modules.unshift('common');
  }

  // Загружаем каждый модуль
  for (const mod of modules) {
    const cacheKey = `${locale}-${mod}`;

    // Проверяем кэш
    if (moduleCache.has(cacheKey)) {
      Object.assign(messages, moduleCache.get(cacheKey));
      continue;
    }

    try {
      // Динамический импорт модуля
      let moduleData: any;

      switch (mod) {
        case 'common':
          moduleData = await import(`@/messages/${locale}/common.json`);
          break;
        case 'auth':
          moduleData = await import(`@/messages/${locale}/auth.json`);
          break;
        case 'marketplace':
          moduleData = await import(`@/messages/${locale}/marketplace.json`);
          break;
        case 'admin':
          moduleData = await import(`@/messages/${locale}/admin.json`);
          break;
        case 'storefront':
          moduleData = await import(`@/messages/${locale}/storefront.json`);
          break;
        case 'cars':
          moduleData = await import(`@/messages/${locale}/cars.json`);
          break;
        case 'chat':
          moduleData = await import(`@/messages/${locale}/chat.json`);
          break;
        case 'cart':
          moduleData = await import(`@/messages/${locale}/cart.json`);
          break;
        case 'realEstate':
          moduleData = await import(`@/messages/${locale}/realEstate.json`);
          break;
        case 'services':
          moduleData = await import(`@/messages/${locale}/services.json`);
          break;
        case 'misc':
          moduleData = await import(`@/messages/${locale}/misc.json`);
          break;
        case 'map':
          moduleData = await import(`@/messages/${locale}/map.json`);
          break;
        default:
          console.warn(`Unknown translation module: ${mod}`);
          continue;
      }

      // Сохраняем в кэш
      const data = moduleData.default || moduleData;
      moduleCache.set(cacheKey, data);

      // Добавляем к общим переводам
      Object.assign(messages, data);
      
      // Для модулей map и cars также добавляем их содержимое под соответствующими ключами
      // чтобы поддержать обращения типа map.cluster.* и cars.filters.*
      if (mod === 'map' || mod === 'cars') {
        messages[mod] = data;
      }
    } catch (error) {
      console.error(
        `Failed to load translation module ${mod} for locale ${locale}:`,
        error
      );
    }
  }

  return messages;
}

/**
 * Очищает кэш модулей (полезно при смене языка)
 */
export function clearModuleCache() {
  moduleCache.clear();
}

/**
 * Предзагружает модули для улучшения производительности
 * @param locale - Локаль
 * @param modules - Модули для предзагрузки
 */
export async function preloadModules(
  locale: Locale,
  modules: TranslationModule[]
): Promise<void> {
  await loadMessages(locale, modules);
}

/**
 * Хелпер для определения необходимых модулей по пути страницы
 * @param pathname - Путь страницы
 * @returns Массив необходимых модулей
 */
export function getRequiredModules(pathname: string): TranslationModule[] {
  const modules: TranslationModule[] = ['common'];

  // Определяем модули по пути
  if (pathname.includes('/admin')) {
    modules.push('admin');
  }
  if (pathname.includes('/marketplace') || pathname.includes('/listing')) {
    modules.push('marketplace');
  }
  if (pathname.includes('/store') || pathname.includes('/storefront')) {
    modules.push('storefront');
  }
  if (pathname.includes('/cars') || pathname.includes('/automotive')) {
    modules.push('cars');
  }
  if (pathname.includes('/chat') || pathname.includes('/messages')) {
    modules.push('chat');
  }
  if (pathname.includes('/cart') || pathname.includes('/checkout')) {
    modules.push('cart');
  }
  if (
    pathname.includes('/auth') ||
    pathname.includes('/login') ||
    pathname.includes('/profile')
  ) {
    modules.push('auth');
  }
  if (pathname.includes('/real-estate') || pathname.includes('/property')) {
    modules.push('realEstate');
  }
  if (pathname.includes('/services')) {
    modules.push('services');
  }
  if (pathname.includes('/map')) {
    modules.push('map', 'marketplace'); // map страница также использует marketplace переводы
  }

  // Добавляем misc и cars модули для главной страницы и общих компонентов
  if (pathname === '/' || pathname.match(/^\/[a-z]{2}$/)) {
    modules.push('misc', 'cars'); // cars нужен для автомобильных фильтров на главной
  }

  return modules;
}
