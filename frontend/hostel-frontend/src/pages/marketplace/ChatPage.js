// frontend/hostel-frontend/src/pages/marketplace/ChatPage.js
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

    const chatServiceRef = useRef(null);
    const messageEndRef = useRef(null);

    // Инициализация чат-сервиса
    useEffect(() => {
        if (!user?.id) {
            setError('Необходима авторизация');
            setLoading(false);
            return;
        }

        chatServiceRef.current = new ChatService(user.id);

        return () => {
            if (chatServiceRef.current) {
                chatServiceRef.current.disconnect();
            }
        };
    }, [user?.id]);

    // Загрузка списка чатов
    const fetchChats = useCallback(async () => {
        try {
            const response = await axios.get('/api/v1/marketplace/chat');
            const chatsData = response.data?.data || [];
            setChats(chatsData);

            // Если есть listingId в URL, выбираем соответствующий чат
            if (listingId) {
                const chat = chatsData.find(c => c.listing_id === parseInt(listingId));
                if (chat) {
                    setSelectedChat(chat);
                }
            }
        } catch (error) {
            console.error('Ошибка загрузки чатов:', error);
            setError('Не удалось загрузить список чатов');
        } finally {
            setLoading(false);
        }
    }, [listingId]);

    // Загрузка сообщений выбранного чата
    const fetchMessages = useCallback(async (chatId) => {
        if (!chatServiceRef.current) {
            throw new Error('ChatService не инициализирован');
        }

        try {
            const messages = await chatServiceRef.current.getMessageHistory(chatId);
            return messages;
        } catch (error) {
            console.error('Ошибка загрузки сообщений:', error);
            throw error;
        }
    }, []);

    // Обработка выбора чата
    const handleSelectChat = useCallback((chat) => {
        console.log('Выбран чат:', chat);
        setSelectedChat(chat);
        if (chat?.listing_id) { // Используем listing_id
            setLoading(true);
            setMessages([]);

            fetchMessages(chat.id, chat.listing_id) // Передаем оба параметра
                .then(messages => {
                    console.log('Загружены сообщения:', messages);
                    if (Array.isArray(messages)) {
                        setMessages(messages);
                    }
                })
                .catch(error => {
                    console.error('Ошибка загрузки сообщений:', error);
                    setError('Не удалось загрузить сообщения');
                })
                .finally(() => {
                    setLoading(false);
                });
        }
    }, [fetchMessages]);

    // Инициализация WebSocket и загрузка данных
    useEffect(() => {
        if (chatServiceRef.current) {
            chatServiceRef.current.connect();

            const unsubscribe = chatServiceRef.current.onMessage((message) => {
                console.log('Получено новое сообщение:', message);

                // Проверяем, относится ли сообщение к текущему чату
                if (selectedChat && message.chat_id === selectedChat.id) {
                    setMessages(prev => {
                        // Проверяем наличие дубликата
                        if (prev.some(m => m.id === message.id)) {
                            return prev;
                        }

                        // Добавляем новое сообщение и сортируем
                        const updatedMessages = [...prev, {
                            ...message,
                            sender: message.sender || {},
                            receiver: message.receiver || {},
                            is_read: message.is_read || false,
                            created_at: message.created_at || new Date().toISOString()
                        }];

                        return updatedMessages.sort((a, b) =>
                            new Date(a.created_at) - new Date(b.created_at)
                        );
                    });
                }
            });

            return () => unsubscribe();
        }
    }, [selectedChat]);

    // Загрузка чатов при монтировании
    useEffect(() => {
        fetchChats();
    }, [fetchChats]);
    useEffect(() => {
        if (selectedChat?.id && selectedChat?.listing_id && messages.length === 0) {
            console.log('Начальная загрузка сообщений для чата:', selectedChat.id);
            
            setLoading(true);
            
            chatServiceRef.current.getMessageHistory(selectedChat.id, selectedChat.listing_id)
                .then(loadedMessages => {
                    console.log('Загружены начальные сообщения:', loadedMessages);
                    setMessages(loadedMessages || []);
                })
                .catch(error => {
                    console.error('Ошибка начальной загрузки сообщений:', error);
                    setError('Не удалось загрузить сообщения');
                })
                .finally(() => {
                    setLoading(false);
                });
        }
    }, [selectedChat?.id]);

    useEffect(() => {
        if (selectedChat && messages.length > 0) {
            const unreadMessages = messages.filter(
                msg => !msg.is_read && msg.receiver_id === user?.id
            );

            if (unreadMessages.length > 0) {
                const messageIds = unreadMessages.map(msg => msg.id);
                chatServiceRef.current?.markMessagesAsRead(messageIds);
            }
        }
    }, [selectedChat, messages, user?.id]);

    // Прокрутка к последнему сообщению
    useEffect(() => {
        messageEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);

    // Отправка сообщения
    const handleSendMessage = async (content) => {
        if (!selectedChat || !user?.id) {
            setError('Недостаточно данных для отправки сообщения');
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
        } catch (error) {
            console.error('Ошибка отправки сообщения:', error);
            setError('Не удалось отправить сообщение');
        }
    };

    // Архивация чата
    const handleArchiveChat = async (chatId) => {
        try {
            await axios.post(`/api/v1/marketplace/chat/${chatId}/archive`);
            await fetchChats();
            if (selectedChat?.id === chatId) {
                setSelectedChat(null);
                setMessages([]);
            }
        } catch (error) {
            console.error('Ошибка архивации чата:', error);
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

    if (!user) {
        return (
            <Container maxWidth="md" sx={{ mt: 4 }}>
                <Alert severity="warning">
                    Для доступа к чату необходимо авторизоваться
                </Alert>
            </Container>
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
                        {!loading && chats.length === 0 && (
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
                    {!loading && chats.length === 0 && (
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
                                <div ref={messageEndRef} />
                            </Box>
                        </Box>
                    ) : (
                        <EmptyState text="Выберите чат для начала общения" />
                    )}
                </Grid>
            </Grid>

            {error && (
                <Alert
                    severity="error"
                    sx={{
                        position: 'fixed',
                        bottom: 16,
                        right: 16,
                        maxWidth: 'calc(100% - 32px)'
                    }}
                    onClose={() => setError(null)}
                >
                    {error}
                </Alert>
            )}
        </Container>
    );
};

export default ChatPage;