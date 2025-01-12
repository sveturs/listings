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

    const [chats, setChats] = useState([]);  // Инициализируем пустым массивом вместо null
    const [selectedChat, setSelectedChat] = useState(null);
    const [messages, setMessages] = useState([]); // Инициализируем пустым массивом
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const chatServiceRef = useRef(null);

    useEffect(() => {
        if (!user?.id) {
            setError('Необходима авторизация');
            setLoading(false);
            return;
        }

        chatServiceRef.current = new ChatService(user.id);
        
        const initChat = async () => {
            try {
                await fetchChats();
                chatServiceRef.current.connect();

                const messageHandler = chatServiceRef.current.onMessage((message) => {
                    if (message.error) {
                        console.error('Ошибка сообщения:', message.error);
                        return;
                    }
                    
                    setMessages(prev => {
                        // Проверяем, нет ли уже такого сообщения
                        if (prev.some(m => m.id === message.id)) {
                            return prev;
                        }
                        // Добавляем новое сообщение и сортируем по дате
                        return [...prev, message].sort((a, b) =>
                            new Date(a.created_at) - new Date(b.created_at)
                        );
                    });
                });

                return () => messageHandler(); // Отписываемся при размонтировании
            } catch (error) {
                console.error('Error initializing chat:', error);
                setError('Ошибка при инициализации чата');
            }
        };

        initChat();

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

            if (listingId && chatsData.length > 0) {
                const chat = chatsData.find(c => c.listing_id === parseInt(listingId));
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
                msg => !msg.is_read && msg.receiver_id === user?.id  
            );
            if (unreadMessages.length > 0 && user?.id) {  
                const messageIds = unreadMessages.map(msg => msg.id);
                markMessageAsRead(messageIds);
            }
        }
    }, [selectedChat, messages, user?.id]);  
    useEffect(() => {
        if (selectedChat) {
            fetchMessages(selectedChat.id);
        }
    }, [selectedChat, fetchMessages]);

    const handleSendMessage = async (content) => {
        if (!chatServiceRef.current) {
            setError('Соединение не установлено');
            return;
        }

        if (!selectedChat || !user?.id) {  // Добавляем проверку user?.id
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

            const newMessage = {
                ...message,
                sender_id: user.id,
                created_at: new Date().toISOString(),
                is_read: false,
                sender: {  // Добавляем информацию об отправителе
                    id: user.id,
                    name: user.name,
                    picture_url: user.picture_url
                }
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

    if (error) {
        return (
            <Alert severity="error" sx={{ m: 2 }} onClose={() => setError(null)}>
                {error}
            </Alert>
        );
    }

    if (isMobile) {
        return (
            <Box sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
                {selectedChat ? (
                    <>
                        <ChatHeader chat={selectedChat} onBack={() => setSelectedChat(null)} />
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
