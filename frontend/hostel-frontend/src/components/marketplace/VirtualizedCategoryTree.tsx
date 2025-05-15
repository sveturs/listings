import React, { useCallback, useEffect, useState } from 'react';
import { FixedSizeList as List, ListChildComponentProps } from 'react-window';
import { useTranslation } from 'react-i18next';
import { useInfiniteQuery } from 'react-query';
import axios from '../../api/axios';

import {
  Box,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Typography,
  CircularProgress
} from '@mui/material';
import { ChevronRight, ChevronDown } from 'lucide-react';
import { Category as BaseCategory } from './HierarchicalCategorySelect';

// Constants
const ITEM_SIZE = 40;
const PAGE_SIZE = 20;

// Interfaces
interface Category extends BaseCategory {
  level: number;
  hasChildren?: boolean;
  children: Category[];
  listing_count?: number;
}

interface CategoryData {
  items: Category[];
  expandedItems: Set<string | number>;
  selectedId: string | number | null;
  onToggle: (id: string | number) => void;
  onSelect: (id: string | number) => void;
  currentLanguage: string;
}

interface CategoryItemProps {
  data: CategoryData;
  index: number;
  style: React.CSSProperties;
}

interface VirtualizedCategoryTreeProps {
  selectedId: string | number | null;
  onSelectCategory: (id: string | number) => void;
}

// Helper function to build tree from flat list
const buildTree = (flatList: BaseCategory[], i18n: any): Category[] => {
  const categoryMap = new Map<string | number, Category>();

  // First add all categories to Map
  flatList.forEach(category => {
    categoryMap.set(category.id, {
      ...category,
      level: 0,
      children: []
    });
  });

  // Root categories
  const rootCategories: Category[] = [];

  // Build tree
  flatList.forEach(category => {
    const categoryWithChildren = categoryMap.get(category.id);
    
    if (!categoryWithChildren) return;

    if (!category.parent_id) {
      rootCategories.push(categoryWithChildren);
    } else {
      const parent = categoryMap.get(category.parent_id);
      if (parent) {
        categoryWithChildren.level = parent.level + 1;
        parent.children.push(categoryWithChildren);
      } else {
        // If parent not found, add as root
        rootCategories.push(categoryWithChildren);
      }
    }
  });

  // Optimized sorting
  const sortCategoriesByName = (categories: Category[]): void => {
    categories.sort((a, b) => {
      const nameA = a.translations?.[i18n.language] || a.name;
      const nameB = b.translations?.[i18n.language] || b.name;
      return nameA.localeCompare(nameB);
    });

    categories.forEach(category => {
      if (category.children.length > 0) {
        sortCategoriesByName(category.children);
      }
    });
  };

  sortCategoriesByName(rootCategories);

  return rootCategories;
};

// Category Item component
const CategoryItem = React.memo<CategoryItemProps>(({ data, index, style }) => {
  const { i18n } = useTranslation();
  const {
    items,
    expandedItems,
    selectedId,
    onToggle,
    onSelect,
    currentLanguage
  } = data;

  const item = items[index];
  if (!item) return null;

  const isExpanded = expandedItems.has(item.id);
  const isSelected = selectedId === String(item.id) || selectedId === item.id;
  const hasChildren = item.children?.length > 0;
  const paddingLeft = item.level * 24;

  // Function to get translated category name
  const getTranslatedName = (category: Category): string => {
    if (!category) return '';

    // Check for translations
    if (category.translations && typeof category.translations === 'object') {
      // If there's a direct translation for current language
      if (category.translations[i18n.language]) {
        return category.translations[i18n.language];
      }

      // If no direct translation, try to find by priority
      const langPriority = [i18n.language, 'ru', 'sr', 'en'];
      for (const lang of langPriority) {
        if (category.translations[lang]) {
          return category.translations[lang];
        }
      }
    }

    // If no translations or they don't fit, return original name
    return category.name;
  };

  const categoryName = getTranslatedName(item);

  return (
    <ListItemButton
      style={{
        ...style,
        paddingLeft,
        backgroundColor: isSelected ? 'rgba(0, 0, 0, 0.04)' : 'transparent'
      }}
      onClick={() => onSelect(item.id)}
      dense
    >
      <ListItemIcon sx={{ minWidth: 24 }}>
        {hasChildren && (
          <Box
            component="span"
            onClick={(e: React.MouseEvent) => {
              e.stopPropagation();
              onToggle(item.id);
            }}
            sx={{ cursor: 'pointer' }}
          >
            {isExpanded ? <ChevronDown size={18} /> : <ChevronRight size={18} />}
          </Box>
        )}
      </ListItemIcon>
      <ListItemText
        primary={
          <Typography variant="body2" noWrap>
            {categoryName}
            {item.listing_count > 0 && (
              <Typography
                component="span"
                variant="caption"
                sx={{ ml: 1, color: 'text.secondary' }}
              >
                ({item.listing_count})
              </Typography>
            )}
          </Typography>
        }
      />
    </ListItemButton>
  );
});

// Main component
const VirtualizedCategoryTree: React.FC<VirtualizedCategoryTreeProps> = ({ selectedId, onSelectCategory }) => {
  const { i18n } = useTranslation();
  const [expandedItems, setExpandedItems] = useState<Set<string | number>>(new Set());
  const [flattenedItems, setFlattenedItems] = useState<Category[]>([]);
  const [treeData, setTreeData] = useState<Category[] | null>(null);
  const [currentLanguage, setCurrentLanguage] = useState<string>(i18n.language);

  // Track language changes
  useEffect(() => {
    if (currentLanguage !== i18n.language) {
      console.log(`Язык изменился: ${currentLanguage} -> ${i18n.language}`);
      setCurrentLanguage(i18n.language);
    }
  }, [i18n.language, currentLanguage]);

  const buildFlattenedList = useCallback((items: Category[] | null, level = 0, result: Category[] = []): Category[] => {
    if (!items || !Array.isArray(items)) return result;

    for (const item of items) {
      // Add item with its level
      const itemCopy: Category = {
        ...item,
        level: level,
        hasChildren: Array.isArray(item.children) && item.children.length > 0
      };
      result.push(itemCopy);

      // If item is expanded and has children, add them with increased level
      if (expandedItems.has(item.id) && Array.isArray(item.children) && item.children.length > 0) {
        buildFlattenedList(item.children, level + 1, result);
      }
    }

    return result;
  }, [expandedItems]);

  // Modify data query
  const { data: queryResult, isLoading, refetch } = useInfiniteQuery(
    ['categories', i18n.language],
    async () => {
      console.log(`Загрузка категорий для языка: ${i18n.language}`);
      const response = await axios.get('/api/v1/marketplace/category-tree');
      return response.data;
    },
    {
      getNextPageParam: () => undefined, // Функция должна возвращать undefined вместо false для предотвращения pagination
      staleTime: 5 * 60 * 1000,
    }
  );

  // Reload data on language change
  useEffect(() => {
    refetch();
  }, [i18n.language, refetch]);

  // Save data to ref on first load or update
  useEffect(() => {
    if (queryResult?.pages?.[0]?.data) {
      const flatData = queryResult.pages[0].data as BaseCategory[];

      // Add logging for translations structure for debugging
      if (process.env.NODE_ENV === 'development') {
        if (flatData.length > 0) {
          console.log(`Пример переводов категории:`,
            flatData[0].translations ? flatData[0].translations : 'Нет переводов');
        }
      }

      const treeStructure = buildTree(flatData, i18n);

      setTreeData(treeStructure);
      const initialFlatList = buildFlattenedList(treeStructure);
      setFlattenedItems(initialFlatList);
    }
  }, [queryResult, buildFlattenedList, i18n]);

  // Update on expanded items or language change
  useEffect(() => {
    if (treeData) {
      const flatList = buildFlattenedList(treeData);
      setFlattenedItems(flatList);
    }
  }, [expandedItems, buildFlattenedList, treeData, i18n.language]);

  const handleToggle = useCallback((id: string | number) => {
    setExpandedItems(prev => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  }, []);

  if (isLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 2 }}>
        <CircularProgress size={24} />
      </Box>
    );
  }

  return (
    <Box sx={{ height: '100%', maxHeight: 400 }}>
      <List
        height={400}
        itemCount={flattenedItems.length}
        itemSize={ITEM_SIZE}
        width="100%"
        itemData={{
          items: flattenedItems,
          expandedItems,
          selectedId,
          onToggle: handleToggle,
          onSelect: onSelectCategory,
          currentLanguage // Pass current language
        }}
      >
        {CategoryItem}
      </List>
    </Box>
  );
};

export default React.memo(VirtualizedCategoryTree);