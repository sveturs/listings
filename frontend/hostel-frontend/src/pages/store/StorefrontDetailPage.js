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
    TableRow
} from '@mui/material';
import { Edit, Trash2, Plus, Database, Upload, Settings, BarChart, RefreshCw, FileType, AlertTriangle } from 'lucide-react';
import axios from '../../api/axios';
import { useAuth } from '../../contexts/AuthContext';
import ListingCard from '../../components/marketplace/ListingCard';
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

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                const [storefrontResponse, sourcesResponse, listingsResponse] = await Promise.all([
                    axios.get(`/api/v1/storefronts/${id}`),
                    axios.get(`/api/v1/storefronts/${id}/import-sources`),
                    axios.get('/api/v1/marketplace/listings', { 
                        params: { storefront_id: id } 
                    })
                ]);
                
                setStorefront(storefrontResponse.data.data);
                setImportSources(sourcesResponse.data.data || []);
                setStoreListings(listingsResponse.data.data?.data || []);
            } catch (err) {
                console.error('Error fetching storefront data:', err);
                setError('Не удалось загрузить данные витрины');
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [id]);

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
                setImportError('Выберите тип импорта');
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

            // Покажем успешное сообщение с инструкцией
            alert('Источник импорта успешно создан! Теперь вы можете загрузить файл через кнопку "Запустить импорт".');
        } catch (err) {
            console.error('Error creating import source:', err);
            setImportError(`Не удалось создать источник импорта: ${err.response?.data?.error || err.message}`);
        } finally {
            setImportLoading(false);
        }
    };

    const handleDeleteImportSource = async (sourceId) => {
        if (!window.confirm('Вы уверены, что хотите удалить этот источник импорта?')) {
            return;
        }

        try {
            await axios.delete(`/api/v1/storefronts/import-sources/${sourceId}`);
            setImportSources(prev => prev.filter(source => source.id !== sourceId));
        } catch (err) {
            console.error(`Error deleting import source ${sourceId}:`, err);
            alert('Не удалось удалить источник импорта');
        }
    };

    const handleRunImport = async (sourceId) => {
        try {
            setRunningImport(sourceId);

            console.log(`Starting import for source ID: ${sourceId}, file: ${importFile ? importFile.name : 'none'}`);

            // Если нет ни файла, ни URL, показываем ошибку
            const source = importSources.find(s => s.id === sourceId);
            if (!importFile && (!source || !source.url)) {
                alert("Для импорта нужно либо загрузить файл, либо настроить URL источника");
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

            // Сбрасываем файл после успешного импорта
            setImportFile(null);

            // Показываем результат импорта
            const importResult = response.data.data;
            alert(`Импорт завершен. Статус: ${importResult.status}.\nОбработано: ${importResult.items_imported} из ${importResult.items_total}`);

        } catch (err) {
            console.error(`Error running import for source ${sourceId}:`, err);
            alert(`Ошибка при импорте: ${err.response?.data?.error || err.message}`);
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
                    {error || 'Витрина не найдена'}
                </Alert>
                <Button variant="outlined" onClick={() => navigate('/storefronts')}>
                    Вернуться к списку витрин
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
                    Редактировать
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
                    <Tab label="Импорт товаров" icon={<Database />} iconPosition="start" />
                    <Tab label="Настройки" icon={<Settings />} iconPosition="start" />
                    <Tab label="Статистика" icon={<BarChart />} iconPosition="start" />
                </Tabs>

                <Divider />

                <TabPanel value={activeTab} index={0}>
                    <Box mt={3}>
                        <Typography variant="h6" gutterBottom>
                            Товары в витрине
                        </Typography>

                        <Grid container spacing={3} mt={1}>
                            {loading ? (
                                <Box display="flex" justifyContent="center" width="100%" p={4}>
                                    <CircularProgress />
                                </Box>
                            ) : (
                                <>
                                    {storeListings && storeListings.length > 0 ? (
                                        storeListings.map((listing) => (
                                            <Grid item xs={12} sm={6} md={4} key={listing.id}>
                                                <Link
                                                    to={`/marketplace/listings/${listing.id}`}
                                                    style={{ textDecoration: 'none' }}
                                                >
                                                    <ListingCard listing={listing} />
                                                </Link>
                                            </Grid>
                                        ))
                                    ) : (
                                        <Box width="100%" p={4} textAlign="center">
                                            <Typography>В этой витрине пока нет товаров</Typography>
                                        </Box>
                                    )}
                                </>
                            )}
                        </Grid>
                    </Box>
                    <Box mb={3}>
                        <Typography variant="h6" gutterBottom>
                            Источники импорта
                        </Typography>
                        <Typography variant="body2" color="text.secondary" paragraph>
                            Создайте источник данных для импорта товаров в вашу витрину.
                            Поддерживаются форматы CSV, XML и JSON.
                        </Typography>
                    </Box>

                    <Box mb={2} display="flex" justifyContent="flex-end">
                        <Button
                            variant="contained"
                            startIcon={<Plus />}
                            onClick={() => setOpenImportModal(true)}
                        >
                            Добавить источник
                        </Button>
                    </Box>

                    {importSources.length === 0 ? (
                        <Paper sx={{ p: 4, textAlign: 'center' }}>
                            <Database size={64} stroke={1} style={{ margin: '20px auto', opacity: 0.5 }} />
                            <Typography variant="h6" gutterBottom>
                                Источники импорта не созданы
                            </Typography>
                            <Typography variant="body1" color="text.secondary" paragraph>
                                Создайте источник импорта, чтобы загружать товары в свою витрину
                            </Typography>
                            <Button
                                variant="contained"
                                startIcon={<Plus />}
                                onClick={() => setOpenImportModal(true)}
                            >
                                Добавить источник
                            </Button>
                        </Paper>
                    ) : (
                        <Grid container spacing={3}>
                            {importSources.map((source) => (
                                <Grid item xs={12} md={6} key={source.id}>
                                    <Card sx={{ height: '100%' }}>
                                        <CardHeader
                                            title={`Источник: ${source.type.toUpperCase()}`}
                                            subheader={`Создан: ${new Date(source.created_at).toLocaleDateString()}`}
                                            action={
                                                <IconButton onClick={() => handleDeleteImportSource(source.id)}>
                                                    <Trash2 size={18} />
                                                </IconButton>
                                            }
                                        />
                                        <Divider />
                                        <CardContent>
                                            <Typography variant="subtitle1" gutterBottom>
                                                {source.url ? `URL: ${source.url}` : 'Ручная загрузка файлов'}
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary" paragraph>
                                                Последний импорт: {source.last_import_at ?
                                                    new Date(source.last_import_at).toLocaleString() : 'Никогда'}
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                Статус: {source.last_import_status || 'Не запускался'}
                                            </Typography>

                                            <Box sx={{ mt: 2 }}>
                                                <input
                                                    type="file"
                                                    id={`file-upload-${source.id}`}
                                                    accept=".csv,.xml,.json"
                                                    style={{ display: 'none' }}
                                                    onChange={(e) => {
                                                        const files = e.target.files;
                                                        if (files && files.length > 0) {
                                                            setImportFile(files[0]);
                                                            // Не запускаем импорт автоматически, даем пользователю нажать кнопку
                                                        }
                                                    }}
                                                />

                                                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                                    <Button
                                                        variant="outlined"
                                                        size="small"
                                                        startIcon={<BarChart />}
                                                        onClick={() => fetchImportHistory(source.id)}
                                                    >
                                                        История
                                                    </Button>

                                                    <Box sx={{ display: 'flex', gap: 1 }}>
                                                        <Button
                                                            variant="outlined"
                                                            size="small"
                                                            component="label"
                                                            htmlFor={`file-upload-${source.id}`}
                                                            startIcon={<Upload />}
                                                        >
                                                            Выбрать файл
                                                        </Button>

                                                        <Button
                                                            variant="contained"
                                                            size="small"
                                                            startIcon={runningImport === source.id ? <CircularProgress size={18} /> : <RefreshCw />}
                                                            onClick={() => handleRunImport(source.id)}
                                                            disabled={runningImport === source.id}
                                                        >
                                                            {runningImport === source.id ? 'Импорт...' : importFile ? 'Загрузить и импорт.' : 'Запустить'}
                                                        </Button>
                                                    </Box>
                                                </Box>

                                                {/* Показываем выбранный файл */}
                                                {importFile && (
                                                    <Box sx={{ mt: 1, p: 1, bgcolor: 'background.paper', borderRadius: 1 }}>
                                                        <Typography variant="body2" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                                            <FileType size={16} />
                                                            {importFile.name} ({Math.round(importFile.size / 1024)} КБ)
                                                        </Typography>
                                                    </Box>
                                                )}
                                            </Box>

                                            {importHistories[source.id] && importHistories[source.id].length > 0 && (
                                                <Box mt={3}>
                                                    <Typography variant="subtitle2" gutterBottom>
                                                        История импорта
                                                    </Typography>
                                                    <TableContainer component={Paper} variant="outlined">
                                                        <Table size="small">
                                                            <TableHead>
                                                                <TableRow>
                                                                    <TableCell>Дата</TableCell>
                                                                    <TableCell>Статус</TableCell>
                                                                    <TableCell align="right">Импортировано</TableCell>
                                                                </TableRow>
                                                            </TableHead>
                                                            <TableBody>
                                                                {importHistories[source.id].map((history) => (
                                                                    <TableRow key={history.id}>
                                                                        <TableCell>
                                                                            {new Date(history.started_at).toLocaleDateString()}
                                                                        </TableCell>
                                                                        <TableCell>
                                                                            {history.status === 'success' ? 'Успешно' :
                                                                                history.status === 'partial' ? 'Частично' : 'Ошибка'}
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
                            Настройки витрины
                        </Typography>
                    </Box>

                    <Grid container spacing={3}>
                        <Grid item xs={12} md={6}>
                            <Card>
                                <CardHeader title="Основные настройки" />
                                <Divider />
                                <CardContent>
                                    <Stack spacing={2}>
                                        <TextField
                                            label="Название витрины"
                                            fullWidth
                                            value={storefront.name}
                                            disabled
                                        />
                                        <TextField
                                            label="Описание"
                                            fullWidth
                                            multiline
                                            rows={3}
                                            value={storefront.description || ''}
                                            disabled
                                        />
                                        <TextField
                                            label="Slug (URL)"
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
                                <CardHeader title="Дополнительные настройки" />
                                <Divider />
                                <CardContent>
                                    <Typography variant="body1" paragraph>
                                        Статус: {storefront.status === 'active' ? 'Активна' : 'Не активна'}
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Создана: {new Date(storefront.created_at).toLocaleString()}
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Обновлена: {new Date(storefront.updated_at).toLocaleString()}
                                    </Typography>
                                </CardContent>
                            </Card>
                        </Grid>
                    </Grid>
                </TabPanel>

                <TabPanel value={activeTab} index={2}>
                    <Box mb={3}>
                        <Typography variant="h6" gutterBottom>
                            Статистика витрины
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Здесь будет отображаться статистика по продажам и посещаемости вашей витрины.
                        </Typography>
                    </Box>

                    <Alert severity="info">
                        Статистика будет доступна после добавления товаров в витрину.
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
                        Создание источника импорта
                    </Typography>

                    {importError && (
                        <Alert severity="error" sx={{ mb: 2 }}>
                            {importError}
                        </Alert>
                    )}

                    <Stack spacing={2} mb={3}>
                        <TextField
                            select
                            label="Тип импорта"
                            fullWidth
                            required
                            value={importForm.type}
                            onChange={(e) => setImportForm({ ...importForm, type: e.target.value })}
                            disabled={importLoading}
                            SelectProps={{
                                native: true
                            }}
                        >
                            <option value="csv">CSV - разделенные запятыми/точкой с запятой значения</option>
                            <option value="xml">XML - структурированные данные</option>
                            <option value="json">JSON - JavaScript Object Notation</option>
                        </TextField>

                        <TextField
                            label="URL источника (необязательно)"
                            fullWidth
                            value={importForm.url}
                            onChange={(e) => setImportForm({ ...importForm, url: e.target.value })}
                            disabled={importLoading}
                            placeholder="Можно оставить пустым для ручной загрузки файлов"
                            helperText="Укажите URL, если хотите импортировать данные из внешнего источника"
                        />
                    </Stack>

                    <Box display="flex" justifyContent="space-between">
                        <Button
                            variant="outlined"
                            onClick={() => setOpenImportModal(false)}
                            disabled={importLoading}
                        >
                            Отмена
                        </Button>

                        <Button
                            variant="contained"
                            onClick={handleCreateImportSource}
                            disabled={importLoading}
                            startIcon={importLoading ? <CircularProgress size={20} /> : <Plus />}
                        >
                            {importLoading ? 'Создание...' : 'Создать источник'}
                        </Button>
                    </Box>
                </Paper>
            </Modal>
        </Container>
    );
};

export default StorefrontDetailPage;