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
import { ArrowLeft } from '@mui/icons-material';
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

    // Дедупликация сообщений по ID
    useEffect(() => {
        const uniqueMessages = Object.values(
            messages.reduce((acc, message) => {
                // Используем ID сообщения как ключ для дедупликации
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

    return (
        <Paper sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            <Box sx={{ flex: 1, overflowY: 'auto', p: 2 }}>
                {processedMessages.map((message) => (
                    <Box
                        key={message.id}
                        sx={{
                            mb: 2,
                            display: 'flex',
                            flexDirection: message.sender_id === currentUser.id ? 'row-reverse' : 'row',
                            alignItems: 'flex-start',
                        }}
                    >
                        <Box sx={{ mx: 1 }}>
                            <Avatar
                                src={message.sender?.picture_url}
                                sx={{ width: 32, height: 32 }}
                            />
                            <Typography
                                variant="caption"
                                sx={{
                                    display: 'block',
                                    textAlign: message.sender_id === currentUser.id ? 'right' : 'left',
                                    mt: 0.5,
                                }}
                            >
                                {message.sender?.name || 'Пользователь'}
                            </Typography>
                        </Box>

                        <Box
                            sx={{
                                maxWidth: '70%',
                                bgcolor: message.sender_id === currentUser.id ? 'primary.main' : 'grey.100',
                                color: message.sender_id === currentUser.id ? 'white' : 'text.primary',
                                borderRadius: 2,
                                p: 1.5,
                                position: 'relative',
                            }}
                        >
                            <Typography variant="body1">{message.content}</Typography>
                            <Typography
                                variant="caption"
                                sx={{
                                    display: 'block',
                                    mt: 0.5,
                                    opacity: 0.8,
                                }}
                            >
                                {formatTime(message.created_at)}
                            </Typography>
                        </Box>
                    </Box>
                ))}
                <div ref={messagesEndRef} />
            </Box>

            <Box
                component="form"
                onSubmit={handleSend}
                sx={{
                    p: 2,
                    borderTop: 1,
                    borderColor: 'divider',
                    bgcolor: 'background.default',
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
                    />
                    <IconButton
                        color="primary"
                        type="submit"
                        disabled={!newMessage.trim()}
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
    return (
        <Paper sx={{
            height: '100%',
            border: 1,
            borderColor: 'divider',
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
                            '&:last-child': {
                                borderBottom: 0,
                            },
                        }}
                        secondaryAction={
                            <IconButton
                                edge="end"
                                onClick={(e) => {
                                    e.stopPropagation();
                                    onArchiveChat(chat.id);
                                }}
                            >
                                <ArchiveIcon />
                            </IconButton>
                        }
                    >
                        <ListItemAvatar>
                            <Badge
                                badgeContent={chat.unread_count}
                                color="primary"
                                invisible={chat.unread_count === 0}
                            >
                                <Avatar src={chat.listing?.images?.[0]?.file_path} />
                            </Badge>
                        </ListItemAvatar>
                        <ListItemText
                            primary={
                                <Typography
                                    variant="subtitle2"
                                    component="div"
                                    noWrap
                                    sx={{ mb: 0.5 }}
                                >
                                    {chat.listing?.title}
                                </Typography>
                            }
                            secondary={
                                <Box component="div" sx={{ mt: 0.5 }}>
                                    <Typography
                                        variant="body2"
                                        component="div"
                                        color="text.secondary"
                                        noWrap
                                        sx={{ mb: 0.5 }}
                                    >
                                        {chat.last_message?.content}
                                    </Typography>
                                    <Typography
                                        variant="caption"
                                        component="div"
                                        color="text.secondary"
                                    >
                                        {formatDistanceToNow(new Date(chat.last_message_at), {
                                            addSuffix: true,
                                            locale: ru
                                        })}
                                    </Typography>
                                </Box>
                            }
                        />

                    </ListItem>
                ))}
            </List>
        </Paper>
    );
};

// Компонент заголовка чата
export const ChatHeader = ({ chat, onBack }) => {
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

    return (
        <Box
            sx={{
                p: 2,
                borderBottom: 1,
                borderColor: 'divider',
                bgcolor: 'background.paper',
                display: 'flex',
                alignItems: 'center',
                gap: 2,
            }}
        >
            {isMobile && (
                <IconButton onClick={onBack}>
                    <ArrowLeft />
                </IconButton>
            )}
            <Avatar
                src={chat.listing?.images?.[0]?.file_path}
                sx={{ width: 40, height: 40 }}
            />
            <Box sx={{ flex: 1 }}>
                <Typography variant="subtitle1" noWrap>
                    {chat.listing?.title}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    {chat.listing?.price}₽
                </Typography>
            </Box>
        </Box>
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