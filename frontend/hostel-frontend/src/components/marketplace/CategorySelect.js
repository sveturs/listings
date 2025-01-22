import React, { useState } from 'react';
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

const CategorySelect = ({ categories, value, onChange, error }) => {
    const [anchorEl, setAnchorEl] = useState(null);
    const [currentPath, setCurrentPath] = useState([]);
    const [selectedCategory, setSelectedCategory] = useState(null);

    const handleClick = (event) => {
        setAnchorEl(event.currentTarget);
        setCurrentPath([]);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const getCurrentCategories = () => {
        if (currentPath.length === 0) {
            return categories.filter(cat => !cat.parent_id);
        }
        const currentCategory = currentPath[currentPath.length - 1];
        return categories.filter(cat => cat.parent_id === currentCategory.id);
    };

    const handleCategoryClick = (category) => {
        const hasChildren = categories.some(cat => cat.parent_id === category.id);
        
        if (hasChildren) {
            setCurrentPath([...currentPath, category]);
        } else {
            setSelectedCategory(category);
            onChange(category.id);
            handleClose();
        }
    };

    const handleBack = () => {
        setCurrentPath(currentPath.slice(0, -1));
    };

    // Находим полный путь для выбранной категории
    const findCategoryPath = (categoryId) => {
        const path = [];
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
                            {selectedPath[selectedPath.length - 1].name}
                        </Typography>
                        <Typography variant="caption" color="text.secondary" noWrap>
                            {selectedPath.map(cat => cat.name).join(' > ')}
                        </Typography>
                    </Stack>
                ) : (
                    'Выберите категорию'
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
                                primary={currentPath[currentPath.length - 1].name}
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
                                        primary={category.name}
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