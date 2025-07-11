import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { cartService } from '@/services/cart';
import type { components } from '@/types/generated/api';

type ShoppingCart =
  components['schemas']['backend_internal_domain_models.ShoppingCart'];
type AddToCartRequest =
  components['schemas']['backend_internal_domain_models.AddToCartRequest'];
type UpdateCartItemRequest =
  components['schemas']['backend_internal_domain_models.UpdateCartItemRequest'];

interface CartState {
  cart: ShoppingCart | null;
  loading: boolean;
  error: string | null;
}

const initialState: CartState = {
  cart: null,
  loading: false,
  error: null,
};

// Async thunks
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
      state.error = null;
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
        state.cart = action.payload;
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
      .addCase(clearCart.fulfilled, (state) => {
        state.loading = false;
        state.cart = null;
      })
      .addCase(clearCart.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to clear cart';
      });
  },
});

export const { clearError, resetCart } = cartSlice.actions;

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
export const selectCartTotal = (state: { cart: CartState }) =>
  state.cart.cart?.items?.reduce(
    (total, item) => total + Number(item.total_price || 0),
    0
  ) || 0;

export default cartSlice.reducer;
