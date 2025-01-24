// frontend/hostel-frontend/src/components/marketplace/chat/ChatButton.js
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
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

const ChatButton = ({ listing, isMobile }) => {
    const navigate = useNavigate();
    const { user, login } = useAuth();
    const [open, setOpen] = useState(false);
    const [message, setMessage] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleClick = () => {
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

    const handleSend = async () => {
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
        } catch (error) {
            console.error('Error sending message:', error);
            setError(
                error.response?.data?.message || 
                'Не удалось отправить сообщение. Пожалуйста, попробуйте позже.'
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
                {isMobile ? <MessageCircle size={20} /> : 'Написать'}
            </Button>

            <Dialog 
                open={open} 
                onClose={() => !loading && setOpen(false)} 
                maxWidth="sm" 
                fullWidth
            >
                <DialogTitle>Написать продавцу</DialogTitle>
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
                        placeholder="Введите сообщение..."
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        sx={{ mt: 2 }}
                        disabled={loading}
                    />
                </DialogContent>
                <DialogActions>
                    <Button 
                        onClick={() => setOpen(false)}
                        disabled={loading}
                    >
                        Отмена
                    </Button>
                    <Button
                        variant="contained"
                        onClick={handleSend}
                        disabled={!message.trim() || loading}
                    >
                        {loading ? 'Отправка...' : 'Отправить'}
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    );
};

export default ChatButton;