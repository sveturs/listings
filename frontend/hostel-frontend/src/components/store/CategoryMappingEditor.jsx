// frontend/hostel-frontend/src/pages/store/CategoryMappingEditor.jsx
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
  Divider,
  CircularProgress,
  Alert,
  InputAdornment,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  ListSubheader,
  Chip,
  Tooltip
} from '@mui/material';
import {
  Plus,
  Edit2,
  Trash2,
  Save,
  X,
  Search,
  RefreshCw,
  ChevronRight,
  ChevronDown
} from 'lucide-react';
import axios from '../../api/axios';
import HierarchicalCategorySelect from '../../components/marketplace/HierarchicalCategorySelect';

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

  // Состояние для импортированных категорий
  const [importedCategories, setImportedCategories] = useState([]);
  const [organizedCategories, setOrganizedCategories] = useState({});
  const [importedCategoriesLoading, setImportedCategoriesLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  // Состояние для применения сопоставлений
  const [applyingMappings, setApplyingMappings] = useState(false);
  const [applyResult, setApplyResult] = useState(null);

  // Загружаем данные при инициализации
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

  // Организуем импортированные категории в иерархию
  useEffect(() => {
    if (importedCategories.length > 0 && Object.keys(mappings).length >= 0) {
      const categoriesTree = organizeImportedCategories(importedCategories);
      const markedTree = markMappedCategories(categoriesTree, mappings);
      setOrganizedCategories(markedTree);
    }
  }, [importedCategories, mappings]);

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

  // Организуем категории по иерархии
  const organizeImportedCategories = (importedCategories) => {
    const categoryTree = {};

    // Сортируем категории для удобства
    const sortedCategories = [...importedCategories].sort();

    sortedCategories.forEach(category => {
      // Разбиваем строку категории на уровни, разделенные "|"
      const levels = category.split('|');

      if (levels.length === 1) {
        // Это одиночная категория верхнего уровня
        if (!categoryTree[levels[0]]) {
          categoryTree[levels[0]] = { children: {}, mapped: false };
        }
      } else if (levels.length === 2) {
        // Это категория второго уровня с родителем
        const [parent, child] = levels;

        if (!categoryTree[parent]) {
          categoryTree[parent] = { children: {}, mapped: false };
        }

        categoryTree[parent].children[child] = { children: {}, mapped: false };

      } else if (levels.length === 3) {
        // Это категория третьего уровня с двумя родителями
        const [grandparent, parent, child] = levels;

        if (!categoryTree[grandparent]) {
          categoryTree[grandparent] = { children: {}, mapped: false };
        }

        if (!categoryTree[grandparent].children[parent]) {
          categoryTree[grandparent].children[parent] = { children: {}, mapped: false };
        }

        categoryTree[grandparent].children[parent].children[child] = { mapped: false };
      }
    });

    return categoryTree;
  };

  // Отмечаем сопоставленные категории в дереве
  const markMappedCategories = (categoryTree, mappings) => {
    const result = { ...categoryTree };

    // Проходимся по каждому отображению
    Object.keys(mappings).forEach(sourceCategory => {
      const levels = sourceCategory.split('|');

      if (levels.length === 1 && result[levels[0]]) {
        result[levels[0]].mapped = true;
        result[levels[0]].mappedTo = mappings[sourceCategory];
      } else if (levels.length === 2) {
        const [parent, child] = levels;
        if (result[parent]?.children[child]) {
          result[parent].children[child].mapped = true;
          result[parent].children[child].mappedTo = mappings[sourceCategory];
        }
      } else if (levels.length === 3) {
        const [grandparent, parent, child] = levels;
        if (result[grandparent]?.children[parent]?.children[child]) {
          result[grandparent].children[parent].children[child].mapped = true;
          result[grandparent].children[parent].children[child].mappedTo = mappings[sourceCategory];
        }
      }
    });

    return result;
  };

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

  // Применение сопоставлений к товарам
  const handleApplyMappings = async () => {
    if (Object.keys(mappings).length === 0) {
      setError(t('marketplace:store.categoryMapping.noMappingsToApply'));
      return;
    }

    setApplyingMappings(true);
    setApplyResult(null);
    setError(null);

    try {
      const response = await axios.post(
        `/api/v1/storefronts/import-sources/${sourceId}/apply-category-mappings`
      );

      if (response.data.success) {
        setApplyResult({
          message: t('marketplace:store.categoryMapping.applySuccess', {
            count: response.data.data.updated_count,
            defaultValue: 'Successfully updated categories for {{count}} listings'
          }),
          count: response.data.data.updated_count
        });
      } else {
        setError(t('marketplace:store.categoryMapping.applyError'));
      }
    } catch (err) {
      console.error('Error applying mappings:', err);
      setError(err.response?.data?.message || t('marketplace:store.categoryMapping.applyError'));
    } finally {
      setApplyingMappings(false);
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
  const buildCategoryTree = (flatCategories) => {
    // Если категории уже имеют иерархическую структуру (с children),
    // просто вернем их
    if (flatCategories.some(c => c.children && c.children.length > 0)) {
      return flatCategories;
    }
    
    // Создаем словарь категорий по ID
    const categoryMap = {};
    flatCategories.forEach(cat => {
      categoryMap[cat.id] = { ...cat, children: [] };
    });
    
    // Формируем дерево категорий
    const rootCategories = [];
    flatCategories.forEach(cat => {
      if (cat.parent_id) {
        // Это дочерняя категория
        if (categoryMap[cat.parent_id]) {
          categoryMap[cat.parent_id].children.push(categoryMap[cat.id]);
        } else {
          // Если родительская категория не найдена, добавляем в корень
          rootCategories.push(categoryMap[cat.id]);
        }
      } else {
        // Это корневая категория
        rootCategories.push(categoryMap[cat.id]);
      }
    });
    
    return rootCategories;
  };
  // Удаление сопоставления
  const handleDeleteMapping = (sourceCategory) => {
    const newMappings = { ...mappings };
    delete newMappings[sourceCategory];
    setMappings(newMappings);
  };

  // Сохранение изменений в диалоге
  const handleSaveMapping = () => {
    if (!currentMapping.source.trim() || !currentMapping.target) {
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

      {applyResult && (
        <Alert severity="success" sx={{ mb: 2 }}>
          {applyResult.message}
        </Alert>
      )}

      <Paper sx={{ p: 2, mb: 3 }}>
        <Typography variant="subtitle1" gutterBottom>
          {t('marketplace:store.categoryMapping.description')}
        </Typography>

        <Alert severity="info" sx={{ mb: 2 }}>
          {t('marketplace:store.categoryMapping.info')}
        </Alert>

        {/* Статистика сопоставлений */}
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="subtitle2">
            {t('marketplace:store.categoryMapping.statistics', { defaultValue: 'Статистика сопоставлений:' })}
          </Typography>

          <Box display="flex" gap={1}>
            <Chip
              size="small"
              color="primary"
              label={t('marketplace:store.categoryMapping.mappingsCount', {
                count: Object.keys(mappings).length,
                defaultValue: 'Сопоставлений: {{count}}'
              })}
            />

            <Chip
              size="small"
              color="secondary"
              label={t('marketplace:store.categoryMapping.importedCategoriesCount', {
                count: importedCategories.length,
                defaultValue: 'Импортировано категорий: {{count}}'
              })}
            />
          </Box>
        </Box>

        <Divider sx={{ mb: 3 }} />

        {/* Поиск по категориям (если необходимо) */}
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

        {/* Импортированные категории с иерархической структурой */}
        <Box sx={{ mt: 3 }}>
          <Typography variant="subtitle2" gutterBottom>
            {t('marketplace:store.categoryMapping.importedCategories', { defaultValue: 'Импортированные категории' })}
          </Typography>

          {importedCategoriesLoading ? (
            <Box display="flex" alignItems="center" gap={1}>
              <CircularProgress size={20} />
              <Typography variant="body2">
                {t('marketplace:store.categoryMapping.loadingImportedCategories', { defaultValue: 'Загрузка категорий...' })}
              </Typography>
            </Box>
          ) : (
            <>
              {importedCategories.length === 0 ? (
                <Alert severity="info">
                  {t('marketplace:store.categoryMapping.noImportedCategories', { defaultValue: 'Нет импортированных категорий' })}
                </Alert>
              ) : (
                <Paper variant="outlined" sx={{ maxHeight: 350, overflow: 'auto' }}>
                  {/* Обработка организованных категорий */}
                  {Object.entries(organizedCategories).map(([parentCategory, parentData]) => (
                    <Box key={parentCategory} sx={{ mb: 0.5, borderBottom: '1px solid', borderColor: 'divider' }}>
                      {/* Родительская категория первого уровня */}
                      <Box
                        sx={{
                          p: 1.5,
                          pl: 2,
                          display: 'flex',
                          justifyContent: 'space-between',
                          alignItems: 'center',
                          bgcolor: parentData.mapped ? 'action.selected' : 'transparent'
                        }}
                      >
                        <Box>
                          <Typography variant="subtitle2" component="div">
                            {parentCategory}
                          </Typography>
                          {parentData.mapped && (
                            <Typography variant="caption" color="primary">
                              → {formatCategoryName(parentData.mappedTo)} (ID: {parentData.mappedTo})
                            </Typography>
                          )}
                        </Box>

                        <Box>
                          {parentData.mapped ? (
                            <IconButton
                              size="small"
                              onClick={() => handleEditMapping(parentCategory)}
                              aria-label="edit mapping"
                            >
                              <Edit2 size={18} />
                            </IconButton>
                          ) : (
                            <Button
                              size="small"
                              variant="outlined"
                              onClick={() => handleEditMapping(parentCategory)}
                              startIcon={<Plus size={16} />}
                            >
                              {t('marketplace:store.categoryMapping.mapCategory', { defaultValue: 'Сопоставить' })}
                            </Button>
                          )}
                        </Box>
                      </Box>

                      {/* Категории второго уровня */}
                      {Object.keys(parentData.children).length > 0 && (
                        <Box sx={{ pl: 4 }}>
                          {Object.entries(parentData.children).map(([childCategory, childData]) => (
                            <Box key={childCategory}>
                              {/* Категория второго уровня */}
                              <Box
                                sx={{
                                  p: 1,
                                  borderTop: '1px dashed',
                                  borderColor: 'divider',
                                  display: 'flex',
                                  justifyContent: 'space-between',
                                  alignItems: 'center',
                                  bgcolor: childData.mapped ? 'action.selected' : 'transparent'
                                }}
                              >
                                <Box sx={{ display: 'flex', alignItems: 'center' }}>
                                  <ChevronRight size={16} style={{ marginRight: 8, opacity: 0.6 }} />
                                  <Box>
                                    <Typography variant="body2" component="div">
                                      {childCategory}
                                    </Typography>
                                    {childData.mapped && (
                                      <Typography variant="caption" color="primary">
                                        → {formatCategoryName(childData.mappedTo)} (ID: {childData.mappedTo})
                                      </Typography>
                                    )}
                                  </Box>
                                </Box>

                                <Box>
                                  {childData.mapped ? (
                                    <IconButton
                                      size="small"
                                      onClick={() => handleEditMapping(`${parentCategory}|${childCategory}`)}
                                      aria-label="edit mapping"
                                    >
                                      <Edit2 size={16} />
                                    </IconButton>
                                  ) : (
                                    <Button
                                      size="small"
                                      variant="outlined"
                                      onClick={() => handleEditMapping(`${parentCategory}|${childCategory}`)}
                                      startIcon={<Plus size={14} />}
                                    >
                                      {t('marketplace:store.categoryMapping.mapCategory', { defaultValue: 'Сопоставить' })}
                                    </Button>
                                  )}
                                </Box>
                              </Box>

                              {/* Категории третьего уровня */}
                              {Object.keys(childData.children).length > 0 && (
                                <Box sx={{ pl: 4 }}>
                                  {Object.entries(childData.children).map(([grandchildCategory, grandchildData]) => (
                                    <Box
                                      key={grandchildCategory}
                                      sx={{
                                        p: 1,
                                        borderTop: '1px dotted',
                                        borderColor: 'divider',
                                        display: 'flex',
                                        justifyContent: 'space-between',
                                        alignItems: 'center',
                                        bgcolor: grandchildData.mapped ? 'action.selected' : 'transparent'
                                      }}
                                    >
                                      <Box sx={{ display: 'flex', alignItems: 'center' }}>
                                        <ChevronRight size={14} style={{ marginRight: 8, opacity: 0.5 }} />
                                        <Box>
                                          <Typography variant="body2" color="text.secondary" component="div">
                                            {grandchildCategory}
                                          </Typography>
                                          {grandchildData.mapped && (
                                            <Typography variant="caption" color="primary">
                                              → {formatCategoryName(grandchildData.mappedTo)} (ID: {grandchildData.mappedTo})
                                            </Typography>
                                          )}
                                        </Box>
                                      </Box>

                                      <Box>
                                        {grandchildData.mapped ? (
                                          <IconButton
                                            size="small"
                                            onClick={() => handleEditMapping(`${parentCategory}|${childCategory}|${grandchildCategory}`)}
                                            aria-label="edit mapping"
                                          >
                                            <Edit2 size={14} />
                                          </IconButton>
                                        ) : (
                                          <Button
                                            size="small"
                                            variant="outlined"
                                            onClick={() => handleEditMapping(`${parentCategory}|${childCategory}|${grandchildCategory}`)}
                                            startIcon={<Plus size={12} />}
                                          >
                                            {t('marketplace:store.categoryMapping.mapCategory', { defaultValue: 'Сопоставить' })}
                                          </Button>
                                        )}
                                      </Box>
                                    </Box>
                                  ))}
                                </Box>
                              )}
                            </Box>
                          ))}
                        </Box>
                      )}
                    </Box>
                  ))}
                </Paper>
              )}
            </>
          )}
        </Box>

        <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 3 }}>
          <Button
            variant="outlined"
            startIcon={<Plus />}
            onClick={handleAddMapping}
          >
            {t('marketplace:store.categoryMapping.addMapping', { defaultValue: 'Добавить сопоставление' })}
          </Button>

          <Box>
            <Button
              variant="outlined"
              onClick={onClose}
              sx={{ mr: 1 }}
              disabled={saving || applyingMappings}
            >
              {t('common:buttons.cancel', { defaultValue: 'Отмена' })}
            </Button>

            <Tooltip
              title={t('marketplace:store.categoryMapping.applyHelp', {
                defaultValue: 'Эта кнопка обновит категории всех товаров, которые были импортированы из этого источника, согласно настроенным сопоставлениям.'
              })}
              placement="top"
            >
              <Button
                variant="outlined"
                color="secondary"
                startIcon={applyingMappings ? <CircularProgress size={20} /> : <RefreshCw />}
                onClick={handleApplyMappings}
                disabled={saving || applyingMappings || Object.keys(mappings).length === 0}
                sx={{ mr: 1 }}
              >
                {applyingMappings
                  ? t('marketplace:store.categoryMapping.applyingMappings', { defaultValue: 'Применение...' })
                  : t('marketplace:store.categoryMapping.applyMappings', { defaultValue: 'Применить к товарам' })
                }
              </Button>
            </Tooltip>

            <Button
              variant="contained"
              startIcon={saving ? <CircularProgress size={20} /> : <Save />}
              onClick={handleSave}
              disabled={saving || applyingMappings}
            >
              {t('common:buttons.save', { defaultValue: 'Сохранить' })}
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
            t('marketplace:store.categoryMapping.editMappingTitle', { defaultValue: 'Редактирование сопоставления' }) :
            t('marketplace:store.categoryMapping.addMappingTitle', { defaultValue: 'Добавление сопоставления' })}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 1 }}>
            <TextField
              fullWidth
              label={t('marketplace:store.categoryMapping.sourceCategory', { defaultValue: 'Категория источника' })}
              value={currentMapping.source}
              onChange={(e) => setCurrentMapping({ ...currentMapping, source: e.target.value })}
              margin="dense"
              helperText={t('marketplace:store.categoryMapping.sourceCategoryHelp', { defaultValue: 'Введите категорию из источника импорта' })}
            />

            <Box sx={{ mt: 3, mb: 1 }}>
              <Typography variant="subtitle2" gutterBottom>
                {t('marketplace:store.categoryMapping.targetCategoryLabel', { defaultValue: 'Выберите категорию для сопоставления:' })}
              </Typography>

              <HierarchicalCategorySelect
                categories={buildCategoryTree(categories)}
                value={currentMapping.target}
                onChange={(value) => setCurrentMapping({ ...currentMapping, target: value })}
                placeholder={t('marketplace:store.categoryMapping.selectCategory', { defaultValue: 'Выберите категорию' })}
              />
            </Box>

            {/* Дополнительная опция для категории "Прочее" */}
            <Box sx={{ mt: 2 }}>
              <Button
                variant="outlined"
                color="primary"
                size="small"
                onClick={() => setCurrentMapping({ ...currentMapping, target: 9999 })}
                startIcon={<Plus size={16} />}
              >
                {t('marketplace:store.categoryMapping.useDefaultCategory', { defaultValue: 'Использовать категорию "Прочее"' })}
              </Button>
            </Box>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialogOpen(false)}>
            {t('common:buttons.cancel', { defaultValue: 'Отмена' })}
          </Button>
          <Button
            onClick={handleSaveMapping}
            color="primary"
            disabled={!currentMapping.source.trim() || !currentMapping.target}
          >
            {t('common:buttons.save', { defaultValue: 'Сохранить' })}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default CategoryMappingEditor;