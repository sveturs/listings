// frontend/hostel-frontend/src/components/balance/DepositDialog.tsx

import React, { useState, useEffect } from 'react';

import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Box,
  RadioGroup,
  FormControlLabel,
  Radio,
  Typography,
  Alert,
  InputAdornment,
  CircularProgress
} from '@mui/material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';

interface PaymentMethod {
  code: string;
  name: string;
  fixed_fee: number;
  fee_percentage: number;
  minimum_amount: number;
  maximum_amount: number;
}

interface PaymentSession {
  payment_url: string;
  session_id: string;
}

interface DepositDialogProps {
  open: boolean;
  onClose: () => void;
  onBalanceUpdate?: () => void;
}

const DepositDialog: React.FC<DepositDialogProps> = ({ open, onClose, onBalanceUpdate }) => {
  const { t } = useTranslation('common');
  const [amount, setAmount] = useState<string>('');
  const [paymentMethod, setPaymentMethod] = useState<string>('');
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
  const [error, setError] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    const fetchPaymentMethods = async (): Promise<void> => {
      try {
        const response = await axios.get('/api/v1/balance/payment-methods');
        setPaymentMethods(response.data.data);
      } catch (err) {
        setError(t('balance.errors.loadPaymentMethods'));
      }
    };

    if (open) {
      fetchPaymentMethods();
    }
  }, [open, t]);

  const handleDeposit = async (): Promise<void> => {
    setError('');
    setLoading(true);

    try {
      const response = await axios.post('/api/v1/balance/deposit', {
        amount: parseFloat(amount),
        payment_method: paymentMethod
      });

      if (!response.data?.success) {
        throw new Error(response.data?.message || 'Unknown error');
      }

      // Получаем URL для оплаты от Stripe
      const paymentSession: PaymentSession = response.data.data;
      
      if (paymentSession.payment_url) {
        // Перенаправляем на страницу оплаты Stripe
        window.location.href = paymentSession.payment_url;
      } else {
        throw new Error('Payment URL not provided');
      }
    } catch (err: any) {
      console.error('Deposit error:', err);
      setError(err.response?.data?.message || t('balance.errors.depositFailed'));
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{t('balance.deposit')}</DialogTitle>
      <DialogContent>
        {error && (
          <Alert 
            severity="error" 
            sx={{ mb: 2 }}
            onClose={() => setError('')}
          >
            {error}
          </Alert>
        )}

        <TextField
          label={t('balance.amount')}
          type="number"
          fullWidth
          value={amount}
          onChange={(e) => {
            const value = parseFloat(e.target.value);
            if (!isNaN(value) && value > 0) {
              setAmount(e.target.value);
            }
          }}
          error={!!error}
          disabled={loading}
          sx={{ mb: 2, mt: 1 }}
          InputProps={{
            endAdornment: <InputAdornment position="end">RSD</InputAdornment>,
          }}
        />

        <Typography variant="subtitle2" gutterBottom>
          {t('balance.paymentMethod')}
        </Typography>
        <RadioGroup
          value={paymentMethod}
          onChange={(e) => setPaymentMethod(e.target.value)}
        >
          {paymentMethods.map((method) => (
            <FormControlLabel
              key={method.code}
              value={method.code}
              control={<Radio disabled={loading} />}
              label={
                <Box>
                  <Typography>{method.name}</Typography>
                  <Typography variant="caption" color="text.secondary">
                    {t('balance.fee', {
                      fixed: method.fixed_fee,
                      percent: method.fee_percentage
                    })}
                    {method.minimum_amount > 0 && ` • ${t('balance.minAmount')}: ${method.minimum_amount} RSD`}
                    {method.maximum_amount > 0 && ` • ${t('balance.maxAmount')}: ${method.maximum_amount} RSD`}
                  </Typography>
                </Box>
              }
            />
          ))}
        </RadioGroup>
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose} disabled={loading}>
          {t('buttons.cancel')}
        </Button>
        <Button
          variant="contained"
          onClick={handleDeposit}
          disabled={!amount || !paymentMethod || loading}
        >
          {loading ? (
            <CircularProgress size={24} color="inherit" />
          ) : (
            t('balance.confirmDeposit')
          )}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default DepositDialog;