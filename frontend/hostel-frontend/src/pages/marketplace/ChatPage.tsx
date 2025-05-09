// frontend/hostel-frontend/src/pages/marketplace/ChatPage.tsx
import React, { useState, useEffect, useCallback, useRef } from 'react';
import { useTranslation } from 'react-i18next';

import { useParams, useNavigate } from 'react-router-dom';
import {
    Container,
    Grid,
    Box,
    useTheme,
    useMediaQuery,
    CircularProgress,
    Alert,
    Paper,
    Typography,
    Button
} from '@mui/material';
import {
    ChatWindow,
    ChatList,
    ChatHeader,
    EmptyState,
} from '../../components/marketplace/chat/ChatComponents';
import ChatService, { ChatMessage, ChatItem } from '../../components/marketplace/chat/ChatService';
import axios from '../../api/axios';
import { useAuth, User } from '../../contexts/AuthContext';
import { useChat } from '../../contexts/ChatContext';

const ChatPage: React.FC = () => {
    const { t } = useTranslation('marketplace');

    const { listingId } = useParams<{ listingId?: string }>(); // Fix params type
    const navigate = useNavigate();
    const { user, loading: authLoading, login } = useAuth();
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));

    const [chats, setChats] = useState<ChatItem[]>([]);
    const [selectedChat, setSelectedChat] = useState<ChatItem | null>(null);
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    const chatServiceRef = useRef<ChatService | null>(null);
    const messageEndRef = useRef<HTMLDivElement | null>(null);
    const { getChatService } = useChat();

    // Initialize chat service
    useEffect(() => {
        if (user?.id) {
            chatServiceRef.current = getChatService(Number(user.id));
        }
    }, [user?.id, getChatService]);

    // Load chat list
    const fetchChats = useCallback(async () => {
        try {
            const response = await axios.get('/api/v1/marketplace/chat');
            const chatsData: ChatItem[] = response.data?.data || [];
            setChats(chatsData);

            // If listingId is in URL, select the corresponding chat
            if (listingId) {
                const chat = chatsData.find(c => c.listing_id === parseInt(listingId));
                if (chat) {
                    setSelectedChat(chat);
                }
            }
        } catch (error) {
            console.error('Error loading chats:', error);
            setError('Не удалось загрузить список чатов');
        } finally {
            setLoading(false);
        }
    }, [listingId]);

    // Load messages for selected chat
    const fetchMessages = useCallback(async (chatId: number) => {
        if (!chatServiceRef.current) {
            throw new Error('ChatService not initialized');
        }

        try {
            const messages = await chatServiceRef.current.getMessageHistory(chatId, selectedChat?.listing_id || 0);
            return messages;
        } catch (error) {
            console.error('Error loading messages:', error);
            throw error;
        }
    }, [selectedChat?.listing_id]);

    // Handle chat selection
    // Функция для выбора чата и загрузки сообщений
    const handleSelectChat = useCallback(async (chat: ChatItem) => {
        if (!chat?.id) {
            console.error('Invalid chat data - missing chat ID:', chat);
            setError('Ошибка чата: отсутствует ID чата');
            return;
        }

        if (!chat?.listing_id) {
            console.error('Invalid chat data - missing listing ID:', chat);
            setError('Ошибка чата: отсутствует ID объявления');
            return;
        }

        console.log('Selected chat with data:', chat);

        setSelectedChat(chat);
        setLoading(true);
        setMessages([]);

        try {
            const chatService = chatServiceRef.current;
            if (!chatService) {
                throw new Error('ChatService not initialized');
            }

            const chatId = Number(chat.id);
            const listingId = Number(chat.listing_id);

            if (isNaN(chatId) || chatId <= 0) {
                throw new Error(`Invalid chat ID: ${chat.id}`);
            }

            if (isNaN(listingId) || listingId <= 0) {
                throw new Error(`Invalid listing ID: ${chat.listing_id}`);
            }

            // Отладочное сообщение
            if (process.env.NODE_ENV === 'development') {
                console.debug(`Fetching messages for chat ${chatId}`);
            }
            const loadedMessages = await chatService.getMessageHistory(chatId, listingId, 100); // Увеличиваем лимит сообщений

            // Отладочное сообщение
            if (process.env.NODE_ENV === 'development') {
                console.debug(`Received ${loadedMessages.length} messages from history`);
            }
            if (Array.isArray(loadedMessages) && loadedMessages.length > 0) {
                setMessages(loadedMessages);

                // Mark unread messages as read
                const unreadMessages = loadedMessages.filter(
                    msg => !msg.is_read && msg.receiver_id === user?.id
                );

                if (unreadMessages.length > 0) {
                    // Отключаем логирование - это сообщение вызывается слишком часто
                    const messageIds = unreadMessages.map(msg => Number(msg.id));
                    const marked = await chatService.markMessagesAsRead(messageIds);

                    if (marked) {
                        // После успешной отметки сообщений, обновляем списки чатов несколько раз с задержкой
                        // Отключаем логирование - это сообщение вызывается слишком часто
                        try {
                            // Обновляем немедленно
                            await chatService.updateChatsList();

                            // И с несколькими задержками для гарантии получения актуальных данных с сервера
                            setTimeout(() => {
                                chatService.updateChatsList().catch(err => {
                                    // Подавляем большинство ошибок логирования
                                    if (process.env.NODE_ENV === 'development' && Math.random() < 0.1) {
                                        console.debug('Error updating chats list (delayed)');
                                    }
                                });
                            }, 500);

                            setTimeout(() => {
                                chatService.updateChatsList().catch(err => {
                                    // Подавляем большинство ошибок логирования
                                    if (process.env.NODE_ENV === 'development' && Math.random() < 0.1) {
                                        console.debug('Error updating chats list (delayed)');
                                    }
                                });
                            }, 2000);
                        } catch (error) {
                            // Логируем только критические ошибки
                            console.error('Failed to update chats list');
                        }
                    }
                }

                // Update unread count in chat list
                setChats(prevChats =>
                    prevChats.map(c => {
                        if (c.id === chat.id) {
                            return { ...c, unread_count: 0 };
                        }
                        return c;
                    })
                );
            } else {
                console.log('No messages in history or invalid response');
            }
        } catch (error) {
            console.error('Error loading messages:', error);
            setError('Не удалось загрузить сообщения: ' + (error instanceof Error ? error.message : String(error)));
        } finally {
            setLoading(false);
        }
    }, [user?.id]);

    // Initialize WebSocket and load data
    useEffect(() => {
        const chatService = chatServiceRef.current;
        if (chatService && user?.id) {
            // Subscribe to chat list updates
            const unsubscribeChatList = chatService.onChatListUpdate((updatedChats) => {
                // Обновляем список чатов и сохраняем выбранный чат
                setChats(updatedChats);

                // Если есть выбранный чат, обновляем его данные из полученного списка
                if (selectedChat) {
                    const updatedSelectedChat = updatedChats.find(c => c.id === selectedChat.id);
                    if (updatedSelectedChat) {
                        setSelectedChat(updatedSelectedChat);
                    }
                }
            });

            // Handle new messages
            const unsubscribeMessages = chatService.onMessage((message) => {
                // Отладочное сообщение
            if (process.env.NODE_ENV === 'development') {
                console.debug('New message received');
            }

                // Проверяем, является ли сообщение сообщением об ошибке
                if (message.error) {
                    console.warn('Error message received:', message.error);
                    return; // Не обрабатываем сообщения с ошибками
                }

                // Update messages in current chat
                if (selectedChat && message.chat_id === selectedChat.id) {
                    setMessages(prev => {
                        // Улучшенная проверка дубликатов с расширенным временным окном
                        const isDuplicate = prev.some(m => {
                            // Проверка по ID (если есть)
                            if (m.id === message.id && message.id) {
                                if (process.env.NODE_ENV === 'development') {
                            console.debug('Duplicate message detected by ID');
                        }
                                return true;
                            }

                            // Проверка по клиентскому ID (если есть)
                            if (m.client_message_id && message.client_message_id &&
                                m.client_message_id === message.client_message_id) {
                                if (process.env.NODE_ENV === 'development') {
                                    console.debug('Duplicate message detected by client_message_id');
                                }
                                return true;
                            }

                            // Проверка по содержимому, отправителю и расширенному временному окну (10 секунд вместо 5)
                            if (m.content === (message.content || '') &&
                                m.sender_id === message.sender_id &&
                                Math.abs(new Date(m.created_at || '').getTime() -
                                    new Date(message.created_at || '').getTime()) < 10000) {
                                // Отладочное сообщение
                                if (process.env.NODE_ENV === 'development') {
                                    console.debug('Duplicate message detected');
                                }
                                return true;
                            }

                            return false;
                        });

                        if (isDuplicate) {
                            // Отладочное сообщение
                        if (process.env.NODE_ENV === 'development') {
                            console.debug('Duplicate message, not adding to UI');
                        }
                            return prev;
                        }

                        // Создаем сообщение с гарантированными полями
                        const newMessage = {
                            ...message,
                            content: message.content || '',  // Используем только поле content
                            sender: message.sender || {},
                            receiver: message.receiver || {},
                            is_read: message.is_read || false,
                            created_at: message.created_at || new Date().toISOString()
                        };

                        // Отладочное сообщение
                        if (process.env.NODE_ENV === 'development') {
                            console.debug('Adding new message to UI');
                        }

                        // Добавляем сообщение и сортируем по времени
                        const updatedMessages = [...prev, newMessage];
                        return updatedMessages.sort((a, b) =>
                            new Date(a.created_at || '').getTime() - new Date(b.created_at || '').getTime()
                        );
                    });

                    // Mark messages as read
                    if (!message.is_read && message.receiver_id === user.id && message.id) {
                        chatService.markMessagesAsRead([Number(message.id)]);
                    }
                }
    
                // Update chat list to display new messages
                // Обновляем только если сообщение имеет chat_id
                if (message.chat_id) {
                    setChats(prevChats => {
                        const updatedChats = prevChats.map(chat => {
                            if (chat.id === message.chat_id) {
                                return {
                                    ...chat,
                                    last_message: {
                                        ...message,
                                        content: message.content || ''  // Используем только поле content
                                    },
                                    unread_count: chat.id !== selectedChat?.id && message.receiver_id === user.id
                                        ? (chat.unread_count || 0) + 1
                                        : chat.unread_count
                                };
                            }
                            return chat;
                        });
                        return updatedChats;
                    });
                }
            });
    
            // Initialize WebSocket connection
            chatService.connect();
    
            return () => {
                unsubscribeChatList();
                unsubscribeMessages();
            };
        }
    }, [user?.id, selectedChat]);

    // Load chats on mount
    useEffect(() => {
        fetchChats();

        // Периодически обновляем список чатов для получения актуальных счетчиков непрочитанных сообщений
        const intervalId = setInterval(() => {
            if (chatServiceRef.current) {
                chatServiceRef.current.updateChatsList();
            }
        }, 10000); // Обновление каждые 10 секунд

        return () => {
            clearInterval(intervalId);
        };
    }, [fetchChats]);

    useEffect(() => {
        if (selectedChat?.id && chatServiceRef.current) {
            // При выборе чата подключаемся к WebSocket
            chatServiceRef.current.connect();

            // Также обновляем список чатов, чтобы получить актуальные счетчики
            chatServiceRef.current.updateChatsList();
        }
    }, [selectedChat?.id]);

    useEffect(() => {
        if (selectedChat && messages.length > 0) {
            const unreadMessages = messages.filter(
                msg => !msg.is_read && msg.receiver_id === user?.id
            );

            if (unreadMessages.length > 0) {
                const messageIds = unreadMessages.map(msg => Number(msg.id));
                const markAndUpdateAsync = async () => {
                    const marked = await chatServiceRef.current?.markMessagesAsRead(messageIds);
                    if (marked && chatServiceRef.current) {
                        // Обновляем список чатов для актуализации счетчика непрочитанных
                        await chatServiceRef.current.updateChatsList();
                    }
                };
                markAndUpdateAsync();
            }
        }
    }, [selectedChat, messages, user?.id]);

    // Scroll to last message
    useEffect(() => {
        messageEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);

    // Send message
    const handleSendMessage = async (content: string) => {
        if (!selectedChat) {
            setError('Выберите чат для отправки сообщения');
            return;
        }

        if (!user?.id) {
            setError('Требуется авторизация для отправки сообщения');
            return;
        }

        // Проверка корректности ID чата и объявления
        const chatId = Number(selectedChat.id);
        const listingId = Number(selectedChat.listing_id);

        if (isNaN(chatId) || chatId <= 0) {
            console.error(`Invalid chat ID: ${selectedChat.id}`);
            setError('Ошибка: некорректный ID чата');
            return;
        }

        if (isNaN(listingId) || listingId <= 0) {
            console.error(`Invalid listing ID: ${selectedChat.listing_id}`);
            setError('Ошибка: некорректный ID объявления');
            return;
        }

        if (!content.trim()) {
            setError('Текст сообщения не может быть пустым');
            return;
        }

        try {
            // Определяем receiver_id: если текущий пользователь - продавец, то получатель - покупатель, и наоборот
            let receiver_id;

            if (selectedChat.seller_id === user.id) {
                receiver_id = Number(selectedChat.buyer_id);
            } else {
                receiver_id = Number(selectedChat.seller_id);
            }

            if (isNaN(receiver_id) || receiver_id <= 0) {
                console.error(`Invalid receiver ID: ${receiver_id}`);
                setError('Ошибка: не удалось определить получателя');
                return;
            }

            const message: ChatMessage = {
                chat_id: chatId,
                listing_id: listingId,
                receiver_id: receiver_id,
                content: content.trim() // Используем поле content и удаляем лишние пробелы
            };

            // Отключаем логирование - это сообщение вызывается слишком часто
            const response = await chatServiceRef.current?.sendMessage(message);
            // Отключаем логирование - это сообщение вызывается слишком часто

            // Генерируем уникальный клиентский ID для дедупликации
            const clientMessageId = `client_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

            // Создаем сообщение для оптимистичного UI с клиентским ID
            const newMessage: ChatMessage = {
                ...message,
                id: Date.now(),  // Временный ID для оптимистичного UI
                sender_id: user.id,
                created_at: new Date().toISOString(),
                is_read: false,
                content: message.content,  // Гарантируем, что поле content существует и заполнено
                client_message_id: clientMessageId  // Добавляем клиентский ID для надежной дедупликации
            };

            // Добавляем сообщение в UI для оптимистичного опыта с улучшенной проверкой на дубликаты
            setMessages(prev => {
                // Улучшенная проверка дубликатов с использованием client_message_id и расширенным временным окном
                const isDuplicate = prev.some(msg => {
                    // Проверка по клиентскому ID
                    if (msg.client_message_id === newMessage.client_message_id && newMessage.client_message_id) {
                        console.log('Duplicate message detected by client_message_id:', newMessage.client_message_id);
                        return true;
                    }

                    // Проверка по содержимому и отправителю с расширенным временным окном (10 секунд)
                    if (msg.content === newMessage.content &&
                        msg.sender_id === newMessage.sender_id &&
                        new Date().getTime() - new Date(msg.created_at || "").getTime() < 10000) {
                        console.log('Duplicate message detected by content and time (10s window)');
                        return true;
                    }

                    return false;
                });

                if (isDuplicate) {
                    if (process.env.NODE_ENV === 'development') {
                        console.debug('Duplicate message detected, not adding to UI');
                    }
                    return prev;
                }

                if (process.env.NODE_ENV === 'development') {
                    console.debug('Adding new message to UI');
                }
                return [...prev, newMessage];
            });

            // Обновляем данные в списке чатов
            setChats(prevChats => {
                return prevChats.map(chat => {
                    if (chat.id === chatId) {
                        return {
                            ...chat,
                            last_message: newMessage,
                            last_message_at: new Date().toISOString()
                        };
                    }
                    return chat;
                });
            });

            // Запрашиваем обновление списка чатов с сервера
            if (chatServiceRef.current) {
                // Сначала сразу пытаемся обновить список чатов
                try {
                    await chatServiceRef.current.updateChatsList();
                } catch (err) {
                    console.error('Error updating chats list immediately after send:', err);
                }

                // Затем с задержкой для гарантии
                // Обновляем только список чатов без перезагрузки сообщений
                setTimeout(async () => {
                    try {
                        // Обновляем только список чатов для получения актуальных счетчиков
                        await chatServiceRef.current?.updateChatsList();
                    } catch (err) {
                        console.error('Error updating chats list after send:', err);
                    }
                }, 2000); // Задержка 2 секунды для гарантии обработки
            }
        } catch (error) {
            console.error('Error sending message:', error);
            setError('Не удалось отправить сообщение');
        }
    };

    // Archive chat
    const handleArchiveChat = async (chatId: string | number) => {
        try {
            await axios.post(`/api/v1/marketplace/chat/${chatId}/archive`);
            await fetchChats();

            // Обновляем список чатов через сервис для актуализации счетчика
            if (chatServiceRef.current) {
                await chatServiceRef.current.updateChatsList();
            }

            if (selectedChat?.id === chatId) {
                setSelectedChat(null);
                setMessages([]);
            }
        } catch (error) {
            console.error('Error archiving chat:', error);
            setError('Не удалось архивировать чат');
        }
    };

    if (authLoading) {
        return (
            <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
                <CircularProgress />
            </Box>
        );
    }

    if (!user) {
        return (
            <Container maxWidth="md" sx={{ mt: 4 }}>
                <Paper sx={{ p: 3, textAlign: 'center' }}>
                    <Typography variant="h6" gutterBottom>
                        Необходима авторизация
                    </Typography>
                    <Typography color="text.secondary" paragraph>
                        Для доступа к чату необходимо войти в систему
                    </Typography>
                    <Button 
                        variant="contained" 
                        onClick={() => {
                            const returnUrl = window.location.pathname;
                            const encodedReturnUrl = encodeURIComponent(returnUrl);
                            login(`?returnTo=${encodedReturnUrl}`);
                        }}
                    >
                        Войти
                    </Button>
                </Paper>
            </Container>
        );
    }

    // Mobile version
    if (isMobile) {
        return (
            <Box sx={{ 
                height: `calc(100vh - 112px)`, // 56px top bar + 56px top margin
                display: 'flex',
                flexDirection: 'column',
                bgcolor: 'background.default',
                position: 'fixed',
                top: 56,
                left: 0,
                right: 0,
                bottom: 0,
                zIndex: 1000
            }}>
                {selectedChat ? (
                    <Box sx={{ 
                        display: 'flex', 
                        flexDirection: 'column',
                        height: '100%',
                        overflow: 'hidden'
                    }}>
                        <ChatHeader
                            chat={selectedChat}
                            onBack={() => setSelectedChat(null)}
                            onArchive={handleArchiveChat}
                        />
                        <Box sx={{ 
                            flex: 1,
                            minHeight: 0,
                            position: 'relative',
                            overflow: 'hidden'
                        }}>
                            <ChatWindow
                                messages={messages}
                                onSendMessage={handleSendMessage}
                                currentUser={user}
                                chat={selectedChat}
                            />
                        </Box>
                    </Box>
                ) : (
                    <Box sx={{ 
                        flex: 1,
                        overflow: 'auto',
                        height: '100%'
                    }}>
                        <ChatList
                            chats={chats}
                            selectedChatId={selectedChat?.id}
                            onSelectChat={handleSelectChat}
                            onArchiveChat={handleArchiveChat}
                        />
                        {!loading && chats.length === 0 && (
                            <EmptyState text={t('chat.noMessages')} />
                        )}
                    </Box>
                )}
            </Box>
        );
    }

    // Desktop version
    return (
        <Container 
            maxWidth="xl" 
            sx={{ 
                py: 2,
                height: 'calc(100vh - 80px)', // 64px top bar + margins
                mt: 1
            }}
        >
            <Grid container spacing={2} sx={{ height: '100%' }}>
                <Grid item xs={12} md={4} sx={{ height: '100%' }}>
                    <ChatList
                        chats={chats}
                        selectedChatId={selectedChat?.id}
                        onSelectChat={handleSelectChat}
                        onArchiveChat={handleArchiveChat}
                    />
                    {!loading && chats.length === 0 && (
                        <EmptyState text={t('chat.noMessages')} />
                    )}
                </Grid>

                <Grid item xs={12} md={8} sx={{ height: '100%' }}>
                    {selectedChat ? (
                        <Box sx={{ 
                            height: '100%', 
                            display: 'flex', 
                            flexDirection: 'column',
                            overflow: 'hidden'
                        }}>
                            <ChatHeader chat={selectedChat} onArchive={handleArchiveChat} />
                            <Box sx={{ flex: 1, overflow: 'hidden' }}>
                                <ChatWindow
                                    messages={messages}
                                    onSendMessage={handleSendMessage}
                                    currentUser={user}
                                    chat={selectedChat}
                                />
                            </Box>
                        </Box>
                    ) : (
                        <EmptyState text={t('chat.empty')} />
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
                        maxWidth: 'calc(100% - 32px)',
                        zIndex: 1200
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