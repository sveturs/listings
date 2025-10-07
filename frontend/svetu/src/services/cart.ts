import { apiClient } from './api-client';
import type { components } from '@/types/generated/api';

type ShoppingCart =
  components['schemas']['models.ShoppingCart'];
type AddToCartRequest =
  components['schemas']['models.AddToCartRequest'];
type UpdateCartItemRequest =
  components['schemas']['models.UpdateCartItemRequest'];

export const cartService = {
  // Получить все корзины пользователя
  async getUserCarts(): Promise<ShoppingCart[]> {
    const response = await apiClient.get('/api/v1/user/carts');
    console.log('[CartService] getUserCarts response:', response.data);
    const carts = response.data?.data || [];
    console.log('[CartService] getUserCarts carts:', carts);
    return carts;
  },

  // Получить корзину витрины
  async getCart(storefrontId: number): Promise<ShoppingCart> {
    const response = await apiClient.get(
      `/api/v1/storefronts/${storefrontId}/cart`
    );
    return response.data.data;
  },

  // Добавить товар в корзину
  async addToCart(
    storefrontId: number,
    item: AddToCartRequest
  ): Promise<ShoppingCart> {
    const response = await apiClient.post(
      `/api/v1/storefronts/${storefrontId}/cart/items`,
      item
    );
    return response.data.data;
  },

  // Обновить количество товара в корзине
  async updateCartItem(
    storefrontId: number,
    itemId: number,
    data: UpdateCartItemRequest
  ): Promise<ShoppingCart> {
    const response = await apiClient.put(
      `/api/v1/storefronts/${storefrontId}/cart/items/${itemId}`,
      data
    );
    return response.data.data;
  },

  // Удалить товар из корзины
  async removeFromCart(
    storefrontId: number,
    itemId: number
  ): Promise<ShoppingCart> {
    const response = await apiClient.delete(
      `/api/v1/storefronts/${storefrontId}/cart/items/${itemId}`
    );
    return response.data.data;
  },

  // Очистить корзину
  async clearCart(storefrontId: number): Promise<void> {
    await apiClient.delete(`/api/v1/storefronts/${storefrontId}/cart`);
  },
};
