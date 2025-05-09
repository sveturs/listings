// frontend/hostel-frontend/src/pages/marketplace/MyListingsPage.tsx
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
import ListingCard, { Listing } from '../../components/marketplace/ListingCard';
import MarketplaceListingsList from '../../components/marketplace/MarketplaceListingsList';
import InfiniteScroll from '../../components/marketplace/InfiniteScroll';
import axios from '../../api/axios';
import DepositDialog from '../../components/balance/DepositDialog';

// Define type for status filter values
type StatusFilter = 'active' | 'inactive' | 'all';

// Define type for view mode values
type ViewMode = 'grid' | 'list';

// Define type for sort fields
type SortField = 'created_at' | 'title' | 'price' | 'reviews';

// Define type for sort directions
type SortDirection = 'asc' | 'desc';

// Define type for confirmation dialog actions
type ConfirmDialogAction = 'activate' | 'deactivate' | 'delete' | null;

// Define type for promotion types
type PromotionType = 'toplist' | 'highlight' | 'vip' | 'auto_refresh' | 'turbo';

// Define interface for API response
interface ListingsResponse {
    data?: {
        data: Listing[];
        meta?: {
            total: number;
            page: number;
            size: number;
        }
    }
}

// Define interface for balance response
interface BalanceResponse {
    data?: {
        balance: number;
    }
}

// Define interface for promotion response
interface PromotionResponse {
    success: boolean;
    cost?: number;
}

// Define interface for translation response
interface TranslationResponse {
    success: boolean;
    cost?: number;
}

const MyListingsPage: React.FC = () => {
    const { t } = useTranslation('marketplace');
    const theme = useTheme();
    const isMobile = useMediaQuery(theme.breakpoints.down('md'));
    const { user } = useAuth();
    
    // Pagination and infinite scroll states
    const [page, setPage] = useState<number>(1);
    const [hasMoreListings, setHasMoreListings] = useState<boolean>(true);
    const [loadingMore, setLoadingMore] = useState<boolean>(false);
    const [totalListings, setTotalListings] = useState<number>(0);
    
    // Main data states
    const [listings, setListings] = useState<Listing[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    
    // Sorting states
    const [sortField, setSortField] = useState<SortField>('created_at');
    const [sortDirection, setSortDirection] = useState<SortDirection>('desc');
    
    // View mode state (grid/list)
    const [viewMode, setViewMode] = useState<ViewMode>(() => {
        // Load saved view mode from localStorage
        const savedViewMode = localStorage.getItem('my-listings-view-mode');
        return (savedViewMode === 'list' ? 'list' : 'grid') as ViewMode;
    });

    // Listing selection management states
    const [selectedListings, setSelectedListings] = useState<string[]>([]);
    const [statusFilter, setStatusFilter] = useState<StatusFilter>('active');
    
    // Dialog and menu management states
    const [menuAnchorEl, setMenuAnchorEl] = useState<HTMLElement | null>(null);
    const [confirmDialogOpen, setConfirmDialogOpen] = useState<boolean>(false);
    const [confirmDialogAction, setConfirmDialogAction] = useState<ConfirmDialogAction>(null);
    const [depositDialogOpen, setDepositDialogOpen] = useState<boolean>(false);
    const [promotionDialogOpen, setPromotionDialogOpen] = useState<boolean>(false);
    const [snackbarOpen, setSnackbarOpen] = useState<boolean>(false);
    const [snackbarMessage, setSnackbarMessage] = useState<string>('');

    // Selected promotion type state
    const [selectedPromotionType, setSelectedPromotionType] = useState<PromotionType | null>(null);
    const [balance, setBalance] = useState<number>(0);

    // Function to fetch initial data with sort and filter parameters
    const fetchInitialData = async (sortParams: { field?: SortField, direction?: SortDirection } = {}): Promise<void> => {
        try {
            setLoading(true);
            setError(null);
            
            // Format API sort parameters
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
            
            const params: Record<string, any> = { 
                page: 1,
                size: 20,
                sort_by: apiSort,
                user_id: user?.id
            };
            
            // Add status filter if not 'all'
            if (statusFilter !== 'all') {
                params.status = statusFilter;
            }
            
            const response = await axios.get<ListingsResponse>('/api/v1/marketplace/listings', {
                params,
                withCredentials: true
            });
            
            // Get user balance in parallel
            try {
                const balanceResponse = await axios.get<BalanceResponse>('/api/v1/balance');
                if (balanceResponse.data?.data?.balance) {
                    setBalance(balanceResponse.data.data.balance);
                }
            } catch (err) {
                console.error('Ошибка получения баланса:', err);
            }
            
            // Process response
            if (response.data?.data?.data && Array.isArray(response.data.data.data)) {
                // Filter listings by user ID
                const userListings = response.data.data.data.filter(listing =>
                    String(listing.user_id) === String(user?.id)
                );
                setListings(userListings);
                
                // Update pagination information
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
            
            // Update sort state if new parameters are provided
            if (sortParams.field) setSortField(sortParams.field);
            if (sortParams.direction) setSortDirection(sortParams.direction);
            
            setPage(1);
            
            // Reset selected listings when data is refreshed
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

    // Function to fetch additional listings
    const fetchMoreListings = async (): Promise<void> => {
        if (!hasMoreListings || loadingMore || !user?.id) return;
        
        try {
            setLoadingMore(true);
            const nextPage = page + 1;
            
            // Format API sort parameters
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
            
            const params: Record<string, any> = { 
                page: nextPage,
                size: 20,
                sort_by: apiSort,
                user_id: user.id
            };
            
            // Add status filter if not 'all'
            if (statusFilter !== 'all') {
                params.status = statusFilter;
            }
            
            console.log(`Загрузка дополнительных объявлений с сортировкой: ${apiSort}`);
            
            const response = await axios.get<ListingsResponse>('/api/v1/marketplace/listings', {
                params,
                withCredentials: true
            });
            
            // Process response
            if (response.data?.data?.data && Array.isArray(response.data.data.data)) {
                // Filter listings by user ID
                const newListings = response.data.data.data.filter(listing =>
                    String(listing.user_id) === String(user.id)
                );
                
                if (newListings.length > 0) {
                    setListings(prev => [...prev, ...newListings]);
                    setPage(nextPage);
                    
                    // Check if there are more listings to load
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

    // Handler for infinite scroll
    const handleLoadMore = (): void => {
        fetchMoreListings();
    };
    
    // Handler for view mode toggle (grid/list)
    const handleViewModeChange = (event: React.MouseEvent<HTMLElement>, newMode: ViewMode | null): void => {
        if (newMode !== null) {
            setViewMode(newMode);
            // Save user preference to localStorage
            localStorage.setItem('my-listings-view-mode', newMode);
        }
    };
    
    // Handler for sort change
    const handleSortChange = (field: SortField, direction: SortDirection): void => {
        console.log(`Обработчик сортировки получил: поле=${field}, направление=${direction}`);
        
        // Reset previous listings and pagination state
        setListings([]);
        setPage(1);
        setHasMoreListings(true);
        
        // Request data with new sort parameters
        fetchInitialData({ field, direction });
    };
    
    // Handler for listing selection
    const handleSelectListing = (listingId: string): void => {
        setSelectedListings(prev => {
            if (prev.includes(listingId)) {
                return prev.filter(id => id !== listingId);
            } else {
                return [...prev, listingId];
            }
        });
    };
    
    // Handler for selecting all listings
    const handleSelectAllListings = (isSelected: boolean): void => {
        if (isSelected) {
            setSelectedListings(listings.map(listing => listing.id.toString()));
        } else {
            setSelectedListings([]);
        }
    };
    
    // Handler for status filter change
    const handleStatusFilterChange = (newStatus: StatusFilter): void => {
        if (newStatus !== statusFilter) {
            setStatusFilter(newStatus);
            // Reload data with new filter
            setListings([]);
            setPage(1);
            setHasMoreListings(true);
            
            // Schedule fetch after state update
            setTimeout(() => {
                fetchInitialData({});
            }, 0);
        }
    };
    
    // Handler for opening actions menu
    const handleOpenActionsMenu = (event: React.MouseEvent<HTMLButtonElement>): void => {
        setMenuAnchorEl(event.currentTarget);
    };
    
    // Handler for closing actions menu
    const handleCloseActionsMenu = (): void => {
        setMenuAnchorEl(null);
    };
    
    // Handler for bulk actions on listings
    const handleBulkAction = (action: string): void => {
        handleCloseActionsMenu();
        
        // Different actions based on selected menu item
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
                // Copy the first selected listing
                if (selectedListings.length > 0) {
                    handleCopyListing(selectedListings[0]);
                }
                break;
            case 'promote':
                // Open promotion type selection dialog
                setPromotionDialogOpen(true);
                break;
            case 'translate':
                // Start translation of selected listings
                handleTranslateListings();
                break;
            default:
                break;
        }
    };
    
    // Handler for changing listing status
    const handleChangeStatus = async (status: string): Promise<void> => {
        try {
            setLoading(true);
            
            // Update status for each selected listing sequentially
            const updatePromises = selectedListings.map(listingId => {
                // First get current listing to preserve its category
                return axios.get(`/api/v1/marketplace/listings/${listingId}`)
                    .then(response => {
                        const listing = response.data;
                        return axios.put(`/api/v1/marketplace/listings/${listingId}`, {
                            status: status,
                            category_id: listing.category_id // Explicitly send category ID
                        });
                    });
            });
            
            await Promise.all(updatePromises);
            
            // Show success message
            setSnackbarMessage(t('listings.actions.statusUpdateSuccess', {
                count: selectedListings.length
            }));
            setSnackbarOpen(true);
            
            // Refresh listings list
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
    
    // Handler for deleting listings
    const handleDeleteListings = async (): Promise<void> => {
        try {
            setLoading(true);
            
            // Delete each selected listing sequentially
            const deletePromises = selectedListings.map(listingId => {
                return axios.delete(`/api/v1/marketplace/listings/${listingId}`);
            });
            
            await Promise.all(deletePromises);
            
            // Show success message
            setSnackbarMessage(t('listings.actions.deleteSuccess', {
                count: selectedListings.length
            }));
            setSnackbarOpen(true);
            
            // Refresh listings list
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
    
    // Handler for copying listing
    const handleCopyListing = async (listingId: string): Promise<void> => {
        try {
            setLoading(true);
            
            // Get listing data
            const response = await axios.get<{ data: Listing }>(`/api/v1/marketplace/listings/${listingId}`);
            const listing = response.data.data;
            
            // Create new listing based on existing one
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
            
            // Send request to create new listing
            const createResponse = await axios.post('/api/v1/marketplace/listings', newListingData);
            const newListingId = createResponse.data.data.id;
            
            // Show success message
            setSnackbarMessage(t('listings.actions.copySuccess'));
            setSnackbarOpen(true);
            
            // Refresh listings list
            fetchInitialData({});
            
        } catch (error) {
            console.error('Ошибка копирования объявления:', error);
            setSnackbarMessage(t('listings.actions.copyError'));
            setSnackbarOpen(true);
        } finally {
            setLoading(false);
        }
    };
    
    // Handler for translating listings
    const handleTranslateListings = async (): Promise<void> => {
        try {
            setLoading(true);
            
            // Request to translate selected listings
            const response = await axios.post<TranslationResponse>('/api/v1/marketplace/translations/batch', {
                listing_ids: selectedListings,
                target_languages: ["en", "ru", "sr"]
            });
            
            // Show success message or insufficient funds message
            if (response.data.success) {
                setSnackbarMessage(t('listings.actions.translateSuccess', {
                    count: selectedListings.length,
                    cost: response.data.cost
                }));
                // Update balance
                setBalance(prev => prev - (response.data.cost || 0));
            } else {
                setSnackbarMessage(t('listings.actions.translateError'));
            }
            setSnackbarOpen(true);
            
        } catch (error) {
            console.error('Ошибка перевода объявлений:', error);
            // Check if error is insufficient funds (HTTP 402)
            if (error && typeof error === 'object' && 'response' in error && error.response?.status === 402) {
                setSnackbarMessage(t('listings.actions.insufficientFunds'));
                // Offer to deposit funds
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
    
    // Handler for promoting listings
    const handlePromoteListings = async (promotionType: PromotionType): Promise<void> => {
        try {
            setLoading(true);
            setPromotionDialogOpen(false);
            
            // Prices for different promotion types
            const promotionPrices: Record<PromotionType, number> = {
                'toplist': 50,
                'highlight': 100,
                'vip': 200,
                'auto_refresh': 150,
                'turbo': 300
            };
            
            const price = promotionPrices[promotionType] * selectedListings.length;
            
            // Check if there are sufficient funds
            if (balance < price) {
                setSnackbarMessage(t('listings.actions.insufficientFunds'));
                setSnackbarOpen(true);
                // Offer to deposit funds
                setTimeout(() => {
                    setDepositDialogOpen(true);
                }, 2000);
                return;
            }
            
            // Send request to promote listings
            const response = await axios.post<PromotionResponse>('/api/v1/marketplace/promotions', {
                listing_ids: selectedListings,
                promotion_type: promotionType
            });
            
            // Show success message
            setSnackbarMessage(t('listings.actions.promoteSuccess', {
                count: selectedListings.length,
                type: t(`listings.promotions.${promotionType}`),
                cost: price
            }));
            setSnackbarOpen(true);
            
            // Update balance and listings list
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
    
    // Handler for balance update
    const handleBalanceUpdate = (newBalance: number): void => {
        setBalance(newBalance);
        setDepositDialogOpen(false);
    };
    
    // Handler for dialog action confirmation
    const handleConfirmAction = (): void => {
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
            {/* Title and tools */}
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
                <Typography variant="h4" component="h1">
                    {t('listings.myListings')}
                </Typography>
                <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
                    {/* View mode toggle */}
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

            {/* Control and filter panel */}
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
                    {/* Listing status selection */}
                    <Box sx={{ mr: 2 }}>
                        <ToggleButtonGroup
                            value={statusFilter}
                            exclusive
                            onChange={(e, value) => value && handleStatusFilterChange(value as StatusFilter)}
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
                    
                    {/* Chip with selected listings count */}
                    {selectedListings.length > 0 && (
                        <Chip 
                            label={t('listings.selected', { count: selectedListings.length })}
                            color="primary"
                            variant="outlined"
                            onDelete={() => setSelectedListings([])}
                        />
                    )}
                </Box>
                
                {/* Button for actions with selected listings */}
                <Button
                    variant="outlined"
                    disabled={selectedListings.length === 0}
                    onClick={handleOpenActionsMenu}
                    sx={{ ml: 1 }}
                >
                    {t('listings.actions.perform')}
                </Button>
            </Toolbar>

            {/* Main content: listing grid or list */}
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
                    autoLoad={!isMobile} // Auto-load only on desktop
                    loadingMessage={t('marketplace:listings.loading', { defaultValue: 'Загрузка...' })}
                    loadMoreButtonText={t('marketplace:listings.loadMore', { defaultValue: 'Показать ещё' })}
                    noMoreItemsText={t('marketplace:listings.noMoreListings', { defaultValue: 'Больше нет объявлений' })}
                >
                    {viewMode === 'grid' ? (
                        <Grid container spacing={3}>
                            {listings.map((listing) => (
                                <Grid item xs={12} sm={6} md={4} key={listing.id}>
                                    <Box sx={{ position: 'relative' }}>
                                        {/* Checkbox for listing selection */}
                                        <Checkbox
                                            checked={selectedListings.includes(listing.id.toString())}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                handleSelectListing(listing.id.toString());
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
                                            <ListingCard listing={listing} showStatus={true} />
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
            
            {/* Actions menu for selected listings */}
            <Menu
                anchorEl={menuAnchorEl}
                open={Boolean(menuAnchorEl)}
                onClose={handleCloseActionsMenu}
            >
                {/* Free actions */}
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
                
                {/* Paid actions */}
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
            
            {/* Confirmation dialog */}
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
            
            {/* Promotion type selection dialog */}
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
                        {/* Top list placement */}
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
                        
                        {/* Highlight color */}
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
                        
                        {/* VIP status */}
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
                        
                        {/* Auto refresh */}
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
                        
                        {/* Turbo sale */}
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
            
            {/* Balance deposit dialog */}
            <DepositDialog
                open={depositDialogOpen}
                onClose={() => setDepositDialogOpen(false)}
                onBalanceUpdate={handleBalanceUpdate}
            />
            
            {/* Snackbar notification */}
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