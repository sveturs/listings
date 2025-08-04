import { describe, it, expect, beforeEach, vi } from 'vitest';
import { loadMessages, clearModuleCache, getRequiredModules, preloadModules } from '../loadMessages';

// Мокаем динамические импорты
vi.mock('@/messages/ru/common.json', () => ({
  default: {
    loading: 'Загрузка...',
    save: 'Сохранить',
    cancel: 'Отмена'
  }
}));

vi.mock('@/messages/ru/marketplace.json', () => ({
  default: {
    title: 'Маркетплейс',
    createListing: 'Создать объявление'
  }
}));

vi.mock('@/messages/ru/admin.json', () => ({
  default: {
    title: 'Админ панель',
    users: 'Пользователи'
  }
}));

describe('loadMessages', () => {
  beforeEach(() => {
    clearModuleCache();
  });

  it('должен загружать базовый модуль common по умолчанию', async () => {
    const messages = await loadMessages('ru', []);
    
    expect(messages).toHaveProperty('loading', 'Загрузка...');
    expect(messages).toHaveProperty('save', 'Сохранить');
  });

  it('должен загружать несколько модулей', async () => {
    const messages = await loadMessages('ru', ['common', 'marketplace']);
    
    expect(messages).toHaveProperty('loading', 'Загрузка...');
    expect(messages).toHaveProperty('title', 'Маркетплейс');
    expect(messages).toHaveProperty('createListing', 'Создать объявление');
  });

  it('должен кэшировать загруженные модули', async () => {
    // Первая загрузка
    const messages1 = await loadMessages('ru', ['admin']);
    expect(messages1).toHaveProperty('title', 'Админ панель');
    
    // Вторая загрузка должна использовать кэш
    const messages2 = await loadMessages('ru', ['admin']);
    expect(messages2).toEqual(messages1);
  });

  it('должен обрабатывать ошибки загрузки модулей', async () => {
    const consoleError = vi.spyOn(console, 'error').mockImplementation(() => {});
    
    // Пытаемся загрузить несуществующий модуль
    const messages = await loadMessages('ru', ['nonexistent' as any]);
    
    // Должны получить только common модуль
    expect(messages).toHaveProperty('loading');
    expect(consoleError).toHaveBeenCalled();
    
    consoleError.mockRestore();
  });
});

describe('getRequiredModules', () => {
  it('должен определять модули для админки', () => {
    const modules = getRequiredModules('/ru/admin/users');
    
    expect(modules).toContain('common');
    expect(modules).toContain('admin');
  });

  it('должен определять модули для маркетплейса', () => {
    const modules = getRequiredModules('/ru/marketplace/listings');
    
    expect(modules).toContain('common');
    expect(modules).toContain('marketplace');
  });

  it('должен определять модули для корзины', () => {
    const modules = getRequiredModules('/ru/cart');
    
    expect(modules).toContain('common');
    expect(modules).toContain('cart');
  });

  it('должен определять несколько модулей для сложных путей', () => {
    const modules = getRequiredModules('/ru/profile/listings');
    
    expect(modules).toContain('common');
    expect(modules).toContain('auth');
    expect(modules).toContain('marketplace');
  });
});

describe('preloadModules', () => {
  it('должен предзагружать модули без ошибок', async () => {
    await expect(preloadModules('ru', ['marketplace', 'admin'])).resolves.not.toThrow();
  });
});

describe('clearModuleCache', () => {
  it('должен очищать кэш модулей', async () => {
    // Загружаем модуль
    await loadMessages('ru', ['admin']);
    
    // Очищаем кэш
    clearModuleCache();
    
    // При следующей загрузке должен загрузиться заново
    const messages = await loadMessages('ru', ['admin']);
    expect(messages).toHaveProperty('title', 'Админ панель');
  });
});