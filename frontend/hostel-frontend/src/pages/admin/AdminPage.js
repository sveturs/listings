// frontend/hostel-frontend/src/pages/admin/AdminPage.js
import React, { useState } from 'react';
import { Box, Button, Paper, Typography, Alert, CircularProgress, Grid } from '@mui/material';
import { People as PeopleIcon } from '@mui/icons-material';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import axios from '../../api/axios';
import { useAuth } from '../../contexts/AuthContext';
import { getAdminEmails } from '../../utils/adminUtils';

const AdminPage = () => {
  const { t } = useTranslation(['common', 'marketplace']);
  const { user } = useAuth();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [syncLoading, setSyncLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [syncResult, setSyncResult] = useState(null);
  const [error, setError] = useState(null);
  const [syncError, setSyncError] = useState(null);
  const [ratingsLoading, setRatingsLoading] = useState(false);
  const [ratingsResult, setRatingsResult] = useState(null);
  const [ratingsError, setRatingsError] = useState(null);

  const handleReindexRatings = async () => {
    setRatingsLoading(true);
    setRatingsResult(null);
    setRatingsError(null);

    try {
      const response = await axios.post('/api/v1/admin/reindex-ratings');
      setRatingsResult(response.data);
    } catch (err) {
      console.error('Ошибка при обновлении рейтингов:', err);
      setRatingsError(err.response?.data?.message || err.message || 'Произошла ошибка при обновлении рейтингов');
    } finally {
      setRatingsLoading(false);
    }
  };
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
          Привет! {user?.name || 'Администратор'}.
        </Typography>

        <Alert severity="info" sx={{ mb: 3 }}>
          Эта страница доступна только для администраторов! ({getAdminEmails().join(', ')})
        </Alert>

        <Box sx={{ mt: 4, mb: 4 }}>
          <Typography variant="h5" gutterBottom>
            Управление системой
          </Typography>
          <Grid container spacing={2} sx={{ mt: 2 }}>
            <Grid item xs={12} sm={6} md={4} lg={3}>
              <Paper
                elevation={3}
                sx={{
                  p: 3,
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  height: '100%',
                  cursor: 'pointer',
                  transition: 'all 0.3s',
                  '&:hover': { transform: 'translateY(-5px)', boxShadow: 6 }
                }}
                onClick={() => navigate('/admin/users')}
              >
                <PeopleIcon sx={{ fontSize: 48, color: 'primary.main', mb: 2 }} />
                <Typography variant="h6" align="center" gutterBottom>
                  Пользователи
                </Typography>
                <Typography variant="body2" align="center" color="text.secondary">
                  Управление пользователями, блокировка, редактирование профилей
                </Typography>
              </Paper>
            </Grid>
          </Grid>
        </Box>

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
      <Box sx={{ mt: 4 }}>
  <Typography variant="h5" gutterBottom>
    Переиндексация рейтингов
  </Typography>
  <Typography variant="body2" paragraph color="text.secondary">
    Нажмите кнопку ниже, чтобы запустить обновление рейтингов объявлений.
    Это обеспечит корректную сортировку объявлений по рейтингу.
  </Typography>

  <Button
    variant="contained"
    color="primary"
    onClick={handleReindexRatings} // Добавим новую функцию
    disabled={ratingsLoading}
    startIcon={ratingsLoading && <CircularProgress size={20} color="inherit" />}
    sx={{ mt: 2 }}
  >
    {ratingsLoading ? 'Обновление рейтингов...' : 'Обновить рейтинги'}
  </Button>

  {ratingsResult && (
    <Alert severity="success" sx={{ mt: 2 }}>
      Обновление рейтингов успешно запущено! {ratingsResult.message}
    </Alert>
  )}

  {ratingsError && (
    <Alert severity="error" sx={{ mt: 2 }}>
      {ratingsError}
    </Alert>
  )}
</Box>
        </Typography>
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