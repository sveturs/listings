import React, { useCallback, useEffect, useState, useRef } from 'react';
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

const buildTree = (flatList, i18n) => {
    const categoryMap = new Map();

    // Сначала добавляем все категории в Map
    flatList.forEach(category => {
        categoryMap.set(category.id, {
            ...category,
            level: 0,
            children: []
        });
    });

    // Корневые категории
    const rootCategories = [];

    // Строим дерево
    flatList.forEach(category => {
        const categoryWithChildren = categoryMap.get(category.id);

        if (!category.parent_id) {
            rootCategories.push(categoryWithChildren);
        } else {
            const parent = categoryMap.get(category.parent_id);
            if (parent) {
                categoryWithChildren.level = parent.level + 1;
                parent.children.push(categoryWithChildren);
            } else {
                // Если родитель не найден, добавляем как корневую
                rootCategories.push(categoryWithChildren);
            }
        }
    });

    // Оптимизированная сортировка
    const sortCategoriesByName = (categories) => {
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

const CategoryItem = React.memo(({ data, index, style }) => {
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

    // Функция для получения переведенного имени категории
    // Исправленная функция getTranslatedName в CategoryItem
    const getTranslatedName = (category) => {
        if (!category) return '';

        // Проверяем наличие переводов
        if (category.translations && typeof category.translations === 'object') {
            // Если есть прямой перевод на текущий язык
            if (category.translations[i18n.language]) {
                return category.translations[i18n.language];
            }

            // Если прямого перевода нет, пробуем найти по приоритету
            const langPriority = [i18n.language, 'ru', 'sr', 'en'];
            for (const lang of langPriority) {
                if (category.translations[lang]) {
                    return category.translations[lang];
                }
            }
        }

        // Если переводов нет или они не подходят, возвращаем исходное имя
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

const VirtualizedCategoryTree = ({ selectedId, onSelectCategory }) => {
    const { i18n } = useTranslation();
    const [expandedItems, setExpandedItems] = useState(new Set());
    const [flattenedItems, setFlattenedItems] = useState([]);
    const [treeData, setTreeData] = useState(null);
    const [currentLanguage, setCurrentLanguage] = useState(i18n.language);

    // Отслеживаем изменение языка
    useEffect(() => {
        if (currentLanguage !== i18n.language) {
            console.log(`Язык изменился: ${currentLanguage} -> ${i18n.language}`);
            setCurrentLanguage(i18n.language);
        }
    }, [i18n.language, currentLanguage]);

    const buildFlattenedList = useCallback((items, level = 0, result = []) => {
        if (!items || !Array.isArray(items)) return result;

        for (const item of items) {
            // Добавляем элемент с его уровнем
            const itemCopy = {
                ...item,
                level: level,
                hasChildren: Array.isArray(item.children) && item.children.length > 0
            };
            result.push(itemCopy);

            // Если элемент развёрнут и имеет дочерние элементы, добавляем их с увеличенным уровнем
            if (expandedItems.has(item.id) && Array.isArray(item.children) && item.children.length > 0) {
                buildFlattenedList(item.children, level + 1, result);
            }
        }

        return result;
    }, [expandedItems]);

    // Изменяем запрос данных:
    const { data: queryResult, isLoading, refetch } = useInfiniteQuery(
        ['categories', i18n.language],
        async () => {
            console.log(`Загрузка категорий для языка: ${i18n.language}`);
            const response = await axios.get('/api/v1/marketplace/category-tree');
            return response.data;
        },
        {
            getNextPageParam: false,
            staleTime: 5 * 60 * 1000,
        }
    );

    // Повторно загружаем данные при смене языка
    useEffect(() => {
        refetch();
    }, [i18n.language, refetch]);

    // Сохраняем данные в ref при первой загрузке или обновлении
    useEffect(() => {
        if (queryResult?.pages?.[0]?.data) {
            const flatData = queryResult.pages[0].data;

            // Добавляем логирование структуры переводов для отладки
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
    }, [queryResult, buildFlattenedList]);

    // Обновляем при изменении expanded items или языка
    useEffect(() => {
        if (treeData) {
            const flatList = buildFlattenedList(treeData);
            setFlattenedItems(flatList);
        }
    }, [expandedItems, buildFlattenedList, treeData, i18n.language]);

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
                    onSelect: onSelectCategory,
                    currentLanguage // Передаем текущий язык
                }}
            >
                {CategoryItem}
            </List>
        </Box>
    );
};

export default React.memo(VirtualizedCategoryTree);