// frontend/hostel-frontend/src/pages/store/StorefrontPage.tsx
import React, { useState, useEffect, ChangeEvent } from 'react';
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
import LocationPicker from "../../components/global/LocationPicker";

interface Storefront {
  id: number;
  name: string;
  description: string;
  status: string;
  created_at: string;
  [key: string]: any;
}

interface CreateForm {
  name: string;
  description: string;
  phone: string;
  email: string;
  website: string;
  address: string;
  city: string;
  country: string;
  latitude?: number | null;
  longitude?: number | null;
  [key: string]: any;
}

interface LocationData {
  latitude: number | null;
  longitude: number | null;
  formatted_address: string;
  address_components?: {
    city: string;
    country: string;
  };
}

const StorefrontPage: React.FC = () => {
  const { t } = useTranslation(['common', 'marketplace']);
  const navigate = useNavigate();
  const { user } = useAuth();

  const [loading, setLoading] = useState<boolean>(true);
  const [storefronts, setStorefronts] = useState<Storefront[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [balance, setBalance] = useState<number>(0);
  const [openCreateModal, setOpenCreateModal] = useState<boolean>(false);
  const [createForm, setCreateForm] = useState<CreateForm>({
    name: '',
    description: '',
    phone: '',
    email: '',
    website: '',
    address: '',
    city: '',
    country: ''
  });
  const [creationError, setCreationError] = useState<string | null>(null);
  const [creationLoading, setCreationLoading] = useState<boolean>(false);

  useEffect(() => {
    const fetchData = async (): Promise<void> => {
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
        setError(t('marketplace:store.errors.loadFailed'));
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [t]);

  const handleCreateStorefront = async (): Promise<void> => {
    try {
      setCreationLoading(true);
      setCreationError(null);

      if (!createForm.name.trim()) {
        setCreationError(t('marketplace:store.create.nameRequired'));
        return;
      }

      const response = await axios.post('/api/v1/storefronts', createForm);
      setStorefronts(prev => [...prev, response.data.data]);
      setOpenCreateModal(false);
      setCreateForm({ 
        name: '', 
        description: '', 
        phone: '',
        email: '',
        website: '',
        address: '',
        city: '',
        country: '' 
      });

      // Обновляем баланс
      const balanceResponse = await axios.get('/api/v1/balance');
      setBalance(balanceResponse.data.data?.balance || 0);
    } catch (err: any) {
      console.error('Error creating storefront:', err);
      if (err.response?.status === 402) {
        setCreationError(t('marketplace:store.create.insufficientFunds'));
      } else {
        setCreationError(t('marketplace:store.create.error'));
      }
    } finally {
      setCreationLoading(false);
    }
  };

  const navigateToStorefront = (id: number): void => {
    navigate(`/storefronts/${id}`);
  };

  const handleLocationSelect = (location: LocationData): void => {
    setCreateForm(prev => ({
      ...prev,
      address: location.formatted_address,
      city: location.address_components?.city || '',
      country: location.address_components?.country || '',
      latitude: location.latitude,
      longitude: location.longitude
    }));
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
          {t('marketplace:store.myStorefronts')}
        </Typography>
        <Button
          variant="contained"
          startIcon={<Plus />}
          onClick={() => setOpenCreateModal(true)}
          disabled={balance < 15000}
        >
          {t('marketplace:store.create.button')}
        </Button>
      </Box>

      {balance < 15000 && (
        <Alert severity="warning" sx={{ mb: 3 }}>
          {t('marketplace:store.create.balanceWarning', { balance: balance.toLocaleString() })}
          <Button
            variant="outlined"
            size="small"
            sx={{ ml: 2 }}
            onClick={() => navigate('/balance')}
          >
            {t('common:balance.deposit')}
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
          <Store size={64} stroke="1" style={{ margin: '20px auto', opacity: 0.5 }} />
          <Typography variant="h6" gutterBottom>
            {t('marketplace:store.noStorefronts')}
          </Typography>
          <Typography variant="body1" color="text.secondary" paragraph>
            {t('marketplace:store.createFirstStorefront')}
          </Typography>
          <Button
            variant="contained"
            startIcon={<Plus />}
            onClick={() => setOpenCreateModal(true)}
            disabled={balance < 15000}
          >
            {t('marketplace:store.create.button')}
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
                    {storefront.description || t('marketplace:store.noDescription')}
                  </Typography>
                  <Box display="flex" justifyContent="space-between" alignItems="center">
                    <Typography variant="caption" color="text.secondary">
                      {t('marketplace:store.created')}: {new Date(storefront.created_at).toLocaleDateString()}
                    </Typography>
                    <Typography
                      variant="caption"
                      color={storefront.status === 'active' ? 'success.main' : 'error.main'}
                    >
                      {storefront.status === 'active'
                        ? t('marketplace:store.statusActive')
                        : t('marketplace:store.statusInactive')}
                    </Typography>
                  </Box>
                </CardContent>
                <Divider />
                <CardActions>
                  <Button size="small" onClick={() => navigateToStorefront(storefront.id)}>
                    {t('marketplace:store.manage')}
                  </Button>
                  <Button size="small" startIcon={<Upload />}>
                    {t('marketplace:store.import.button')}
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
        keepMounted
        children={
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
            {t('marketplace:store.create.title')}
          </Typography>

          <Typography variant="body2" color="text.secondary" paragraph>
            {t('marketplace:store.create.costNote')}
          </Typography>

          <Box mb={2}>
            <Typography variant="body2" fontWeight="bold">
              {t('marketplace:store.create.yourBalance')}: {balance.toLocaleString()} RSD
            </Typography>
          </Box>

          {creationError && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {creationError}
            </Alert>
          )}

          <Stack spacing={2} mb={3}>
            {/* Существующие поля name и description */}
            <TextField
              label={t('marketplace:store.settings.name')}
              fullWidth
              required
              value={createForm.name}
              onChange={(e: ChangeEvent<HTMLInputElement>) => 
                setCreateForm({ ...createForm, name: e.target.value })}
              disabled={creationLoading}
            />

            <TextField
              label={t('common:common.description')}
              fullWidth
              multiline
              rows={3}
              value={createForm.description}
              onChange={(e: ChangeEvent<HTMLInputElement>) => 
                setCreateForm({ ...createForm, description: e.target.value })}
              disabled={creationLoading}
            />
            <Box mb={3}>
              <Typography variant="subtitle1" gutterBottom>
                {t('marketplace:store.settings.location')}
              </Typography>

              <LocationPicker
                onLocationSelect={handleLocationSelect}
                initialLocation={{
                  latitude: createForm.latitude || null,
                  longitude: createForm.longitude || null,
                  formatted_address: createForm.address,
                  address_components: {
                    city: createForm.city,
                    country: createForm.country
                  }
                }}
              />
            </Box>
            <TextField
              label={t('marketplace:store.settings.phone')}
              fullWidth
              value={createForm.phone}
              onChange={(e: ChangeEvent<HTMLInputElement>) => 
                setCreateForm({ ...createForm, phone: e.target.value })}
              disabled={creationLoading}
            />

            <TextField
              label={t('marketplace:store.settings.email')}
              fullWidth
              value={createForm.email}
              onChange={(e: ChangeEvent<HTMLInputElement>) => 
                setCreateForm({ ...createForm, email: e.target.value })}
              disabled={creationLoading}
            />

            <TextField
              label={t('marketplace:store.settings.website')}
              fullWidth
              value={createForm.website}
              onChange={(e: ChangeEvent<HTMLInputElement>) => 
                setCreateForm({ ...createForm, website: e.target.value })}
              disabled={creationLoading}
            />
          </Stack>

          <Box display="flex" justifyContent="space-between">
            <Button
              variant="outlined"
              onClick={() => setOpenCreateModal(false)}
              disabled={creationLoading}
            >
              {t('common:buttons.cancel')}
            </Button>

            <Button
              variant="contained"
              onClick={handleCreateStorefront}
              disabled={creationLoading || !createForm.name.trim() || balance < 15000}
              startIcon={creationLoading ? <CircularProgress size={20} /> : <Store />}
            >
              {creationLoading ? t('marketplace:store.create.creating') : t('marketplace:store.create.button')}
            </Button>
          </Box>

          {balance < 15000 && (
            <Alert severity="warning" sx={{ mt: 2 }} icon={<AlertTriangle />}>
              {t('marketplace:store.create.insufficientFunds')}
            </Alert>
          )}
        </Paper>
      }
      />
    </Container>
  );
};

export default StorefrontPage;