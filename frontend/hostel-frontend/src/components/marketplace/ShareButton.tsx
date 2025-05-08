import React, { useState, MouseEvent } from 'react';
import { useTranslation } from 'react-i18next';

import {
    Button,
    Menu,
    MenuItem,
    ListItemIcon,
    ListItemText,
    Snackbar,
    Alert
} from '@mui/material';
import { Share2, Copy, Twitter, MessageCircle } from 'lucide-react';
import { Telegram, Facebook as FacebookIcon } from '@mui/icons-material';

interface ShareButtonProps {
    url: string;
    title: string;
    isMobile?: boolean;
}

type SharePlatform = 'facebook' | 'twitter' | 'telegram' | 'viber';

const ShareButton: React.FC<ShareButtonProps> = ({ url, title, isMobile = false }) => {
    const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);
    const [showCopyAlert, setShowCopyAlert] = useState<boolean>(false);
    const { t } = useTranslation('marketplace') as any; 
    
    const handleClick = (event: MouseEvent<HTMLButtonElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const handleShare = (platform: SharePlatform) => {
        const shareUrls: Record<SharePlatform, string> = {
            facebook: `https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(url)}`,
            twitter: `https://twitter.com/intent/tweet?url=${encodeURIComponent(url)}&text=${encodeURIComponent(title)}`,
            telegram: `https://t.me/share/url?url=${encodeURIComponent(url)}&text=${encodeURIComponent(title)}`,
            viber: `viber://forward?text=${encodeURIComponent(title + ' ' + url)}`
        };

        if (shareUrls[platform]) {
            window.open(shareUrls[platform], '_blank');
        }
        handleClose();
    };

    const copyToClipboard = async (): Promise<void> => {
        try {
            await navigator.clipboard.writeText(url);
            setShowCopyAlert(true);
        } catch (err) {
            console.error('Failed to copy:', err);
        }
        handleClose();
    };

    return (
        <>
            <Button
                id="shareButton2"
                variant="outlined"
                fullWidth
                startIcon={!isMobile && <Share2 />}
                onClick={handleClick}
            >
                {isMobile ? <Share2 size={20} /> : t('listings.details.contact.share')}
            </Button>
            
            <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleClose}
                PaperProps={{
                    elevation: 3,
                    sx: { width: 220 }
                }}
            >
                <MenuItem onClick={() => handleShare('facebook')}>
                    <ListItemIcon>
                        <FacebookIcon fontSize="small" color="primary" />
                    </ListItemIcon>
                    <ListItemText>Facebook</ListItemText>
                </MenuItem>

                <MenuItem onClick={copyToClipboard}>
                    <ListItemIcon>
                        <Copy size={20} />
                    </ListItemIcon>
                    <ListItemText>Копировать ссылку</ListItemText>
                </MenuItem>

                <MenuItem onClick={() => handleShare('telegram')}>
                    <ListItemIcon>
                        <Telegram fontSize="small" color="primary" />
                    </ListItemIcon>
                    <ListItemText>Telegram</ListItemText>
                </MenuItem>

                <MenuItem onClick={() => handleShare('viber')}>
                    <ListItemIcon>
                        <MessageCircle size={20} color="#665CAC" />
                    </ListItemIcon>
                    <ListItemText>Viber</ListItemText>
                </MenuItem>

                <MenuItem onClick={() => handleShare('twitter')}>
                    <ListItemIcon>
                        <Twitter size={20} color="#1DA1F2" />
                    </ListItemIcon>
                    <ListItemText>Twitter</ListItemText>
                </MenuItem>
            </Menu>

            <Snackbar
                open={showCopyAlert}
                autoHideDuration={3000}
                onClose={() => setShowCopyAlert(false)}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
            >
                <Alert severity="success" variant="filled">
                    Ссылка скопирована
                </Alert>
            </Snackbar>
        </>
    );
};

export default ShareButton;