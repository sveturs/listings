import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../../../contexts/AuthContext';
import {
    Button,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    TextField,
    Alert,
} from '@mui/material';
import { MessageCircle } from 'lucide-react';
import axios from '../../../api/axios';

interface Listing {
    id: string | number;
    user_id: string | number;
    title?: string;
    [key: string]: any;
}

interface ChatButtonProps {
    listing: Listing;
    isMobile?: boolean;
}

const ChatButton: React.FC<ChatButtonProps> = ({ listing, isMobile = false }) => {
    const navigate = useNavigate();
    const { user, login } = useAuth();
    const { t } = useTranslation('marketplace');
    const [open, setOpen] = useState<boolean>(false);
    const [message, setMessage] = useState<string>('');
    const [error, setError] = useState<string>('');
    const [loading, setLoading] = useState<boolean>(false);

    const handleClick = (): void => {
        if (!user) {
            const returnUrl = window.location.pathname;
            const encodedReturnUrl = encodeURIComponent(returnUrl);
            login(`?returnTo=${encodedReturnUrl}`);
            return;
        }

        if (user.id === listing.user_id) {
            navigate('/marketplace/chat');
            return;
        }

        setOpen(true);
    };

    const handleSend = async (): Promise<void> => {
        if (!message.trim()) return;
        
        setLoading(true);
        setError('');
    
        try {
            await axios.post('/api/v1/marketplace/chat/messages', {
                listing_id: listing.id,
                receiver_id: listing.user_id,
                content: message.trim()
            }, {
                withCredentials: true  
            });
            
            setOpen(false);
            setMessage('');
            navigate('/marketplace/chat');
        } catch (error: any) {
            console.error('Error sending message:', error);
            setError(
                error.response?.data?.message || 
                t('chat.errors.sendFailed')
            );
        } finally {
            setLoading(false);
        }
    };

    return (
        <>
            <Button
                id="chatButton"
                variant="outlined"
                fullWidth
                startIcon={!isMobile && <MessageCircle />}
                onClick={handleClick}
            >
                {isMobile ? <MessageCircle size={20} /> : t('listings.details.contact.message')}
            </Button>

            <Dialog 
                open={open} 
                onClose={() => !loading && setOpen(false)} 
                maxWidth="sm" 
                fullWidth
            >
                <DialogTitle>{t('chat.newMessage')}</DialogTitle>
                <DialogContent>
                    {error && (
                        <Alert severity="error" sx={{ mb: 2 }}>
                            {error}
                        </Alert>
                    )}
                    <TextField
                        autoFocus
                        fullWidth
                        multiline
                        rows={4}
                        placeholder={t('chat.placeholder')}
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        sx={{ mt: 2 }}
                        disabled={loading}
                    />
                </DialogContent>
                <DialogActions>
                    <Button
                        id="cancelMessageButton"
                        onClick={() => setOpen(false)}
                        disabled={loading}
                    >
                        {t('reviews.cancel')}
                    </Button>
                    <Button
                        id="sendMessageButton"
                        variant="contained"
                        onClick={handleSend}
                        disabled={!message.trim() || loading}
                    >
                        {loading ? t('chat.sending') : t('chat.send')}
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    );
};

export default ChatButton;