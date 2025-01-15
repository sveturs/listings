// frontend/hostel-frontend/src/components/marketplace/chat/ChatComponents.js

import React, { useState, useRef, useEffect } from 'react';
import {
    Box,
    Paper,
    Typography,
    Avatar,
    TextField,
    IconButton,
    List,
    ListItem,
    ListItemAvatar,
    ListItemText,
    Badge,
    Stack,
     
    Chip,
    Button,
    useTheme,
    useMediaQuery,
} from '@mui/material';
import {Phone, ArrowLeft } from '@mui/icons-material';
import {
    Send as SendIcon,
    Archive as ArchiveIcon,
} from '@mui/icons-material';
import { formatDistanceToNow } from 'date-fns';
import { ru } from 'date-fns/locale';
export const ChatWindow = ({ messages = [], onSendMessage, currentUser }) => {

    const [newMessage, setNewMessage] = useState('');
    const messagesEndRef = useRef(null);
    const [processedMessages, setProcessedMessages] = useState([]);

    useEffect(() => {
        const uniqueMessages = Object.values(
            messages.reduce((acc, message) => {
                acc[message.id] = message;
                return acc;
            }, {})
        ).sort((a, b) => new Date(a.created_at) - new Date(b.created_at));

        setProcessedMessages(uniqueMessages);
    }, [messages]);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    useEffect(() => {
        scrollToBottom();
    }, [processedMessages]);

    const handleSend = (e) => {
        e.preventDefault();
        if (newMessage.trim()) {
            onSendMessage(newMessage.trim());
            setNewMessage('');
        }
    };

    const formatTime = (dateString) => {
        const date = new Date(dateString);
        return date.toLocaleTimeString('ru-RU', {
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    const formatDate = (dateString) => {
        const date = new Date(dateString);
        const today = new Date();
        const yesterday = new Date(today);
        yesterday.setDate(yesterday.getDate() - 1);

        if (date.toDateString() === today.toDateString()) {
            return 'Сегодня';
        } else if (date.toDateString() === yesterday.toDateString()) {
            return 'Вчера';
        }
        return date.toLocaleDateString('ru-RU', {
            day: 'numeric',
            month: 'long'
        });
    };

    // Группируем сообщения по датам
    const messagesByDate = processedMessages.reduce((acc, message) => {
        const date = formatDate(message.created_at);
        if (!acc[date]) {
            acc[date] = [];
        }
        acc[date].push(message);
        return acc;
    }, {});

    return (
        <Paper
            elevation={0}
            sx={{
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                bgcolor: 'grey.50',
                borderRadius: 2,
                overflow: 'hidden'
            }}
        >
            <Box
                sx={{
                    flex: 1,
                    overflowY: 'auto',
                    p: 2,
                    '&::-webkit-scrollbar': {
                        width: '8px',
                    },
                    '&::-webkit-scrollbar-track': {
                        background: 'transparent'
                    },
                    '&::-webkit-scrollbar-thumb': {
                        background: 'rgba(0,0,0,0.1)',
                        borderRadius: '4px',
                    },
                }}
            >
                {Object.entries(messagesByDate).map(([date, dateMessages]) => (
                    <Box key={date}>
                        <Box
                            sx={{
                                display: 'flex',
                                justifyContent: 'center',
                                my: 2,
                            }}
                        >
                            <Typography
                                variant="caption"
                                sx={{
                                    px: 2,
                                    py: 0.5,
                                    bgcolor: 'grey.200',
                                    borderRadius: 5,
                                    color: 'text.secondary'
                                }}
                            >
                                {date}
                            </Typography>
                        </Box>
                        {dateMessages.map((message) => (
                            <Box
                                key={message.id}
                                sx={{
                                    mb: 1,
                                    display: 'flex',
                                    justifyContent: message.sender_id === currentUser.id ? 'flex-end' : 'flex-start',
                                }}
                            >
                                <Box
                                    sx={{
                                        maxWidth: '70%',
                                        bgcolor: message.sender_id === currentUser.id ? 'primary.main' : 'background.paper',
                                        color: message.sender_id === currentUser.id ? 'white' : 'text.primary',
                                        borderRadius: 2,
                                        boxShadow: 1,
                                        p: 1.5,
                                    }}
                                >
                                    <Typography variant="body1" sx={{ mb: 0.5 }}>
                                        {message.content}
                                    </Typography>
                                    <Typography
                                        variant="caption"
                                        sx={{
                                            display: 'block',
                                            textAlign: 'right',
                                            opacity: 0.8
                                        }}
                                    >
                                        {formatTime(message.created_at)}
                                    </Typography>
                                </Box>
                            </Box>
                        ))}
                    </Box>
                ))}
                <div ref={messagesEndRef} />
            </Box>

            <Box
                component="form"
                onSubmit={handleSend}
                sx={{
                    p: 2,
                    bgcolor: 'background.paper',
                    borderTop: 1,
                    borderColor: 'divider',
                }}
            >
                <Stack direction="row" spacing={1}>
                    <TextField
                        fullWidth
                        size="small"
                        placeholder="Введите сообщение..."
                        value={newMessage}
                        onChange={(e) => setNewMessage(e.target.value)}
                        multiline
                        maxRows={4}
                        sx={{
                            '& .MuiOutlinedInput-root': {
                                borderRadius: 2,
                                bgcolor: 'grey.50'
                            }
                        }}
                    />
                    <IconButton
                        color="primary"
                        type="submit"
                        disabled={!newMessage.trim()}
                        sx={{
                            bgcolor: newMessage.trim() ? 'primary.main' : 'grey.200',
                            color: newMessage.trim() ? 'white' : 'grey.400',
                            '&:hover': {
                                bgcolor: newMessage.trim() ? 'primary.dark' : 'grey.300',
                            }
                        }}
                    >
                        <SendIcon />
                    </IconButton>
                </Stack>
            </Box>
        </Paper>
    );
};

// Компонент списка чатов
export const ChatList = ({ chats, selectedChatId, onSelectChat, onArchiveChat }) => {
    const formatPrice = (price) => {
        return new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'RUB',
            maximumFractionDigits: 0
        }).format(price);
    };

    return (
        <Paper sx={{
            height: '100%',
            border: 1,
            borderColor: 'divider',
            borderRadius: 2,
            overflow: 'hidden'
        }}>
            <List sx={{ p: 0 }}>
                {chats.map((chat) => (
                    <ListItem
                        key={chat.id}
                        button
                        selected={selectedChatId === chat.id}
                        onClick={() => onSelectChat(chat)}
                        sx={{
                            borderBottom: 1,
                            borderColor: 'divider',
                            '&:last-child': { borderBottom: 0 },
                            '&.Mui-selected': {
                                bgcolor: 'primary.light',
                                '&:hover': {
                                    bgcolor: 'primary.light',
                                }
                            }
                        }}
                    >
                        <Box sx={{ width: '100%' }}>
                            <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                <Avatar
                                    src={chat.listing?.images?.[0]?.file_path}
                                    variant="rounded"
                                    sx={{ width: 48, height: 48, mr: 1.5 }}
                                />
                                <Box sx={{ flex: 1 }}>
                                    <Typography variant="subtitle2" noWrap>
                                        {chat.listing?.title}
                                    </Typography>
                                    <Typography
                                        variant="body2"
                                        color="primary"
                                        sx={{ fontWeight: 'medium' }}
                                    >
                                        {formatPrice(chat.listing?.price)}
                                    </Typography>
                                </Box>
                                {chat.unread_count > 0 && (
                                    <Chip
                                        size="small"
                                        label={chat.unread_count}
                                        color="primary"
                                        sx={{ ml: 1 }}
                                    />
                                )}
                            </Box>

                            <Box sx={{
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'space-between'
                            }}>
                                <Typography
                                    variant="caption"
                                    sx={{
                                        color: 'text.secondary',
                                        display: 'flex',
                                        alignItems: 'center'
                                    }}
                                >
                                    {chat.other_user?.name || 'Пользователь'}
                                </Typography>
                                <Typography variant="caption" color="text.secondary">
                                    {formatDistanceToNow(new Date(chat.last_message_at), {
                                        addSuffix: true,
                                        locale: ru
                                    })}
                                </Typography>
                            </Box>

                            {chat.last_message && (
                                <Typography
                                    variant="body2"
                                    color="text.secondary"
                                    sx={{
                                        mt: 0.5,
                                        overflow: 'hidden',
                                        textOverflow: 'ellipsis',
                                        display: '-webkit-box',
                                        WebkitLineClamp: 1,
                                        WebkitBoxOrient: 'vertical',
                                    }}
                                >
                                    {chat.last_message.content}
                                </Typography>
                            )}
                        </Box>
                    </ListItem>
                ))}
            </List>
        </Paper>
    );
};

// Компонент заголовка чата
export const ChatHeader = ({ chat, onBack, onArchive }) => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    const formatPrice = (price) => {
        return new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'RUB',
            maximumFractionDigits: 0
        }).format(price || 0);
    };

    const formatLastSeen = (date) => {
        if (!date) return 'не в сети';
        const lastSeen = new Date(date);
        const now = new Date();
        const diffMinutes = Math.floor((now - lastSeen) / (1000 * 60));

        if (diffMinutes < 1) return 'в сети';
        if (diffMinutes < 60) return `был(а) ${diffMinutes} мин. назад`;
        if (diffMinutes < 1440) {
            const hours = Math.floor(diffMinutes / 60);
            return `был(а) ${hours} ч. назад`;
        }
        return 'был(а) давно';
    };

    return (
        <Paper
            elevation={1}
            sx={{
                bgcolor: 'background.paper',
                borderBottom: 1,
                borderColor: 'divider',
            }}
        >
            {/* Основная информация */}
            <Box sx={{ p: 2 }}>
                <Stack direction="row" spacing={2} alignItems="center">
                    {/* Кнопка "Назад" для мобильной версии */}
                    {isMobile && (
                        <IconButton onClick={onBack} edge="start" sx={{ mr: 1 }}>
                            <ArrowLeft />
                        </IconButton>
                    )}

                    {/* Изображение товара */}
                    <Avatar
                        variant="rounded"
                        src={chat.listing?.images?.[0]?.file_path}
                        sx={{
                            width: 48,
                            height: 48,
                            borderRadius: 1,
                            border: 1,
                            borderColor: 'divider'
                        }}
                    />

                    {/* Информация о товаре и продавце */}
                    <Box sx={{ flex: 1, minWidth: 0 }}>
                        <Typography variant="subtitle1" noWrap>
                            {chat.listing?.title}
                        </Typography>
                        <Stack direction="row" spacing={2} alignItems="center">
                            <Typography
                                variant="subtitle2"
                                color="primary.main"
                                sx={{ fontWeight: 500 }}
                            >
                                {formatPrice(chat.listing?.price)}
                            </Typography>
                            <Box
                                component="span"
                                sx={{
                                    width: 4,
                                    height: 4,
                                    borderRadius: '50%',
                                    bgcolor: 'grey.400'
                                }}
                            />
                            <Typography variant="body2" color="text.secondary" noWrap>
                                {chat.other_user?.name}
                            </Typography>
                        </Stack>
                    </Box>

                    {/* Действия */}
                    <Stack direction="row" spacing={1}>
                        {/* Кнопка архивации */}
                        <IconButton
                            onClick={() => onArchive?.(chat.id)}
                            sx={{
                                color: 'grey.600',
                                '&:hover': {
                                    color: 'warning.main',
                                    bgcolor: 'warning.lighter'
                                }
                            }}
                        >
                            <ArchiveIcon fontSize="small" />
                        </IconButton>

                        {/* Кнопка звонка */}
                        <IconButton
                            href={`tel:${chat.other_user?.phone}`}
                            sx={{
                                color: 'grey.600',
                                '&:hover': {
                                    color: 'success.main',
                                    bgcolor: 'success.lighter'
                                }
                            }}
                        >
                            <Phone fontSize="small" />
                        </IconButton>

                        {/* Переход к объявлению */}
                        <Button
                            variant="outlined"
                            size="small"
                            onClick={() => window.open(`/marketplace/listings/${chat.listing?.id}`, '_blank')}
                            sx={{
                                minWidth: 'auto',
                                px: 2,
                                borderRadius: 1,
                                display: { xs: 'none', sm: 'inline-flex' }
                            }}
                        >
                            Открыть объявление
                        </Button>
                    </Stack>
                </Stack>
            </Box>
        </Paper>
    );
};

// Компонент пустого состояния
export const EmptyState = ({ text }) => (
    <Box
        sx={{
            height: '100%',
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            p: 3,
            textAlign: 'center',
        }}
    >
        <Typography variant="h6" color="text.secondary" gutterBottom>
            {text}
        </Typography>
    </Box>
);