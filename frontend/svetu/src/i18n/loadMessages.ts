import { Locale } from './config';

type Messages = Record<string, any>;

export async function loadMessages(
  locale: Locale,
  modules: string[] = []
): Promise<Messages> {
  try {
    // Загрузка базовых модулей
    const baseModules = ['common', 'nav', 'home'];

    // Объединение с дополнительными модулями
    const allModules = [...baseModules, ...modules];
    const uniqueModules = [...new Set(allModules)];

    // Загрузка всех модулей
    const messagePromises = uniqueModules.map(async (module) => {
      try {
        const messages = await import(`../messages/${locale}/${module}.json`);
        return { [module]: messages.default };
      } catch (error) {
        console.warn(
          `Failed to load ${module} messages for locale ${locale}:`,
          error
        );
        return { [module]: {} };
      }
    });

    const moduleResults = await Promise.all(messagePromises);

    // Объединение всех сообщений
    const allMessages = moduleResults.reduce((acc, moduleMessages) => {
      return { ...acc, ...moduleMessages };
    }, {});

    return allMessages;
  } catch (error) {
    console.error('Error loading messages:', error);
    return {};
  }
}

export function getRequiredModules(pathname: string): string[] {
  // Определение необходимых модулей на основе пути
  const modules: string[] = [];

  if (pathname.includes('/admin')) {
    modules.push('admin');
  }

  if (pathname.includes('/profile')) {
    modules.push('profile');
  }

  if (pathname.includes('/marketplace')) {
    modules.push('marketplace');
  }

  if (pathname.includes('/chat')) {
    modules.push('chat');
  }

  if (pathname.includes('/storefront')) {
    modules.push('storefront');
  }

  return modules;
}
