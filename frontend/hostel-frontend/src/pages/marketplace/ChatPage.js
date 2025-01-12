import React, { useState, useEffect, useCallback, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Container,
    Grid,
    Box,
    useTheme,
    useMediaQuery,
    CircularProgress,
    Alert,
} from '@mui/material';
import {
    ChatWindow,
    ChatList,
    ChatHeader,
    EmptyState,
} from '../../components/marketplace/chat/ChatComponents';
import ChatService from '../../components/marketplace/chat/ChatService';
import axios from '../../api/axios';
import { useAuth } from '../../contexts/AuthContext';

const ChatPage = () => {
    const { listingId } = useParams();
    const navigate = useNavigate();
    const { user } = useAuth();
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));

    const [chats, setChats] = useState([]);

    const [selectedChat, setSelectedChat] = useState(null);
    const [messages, setMessages] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    // Используем useRef для хранения экземпляра ChatService
    const chatServiceRef = useRef(null);

    // Инициализация ChatService
    useEffect(() => {
        if (!user?.id) {
            setError('Необходима авторизация');
            setLoading(false);
            return;
        }

        chatServiceRef.current = new ChatService(user.id);
        let messageHandler = null; // Функция-обработчик сообщений

        const initChat = async () => {
            try {
                await fetchChats();
                chatServiceRef.current.connect();

                // Сохраняем функцию отписки
                messageHandler = chatServiceRef.current.onMessage((message) => {
                    if (message.error) {
                        console.error('Ошибка сообщения:', message.error);
                        return;
                    }
                    console.log('Получено новое сообщение:', message);
                    setMessages(prev => [...prev, message].sort((a, b) =>
                        new Date(a.created_at) - new Date(b.created_at)
                    ));
                });

            } catch (error) {
                console.error('Error initializing chat:', error);
                setError('Ошибка при инициализации чата');
            }
        };

        initChat();

        // Очистка при размонтировании
        return () => {
            if (messageHandler) {
                messageHandler(); // Отписываемся от сообщений
            }
            if (chatServiceRef.current) {
                chatServiceRef.current.disconnect();
            }
        };
    }, [user?.id]);

    // Загрузка списка чатов
    const fetchChats = useCallback(async () => {
        try {
            const response = await axios.get('/api/v1/marketplace/chat');
            setChats(response.data.data);

            // Если есть listingId в URL, находим соответствующий чат
            if (listingId) {
                const chat = response.data.data.find(c => c.listing_id === parseInt(listingId));
                if (chat) {
                    setSelectedChat(chat);
                }
            }
        } catch (error) {
            console.error('Error fetching chats:', error);
            setError('Не удалось загрузить список чатов');
        } finally {
            setLoading(false);
        }
    }, [listingId]);

    // Загрузка сообщений для выбранного чата
    const fetchMessages = useCallback(async (chatId) => {
        try {
            const response = await axios.get(`/api/v1/marketplace/chat/${chatId}/messages`);
            const messages = response.data?.data || [];
            console.log('Получены сообщения с сервера:', messages);
            if (Array.isArray(messages)) {
                setMessages(messages.sort((a, b) =>
                    new Date(a.created_at) - new Date(b.created_at)
                ));
            } else {
                console.error('Неверный формат данных:', messages);
                setMessages([]);
            }
        } catch (error) {
            console.error('Error fetching messages:', error);
            setError('Не удалось загрузить сообщения');
            setMessages([]);
        }
    }, []);

    const handleSelectChat = useCallback((chat) => {
        setSelectedChat(chat);
        if (chatServiceRef.current) {
            chatServiceRef.current.setCurrentChat(chat.id);
        }
        fetchMessages(chat.id);
    }, [fetchMessages]);


    useEffect(() => {
        fetchChats();
    }, [fetchChats]);
    useEffect(() => {
        if (selectedChat && messages.length > 0) {
            const unreadMessages = messages.filter(
                msg => !msg.is_read && msg.receiver_id === user.id
            );
            if (unreadMessages.length > 0) {
                const messageIds = unreadMessages.map(msg => msg.id);
                markMessageAsRead(messageIds);
            }
        }
    }, [selectedChat, messages, user.id]);
    useEffect(() => {
        if (selectedChat) {
            fetchMessages(selectedChat.id);
        }
    }, [selectedChat, fetchMessages]);

    // Отправка сообщения через ChatService
    const handleSendMessage = async (content) => {
        if (!chatServiceRef.current) {
            setError('Соединение не установлено');
            return;
        }

        try {
            const message = {
                chat_id: selectedChat.id,
                listing_id: selectedChat.listing_id,
                receiver_id: selectedChat.seller_id === user.id ?
                    selectedChat.buyer_id : selectedChat.seller_id,
                content: content
            };

            await chatServiceRef.current.sendMessage(message);

            // Добавляем сообщение локально
            const newMessage = {
                ...message,
                sender_id: user.id,
                created_at: new Date().toISOString(),
                is_read: false
            };
            setMessages(prev => [...prev, newMessage]);

        } catch (error) {
            console.error('Error sending message:', error);
            setError('Не удалось отправить сообщение');
        }
    };

    const markMessageAsRead = async (messageIds) => {
        try {
            await axios.put('/api/v1/marketplace/chat/messages/read', {
                message_ids: messageIds
            });
            setMessages(prev => prev.map(msg =>
                messageIds.includes(msg.id) ? { ...msg, is_read: true } : msg
            ));
        } catch (error) {
            console.error('Error marking messages as read:', error);
        }
    };

    const handleArchiveChat = async (chatId) => {
        try {
            await axios.post(`/api/v1/marketplace/chat/${chatId}/archive`);
            fetchChats();
            if (selectedChat?.id === chatId) {
                setSelectedChat(null);
                setMessages([]);
            }
        } catch (error) {
            console.error('Error archiving chat:', error);
            setError('Не удалось архивировать чат');
        }
    };

    if (loading) {
        return (
            <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
                <CircularProgress />
            </Box>
        );
    }

    // Отображение ошибок
    if (error) {
        return (
            <Box p={2}>
                <Alert severity="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            </Box>
        );
    }


    // Мобильная версия
    if (isMobile) {
        return (
            <Box sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
                {selectedChat ? (
                    <>
                        <ChatHeader
                            chat={selectedChat}
                            onBack={() => setSelectedChat(null)}
                        />
                        <Box sx={{ flex: 1, overflow: 'hidden' }}>
                            <ChatWindow
                                messages={messages}
                                onSendMessage={handleSendMessage}
                                currentUser={user}
                            />
                        </Box>
                    </>
                ) : (
                    <>
                        <ChatList
                            chats={chats}
                            selectedChatId={selectedChat?.id}
                            onSelectChat={handleSelectChat}
                            onArchiveChat={handleArchiveChat}
                        />
                        {chats.length === 0 && (
                            <EmptyState text="У вас пока нет сообщений" />
                        )}
                    </>
                )}
            </Box>
        );
    }

    // Десктопная версия
    return (
        <Container maxWidth="xl" sx={{ py: 4, height: 'calc(100vh - 64px)' }}>
            <Grid container spacing={2} sx={{ height: '100%' }}>
                {/* Список чатов */}
                <Grid item xs={12} md={4} sx={{ height: '100%' }}>
                    <ChatList
                        chats={chats}
                        selectedChatId={selectedChat?.id}
                        onSelectChat={handleSelectChat}
                        onArchiveChat={handleArchiveChat}
                    />
                    {chats.length === 0 && (
                        <EmptyState text="У вас пока нет сообщений" />
                    )}
                </Grid>

                {/* Окно чата */}
                <Grid item xs={12} md={8} sx={{ height: '100%' }}>
                    {selectedChat ? (
                        <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                            <ChatHeader chat={selectedChat} />
                            <Box sx={{ flex: 1, overflow: 'hidden' }}>
                                <ChatWindow
                                    messages={messages}
                                    onSendMessage={handleSendMessage}
                                    currentUser={user}
                                />
                            </Box>
                        </Box>
                    ) : (
                        <EmptyState text="Выберите чат для начала общения" />
                    )}
                </Grid>
            </Grid>
        </Container>
    );
};

export default ChatPage;