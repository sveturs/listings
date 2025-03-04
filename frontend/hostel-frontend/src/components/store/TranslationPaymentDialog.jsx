// frontend/hostel-frontend/src/components/store/TranslationPaymentDialog.jsx
import React from 'react';
import { useTranslation } from 'react-i18next';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Alert,
  Box,
  CircularProgress
} from '@mui/material';
import { Wallet, Info, AlertTriangle } from 'lucide-react';

/**
 * Диалог подтверждения оплаты перевода объявлений
 * 
 * @param {Object} props - Свойства компонента
 * @param {boolean} props.open - Флаг открытия диалога
 * @param {Function} props.onClose - Функция закрытия диалога
 * @param {Array} props.selectedListings - Массив выбранных объявлений
 * @param {number} props.balance - Текущий баланс пользователя
 * @param {Function} props.onConfirm - Функция подтверждения оплаты
 * @param {boolean} props.loading - Флаг загрузки
 * @param {number} props.costPerListing - Стоимость перевода одного объявления
 * @returns {JSX.Element}
 */
const TranslationPaymentDialog = ({ 
  open, 
  onClose, 
  selectedListings = [], 
  balance, 
  onConfirm, 
  loading = false,
  costPerListing = 25
}) => {
  const { t } = useTranslation(['marketplace', 'common']);
  
  // Форматирование цены
  const formatPrice = (price) => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RSD',
      maximumFractionDigits: 0
    }).format(price);
  };
  
  // Рассчитываем общую стоимость
  const totalCost = selectedListings.length * costPerListing;
  
  // Проверяем достаточно ли средств
  const hasEnoughFunds = balance >= totalCost;
  
  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="sm"
      fullWidth
    >
      <DialogTitle>
        {t('store.listings.translateConfirmTitle', {defaultValue: 'Confirm Translation'})}
      </DialogTitle>
      
      <DialogContent>
        <Typography variant="body1" paragraph>
          {t('store.listings.translateConfirmText', {
            count: selectedListings.length,
            defaultValue: 'You are about to translate {{count}} listings to all available languages.'
          })}
        </Typography>
        
        <Box 
          sx={{ 
            display: 'flex', 
            alignItems: 'center', 
            justifyContent: 'space-between',
            p: 2,
            mb: 2,
            bgcolor: 'background.paper',
            border: 1,
            borderColor: 'divider',
            borderRadius: 1
          }}
        >
          <Box>
            <Typography variant="subtitle2">
              {t('store.listings.translationCost', {defaultValue: 'Translation cost'})}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {t('store.listings.costPerListing', {
                cost: formatPrice(costPerListing),
                defaultValue: '{{cost}} per listing'
              })}
            </Typography>
          </Box>
          <Typography variant="h6" color="primary.main" fontWeight="bold">
            {formatPrice(totalCost)}
          </Typography>
        </Box>
        
        <Box 
          sx={{ 
            display: 'flex', 
            alignItems: 'center', 
            p: 2,
            mb: 2,
            bgcolor: 'background.paper',
            border: 1,
            borderColor: 'divider',
            borderRadius: 1
          }}
        >
          <Wallet size={24} style={{ marginRight: 12, opacity: 0.7 }} />
          <Box sx={{ flex: 1 }}>
            <Typography variant="subtitle2">
              {t('common:balance.available', {defaultValue: 'Available balance'})}
            </Typography>
          </Box>
          <Typography 
            variant="h6" 
            color={hasEnoughFunds ? 'success.main' : 'error.main'}
            fontWeight="bold"
          >
            {formatPrice(balance)}
          </Typography>
        </Box>
        
        {!hasEnoughFunds && (
          <Alert 
            severity="error" 
            icon={<AlertTriangle />}
            sx={{ mb: 2 }}
          >
            {t('store.listings.insufficientFunds', {
              required: formatPrice(totalCost),
              available: formatPrice(balance),
              defaultValue: 'Insufficient funds for translation. Required: {{required}}, Available: {{available}}'
            })}
          </Alert>
        )}
        
        <Alert 
          severity="info" 
          icon={<Info />}
        >
          {t('store.listings.translationInfo', {
            defaultValue: 'Translation will be performed using machine translation. The quality may vary depending on the language and text complexity.'
          })}
        </Alert>
      </DialogContent>
      
      <DialogActions>
        <Button 
          variant="outlined" 
          onClick={onClose}
          disabled={loading}
        >
          {t('common:buttons.cancel', {defaultValue: 'Cancel'})}
        </Button>
        
        <Button
          variant="contained"
          color="primary"
          onClick={onConfirm}
          disabled={!hasEnoughFunds || loading}
          startIcon={loading ? <CircularProgress size={20} color="inherit" /> : null}
        >
          {loading 
            ? t('store.listings.translating', {defaultValue: 'Translating...'})
            : t('store.listings.confirmTranslation', {defaultValue: 'Confirm Translation'})
          }
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default TranslationPaymentDialog;