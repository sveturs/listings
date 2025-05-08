// frontend/hostel-frontend/src/components/marketplace/CategorySelect.tsx
import React, { useState, useEffect, MouseEvent } from 'react';  
import { useTranslation } from 'react-i18next';
import {
    Box,
    Button,
    Popover,
    List,
    ListItemButton,
    ListItemText,
    ListItemIcon,
    Typography,
    Paper,
    Stack
} from '@mui/material';
import { ChevronRight, ChevronLeft } from 'lucide-react';

export interface Category {
    id: number | string;
    name: string;
    parent_id?: number | string | null;
    translations?: {
        [lang: string]: string;
    };
    [key: string]: any;
}

interface CategorySelectProps {
    categories: Category[];
    value?: number | string | null;
    onChange: (categoryId: number | string) => void;
    error?: boolean;
}

const CategorySelect: React.FC<CategorySelectProps> = ({ categories, value, onChange, error }) => {
    const { t, i18n } = useTranslation('marketplace') as any;
    const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);
    const [currentPath, setCurrentPath] = useState<Category[]>([]);
    const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);

    useEffect(() => {
        console.log('Current language:', i18n.language);
        console.log('Categories with translations:', categories);
        // Проверим одну категорию для примера
        if (categories.length > 0) {
            console.log('First category:', categories[0]);
            console.log('Translations for first category:', categories[0].translations);
        }
    }, [categories, i18n.language]);

    // Получение переведенного имени категории
    const getTranslatedName = (category: Category): string => {
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

    const handleClick = (event: MouseEvent<HTMLButtonElement>): void => {
        setAnchorEl(event.currentTarget);
        setCurrentPath([]);
    };

    const handleClose = (): void => {
        setAnchorEl(null);
    };

    const getCurrentCategories = (): Category[] => {
        if (currentPath.length === 0) {
            return categories.filter(cat => !cat.parent_id);
        }
        const currentCategory = currentPath[currentPath.length - 1];
        return categories.filter(cat => cat.parent_id === currentCategory.id);
    };

    const handleCategoryClick = (category: Category): void => {
        const hasChildren = categories.some(cat => cat.parent_id === category.id);

        if (hasChildren) {
            setCurrentPath([...currentPath, category]);
        } else {
            setSelectedCategory(category);
            onChange(category.id);
            handleClose();
        }
    };

    const handleBack = (): void => {
        setCurrentPath(currentPath.slice(0, -1));
    };

    // Находим полный путь для выбранной категории
    const findCategoryPath = (categoryId: number | string | null | undefined): Category[] => {
        if (!categoryId) return [];
        
        const path: Category[] = [];
        let current = categories.find(c => c.id === categoryId);

        while (current) {
            path.unshift(current);
            current = categories.find(c => c.id === current.parent_id);
        }

        return path;
    };

    const selectedPath = value ? findCategoryPath(value) : [];

    return (
        <>
            <Button
                onClick={handleClick}
                variant="outlined"
                fullWidth
                sx={{
                    justifyContent: 'flex-start',
                    textAlign: 'left',
                    color: error ? 'error.main' : 'text.primary',
                    borderColor: error ? 'error.main' : 'inherit'
                }}
            >
                {selectedPath.length > 0 ? (
                    <Stack spacing={0.5}>
                        <Typography variant="body1" noWrap>
                            {getTranslatedName(selectedPath[selectedPath.length - 1])}
                        </Typography>
                        <Typography variant="caption" color="text.secondary" noWrap>
                            {selectedPath.map(cat => getTranslatedName(cat)).join(' > ')}
                        </Typography>
                    </Stack>
                ) : (
                    t('listings.details.ChooseAcategory')
                )}
            </Button>

            <Popover
                open={Boolean(anchorEl)}
                anchorEl={anchorEl}
                onClose={handleClose}
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'left',
                }}
                transformOrigin={{
                    vertical: 'top',
                    horizontal: 'left',
                }}
                PaperProps={{
                    sx: {
                        width: 320,
                        maxHeight: 400,
                    }
                }}
            >
                <Paper elevation={0}>
                    {currentPath.length > 0 && (
                        <ListItemButton onClick={handleBack}>
                            <ListItemIcon sx={{ minWidth: 32 }}>
                                <ChevronLeft size={20} />
                            </ListItemIcon>
                            <ListItemText
                                primary={getTranslatedName(currentPath[currentPath.length - 1])}
                                primaryTypographyProps={{
                                    variant: 'subtitle2',
                                    color: 'text.secondary'
                                }}
                            />
                        </ListItemButton>
                    )}

                    <List sx={{ py: 0 }}>
                        {getCurrentCategories().map((category) => {
                            const hasChildren = categories.some(cat => cat.parent_id === category.id);
                            const isSelected = value === category.id;

                            return (
                                <ListItemButton
                                    key={category.id}
                                    onClick={() => handleCategoryClick(category)}
                                    selected={isSelected}
                                >
                                    <ListItemText
                                        primary={getTranslatedName(category)}
                                        primaryTypographyProps={{
                                            variant: 'body2',
                                            color: isSelected ? 'primary.main' : 'text.primary'
                                        }}
                                    />
                                    {hasChildren && (
                                        <ChevronRight size={20} />
                                    )}
                                </ListItemButton>
                            );
                        })}
                    </List>
                </Paper>
            </Popover>
        </>
    );
};

export default CategorySelect;