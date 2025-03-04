// frontend/hostel-frontend/src/pages/store/StorefrontDetailPage.js
import React, { useState, useEffect } from 'react';
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

const TabPanel = (props) => {
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

const StorefrontDetailPage = () => {
    const { t } = useTranslation(['common', 'marketplace']);
    const [storeListings, setStoreListings] = useState([]);
    const navigate = useNavigate();
    const { id } = useParams();
    const { user } = useAuth();

    const [loading, setLoading] = useState(true);
    const [storefront, setStorefront] = useState(null);
    const [error, setError] = useState(null);
    const [importSources, setImportSources] = useState([]);
    const [activeTab, setActiveTab] = useState(0);
    const [importHistories, setImportHistories] = useState({});
    const [openImportModal, setOpenImportModal] = useState(false);
    const [importForm, setImportForm] = useState({
        type: 'csv',
        url: '',
        storefront_id: Number(id)
    });
    const [importFile, setImportFile] = useState(null);
    const [importError, setImportError] = useState(null);
    const [importLoading, setImportLoading] = useState(false);
    const [runningImport, setRunningImport] = useState(null);

    // Новые состояния для списочного режима и выбора элементов
    const [viewMode, setViewMode] = useState('grid');
    const [selectedListings, setSelectedListings] = useState([]);
    const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);
    const [translationLoading, setTranslationLoading] = useState(false);
    const [batchActionSuccess, setBatchActionSuccess] = useState(null);
    
    // Состояния для пагинации
    const [page, setPage] = useState(1);
    const [limit, setLimit] = useState(20);
    const [totalItems, setTotalItems] = useState(0);

    // Функция для загрузки данных с учетом пагинации
    const fetchData = async () => {
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
        fetchData();
    }, [id, page, limit, t]);

    const fetchImportHistory = async (sourceId) => {
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

    const handleCreateImportSource = async () => {
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
        } catch (err) {
            console.error('Error creating import source:', err);
            setImportError(`${t('marketplace:store.import.createError')}: ${err.response?.data?.error || err.message}`);
        } finally {
            setImportLoading(false);
        }
    };

    const handleDeleteImportSource = async (sourceId) => {
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

    const handleRunImport = async (sourceId) => {
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

        } catch (err) {
            console.error(`Error running import for source ${sourceId}:`, err);
            alert(t('marketplace:store.import.error', {
                message: err.response?.data?.error || err.message
            }));
        } finally {
            setRunningImport(null);
        }
    };

    const handleFileChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            setImportFile(file);
        }
    };

    const handleTabChange = (event, newValue) => {
        setActiveTab(newValue);
    };

    // Обработчик изменения режима отображения (плитка/список)
    const handleViewModeChange = (event, newMode) => {
        if (newMode !== null) {
            setViewMode(newMode);
        }
    };

    // Обработчик выбора одного объявления
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
    const handleSelectAllListings = (checked) => {
        if (checked) {
            setSelectedListings(storeListings.map(listing => listing.id));
        } else {
            setSelectedListings([]);
        }
    };

    // Обработчики пагинации
    const handlePageChange = (newPage) => {
        setPage(newPage);
        // При смене страницы сбрасываем выбранные элементы
        setSelectedListings([]);
    };

    const handleLimitChange = (newLimit) => {
        setLimit(newLimit);
        // При изменении лимита возвращаемся на первую страницу и сбрасываем выбор
        setPage(1);
        setSelectedListings([]);
    };

    // Обработчик для группового удаления объявлений
    const handleDeleteSelectedListings = async () => {
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
    const handleTranslateSelectedListings = async () => {
        if (selectedListings.length === 0) return;

        try {
            setTranslationLoading(true);

            // Используем новый API-метод для группового перевода
            const response = await axios.post('/api/v1/marketplace/translations/batch', {
                listing_ids: selectedListings,
                target_languages: ['sr', 'en', 'ru']
            });

            console.log('Translation response:', response.data);

            // Показываем сообщение об успехе
            setBatchActionSuccess(t('marketplace:store.listings.translateSuccess', { count: selectedListings.length }));

            // Скрываем сообщение через 3 секунды
            setTimeout(() => setBatchActionSuccess(null), 3000);
        } catch (err) {
            console.error('Error translating listings:', err);
            alert(t('marketplace:store.listings.translateError'));
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
                    
                    <Box mb={3} mt={5}>
                        <Typography variant="h6" gutterBottom>
                            {t('marketplace:store.import.sources')}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" paragraph>
                            {t('marketplace:store.import.createSourceDescription')}
                        </Typography>
                    </Box>

                    <Box mb={2} display="flex" justifyContent="flex-end">
                        <Button
                            variant="contained"
                            startIcon={<Plus />}
                            onClick={() => setOpenImportModal(true)}
                        >
                            {t('marketplace:store.import.addSource')}
                        </Button>
                    </Box>

                    {importSources.length === 0 ? (
                        <Paper sx={{ p: 4, textAlign: 'center' }}>
                            <Database size={64} stroke={1} style={{ margin: '20px auto', opacity: 0.5 }} />
                            <Typography variant="h6" gutterBottom>
                                {t('marketplace:store.import.noSources')}
                            </Typography>
                            <Typography variant="body1" color="text.secondary" paragraph>
                                {t('marketplace:store.import.createFirstSource')}
                            </Typography>
                            <Button
                                variant="contained"
                                startIcon={<Plus />}
                                onClick={() => setOpenImportModal(true)}
                            >
                                {t('marketplace:store.import.addSource')}
                            </Button>
                        </Paper>
                    ) : (
                        <Grid container spacing={3}>
                            {importSources.map((source) => (
                                <Grid item xs={12} md={6} key={source.id}>
                                    <Card sx={{ height: '100%' }}>
                                        <CardHeader
                                            title={t('marketplace:store.import.sourceType', { type: source.type.toUpperCase() })}
                                            subheader={t('marketplace:store.import.created', { date: new Date(source.created_at).toLocaleDateString() })}
                                            action={
                                                <IconButton onClick={() => handleDeleteImportSource(source.id)}>
                                                    <Trash2 size={18} />
                                                </IconButton>
                                            }
                                        />
                                        <Divider />
                                        <CardContent>
                                            <Typography variant="subtitle1" gutterBottom>
                                                {source.url 
                                                    ? t('marketplace:store.import.urlSource', { url: source.url }) 
                                                    : t('marketplace:store.import.manualUpload')
                                                }
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary" paragraph>
                                                {t('marketplace:store.import.lastImport')}: {source.last_import_at ?
                                                    new Date(source.last_import_at).toLocaleString() : t('marketplace:store.import.never')}
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                {t('common:balance.status')}: {source.last_import_status || t('marketplace:store.import.notRun')}
                                            </Typography>

                                            <Box sx={{ mt: 2 }}>
                                                <input
                                                    type="file"
                                                    id={`file-upload-${source.id}`}
                                                    accept=".csv,.xml,.json"
                                                    style={{ display: 'none' }}
                                                    onChange={handleFileChange}
                                                />

                                                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                                    <Button
                                                        variant="outlined"
                                                        size="small"
                                                        startIcon={<BarChart />}
                                                        onClick={() => fetchImportHistory(source.id)}
                                                    >
                                                        {t('marketplace:store.import.history')}
                                                    </Button>

                                                    <Box sx={{ display: 'flex', gap: 1 }}>
                                                        <Button
                                                            variant="outlined"
                                                            size="small"
                                                            component="label"
                                                            htmlFor={`file-upload-${source.id}`}
                                                            startIcon={<Upload />}
                                                        >
                                                            {t('marketplace:store.import.selectFile')}
                                                        </Button>

                                                        <Button
                                                            variant="contained"
                                                            size="small"
                                                            startIcon={runningImport === source.id ? <CircularProgress size={18} /> : <RefreshCw />}
                                                            onClick={() => handleRunImport(source.id)}
                                                            disabled={runningImport === source.id}
                                                        >
                                                            {runningImport === source.id 
                                                                ? t('marketplace:store.import.importing') 
                                                                : importFile 
                                                                    ? t('marketplace:store.import.uploadAndImport') 
                                                                    : t('marketplace:store.import.runImport')}
                                                        </Button>
                                                    </Box>
                                                </Box>

                                                {/* Показываем выбранный файл */}
                                                {importFile && (
                                                    <Box sx={{ mt: 1, p: 1, bgcolor: 'background.paper', borderRadius: 1 }}>
                                                        <Typography variant="body2" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                            <FileType size={16} />
                                                            {importFile.name} ({Math.round(importFile.size / 1024)} {t('marketplace:store.import.kb')})
                                                        </Typography>
                                                    </Box>
                                                )}
                                            </Box>

                                            {importHistories[source.id] && importHistories[source.id].length > 0 && (
                                                <Box mt={3}>
                                                    <Typography variant="subtitle2" gutterBottom>
                                                        {t('marketplace:store.import.historyTitle')}
                                                    </Typography>
                                                    <TableContainer component={Paper} variant="outlined">
                                                        <Table size="small">
                                                            <TableHead>
                                                                <TableRow>
                                                                    <TableCell>{t('common:balance.date')}</TableCell>
                                                                    <TableCell>{t('common:balance.status')}</TableCell>
                                                                    <TableCell align="right">{t('marketplace:store.import.imported')}</TableCell>
                                                                </TableRow>
                                                            </TableHead>
                                                            <TableBody>
                                                                {importHistories[source.id].map((history) => (
                                                                    <TableRow key={history.id}>
                                                                        <TableCell>
                                                                            {new Date(history.started_at).toLocaleDateString()}
                                                                        </TableCell>
                                                                        <TableCell>
                                                                            {history.status === 'success' 
                                                                                ? t('marketplace:store.import.statusSuccess')
                                                                                : history.status === 'partial' 
                                                                                    ? t('marketplace:store.import.statusPartial') 
                                                                                    : t('marketplace:store.import.statusError')}
                                                                        </TableCell>
                                                                        <TableCell align="right">
                                                                            {history.items_imported}/{history.items_total}
                                                                        </TableCell>
                                                                    </TableRow>
                                                                ))}
                                                            </TableBody>
                                                        </Table>
                                                    </TableContainer>
                                                </Box>
                                            )}
                                        </CardContent>
                                    </Card>
                                </Grid>
                            ))}
                        </Grid>
                    )}
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
                            onChange={(e) => setImportForm({ ...importForm, type: e.target.value })}
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
                            onChange={(e) => setImportForm({ ...importForm, url: e.target.value })}
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
            </Modal>

            {/* Модальное окно подтверждения удаления выбранных объявлений */}
            <Modal
                open={deleteConfirmOpen}
                onClose={() => setDeleteConfirmOpen(false)}
                aria-labelledby="delete-confirm-modal"
            >
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
            </Modal>
            
            {/* Добавляем плавающую панель для групповых действий */}
            <BatchActionsBar 
                selectedItems={selectedListings}
                onClearSelection={() => setSelectedListings([])}
                onDelete={() => setDeleteConfirmOpen(true)}
                onTranslate={handleTranslateSelectedListings}
                isTranslating={translationLoading}
            />
        </Container>
    );
};

export default StorefrontDetailPage;