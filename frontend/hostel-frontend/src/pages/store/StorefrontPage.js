import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  CardActions,
  Box,
  CircularProgress,
  Alert,
  TextField,
  Modal,
  Paper,
  Stack,
  Divider
} from '@mui/material';
import { Plus, Store, Upload, Database, AlertTriangle } from 'lucide-react';
import axios from '../../api/axios';
import { useAuth } from '../../contexts/AuthContext';

const StorefrontPage = () => {
  const { t } = useTranslation(['common', 'marketplace']);
  const navigate = useNavigate();
  const { user } = useAuth();
  
  const [loading, setLoading] = useState(true);
  const [storefronts, setStorefronts] = useState([]);
  const [error, setError] = useState(null);
  const [balance, setBalance] = useState(0);
  const [openCreateModal, setOpenCreateModal] = useState(false);
  const [createForm, setCreateForm] = useState({ name: '', description: '' });
  const [creationError, setCreationError] = useState(null);
  const [creationLoading, setCreationLoading] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [storefrontsResponse, balanceResponse] = await Promise.all([
          axios.get('/api/v1/storefronts'),
          axios.get('/api/v1/balance')
        ]);
        
        setStorefronts(storefrontsResponse.data.data || []);
        setBalance(balanceResponse.data.data?.balance || 0);
      } catch (err) {
        console.error('Error fetching data:', err);
        setError('Не удалось загрузить данные');
      } finally {
        setLoading(false);
      }
    };
    
    fetchData();
  }, []);

  const handleCreateStorefront = async () => {
    try {
      setCreationLoading(true);
      setCreationError(null);
      
      if (!createForm.name.trim()) {
        setCreationError('Укажите название магазина');
        return;
      }
      
      const response = await axios.post('/api/v1/storefronts', createForm);
      setStorefronts(prev => [...prev, response.data.data]);
      setOpenCreateModal(false);
      setCreateForm({ name: '', description: '' });
      
      // Обновляем баланс
      const balanceResponse = await axios.get('/api/v1/balance');
      setBalance(balanceResponse.data.data?.balance || 0);
    } catch (err) {
      console.error('Error creating storefront:', err);
      if (err.response?.status === 402) {
        setCreationError('Недостаточно средств для создания витрины. Требуется 15000 RSD.');
      } else {
        setCreationError('Не удалось создать витрину');
      }
    } finally {
      setCreationLoading(false);
    }
  };

  const navigateToStorefront = (id) => {
    navigate(`/storefronts/${id}`);
  };

  if (loading) {
    return (
      <Container maxWidth="lg" sx={{ py: 4 }}>
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="50vh">
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={4}>
        <Typography variant="h4" component="h1">
          Мои витрины
        </Typography>
        <Button
          variant="contained"
          startIcon={<Plus />}
          onClick={() => setOpenCreateModal(true)}
          disabled={balance < 15000}
        >
          Создать витрину
        </Button>
      </Box>

      {balance < 15000 && (
        <Alert severity="warning" sx={{ mb: 3 }}>
          Для создания витрины необходимо иметь на балансе не менее 15000 RSD. 
          Ваш текущий баланс: {balance.toLocaleString()} RSD.
          <Button 
            variant="outlined" 
            size="small" 
            sx={{ ml: 2 }}
            onClick={() => navigate('/balance')}
          >
            Пополнить баланс
          </Button>
        </Alert>
      )}

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {storefronts.length === 0 ? (
        <Card sx={{ p: 4, textAlign: 'center' }}>
          <Store size={64} stroke={1} style={{ margin: '20px auto', opacity: 0.5 }} />
          <Typography variant="h6" gutterBottom>
            У вас пока нет витрин
          </Typography>
          <Typography variant="body1" color="text.secondary" paragraph>
            Создайте витрину, чтобы начать продавать товары в своем магазине
          </Typography>
          <Button
            variant="contained"
            startIcon={<Plus />}
            onClick={() => setOpenCreateModal(true)}
            disabled={balance < 15000}
          >
            Создать витрину
          </Button>
        </Card>
      ) : (
        <Grid container spacing={3}>
          {storefronts.map((storefront) => (
            <Grid item xs={12} md={6} lg={4} key={storefront.id}>
              <Card sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
                <CardContent sx={{ flexGrow: 1 }}>
                  <Typography variant="h5" component="h2" gutterBottom>
                    {storefront.name}
                  </Typography>
                  <Typography variant="body2" color="text.secondary" paragraph>
                    {storefront.description || 'Нет описания'}
                  </Typography>
                  <Box display="flex" justifyContent="space-between" alignItems="center">
                    <Typography variant="caption" color="text.secondary">
                      Создана: {new Date(storefront.created_at).toLocaleDateString()}
                    </Typography>
                    <Typography
                      variant="caption"
                      color={storefront.status === 'active' ? 'success.main' : 'error.main'}
                    >
                      {storefront.status === 'active' ? 'Активна' : 'Не активна'}
                    </Typography>
                  </Box>
                </CardContent>
                <Divider />
                <CardActions>
                  <Button size="small" onClick={() => navigateToStorefront(storefront.id)}>
                    Управление
                  </Button>
                  <Button size="small" startIcon={<Upload />}>
                    Импорт
                  </Button>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}

      {/* Модальное окно создания витрины */}
      <Modal
        open={openCreateModal}
        onClose={() => setOpenCreateModal(false)}
        aria-labelledby="create-storefront-modal"
      >
        <Paper
          sx={{
            position: 'absolute',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
            width: { xs: '90%', sm: 500 },
            p: 4,
            maxHeight: '90vh',
            overflow: 'auto'
          }}
        >
          <Typography variant="h5" component="h2" gutterBottom>
            Создание витрины
          </Typography>
          
          <Typography variant="body2" color="text.secondary" paragraph>
            Стоимость создания витрины составляет 15000 RSD. Эта сумма будет списана с вашего баланса.
          </Typography>
          
          <Box mb={2}>
            <Typography variant="body2" fontWeight="bold">
              Ваш баланс: {balance.toLocaleString()} RSD
            </Typography>
          </Box>
          
          {creationError && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {creationError}
            </Alert>
          )}
          
          <Stack spacing={2} mb={3}>
            <TextField
              label="Название магазина"
              fullWidth
              required
              value={createForm.name}
              onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
              disabled={creationLoading}
            />
            
            <TextField
              label="Описание"
              fullWidth
              multiline
              rows={3}
              value={createForm.description}
              onChange={(e) => setCreateForm({ ...createForm, description: e.target.value })}
              disabled={creationLoading}
            />
          </Stack>
          
          <Box display="flex" justifyContent="space-between">
            <Button
              variant="outlined"
              onClick={() => setOpenCreateModal(false)}
              disabled={creationLoading}
            >
              Отмена
            </Button>
            
            <Button
              variant="contained"
              onClick={handleCreateStorefront}
              disabled={creationLoading || !createForm.name.trim() || balance < 15000}
              startIcon={creationLoading ? <CircularProgress size={20} /> : <Store />}
            >
              {creationLoading ? 'Создание...' : 'Создать витрину'}
            </Button>
          </Box>
          
          {balance < 15000 && (
            <Alert severity="warning" sx={{ mt: 2 }} icon={<AlertTriangle />}>
              Недостаточно средств для создания витрины. Пополните баланс.
            </Alert>
          )}
        </Paper>
      </Modal>
    </Container>
  );
};

export default StorefrontPage;