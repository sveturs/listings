import React, { useState, useEffect, Fragment } from 'react';
import {
  Box,
  Paper,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Chip,
  Tooltip,
  Alert,
  CircularProgress,
  InputAdornment,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Checkbox,
  Divider,
  Grid,
  FormControlLabel,
  Switch,
  Autocomplete,
  Tab,
  Tabs,
  Snackbar
} from '@mui/material';
import {
  Edit as EditIcon,
  Delete as DeleteIcon,
  Add as AddIcon,
  DragIndicator as DragIcon,
  Link as LinkIcon,
  LinkOff as UnlinkIcon,
  Settings as SettingsIcon,
  CheckBox as CheckBoxIcon,
  CheckBoxOutlineBlank as CheckBoxOutlineBlankIcon,
  Search as SearchIcon,
  Clear as ClearIcon,
  ExpandMore as ExpandMoreIcon,
  ChevronRight as ChevronRightIcon
} from '@mui/icons-material';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';
import axios from '../../api/axios';

interface AttributeGroup {
  id: number;
  name: string;
  display_name: string;
  description: string;
  icon: string;
  sort_order: number;
  is_active: boolean;
  is_system: boolean;
  created_at: string;
  updated_at: string;
}

interface CategoryAttribute {
  id: number;
  name: string;
  display_name: string;
  description: string;
  data_type: string;
  required: boolean;
  icon: string;
}

interface AttributeGroupItem {
  id: number;
  group_id: number;
  attribute_id: number;
  icon: string;
  sort_order: number;
  custom_display_name: string;
  visibility_condition: any;
  created_at: string;
  attribute?: CategoryAttribute;
}

interface GroupWithItems {
  id: number;
  name: string;
  display_name: string;
  description: string;
  icon: string;
  sort_order: number;
  is_active: boolean;
  is_system: boolean;
  items: AttributeGroupItem[];
}

interface Category {
  id: number;
  name: string;
  display_name?: string;
  parent_id: number | null;
  children?: Category[];
  [key: string]: any; // для других полей
}

interface CategoryAttributeGroup {
  id: number;
  category_id: number;
  group_id: number;
  component_id: number | null;
  sort_order: number;
  is_active: boolean;
  display_mode: string;
  collapsed_by_default: boolean;
  configuration: any;
  created_at: string;
  group?: AttributeGroup;
}

const AttributeGroupsPage: React.FC = () => {
  const [groups, setGroups] = useState<AttributeGroup[]>([]);
  const [attributes, setAttributes] = useState<CategoryAttribute[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedGroup, setSelectedGroup] = useState<AttributeGroup | null>(null);
  const [selectedGroupItems, setSelectedGroupItems] = useState<AttributeGroupItem[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);
  const [categoryGroups, setCategoryGroups] = useState<CategoryAttributeGroup[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingCategories, setLoadingCategories] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState(0);
  const [categorySearchQuery, setCategorySearchQuery] = useState<string>('');
  const [expandedCategories, setExpandedCategories] = useState<Set<number>>(new Set());
  const [showOnlyLinked, setShowOnlyLinked] = useState(false);

  // Диалоги
  const [openGroupDialog, setOpenGroupDialog] = useState(false);
  const [openItemsDialog, setOpenItemsDialog] = useState(false);
  const [openCategoryDialog, setOpenCategoryDialog] = useState(false);
  const [groupForm, setGroupForm] = useState({
    name: '',
    display_name: '',
    description: '',
    icon: '',
    sort_order: 0,
    is_active: true
  });

  useEffect(() => {
    console.log('Initial useEffect running...');
    fetchGroups();
    fetchAttributes();
    fetchCategories();
  }, []);

  const fetchGroups = async () => {
    setLoading(true);
    try {
      const response = await axios.get('/api/admin/attribute-groups');
      const data = response.data;
      // API возвращает объект с полем groups
      setGroups(data.groups || []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка загрузки групп');
    } finally {
      setLoading(false);
    }
  };

  const fetchAttributes = async () => {
    try {
      const response = await axios.get('/api/admin/attributes');
      const data = response.data;
      // API возвращает объект с полем data
      setAttributes(data.data || []);
    } catch (err) {
      console.error('Ошибка загрузки атрибутов:', err);
    }
  };

  const fetchCategories = async () => {
    try {
      const response = await axios.get('/api/admin/categories');
      const data = response.data;
      // API возвращает объект с полем data
      const flatCategories = data.data || [];
      
      console.log('Flat categories from API:', flatCategories);
      
      // Преобразуем плоский список категорий в иерархическое дерево
      const categoryTree = buildCategoryTree(flatCategories);
      console.log('Built category tree:', categoryTree);
      setCategories(categoryTree);
    } catch (err) {
      console.error('Ошибка загрузки категорий:', err);
    }
  };
  
  const buildCategoryTree = (flatCategories: Category[]): Category[] => {
    const categoryMap = new Map<number, Category>();
    const rootCategories: Category[] = [];
    
    // Сначала создаем карту всех категорий по ID
    flatCategories.forEach(cat => {
      categoryMap.set(cat.id, { ...cat, children: [] });
    });
    
    // Затем строим иерархическую структуру
    flatCategories.forEach(cat => {
      const category = categoryMap.get(cat.id);
      if (!category) return;
      
      if (cat.parent_id && categoryMap.get(cat.parent_id)) {
        const parent = categoryMap.get(cat.parent_id);
        if (parent && parent.children) {
          parent.children.push(category);
        }
      } else {
        rootCategories.push(category);
      }
    });
    
    // Сортируем категории на каждом уровне
    const sortCategories = (categories: Category[]): Category[] => {
      return categories.sort((a, b) => 
        (a.display_name || a.name).localeCompare(b.display_name || b.name, 'ru')
      ).map(cat => ({
        ...cat,
        children: cat.children ? sortCategories(cat.children) : []
      }));
    };
    
    return sortCategories(rootCategories);
  };

  const fetchGroupItems = async (groupId: number) => {
    try {
      const response = await axios.get(`/api/admin/attribute-groups/${groupId}/items`);
      const data = response.data?.data || {};
      // API возвращает объект с полем data.items, items может быть null
      setSelectedGroupItems(data.items || []);
    } catch (err) {
      console.error('Ошибка загрузки атрибутов группы:', err);
    }
  };

  const fetchCategoryGroups = async (categoryId: number) => {
    try {
      const response = await axios.get(`/api/admin/categories/${categoryId}/attribute-groups`);
      const data = response.data;
      // API возвращает объект с полем groups
      setCategoryGroups(data.groups || []);
    } catch (err) {
      console.error('Ошибка загрузки групп категории:', err);
    }
  };
  
  const fetchAllCategoryGroups = async () => {
    setLoadingCategories(true);
    try {
      console.log('fetchAllCategoryGroups - categories:', categories);
      const groupMappings = [];
      
      // Рекурсивная функция для обхода всех категорий
      const fetchForCategory = async (category: Category) => {
        const response = await axios.get(`/api/admin/categories/${category.id}/attribute-groups`);
        const data = response.data;
        const groups = data.groups || [];
        for (const group of groups) {
          groupMappings.push({
            ...group,
            category_id: category.id
          });
        }
        
        // Обходим дочерние категории
        if (category.children) {
          for (const child of category.children) {
            await fetchForCategory(child);
          }
        }
      };
      
      // Обходим все корневые категории
      for (const category of categories) {
        await fetchForCategory(category);
      }
      
      console.log('Found category groups:', groupMappings);
      setCategoryGroups(groupMappings);
    } catch (err) {
      console.error('Ошибка загрузки привязок групп к категориям:', err);
    } finally {
      setLoadingCategories(false);
    }
  };

  const handleCreateGroup = () => {
    setSelectedGroup(null);
    setGroupForm({
      name: '',
      display_name: '',
      description: '',
      icon: '',
      sort_order: groups.length + 1,
      is_active: true
    });
    setOpenGroupDialog(true);
  };

  const handleEditGroup = (group: AttributeGroup) => {
    setSelectedGroup(group);
    setGroupForm({
      name: group.name,
      display_name: group.display_name,
      description: group.description,
      icon: group.icon,
      sort_order: group.sort_order,
      is_active: group.is_active
    });
    setOpenGroupDialog(true);
  };

  const handleDeleteGroup = async (id: number) => {
    if (window.confirm('Удалить группу атрибутов?')) {
      try {
        await axios.delete(`/api/admin/attribute-groups/${id}`);
        fetchGroups();
      } catch (err: any) {
        setError(err.response?.data?.error || 'Ошибка удаления группы');
      }
    }
  };

  const handleSaveGroup = async () => {
    try {
      if (selectedGroup) {
        await axios.put(`/api/admin/attribute-groups/${selectedGroup.id}`, groupForm);
      } else {
        await axios.post('/api/admin/attribute-groups', groupForm);
      }
      setOpenGroupDialog(false);
      fetchGroups();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Ошибка сохранения группы');
    }
  };

  const handleManageItems = async (group: AttributeGroup) => {
    setSelectedGroup(group);
    await fetchGroupItems(group.id);
    setOpenItemsDialog(true);
  };

  const handleToggleAttribute = async (attributeId: number, isChecked: boolean) => {
    if (!selectedGroup) return;

    try {
      if (isChecked) {
        await axios.post(`/api/admin/attribute-groups/${selectedGroup.id}/items`, {
          attribute_id: attributeId,
          sort_order: selectedGroupItems.length + 1
        });
      } else {
        await axios.delete(`/api/admin/attribute-groups/${selectedGroup.id}/items/${attributeId}`);
      }
      await fetchGroupItems(selectedGroup.id);
    } catch (err) {
      console.error('Ошибка изменения атрибута в группе:', err);
    }
  };

  const handleReorderItems = async (result: any) => {
    if (!result.destination || !selectedGroup) return;

    const items = Array.from(selectedGroupItems);
    const [reorderedItem] = items.splice(result.source.index, 1);
    items.splice(result.destination.index, 0, reorderedItem);

    // Обновляем порядок сортировки
    const updatedItems = items.map((item, index) => ({
      ...item,
      sort_order: index + 1
    }));

    setSelectedGroupItems(updatedItems);

    // Здесь можно добавить API-вызов для сохранения нового порядка
  };

  const handleCategoryGroupToggle = async (categoryId: number, isChecked: boolean) => {
    if (!selectedGroup) return;

    try {
      if (isChecked) {
        await axios.post(`/api/admin/categories/${categoryId}/attribute-groups`, {
          group_id: selectedGroup.id,
          sort_order: categoryGroups.length + 1,
          is_active: true,
          display_mode: 'list',
          collapsed_by_default: false,
          configuration: {}
        });
        setSuccess('Группа привязана к категории');
      } else {
        await axios.delete(`/api/admin/categories/${categoryId}/attribute-groups/${selectedGroup.id}`);
        setSuccess('Группа отвязана от категории');
      }
      await fetchAllCategoryGroups();
    } catch (err: any) {
      console.error('Ошибка изменения привязки группы к категории:', err);
      setError(err.response?.data?.error || 'Ошибка изменения привязки');
    }
  };

  const handleManageCategories = async (group: AttributeGroup) => {
    setSelectedGroup(group);
    setLoadingCategories(true);
    setOpenCategoryDialog(true);
    
    // Убедимся, что категории загружены
    if (categories.length === 0) {
      await fetchCategories();
    }
    
    await fetchAllCategoryGroups();
  };

  const handleSelectCategory = async (category: Category) => {
    setSelectedCategory(category);
    await fetchCategoryGroups(category.id);
  };

  const handleToggleGroupInCategory = async (groupId: number, isChecked: boolean) => {
    if (!selectedCategory) return;

    try {
      if (isChecked) {
        await axios.post(`/api/admin/categories/${selectedCategory.id}/attribute-groups`, {
          group_id: groupId,
          sort_order: categoryGroups.length + 1,
          is_active: true,
          display_mode: 'list',
          collapsed_by_default: false,
          configuration: {}
        });
      } else {
        await axios.delete(`/api/admin/categories/${selectedCategory.id}/attribute-groups/${groupId}`);
      }
      await fetchCategoryGroups(selectedCategory.id);
    } catch (err) {
      console.error('Ошибка изменения группы в категории:', err);
    }
  };

  const toggleCategoryExpansion = (categoryId: number) => {
    const newExpanded = new Set(expandedCategories);
    if (newExpanded.has(categoryId)) {
      newExpanded.delete(categoryId);
    } else {
      newExpanded.add(categoryId);
    }
    setExpandedCategories(newExpanded);
  };

  const filterCategories = (category: Category, query: string): boolean => {
    const lowerQuery = query.toLowerCase();
    const displayName = category.display_name || category.name || '';
    const name = category.name || '';
    const nameMatch = displayName.toLowerCase().includes(lowerQuery) || 
                     name.toLowerCase().includes(lowerQuery);
    
    if (nameMatch) return true;
    
    if (category.children) {
      return category.children.some(child => filterCategories(child, query));
    }
    
    return false;
  };

  const renderCategoryTree = (categoriesList: Category[], level: number = 0): React.ReactNode => {
    console.log('renderCategoryTree called with:', categoriesList);
    if (!categoriesList || categoriesList.length === 0) {
      return <Typography variant="body2" sx={{ p: 2 }}>Нет категорий для отображения</Typography>;
    }
    
    const filteredCategories = categoriesList
      .filter(category => !categorySearchQuery || filterCategories(category, categorySearchQuery))
      .filter(category => {
        if (showOnlyLinked && selectedGroup) {
          // Проверяем привязана ли категория или любая из ее подкатегорий
          const isLinkedRecursive = (cat: Category): boolean => {
            const isLinked = categoryGroups.some(
              cg => cg.group_id === selectedGroup.id && cg.category_id === cat.id
            );
            if (isLinked) return true;
            
            if (cat.children) {
              return cat.children.some(child => isLinkedRecursive(child));
            }
            return false;
          };
          return isLinkedRecursive(category);
        }
        return true;
      });
    
    console.log('Filtered categories:', filteredCategories);
    return filteredCategories.map(category => renderCategoryItem(category, level));
  };
  
  const renderCategoryItem = (category: Category, level: number = 0): React.ReactNode => {
    const isLinked = categoryGroups.some(
      cg => cg.group_id === selectedGroup?.id && cg.category_id === category.id
    );
    const hasChildren = category.children && category.children.length > 0;
    const isExpanded = expandedCategories.has(category.id);
    const childrenCount = category.children?.length || 0;

    return (
      <Fragment key={category.id}>
        <ListItem 
          sx={{ 
            pl: level * 2,
            '.MuiListItemIcon-root': {
              minWidth: hasChildren ? 36 : 42
            }
          }}
        >
          {hasChildren && (
            <IconButton 
              size="small" 
              onClick={() => toggleCategoryExpansion(category.id)}
              sx={{ p: 0.5 }}
            >
              {isExpanded ? <ExpandMoreIcon fontSize="small" /> : <ChevronRightIcon fontSize="small" />}
            </IconButton>
          )}
          <ListItemIcon>
            <Checkbox
              checked={isLinked}
              onChange={() => handleCategoryGroupToggle(category.id, !isLinked)}
              size="small"
            />
          </ListItemIcon>
          <ListItemText 
            primary={
              <Box display="flex" alignItems="center">
                <Typography variant="body2">
                  {category.display_name || category.name}
                </Typography>
                {childrenCount > 0 && (
                  <Chip
                    label={childrenCount}
                    size="small"
                    sx={{ ml: 1, height: 20, fontSize: '0.75rem' }}
                  />
                )}
              </Box>
            } 
            secondary={
              <Typography variant="caption" color="text.secondary">
                {category.name !== (category.display_name || category.name) ? category.name : ''}
              </Typography>
            }
          />
        </ListItem>
        {hasChildren && isExpanded && (
          <Box sx={{ ml: 2 }}>
            {renderCategoryTree(category.children!, level + 1)}
          </Box>
        )}
      </Fragment>
    );
  };

  const handleSelectAllCategories = async () => {
    if (!selectedGroup) return;
    
    const getAllCategoryIds = (cats: Category[]): number[] => {
      let ids: number[] = [];
      cats.forEach(cat => {
        ids.push(cat.id);
        if (cat.children) {
          ids = [...ids, ...getAllCategoryIds(cat.children)];
        }
      });
      return ids;
    };
    
    const allIds = getAllCategoryIds(categories);
    const linkedIds = categoryGroups
      .filter(cg => cg.group_id === selectedGroup.id)
      .map(cg => cg.category_id);
    
    const unlinkedIds = allIds.filter(id => !linkedIds.includes(id));
    
    for (const categoryId of unlinkedIds) {
      try {
        await axios.post(`/api/admin/categories/${categoryId}/attribute-groups`, {
          group_id: selectedGroup.id,
          sort_order: categoryGroups.length + 1,
          is_active: true,
          display_mode: 'list',
          collapsed_by_default: false,
          configuration: {}
        });
      } catch (err) {
        console.error('Ошибка привязки группы к категории:', err);
      }
    }
    
    await fetchAllCategoryGroups();
    setSuccess('Группа привязана ко всем категориям');
  };

  const handleDeselectAllCategories = async () => {
    if (!selectedGroup) return;
    
    const linkedCategories = categoryGroups.filter(cg => cg.group_id === selectedGroup.id);
    
    for (const cg of linkedCategories) {
      try {
        await axios.delete(`/api/admin/categories/${cg.category_id}/attribute-groups/${selectedGroup.id}`);
      } catch (err) {
        console.error('Ошибка отвязки группы от категории:', err);
      }
    }
    
    await fetchAllCategoryGroups();
    setSuccess('Группа отвязана от всех категорий');
  };

  const GroupsTab = () => (
    <>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h5">Группы атрибутов</Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={handleCreateGroup}
        >
          Создать группу
        </Button>
      </Box>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Название</TableCell>
              <TableCell>Отображаемое имя</TableCell>
              <TableCell>Описание</TableCell>
              <TableCell>Статус</TableCell>
              <TableCell>Системная</TableCell>
              <TableCell>Действия</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {(groups || []).map((group) => (
              <TableRow key={group.id}>
                <TableCell>{group.name}</TableCell>
                <TableCell>{group.display_name}</TableCell>
                <TableCell>{group.description}</TableCell>
                <TableCell>
                  <Chip
                    label={group.is_active ? 'Активна' : 'Неактивна'}
                    color={group.is_active ? 'success' : 'default'}
                    size="small"
                  />
                </TableCell>
                <TableCell>
                  <Chip
                    label={group.is_system ? 'Да' : 'Нет'}
                    color={group.is_system ? 'info' : 'default'}
                    size="small"
                  />
                </TableCell>
                <TableCell>
                  <Tooltip title="Редактировать">
                    <IconButton onClick={() => handleEditGroup(group)}>
                      <EditIcon />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Управление атрибутами">
                    <IconButton onClick={() => handleManageItems(group)}>
                      <SettingsIcon />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Привязка к категориям">
                    <IconButton onClick={() => handleManageCategories(group)}>
                      <LinkIcon />
                    </IconButton>
                  </Tooltip>
                  {!group.is_system && (
                    <Tooltip title="Удалить">
                      <IconButton onClick={() => handleDeleteGroup(group.id)}>
                        <DeleteIcon />
                      </IconButton>
                    </Tooltip>
                  )}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </>
  );

  const CategoriesTab = () => (
    <Grid container spacing={3}>
      <Grid item xs={4}>
        <Typography variant="h6" gutterBottom>Категории</Typography>
        <List>
          {(categories || []).map((category) => (
            <ListItem
              key={category.id}
              button
              selected={selectedCategory?.id === category.id}
              onClick={() => handleSelectCategory(category)}
            >
              <ListItemText primary={category.display_name} />
            </ListItem>
          ))}
        </List>
      </Grid>
      <Grid item xs={8}>
        {selectedCategory && (
          <>
            <Typography variant="h6" gutterBottom>
              Группы атрибутов для: {selectedCategory.display_name}
            </Typography>
            <List>
              {(groups || []).map((group) => {
                const isAttached = categoryGroups.some(cg => cg.group_id === group.id);
                return (
                  <ListItem key={group.id}>
                    <ListItemIcon>
                      <Checkbox
                        checked={isAttached}
                        onChange={(e) => handleToggleGroupInCategory(group.id, e.target.checked)}
                      />
                    </ListItemIcon>
                    <ListItemText
                      primary={group.display_name}
                      secondary={group.description}
                    />
                  </ListItem>
                );
              })}
            </List>
          </>
        )}
      </Grid>
    </Grid>
  );

  return (
    <Box>
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography variant="h4" gutterBottom>
          Управление группами атрибутов
        </Typography>

        {error && (
          <Alert severity="error" onClose={() => setError(null)} sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <Tabs value={activeTab} onChange={(e, v) => setActiveTab(v)} sx={{ mb: 3 }}>
          <Tab label="Группы атрибутов" />
          <Tab label="Привязка к категориям" />
        </Tabs>

        {loading ? (
          <Box display="flex" justifyContent="center" alignItems="center" minHeight={200}>
            <CircularProgress />
          </Box>
        ) : (
          <>
            {activeTab === 0 && <GroupsTab />}
            {activeTab === 1 && <CategoriesTab />}
          </>
        )}
      </Paper>

      {/* Диалог создания/редактирования группы */}
      <Dialog open={openGroupDialog} onClose={() => setOpenGroupDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle>
          {selectedGroup ? 'Редактировать группу' : 'Создать группу'}
        </DialogTitle>
        <DialogContent>
          <Box display="flex" flexDirection="column" gap={2} sx={{ mt: 2 }}>
            <TextField
              label="Системное имя"
              value={groupForm.name}
              onChange={(e) => setGroupForm({ ...groupForm, name: e.target.value })}
              fullWidth
              required
            />
            <TextField
              label="Отображаемое имя"
              value={groupForm.display_name}
              onChange={(e) => setGroupForm({ ...groupForm, display_name: e.target.value })}
              fullWidth
              required
            />
            <TextField
              label="Описание"
              value={groupForm.description}
              onChange={(e) => setGroupForm({ ...groupForm, description: e.target.value })}
              fullWidth
              multiline
              rows={2}
            />
            <TextField
              label="Иконка"
              value={groupForm.icon}
              onChange={(e) => setGroupForm({ ...groupForm, icon: e.target.value })}
              fullWidth
            />
            <TextField
              label="Порядок сортировки"
              type="number"
              value={groupForm.sort_order}
              onChange={(e) => setGroupForm({ ...groupForm, sort_order: parseInt(e.target.value) })}
              fullWidth
            />
            <FormControlLabel
              control={
                <Switch
                  checked={groupForm.is_active}
                  onChange={(e) => setGroupForm({ ...groupForm, is_active: e.target.checked })}
                />
              }
              label="Активна"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenGroupDialog(false)}>Отмена</Button>
          <Button onClick={handleSaveGroup} variant="contained" color="primary">
            Сохранить
          </Button>
        </DialogActions>
      </Dialog>

      {/* Диалог управления атрибутами в группе */}
      <Dialog open={openItemsDialog} onClose={() => setOpenItemsDialog(false)} maxWidth="md" fullWidth>
        <DialogTitle>
          Атрибуты в группе: {selectedGroup?.display_name}
        </DialogTitle>
        <DialogContent>
          <DragDropContext onDragEnd={handleReorderItems}>
            <Droppable droppableId="attributes">
              {(provided) => (
                <List {...provided.droppableProps} ref={provided.innerRef}>
                  {(attributes || []).map((attribute, index) => {
                    const isInGroup = selectedGroupItems.some(item => item.attribute_id === attribute.id);
                    return (
                      <Draggable key={attribute.id} draggableId={`attr-${attribute.id}`} index={index}>
                        {(provided) => (
                          <ListItem
                            ref={provided.innerRef}
                            {...provided.draggableProps}
                            {...provided.dragHandleProps}
                          >
                            <ListItemIcon>
                              <Checkbox
                                checked={isInGroup}
                                onChange={(e) => handleToggleAttribute(attribute.id, e.target.checked)}
                              />
                            </ListItemIcon>
                            <ListItemIcon>
                              <DragIcon />
                            </ListItemIcon>
                            <ListItemText
                              primary={attribute.display_name}
                              secondary={`${attribute.name} (${attribute.data_type})`}
                            />
                          </ListItem>
                        )}
                      </Draggable>
                    );
                  })}
                  {provided.placeholder}
                </List>
              )}
            </Droppable>
          </DragDropContext>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenItemsDialog(false)}>Закрыть</Button>
        </DialogActions>
      </Dialog>

      {/* Диалог привязки к категориям */}
      <Dialog open={openCategoryDialog} onClose={() => setOpenCategoryDialog(false)} maxWidth="md" fullWidth>
        <DialogTitle>
          Выберите категории для группы "{selectedGroup?.display_name}"
        </DialogTitle>
        <DialogContent>
          <Typography variant="body2" color="textSecondary" gutterBottom>
            Выберите категории, в которых должна отображаться эта группа атрибутов
          </Typography>
          {loadingCategories ? (
            <Box display="flex" justifyContent="center" alignItems="center" minHeight={200}>
              <CircularProgress />
            </Box>
          ) : (
            <>
              <Box mb={2}>
                <TextField
                  fullWidth
                  variant="outlined"
                  placeholder="Поиск категорий..."
                  value={categorySearchQuery}
                  onChange={(e) => setCategorySearchQuery(e.target.value)}
                  InputProps={{
                    startAdornment: (
                      <InputAdornment position="start">
                        <SearchIcon />
                      </InputAdornment>
                    ),
                    endAdornment: categorySearchQuery && (
                      <InputAdornment position="end">
                        <IconButton size="small" onClick={() => setCategorySearchQuery('')}>
                          <ClearIcon />
                        </IconButton>
                      </InputAdornment>
                    )
                  }}
                />
              </Box>
              
              {/* Фильтр "Только привязанные" */}
              <Box mb={2}>
                <FormControlLabel
                  control={
                    <Switch
                      checked={showOnlyLinked}
                      onChange={(e) => setShowOnlyLinked(e.target.checked)}
                    />
                  }
                  label="Показывать только привязанные категории"
                />
              </Box>
              <Box mb={2} display="flex" justifyContent="space-between" alignItems="center">
                <Typography variant="body2">
                  Выбрано категорий: {
                    categoryGroups.filter(cg => cg.group_id === selectedGroup?.id).length
                  }
                </Typography>
                <Box>
                  <Button 
                    size="small" 
                    onClick={() => handleSelectAllCategories()}
                    sx={{ mr: 1 }}
                  >
                    Выбрать все
                  </Button>
                  <Button 
                    size="small" 
                    onClick={() => handleDeselectAllCategories()}
                  >
                    Снять все
                  </Button>
                </Box>
              </Box>
              <List sx={{ maxHeight: 400, overflow: 'auto' }}>
                {(() => {
                  console.log('Rendering categories:', categories);
                  return renderCategoryTree(categories);
                })()}
              </List>
            </>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenCategoryDialog(false)}>Закрыть</Button>
        </DialogActions>
      </Dialog>

      {/* Уведомления об ошибках */}
      <Snackbar
        open={!!error}
        autoHideDuration={6000}
        onClose={() => setError(null)}
      >
        <Alert onClose={() => setError(null)} severity="error" sx={{ width: '100%' }}>
          {error}
        </Alert>
      </Snackbar>

      {/* Уведомления об успехе */}
      <Snackbar
        open={!!success}
        autoHideDuration={3000}
        onClose={() => setSuccess(null)}
      >
        <Alert onClose={() => setSuccess(null)} severity="success" sx={{ width: '100%' }}>
          {success}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default AttributeGroupsPage;