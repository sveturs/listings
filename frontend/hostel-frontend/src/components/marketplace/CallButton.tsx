import React, { useState, MouseEvent } from 'react';
import { useTranslation } from 'react-i18next';
import { Button, Box, CircularProgress } from '@mui/material';
import { Phone } from 'lucide-react';
import PhonePopup from './PhonePopup';

interface CallButtonProps {
  phone: string;
  isMobile?: boolean;
}

const CallButton: React.FC<CallButtonProps> = ({ phone, isMobile = false }) => {
  const { t } = useTranslation('marketplace') as any;
  const [showPhone, setShowPhone] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(false);

  const formatPhone = (phoneNumber: string): string => {
    const cleaned = phoneNumber.replace(/[^\d]/g, '');
    return cleaned.startsWith('381') 
      ? cleaned.replace(/(\d{3})(\d{2})(\d{3})(\d{2})(\d{2})/, '+$1 $2 $3 $4 $5')
      : cleaned.replace(/(\d{1})(\d{3})(\d{3})(\d{2})(\d{2})/, '+$1 $2 $3 $4 $5');
  };

  const handleMouseEnter = (): void => !isMobile && setShowPhone(true);
  const handleMouseLeave = (): void => !isMobile && setShowPhone(false);

  const handleCallClick = (): void => {
    if (!phone) {
      alert(t('listings.details.contact.errors.noPhone'));
      return;
    }

    if (isMobile) {
      setShowPhone(true);
      setLoading(true);
      setTimeout(() => {
        window.location.href = `tel:${phone}`;
        setShowPhone(false);
        setLoading(false);
      }, 1500);
    } else {
      window.location.href = `tel:${phone}`;
    }
  };

  return (
    <Box sx={{ position: 'relative', width: '100%' }}>
      <Button
        variant="contained"
        fullWidth
        startIcon={<Phone size={20} />}
        onClick={handleCallClick}
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
        sx={{ 
          height: 48,
          backgroundColor: 'success.main',
          '&:hover': { backgroundColor: 'success.dark' }
        }}
      >
        {isMobile && showPhone ? (
          <Box sx={{ 
            display: 'flex', 
            alignItems: 'center', 
            gap: 1,
            fontSize: '0.875rem',
            whiteSpace: 'nowrap'
          }}>
            {formatPhone(phone)}
            {loading && (
              <CircularProgress
                size={16}
                sx={{ color: 'common.white' }}
              />
            )}
          </Box>
        ) : (
          t('listings.details.contact.call')
        )}
      </Button>

      {!isMobile && <PhonePopup phone={phone} visible={showPhone} onClose={() => setShowPhone(false)} />}
    </Box>
  );
};

export default CallButton;