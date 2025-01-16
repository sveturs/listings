// frontend/hostel-frontend/src/components/marketplace/chat/ChatComponents.js
// frontend/hostel-frontend/src/components/marketplace/chat/ChatComponents.js

import React, { useState, useRef, useEffect } from 'react';
import EmojiPicker from 'emoji-picker-react';
import { MessageCircle, Send, Smile, Paperclip, ChevronLeft, Phone } from 'lucide-react';
import {
    Box,
    Paper,
    Typography,
    Avatar,
    TextField,
    IconButton,
    List,
    ListItem,
    Badge,
    Stack,
    Popover,
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

const MessageContent = ({ content }) => {
    const isOnlyEmoji = (text) => {
        const emojiRegex = /^(?:\p{Extended_Pictographic}|\p{Emoji_Presentation}|\p{Emoji}\uFE0F|\p{Emoji_Modifier})+$/u;
        return emojiRegex.test(text.trim());
    };

    const parseMessage = (text) => {
        const emojiRegex = /((?:\p{Extended_Pictographic}|\p{Emoji_Presentation}|\p{Emoji}\uFE0F|\p{Emoji_Modifier})+)/u;
        return text.split(emojiRegex).filter(Boolean);
    };

    const onlyEmoji = isOnlyEmoji(content);
    const parts = parseMessage(content);

    return (
        <Typography
            variant="body2"
            component="div"
            sx={{
                '& .emoji': {
                    // Используем стили для Apple эмодзи
                    fontFamily: 'Apple Color Emoji',  // Основной шрифт для Apple эмодзи
                    fontSize: onlyEmoji ? '4rem' : '1.5rem',
                    lineHeight: 1,
                    verticalAlign: 'middle',
                    fontStyle: 'normal', // Важно для корректного отображения
                    WebkitFontSmoothing: 'antialiased', // Улучшает отображение на webkit браузерах
                    textRendering: 'optimizeLegibility', // Улучшает четкость отображения
                    fontFamily: 'Apple Color Emoji, -apple-system-emoji, "Segoe UI Emoji", "Noto Color Emoji", sans-serif',
                }
            }}
        >
            {parts.map((part, index) => {
                const isEmoji = /(?:\p{Extended_Pictographic}|\p{Emoji_Presentation}|\p{Emoji}\uFE0F|\p{Emoji_Modifier})+/u.test(part);
                return (
                    <span
                        key={index}
                        className={isEmoji ? 'emoji' : undefined}
                    >
                        {part}
                    </span>
                );
            })}
        </Typography>
    );
};

export const ChatWindow = ({ messages = [], onSendMessage, currentUser, chat, onBack }) => {
    const [newMessage, setNewMessage] = useState('');
    const messagesEndRef = useRef(null);
    const [processedMessages, setProcessedMessages] = useState([]);
    const [isTyping, setIsTyping] = useState(false);
    const [anchorEl, setAnchorEl] = useState(null);

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

    const handleEmojiClick = (emojiData) => {
        setNewMessage((prevMessage) => prevMessage + emojiData.emoji);
        setAnchorEl(null);
    };

    const handleEmojiButtonClick = (event) => {
        setAnchorEl(anchorEl ? null : event.currentTarget);
    };

    return (
        <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%', bgcolor: 'grey.50' }}>
            {/* Шапка чата */}
            <Box sx={{
                display: 'flex',
                alignItems: 'center',
                px: 2,
                py: 1.5,
                bgcolor: 'white',
                borderBottom: 1,
                borderColor: 'divider'
            }}>
                <IconButton onClick={onBack} sx={{ mr: 1 }}>
                    <ChevronLeft />
                </IconButton>
                <Box sx={{ flex: 1 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        <Avatar
                            variant="rounded"
                            src={chat?.listing?.images?.[0]?.file_path}
                            sx={{
                                width: 40,
                                height: 40,
                                mr: 1.5,
                                bgcolor: 'primary.light'
                            }}
                        />
                        <Box>
                            <Typography variant="subtitle1">
                                {chat?.listing?.title || 'Чат с продавцом'}
                            </Typography>
                            {isTyping && (
                                <Typography variant="caption" color="success.main">
                                    печатает...
                                </Typography>
                            )}
                        </Box>
                    </Box>
                </Box>
                <IconButton>
                    <Phone />
                </IconButton>
            </Box>

            {/* Область сообщений */}
            <Box sx={{
                flex: 1,
                overflowY: 'auto',
                p: 2,
                '&::-webkit-scrollbar': {
                    width: 8,
                    borderRadius: 4,
                },
                '&::-webkit-scrollbar-track': {
                    backgroundColor: 'transparent'
                },
                '&::-webkit-scrollbar-thumb': {
                    backgroundColor: 'rgba(0,0,0,0.1)',
                    borderRadius: 4
                }
            }}>
                {processedMessages.map((message) => (
                    <Box
                        key={message.id}
                        sx={{
                            display: 'flex',
                            justifyContent: message.sender_id === currentUser.id ? 'flex-end' : 'flex-start',
                            mb: 1.5
                        }}
                    >
                        <Box sx={{
                            maxWidth: '70%',
                            px: 2,
                            py: 1.5,
                            bgcolor: message.sender_id === currentUser.id ? 'primary.main' : 'white',
                            color: message.sender_id === currentUser.id ? 'white' : 'text.primary',
                            borderRadius: 3,
                            ...(message.sender_id === currentUser.id ? {
                                borderBottomRightRadius: 0,
                            } : {
                                borderBottomLeftRadius: 0,
                                boxShadow: 1
                            })
                        }}>
                            <MessageContent content={message.content} />
                            <Typography
                                variant="caption"
                                sx={{
                                    display: 'block',
                                    mt: 0.5,
                                    textAlign: 'right',
                                    opacity: 0.7
                                }}
                            >
                                {new Date(message.created_at).toLocaleTimeString('ru-RU', {
                                    hour: '2-digit',
                                    minute: '2-digit'
                                })}
                            </Typography>
                        </Box>
                    </Box>
                ))}
                <div ref={messagesEndRef} />
            </Box>

            {/* Форма отправки */}
            <Box
                component="form"
                onSubmit={handleSend}
                sx={{
                    p: 2,
                    bgcolor: 'white',
                    borderTop: 1,
                    borderColor: 'divider'
                }}
            >
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <IconButton>
                        <Paperclip />
                    </IconButton>
                    <Box sx={{
                        flex: 1,
                        position: 'relative'
                    }}>
                        <TextField
                            fullWidth
                            size="small"
                            placeholder="Введите сообщение..."
                            value={newMessage}
                            onChange={(e) => setNewMessage(e.target.value)}
                            sx={{
                                '& .MuiOutlinedInput-root': {
                                    borderRadius: '24px',
                                    bgcolor: 'grey.50',
                                    pr: 5
                                }
                            }}
                        />
                        <IconButton
                            onClick={handleEmojiButtonClick}
                            sx={{
                                position: 'absolute',
                                right: 8,
                                top: '50%',
                                transform: 'translateY(-50%)'
                            }}
                        >
                            <Smile />
                        </IconButton>
                        <Popover
                            open={Boolean(anchorEl)}
                            anchorEl={anchorEl}
                            onClose={() => setAnchorEl(null)}
                            anchorOrigin={{
                                vertical: 'top',
                                horizontal: 'right',
                            }}
                            transformOrigin={{
                                vertical: 'bottom',
                                horizontal: 'right',
                            }}
                        >
                            <EmojiPicker
                                onEmojiClick={handleEmojiClick}
                                width={320}
                                height={400}
                                searchDisabled={true}
                                skinTonesDisabled={true}
                                emojiStyle="native"
                                style={{
                                    '--epr-bg-color': 'white',
                                    '--epr-category-label-bg-color': 'white'
                                }}
                            />
                        </Popover>
                    </Box>
                    <IconButton
                        type="submit"
                        disabled={!newMessage.trim()}
                        sx={{
                            bgcolor: newMessage.trim() ? 'primary.main' : 'grey.200',
                            color: newMessage.trim() ? 'white' : 'grey.400',
                            '&:hover': {
                                bgcolor: newMessage.trim() ? 'primary.dark' : 'grey.300'
                            }
                        }}
                    >
                        <Send />
                    </IconButton>
                </Box>
            </Box>
        </Box>
    );
};


// Компонент списка чатов
export const ChatList = ({ chats, selectedChatId, onSelectChat, onArchiveChat }) => {
    const formatPrice = (price) => {
        if (!price || price === undefined) return 'Цена не указана';

        return new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'RUB',
            maximumFractionDigits: 0
        }).format(price || 0);
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
                                        {chat.listing && formatPrice(chat.listing.price)}
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
        if (!price || price === undefined) return 'Цена не указана';

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
                                {chat.listing && formatPrice(chat.listing.price)}
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
                            onClick={() => {
                                if (chat.listing_id) {
                                    window.open(`/marketplace/listings/${chat.listing_id}`, '_blank');
                                }
                            }}
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