// frontend/hostel-frontend/src/components/store/CategoryMappingEditor.jsx
import React, { useState, useEffect, useMemo, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Box,
  Typography,
  TextField,
  Button,
  IconButton,
  Divider,
  CircularProgress,
  Alert,
  InputAdornment,
  Chip,
  Tooltip,
  Paper,
  Collapse,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText
} from '@mui/material';
import {
  Save,
  X,
  Search,
  RefreshCw,
  ChevronRight,
  ChevronDown,
  FolderClosed,
  FolderOpen,
  Tag,
  Check,
  Move,
  Trash2
} from 'lucide-react';
import axios from '../../api/axios';

const DraggableSystemCategory = ({ category, onDragStart }) => {
  // Используем встроенный HTML5 Drag and Drop API
  const handleDragStart = (e) => {
    // Передаем ID категории через dataTransfer
    e.dataTransfer.setData('text/plain', category.id);
    e.dataTransfer.setData('categoryName', category.name);
    
    // Добавляем эффект перемещения
    e.dataTransfer.effectAllowed = 'move';
    
    // Вызываем колбэк onDragStart
    onDragStart?.(category);
  };

  return (
    <Box
      draggable
      onDragStart={handleDragStart}
      sx={{
        cursor: 'grab',
        p: 1,
        display: 'flex',
        alignItems: 'center',
        flexDirection: 'column',
        alignItems: 'flex-start',
        borderRadius: 1,
        '&:hover': {
          bgcolor: 'action.hover'
        }
      }}
    >
      <Typography variant="body2">
        {category.name}
      </Typography>
      
      {/* Показываем полный путь, если он есть */}
      {category.pathLabel && (
        <Typography 
          variant="caption" 
          color="text.secondary"
          sx={{ mt: 0.5, fontSize: '0.7rem' }}
        >
          {category.pathLabel}
        </Typography>
      )}
    </Box>
  );
};


// Компонент для целевой зоны перетаскивания (импортированная категория)
const DroppableImportCategory = ({ sourceCategory, sourcePath, mappedTo, formatCategoryName, onDrop }) => {
  const { t } = useTranslation(['marketplace', 'common']);
  const [isOver, setIsOver] = useState(false);
  
  // Обрабатываем события перетаскивания
  const handleDragOver = (e) => {
    // Предотвращаем стандартное поведение, чтобы разрешить drop
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
    if (!isOver) setIsOver(true);
  };
  
  const handleDragLeave = () => {
    setIsOver(false);
  };
  
  const handleDrop = (e) => {
    e.preventDefault();
    setIsOver(false);
    
    // Получаем ID категории из данных перетаскивания
    const categoryId = parseInt(e.dataTransfer.getData('text/plain'), 10);
    if (!isNaN(categoryId)) {
      onDrop(sourcePath, categoryId);
    }
  };
  
  const handleRemoveMapping = (e) => {
    e.stopPropagation();
    onDrop(sourcePath, null);
  };

  return (
    <Box
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      sx={{
        p: 1.5,
        borderRadius: 1,
        bgcolor: isOver ? 'action.selected' : (mappedTo ? 'action.selected' : 'transparent'),
        border: isOver ? '1px dashed' : '1px solid transparent',
        borderColor: isOver ? 'primary.main' : 'divider',
        transition: 'all 0.2s',
        '&:hover': {
          bgcolor: isOver ? 'action.selected' : (mappedTo ? 'action.selected' : 'action.hover')
        }
      }}
    >
      <Box display="flex" alignItems="center" justifyContent="space-between">
        <Typography variant="body2" component="div">
          {sourceCategory}
        </Typography>
        
        {isOver && (
          <Chip 
            label={t('marketplace:store.categoryMapping.dropHere', { defaultValue: 'Перетащите сюда' })} 
            color="primary" 
            size="small"
            variant="outlined"
          />
        )}
      </Box>
      
      {mappedTo && (
        <Box sx={{ mt: 0.5, display: 'flex', alignItems: 'center' }}>
          <ChevronRight size={14} style={{ opacity: 0.6 }} />
          <Typography variant="caption" color="primary" sx={{ fontWeight: 'medium' }}>
            {formatCategoryName(mappedTo)}
          </Typography>
          <IconButton 
            size="small"
            sx={{ ml: 'auto' }}
            onClick={handleRemoveMapping}
          >
            <X size={14} />
          </IconButton>
        </Box>
      )}
    </Box>
  );
};

// Компонент для иерархического отображения категорий системы
const SystemCategoryTree = ({ categories, onCategoryDragStart }) => {
  const [expanded, setExpanded] = useState({});
  const { t } = useTranslation(['marketplace', 'common']);

  // Функция для построения дерева категорий
  const buildCategoryTree = (categories) => {
    // Создаем словарь категорий
    const categoriesMap = {};
    categories.forEach(cat => {
      categoriesMap[cat.id] = { ...cat, children: [] };
    });

    // Строим древовидную структуру
    const rootCategories = [];
    categories.forEach(cat => {
      if (cat.parent_id) {
        if (categoriesMap[cat.parent_id]) {
          categoriesMap[cat.parent_id].children.push(categoriesMap[cat.id]);
        } else {
          rootCategories.push(categoriesMap[cat.id]);
        }
      } else {
        rootCategories.push(categoriesMap[cat.id]);
      }
    });

    return rootCategories;
  };

  // Строим дерево категорий
  const categoryTree = useMemo(() => buildCategoryTree(categories), [categories]);

  // Функция для обработки клика по категории
  const handleToggle = (categoryId, e) => {
    e.stopPropagation();
    setExpanded(prev => ({
      ...prev,
      [categoryId]: !prev[categoryId]
    }));
  };

  // Рекурсивный рендер категорий и их потомков
  const renderCategoryItem = (category, level = 0) => {
    const isExpanded = expanded[category.id];
    const hasChildren = category.children && category.children.length > 0;

    return (
      <React.Fragment key={category.id}>
        <ListItem
          disablePadding
          sx={{
            pl: level * 2,
            borderRadius: 1,
            mb: 0.5
          }}
        >
          <ListItemButton
            sx={{
              py: 0.5,
              borderRadius: 1
            }}
            dense
          >
            <ListItemIcon sx={{ minWidth: 36 }}>
              {hasChildren ? (
                <IconButton
                  edge="start"
                  size="small"
                  onClick={(e) => handleToggle(category.id, e)}
                  sx={{ mr: 1 }}
                >
                  {isExpanded ? <ChevronDown size={18} /> : <ChevronRight size={18} />}
                </IconButton>
              ) : (
                <Box sx={{ width: 18, mr: 1 }} /> // Пустой элемент для выравнивания
              )}
            </ListItemIcon>
            
            <DraggableSystemCategory 
              category={category} 
              onDragStart={onCategoryDragStart}
            />
          </ListItemButton>
        </ListItem>
        
        {hasChildren && (
          <Collapse in={isExpanded} timeout="auto" unmountOnExit>
            <List disablePadding>
              {category.children.map(child => renderCategoryItem(child, level + 1))}
            </List>
          </Collapse>
        )}
      </React.Fragment>
    );
  };

  // Рендерим корневые категории
  return (
    <List disablePadding>
      {categoryTree.map(category => renderCategoryItem(category))}
    </List>
  );
};

// Компонент для иерархического отображения импортированных категорий
const ImportedCategoryTree = ({ categories, mappings, formatCategoryName, onDrop }) => {
  const [expanded, setExpanded] = useState({});
  const { t } = useTranslation(['marketplace', 'common']);

  // Функция для обработки клика по категории
  const handleToggle = (categoryKey, e) => {
    e.stopPropagation();
    setExpanded(prev => ({
      ...prev,
      [categoryKey]: !prev[categoryKey]
    }));
  };

  if (Object.keys(categories).length === 0) {
    return (
      <Typography variant="body2" color="text.secondary" sx={{ p: 2, textAlign: 'center' }}>
        {t('marketplace:store.categoryMapping.noCategories')}
      </Typography>
    );
  }

  return (
    <List disablePadding>
      {Object.entries(categories).map(([parentCategory, parentData]) => (
        <React.Fragment key={parentCategory}>
          <ListItem
            disablePadding
            sx={{ mb: 0.5 }}
          >
            <ListItemButton 
              dense
              onClick={(e) => handleToggle(parentCategory, e)}
              sx={{ py: 0.5, px: 1, borderRadius: 1 }}
            >
              <ListItemIcon sx={{ minWidth: 36 }}>
                {Object.keys(parentData.children).length > 0 ? (
                  expanded[parentCategory] ? <FolderOpen size={18} /> : <FolderClosed size={18} />
                ) : (
                  <Tag size={18} />
                )}
              </ListItemIcon>
              
              <DroppableImportCategory 
                sourceCategory={parentCategory} 
                sourcePath={parentCategory}
                mappedTo={parentData.mappedTo}
                formatCategoryName={formatCategoryName}
                onDrop={onDrop}
              />
            </ListItemButton>
          </ListItem>
          
          {Object.keys(parentData.children).length > 0 && (
            <Collapse in={expanded[parentCategory]} timeout="auto" unmountOnExit>
              <List disablePadding sx={{ pl: 4 }}>
                {Object.entries(parentData.children).map(([childCategory, childData]) => (
                  <React.Fragment key={`${parentCategory}|${childCategory}`}>
                    <ListItem
                      disablePadding
                      sx={{ mb: 0.5 }}
                    >
                      <ListItemButton 
                        dense
                        onClick={(e) => handleToggle(`${parentCategory}|${childCategory}`, e)}
                        sx={{ py: 0.5, px: 1, borderRadius: 1 }}
                      >
                        <ListItemIcon sx={{ minWidth: 36 }}>
                          {Object.keys(childData.children).length > 0 ? (
                            expanded[`${parentCategory}|${childCategory}`] ? <FolderOpen size={16} /> : <FolderClosed size={16} />
                          ) : (
                            <Tag size={16} />
                          )}
                        </ListItemIcon>
                        
                        <DroppableImportCategory 
                          sourceCategory={childCategory} 
                          sourcePath={`${parentCategory}|${childCategory}`}
                          mappedTo={childData.mappedTo}
                          formatCategoryName={formatCategoryName}
                          onDrop={onDrop}
                        />
                      </ListItemButton>
                    </ListItem>
                    
                    {Object.keys(childData.children).length > 0 && (
                      <Collapse in={expanded[`${parentCategory}|${childCategory}`]} timeout="auto" unmountOnExit>
                        <List disablePadding sx={{ pl: 4 }}>
                          {Object.entries(childData.children).map(([grandchildCategory, grandchildData]) => (
                            <ListItem
                              key={`${parentCategory}|${childCategory}|${grandchildCategory}`}
                              disablePadding
                              sx={{ mb: 0.5 }}
                            >
                              <ListItemButton 
                                dense
                                sx={{ py: 0.5, px: 1, borderRadius: 1 }}
                              >
                                <ListItemIcon sx={{ minWidth: 36 }}>
                                  <Tag size={14} />
                                </ListItemIcon>
                                
                                <DroppableImportCategory 
                                  sourceCategory={grandchildCategory} 
                                  sourcePath={`${parentCategory}|${childCategory}|${grandchildCategory}`}
                                  mappedTo={grandchildData.mappedTo}
                                  formatCategoryName={formatCategoryName}
                                  onDrop={onDrop}
                                />
                              </ListItemButton>
                            </ListItem>
                          ))}
                        </List>
                      </Collapse>
                    )}
                  </React.Fragment>
                ))}
              </List>
            </Collapse>
          )}
        </React.Fragment>
      ))}
    </List>
  );
};

// Основной компонент редактора сопоставления категорий
const CategoryMappingEditor = ({ sourceId, onClose, onSave }) => {
  const { t } = useTranslation(['marketplace', 'common']);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [mappings, setMappings] = useState({});
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);
  const [categories, setCategories] = useState([]);
  
  // Состояние для импортированных категорий
  const [importedCategories, setImportedCategories] = useState([]);
  const [organizedCategories, setOrganizedCategories] = useState({});
  const [importedCategoriesLoading, setImportedCategoriesLoading] = useState(false);
  const [searchSystemTerm, setSearchSystemTerm] = useState('');
  const [searchImportTerm, setSearchImportTerm] = useState('');
  const [filteredCategories, setFilteredCategories] = useState([]);
  const [filteredImportCategories, setFilteredImportCategories] = useState({});

  // Состояние для применения сопоставлений
  const [applyingMappings, setApplyingMappings] = useState(false);
  const [applyResult, setApplyResult] = useState(null);
  
  // Состояние для отслеживания перетаскивания
  const [draggedCategory, setDraggedCategory] = useState(null);

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

  // Организуем импортированные категории в иерархию
  useEffect(() => {
    if (importedCategories.length > 0) {
      const categoriesTree = organizeImportedCategories(importedCategories);
      const markedTree = markMappedCategories(categoriesTree, mappings);
      setOrganizedCategories(markedTree);
      setFilteredImportCategories(markedTree);
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

  // Фильтрация системных категорий при поиске
  useEffect(() => {
    if (searchSystemTerm) {
      // Расширенная фильтрация категорий с сохранением информации о пути
      let filtered = [];
      
      // Построим полные пути для всех категорий
      const buildCategoryPaths = (categories) => {
        const categoryPaths = {};
        const categoriesMap = {};
        
        // Создаем карту категорий по ID
        categories.forEach(cat => {
          categoriesMap[cat.id] = cat;
        });
        
        // Рекурсивная функция для построения пути
        const buildPath = (categoryId, path = []) => {
          const category = categoriesMap[categoryId];
          if (!category) return [...path];
          
          const newPath = [...path, category];
          
          // Если есть родитель, продолжаем строить путь
          if (category.parent_id) {
            return buildPath(category.parent_id, newPath);
          }
          
          return newPath.reverse(); // Разворачиваем путь от корня к листу
        };
        
        // Строим пути для всех категорий
        categories.forEach(cat => {
          categoryPaths[cat.id] = buildPath(cat.id);
        });
        
        return categoryPaths;
      };
      
      const categoryPaths = buildCategoryPaths(categories);
      
      // Фильтруем категории по поисковому запросу
      categories.forEach(cat => {
        if (cat.name.toLowerCase().includes(searchSystemTerm.toLowerCase())) {
          // Добавляем элемент с полным путем
          const pathWithLabels = categoryPaths[cat.id].map(c => c.name).join(' > ');
          cat.pathLabel = pathWithLabels;
          filtered.push(cat);
        }
      });
      
      setFilteredCategories(filtered);
    } else {
      // Если поиск пуст, сбрасываем состояние
      setFilteredCategories(categories.map(cat => ({...cat, pathLabel: null})));
    }
  }, [searchSystemTerm, categories]);

  // Фильтрация импортированных категорий при поиске
  useEffect(() => {
    if (!searchImportTerm || searchImportTerm === '') {
      setFilteredImportCategories(organizedCategories);
      return;
    }
    
    const searchLower = searchImportTerm.toLowerCase();
    const filtered = {};
    
    // Поиск по всем уровням категорий
    Object.entries(organizedCategories).forEach(([parentKey, parentData]) => {
      // Если родитель соответствует поиску
      if (parentKey.toLowerCase().includes(searchLower)) {
        filtered[parentKey] = { ...parentData };
        return;
      }
      
      // Поиск по дочерним категориям
      const filteredChildren = {};
      let hasMatchingChildren = false;
      
      Object.entries(parentData.children).forEach(([childKey, childData]) => {
        // Если ребенок соответствует поиску
        if (childKey.toLowerCase().includes(searchLower)) {
          filteredChildren[childKey] = { ...childData };
          hasMatchingChildren = true;
          return;
        }
        
        // Поиск по внукам
        const filteredGrandchildren = {};
        let hasMatchingGrandchildren = false;
        
        Object.entries(childData.children).forEach(([grandchildKey, grandchildData]) => {
          if (grandchildKey.toLowerCase().includes(searchLower)) {
            filteredGrandchildren[grandchildKey] = { ...grandchildData };
            hasMatchingGrandchildren = true;
          }
        });
        
        if (hasMatchingGrandchildren) {
          filteredChildren[childKey] = {
            ...childData,
            children: filteredGrandchildren
          };
          hasMatchingChildren = true;
        }
      });
      
      if (hasMatchingChildren) {
        filtered[parentKey] = {
          ...parentData,
          children: filteredChildren
        };
      }
    });
    
    setFilteredImportCategories(filtered);
  }, [searchImportTerm, organizedCategories]);

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

  // Обработка перетаскивания категории
  const handleCategoryDrop = (sourcePath, targetCategoryId) => {
    const newMappings = { ...mappings };
    
    if (targetCategoryId === null) {
      // Если targetCategoryId === null, то это удаление сопоставления
      delete newMappings[sourcePath];
    } else {
      // Добавляем или обновляем сопоставление
      newMappings[sourcePath] = targetCategoryId;
    }
    
    setMappings(newMappings);
  };

  // Отслеживаем начало перетаскивания категории
  const handleCategoryDragStart = (category) => {
    setDraggedCategory(category);
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
          {t('marketplace:store.categoryMapping.dragDropInfo', { 
            defaultValue: 'Перетащите категорию из правой панели на импортированную категорию слева, чтобы создать сопоставление. Для удаления сопоставления нажмите на значок X рядом с сопоставленной категорией.' 
          })}
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

        <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, gap: 2 }}>
          {/* Левая колонка - импортированные категории */}
          <Box sx={{ width: { xs: '100%', md: '50%' } }}>
            <Typography variant="subtitle2" gutterBottom>
              {t('marketplace:store.categoryMapping.importedCategories', { defaultValue: 'Импортированные категории' })}
              <Typography variant="caption" sx={{ ml: 1, color: 'text.secondary' }}>
                ({t('marketplace:store.categoryMapping.dropTarget', { defaultValue: 'Целевая зона для перетаскивания' })})
              </Typography>
            </Typography>

            <TextField
              fullWidth
              variant="outlined"
              size="small"
              placeholder={t('marketplace:store.categoryMapping.searchImportedCategories')}
              value={searchImportTerm}
              onChange={(e) => setSearchImportTerm(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Search size={18} />
                  </InputAdornment>
                ),
                endAdornment: searchImportTerm ? (
                  <InputAdornment position="end">
                    <IconButton
                      size="small"
                      onClick={() => setSearchImportTerm('')}
                    >
                      <X size={16} />
                    </IconButton>
                  </InputAdornment>
                ) : null
              }}
              sx={{ mb: 2 }}
            />

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
                  <Paper 
                    variant="outlined" 
                    sx={{ 
                      height: 400, 
                      overflow: 'auto',
                      p: 1,
                      bgcolor: 'background.default'
                    }}
                  >
                    <ImportedCategoryTree 
                      categories={filteredImportCategories} 
                      mappings={mappings}
                      formatCategoryName={formatCategoryName}
                      onDrop={handleCategoryDrop}
                    />
                  </Paper>
                )}
              </>
            )}
          </Box>

          {/* Правая колонка - категории системы для перетаскивания */}
          <Box sx={{ width: { xs: '100%', md: '50%' } }}>
            <Typography variant="subtitle2" gutterBottom>
              {t('marketplace:store.categoryMapping.systemCategories', { defaultValue: 'Категории системы' })}
              <Typography variant="caption" sx={{ ml: 1, color: 'text.secondary' }}>
                ({t('marketplace:store.categoryMapping.dragSource', { defaultValue: 'Перетащите на импортированную категорию' })})
              </Typography>
            </Typography>

            <TextField
              fullWidth
              variant="outlined"
              size="small"
              placeholder={t('marketplace:store.categoryMapping.searchCategories')}
              value={searchSystemTerm}
              onChange={(e) => setSearchSystemTerm(e.target.value)}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Search size={18} />
                  </InputAdornment>
                ),
                endAdornment: searchSystemTerm ? (
                  <InputAdornment position="end">
                    <IconButton
                      size="small"
                      onClick={() => setSearchSystemTerm('')}
                    >
                      <X size={16} />
                    </IconButton>
                  </InputAdornment>
                ) : null
              }}
              sx={{ mb: 2 }}
            />

            <Paper 
              variant="outlined" 
              sx={{ 
                height: 400, 
                overflow: 'auto',
                p: 1,
                bgcolor: 'background.default'
              }}
            >
              <SystemCategoryTree 
                categories={filteredCategories}
                onCategoryDragStart={handleCategoryDragStart}
              />
            </Paper>
          </Box>
        </Box>

        <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 3 }}>
          <Button
            variant="outlined"
            onClick={onClose}
            sx={{ mr: 1 }}
            disabled={saving || applyingMappings}
          >
            {t('common:buttons.cancel', { defaultValue: 'Отмена' })}
          </Button>

          <Box>
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
    </Box>
  );
};

export default CategoryMappingEditor;