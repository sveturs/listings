// frontend/hostel-frontend/src/pages/balance/TransactionsPage.js
import React, { useState, useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import {
  Container,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  Chip,
  Box,
  Alert,
  CircularProgress,
  Snackbar
} from '@mui/material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';
import BalanceWidget from '../../components/balance/BalanceWidget';
import DepositDialog from '../../components/balance/DepositDialog';

const TransactionsPage = () => {
  const { t, i18n } = useTranslation('common');
  const { checkAuth } = useAuth(); 
  const [transactions, setTransactions] = useState([]);
  const [balance, setBalance] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [depositDialogOpen, setDepositDialogOpen] = useState(false);
  const [paymentStatus, setPaymentStatus] = useState(null);
  const location = useLocation();
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState('');

  const handleBalanceUpdate = (newBalance) => {
    setBalance(newBalance);
    fetchData(); // Обновляем список транзакций
  };

  const fetchData = async () => {
    try {
      console.log("Fetching balance data...");
      const [balanceRes, transactionsRes] = await Promise.all([
        axios.get('/api/v1/balance'),
        axios.get('/api/v1/balance/transactions', {
          params: {
            limit: rowsPerPage,
            offset: page * rowsPerPage
          }
        })
      ]);
      
      // Добавляем отладочный вывод
      console.log("Balance response:", balanceRes.data);
      console.log("Transactions response:", transactionsRes.data);

      // Проверяем что данные существуют перед установкой
      if (balanceRes.data?.data) {
        setBalance(balanceRes.data.data.balance);
      }
      
      // Проверяем что транзакции существуют
      if (transactionsRes.data?.data) {
        setTransactions(transactionsRes.data.data);
      } else {
        setTransactions([]); // Если данных нет - устанавливаем пустой массив
      }
    } catch (err) {
      console.error('Error fetching data:', err);
      setError(t('balance.errors.loadFailed'));
      setTransactions([]); // При ошибке устанавливаем пустой массив
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [page, rowsPerPage]);

  useEffect(() => {
    // Проверяем параметры URL
    const params = new URLSearchParams(location.search);
    if (params.has('success') && params.get('success') === 'true') {
      setPaymentStatus('success');
      
      // Проверяем, есть ли токен сессии в URL
      const sessionToken = params.get('session_token');
      if (sessionToken && sessionToken !== "{CHECKOUT_SESSION_METADATA.session_token}") {
        // Если есть валидный токен сессии в URL, обновляем аутентификацию
        console.log("Found session token in URL, updating authentication");
        localStorage.setItem('user_session', sessionToken);
      }
      
      // После успешной оплаты обновляем сессию и затем данные
      checkAuth().then(() => {
        fetchData();
      });
      
      // Показываем уведомление
      setSnackbarMessage(t('balance.paymentSuccess'));
      setSnackbarOpen(true);
    } else if (params.has('canceled') && params.get('canceled') === 'true') {
      setPaymentStatus('canceled');
    }
    
    // Очищаем параметры URL
    if (params.has('success') || params.has('canceled')) {
      const cleanUrl = window.location.pathname;
      window.history.replaceState({}, document.title, cleanUrl);
    }
  }, [location, t, checkAuth]);
    
  const formatAmount = (amount) => {
    return new Intl.NumberFormat('sr-RS', {
      style: 'currency',
      currency: 'RSD',
      maximumFractionDigits: 0
    }).format(amount);
  };

  // Функция для перевода описаний транзакций
  const translateDescription = (description) => {
    // Проверка на создание витрины
    if (description === "Создание витрины магазина") {
      return t('balance.descriptions.storeCreation');
    }
    
    // Проверка на перевод объявлений
    const translationMatch = description.match(/Перевод (\d+) объявлений/);
    if (translationMatch) {
      return t('balance.descriptions.listingTranslation', { count: translationMatch[1] });
    }
    
    // Если нет соответствия, возвращаем оригинал
    return description;
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" p={4}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Snackbar
        open={snackbarOpen}
        autoHideDuration={6000}
        onClose={() => setSnackbarOpen(false)}
        message={snackbarMessage}
      />

      {paymentStatus === 'success' && (
        <Alert severity="success" sx={{ mb: 2 }}>
          {t('balance.paymentSuccess')}
        </Alert>
      )}
      
      {paymentStatus === 'canceled' && (
        <Alert severity="info" sx={{ mb: 2 }}>
          {t('balance.paymentCanceled')}
        </Alert>
      )}

      <Typography variant="h4" gutterBottom>
        {t('balance.title')}
      </Typography>

      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <BalanceWidget
        balance={balance}
        onDeposit={() => setDepositDialogOpen(true)}
      />

      <Paper sx={{ mt: 4 }}>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>{t('balance.date')}</TableCell>
                <TableCell>{t('balance.type')}</TableCell>
                <TableCell>{t('balance.description')}</TableCell>
                <TableCell align="right">{t('balance.amount')}</TableCell>
                <TableCell>{t('balance.status')}</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {transactions.map((transaction) => (
                <TableRow key={transaction.id}>
                  <TableCell>
                    {new Date(transaction.created_at).toLocaleDateString(i18n.language)}
                  </TableCell>
                  <TableCell>
                    {t(`balance.types.${transaction.type}`)}
                  </TableCell>
                  <TableCell>
                    {translateDescription(transaction.description)}
                  </TableCell>
                  <TableCell align="right">
                    {formatAmount(transaction.amount)}
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={t(`balance.statuses.${transaction.status}`)}
                      color={transaction.status === 'completed' ? 'success' : 'default'}
                      size="small"
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
        <TablePagination
          component="div"
          count={-1}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={(e, newPage) => setPage(newPage)}
          onRowsPerPageChange={(e) => {
            setRowsPerPage(parseInt(e.target.value, 10));
            setPage(0);
          }}
        />
      </Paper>
      <DepositDialog
        open={depositDialogOpen}
        onClose={() => setDepositDialogOpen(false)}
        onBalanceUpdate={handleBalanceUpdate}
      />
    </Container>
  );
};

export default TransactionsPage;