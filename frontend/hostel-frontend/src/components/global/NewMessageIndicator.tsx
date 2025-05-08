import React, { useState, useEffect } from 'react';
import { Badge } from '@mui/material';
import { MessageCircle } from 'lucide-react';
import { keyframes } from '@emotion/react';

interface NewMessageIndicatorProps {
  unreadCount?: number;
}

const pulseAnimation = keyframes`
  0% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.2);
    opacity: 0.7;
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
`;

const NewMessageIndicator: React.FC<NewMessageIndicatorProps> = ({ unreadCount = 0 }) => {
  const [shouldPulse, setShouldPulse] = useState<boolean>(false);

  useEffect(() => {
    if (unreadCount > 0) {
      setShouldPulse(true);
      const timer = setTimeout(() => setShouldPulse(false), 1000);
      return () => clearTimeout(timer);
    }
  }, [unreadCount]);

  return (
    <Badge
      badgeContent={unreadCount}
      color="error"
      overlap="circular"
      sx={{
        '& .MuiBadge-badge': {
          animation: shouldPulse ? `${pulseAnimation} 1s ease-in-out` : 'none',
        }
      }}
    >
      <MessageCircle size={20} />
    </Badge>
  );
};

export default NewMessageIndicator;