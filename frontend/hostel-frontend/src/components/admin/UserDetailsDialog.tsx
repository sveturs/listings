// frontend/hostel-frontend/src/components/admin/UserDetailsDialog.tsx
import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Grid,
  Box,
  Avatar,
  Divider,
  Chip,
  CircularProgress,
  Alert,
  Paper,
  Tabs,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow
} from '@mui/material';
import { useTranslation } from 'react-i18next';
import axios from '../../api/axios';
import { formatDate, formatRelativeDate } from '../../utils/dateUtils';

// Интерфейсы для типизации данных
interface TabPanelProps {
  children?: React.ReactNode;
  value: number;
  index: number;
  [key: string]: any;
}

interface UserDetailsDialogProps {
  open: boolean;
  onClose: () => void;
  userId: number | string;
}

interface UserData {
  id: number | string;
  name: string;
  email: string;
  phone?: string;
  city?: string;
  country?: string;
  picture_url?: string;
  account_status?: string;
  created_at: string;
  last_seen?: string;
  google_id?: string;
  timezone?: string;
  notification_email?: boolean;
  bio?: string;
  settings?: UserSettings | string;
  [key: string]: any;
}

interface UserSettings {
  [key: string]: any;
}

interface Review {
  id: number | string;
  rating: number;
  comment?: string;
  entity_type: string;
  entity_id: number | string;
  created_at: string;
  [key: string]: any;
}

interface Listing {
  id: number | string;
  title: string;
  price: number;
  currency?: string;
  status: string;
  views_count?: number;
  created_at: string;
  [key: string]: any;
}

interface Balance {
  balance: number;
  frozen_balance: number;
  currency: string;
  updated_at: string;
  [key: string]: any;
}

interface Transaction {
  id: number | string;
  type: string;
  amount: number;
  currency: string;
  status: string;
  payment_method?: string;
  created_at: string;
  [key: string]: any;
}

interface LoadingData {
  reviews: boolean;
  listings: boolean;
  balance: boolean;
  transactions: boolean;
}

// Компонент для отображения панелей с разной информацией
function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
      <div
          role="tabpanel"
          hidden={value !== index}
          id={`user-tabpanel-${index}`}
          aria-labelledby={`user-tab-${index}`}
          {...other}
      >
        {value === index && (
            <Box sx={{ p: 3 }}>
              {children}
            </Box>
        )}
      </div>
  );
}

const UserDetailsDialog: React.FC<UserDetailsDialogProps> = ({ open, onClose, userId }) => {
  const { t } = useTranslation(['common', 'admin']);
  const [userData, setUserData] = useState<UserData | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [tabValue, setTabValue] = useState<number>(0);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [listings, setListings] = useState<Listing[]>([]);
  const [balance, setBalance] = useState<Balance | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loadingData, setLoadingData] = useState<LoadingData>({
    reviews: false,
    listings: false,
    balance: false,
    transactions: false
  });

  // Загрузка данных пользователя
  useEffect(() => {
    if (open && userId) {
      loadUserDetails();
    }
  }, [open, userId]);

  // Загрузка дополнительных данных при переключении вкладок
  useEffect(() => {
    if (open && userId) {
      if (tabValue === 1) {
        loadUserReviews();
      } else if (tabValue === 2) {
        loadUserListings();
      } else if (tabValue === 3) {
        loadUserBalance();
        loadUserTransactions();
      }
    }
  }, [tabValue, open, userId]);

  const loadUserDetails = async (): Promise<void> => {
    setLoading(true);
    setError(null);

    try {
      const response = await axios.get(`/api/v1/admin/users/${userId}`);
      setUserData(response.data.data);
    } catch (err: any) {
      console.error('Ошибка при загрузке данных пользователя:', err);
      setError(err.response?.data?.message || err.message || 'Не удалось загрузить данные пользователя');
    } finally {
      setLoading(false);
    }
  };

  const loadUserReviews = async (): Promise<void> => {
    setLoadingData(prev => ({ ...prev, reviews: true }));

    try {
      const response = await axios.get(`/api/v1/users/${userId}/reviews`);
      setReviews(response.data.data || []);
    } catch (err) {
      console.error('Ошибка при загрузке отзывов пользователя:', err);
    } finally {
      setLoadingData(prev => ({ ...prev, reviews: false }));
    }
  };

  const loadUserListings = async (): Promise<void> => {
    setLoadingData(prev => ({ ...prev, listings: true }));

    try {
      // Используем фильтр по user_id для получения объявлений пользователя
      const response = await axios.get(`/api/v1/marketplace/listings?user_id=${userId}`);
      console.log('Получены объявления пользователя:', response.data);
      console.log('Структура данных:', JSON.stringify(response.data));
      
      // Извлекаем данные из структуры ответа
      let listingsData: Listing[] = [];
      if (response.data && response.data.data && Array.isArray(response.data.data.data)) {
        // Если структура: response.data.data.data[]
        listingsData = response.data.data.data;
        console.log('Найдены данные в response.data.data.data:', listingsData.length, 'объявлений');
      } else if (response.data && Array.isArray(response.data.data)) {
        // Если структура: response.data.data[]
        listingsData = response.data.data;
        console.log('Найдены данные в response.data.data:', listingsData.length, 'объявлений');
      } else if (response.data && Array.isArray(response.data)) {
        // Если структура: response.data[]
        listingsData = response.data;
        console.log('Найдены данные в response.data:', listingsData.length, 'объявлений');
      }
      
      console.log('Финальные данные объявлений:', listingsData);
      setListings(listingsData);
    } catch (err) {
      console.error('Ошибка при загрузке объявлений пользователя:', err);
    } finally {
      setLoadingData(prev => ({ ...prev, listings: false }));
    }
  };

  const loadUserBalance = async (): Promise<void> => {
    setLoadingData(prev => ({ ...prev, balance: true }));

    try {
      // Использование admin эндпоинта для получения баланса
      const response = await axios.get(`/api/v1/admin/users/${userId}/balance`);
      setBalance(response.data.data || null);
    } catch (err) {
      console.error('Ошибка при загрузке баланса пользователя:', err);
    } finally {
      setLoadingData(prev => ({ ...prev, balance: false }));
    }
  };

  const loadUserTransactions = async (): Promise<void> => {
    setLoadingData(prev => ({ ...prev, transactions: true }));

    try {
      // Использование admin эндпоинта для получения транзакций
      const response = await axios.get(`/api/v1/admin/users/${userId}/transactions`);
      setTransactions(response.data.data || []);
    } catch (err) {
      console.error('Ошибка при загрузке транзакций пользователя:', err);
    } finally {
      setLoadingData(prev => ({ ...prev, transactions: false }));
    }
  };

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number): void => {
    setTabValue(newValue);
  };

  // Функция для получения цвета статуса
  const getStatusColor = (status: string): "success" | "error" | "warning" | "default" => {
    switch (status) {
      case 'active':
        return 'success';
      case 'blocked':
        return 'error';
      case 'pending':
        return 'warning';
      default:
        return 'default';
    }
  };

  // Безопасный рендеринг JSON данных
  const renderJsonData = (data: any): string => {
    try {
      // Если данные уже являются объектом, просто преобразуем в строку
      // Если строка - попробуем распарсить
      const parsedData = typeof data === 'string'
          ? JSON.parse(data)
          : data;

      return JSON.stringify(parsedData, null, 2);
    } catch (err) {
      console.error("Ошибка при парсинге JSON:", err);
      return "Ошибка при отображении данных";
    }
  };

  // Отображаем лоадер при загрузке
  if (loading) {
    return (
        <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
          <DialogContent>
            <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
              <CircularProgress />
            </Box>
          </DialogContent>
          <DialogActions>
            <Button onClick={onClose}>Закрыть</Button>
          </DialogActions>
        </Dialog>
    );
  }

  return (
      <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
        <DialogTitle>
          Информация о пользователе
          {userData && (
              <Typography variant="subtitle1" color="text.secondary">
                ID: {userData.id}
              </Typography>
          )}
        </DialogTitle>

        <DialogContent dividers>
          {error && (
              <Alert severity="error" sx={{ mb: 3 }}>
                {error}
              </Alert>
          )}

          {userData && (
              <>
                {/* Верхняя секция с основной информацией */}
                <Box sx={{ mb: 3 }}>
                  <Grid container spacing={3}>
                    <Grid item xs={12} sm={4} md={3}>
                      <Box display="flex" flexDirection="column" alignItems="center">
                        <Avatar
                            src={userData.picture_url}
                            alt={userData.name}
                            sx={{ width: 120, height: 120, mb: 2 }}
                        />
                        <Chip
                            label={userData.account_status || 'active'}
                            color={getStatusColor(userData.account_status || 'active')}
                            sx={{ mb: 1 }}
                        />
                        <Typography variant="body2" align="center">
                          Зарегистрирован: {formatDate(userData.created_at)}
                        </Typography>
                        {userData.last_seen && (
                            <Typography variant="body2" align="center">
                              Последний вход: {formatRelativeDate(userData.last_seen)}
                            </Typography>
                        )}
                      </Box>
                    </Grid>

                    <Grid item xs={12} sm={8} md={9}>
                      <Typography variant="h5" gutterBottom>{userData.name}</Typography>

                      <Box sx={{ mb: 2 }}>
                        <Typography variant="subtitle1" gutterBottom>Контактная информация:</Typography>
                        <Grid container spacing={1}>
                          <Grid item xs={4} sm={3}>
                            <Typography variant="body2" color="text.secondary">Email:</Typography>
                          </Grid>
                          <Grid item xs={8} sm={9}>
                            <Typography variant="body2">{userData.email}</Typography>
                          </Grid>

                          <Grid item xs={4} sm={3}>
                            <Typography variant="body2" color="text.secondary">Телефон:</Typography>
                          </Grid>
                          <Grid item xs={8} sm={9}>
                            <Typography variant="body2">{userData.phone || '-'}</Typography>
                          </Grid>

                          <Grid item xs={4} sm={3}>
                            <Typography variant="body2" color="text.secondary">Город:</Typography>
                          </Grid>
                          <Grid item xs={8} sm={9}>
                            <Typography variant="body2">{userData.city || '-'}</Typography>
                          </Grid>

                          <Grid item xs={4} sm={3}>
                            <Typography variant="body2" color="text.secondary">Страна:</Typography>
                          </Grid>
                          <Grid item xs={8} sm={9}>
                            <Typography variant="body2">{userData.country || '-'}</Typography>
                          </Grid>
                        </Grid>
                      </Box>

                      <Divider sx={{ my: 2 }} />

                      <Box>
                        <Typography variant="subtitle1" gutterBottom>Системная информация:</Typography>
                        <Grid container spacing={1}>
                          <Grid item xs={4} sm={3}>
                            <Typography variant="body2" color="text.secondary">Google ID:</Typography>
                          </Grid>
                          <Grid item xs={8} sm={9}>
                            <Typography variant="body2">{userData.google_id || '-'}</Typography>
                          </Grid>

                          <Grid item xs={4} sm={3}>
                            <Typography variant="body2" color="text.secondary">Часовой пояс:</Typography>
                          </Grid>
                          <Grid item xs={8} sm={9}>
                            <Typography variant="body2">{userData.timezone || 'UTC'}</Typography>
                          </Grid>

                          <Grid item xs={4} sm={3}>
                            <Typography variant="body2" color="text.secondary">Уведомления:</Typography>
                          </Grid>
                          <Grid item xs={8} sm={9}>
                            <Typography variant="body2">
                              {userData.notification_email ? 'Включены' : 'Отключены'}
                            </Typography>
                          </Grid>
                        </Grid>
                      </Box>
                    </Grid>
                  </Grid>
                </Box>

                <Divider sx={{ my: 3 }} />

                {/* Вкладки с дополнительной информацией */}
                <Box sx={{ width: '100%' }}>
                  <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                    <Tabs value={tabValue} onChange={handleTabChange} aria-label="user details tabs">
                      <Tab label="Профиль" id="user-tab-0" />
                      <Tab label="Отзывы" id="user-tab-1" />
                      <Tab label="Объявления" id="user-tab-2" />
                      <Tab label="Финансы" id="user-tab-3" />
                    </Tabs>
                  </Box>

                  {/* Вкладка Профиль */}
                  <TabPanel value={tabValue} index={0}>
                    {(() => {
                      return (
                        <>
                          <Box sx={{ mb: 3 }}>
                            <Typography variant="h6" gutterBottom>О пользователе</Typography>
                            <Paper variant="outlined" sx={{ p: 2 }}>
                              <Typography variant="body1">
                                {userData.bio || 'Информация отсутствует'}
                              </Typography>
                            </Paper>
                          </Box>

                          {userData.settings && (
                              <Box sx={{ mb: 3 }}>
                                <Typography variant="h6" gutterBottom>Настройки пользователя</Typography>
                                <Paper variant="outlined" sx={{ p: 2 }}>
                            <pre style={{ whiteSpace: 'pre-wrap', margin: 0 }}>
                              {renderJsonData(userData.settings)}
                            </pre>
                                </Paper>
                              </Box>
                          )}
                        </>
                      );
                    })()}
                  </TabPanel>

                  {/* Вкладка Отзывы */}
                  <TabPanel value={tabValue} index={1}>
                    {(() => {
                      return loadingData.reviews ? (
                          <Box display="flex" justifyContent="center" p={3}>
                            <CircularProgress />
                          </Box>
                      ) : reviews.length > 0 ? (
                          <TableContainer component={Paper}>
                            <Table>
                              <TableHead>
                                <TableRow>
                                  <TableCell>ID</TableCell>
                                  <TableCell>Рейтинг</TableCell>
                                  <TableCell>Комментарий</TableCell>
                                  <TableCell>Объект</TableCell>
                                  <TableCell>Дата</TableCell>
                                </TableRow>
                              </TableHead>
                              <TableBody>
                                {reviews.map((review) => (
                                    <TableRow key={review.id}>
                                      <TableCell>{review.id}</TableCell>
                                      <TableCell>{review.rating}/5</TableCell>
                                      <TableCell>
                                        {review.comment ? review.comment.substring(0, 100) + (review.comment.length > 100 ? '...' : '') : '-'}
                                      </TableCell>
                                      <TableCell>
                                        {review.entity_type} #{review.entity_id}
                                      </TableCell>
                                      <TableCell>{formatDate(review.created_at)}</TableCell>
                                    </TableRow>
                                ))}
                              </TableBody>
                            </Table>
                          </TableContainer>
                      ) : (
                          <Typography>У пользователя нет отзывов</Typography>
                      );
                    })()}
                  </TabPanel>

                  {/* Вкладка Объявления */}
                  <TabPanel value={tabValue} index={2}>
                    {(() => {
                      console.log('Отображение вкладки объявлений. Данные:', listings);
                      console.log('Количество объявлений:', listings ? listings.length : 0);

                      return loadingData.listings ? (
                        <Box display="flex" justifyContent="center" p={3}>
                          <CircularProgress />
                        </Box>
                    ) : listings && listings.length > 0 ? (
                        <TableContainer component={Paper}>
                          <Table>
                            <TableHead>
                              <TableRow>
                                <TableCell>ID</TableCell>
                                <TableCell>Заголовок</TableCell>
                                <TableCell>Цена</TableCell>
                                <TableCell>Статус</TableCell>
                                <TableCell>Прогледалы</TableCell>
                                <TableCell>Дата</TableCell>
                              </TableRow>
                            </TableHead>
                            <TableBody>
                              {listings.map((listing) => {
                                  console.log('Отображение объявления:', listing);
                                  return (
                                    <TableRow key={listing.id}>
                                      <TableCell>{listing.id}</TableCell>
                                      <TableCell>{listing.title}</TableCell>
                                      <TableCell>{listing.price} {listing.currency || 'RSD'}</TableCell>
                                      <TableCell>
                                        <Chip
                                            label={listing.status}
                                            color={listing.status === 'active' ? 'success' : listing.status === 'draft' ? 'warning' : 'default'}
                                            size="small"
                                        />
                                      </TableCell>
                                      <TableCell>{listing.views_count}</TableCell>
                                      <TableCell>{formatDate(listing.created_at)}</TableCell>
                                    </TableRow>
                                  );
                              })}
                            </TableBody>
                          </Table>
                        </TableContainer>
                    ) : (
                        <Typography>У пользователя нет объявлений</Typography>
                    );
                    })()}
                  </TabPanel>

                  {/* Вкладка Финансы */}
                  <TabPanel value={tabValue} index={3}>
                    {(() => {
                      return loadingData.balance || loadingData.transactions ? (
                          <Box display="flex" justifyContent="center" p={3}>
                            <CircularProgress />
                          </Box>
                      ) : (
                          <>
                            {balance ? (
                                <Box sx={{ mb: 4 }}>
                                  <Typography variant="h6" gutterBottom>Баланс</Typography>
                                  <Paper variant="outlined" sx={{ p: 2 }}>
                                    <Grid container spacing={2}>
                                      <Grid item xs={6} md={4}>
                                        <Typography variant="subtitle2" color="text.secondary">Доступно</Typography>
                                        <Typography variant="h5">{balance.balance} {balance.currency}</Typography>
                                      </Grid>
                                      <Grid item xs={6} md={4}>
                                        <Typography variant="subtitle2" color="text.secondary">Заморожено</Typography>
                                        <Typography variant="h5">{balance.frozen_balance} {balance.currency}</Typography>
                                      </Grid>
                                      <Grid item xs={12} md={4}>
                                        <Typography variant="subtitle2" color="text.secondary">Обновлено</Typography>
                                        <Typography variant="body2">{formatDate(balance.updated_at)}</Typography>
                                      </Grid>
                                    </Grid>
                                  </Paper>
                                </Box>
                            ) : (
                                <Alert severity="info" sx={{ mb: 3 }}>
                                  Информация о балансе недоступна
                                </Alert>
                            )}

                            {transactions.length > 0 ? (
                                <Box>
                                  <Typography variant="h6" gutterBottom>История транзакций</Typography>
                                  <TableContainer component={Paper}>
                                    <Table>
                                      <TableHead>
                                        <TableRow>
                                          <TableCell>ID</TableCell>
                                          <TableCell>Тип</TableCell>
                                          <TableCell>Сумма</TableCell>
                                          <TableCell>Статус</TableCell>
                                          <TableCell>Метод</TableCell>
                                          <TableCell>Дата</TableCell>
                                        </TableRow>
                                      </TableHead>
                                      <TableBody>
                                        {transactions.map((tx) => (
                                            <TableRow key={tx.id}>
                                              <TableCell>{tx.id}</TableCell>
                                              <TableCell>{tx.type}</TableCell>
                                              <TableCell>
                                                {tx.amount} {tx.currency}
                                              </TableCell>
                                              <TableCell>
                                                <Chip
                                                    label={tx.status}
                                                    color={
                                                      tx.status === 'completed' ? 'success' :
                                                          tx.status === 'pending' ? 'warning' :
                                                              'error'
                                                    }
                                                    size="small"
                                                />
                                              </TableCell>
                                              <TableCell>{tx.payment_method || '-'}</TableCell>
                                              <TableCell>{formatDate(tx.created_at)}</TableCell>
                                            </TableRow>
                                        ))}
                                      </TableBody>
                                    </Table>
                                  </TableContainer>
                                </Box>
                            ) : (
                                <Typography>История транзакций отсутствует</Typography>
                            )}
                          </>
                      );
                    })()}
                  </TabPanel>
                </Box>
              </>
          )}
        </DialogContent>

        <DialogActions>
          <Button onClick={onClose}>Закрыть</Button>
        </DialogActions>
      </Dialog>
  );
};

export default UserDetailsDialog;