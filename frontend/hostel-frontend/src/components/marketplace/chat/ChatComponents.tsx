// frontend/hostel-frontend/src/components/marketplace/chat/ChatComponents.tsx
import React, { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import EmojiPicker, { EmojiClickData, EmojiStyle } from 'emoji-picker-react';
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
    Menu,
    MenuItem,
} from '@mui/material';
import {
    ArrowLeft,
    Send as SendIcon,
    Archive as ArchiveIcon,
    ContentCopy,
    Phone as PhoneIcon,
} from '@mui/icons-material';
import { formatDistanceToNow } from 'date-fns';
import { enUS, ru, sr } from 'date-fns/locale';

interface PhonePart {
    type: 'phone';
    content: string;
    phoneNumber: string;
}

interface TextPart {
    type: 'text' | 'emoji';
    content: string;
}

type MessagePart = PhonePart | TextPart;

interface MessageContentProps {
    content: string;
}

interface ChatMessage {
    id: string | number;
    chat_id?: string | number;
    listing_id?: string | number;
    sender_id?: string | number;
    receiver_id?: string | number;
    content: string;
    created_at?: string;
    is_read?: boolean;
}

interface User {
    id: string | number;
    name: string;
    email?: string;
    pictureUrl?: string;
    phone?: string;
    [key: string]: any;
}

interface ListingImage {
    file_path?: string;
    public_url?: string;
    is_main?: boolean;
    [key: string]: any;
}

interface Listing {
    id: string | number;
    title: string;
    price: number;
    images?: ListingImage[];
    [key: string]: any;
}

interface Chat {
    id: string | number;
    listing_id: string | number;
    listing?: Listing;
    user_id?: string | number;
    participant_id?: string | number;
    other_user?: User;
    unread_count?: number;
    last_message_at?: string;
    last_message?: ChatMessage;
    [key: string]: any;
}

interface ChatWindowProps {
    messages?: ChatMessage[];
    onSendMessage: (message: string) => void;
    currentUser: User;
    chat: Chat;
    onBack?: () => void;
}

interface ChatListProps {
    chats: Chat[];
    selectedChatId?: string | number;
    onSelectChat: (chat: Chat) => void;
    onArchiveChat?: (chatId: string | number) => void;
}

interface ChatHeaderProps {
    chat: Chat;
    onBack?: () => void;
    onArchive?: (chatId: string | number) => void;
}

interface EmptyStateProps {
    text: string;
}

const getLocale = (language: string) => {
    switch (language) {
        case 'ru':
            return ru;
        case 'sr':
            return sr;
        default:
            return enUS;
    }
};

const formatMessageTime = (date: string | undefined, language: string): string => {
    if (!date) return '';
    return formatDistanceToNow(new Date(date), {
        addSuffix: true,
        locale: getLocale(language)
    });
};

const MessageContent: React.FC<MessageContentProps> = ({ content }) => {
    const { t } = useTranslation('marketplace');
    const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);
    const [selectedPhone, setSelectedPhone] = useState<PhonePart | null>(null);

    const isOnlyEmoji = (text: string): boolean => {
        const emojiRegex = /^(?:\p{Extended_Pictographic}|\p{Emoji_Presentation}|\p{Emoji}\uFE0F|\p{Emoji_Modifier})+$/u;
        return emojiRegex.test(text.trim());
    };

    const parsePhoneNumbers = (text: string): MessagePart[] => {
        // Обновленное регулярное выражение для сербских номеров
        const phoneRegex = /(?:(?:\+381|0)[\s.-]?(?:6[0-9])[\s.-]?[0-9]{3}[\s.-]?[0-9]{3,4})|(?:\+?\d{1,3}[-.\s]?)?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}/g;
        const parts: MessagePart[] = [];
        let lastIndex = 0;
        let match;

        while ((match = phoneRegex.exec(text)) !== null) {
            // Добавляем текст до телефона
            if (match.index > lastIndex) {
                parts.push({
                    type: 'text',
                    content: text.slice(lastIndex, match.index)
                });
            }

            // Нормализуем номер телефона
            let phoneNumber = match[0].replace(/[-\s()]/g, '');

            // Преобразуем номер в международный формат если начинается с 0
            if (phoneNumber.startsWith('0')) {
                phoneNumber = '+381' + phoneNumber.substring(1);
            }

            // Добавляем телефон
            parts.push({
                type: 'phone',
                content: match[0],
                phoneNumber: phoneNumber
            });

            lastIndex = match.index + match[0].length;
        }

        // Добавляем оставшийся текст после последнего телефона
        if (lastIndex < text.length) {
            parts.push({
                type: 'text',
                content: text.slice(lastIndex)
            });
        }

        return parts.length > 0 ? parts : [{ type: 'text', content: text }];
    };

    const parseMessage = (text: string): MessagePart[] => {
        const emojiRegex = /((?:\p{Extended_Pictographic}|\p{Emoji_Presentation}|\p{Emoji}\uFE0F|\p{Emoji_Modifier})+)/u;
        const parts: MessagePart[] = [];

        const phoneParts = parsePhoneNumbers(text);

        phoneParts.forEach(part => {
            if (part.type === 'phone') {
                parts.push(part);
            } else {
                const emojiParts = part.content.split(emojiRegex).filter(Boolean);
                emojiParts.forEach(textPart => {
                    const isEmoji = emojiRegex.test(textPart);
                    parts.push({
                        type: isEmoji ? 'emoji' : 'text',
                        content: textPart
                    });
                });
            }
        });

        return parts;
    };

    const handlePhoneClick = (event: React.MouseEvent<HTMLElement>, phone: PhonePart): void => {
        event.preventDefault();
        setSelectedPhone(phone);
        setAnchorEl(event.currentTarget);
    };

    const handleClose = (): void => {
        setAnchorEl(null);
    };

    const handleCopy = (): void => {
        if (selectedPhone) {
            navigator.clipboard.writeText(selectedPhone.phoneNumber);
        }
        handleClose();
    };

    const handleCall = (): void => {
        if (selectedPhone) {
            window.location.href = `tel:${selectedPhone.phoneNumber}`;
        }
        handleClose();
    };

    const onlyEmoji = isOnlyEmoji(content);
    const parts = parseMessage(content);

    return (
        <Typography
            variant="body2"
            component="div"
            sx={{
                '& .emoji': {
                    fontFamily: 'Apple Color Emoji, -apple-system-emoji, "Segoe UI Emoji", "Noto Color Emoji", sans-serif',
                    fontSize: onlyEmoji ? '4rem' : '1.5rem',
                    lineHeight: 1,
                    verticalAlign: 'middle',
                    fontStyle: 'normal',
                    WebkitFontSmoothing: 'antialiased',
                    textRendering: 'optimizeLegibility',
                },
                '& .phone-container': {
                    display: 'inline-flex',
                    alignItems: 'center',
                    gap: 0.5,
                    bgcolor: 'action.hover',
                    borderRadius: 1,
                    px: 0.5,
                    cursor: 'pointer',
                    '&:hover': {
                        bgcolor: 'action.selected'
                    }
                }
            }}
        >
            {parts.map((part, index) => {
                switch (part.type) {
                    case 'emoji':
                        return (
                            <span key={index} className="emoji">
                                {part.content}
                            </span>
                        );
                    case 'phone':
                        return (
                            <Box
                                key={index}
                                component="span"
                                className="phone-container"
                                onClick={(e) => handlePhoneClick(e, part)}
                            >
                                {part.content}
                                <IconButton
                                    size="small"
                                    sx={{
                                        ml: 0.5,
                                        p: 0.3,
                                        color: 'success.main',
                                        '&:hover': {
                                            bgcolor: 'success.lighter'
                                        }
                                    }}
                                >
                                    <Phone fontSize="small" />
                                </IconButton>
                            </Box>
                        );
                    default:
                        return <span key={index}>{part.content}</span>;
                }
            })}

            <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleClose}
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'left',
                }}
                transformOrigin={{
                    vertical: 'top',
                    horizontal: 'left',
                }}
            >
                <MenuItem onClick={handleCall}>
                    <PhoneIcon sx={{ mr: 1, color: 'success.main' }} fontSize="small" />
                    {t('chat.call')}
                </MenuItem>
                <MenuItem onClick={handleCopy}>
                    <ContentCopy sx={{ mr: 1 }} fontSize="small" />
                    {t('chat.copyPhone')}
                </MenuItem>
            </Menu>
        </Typography>
    );
};

export const ChatWindow: React.FC<ChatWindowProps> = ({ messages = [], onSendMessage, currentUser, chat, onBack }) => {
    const { t, i18n } = useTranslation('marketplace');
    const [newMessage, setNewMessage] = useState('');
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const [processedMessages, setProcessedMessages] = useState<ChatMessage[]>([]);
    const [isTyping, setIsTyping] = useState(false);
    const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);

    useEffect(() => {
        const uniqueMessages = Object.values(
            messages.reduce((acc: Record<string | number, ChatMessage>, message) => {
                acc[message.id] = message;
                return acc;
            }, {})
        ).sort((a: ChatMessage, b: ChatMessage) => new Date(a.created_at || '').getTime() - new Date(b.created_at || '').getTime());

        setProcessedMessages(uniqueMessages);
    }, [messages]);

    const scrollToBottom = (): void => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    useEffect(() => {
        scrollToBottom();
    }, [processedMessages]);

    const handleSend = (e: React.FormEvent): void => {
        e.preventDefault();
        if (newMessage.trim()) {
            onSendMessage(newMessage.trim());
            setNewMessage('');
        }
    };

    const handleEmojiClick = (emojiData: EmojiClickData): void => {
        setNewMessage((prevMessage) => prevMessage + emojiData.emoji);
        setAnchorEl(null);
    };

    const handleEmojiButtonClick = (event: React.MouseEvent<HTMLElement>): void => {
        setAnchorEl(anchorEl ? null : event.currentTarget);
    };

    return (
        <Box
            sx={{
                display: 'flex',
                flexDirection: 'column',
                height: '100%',
                overflow: 'hidden',
                bgcolor: 'grey.50'
            }}
        >
            {/* Область сообщений */}
            <Box sx={{
                flex: 1,
                overflowY: 'auto',
                overflowX: 'hidden',
                p: 2,
                WebkitOverflowScrolling: 'touch',
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
                                {formatMessageTime(message.created_at, i18n.language)}
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
                    borderColor: 'divider',
                    position: 'sticky',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    zIndex: 2
                }}
            >
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <IconButton type="button">
                        <Paperclip />
                    </IconButton>
                    <Box sx={{
                        flex: 1,
                        position: 'relative'
                    }}>
                        <TextField
                            fullWidth
                            size="small"
                            placeholder={t('chat.placeholder')}
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
                            type="button"
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
                                vertical: 'bottom',
                                horizontal: 'right',
                            }}
                            transformOrigin={{
                                vertical: 'top',
                                horizontal: 'right',
                            }}
                        >
                            <EmojiPicker
                            onEmojiClick={handleEmojiClick}
                            width={320}
                            height={400}
                            searchDisabled={true}
                            skinTonesDisabled={true}
                            emojiStyle={EmojiStyle.NATIVE}
                            lazyLoadEmojis={true}
                            style={{
                            '--epr-bg-color': 'white',
                            '--epr-category-label-bg-color': 'white'
                        } as React.CSSProperties}
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
export const ChatList: React.FC<ChatListProps> = ({ chats, selectedChatId, onSelectChat, onArchiveChat }) => {
    const { t, i18n } = useTranslation('marketplace');
    
    const formatPrice = (price?: number): string => {
        if (!price || price === undefined) return 'Цена не указана';

        return new Intl.NumberFormat('sr-RS', {
            style: 'currency',
            currency: 'RSD',
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
                                        {chat.listing && formatPrice(chat.listing.price)}
                                    </Typography>
                                </Box>
                                {chat.unread_count && chat.unread_count > 0 && (
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
                                    {formatMessageTime(chat.last_message_at, i18n.language)}
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
export const ChatHeader: React.FC<ChatHeaderProps> = ({ chat, onBack, onArchive }) => {
    const { t } = useTranslation('marketplace');

    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
    
    const formatPrice = (price?: number): string => {
        if (!price || price === undefined) return 'Цена не указана';

        return new Intl.NumberFormat('sr-RS', {
            style: 'currency',
            currency: 'RSD',
            maximumFractionDigits: 0
        }).format(price);
    };

    const formatLastSeen = (date?: string): string => {
        if (!date) return 'не в сети';
        const lastSeen = new Date(date);
        const now = new Date();
        const diffMinutes = Math.floor((now.getTime() - lastSeen.getTime()) / (1000 * 60));

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
                    {isMobile && onBack && (
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
                        {onArchive && (
                            <IconButton
                                onClick={() => onArchive(chat.id)}
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
                        )}

                        {/* Кнопка звонка */}
                        {chat.other_user?.phone && (
                            <IconButton
                                href={`tel:${chat.other_user.phone}`}
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
                        )}

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
                            {t('chat.openlisting')}
                        </Button>
                    </Stack>
                </Stack>
            </Box>
        </Paper>
    );
};

// Компонент пустого состояния
export const EmptyState: React.FC<EmptyStateProps> = ({ text }) => (
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