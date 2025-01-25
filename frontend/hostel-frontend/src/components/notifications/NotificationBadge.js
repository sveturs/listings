import React from 'react';
import { Badge, IconButton } from '@mui/material';
import { Bell } from 'lucide-react';
import { useNotifications } from '../../hooks/useNotifications';

const NotificationBadge = ({ onClick }) => {
    const { unreadCount } = useNotifications();

    return (
        <IconButton onClick={onClick}>
            <Badge 
                badgeContent={unreadCount} 
                color="error"
                max={99}
                overlap="circular"
                sx={{
                    '& .MuiBadge-badge': {
                        fontSize: '0.75rem',
                        height: '20px',
                        minWidth: '20px',
                    }
                }}
            >
                <Bell size={20} />
            </Badge>
        </IconButton>
    );
};

export default NotificationBadge;