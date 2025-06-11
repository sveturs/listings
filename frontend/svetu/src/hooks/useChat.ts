import { useCallback } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  loadChats,
  loadMessages,
  sendMessage,
  markMessagesAsRead,
  archiveChat,
  setCurrentChat,
  selectLatestChat,
  uploadFiles,
  deleteAttachment,
  removeUploadingFile,
  setUserTyping,
  selectChats,
  selectCurrentChat,
  selectUnreadCount,
  selectIsLoading,
  selectError,
  selectOnlineUsers,
  selectUploadingFiles,
} from '@/store/slices/chatSlice';
import type {
  MarketplaceChat,
  SendMessagePayload,
  GetMessagesParams,
} from '@/types/chat';

export function useChat() {
  const dispatch = useAppDispatch();

  // Selectors
  const chats = useAppSelector(selectChats);
  const currentChat = useAppSelector(selectCurrentChat);
  const unreadCount = useAppSelector(selectUnreadCount);
  const isLoading = useAppSelector(selectIsLoading);
  const error = useAppSelector(selectError);
  const onlineUsers = useAppSelector(selectOnlineUsers);
  const uploadingFiles = useAppSelector(selectUploadingFiles);

  // Actions
  const loadChatsAction = useCallback(
    (page?: number) => {
      return dispatch(loadChats(page || 1));
    },
    [dispatch]
  );

  const loadMessagesAction = useCallback(
    (params: GetMessagesParams) => {
      return dispatch(loadMessages(params));
    },
    [dispatch]
  );

  const sendMessageAction = useCallback(
    (payload: SendMessagePayload) => {
      return dispatch(sendMessage(payload));
    },
    [dispatch]
  );

  const markMessagesAsReadAction = useCallback(
    (chatId: number, messageIds: number[]) => {
      return dispatch(markMessagesAsRead({ chatId, messageIds }));
    },
    [dispatch]
  );

  const archiveChatAction = useCallback(
    (chatId: number) => {
      return dispatch(archiveChat(chatId));
    },
    [dispatch]
  );

  const setCurrentChatAction = useCallback(
    (chat: MarketplaceChat | null) => {
      dispatch(setCurrentChat(chat));
    },
    [dispatch]
  );

  const selectLatestChatAction = useCallback(() => {
    dispatch(selectLatestChat());
  }, [dispatch]);

  const uploadFilesAction = useCallback(
    (messageId: number, files: File[]) => {
      return dispatch(uploadFiles({ messageId, files }));
    },
    [dispatch]
  );

  const deleteAttachmentAction = useCallback(
    (attachmentId: number) => {
      return dispatch(deleteAttachment(attachmentId));
    },
    [dispatch]
  );

  const removeUploadingFileAction = useCallback(
    (fileId: string) => {
      dispatch(removeUploadingFile(fileId));
    },
    [dispatch]
  );

  const setUserTypingAction = useCallback(
    (chatId: number, userId: number, isTyping: boolean) => {
      dispatch(setUserTyping({ chatId, userId, isTyping }));
    },
    [dispatch]
  );

  const initWebSocket = useCallback(
    (getCurrentUserId: () => number) => {
      dispatch({ type: 'chat/initWebSocket', payload: { getCurrentUserId } });
    },
    [dispatch]
  );

  const closeWebSocket = useCallback(() => {
    dispatch({ type: 'chat/closeWebSocket' });
  }, [dispatch]);

  // Additional selectors for specific chat data
  const messages = useAppSelector((state) => state.chat.messages);
  const typingUsers = useAppSelector((state) => state.chat.typingUsers);
  const hasMoreMessages = useAppSelector((state) => state.chat.hasMoreMessages);
  const userLastSeen = useAppSelector((state) => state.chat.userLastSeen);
  const hasMoreChats = useAppSelector((state) => state.chat.hasMoreChats);
  const messagesLoaded = useAppSelector((state) => state.chat.messagesLoaded);

  return {
    // State
    chats,
    currentChat,
    unreadCount,
    isLoading,
    error,
    onlineUsers,
    uploadingFiles,
    messages,
    typingUsers,
    hasMoreMessages,
    userLastSeen,
    hasMoreChats,
    messagesLoaded,

    // Actions
    loadChats: loadChatsAction,
    loadMessages: loadMessagesAction,
    sendMessage: sendMessageAction,
    markMessagesAsRead: markMessagesAsReadAction,
    archiveChat: archiveChatAction,
    setCurrentChat: setCurrentChatAction,
    selectLatestChat: selectLatestChatAction,
    uploadFiles: uploadFilesAction,
    deleteAttachment: deleteAttachmentAction,
    removeUploadingFile: removeUploadingFileAction,
    setUserTyping: setUserTypingAction,
    initWebSocket,
    closeWebSocket,

    // Helpers
    getMessages: (chatId?: number) => {
      return {
        messages: chatId ? messages[chatId] || [] : [],
        hasMore: chatId ? hasMoreMessages[chatId] || false : false,
      };
    },
    getTypingUsers: (chatId: number) => {
      return typingUsers[chatId] || [];
    },
  };
}

// Экспорт для обратной совместимости
export const useChatStore = useChat;
