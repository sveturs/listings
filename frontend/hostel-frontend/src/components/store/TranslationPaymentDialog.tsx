// frontend/hostel-frontend/src/components/store/TranslationPaymentDialog.tsx
import React, { useState, useEffect } from 'react';
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
  CircularProgress,
  FormControl,
  FormControlLabel,
  Radio,
  RadioGroup,
  Divider
} from '@mui/material';
import { Wallet, Info, AlertTriangle, Globe, Sparkles } from 'lucide-react';
import axios from '../../api/axios';

type TranslationProvider = 'google' | 'openai';

interface GoogleTranslateInfo {
  used: number;
  limit: number;
}

interface TranslationPaymentDialogProps {
  /** Флаг открытия диалога */
  open: boolean;
  /** Функция закрытия диалога */
  onClose: () => void;
  /** Массив выбранных объявлений */
  selectedListings: Array<number | string>;
  /** Текущий баланс пользователя */
  balance: number;
  /** Функция подтверждения оплаты */
  onConfirm: (provider: TranslationProvider) => void;
  /** Флаг загрузки */
  loading?: boolean;
  /** Стоимость перевода одного объявления через OpenAI */
  costPerListing?: number;
}

/**
 * Диалог подтверждения оплаты перевода объявлений с выбором провайдера
 */
const TranslationPaymentDialog: React.FC<TranslationPaymentDialogProps> = ({ 
  open, 
  onClose, 
  selectedListings = [], 
  balance, 
  onConfirm, 
  loading = false,
  costPerListing = 25
}) => {
  const { t, i18n } = useTranslation(['marketplace', 'common']);
  // Загружаем сохраненный провайдер из localStorage или используем Google по умолчанию
  const [translationProvider, setTranslationProvider] = useState<TranslationProvider>(
    (localStorage.getItem('preferredTranslationProvider') as TranslationProvider) || 'google'
  );
  const [googleTranslateInfo, setGoogleTranslateInfo] = useState<GoogleTranslateInfo>({ used: 0, limit: 100 });
  const [loadingLimits, setLoadingLimits] = useState<boolean>(false);
  
  // Загружаем информацию о лимитах Google Translate при открытии диалога
  useEffect(() => {
    if (open) {
      setLoadingLimits(true);
      // Здесь должен быть запрос к API для получения текущего состояния лимитов
      axios.get('/api/v1/translation/limits')
        .then(response => {
          const { google } = response.data.data;
          setGoogleTranslateInfo({
            used: google.used,
            limit: google.limit
          });
        })
        .catch(error => {
          console.error('Ошибка при получении лимитов перевода:', error);
          // В случае ошибки используем значения по умолчанию
          setGoogleTranslateInfo({ used: 0, limit: 100 });
        })
        .finally(() => {
          setLoadingLimits(false);
        });
    }
  }, [open]);
  
  // Форматирование цены
  const formatPrice = (price: number): string => {
    return new Intl.NumberFormat('ru-RU', {
      style: 'currency',
      currency: 'RSD',
      maximumFractionDigits: 0
    }).format(price);
  };
  
  // Рассчитываем общую стоимость в зависимости от провайдера
  const totalCost = translationProvider === 'openai' 
    ? selectedListings.length * costPerListing 
    : 0; // Google Translate бесплатен в пределах лимитов
  
  // Проверяем достаточно ли средств для OpenAI
  const hasEnoughFunds = balance >= totalCost;
  
  // Проверяем, не превышен ли лимит Google Translate
  const googleLimitExceeded = 
    translationProvider === 'google' && 
    (googleTranslateInfo.used + selectedListings.length > googleTranslateInfo.limit);
  
  // Обработчик подтверждения с указанием провайдера
  const handleConfirm = (): void => {
    // Сохраняем выбранный провайдер в localStorage для использования при создании новых объявлений
    localStorage.setItem('preferredTranslationProvider', translationProvider);
    onConfirm(translationProvider);
  };
  
  // Получаем текущий язык страницы (язык оригинала)
  const originalLanguage = i18n.language;
  
  // Определяем целевые языки для перевода (все языки кроме текущего)
  const targetLanguages = ['ru', 'en', 'sr'].filter(lang => lang !== originalLanguage);
  
  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="sm"
      fullWidth
    >
      <DialogTitle>
        {t('store.listings.translateConfirmTitle', {defaultValue: 'Подтверждение перевода'})}
      </DialogTitle>
      
      <DialogContent>
        <Typography variant="body1" paragraph>
          {t('store.listings.translateConfirmText', {
            count: selectedListings.length,
            defaultValue: 'Вы собираетесь перевести {{count}} объявлений на все доступные языки.'
          })}
        </Typography>
        
        <Alert 
          severity="info" 
          icon={<Info />}
          sx={{ mb: 2 }}
        >
          {t('store.listings.translateLanguageInfo', {
            sourceLanguage: originalLanguage === 'ru' ? 'Русский' : 
                            originalLanguage === 'en' ? 'English' : 'Српски',
            targetLanguages: targetLanguages.map(lang => 
              lang === 'ru' ? 'Русский' : 
              lang === 'en' ? 'English' : 'Српски'
            ).join(', '),
            defaultValue: 'Перевод будет выполнен с языка: {{sourceLanguage}} на языки: {{targetLanguages}}.'
          })}
        </Alert>
        
        <Box sx={{ mb: 3 }}>
          <Typography variant="subtitle1" gutterBottom>
            {t('store.listings.selectTranslationProvider', {defaultValue: 'Выберите сервис перевода'})}
          </Typography>
          
          <FormControl component="fieldset">
            <RadioGroup
              value={translationProvider}
              onChange={(e) => setTranslationProvider(e.target.value as TranslationProvider)}
            >
              <FormControlLabel 
                value="google" 
                control={<Radio />} 
                label={
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Globe size={18} />
                    <Typography>
                      {t('store.listings.googleTranslate', {defaultValue: 'Google Translate'})}
                    </Typography>
                    <Typography variant="caption" color="success.main">
                      ({t('store.listings.freeWithinLimits', {defaultValue: 'Бесплатно в пределах лимита'})})
                    </Typography>
                  </Box>
                }
              />
              
              <Box pl={4} mb={1}>
                <Typography variant="caption" color="text.secondary">
                  {loadingLimits ? (
                    <CircularProgress size={12} />
                  ) : (
                    t('store.listings.googleTranslateLimits', {
                      used: googleTranslateInfo.used,
                      limit: googleTranslateInfo.limit,
                      defaultValue: 'Использовано {{used}} из {{limit}} в этом месяце.'
                    })
                  )}
                </Typography>
                
                {googleLimitExceeded && (
                  <Alert 
                    severity="warning" 
                    icon={<AlertTriangle size={16} />}
                    sx={{ mt: 1, py: 0.5 }}
                  >
                    {t('store.listings.googleLimitExceeded', {
                      required: selectedListings.length,
                      available: googleTranslateInfo.limit - googleTranslateInfo.used,
                      defaultValue: 'Недостаточно бесплатных переводов. Требуется: {{required}}, доступно: {{available}}.'
                    })}
                  </Alert>
                )}
              </Box>
              
              <FormControlLabel 
                value="openai" 
                control={<Radio />} 
                label={
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Sparkles size={18} />
                    <Typography>
                      {t('store.listings.openaiTranslate', {defaultValue: 'OpenAI (GPT)'})}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      ({formatPrice(costPerListing)} {t('store.listings.perListing', {defaultValue: 'за объявление'})})
                    </Typography>
                  </Box>
                }
              />
              
              <Box pl={4}>
                <Typography variant="caption" color="text.secondary">
                  {t('store.listings.openaiInfo', {
                    defaultValue: 'Более высокое качество перевода, особенно для сложных текстов.'
                  })}
                </Typography>
              </Box>
            </RadioGroup>
          </FormControl>
        </Box>
        
        <Divider sx={{ my: 2 }} />
        
        {translationProvider === 'openai' && (
          <>
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
                  {t('store.listings.translationCost', {defaultValue: 'Стоимость перевода'})}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  {t('store.listings.costPerListing', {
                    cost: formatPrice(costPerListing),
                    defaultValue: '{{cost}} за объявление'
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
                  {t('common:balance.available', {defaultValue: 'Доступный баланс'})}
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
                  defaultValue: 'Недостаточно средств для перевода. Требуется: {{required}}, доступно: {{available}}'
                })}
              </Alert>
            )}
          </>
        )}
        
        <Alert 
          severity="info" 
          icon={<Info />}
        >
          {translationProvider === 'google' 
            ? t('store.listings.googleTranslationInfo', {
                defaultValue: 'Перевод будет выполнен с помощью Google Translate. Качество может варьироваться в зависимости от языка и сложности текста.'
              })
            : t('store.listings.openaiTranslationInfo', {
                defaultValue: 'Перевод будет выполнен с помощью нейросети OpenAI GPT, что обеспечивает высокое качество перевода для большинства текстов.'
              })
          }
        </Alert>
      </DialogContent>
      
      <DialogActions>
        <Button 
          variant="outlined" 
          onClick={onClose}
          disabled={loading}
        >
          {t('common:buttons.cancel', {defaultValue: 'Отмена'})}
        </Button>
        
        <Button
          variant="contained"
          color="primary"
          onClick={handleConfirm}
          disabled={
            loading || 
            (translationProvider === 'openai' && !hasEnoughFunds) ||
            (translationProvider === 'google' && googleLimitExceeded)
          }
          startIcon={loading ? <CircularProgress size={20} color="inherit" /> : null}
        >
          {loading 
            ? t('store.listings.translating', {defaultValue: 'Перевод...'})
            : t('store.listings.confirmTranslation', {defaultValue: 'Подтвердить перевод'})
          }
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default TranslationPaymentDialog;