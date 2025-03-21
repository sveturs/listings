// frontend/hostel-frontend/src/pages/admin/AdminPage.js
import React, { useState } from 'react';
import { Box, Button, Paper, Typography, Alert, CircularProgress } from '@mui/material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';
import { useAuth } from '../../contexts/AuthContext';

const AdminPage = () => {
  const { t } = useTranslation(['common', 'marketplace']);
  const { user } = useAuth();
  const [loading, setLoading] = useState(false);
  const [syncLoading, setSyncLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [syncResult, setSyncResult] = useState(null);
  const [error, setError] = useState(null);
  const [syncError, setSyncError] = useState(null);

  const handleReindex = async () => {
    setLoading(true);
    setResult(null);
    setError(null);

    try {
      const response = await axios.post('/api/v1/admin/reindex-listings');
      setResult(response.data);
    } catch (err) {
      console.error('Ошибка при переиндексации:', err);
      setError(err.response?.data?.message || err.message || 'Произошла ошибка при переиндексации');
    } finally {
      setLoading(false);
    }
  };

  const handleSyncDiscounts = async () => {
    setSyncLoading(true);
    setSyncResult(null);
    setSyncError(null);

    try {
      const response = await axios.post('/api/v1/admin/sync-discounts');
      setSyncResult(response.data);
    } catch (err) {
      console.error('Ошибка при синхронизации скидок:', err);
      setSyncError(err.response?.data?.message || err.message || 'Произошла ошибка при синхронизации скидок');
    } finally {
      setSyncLoading(false);
    }
  };

  return (
    <Box sx={{ py: 4 }}>
      <Paper sx={{ p: 3, mb: 4 }}>
        <Typography variant="h4" gutterBottom>
          Панель администратора
        </Typography>
        <Typography variant="body1" paragraph>
          Добро пожаловать, {user?.name || 'Администратор'}!
        </Typography>
        
        <Alert severity="info" sx={{ mb: 3 }}>
          Эта страница доступна только для администратора!  (voroshilovdo@gmail.com)
        </Alert>

        <Box sx={{ mt: 4 }}>
          <Typography variant="h5" gutterBottom>
            Управление поисковым индексом
          </Typography>
          <Typography variant="body2" paragraph color="text.secondary">
            Нажмите кнопку ниже, чтобы запустить полную переиндексацию всех объявлений в OpenSearch.
            Это может занять некоторое время, особенно если в системе много объявлений.
          </Typography>

          <Button 
            variant="contained" 
            color="primary" 
            onClick={handleReindex}
            disabled={loading}
            startIcon={loading && <CircularProgress size={20} color="inherit" />}
            sx={{ mt: 2, mr: 2 }}
          >
            {loading ? 'Выполняется переиндексация...' : 'Запустить переиндексацию'}
          </Button>

          {result && (
            <Alert severity="success" sx={{ mt: 2 }}>
              Переиндексация успешно запущена! {result.message}
            </Alert>
          )}

          {error && (
            <Alert severity="error" sx={{ mt: 2 }}>
              {error}
            </Alert>
          )}
        </Box>

        <Box sx={{ mt: 4 }}>
          <Typography variant="h5" gutterBottom>
            Синхронизация данных о скидках
          </Typography>
          <Typography variant="body2" paragraph color="text.secondary">
            Нажмите кнопку ниже, чтобы запустить синхронизацию данных о скидках.
            Это обновит информацию о скидках для всех объявлений и обеспечит корректное отображение скидок.
          </Typography>

          <Button 
            variant="contained" 
            color="secondary" 
            onClick={handleSyncDiscounts}
            disabled={syncLoading}
            startIcon={syncLoading && <CircularProgress size={20} color="inherit" />}
            sx={{ mt: 2 }}
          >
            {syncLoading ? 'Выполняется синхронизация...' : 'Синхронизировать скидки'}
          </Button>

          {syncResult && (
            <Alert severity="success" sx={{ mt: 2 }}>
              Синхронизация скидок успешно запущена! {syncResult.message}
            </Alert>
          )}

          {syncError && (
            <Alert severity="error" sx={{ mt: 2 }}>
              {syncError}
            </Alert>
          )}
        </Box>
      </Paper>

      <Paper sx={{ p: 3 }}>
        <Typography variant="h5" gutterBottom>
          Информация о системе
        </Typography>
        <Typography variant="body2">
          • Версия API: 1.0.0<br />
          • Текущая роль пользователя: {user?.is_admin ? 'Администратор' : 'Стандартный пользователь'}<br />
          • ID пользователя: {user?.id || 'Не авторизован'}<br />
          • Индекс OpenSearch: marketplace<br />
        </Typography>
      </Paper>
    </Box>
  );
};

export default AdminPage;