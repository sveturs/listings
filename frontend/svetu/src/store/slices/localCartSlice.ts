import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface CartItem {
  productId: number;
  variantId?: number;
  name: string;
  variantName?: string;
  price: number;
  currency: string;
  quantity: number;
  stockQuantity?: number;
  image?: string;
  storefrontId: number;
  storefrontName?: string;
  storefrontSlug?: string;
}

interface LocalCartState {
  items: CartItem[];
  lastUpdated: number;
}

const CART_STORAGE_KEY = 'svetu_cart';

// Load cart from localStorage
const loadCartFromStorage = (): CartItem[] => {
  if (typeof window === 'undefined') return [];

  try {
    const stored = localStorage.getItem(CART_STORAGE_KEY);
    if (stored) {
      const { items, lastUpdated } = JSON.parse(stored);
      // Clear cart if it's older than 7 days
      if (Date.now() - lastUpdated > 7 * 24 * 60 * 60 * 1000) {
        localStorage.removeItem(CART_STORAGE_KEY);
        return [];
      }
      return items;
    }
  } catch (error) {
    console.error('Failed to load cart from storage:', error);
  }
  return [];
};

// Save cart to localStorage
const saveCartToStorage = (items: CartItem[]) => {
  if (typeof window === 'undefined') return;

  try {
    localStorage.setItem(
      CART_STORAGE_KEY,
      JSON.stringify({
        items,
        lastUpdated: Date.now(),
      })
    );
  } catch (error) {
    console.error('Failed to save cart to storage:', error);
  }
};

const initialState: LocalCartState = {
  items: loadCartFromStorage(),
  lastUpdated: Date.now(),
};

const localCartSlice = createSlice({
  name: 'localCart',
  initialState,
  reducers: {
    addItem: (state, action: PayloadAction<CartItem>) => {
      const newItem = action.payload;
      const existingItemIndex = state.items.findIndex(
        (item) =>
          item.productId === newItem.productId &&
          item.variantId === newItem.variantId &&
          item.storefrontId === newItem.storefrontId
      );

      if (existingItemIndex >= 0) {
        // Update quantity if item already exists
        state.items[existingItemIndex].quantity += newItem.quantity;
      } else {
        // Add new item
        state.items.push(newItem);
      }

      state.lastUpdated = Date.now();
      saveCartToStorage(state.items);
    },

    updateQuantity: (
      state,
      action: PayloadAction<{
        productId: number;
        variantId?: number;
        quantity: number;
      }>
    ) => {
      const { productId, variantId, quantity } = action.payload;
      const itemIndex = state.items.findIndex(
        (item) => item.productId === productId && item.variantId === variantId
      );

      if (itemIndex >= 0 && quantity > 0) {
        state.items[itemIndex].quantity = quantity;
        state.lastUpdated = Date.now();
        saveCartToStorage(state.items);
      }
    },

    removeItem: (
      state,
      action: PayloadAction<{ productId: number; variantId?: number }>
    ) => {
      const { productId, variantId } = action.payload;
      state.items = state.items.filter(
        (item) =>
          !(item.productId === productId && item.variantId === variantId)
      );
      state.lastUpdated = Date.now();
      saveCartToStorage(state.items);
    },

    clearCart: (state) => {
      state.items = [];
      state.lastUpdated = Date.now();
      saveCartToStorage([]);
    },

    syncCart: (state) => {
      // Reload cart from storage (useful after login)
      state.items = loadCartFromStorage();
      state.lastUpdated = Date.now();
    },
  },
});

export const { addItem, updateQuantity, removeItem, clearCart, syncCart } =
  localCartSlice.actions;

// Re-export clearCart as clearLocalCart for consistency with imports
export const clearLocalCart = clearCart;

// Selectors
export const selectLocalCart = (state: { localCart: LocalCartState }) =>
  state.localCart;

export const selectCartItems = (state: { localCart: LocalCartState }) =>
  state.localCart.items;

export const selectCartItemsCount = (state: { localCart: LocalCartState }) =>
  state.localCart.items.reduce((total, item) => total + item.quantity, 0);

export const selectCartTotal = (state: { localCart: LocalCartState }) =>
  state.localCart.items.reduce(
    (total, item) => total + item.price * item.quantity,
    0
  );

export const selectCartByStorefront = (state: {
  localCart: LocalCartState;
}) => {
  const itemsByStorefront: Record<number, CartItem[]> = {};

  state.localCart.items.forEach((item) => {
    if (!itemsByStorefront[item.storefrontId]) {
      itemsByStorefront[item.storefrontId] = [];
    }
    itemsByStorefront[item.storefrontId].push(item);
  });

  return itemsByStorefront;
};

export default localCartSlice.reducer;
