// frontend/hostel-frontend/src/components/balance/BalanceWidget.tsx

import React from 'react';
import { Box, Typography, Button, Chip } from '@mui/material';
import { Wallet } from 'lucide-react';
import { useTranslation } from 'react-i18next';

interface BalanceWidgetProps {
  balance: number;
  onDeposit: () => void;
}

const BalanceWidget: React.FC<BalanceWidgetProps> = ({ balance, onDeposit }) => {
  const { t } = useTranslation('common');

  const formatBalance = (amount: number): string => {
    return new Intl.NumberFormat('sr-RS', {
      style: 'currency',
      currency: 'RSD',
      maximumFractionDigits: 0
    }).format(amount);
  };

  return (
    <Box sx={{ 
      display: 'flex', 
      alignItems: 'center', 
      gap: 2, 
      p: 2, 
      bgcolor: 'background.paper',
      borderRadius: 1,
      boxShadow: 1
    }}>
      <Wallet size={24} />
      <Box sx={{ flex: 1 }}>
        <Typography variant="body2" color="text.secondary">
          {t('balance.available')}
        </Typography>
        <Typography variant="h6">
          {formatBalance(balance)}
        </Typography>
      </Box>
      <Button 
        variant="contained" 
        color="primary"
        onClick={onDeposit}
      >
        {t('balance.deposit')}
      </Button>
    </Box>
  );
};

export default BalanceWidget;