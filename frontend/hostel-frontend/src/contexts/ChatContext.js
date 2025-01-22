// src/contexts/ChatContext.js
import React, { createContext, useContext, useRef, useCallback } from 'react';
import ChatService from '../components/marketplace/chat/ChatService';

const ChatContext = createContext(null);

export const ChatProvider = ({ children }) => {
    const serviceRef = useRef(null);
    const userIdRef = useRef(null);

    const getChatService = useCallback((userId) => {
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

export const useChat = () => {
    const context = useContext(ChatContext);
    if (!context) {
        throw new Error('useChat must be used within a ChatProvider');
    }
    return context;
};