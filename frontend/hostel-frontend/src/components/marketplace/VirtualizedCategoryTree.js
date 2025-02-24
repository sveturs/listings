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
const buildTree = (flatList) => {
    const rootCategories = flatList.filter(item => !item.parent_id);
    const nonRootCategories = flatList.filter(item => item.parent_id);
    
    // Сначала создаём корневые категории
    const tree = rootCategories.map(item => ({
        ...item,
        level: 0,
        children: []
    }));
    
    // Сортируем корневые категории по имени
    tree.sort((a, b) => a.name.localeCompare(b.name));
    
    // Создаём Map для быстрого доступа к категориям
    const categoryMap = new Map();
    tree.forEach(category => categoryMap.set(category.id, category));
    
    // Добавляем все остальные категории к их родителям
    nonRootCategories.forEach(item => {
        const parent = categoryMap.get(item.parent_id);
        if (parent) {
            const child = {
                ...item,
                level: parent.level + 1,
                children: []
            };
            parent.children.push(child);
            categoryMap.set(child.id, child);
            
            // Сортируем дочерние элементы по имени
            parent.children.sort((a, b) => a.name.localeCompare(b.name));
        }
    });
    
    return tree;
};

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
    if (!item) return null;

    const isExpanded = expandedItems.has(item.id);
    const isSelected = selectedId === item.id;
    const hasChildren = item.children?.length > 0;
    const paddingLeft = item.level * 24;

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
                        onClick={(e) => {
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
                        {item.name}
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
const VirtualizedCategoryTree = ({ selectedId, onSelectCategory }) => {
    const { i18n } = useTranslation();
    const [expandedItems, setExpandedItems] = useState(new Set());
    const [flattenedItems, setFlattenedItems] = useState([]);
    const [treeData, setTreeData] = useState(null);

    const buildFlattenedList = useCallback((items, level = 0, result = []) => {
        if (!items || !Array.isArray(items)) return result;
        
        items.forEach(item => {
            // Добавляем элемент с его уровнем
            const itemCopy = {
                ...item,
                level: level,
                hasChildren: item.children?.length > 0
            };
            result.push(itemCopy);
            
            // Если элемент развёрнут и имеет дочерние элементы, добавляем их с увеличенным уровнем
            if (expandedItems.has(item.id) && item.children && item.children.length > 0) {
                buildFlattenedList(item.children, level + 1, result);
            }
        });
        
        return result;
    }, [expandedItems]);

    // Изменяем запрос данных:
    const { data: queryResult, isLoading } = useInfiniteQuery(
        ['categories', i18n.language],
        async () => {
            const response = await axios.get('/api/v1/marketplace/category-tree');
            console.log('Full category tree:', JSON.stringify(response.data.data, null, 2));
            return response.data;
        },
        {
            getNextPageParam: false,
            staleTime: 5 * 60 * 1000,
        }
    );

    // Сохраняем данные в ref при первой загрузке
    useEffect(() => {
        if (queryResult?.pages?.[0]?.data) {
            const flatData = queryResult.pages[0].data;
            console.log('Received flat data:', flatData);
            
            const treeStructure = buildTree(flatData);
            console.log('Built tree structure:', treeStructure);
            
            setTreeData(treeStructure);
            const initialFlatList = buildFlattenedList(treeStructure);
            console.log('Initial flattened list:', initialFlatList);
            setFlattenedItems(initialFlatList);
        }
    }, [queryResult, buildFlattenedList]);

    // Обновляем при изменении expanded items
    useEffect(() => {
        if (treeData) {
            console.log('Building flattened list from:', treeData);
            const flatList = buildFlattenedList(treeData);
            console.log('Resulting flattened list:', flatList);
            setFlattenedItems(flatList);
        }
    }, [expandedItems, buildFlattenedList, treeData]);

    const handleToggle = useCallback((id) => {
        console.log('Toggling category:', id);
        setExpandedItems(prev => {
            const next = new Set(prev);
            if (next.has(id)) {
                console.log('Collapsing category:', id);
                next.delete(id);
            } else {
                console.log('Expanding category:', id);
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