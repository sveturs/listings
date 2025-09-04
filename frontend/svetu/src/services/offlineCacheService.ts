// Offline Cache Service для мобильной оптимизации
// День 25: Offline Caching Стратегии

import { openDB, DBSchema, IDBPDatabase } from 'idb';

interface CacheDBSchema extends DBSchema {
  attributes: {
    key: string;
    value: {
      id: string;
      data: any;
      timestamp: number;
      version: string;
      categoryId?: number;
    };
    indexes: { 'by-category': number; 'by-timestamp': number };
  };
  searchResults: {
    key: string;
    value: {
      query: string;
      results: any[];
      timestamp: number;
      filters?: Record<string, any>;
    };
    indexes: { 'by-timestamp': number };
  };
  userPreferences: {
    key: string;
    value: {
      key: string;
      value: any;
      timestamp: number;
    };
  };
  pendingSync: {
    key: string;
    value: {
      id: string;
      action: 'create' | 'update' | 'delete';
      endpoint: string;
      data: any;
      timestamp: number;
      retries: number;
    };
    indexes: { 'by-timestamp': number };
  };
}

class OfflineCacheService {
  private db: IDBPDatabase<CacheDBSchema> | null = null;
  private readonly DB_NAME = 'svetu-offline-cache';
  private readonly DB_VERSION = 1;
  private readonly CACHE_DURATION = 7 * 24 * 60 * 60 * 1000; // 7 дней
  private readonly MAX_CACHE_SIZE = 50 * 1024 * 1024; // 50MB
  private syncInProgress = false;

  async init(): Promise<void> {
    try {
      this.db = await openDB<CacheDBSchema>(this.DB_NAME, this.DB_VERSION, {
        upgrade(db) {
          // Attributes store
          if (!db.objectStoreNames.contains('attributes')) {
            const attrStore = db.createObjectStore('attributes', {
              keyPath: 'id',
            });
            attrStore.createIndex('by-category', 'categoryId');
            attrStore.createIndex('by-timestamp', 'timestamp');
          }

          // Search results store
          if (!db.objectStoreNames.contains('searchResults')) {
            const searchStore = db.createObjectStore('searchResults', {
              keyPath: 'query',
            });
            searchStore.createIndex('by-timestamp', 'timestamp');
          }

          // User preferences store
          if (!db.objectStoreNames.contains('userPreferences')) {
            db.createObjectStore('userPreferences', { keyPath: 'key' });
          }

          // Pending sync store
          if (!db.objectStoreNames.contains('pendingSync')) {
            const syncStore = db.createObjectStore('pendingSync', {
              keyPath: 'id',
              autoIncrement: true,
            });
            syncStore.createIndex('by-timestamp', 'timestamp');
          }
        },
      });

      // Очистка старых данных при инициализации
      await this.cleanExpiredCache();

      // Регистрация background sync если поддерживается
      if ('serviceWorker' in navigator && 'SyncManager' in window) {
        await this.registerBackgroundSync();
      }

      // Слушаем изменения сети
      this.setupNetworkListener();
    } catch (error) {
      console.error('Failed to initialize offline cache:', error);
    }
  }

  // Кеширование атрибутов
  async cacheAttributes(categoryId: number, attributes: any[]): Promise<void> {
    if (!this.db) return;

    const tx = this.db.transaction('attributes', 'readwrite');
    const store = tx.objectStore('attributes');

    for (const attr of attributes) {
      await store.put({
        id: `${categoryId}-${attr.id}`,
        data: attr,
        timestamp: Date.now(),
        version: this.getCurrentVersion(),
        categoryId,
      });
    }

    await tx.complete;
  }

  // Получение атрибутов из кеша
  async getCachedAttributes(categoryId: number): Promise<any[] | null> {
    if (!this.db) return null;

    try {
      const index = this.db
        .transaction('attributes')
        .store.index('by-category');
      const items = await index.getAll(categoryId);

      if (items.length === 0) return null;

      // Проверяем актуальность кеша
      const isExpired = items.some(
        (item) => Date.now() - item.timestamp > this.CACHE_DURATION
      );

      if (isExpired) {
        // Помечаем для обновления в фоне
        this.scheduleBackgroundUpdate('attributes', categoryId);
        return null;
      }

      return items.map((item) => item.data);
    } catch (error) {
      console.error('Failed to get cached attributes:', error);
      return null;
    }
  }

  // Кеширование результатов поиска
  async cacheSearchResults(
    query: string,
    results: any[],
    filters?: Record<string, any>
  ): Promise<void> {
    if (!this.db) return;

    const cacheKey = this.generateSearchCacheKey(query, filters);

    await this.db.put('searchResults', {
      query: cacheKey,
      results,
      timestamp: Date.now(),
      filters,
    });
  }

  // Получение результатов поиска из кеша
  async getCachedSearchResults(
    query: string,
    filters?: Record<string, any>
  ): Promise<any[] | null> {
    if (!this.db) return null;

    const cacheKey = this.generateSearchCacheKey(query, filters);

    try {
      const cached = await this.db.get('searchResults', cacheKey);

      if (!cached) return null;

      // Проверка актуальности (поиск кешируется на меньшее время)
      const searchCacheDuration = 60 * 60 * 1000; // 1 час
      if (Date.now() - cached.timestamp > searchCacheDuration) {
        return null;
      }

      return cached.results;
    } catch (error) {
      console.error('Failed to get cached search results:', error);
      return null;
    }
  }

  // Сохранение пользовательских настроек
  async saveUserPreference(key: string, value: any): Promise<void> {
    if (!this.db) return;

    await this.db.put('userPreferences', {
      key,
      value,
      timestamp: Date.now(),
    });
  }

  // Получение пользовательских настроек
  async getUserPreference(key: string): Promise<any> {
    if (!this.db) return null;

    try {
      const pref = await this.db.get('userPreferences', key);
      return pref?.value || null;
    } catch (error) {
      console.error('Failed to get user preference:', error);
      return null;
    }
  }

  // Добавление в очередь синхронизации
  async addToPendingSync(
    action: 'create' | 'update' | 'delete',
    endpoint: string,
    data: any
  ): Promise<void> {
    if (!this.db) return;

    await this.db.add('pendingSync', {
      id: '', // auto-increment
      action,
      endpoint,
      data,
      timestamp: Date.now(),
      retries: 0,
    });

    // Попытка немедленной синхронизации если есть сеть
    if (navigator.onLine) {
      await this.syncPendingChanges();
    }
  }

  // Синхронизация ожидающих изменений
  async syncPendingChanges(): Promise<void> {
    if (!this.db || this.syncInProgress) return;

    this.syncInProgress = true;

    try {
      const tx = this.db.transaction('pendingSync', 'readwrite');
      const store = tx.objectStore('pendingSync');
      const items = await store.getAll();

      for (const item of items) {
        try {
          // Отправка запроса на сервер
          const response = await fetch(item.endpoint, {
            method:
              item.action === 'delete'
                ? 'DELETE'
                : item.action === 'create'
                  ? 'POST'
                  : 'PUT',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(item.data),
          });

          if (response.ok) {
            // Удаляем из очереди после успешной синхронизации
            await store.delete(item.id);
          } else {
            // Увеличиваем счетчик попыток
            item.retries++;
            if (item.retries < 3) {
              await store.put(item);
            } else {
              // После 3 попыток удаляем из очереди
              await store.delete(item.id);
              console.error('Failed to sync after 3 retries:', item);
            }
          }
        } catch (error) {
          console.error('Sync error for item:', item, error);
          // Оставляем в очереди для повторной попытки
        }
      }

      await tx.complete;
    } finally {
      this.syncInProgress = false;
    }
  }

  // Очистка устаревшего кеша
  async cleanExpiredCache(): Promise<void> {
    if (!this.db) return;

    const now = Date.now();

    // Очистка атрибутов
    const attrTx = this.db.transaction('attributes', 'readwrite');
    const attrIndex = attrTx.store.index('by-timestamp');
    const expiredAttrs = await attrIndex.getAllKeys(
      IDBKeyRange.upperBound(now - this.CACHE_DURATION)
    );
    for (const key of expiredAttrs) {
      await attrTx.store.delete(key);
    }
    await attrTx.complete;

    // Очистка результатов поиска
    const searchTx = this.db.transaction('searchResults', 'readwrite');
    const searchIndex = searchTx.store.index('by-timestamp');
    const expiredSearches = await searchIndex.getAllKeys(
      IDBKeyRange.upperBound(now - 60 * 60 * 1000) // 1 час для поиска
    );
    for (const key of expiredSearches) {
      await searchTx.store.delete(key);
    }
    await searchTx.complete;
  }

  // Проверка размера кеша
  async getCacheSize(): Promise<number> {
    if (!('estimate' in navigator.storage)) {
      return 0;
    }

    const estimate = await navigator.storage.estimate();
    return estimate.usage || 0;
  }

  // Очистка всего кеша
  async clearCache(): Promise<void> {
    if (!this.db) return;

    const stores: Array<keyof CacheDBSchema> = [
      'attributes',
      'searchResults',
      'userPreferences',
    ];

    for (const storeName of stores) {
      await this.db.clear(storeName);
    }
  }

  // Экспорт данных для backup
  async exportData(): Promise<any> {
    if (!this.db) return null;

    const data: any = {};

    data.attributes = await this.db.getAll('attributes');
    data.searchResults = await this.db.getAll('searchResults');
    data.userPreferences = await this.db.getAll('userPreferences');
    data.pendingSync = await this.db.getAll('pendingSync');

    return data;
  }

  // Импорт данных из backup
  async importData(data: any): Promise<void> {
    if (!this.db) return;

    if (data.attributes) {
      const tx = this.db.transaction('attributes', 'readwrite');
      for (const item of data.attributes) {
        await tx.store.put(item);
      }
      await tx.complete;
    }

    if (data.userPreferences) {
      const tx = this.db.transaction('userPreferences', 'readwrite');
      for (const item of data.userPreferences) {
        await tx.store.put(item);
      }
      await tx.complete;
    }
  }

  // Вспомогательные методы
  private generateSearchCacheKey(
    query: string,
    filters?: Record<string, any>
  ): string {
    const filterStr = filters ? JSON.stringify(filters) : '';
    return `${query.toLowerCase().trim()}-${filterStr}`;
  }

  private getCurrentVersion(): string {
    return '1.0.0'; // Можно получать из конфига или API
  }

  private async scheduleBackgroundUpdate(type: string, id: any): Promise<void> {
    // Планирование обновления в фоне
    if ('requestIdleCallback' in window) {
      requestIdleCallback(() => {
        // Обновление данных когда браузер не занят
        console.log(`Scheduling background update for ${type}:${id}`);
      });
    }
  }

  private setupNetworkListener(): void {
    window.addEventListener('online', () => {
      console.log('Network restored, syncing pending changes...');
      this.syncPendingChanges();
    });

    window.addEventListener('offline', () => {
      console.log('Network lost, switching to offline mode');
    });
  }

  private async registerBackgroundSync(): Promise<void> {
    try {
      const registration = await navigator.serviceWorker.ready;
      if ('sync' in registration) {
        await (registration as any).sync.register('sync-pending-changes');
        console.log('Background sync registered');
      }
    } catch (error) {
      console.error('Failed to register background sync:', error);
    }
  }
}

// Singleton экземпляр
const offlineCacheService = new OfflineCacheService();

export default offlineCacheService;

// React hook для работы с offline cache
export const useOfflineCache = () => {
  const [isOffline, setIsOffline] = useState(!navigator.onLine);
  const [cacheSize, setCacheSize] = useState(0);

  useEffect(() => {
    const handleOnline = () => setIsOffline(false);
    const handleOffline = () => setIsOffline(true);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    // Проверка размера кеша
    offlineCacheService.getCacheSize().then(setCacheSize);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
    };
  }, []);

  return {
    isOffline,
    cacheSize,
    cacheService: offlineCacheService,
  };
};
