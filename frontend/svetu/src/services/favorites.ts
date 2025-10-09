import { apiClient } from './api-client';
import { toast } from 'react-hot-toast';

export interface FavoriteItem {
  id: number | string; // Может быть string для товаров витрин с префиксом sp_
  title: string;
  price: number;
  currency: string;
  location?: string;
  image?: string;
  created_at?: string;
  user_id?: number;
  category?: {
    id: number;
    name: string;
  };
  is_storefront_product?: boolean; // Флаг для товаров витрин
}

class FavoritesService {
  private favoriteIds: Set<number> = new Set();
  private isInitialized = false;
  private lastFetchTime: number = 0;
  private cachedFavorites: FavoriteItem[] | null = null;
  private readonly CACHE_DURATION = 60000; // 60 секунд кеш

  // Инициализация сервиса - загрузка списка избранных
  async initialize() {
    if (this.isInitialized) return;

    try {
      const response = await apiClient.get('/c2c/favorites');
      if (response.data?.success && Array.isArray(response.data.data)) {
        this.favoriteIds = new Set(
          response.data.data.map((item: any) => item.id)
        );
        this.isInitialized = true;
        // Сбрасываем кеш при инициализации
        this.cachedFavorites = null;
        this.lastFetchTime = 0;
      }
    } catch (error) {
      console.error('Failed to initialize favorites:', error);
      // Не показываем ошибку если пользователь не авторизован
    }
  }

  // Получить список всех избранных объявлений
  async getFavorites(): Promise<FavoriteItem[]> {
    // Проверяем кеш
    const now = Date.now();
    if (
      this.cachedFavorites &&
      this.lastFetchTime &&
      now - this.lastFetchTime < this.CACHE_DURATION
    ) {
      return this.cachedFavorites;
    }

    try {
      // Получаем избранные объявления (включая товары витрин)
      const response = await apiClient.get('/c2c/favorites');
      if (response.data?.success) {
        const rawFavorites = response.data.data || [];
        // Преобразуем backend данные в FavoriteItem
        const favorites = rawFavorites.map((item: any) => {
          // Получаем первое изображение или главное изображение
          let mainImage = null;
          if (
            item.images &&
            Array.isArray(item.images) &&
            item.images.length > 0
          ) {
            // Ищем главное изображение или берем первое
            const mainImg =
              item.images.find((img: any) => img.is_main) || item.images[0];
            mainImage =
              mainImg.image_url || mainImg.public_url || mainImg.file_path;
          }

          // Для товаров витрин добавляем префикс sp_ к ID
          const id = item.is_storefront_product ? `sp_${item.id}` : item.id;

          return {
            id: id,
            title: item.title,
            price: item.price,
            currency: item.currency || 'EUR',
            location: item.location || item.city,
            image: mainImage,
            created_at: item.created_at,
            user_id: item.user_id,
            category: item.category,
            is_storefront_product: item.is_storefront_product || false,
          } as FavoriteItem;
        });

        // Обновляем локальный кеш
        this.favoriteIds.clear();
        favorites.forEach((item: FavoriteItem) => {
          // Для товаров витрин с префиксом sp_ извлекаем числовой ID
          if (typeof item.id === 'string' && item.id.startsWith('sp_')) {
            const numericId = parseInt(item.id.replace('sp_', ''));
            if (!isNaN(numericId)) {
              this.favoriteIds.add(numericId);
            }
          } else {
            this.favoriteIds.add(item.id as number);
          }
        });

        // Сохраняем в кеш
        this.cachedFavorites = favorites;
        this.lastFetchTime = Date.now();

        return favorites;
      }
      return [];
    } catch (error: any) {
      if (error.response?.status === 401) {
        // Пользователь не авторизован
        return [];
      }
      console.error('Failed to get favorites:', error);
      throw error;
    }
  }

  // Добавить в избранное
  async addToFavorites(listingId: number | string): Promise<boolean> {
    try {
      // Определяем тип товара по ID
      const isB2CProduct =
        typeof listingId === 'string' && listingId.startsWith('sp_');

      // Для товаров витрин передаем специальный параметр type
      const endpoint = isB2CProduct
        ? `/c2c/favorites/${listingId.replace('sp_', '')}?type=storefront`
        : `/c2c/favorites/${listingId}`;

      const response = await apiClient.post(endpoint);
      if (response.data?.success) {
        // Для хранения используем числовой ID
        const numericId = isB2CProduct
          ? parseInt(String(listingId).replace('sp_', ''))
          : Number(listingId);
        this.favoriteIds.add(numericId);
        // Сбрасываем кеш при изменении
        this.cachedFavorites = null;
        toast.success('Добавлено в избранное');
        // Отправляем событие об изменении избранного
        window.dispatchEvent(new Event('favoritesChanged'));
        return true;
      }
      return false;
    } catch (error: any) {
      if (error.response?.status === 401) {
        toast.error('Войдите, чтобы добавить в избранное');
      } else {
        toast.error('Ошибка при добавлении в избранное');
      }
      console.error('Failed to add to favorites:', error);
      return false;
    }
  }

  // Удалить из избранного
  async removeFromFavorites(listingId: number | string): Promise<boolean> {
    try {
      // Определяем тип товара по ID
      const isB2CProduct =
        typeof listingId === 'string' && listingId.startsWith('sp_');
      const endpoint = isB2CProduct
        ? `/c2c/favorites/${listingId.replace('sp_', '')}?type=storefront`
        : `/c2c/favorites/${listingId}`;

      const response = await apiClient.delete(endpoint);
      if (response.data?.success) {
        // Для хранения используем числовой ID
        const numericId = isB2CProduct
          ? parseInt(String(listingId).replace('sp_', ''))
          : Number(listingId);
        this.favoriteIds.delete(numericId);
        // Сбрасываем кеш при изменении
        this.cachedFavorites = null;
        toast.success('Удалено из избранного');
        // Отправляем событие об изменении избранного
        window.dispatchEvent(new Event('favoritesChanged'));
        return true;
      }
      return false;
    } catch (error: any) {
      if (error.response?.status === 401) {
        toast.error('Войдите, чтобы управлять избранным');
      } else {
        toast.error('Ошибка при удалении из избранного');
      }
      console.error('Failed to remove from favorites:', error);
      return false;
    }
  }

  // Переключить статус избранного
  async toggleFavorite(listingId: number | string): Promise<boolean> {
    const numericId =
      typeof listingId === 'string' && listingId.startsWith('sp_')
        ? parseInt(listingId.replace('sp_', ''))
        : Number(listingId);

    if (this.isInFavorites(numericId)) {
      return this.removeFromFavorites(listingId);
    } else {
      return this.addToFavorites(listingId);
    }
  }

  // Проверить, находится ли объявление в избранном (локальная проверка)
  isInFavorites(listingId: number | string): boolean {
    const numericId =
      typeof listingId === 'string' && listingId.startsWith('sp_')
        ? parseInt(listingId.replace('sp_', ''))
        : Number(listingId);
    return this.favoriteIds.has(numericId);
  }

  // Проверить статус на сервере
  async checkFavoriteStatus(listingId: number): Promise<boolean> {
    try {
      const response = await apiClient.get(`/c2c/favorites/${listingId}/check`);
      if (response.data?.success) {
        const isInFavorites = response.data.data?.isInFavorites || false;
        if (isInFavorites) {
          this.favoriteIds.add(listingId);
        } else {
          this.favoriteIds.delete(listingId);
        }
        return isInFavorites;
      }
      return false;
    } catch (error) {
      console.error('Failed to check favorite status:', error);
      return false;
    }
  }

  // Получить количество избранных
  async getFavoritesCount(): Promise<number> {
    try {
      const response = await apiClient.get('/c2c/favorites/count');
      if (response.data?.success) {
        return response.data.data?.count || 0;
      }
      return 0;
    } catch (error) {
      console.error('Failed to get favorites count:', error);
      return 0;
    }
  }

  // Получить Set с ID избранных объявлений
  getFavoritesIds(): Set<number> {
    return this.favoriteIds;
  }

  // Очистить кеш (при выходе пользователя)
  clearCache() {
    this.favoriteIds.clear();
    this.isInitialized = false;
    // Отправляем событие об очистке избранного
    window.dispatchEvent(new Event('favoritesChanged'));
  }
}

// Экспортируем singleton
export const favoritesService = new FavoritesService();
