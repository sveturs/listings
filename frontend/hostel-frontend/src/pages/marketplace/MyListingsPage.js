// frontend/hostel-frontend/src/pages/marketplace/MyListingsPage.js
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

import {
    Container,
    Typography,
    Grid,
    Box,
    CircularProgress,
    Alert,
    Button,
    ToggleButtonGroup,
    ToggleButton,
    useTheme,
    useMediaQuery,
    Toolbar,
    Chip,
    Menu,
    MenuItem,
    ListItemIcon,
    ListItemText,
    Divider,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    Snackbar
} from '@mui/material';
import { 
    Plus, 
    Grid as GridIcon, 
    List, 
    Check, 
    Filter, 
    ArrowUp, 
    Zap, 
    AlertTriangle, 
    Trash2, 
    Copy, 
    Star, 
    Percent
} from 'lucide-react';
import Checkbox from '@mui/material/Checkbox';
import { Link } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import ListingCard from '../../components/marketplace/ListingCard';
import MarketplaceListingsList from '../../components/marketplace/MarketplaceListingsList';
import InfiniteScroll from '../../components/marketplace/InfiniteScroll';
import axios from '../../api/axios';
import DepositDialog from '../../components/balance/DepositDialog';

const MyListingsPage = () => {
    const { t } = useTranslation('marketplace');
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const { user } = useAuth();
    
    // Состояния для пагинации и бесконечной прокрутки
    const [page, setPage] = useState(1);
    const [hasMoreListings, setHasMoreListings] = useState(true);
    const [loadingMore, setLoadingMore] = useState(false);
    const [totalListings, setTotalListings] = useState(0);
    
    // Состояния для основных данных
    const [listings, setListings] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    
    // Состояния для сортировки
    const [sortField, setSortField] = useState('created_at');
    const [sortDirection, setSortDirection] = useState('desc');
    
    // Состояние для режима отображения (grid/list)
    const [viewMode, setViewMode] = useState(() => {
        // Загружаем сохраненный режим просмотра из localStorage
        const savedViewMode = localStorage.getItem('my-listings-view-mode');
        return savedViewMode === 'list' ? 'list' : 'grid';
    });

    // Состояния для управления выбором объявлений
    const [selectedListings, setSelectedListings] = useState([]);
    const [statusFilter, setStatusFilter] = useState('active'); // 'active', 'inactive', 'all'
    
    // Состояния для управления диалоговыми окнами
    const [menuAnchorEl, setMenuAnchorEl] = useState(null);
    const [confirmDialogOpen, setConfirmDialogOpen] = useState(false);
    const [confirmDialogAction, setConfirmDialogAction] = useState(null);
    const [depositDialogOpen, setDepositDialogOpen] = useState(false);
    const [promotionDialogOpen, setPromotionDialogOpen] = useState(false);
    const [snackbarOpen, setSnackbarOpen] = useState(false);
    const [snackbarMessage, setSnackbarMessage] = useState('');

    // Состояние для выбранного платного действия
    const [selectedPromotionType, setSelectedPromotionType] = useState(null);
    const [balance, setBalance] = useState(0);

    // Функция для загрузки начальных данных с учетом сортировки и фильтрации
    const fetchInitialData = async (sortParams = {}) => {
        try {
            setLoading(true);
            setError(null);
            
            // Формируем параметры API сортировки
            let apiSort = '';
            const currentSortField = sortParams.field || sortField;
            const currentSortDirection = sortParams.direction || sortDirection;
            
            switch (currentSortField) {
                case 'created_at':
                    apiSort = `date_${currentSortDirection}`;
                    break;
                case 'title':
                    apiSort = `title_${currentSortDirection}`;
                    break;
                case 'price':
                    apiSort = `price_${currentSortDirection}`;
                    break;
                case 'reviews':
                    apiSort = `rating_${currentSortDirection}`;
                    break;
                default:
                    apiSort = `date_${currentSortDirection}`;
            }
            
            console.log(`Запрос моих объявлений с сортировкой: ${apiSort} и фильтром статуса: ${statusFilter}`);
            
            const params = { 
                page: 1,
                size: 20,
                sort_by: apiSort,
                user_id: user?.id
            };
            
            // Добавляем фильтр статуса, если выбран не "all"
            if (statusFilter !== 'all') {
                params.status = statusFilter;
            }
            
            const response = await axios.get('/api/v1/marketplace/listings', {
                params,
                withCredentials: true
            });
            
            // Параллельно получаем баланс пользователя
            try {
                const balanceResponse = await axios.get('/api/v1/balance');
                if (balanceResponse.data?.data?.balance) {
                    setBalance(balanceResponse.data.data.balance);
                }
            } catch (err) {
                console.error('Ошибка получения баланса:', err);
            }
            
            // Обработка ответа
            if (response.data?.data?.data && Array.isArray(response.data.data.data)) {
                // Фильтр объявлений по ID пользователя
                const userListings = response.data.data.data.filter(listing =>
                    String(listing.user_id) === String(user?.id)
                );
                setListings(userListings);
                
                // Обновляем информацию о пагинации
                if (response.data.data.meta) {
                    setTotalListings(response.data.data.meta.total || 0);
                    setHasMoreListings(userListings.length < response.data.data.meta.total);
                } else {
                    setHasMoreListings(userListings.length >= 20);
                }
            } else {
                console.log('Unexpected listings data structure:', response.data);
                setListings([]);
                setHasMoreListings(false);
            }
            
            // Обновляем состояние сортировки если переданы новые параметры
            if (sortParams.field) setSortField(sortParams.field);
            if (sortParams.direction) setSortDirection(sortParams.direction);
            
            setPage(1);
            
            // Сбрасываем выбранные объявления при обновлении данных
            setSelectedListings([]);
            
        } catch (err) {
            console.error('Error fetching listings:', err);
            setError(t('listings.errors.loadFailed', { defaultValue: 'Не удалось загрузить объявления' }));
            setListings([]);
            setHasMoreListings(false);
        } finally {
            setLoading(false);
        }
    };

    // Функция для загрузки дополнительных объявлений
    const fetchMoreListings = async () => {
        if (!hasMoreListings || loadingMore || !user?.id) return;
        
        try {
            setLoadingMore(true);
            const nextPage = page + 1;
            
            // Формируем параметры API сортировки
            let apiSort = '';
            switch (sortField) {
                case 'created_at':
                    apiSort = `date_${sortDirection}`;
                    break;
                case 'title':
                    apiSort = `title_${sortDirection}`;
                    break;
                case 'price':
                    apiSort = `price_${sortDirection}`;
                    break;
                case 'reviews':
                    apiSort = `rating_${sortDirection}`;
                    break;
                default:
                    apiSort = `date_${sortDirection}`;
            }
            
            const params = { 
                page: nextPage,
                size: 20,
                sort_by: apiSort,
                user_id: user.id
            };
            
            // Добавляем фильтр статуса, если выбран не "all"
            if (statusFilter !== 'all') {
                params.status = statusFilter;
            }
            
            console.log(`Загрузка дополнительных объявлений с сортировкой: ${apiSort}`);
            
            const response = await axios.get('/api/v1/marketplace/listings', {
                params,
                withCredentials: true
            });
            
            // Обработка ответа
            if (response.data?.data?.data && Array.isArray(response.data.data.data)) {
                // Фильтр объявлений по ID пользователя
                const newListings = response.data.data.data.filter(listing =>
                    String(listing.user_id) === String(user.id)
                );
                
                if (newListings.length > 0) {
                    setListings(prev => [...prev, ...newListings]);
                    setPage(nextPage);
                    
                    // Проверяем, есть ли еще объявления для загрузки
                    if (response.data.data.meta) {
                        const total = response.data.data.meta.total || 0;
                        setHasMoreListings(listings.length + newListings.length < total);
                    } else {
                        setHasMoreListings(newListings.length >= 20);
                    }
                } else {
                    setHasMoreListings(false);
                }
            } else {
                setHasMoreListings(false);
            }
        } catch (err) {
            console.error('Error fetching more listings:', err);
        } finally {
            setLoadingMore(false);
        }
    };

    // Обработчик для бесконечной прокрутки
    const handleLoadMore = () => {
        fetchMoreListings();
    };
    
    // Обработчик переключения режима отображения (сетка/список)
    const handleViewModeChange = (event, newMode) => {
        if (newMode !== null) {
            setViewMode(newMode);
            // Сохраняем предпочтение пользователя в localStorage
            localStorage.setItem('my-listings-view-mode', newMode);
        }
    };
    
    // Обработчик изменения сортировки
    const handleSortChange = (field, direction) => {
        console.log(`Обработчик сортировки получил: поле=${field}, направление=${direction}`);
        
        // Сбрасываем предыдущие объявления и состояние пагинации
        setListings([]);
        setPage(1);
        setHasMoreListings(true);
        
        // Запрашиваем данные с новыми параметрами сортировки
        fetchInitialData({ field, direction });
    };
    
    // Обработчик выбора объявления
    const handleSelectListing = (listingId) => {
        setSelectedListings(prev => {
            if (prev.includes(listingId)) {
                return prev.filter(id => id !== listingId);
            } else {
                return [...prev, listingId];
            }
        });
    };
    
    // Обработчик выбора всех объявлений
    const handleSelectAllListings = (isSelected) => {
        if (isSelected) {
            setSelectedListings(listings.map(listing => listing.id));
        } else {
            setSelectedListings([]);
        }
    };
    
    // Обработчик фильтрации по статусу
    const handleStatusFilterChange = (newStatus) => {
        if (newStatus !== statusFilter) {
            setStatusFilter(newStatus);
            // Перезагрузка данных с новым фильтром
            setListings([]);
            setPage(1);
            setHasMoreListings(true);
            
            // После изменения состояния выполняем запрос
            setTimeout(() => {
                fetchInitialData({});
            }, 0);
        }
    };
    
    // Обработчик открытия меню действий
    const handleOpenActionsMenu = (event) => {
        setMenuAnchorEl(event.currentTarget);
    };
    
    // Обработчик закрытия меню действий
    const handleCloseActionsMenu = () => {
        setMenuAnchorEl(null);
    };
    
    // Обработчик для множественного действия с объявлениями
    const handleBulkAction = (action) => {
        handleCloseActionsMenu();
        
        // Различные действия в зависимости от выбранного пункта меню
        switch (action) {
            case 'activate':
                setConfirmDialogAction('activate');
                setConfirmDialogOpen(true);
                break;
            case 'deactivate':
                setConfirmDialogAction('deactivate');
                setConfirmDialogOpen(true);
                break;
            case 'delete':
                setConfirmDialogAction('delete');
                setConfirmDialogOpen(true);
                break;
            case 'copy':
                // Копируем первое выбранное объявление
                if (selectedListings.length > 0) {
                    handleCopyListing(selectedListings[0]);
                }
                break;
            case 'promote':
                // Открываем диалог выбора типа продвижения
                setPromotionDialogOpen(true);
                break;
            case 'translate':
                // Запускаем перевод выбранных объявлений
                handleTranslateListings();
                break;
            default:
                break;
        }
    };
    
    // Обработчик для изменения статуса объявлений
    const handleChangeStatus = async (status) => {
        try {
            setLoading(true);
            
            // Последовательно обновляем статус каждого выбранного объявления
            const updatePromises = selectedListings.map(listingId => {
                return axios.put(`/api/v1/marketplace/listings/${listingId}`, {
                    status: status
                });
            });
            
            await Promise.all(updatePromises);
            
            // Показываем сообщение об успехе
            setSnackbarMessage(t('listings.actions.statusUpdateSuccess', {
                count: selectedListings.length
            }));
            setSnackbarOpen(true);
            
            // Обновляем список объявлений
            fetchInitialData({});
            
        } catch (error) {
            console.error('Ошибка обновления статуса:', error);
            setSnackbarMessage(t('listings.actions.statusUpdateError'));
            setSnackbarOpen(true);
        } finally {
            setLoading(false);
            setConfirmDialogOpen(false);
        }
    };
    
    // Обработчик для удаления объявлений
    const handleDeleteListings = async () => {
        try {
            setLoading(true);
            
            // Последовательно удаляем каждое выбранное объявление
            const deletePromises = selectedListings.map(listingId => {
                return axios.delete(`/api/v1/marketplace/listings/${listingId}`);
            });
            
            await Promise.all(deletePromises);
            
            // Показываем сообщение об успехе
            setSnackbarMessage(t('listings.actions.deleteSuccess', {
                count: selectedListings.length
            }));
            setSnackbarOpen(true);
            
            // Обновляем список объявлений
            fetchInitialData({});
            
        } catch (error) {
            console.error('Ошибка удаления объявлений:', error);
            setSnackbarMessage(t('listings.actions.deleteError'));
            setSnackbarOpen(true);
        } finally {
            setLoading(false);
            setConfirmDialogOpen(false);
        }
    };
    
    // Обработчик для копирования объявления
    const handleCopyListing = async (listingId) => {
        try {
            setLoading(true);
            
            // Получаем данные объявления
            const response = await axios.get(`/api/v1/marketplace/listings/${listingId}`);
            const listing = response.data.data;
            
            // Создаем новое объявление на основе существующего
            const newListingData = {
                title: `${listing.title} (${t('listings.actions.copy')})`,
                description: listing.description,
                price: listing.price,
                category_id: listing.category_id,
                condition: listing.condition,
                location: listing.location,
                latitude: listing.latitude,
                longitude: listing.longitude,
                city: listing.city,
                country: listing.country,
                attributes: listing.attributes
            };
            
            // Отправляем запрос на создание нового объявления
            const createResponse = await axios.post('/api/v1/marketplace/listings', newListingData);
            const newListingId = createResponse.data.data.id;
            
            // Показываем сообщение об успехе
            setSnackbarMessage(t('listings.actions.copySuccess'));
            setSnackbarOpen(true);
            
            // Обновляем список объявлений
            fetchInitialData({});
            
        } catch (error) {
            console.error('Ошибка копирования объявления:', error);
            setSnackbarMessage(t('listings.actions.copyError'));
            setSnackbarOpen(true);
        } finally {
            setLoading(false);
        }
    };
    
    // Обработчик для перевода объявлений
    const handleTranslateListings = async () => {
        try {
            setLoading(true);
            
            // Запрос на перевод выбранных объявлений
            const response = await axios.post('/api/v1/marketplace/translations/batch', {
                listing_ids: selectedListings,
                target_languages: ["en", "ru", "sr"]
            });
            
            // Показываем сообщение об успехе или о недостатке средств
            if (response.data.success) {
                setSnackbarMessage(t('listings.actions.translateSuccess', {
                    count: selectedListings.length,
                    cost: response.data.cost
                }));
                // Обновляем баланс
                setBalance(prev => prev - response.data.cost);
            } else {
                setSnackbarMessage(t('listings.actions.translateError'));
            }
            setSnackbarOpen(true);
            
        } catch (error) {
            console.error('Ошибка перевода объявлений:', error);
            // Проверяем, является ли ошибка недостатком средств (HTTP 402)
            if (error.response && error.response.status === 402) {
                setSnackbarMessage(t('listings.actions.insufficientFunds'));
                // Предлагаем пополнить баланс
                setTimeout(() => {
                    setDepositDialogOpen(true);
                }, 2000);
            } else {
                setSnackbarMessage(t('listings.actions.translateError'));
            }
            setSnackbarOpen(true);
        } finally {
            setLoading(false);
        }
    };
    
    // Обработчик для продвижения объявлений
    const handlePromoteListings = async (promotionType) => {
        try {
            setLoading(true);
            setPromotionDialogOpen(false);
            
            // Цены за различные типы продвижения
            const promotionPrices = {
                'toplist': 50,
                'highlight': 100,
                'vip': 200,
                'auto_refresh': 150,
                'turbo': 300
            };
            
            const price = promotionPrices[promotionType] * selectedListings.length;
            
            // Проверяем, достаточно ли средств
            if (balance < price) {
                setSnackbarMessage(t('listings.actions.insufficientFunds'));
                setSnackbarOpen(true);
                // Предлагаем пополнить баланс
                setTimeout(() => {
                    setDepositDialogOpen(true);
                }, 2000);
                return;
            }
            
            // Отправляем запрос на продвижение объявлений
            const response = await axios.post('/api/v1/marketplace/promotions', {
                listing_ids: selectedListings,
                promotion_type: promotionType
            });
            
            // Показываем сообщение об успехе
            setSnackbarMessage(t('listings.actions.promoteSuccess', {
                count: selectedListings.length,
                type: t(`listings.promotions.${promotionType}`),
                cost: price
            }));
            setSnackbarOpen(true);
            
            // Обновляем баланс и список объявлений
            setBalance(prev => prev - price);
            fetchInitialData({});
            
        } catch (error) {
            console.error('Ошибка продвижения объявлений:', error);
            setSnackbarMessage(t('listings.actions.promoteError'));
            setSnackbarOpen(true);
        } finally {
            setLoading(false);
        }
    };
    
    // Обработчик для пополнения баланса
    const handleBalanceUpdate = (newBalance) => {
        setBalance(newBalance);
        setDepositDialogOpen(false);
    };
    
    // Обработчик подтверждения действия в диалоге
    const handleConfirmAction = () => {
        if (confirmDialogAction === 'activate') {
            handleChangeStatus('active');
        } else if (confirmDialogAction === 'deactivate') {
            handleChangeStatus('inactive');
        } else if (confirmDialogAction === 'delete') {
            handleDeleteListings();
        }
    };

    useEffect(() => {
        if (user?.id) {
            fetchInitialData();
        } else {
            setLoading(false);
        }
    }, [user]);

    if (loading) {
        return (
            <Box display="flex" justifyContent="center" p={4}>
                <CircularProgress />
            </Box>
        );
    }

    if (error) {
        return (
            <Container>
                <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>
            </Container>
        );
    }

    return (
        <Container sx={{ py: 4 }}>
            {/* Заголовок и инструменты */}
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
                <Typography variant="h4" component="h1">
                    {t('listings.myListings')}
                </Typography>
                <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                    {/* Переключатель режима отображения */}
                    <ToggleButtonGroup
                        value={viewMode}
                        exclusive
                        onChange={handleViewModeChange}
                        aria-label="view mode"
                        size="small"
                    >
                        <ToggleButton value="grid" aria-label="grid view">
                            <GridIcon size={18} />
                        </ToggleButton>
                        <ToggleButton value="list" aria-label="list view">
                            <List size={18} />
                        </ToggleButton>
                    </ToggleButtonGroup>
                    
                    <Button
                        id="createAnnouncementButton"
                        component={Link}
                        to="/marketplace/create"
                        variant="contained"
                        startIcon={<Plus />}
                    >
                        {t('listings.create.title')}
                    </Button>
                </Box>
            </Box>

            {/* Панель управления и фильтрации */}
            <Toolbar
                sx={{
                    pl: { sm: 2 },
                    pr: { xs: 1, sm: 1 },
                    mb: 2,
                    bgcolor: 'background.paper',
                    borderRadius: 1,
                    display: 'flex',
                    justifyContent: 'space-between'
                }}
                variant="dense"
            >
                <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    {/* Выбор статуса объявлений */}
                    <Box sx={{ mr: 2 }}>
                        <ToggleButtonGroup
                            value={statusFilter}
                            exclusive
                            onChange={(e, value) => value && handleStatusFilterChange(value)}
                            size="small"
                        >
                            <ToggleButton value="active">
                                <Check size={16} style={{ marginRight: 4 }} />
                                {t('listings.status.active')}
                            </ToggleButton>
                            <ToggleButton value="inactive">
                                <AlertTriangle size={16} style={{ marginRight: 4 }} />
                                {t('listings.status.inactive')}
                            </ToggleButton>
                            <ToggleButton value="all">
                                <Filter size={16} style={{ marginRight: 4 }} />
                                {t('listings.status.all')}
                            </ToggleButton>
                        </ToggleButtonGroup>
                    </Box>
                    
                    {/* Чип с количеством выбранных объявлений */}
                    {selectedListings.length > 0 && (
                        <Chip 
                            label={t('listings.selected', { count: selectedListings.length })}
                            color="primary"
                            variant="outlined"
                            onDelete={() => setSelectedListings([])}
                        />
                    )}
                </Box>
                
                {/* Кнопка действий с выбранными объявлениями */}
                <Button
                    variant="outlined"
                    disabled={selectedListings.length === 0}
                    onClick={handleOpenActionsMenu}
                    sx={{ ml: 1 }}
                >
                    {t('listings.actions.perform')}
                </Button>
            </Toolbar>

            {/* Основной контент: список или сетка объявлений */}
            {listings.length === 0 ? (
                <Alert severity="info">
                    {statusFilter === 'active' 
                        ? t('listings.noActiveListings')
                        : statusFilter === 'inactive'
                            ? t('listings.noInactiveListings')
                            : t('listings.Youdonthave')
                    }
                </Alert>
            ) : (
                <InfiniteScroll
                    hasMore={hasMoreListings}
                    loading={loadingMore}
                    onLoadMore={handleLoadMore}
                    autoLoad={!isMobile} // Автозагрузка только на десктопе
                    loadingMessage={t('marketplace:listings.loading', { defaultValue: 'Загрузка...' })}
                    loadMoreButtonText={t('marketplace:listings.loadMore', { defaultValue: 'Показать ещё' })}
                    noMoreItemsText={t('marketplace:listings.noMoreListings', { defaultValue: 'Больше нет объявлений' })}
                >
                    {viewMode === 'grid' ? (
                        <Grid container spacing={3}>
                            {listings.map((listing) => (
                                <Grid item xs={12} sm={6} md={4} key={listing.id}>
                                    <Box sx={{ position: 'relative' }}>
                                        {/* Чекбокс для выбора объявления */}
                                        <Checkbox
                                            checked={selectedListings.includes(listing.id)}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                handleSelectListing(listing.id);
                                            }}
                                            sx={{
                                                position: 'absolute',
                                                top: 8,
                                                left: 8,
                                                zIndex: 10,
                                                bgcolor: 'rgba(255,255,255,0.8)',
                                                borderRadius: '50%'
                                            }}
                                        />
                                        <Link
                                            to={`/marketplace/listings/${listing.id}`}
                                            style={{ textDecoration: 'none' }}
                                        >
                                            <ListingCard listing={listing} />
                                        </Link>
                                    </Box>
                                </Grid>
                            ))}
                        </Grid>
                    ) : (
                        <MarketplaceListingsList
                            listings={listings}
                            showSelection={true}
                            selectedItems={selectedListings}
                            onSelectItem={handleSelectListing}
                            onSelectAll={handleSelectAllListings}
                            initialSortField={sortField}
                            initialSortOrder={sortDirection}
                            onSortChange={handleSortChange}
                            filters={{
                                sort_by: `${sortField === 'created_at' ? 'date' : sortField}_${sortDirection}`
                            }}
                        />
                    )}
                </InfiniteScroll>
            )}
            
            {/* Меню действий с выбранными объявлениями */}
            <Menu
                anchorEl={menuAnchorEl}
                open={Boolean(menuAnchorEl)}
                onClose={handleCloseActionsMenu}
            >
                {/* Бесплатные действия */}
                <MenuItem onClick={() => handleBulkAction('activate')} disabled={statusFilter === 'active'}>
                    <ListItemIcon>
                        <Check size={16} color={theme.palette.success.main} />
                    </ListItemIcon>
                    <ListItemText primary={t('listings.actions.activate')} />
                </MenuItem>
                <MenuItem onClick={() => handleBulkAction('deactivate')} disabled={statusFilter === 'inactive'}>
                    <ListItemIcon>
                        <AlertTriangle size={16} color={theme.palette.warning.main} />
                    </ListItemIcon>
                    <ListItemText primary={t('listings.actions.deactivate')} />
                </MenuItem>
                <MenuItem onClick={() => handleBulkAction('delete')}>
                    <ListItemIcon>
                        <Trash2 size={16} color={theme.palette.error.main} />
                    </ListItemIcon>
                    <ListItemText primary={t('listings.actions.delete')} />
                </MenuItem>
                <MenuItem onClick={() => handleBulkAction('copy')} disabled={selectedListings.length !== 1}>
                    <ListItemIcon>
                        <Copy size={16} />
                    </ListItemIcon>
                    <ListItemText primary={t('listings.actions.copy')} />
                </MenuItem>
                
                <Divider />
                
                {/* Платные действия */}
                <MenuItem onClick={() => handleBulkAction('promote')}>
                    <ListItemIcon>
                        <ArrowUp size={16} color={theme.palette.primary.main} />
                    </ListItemIcon>
                    <ListItemText primary={t('listings.actions.promote')} />
                </MenuItem>
                <MenuItem onClick={() => handleBulkAction('translate')}>
                    <ListItemIcon>
                        <Zap size={16} color={theme.palette.secondary.main} />
                    </ListItemIcon>
                    <ListItemText primary={t('listings.actions.translate')} />
                </MenuItem>
            </Menu>
            
            {/* Диалог подтверждения действия */}
            <Dialog
                open={confirmDialogOpen}
                onClose={() => setConfirmDialogOpen(false)}
            >
                <DialogTitle>
                    {confirmDialogAction === 'activate' && t('listings.confirmations.activateTitle')}
                    {confirmDialogAction === 'deactivate' && t('listings.confirmations.deactivateTitle')}
                    {confirmDialogAction === 'delete' && t('listings.confirmations.deleteTitle')}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        {confirmDialogAction === 'activate' && t('listings.confirmations.activateText', { count: selectedListings.length })}
                        {confirmDialogAction === 'deactivate' && t('listings.confirmations.deactivateText', { count: selectedListings.length })}
                        {confirmDialogAction === 'delete' && t('listings.confirmations.deleteText', { count: selectedListings.length })}
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setConfirmDialogOpen(false)} color="primary">
                        {t('common:buttons.cancel')}
                    </Button>
                    <Button 
                        onClick={handleConfirmAction} 
                        color={confirmDialogAction === 'delete' ? 'error' : 'primary'}
                        variant="contained"
                        autoFocus
                    >
                        {confirmDialogAction === 'activate' && t('listings.actions.activate')}
                        {confirmDialogAction === 'deactivate' && t('listings.actions.deactivate')}
                        {confirmDialogAction === 'delete' && t('listings.actions.delete')}
                    </Button>
                </DialogActions>
            </Dialog>
            
            {/* Диалог выбора типа продвижения */}
            <Dialog
                open={promotionDialogOpen}
                onClose={() => setPromotionDialogOpen(false)}
                maxWidth="sm"
                fullWidth
            >
                <DialogTitle>{t('listings.promotions.title')}</DialogTitle>
                <DialogContent>
                    <DialogContentText sx={{ mb: 2 }}>
                        {t('listings.promotions.description')}
                    </DialogContentText>
                    
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                        {/* Поднятие в топ */}
                        <Button
                            variant="outlined"
                            startIcon={<ArrowUp />}
                            onClick={() => handlePromoteListings('toplist')}
                            sx={{ justifyContent: 'space-between' }}
                        >
                            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start' }}>
                                <Typography variant="subtitle1">{t('listings.promotions.toplist')}</Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {t('listings.promotions.toplistDescription')}
                                </Typography>
                            </Box>
                            <Chip label="50 RSD" color="primary" />
                        </Button>
                        
                        {/* Выделение цветом */}
                        <Button
                            variant="outlined"
                            startIcon={<Star color={theme.palette.warning.main} />}
                            onClick={() => handlePromoteListings('highlight')}
                            sx={{ justifyContent: 'space-between' }}
                        >
                            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start' }}>
                                <Typography variant="subtitle1">{t('listings.promotions.highlight')}</Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {t('listings.promotions.highlightDescription')}
                                </Typography>
                            </Box>
                            <Chip label="100 RSD" color="primary" />
                        </Button>
                        
                        {/* VIP-статус */}
                        <Button
                            variant="outlined"
                            startIcon={<Star color={theme.palette.primary.main} />}
                            onClick={() => handlePromoteListings('vip')}
                            sx={{ justifyContent: 'space-between' }}
                        >
                            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start' }}>
                                <Typography variant="subtitle1">{t('listings.promotions.vip')}</Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {t('listings.promotions.vipDescription')}
                                </Typography>
                            </Box>
                            <Chip label="200 RSD" color="primary" />
                        </Button>
                        
                        {/* Автообновление */}
                        <Button
                            variant="outlined"
                            startIcon={<ArrowUp />}
                            onClick={() => handlePromoteListings('auto_refresh')}
                            sx={{ justifyContent: 'space-between' }}
                        >
                            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start' }}>
                                <Typography variant="subtitle1">{t('listings.promotions.autoRefresh')}</Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {t('listings.promotions.autoRefreshDescription')}
                                </Typography>
                            </Box>
                            <Chip label="150 RSD" color="primary" />
                        </Button>
                        
                        {/* Турбо-продажа */}
                        <Button
                            variant="outlined"
                            startIcon={<Zap color={theme.palette.error.main} />}
                            onClick={() => handlePromoteListings('turbo')}
                            sx={{ justifyContent: 'space-between' }}
                        >
                            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start' }}>
                                <Typography variant="subtitle1">{t('listings.promotions.turbo')}</Typography>
                                <Typography variant="body2" color="text.secondary">
                                    {t('listings.promotions.turboDescription')}
                                </Typography>
                            </Box>
                            <Chip label="300 RSD" color="primary" />
                        </Button>
                    </Box>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setPromotionDialogOpen(false)}>
                        {t('common:buttons.cancel')}
                    </Button>
                </DialogActions>
            </Dialog>
            
            {/* Диалог пополнения баланса */}
            <DepositDialog
                open={depositDialogOpen}
                onClose={() => setDepositDialogOpen(false)}
                onBalanceUpdate={handleBalanceUpdate}
            />
            
            {/* Уведомление-снекбар */}
            <Snackbar
                open={snackbarOpen}
                autoHideDuration={6000}
                onClose={() => setSnackbarOpen(false)}
                message={snackbarMessage}
            />
        </Container>
    );
};

export default MyListingsPage;