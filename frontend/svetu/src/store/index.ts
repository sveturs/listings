import { configureStore } from '@reduxjs/toolkit';
import chatReducer from './slices/chatSlice';
import reviewsReducer from './slices/reviewsSlice';
import storefrontsReducer from './slices/storefrontSlice';
import importReducer from './slices/importSlice';
import productReducer from './slices/productSlice';
import paymentReducer from './slices/paymentSlice';
import cartReducer from './slices/cartSlice';
import localCartReducer from './slices/localCartSlice';
import categoriesReducer from './slices/categoriesSlice';
import compareReducer from './slices/compareSlice';
import favoritesReducer from './slices/favoritesSlice';
import savedSearchesReducer from './slices/savedSearchesSlice';
import { websocketMiddleware } from './middleware/websocketMiddleware';

export const store = configureStore({
  reducer: {
    chat: chatReducer,
    reviews: reviewsReducer,
    storefronts: storefrontsReducer,
    import: importReducer,
    products: productReducer,
    payment: paymentReducer,
    cart: cartReducer,
    localCart: localCartReducer,
    categories: categoriesReducer,
    compare: compareReducer,
    favorites: favoritesReducer,
    savedSearches: savedSearchesReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // Игнорируем проверку для WebSocket, Set и File объектов
        ignoredActions: [
          'chat/initWebSocket',
          'chat/setWebSocket',
          'chat/handleUserOnline',
          'chat/handleUserOffline',
          'import/setSelectedFiles',
          'import/importFromFile/pending',
          'import/importFromFile/fulfilled',
        ],
        ignoredPaths: [
          'chat.ws',
          'chat.onlineUsers',
          'chat.typingUsers',
          'import.selectedFiles',
          'products.selectedIds',
          'favorites.itemIds',
        ],
      },
    }).concat(websocketMiddleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
