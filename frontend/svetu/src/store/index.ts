import { configureStore } from '@reduxjs/toolkit';
import chatReducer from './slices/chatSlice';
import reviewsReducer from './slices/reviewsSlice';
import { websocketMiddleware } from './middleware/websocketMiddleware';

export const store = configureStore({
  reducer: {
    chat: chatReducer,
    reviews: reviewsReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // Игнорируем проверку для WebSocket и Set
        ignoredActions: [
          'chat/initWebSocket',
          'chat/setWebSocket',
          'chat/handleUserOnline',
          'chat/handleUserOffline',
        ],
        ignoredPaths: ['chat.ws', 'chat.onlineUsers', 'chat.typingUsers'],
      },
    }).concat(websocketMiddleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
