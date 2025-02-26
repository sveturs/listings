import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams, useNavigate } from 'react-router-dom';
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
          const [storefrontResponse, sourcesResponse] = await Promise.all([
            axios.get(`/api/v1/storefronts/${id}`),
            axios.get(`/api/v1/storefronts/${id}/import-sources`)
          ]);
          
          setStorefront(storefrontResponse.data.data);
          setImportSources(sourcesResponse.data.data || []);
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
        
        const response = await axios.post('/api/v1/storefronts/import-sources', importForm);
        setImportSources(prev => [...prev, response.data.data]);
        setOpenImportModal(false);
        setImportForm({
          type: 'csv',
          url: '',
          storefront_id: Number(id)
        });
        setImportFile(null);
      } catch (err) {
        console.error('Error creating import source:', err);
        setImportError('Не удалось создать источник импорта');
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
        
        let response;
        if (importFile && sourceId === 'new') {
          // Если это новый файл для импорта
          const formData = new FormData();
          formData.append('file', importFile);
          
          response = await axios.post(`/api/v1/storefronts/import-sources/${importSources[0].id}/run`, formData, {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          });
        } else {
          // Если это запуск существующего источника
          response = await axios.post(`/api/v1/storefronts/import-sources/${sourceId}/run`);
        }
        
        // Обновляем источники после импорта
        const sourcesResponse = await axios.get(`/api/v1/storefronts/${id}/import-sources`);
        setImportSources(sourcesResponse.data.data || []);
        
        // Обновляем историю для этого источника
        await fetchImportHistory(sourceId);
        
        if (sourceId === 'new') {
          setOpenImportModal(false);
          setImportFile(null);
        }
      } catch (err) {
        console.error(`Error running import for source ${sourceId}:`, err);
        alert('Не удалось запустить импорт');
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
                    <Card>
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
                          URL: {source.url || 'Не указан'}
                        </Typography>
                        <Typography variant="body2" color="text.secondary" paragraph>
                          Последний импорт: {source.last_import_at ? 
                            new Date(source.last_import_at).toLocaleString() : 'Никогда'}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          Статус: {source.last_import_status || 'Не запускался'}
                        </Typography>
                        
                        <Box display="flex" justifyContent="space-between" mt={2}>
                          <Button
                            variant="outlined"
                            size="small"
                            startIcon={<BarChart />}
                            onClick={() => fetchImportHistory(source.id)}
                          >
                            История
                          </Button>
                          <Button
                            variant="contained"
                            size="small"
                            startIcon={runningImport === source.id ? <CircularProgress size={18} /> : <RefreshCw />}
                            onClick={() => handleRunImport(source.id)}
                            disabled={runningImport === source.id}
                          >
                            {runningImport === source.id ? 'Импорт...' : 'Запустить импорт'}
                          </Button>
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
                <option value="csv">CSV - разделенные запятыми значения</option>
                <option value="xml">XML - структурированные данные</option>
                <option value="json">JSON - JavaScript Object Notation</option>
              </TextField>
              
              <TextField
                label="URL источника (необязательно)"
                fullWidth
                value={importForm.url}
                onChange={(e) => setImportForm({ ...importForm, url: e.target.value })}
                disabled={importLoading}
                placeholder="https://example.com/products.csv"
                helperText="Укажите URL, если хотите импортировать данные из внешнего источника"
              />
              
              <Divider sx={{ my: 1 }}>
                <Typography variant="body2" color="text.secondary">или</Typography>
              </Divider>
              
              <Box>
                <Button
                  variant="outlined"
                  component="label"
                  startIcon={<Upload />}
                  fullWidth
                >
                  Загрузить файл
                  <input
                    type="file"
                    hidden
                    onChange={handleFileChange}
                    accept=".csv,.xml,.json"
                  />
                </Button>
                {importFile && (
                  <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                    Выбран файл: {importFile.name}
                  </Typography>
                )}
              </Box>
            </Stack>
            
            <Box display="flex" justifyContent="space-between">
              <Button
                variant="outlined"
                onClick={() => setOpenImportModal(false)}
                disabled={importLoading}
              >
                Отмена
              </Button>
              
              {importFile ? (
                <Button
                  variant="contained"
                  onClick={() => handleRunImport('new')}
                  disabled={importLoading || !importFile || importSources.length === 0}
                  startIcon={importLoading ? <CircularProgress size={20} /> : <Upload />}
                >
                  {importLoading ? 'Загрузка...' : 'Загрузить и импортировать'}
                </Button>
              ) : (
                <Button
                  variant="contained"
                  onClick={handleCreateImportSource}
                  disabled={importLoading || (!importForm.url && !importFile)}
                  startIcon={importLoading ? <CircularProgress size={20} /> : <Plus />}
                >
                  {importLoading ? 'Создание...' : 'Создать источник'}
                </Button>
              )}
            </Box>
            
            {importSources.length === 0 && importFile && (
              <Alert severity="warning" sx={{ mt: 2 }} icon={<AlertTriangle />}>
                Сначала создайте источник импорта без файла, затем вы сможете загрузить файл.
              </Alert>
            )}
          </Paper>
        </Modal>
      </Container>
    );
  };
  
  export default StorefrontDetailPage;