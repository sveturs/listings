// frontend/hostel-frontend/src/pages/store/StorefrontDetailPage.tsx
import React, { useState, useEffect, ChangeEvent, ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams, useNavigate, Link } from 'react-router-dom';
import {
    Container,
    Typography,
    Button,
    Box,
    CircularProgress,
    Alert,
    Tabs,
    Tab,
    Paper,
    Divider,
    Grid,
    Card,
    CardContent,
    CardHeader,
    TextField,
    IconButton,
    Modal,
    Stack,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    ButtonGroup,
    Tooltip,
    ToggleButtonGroup,
    ToggleButton,
    Checkbox
} from '@mui/material';
import {
    Edit,
    Trash2,
    Plus,
    Database,
    Upload,
    Settings,
    BarChart,
    RefreshCw,
    FileType,
    AlertTriangle,
    Grid as GridIcon,
    List,
    LanguagesIcon, // Используем LanguagesIcon вместо Translate
    Check,
    Info
} from 'lucide-react';
import axios from '../../api/axios';
import { useAuth } from '../../contexts/AuthContext';
import ListingCard from '../../components/marketplace/ListingCard';
import StorefrontListingsList from '../../components/store/StorefrontListingsList';
import BatchActionsBar from '../../components/store/BatchActionsBar';
import ListingsPagination from '../../components/store/ListingsPagination';
import TranslationPaymentDialog from '../../components/store/TranslationPaymentDialog';
import ImportModal from '../../components/store/ImportModal';
import ImportSourceList from '../../components/store/ImportSourceList';

interface TabPanelProps {
    children: ReactNode;
    value: number;
    index: number;
    [key: string]: any;
}

interface Storefront {
    id: number;
    name: string;
    description: string;
    slug: string;
    status: string;
    phone: string;
    email: string;
    website: string;
    address: string;
    city: string;
    country: string;
    latitude: number | null;
    longitude: number | null;
    logo_path: string | null;
    created_at: string;
    updated_at: string;
    [key: string]: any;
}

interface Listing {
    id: number;
    title: string;
    description: string;
    price: number;
    images: any[];
    status: string;
    created_at: string;
    category_id: number;
    user_id: number;
    [key: string]: any;
}

interface ImportSource {
    id: number;
    type: string;
    url: string;
    storefront_id: number;
    created_at: string;
    updated_at: string;
    [key: string]: any;
}

interface ImportForm {
    type: 'csv' | 'xml' | 'json';
    url: string;
    storefront_id: number;
}

type RouteParams = {
    id: string;
}

interface ImportHistories {
    [key: number]: any[];
}

const TabPanel: React.FC<TabPanelProps> = (props) => {
    const { children, value, index, ...other } = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`storefront-tabpanel-${index}`}
            aria-labelledby={`storefront-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box sx={{ pt: 3 }}>
                    {children}
                </Box>
            )}
        </div>
    );
};

const StorefrontDetailPage: React.FC = () => {
    const { t } = useTranslation(['common', 'marketplace']);
    const [storeListings, setStoreListings] = useState<Listing[]>([]);
    const navigate = useNavigate();
    const { id } = useParams<keyof RouteParams>();
    const { user } = useAuth();

    const [loading, setLoading] = useState<boolean>(true);
    const [storefront, setStorefront] = useState<Storefront | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [importSources, setImportSources] = useState<ImportSource[]>([]);
    const [activeTab, setActiveTab] = useState<number>(0);
    const [importHistories, setImportHistories] = useState<ImportHistories>({});
    const [openImportModal, setOpenImportModal] = useState<boolean>(false);
    const [importForm, setImportForm] = useState<ImportForm>({
        type: 'csv',
        url: '',
        storefront_id: Number(id)
    });
    const [importFile, setImportFile] = useState<File | null>(null);
    const [importError, setImportError] = useState<string | null>(null);
    const [importLoading, setImportLoading] = useState<boolean>(false);
    const [runningImport, setRunningImport] = useState<number | null>(null);

    // Новые состояния для списочного режима и выбора элементов
    const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
    const [selectedListings, setSelectedListings] = useState<number[]>([]);
    const [deleteConfirmOpen, setDeleteConfirmOpen] = useState<boolean>(false);
    const [translationLoading, setTranslationLoading] = useState<boolean>(false);
    const [batchActionSuccess, setBatchActionSuccess] = useState<string | null>(null);

    // Состояния для пагинации
    const [page, setPage] = useState<number>(1);
    const [limit, setLimit] = useState<number>(20);
    const [totalItems, setTotalItems] = useState<number>(0);
    const [openTranslationDialog, setOpenTranslationDialog] = useState<boolean>(false);
    const [userBalance, setUserBalance] = useState<number>(0);
    const [importZipFile, setImportZipFile] = useState<File | null>(null);
    const [selectedSourceId, setSelectedSourceId] = useState<number | null>(null);
    
    const handleZipFileChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
        const file = e.target.files?.[0];
        if (file) {
            setImportZipFile(file);
        }
    };
    
    // Функция для загрузки данных с учетом пагинации
    const fetchData = async (): Promise<void> => {
        try {
            setLoading(true);
            const [storefrontResponse, sourcesResponse, listingsResponse] = await Promise.all([
                axios.get(`/api/v1/storefronts/${id}`),
                axios.get(`/api/v1/storefronts/${id}/import-sources`),
                axios.get('/api/v1/marketplace/listings', {
                    params: {
                        storefront_id: id,
                        page: page,
                        limit: limit
                    }
                })
            ]);

            setStorefront(storefrontResponse.data.data);
            setImportSources(sourcesResponse.data.data || []);

            // Используем data и meta из ответа
            const listings = listingsResponse.data.data?.data || [];
            const meta = listingsResponse.data.data?.meta || {};

            setStoreListings(listings);
            setTotalItems(meta.total || listings.length);
        } catch (err) {
            console.error('Error fetching storefront data:', err);
            setError(t('marketplace:store.errors.loadFailed'));
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (id) {
            fetchData();
        }
    }, [id, page, limit, t]);
    
    // Добавить useEffect для получения баланса пользователя
    useEffect(() => {
        const fetchUserBalance = async (): Promise<void> => {
            try {
                const response = await axios.get('/api/v1/balance');
                if (response.data?.data?.balance) {
                    setUserBalance(response.data.data.balance);
                }
            } catch (err) {
                console.error('Error fetching user balance:', err);
            }
        };

        fetchUserBalance();
    }, []);

    const fetchImportHistory = async (sourceId: number): Promise<void> => {
        try {
            const response = await axios.get(`/api/v1/storefronts/import-sources/${sourceId}/history`);
            setImportHistories(prev => ({
                ...prev,
                [sourceId]: response.data.data || []
            }));
        } catch (err) {
            console.error(`Error fetching import history for source ${sourceId}:`, err);
        }
    };

    const handleCreateImportSource = async (): Promise<void> => {
        try {
            setImportLoading(true);
            setImportError(null);

            if (!importForm.type) {
                setImportError(t('marketplace:store.import.selectTypeError'));
                return;
            }

            // Создаем минимальный формат для отправки
            const formData = {
                ...importForm,
                storefront_id: Number(id),
                // Если URL не указан, используем пустую строку (не null)
                url: importForm.url || ''
            };

            const response = await axios.post('/api/v1/storefronts/import-sources', formData);
            setImportSources(prev => [...prev, response.data.data]);
            setOpenImportModal(false);
            setImportForm({
                type: 'csv',
                url: '',
                storefront_id: Number(id)
            });
            setImportFile(null);

            // Показываем успешное сообщение с инструкцией
            alert(t('marketplace:store.import.sourceCreated'));
        } catch (err: any) {
            console.error('Error creating import source:', err);
            setImportError(`${t('marketplace:store.import.createError')}: ${err.response?.data?.error || err.message}`);
        } finally {
            setImportLoading(false);
        }
    };

    const handleDeleteImportSource = async (sourceId: number): Promise<void> => {
        if (!window.confirm(t('marketplace:store.import.deleteConfirm'))) {
            return;
        }

        try {
            await axios.delete(`/api/v1/storefronts/import-sources/${sourceId}`);
            setImportSources(prev => prev.filter(source => source.id !== sourceId));
        } catch (err) {
            console.error(`Error deleting import source ${sourceId}:`, err);
            alert(t('marketplace:store.import.deleteError'));
        }
    };

    const handleRunImport = async (sourceId: number): Promise<void> => {
        try {
            setRunningImport(sourceId);

            console.log(`Starting import for source ID: ${sourceId}, file: ${importFile ? importFile.name : 'none'}`);

            // Если нет ни файла, ни URL, показываем ошибку
            const source = importSources.find(s => s.id === sourceId);
            if (!importFile && (!source || !source.url)) {
                alert(t('marketplace:store.import.needFileOrUrl'));
                setRunningImport(null);
                return;
            }

            let response;

            if (importFile) {
                // Если есть файл для загрузки
                console.log(`Uploading file: ${importFile.name}`);

                // Подготовка формы для загрузки файла
                const formData = new FormData();
                formData.append('file', importFile);

                // Добавляем ZIP-архив, если он выбран
                if (importZipFile) {
                    formData.append('images_zip', importZipFile);
                }

                // Отправляем запрос
                response = await axios.post(
                    `/api/v1/storefronts/import-sources/${sourceId}/run`,
                    formData,
                    {
                        headers: {
                            'Content-Type': 'multipart/form-data'
                        }
                    }
                );
            } else {
                // Запуск импорта по URL
                console.log(`Running import from URL for source ID: ${sourceId}`);
                response = await axios.post(`/api/v1/storefronts/import-sources/${sourceId}/run`);
            }

            console.log("Import response:", response.data);

            // Обновляем источники после импорта
            const sourcesResponse = await axios.get(`/api/v1/storefronts/${id}/import-sources`);
            setImportSources(sourcesResponse.data.data || []);

            // Обновляем историю для этого источника
            await fetchImportHistory(sourceId);

            // Обновляем список объявлений
            fetchData();

            // Сбрасываем файл после успешного импорта
            setImportFile(null);

            // Показываем результат импорта
            const importResult = response.data.data;
            alert(t('marketplace:store.import.completed', {
                status: importResult.status,
                imported: importResult.items_imported,
                total: importResult.items_total
            }));

        } catch (err: any) {
            console.error(`Error running import for source ${sourceId}:`, err);
            alert(t('marketplace:store.import.error', {
                message: err.response?.data?.error || err.message
            }));
        } finally {
            setRunningImport(null);
        }
    };

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
        const file = e.target.files?.[0];
        if (file) {
            setImportFile(file);
        }
    };

    const handleTabChange = (_event: React.SyntheticEvent, newValue: number): void => {
        setActiveTab(newValue);
    };

    // Обработчик изменения режима отображения (плитка/список)
    const handleViewModeChange = (_event: React.MouseEvent<HTMLElement>, newMode: 'grid' | 'list' | null): void => {
        if (newMode !== null) {
            setViewMode(newMode);
        }
    };

    // Обработчик выбора одного объявления
    const handleSelectListing = (listingId: number): void => {
        setSelectedListings(prev => {
            if (prev.includes(listingId)) {
                return prev.filter(id => id !== listingId);
            } else {
                return [...prev, listingId];
            }
        });
    };

    // Обработчик выбора всех объявлений
    const handleSelectAllListings = (checked: boolean): void => {
        if (checked) {
            setSelectedListings(storeListings.map(listing => listing.id));
        } else {
            setSelectedListings([]);
        }
    };

    // Обработчики пагинации
    const handlePageChange = (newPage: number): void => {
        setPage(newPage);
        // При смене страницы сбрасываем выбранные элементы
        setSelectedListings([]);
    };

    const handleLimitChange = (newLimit: number): void => {
        setLimit(newLimit);
        // При изменении лимита возвращаемся на первую страницу и сбрасываем выбор
        setPage(1);
        setSelectedListings([]);
    };

    // Обработчик для группового удаления объявлений
    const handleDeleteSelectedListings = async (): Promise<void> => {
        if (selectedListings.length === 0) return;

        try {
            // Последовательно удаляем каждое объявление
            for (const listingId of selectedListings) {
                await axios.delete(`/api/v1/marketplace/listings/${listingId}`);
            }

            // Обновляем список объявлений через нашу новую функцию
            fetchData();

            // Очищаем выбранные элементы
            setSelectedListings([]);

            // Показываем сообщение об успехе
            setBatchActionSuccess(t('marketplace:store.listings.deleteSuccess', { count: selectedListings.length }));

            // Скрываем сообщение через 3 секунды
            setTimeout(() => setBatchActionSuccess(null), 3000);

            // Закрываем модальное окно подтверждения
            setDeleteConfirmOpen(false);
        } catch (err) {
            console.error('Error deleting listings:', err);
            alert(t('marketplace:store.listings.deleteError'));
            setDeleteConfirmOpen(false);
        }
    };

    // Обработчик для группового перевода объявлений
    const handleTranslateSelectedListings = async (): Promise<void> => {
        if (selectedListings.length === 0) return;

        // Открываем диалог подтверждения оплаты
        setOpenTranslationDialog(true);
    };

    // Обработчик для прямого запуска синхронизации
    const handleRunDirectSync = async (sourceId: number): Promise<void> => {
        try {
            setLoading(true);

            // Находим источник по ID
            const source = importSources.find(s => s.id === sourceId);
            if (!source || !source.url) {
                setError(t('marketplace:store.import.noUrlConfigured'));
                return;
            }

            // Показываем индикатор запуска синхронизации
            setBatchActionSuccess(t('marketplace:store.import.syncStarting', {
                defaultValue: 'Запуск синхронизации...'
            }));

            // Запускаем импорт
            const response = await axios.post(`/api/v1/storefronts/import-sources/${sourceId}/run`);

            // Обновляем данные
            fetchData();
            fetchImportHistory(sourceId);

            // Показываем результат синхронизации
            const result = response.data?.data;
            if (result) {
                setBatchActionSuccess(t('marketplace:store.import.syncCompleted', {
                    status: result.status,
                    imported: result.items_imported || 0,
                    total: result.items_total || 0,
                    defaultValue: `Синхронизация завершена. Статус: ${result.status}. Импортировано: ${result.items_imported || 0}/${result.items_total || 0} товаров.`
                }));
            } else {
                setBatchActionSuccess(t('marketplace:store.import.syncSuccess', {
                    defaultValue: 'Синхронизация успешно запущена'
                }));
            }

            // Убираем сообщение через 5 секунд
            setTimeout(() => setBatchActionSuccess(null), 5000);
        } catch (err: any) {
            console.error('Error running direct sync:', err);
            setError(err.response?.data?.error || t('marketplace:store.import.syncError', {
                defaultValue: 'Ошибка при запуске синхронизации'
            }));

            // Убираем сообщение об ошибке через 5 секунд
            setTimeout(() => setError(null), 5000);
        } finally {
            setLoading(false);
        }
    };
    //  метод для подтверждения и выполнения перевода
    const confirmAndExecuteTranslation = async (): Promise<void> => {
        if (selectedListings.length === 0) return;

        try {
            setTranslationLoading(true);
            setOpenTranslationDialog(false);

            // Используем API-метод для группового перевода
            const response = await axios.post('/api/v1/marketplace/translations/batch', {
                listing_ids: selectedListings,
                target_languages: ['sr', 'en', 'ru']
            });

            console.log('Translation response:', response.data);

            // Обновляем баланс пользователя
            const balanceResponse = await axios.get('/api/v1/balance');
            if (balanceResponse.data?.data?.balance) {
                setUserBalance(balanceResponse.data.data.balance);
            }

            // Показываем сообщение об успехе
            setBatchActionSuccess(
                t('marketplace:store.listings.translateSuccessWithCost', {
                    count: selectedListings.length,
                    cost: response.data.data?.cost || (selectedListings.length * 25),
                    defaultValue: 'Successfully translated {{count}} listings for {{cost}} RSD'
                })
            );

            // Скрываем сообщение через 3 секунды
            setTimeout(() => setBatchActionSuccess(null), 3000);
        } catch (err: any) {
            console.error('Error translating listings:', err);

            let errorMessage = t('marketplace:store.listings.translateError');

            // Проверяем на ошибку недостаточно средств
            if (err.response?.status === 402) {
                errorMessage = t('marketplace:store.listings.insufficientFunds', {
                    required: err.response.data?.error || '',
                    available: userBalance,
                    defaultValue: 'Insufficient funds for translation.'
                });
            }

            alert(errorMessage);
        } finally {
            setTranslationLoading(false);
        }
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

    if (error || !storefront) {
        return (
            <Container maxWidth="lg" sx={{ py: 4 }}>
                <Alert severity="error" sx={{ mb: 3 }}>
                    {error || t('marketplace:store.notFound')}
                </Alert>
                <Button variant="outlined" onClick={() => navigate('/storefronts')}>
                    {t('marketplace:store.backToList')}
                </Button>
            </Container>
        );
    }

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={4}>
                <Typography variant="h4" component="h1">
                    {storefront.name}
                </Typography>
                <Button
                    variant="outlined"
                    startIcon={<Edit />}
                    onClick={() => navigate(`/storefronts/${id}/edit`)}
                >
                    {t('common:buttons.edit')}
                </Button>
            </Box>

            <Paper sx={{ mb: 4 }}>
                <Tabs
                    value={activeTab}
                    onChange={handleTabChange}
                    aria-label="storefront tabs"
                    variant="scrollable"
                    scrollButtons="auto"
                >
                    <Tab label={t('marketplace:store.tabs.import')} icon={<Database />} iconPosition="start" />
                    <Tab label={t('marketplace:store.tabs.settings')} icon={<Settings />} iconPosition="start" />
                    <Tab label={t('marketplace:store.tabs.stats')} icon={<BarChart />} iconPosition="start" />
                </Tabs>

                <Divider />

                <TabPanel value={activeTab} index={0}>
                    <Box mt={3}>
                        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
                            <Typography variant="h6">
                                {t('marketplace:store.storeProducts')}
                            </Typography>

                            <Box display="flex" alignItems="center" gap={2}>
                                {/* Групповые действия */}
                                {selectedListings.length > 0 && (
                                    <ButtonGroup size="small" variant="outlined">
                                        <Tooltip title={t('marketplace:store.listings.translateSelected')}>
                                            <Button
                                                startIcon={<LanguagesIcon />}
                                                onClick={handleTranslateSelectedListings}
                                                disabled={translationLoading}
                                            >
                                                {translationLoading ? (
                                                    <CircularProgress size={20} />
                                                ) : (
                                                    t('marketplace:store.listings.translate')
                                                )}
                                            </Button>
                                        </Tooltip>
                                        <Tooltip title={t('marketplace:store.listings.deleteSelected')}>
                                            <Button
                                                startIcon={<Trash2 />}
                                                color="error"
                                                onClick={() => setDeleteConfirmOpen(true)}
                                            >
                                                {t('marketplace:store.listings.delete')}
                                            </Button>
                                        </Tooltip>
                                    </ButtonGroup>
                                )}

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
                            </Box>
                        </Box>

                        {/* Отображение сообщения об успехе */}
                        {batchActionSuccess && (
                            <Alert
                                icon={<Check />}
                                severity="success"
                                sx={{ mb: 2 }}
                                action={
                                    <IconButton
                                        aria-label="close" color="inherit"
                                        size="small"
                                        onClick={() => setBatchActionSuccess(null)}
                                    >
                                        <Trash2 fontSize="inherit" />
                                    </IconButton>
                                }
                            >
                                {batchActionSuccess}
                            </Alert>
                        )}

                        {storeListings.length === 0 ? (
                            <Alert severity="info" icon={<Info />}>
                                {t('marketplace:store.noProducts')}
                            </Alert>
                        ) : viewMode === 'grid' ? (
                            <>
                                <Grid container spacing={3}>
                                    {storeListings.map((listing) => (
                                        <Grid item xs={12} sm={6} md={4} key={listing.id}>
                                            <Box position="relative">
                                                <Checkbox
                                                    checked={selectedListings.includes(listing.id)}
                                                    onChange={() => handleSelectListing(listing.id)}
                                                    sx={{
                                                        position: 'absolute',
                                                        top: 8,
                                                        left: 8,
                                                        zIndex: 1,
                                                        bgcolor: 'rgba(255,255,255,0.8)',
                                                        borderRadius: '50%',
                                                        '&:hover': { bgcolor: 'rgba(255,255,255,0.9)' }
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

                                {/* Добавляем пагинацию после грида */}
                                <ListingsPagination
                                    totalItems={totalItems}
                                    page={page}
                                    limit={limit}
                                    onPageChange={handlePageChange}
                                    onLimitChange={handleLimitChange}
                                />
                            </>
                        ) : (
                            <>
                                <StorefrontListingsList
                                    listings={storeListings}
                                    selectedItems={selectedListings}
                                    onSelectItem={handleSelectListing}
                                    onSelectAll={handleSelectAllListings}
                                />

                                {/* Добавляем пагинацию после списка */}
                                <ListingsPagination
                                    totalItems={totalItems}
                                    page={page}
                                    limit={limit}
                                    onPageChange={handlePageChange}
                                    onLimitChange={handleLimitChange}
                                />
                            </>
                        )}
                    </Box>

                    <ImportSourceList
                        sources={importSources}
                        storefrontId={Number(id)}
                        onUpdate={fetchData}
                        onDelete={async (sourceId: number) => {
                            try {
                                await axios.delete(`/api/v1/storefronts/import-sources/${sourceId}`);
                                setImportSources(prev => prev.filter(source => source.id !== sourceId));
                            } catch (err) {
                                console.error(`Error deleting import source ${sourceId}:`, err);
                                alert(t('marketplace:store.import.deleteError'));
                            }
                        }}
                        onFetchHistory={fetchImportHistory}
                        onRunDirectSync={handleRunDirectSync}
                    />

                </TabPanel>

                <TabPanel value={activeTab} index={1}>
                    <Box mb={3}>
                        <Typography variant="h6" gutterBottom>
                            {t('marketplace:store.settings.title')}
                        </Typography>
                    </Box>

                    <Grid container spacing={3}>
                        <Grid item xs={12} md={6}>
                            <Card>
                                <CardHeader title={t('marketplace:store.settings.basic')} />
                                <Divider />
                                <CardContent>
                                    <Stack spacing={2}>
                                        <TextField
                                            label={t('marketplace:store.settings.name')}
                                            fullWidth
                                            value={storefront.name}
                                            disabled
                                        />
                                        <TextField
                                            label={t('common:common.description')}
                                            fullWidth
                                            multiline
                                            rows={3}
                                            value={storefront.description || ''}
                                            disabled
                                        />
                                        <TextField
                                            label={t('marketplace:store.settings.slug')}
                                            fullWidth
                                            value={storefront.slug || ''}
                                            disabled
                                        />
                                    </Stack>
                                </CardContent>
                            </Card>
                        </Grid>

                        <Grid item xs={12} md={6}>
                            <Card>
                                <CardHeader title={t('marketplace:store.settings.additional')} />
                                <Divider />
                                <CardContent>
                                    <Typography variant="body1" paragraph>
                                        {t('common:balance.status')}: {storefront.status === 'active'
                                            ? t('marketplace:store.settings.statusActive')
                                            : t('marketplace:store.settings.statusInactive')}
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        {t('marketplace:store.settings.created')}: {new Date(storefront.created_at).toLocaleString()}
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        {t('marketplace:store.settings.updated')}: {new Date(storefront.updated_at).toLocaleString()}
                                    </Typography>
                                </CardContent>
                            </Card>
                        </Grid>
                    </Grid>
                </TabPanel>

                <TabPanel value={activeTab} index={2}>
                    <Box mb={3}>
                        <Typography variant="h6" gutterBottom>
                            {t('marketplace:store.stats.title')}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            {t('marketplace:store.stats.description')}
                        </Typography>
                    </Box>

                    <Alert severity="info">
                        {t('marketplace:store.stats.notAvailable')}
                    </Alert>
                </TabPanel>
            </Paper>

            {/* Модальное окно создания источника импорта */}
            <Modal
                open={openImportModal}
                onClose={() => setOpenImportModal(false)}
                aria-labelledby="import-source-modal"
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
                        {t('marketplace:store.import.createSourceTitle')}
                    </Typography>

                    {importError && (
                        <Alert severity="error" sx={{ mb: 2 }}>
                            {importError}
                        </Alert>
                    )}

                    <Stack spacing={2} mb={3}>
                        <TextField
                            select
                            label={t('marketplace:store.import.typeLabel')}
                            fullWidth
                            required
                            value={importForm.type}
                            onChange={(e: ChangeEvent<HTMLInputElement>) => 
                                setImportForm({ ...importForm, type: e.target.value as 'csv' | 'xml' | 'json' })}
                            disabled={importLoading}
                            SelectProps={{
                                native: true
                            }}
                        >
                            <option value="csv">{t('marketplace:store.import.typeCSV')}</option>
                            <option value="xml">{t('marketplace:store.import.typeXML')}</option>
                            <option value="json">{t('marketplace:store.import.typeJSON')}</option>
                        </TextField>

                        <TextField
                            label={t('marketplace:store.import.urlLabel')}
                            fullWidth
                            value={importForm.url}
                            onChange={(e: ChangeEvent<HTMLInputElement>) => 
                                setImportForm({ ...importForm, url: e.target.value })}
                            disabled={importLoading}
                            placeholder={t('marketplace:store.import.urlPlaceholder')}
                            helperText={t('marketplace:store.import.urlHelp')}
                        />
                    </Stack>

                    <Box display="flex" justifyContent="space-between">
                        <Button
                            variant="outlined"
                            onClick={() => setOpenImportModal(false)}
                            disabled={importLoading}
                        >
                            {t('common:buttons.cancel')}
                        </Button>

                        <Button
                            variant="contained"
                            onClick={handleCreateImportSource}
                            disabled={importLoading}
                            startIcon={importLoading ? <CircularProgress size={20} /> : <Plus />}
                        >
                            {importLoading ? t('marketplace:store.import.creating') : t('marketplace:store.import.createSource')}
                        </Button>
                    </Box>
                </Paper>
            }
            />

            {/* Модальное окно подтверждения удаления выбранных объявлений */}
            <Modal
                open={deleteConfirmOpen}
                onClose={() => setDeleteConfirmOpen(false)}
                aria-labelledby="delete-confirm-modal"
                keepMounted
                children={
                <Paper
                    sx={{
                        position: 'absolute',
                        top: '50%',
                        left: '50%',
                        transform: 'translate(-50%, -50%)',
                        width: { xs: '90%', sm: 400 },
                        p: 3,
                    }}
                >
                    <Typography variant="h6" gutterBottom>
                        {t('marketplace:store.listings.deleteConfirmTitle')}
                    </Typography>

                    <Typography variant="body1" paragraph>
                        {t('marketplace:store.listings.deleteConfirmText', { count: selectedListings.length })}
                    </Typography>

                    <Alert severity="warning" sx={{ mb: 3 }}>
                        {t('marketplace:store.listings.deleteWarning')}
                    </Alert>

                    <Box display="flex" justifyContent="flex-end" gap={2}>
                        <Button variant="outlined" onClick={() => setDeleteConfirmOpen(false)}>
                            {t('common:buttons.cancel')}
                        </Button>
                        <Button
                            variant="contained"
                            color="error"
                            onClick={handleDeleteSelectedListings}
                        >
                            {t('common:buttons.delete')}
                        </Button>
                    </Box>
                </Paper>
            }
            />

            {/* Добавляем плавающую панель для групповых действий */}
            <BatchActionsBar
                selectedItems={selectedListings}
                onClearSelection={() => setSelectedListings([])}
                onDelete={() => setDeleteConfirmOpen(true)}
                onTranslate={handleTranslateSelectedListings}
                isTranslating={translationLoading}
            />
            <TranslationPaymentDialog
                open={openTranslationDialog}
                onClose={() => setOpenTranslationDialog(false)}
                selectedListings={selectedListings}
                balance={userBalance}
                onConfirm={confirmAndExecuteTranslation}
                loading={translationLoading}
                costPerListing={25}
            />
            <ImportModal
                open={openImportModal && selectedSourceId !== null}
                onClose={() => {
                    setOpenImportModal(false);
                    setSelectedSourceId(null);
                }}
                sourceId={selectedSourceId}
                onSuccess={(result) => {
                    // Обновляем источники и список объявлений после успешного импорта
                    fetchData();
                    if (selectedSourceId) {
                        fetchImportHistory(selectedSourceId);
                    }
                }}
            />
        </Container>
    );
};

export default StorefrontDetailPage;