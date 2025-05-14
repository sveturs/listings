import React, { useState, useRef } from 'react';
import {
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  ListItemSecondaryAction,
  IconButton,
  Collapse,
  Paper,
  Typography,
  Box,
  Tooltip,
  Chip,
  Divider,
} from '@mui/material';
import {
  ExpandMore,
  ExpandLess,
  Edit as EditIcon,
  Delete as DeleteIcon,
  ArrowUpward as ArrowUpwardIcon,
  ArrowDownward as ArrowDownwardIcon,
  Code as CodeIcon,
  DragIndicator as DragIndicatorIcon,
} from '@mui/icons-material';
import { useTranslation } from 'react-i18next';

interface Category {
  id: number;
  name: string;
  slug: string;
  parent_id?: number | null;
  icon?: string;
  translations?: Record<string, string>;
  listing_count: number;
  has_custom_ui: boolean;
  custom_ui_component?: string;
}

interface CategoryWithChildren extends Category {
  children?: CategoryWithChildren[];
}

interface CategoryTreeViewProps {
  categories: Category[];
  onEdit: (category: Category) => void;
  onDelete: (categoryId: number) => void;
  onReorder: (orderedIds: number[]) => void;
  onMove: (categoryId: number, newParentId: number) => void;
}

const CategoryTreeView: React.FC<CategoryTreeViewProps> = ({
  categories,
  onEdit,
  onDelete,
  onReorder,
  onMove
}) => {
  const { t } = useTranslation();
  const [expanded, setExpanded] = useState<Record<number, boolean>>({});
  const [draggedCategory, setDraggedCategory] = useState<Category | null>(null);
  const [isDraggingOver, setIsDraggingOver] = useState<Record<string, boolean>>({});
  const draggedIdRef = useRef<number | null>(null);
  const draggedParentIdRef = useRef<number | null>(null);

  // Преобразуем плоский список категорий в иерархическое дерево
  const buildCategoryTree = (cats: Category[] = []): CategoryWithChildren[] => {
    if (!cats || !Array.isArray(cats)) {
      return [];
    }

    const categoryMap: Record<number, CategoryWithChildren> = {};
    const rootCategories: CategoryWithChildren[] = [];

    // Сначала создаем карту всех категорий по ID
    cats.forEach(cat => {
      categoryMap[cat.id] = { ...cat, children: [] };
    });

    // Затем строим иерархическую структуру
    cats.forEach(cat => {
      const category = categoryMap[cat.id];
      if (cat.parent_id && categoryMap[cat.parent_id]) {
        if (!categoryMap[cat.parent_id].children) {
          categoryMap[cat.parent_id].children = [];
        }
        categoryMap[cat.parent_id].children!.push(category);
      } else {
        rootCategories.push(category);
      }
    });

    // Сортируем категории
    return rootCategories;
  };

  const handleToggle = (categoryId: number) => {
    setExpanded(prev => ({
      ...prev,
      [categoryId]: !prev[categoryId]
    }));
  };

  // Обработчик начала перетаскивания
  const handleDragStart = (e: React.DragEvent<HTMLDivElement>, category: Category) => {
    // Сохраняем данные категории в dataTransfer
    e.dataTransfer.setData('text/plain', category.id.toString());
    e.dataTransfer.setData('application/json', JSON.stringify(category));
    e.dataTransfer.effectAllowed = 'move';

    // Сохраняем информацию о перетаскиваемой категории
    setDraggedCategory(category);
    draggedIdRef.current = category.id;
    draggedParentIdRef.current = category.parent_id || null;
  };

  // Обработчик разрешения перетаскивания над целью
  const handleDragOver = (e: React.DragEvent<HTMLDivElement>, targetId: string) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';

    // Обновляем состояние перетаскивания, чтобы показать визуальный эффект
    if (!isDraggingOver[targetId]) {
      setIsDraggingOver(prev => ({
        ...prev,
        [targetId]: true
      }));
    }
  };

  // Обработчик выхода за пределы цели при перетаскивании
  const handleDragLeave = (e: React.DragEvent<HTMLDivElement>, targetId: string) => {
    setIsDraggingOver(prev => ({
      ...prev,
      [targetId]: false
    }));
  };

  // Обработчик завершения перетаскивания
  const handleDragEnd = (e: React.DragEvent<HTMLDivElement>) => {
    setDraggedCategory(null);
    draggedIdRef.current = null;
    draggedParentIdRef.current = null;
    setIsDraggingOver({});
  };

  // Обработчик сброса элемента
  const handleDrop = (e: React.DragEvent<HTMLDivElement>, targetParentId: number | 'root', targetIndex: number) => {
    e.preventDefault();
    setIsDraggingOver({});

    const draggedId = parseInt(e.dataTransfer.getData('text/plain'));
    if (isNaN(draggedId)) return;

    const actualParentId = targetParentId === 'root' ? null : targetParentId;
    const draggedParentId = draggedParentIdRef.current;

    // Если категория перетаскивается в тот же уровень (для изменения порядка)
    if (draggedParentId === actualParentId) {
      const parentId = draggedParentId;
      const childCategories = categories.filter(cat =>
        (parentId === null && !cat.parent_id) || cat.parent_id === parentId
      );

      // Находим индекс перетаскиваемой категории в текущем списке
      const draggedIndex = childCategories.findIndex(cat => cat.id === draggedId);
      if (draggedIndex === -1) return;

      // Создаем новый порядок
      const reordered = Array.from(childCategories);
      const [moved] = reordered.splice(draggedIndex, 1);
      reordered.splice(targetIndex, 0, moved);

      // Вызываем callback с новым порядком ID
      const orderedIds = reordered.map(cat => cat.id);
      onReorder(orderedIds);
    }
    // Если категория перетаскивается между разными уровнями (для изменения родителя)
    else {
      const newParentId = actualParentId || 0;
      onMove(draggedId, newParentId);
    }
  };

  const renderCategoryTree = (categoryTree: CategoryWithChildren[] = [], level = 0) => {
    if (!categoryTree || !Array.isArray(categoryTree) || categoryTree.length === 0) {
      return null;
    }

    const parentId = level === 0 ? 'root' : String(categoryTree[0]?.parent_id || 0);

    return (
      <List
        component="div"
        disablePadding
        sx={{ ml: level * 2 }}
        onDragOver={(e) => handleDragOver(e, parentId)}
        onDragLeave={(e) => handleDragLeave(e, parentId)}
        onDrop={(e) => handleDrop(e, parentId === 'root' ? 'root' : parseInt(parentId), 0)}
      >
        {categoryTree.map((category, index) => {
          // Создаем уникальный ID для этой позиции перетаскивания
          const dropId = `${parentId}-${index}`;

          return (
            <div key={category.id}>
              <div
                draggable
                onDragStart={(e) => handleDragStart(e, category)}
                onDragEnd={handleDragEnd}
              >
                <ListItem
                  button
                  divider
                  onClick={() => handleToggle(category.id)}
                  sx={{
                    pl: level * 2,
                    backgroundColor: isDraggingOver[dropId] ? 'rgba(25, 118, 210, 0.12)' :
                                     level % 2 === 0 ? 'rgba(0, 0, 0, 0.03)' : 'inherit',
                    position: 'relative',
                    '&:hover': {
                      '& .drag-indicator': {
                        opacity: 1
                      }
                    }
                  }}
                >
                  <ListItemIcon sx={{ minWidth: 36 }}>
                    <DragIndicatorIcon
                      className="drag-indicator"
                      sx={{
                        opacity: 0.3,
                        cursor: 'grab',
                        '&:active': {
                          cursor: 'grabbing'
                        }
                      }}
                    />
                  </ListItemIcon>

                  <ListItemText
                    primary={
                      <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        <Typography variant="subtitle1">
                          {category.name}
                        </Typography>
                        {category.has_custom_ui && (
                          <Tooltip title={t('admin.categories.hasCustomUi')}>
                            <Chip
                              icon={<CodeIcon />}
                              label={category.custom_ui_component || t('admin.categories.custom')}
                              size="small"
                              color="primary"
                              variant="outlined"
                              sx={{ ml: 1 }}
                            />
                          </Tooltip>
                        )}
                        <Chip
                          label={t('admin.categories.listings', { count: category.listing_count })}
                          size="small"
                          color="default"
                          sx={{ ml: 1 }}
                        />
                      </Box>
                    }
                    secondary={`ID: ${category.id} | Slug: ${category.slug}`}
                  />

                  <ListItemSecondaryAction>
                    <IconButton
                      edge="end"
                      onClick={(e) => {
                        e.stopPropagation();
                        onEdit(category);
                      }}
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      edge="end"
                      onClick={(e) => {
                        e.stopPropagation();
                        onDelete(category.id);
                      }}
                      disabled={category.listing_count > 0 || (category.children && category.children.length > 0)}
                    >
                      <DeleteIcon />
                    </IconButton>

                    {category.children && category.children.length > 0 ? (
                      expanded[category.id] ? <ExpandLess /> : <ExpandMore />
                    ) : null}
                  </ListItemSecondaryAction>
                </ListItem>
              </div>

              {/* Область для перетаскивания между элементами */}
              {index < categoryTree.length - 1 && (
                <div
                  style={{
                    height: '6px',
                    backgroundColor: isDraggingOver[`${dropId}-after`] ? 'rgba(25, 118, 210, 0.3)' : 'transparent',
                    transition: 'background-color 0.2s'
                  }}
                  onDragOver={(e) => handleDragOver(e, `${dropId}-after`)}
                  onDragLeave={(e) => handleDragLeave(e, `${dropId}-after`)}
                  onDrop={(e) => handleDrop(e, parentId === 'root' ? 'root' : parseInt(parentId), index + 1)}
                ></div>
              )}

              {category.children && category.children.length > 0 && (
                <Collapse in={expanded[category.id]} timeout="auto" unmountOnExit>
                  {renderCategoryTree(category.children, level + 1)}
                </Collapse>
              )}
            </div>
          );
        })}
      </List>
    );
  };

  const categoryTree = buildCategoryTree(categories || []);

  return (
    <Paper
      elevation={0}
      variant="outlined"
      sx={{
        maxHeight: '60vh',
        overflow: 'auto',
        bgcolor: 'background.paper',
        borderRadius: 1
      }}
    >
      <Box sx={{ p: 2, bgcolor: 'primary.light', color: 'primary.contrastText' }}>
        <Typography variant="subtitle1" fontWeight="bold">
          {t('admin.categories.structure')}
        </Typography>
        <Typography variant="caption">
          {t('admin.categories.dragHint')}
        </Typography>
      </Box>

      <Divider />

      {categoryTree && categoryTree.length > 0 ? (
        renderCategoryTree(categoryTree)
      ) : (
        <Box sx={{ p: 3, textAlign: 'center' }}>
          <Typography color="textSecondary">
            {t('admin.categories.noCategories')}
          </Typography>
        </Box>
      )}
    </Paper>
  );
};

export default CategoryTreeView;