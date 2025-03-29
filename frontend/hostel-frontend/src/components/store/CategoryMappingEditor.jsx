// frontend/hostel-frontend/src/components/store/CategoryMappingEditor.jsx
import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Box,
  Typography,
  Paper,
  Button,
  TextField,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Divider,
  CircularProgress,
  Alert,
  Autocomplete,
  InputAdornment,
  MenuItem,
  Select,
  FormControl,
  InputLabel
} from '@mui/material';
import {
  Plus,
  Edit2,
  Trash2,
  Save,
  X,
  FileType,
  Search
} from 'lucide-react';
import axios from '../../api/axios';

const CategoryMappingEditor = ({ sourceId, onClose, onSave }) => {
  const { t } = useTranslation(['marketplace', 'common']);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [mappings, setMappings] = useState({});
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);
  const [categories, setCategories] = useState([]);
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [currentMapping, setCurrentMapping] = useState({ source: '', target: 0 });
  
  // Добавляем состояние для отслеживания импортированных категорий
  const [importedCategories, setImportedCategories] = useState([]);
  const [importedCategoriesLoading, setImportedCategoriesLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [filteredCategories, setFilteredCategories] = useState([]);

  // Загружаем существующие сопоставления и категории системы
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);
      
      try {
        // Получаем существующие сопоставления
        const mappingsResponse = await axios.get(`/api/v1/storefronts/import-sources/${sourceId}/category-mappings`);
        
        // Получаем категории системы
        const categoriesResponse = await axios.get('/api/v1/marketplace/categories');
        
        // Загружаем импортированные категории
        await fetchImportedCategories();
        
        if (mappingsResponse.data.success) {
          setMappings(mappingsResponse.data.data || {});
        }
        
        if (categoriesResponse.data.success) {
          // Сортируем категории по имени для удобства
          const sortedCategories = [...categoriesResponse.data.data].sort((a, b) => 
            a.name.localeCompare(b.name)
          );
          setCategories(sortedCategories);
          setFilteredCategories(sortedCategories);
        }
      } catch (err) {
        console.error('Error fetching data:', err);
        setError(t('marketplace:store.categoryMapping.fetchError'));
      } finally {
        setLoading(false);
      }
    };
    
    fetchData();
  }, [sourceId, t]);
  
  // Получаем список импортированных категорий
  const fetchImportedCategories = async () => {
    setImportedCategoriesLoading(true);
    try {
      const response = await axios.get(`/api/v1/storefronts/import-sources/${sourceId}/imported-categories`);
      if (response.data.success) {
        setImportedCategories(response.data.data || []);
      }
    } catch (err) {
      console.error('Error fetching imported categories:', err);
      // Не показываем ошибку пользователю, просто логируем
    } finally {
      setImportedCategoriesLoading(false);
    }
  };
  
  // Фильтрация категорий по поисковому запросу
  useEffect(() => {
    if (searchTerm.trim() === '') {
      setFilteredCategories(categories);
    } else {
      const lowerCaseSearchTerm = searchTerm.toLowerCase();
      const filtered = categories.filter(category => 
        category.name.toLowerCase().includes(lowerCaseSearchTerm) ||
        (category.translations && 
         Object.values(category.translations).some(
           translation => translation.toLowerCase().includes(lowerCaseSearchTerm)
         ))
      );
      setFilteredCategories(filtered);
    }
  }, [searchTerm, categories]);

  // Сохранение сопоставлений
  const handleSave = async () => {
    setSaving(true);
    setError(null);
    setSuccess(false);
    
    try {
      const response = await axios.put(
        `/api/v1/storefronts/import-sources/${sourceId}/category-mappings`,
        mappings
      );
      
      if (response.data.success) {
        setSuccess(true);
        if (onSave) {
          onSave(mappings);
        }
        
        // Закрываем окно через 2 секунды после успешного сохранения
        setTimeout(() => {
          onClose();
        }, 2000);
      }
    } catch (err) {
      console.error('Error saving mappings:', err);
      setError(t('marketplace:store.categoryMapping.saveError'));
    } finally {
      setSaving(false);
    }
  };

  // Открытие диалога для добавления нового сопоставления
  const handleAddMapping = () => {
    setCurrentMapping({ source: '', target: 0 });
    setEditDialogOpen(true);
  };

  // Открытие диалога для редактирования существующего сопоставления
  const handleEditMapping = (sourceCategory) => {
    setCurrentMapping({
      source: sourceCategory,
      target: mappings[sourceCategory] || 0
    });
    setEditDialogOpen(true);
  };

  // Удаление сопоставления
  const handleDeleteMapping = (sourceCategory) => {
    const newMappings = { ...mappings };
    delete newMappings[sourceCategory];
    setMappings(newMappings);
  };

  // Сохранение изменений в диалоге
  const handleSaveMapping = () => {
    if (!currentMapping.source.trim()) {
      return;
    }
    
    const newMappings = { ...mappings };
    newMappings[currentMapping.source] = currentMapping.target;
    setMappings(newMappings);
    setEditDialogOpen(false);
  };

  // Форматирование названия категории для отображения
  const formatCategoryName = (categoryId) => {
    const category = categories.find(c => c.id === categoryId);
    return category ? category.name : t('marketplace:store.categoryMapping.unknownCategory');
  };

  // Отображение загрузки
  if (loading) {
    return (
      <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" p={4}>
        <CircularProgress />
        <Typography variant="body1" sx={{ mt: 2 }}>
          {t('common:common.loading')}
        </Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ width: '100%', maxWidth: 800, mx: 'auto', p: 2 }}>
      <Typography variant="h5" gutterBottom>
        {t('marketplace:store.categoryMapping.title')}
      </Typography>
      
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      
      {success && (
        <Alert severity="success" sx={{ mb: 2 }}>
          {t('marketplace:store.categoryMapping.saveSuccess')}
        </Alert>
      )}
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <Typography variant="subtitle1" gutterBottom>
          {t('marketplace:store.categoryMapping.description')}
        </Typography>
        
        <Alert severity="info" sx={{ mb: 2 }}>
          {t('marketplace:store.categoryMapping.info')}
        </Alert>
        
        {/* Поиск по категориям */}
        <TextField
          fullWidth
          variant="outlined"
          label={t('marketplace:store.categoryMapping.searchCategories')}
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          sx={{ mb: 2 }}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Search size={20} />
              </InputAdornment>
            ),
            endAdornment: searchTerm ? (
              <InputAdornment position="end">
                <IconButton
                  aria-label="clear search"
                  onClick={() => setSearchTerm('')}
                  edge="end"
                >
                  <X size={16} />
                </IconButton>
              </InputAdornment>
            ) : null
          }}
        />
        
        {/* Текущие сопоставления */}
        <Typography variant="subtitle2" gutterBottom sx={{ mt: 2 }}>
          {t('marketplace:store.categoryMapping.currentMappings')}
        </Typography>
        
        {Object.keys(mappings).length === 0 ? (
          <Alert severity="warning" sx={{ mb: 2 }}>
            {t('marketplace:store.categoryMapping.noMappings')}
          </Alert>
        ) : (
          <List>
            {Object.entries(mappings).map(([sourceCategory, targetCategoryId]) => (
              <React.Fragment key={sourceCategory}>
                <ListItem>
                  <ListItemText 
                    primary={sourceCategory} 
                    secondary={`→ ${formatCategoryName(targetCategoryId)} (ID: ${targetCategoryId})`}
                  />
                  <ListItemSecondaryAction>
                    <IconButton 
                      edge="end" 
                      aria-label="edit"
                      onClick={() => handleEditMapping(sourceCategory)}
                    >
                      <Edit2 size={18} />
                    </IconButton>
                    <IconButton 
                      edge="end" 
                      aria-label="delete"
                      onClick={() => handleDeleteMapping(sourceCategory)}
                    >
                      <Trash2 size={18} />
                    </IconButton>
                  </ListItemSecondaryAction>
                </ListItem>
                <Divider />
              </React.Fragment>
            ))}
          </List>
        )}
        
        {/* Импортированные категории */}
        <Box sx={{ mt: 3 }}>
          <Typography variant="subtitle2" gutterBottom>
            {t('marketplace:store.categoryMapping.importedCategories')}
          </Typography>
          
          {importedCategoriesLoading ? (
            <Box display="flex" alignItems="center" gap={1}>
              <CircularProgress size={20} />
              <Typography variant="body2">
                {t('marketplace:store.categoryMapping.loadingImportedCategories')}
              </Typography>
            </Box>
          ) : importedCategories.length === 0 ? (
            <Alert severity="info">
              {t('marketplace:store.categoryMapping.noImportedCategories')}
            </Alert>
          ) : (
            <List dense sx={{ maxHeight: 250, overflow: 'auto', bgcolor: 'background.paper' }}>
              {importedCategories.map((category, index) => (
                <ListItem key={index} divider>
                  <ListItemText
                    primary={category}
                    secondary={
                      mappings[category] ? 
                      `→ ${formatCategoryName(mappings[category])} (ID: ${mappings[category]})` :
                      t('marketplace:store.categoryMapping.notMappedYet')
                    }
                  />
                  <ListItemSecondaryAction>
                    {mappings[category] ? (
                      <IconButton 
                        edge="end" 
                        aria-label="edit"
                        onClick={() => handleEditMapping(category)}
                      >
                        <Edit2 size={18} />
                      </IconButton>
                    ) : (
                      <Button
                        size="small"
                        onClick={() => handleEditMapping(category)}
                        startIcon={<Plus size={16} />}
                      >
                        {t('marketplace:store.categoryMapping.mapCategory')}
                      </Button>
                    )}
                  </ListItemSecondaryAction>
                </ListItem>
              ))}
            </List>
          )}
        </Box>
        
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 3 }}>
          <Button
            variant="outlined"
            startIcon={<Plus />}
            onClick={handleAddMapping}
          >
            {t('marketplace:store.categoryMapping.addMapping')}
          </Button>
          
          <Box>
            <Button
              variant="outlined"
              onClick={onClose}
              sx={{ mr: 1 }}
              disabled={saving}
            >
              {t('common:buttons.cancel')}
            </Button>
            
            <Button
              variant="contained"
              startIcon={saving ? <CircularProgress size={20} /> : <Save />}
              onClick={handleSave}
              disabled={saving}
            >
              {t('common:buttons.save')}
            </Button>
          </Box>
        </Box>
      </Paper>
      
      {/* Диалог добавления/редактирования сопоставления */}
      <Dialog 
        open={editDialogOpen} 
        onClose={() => setEditDialogOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>
          {currentMapping.source ? 
            t('marketplace:store.categoryMapping.editMappingTitle') :
            t('marketplace:store.categoryMapping.addMappingTitle')}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 1 }}>
            <TextField
              fullWidth
              label={t('marketplace:store.categoryMapping.sourceCategory')}
              value={currentMapping.source}
              onChange={(e) => setCurrentMapping({...currentMapping, source: e.target.value})}
              margin="dense"
              helperText={t('marketplace:store.categoryMapping.sourceCategoryHelp')}
            />
            
            <FormControl fullWidth margin="dense" sx={{ mt: 2 }}>
              <InputLabel id="target-category-label">
                {t('marketplace:store.categoryMapping.targetCategory')}
              </InputLabel>
              <Select
                labelId="target-category-label"
                value={currentMapping.target}
                onChange={(e) => setCurrentMapping({...currentMapping, target: e.target.value})}
                label={t('marketplace:store.categoryMapping.targetCategory')}
              >
                <MenuItem value={9999}>
                  {t('marketplace:categories.other')} (ID: 9999)
                </MenuItem>
                
                <Divider />
                
                {filteredCategories.map(category => (
                  <MenuItem key={category.id} value={category.id}>
                    {category.name} (ID: {category.id})
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialogOpen(false)}>
            {t('common:buttons.cancel')}
          </Button>
          <Button 
            onClick={handleSaveMapping} 
            color="primary"
            disabled={!currentMapping.source.trim()}
          >
            {t('common:buttons.save')}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default CategoryMappingEditor;