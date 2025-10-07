import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { cartService } from '@/services/cart';
import type { components } from '@/types/generated/api';

type ShoppingCart =
  components['schemas']['models.ShoppingCart'];
type AddToCartRequest =
  components['schemas']['models.AddToCartRequest'];
type UpdateCartItemRequest =
  components['schemas']['models.UpdateCartItemRequest'];

interface CartState {
  cart: ShoppingCart | null; // Текущая активная корзина
  allCarts: ShoppingCart[]; // Все корзины пользователя
  loading: boolean;
  error: string | null;
}

// Утилиты для работы с localStorage
const CART_STORAGE_KEY_PREFIX = 'svetu_cart';

const getCartStorageKey = (userId?: number): string => {
  // Если есть userId, используем уникальный ключ для пользователя
  // Иначе используем общий ключ для анонимной корзины
  return userId
    ? `${CART_STORAGE_KEY_PREFIX}_user_${userId}`
    : `${CART_STORAGE_KEY_PREFIX}_anon`;
};

const _loadCartFromStorage = (userId?: number): ShoppingCart | null => {
  if (typeof window === 'undefined') return null;
  try {
    const key = getCartStorageKey(userId);
    const savedCart = localStorage.getItem(key);
    if (savedCart) {
      return JSON.parse(savedCart);
    }
  } catch (error) {
    console.error('Failed to load cart from localStorage:', error);
  }
  return null;
};

const saveCartToStorage = (cart: ShoppingCart | null, userId?: number) => {
  if (typeof window === 'undefined') return;
  try {
    const key = getCartStorageKey(userId);
    if (cart) {
      localStorage.setItem(key, JSON.stringify(cart));
    } else {
      localStorage.removeItem(key);
    }
  } catch (error) {
    console.error('Failed to save cart to localStorage:', error);
  }
};

// Очистка старой корзины при смене пользователя
const clearOldCartStorage = () => {
  if (typeof window === 'undefined') return;
  try {
    // Очищаем общий ключ (для обратной совместимости)
    localStorage.removeItem('svetu_cart');

    // Очищаем анонимную корзину при входе пользователя
    localStorage.removeItem(`${CART_STORAGE_KEY_PREFIX}_anon`);
  } catch (error) {
    console.error('Failed to clear old cart from localStorage:', error);
  }
};

const initialState: CartState = {
  cart: null, // Не загружаем из localStorage при инициализации - будет загружено в useCartSync
  allCarts: [], // Все корзины пользователя
  loading: false,
  error: null,
};

// Async thunks
export const fetchUserCarts = createAsyncThunk(
  'cart/fetchUserCarts',
  async (userId: number) => {
    // Сначала очищаем старые корзины при входе нового пользователя
    clearOldCartStorage();

    const carts = await cartService.getUserCarts();
    // Если есть корзины, берем последнюю обновленную для текущего storefront
    // или первую из списка как активную
    return { carts, userId };
  }
);

export const fetchCart = createAsyncThunk(
  'cart/fetchCart',
  async (storefrontId: number) => {
    return await cartService.getCart(storefrontId);
  }
);

export const addToCart = createAsyncThunk(
  'cart/addToCart',
  async ({
    storefrontId,
    item,
  }: {
    storefrontId: number;
    item: AddToCartRequest;
  }) => {
    return await cartService.addToCart(storefrontId, item);
  }
);

export const updateCartItem = createAsyncThunk(
  'cart/updateCartItem',
  async ({
    storefrontId,
    itemId,
    data,
  }: {
    storefrontId: number;
    itemId: number;
    data: UpdateCartItemRequest;
  }) => {
    return await cartService.updateCartItem(storefrontId, itemId, data);
  }
);

export const removeFromCart = createAsyncThunk(
  'cart/removeFromCart',
  async ({
    storefrontId,
    itemId,
  }: {
    storefrontId: number;
    itemId: number;
  }) => {
    return await cartService.removeFromCart(storefrontId, itemId);
  }
);

export const clearCart = createAsyncThunk(
  'cart/clearCart',
  async (storefrontId: number) => {
    await cartService.clearCart(storefrontId);
    return null;
  }
);

const cartSlice = createSlice({
  name: 'cart',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    resetCart: (state) => {
      state.cart = null;
      state.allCarts = [];
      state.error = null;
      saveCartToStorage(null);
    },
    clearCartOnLogout: (state) => {
      state.cart = null;
      state.allCarts = [];
      state.error = null;
      // Очищаем все возможные корзины из localStorage при выходе
      clearOldCartStorage();
      // Также очищаем корзины для всех пользователей (на случай если есть старые ключи)
      if (typeof window !== 'undefined') {
        const keys = Object.keys(localStorage);
        keys.forEach((key) => {
          if (key.startsWith(CART_STORAGE_KEY_PREFIX)) {
            localStorage.removeItem(key);
          }
        });
      }
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch cart
      .addCase(fetchCart.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchCart.fulfilled, (state, action) => {
        state.loading = false;
        state.cart = action.payload;
        saveCartToStorage(action.payload);
      })
      .addCase(fetchCart.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch cart';
      })

      // Add to cart
      .addCase(addToCart.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(addToCart.fulfilled, (state, action) => {
        state.loading = false;
        const updatedCart = action.payload;
        state.cart = updatedCart;

        // Обновляем корзину в allCarts
        if (updatedCart) {
          const existingCartIndex = state.allCarts.findIndex(
            (cart) => cart.storefront_id === updatedCart.storefront_id
          );
          if (existingCartIndex >= 0) {
            // Обновляем существующую корзину
            state.allCarts[existingCartIndex] = updatedCart;
          } else {
            // Добавляем новую корзину если её не было
            state.allCarts.push(updatedCart);
          }
        }

        saveCartToStorage(updatedCart);
      })
      .addCase(addToCart.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to add item to cart';
      })

      // Update cart item
      .addCase(updateCartItem.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(updateCartItem.fulfilled, (state, action) => {
        state.loading = false;
        state.cart = action.payload;
        saveCartToStorage(action.payload);

        // Обновляем соответствующую корзину в allCarts
        if (action.payload && action.payload.storefront_id) {
          const cartIndex = state.allCarts.findIndex(
            (cart) => cart.storefront_id === action.payload.storefront_id
          );
          if (cartIndex !== -1) {
            state.allCarts[cartIndex] = action.payload;
          } else {
            // Добавляем корзину если её не было
            state.allCarts.push(action.payload);
          }
        }
      })
      .addCase(updateCartItem.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to update cart item';
      })

      // Remove from cart
      .addCase(removeFromCart.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(removeFromCart.fulfilled, (state, action) => {
        state.loading = false;
        state.cart = action.payload;
        saveCartToStorage(action.payload);

        // Обновляем соответствующую корзину в allCarts
        if (action.payload && action.payload.storefront_id) {
          const cartIndex = state.allCarts.findIndex(
            (cart) => cart.storefront_id === action.payload.storefront_id
          );
          if (cartIndex !== -1) {
            state.allCarts[cartIndex] = action.payload;
          } else {
            // Добавляем корзину если её не было
            state.allCarts.push(action.payload);
          }
        }
      })
      .addCase(removeFromCart.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to remove item from cart';
      })

      // Clear cart
      .addCase(clearCart.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(clearCart.fulfilled, (state, _action) => {
        state.loading = false;
        state.cart = null;
        saveCartToStorage(null);

        // Удаляем соответствующую корзину из allCarts или обновляем её как пустую
        // Нужно знать storefrontId чтобы обновить правильную корзину
        // Поскольку clearCart не возвращает storefrontId, будем очищать текущую активную корзину
        if (state.allCarts.length > 0) {
          // Если была активная корзина, найдем её в allCarts и очистим
          const currentCartIndex = state.allCarts.findIndex(
            (cart) => cart.items && cart.items.length > 0
          );
          if (currentCartIndex !== -1) {
            state.allCarts[currentCartIndex] = {
              ...state.allCarts[currentCartIndex],
              items: [],
            };
          }
        }
      })
      .addCase(clearCart.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to clear cart';
      })

      // Fetch user carts
      .addCase(fetchUserCarts.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchUserCarts.fulfilled, (state, action) => {
        state.loading = false;
        const { carts, userId } = action.payload;

        // Сохраняем все корзины
        state.allCarts = carts || [];

        // Если есть корзины, берем первую как активную
        if (carts && carts.length > 0) {
          // Находим корзину с наибольшим количеством товаров или последнюю обновленную
          const mostRecentCart = carts.reduce((prev, current) => {
            const prevTime = new Date(prev.updated_at || 0).getTime();
            const currentTime = new Date(current.updated_at || 0).getTime();
            return currentTime > prevTime ? current : prev;
          });
          state.cart = mostRecentCart;
          saveCartToStorage(mostRecentCart, userId);
        } else {
          // Если корзин нет, очищаем localStorage для этого пользователя
          state.cart = null;
          saveCartToStorage(null, userId);
        }
      })
      .addCase(fetchUserCarts.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch user carts';
      });
  },
});

export const { clearError, resetCart, clearCartOnLogout } = cartSlice.actions;

// Selectors
export const selectCart = (state: { cart: CartState }) => state.cart.cart;
export const selectCartLoading = (state: { cart: CartState }) =>
  state.cart.loading;
export const selectCartError = (state: { cart: CartState }) => state.cart.error;
export const selectCartItemsCount = (state: { cart: CartState }) =>
  state.cart.cart?.items?.reduce(
    (total, item) => total + (item.quantity || 0),
    0
  ) || 0;

// Селектор для всех корзин
export const selectAllCarts = (state: { cart: CartState }) =>
  state.cart.allCarts;

// Селектор для подсчета всех товаров во всех корзинах
export const selectAllCartsItemsCount = (state: { cart: CartState }) =>
  state.cart.allCarts.reduce(
    (total, cart) =>
      total +
      (cart.items?.reduce(
        (cartTotal, item) => cartTotal + (item.quantity || 0),
        0
      ) || 0),
    0
  );
export const selectCartTotal = (state: { cart: CartState }) =>
  state.cart.cart?.items?.reduce(
    (total, item) => total + Number(item.total_price || 0),
    0
  ) || 0;

// Алиас для fetchUserCarts
export const fetchAllCarts = fetchUserCarts;

export default cartSlice.reducer;
