import { configureStore } from '@reduxjs/toolkit';
import chatReducer from './slices/chatSlice';
import reviewsReducer from './slices/reviewsSlice';
import b2cStoresReducer from './slices/b2cStoreSlice';
import importReducer from './slices/importSlice';
import productReducer from './slices/productSlice';
import paymentReducer from './slices/paymentSlice';
import cartReducer from './slices/cartSlice';
import localCartReducer from './slices/localCartSlice';
import categoriesReducer from './slices/categoriesSlice';
import universalCompareReducer from './slices/universalCompareSlice';
import favoritesReducer from './slices/favoritesSlice';
import savedSearchesReducer from './slices/savedSearchesSlice';
import categoryProposalsReducer from './slices/categoryProposalsSlice';
import { websocketMiddleware } from './middleware/websocketMiddleware';

export const store = configureStore({
  reducer: {
    chat: chatReducer,
    reviews: reviewsReducer,
    b2cStores: b2cStoresReducer,
    import: importReducer,
    products: productReducer,
    payment: paymentReducer,
    cart: cartReducer,
    localCart: localCartReducer,
    categories: categoriesReducer,
    universalCompare: universalCompareReducer,
    favorites: favoritesReducer,
    savedSearches: savedSearchesReducer,
    categoryProposals: categoryProposalsReducer,
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
          'import/setAnalysisFile',
          'import/setAnalysisFileType',
          'import/importFromFile/pending',
          'import/importFromFile/fulfilled',
          'import/analyzeImportFile/pending',
          'import/analyzeCategories/pending',
          'import/analyzeAttributes/pending',
          'import/detectVariants/pending',
          'import/analyzeClientCategories/pending',
          'import/fetchJobStatus/pending',
          'import/fetchJobStatus/fulfilled',
        ],
        ignoredPaths: [
          'chat.ws',
          'chat.onlineUsers',
          'chat.typingUsers',
          'import.selectedFiles',
          'import.analysisFile',
          'products.selectedIds',
          'favorites.itemIds',
        ],
      },
    }).concat(websocketMiddleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
