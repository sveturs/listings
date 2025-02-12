//frontend/hostel-frontend/src/components/marketplace/VirtualizedCategoryTree.js
import React, { useCallback, useEffect, useState, useRef, useMemo } from 'react';
import { FixedSizeList as List } from 'react-window';
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

const ITEM_SIZE = 40;
const PAGE_SIZE = 20;


const CategoryItem = React.memo(({ data, index, style }) => {
    const { i18n } = useTranslation();
    const {
        items,
        expandedItems,
        selectedId,
        onToggle,
        onSelect,
    } = data;

    const item = items[index];
    // Перемещаем проверку в конец, после всех хуков
    const renderKey = `${item?.id}-${item ? expandedItems.has(item.id) : false}`;
    const memoizedItem = useMemo(() => {
        if (!item) return null;
        // Оставляем только одну группу логов
        if (process.env.NODE_ENV === 'development') {
            console.log('Rendering category:', item);
            console.log('Children:', item.children);
        }
        return item;
    }, [renderKey, item]);

    // Теперь проверяем memoizedItem
    if (!memoizedItem) return null;

    const isExpanded = expandedItems.has(memoizedItem.id);
    const isSelected = selectedId === memoizedItem.id;
    const hasChildren = memoizedItem.children?.length > 0;
    const paddingLeft = memoizedItem.level * 16 + 8;

    // Добавляем функцию получения переведенного имени
    const getTranslatedName = (category) => {
        if (category.translations && category.translations[i18n.language]) {
            return category.translations[i18n.language];
        }
        return category.name;
    };

    return (
        <ListItemButton
            style={{
                ...style,
                paddingLeft,
                backgroundColor: isSelected ? 'rgba(0, 0, 0, 0.04)' : 'transparent'
            }}
            onClick={() => onSelect(memoizedItem.id)}
            dense
        >
            <ListItemIcon sx={{ minWidth: 24 }}>
                {hasChildren && (
                    <Box
                        component="span"
                        onClick={(e) => {
                            e.stopPropagation();
                            onToggle(memoizedItem.id);
                        }}
                    >
                        {isExpanded ? <ChevronDown size={18} /> : <ChevronRight size={18} />}
                    </Box>
                )}
            </ListItemIcon>
            <ListItemText
                primary={
                    <Typography variant="body2" noWrap>
                        {getTranslatedName(memoizedItem)}
                        <Typography
                            component="span"
                            variant="caption"
                            sx={{ ml: 1, color: 'text.secondary' }}
                        >
                            ({memoizedItem.listing_count || 0})
                        </Typography>
                    </Typography>
                }
            />
        </ListItemButton>
    );
});

const VirtualizedCategoryTree = ({ selectedId, onSelectCategory }) => {
    const { i18n } = useTranslation();
    const [expandedItems, setExpandedItems] = useState(new Set());
    const [flattenedItems, setFlattenedItems] = useState([]);
    const [treeData, setTreeData] = useState(null); 

    const buildFlattenedList = useCallback((items, level = 0, result = []) => {
        if (!items || !Array.isArray(items)) return result;
        
        // Сначала сортируем элементы
        const sortedItems = [...items].sort((a, b) => a.name.localeCompare(b.name));
        
        sortedItems.forEach(item => {
            const itemCopy = {
                ...item,
                level,
                children: Array.isArray(item.children) ? [...item.children].sort((a, b) => 
                    a.name.localeCompare(b.name)) : []
            };
            result.push(itemCopy);
            
            if (expandedItems.has(item.id) && itemCopy.children.length > 0) {
                buildFlattenedList(itemCopy.children, level + 1, result);
            }
        });
        
        return result;
    }, [expandedItems]);

    // Изменяем запрос данных:
    const { data: queryResult, isLoading } = useInfiniteQuery(
        ['categories', i18n.language],
        async () => {
            console.log('Fetching categories...');
            const response = await axios.get('/api/v1/marketplace/category-tree');
            console.log('Categories response:', response.data);
            return response.data;  // Убираем .data, так как данные уже в response.data
        },
        {
            getNextPageParam: false,
            staleTime: 5 * 60 * 1000,
        }
    );

    // Сохраняем данные в ref при первой загрузке
    useEffect(() => {
        console.log('Query result:', queryResult);
        if (queryResult?.pages?.[0]?.data && !treeData) {  // Заменяем проверку treeRef на treeData
            console.log('Setting tree data:', queryResult.pages[0].data);
            setTreeData(queryResult.pages[0].data);  // Используем setTreeData вместо присвоения treeRef.current
            const flatList = buildFlattenedList(queryResult.pages[0].data);
            setFlattenedItems(flatList);
        }
    }, [queryResult, buildFlattenedList]);

    // Обновляем при изменении expanded items
    useEffect(() => {
        if (treeData) {  // Заменяем проверку treeRef.current на treeData
            const flatList = buildFlattenedList(treeData);
            setFlattenedItems(flatList);
        }
    }, [expandedItems, buildFlattenedList, treeData]);

    const handleToggle = useCallback((id) => {
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
                itemSize={40}
                width="100%"
                itemData={{
                    items: flattenedItems,
                    expandedItems,
                    selectedId,
                    onToggle: handleToggle,
                    onSelect: onSelectCategory
                }}
            >
                {CategoryItem}
            </List>
        </Box>
    );
};

export default React.memo(VirtualizedCategoryTree);