import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import {
  MarketplaceChat,
  MarketplaceMessage,
  SendMessagePayload,
  GetMessagesParams,
  ChatAttachment,
  UploadingFile,
} from '@/types/chat';
import { chatService } from '@/services/chat';
import { tokenManager } from '@/utils/tokenManager';
import type { RootState } from '../index';

interface ChatState {
  // Состояние
  chats: MarketplaceChat[];
  currentChat: MarketplaceChat | null;
  messages: Record<number, MarketplaceMessage[]>; // chatId -> messages
  unreadCount: number;
  isLoading: boolean;
  error: string | null;

  // WebSocket
  ws: WebSocket | null;
  typingUsers: Record<number, Set<number>>; // chatId -> Set of userIds
  onlineUsers: Set<number>;
  userLastSeen: Record<number, string>; // userId -> last seen timestamp
  getCurrentUserId?: () => number; // Функция для получения текущего userId

  // Пагинация
  chatsPage: number;
  messagesPage: Record<number, number>; // chatId -> current page
  hasMoreChats: boolean;
  hasMoreMessages: Record<number, boolean>; // chatId -> hasMore
  messagesLoaded: Record<number, boolean>; // chatId -> были ли загружены сообщения с сервера

  // Загрузка файлов
  uploadingFiles: Record<string, UploadingFile>; // fileId -> uploadingFile
  attachments: Record<number, ChatAttachment[]>; // messageId -> attachments
}

const initialState: ChatState = {
  chats: [],
  currentChat: null,
  messages: {},
  unreadCount: 0,
  isLoading: false,
  error: null,
  ws: null,
  typingUsers: {},
  onlineUsers: new Set(),
  userLastSeen: {},
  getCurrentUserId: undefined,
  chatsPage: 1,
  messagesPage: {},
  hasMoreChats: true,
  hasMoreMessages: {},
  messagesLoaded: {},
  uploadingFiles: {},
  attachments: {},
};

// Async thunks
export const loadChats = createAsyncThunk(
  'chat/loadChats',
  async (page: number = 1) => {
    const hasToken = tokenManager.getAccessToken() !== null;
    if (!hasToken) {
      console.log('[ChatStore] No access token, skipping chat load');
      return { chats: [], total: 0, page: 1, limit: 20 };
    }

    const response = await chatService.getChats(page);
    return response;
  }
);

export const loadMessages = createAsyncThunk(
  'chat/loadMessages',
  async (params: GetMessagesParams) => {
    const response = await chatService.getMessages(params);
    return { ...response, chatId: params.chat_id };
  }
);

export const sendMessage = createAsyncThunk(
  'chat/sendMessage',
  async (payload: SendMessagePayload) => {
    const response = await chatService.sendMessage(payload);
    return response;
  }
);

export const markMessagesAsRead = createAsyncThunk(
  'chat/markMessagesAsRead',
  async ({ chatId, messageIds }: { chatId: number; messageIds: number[] }) => {
    await chatService.markMessagesAsRead({
      chat_id: chatId,
      message_ids: messageIds,
    });
    return { chatId, messageIds };
  }
);

export const archiveChat = createAsyncThunk(
  'chat/archiveChat',
  async (chatId: number) => {
    await chatService.archiveChat(chatId);
    return chatId;
  }
);

export const uploadFiles = createAsyncThunk(
  'chat/uploadFiles',
  async (
    { messageId, files }: { messageId: number; files: File[] },
    { dispatch }
  ) => {
    const uploadingFileIds: string[] = [];

    // Добавляем файлы в состояние загрузки
    files.forEach((file) => {
      const fileId = `${messageId}-${file.name}-${Date.now()}`;
      uploadingFileIds.push(fileId);
      dispatch(
        chatSlice.actions.addUploadingFile({
          fileId,
          file: {
            id: fileId,
            file,
            progress: 0,
            status: 'pending',
          },
        })
      );
    });

    try {
      // Загружаем файлы
      const response = await chatService.uploadFiles(
        messageId,
        files,
        (progress) => {
          uploadingFileIds.forEach((fileId) => {
            dispatch(
              chatSlice.actions.updateUploadProgress({ fileId, progress })
            );
          });
        }
      );

      // Удаляем файлы из загрузки
      uploadingFileIds.forEach((fileId) => {
        dispatch(chatSlice.actions.removeUploadingFile(fileId));
      });

      return { messageId, attachments: response.attachments };
    } catch (error) {
      // Помечаем файлы как ошибочные
      uploadingFileIds.forEach((fileId) => {
        dispatch(
          chatSlice.actions.setUploadError({
            fileId,
            error: error instanceof Error ? error.message : 'Upload failed',
          })
        );
      });
      throw error;
    }
  }
);

export const deleteAttachment = createAsyncThunk(
  'chat/deleteAttachment',
  async (attachmentId: number) => {
    await chatService.deleteAttachment(attachmentId);
    return attachmentId;
  }
);

const chatSlice = createSlice({
  name: 'chat',
  initialState,
  reducers: {
    setCurrentChat: (state, action: PayloadAction<MarketplaceChat | null>) => {
      state.currentChat = action.payload;
    },

    selectLatestChat: (state) => {
      if (state.chats.length > 0) {
        state.currentChat = state.chats[0];
      }
    },

    setWebSocket: (state, action: PayloadAction<WebSocket | null>) => {
      state.ws = action.payload;
    },

    setUserTyping: (
      state,
      action: PayloadAction<{
        chatId: number;
        userId: number;
        isTyping: boolean;
      }>
    ) => {
      const { chatId, userId, isTyping } = action.payload;
      if (!state.typingUsers[chatId]) {
        state.typingUsers[chatId] = new Set();
      }
      if (isTyping) {
        state.typingUsers[chatId].add(userId);
      } else {
        state.typingUsers[chatId].delete(userId);
      }
    },

    addUploadingFile: (
      state,
      action: PayloadAction<{ fileId: string; file: UploadingFile }>
    ) => {
      state.uploadingFiles[action.payload.fileId] = action.payload.file;
    },

    updateUploadProgress: (
      state,
      action: PayloadAction<{ fileId: string; progress: number }>
    ) => {
      if (state.uploadingFiles[action.payload.fileId]) {
        state.uploadingFiles[action.payload.fileId].progress =
          action.payload.progress;
        state.uploadingFiles[action.payload.fileId].status = 'uploading';
      }
    },

    setUploadError: (
      state,
      action: PayloadAction<{ fileId: string; error: string }>
    ) => {
      if (state.uploadingFiles[action.payload.fileId]) {
        state.uploadingFiles[action.payload.fileId].status = 'error';
        state.uploadingFiles[action.payload.fileId].error =
          action.payload.error;
      }
    },

    removeUploadingFile: (state, action: PayloadAction<string>) => {
      delete state.uploadingFiles[action.payload];
    },

    setGetCurrentUserId: (state, action: PayloadAction<() => number>) => {
      state.getCurrentUserId = action.payload;
    },

    // WebSocket события
    handleNewMessage: (state, action: PayloadAction<MarketplaceMessage>) => {
      const message = action.payload;

      // Добавляем сообщение в список
      if (!state.messages[message.chat_id]) {
        state.messages[message.chat_id] = [];
      }

      // Проверяем, что сообщение еще не добавлено
      const exists = state.messages[message.chat_id].some(
        (msg) => msg.id === message.id
      );
      if (!exists) {
        state.messages[message.chat_id].push(message);
      }

      // Обновляем чат
      const chatIndex = state.chats.findIndex(
        (chat) => chat.id === message.chat_id
      );
      if (chatIndex !== -1) {
        state.chats[chatIndex].last_message = message;
        state.chats[chatIndex].last_message_at = message.created_at;

        // Увеличиваем счетчик непрочитанных, если сообщение не от текущего пользователя
        if (
          state.getCurrentUserId &&
          message.sender_id !== state.getCurrentUserId()
        ) {
          state.chats[chatIndex].unread_count =
            (state.chats[chatIndex].unread_count || 0) + 1;
        }
      }

      // Сортируем чаты по дате последнего сообщения
      state.chats.sort((a, b) => {
        const aTime = new Date(a.last_message_at || a.created_at).getTime();
        const bTime = new Date(b.last_message_at || b.created_at).getTime();
        return bTime - aTime;
      });
    },

    handleMessageRead: (
      state,
      action: PayloadAction<{
        chat_id: number;
        message_ids: number[];
        reader_id: number;
      }>
    ) => {
      const { chat_id, message_ids } = action.payload;

      // Обновляем статус прочтения сообщений
      if (state.messages[chat_id]) {
        state.messages[chat_id] = state.messages[chat_id].map((msg) =>
          message_ids.includes(msg.id) ? { ...msg, is_read: true } : msg
        );
      }

      // Обновляем счетчик непрочитанных
      const chatIndex = state.chats.findIndex((chat) => chat.id === chat_id);
      if (chatIndex !== -1) {
        state.chats[chatIndex].unread_count = 0;
      }
    },

    handleUserOnline: (state, action: PayloadAction<{ user_id: number }>) => {
      state.onlineUsers.add(action.payload.user_id);
    },

    handleUserOffline: (
      state,
      action: PayloadAction<{ user_id: number; last_seen?: string }>
    ) => {
      state.onlineUsers.delete(action.payload.user_id);
      if (action.payload.last_seen) {
        state.userLastSeen[action.payload.user_id] = action.payload.last_seen;
      }
    },

    reset: (state) => {
      Object.assign(state, initialState);
    },
  },

  extraReducers: (builder) => {
    // loadChats
    builder
      .addCase(loadChats.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(loadChats.fulfilled, (state, action) => {
        state.isLoading = false;
        if (action.meta.arg === 1) {
          state.chats = action.payload.chats;
        } else {
          state.chats = [...state.chats, ...action.payload.chats];
        }
        state.chatsPage = action.payload.page;
        state.hasMoreChats =
          action.payload.chats.length === action.payload.limit;
        state.unreadCount = action.payload.chats.reduce(
          (sum, chat) => sum + (chat.unread_count || 0),
          0
        );
      })
      .addCase(loadChats.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to load chats';
      });

    // loadMessages
    builder
      .addCase(loadMessages.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(loadMessages.fulfilled, (state, action) => {
        state.isLoading = false;
        const chatId = action.payload.chatId!;

        if (!state.messagesPage[chatId] || state.messagesPage[chatId] === 1) {
          state.messages[chatId] = action.payload.messages;
        } else {
          const existingIds = new Set(
            state.messages[chatId]?.map((m) => m.id) || []
          );
          const newMessages = action.payload.messages.filter(
            (m) => !existingIds.has(m.id)
          );
          state.messages[chatId] = [
            ...newMessages,
            ...(state.messages[chatId] || []),
          ];
        }

        state.messagesPage[chatId] = action.payload.page;
        state.hasMoreMessages[chatId] =
          action.payload.messages.length === action.payload.limit;
        state.messagesLoaded[chatId] = true;
      })
      .addCase(loadMessages.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.error.message || 'Failed to load messages';
      });

    // sendMessage
    builder.addCase(sendMessage.fulfilled, (state, action) => {
      const message = action.payload;

      // Добавляем сообщение
      if (!state.messages[message.chat_id]) {
        state.messages[message.chat_id] = [];
      }
      state.messages[message.chat_id].push(message);

      // Обновляем чат
      const chatIndex = state.chats.findIndex(
        (chat) => chat.id === message.chat_id
      );
      if (chatIndex !== -1) {
        state.chats[chatIndex].last_message = message;
        state.chats[chatIndex].last_message_at = message.created_at;
      }
    });

    // markMessagesAsRead
    builder.addCase(markMessagesAsRead.fulfilled, (state, action) => {
      const { chatId, messageIds } = action.payload;

      if (state.messages[chatId]) {
        state.messages[chatId] = state.messages[chatId].map((msg) =>
          messageIds.includes(msg.id) ? { ...msg, is_read: true } : msg
        );
      }

      const chatIndex = state.chats.findIndex((chat) => chat.id === chatId);
      if (chatIndex !== -1) {
        state.chats[chatIndex].unread_count = 0;
      }
    });

    // archiveChat
    builder.addCase(archiveChat.fulfilled, (state, action) => {
      state.chats = state.chats.filter((chat) => chat.id !== action.payload);
      if (state.currentChat?.id === action.payload) {
        state.currentChat = null;
      }
    });

    // uploadFiles
    builder.addCase(uploadFiles.fulfilled, (state, action) => {
      const { messageId, attachments } = action.payload;
      state.attachments[messageId] = attachments;

      // Обновляем сообщение с флагом has_attachments
      Object.values(state.messages).forEach((messages) => {
        const message = messages.find((m) => m.id === messageId);
        if (message) {
          message.has_attachments = true;
          message.attachments = attachments;
          message.attachments_count = attachments.length;
        }
      });
    });

    // deleteAttachment
    builder.addCase(deleteAttachment.fulfilled, (state, action) => {
      // Удаляем вложение из всех сообщений
      Object.values(state.messages).forEach((messages) => {
        messages.forEach((message) => {
          if (message.attachments) {
            message.attachments = message.attachments.filter(
              (att) => att.id !== action.payload
            );
            message.attachments_count = message.attachments.length;
            message.has_attachments = message.attachments.length > 0;
          }
        });
      });
    });
  },
});

export const {
  setCurrentChat,
  selectLatestChat,
  setWebSocket,
  setUserTyping,
  addUploadingFile,
  updateUploadProgress,
  setUploadError,
  removeUploadingFile,
  setGetCurrentUserId,
  handleNewMessage,
  handleMessageRead,
  handleUserOnline,
  handleUserOffline,
  reset,
} = chatSlice.actions;

// Selectors
export const selectChats = (state: RootState) => state.chat.chats;
export const selectCurrentChat = (state: RootState) => state.chat.currentChat;
export const selectMessages = (state: RootState, chatId?: number) =>
  chatId ? state.chat.messages[chatId] || [] : [];
export const selectUnreadCount = (state: RootState) => state.chat.unreadCount;
export const selectIsLoading = (state: RootState) => state.chat.isLoading;
export const selectError = (state: RootState) => state.chat.error;
export const selectOnlineUsers = (state: RootState) => state.chat.onlineUsers;
export const selectTypingUsers = (state: RootState, chatId: number) =>
  state.chat.typingUsers[chatId] || new Set();
export const selectUploadingFiles = (state: RootState) =>
  state.chat.uploadingFiles;

export default chatSlice.reducer;
