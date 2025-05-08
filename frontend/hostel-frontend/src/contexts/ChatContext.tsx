// src/contexts/ChatContext.tsx
import React, { createContext, useContext, useRef, useCallback, ReactNode } from 'react';
import ChatService from '../components/marketplace/chat/ChatService';

interface ChatContextType {
  getChatService: (userId: number) => ChatService | null;
}

interface ChatProviderProps {
  children: ReactNode;
}

const ChatContext = createContext<ChatContextType | null>(null);

export const ChatProvider: React.FC<ChatProviderProps> = ({ children }) => {
  const serviceRef = useRef<ChatService | null>(null);
  const userIdRef = useRef<number | null>(null);

  const getChatService = useCallback((userId: number) => {
    if (!userId) return null;
    
    // Если сервис уже существует для данного пользователя, возвращаем его
    if (serviceRef.current && userIdRef.current === userId) {
      return serviceRef.current;
    }

    // Если сервис существует, но для другого пользователя, отключаем его
    if (serviceRef.current) {
      serviceRef.current.disconnect();
    }

    // Создаем новый сервис
    serviceRef.current = new ChatService(userId);
    userIdRef.current = userId;
    
    return serviceRef.current;
  }, []);

  return (
    <ChatContext.Provider value={{ getChatService }}>
      {children}
    </ChatContext.Provider>
  );
};

export const useChat = (): ChatContextType => {
  const context = useContext(ChatContext);
  if (!context) {
    throw new Error('useChat must be used within a ChatProvider');
  }
  return context;
};